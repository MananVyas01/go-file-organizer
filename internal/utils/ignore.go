package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IgnoreManager handles file and directory ignore patterns
type IgnoreManager struct {
	patterns []string
	rootPath string
}

// NewIgnoreManager creates a new ignore manager
func NewIgnoreManager(rootPath string) *IgnoreManager {
	return &IgnoreManager{
		patterns: make([]string, 0),
		rootPath: rootPath,
	}
}

// LoadIgnoreFile loads patterns from .organizerignore file
func (im *IgnoreManager) LoadIgnoreFile(ignoreFilePath string) error {
	// Check if ignore file exists
	if _, err := os.Stat(ignoreFilePath); os.IsNotExist(err) {
		// Ignore file doesn't exist, that's okay
		return nil
	}
	
	file, err := os.Open(ignoreFilePath)
	if err != nil {
		return fmt.Errorf("failed to open ignore file: %v", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	lineCount := 0
	patternCount := 0
	
	for scanner.Scan() {
		lineCount++
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		im.patterns = append(im.patterns, line)
		patternCount++
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading ignore file: %v", err)
	}
	
	if patternCount > 0 {
		fmt.Printf("Loaded %d ignore patterns from .organizerignore\n", patternCount)
	}
	
	return nil
}

// ShouldIgnore checks if a file path should be ignored based on patterns
func (im *IgnoreManager) ShouldIgnore(filePath string) bool {
	// Convert to relative path from root for consistent matching
	relPath, err := filepath.Rel(im.rootPath, filePath)
	if err != nil {
		relPath = filePath
	}
	
	// Normalize path separators for cross-platform compatibility
	relPath = filepath.ToSlash(relPath)
	fileName := filepath.Base(filePath)
	
	for _, pattern := range im.patterns {
		if im.matchPattern(pattern, relPath, fileName) {
			return true
		}
	}
	
	return false
}

// matchPattern checks if a file path matches an ignore pattern
func (im *IgnoreManager) matchPattern(pattern, relPath, fileName string) bool {
	// Normalize pattern
	pattern = strings.TrimSpace(pattern)
	pattern = filepath.ToSlash(pattern)
	
	// Handle different pattern types
	
	// 1. Exact filename match
	if !strings.Contains(pattern, "/") && !strings.Contains(pattern, "*") {
		return fileName == pattern
	}
	
	// 2. Directory pattern (ends with /)
	if strings.HasSuffix(pattern, "/") {
		dirPattern := strings.TrimSuffix(pattern, "/")
		pathParts := strings.Split(relPath, "/")
		
		// Check if any directory in the path matches
		for _, part := range pathParts[:len(pathParts)-1] { // exclude filename
			if im.simpleMatch(dirPattern, part) {
				return true
			}
		}
		
		// Check if the relative path starts with this directory
		return strings.HasPrefix(relPath+"/", pattern)
	}
	
	// 3. Path pattern starting with / (from root)
	if strings.HasPrefix(pattern, "/") {
		pattern = strings.TrimPrefix(pattern, "/")
		return im.simpleMatch(pattern, relPath)
	}
	
	// 4. Wildcard patterns
	if strings.Contains(pattern, "*") {
		return im.wildcardMatch(pattern, relPath) || im.wildcardMatch(pattern, fileName)
	}
	
	// 5. Substring match anywhere in path
	return strings.Contains(relPath, pattern) || strings.Contains(fileName, pattern)
}

// simpleMatch performs simple string matching with basic wildcard support
func (im *IgnoreManager) simpleMatch(pattern, text string) bool {
	if !strings.Contains(pattern, "*") {
		return pattern == text
	}
	
	return im.wildcardMatch(pattern, text)
}

// wildcardMatch performs wildcard matching (* and ?)
func (im *IgnoreManager) wildcardMatch(pattern, text string) bool {
	// Simple wildcard implementation
	// This is a basic version - for production, consider using filepath.Match or a more robust library
	
	// Handle the simple case of * at the end
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(text, prefix)
	}
	
	// Handle the simple case of * at the beginning
	if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(text, suffix)
	}
	
	// For more complex patterns, use filepath.Match
	matched, err := filepath.Match(pattern, text)
	if err != nil {
		// If pattern is invalid, fall back to substring match
		return strings.Contains(text, strings.ReplaceAll(pattern, "*", ""))
	}
	
	return matched
}

// GetPatterns returns all loaded ignore patterns
func (im *IgnoreManager) GetPatterns() []string {
	result := make([]string, len(im.patterns))
	copy(result, im.patterns)
	return result
}

// PrintSummary prints a summary of loaded ignore patterns
func (im *IgnoreManager) PrintSummary() {
	if len(im.patterns) > 0 {
		fmt.Printf("  🚫 Ignore patterns: %d\n", len(im.patterns))
	}
}
