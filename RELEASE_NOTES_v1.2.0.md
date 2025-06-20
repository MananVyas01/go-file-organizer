# 🎉 go-file-organizer v1.2.0 - Watch Mode Release

**Release Date:** June 20, 2025

## 🚀 What's New

### 👀 Watch Mode - The Game Changer!

We're excited to introduce **Watch Mode** - a powerful new feature that automatically organizes files as they appear in your directories!

```bash
# Start watching a directory and organize files automatically
go-file-organizer --path ./Downloads --watch

# Watch mode with dry-run to see what would happen
go-file-organizer --path ./Downloads --watch --dry-run

# Combine with other features for maximum flexibility
go-file-organizer --path ./Downloads --watch --progress --map .py=Scripts
```

### ✨ Key Features

- **🔄 Real-time Monitoring**: Uses efficient file system events to detect new files instantly
- **⏰ Event Debouncing**: Intelligently handles duplicate file operations to prevent conflicts
- **🛑 Graceful Shutdown**: Clean exit with `Ctrl+C` - no forced termination needed
- **📋 Comprehensive Logging**: All watch events are logged for complete traceability
- **🔮 Dry-run Support**: Preview watch mode actions without actually moving files

### 🎯 Perfect For

- **Downloads Folder**: Automatically organize downloads as they complete
- **Camera Imports**: Organize photos and videos as they're imported
- **Project Directories**: Keep development folders clean as files are added
- **Shared Folders**: Maintain organization in collaborative environments

## 🔧 Technical Improvements

- **New Dependency**: Added `github.com/fsnotify/fsnotify` for cross-platform file monitoring
- **Enhanced CLI**: Updated help text and usage examples
- **Test Coverage**: Added comprehensive tests for watch functionality
- **Documentation**: Updated README with watch mode examples and best practices

## 📊 Version History Recap

| Version | Key Feature | Release Date |
|---------|-------------|--------------|
| v1.2.0  | **Watch Mode** - Automatic file monitoring | June 20, 2025 |
| v1.1.0  | Progress Bar - Visual progress tracking | June 20, 2025 |
| v1.0.0  | Initial Release - Core organization features | June 20, 2025 |

## 📦 Download Options

### Pre-built Binaries

Download the appropriate binary for your platform:

- **Linux (x64)**: `go-file-organizer_linux_amd64.tar.gz`
- **Linux (ARM64)**: `go-file-organizer_linux_arm64.tar.gz`
- **Windows (x64)**: `go-file-organizer_windows_amd64.zip`
- **Windows (ARM64)**: `go-file-organizer_windows_arm64.zip`
- **macOS (Intel)**: `go-file-organizer_darwin_amd64.tar.gz`
- **macOS (Apple Silicon)**: `go-file-organizer_darwin_arm64.tar.gz`

### Installation

1. Download the appropriate binary for your platform
2. Extract the archive: `tar -xzf go-file-organizer_*.tar.gz` (Unix) or unzip (Windows)
3. Make executable (Unix): `chmod +x go-file-organizer_*`
4. Run: `./go-file-organizer --version` to verify installation

## 🛠️ Build from Source

```bash
git clone https://github.com/MananVyas01/go-file-organizer.git
cd go-file-organizer
go build -o go-file-organizer .
```

## 📖 Quick Start with Watch Mode

```bash
# 1. Test with dry-run first
go-file-organizer --path ./Downloads --watch --dry-run

# 2. If satisfied, run for real
go-file-organizer --path ./Downloads --watch

# 3. Watch mode will show:
👀 Starting watch mode for directory: ./Downloads
Press Ctrl+C to stop watching...
✅ [WATCH] Moved: document.pdf → Documents/document.pdf
✅ [WATCH] Moved: image.jpg → Images/image.jpg
```

## 🤝 Community & Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/MananVyas01/go-file-organizer/issues)
- **Discussions**: [Join the community](https://github.com/MananVyas01/go-file-organizer/discussions)
- **Documentation**: [Complete guide in README](https://github.com/MananVyas01/go-file-organizer#readme)

## 🔮 What's Next?

We're continuously improving go-file-organizer! Future considerations include:

- 🌐 Web UI for visual management
- 🗑️ Trash mode for safe file handling
- 🔗 Integration with cloud storage services
- 📱 Mobile companion app

---

**Full Changelog**: [View all changes](https://github.com/MananVyas01/go-file-organizer/blob/main/CHANGELOG.md)

**Previous Release**: [v1.1.0](https://github.com/MananVyas01/go-file-organizer/releases/tag/v1.1.0)
