package main

import "github.com/shirou/gopsutil/v4/process"

type ProcInfo struct {
	PID     int32  `json:"pid"`
	Name    string `json:"name"`
	Cmdline string `json:"cmdline"`
	User    string `json:"user"`
	Exe     string `json:"exe"`
}

func (a *App) GetAllProcesses() []ProcInfo {
	procs, err := process.Processes()
	if err != nil {
		return nil
	}
	var result []ProcInfo
	for _, p := range procs {
		name, _ := p.Name()
		cmd, _ := p.Cmdline()
		user, _ := p.Username()
		exe, _ := p.Exe()
		result = append(result, ProcInfo{
			PID:     p.Pid,
			Name:    name,
			Cmdline: cmd,
			User:    user,
			Exe:     exe,
		})
	}
	return result
}
