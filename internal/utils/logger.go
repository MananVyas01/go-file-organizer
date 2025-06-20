package utils

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	file   *os.File
	logger *log.Logger
}

// NewLogger creates a new logger that writes to organizer.log
func NewLogger(logPath string) (*Logger, error) {
	// Open log file in append mode, create if doesn't exist
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	logger := log.New(file, "", log.LstdFlags)
	
	// Write session separator
	logger.Println("=== New Organizer Session Started ===")
	
	return &Logger{
		file:   file,
		logger: logger,
	}, nil
}

// LogDryRun logs a dry-run action
func (l *Logger) LogDryRun(source, destination string) {
	l.logger.Printf("[DRY-RUN] Would move: %s -> %s", source, destination)
}

// LogMove logs an actual file move
func (l *Logger) LogMove(source, destination string) {
	l.logger.Printf("[MOVE] Moved: %s -> %s", source, destination)
}

// LogFolderCreation logs folder creation
func (l *Logger) LogFolderCreation(folderPath string, isDryRun bool) {
	if isDryRun {
		l.logger.Printf("[DRY-RUN] Would create folder: %s", folderPath)
	} else {
		l.logger.Printf("[FOLDER] Created folder: %s", folderPath)
	}
}

// LogError logs an error
func (l *Logger) LogError(operation, filePath string, err error) {
	l.logger.Printf("[ERROR] %s failed for %s: %v", operation, filePath, err)
}

// LogSummary logs the final summary statistics
func (l *Logger) LogSummary(stats Summary) {
	l.logger.Printf("[SUMMARY] Files scanned: %d, moved: %d, folders created: %d, skipped: %d",
		stats.FilesScanned, stats.FilesMoved, stats.FoldersCreated, stats.FilesSkipped)
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.file != nil {
		l.logger.Println("=== Session Ended ===")
		return l.file.Close()
	}
	return nil
}

// Summary holds statistics about the organization process
type Summary struct {
	FilesScanned   int
	FilesMoved     int
	FoldersCreated int
	FilesSkipped   int
}
