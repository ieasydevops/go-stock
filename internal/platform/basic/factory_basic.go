package basic

import (
	"log"

	"go-stock/internal/platform/interfaces"
)

// ensure basicServiceFactory implements the interface.
var _ interfaces.PlatformServiceFactory = (*basicServiceFactory)(nil)

type basicServiceFactory struct{}

// NewServiceFactory creates a new factory for basic/default platform services.
func NewServiceFactory() interfaces.PlatformServiceFactory {
	log.Println("Warning: Using basic platform service factory. Platform-specific features may be limited or unavailable.")
	return &basicServiceFactory{}
}

// Implement PlatformServiceFactory methods (will return nil or very basic implementations)

type basicLogger struct{}

func (l *basicLogger) Debug(args ...interface{}) {
	log.Println(append([]interface{}{"[DEBUG]"}, args...)...)
}
func (l *basicLogger) Info(args ...interface{}) {
	log.Println(append([]interface{}{"[INFO]"}, args...)...)
}
func (l *basicLogger) Warn(args ...interface{}) {
	log.Println(append([]interface{}{"[WARN]"}, args...)...)
}
func (l *basicLogger) Error(args ...interface{}) {
	log.Println(append([]interface{}{"[ERROR]"}, args...)...)
}
func (l *basicLogger) Fatal(args ...interface{}) {
	log.Fatal(append([]interface{}{"[FATAL]"}, args...)...)
}
func (l *basicLogger) Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}
func (l *basicLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}
func (l *basicLogger) Warnf(format string, args ...interface{}) {
	log.Printf("[WARN] "+format, args...)
}
func (l *basicLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}
func (l *basicLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf("[FATAL] "+format, args...)
}
func (l *basicLogger) Init(config map[string]interface{}) error { return nil }

func (f *basicServiceFactory) CreateLogger() interfaces.LoggerService {
	return &basicLogger{}
}

func (f *basicServiceFactory) CreateFileSystemService() interfaces.FileSystemService {
	log.Println("Warning: Basic FileSystemService returning nil. Filesystem operations will fail.")
	return nil // Or a very basic in-memory version if essential for non-platform-specific logic
}

func (f *basicServiceFactory) CreateNetworkService() interfaces.NetworkService {
	log.Println("Warning: Basic NetworkService returning nil. Network operations will fail.")
	return nil
}

func (f *basicServiceFactory) CreateNotificationService() interfaces.NotificationService {
	log.Println("Warning: Basic NotificationService returning nil. Notifications will not be shown.")
	return nil
}

func (f *basicServiceFactory) CreateDialogService() interfaces.DialogService {
	log.Println("Warning: Basic DialogService returning nil. Dialogs will not be shown.")
	return nil
}

func (f *basicServiceFactory) CreateSystemTrayService() interfaces.SystemTrayService {
	log.Println("Warning: Basic SystemTrayService returning nil. System tray will not be available.")
	return nil
}

func (f *basicServiceFactory) CreateScreenInfoService() interfaces.ScreenInfoService {
	log.Println("Warning: Basic ScreenInfoService returning nil. Screen information will not be available.")
	return nil
}

func (f *basicServiceFactory) CreateProcessManagerService() interfaces.ProcessManagerService {
	log.Println("Warning: Basic ProcessManagerService returning nil. Process management will be limited.")
	return nil
}

func (f *basicServiceFactory) CreateBrowserDetectorService() interfaces.BrowserDetectorService {
	log.Println("Warning: Basic BrowserDetectorService returning nil. Browser detection will not be available.")
	return nil
}

func (f *basicServiceFactory) CreateAppLifecycleService() interfaces.AppLifecycleService {
	log.Println("Warning: Basic AppLifecycleService returning nil. App lifecycle management will be limited.")
	return nil
}
