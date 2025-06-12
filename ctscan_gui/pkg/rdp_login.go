package pkg

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
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
	var logs []RDPLoginInfo

	// 使用wevtutil命令获取Windows事件日志
	cmd := exec.Command("wevtutil", "qe", "Microsoft-Windows-RemoteDesktopServices-RdpCoreTS/Operational", "/f:xml")
	output, err := cmd.Output()
	if err != nil {
		return logs
	}

	// 解析XML格式的事件日志
	type Event struct {
		EventData struct {
			Data []struct {
				Name  string `xml:"Name,attr"`
				Value string `xml:",chardata"`
			} `xml:"Data"`
		} `xml:"EventData"`
		System struct {
			TimeCreated struct {
				SystemTime string `xml:"SystemTime,attr"`
			} `xml:"TimeCreated"`
		} `xml:"System"`
	}

	var events struct {
		Events []Event `xml:"Event"`
	}

	if err := xml.Unmarshal(output, &events); err != nil {
		return logs
	}

	for _, event := range events.Events {
		var username, ip, status string
		for _, data := range event.EventData.Data {
			switch data.Name {
			case "User":
				username = data.Value
			case "Address":
				ip = data.Value
			case "Status":
				status = data.Value
			}
		}

		// 解析时间
		t, err := time.Parse(time.RFC3339, event.System.TimeCreated.SystemTime)
		if err != nil {
			continue
		}

		logs = append(logs, RDPLoginInfo{
			Time:        t.Format("2006-01-02 15:04:05"),
			Username:    username,
			IP:          ip,
			Status:      status,
			Description: fmt.Sprintf("用户 %s 从 %s 尝试登录", username, ip),
		})
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
	cmd := exec.Command("grep", "-i", "-E", "(xrdp|rdp|RemoteDesktop)", logFile)
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
