package pkg

import (
	"os"
	"runtime"
)

type SystemInfo struct {
	Hostname  string `json:"hostname"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	CPUCores  int    `json:"cpu_cores"`
	GoVersion string `json:"go_version"`
}

func (a *App) GetSystemInfo() SystemInfo {
	host, _ := os.Hostname()
	return SystemInfo{
		Hostname:  host,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		CPUCores:  runtime.NumCPU(),
		GoVersion: runtime.Version(),
	}
}
