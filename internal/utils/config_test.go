package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExtensionMapping(t *testing.T) {
	defaultMappings := map[string]string{
		".txt": "Documents",
		".jpg": "Images",
		".mp3": "Audio",
	}
	
	mapping := NewExtensionMapping(defaultMappings)
	
	// Test default mappings are loaded
	category, exists := mapping.GetMapping(".txt")
	assert.True(t, exists)
	assert.Equal(t, "Documents", category)
	
	category, exists = mapping.GetMapping(".jpg")
	assert.True(t, exists)
	assert.Equal(t, "Images", category)
	
	// Test case insensitive lookup
	category, exists = mapping.GetMapping(".TXT")
	assert.True(t, exists)
	assert.Equal(t, "Documents", category)
}

func TestLoadConfig(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "config-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create valid config file
	configPath := filepath.Join(tempDir, "config.json")
	configContent := `{
		"customMappings": {
			".md": "Notes",
			".log": "Logs",
			".bak": "Backups"
		},
		"description": "Test config"
	}`
	
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)
	
	// Test loading config
	mapping := NewExtensionMapping(map[string]string{".txt": "Documents"})
	err = mapping.LoadConfig(configPath)
	assert.NoError(t, err)
	
	// Verify custom mappings are loaded
	category, exists := mapping.GetMapping(".md")
	assert.True(t, exists)
	assert.Equal(t, "Notes", category)
	
	category, exists = mapping.GetMapping(".log")
	assert.True(t, exists)
	assert.Equal(t, "Logs", category)
	
	// Verify default mappings still exist
	category, exists = mapping.GetMapping(".txt")
	assert.True(t, exists)
	assert.Equal(t, "Documents", category)
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "config-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create invalid config file
	configPath := filepath.Join(tempDir, "config.json")
	configContent := `{
		"customMappings": {
			".md": "Notes"
		// Invalid JSON - missing closing brace
	`
	
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)
	
	// Test loading invalid config
	mapping := NewExtensionMapping(map[string]string{".txt": "Documents"})
	err = mapping.LoadConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse config JSON")
}

func TestLoadConfigNonExistentFile(t *testing.T) {
	mapping := NewExtensionMapping(map[string]string{".txt": "Documents"})
	
	// Loading non-existent config should not error
	err := mapping.LoadConfig("/nonexistent/config.json")
	assert.NoError(t, err)
	
	// Default mappings should still be present
	category, exists := mapping.GetMapping(".txt")
	assert.True(t, exists)
	assert.Equal(t, "Documents", category)
}

func TestApplyCLIMappings(t *testing.T) {
	mapping := NewExtensionMapping(map[string]string{
		".txt": "Documents",
		".jpg": "Images",
	})
	
	// Test valid CLI mappings
	cliMappings := []string{
		".md=Notes",
		".log=Logs",
		".txt=TextFiles", // Override default
	}
	
	err := mapping.ApplyCLIMappings(cliMappings)
	assert.NoError(t, err)
	
	// Verify CLI mappings are applied
	category, exists := mapping.GetMapping(".md")
	assert.True(t, exists)
	assert.Equal(t, "Notes", category)
	
	category, exists = mapping.GetMapping(".log")
	assert.True(t, exists)
	assert.Equal(t, "Logs", category)
	
	// Verify override works
	category, exists = mapping.GetMapping(".txt")
	assert.True(t, exists)
	assert.Equal(t, "TextFiles", category)
}

func TestApplyCLIMappingsInvalidFormat(t *testing.T) {
	mapping := NewExtensionMapping(map[string]string{".txt": "Documents"})
	
	// Test invalid format
	invalidMappings := []string{
		"invalid-format",
		".md", // Missing =Category
	}
	
	for _, invalidMapping := range invalidMappings {
		err := mapping.ApplyCLIMappings([]string{invalidMapping})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid mapping format")
	}
	
	// Test empty extension (different error)
	err := mapping.ApplyCLIMappings([]string{"=Category"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid extension")
}

func TestValidateExtension(t *testing.T) {
	mapping := NewExtensionMapping(map[string]string{})
	
	// Test valid extensions
	validExtensions := []string{".txt", ".md", ".jpg", ".MP3"}
	for _, ext := range validExtensions {
		err := mapping.validateExtension(ext)
		assert.NoError(t, err, "Extension %s should be valid", ext)
	}
	
	// Test invalid extensions
	invalidExtensions := []string{"", ".", "txt", "md"}
	for _, ext := range invalidExtensions {
		err := mapping.validateExtension(ext)
		assert.Error(t, err, "Extension %s should be invalid", ext)
	}
}

func TestValidateCategory(t *testing.T) {
	mapping := NewExtensionMapping(map[string]string{})
	
	// Test valid categories
	validCategories := []string{"Documents", "Images", "Audio_Files", "Code-Files"}
	for _, category := range validCategories {
		err := mapping.validateCategory(category)
		assert.NoError(t, err, "Category %s should be valid", category)
	}
	
	// Test invalid categories
	invalidCategories := []string{
		"", // Empty
		" Documents ", // Leading/trailing spaces
		"Docs/Images", // Contains slash
		"Docs\\Images", // Contains backslash
		"Docs:Images", // Contains colon
		"Docs*Images", // Contains asterisk
		"Docs?Images", // Contains question mark
		"Docs\"Images", // Contains quote
		"Docs<Images", // Contains less than
		"Docs>Images", // Contains greater than
		"Docs|Images", // Contains pipe
	}
	
	for _, category := range invalidCategories {
		err := mapping.validateCategory(category)
		assert.Error(t, err, "Category %s should be invalid", category)
	}
}

func TestGetMappings(t *testing.T) {
	originalMappings := map[string]string{
		".txt": "Documents",
		".jpg": "Images",
	}
	
	mapping := NewExtensionMapping(originalMappings)
	
	// Get all mappings
	allMappings := mapping.GetMappings()
	
	// Verify we get all mappings
	assert.Equal(t, "Documents", allMappings[".txt"])
	assert.Equal(t, "Images", allMappings[".jpg"])
	
	// Verify it's a copy (not reference)
	allMappings[".test"] = "Test"
	_, exists := mapping.GetMapping(".test")
	assert.False(t, exists, "Original mapping should not be modified")
}

func TestPriorityOrder(t *testing.T) {
	// Test that CLI mappings override config mappings which override defaults
	tempDir, err := os.MkdirTemp("", "config-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create config file
	configPath := filepath.Join(tempDir, "config.json")
	configContent := `{
		"customMappings": {
			".txt": "ConfigDocuments"
		}
	}`
	
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)
	
	// Start with default mapping
	mapping := NewExtensionMapping(map[string]string{".txt": "DefaultDocuments"})
	
	// Load config (should override default)
	err = mapping.LoadConfig(configPath)
	assert.NoError(t, err)
	
	category, exists := mapping.GetMapping(".txt")
	assert.True(t, exists)
	assert.Equal(t, "ConfigDocuments", category)
	
	// Apply CLI mapping (should override config)
	err = mapping.ApplyCLIMappings([]string{".txt=CLIDocuments"})
	assert.NoError(t, err)
	
	category, exists = mapping.GetMapping(".txt")
	assert.True(t, exists)
	assert.Equal(t, "CLIDocuments", category)
}
