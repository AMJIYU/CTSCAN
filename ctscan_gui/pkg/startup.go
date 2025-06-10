package pkg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type StartupItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (a *App) GetStartupItems() []StartupItem {
	var items []StartupItem
	if runtime.GOOS != "darwin" {
		return items // 仅macOS实现
	}
	paths := []string{
		"/Library/LaunchAgents",
		"/Library/LaunchDaemons",
		filepath.Join(os.Getenv("HOME"), "Library/LaunchAgents"),
	}
	for _, dir := range paths {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, f := range files {
			items = append(items, StartupItem{
				Name: f.Name(),
				Path: filepath.Join(dir, f.Name()),
			})
		}
	}
	return items
}
