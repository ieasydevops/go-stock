//go:build windows
// +build windows

package data

import (
	"golang.org/x/sys/windows/registry"
)

// checkChromeOnWindows 在 Windows 上的实现
func checkChromeOnWindows() (string, bool) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe`, registry.QUERY_VALUE)
	if err != nil {
		return "", false
	}
	defer k.Close()

	s, _, err := k.GetStringValue("")
	if err != nil {
		return "", false
	}
	return s, true
}

// CheckBrowserOnWindows 在 Windows 上的实现
func CheckBrowserOnWindows() (string, bool) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\msedge.exe`, registry.QUERY_VALUE)
	if err != nil {
		return "", false
	}
	defer k.Close()

	s, _, err := k.GetStringValue("")
	if err != nil {
		return "", false
	}
	return s, true
}
