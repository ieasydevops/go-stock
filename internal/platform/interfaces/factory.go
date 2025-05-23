package interfaces

// PlatformServiceFactory defines the contract for a factory that creates
// platform-specific service instances.
type PlatformServiceFactory interface {
	CreateLogger() LoggerService
	CreateFileSystemService() FileSystemService
	CreateNetworkService() NetworkService
	CreateNotificationService() NotificationService
	CreateDialogService() DialogService
	CreateSystemTrayService() SystemTrayService
	CreateScreenInfoService() ScreenInfoService
	CreateProcessManagerService() ProcessManagerService
	CreateBrowserDetectorService() BrowserDetectorService
	CreateAppLifecycleService() AppLifecycleService
}
