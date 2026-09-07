//go:build !windows

package pkg

import "os/exec"

func backgroundCommand(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}
