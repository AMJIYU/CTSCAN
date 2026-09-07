//go:build !windows

package pkg

func (a *App) getWindowsStartupItems() []StartupItem {
	return nil
}
