//go:build !windows

package pkg

import "fmt"

func launchWindowsTool(path string, elevated bool) error {
	return fmt.Errorf("当前系统不能运行 Windows EXE: %s", path)
}
