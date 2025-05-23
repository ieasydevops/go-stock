//go:build linux
// +build linux

package linux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.BrowserDetectorService = (*browserDetectorService)(nil)

type browserDetectorService struct{}

// NewBrowserDetectorService creates a new service for Linux browser detection.
func NewBrowserDetectorService() interfaces.BrowserDetectorService {
	return &browserDetectorService{}
}

func (s *browserDetectorService) GetInstalledBrowsers() ([]interfaces.BrowserInfo, error) {
	browsers := []interfaces.BrowserInfo{}

	// Common browser executables on Linux
	browserNames := []string{
		"firefox",
		"google-chrome",
		"chromium-browser",
		"chromium",
		"opera",
		"brave-browser",
		"vivaldi-stable",
		"epiphany", // GNOME Web
		"konqueror",
		"falkon",
		"midori",
	}

	for _, name := range browserNames {
		// Use 'which' to find the executable in PATH
		path, err := exec.Command("which", name).Output()
		if err == nil && len(path) > 0 {
			// Trim newline from path
			browserPath := strings.TrimSpace(string(path))
			version := s.getBrowserVersion(name, browserPath)

			// Determine proper display name based on executable
			displayName := s.getDisplayName(name)

			browsers = append(browsers, interfaces.BrowserInfo{
				Name:    displayName,
				Path:    browserPath,
				Version: version,
			})
		}
	}

	return browsers, nil
}

func (s *browserDetectorService) GetDefaultBrowser() (interfaces.BrowserInfo, error) {
	// On Linux, we can use xdg-mime query to get the default browser
	cmd := exec.Command("xdg-mime", "query", "default", "x-scheme-handler/http")
	output, err := cmd.Output()
	if err != nil {
		return interfaces.BrowserInfo{}, fmt.Errorf("failed to get default browser: %w", err)
	}

	// The output will be something like "firefox.desktop"
	defaultBrowserDesktop := strings.TrimSpace(string(output))

	// Extract the browser name from the .desktop file
	browserName := strings.TrimSuffix(defaultBrowserDesktop, ".desktop")

	// Map common .desktop file names to executable names
	executableMapping := map[string]string{
		"firefox":       "firefox",
		"google-chrome": "google-chrome",
		"chromium":      "chromium",
		"brave-browser": "brave-browser",
		"opera":         "opera",
	}

	executable, ok := executableMapping[browserName]
	if !ok {
		// If we don't have a mapping, just use the browser name
		executable = browserName
	}

	// Find the path using which
	pathCmd := exec.Command("which", executable)
	pathOutput, err := pathCmd.Output()
	if err != nil {
		return interfaces.BrowserInfo{}, fmt.Errorf("failed to find path for default browser: %w", err)
	}

	browserPath := strings.TrimSpace(string(pathOutput))
	displayName := s.getDisplayName(browserName)
	version := s.getBrowserVersion(browserName, browserPath)

	return interfaces.BrowserInfo{
		Name:    displayName,
		Path:    browserPath,
		Version: version,
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

	// Try direct which lookup
	cmd := exec.Command("which", browserNameLower)
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		path := strings.TrimSpace(string(output))
		if _, err := os.Stat(path); err == nil {
			return path, true, nil
		}
	}

	return "", false, nil
}

// Helper function to get browser version
func (s *browserDetectorService) getBrowserVersion(name, path string) string {
	var cmd *exec.Cmd

	switch name {
	case "firefox":
		cmd = exec.Command(path, "--version")
	case "google-chrome", "chromium", "chromium-browser":
		cmd = exec.Command(path, "--version")
	case "brave-browser":
		cmd = exec.Command(path, "--version")
	default:
		// Generic attempt with --version
		cmd = exec.Command(path, "--version")
	}

	output, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	version := strings.TrimSpace(string(output))

	// Many browsers output "BrowserName Version x.y.z"
	// Try to extract just the version part
	parts := strings.Split(version, " ")
	if len(parts) > 1 {
		// Try to find a part that looks like a version (contains a dot)
		for _, part := range parts {
			if strings.Contains(part, ".") {
				return part
			}
		}
		// If we can't find a part with a dot, return the last part
		return parts[len(parts)-1]
	}

	return version
}

// getDisplayName returns a friendly display name for the browser executable
func (s *browserDetectorService) getDisplayName(executable string) string {
	displayNames := map[string]string{
		"firefox":          "Mozilla Firefox",
		"google-chrome":    "Google Chrome",
		"chromium-browser": "Chromium",
		"chromium":         "Chromium",
		"opera":            "Opera",
		"brave-browser":    "Brave Browser",
		"vivaldi-stable":   "Vivaldi",
		"epiphany":         "GNOME Web",
		"konqueror":        "Konqueror",
		"falkon":           "Falkon",
		"midori":           "Midori",
	}

	if name, ok := displayNames[executable]; ok {
		return name
	}

	// Capitalize first letter and convert dashes to spaces for unknown browsers
	name := strings.ReplaceAll(executable, "-", " ")
	if len(name) > 0 {
		return strings.ToUpper(name[:1]) + name[1:]
	}

	return executable
}
