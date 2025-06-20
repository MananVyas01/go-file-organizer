package organizer

import (
	"fmt"
	"go-file-organizer/internal/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// OrganizeFiles organizes files in the given directory by their categories
func OrganizeFiles(rootPath string, isDryRun bool, logger *utils.Logger) (*utils.Summary, error) {
	return OrganizeFilesWithConfig(rootPath, isDryRun, logger, nil, nil, false)
}

// OrganizeFilesWithConfig organizes files with custom configuration and ignore rules
func OrganizeFilesWithConfig(rootPath string, isDryRun bool, logger *utils.Logger, extensionMapping *utils.ExtensionMapping, ignoreManager *utils.IgnoreManager, showProgress bool) (*utils.Summary, error) {
	summary := &utils.Summary{}

	// First, scan all files to get categories
	categories, err := ScanFilesWithConfig(rootPath, extensionMapping, ignoreManager)
	if err != nil {
		return summary, fmt.Errorf("failed to scan files: %v", err)
	}

	// Count total files and prepare progress bar
	totalFiles := 0
	for _, files := range categories {
		totalFiles += len(files)
	}
	summary.FilesScanned = totalFiles

	var bar *progressbar.ProgressBar
	if showProgress && totalFiles > 0 {
		bar = progressbar.NewOptions(totalFiles,
			progressbar.OptionSetDescription("Organizing files"),
			progressbar.OptionSetWidth(40),
			progressbar.OptionShowCount(),
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)
		defer func() {
			if bar != nil {
				bar.Finish()
				fmt.Println() // Add newline after progress bar
			}
		}()
	}

	// Process each category
	for category, files := range categories {
		// Skip categories that shouldn't be organized (like "Unknown" or "No Extension")
		if shouldSkipCategory(category) {
			summary.FilesSkipped += len(files)
			// Update progress bar for skipped files
			if bar != nil {
				for range files {
					bar.Add(1)
				}
			}
			continue
		}

		// Create category folder
		categoryPath := filepath.Join(rootPath, category)
		if err := createCategoryFolder(categoryPath, isDryRun, logger); err != nil {
			logger.LogError("Folder creation", categoryPath, err)
			continue
		}
		summary.FoldersCreated++

		// Move files to category folder
		for _, filePath := range files {
			fileName := filepath.Base(filePath)
			destPath := filepath.Join(categoryPath, fileName)

			// Skip if file is already in the target directory
			if filepath.Dir(filePath) == categoryPath {
				if bar != nil {
					bar.Add(1)
				}
				continue
			}

			if isDryRun {
				logger.LogDryRun(filePath, destPath)
				if !showProgress {
					fmt.Printf("  [DRY-RUN] Would move: %s -> %s\n", filePath, destPath)
				}
			} else {
				if err := moveFile(filePath, destPath); err != nil {
					logger.LogError("Move", filePath, err)
					if !showProgress {
						fmt.Printf("  [ERROR] Failed to move %s: %v\n", filePath, err)
					}
					if bar != nil {
						bar.Add(1)
					}
					continue
				}
				logger.LogMove(filePath, destPath)
				if !showProgress {
					fmt.Printf("  [MOVED] %s -> %s\n", filePath, destPath)
				}
			}
			summary.FilesMoved++
			
			// Update progress bar
			if bar != nil {
				bar.Add(1)
			}
		}
	}

	// Log summary
	logger.LogSummary(*summary)

	return summary, nil
}

// shouldSkipCategory determines if a category should be skipped during organization
func shouldSkipCategory(category string) bool {
	skipCategories := map[string]bool{
		"Unknown":      true,
		"No Extension": true,
	}
	return skipCategories[category]
}

// createCategoryFolder creates a folder for the category if it doesn't exist
func createCategoryFolder(folderPath string, isDryRun bool, logger *utils.Logger) error {
	// Check if folder already exists
	if _, err := os.Stat(folderPath); err == nil {
		return nil // Folder already exists
	}

	if isDryRun {
		logger.LogFolderCreation(folderPath, true)
		return nil
	}

	// Create the folder
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return fmt.Errorf("failed to create folder %s: %v", folderPath, err)
	}

	logger.LogFolderCreation(folderPath, false)
	return nil
}

// moveFile moves a file from source to destination
func moveFile(source, destination string) error {
	// Check if destination already exists
	if _, err := os.Stat(destination); err == nil {
		return fmt.Errorf("destination file already exists: %s", destination)
	}

	// Ensure destination directory exists
	destDir := filepath.Dir(destination)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Move the file
	if err := os.Rename(source, destination); err != nil {
		return fmt.Errorf("failed to move file: %v", err)
	}

	return nil
}

// PrintSummary prints a clean summary of the organization process
func PrintSummary(summary *utils.Summary, isDryRun bool) {
	separator := strings.Repeat("=", 50)

	fmt.Println("\n" + separator)
	if isDryRun {
		fmt.Println("📋 DRY-RUN SUMMARY")
	} else {
		fmt.Println("📋 ORGANIZATION SUMMARY")
	}
	fmt.Println(separator)

	fmt.Printf("✅  Total files scanned: %d\n", summary.FilesScanned)

	if isDryRun {
		fmt.Printf("🔮  Files that would be moved: %d\n", summary.FilesMoved)
		fmt.Printf("📁  Folders that would be created: %d\n", summary.FoldersCreated)
	} else {
		fmt.Printf("🔀  Files moved: %d\n", summary.FilesMoved)
		fmt.Printf("📁  Folders created: %d\n", summary.FoldersCreated)
	}

	fmt.Printf("🚫  Skipped (unknown/no extension): %d\n", summary.FilesSkipped)
	fmt.Println(separator)
}
