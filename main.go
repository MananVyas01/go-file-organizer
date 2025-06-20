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

    // TODO: Call organizer package to handle organization
}
