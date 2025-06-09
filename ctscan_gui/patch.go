package main

import (
	"os/exec"
	"strings"
)

type PatchInfo struct {
	Name   string `json:"name"`
	Date   string `json:"date"`
	Status string `json:"status"`
}

func (a *App) GetPatchInfo() []PatchInfo {
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
