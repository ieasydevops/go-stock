//go:build windows
// +build windows

package windows

import (
	"go-stock/internal/platform/interfaces"
)

// ensure windowsServiceFactory implements the interface.
var _ interfaces.PlatformServiceFactory = (*windowsServiceFactory)(nil)

type windowsServiceFactory struct{}

// NewServiceFactory creates a new factory for Windows platform services.
func NewServiceFactory() interfaces.PlatformServiceFactory {
	return &windowsServiceFactory{}
}

// Implement PlatformServiceFactory methods (will return nil for now)
func (f *windowsServiceFactory) CreateLogger() interfaces.LoggerService {
	// TODO: Implement in T003 or later task for Windows logger
	return NewLoggerService()
}

func (f *windowsServiceFactory) CreateFileSystemService() interfaces.FileSystemService {
	// TODO: Implement in T003 or later task for Windows filesystem
	return NewFileSystemService()
}

func (f *windowsServiceFactory) CreateNetworkService() interfaces.NetworkService {
	// TODO: Implement in T003 or later task for Windows network
	return NewNetworkService()
}

func (f *windowsServiceFactory) CreateNotificationService() interfaces.NotificationService {
	// TODO: Implement in T003 or T010 for Windows notifications
	return NewNotificationService()
}

func (f *windowsServiceFactory) CreateDialogService() interfaces.DialogService {
	// TODO: Implement in T003 or later task for Windows dialogs
	return NewDialogService()
}

func (f *windowsServiceFactory) CreateSystemTrayService() interfaces.SystemTrayService {
	// TODO: Implement in T003 or T011 for Windows system tray
	return NewSystemTrayService()
}

func (f *windowsServiceFactory) CreateScreenInfoService() interfaces.ScreenInfoService {
	// TODO: Implement in T003 or later task for Windows screen info
	return NewScreenInfoService()
}

func (f *windowsServiceFactory) CreateProcessManagerService() interfaces.ProcessManagerService {
	// TODO: Implement in T003 or later task for Windows process manager
	return NewProcessManagerService()
}

func (f *windowsServiceFactory) CreateBrowserDetectorService() interfaces.BrowserDetectorService {
	// TODO: Implement in T003 or T013 for Windows browser detector
	return NewBrowserDetectorService()
}

func (f *windowsServiceFactory) CreateAppLifecycleService() interfaces.AppLifecycleService {
	// TODO: Implement in T003 or later task for Windows app lifecycle
	// TODO: appName and appVersion should come from a reliable source (config, build flags)
	return NewAppLifecycleService("GoStockApp", "0.0.1")
}
