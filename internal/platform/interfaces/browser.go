package interfaces

// BrowserInfo holds information about an installed browser.
type BrowserInfo struct {
	Name    string // e.g., "Google Chrome", "Mozilla Firefox"
	Path    string // Path to the executable
	Version string // Version string, if available
}

// BrowserDetectorService defines the contract for detecting installed web browsers.
type BrowserDetectorService interface {
	// GetInstalledBrowsers returns a list of detected browsers on the system.
	GetInstalledBrowsers() ([]BrowserInfo, error)

	// GetDefaultBrowser returns information about the system's default web browser.
	GetDefaultBrowser() (BrowserInfo, error)

	// CheckBrowser checks for a specific browser (e.g., by name or a known path).
	// This was in the original design doc, a more generic GetInstalledBrowsers might be better.
	// Keep for now if specific checks are needed.
	// Returns path and existence status.
	CheckBrowser(browserName string) (path string, exists bool, err error)
}
