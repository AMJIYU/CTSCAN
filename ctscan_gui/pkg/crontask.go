package pkg

import (
	"os/exec"
	"runtime"
	"strings"
)

type CronTask struct {
	Line string `json:"line"`
}

func (a *App) GetCronTasks() []CronTask {
	switch runtime.GOOS {
	case "windows":
		return a.getWindowsTasks()
	case "linux", "darwin":
		return a.getUnixTasks()
	default:
		return nil
	}
}

func (a *App) getWindowsTasks() []CronTask {
	// 使用 schtasks 命令获取 Windows 任务
	cmd := exec.Command("schtasks", "/query", "/fo", "list", "/v")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	var tasks []CronTask
	lines := strings.Split(string(out), "\n")
	var currentTask strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if currentTask.Len() > 0 {
				tasks = append(tasks, CronTask{Line: currentTask.String()})
				currentTask.Reset()
			}
			continue
		}
		if currentTask.Len() > 0 {
			currentTask.WriteString("\n")
		}
		currentTask.WriteString(line)
	}

	// 添加最后一个任务
	if currentTask.Len() > 0 {
		tasks = append(tasks, CronTask{Line: currentTask.String()})
	}

	return tasks
}

func (a *App) getUnixTasks() []CronTask {
	out, err := exec.Command("crontab", "-l").Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\n")
	var tasks []CronTask
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		tasks = append(tasks, CronTask{Line: line})
	}
	return tasks
}
