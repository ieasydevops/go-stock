//go:build darwin
// +build darwin

package darwin

import (
	"go-stock/internal/platform/interfaces"
)

// ensure darwinServiceFactory implements the interface.
var _ interfaces.PlatformServiceFactory = (*darwinServiceFactory)(nil)

type darwinServiceFactory struct {
	appName    string
	appVersion string
}

// NewServiceFactory creates a new factory for Darwin (macOS) platform services.
func NewServiceFactory(appName, appVersion string) interfaces.PlatformServiceFactory {
	return &darwinServiceFactory{
		appName:    appName,
		appVersion: appVersion,
	}
}

func (f *darwinServiceFactory) CreateLogger() interfaces.LoggerService {
	return NewLoggerService(f.appName)
}

func (f *darwinServiceFactory) CreateFileSystemService() interfaces.FileSystemService {
	return NewFileSystemService(f.appName)
}

func (f *darwinServiceFactory) CreateNetworkService() interfaces.NetworkService {
	return NewNetworkService()
}

func (f *darwinServiceFactory) CreateNotificationService() interfaces.NotificationService {
	return NewNotificationService(f.appName)
}

func (f *darwinServiceFactory) CreateDialogService() interfaces.DialogService {
	return NewDialogService()
}

func (f *darwinServiceFactory) CreateSystemTrayService() interfaces.SystemTrayService {
	return NewSystemTrayService()
}

func (f *darwinServiceFactory) CreateScreenInfoService() interfaces.ScreenInfoService {
	return NewScreenInfoService()
}

func (f *darwinServiceFactory) CreateProcessManagerService() interfaces.ProcessManagerService {
	return NewProcessManagerService()
}

func (f *darwinServiceFactory) CreateBrowserDetectorService() interfaces.BrowserDetectorService {
	return NewBrowserDetectorService()
}

func (f *darwinServiceFactory) CreateAppLifecycleService() interfaces.AppLifecycleService {
	return NewAppLifecycleService(f.appName, f.appVersion)
}
