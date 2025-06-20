package organizer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-file-organizer/internal/utils"
)

func TestGetDefaultExtensionCategories(t *testing.T) {
	categories := GetDefaultExtensionCategories()

	// Test that we get expected categories
	assert.Equal(t, "Images", categories[".jpg"])
	assert.Equal(t, "Images", categories[".png"])
	assert.Equal(t, "Documents", categories[".pdf"])
	assert.Equal(t, "Code", categories[".go"])
	assert.Equal(t, "Audio", categories[".mp3"])
	assert.Equal(t, "Video", categories[".mp4"])

	// Test that the map is a copy (not reference)
	categories[".test"] = "TestCategory"
	originalCategories := GetDefaultExtensionCategories()
	_, exists := originalCategories[".test"]
	assert.False(t, exists, "Default categories should not be modified")
}

func TestScanFiles(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "go-file-organizer-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{
		"document.pdf",
		"image.jpg",
		"script.py",
		"music.mp3",
		"unknown.xyz",
		"no-extension",
		"README.md",
	}

	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		file, err := os.Create(filePath)
		assert.NoError(t, err)
		file.Close()
	}

	// Test scanning
	categories, err := ScanFiles(tempDir)
	assert.NoError(t, err)

	// Verify categorization
	assert.Contains(t, categories["Documents"], filepath.Join(tempDir, "document.pdf"))
	assert.Contains(t, categories["Images"], filepath.Join(tempDir, "image.jpg"))
	assert.Contains(t, categories["Code"], filepath.Join(tempDir, "script.py"))
	assert.Contains(t, categories["Audio"], filepath.Join(tempDir, "music.mp3"))
	assert.Contains(t, categories["Unknown"], filepath.Join(tempDir, "unknown.xyz"))
	assert.Contains(t, categories["No Extension"], filepath.Join(tempDir, "no-extension"))

	// Note: .md files are not in default mappings, so they go to "Unknown"
	assert.Contains(t, categories["Unknown"], filepath.Join(tempDir, "README.md"))

	// Test file counts
	assert.Len(t, categories["Documents"], 1) // Only PDF
	assert.Len(t, categories["Images"], 1)
	assert.Len(t, categories["Code"], 1)
	assert.Len(t, categories["Audio"], 1)
	assert.Len(t, categories["Unknown"], 2) // README.md and unknown.xyz
	assert.Len(t, categories["No Extension"], 1)
}

func TestScanFilesWithConfig(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "go-file-organizer-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{
		"test.md",
		"backup.bak",
		"config.env",
	}

	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		file, err := os.Create(filePath)
		assert.NoError(t, err)
		file.Close()
	}

	// Create custom extension mapping
	extensionMapping := utils.NewExtensionMapping(GetDefaultExtensionCategories())
	customMappings := []string{
		".md=Notes",
		".bak=Backups",
		".env=Configuration",
	}
	err = extensionMapping.ApplyCLIMappings(customMappings)
	assert.NoError(t, err)

	// Test scanning with custom config
	categories, err := ScanFilesWithConfig(tempDir, extensionMapping, nil)
	assert.NoError(t, err)

	// Verify custom categorization
	assert.Contains(t, categories["Notes"], filepath.Join(tempDir, "test.md"))
	assert.Contains(t, categories["Backups"], filepath.Join(tempDir, "backup.bak"))
	assert.Contains(t, categories["Configuration"], filepath.Join(tempDir, "config.env"))
}

func TestScanFilesWithIgnore(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "go-file-organizer-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{
		"document.pdf",
		"image.jpg",
		"temp.tmp",
		".hidden",
		"log.log",
	}

	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		file, err := os.Create(filePath)
		assert.NoError(t, err)
		file.Close()
	}

	// Create ignore manager with test patterns
	ignoreManager := utils.NewIgnoreManager(tempDir)

	// Create temporary ignore file
	ignoreFile := filepath.Join(tempDir, ".testignore")
	ignoreContent := "*.tmp\n.hidden\n*.log\n"
	err = os.WriteFile(ignoreFile, []byte(ignoreContent), 0644)
	assert.NoError(t, err)

	err = ignoreManager.LoadIgnoreFile(ignoreFile)
	assert.NoError(t, err)

	// Test scanning with ignore patterns
	categories, err := ScanFilesWithConfig(tempDir, nil, ignoreManager)
	assert.NoError(t, err)

	// Verify ignored files are not in results
	for _, fileList := range categories {
		for _, filePath := range fileList {
			filename := filepath.Base(filePath)
			assert.NotEqual(t, "temp.tmp", filename)
			assert.NotEqual(t, ".hidden", filename)
			assert.NotEqual(t, "log.log", filename)
		}
	}

	// Verify non-ignored files are present
	assert.Contains(t, categories["Documents"], filepath.Join(tempDir, "document.pdf"))
	assert.Contains(t, categories["Images"], filepath.Join(tempDir, "image.jpg"))
}

func TestScanFilesNonExistentDirectory(t *testing.T) {
	_, err := ScanFiles("/nonexistent/directory")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "directory does not exist")
}

func TestShouldSkipCategory(t *testing.T) {
	// Test categories that should be skipped
	assert.True(t, shouldSkipCategory("Unknown"))
	assert.True(t, shouldSkipCategory("No Extension"))

	// Test categories that should not be skipped
	assert.False(t, shouldSkipCategory("Documents"))
	assert.False(t, shouldSkipCategory("Images"))
	assert.False(t, shouldSkipCategory("Code"))
}

func TestOrganizeFilesDryRun(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "go-file-organizer-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{
		"document.pdf",
		"image.jpg",
		"script.py",
	}

	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		file, err := os.Create(filePath)
		assert.NoError(t, err)
		file.Close()
	}

	// Create a temporary log file in a separate directory to avoid scanning it
	logDir, err := os.MkdirTemp("", "log-test")
	assert.NoError(t, err)
	defer os.RemoveAll(logDir)

	logFile := filepath.Join(logDir, "test.log")
	logger, err := utils.NewLogger(logFile)
	assert.NoError(t, err)
	defer logger.Close()

	// Run dry-run organization
	summary, err := OrganizeFiles(tempDir, true, logger)
	assert.NoError(t, err)

	// Verify summary
	assert.Equal(t, 3, summary.FilesScanned)
	assert.Equal(t, 3, summary.FilesMoved) // In dry-run, this counts "would be moved"
	assert.Equal(t, 3, summary.FoldersCreated)
	assert.Equal(t, 0, summary.FilesSkipped)

	// Verify files were NOT actually moved (dry-run)
	assert.FileExists(t, filepath.Join(tempDir, "document.pdf"))
	assert.FileExists(t, filepath.Join(tempDir, "image.jpg"))
	assert.FileExists(t, filepath.Join(tempDir, "script.py"))

	// Verify category folders were NOT created in dry-run
	assert.NoDirExists(t, filepath.Join(tempDir, "Documents"))
	assert.NoDirExists(t, filepath.Join(tempDir, "Images"))
	assert.NoDirExists(t, filepath.Join(tempDir, "Code"))
}
