package pkg

import (
	"bufio"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type LoginSuccess struct {
	Time      string `json:"time"`
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	Source    string `json:"source"`
	Username  string `json:"username"`
	IPAddress string `json:"ip_address"`
}

func (a *App) GetLoginSuccessRecords() []LoginSuccess {
	records := []LoginSuccess{}
	switch runtime.GOOS {
	case "windows":
		records = a.getWindowsLoginSuccessRecords()
	case "linux":
		out, err := exec.Command("grep", "Accepted", "/var/log/auth.log").CombinedOutput()
		if err != nil {
			return records
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			records = append(records, LoginSuccess{Time: line[:15]})
		}
	case "darwin":
		cmd := exec.Command("last")
		output, err := cmd.Output()
		if err != nil {
			return records
		}

		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		// 匹配格式：eleven     ttys001                         四  5 29 10:36   still logged in
		// 或者：eleven     ttys014                         二  5 27 14:10 - 14:10  (00:00)
		re := regexp.MustCompile(`^(\S+)\s+(\S+)\s+(?:\S+\s+){1,2}(\d+)\s+(\d+):(\d+)(?:\s+-\s+\d+:\d+\s+\((\d+:\d+)\))?(?:\s+still logged in)?$`)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" || strings.Contains(line, "wtmp begins") {
				continue
			}

			matches := re.FindStringSubmatch(line)
			if len(matches) >= 5 {
				username := matches[1]
				terminal := matches[2] // 获取终端信息
				day := matches[3]
				hour := matches[4]
				minute := matches[5]

				// 构建时间字符串
				now := time.Now()
				year := now.Year()
				monthNum := now.Month()
				dayNum, _ := strconv.Atoi(day)
				hourNum, _ := strconv.Atoi(hour)
				minuteNum, _ := strconv.Atoi(minute)

				// 创建时间对象
				t := time.Date(year, monthNum, dayNum, hourNum, minuteNum, 0, 0, time.Local)

				// 检查是否是当前登录
				status := "已登出"
				if strings.Contains(line, "still logged in") {
					status = "当前登录"
				}

				// 获取终端对应的设备路径
				devicePath := "/dev/" + terminal

				records = append(records, LoginSuccess{
					Time:      t.Format(time.RFC3339),
					EventType: status,
					Source:    devicePath,
					Username:  username,
				})
			}
		}
	}
	return records
}

func (a *App) getWindowsLoginSuccessRecords() []LoginSuccess {
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

	// 查询安全日志中的登录成功事件
	query := "SELECT * FROM Win32_NTLogEvent WHERE LogFile='Security' AND EventCode=4624"
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

	var records []LoginSuccess
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

		record := LoginSuccess{
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
