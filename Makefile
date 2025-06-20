# Makefile for go-file-organizer

APP_NAME = go-file-organizer
BUILD_DIR = build
VERSION = $(shell go run -ldflags="-s -w" ./internal/version/version.go 2>/dev/null || echo "v1.0.0")

# Default target
.PHONY: all
all: clean build

# Clean build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)

# Build for current platform
.PHONY: build
build:
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) .

# Build for all platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)_linux_amd64 .
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)_linux_arm64 .
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)_windows_amd64.exe .
	
	# Windows ARM64
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)_windows_arm64.exe .
	
	# macOS AMD64 (Intel)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)_darwin_amd64 .
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)_darwin_arm64 .
	
	@echo "✅ Build complete! Binaries are in $(BUILD_DIR)/"
	@ls -la $(BUILD_DIR)/

# Create release archives
.PHONY: package
package: build-all
	@echo "Creating release packages..."
	cd $(BUILD_DIR) && \
	for binary in $(APP_NAME)_*; do \
		case "$$binary" in \
			*.exe) \
				zip "$${binary%.exe}.zip" "$$binary" ;; \
			*) \
				tar -czf "$$binary.tar.gz" "$$binary" ;; \
		esac; \
	done
	@echo "✅ Packages created!"
	@ls -la $(BUILD_DIR)/*.{zip,tar.gz} 2>/dev/null || true

# Test
.PHONY: test
test:
	go test ./... -v

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Vet code
.PHONY: vet
vet:
	go vet ./...

# Run quality checks
.PHONY: check
check: fmt vet test

# Install locally
.PHONY: install
install:
	go install .

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all         - Clean and build for current platform"
	@echo "  build       - Build for current platform"
	@echo "  build-all   - Build for all supported platforms"
	@echo "  package     - Create release packages (zip/tar.gz)"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format code"
	@echo "  vet         - Vet code"
	@echo "  check       - Run fmt, vet, and test"
	@echo "  install     - Install locally"
	@echo "  clean       - Clean build directory"
	@echo "  help        - Show this help"
