package domain

import (
	"fmt"
	"log"
	"os"
)

type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
	}
}

func (l *SimpleLogger) Info(format string, args ...interface{}) {
	if err := l.infoLogger.Output(2, fmt.Sprintf(format, args...)); err != nil {
		// We can't use the logger itself to log this error as it would cause a recursive call
		fmt.Fprintf(os.Stderr, "Failed to log info message: %v\n", err)
	}
}

func (l *SimpleLogger) Error(format string, args ...interface{}) {
	if err := l.errorLogger.Output(2, fmt.Sprintf(format, args...)); err != nil {
		// We can't use the logger itself to log this error as it would cause a recursive call
		fmt.Fprintf(os.Stderr, "Failed to log error message: %v\n", err)
	}
}

func (l *SimpleLogger) Debug(format string, args ...interface{}) {
	if err := l.debugLogger.Output(2, fmt.Sprintf(format, args...)); err != nil {
		// We can't use the logger itself to log this error as it would cause a recursive call
		fmt.Fprintf(os.Stderr, "Failed to log debug message: %v\n", err)
	}
}
