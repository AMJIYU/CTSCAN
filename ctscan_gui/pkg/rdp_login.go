package pkg

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// RDPLoginInfo 表示RDP登录信息
type RDPLoginInfo struct {
	Time        string `json:"time"`        // 登录时间
	Username    string `json:"username"`    // 用户名
	IP          string `json:"ip"`          // 登录IP
	Status      string `json:"status"`      // 登录状态
	Description string `json:"description"` // 描述信息
}

// 获取RDP登录日志
func (a *App) GetRDPLoginLogs() []RDPLoginInfo {
	switch runtime.GOOS {
	case "windows":
		return getWindowsRDPLogs()
	case "linux", "darwin":
		return getUnixRDPLogs()
	default:
		return []RDPLoginInfo{}
	}
}

// 获取Windows系统的RDP日志
func getWindowsRDPLogs() []RDPLoginInfo {
	records := make([]windowsRDPRecord, 0, 128)
	records = append(records, queryWindowsSecurityRDPLogs()...)
	records = append(records, queryWindowsTerminalServicesRDPLogs()...)
	return normalizeWindowsRDPRecords(records)
}

type windowsRDPRecord struct {
	RDPLoginInfo
	sortTime time.Time
}

type windowsRDPEventXML struct {
	System struct {
		Provider struct {
			Name string `xml:"Name,attr"`
		} `xml:"Provider"`
		EventID     string `xml:"EventID"`
		Channel     string `xml:"Channel"`
		Computer    string `xml:"Computer"`
		TimeCreated struct {
			SystemTime string `xml:"SystemTime,attr"`
		} `xml:"TimeCreated"`
	} `xml:"System"`
	EventData struct {
		Data []windowsRDPEventDataXML `xml:"Data"`
	} `xml:"EventData"`
	UserData struct {
		InnerXML string `xml:",innerxml"`
	} `xml:"UserData"`
	RenderingInfo struct {
		Message string `xml:"Message"`
	} `xml:"RenderingInfo"`
}

type windowsRDPEventDataXML struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:",chardata"`
}

func queryWindowsSecurityRDPLogs() []windowsRDPRecord {
	query := "*[System[(EventID=4624 or EventID=4625)] and EventData[(Data[@Name='LogonType']='10' or Data[@Name='LogonType']='7' or Data[@Name='LogonType']='3')]]"
	events := queryWindowsEventLogXML("Security", query, 500)
	records := make([]windowsRDPRecord, 0, len(events))

	for _, event := range events {
		eventID := event.eventID()
		values := event.values()
		logonType := firstWindowsEventValue(values, "LogonType")
		if eventID == 4624 && logonType != "10" && logonType != "7" {
			continue
		}
		if eventID == 4625 && logonType != "10" && logonType != "7" && logonType != "3" {
			continue
		}

		domain := firstWindowsEventValue(values, "TargetDomainName", "AccountDomain")
		username := normalizeWindowsRDPUser(domain, firstWindowsEventValue(values, "TargetUserName", "AccountName", "User"))
		ip := normalizeWindowsRDPIP(firstWindowsEventValue(values, "IpAddress", "SourceNetworkAddress", "ClientAddress", "Address"))
		if ip == "" {
			ip = normalizeWindowsRDPIP(firstWindowsEventValue(values, "WorkstationName"))
		}
		if username == "" && ip == "" {
			continue
		}

		status := "成功"
		if eventID == 4625 {
			status = "失败"
		}

		localTime, sortTime := event.localTime()
		description := fmt.Sprintf("安全日志 %d，%s", eventID, getLogonTypeDescription(logonType))
		if eventID == 4625 {
			if reason := firstWindowsEventValue(values, "FailureReason", "Status", "SubStatus"); reason != "" {
				description += "，原因：" + reason
			}
		}
		if logonType == "3" {
			description += "；网络级认证失败可能出现在 RDP/NLA 场景，也可能来自其他远程网络访问"
		}

		records = append(records, windowsRDPRecord{
			RDPLoginInfo: RDPLoginInfo{
				Time:        localTime,
				Username:    username,
				IP:          ip,
				Status:      status,
				Description: description,
			},
			sortTime: sortTime,
		})
	}

	return records
}

func queryWindowsTerminalServicesRDPLogs() []windowsRDPRecord {
	queries := []struct {
		channel string
		query   string
	}{
		{
			channel: "Microsoft-Windows-TerminalServices-RemoteConnectionManager/Operational",
			query:   "*[System[(EventID=1149)]]",
		},
		{
			channel: "Microsoft-Windows-TerminalServices-LocalSessionManager/Operational",
			query:   "*[System[(EventID=21 or EventID=23 or EventID=24 or EventID=25)]]",
		},
	}

	records := make([]windowsRDPRecord, 0, 128)
	for _, query := range queries {
		events := queryWindowsEventLogXML(query.channel, query.query, 500)
		for _, event := range events {
			eventID := event.eventID()
			values := event.values()
			status, action := terminalServicesRDPStatus(eventID)
			if status == "" {
				continue
			}

			domain := firstWindowsEventValue(values, "Domain", "Param2")
			username := normalizeWindowsRDPUser(domain, firstWindowsEventValue(values, "User", "UserName", "Param1", "AccountName"))
			ip := normalizeWindowsRDPIP(firstWindowsEventValue(values, "Address", "ClientAddress", "Param3", "SourceNetworkAddress", "IP"))
			if username == "" && ip == "" {
				continue
			}

			localTime, sortTime := event.localTime()
			records = append(records, windowsRDPRecord{
				RDPLoginInfo: RDPLoginInfo{
					Time:        localTime,
					Username:    username,
					IP:          ip,
					Status:      status,
					Description: fmt.Sprintf("TerminalServices %d，%s", eventID, action),
				},
				sortTime: sortTime,
			})
		}
	}

	return records
}

func terminalServicesRDPStatus(eventID int) (string, string) {
	switch eventID {
	case 1149:
		return "成功", "远程桌面用户认证成功"
	case 21:
		return "成功", "RDP 会话登录成功"
	case 25:
		return "重连", "RDP 会话重新连接"
	case 24:
		return "断开", "RDP 会话断开"
	case 23:
		return "登出", "RDP 会话注销"
	default:
		return "", ""
	}
}

func queryWindowsEventLogXML(channel, query string, count int) []windowsRDPEventXML {
	args := []string{"qe", channel, "/f:xml", "/rd:true", fmt.Sprintf("/c:%d", count)}
	if strings.TrimSpace(query) != "" {
		args = append(args, "/q:"+query)
	}

	output, err := backgroundCommand("wevtutil.exe", args...).Output()
	if err != nil || len(output) == 0 {
		return nil
	}

	events, err := parseWindowsRDPEventsXML(output)
	if err != nil {
		return nil
	}
	return events
}

func parseWindowsRDPEventsXML(output []byte) ([]windowsRDPEventXML, error) {
	xmlText := strings.TrimSpace(strings.TrimPrefix(string(output), "\ufeff"))
	xmlText = regexp.MustCompile(`(?s)<\?xml[^>]*\?>`).ReplaceAllString(xmlText, "")
	if xmlText == "" {
		return nil, nil
	}

	type eventsWrapper struct {
		Events []windowsRDPEventXML `xml:"Event"`
	}

	var wrapper eventsWrapper
	payload := xmlText
	if !strings.HasPrefix(strings.TrimSpace(xmlText), "<Events") {
		payload = "<Events>" + xmlText + "</Events>"
	}
	if err := xml.Unmarshal([]byte(payload), &wrapper); err != nil {
		var event windowsRDPEventXML
		if singleErr := xml.Unmarshal([]byte(xmlText), &event); singleErr == nil && event.eventID() != 0 {
			return []windowsRDPEventXML{event}, nil
		}
		return nil, err
	}

	return wrapper.Events, nil
}

func (event windowsRDPEventXML) eventID() int {
	value, _ := strconv.Atoi(strings.TrimSpace(event.System.EventID))
	return value
}

func (event windowsRDPEventXML) localTime() (string, time.Time) {
	raw := strings.TrimSpace(event.System.TimeCreated.SystemTime)
	if raw == "" {
		return "", time.Time{}
	}
	timestamp, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return raw, time.Time{}
	}
	local := timestamp.Local()
	return local.Format("2006-01-02 15:04:05"), local
}

func (event windowsRDPEventXML) values() map[string]string {
	values := make(map[string]string)
	for index, data := range event.EventData.Data {
		key := strings.TrimSpace(data.Name)
		if key == "" {
			key = fmt.Sprintf("Data%d", index+1)
		}
		addWindowsEventValue(values, key, data.Value)
	}
	for key, value := range extractWindowsUserDataValues(event.UserData.InnerXML) {
		addWindowsEventValue(values, key, value)
	}
	return values
}

func extractWindowsUserDataValues(innerXML string) map[string]string {
	values := make(map[string]string)
	if strings.TrimSpace(innerXML) == "" {
		return values
	}

	decoder := xml.NewDecoder(strings.NewReader("<Root>" + innerXML + "</Root>"))
	stack := make([]string, 0, 8)
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch item := token.(type) {
		case xml.StartElement:
			stack = append(stack, item.Name.Local)
		case xml.CharData:
			if len(stack) == 0 {
				continue
			}
			value := strings.TrimSpace(string(item))
			key := stack[len(stack)-1]
			if value != "" && key != "Root" {
				addWindowsEventValue(values, key, value)
			}
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}
	return values
}

func addWindowsEventValue(values map[string]string, key, value string) {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if key == "" || value == "" {
		return
	}
	if existing := strings.TrimSpace(values[key]); existing != "" {
		if !strings.Contains(existing, value) {
			values[key] = existing + ", " + value
		}
		return
	}
	values[key] = value
}

func firstWindowsEventValue(values map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(values[key]); value != "" && value != "-" {
			return value
		}
	}
	for _, key := range keys {
		for actualKey, value := range values {
			if strings.EqualFold(actualKey, key) && strings.TrimSpace(value) != "" && strings.TrimSpace(value) != "-" {
				return strings.TrimSpace(value)
			}
		}
	}
	return ""
}

func normalizeWindowsRDPUser(domain, username string) string {
	username = strings.TrimSpace(username)
	domain = strings.TrimSpace(domain)
	if username == "" || username == "-" {
		return ""
	}
	if strings.Contains(username, "\\") || domain == "" || domain == "-" || domain == "." {
		return username
	}
	return domain + "\\" + username
}

func normalizeWindowsRDPIP(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || value == "-" {
		return ""
	}
	return value
}

func normalizeWindowsRDPRecords(records []windowsRDPRecord) []RDPLoginInfo {
	sort.SliceStable(records, func(i, j int) bool {
		return records[i].sortTime.After(records[j].sortTime)
	})

	logs := make([]RDPLoginInfo, 0, len(records))
	seen := make(map[string]struct{}, len(records))
	for _, record := range records {
		key := strings.ToLower(strings.Join([]string{
			record.Time,
			record.Username,
			record.IP,
			record.Status,
			record.Description,
		}, "|"))
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		logs = append(logs, record.RDPLoginInfo)
	}
	return logs
}

// 获取Linux/macOS系统的RDP日志
func getUnixRDPLogs() []RDPLoginInfo {
	var logs []RDPLoginInfo

	// 根据操作系统选择对应的日志文件
	var logFile string
	switch runtime.GOOS {
	case "darwin":
		// 在 macOS 上，检查多个可能的日志文件位置
		possibleLogFiles := []string{
			"/var/log/system.log",
			"/var/log/asl.log",
			"/Library/Logs/DiagnosticReports/RemoteDesktop*.log",
		}
		for _, file := range possibleLogFiles {
			if _, err := os.Stat(file); err == nil {
				logFile = file
				break
			}
		}
	case "linux":
		// 检查常见的日志文件
		possibleLogFiles := []string{
			"/var/log/auth.log", // Debian/Ubuntu
			"/var/log/secure",   // RHEL/CentOS
		}
		for _, file := range possibleLogFiles {
			if _, err := os.Stat(file); err == nil {
				logFile = file
				break
			}
		}
	}

	// 如果没有找到日志文件，直接返回
	if logFile == "" {
		fmt.Printf("未找到RDP日志文件\n")
		return logs
	}

	fmt.Printf("正在检查日志文件: %s\n", logFile)

	// 使用 grep 命令直接过滤出包含 RDP 相关内容的行
	cmd := backgroundCommand("grep", "-i", "-E", "(xrdp|rdp|RemoteDesktop)", logFile)
	output, err := cmd.Output()
	if err != nil {
		// grep 在没有找到匹配项时会返回错误，这是正常的
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			fmt.Printf("未找到RDP相关日志记录\n")
			return logs
		}
		fmt.Printf("执行grep命令失败: %v\n", err)
		return logs
	}

	// 匹配RDP相关的日志行
	lines := strings.Split(string(output), "\n")
	fmt.Printf("找到 %d 行相关日志\n", len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		// 提取时间
		timeRegex := regexp.MustCompile(`(\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})`)
		timeMatch := timeRegex.FindStringSubmatch(line)
		if len(timeMatch) < 2 {
			continue
		}

		// 提取用户名
		userRegex := regexp.MustCompile(`user\s+(\w+)`)
		userMatch := userRegex.FindStringSubmatch(line)
		username := "unknown"
		if len(userMatch) >= 2 {
			username = userMatch[1]
		}

		// 提取IP地址
		ipRegex := regexp.MustCompile(`from\s+([\d\.]+)`)
		ipMatch := ipRegex.FindStringSubmatch(line)
		ip := "unknown"
		if len(ipMatch) >= 2 {
			ip = ipMatch[1]
		}

		// 判断登录状态
		status := "失败"
		if strings.Contains(line, "successful") || strings.Contains(line, "Accepted") || strings.Contains(line, "connected") {
			status = "成功"
		}

		logs = append(logs, RDPLoginInfo{
			Time:        timeMatch[1],
			Username:    username,
			IP:          ip,
			Status:      status,
			Description: fmt.Sprintf("用户 %s 从 %s 尝试登录", username, ip),
		})
	}

	fmt.Printf("成功解析 %d 条RDP登录记录\n", len(logs))
	return logs
}
