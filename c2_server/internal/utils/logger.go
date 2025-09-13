package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	// DEBUG level for detailed information
	DEBUG LogLevel = iota
	// INFO level for general operational information
	INFO
	// WARNING level for concerning but non-critical issues
	WARNING
	// ERROR level for errors that should be addressed
	ERROR
	// FATAL level for critical errors that require immediate attention
	FATAL
)

var logLevelNames = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

type Logger struct {
	minLevel LogLevel
	logFile  *os.File
}

func NewLogger(minLevel LogLevel, logFilePath string) (*Logger, error) {
	var logFile *os.File
	var err error

	if logFilePath != "" {

		logDir := filepath.Dir(logFilePath)
		if mkdirErr := os.MkdirAll(logDir, 0755); mkdirErr != nil {
			return nil, fmt.Errorf("failed to create log directory: %v", err)
		}

		logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)

		}

		// Set log output to file
		log.SetOutput(logFile)
	}

	return &Logger{
		minLevel: minLevel,
		logFile:  logFile,
	}, nil
}

func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// Log logs a message with the specified level
func (l *Logger) Log(level LogLevel, format string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}

	// Format the log message
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logEntry := fmt.Sprintf("%s [%s] %s:%d - %s", timestamp, logLevelNames[level], filepath.Base(file), line, message)

	// Log the message
	log.Println(logEntry)

	// If fatal, exit the program
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.Log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.Log(INFO, format, args...)
}

// Warning logs a warning message
func (l *Logger) Warning(format string, args ...interface{}) {
	l.Log(WARNING, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.Log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.Log(FATAL, format, args...)
}
