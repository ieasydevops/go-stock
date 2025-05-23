//go:build windows
// +build windows

package windows

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.NetworkService = (*networkService)(nil)

type networkService struct {
	httpClient *http.Client
}

// NewNetworkService creates a new service for Windows network operations.
func NewNetworkService() interfaces.NetworkService {
	// TODO: Consider platform-specific http client configurations if needed (e.g., system proxy)
	// For Windows, this might involve querying registry for proxy settings, but that's complex.
	// For now, a default client is provided.
	return &networkService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // Default timeout
		},
	}
}

// GetHTTPClient returns a configured HTTP client.
func (s *networkService) GetHTTPClient() (*http.Client, error) {
	return s.httpClient, nil
}

// CheckConnectivity checks if there is an active internet connection.
// It attempts a DNS lookup on a reliable host (e.g., google.com) and then an HTTP HEAD request.
func (s *networkService) CheckConnectivity(hostsToTest ...string) (bool, error) {
	hosts := hostsToTest
	if len(hosts) == 0 {
		hosts = []string{"www.google.com", "www.bing.com"} // Default reliable hosts
	}

	var lastErr error
	for _, host := range hosts {
		// First, try DNS lookup
		_, err := net.LookupHost(host)
		if err == nil {
			// If DNS lookup succeeds, try a lightweight HTTP HEAD request to confirm HTTP connectivity
			// Ensure to use https for these well-known hosts.
			protocolHost := host
			if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
				protocolHost = "https://" + host
			}
			req, reqErr := http.NewRequest(http.MethodHead, protocolHost, nil)
			if reqErr != nil {
				lastErr = fmt.Errorf("failed to create HEAD request for %s: %w", host, reqErr)
				continue // Try next host
			}

			resp, err := s.httpClient.Do(req)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode >= 200 && resp.StatusCode < 300 { // Check for 2xx success
					return true, nil // Connectivity confirmed
				} else {
					lastErr = fmt.Errorf("HEAD request to %s returned non-2xx status: %d", host, resp.StatusCode)
				}
			} else {
				lastErr = fmt.Errorf("HEAD request to %s failed: %w", host, err)
			}
		} else {
			lastErr = fmt.Errorf("DNS lookup for %s failed: %w", host, err)
		}
	}

	if lastErr != nil {
		return false, fmt.Errorf("connectivity check failed for all hosts. Last error: %w", lastErr)
	}
	return false, nil // Should not be reached if hosts list is not empty and logic is correct
}

// DownloadFile downloads a file from a URL to a local path.
func (s *networkService) DownloadFile(url, localPath string, progressCallback func(current, total int64)) error {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to start download from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download from %s failed with status: %s", url, resp.Status)
	}

	out, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	defer out.Close()

	totalSize := resp.ContentLength // Can be -1 if server doesn't send it
	var currentSize int64

	// Use a buffer for copying
	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			written, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("error writing to file %s: %w", localPath, writeErr)
			}
			currentSize += int64(written)
			if progressCallback != nil {
				progressCallback(currentSize, totalSize)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading from download stream %s: %w", url, err)
		}
	}

	// Final progress update if totalSize was known and we reached it (or if totalSize was -1)
	if progressCallback != nil {
		if totalSize == -1 || currentSize == totalSize { // Ensure final call, especially if totalSize was unknown
			progressCallback(currentSize, totalSize)
		}
	}

	return nil
}
