package pkgs

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// SetUpLogger sets up logging to a file and stdout.
// It returns the log file for deferred closing and an error if any occurs.
func SetUpLogger(logFileName string) (*os.File, error) {
	// Ensure the logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Get the current date and time for the log file prefix
	currentTime := time.Now().Format("20060102_150405")
	logFileName = fmt.Sprintf("%s_%s", currentTime, logFileName)

	// Open the log file
	logFilePath := "logs/" + logFileName
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Error opening log file %s, falling back to stdout only: %v", logFilePath, err)
		logFile = nil // Continue without a log file
	} else {
		log.Printf("Logging initialized. Log file: %s", logFilePath)
	}

	// Set up multi-writer to write to stdout and file if possible
	var multiWriter io.Writer
	if logFile != nil {
		multiWriter = io.MultiWriter(os.Stdout, logFile)
	} else {
		multiWriter = os.Stdout
	}

	// Configure log format to include timestamp, log level, and file location
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	return logFile, nil
}
