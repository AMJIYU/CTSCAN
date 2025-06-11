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

// GetStartupItems 获取系统启动项列表
// 这个方法被绑定了，所以可以从前端JS调用
func (a *App) GetStartupItems() []StartupItem {
	var items []StartupItem

	switch runtime.GOOS {
	case "darwin":
		items = a.getMacStartupItems()
	case "windows":
		items = a.getWindowsStartupItems()
	case "linux":
		items = a.getLinuxStartupItems()
	}

	return items
}

func (a *App) getMacStartupItems() []StartupItem {
	var items []StartupItem
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

func (a *App) getWindowsStartupItems() []StartupItem {
	var items []StartupItem

	if runtime.GOOS == "windows" {
		// 检查启动文件夹
		startupFolders := []string{
			filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\Start Menu\\Programs\\Startup"),
			filepath.Join(os.Getenv("ProgramData"), "Microsoft\\Windows\\Start Menu\\Programs\\Startup"),
		}

		for _, folder := range startupFolders {
			files, err := ioutil.ReadDir(folder)
			if err != nil {
				continue
			}
			for _, f := range files {
				items = append(items, StartupItem{
					Name: f.Name(),
					Path: filepath.Join(folder, f.Name()),
				})
			}
		}
	}

	return items
}

func (a *App) getLinuxStartupItems() []StartupItem {
	var items []StartupItem

	// 检查用户级自动启动目录
	userAutostart := filepath.Join(os.Getenv("HOME"), ".config/autostart")
	if files, err := ioutil.ReadDir(userAutostart); err == nil {
		for _, f := range files {
			items = append(items, StartupItem{
				Name: f.Name(),
				Path: filepath.Join(userAutostart, f.Name()),
			})
		}
	}

	// 检查系统级自动启动目录
	systemAutostart := "/etc/xdg/autostart"
	if files, err := ioutil.ReadDir(systemAutostart); err == nil {
		for _, f := range files {
			items = append(items, StartupItem{
				Name: f.Name(),
				Path: filepath.Join(systemAutostart, f.Name()),
			})
		}
	}

	return items
}
