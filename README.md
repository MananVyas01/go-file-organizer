# Go File Organizer

[![Go Report Card](https://goreportcard.com/badge/github.com/MananVyas01/go-file-organizer)](https://goreportcard.com/report/github.com/MananVyas01/go-file-organizer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/MananVyas01/go-file-organizer.svg)](https://github.com/MananVyas01/go-file-organizer/releases)

A powerful CLI tool to automatically organize files in any directory by their file types. Keep your folders clean and organized with intelligent file categorization, custom rules, and ignore patterns.

## ✨ Features

- 🗂️ **Smart Organization**: Automatically categorizes files by extension into logical folders
- 🔍 **Dry-Run Mode**: Preview changes before applying them
- ⚙️ **Configurable**: Custom file extension mappings via JSON config
- 🚫 **Ignore Patterns**: Skip files and directories using `.organizerignore`
- 📝 **Detailed Logging**: Complete operation logs with summary reports
- 🌍 **Cross-Platform**: Available for Windows, macOS, and Linux
- 🎯 **CLI Overrides**: Quick extension mapping changes via command line
- 📊 **Summary Reports**: See what was organized at a glance
- 📈 **Progress Tracking**: Optional progress bar for large operations
- 👀 **Watch Mode**: Automatically organize new files as they appear in the directory

## 🚀 Quick Start

### Installation

#### Download Pre-built Binaries

1. Go to the [Releases page](https://github.com/MananVyas01/go-file-organizer/releases)
2. Download the binary for your platform:
   - `go-file-organizer_linux_amd64.tar.gz` - Linux (x64)
   - `go-file-organizer_darwin_amd64.tar.gz` - macOS (Intel)
   - `go-file-organizer_darwin_arm64.tar.gz` - macOS (Apple Silicon)
   - `go-file-organizer_windows_amd64.zip` - Windows (x64)

3. Extract and move to your PATH:

**Linux/macOS:**
```bash
tar -xzf go-file-organizer_*.tar.gz
sudo mv go-file-organizer /usr/local/bin/
```

**Windows:**
```powershell
# Extract the .exe and add to your PATH or run directly
```

#### Build from Source

```bash
git clone https://github.com/MananVyas01/go-file-organizer.git
cd go-file-organizer
go build -o go-file-organizer .
```

### Basic Usage

```bash
# Organize files in current directory (dry-run first!)
go-file-organizer --path . --dry-run

# Actually organize files
go-file-organizer --path .

# Organize a specific folder
go-file-organizer --path /path/to/messy/folder
```

## 📖 Usage

### Command Line Options

```bash
go-file-organizer [OPTIONS]

Options:
  --path string       Path to the folder to organize (required)
  --dry-run          Preview actions without moving files
  --version          Show version information
  --progress         Show progress bar during organization
  --watch            Watch directory for new files and organize them automatically
  --map string       Override extension mappings (format: .ext=Category)
  --help             Show usage information
```

### Examples

#### Basic Organization
```bash
# Always start with a dry-run to see what will happen
go-file-organizer --path ./Downloads --dry-run

# If you're happy with the preview, run it for real
go-file-organizer --path ./Downloads
```

#### Custom Extension Mappings
```bash
# Override specific extensions
go-file-organizer --path ./Downloads --map .py=Scripts --map .txt=Notes

# Multiple mappings in one command
go-file-organizer --path ./Downloads --map .py=Scripts,.txt=Notes,.log=Logs
```

#### Progress Bar for Large Operations
```bash
# Show progress bar while organizing many files
go-file-organizer --path ./Downloads --progress

# Combine with dry-run to see progress of
go-file-organizer --path ./Downloads --dry-run --progress
```

#### Watch Mode for Automatic Organization
```bash
# Watch a directory and automatically organize new files
go-file-organizer --path ./Downloads --watch

# Watch mode with dry-run to see what would happen
go-file-organizer --path ./Downloads --watch --dry-run

# Combine watch mode with progress bar and custom mappings
go-file-organizer --path ./Downloads --watch --progress --map .py=Scripts
```

**Note:** In watch mode, press `Ctrl+C` to stop monitoring the directory.

### Sample Output

**Standard Mode:**
```
🚀 ORGANIZING FILES...

📁 Created: Documents/
📁 Created: Images/
📁 Created: Code/

✅ Moved: report.pdf → Documents/report.pdf
✅ Moved: photo.jpg → Images/photo.jpg
✅ Moved: script.py → Code/script.py

📊 ORGANIZATION COMPLETE
Files processed: 15
Files moved: 12
Files skipped: 3
Directories created: 4

📝 Detailed log written to: organizer.log
```

**With Progress Bar (`--progress`):**
```
🚀 ORGANIZING FILES...
Organizing files  80% [===============================>        ] (12/15)

📊 ORGANIZATION COMPLETE
Files processed: 15
Files moved: 12
Files skipped: 3
Directories created: 4

📝 Detailed log written to: organizer.log
```

## ⚙️ Configuration

### Custom Extension Mappings

Create a `config/config.json` file to customize how file extensions are categorized:

```json
{
  "extensions": {
    "py": "Code",
    "js": "Code", 
    "html": "Code",
    "css": "Code",
    "txt": "Documents",
    "pdf": "Documents",
    "md": "Documents",
    "jpg": "Images",
    "png": "Images",
    "gif": "Images",
    "mp4": "Videos",
    "mp3": "Audio",
    "zip": "Archives"
  }
}
```

You can also use the example configuration file:
```bash
cp config/example.config.json config/config.json
# Edit config/config.json to your preferences
```

### File Categories

Default categories include:
- **Documents**: PDF, DOC, TXT, MD, etc.
- **Images**: JPG, PNG, GIF, SVG, etc. 
- **Videos**: MP4, AVI, MKV, MOV, etc.
- **Audio**: MP3, WAV, FLAC, etc.
- **Code**: JS, PY, GO, HTML, CSS, etc.
- **Archives**: ZIP, RAR, 7Z, TAR, etc.
- **Unknown**: Files with unrecognized extensions

### Ignore Patterns

Create a `.organizerignore` file to specify files and patterns to skip during organization:

```bash
# Copy the example ignore file
cp .organizerignore.example .organizerignore
# Edit .organizerignore to your needs
```

**Example `.organizerignore`:**
```
# System files
.DS_Store
desktop.ini

# Temporary files  
*.tmp
*.log

# Important files that should stay in place
README.md
LICENSE
Makefile

# Directories to skip
.git/
node_modules/
build/

# Files at root level only (leading slash)
/important-config.json
```

**Pattern Types:**
- `filename.ext` - Exact filename match
- `*.ext` - Wildcard for any filename with extension
- `directory/` - Skip entire directories 
- `/file.ext` - Files only at root level

## 🛠️ Development & Contributing

### Building from Source

```bash
# Clone the repository
git clone https://github.com/MananVyas01/go-file-organizer.git
cd go-file-organizer

# Install dependencies
go mod tidy

# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Format and vet code
make check
```

### Cross-Platform Builds

Use the included Makefile for easy cross-platform compilation:

```bash
# Build for all supported platforms
make build-all

# Create release packages
make package
```

This creates binaries for:
- Linux (amd64, arm64)
- Windows (amd64, arm64)  
- macOS (amd64, arm64)

### Testing

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

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code follows the project standards:
- Run `go fmt ./...` before committing
- Run `go vet ./...` to check for issues
- Add tests for new functionality
- Update documentation as needed

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

## � Future Roadmap

We're continuously improving go-file-organizer! Here's what's planned for upcoming releases:

### 🌐 v1.3.0 - Web UI (Planned)
- **Visual Interface**: Beautiful web-based GUI for file organization
- **Drag & Drop**: Intuitive file management with drag-and-drop support
- **Real-time Preview**: Live preview of organization changes before applying
- **Configuration Management**: Visual config editor for extension mappings

### 🗑️ v1.4.0 - Smart Safety Features (Planned)
- **Trash Mode**: Move files to trash/recycle bin instead of permanent moves
- **Undo Operations**: Rollback file organization changes
- **Backup Creation**: Automatic backups before major operations
- **Conflict Resolution**: Smart handling of duplicate file names

### ☁️ v1.5.0 - Cloud Integration (not Planned yet)
- **Cloud Storage**: Support for Google Drive, Dropbox, OneDrive
- **Sync Across Devices**: Keep organization rules synchronized
- **Remote Organization**: Organize cloud folders directly

### 🔧 Advanced Features (Future)
- **Custom Scripts**: Run custom scripts on organized files
- **File Content Analysis**: Organization based on file content
- **API Server**: REST API for integration with other tools
- **Mobile Apps**: iOS and Android companion apps

> **Want to contribute?** Join our [GitHub Discussions](https://github.com/MananVyas01/go-file-organizer/discussions) or [Request Features](https://github.com/MananVyas01/go-file-organizer/issues)

## �📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [testify](https://github.com/stretchr/testify) for testing
- Inspired by file organization tools and the need for a cross-platform solution

## 📧 Support

- 🐛 [Report Issues](https://github.com/MananVyas01/go-file-organizer/issues)
- 💡 [Request Features](https://github.com/MananVyas01/go-file-organizer/issues)
- 📖 [Documentation](https://github.com/MananVyas01/go-file-organizer/wiki)

---

**Made with ❤️ by MananVyas01 for organizing messy folders everywhere!**
