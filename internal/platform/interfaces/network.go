package interfaces

import "net/http"

// NetworkService defines the contract for network-related operations.
type NetworkService interface {
	// GetHTTPClient returns a configured HTTP client, potentially with platform-specific settings
	// (e.g., proxy, CA certificates).
	GetHTTPClient() (*http.Client, error)

	// CheckConnectivity checks if there is an active internet connection.
	CheckConnectivity(hostToCheck ...string) (bool, error)

	// DownloadFile downloads a file from a URL to a local path.
	// It should handle progress reporting if possible (e.g., via a callback).
	DownloadFile(url, localPath string, progressCallback func(current, total int64)) error
}
