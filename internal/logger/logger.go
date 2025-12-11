package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	errorLog *log.Logger
	warnLog  *log.Logger
	debugLog *log.Logger
	infoLog  *log.Logger
	logDir   string
	debug    bool
)

// Init initializes the logger with the logs directory
func Init(logDirPath string, debugMode bool) error {
	logDir = logDirPath
	debug = debugMode

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de logs: %w", err)
	}

	// Error log file
	errorFile, err := os.OpenFile(
		filepath.Join(logDir, "errors.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error abriendo archivo de errores: %w", err)
	}

	// Warning log file
	warnFile, err := os.OpenFile(
		filepath.Join(logDir, "warnings.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error abriendo archivo de warnings: %w", err)
	}

	// Debug log file
	debugFile, err := os.OpenFile(
		filepath.Join(logDir, "debug.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error abriendo archivo de debug: %w", err)
	}

	// Create loggers
	errorLog = log.New(errorFile, "", log.LstdFlags)
	warnLog = log.New(warnFile, "", log.LstdFlags)
	debugLog = log.New(debugFile, "", log.LstdFlags)
	infoLog = log.New(debugFile, "", log.LstdFlags)

	return nil
}

// SetDebug enables or disables debug mode
func SetDebug(debugMode bool) {
	debug = debugMode
}

// IsDebug returns whether debug mode is enabled
func IsDebug() bool {
	return debug
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[ERROR] %s - %s", timestamp, msg)
	if errorLog != nil {
		errorLog.Println(logMsg)
	}
	if debug {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[0m\n", msg)
	}
}

// Warn logs a warning message
func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[WARN] %s - %s", timestamp, msg)
	if warnLog != nil {
		warnLog.Println(logMsg)
	}
	if debug {
		fmt.Fprintf(os.Stderr, "\033[33mWARN: %s\033[0m\n", msg)
	}
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[DEBUG] %s - %s", timestamp, msg)
	if debugLog != nil {
		debugLog.Println(logMsg)
	}
	if debug {
		fmt.Fprintf(os.Stdout, "\033[36mDEBUG: %s\033[0m\n", msg)
	}
}

// Info logs an informational message (only in debug mode to console)
func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[INFO] %s - %s", timestamp, msg)
	if infoLog != nil {
		infoLog.Println(logMsg)
	}
	if debug {
		fmt.Fprintf(os.Stdout, "\033[34mINFO: %s\033[0m\n", msg)
	}
}

// Success logs a success message (always to file, console only in debug)
func Success(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[SUCCESS] %s - %s", timestamp, msg)
	if infoLog != nil {
		infoLog.Println(logMsg)
	}
	if debug {
		fmt.Fprintf(os.Stdout, "\033[32mSUCCESS: %s\033[0m\n", msg)
	}
}

