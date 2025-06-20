// Entry point of the go-file-organizer CLI tool.
// This tool organizes files in a given directory by file type.
// Supports flags like --path, --dry-run, and --help.

// 1. Use the flag package to parse command-line arguments:
//    --path string: the target directory
//    --dry-run bool: if true, show what would be done without moving files

// 2. Print usage instructions when no path is provided.

// 3. Call the internal organizer logic (to be implemented in next stage).

package main

import (
    "flag"
    "fmt"
    "os"
    "go-file-organizer/internal/organizer"
    "go-file-organizer/internal/utils"
)

func main() {
    // Define flags
    path := flag.String("path", "", "Path to the folder to organize")
    dryRun := flag.Bool("dry-run", false, "Preview actions without moving files")
    help := flag.Bool("help", false, "Show usage")

    flag.Parse()

    if *help || *path == "" {
        fmt.Println("Usage: go-file-organizer --path <directory> [--dry-run]")
        flag.PrintDefaults()
        os.Exit(0)
    }

    fmt.Println("Organizing path:", *path)
    fmt.Println("Dry run mode:", *dryRun)

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

    summary, err := organizer.OrganizeFiles(*path, *dryRun, logger)
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
