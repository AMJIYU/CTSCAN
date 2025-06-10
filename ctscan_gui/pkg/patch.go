package pkg

import (
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type PatchInfo struct {
	Name   string `json:"name"`
	Date   string `json:"date"`
	Status string `json:"status"`
}

func (a *App) GetPatchInfo() []PatchInfo {
	switch runtime.GOOS {
	case "windows":
		return a.getWindowsPatches()
	case "linux":
		return a.getLinuxPatches()
	case "darwin":
		return a.getMacOSPatches()
	default:
		return nil
	}
}

func (a *App) getWindowsPatches() []PatchInfo {
	// PowerShell命令获取Windows Update历史
	cmd := `Get-HotFix | Select-Object HotFixID, InstalledOn, Description | Format-List`
	out, err := exec.Command("powershell", "-Command", cmd).CombinedOutput()
	if err != nil {
		return nil
	}

	var patches []PatchInfo
	lines := strings.Split(string(out), "\n")
	var currentPatch PatchInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if currentPatch.Name != "" {
				patches = append(patches, currentPatch)
				currentPatch = PatchInfo{}
			}
			continue
		}

		if strings.HasPrefix(line, "HotFixID") {
			currentPatch.Name = strings.TrimSpace(strings.TrimPrefix(line, "HotFixID :"))
		} else if strings.HasPrefix(line, "InstalledOn") {
			dateStr := strings.TrimSpace(strings.TrimPrefix(line, "InstalledOn :"))
			if t, err := time.Parse("1/2/2006 3:04:05 PM", dateStr); err == nil {
				currentPatch.Date = t.Format("2006-01-02 15:04:05")
			}
		} else if strings.HasPrefix(line, "Description") {
			currentPatch.Status = strings.TrimSpace(strings.TrimPrefix(line, "Description :"))
		}
	}

	if currentPatch.Name != "" {
		patches = append(patches, currentPatch)
	}

	return patches
}

func (a *App) getLinuxPatches() []PatchInfo {
	var patches []PatchInfo

	// 尝试使用 apt 命令（Debian/Ubuntu）
	aptCmd := `apt-history list | grep -E "^(Start-Date|Commandline)" | grep -v "apt-history"`
	aptOut, err := exec.Command("bash", "-c", aptCmd).CombinedOutput()
	if err == nil {
		lines := strings.Split(string(aptOut), "\n")
		var currentPatch PatchInfo
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Start-Date") {
				if currentPatch.Name != "" {
					patches = append(patches, currentPatch)
				}
				dateStr := strings.TrimSpace(strings.TrimPrefix(line, "Start-Date:"))
				if t, err := time.Parse("2006-01-02 15:04:05", dateStr); err == nil {
					currentPatch.Date = t.Format("2006-01-02 15:04:05")
				}
			} else if strings.HasPrefix(line, "Commandline") {
				currentPatch.Name = strings.TrimSpace(strings.TrimPrefix(line, "Commandline:"))
				currentPatch.Status = "已安装"
			}
		}
		if currentPatch.Name != "" {
			patches = append(patches, currentPatch)
		}
		return patches
	}

	// 尝试使用 yum 命令（RHEL/CentOS）
	yumCmd := `yum history list | grep -E "^[0-9]" | head -n 10`
	yumOut, err := exec.Command("bash", "-c", yumCmd).CombinedOutput()
	if err == nil {
		lines := strings.Split(string(yumOut), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				patches = append(patches, PatchInfo{
					Name:   fields[2],
					Date:   fields[1],
					Status: "已安装",
				})
			}
		}
		return patches
	}

	return nil
}

func (a *App) getMacOSPatches() []PatchInfo {
	out, err := exec.Command("softwareupdate", "--history").Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\n")
	var patches []PatchInfo
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Name") || strings.HasPrefix(line, "-") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		patches = append(patches, PatchInfo{
			Name:   fields[0],
			Date:   fields[len(fields)-2],
			Status: fields[len(fields)-1],
		})
	}
	return patches
}
