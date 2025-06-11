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

// 获取当前 shell
func getCurrentShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return ""
	}
	// 从完整路径中提取 shell 名称
	parts := strings.Split(shell, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func (a *App) GetShellHistory() []ShellHistory {
	var records []ShellHistory
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME") // Windows 环境
	}

	// 获取当前 shell
	currentShell := getCurrentShell()

	// 用于去重的 map
	seen := make(map[string]bool)

	// 用于存储不同 shell 的记录
	var currentShellRecords []ShellHistory
	var otherShellRecords []ShellHistory

	switch runtime.GOOS {
	case "darwin", "linux":
		// 获取 zsh 历史记录
		zshHistory := filepath.Join(os.Getenv("HOME"), ".zsh_history")
		if content, err := os.ReadFile(zshHistory); err == nil {
			lines := strings.Split(string(content), "\n")
			// 倒序遍历，这样最新的记录会先被处理
			for i := len(lines) - 1; i >= 0; i-- {
				line := lines[i]
				if line == "" {
					continue
				}
				// zsh 历史记录格式：: 时间戳:0;命令
				if strings.HasPrefix(line, ": ") {
					parts := strings.SplitN(line[2:], ":", 2)
					if len(parts) == 2 {
						cmd := strings.TrimSpace(parts[1])
						// 去重
						key := cmd + "|" + user + "|zsh"
						if !seen[key] {
							seen[key] = true
							record := ShellHistory{
								Time:    time.Now().Format("2006-01-02 15:04:05"),
								Command: cmd,
								User:    user,
								Shell:   "zsh",
							}
							if currentShell == "zsh" {
								currentShellRecords = append(currentShellRecords, record)
							} else {
								otherShellRecords = append(otherShellRecords, record)
							}
						}
					}
				}
			}
		}

		// 获取 bash 历史记录
		bashHistory := filepath.Join(os.Getenv("HOME"), ".bash_history")
		if content, err := os.ReadFile(bashHistory); err == nil {
			lines := strings.Split(string(content), "\n")
			// 倒序遍历
			for i := len(lines) - 1; i >= 0; i-- {
				line := lines[i]
				if line == "" {
					continue
				}
				cmd := strings.TrimSpace(line)
				// 去重
				key := cmd + "|" + user + "|bash"
				if !seen[key] {
					seen[key] = true
					record := ShellHistory{
						Time:    time.Now().Format("2006-01-02 15:04:05"),
						Command: cmd,
						User:    user,
						Shell:   "bash",
					}
					if currentShell == "bash" {
						currentShellRecords = append(currentShellRecords, record)
					} else {
						otherShellRecords = append(otherShellRecords, record)
					}
				}
			}
		}

		// 获取 fish 历史记录
		fishHistory := filepath.Join(os.Getenv("HOME"), ".local/share/fish/fish_history")
		if content, err := os.ReadFile(fishHistory); err == nil {
			lines := strings.Split(string(content), "\n")
			// 倒序遍历
			for i := len(lines) - 1; i >= 0; i-- {
				line := lines[i]
				if line == "" {
					continue
				}
				// fish 历史记录格式：- cmd: 命令
				if strings.HasPrefix(line, "- cmd: ") {
					cmd := strings.TrimSpace(strings.TrimPrefix(line, "- cmd: "))
					// 去重
					key := cmd + "|" + user + "|fish"
					if !seen[key] {
						seen[key] = true
						record := ShellHistory{
							Time:    time.Now().Format("2006-01-02 15:04:05"),
							Command: cmd,
							User:    user,
							Shell:   "fish",
						}
						if currentShell == "fish" {
							currentShellRecords = append(currentShellRecords, record)
						} else {
							otherShellRecords = append(otherShellRecords, record)
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
			// 倒序遍历
			for i := len(lines) - 1; i >= 0; i-- {
				line := lines[i]
				if line == "" {
					continue
				}
				cmd := strings.TrimSpace(line)
				// 去重
				key := cmd + "|" + user + "|powershell"
				if !seen[key] {
					seen[key] = true
					record := ShellHistory{
						Time:    time.Now().Format("2006-01-02 15:04:05"),
						Command: cmd,
						User:    user,
						Shell:   "powershell",
					}
					if currentShell == "powershell" {
						currentShellRecords = append(currentShellRecords, record)
					} else {
						otherShellRecords = append(otherShellRecords, record)
					}
				}
			}
		}

		// 获取 CMD 历史记录
		cmdHistory := filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\PowerShell\\PSReadline\\ConsoleHost_history.txt")
		if content, err := os.ReadFile(cmdHistory); err == nil {
			lines := strings.Split(string(content), "\n")
			// 倒序遍历
			for i := len(lines) - 1; i >= 0; i-- {
				line := lines[i]
				if line == "" {
					continue
				}
				cmd := strings.TrimSpace(line)
				// 去重
				key := cmd + "|" + user + "|cmd"
				if !seen[key] {
					seen[key] = true
					record := ShellHistory{
						Time:    time.Now().Format("2006-01-02 15:04:05"),
						Command: cmd,
						User:    user,
						Shell:   "cmd",
					}
					if currentShell == "cmd" {
						currentShellRecords = append(currentShellRecords, record)
					} else {
						otherShellRecords = append(otherShellRecords, record)
					}
				}
			}
		}
	}

	// 合并记录，当前 shell 的记录在前
	records = append(currentShellRecords, otherShellRecords...)

	return records
}
