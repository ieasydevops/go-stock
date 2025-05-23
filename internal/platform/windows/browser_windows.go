//go:build windows
// +build windows

package windows

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.BrowserDetectorService = (*browserDetectorService)(nil)

type browserDetectorService struct{}

// NewBrowserDetectorService creates a new service for Windows browser detection.
func NewBrowserDetectorService() interfaces.BrowserDetectorService {
	return &browserDetectorService{}
}

// GetInstalledBrowsers attempts to find common installed browsers on Windows.
// This is a best-effort approach and might not find all browsers or portable versions.
func (s *browserDetectorService) GetInstalledBrowsers() ([]interfaces.BrowserInfo, error) {
	var browsers []interfaces.BrowserInfo

	// Common browser paths and registry checks
	// Google Chrome
	if path, err := s.getBrowserPathFromRegistry(`SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe`, ""); err == nil && path != "" {
		browsers = append(browsers, interfaces.BrowserInfo{Name: "Google Chrome", Path: path})
	} else if path, err := s.getBrowserPathFromKey(registry.LOCAL_MACHINE, `SOFTWARE\Clients\StartMenuInternet\Google Chrome\shell\open\command`, ""); err == nil && path != "" {
		browsers = append(browsers, interfaces.BrowserInfo{Name: "Google Chrome", Path: path})
	}

	// Mozilla Firefox
	if path, err := s.getBrowserPathFromKey(registry.LOCAL_MACHINE, `SOFTWARE\Clients\StartMenuInternet\FIREFOX.EXE\shell\open\command`, ""); err == nil && path != "" {
		browsers = append(browsers, interfaces.BrowserInfo{Name: "Mozilla Firefox", Path: path})
	} else if path, err := s.getBrowserPathFromKey(registry.CLASSES_ROOT, `FirefoxHTML\shell\open\command`, ""); err == nil && path != "" {
		browsers = append(browsers, interfaces.BrowserInfo{Name: "Mozilla Firefox", Path: path})
	}

	// Microsoft Edge (New Chromium-based)
	if path, err := s.getBrowserPathFromRegistry(`SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\msedge.exe`, ""); err == nil && path != "" {
		browsers = append(browsers, interfaces.BrowserInfo{Name: "Microsoft Edge", Path: path})
	} else if path, err := s.getBrowserPathFromKey(registry.LOCAL_MACHINE, `SOFTWARE\Clients\StartMenuInternet\Microsoft Edge\shell\open\command`, ""); err == nil && path != "" {
		browsers = append(browsers, interfaces.BrowserInfo{Name: "Microsoft Edge", Path: path})
	}

	// TODO: Add checks for other browsers like Opera, Vivaldi, Brave if necessary

	// Deduplicate and ensure paths exist
	uniqueBrowsers := []interfaces.BrowserInfo{}
	seenPaths := make(map[string]bool)
	for _, b := range browsers {
		cleanedPath := cleanBrowserPath(b.Path)
		if cleanedPath == "" {
			continue
		}
		if _, err := os.Stat(cleanedPath); err == nil {
			if !seenPaths[cleanedPath] {
				b.Path = cleanedPath
				// TODO: Get browser version if possible (platform-specific and browser-specific)
				uniqueBrowsers = append(uniqueBrowsers, b)
				seenPaths[cleanedPath] = true
			}
		}
	}

	return uniqueBrowsers, nil
}

// GetDefaultBrowser retrieves the system's default web browser.
func (s *browserDetectorService) GetDefaultBrowser() (interfaces.BrowserInfo, error) {
	// The default browser is typically determined by the handler for the http protocol.
	k, err := registry.OpenKey(registry.CLASSES_ROOT, `HTTP\shell\open\command`, registry.QUERY_VALUE)
	if err != nil {
		// Fallback: Check HTTPS as well
		k, err = registry.OpenKey(registry.CLASSES_ROOT, `HTTPS\shell\open\command`, registry.QUERY_VALUE)
		if err != nil {
			return interfaces.BrowserInfo{}, fmt.Errorf("failed to open registry key for HTTP/HTTPS protocol: %w", err)
		}
	}
	defer k.Close()

	cmd, _, err := k.GetStringValue("")
	if err != nil {
		return interfaces.BrowserInfo{}, fmt.Errorf("failed to read default browser command from registry: %w", err)
	}

	path := cleanBrowserPath(cmd)
	if path == "" {
		return interfaces.BrowserInfo{}, fmt.Errorf("could not determine default browser path from command: %s", cmd)
	}

	// Attempt to get a more friendly name
	browserName := "Default Browser"
	// This is a heuristic. We could try to match `path` against known browser executable names.
	if strings.Contains(strings.ToLower(filepath.Base(path)), "chrome.exe") {
		browserName = "Google Chrome"
	} else if strings.Contains(strings.ToLower(filepath.Base(path)), "firefox.exe") {
		browserName = "Mozilla Firefox"
	} else if strings.Contains(strings.ToLower(filepath.Base(path)), "msedge.exe") {
		browserName = "Microsoft Edge"
	}

	return interfaces.BrowserInfo{Name: browserName, Path: path}, nil
}

// CheckBrowser is kept for compatibility with the interface, but GetInstalledBrowsers and GetDefaultBrowser are generally more useful.
func (s *browserDetectorService) CheckBrowser(browserName string) (path string, exists bool, err error) {
	browsers, err := s.GetInstalledBrowsers()
	if err != nil {
		return "", false, err
	}
	for _, b := range browsers {
		// Simple name check, might need to be more robust (e.g., case-insensitive, check for aliases)
		if strings.Contains(strings.ToLower(b.Name), strings.ToLower(browserName)) {
			return b.Path, true, nil
		}
	}
	return "", false, nil
}

// getBrowserPathFromRegistry reads a browser path from a specific registry key and value.
func (s *browserDetectorService) getBrowserPathFromRegistry(regPath, valueName string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.QUERY_VALUE)
	if err != nil {
		// Try CURRENT_USER as well for some App Paths
		k, err = registry.OpenKey(registry.CURRENT_USER, regPath, registry.QUERY_VALUE)
		if err != nil {
			return "", err
		}
	}
	defer k.Close()

	path, _, err := k.GetStringValue(valueName)
	if err != nil {
		return "", err
	}
	return cleanBrowserPath(path), nil
}

func (s *browserDetectorService) getBrowserPathFromKey(baseKey registry.Key, regPath, valueName string) (string, error) {
	k, err := registry.OpenKey(baseKey, regPath, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	path, _, err := k.GetStringValue(valueName)
	if err != nil {
		return "", err
	}
	return cleanBrowserPath(path), nil
}

// cleanBrowserPath extracts the executable path from a command string (e.g., from registry).
// It handles quotes and arguments like "%1".
func cleanBrowserPath(cmd string) string {
	if cmd == "" {
		return ""
	}
	path := strings.TrimSpace(cmd)

	// Remove common placeholders like "%1", "-- "%L" etc.
	path = strings.ReplaceAll(path, "\"%1\"", "")
	path = strings.ReplaceAll(path, "%1", "")
	path = strings.ReplaceAll(path, "\"%L\"", "")
	path = strings.ReplaceAll(path, "%L", "")
	path = strings.ReplaceAll(path, "--", "") // simple removal of -- argument prefix

	// Handle paths enclosed in quotes
	if strings.HasPrefix(path, "\"") {
		idx := strings.Index(path[1:], "\"")
		if idx != -1 {
			path = path[1 : idx+1]
		} else {
			// Unmatched quote, maybe it's just the exe name with args
			parts := strings.Fields(path)
			if len(parts) > 0 {
				path = strings.Trim(parts[0], "\"")
			}
		}
	} else {
		// If not quoted, the path might be the first part before a space (if arguments follow)
		parts := strings.Fields(path)
		if len(parts) > 0 {
			path = parts[0]
		}
	}
	return strings.TrimSpace(path)
}
