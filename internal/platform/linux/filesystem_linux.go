//go:build linux
// +build linux

package linux

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.FileSystemService = (*fileSystemService)(nil)

type fileSystemService struct {
	appName string
}

// NewFileSystemService creates a new filesystem service for Linux.
// AppName is used to create app-specific subdirectories in standard locations.
func NewFileSystemService(appName string) interfaces.FileSystemService {
	return &fileSystemService{appName: appName}
}

// getXDGDirectory retrieves XDG directory paths with fallbacks
func (s *fileSystemService) getXDGDirectory(envVar, fallbackRelPath string) (string, error) {
	// Check XDG_* environment variable first
	if dir := os.Getenv(envVar); dir != "" {
		return filepath.Join(dir, s.appName), nil
	}

	// Fallback to standard location
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, fallbackRelPath, s.appName), nil
}

// EnsureDataDirectory ensures the application-specific data directory exists.
// For Linux, this follows XDG Base Directory spec (~/.local/share/AppName)
func (s *fileSystemService) EnsureDataDirectory() (string, error) {
	dataDir, err := s.GetDataHome()
	if err != nil {
		return "", fmt.Errorf("failed to determine data directory path: %w", err)
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create data directory %s: %w", dataDir, err)
	}
	return dataDir, nil
}

// GetDataHome returns the base directory for user-specific data files.
// On Linux, this is $XDG_DATA_HOME/AppName or ~/.local/share/AppName
func (s *fileSystemService) GetDataHome() (string, error) {
	return s.getXDGDirectory("XDG_DATA_HOME", ".local/share")
}

// GetConfigDirectory returns the path to the application's configuration directory.
// On Linux, this is $XDG_CONFIG_HOME/AppName or ~/.config/AppName
func (s *fileSystemService) GetConfigDirectory() (string, error) {
	configDir, err := s.getXDGDirectory("XDG_CONFIG_HOME", ".config")
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}
	return configDir, nil
}

// GetCacheDirectory returns the path to the application's cache directory.
// On Linux, this is $XDG_CACHE_HOME/AppName or ~/.cache/AppName
func (s *fileSystemService) GetCacheDirectory() (string, error) {
	cacheDir, err := s.getXDGDirectory("XDG_CACHE_HOME", ".cache")
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory %s: %w", cacheDir, err)
	}
	return cacheDir, nil
}

// ReadFile reads the content of a file at the given path.
func (s *fileSystemService) ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return data, nil
}

// WriteFile writes data to a file at the given path, creating it if necessary.
func (s *fileSystemService) WriteFile(path string, data []byte, perm fs.FileMode) error {
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	return nil
}

// Stat returns a FileInfo describing the named file.
func (s *fileSystemService) Stat(path string) (fs.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file information for %s: %w", path, err)
	}
	return info, nil
}

// MkdirAll creates a directory named path, along with any necessary parents.
func (s *fileSystemService) MkdirAll(path string, perm fs.FileMode) error {
	if err := os.MkdirAll(path, perm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

// Remove deletes the named file or (empty) directory.
func (s *fileSystemService) Remove(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove %s: %w", path, err)
	}
	return nil
}

// RemoveAll removes path and any children it contains.
func (s *fileSystemService) RemoveAll(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to recursively remove %s: %w", path, err)
	}
	return nil
}

// Exists checks if a file or directory exists at the given path.
func (s *fileSystemService) Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		// Some other error occurred (e.g., permission issue)
		return false, fmt.Errorf("failed to check existence of %s: %w", path, err)
	}
}
