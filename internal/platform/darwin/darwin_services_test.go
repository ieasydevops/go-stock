//go:build darwin
// +build darwin

package darwin

import (
	"testing"
)

func TestNewServiceFactory(t *testing.T) {
	factory := NewServiceFactory("TestApp", "1.0.0")

	// Test logger service
	logger := factory.CreateLogger()
	if logger == nil {
		t.Error("Logger service should not be nil")
	}

	// Test file system service
	fs := factory.CreateFileSystemService()
	if fs == nil {
		t.Error("FileSystem service should not be nil")
	}

	// Test notification service
	notification := factory.CreateNotificationService()
	if notification == nil {
		t.Error("Notification service should not be nil")
	}

	// Test dialog service
	dialog := factory.CreateDialogService()
	if dialog == nil {
		t.Error("Dialog service should not be nil")
	}

	// Test system tray service
	tray := factory.CreateSystemTrayService()
	if tray == nil {
		t.Error("SystemTray service should not be nil")
	}

	// Test screen info service
	screen := factory.CreateScreenInfoService()
	if screen == nil {
		t.Error("ScreenInfo service should not be nil")
	}

	// Test process manager service
	process := factory.CreateProcessManagerService()
	if process == nil {
		t.Error("ProcessManager service should not be nil")
	}

	// Test browser detector service
	browser := factory.CreateBrowserDetectorService()
	if browser == nil {
		t.Error("BrowserDetector service should not be nil")
	}

	// Test app lifecycle service
	app := factory.CreateAppLifecycleService()
	if app == nil {
		t.Error("AppLifecycle service should not be nil")
	}
}

func TestNetworkService(t *testing.T) {
	network := NewNetworkService()

	// Test GetHTTPClient
	client, err := network.GetHTTPClient()
	if err != nil {
		t.Errorf("GetHTTPClient failed: %v", err)
	}
	if client == nil {
		t.Error("HTTP client should not be nil")
	}

	// We're not testing connectivity here as it depends on external factors
}

func TestFileSystemService(t *testing.T) {
	fs := NewFileSystemService("TestApp")

	// Test EnsureDataDirectory
	dataDir, err := fs.EnsureDataDirectory()
	if err != nil {
		t.Errorf("EnsureDataDirectory failed: %v", err)
	}
	if dataDir == "" {
		t.Error("Data directory should not be empty")
	}
}

// Additional tests for other services would be implemented here
// They are omitted for brevity and to avoid external dependencies
