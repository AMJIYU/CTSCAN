package pkg

import (
	"os/exec"
	"runtime"
	"strings"
)

type LoginFailed struct {
	Time          string `json:"time"`
	EventID       string `json:"event_id"`
	LoginType     string `json:"login_type"`
	SourceIP      string `json:"source_ip"`
	Username      string `json:"username"`
	Workstation   string `json:"workstation"`
	SubjectUser   string `json:"subject_user"`
	SubjectDomain string `json:"subject_domain"`
	Process       string `json:"process"`
	FailureReason string `json:"failure_reason"`
}

func (a *App) GetLoginFailedRecords() []LoginFailed {
	records := []LoginFailed{}
	switch runtime.GOOS {
	case "windows":
		// PowerShell命令获取4625事件（登录失败）
		cmd := `Get-WinEvent -LogName Security | Where-Object { $_.Id -eq 4625 } | Select-Object TimeCreated,Id,Properties | Format-List`
		out, err := exec.Command("powershell", "-Command", cmd).CombinedOutput()
		if err != nil {
			return records
		}
		lines := strings.Split(string(out), "\n")
		var currentRecord LoginFailed
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				if currentRecord.Time != "" {
					records = append(records, currentRecord)
					currentRecord = LoginFailed{}
				}
				continue
			}
			if strings.HasPrefix(line, "TimeCreated") {
				currentRecord.Time = strings.TrimPrefix(line, "TimeCreated : ")
			} else if strings.HasPrefix(line, "Id") {
				currentRecord.EventID = strings.TrimPrefix(line, "Id : ")
			}
		}
		if currentRecord.Time != "" {
			records = append(records, currentRecord)
		}

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
			parts := strings.Fields(line)
			if len(parts) >= 15 {
				record := LoginFailed{
					Time:          strings.Join(parts[0:3], " "),
					Username:      parts[8],
					SourceIP:      parts[10],
					FailureReason: "Invalid password",
				}
				records = append(records, record)
			}
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
			record := LoginFailed{
				Time:          line[:15],
				FailureReason: "Authentication failed",
			}
			records = append(records, record)
		}
	}
	return records
}
