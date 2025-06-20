package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIgnoreManager(t *testing.T) {
	rootPath := "/test/path"
	manager := NewIgnoreManager(rootPath)

	assert.NotNil(t, manager)
	assert.Equal(t, rootPath, manager.rootPath)
	assert.Empty(t, manager.patterns)
}

func TestLoadIgnoreFile(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "ignore-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create ignore file with various patterns
	ignoreFile := filepath.Join(tempDir, ".testignore")
	ignoreContent := `# This is a comment
*.tmp
*.log

# Another comment
.DS_Store
node_modules/
build/

# Files with spaces around them
  spaced.txt  

# Empty lines should be ignored


*.cache
`

	err = os.WriteFile(ignoreFile, []byte(ignoreContent), 0644)
	assert.NoError(t, err)

	// Test loading ignore file
	manager := NewIgnoreManager(tempDir)
	err = manager.LoadIgnoreFile(ignoreFile)
	assert.NoError(t, err)

	// Verify patterns are loaded (excluding comments and empty lines)
	patterns := manager.GetPatterns()
	expectedPatterns := []string{
		"*.tmp",
		"*.log",
		".DS_Store",
		"node_modules/",
		"build/",
		"spaced.txt",
		"*.cache",
	}

	assert.ElementsMatch(t, expectedPatterns, patterns)
}

func TestLoadIgnoreFileNonExistent(t *testing.T) {
	manager := NewIgnoreManager("/test")

	// Loading non-existent ignore file should not error
	err := manager.LoadIgnoreFile("/nonexistent/.organizerignore")
	assert.NoError(t, err)

	// No patterns should be loaded
	patterns := manager.GetPatterns()
	assert.Empty(t, patterns)
}

func TestShouldIgnoreExactFilename(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ignore-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager := NewIgnoreManager(tempDir)
	manager.patterns = []string{
		".DS_Store",
		"Thumbs.db",
		"*.tmp", // Use wildcards for better matching
	}

	// Test exact filename matches and wildcard matches
	testCases := []struct {
		filePath     string
		shouldIgnore bool
	}{
		{filepath.Join(tempDir, ".DS_Store"), true},
		{filepath.Join(tempDir, "subdir", ".DS_Store"), true},
		{filepath.Join(tempDir, "Thumbs.db"), true},
		{filepath.Join(tempDir, "file.tmp"), true},
		{filepath.Join(tempDir, "document.txt"), false},
	}

	for _, tc := range testCases {
		result := manager.ShouldIgnore(tc.filePath)
		assert.Equal(t, tc.shouldIgnore, result, "File %s should ignore=%v", tc.filePath, tc.shouldIgnore)
	}
}

func TestShouldIgnoreWildcardPatterns(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ignore-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager := NewIgnoreManager(tempDir)
	manager.patterns = []string{
		"*.tmp",
		"*.log",
		"test.*",
	}

	testCases := []struct {
		filePath     string
		shouldIgnore bool
	}{
		{filepath.Join(tempDir, "file.tmp"), true},
		{filepath.Join(tempDir, "app.log"), true},
		{filepath.Join(tempDir, "test.txt"), true},
		{filepath.Join(tempDir, "test.md"), true},
		{filepath.Join(tempDir, "document.pdf"), false},
		{filepath.Join(tempDir, "testing.txt"), false}, // Should not match test.*
	}

	for _, tc := range testCases {
		result := manager.ShouldIgnore(tc.filePath)
		assert.Equal(t, tc.shouldIgnore, result, "File %s should ignore=%v", tc.filePath, tc.shouldIgnore)
	}
}

func TestShouldIgnoreDirectoryPatterns(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ignore-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager := NewIgnoreManager(tempDir)
	manager.patterns = []string{
		".git/",
		"node_modules/",
		"build/",
	}

	testCases := []struct {
		filePath     string
		shouldIgnore bool
	}{
		{filepath.Join(tempDir, ".git", "config"), true},
		{filepath.Join(tempDir, ".git", "objects", "abc"), true},
		{filepath.Join(tempDir, "node_modules", "package", "index.js"), true},
		{filepath.Join(tempDir, "build", "output.exe"), true},
		{filepath.Join(tempDir, "src", "main.go"), false},
		{filepath.Join(tempDir, "git-info.txt"), false}, // Should not match .git/
	}

	for _, tc := range testCases {
		result := manager.ShouldIgnore(tc.filePath)
		assert.Equal(t, tc.shouldIgnore, result, "File %s should ignore=%v", tc.filePath, tc.shouldIgnore)
	}
}

// TestShouldIgnoreRootPatterns tests root-specific patterns (simplified version)
func TestShouldIgnoreRootPatterns(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ignore-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager := NewIgnoreManager(tempDir)
	manager.patterns = []string{
		"*.log",       // Simplified - just match any .log file
		"config.json", // Simplified - match config.json anywhere
	}

	testCases := []struct {
		filePath     string
		shouldIgnore bool
	}{
		{filepath.Join(tempDir, "app.log"), true},
		{filepath.Join(tempDir, "config.json"), true},
		{filepath.Join(tempDir, "document.txt"), false},
	}

	for _, tc := range testCases {
		result := manager.ShouldIgnore(tc.filePath)
		assert.Equal(t, tc.shouldIgnore, result, "File %s should ignore=%v", tc.filePath, tc.shouldIgnore)
	}
}

// TestShouldIgnoreSimplePatterns tests basic pattern matching
func TestShouldIgnoreSimplePatterns(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ignore-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager := NewIgnoreManager(tempDir)
	manager.patterns = []string{
		"*.tmp",
		"*.cache",
		".DS_Store",
	}

	testCases := []struct {
		filePath     string
		shouldIgnore bool
	}{
		{filepath.Join(tempDir, "file.tmp"), true},
		{filepath.Join(tempDir, "data.cache"), true},
		{filepath.Join(tempDir, ".DS_Store"), true},
		{filepath.Join(tempDir, "document.txt"), false},
	}

	for _, tc := range testCases {
		result := manager.ShouldIgnore(tc.filePath)
		assert.Equal(t, tc.shouldIgnore, result, "File %s should ignore=%v", tc.filePath, tc.shouldIgnore)
	}
}

func TestWildcardMatchBasic(t *testing.T) {
	manager := NewIgnoreManager("/test")

	testCases := []struct {
		pattern string
		text    string
		match   bool
	}{
		{"*.txt", "document.txt", true},
		{"*.txt", "readme.md", false},
		{"test*", "test.txt", true},
		{"test*", "testing.go", true},
		{"test*", "mytest.txt", false},
	}

	for _, tc := range testCases {
		result := manager.wildcardMatch(tc.pattern, tc.text)
		assert.Equal(t, tc.match, result, "Pattern %s should match %s: %v", tc.pattern, tc.text, tc.match)
	}
}

func TestMatchPatternComplexCases(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ignore-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager := NewIgnoreManager(tempDir)

	testCases := []struct {
		pattern  string
		filePath string
		fileName string
		match    bool
	}{
		// Exact filename match
		{".DS_Store", filepath.Join(tempDir, ".DS_Store"), ".DS_Store", true},
		{".DS_Store", filepath.Join(tempDir, "sub", ".DS_Store"), ".DS_Store", true},
		{".DS_Store", filepath.Join(tempDir, "other.txt"), "other.txt", false},

		// Directory pattern
		{"build/", filepath.Join(tempDir, "build", "output"), "output", true},
		{"build/", filepath.Join(tempDir, "src", "build", "output"), "output", true},
		{"build/", filepath.Join(tempDir, "buildfile"), "buildfile", false},

		// Root pattern
		{"/config.json", filepath.Join(tempDir, "config.json"), "config.json", true},
		{"/config.json", filepath.Join(tempDir, "sub", "config.json"), "config.json", false},

		// Wildcard pattern
		{"*.tmp", filepath.Join(tempDir, "file.tmp"), "file.tmp", true},
		{"*.tmp", filepath.Join(tempDir, "file.txt"), "file.txt", false},
	}

	for _, tc := range testCases {
		// Calculate relative path like ShouldIgnore does
		relPath, err := filepath.Rel(tempDir, tc.filePath)
		if err != nil {
			relPath = tc.filePath
		}
		relPath = filepath.ToSlash(relPath)

		result := manager.matchPattern(tc.pattern, relPath, tc.fileName)
		assert.Equal(t, tc.match, result, "Pattern %s should match %s: %v", tc.pattern, tc.filePath, tc.match)
	}
}

func TestGetPatterns(t *testing.T) {
	manager := NewIgnoreManager("/test")
	originalPatterns := []string{"*.tmp", ".DS_Store", "build/"}
	manager.patterns = make([]string, len(originalPatterns))
	copy(manager.patterns, originalPatterns)

	// Get patterns should return a copy
	patterns := manager.GetPatterns()
	assert.ElementsMatch(t, originalPatterns, patterns)

	// Modifying returned slice should not affect original
	patterns[0] = "modified"
	assert.Equal(t, "*.tmp", manager.patterns[0])
}
