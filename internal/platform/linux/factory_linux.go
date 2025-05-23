//go:build linux
// +build linux

package linux

import (
	"go-stock/internal/platform/interfaces"
)

// ensure linuxServiceFactory implements the interface.
var _ interfaces.PlatformServiceFactory = (*linuxServiceFactory)(nil)

type linuxServiceFactory struct {
	appName string // Store app name for services that need it
}

// NewServiceFactory creates a new factory for Linux platform services.
func NewServiceFactory() interfaces.PlatformServiceFactory {
	return &linuxServiceFactory{
		appName: "go-stock", // Default app name
	}
}

// Implement PlatformServiceFactory methods (will return nil for now)
func (f *linuxServiceFactory) CreateLogger() interfaces.LoggerService {
	// TODO: Implement in T005 or later task for Linux logger
	return nil
}

func (f *linuxServiceFactory) CreateFileSystemService() interfaces.FileSystemService {
	// Return our new Linux filesystem service
	return NewFileSystemService(f.appName)
}

func (f *linuxServiceFactory) CreateNetworkService() interfaces.NetworkService {
	// TODO: Implement in T005 or later task for Linux network
	return nil
}

func (f *linuxServiceFactory) CreateNotificationService() interfaces.NotificationService {
	// Return our new Linux notification service
	return NewNotificationService(f.appName)
}

func (f *linuxServiceFactory) CreateDialogService() interfaces.DialogService {
	// TODO: Implement in T005 or later task for Linux dialogs
	return nil
}

func (f *linuxServiceFactory) CreateSystemTrayService() interfaces.SystemTrayService {
	// Return our new Linux system tray service
	return NewSystemTrayService()
}

func (f *linuxServiceFactory) CreateScreenInfoService() interfaces.ScreenInfoService {
	// TODO: Implement in T005 or later task for Linux screen info
	return nil
}

func (f *linuxServiceFactory) CreateProcessManagerService() interfaces.ProcessManagerService {
	// TODO: Implement in T005 or later task for Linux process manager
	return nil
}

func (f *linuxServiceFactory) CreateBrowserDetectorService() interfaces.BrowserDetectorService {
	// Return our new Linux browser detector service
	return NewBrowserDetectorService()
}

func (f *linuxServiceFactory) CreateAppLifecycleService() interfaces.AppLifecycleService {
	// TODO: Implement in T005 or later task for Linux app lifecycle
	return nil
}
