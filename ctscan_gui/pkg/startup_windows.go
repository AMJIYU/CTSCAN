//go:build windows

package pkg

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const startupRegistryReadAccess = registry.READ

var (
	versionDLL                   = windows.NewLazySystemDLL("version.dll")
	procGetFileVersionInfoSizeW  = versionDLL.NewProc("GetFileVersionInfoSizeW")
	procGetFileVersionInfoW      = versionDLL.NewProc("GetFileVersionInfoW")
	procVerQueryValueW           = versionDLL.NewProc("VerQueryValueW")
	windowsVersionFieldLanguages = []string{"040904b0", "040904e4", "080404b0"}
)

type windowsStartupCollector struct {
	items []StartupItem
	seen  map[string]struct{}
}

type windowsRegistryAutorunSpec struct {
	root       registry.Key
	rootName   string
	subKey     string
	category   string
	entryType  string
	valueNames []string
	allValues  bool
}

type scheduledTaskXML struct {
	RegistrationInfo struct {
		Description string `xml:"Description"`
		Author      string `xml:"Author"`
	} `xml:"RegistrationInfo"`
	Triggers struct {
		Boot                 []scheduledTaskTrigger `xml:"BootTrigger"`
		Logon                []scheduledTaskTrigger `xml:"LogonTrigger"`
		SessionStateChange   []scheduledTaskTrigger `xml:"SessionStateChangeTrigger"`
		Time                 []scheduledTaskTrigger `xml:"TimeTrigger"`
		Calendar             []scheduledTaskTrigger `xml:"CalendarTrigger"`
		Idle                 []scheduledTaskTrigger `xml:"IdleTrigger"`
		Registration         []scheduledTaskTrigger `xml:"RegistrationTrigger"`
		Event                []scheduledTaskTrigger `xml:"EventTrigger"`
		CustomTrigger        []scheduledTaskTrigger `xml:"CustomTrigger"`
		CustomTriggerOnEvent []scheduledTaskTrigger `xml:"CustomTriggerOnEvent"`
	} `xml:"Triggers"`
	Actions struct {
		Exec []struct {
			Command          string `xml:"Command"`
			Arguments        string `xml:"Arguments"`
			WorkingDirectory string `xml:"WorkingDirectory"`
		} `xml:"Exec"`
		ComHandler []struct {
			ClassID string `xml:"ClassId"`
			Data    string `xml:"Data"`
		} `xml:"ComHandler"`
	} `xml:"Actions"`
	Settings struct {
		Enabled string `xml:"Enabled"`
		Hidden  string `xml:"Hidden"`
	} `xml:"Settings"`
}

type scheduledTaskTrigger struct {
	Enabled string `xml:"Enabled"`
}

type wmiEventConsumer struct {
	Class               string `json:"Class"`
	Name                string `json:"Name"`
	CommandLineTemplate string `json:"CommandLineTemplate"`
	ExecutablePath      string `json:"ExecutablePath"`
	ScriptText          string `json:"ScriptText"`
}

func (a *App) getWindowsStartupItems() []StartupItem {
	collector := &windowsStartupCollector{
		items: make([]StartupItem, 0, 256),
		seen:  make(map[string]struct{}, 256),
	}

	collectWindowsStartupFolders(collector)
	collectWindowsRegistryAutoruns(collector)
	collectWindowsLoadedUserRunKeys(collector)
	collectWindowsServicesAndDrivers(collector)
	collectWindowsScheduledTasks(collector)
	collectWindowsWMIEventConsumers(collector)

	sort.SliceStable(collector.items, func(i, j int) bool {
		if collector.items[i].Category == collector.items[j].Category {
			return strings.ToLower(collector.items[i].Name) < strings.ToLower(collector.items[j].Name)
		}
		return startupCategoryRank(collector.items[i].Category) < startupCategoryRank(collector.items[j].Category)
	})

	return collector.items
}

func collectWindowsStartupFolders(collector *windowsStartupCollector) {
	folders := []struct {
		path      string
		entryType string
	}{
		{filepath.Join(os.Getenv("APPDATA"), `Microsoft\Windows\Start Menu\Programs\Startup`), "当前用户 Startup 文件夹"},
		{filepath.Join(os.Getenv("ProgramData"), `Microsoft\Windows\Start Menu\Programs\Startup`), "所有用户 Startup 文件夹"},
	}

	if systemDrive := os.Getenv("SystemDrive"); systemDrive != "" {
		userRoot := filepath.Join(systemDrive+`\`, "Users")
		if entries, err := os.ReadDir(userRoot); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				folder := filepath.Join(userRoot, entry.Name(), `AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`)
				folders = append(folders, struct {
					path      string
					entryType string
				}{folder, "用户 Startup 文件夹"})
			}
		}
	}

	for _, folder := range folders {
		if strings.TrimSpace(folder.path) == "" {
			continue
		}
		entries, err := os.ReadDir(folder.path)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			filePath := filepath.Join(folder.path, entry.Name())
			info, _ := entry.Info()
			item := StartupItem{
				Name:        entry.Name(),
				Path:        filePath,
				Type:        folder.entryType,
				Enabled:     true,
				Description: "启动文件夹项目；若为 .lnk，请在属性中查看目标程序",
				Category:    "Logon",
				Location:    folder.path,
				ImagePath:   filePath,
			}
			if info != nil {
				item.LastModTime = info.ModTime()
				item.Size = info.Size()
			}
			applyStartupFileInfo(&item, filePath)
			collector.add(item)
		}
	}
}

func collectWindowsRegistryAutoruns(collector *windowsStartupCollector) {
	specs := []windowsRegistryAutorunSpec{
		{registry.CURRENT_USER, "HKCU", `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, "Logon", "HKCU Run", nil, true},
		{registry.CURRENT_USER, "HKCU", `SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce`, "Logon", "HKCU RunOnce", nil, true},
		{registry.CURRENT_USER, "HKCU", `SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnceEx`, "Logon", "HKCU RunOnceEx", nil, true},
		{registry.CURRENT_USER, "HKCU", `SOFTWARE\Microsoft\Windows\CurrentVersion\RunServices`, "Logon", "HKCU RunServices", nil, true},
		{registry.CURRENT_USER, "HKCU", `SOFTWARE\Microsoft\Windows\CurrentVersion\RunServicesOnce`, "Logon", "HKCU RunServicesOnce", nil, true},
		{registry.CURRENT_USER, "HKCU", `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer\Run`, "Logon", "HKCU Policies Explorer Run", nil, true},
		{registry.CURRENT_USER, "HKCU", `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Windows`, "Logon", "HKCU Windows", []string{"Load", "Run"}, false},

		{registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, "Logon", "HKLM Run", nil, true},
		{registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce`, "Logon", "HKLM RunOnce", nil, true},
		{registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnceEx`, "Logon", "HKLM RunOnceEx", nil, true},
		{registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows\CurrentVersion\RunServices`, "Logon", "HKLM RunServices", nil, true},
		{registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows\CurrentVersion\RunServicesOnce`, "Logon", "HKLM RunServicesOnce", nil, true},
		{registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer\Run`, "Logon", "HKLM Policies Explorer Run", nil, true},
		{registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Windows`, "AppInit", "AppInit/Windows", []string{"AppInit_DLLs", "LoadAppInit_DLLs"}, false},
		{registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Control\Session Manager\AppCertDlls`, "AppInit", "AppCertDlls", nil, true},
		{registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`, "Winlogon", "Winlogon", []string{"Shell", "Userinit", "VMApplet", "Taskman", "GinaDLL", "AppSetup"}, false},
		{registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Control\Session Manager`, "Boot Execute", "Session Manager", []string{"BootExecute", "SetupExecute", "Execute"}, false},
		{registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Control\SafeBoot\AlternateShell`, "Boot Execute", "SafeBoot AlternateShell", []string{""}, false},
		{registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Control\Lsa`, "LSA Providers", "LSA Providers", []string{"Authentication Packages", "Notification Packages", "Security Packages"}, false},
		{registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Control\Lsa\OSConfig`, "LSA Providers", "LSA OSConfig", []string{"Security Packages"}, false},
		{registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Control\SecurityProviders`, "LSA Providers", "Security Providers", []string{"SecurityProviders"}, false},
		{registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Control\NetworkProvider\Order`, "Network Providers", "Network Provider Order", []string{"ProviderOrder"}, false},
		{registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Control\Session Manager\KnownDLLs`, "Known DLLs", "KnownDLLs", nil, true},
		{registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Drivers32`, "Codecs", "Drivers32", nil, true},
		{registry.CURRENT_USER, "HKCU", `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Drivers32`, "Codecs", "HKCU Drivers32", nil, true},
		{registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows NT\CurrentVersion\drivers.desc`, "Codecs", "Drivers Description", nil, true},
	}

	for _, spec := range specs {
		collectWindowsRegistryValues(collector, spec)
	}

	collectWindowsRegistrySubkeyValues(collector, registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\Notify`, "Winlogon", "Winlogon Notify", []string{"DLLName"})
	collectWindowsRegistrySubkeyAllValues(collector, registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnceEx`, "Logon", "HKLM RunOnceEx")
	collectWindowsRegistrySubkeyAllValues(collector, registry.CURRENT_USER, "HKCU", `SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnceEx`, "Logon", "HKCU RunOnceEx")
	collectWindowsRegistrySubkeyValues(collector, registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Image File Execution Options`, "Image Hijacks", "IFEO Debugger", []string{"Debugger"})
	collectWindowsRegistrySubkeyValues(collector, registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Windows NT\CurrentVersion\SilentProcessExit`, "Image Hijacks", "SilentProcessExit", []string{"MonitorProcess"})
	collectWindowsRegistrySubkeyValues(collector, registry.LOCAL_MACHINE, "HKLM", `SOFTWARE\Microsoft\Active Setup\Installed Components`, "Active Setup", "HKLM Active Setup", []string{"StubPath"})
	collectWindowsRegistrySubkeyValues(collector, registry.CURRENT_USER, "HKCU", `SOFTWARE\Microsoft\Active Setup\Installed Components`, "Active Setup", "HKCU Active Setup", []string{"StubPath"})
	collectWindowsRegistrySubkeyValues(collector, registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Control\Print\Monitors`, "Print Monitors", "Print Monitor", []string{"Driver"})
	collectWindowsRegistrySubkeyValues(collector, registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Services`, "Network Providers", "Network Provider", []string{`NetworkProvider\ProviderPath`})
	collectWindowsRegistrySubkeyValues(collector, registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Services\WinSock2\Parameters\Protocol_Catalog9\Catalog_Entries`, "Winsock Providers", "Winsock Catalog", []string{"PackedCatalogItem", "LibraryPath"})
	collectWindowsRegistrySubkeyValues(collector, registry.LOCAL_MACHINE, "HKLM", `SYSTEM\CurrentControlSet\Services\WinSock2\Parameters\Protocol_Catalog9\Catalog_Entries64`, "Winsock Providers", "Winsock Catalog 64-bit", []string{"PackedCatalogItem", "LibraryPath"})
	collectWindowsCOMAutoruns(collector)
	collectWindowsOfficeAddins(collector)
}

func collectWindowsCOMAutoruns(collector *windowsStartupCollector) {
	roots := []struct {
		root     registry.Key
		rootName string
	}{
		{registry.LOCAL_MACHINE, "HKLM"},
		{registry.CURRENT_USER, "HKCU"},
	}

	comValueKeys := []struct {
		subKey    string
		category  string
		entryType string
	}{
		{`SOFTWARE\Microsoft\Windows\CurrentVersion\ShellServiceObjectDelayLoad`, "Explorer", "ShellServiceObjectDelayLoad"},
		{`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\ShellExecuteHooks`, "Explorer", "ShellExecuteHooks"},
		{`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\SharedTaskScheduler`, "Explorer", "SharedTaskScheduler"},
		{`SOFTWARE\Microsoft\Internet Explorer\Toolbar`, "Internet Explorer", "IE Toolbar"},
		{`SOFTWARE\Microsoft\Internet Explorer\URLSearchHooks`, "Internet Explorer", "IE URLSearchHooks"},
	}
	for _, rootInfo := range roots {
		for _, item := range comValueKeys {
			collectWindowsCOMRegistryValues(collector, rootInfo.root, rootInfo.rootName, item.subKey, item.category, item.entryType)
		}
	}

	comSubkeyKeys := []struct {
		subKey    string
		category  string
		entryType string
	}{
		{`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Browser Helper Objects`, "Internet Explorer", "Browser Helper Object"},
		{`SOFTWARE\Microsoft\Internet Explorer\Explorer Bars`, "Internet Explorer", "IE Explorer Bar"},
	}
	for _, rootInfo := range roots {
		for _, item := range comSubkeyKeys {
			collectWindowsCOMSubkeys(collector, rootInfo.root, rootInfo.rootName, item.subKey, item.category, item.entryType)
		}
		collectWindowsIESubkeyExtensions(collector, rootInfo.root, rootInfo.rootName, `SOFTWARE\Microsoft\Internet Explorer\Extensions`)
	}

	shellHandlerKeys := []struct {
		subKey    string
		entryType string
	}{
		{`SOFTWARE\Classes\*\shellex\ContextMenuHandlers`, "ContextMenuHandler"},
		{`SOFTWARE\Classes\AllFileSystemObjects\shellex\ContextMenuHandlers`, "ContextMenuHandler"},
		{`SOFTWARE\Classes\Directory\shellex\ContextMenuHandlers`, "Directory ContextMenuHandler"},
		{`SOFTWARE\Classes\Directory\Background\shellex\ContextMenuHandlers`, "Directory Background ContextMenuHandler"},
		{`SOFTWARE\Classes\Drive\shellex\ContextMenuHandlers`, "Drive ContextMenuHandler"},
		{`SOFTWARE\Classes\Folder\shellex\ContextMenuHandlers`, "Folder ContextMenuHandler"},
		{`SOFTWARE\Classes\*\shellex\CopyHookHandlers`, "CopyHookHandler"},
		{`SOFTWARE\Classes\Directory\shellex\CopyHookHandlers`, "Directory CopyHookHandler"},
		{`SOFTWARE\Classes\Directory\shellex\DragDropHandlers`, "Directory DragDropHandler"},
		{`SOFTWARE\Classes\Folder\shellex\ColumnHandlers`, "ColumnHandler"},
		{`SOFTWARE\Classes\lnkfile\shellex\ContextMenuHandlers`, "Shortcut ContextMenuHandler"},
		{`SOFTWARE\Classes\exefile\shellex\ContextMenuHandlers`, "Executable ContextMenuHandler"},
	}
	for _, rootInfo := range roots {
		for _, item := range shellHandlerKeys {
			collectWindowsCOMHandlerSubkeys(collector, rootInfo.root, rootInfo.rootName, item.subKey, "Explorer", item.entryType)
		}
	}
}

func collectWindowsCOMRegistryValues(collector *windowsStartupCollector, root registry.Key, rootName, subKey, category, entryType string) {
	for _, access := range registryAccessVariants(root, subKey) {
		key, err := registry.OpenKey(root, subKey, access)
		if err != nil {
			continue
		}
		valueNames, err := key.ReadValueNames(-1)
		if err != nil {
			key.Close()
			continue
		}

		for _, valueName := range valueNames {
			value, ok := getRegistryValueAsString(key, valueName)
			if !ok || strings.TrimSpace(value) == "" {
				continue
			}
			location := fmt.Sprintf(`%s\%s`, rootName, subKey)
			addWindowsCOMAutorunItem(collector, category, firstNonEmpty(valueName, value), entryType, location, value, value)
		}
		key.Close()
	}
}

func collectWindowsCOMSubkeys(collector *windowsStartupCollector, root registry.Key, rootName, baseKey, category, entryType string) {
	for _, access := range registryAccessVariants(root, baseKey) {
		base, err := registry.OpenKey(root, baseKey, access)
		if err != nil {
			continue
		}
		subKeys, err := base.ReadSubKeyNames(-1)
		base.Close()
		if err != nil {
			continue
		}

		for _, subKeyName := range subKeys {
			fullKey := baseKey + `\` + subKeyName
			location := fmt.Sprintf(`%s\%s`, rootName, fullKey)
			addWindowsCOMAutorunItem(collector, category, subKeyName, entryType, location, subKeyName, "")
		}
	}
}

func collectWindowsCOMHandlerSubkeys(collector *windowsStartupCollector, root registry.Key, rootName, baseKey, category, entryType string) {
	for _, access := range registryAccessVariants(root, baseKey) {
		base, err := registry.OpenKey(root, baseKey, access)
		if err != nil {
			continue
		}
		subKeys, err := base.ReadSubKeyNames(-1)
		base.Close()
		if err != nil {
			continue
		}

		for _, subKeyName := range subKeys {
			fullKey := baseKey + `\` + subKeyName
			key, err := registry.OpenKey(root, fullKey, access)
			if err != nil {
				continue
			}
			defaultValue, _ := getRegistryValueAsString(key, "")
			key.Close()

			location := fmt.Sprintf(`%s\%s`, rootName, fullKey)
			addWindowsCOMAutorunItem(collector, category, subKeyName, entryType, location, firstNonEmpty(defaultValue, subKeyName), defaultValue)
		}
	}
}

func collectWindowsIESubkeyExtensions(collector *windowsStartupCollector, root registry.Key, rootName, baseKey string) {
	for _, access := range registryAccessVariants(root, baseKey) {
		base, err := registry.OpenKey(root, baseKey, access)
		if err != nil {
			continue
		}
		subKeys, err := base.ReadSubKeyNames(-1)
		base.Close()
		if err != nil {
			continue
		}

		for _, subKeyName := range subKeys {
			fullKey := baseKey + `\` + subKeyName
			key, err := registry.OpenKey(root, fullKey, access)
			if err != nil {
				continue
			}
			execValue, _ := getRegistryValueAsString(key, "Exec")
			scriptValue, _ := getRegistryValueAsString(key, "Script")
			clsidValue, _ := getRegistryValueAsString(key, "CLSID")
			clsidExtValue, _ := getRegistryValueAsString(key, "ClsidExtension")
			buttonText, _ := getRegistryValueAsString(key, "ButtonText")
			menuText, _ := getRegistryValueAsString(key, "MenuText")
			key.Close()

			command := firstNonEmpty(execValue, scriptValue, resolveWindowsCOMServerPath(firstNonEmpty(clsidValue, clsidExtValue)), clsidValue, clsidExtValue, subKeyName)
			location := fmt.Sprintf(`%s\%s`, rootName, fullKey)
			collector.add(newWindowsAutorunItem("Internet Explorer", firstNonEmpty(buttonText, menuText, subKeyName), "IE Extension", location, command, true, firstNonEmpty(menuText, buttonText)))
		}
	}
}

func collectWindowsOfficeAddins(collector *windowsStartupCollector) {
	roots := []struct {
		root     registry.Key
		rootName string
	}{
		{registry.LOCAL_MACHINE, "HKLM"},
		{registry.CURRENT_USER, "HKCU"},
	}
	for _, rootInfo := range roots {
		collectWindowsOfficeAddinsFromRoot(collector, rootInfo.root, rootInfo.rootName)
	}
}

func collectWindowsOfficeAddinsFromRoot(collector *windowsStartupCollector, root registry.Key, rootName string) {
	const officeRoot = `SOFTWARE\Microsoft\Office`

	for _, access := range registryAccessVariants(root, officeRoot) {
		office, err := registry.OpenKey(root, officeRoot, access)
		if err != nil {
			continue
		}
		versions, err := office.ReadSubKeyNames(-1)
		office.Close()
		if err != nil {
			continue
		}

		for _, version := range versions {
			versionKey := officeRoot + `\` + version
			versionHandle, err := registry.OpenKey(root, versionKey, access)
			if err != nil {
				continue
			}
			appNames, err := versionHandle.ReadSubKeyNames(-1)
			versionHandle.Close()
			if err != nil {
				continue
			}

			for _, appName := range appNames {
				addinsKey := versionKey + `\` + appName + `\Addins`
				addinsHandle, err := registry.OpenKey(root, addinsKey, access)
				if err != nil {
					continue
				}
				addinNames, err := addinsHandle.ReadSubKeyNames(-1)
				addinsHandle.Close()
				if err != nil {
					continue
				}

				for _, addinName := range addinNames {
					collectWindowsOfficeAddin(collector, root, rootName, addinsKey+`\`+addinName, version, appName, addinName, access)
				}
			}
		}
	}
}

func collectWindowsOfficeAddin(collector *windowsStartupCollector, root registry.Key, rootName, fullKey, version, appName, addinName string, access uint32) {
	key, err := registry.OpenKey(root, fullKey, access)
	if err != nil {
		return
	}
	defer key.Close()

	friendlyName, _ := getRegistryValueAsString(key, "FriendlyName")
	description, _ := getRegistryValueAsString(key, "Description")
	manifest, _ := getRegistryValueAsString(key, "Manifest")
	loadBehavior := getRegistryDWORD(key, "LoadBehavior", 0)
	clsid := resolveWindowsProgIDCLSID(addinName)
	command := firstNonEmpty(manifest, resolveWindowsCOMServerPath(clsid), clsid, addinName)
	enabled := loadBehavior != 0
	entryType := strings.TrimSpace(appName + " Addin")
	if version != "" {
		entryType += " " + version
	}
	location := fmt.Sprintf(`%s\%s`, rootName, fullKey)
	collector.add(newWindowsAutorunItem("Office", firstNonEmpty(friendlyName, addinName), entryType, location, command, enabled, description))
}

func collectWindowsRegistrySubkeyAllValues(collector *windowsStartupCollector, root registry.Key, rootName, baseKey, category, entryType string) {
	for _, access := range registryAccessVariants(root, baseKey) {
		base, err := registry.OpenKey(root, baseKey, access)
		if err != nil {
			continue
		}
		subKeys, err := base.ReadSubKeyNames(-1)
		base.Close()
		if err != nil {
			continue
		}

		for _, subKeyName := range subKeys {
			fullKey := baseKey + `\` + subKeyName
			key, err := registry.OpenKey(root, fullKey, access)
			if err != nil {
				continue
			}
			valueNames, err := key.ReadValueNames(-1)
			if err != nil {
				key.Close()
				continue
			}
			for _, valueName := range valueNames {
				value, ok := getRegistryValueAsString(key, valueName)
				if !ok || strings.TrimSpace(value) == "" {
					continue
				}
				name := subKeyName
				if valueName != "" {
					name += `\` + valueName
				}
				location := fmt.Sprintf(`%s\%s`, rootName, fullKey)
				collector.add(newWindowsAutorunItem(category, name, entryType, location, value, true, "注册表子项启动/加载位置"))
			}
			key.Close()
		}
	}
}

func collectWindowsLoadedUserRunKeys(collector *windowsStartupCollector) {
	hku, err := registry.OpenKey(registry.USERS, ``, startupRegistryReadAccess)
	if err != nil {
		return
	}
	defer hku.Close()

	users, err := hku.ReadSubKeyNames(-1)
	if err != nil {
		return
	}

	relativeKeys := []struct {
		key       string
		entryType string
	}{
		{`SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, "HKU Run"},
		{`SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce`, "HKU RunOnce"},
		{`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer\Run`, "HKU Policies Explorer Run"},
	}

	for _, sid := range users {
		if strings.HasSuffix(strings.ToUpper(sid), "_CLASSES") {
			continue
		}
		for _, relative := range relativeKeys {
			collectWindowsRegistryValues(collector, windowsRegistryAutorunSpec{
				root:      registry.USERS,
				rootName:  "HKU",
				subKey:    sid + `\` + relative.key,
				category:  "Logon",
				entryType: relative.entryType,
				allValues: true,
			})
		}
	}
}

func collectWindowsRegistryValues(collector *windowsStartupCollector, spec windowsRegistryAutorunSpec) {
	for _, access := range registryAccessVariants(spec.root, spec.subKey) {
		key, err := registry.OpenKey(spec.root, spec.subKey, access)
		if err != nil {
			continue
		}

		valueNames := spec.valueNames
		if spec.allValues {
			valueNames, err = key.ReadValueNames(-1)
			if err != nil {
				key.Close()
				continue
			}
		}

		for _, valueName := range valueNames {
			value, ok := getRegistryValueAsString(key, valueName)
			if !ok || strings.TrimSpace(value) == "" {
				continue
			}
			name := valueName
			if name == "" {
				name = "(Default)"
			}
			location := fmt.Sprintf(`%s\%s`, spec.rootName, spec.subKey)
			collector.add(newWindowsAutorunItem(spec.category, name, spec.entryType, location, value, true, "注册表启动项"))
		}
		key.Close()
	}
}

func collectWindowsRegistrySubkeyValues(collector *windowsStartupCollector, root registry.Key, rootName, baseKey, category, entryType string, valueNames []string) {
	for _, access := range registryAccessVariants(root, baseKey) {
		base, err := registry.OpenKey(root, baseKey, access)
		if err != nil {
			continue
		}
		subKeys, err := base.ReadSubKeyNames(-1)
		base.Close()
		if err != nil {
			continue
		}

		for _, subKeyName := range subKeys {
			fullKey := baseKey + `\` + subKeyName
			key, err := registry.OpenKey(root, fullKey, access)
			if err != nil {
				continue
			}
			for _, valueName := range valueNames {
				value, ok := getRegistryValueByPath(key, valueName)
				if !ok || strings.TrimSpace(value) == "" {
					continue
				}
				location := fmt.Sprintf(`%s\%s`, rootName, fullKey)
				name := subKeyName
				if valueName != "" && !strings.Contains(valueName, `\`) {
					name += `\` + valueName
				}
				collector.add(newWindowsAutorunItem(category, name, entryType, location, value, true, "注册表子项启动/加载位置"))
			}
			key.Close()
		}
	}
}

func collectWindowsServicesAndDrivers(collector *windowsStartupCollector) {
	base, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services`, startupRegistryReadAccess)
	if err != nil {
		return
	}
	defer base.Close()

	serviceNames, err := base.ReadSubKeyNames(-1)
	if err != nil {
		return
	}

	for _, serviceName := range serviceNames {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\`+serviceName, startupRegistryReadAccess)
		if err != nil {
			continue
		}

		startType := getRegistryDWORD(key, "Start", 3)
		if startType > 2 {
			key.Close()
			continue
		}
		serviceType := getRegistryDWORD(key, "Type", 0)
		displayName, _, _ := key.GetStringValue("DisplayName")
		description, _, _ := key.GetStringValue("Description")
		imagePath, _, _ := key.GetStringValue("ImagePath")
		if strings.TrimSpace(imagePath) == "" {
			if parameterKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\`+serviceName+`\Parameters`, startupRegistryReadAccess); err == nil {
				imagePath, _, _ = parameterKey.GetStringValue("ServiceDll")
				parameterKey.Close()
			}
		}
		key.Close()

		category := "Services"
		if serviceType&0x3 != 0 {
			category = "Drivers"
		}

		name := serviceName
		if strings.TrimSpace(displayName) != "" {
			name = strings.TrimSpace(displayName)
		}
		entryType := windowsServiceStartType(startType)
		if category == "Drivers" {
			entryType += " Driver"
		} else {
			entryType += " Service"
		}
		location := `HKLM\SYSTEM\CurrentControlSet\Services\` + serviceName
		item := newWindowsAutorunItem(category, name, entryType, location, imagePath, true, strings.TrimSpace(description))
		if item.ImagePath == "" && category == "Drivers" && serviceName != "" {
			item.ImagePath = `%SystemRoot%\System32\drivers\` + serviceName + ".sys"
			item.Path = expandWindowsEnvironmentStrings(item.ImagePath)
			applyStartupFileInfo(&item, item.Path)
		}
		collector.add(item)
	}
}

func collectWindowsScheduledTasks(collector *windowsStartupCollector) {
	systemRoot := os.Getenv("SystemRoot")
	if systemRoot == "" {
		systemRoot = `C:\Windows`
	}
	taskRoot := filepath.Join(systemRoot, `System32\Tasks`)

	_ = filepath.WalkDir(taskRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry == nil || entry.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil || len(data) == 0 {
			return nil
		}

		var task scheduledTaskXML
		if err := xml.Unmarshal(data, &task); err != nil {
			return nil
		}

		actions := scheduledTaskActions(task)
		if len(actions) == 0 {
			return nil
		}
		triggers := scheduledTaskTriggers(task)
		if len(triggers) == 0 {
			triggers = []string{"无显式触发器"}
		}

		relName, err := filepath.Rel(taskRoot, path)
		if err != nil {
			relName = path
		}
		relName = strings.TrimPrefix(filepath.ToSlash(relName), "/")
		relName = strings.ReplaceAll(relName, "/", `\`)

		info, _ := entry.Info()
		enabled := !strings.EqualFold(strings.TrimSpace(task.Settings.Enabled), "false")
		for _, action := range actions {
			item := newWindowsAutorunItem("Scheduled Tasks", relName, strings.Join(triggers, ", "), path, action, enabled, strings.TrimSpace(task.RegistrationInfo.Description))
			item.Publisher = strings.TrimSpace(task.RegistrationInfo.Author)
			if info != nil {
				item.LastModTime = info.ModTime()
				item.Size = info.Size()
			}
			collector.add(item)
		}
		return nil
	})
}

func collectWindowsWMIEventConsumers(collector *windowsStartupCollector) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	script := `$ErrorActionPreference='SilentlyContinue'; Get-CimInstance -Namespace root/subscription -ClassName __EventConsumer | ForEach-Object { [pscustomobject]@{ Class=$_.CimClass.CimClassName; Name=$_.Name; CommandLineTemplate=$_.CommandLineTemplate; ExecutablePath=$_.ExecutablePath; ScriptText=$_.ScriptText } } | ConvertTo-Json -Compress`
	output, err := backgroundCommandContext(ctx, "powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script).Output()
	if err != nil || len(output) == 0 {
		return
	}

	consumers := parseWMIConsumers(output)
	for _, consumer := range consumers {
		command := strings.TrimSpace(consumer.ExecutablePath)
		if command == "" {
			command = strings.TrimSpace(consumer.CommandLineTemplate)
		}
		if command == "" && strings.TrimSpace(consumer.ScriptText) != "" {
			command = "WMI Script Consumer"
		}
		if command == "" {
			continue
		}

		description := strings.TrimSpace(consumer.CommandLineTemplate)
		if description == "" {
			description = truncateString(strings.TrimSpace(consumer.ScriptText), 500)
		}
		collector.add(newWindowsAutorunItem(
			"WMI",
			firstNonEmpty(consumer.Name, consumer.Class),
			consumer.Class,
			`root\subscription\__EventConsumer`,
			command,
			true,
			description,
		))
	}
}

func parseWMIConsumers(output []byte) []wmiEventConsumer {
	text := strings.TrimSpace(string(output))
	if text == "" || text == "null" {
		return nil
	}

	var list []wmiEventConsumer
	if err := json.Unmarshal([]byte(text), &list); err == nil {
		return list
	}

	var single wmiEventConsumer
	if err := json.Unmarshal([]byte(text), &single); err == nil {
		if single.Name != "" || single.Class != "" {
			return []wmiEventConsumer{single}
		}
	}
	return nil
}

func scheduledTaskActions(task scheduledTaskXML) []string {
	actions := make([]string, 0, len(task.Actions.Exec)+len(task.Actions.ComHandler))
	for _, action := range task.Actions.Exec {
		command := strings.TrimSpace(action.Command)
		if command == "" {
			continue
		}
		if args := strings.TrimSpace(action.Arguments); args != "" {
			command += " " + args
		}
		actions = append(actions, command)
	}
	for _, action := range task.Actions.ComHandler {
		if classID := strings.TrimSpace(action.ClassID); classID != "" {
			actions = append(actions, "COM Handler "+classID)
		}
	}
	return actions
}

func scheduledTaskTriggers(task scheduledTaskXML) []string {
	triggers := make([]string, 0, 8)
	add := func(label string, values []scheduledTaskTrigger) {
		for _, trigger := range values {
			if strings.EqualFold(strings.TrimSpace(trigger.Enabled), "false") {
				continue
			}
			triggers = append(triggers, label)
			return
		}
	}

	add("Boot", task.Triggers.Boot)
	add("Logon", task.Triggers.Logon)
	add("Session", task.Triggers.SessionStateChange)
	add("Time", task.Triggers.Time)
	add("Calendar", task.Triggers.Calendar)
	add("Idle", task.Triggers.Idle)
	add("Registration", task.Triggers.Registration)
	add("Event", task.Triggers.Event)
	add("Custom", task.Triggers.CustomTrigger)
	add("CustomEvent", task.Triggers.CustomTriggerOnEvent)
	return uniqueStrings(triggers)
}

func newWindowsAutorunItem(category, name, entryType, location, command string, enabled bool, description string) StartupItem {
	command = strings.TrimSpace(command)
	expandedCommand := expandWindowsEnvironmentStrings(command)
	imagePath := extractWindowsImagePath(expandedCommand)
	path := imagePath
	if path == "" {
		path = expandedCommand
	}

	item := StartupItem{
		Name:        strings.TrimSpace(name),
		Path:        path,
		Type:        strings.TrimSpace(entryType),
		Enabled:     enabled,
		Description: strings.TrimSpace(description),
		Category:    strings.TrimSpace(category),
		Location:    strings.TrimSpace(location),
		ImagePath:   imagePath,
	}
	if item.Description == "" && expandedCommand != "" && expandedCommand != imagePath {
		item.Description = expandedCommand
	}
	applyStartupFileInfo(&item, path)
	return item
}

func (collector *windowsStartupCollector) add(item StartupItem) {
	item.Name = strings.TrimSpace(item.Name)
	item.Category = strings.TrimSpace(item.Category)
	item.Type = strings.TrimSpace(item.Type)
	item.Path = strings.TrimSpace(item.Path)
	item.Location = strings.TrimSpace(item.Location)
	item.ImagePath = strings.TrimSpace(item.ImagePath)
	if item.Name == "" {
		item.Name = item.Path
	}
	if item.Category == "" {
		item.Category = "Other"
	}

	key := strings.ToLower(strings.Join([]string{item.Category, item.Type, item.Location, item.Name, item.Path}, "|"))
	if _, ok := collector.seen[key]; ok {
		return
	}
	collector.seen[key] = struct{}{}
	collector.items = append(collector.items, item)
}

func getRegistryValueAsString(key registry.Key, valueName string) (string, bool) {
	if value, _, err := key.GetStringValue(valueName); err == nil {
		return strings.TrimSpace(value), true
	}
	if values, _, err := key.GetStringsValue(valueName); err == nil {
		return strings.TrimSpace(strings.Join(values, " ")), true
	}
	if value, _, err := key.GetBinaryValue(valueName); err == nil {
		text := extractWindowsRegistryBinaryText(value)
		if text != "" {
			return text, true
		}
	}
	if value, _, err := key.GetIntegerValue(valueName); err == nil {
		return strconv.FormatUint(value, 10), true
	}
	return "", false
}

func addWindowsCOMAutorunItem(collector *windowsStartupCollector, category, name, entryType, location, clsidOrCommand, description string) {
	clsid := firstNonEmpty(findWindowsCLSID(clsidOrCommand), findWindowsCLSID(name))
	displayName := resolveWindowsCOMDisplayName(clsid)
	serverPath := resolveWindowsCOMServerPath(clsid)
	command := firstNonEmpty(serverPath, clsidOrCommand)
	itemName := firstNonEmpty(name, displayName, clsid, clsidOrCommand)
	itemDescription := strings.TrimSpace(description)
	if itemDescription == "" || strings.EqualFold(itemDescription, clsid) {
		itemDescription = displayName
	}

	collector.add(newWindowsAutorunItem(category, itemName, entryType, location, command, true, itemDescription))
}

func resolveWindowsCOMServerPath(clsid string) string {
	clsid = findWindowsCLSID(clsid)
	if clsid == "" {
		return ""
	}

	baseKeys := []struct {
		root   registry.Key
		subKey string
	}{
		{registry.CURRENT_USER, `SOFTWARE\Classes\CLSID\` + clsid},
		{registry.LOCAL_MACHINE, `SOFTWARE\Classes\CLSID\` + clsid},
		{registry.CLASSES_ROOT, `CLSID\` + clsid},
	}
	serverSubKeys := []string{"InprocServer32", "LocalServer32", "InprocHandler32"}

	for _, base := range baseKeys {
		for _, serverSubKey := range serverSubKeys {
			if value, ok := readWindowsRegistryString(base.root, base.subKey+`\`+serverSubKey, ""); ok {
				return value
			}
		}
	}
	return ""
}

func resolveWindowsCOMDisplayName(clsid string) string {
	clsid = findWindowsCLSID(clsid)
	if clsid == "" {
		return ""
	}

	baseKeys := []struct {
		root   registry.Key
		subKey string
	}{
		{registry.CURRENT_USER, `SOFTWARE\Classes\CLSID\` + clsid},
		{registry.LOCAL_MACHINE, `SOFTWARE\Classes\CLSID\` + clsid},
		{registry.CLASSES_ROOT, `CLSID\` + clsid},
	}
	for _, base := range baseKeys {
		if value, ok := readWindowsRegistryString(base.root, base.subKey, ""); ok {
			return value
		}
	}
	return ""
}

func resolveWindowsProgIDCLSID(progID string) string {
	progID = strings.TrimSpace(progID)
	if progID == "" {
		return ""
	}

	baseKeys := []struct {
		root   registry.Key
		subKey string
	}{
		{registry.CURRENT_USER, `SOFTWARE\Classes\` + progID + `\CLSID`},
		{registry.LOCAL_MACHINE, `SOFTWARE\Classes\` + progID + `\CLSID`},
		{registry.CLASSES_ROOT, progID + `\CLSID`},
	}
	for _, base := range baseKeys {
		if value, ok := readWindowsRegistryString(base.root, base.subKey, ""); ok {
			return findWindowsCLSID(value)
		}
	}
	return ""
}

func readWindowsRegistryString(root registry.Key, subKey, valueName string) (string, bool) {
	for _, access := range registryAccessVariants(root, subKey) {
		key, err := registry.OpenKey(root, subKey, access)
		if err != nil {
			continue
		}
		value, ok := getRegistryValueAsString(key, valueName)
		key.Close()
		if ok && strings.TrimSpace(value) != "" {
			return value, true
		}
	}
	return "", false
}

func findWindowsCLSID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if normalized := normalizeWindowsGUID(value); normalized != "" {
		return normalized
	}

	for offset := 0; offset < len(value); {
		left := strings.Index(value[offset:], "{")
		if left < 0 {
			return ""
		}
		left += offset
		right := strings.Index(value[left:], "}")
		if right < 0 {
			return ""
		}
		right += left
		if normalized := normalizeWindowsGUID(value[left : right+1]); normalized != "" {
			return normalized
		}
		offset = right + 1
	}
	return ""
}

func normalizeWindowsGUID(value string) string {
	trimmed := strings.Trim(strings.TrimSpace(value), "{}")
	if len(trimmed) != 36 {
		return ""
	}
	for index, char := range trimmed {
		switch index {
		case 8, 13, 18, 23:
			if char != '-' {
				return ""
			}
		default:
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
				return ""
			}
		}
	}
	return "{" + strings.ToUpper(trimmed) + "}"
}

func extractWindowsRegistryBinaryText(data []byte) string {
	values := uniqueStrings(append(extractASCIITextFromBinary(data), extractUTF16TextFromBinary(data)...))
	return strings.Join(values, " ")
}

func extractASCIITextFromBinary(data []byte) []string {
	values := make([]string, 0)
	buffer := make([]byte, 0, 64)
	flush := func() {
		if len(buffer) >= 4 {
			value := strings.TrimSpace(string(buffer))
			if isUsefulRegistryBinaryText(value) {
				values = append(values, value)
			}
		}
		buffer = buffer[:0]
	}

	for _, char := range data {
		if (char >= 32 && char <= 126) || char >= 128 {
			buffer = append(buffer, char)
			continue
		}
		flush()
	}
	flush()
	return values
}

func extractUTF16TextFromBinary(data []byte) []string {
	values := make([]string, 0)
	buffer := make([]uint16, 0, 64)
	flush := func() {
		if len(buffer) >= 4 {
			value := strings.TrimSpace(string(utf16.Decode(buffer)))
			if isUsefulRegistryBinaryText(value) {
				values = append(values, value)
			}
		}
		buffer = buffer[:0]
	}

	for index := 0; index+1 < len(data); index += 2 {
		char := binary.LittleEndian.Uint16(data[index : index+2])
		if char >= 32 && char != 0xffff {
			buffer = append(buffer, char)
			continue
		}
		flush()
	}
	flush()
	return values
}

func isUsefulRegistryBinaryText(value string) bool {
	value = strings.TrimSpace(value)
	if len(value) < 4 {
		return false
	}
	lowerValue := strings.ToLower(value)
	return strings.ContainsAny(value, `\/:`) ||
		strings.Contains(lowerValue, ".dll") ||
		strings.Contains(lowerValue, ".exe") ||
		strings.Contains(lowerValue, ".sys") ||
		strings.Contains(lowerValue, ".ocx")
}

func getRegistryValueByPath(key registry.Key, valuePath string) (string, bool) {
	if !strings.Contains(valuePath, `\`) {
		return getRegistryValueAsString(key, valuePath)
	}
	parts := strings.Split(valuePath, `\`)
	if len(parts) < 2 {
		return "", false
	}
	subKeyPath := strings.Join(parts[:len(parts)-1], `\`)
	valueName := parts[len(parts)-1]
	subKey, err := registry.OpenKey(key, subKeyPath, startupRegistryReadAccess)
	if err != nil {
		return "", false
	}
	defer subKey.Close()
	return getRegistryValueAsString(subKey, valueName)
}

func getRegistryDWORD(key registry.Key, valueName string, fallback uint64) uint64 {
	value, _, err := key.GetIntegerValue(valueName)
	if err != nil {
		return fallback
	}
	return value
}

func registryAccessVariants(root registry.Key, subKey string) []uint32 {
	if (root == registry.LOCAL_MACHINE || root == registry.CURRENT_USER) && strings.HasPrefix(strings.ToUpper(subKey), `SOFTWARE\`) {
		return []uint32{
			startupRegistryReadAccess | registry.WOW64_64KEY,
			startupRegistryReadAccess | registry.WOW64_32KEY,
			startupRegistryReadAccess,
		}
	}
	return []uint32{startupRegistryReadAccess}
}

func windowsServiceStartType(startType uint64) string {
	switch startType {
	case 0:
		return "Boot Start"
	case 1:
		return "System Start"
	case 2:
		return "Auto Start"
	default:
		return "Start " + strconv.FormatUint(startType, 10)
	}
}

func expandWindowsEnvironmentStrings(value string) string {
	if strings.TrimSpace(value) == "" {
		return value
	}
	if strings.HasPrefix(strings.ToLower(value), `\systemroot\`) {
		if systemRoot := os.Getenv("SystemRoot"); systemRoot != "" {
			return systemRoot + value[len(`\SystemRoot`):]
		}
	}
	src, err := windows.UTF16PtrFromString(value)
	if err != nil {
		return value
	}
	required, _ := windows.ExpandEnvironmentStrings(src, nil, 0)
	if required == 0 {
		return value
	}
	buffer := make([]uint16, required)
	if expanded, err := windows.ExpandEnvironmentStrings(src, &buffer[0], required); err == nil && expanded > 0 {
		return windows.UTF16ToString(buffer)
	}
	return value
}

func extractWindowsImagePath(command string) string {
	command = strings.TrimSpace(strings.Trim(command, "\x00"))
	if command == "" {
		return ""
	}
	command = strings.TrimPrefix(command, `\??\`)
	if strings.HasPrefix(strings.ToLower(command), "com handler ") {
		return command
	}
	if strings.HasPrefix(command, `"`) {
		rest := strings.TrimPrefix(command, `"`)
		if index := strings.Index(rest, `"`); index >= 0 {
			return strings.TrimSpace(rest[:index])
		}
	}

	cleanedCommand := strings.Trim(command, `"'`)
	parts := strings.Fields(cleanedCommand)
	for i := range parts {
		candidate := strings.Trim(strings.Join(parts[:i+1], " "), `"' ,`)
		if hasExecutableLikeExtension(candidate) {
			return candidate
		}
	}
	if len(parts) > 0 {
		return strings.Trim(parts[0], `"' ,`)
	}
	return command
}

func hasExecutableLikeExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(strings.Trim(path, `"' ,`)))
	switch ext {
	case ".exe", ".dll", ".sys", ".com", ".bat", ".cmd", ".ps1", ".vbs", ".js", ".lnk", ".scr", ".ocx", ".cpl":
		return true
	default:
		return false
	}
}

func applyStartupFileInfo(item *StartupItem, path string) {
	path = strings.Trim(path, `"'`)
	if path == "" || strings.HasPrefix(strings.ToLower(path), "com handler ") {
		return
	}
	if !filepath.IsAbs(path) {
		if systemRoot := os.Getenv("SystemRoot"); systemRoot != "" {
			names := []string{path}
			if filepath.Ext(path) == "" {
				names = append(names, path+".dll", path+".exe", path+".sys")
			}
			candidates := make([]string, 0, len(names)*3)
			for _, name := range names {
				candidates = append(candidates,
					filepath.Join(systemRoot, name),
					filepath.Join(systemRoot, "System32", name),
					filepath.Join(systemRoot, "SysWOW64", name),
				)
			}
			for _, candidate := range candidates {
				if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
					item.LastModTime = info.ModTime()
					item.Size = info.Size()
					item.Path = candidate
					if item.ImagePath == path {
						item.ImagePath = candidate
					}
					enrichStartupItemFromVersionInfo(item, candidate)
					return
				}
			}
		}
		return
	}
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		item.LastModTime = info.ModTime()
		item.Size = info.Size()
		enrichStartupItemFromVersionInfo(item, path)
	}
}

func enrichStartupItemFromVersionInfo(item *StartupItem, path string) {
	if item.Publisher == "" {
		item.Publisher = readWindowsVersionString(path, "CompanyName")
	}
	if item.Description == "" {
		item.Description = readWindowsVersionString(path, "FileDescription")
	}
}

func readWindowsVersionString(path, field string) string {
	path = strings.Trim(path, `"'`)
	if path == "" {
		return ""
	}
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return ""
	}

	size, _, _ := procGetFileVersionInfoSizeW.Call(uintptr(unsafe.Pointer(pathPtr)), 0)
	if size == 0 {
		return ""
	}
	buffer := make([]byte, size)
	ok, _, _ := procGetFileVersionInfoW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		0,
		size,
		uintptr(unsafe.Pointer(&buffer[0])),
	)
	if ok == 0 {
		return ""
	}

	languages := append(readWindowsVersionLanguages(buffer), windowsVersionFieldLanguages...)
	for _, language := range uniqueStrings(languages) {
		value := queryWindowsVersionString(buffer, fmt.Sprintf(`\StringFileInfo\%s\%s`, language, field))
		if value != "" {
			return value
		}
	}
	return ""
}

func readWindowsVersionLanguages(buffer []byte) []string {
	valuePtr, length, ok := queryWindowsVersionValue(buffer, `\VarFileInfo\Translation`)
	if !ok || valuePtr == 0 || length < 4 {
		return nil
	}

	raw := (*[1 << 20]uint16)(unsafe.Pointer(valuePtr))[: int(length)/2 : int(length)/2]
	languages := make([]string, 0, len(raw)/2)
	for index := 0; index+1 < len(raw); index += 2 {
		language := raw[index]
		codePage := raw[index+1]
		if language == 0 || codePage == 0 {
			continue
		}
		languages = append(languages, fmt.Sprintf("%04x%04x", language, codePage))
	}
	return languages
}

func queryWindowsVersionString(buffer []byte, subBlock string) string {
	valuePtr, length, ok := queryWindowsVersionValue(buffer, subBlock)
	if !ok || valuePtr == 0 || length == 0 {
		return ""
	}
	raw := (*[1 << 20]uint16)(unsafe.Pointer(valuePtr))[:length:length]
	return strings.TrimSpace(windows.UTF16ToString(raw))
}

func queryWindowsVersionValue(buffer []byte, subBlock string) (uintptr, uint32, bool) {
	if len(buffer) == 0 {
		return 0, 0, false
	}
	subBlockPtr, err := windows.UTF16PtrFromString(subBlock)
	if err != nil {
		return 0, 0, false
	}
	var (
		valuePtr uintptr
		length   uint32
	)
	ok, _, _ := procVerQueryValueW.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(subBlockPtr)),
		uintptr(unsafe.Pointer(&valuePtr)),
		uintptr(unsafe.Pointer(&length)),
	)
	return valuePtr, length, ok != 0
}

func startupCategoryRank(category string) int {
	order := map[string]int{
		"Logon":             10,
		"Scheduled Tasks":   20,
		"Services":          30,
		"Drivers":           40,
		"Winlogon":          50,
		"Boot Execute":      60,
		"AppInit":           70,
		"Image Hijacks":     80,
		"Explorer":          90,
		"Internet Explorer": 100,
		"LSA Providers":     110,
		"Network Providers": 120,
		"Winsock Providers": 130,
		"Print Monitors":    140,
		"Known DLLs":        150,
		"Active Setup":      160,
		"Codecs":            170,
		"Office":            180,
		"WMI":               190,
	}
	if rank, ok := order[category]; ok {
		return rank
	}
	return 999
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		key := strings.ToLower(value)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, value)
	}
	return result
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func truncateString(value string, maxRunes int) string {
	runes := []rune(value)
	if len(runes) <= maxRunes {
		return value
	}
	return string(runes[:maxRunes]) + "..."
}
