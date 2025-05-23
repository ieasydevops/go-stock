package interfaces

// AppLifecycleService defines the contract for managing core application lifecycle events
// and platform-specific application behaviors.
type AppLifecycleService interface {
	// Initialize performs platform-specific setup when the application starts.
	// This could include setting up single instance locks, environment checks, etc.
	Initialize() error

	// Shutdown performs platform-specific cleanup before the application exits.
	Shutdown() error

	// HandleSecondInstance is called when another instance of the application is launched.
	// It should bring the existing instance to the foreground and potentially process `data`
	// (e.g., command line arguments from the new instance).
	// Returns true if this instance should handle it, false if the new instance should proceed (or error).
	HandleSecondInstance(data interface{}) (bool, error)

	// GetAppVersion returns the application's version string.
	GetAppVersion() (string, error)

	// GetAppName returns the application's name.
	GetAppName() (string, error)

	// RequestExit requests the application to close.
	// This might involve asking the user for confirmation on unsaved work.
	RequestExit()

	// RestartApplication restarts the current application.
	RestartApplication() error
}
