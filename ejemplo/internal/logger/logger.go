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
	logDir   string
	debug    bool
)

// Init inicializa el logger con la carpeta de logs
func Init(logDirPath string, debugMode bool) error {
	logDir = logDirPath
	debug = debugMode

	// Crear directorio de logs si no existe
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de logs: %w", err)
	}

	// Archivo de errores
	errorFile, err := os.OpenFile(
		filepath.Join(logDir, "errors.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error abriendo archivo de errores: %w", err)
	}

	// Archivo de warnings
	warnFile, err := os.OpenFile(
		filepath.Join(logDir, "warnings.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error abriendo archivo de warnings: %w", err)
	}

	// Archivo de debug
	debugFile, err := os.OpenFile(
		filepath.Join(logDir, "debug.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error abriendo archivo de debug: %w", err)
	}

	// Crear loggers
	errorLog = log.New(errorFile, "", log.LstdFlags)
	warnLog = log.New(warnFile, "", log.LstdFlags)
	debugLog = log.New(debugFile, "", log.LstdFlags)

	return nil
}

// Error registra un error
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[ERROR] %s - %s", timestamp, msg)
	errorLog.Println(logMsg)
	if debug {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", msg)
	}
}

// Warn registra un warning
func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[WARN] %s - %s", timestamp, msg)
	warnLog.Println(logMsg)
	if debug {
		fmt.Fprintf(os.Stderr, "WARN: %s\n", msg)
	}
}

// Debug registra un mensaje de debug
func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[DEBUG] %s - %s", timestamp, msg)
	debugLog.Println(logMsg)
	if debug {
		fmt.Fprintf(os.Stdout, "DEBUG: %s\n", msg)
	}
}

// Info registra un mensaje informativo (solo en debug)
func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[INFO] %s - %s", timestamp, msg)
	debugLog.Println(logMsg)
	if debug {
		fmt.Fprintf(os.Stdout, "INFO: %s\n", msg)
	}
}



