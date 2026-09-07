//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

const (
	messageBoxIconError = 0x00000010
	messageBoxOK        = 0x00000000
)

func showStartupError(title, message string) {
	titlePtr, titleErr := syscall.UTF16PtrFromString(title)
	messagePtr, messageErr := syscall.UTF16PtrFromString(message)
	if titleErr != nil || messageErr != nil {
		return
	}

	user32 := syscall.NewLazyDLL("user32.dll")
	messageBoxW := user32.NewProc("MessageBoxW")
	_, _, _ = messageBoxW.Call(
		0,
		uintptr(unsafe.Pointer(messagePtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(messageBoxOK|messageBoxIconError),
	)
}
