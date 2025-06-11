package pkg

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type Patch struct {
	Time        string `json:"time"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	KB          string `json:"kb"`
}

func (a *App) GetPatchRecords() []Patch {
	records := []Patch{}
	switch runtime.GOOS {
	case "windows":
		records = a.getWindowsPatchRecords()
	case "linux":
		out, err := exec.Command("grep", "upgrade", "/var/log/apt/history.log").CombinedOutput()
		if err != nil {
			return records
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			records = append(records, Patch{Time: line[:15]})
		}
	case "darwin":
		out, err := exec.Command("softwareupdate", "--history").CombinedOutput()
		if err != nil {
			return records
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			records = append(records, Patch{Time: line})
		}
	}
	return records
}

func (a *App) getWindowsPatchRecords() []Patch {
	if runtime.GOOS != "windows" {
		return nil
	}

	ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil
	}
	defer unknown.Release()

	locator, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil
	}
	defer locator.Release()

	// 连接到本地 WMI
	service, err := oleutil.CallMethod(locator, "ConnectServer", ".", "root\\cimv2")
	if err != nil {
		return nil
	}
	defer service.ToIDispatch().Release()

	// 查询 Windows Update 历史记录
	query := "SELECT * FROM Win32_QuickFixEngineering ORDER BY InstalledOn DESC"
	updates, err := oleutil.CallMethod(service.ToIDispatch(), "ExecQuery", query)
	if err != nil {
		return nil
	}
	defer updates.ToIDispatch().Release()

	// 获取更新数量
	count, err := oleutil.GetProperty(updates.ToIDispatch(), "Count")
	if err != nil {
		return nil
	}
	defer count.Clear()

	var records []Patch
	for i := 0; i < int(count.Val); i++ {
		update, err := oleutil.CallMethod(updates.ToIDispatch(), "ItemIndex", i)
		if err != nil {
			continue
		}
		defer update.ToIDispatch().Release()

		// 获取更新详细信息
		installedOn, _ := oleutil.GetProperty(update.ToIDispatch(), "InstalledOn")
		hotFixID, _ := oleutil.GetProperty(update.ToIDispatch(), "HotFixID")
		description, _ := oleutil.GetProperty(update.ToIDispatch(), "Description")
		installedBy, _ := oleutil.GetProperty(update.ToIDispatch(), "InstalledBy")

		// 提取 KB 编号
		kb := ""
		if hotFixID != nil {
			kb = strings.TrimPrefix(hotFixID.ToString(), "KB")
		}

		record := Patch{
			Time:        installedOn.ToString(),
			Title:       description.ToString(),
			Description: fmt.Sprintf("Installed by: %s", installedBy.ToString()),
			Status:      "Installed",
			KB:          kb,
		}

		records = append(records, record)

		installedOn.Clear()
		hotFixID.Clear()
		description.Clear()
		installedBy.Clear()
	}

	return records
}
