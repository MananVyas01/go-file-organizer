# Go File Organizer

A CLI tool to organize files in a directory by their file types.

## Features

- Organize files by file type into separate directories
- Dry-run mode to preview changes before applying them
- Command-line interface with intuitive flags

## Usage

```bash
go run main.go --path <directory> [--dry-run]
```

### Flags

- `--path`: Path to the folder to organize (required)
- `--dry-run`: Preview actions without moving files (optional)
- `--help`: Show usage information

## Examples

```bash
# Organize files in the current directory
go run main.go --path .

# Preview what would be organized without making changes
go run main.go --path ./downloads --dry-run

# Show help
go run main.go --help
```

## Project Structure

```
go-file-organizer/
├── cmd/              # CLI entry point
├── internal/
│   ├── organizer/    # File organizing logic
│   └── utils/        # Helpers (e.g., file type detection)
├── assets/           # Sample folder with mixed files
├── config/           # For JSON config/rules
├── main.go
└── README.md
```

## Development

1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Build with `go build` or run directly with `go run main.go`
