//go:build !windows

package pkg

func (a *App) getWindowsUsers() []SystemUser {
	return nil
}
