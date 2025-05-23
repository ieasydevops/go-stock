//go:build darwin
// +build darwin

package darwin

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

// NewFileSystemService creates a new filesystem service for macOS.
// AppName is used to create app-specific subdirectories in standard locations.
func NewFileSystemService(appName string) interfaces.FileSystemService {
	return &fileSystemService{appName: appName}
}

// EnsureDataDirectory ensures the application-specific data directory exists.
// For macOS, this is typically ~/Library/Application Support/AppName.
func (s *fileSystemService) EnsureDataDirectory() (string, error) {
	configDir, err := s.GetDataDirectory() // Corrected from GetAppDataDir
	if err != nil {
		return "", fmt.Errorf("failed to determine app data directory path: %w", err)
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app data directory %s: %w", configDir, err)
	}
	return configDir, nil
}

// GetDataHome returns the base directory for user-specific data files.
// On macOS, this is typically ~/Library/Application Support/AppName.
func (s *fileSystemService) GetDataHome() (string, error) {
	return s.GetConfigDirectory() // Leverages existing logic for this path
}

func (s *fileSystemService) GetCacheDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	// macOS cache directory: ~/Library/Caches/AppName
	cacheDir := filepath.Join(homeDir, "Library", "Caches", s.appName)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory %s: %w", cacheDir, err)
	}
	return cacheDir, nil
}

func (s *fileSystemService) GetConfigDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	// macOS application support directory (for config, data): ~/Library/Application Support/AppName
	configDir := filepath.Join(homeDir, "Library", "Application Support", s.appName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}
	return configDir, nil
}

func (s *fileSystemService) GetDataDirectory() (string, error) {
	// For macOS, AppDataDir is often the same as AppConfigDir or a subdirectory within it.
	return s.GetConfigDirectory() // Or a specific "Data" subfolder if desired: filepath.Join(configDir, "Data")
}

func (s *fileSystemService) GetLogDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	// macOS log directory: ~/Library/Logs/AppName
	logDir := filepath.Join(homeDir, "Library", "Logs", s.appName)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create log directory %s: %w", logDir, err)
	}
	return logDir, nil
}

func (s *fileSystemService) MkdirAll(dirPath string, perm fs.FileMode) error {
	if err := os.MkdirAll(dirPath, perm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}
	return nil
}

func (s *fileSystemService) Exists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		// Some other error occurred (e.g., permission issue)
		return false, fmt.Errorf("failed to check existence of %s: %w", filePath, err)
	}
}

func (s *fileSystemService) ReadFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return data, nil
}

func (s *fileSystemService) WriteFile(filePath string, data []byte, perm fs.FileMode) error {
	if err := os.WriteFile(filePath, data, perm); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}
	return nil
}

func (s *fileSystemService) GetUserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (s *fileSystemService) GetDesktopDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, "Desktop"), nil
}

func (s *fileSystemService) GetDownloadsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, "Downloads"), nil
}

func (s *fileSystemService) GetDocumentsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, "Documents"), nil
}

// Remove removes the named file or directory.
// If path refers to a directory, it removes the directory and its contents recursively.
func (s *fileSystemService) Remove(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove %s: %w", path, err)
	}
	return nil
}

// RemoveAll递归删除指定路径下的所有文件和目录
func (s *fileSystemService) RemoveAll(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to recursively remove %s: %w", path, err)
	}
	return nil
}

// Stat返回文件的信息
func (s *fileSystemService) Stat(path string) (fs.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file information for %s: %w", path, err)
	}
	return info, nil
}
