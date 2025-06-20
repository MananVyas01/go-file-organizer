# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.0.0] - 2025-06-20

### Added
- 🎉 **Initial Release** - Complete file organization CLI tool
- 🗂️ **Smart File Organization** - Automatically categorize files by extension
- 🔍 **Dry-Run Mode** - Preview changes before applying them (`--dry-run`)
- ⚙️ **Custom Configuration** - JSON config file support for extension mappings
- 🚫 **Ignore Patterns** - `.organizerignore` file support with wildcards and patterns
- 📝 **Comprehensive Logging** - Detailed operation logs with summary reports
- 🌍 **Cross-Platform Support** - Binaries for Windows, macOS, and Linux (amd64/arm64)
- 🎯 **CLI Overrides** - Command-line extension mapping overrides (`--map`)
- 📊 **Summary Reports** - Clear organization statistics and results
- 🧪 **Extensive Testing** - Complete unit test coverage (31 tests)

### Features
- **File Categories**: Documents, Images, Videos, Audio, Code, Archives, and more
- **Pattern Matching**: Support for exact filenames, wildcards, and directory patterns
- **Ignore Rules**: Root-level patterns, directory exclusions, and custom rules
- **Version Display**: `--version` flag to show application version
- **Error Handling**: Graceful error handling with helpful messages
- **Makefile**: Automated building and packaging for all platforms

### Documentation
- 📖 **Complete README** with installation and usage instructions
- 🔧 **Example Configuration Files** for easy setup
- 📄 **MIT License** for open-source use
- 🎯 **Sample Ignore Patterns** with common use cases

### Technical Details
- **Built with Go 1.21+**
- **Zero external runtime dependencies**
- **Cross-compiled for 6 platforms**
- **Optimized binaries with stripped symbols**
- **Comprehensive error handling and validation**

### Supported Platforms
- Linux (amd64, arm64)
- Windows (amd64, arm64)
- macOS (Intel, Apple Silicon)

---

## Planned for Future Releases

### v1.1.0 (Planned)
- [ ] GUI application version
- [ ] Batch processing mode
- [ ] Undo functionality
- [ ] Advanced pattern matching (regex support)
- [ ] Plugin system for custom handlers

### v1.2.0 (Planned)
- [ ] File content-based categorization
- [ ] Duplicate file detection
- [ ] Integration with cloud storage
- [ ] Progress bars for large operations
