//go:build windows
// +build windows

package windows

import (
	"fmt"
	"log"
	os "os"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.LoggerService = (*loggerService)(nil)

type loggerService struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	// Potentially add file handles if logging to files
}

func NewLoggerService() interfaces.LoggerService {
	// TODO: Enhance with file logging, rotation, proper output streams based on config
	// For now, basic stdout/stderr logging
	l := &loggerService{
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stderr, "WARN:  ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	return l
}

func (l *loggerService) Init(config map[string]interface{}) error {
	// TODO: Parse config for log file paths, levels, rotation settings for Windows
	log.Println("Windows LoggerService Init called with config:", config)
	// Example: Get log directory from FileSystemService if available and configured
	return nil
}

func (l *loggerService) Debug(args ...interface{}) {
	l.debugLogger.Output(2, fmt.Sprintln(args...))
}

func (l *loggerService) Info(args ...interface{}) {
	l.infoLogger.Output(2, fmt.Sprintln(args...))
}

func (l *loggerService) Warn(args ...interface{}) {
	l.warnLogger.Output(2, fmt.Sprintln(args...))
}

func (l *loggerService) Error(args ...interface{}) {
	l.errorLogger.Output(2, fmt.Sprintln(args...))
}

func (l *loggerService) Fatal(args ...interface{}) {
	l.fatalLogger.Output(2, fmt.Sprintln(args...))
	os.Exit(1)
}

func (l *loggerService) Debugf(format string, args ...interface{}) {
	l.debugLogger.Output(2, fmt.Sprintf(format, args...))
}

func (l *loggerService) Infof(format string, args ...interface{}) {
	l.infoLogger.Output(2, fmt.Sprintf(format, args...))
}

func (l *loggerService) Warnf(format string, args ...interface{}) {
	l.warnLogger.Output(2, fmt.Sprintf(format, args...))
}

func (l *loggerService) Errorf(format string, args ...interface{}) {
	l.errorLogger.Output(2, fmt.Sprintf(format, args...))
}

func (l *loggerService) Fatalf(format string, args ...interface{}) {
	l.fatalLogger.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}
