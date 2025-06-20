package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "logger-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logPath := filepath.Join(tempDir, "test.log")

	// Create logger
	logger, err := NewLogger(logPath)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// Clean up
	err = logger.Close()
	assert.NoError(t, err)

	// Verify log file was created
	assert.FileExists(t, logPath)

	// Verify session header was written
	content, err := os.ReadFile(logPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "=== New Organizer Session Started ===")
}

func TestLoggerOperations(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "logger-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logPath := filepath.Join(tempDir, "test.log")

	// Create logger
	logger, err := NewLogger(logPath)
	assert.NoError(t, err)
	defer logger.Close()

	// Test different log operations
	logger.LogDryRun("source.txt", "dest/source.txt")
	logger.LogMove("file.pdf", "Documents/file.pdf")
	logger.LogFolderCreation("Documents/", false)
	logger.LogFolderCreation("Images/", true)
	logger.LogError("Move", "broken.txt", assert.AnError)

	summary := Summary{
		FilesScanned:   10,
		FilesMoved:     8,
		FoldersCreated: 3,
		FilesSkipped:   2,
	}
	logger.LogSummary(summary)

	// Close to flush
	err = logger.Close()
	assert.NoError(t, err)

	// Read and verify log content
	content, err := os.ReadFile(logPath)
	assert.NoError(t, err)
	logContent := string(content)

	// Verify all operations were logged
	assert.Contains(t, logContent, "[DRY-RUN] Would move: source.txt -> dest/source.txt")
	assert.Contains(t, logContent, "[MOVE] Moved: file.pdf -> Documents/file.pdf")
	assert.Contains(t, logContent, "[FOLDER] Created folder: Documents/")
	assert.Contains(t, logContent, "[DRY-RUN] Would create folder: Images/")
	assert.Contains(t, logContent, "[ERROR] Move failed for broken.txt:")
	assert.Contains(t, logContent, "[SUMMARY] Files scanned: 10, moved: 8, folders created: 3, skipped: 2")
	assert.Contains(t, logContent, "=== Session Ended ===")
}

func TestLoggerAppendMode(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "logger-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logPath := filepath.Join(tempDir, "test.log")

	// Create first logger session
	logger1, err := NewLogger(logPath)
	assert.NoError(t, err)
	logger1.LogMove("file1.txt", "Documents/file1.txt")
	err = logger1.Close()
	assert.NoError(t, err)

	// Create second logger session
	logger2, err := NewLogger(logPath)
	assert.NoError(t, err)
	logger2.LogMove("file2.txt", "Documents/file2.txt")
	err = logger2.Close()
	assert.NoError(t, err)

	// Read log content
	content, err := os.ReadFile(logPath)
	assert.NoError(t, err)
	logContent := string(content)

	// Verify both sessions are in the log
	sessionCount := strings.Count(logContent, "=== New Organizer Session Started ===")
	assert.Equal(t, 2, sessionCount)

	endSessionCount := strings.Count(logContent, "=== Session Ended ===")
	assert.Equal(t, 2, endSessionCount)

	// Verify both file moves are logged
	assert.Contains(t, logContent, "file1.txt -> Documents/file1.txt")
	assert.Contains(t, logContent, "file2.txt -> Documents/file2.txt")
}

func TestSummaryStruct(t *testing.T) {
	summary := Summary{
		FilesScanned:   15,
		FilesMoved:     12,
		FoldersCreated: 4,
		FilesSkipped:   3,
	}

	assert.Equal(t, 15, summary.FilesScanned)
	assert.Equal(t, 12, summary.FilesMoved)
	assert.Equal(t, 4, summary.FoldersCreated)
	assert.Equal(t, 3, summary.FilesSkipped)
}

func TestLoggerNilFile(t *testing.T) {
	logger := &Logger{file: nil}

	// Closing nil file should not panic
	err := logger.Close()
	assert.NoError(t, err)
}
