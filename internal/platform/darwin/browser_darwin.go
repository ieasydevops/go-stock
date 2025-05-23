//go:build darwin
// +build darwin

package darwin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.BrowserDetectorService = (*browserDetectorService)(nil)

type browserDetectorService struct{}

// NewBrowserDetectorService creates a new service for Darwin (macOS) browser detection.
func NewBrowserDetectorService() interfaces.BrowserDetectorService {
	return &browserDetectorService{}
}

func (s *browserDetectorService) GetInstalledBrowsers() ([]interfaces.BrowserInfo, error) {
	browsers := []interfaces.BrowserInfo{}

	// Common browser paths on macOS
	browserPaths := map[string]string{
		"Safari":   "/Applications/Safari.app",
		"Chrome":   "/Applications/Google Chrome.app",
		"Firefox":  "/Applications/Firefox.app",
		"Edge":     "/Applications/Microsoft Edge.app",
		"Opera":    "/Applications/Opera.app",
		"Brave":    "/Applications/Brave Browser.app",
		"Vivaldi":  "/Applications/Vivaldi.app",
		"Chromium": "/Applications/Chromium.app",
	}

	for name, path := range browserPaths {
		if _, err := os.Stat(path); err == nil {
			version := s.getBrowserVersion(name, path)
			browsers = append(browsers, interfaces.BrowserInfo{
				Name:    name,
				Path:    path,
				Version: version,
			})
		}
	}

	return browsers, nil
}

func (s *browserDetectorService) GetDefaultBrowser() (interfaces.BrowserInfo, error) {
	// On macOS, we can use AppleScript to get the default browser
	script := `tell application "System Events" to return name of default web browser`
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return interfaces.BrowserInfo{}, fmt.Errorf("failed to get default browser: %w", err)
	}

	defaultBrowserName := strings.TrimSpace(string(output))
	// Remove ".app" suffix if present
	defaultBrowserName = strings.TrimSuffix(defaultBrowserName, ".app")

	browsers, err := s.GetInstalledBrowsers()
	if err != nil {
		return interfaces.BrowserInfo{}, fmt.Errorf("failed to get installed browsers: %w", err)
	}

	// Find the default browser in our list of installed browsers
	for _, browser := range browsers {
		if strings.Contains(strings.ToLower(browser.Name), strings.ToLower(defaultBrowserName)) {
			return browser, nil
		}
	}

	// If we can't match it with our known browsers, return a basic info with just the name
	return interfaces.BrowserInfo{
		Name: defaultBrowserName,
		Path: fmt.Sprintf("/Applications/%s.app", defaultBrowserName),
	}, nil
}

func (s *browserDetectorService) CheckBrowser(browserName string) (string, bool, error) {
	browsers, err := s.GetInstalledBrowsers()
	if err != nil {
		return "", false, fmt.Errorf("failed to get installed browsers: %w", err)
	}

	// Case-insensitive search for the browser
	browserNameLower := strings.ToLower(browserName)
	for _, browser := range browsers {
		if strings.Contains(strings.ToLower(browser.Name), browserNameLower) {
			return browser.Path, true, nil
		}
	}

	// Try direct path lookup for common browsers
	commonPaths := map[string]string{
		"safari":  "/Applications/Safari.app",
		"chrome":  "/Applications/Google Chrome.app",
		"firefox": "/Applications/Firefox.app",
		"edge":    "/Applications/Microsoft Edge.app",
		"brave":   "/Applications/Brave Browser.app",
	}

	if path, ok := commonPaths[browserNameLower]; ok {
		if _, err := os.Stat(path); err == nil {
			return path, true, nil
		}
	}

	// Check if it might be directly in Applications folder
	possiblePath := filepath.Join("/Applications", browserName+".app")
	if _, err := os.Stat(possiblePath); err == nil {
		return possiblePath, true, nil
	}

	return "", false, nil
}

// Helper function to get browser version
func (s *browserDetectorService) getBrowserVersion(name, path string) string {
	// This is a simplified approach - a full implementation would have browser-specific logic
	// For demonstration, we'll just return a placeholder version
	return "Unknown" // In a real implementation, extract version from Info.plist or other sources
}
