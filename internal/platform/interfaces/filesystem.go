package interfaces

import "io/fs"

// FileSystemService defines the contract for platform-specific filesystem operations.
type FileSystemService interface {
	// EnsureDataDirectory creates the application's data directory if it doesn't exist
	// and returns its path.
	EnsureDataDirectory() (string, error)

	// GetConfigDirectory returns the path to the application's configuration directory.
	GetConfigDirectory() (string, error)

	// GetCacheDirectory returns the path to the application's cache directory.
	GetCacheDirectory() (string, error)

	// GetDataHome returns the base directory relative to which user-specific data files should be stored.
	GetDataHome() (string, error)

	// ReadFile reads the content of a file at the given path.
	ReadFile(path string) ([]byte, error)

	// WriteFile writes data to a file at the given path, creating it if necessary.
	// It uses perm to set file permissions.
	WriteFile(path string, data []byte, perm fs.FileMode) error

	// Stat returns a FileInfo describing the named file.
	Stat(path string) (fs.FileInfo, error)

	// MkdirAll creates a directory named path, along with any necessary parents,
	// and returns nil, or else returns an error.
	MkdirAll(path string, perm fs.FileMode) error

	// Remove deletes the named file or (empty) directory.
	Remove(path string) error

	// RemoveAll removes path and any children it contains.
	RemoveAll(path string) error

	// Exists checks if a file or directory exists at the given path.
	Exists(path string) (bool, error)
}
