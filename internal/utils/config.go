package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Config represents the configuration file structure
type Config struct {
	CustomMappings map[string]string `json:"customMappings"`
	Description    string            `json:"description,omitempty"`
}

// ExtensionMapping holds all extension mappings from various sources
type ExtensionMapping struct {
	mappings map[string]string
	sources  map[string]string // tracks where each mapping came from
}

// NewExtensionMapping creates a new extension mapping with default values
func NewExtensionMapping(defaultMappings map[string]string) *ExtensionMapping {
	mapping := &ExtensionMapping{
		mappings: make(map[string]string),
		sources:  make(map[string]string),
	}
	
	// Add default mappings
	for ext, category := range defaultMappings {
		mapping.mappings[ext] = category
		mapping.sources[ext] = "default"
	}
	
	return mapping
}

// LoadConfig loads configuration from a JSON file and merges with existing mappings
func (em *ExtensionMapping) LoadConfig(configPath string) error {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config file doesn't exist, that's okay
		return nil
	}
	
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	
	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config JSON: %v", err)
	}
	
	// Validate and merge custom mappings
	count := 0
	for ext, category := range config.CustomMappings {
		if err := em.validateExtension(ext); err != nil {
			fmt.Printf("Warning: Invalid extension '%s' in config: %v\n", ext, err)
			continue
		}
		
		if err := em.validateCategory(category); err != nil {
			fmt.Printf("Warning: Invalid category '%s' for extension '%s': %v\n", category, ext, err)
			continue
		}
		
		em.mappings[strings.ToLower(ext)] = category
		em.sources[strings.ToLower(ext)] = "config"
		count++
	}
	
	fmt.Printf("Loaded %d custom mappings from config file\n", count)
	return nil
}

// ApplyCLIMappings applies command-line mapping overrides
func (em *ExtensionMapping) ApplyCLIMappings(cliMappings []string) error {
	count := 0
	for _, mapping := range cliMappings {
		parts := strings.SplitN(mapping, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid mapping format '%s', expected '.ext=Category'", mapping)
		}
		
		ext := strings.TrimSpace(parts[0])
		category := strings.TrimSpace(parts[1])
		
		if err := em.validateExtension(ext); err != nil {
			return fmt.Errorf("invalid extension '%s': %v", ext, err)
		}
		
		if err := em.validateCategory(category); err != nil {
			return fmt.Errorf("invalid category '%s': %v", category, err)
		}
		
		em.mappings[strings.ToLower(ext)] = category
		em.sources[strings.ToLower(ext)] = "cli"
		count++
	}
	
	if count > 0 {
		fmt.Printf("Applied %d CLI mapping overrides\n", count)
	}
	return nil
}

// GetMapping returns the category for a given extension
func (em *ExtensionMapping) GetMapping(ext string) (string, bool) {
	category, exists := em.mappings[strings.ToLower(ext)]
	return category, exists
}

// GetMappings returns all current mappings
func (em *ExtensionMapping) GetMappings() map[string]string {
	result := make(map[string]string)
	for ext, category := range em.mappings {
		result[ext] = category
	}
	return result
}

// PrintSummary prints a summary of applied custom rules
func (em *ExtensionMapping) PrintSummary() {
	configCount := 0
	cliCount := 0
	
	for _, source := range em.sources {
		switch source {
		case "config":
			configCount++
		case "cli":
			cliCount++
		}
	}
	
	if configCount > 0 || cliCount > 0 {
		fmt.Printf("\n📋 Custom Rules Applied:\n")
		if configCount > 0 {
			fmt.Printf("  📄 Config file mappings: %d\n", configCount)
		}
		if cliCount > 0 {
			fmt.Printf("  ⚡ CLI overrides: %d\n", cliCount)
		}
	}
}

// validateExtension validates that an extension is properly formatted
func (em *ExtensionMapping) validateExtension(ext string) error {
	if ext == "" {
		return fmt.Errorf("extension cannot be empty")
	}
	
	if !strings.HasPrefix(ext, ".") {
		return fmt.Errorf("extension must start with '.'")
	}
	
	if len(ext) == 1 {
		return fmt.Errorf("extension must have content after '.'")
	}
	
	return nil
}

// validateCategory validates that a category name is reasonable
func (em *ExtensionMapping) validateCategory(category string) error {
	if category == "" {
		return fmt.Errorf("category cannot be empty")
	}
	
	if strings.TrimSpace(category) != category {
		return fmt.Errorf("category cannot have leading or trailing whitespace")
	}
	
	// Check for invalid characters that might cause filesystem issues
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(category, char) {
			return fmt.Errorf("category contains invalid character '%s'", char)
		}
	}
	
	return nil
}
