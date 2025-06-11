package pkg

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type ShellHistory struct {
	Time    string `json:"time"`
	Command string `json:"command"`
	User    string `json:"user"`
	Shell   string `json:"shell"`
}

func (a *App) GetShellHistory() []ShellHistory {
	var records []ShellHistory

	switch runtime.GOOS {
	case "darwin":
		// 获取当前用户
		user := os.Getenv("USER")

		// 获取 zsh 历史记录
		zshHistory := filepath.Join(os.Getenv("HOME"), ".zsh_history")
		if content, err := os.ReadFile(zshHistory); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				// zsh 历史记录格式：: 时间戳:0;命令
				if strings.HasPrefix(line, ": ") {
					parts := strings.SplitN(line[2:], ":", 2)
					if len(parts) == 2 {
						timestamp, err := time.Parse("2006-01-02 15:04:05", time.Unix(0, 0).Format("2006-01-02 15:04:05"))
						if err == nil {
							records = append(records, ShellHistory{
								Time:    timestamp.Format("2006-01-02 15:04:05"),
								Command: strings.TrimSpace(parts[1]),
								User:    user,
								Shell:   "zsh",
							})
						}
					}
				}
			}
		}

		// 获取 bash 历史记录
		bashHistory := filepath.Join(os.Getenv("HOME"), ".bash_history")
		if content, err := os.ReadFile(bashHistory); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				records = append(records, ShellHistory{
					Time:    time.Now().Format("2006-01-02 15:04:05"), // bash 历史记录没有时间戳
					Command: strings.TrimSpace(line),
					User:    user,
					Shell:   "bash",
				})
			}
		}

		// 获取 fish 历史记录
		fishHistory := filepath.Join(os.Getenv("HOME"), ".local/share/fish/fish_history")
		if content, err := os.ReadFile(fishHistory); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				// fish 历史记录格式：- cmd: 命令
				if strings.HasPrefix(line, "- cmd: ") {
					records = append(records, ShellHistory{
						Time:    time.Now().Format("2006-01-02 15:04:05"), // fish 历史记录没有时间戳
						Command: strings.TrimSpace(strings.TrimPrefix(line, "- cmd: ")),
						User:    user,
						Shell:   "fish",
					})
				}
			}
		}

	case "linux":
		// 获取当前用户
		user := os.Getenv("USER")

		// 获取 bash 历史记录
		bashHistory := filepath.Join(os.Getenv("HOME"), ".bash_history")
		if content, err := os.ReadFile(bashHistory); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				records = append(records, ShellHistory{
					Time:    time.Now().Format("2006-01-02 15:04:05"),
					Command: strings.TrimSpace(line),
					User:    user,
					Shell:   "bash",
				})
			}
		}

		// 获取 zsh 历史记录
		zshHistory := filepath.Join(os.Getenv("HOME"), ".zsh_history")
		if content, err := os.ReadFile(zshHistory); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				if strings.HasPrefix(line, ": ") {
					parts := strings.SplitN(line[2:], ":", 2)
					if len(parts) == 2 {
						timestamp, err := time.Parse("2006-01-02 15:04:05", time.Unix(0, 0).Format("2006-01-02 15:04:05"))
						if err == nil {
							records = append(records, ShellHistory{
								Time:    timestamp.Format("2006-01-02 15:04:05"),
								Command: strings.TrimSpace(parts[1]),
								User:    user,
								Shell:   "zsh",
							})
						}
					}
				}
			}
		}

	case "windows":
		// 获取 PowerShell 历史记录
		psHistory := filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\PowerShell\\PSReadline\\ConsoleHost_history.txt")
		if content, err := os.ReadFile(psHistory); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				records = append(records, ShellHistory{
					Time:    time.Now().Format("2006-01-02 15:04:05"),
					Command: strings.TrimSpace(line),
					User:    os.Getenv("USERNAME"),
					Shell:   "powershell",
				})
			}
		}

		// 获取 CMD 历史记录
		cmdHistory := filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\PowerShell\\PSReadline\\ConsoleHost_history.txt")
		if content, err := os.ReadFile(cmdHistory); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				records = append(records, ShellHistory{
					Time:    time.Now().Format("2006-01-02 15:04:05"),
					Command: strings.TrimSpace(line),
					User:    os.Getenv("USERNAME"),
					Shell:   "cmd",
				})
			}
		}
	}

	return records
}
