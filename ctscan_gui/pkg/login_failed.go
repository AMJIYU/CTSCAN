package pkg

import (
	"encoding/json"
	"os/exec"
	"runtime"
	"strings"
	"time"

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
	Reason    string `json:"reason"`
}

func (a *App) GetLoginFailedRecords() []LoginFailed {
	records := []LoginFailed{}
	switch runtime.GOOS {
	case "windows":
		records = a.getWindowsLoginFailedRecords()
	case "linux":
		out, err := exec.Command("grep", "Failed password", "/var/log/auth.log").CombinedOutput()
		if err != nil {
			println("Linux 获取登录失败记录错误:", err.Error())
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
		println("开始获取 macOS 登录失败记录...")
		// 使用多个条件组合查询，确保能捕获所有登录失败情况
		predicates := []string{
			"subsystem == 'com.apple.security'",
			"subsystem == 'com.apple.authentication'",
			"eventMessage CONTAINS 'Failed to authenticate'",
			"eventMessage CONTAINS 'authentication failed'",
			"eventMessage CONTAINS 'password check failed'",
			"eventMessage CONTAINS 'login failed'",
			"eventMessage CONTAINS 'authentication error'",
		}

		var allRecords []LoginFailed
		for _, predicate := range predicates {
			println("执行查询:", predicate)
			// 增加时间范围到30天，并添加更多日志级别
			out, err := exec.Command("log", "show", "--predicate", predicate, "--info", "--debug", "--last", "30d", "--style", "json").CombinedOutput()
			if err != nil {
				println("查询失败:", err.Error())
				continue
			}

			// 解析 JSON 输出
			var logEntries []struct {
				Timestamp    string `json:"timestamp"`
				EventMessage string `json:"eventMessage"`
			}

			if err := json.Unmarshal(out, &logEntries); err != nil {
				println("JSON 解析失败:", err.Error())
				continue
			}

			println("找到", len(logEntries), "条记录")

			for _, entry := range logEntries {
				// 解析时间戳
				timestamp, err := time.Parse(time.RFC3339, entry.Timestamp)
				if err != nil {
					continue
				}

				// 解析用户名和失败原因
				message := entry.EventMessage
				username := ""
				reason := ""

				// 提取用户名
				if strings.Contains(message, "for user") {
					parts := strings.Split(message, "for user")
					if len(parts) > 1 {
						username = strings.TrimSpace(parts[1])
					}
				} else if strings.Contains(message, "user:") {
					parts := strings.Split(message, "user:")
					if len(parts) > 1 {
						username = strings.TrimSpace(parts[1])
					}
				}

				// 提取失败原因
				if strings.Contains(message, "reason:") {
					parts := strings.Split(message, "reason:")
					if len(parts) > 1 {
						reason = strings.TrimSpace(parts[1])
					}
				} else {
					// 如果没有明确的原因，使用消息本身作为原因
					reason = message
				}

				record := LoginFailed{
					Time:      timestamp.Format("2006-01-02 15:04:05"),
					EventID:   "4625",
					EventType: "登录失败",
					Source:    "System",
					Username:  username,
					IPAddress: "本地",
					Reason:    reason,
				}

				allRecords = append(allRecords, record)
			}
		}

		// 去重
		seen := make(map[string]bool)
		for _, record := range allRecords {
			key := record.Time + record.Username + record.Reason
			if !seen[key] {
				seen[key] = true
				records = append(records, record)
			}
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
