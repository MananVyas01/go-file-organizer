# Go File Organizer

A CLI tool to organize files in a directory by their file types.

## Features

- Organize files by file type into separate directories
- Dry-run mode to preview changes before applying them
- Configurable file extension mappings via `config/config.json`
- Ignore patterns support via `.organizerignore` file
- Comprehensive logging of operations
- Command-line interface with intuitive flags
- Extensive unit test coverage

## Usage

```bash
go run main.go --path <directory> [--dry-run]
```

### Flags

- `--path`: Path to the folder to organize (required)
- `--dry-run`: Preview actions without moving files (optional)
- `--config`: Path to custom config file (optional)
- `--mapping`: Override file extension mappings (e.g., "txt:Documents,py:Code")
- `--help`: Show usage information

## Configuration

### Extension Mapping

Create a `config/config.json` file to customize how file extensions are categorized:

```json
{
  "extensions": {
    "py": "Code",
    "txt": "Documents",
    "jpg": "Images"
  }
}
```

### Ignore Patterns

Create a `.organizerignore` file to specify files and patterns to skip:

```
# Ignore specific files
.DS_Store
desktop.ini

# Ignore file patterns
*.tmp
*.log

# Ignore directories
build/
node_modules/

# Ignore files at root level only
/important.txt
```

## Examples

```bash
# Organize files in the current directory
go run main.go --path .

# Preview what would be organized without making changes
go run main.go --path ./downloads --dry-run

# Use custom config file
go run main.go --path ./downloads --config ./custom-config.json

# Override extension mappings via command line
go run main.go --path ./downloads --mapping "py:Scripts,txt:Notes"

# Show help
go run main.go --help
```

## Project Structure

```
go-file-organizer/
├── cmd/                        # CLI entry point (future)
├── internal/
│   ├── organizer/             # File organizing logic
│   │   ├── organizer.go       # Core organization logic
│   │   ├── organizer_test.go  # Organization tests
│   │   └── scanner.go         # File scanning logic
│   └── utils/                 # Helper utilities
│       ├── config.go          # Configuration management
│       ├── config_test.go     # Config tests
│       ├── ignore.go          # Ignore pattern handling
│       ├── ignore_test.go     # Ignore pattern tests
│       ├── logger.go          # Logging utilities
│       └── logger_test.go     # Logger tests
├── assets/                    # Sample files for testing
├── config/                    # Configuration files
│   └── config.json           # Default extension mappings
├── test_assets/              # Test files (gitignored)
├── .editorconfig             # Editor configuration
├── .gitignore               # Git ignore patterns
├── .organizerignore         # File organization ignore patterns
├── go.mod                   # Go module definition
├── main.go                  # Application entry point
└── README.md               # This file
```

## Development

### Setup

1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Build with `go build` or run directly with `go run main.go`

### Testing

The project includes comprehensive unit tests for all core functionality:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run specific test package
go test ./internal/organizer -v
go test ./internal/utils -v
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```
