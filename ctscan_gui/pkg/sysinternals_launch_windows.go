//go:build windows

package pkg

import (
	"errors"
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func launchWindowsTool(path string, elevated bool) error {
	if !fileExists(path) {
		return fmt.Errorf("工具文件不存在: %s", path)
	}

	verb := "open"
	if elevated {
		verb = "runas"
	}

	if err := shellExecuteWindows(verb, path); err != nil {
		toolName := filepath.Base(path)
		if elevated && errors.Is(err, windows.ERROR_CANCELLED) {
			return fmt.Errorf("用户取消了管理员授权，未启动 %s", toolName)
		}
		if elevated {
			return fmt.Errorf("以管理员权限启动 %s 失败: %w；路径：%s", toolName, err, path)
		}
		return fmt.Errorf("启动 %s 失败: %w；路径：%s", toolName, err, path)
	}
	return nil
}

func shellExecuteWindows(verb string, path string) error {
	verbPtr, err := windows.UTF16PtrFromString(verb)
	if err != nil {
		return fmt.Errorf("初始化启动参数失败: %w", err)
	}
	filePtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return fmt.Errorf("初始化工具路径失败: %w", err)
	}
	workingDirPtr, err := windows.UTF16PtrFromString(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("初始化工作目录失败: %w", err)
	}

	if err := windows.ShellExecute(0, verbPtr, filePtr, nil, workingDirPtr, windows.SW_SHOWNORMAL); err != nil {
		return err
	}
	return nil
}
