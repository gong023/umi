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
	l.infoLogger.Output(2, fmt.Sprintf(format, args...))
}

func (l *SimpleLogger) Error(format string, args ...interface{}) {
	l.errorLogger.Output(2, fmt.Sprintf(format, args...))
}

func (l *SimpleLogger) Debug(format string, args ...interface{}) {
	l.debugLogger.Output(2, fmt.Sprintf(format, args...))
}
