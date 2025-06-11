package pkg

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type LoginFailed struct {
	Time      string `json:"time"`
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	Source    string `json:"source"`
	Username  string `json:"username"`
	IPAddress string `json:"ip_address"`
}

func (a *App) GetLoginFailedRecords() []LoginFailed {
	records := []LoginFailed{}
	switch runtime.GOOS {
	case "windows":
		records = a.getWindowsLoginFailedRecords()
	case "linux":
		out, err := exec.Command("grep", "Failed password", "/var/log/auth.log").CombinedOutput()
		if err != nil {
			return records
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			records = append(records, LoginFailed{Time: line[:15]})
		}
	case "darwin":
		out, err := exec.Command("log", "show", "--predicate", "eventMessage CONTAINS 'Failed login'", "--info", "--last", "1d").CombinedOutput()
		if err != nil {
			return records
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			records = append(records, LoginFailed{Time: line})
		}
	}
	return records
}

func (a *App) getWindowsLoginFailedRecords() []LoginFailed {
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

	// 查询安全日志中的登录失败事件
	query := "SELECT * FROM Win32_NTLogEvent WHERE LogFile='Security' AND EventCode=4625"
	events, err := oleutil.CallMethod(service.ToIDispatch(), "ExecQuery", query)
	if err != nil {
		return nil
	}
	defer events.ToIDispatch().Release()

	// 获取事件数量
	count, err := oleutil.GetProperty(events.ToIDispatch(), "Count")
	if err != nil {
		return nil
	}
	defer count.Clear()

	var records []LoginFailed
	for i := 0; i < int(count.Val); i++ {
		event, err := oleutil.CallMethod(events.ToIDispatch(), "ItemIndex", i)
		if err != nil {
			continue
		}
		defer event.ToIDispatch().Release()

		// 获取事件详细信息
		timeCreated, _ := oleutil.GetProperty(event.ToIDispatch(), "TimeGenerated")
		eventID, _ := oleutil.GetProperty(event.ToIDispatch(), "EventCode")
		eventType, _ := oleutil.GetProperty(event.ToIDispatch(), "Type")
		sourceName, _ := oleutil.GetProperty(event.ToIDispatch(), "SourceName")
		message, _ := oleutil.GetProperty(event.ToIDispatch(), "Message")

		// 解析消息中的用户名和IP地址
		username := ExtractUsername(message.ToString())
		ipAddress := ExtractIPAddress(message.ToString())

		record := LoginFailed{
			Time:      timeCreated.ToString(),
			EventID:   eventID.ToString(),
			EventType: eventType.ToString(),
			Source:    sourceName.ToString(),
			Username:  username,
			IPAddress: ipAddress,
		}

		records = append(records, record)

		timeCreated.Clear()
		eventID.Clear()
		eventType.Clear()
		sourceName.Clear()
		message.Clear()
	}

	return records
}
