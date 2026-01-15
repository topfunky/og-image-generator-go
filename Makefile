.PHONY: help build clean test run install fmt lint

# Default target
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  clean       - Remove build artifacts and test outputs"
	@echo "  test        - Run tests"
	@echo "  run         - Build and run with example arguments"
	@echo "  install     - Build and install binary to GOPATH"
	@echo "  fmt         - Format code with gofmt"
	@echo "  lint        - Run golint checks"
	@echo "  all         - Clean, format, lint, test, and build"

# Build the binary
build:
	go build -ldflags="-X main.commit=$$(git rev-parse --short HEAD)" -o og-image-generator

# Remove build artifacts
clean:
	rm -f og-image-generator
	rm -f test*.png example-*.png dark-image.png wide-image.png long-title.png final-test.png social-image.png
	go clean

# Run tests
test:
	go test -v ./...

# Run with example
run: build
	./og-image-generator -title "Hello World" -url "https://example.com"

# Install binary
install: build
	go install

# Format code
fmt:
	gofmt -w .

# Lint checks
lint:
	@command -v staticcheck >/dev/null || (echo "Installing staticcheck..." && go install honnef.co/go/tools/cmd/staticcheck@latest)
	staticcheck ./...

# Build and clean all test outputs
all: clean fmt lint test build
	@echo "✓ Build complete"
