package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	PPID       int32   `json:"ppid"`
	ParentName string  `json:"parent_name"`
	CreateTime int64   `json:"create_time"`
	Exe        string  `json:"exe"`
	FileCtime  int64   `json:"file_ctime"`
	FileMtime  int64   `json:"file_mtime"`
	MD5        string  `json:"md5"`
	Signature  string  `json:"signature"`
	CPUPercent float64 `json:"cpu_percent"`
	MemPercent float64 `json:"mem_percent"`
}

func getFileMD5(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func getFileTimes(path string) (ctime, mtime int64) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, 0
	}
	mtime = fi.ModTime().Unix()
	return 0, mtime // macOS/Linux下ctime无法直接获取，返回0
}

func getSignature(path string) string {
	out, err := exec.Command("codesign", "-dv", "--verbose=4", path).CombinedOutput()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Authority=") {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

func (a *App) GetAllProcesses() []ProcInfo {
	procs, err := process.Processes()
	if err != nil {
		return nil
	}
	var result []ProcInfo
	for _, p := range procs {
		name, _ := p.Name()
		exe, _ := p.Exe()
		ctime, _ := p.CreateTime()
		ppid, _ := p.Ppid()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		parentName := ""
		if ppid > 0 {
			if parent, err := process.NewProcess(ppid); err == nil {
				parentName, _ = parent.Name()
			}
		}
		fileCtime, fileMtime := getFileTimes(exe)
		md5sum := getFileMD5(exe)
		signature := getSignature(exe)
		result = append(result, ProcInfo{
			PID:        p.Pid,
			Name:       name,
			PPID:       ppid,
			ParentName: parentName,
			CreateTime: ctime,
			Exe:        exe,
			FileCtime:  fileCtime,
			FileMtime:  fileMtime,
			MD5:        md5sum,
			Signature:  signature,
			CPUPercent: cpuPercent,
			MemPercent: float64(memPercent),
		})
	}
	return result
}
