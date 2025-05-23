package platform

import (
	"runtime"
	"sync"

	"go-stock/internal/platform/basic"
	"go-stock/internal/platform/darwin"
	"go-stock/internal/platform/interfaces"
	"go-stock/internal/platform/linux"
	"go-stock/internal/platform/windows"
)

var (
	currentFactory interfaces.PlatformServiceFactory
	once           sync.Once
)

// GetFactory returns the appropriate PlatformServiceFactory for the current operating system.
// It initializes the factory on the first call.
func GetFactory() interfaces.PlatformServiceFactory {
	once.Do(func() {
		switch runtime.GOOS {
		case "windows":
			currentFactory = windows.NewServiceFactory()
		case "darwin":
			currentFactory = darwin.NewServiceFactory()
		case "linux":
			currentFactory = linux.NewServiceFactory()
		default:
			currentFactory = basic.NewServiceFactory()
		}
	})
	return currentFactory
}

// SetFactory (primarily for testing) allows overriding the global factory.
// This should be used with caution in production code.
func SetFactory(factory interfaces.PlatformServiceFactory) {
	once.Do(func() { // Ensure the once block is registered even if not run with default logic
		if currentFactory == nil && factory != nil { // Only if not already initialized by GetFactory's switch
			// This condition helps if SetFactory is called before GetFactory
		}
	})
	currentFactory = factory
}
