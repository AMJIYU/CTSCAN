//go:build windows

package pkg

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	windowsFilterNormalAccount = 0x0002
	windowsMaxPreferredLength  = ^uint32(0)
)

type windowsUserInfo0 struct {
	Name *uint16
}

type windowsUserInfo10 struct {
	Name        *uint16
	Comment     *uint16
	UserComment *uint16
	FullName    *uint16
}

func (a *App) getWindowsUsers() []SystemUser {
	names, err := enumerateWindowsUserNames()
	if err != nil || len(names) == 0 {
		return fallbackWindowsCurrentUser()
	}

	users := make([]SystemUser, 0, len(names))
	for _, name := range names {
		sid, _, _, err := windows.LookupSID("", name)
		sidString := ""
		if err == nil && sid != nil {
			sidString = sid.String()
		}

		users = append(users, SystemUser{
			Username: name,
			Uid:      sidString,
			Gid:      "0",
			HomeDir:  resolveWindowsUserHomeDir(name, sidString),
			Name:     lookupWindowsUserFullName(name),
		})
	}

	return users
}

func enumerateWindowsUserNames() ([]string, error) {
	names := make([]string, 0, 16)
	var resumeHandle uint32

	for {
		var (
			buffer       *byte
			entriesRead  uint32
			totalEntries uint32
		)

		err := windows.NetUserEnum(
			nil,
			0,
			windowsFilterNormalAccount,
			&buffer,
			windowsMaxPreferredLength,
			&entriesRead,
			&totalEntries,
			&resumeHandle,
		)
		if buffer != nil {
			entries := unsafe.Slice((*windowsUserInfo0)(unsafe.Pointer(buffer)), int(entriesRead))
			for _, entry := range entries {
				if entry.Name == nil {
					continue
				}
				if name := strings.TrimSpace(windows.UTF16PtrToString(entry.Name)); name != "" {
					names = append(names, name)
				}
			}
			_ = windows.NetApiBufferFree(buffer)
		}

		if err == nil {
			break
		}
		if err == windows.ERROR_MORE_DATA {
			continue
		}
		return uniqueSortedWindowsUserNames(names), err
	}

	return uniqueSortedWindowsUserNames(names), nil
}

func lookupWindowsUserFullName(name string) string {
	namePtr, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return ""
	}

	var buffer *byte
	if err := windows.NetUserGetInfo(nil, namePtr, 10, &buffer); err != nil || buffer == nil {
		return ""
	}
	defer windows.NetApiBufferFree(buffer)

	info := (*windowsUserInfo10)(unsafe.Pointer(buffer))
	if info.FullName != nil {
		if fullName := strings.TrimSpace(windows.UTF16PtrToString(info.FullName)); fullName != "" {
			return fullName
		}
	}
	if info.Comment != nil {
		return strings.TrimSpace(windows.UTF16PtrToString(info.Comment))
	}
	return ""
}

func resolveWindowsUserHomeDir(username, sidString string) string {
	if sidString != "" {
		keyPath := `SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList\` + sidString
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
		if err == nil {
			defer key.Close()
			if profilePath, _, err := key.GetStringValue("ProfileImagePath"); err == nil {
				if expanded := strings.TrimSpace(os.ExpandEnv(profilePath)); expanded != "" {
					return expanded
				}
			}
		}
	}

	if strings.EqualFold(username, os.Getenv("USERNAME")) {
		if profile := strings.TrimSpace(os.Getenv("USERPROFILE")); profile != "" {
			return profile
		}
	}

	systemDrive := os.Getenv("SystemDrive")
	if systemDrive == "" {
		systemDrive = "C:"
	}
	candidate := filepath.Join(systemDrive+`\`, "Users", username)
	if info, err := os.Stat(candidate); err == nil && info.IsDir() {
		return candidate
	}
	return ""
}

func fallbackWindowsCurrentUser() []SystemUser {
	username := strings.TrimSpace(os.Getenv("USERNAME"))
	if username == "" {
		return nil
	}

	sidString := ""
	if sid, _, _, err := windows.LookupSID("", username); err == nil && sid != nil {
		sidString = sid.String()
	}

	return []SystemUser{
		{
			Username: username,
			Uid:      sidString,
			Gid:      "0",
			HomeDir:  resolveWindowsUserHomeDir(username, sidString),
			Name:     lookupWindowsUserFullName(username),
		},
	}
}

func uniqueSortedWindowsUserNames(names []string) []string {
	seen := make(map[string]struct{}, len(names))
	result := make([]string, 0, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}
