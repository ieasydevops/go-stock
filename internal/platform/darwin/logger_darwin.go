//go:build darwin
// +build darwin

package darwin

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.LoggerService = (*loggerService)(nil)

type loggerService struct {
	logger  *zap.Logger
	sugar   *zap.SugaredLogger
	appName string
}

// NewLoggerService creates a new logger service for macOS.
// AppName is used to determine the log file path within standard macOS log directories.
func NewLoggerService(appName string) interfaces.LoggerService {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get user home directory: %v\n", err)
		// Fallback or panic, for now, we'll use a local log file if home dir is not found.
		homeDir = "."
	}

	// Standard macOS log location: ~/Library/Logs/AppName/
	logDirPath := filepath.Join(homeDir, "Library", "Logs", appName)
	if err := os.MkdirAll(logDirPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log directory %s: %v\n", logDirPath, err)
		// Fallback to current directory if creation fails
		logDirPath = "."
	}
	logFilePath := filepath.Join(logDirPath, appName+".log")

	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Compress:   true,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(lumberjackLogger),
		zap.InfoLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &loggerService{
		logger:  logger,
		sugar:   logger.Sugar(),
		appName: appName,
	}
}

// Init implements interfaces.LoggerService.
func (s *loggerService) Init(config map[string]interface{}) error {
	// The logger is already initialized in NewLoggerService.
	// This method is a no-op for the Darwin implementation.
	return nil
}

// Debug implements interfaces.LoggerService.
func (s *loggerService) Debug(args ...interface{}) {
	s.sugar.Debug(args...)
}

// Debugf implements interfaces.LoggerService.
func (s *loggerService) Debugf(format string, args ...interface{}) {
	s.sugar.Debugf(format, args...)
}

// Info implements interfaces.LoggerService.
func (s *loggerService) Info(args ...interface{}) {
	s.sugar.Info(args...)
}

// Infof implements interfaces.LoggerService.
func (s *loggerService) Infof(format string, args ...interface{}) {
	s.sugar.Infof(format, args...)
}

// Warn implements interfaces.LoggerService.
func (s *loggerService) Warn(args ...interface{}) {
	s.sugar.Warn(args...)
}

// Warnf implements interfaces.LoggerService.
func (s *loggerService) Warnf(format string, args ...interface{}) {
	s.sugar.Warnf(format, args...)
}

// Error implements interfaces.LoggerService.
func (s *loggerService) Error(args ...interface{}) {
	s.sugar.Error(args...)
}

// Errorf implements interfaces.LoggerService.
func (s *loggerService) Errorf(format string, args ...interface{}) {
	s.sugar.Errorf(format, args...)
}

// Fatal implements interfaces.LoggerService.
func (s *loggerService) Fatal(args ...interface{}) {
	s.sugar.Fatal(args...)
}

// Fatalf implements interfaces.LoggerService.
func (s *loggerService) Fatalf(format string, args ...interface{}) {
	s.sugar.Fatalf(format, args...)
}

// GetLogFilePath returns the path to the log file.
func (s *loggerService) GetLogFilePath() string {
	homeDir, _ := os.UserHomeDir()
	logDirPath := filepath.Join(homeDir, "Library", "Logs", s.appName)
	return filepath.Join(logDirPath, s.appName+".log")
}
