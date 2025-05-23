//go:build windows
// +build windows

package windows

import (
	"io/fs"
	"os"
	"path/filepath"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.FileSystemService = (*fileSystemService)(nil)

type fileSystemService struct{}

func NewFileSystemService() interfaces.FileSystemService {
	return &fileSystemService{}
}

func (s *fileSystemService) getDataDirName() string {
	// TODO: Make this configurable, perhaps via AppLifecycleService.GetAppName()
	return "GoStock" // Or a more appropriate name from a central config
}

// EnsureDataDirectory creates the application's data directory if it doesn't exist
// and returns its path. For Windows, this is typically in %APPDATA%.
func (s *fileSystemService) EnsureDataDirectory() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		// Fallback or error if APPDATA is not set, though it usually is.
		// For robust solution, consider using os.UserConfigDir() and then appending app name
		// However, os.UserConfigDir() points to %APPDATA% on windows.
		return "", os.ErrNotExist // Or a more specific error
	}
	dataDir := filepath.Join(appData, s.getDataDirName(), "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}
	return dataDir, nil
}

// GetConfigDirectory returns the path to the application's configuration directory.
// On Windows, this is often the same as EnsureDataDirectory or a subdirectory within it.
func (s *fileSystemService) GetConfigDirectory() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", os.ErrNotExist
	}
	configDir := filepath.Join(appData, s.getDataDirName(), "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return configDir, nil
}

// GetCacheDirectory returns the path to the application's cache directory.
// On Windows, this is typically in %LOCALAPPDATA%.
func (s *fileSystemService) GetCacheDirectory() (string, error) {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		// Fallback if LOCALAPPDATA is not set
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", os.ErrNotExist
		}
		// Attempt to construct Local AppData path from Roaming AppData
		// %APPDATA% is typically C:\Users\<user>\AppData\Roaming
		// %LOCALAPPDATA% is C:\Users\<user>\AppData\Local
		localAppData = filepath.Join(filepath.Dir(appData), "Local")
	}
	cacheDir := filepath.Join(localAppData, s.getDataDirName(), "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}
	return cacheDir, nil
}

// GetDataHome is more relevant for XDG on Linux. For Windows, APPDATA (Roaming) is common for user-specific data.
// We can return the general app data directory here (under APPDATA).
func (s *fileSystemService) GetDataHome() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", os.ErrNotExist
	}
	dataHome := filepath.Join(appData, s.getDataDirName())
	if err := os.MkdirAll(dataHome, 0755); err != nil {
		return "", err
	}
	return dataHome, nil
}

func (s *fileSystemService) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (s *fileSystemService) WriteFile(path string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (s *fileSystemService) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

func (s *fileSystemService) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (s *fileSystemService) Remove(path string) error {
	return os.Remove(path)
}

func (s *fileSystemService) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (s *fileSystemService) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

package p 