package logger

import (
	"github/MahfujulSagor/movies_crud/internals/config"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func Init(cfg *config.Config) {
	logFilePath := cfg.LoggingConfig.File
	if logFilePath == "" {
		logFilePath = "logs/app.log"
	}

	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	logFile, err := os.OpenFile(filepath.Join("logs", "app.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	var writer io.Writer
	if cfg.Env == "development" {
		writer = io.MultiWriter(os.Stdout, logFile)
	} else {
		writer = logFile
	}

	flags := log.Ldate | log.Ltime | log.Lshortfile
	Info = log.New(writer, "INFO: ", flags)
	Error = log.New(writer, "ERROR: ", flags)
}
