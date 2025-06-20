package organizer

import (
	"fmt"
	"go-file-organizer/internal/utils"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
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

// StartWatchMode starts watching the directory for new files and organizes them automatically
func StartWatchMode(rootPath string, isDryRun bool, logger *utils.Logger, extensionMapping *utils.ExtensionMapping, ignoreManager *utils.IgnoreManager, showProgress bool) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %v", err)
	}
	defer watcher.Close()

	// Add the root directory to watch
	err = watcher.Add(rootPath)
	if err != nil {
		return fmt.Errorf("failed to watch directory %s: %v", rootPath, err)
	}

	// Channel to listen for interrupt signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Debounce duplicate events (some file operations trigger multiple events)
	eventDebounce := make(map[string]time.Time)
	debounceDelay := 500 * time.Millisecond

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			// Only process file creation and write events
			if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
				// Debounce events
				now := time.Now()
				if lastTime, exists := eventDebounce[event.Name]; exists && now.Sub(lastTime) < debounceDelay {
					continue
				}
				eventDebounce[event.Name] = now

				// Check if it's a regular file (not a directory)
				fileInfo, err := os.Stat(event.Name)
				if err != nil || fileInfo.IsDir() {
					continue
				}

				// Check if file should be ignored
				if ignoreManager != nil && ignoreManager.ShouldIgnore(event.Name) {
					if logger != nil {
						logger.LogMove(event.Name, "IGNORED: "+event.Name)
					}
					continue
				}

				// Get file extension and category
				extension := strings.ToLower(filepath.Ext(event.Name))
				if extension == "" {
					if logger != nil {
						logger.LogMove(event.Name, "SKIPPED: no extension")
					}
					continue
				}

				var category string
				var exists bool
				if extensionMapping != nil {
					category, exists = extensionMapping.GetMapping(extension)
				} else {
					category, exists = extensionCategories[extension]
				}

				if !exists {
					if logger != nil {
						logger.LogMove(event.Name, "SKIPPED: unknown extension")
					}
					continue
				}

				// Skip categories that shouldn't be organized
				if shouldSkipCategory(category) {
					if logger != nil {
						logger.LogMove(event.Name, "SKIPPED: "+category)
					}
					continue
				}

				// Organize the file
				filename := filepath.Base(event.Name)
				targetDir := filepath.Join(rootPath, category)
				targetPath := filepath.Join(targetDir, filename)

				if isDryRun {
					fmt.Printf("🔮 [WATCH] Would move: %s → %s/%s\n", event.Name, category, filename)
					if logger != nil {
						logger.LogDryRun(event.Name, targetPath)
					}
				} else {
					// Create target directory if it doesn't exist
					if err := createCategoryFolder(targetDir, false, logger); err != nil {
						fmt.Printf("❌ [WATCH] Error creating directory %s: %v\n", targetDir, err)
						if logger != nil {
							logger.LogError("Folder creation", targetDir, err)
						}
						continue
					}

					// Move the file
					if err := moveFile(event.Name, targetPath); err != nil {
						fmt.Printf("❌ [WATCH] Error moving file %s: %v\n", event.Name, err)
						if logger != nil {
							logger.LogError("File move", event.Name, err)
						}
						continue
					}

					fmt.Printf("✅ [WATCH] Moved: %s → %s/%s\n", filename, category, filename)
					if logger != nil {
						logger.LogMove(event.Name, targetPath)
					}
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Printf("⚠️  [WATCH] Watcher error: %v\n", err)
			if logger != nil {
				logger.LogError("Watcher", "filesystem", err)
			}

		case <-interrupt:
			fmt.Println("\n🛑 Watch mode stopped by user")
			return nil
		}
	}
}
