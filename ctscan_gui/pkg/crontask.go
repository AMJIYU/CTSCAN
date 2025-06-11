package pkg

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
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
	if runtime.GOOS != "windows" {
		return nil
	}

	ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("Schedule.Service")
	if err != nil {
		return nil
	}
	defer unknown.Release()

	service, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil
	}
	defer service.Release()

	// 连接到本地任务计划程序
	_, err = oleutil.CallMethod(service, "Connect")
	if err != nil {
		return nil
	}

	// 获取根文件夹
	rootFolder, err := oleutil.CallMethod(service, "GetFolder", "\\")
	if err != nil {
		return nil
	}
	defer rootFolder.ToIDispatch().Release()

	// 获取所有任务
	tasks, err := oleutil.CallMethod(rootFolder.ToIDispatch(), "GetTasks", 0)
	if err != nil {
		return nil
	}
	defer tasks.ToIDispatch().Release()

	// 获取任务数量
	count, err := oleutil.GetProperty(tasks.ToIDispatch(), "Count")
	if err != nil {
		return nil
	}
	defer count.Clear()

	var taskList []CronTask
	for i := 0; i < int(count.Val); i++ {
		task, err := oleutil.CallMethod(tasks.ToIDispatch(), "Item", i+1)
		if err != nil {
			continue
		}
		defer task.ToIDispatch().Release()

		// 获取任务信息
		name, _ := oleutil.GetProperty(task.ToIDispatch(), "Name")
		state, _ := oleutil.GetProperty(task.ToIDispatch(), "State")
		lastRunTime, _ := oleutil.GetProperty(task.ToIDispatch(), "LastRunTime")
		nextRunTime, _ := oleutil.GetProperty(task.ToIDispatch(), "NextRunTime")

		taskInfo := fmt.Sprintf("TaskName: %s; State: %s; LastRunTime: %s; NextRunTime: %s",
			name.ToString(),
			state.ToString(),
			lastRunTime.ToString(),
			nextRunTime.ToString())

		taskList = append(taskList, CronTask{Line: taskInfo})

		name.Clear()
		state.Clear()
		lastRunTime.Clear()
		nextRunTime.Clear()
	}

	return taskList
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
