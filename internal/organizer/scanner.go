package organizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"go-file-organizer/internal/utils"
)

// extensionCategories maps file extensions to their categories
var extensionCategories = map[string]string{
	// Images
	".jpg":  "Images",
	".jpeg": "Images",
	".png":  "Images",
	".gif":  "Images",
	".bmp":  "Images",
	".svg":  "Images",
	".webp": "Images",
	".tiff": "Images",
	".ico":  "Images",

	// Documents
	".pdf":  "Documents",
	".doc":  "Documents",
	".docx": "Documents",
	".txt":  "Documents",
	".rtf":  "Documents",
	".odt":  "Documents",
	".pages": "Documents",

	// Spreadsheets
	".xls":  "Spreadsheets",
	".xlsx": "Spreadsheets",
	".csv":  "Spreadsheets",
	".ods":  "Spreadsheets",
	".numbers": "Spreadsheets",

	// Presentations
	".ppt":  "Presentations",
	".pptx": "Presentations",
	".odp":  "Presentations",
	".key":  "Presentations",

	// Code
	".go":   "Code",
	".js":   "Code",
	".ts":   "Code",
	".py":   "Code",
	".java": "Code",
	".c":    "Code",
	".cpp":  "Code",
	".h":    "Code",
	".hpp":  "Code",
	".cs":   "Code",
	".php":  "Code",
	".rb":   "Code",
	".rs":   "Code",
	".swift": "Code",
	".kt":   "Code",
	".scala": "Code",
	".html": "Code",
	".css":  "Code",
	".scss": "Code",
	".sass": "Code",
	".less": "Code",
	".xml":  "Code",
	".json": "Code",
	".yaml": "Code",
	".yml":  "Code",
	".toml": "Code",
	".ini":  "Code",
	".cfg":  "Code",
	".conf": "Code",

	// Archives
	".zip":  "Archives",
	".rar":  "Archives",
	".7z":   "Archives",
	".tar":  "Archives",
	".gz":   "Archives",
	".bz2":  "Archives",
	".xz":   "Archives",
	".iso":  "Archives",

	// Audio
	".mp3":  "Audio",
	".wav":  "Audio",
	".flac": "Audio",
	".aac":  "Audio",
	".ogg":  "Audio",
	".wma":  "Audio",
	".m4a":  "Audio",

	// Video
	".mp4":  "Video",
	".avi":  "Video",
	".mkv":  "Video",
	".mov":  "Video",
	".wmv":  "Video",
	".flv":  "Video",
	".webm": "Video",
	".m4v":  "Video",
	".3gp":  "Video",

	// Executables
	".exe":  "Executables",
	".msi":  "Executables",
	".deb":  "Executables",
	".rpm":  "Executables",
	".dmg":  "Executables",
	".app":  "Executables",
	".apk":  "Executables",
}

// ScanFiles recursively scans a directory and categorizes files by their extensions
func ScanFiles(rootPath string) (map[string][]string, error) {
	return ScanFilesWithConfig(rootPath, nil, nil)
}

// ScanFilesWithConfig recursively scans a directory with custom configuration and ignore rules
func ScanFilesWithConfig(rootPath string, extensionMapping *utils.ExtensionMapping, ignoreManager *utils.IgnoreManager) (map[string][]string, error) {
	// Initialize the result map
	categories := make(map[string][]string)
	
	// Check if the root path exists and is accessible
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", rootPath)
	}

	// Walk through the directory tree
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		// Handle errors during walk
		if err != nil {
			// Log the error but continue walking
			fmt.Printf("Warning: Could not access %s: %v\n", path, err)
			return nil
		}

		// Skip directories
		if info.IsDir() {
			// Check if this directory should be ignored
			if ignoreManager != nil && ignoreManager.ShouldIgnore(path) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if this file should be ignored
		if ignoreManager != nil && ignoreManager.ShouldIgnore(path) {
			return nil
		}

		// Get the file extension (case-insensitive)
		ext := strings.ToLower(filepath.Ext(info.Name()))
		
		// Determine the category
		var category string
		var exists bool
		
		if extensionMapping != nil {
			category, exists = extensionMapping.GetMapping(ext)
		} else {
			category, exists = extensionCategories[ext]
		}
		
		if !exists {
			// Handle files with no extension or unknown extensions
			if ext == "" {
				category = "No Extension"
			} else {
				category = "Unknown"
			}
		}

		// Add the file path to the appropriate category
		categories[category] = append(categories[category], path)
		
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	return categories, nil
}

// GetDefaultExtensionCategories returns a copy of the default extension mappings
func GetDefaultExtensionCategories() map[string]string {
	result := make(map[string]string)
	for ext, category := range extensionCategories {
		result[ext] = category
	}
	return result
}
