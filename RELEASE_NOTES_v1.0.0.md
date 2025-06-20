# 🎉 go-file-organizer v1.0.0 - Initial Release

A powerful, cross-platform CLI tool for organizing files by type with intelligent categorization and flexible configuration.

## ✨ What's New

### 🗂️ **Smart File Organization**
- Automatically categorizes files by extension into logical folders (Documents, Images, Videos, Audio, Code, Archives, etc.)
- Support for 80+ file extensions out of the box
- Unknown files go to "Unknown" folder for easy review

### 🔍 **Safe Preview Mode** 
- `--dry-run` flag to preview all changes before applying them
- See exactly what files will be moved where
- No accidental file movements

### ⚙️ **Flexible Configuration**
- JSON configuration files for custom extension mappings
- Command-line overrides with `--map` flag
- Example configuration included

### 🚫 **Smart Ignore Patterns**
- `.organizerignore` file support (like .gitignore)
- Wildcard patterns (`*.tmp`, `build/`, etc.)
- Root-level exclusions and directory skipping
- Example ignore file included

### 📝 **Comprehensive Logging**
- Detailed operation logs with timestamps
- Summary reports showing files processed, moved, and skipped
- Easy troubleshooting and audit trails

### 🌍 **Cross-Platform Support**
- Native binaries for Linux, Windows, and macOS
- Both Intel and ARM64 architectures supported
- Zero dependencies - just download and run!

## 🚀 Quick Start

```bash
# Download binary for your platform from release assets below
# Extract and make executable (Linux/macOS):
tar -xzf go-file-organizer_*.tar.gz
sudo mv go-file-organizer /usr/local/bin/

# Always start with a dry-run to see what will happen:
go-file-organizer --path ./Downloads --dry-run

# If you're happy with the preview:
go-file-organizer --path ./Downloads
```

## 📦 Download Options

Choose the binary for your platform:

| Platform | Architecture | Download |
|----------|-------------|----------|
| **Linux** | x64 | `go-file-organizer_linux_amd64.tar.gz` |
| **Linux** | ARM64 | `go-file-organizer_linux_arm64.tar.gz` |
| **Windows** | x64 | `go-file-organizer_windows_amd64.zip` |
| **Windows** | ARM64 | `go-file-organizer_windows_arm64.zip` |
| **macOS** | Intel | `go-file-organizer_darwin_amd64.tar.gz` |
| **macOS** | Apple Silicon | `go-file-organizer_darwin_arm64.tar.gz` |

## 🔧 Usage Examples

```bash
# Basic organization with preview
go-file-organizer --path ./messy-folder --dry-run

# Custom extension mappings
go-file-organizer --path ./Downloads --map .py=Scripts --map .log=Logs

# Show version
go-file-organizer --version

# Get help
go-file-organizer --help
```

## 📋 Sample Output

```
🚀 ORGANIZING FILES...

📁 Created: Documents/
📁 Created: Images/ 
📁 Created: Code/

✅ Moved: report.pdf → Documents/report.pdf
✅ Moved: vacation.jpg → Images/vacation.jpg
✅ Moved: script.py → Code/script.py

📊 ORGANIZATION COMPLETE
Files processed: 15
Files moved: 12
Files skipped: 3
Directories created: 4

📝 Detailed log written to: organizer.log
```

## 🛠️ Building from Source

```bash
git clone https://github.com/MananVyas01/go-file-organizer.git
cd go-file-organizer
go build -o go-file-organizer .
```

## 📚 Documentation

- 📖 [Complete README](https://github.com/MananVyas01/go-file-organizer#readme)
- 🔧 [Configuration Guide](https://github.com/MananVyas01/go-file-organizer#%EF%B8%8F-configuration)
- 🐛 [Report Issues](https://github.com/MananVyas01/go-file-organizer/issues)
- 💡 [Request Features](https://github.com/MananVyas01/go-file-organizer/issues)

## ✅ What's Tested

This release includes comprehensive testing:
- ✅ 31 unit tests covering all core functionality
- ✅ Cross-platform compatibility verified
- ✅ Configuration and ignore pattern validation
- ✅ Dry-run mode accuracy
- ✅ Error handling and edge cases

---

**Perfect for organizing Downloads folders, project directories, or any messy file collections! 🗂️**

Made with ❤️ for keeping folders clean and organized.
