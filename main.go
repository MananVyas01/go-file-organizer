// Entry point of the go-file-organizer CLI tool.
// This tool organizes files in a given directory by file type.
// Supports flags like --path, --dry-run, --map, and --help.

// 1. Use the flag package to parse command-line arguments:
//    --path string: the target directory
//    --dry-run bool: if true, show what would be done without moving files
//    --map .ext=Category: override extension mappings (can be used multiple times)

// 2. Print usage instructions when no path is provided.

// 3. Load configuration from config.json and .organizerignore files.

// 4. Call the internal organizer logic with custom configuration.

package main

import (
	"flag"
	"fmt"
	"go-file-organizer/internal/organizer"
	"go-file-organizer/internal/utils"
	"os"
)

// arrayFlags allows multiple values for the same flag
type arrayFlags []string

func (af *arrayFlags) String() string {
	return fmt.Sprintf("%v", *af)
}

func (af *arrayFlags) Set(value string) error {
	*af = append(*af, value)
	return nil
}

func main() {
	// Define flags
	path := flag.String("path", "", "Path to the folder to organize")
	dryRun := flag.Bool("dry-run", false, "Preview actions without moving files")
	help := flag.Bool("help", false, "Show usage")

	// Define flag for multiple mapping overrides
	var mapOverrides arrayFlags
	flag.Var(&mapOverrides, "map", "Override extension mappings (format: .ext=Category, can be used multiple times)")

	flag.Parse()

	if *help || *path == "" {
		fmt.Println("Usage: go-file-organizer --path <directory> [--dry-run] [--map .ext=Category]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	fmt.Println("Organizing path:", *path)
	fmt.Println("Dry run mode:", *dryRun)

	// Initialize configuration
	extensionMapping := utils.NewExtensionMapping(organizer.GetDefaultExtensionCategories())

	// Load config file if it exists
	configPath := "config/config.json"
	if err := extensionMapping.LoadConfig(configPath); err != nil {
		fmt.Printf("Warning: Could not load config file: %v\n", err)
		fmt.Println("Continuing with default mappings...")
	}

	// Apply CLI mapping overrides
	if len(mapOverrides) > 0 {
		if err := extensionMapping.ApplyCLIMappings(mapOverrides); err != nil {
			fmt.Printf("Error applying CLI mappings: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize ignore manager
	ignoreManager := utils.NewIgnoreManager(*path)
	ignoreFilePath := ".organizerignore"
	if err := ignoreManager.LoadIgnoreFile(ignoreFilePath); err != nil {
		fmt.Printf("Warning: Could not load ignore file: %v\n", err)
		fmt.Println("Continuing without ignore rules...")
	}

	// Print summary of custom rules
	if len(mapOverrides) > 0 {
		extensionMapping.PrintSummary()
		ignoreManager.PrintSummary()
	}

	// Initialize logger
	logger, err := utils.NewLogger("organizer.log")
	if err != nil {
		fmt.Printf("Warning: Could not create log file: %v\n", err)
		fmt.Println("Continuing without logging...")
	}
	defer func() {
		if logger != nil {
			logger.Close()
		}
	}()

	// Organize files
	if *dryRun {
		fmt.Println("\n🔮 DRY-RUN MODE: Simulating file organization...")
	} else {
		fmt.Println("\n🚀 ORGANIZING FILES...")
	}

	summary, err := organizer.OrganizeFilesWithConfig(*path, *dryRun, logger, extensionMapping, ignoreManager)
	if err != nil {
		fmt.Printf("Error organizing files: %v\n", err)
		os.Exit(1)
	}

	// Print summary
	organizer.PrintSummary(summary, *dryRun)

	if logger != nil {
		fmt.Printf("\n📝 Detailed log written to: organizer.log\n")
	}
}
