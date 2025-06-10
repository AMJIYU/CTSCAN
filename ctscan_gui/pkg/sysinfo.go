package pkg

import (
	"os"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type DiskInfo struct {
	MountPoint string  `json:"mount_point"`
	TotalSize  uint64  `json:"total_size"`
	UsedSize   uint64  `json:"used_size"`
	FreeSize   uint64  `json:"free_size"`
	Usage      float64 `json:"usage"`
}

type SystemInfo struct {
	Hostname      string     `json:"hostname"`
	OS            string     `json:"os"`
	Arch          string     `json:"arch"`
	CPUCores      int        `json:"cpu_cores"`
	KernelVersion string     `json:"kernel_version"`
	CPUUsage      float64    `json:"cpu_usage"`
	TotalMemory   uint64     `json:"total_memory"`
	MemoryUsage   float64    `json:"memory_usage"`
	Disks         []DiskInfo `json:"disks"`
}

func (a *App) GetSystemInfo() SystemInfo {
	hostname, _ := os.Hostname()

	// 使用 gopsutil 获取操作系统信息
	kernelVersion := ""
	if hostInfo, err := host.Info(); err == nil {
		kernelVersion = hostInfo.KernelVersion
	}

	// 使用 gopsutil 获取 CPU 使用率
	cpuUsage := 0.0
	if percentages, err := cpu.Percent(0, false); err == nil && len(percentages) > 0 {
		cpuUsage = percentages[0]
	}

	// 使用 gopsutil 获取内存信息
	var totalMemory uint64
	var memoryUsage float64
	if vmStat, err := mem.VirtualMemory(); err == nil {
		totalMemory = vmStat.Total
		memoryUsage = vmStat.UsedPercent
	}

	// 使用 gopsutil 获取磁盘信息
	var disks []DiskInfo
	if partitions, err := disk.Partitions(false); err == nil {
		for _, partition := range partitions {
			// 获取所有挂载点的使用情况
			if usage, err := disk.Usage(partition.Mountpoint); err == nil {
				disks = append(disks, DiskInfo{
					MountPoint: partition.Mountpoint,
					TotalSize:  usage.Total,
					UsedSize:   usage.Used,
					FreeSize:   usage.Free,
					Usage:      usage.UsedPercent,
				})
			}
		}
	}

	return SystemInfo{
		Hostname:      hostname,
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		CPUCores:      runtime.NumCPU(),
		KernelVersion: kernelVersion,
		CPUUsage:      cpuUsage,
		TotalMemory:   totalMemory,
		MemoryUsage:   memoryUsage,
		Disks:         disks,
	}
}
