//go:build darwin
// +build darwin

package data

import (
	"os"
	"os/exec"
	"strings"
)

// CheckBrowserOnDarwin 在 macOS 系统上检查浏览器是否安装，并返回安装路径
// 由于在macOS上浏览器通常位于固定位置，所以这里返回常见位置
func CheckBrowserOnDarwin() (string, bool) {
	// 检查Chrome
	chromePaths := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
	}

	for _, path := range chromePaths {
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}

	// 检查Safari
	safariPath := "/Applications/Safari.app/Contents/MacOS/Safari"
	if _, err := os.Stat(safariPath); err == nil {
		return safariPath, true
	}

	// 检查是否有Edge
	edgePaths := []string{
		"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		"/Applications/Microsoft Edge Canary.app/Contents/MacOS/Microsoft Edge Canary",
		"/Applications/Microsoft Edge Dev.app/Contents/MacOS/Microsoft Edge Dev",
		"/Applications/Microsoft Edge Beta.app/Contents/MacOS/Microsoft Edge Beta",
	}

	for _, path := range edgePaths {
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}

	// 检查Firefox
	firefoxPath := "/Applications/Firefox.app/Contents/MacOS/firefox"
	if _, err := os.Stat(firefoxPath); err == nil {
		return firefoxPath, true
	}

	// 尝试使用which命令找到Chrome
	cmd := exec.Command("which", "google-chrome")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output)), true
	}

	// 如果没有找到任何浏览器，返回空字符串和false
	return "", false
}

// CheckBrowserOnWindows 在 macOS 系统上检查浏览器是否安装，并返回安装路径
// 由于函数名需要保持一致，但实际上实现的是macOS版本的检测
func CheckBrowserOnWindows() (string, bool) {
	// 检查Chrome
	chromePaths := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
	}

	for _, path := range chromePaths {
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}

	// 检查Safari
	safariPath := "/Applications/Safari.app/Contents/MacOS/Safari"
	if _, err := os.Stat(safariPath); err == nil {
		return safariPath, true
	}

	// 检查是否有Edge
	edgePaths := []string{
		"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		"/Applications/Microsoft Edge Canary.app/Contents/MacOS/Microsoft Edge Canary",
	}
	for _, path := range edgePaths {
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}

	// 检查Firefox
	firefoxPath := "/Applications/Firefox.app/Contents/MacOS/firefox"
	if _, err := os.Stat(firefoxPath); err == nil {
		return firefoxPath, true
	}

	// 找不到浏览器，返回false
	return "", false
}

// checkEdgeOnDarwin 检查Edge浏览器是否安装
func checkEdgeOnDarwin() (string, bool) {
	cmd := exec.Command("mdfind", "kMDItemCFBundleIdentifier == 'com.microsoft.edgemac'")
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		return "", false
	}

	paths := strings.Split(string(output), "\n")
	if len(paths) > 0 && paths[0] != "" {
		return paths[0] + "/Contents/MacOS/Microsoft Edge", true
	}

	return "", false
}
