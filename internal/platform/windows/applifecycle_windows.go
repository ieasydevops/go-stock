//go:build windows
// +build windows

package windows

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.AppLifecycleService = (*appLifecycleService)(nil)

type appLifecycleService struct {
	appName             string
	appVersion          string
	singleInstanceMutex windows.Handle
}

// NewAppLifecycleService creates a new service for Windows application lifecycle management.
// appName is used for creating a unique named mutex for single instance check.
func NewAppLifecycleService(appName string, appVersion string) interfaces.AppLifecycleService {
	return &appLifecycleService{
		appName:    appName,
		appVersion: appVersion,
	}
}

// GetAppName returns the application name.
func (s *appLifecycleService) GetAppName() (string, error) {
	return s.appName, nil
}

// GetAppVersion returns the application version.
func (s *appLifecycleService) GetAppVersion() (string, error) {
	return s.appVersion, nil
}

// Initialize performs platform-specific startup tasks.
func (s *appLifecycleService) Initialize() error {
	log.Println("Windows AppLifecycleService: Initialize")
	// Actual single instance lock request should ideally be here.
	// For example:
	// locked, err := s.RequestSingleInstanceLock()
	// if err != nil { return fmt.Errorf("failed instance lock check: %w", err) }
	// if !locked { return fmt.Errorf("another instance is already running") }
	return nil
}

// Shutdown performs platform-specific shutdown tasks.
func (s *appLifecycleService) Shutdown() error {
	log.Println("Windows AppLifecycleService: Shutdown")
	// Release the mutex if it was acquired
	if s.singleInstanceMutex != 0 {
		err := windows.CloseHandle(s.singleInstanceMutex)
		s.singleInstanceMutex = 0 // Mark as closed
		if err != nil {
			return fmt.Errorf("failed to release single instance mutex: %w", err)
		}
	}
	return nil
}

// RequestSingleInstanceLock attempts to acquire a system-wide lock to ensure only one instance runs.
// Returns true if the lock was acquired (this is the first instance), false otherwise.
// This is more of a helper or part of Initialize logic.
func (s *appLifecycleService) RequestSingleInstanceLock() (bool, error) {
	mutexName, err := syscall.UTF16PtrFromString(s.appName + "_InstanceMutex")
	if err != nil {
		return false, fmt.Errorf("failed to create UTF16 string for mutex name: %w", err)
	}

	handle, err := windows.CreateMutex(nil, true, mutexName)

	if err == windows.ERROR_ALREADY_EXISTS {
		if handle != 0 { // If ERROR_ALREADY_EXISTS, a handle to the existing mutex is returned. Close it.
			windows.CloseHandle(handle)
		}
		s.singleInstanceMutex = 0 // We didn't acquire it or it was already set by another instance.
		return false, nil         // Mutex already exists, so another instance is running
	}
	if err != nil {
		if handle != 0 { // Clean up handle if CreateMutex failed after returning one (unlikely for non-ERROR_ALREADY_EXISTS)
			windows.CloseHandle(handle)
		}
		s.singleInstanceMutex = 0
		return false, fmt.Errorf("failed to create/open single instance mutex: %w", err)
	}
	// If err is nil, we created it and have the lock.
	s.singleInstanceMutex = handle // Store the handle
	return true, nil
}

// HandleSecondInstance is called when another instance of the application is launched.
func (s *appLifecycleService) HandleSecondInstance(data interface{}) (bool, error) {
	log.Printf("Windows AppLifecycleService: HandleSecondInstance called with data: %v. This instance should activate.", data)
	// TODO: Implement logic to bring the current window to the foreground.
	// This typically involves finding the window handle (e.g. by class name/title) and using SetForegroundWindow.
	// For a Wails app, runtime.WindowShow() or similar on the main window might be sufficient if the Wails context is available.
	// Since this is a platform service, a more generic window finding mechanism might be needed if not directly tied to Wails window.
	return true, nil // True means this existing instance will handle it.
}

// RequestExit requests the application to close.
func (s *appLifecycleService) RequestExit() {
	log.Println("Windows AppLifecycleService: RequestExit called.")
	// This should ideally trigger a graceful shutdown. For a Wails app, this means calling runtime.Quit(app.ctx).
	// Since this is a platform service, a direct os.Exit might be too abrupt if not coordinated.
	// For now, just log. The actual exit should be managed by the application layer using this signal.
	// Alternatively, if this service has a way to signal the main app loop to quit gracefully, that would be better.
	// For a simple placeholder:
	os.Exit(0)
}

// RestartApplication attempts to restart the application.
// This is platform-specific and can be complex.
func (s *appLifecycleService) RestartApplication() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Release single instance lock if held by this instance before restarting
	if s.singleInstanceMutex != 0 {
		windows.CloseHandle(s.singleInstanceMutex)
		s.singleInstanceMutex = 0
	}

	cmd := exec.Command(exe)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start new application instance: %w", err)
	}

	log.Println("Restarting application...")
	os.Exit(0)
	return nil // Unreachable
}
