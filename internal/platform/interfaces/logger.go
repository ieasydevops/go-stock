package interfaces

// LoggerService defines the contract for logging functionalities.
// It supports different log levels.
// TODO: Define specific log levels and methods (e.g., Debugf, Infof, Warnf, Errorf, Fatalf)
type LoggerService interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	// Init allows for platform-specific logger initialization if needed.
	Init(config map[string]interface{}) error
}
