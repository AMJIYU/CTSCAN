package pkg

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type StartupItem struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Type        string    `json:"type"`        // 启动项类型：LaunchAgent/LaunchDaemon/Startup等
	Enabled     bool      `json:"enabled"`     // 是否启用
	LastModTime time.Time `json:"lastModTime"` // 最后修改时间
	Size        int64     `json:"size"`        // 文件大小
	Description string    `json:"description"` // 描述信息
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
	paths := []struct {
		path string
		typ  string
	}{
		{"/Library/LaunchAgents", "LaunchAgent"},
		{"/Library/LaunchDaemons", "LaunchDaemon"},
		{filepath.Join(os.Getenv("HOME"), "Library/LaunchAgents"), "UserLaunchAgent"},
	}

	for _, p := range paths {
		files, err := ioutil.ReadDir(p.path)
		if err != nil {
			continue
		}
		for _, f := range files {
			if !f.IsDir() && filepath.Ext(f.Name()) == ".plist" {
				filePath := filepath.Join(p.path, f.Name())
				info, err := os.Stat(filePath)
				if err != nil {
					continue
				}

				// 读取plist文件内容
				content, err := ioutil.ReadFile(filePath)
				if err != nil {
					continue
				}

				// 解析plist文件获取描述信息
				description := ""
				if len(content) > 0 {
					// 简单的XML解析，获取Label或Description
					type Plist struct {
						Label       string `xml:"dict>key>Label"`
						Description string `xml:"dict>key>Description"`
					}
					var plist Plist
					if err := xml.Unmarshal(content, &plist); err == nil {
						if plist.Description != "" {
							description = plist.Description
						} else if plist.Label != "" {
							description = plist.Label
						}
					}
				}

				items = append(items, StartupItem{
					Name:        f.Name(),
					Path:        filePath,
					Type:        p.typ,
					Enabled:     true, // 默认启用
					LastModTime: info.ModTime(),
					Size:        info.Size(),
					Description: description,
				})
			}
		}
	}
	return items
}

func (a *App) getWindowsStartupItems() []StartupItem {
	var items []StartupItem

	if runtime.GOOS == "windows" {
		startupFolders := []struct {
			path string
			typ  string
		}{
			{filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\Start Menu\\Programs\\Startup"), "UserStartup"},
			{filepath.Join(os.Getenv("ProgramData"), "Microsoft\\Windows\\Start Menu\\Programs\\Startup"), "SystemStartup"},
		}

		for _, folder := range startupFolders {
			files, err := ioutil.ReadDir(folder.path)
			if err != nil {
				continue
			}
			for _, f := range files {
				if !f.IsDir() {
					filePath := filepath.Join(folder.path, f.Name())
					info, err := os.Stat(filePath)
					if err != nil {
						continue
					}

					items = append(items, StartupItem{
						Name:        f.Name(),
						Path:        filePath,
						Type:        folder.typ,
						Enabled:     true,
						LastModTime: info.ModTime(),
						Size:        info.Size(),
						Description: "", // Windows启动项可能需要额外解析快捷方式
					})
				}
			}
		}
	}

	return items
}

func (a *App) getLinuxStartupItems() []StartupItem {
	var items []StartupItem

	autostartDirs := []struct {
		path string
		typ  string
	}{
		{filepath.Join(os.Getenv("HOME"), ".config/autostart"), "UserAutostart"},
		{"/etc/xdg/autostart", "SystemAutostart"},
	}

	for _, dir := range autostartDirs {
		files, err := ioutil.ReadDir(dir.path)
		if err != nil {
			continue
		}
		for _, f := range files {
			if !f.IsDir() && filepath.Ext(f.Name()) == ".desktop" {
				filePath := filepath.Join(dir.path, f.Name())
				info, err := os.Stat(filePath)
				if err != nil {
					continue
				}

				// 读取.desktop文件内容
				content, err := ioutil.ReadFile(filePath)
				if err != nil {
					continue
				}

				// 解析.desktop文件获取描述信息
				description := ""
				if len(content) > 0 {
					// 简单的文本解析，获取Comment或Name
					lines := strings.Split(string(content), "\n")
					for _, line := range lines {
						if strings.HasPrefix(line, "Comment=") {
							description = strings.TrimPrefix(line, "Comment=")
							break
						} else if strings.HasPrefix(line, "Name=") && description == "" {
							description = strings.TrimPrefix(line, "Name=")
						}
					}
				}

				items = append(items, StartupItem{
					Name:        f.Name(),
					Path:        filePath,
					Type:        dir.typ,
					Enabled:     true,
					LastModTime: info.ModTime(),
					Size:        info.Size(),
					Description: description,
				})
			}
		}
	}

	return items
}
