package pkg

import (
	"os/exec"
	"runtime"
	"strings"
)

type LoginSuccess struct {
	Time          string `json:"time"`
	EventID       string `json:"event_id"`
	LoginType     string `json:"login_type"`
	SourceIP      string `json:"source_ip"`
	Username      string `json:"username"`
	Workstation   string `json:"workstation"`
	SubjectUser   string `json:"subject_user"`
	SubjectDomain string `json:"subject_domain"`
	Process       string `json:"process"`
}

func (a *App) GetLoginSuccessRecords() []LoginSuccess {
	records := []LoginSuccess{}
	switch runtime.GOOS {
	case "windows":
		// PowerShell命令获取4624事件
		cmd := `Get-WinEvent -LogName Security | Where-Object { $_.Id -eq 4624 } | Select-Object TimeCreated,Id | Format-List`
		out, err := exec.Command("powershell", "-Command", cmd).CombinedOutput()
		if err != nil {
			return records
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			records = append(records, LoginSuccess{Time: line})
		}
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
		out, err := exec.Command("log", "show", "--predicate", "eventMessage CONTAINS 'login'", "--info", "--last", "1d").CombinedOutput()
		if err != nil {
			return records
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			records = append(records, LoginSuccess{Time: line})
		}
	}
	return records
}
