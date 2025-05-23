//go:build darwin
// +build darwin

package darwin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.AppLifecycleService = (*appLifecycleService)(nil)

type appLifecycleService struct {
	appName    string
	appVersion string
	lockFile   *os.File
	lockPath   string
}

// NewAppLifecycleService creates a new service for Darwin (macOS) application lifecycle management.
func NewAppLifecycleService(appName, appVersion string) interfaces.AppLifecycleService {
	// Determine a good lock file location (typically in /tmp or application data dir)
	tmpDir := os.TempDir()
	lockPath := filepath.Join(tmpDir, appName+".lock")

	return &appLifecycleService{
		appName:    appName,
		appVersion: appVersion,
		lockPath:   lockPath,
	}
}

func (s *appLifecycleService) Initialize() error {
	// Try to acquire app lock to ensure single instance
	lockFile, err := os.OpenFile(s.lockPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		if os.IsExist(err) {
			// Lock file exists, another instance is running
			return fmt.Errorf("another instance is already running")
		}
		return fmt.Errorf("failed to create lock file: %w", err)
	}

	// Write our PID to the lock file
	pid := os.Getpid()
	if _, err := fmt.Fprintf(lockFile, "%d", pid); err != nil {
		_ = lockFile.Close()
		_ = os.Remove(s.lockPath)
		return fmt.Errorf("failed to write PID to lock file: %w", err)
	}

	// Keep lock file open for the duration of the app
	s.lockFile = lockFile
	return nil
}

func (s *appLifecycleService) Shutdown() error {
	// Cleanup lock file
	if s.lockFile != nil {
		_ = s.lockFile.Close()
		_ = os.Remove(s.lockPath)
	}
	return nil
}

func (s *appLifecycleService) HandleSecondInstance(data interface{}) (bool, error) {
	// On macOS, we'd typically use AppleScript to bring the app to foreground
	// We need to know the bundle ID, which we can derive from appName in a real implementation
	// For now, let's just log a message

	// This sample shows how to activate using AppleScript, but will only work with a real app bundle
	bundleID := fmt.Sprintf("com.example.%s", s.appName) // Example, replace with real bundleID
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`tell application id "%s" to activate`, bundleID))
	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("failed to activate application: %w", err)
	}

	return true, nil
}

func (s *appLifecycleService) GetAppVersion() (string, error) {
	return s.appVersion, nil
}

func (s *appLifecycleService) GetAppName() (string, error) {
	return s.appName, nil
}

func (s *appLifecycleService) RequestExit() {
	// Send interrupt signal to self
	syscall.Kill(os.Getpid(), syscall.SIGINT)
}

func (s *appLifecycleService) RestartApplication() error {
	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Start the application in a detached process
	cmd := exec.Command(execPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // Detach from parent process
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to restart application: %w", err)
	}

	// Exit current process
	go func() {
		// Give new process a moment to start
		// This is a simplistic approach; a more robust solution would
		// involve IPC to confirm the new process is ready before exiting
		s.RequestExit()
	}()

	return nil
}
