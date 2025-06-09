package main

import (
	"os/exec"
	"strings"
)

type CronTask struct {
	Line string `json:"line"`
}

func (a *App) GetCronTasks() []CronTask {
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
