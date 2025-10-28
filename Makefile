.PHONY: build build-windows build-mac build-linux build-all test test-coverage test-short clean install-fyne-cross

# Build for current platform (native)
build:
	go build -o bin/addonprofiles-manager ./cmd/gui

# Install fyne-cross for cross-compilation
install-fyne-cross:
	go install github.com/fyne-io/fyne-cross@latest

# Build for Windows (requires fyne-cross and Docker)
build-windows:
	fyne-cross windows -arch=amd64 -output addonprofiles-manager.exe

# Build for macOS (requires fyne-cross and Docker)
build-mac:
	fyne-cross darwin -arch=amd64,arm64 -output addonprofiles-manager

# Build for Linux (requires fyne-cross and Docker)
build-linux:
	fyne-cross linux -arch=amd64 -output addonprofiles-manager

# Build for all platforms (requires fyne-cross and Docker)
build-all:
	fyne-cross windows linux darwin -arch=* -output addonprofiles-manager

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run short tests (skip slow tests)
test-short:
	go test -v -short ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

