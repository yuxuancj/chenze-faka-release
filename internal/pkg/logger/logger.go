package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
)

func Init(logDir string) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(logDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	infoLogger = log.New(mw, "[INFO] ", log.LstdFlags|log.Lshortfile)
	errorLogger = log.New(mw, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	debugLogger = log.New(mw, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
	return nil
}

func Infof(format string, v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(format, v...)
	} else {
		log.Printf("[INFO] "+format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Printf(format, v...)
	} else {
		log.Printf("[ERROR] "+format, v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if debugLogger != nil {
		debugLogger.Printf(format, v...)
	} else {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func Warnf(format string, v ...interface{}) {
	Infof("[WARN] "+format, v...)
}

func Fatalf(format string, v ...interface{}) {
	Errorf(format, v...)
	os.Exit(1)
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
