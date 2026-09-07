//go:build !windows

package pkg

import (
	"context"
	"os/exec"
)

func backgroundCommand(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

func backgroundCommandContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, args...)
}
