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

    // Scan files in the specified directory
    categories, err := organizer.ScanFiles(*path)
    if err != nil {
        fmt.Printf("Error scanning files: %v\n", err)
        os.Exit(1)
    }

    // Display the results
    fmt.Printf("\nFound files in %d categories:\n", len(categories))
    for category, files := range categories {
        fmt.Printf("\n%s (%d files):\n", category, len(files))
        for _, file := range files {
            if *dryRun {
                fmt.Printf("  [DRY RUN] %s\n", file)
            } else {
                fmt.Printf("  %s\n", file)
            }
        }
    }
}
