//go:build darwin
// +build darwin

package data

import (
	"os/exec"
)

// checkChromeOnWindows 在 Mac 上的实现
func checkChromeOnWindows() (string, bool) {
	// 检查 Chrome 是否安装在默认位置
	chromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	_, err := exec.Command("ls", chromePath).Output()
	if err == nil {
		return chromePath, true
	}
	return "", false
}

// CheckBrowserOnWindows 在 Mac 上的实现
func CheckBrowserOnWindows() (string, bool) {
	// 首先检查 Chrome
	if path, exists := checkChromeOnWindows(); exists {
		return path, true
	}

	// 检查 Safari（Mac 默认浏览器）
	safariPath := "/Applications/Safari.app/Contents/MacOS/Safari"
	_, err := exec.Command("ls", safariPath).Output()
	if err == nil {
		return safariPath, true
	}

	return "", false
}
