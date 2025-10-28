.PHONY: build build-windows build-mac build-linux build-all test test-coverage test-short clean

# Build for current platform
build:
	go build -o bin/addonprofiles-manager ./cmd/gui

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/addonprofiles-manager.exe ./cmd/gui

# Build for macOS
build-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/addonprofiles-manager-amd64 ./cmd/gui
	GOOS=darwin GOARCH=arm64 go build -o bin/addonprofiles-manager-arm64 ./cmd/gui

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/addonprofiles-manager ./cmd/gui

# Build for all platforms
build-all: build-windows build-mac build-linux

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

