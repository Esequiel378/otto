# Makefile for Otto - 3D Graphics Engine
# Supports building for Linux and Windows targets

# Project configuration
PROJECT_NAME := otto
VERSION := 1.0.0
BUILD_DIR := build
ASSETS_DIR := assets

# Go build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -s -w"
CGO_ENABLED := 1

# Default target
.PHONY: all
all: build-linux

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Build for Linux (default)
.PHONY: build-linux
build-linux: $(BUILD_DIR)
	@echo "Building for Linux..."
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-amd64 .

# Build for Windows
.PHONY: build-windows
build-windows: $(BUILD_DIR)
	@echo "Building for Windows..."
	@echo "Note: Cross-compiling with CGO requires MinGW-w64 toolchain"
	@echo "If you get gcc errors, try: make install-toolchain"
	@echo "Or use: make build-windows-nocgo (without CGO)"
	@which x86_64-w64-mingw32-gcc > /dev/null || (echo "MinGW-w64 GCC not found. Run: make install-toolchain" && exit 1)
	@which x86_64-w64-mingw32-g++ > /dev/null || (echo "MinGW-w64 G++ not found. Run: make install-toolchain" && exit 1)
	CGO_ENABLED=$(CGO_ENABLED) CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CGO_LDFLAGS="-static-libgcc -static-libstdc++" GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64.exe .

# Build for Windows without CGO (fallback)
.PHONY: build-windows-nocgo
build-windows-nocgo: $(BUILD_DIR)
	@echo "Building for Windows (without CGO)..."
	@echo "Warning: This may not work with OpenGL/SDL dependencies"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64-nocgo.exe .

# Build for Windows with alternative CGO flags
.PHONY: build-windows-alt
build-windows-alt: $(BUILD_DIR)
	@echo "Building for Windows (alternative CGO approach)..."
	@echo "Trying with different CGO flags to resolve SDL2 compatibility issues"
	@which x86_64-w64-mingw32-gcc > /dev/null || (echo "MinGW-w64 GCC not found. Run: make install-toolchain" && exit 1)
	@which x86_64-w64-mingw32-g++ > /dev/null || (echo "MinGW-w64 G++ not found. Run: make install-toolchain" && exit 1)
	CGO_ENABLED=$(CGO_ENABLED) CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CGO_LDFLAGS="-static-libgcc -static-libstdc++ -lmsvcrt" CGO_CFLAGS="-D_WIN32_WINNT=0x0601" GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64-alt.exe .

# Build for both platforms
.PHONY: build-all
build-all: build-linux build-windows
	@echo "Build completed for all platforms"

# Build for current platform (detected automatically)
.PHONY: build-current
build-current: $(BUILD_DIR)
	@echo "Building for current platform..."
	CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME) .

# Run the application with metrics enabled
.PHONY: run
run:
	@echo "Running application with metrics enabled..."
	@echo "Metrics will be available at http://localhost:8080/metrics"
	@echo "Grafana dashboard: http://localhost:3030 (admin/admin)"
	@echo "Prometheus: http://localhost:9090"
	OTTO_METRICS_ENABLED=true go run ./cmd/playground/main.go

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Install cross-compilation toolchain
.PHONY: install-toolchain
install-toolchain:
	@echo "Installing MinGW-w64 toolchain for Windows cross-compilation..."
	@echo "This requires sudo privileges..."
	sudo apt-get update
	sudo apt-get install -y gcc-mingw-w64 g++-mingw-w64
	@echo "Toolchain installed successfully!"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Create release package (Linux)
.PHONY: release-linux
release-linux: build-linux
	@echo "Creating Linux release package..."
	cd $(BUILD_DIR) && tar -czf $(PROJECT_NAME)-linux-amd64-$(VERSION).tar.gz $(PROJECT_NAME)-linux-amd64

# Create release package (Windows)
.PHONY: release-windows
release-windows: build-windows
	@echo "Creating Windows release package..."
	cd $(BUILD_DIR) && zip $(PROJECT_NAME)-windows-amd64-$(VERSION).zip $(PROJECT_NAME)-windows-amd64.exe

# Create release packages for both platforms
.PHONY: release
release: release-linux release-windows
	@echo "Release packages created for all platforms"

# Development build with debug information
.PHONY: dev-build
dev-build: $(BUILD_DIR)
	@echo "Building development version..."
	CGO_ENABLED=$(CGO_ENABLED) go build -race -gcflags="all=-N -l" -o $(BUILD_DIR)/$(PROJECT_NAME)-dev .

# Monitoring commands
.PHONY: monitor-start
monitor-start:
	@echo "Starting Prometheus and Grafana monitoring stack..."
	@echo "Grafana: http://localhost:3000 (admin/admin)"
	@echo "Prometheus: http://localhost:9090"
	docker-compose up -d

.PHONY: monitor-stop
monitor-stop:
	@echo "Stopping monitoring stack..."
	docker-compose down

.PHONY: monitor-logs
monitor-logs:
	@echo "Showing monitoring stack logs..."
	docker-compose logs -f

.PHONY: monitor-status
monitor-status:
	@echo "Monitoring stack status:"
	docker-compose ps

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build-linux     - Build for Linux (amd64)"
	@echo "  build-windows   - Build for Windows (amd64)"
	@echo "  build-windows-nocgo - Build for Windows without CGO (fallback)"
	@echo "  build-windows-alt - Build for Windows with alternative CGO flags"
	@echo "  build-all       - Build for both platforms"
	@echo "  build-current   - Build for current platform"
	@echo "  run             - Run the application"
	@echo "  clean           - Clean build artifacts"
	@echo "  deps            - Install dependencies"
	@echo "  install-toolchain - Install MinGW-w64 for Windows cross-compilation"
	@echo "  test            - Run tests"
	@echo "  fmt             - Format code"
	@echo "  lint            - Lint code"
	@echo "  release-linux   - Create Linux release package"
	@echo "  release-windows - Create Windows release package"
	@echo "  release         - Create release packages for all platforms"
	@echo "  dev-build       - Build with debug information"
	@echo "  monitor-start   - Start Prometheus and Grafana monitoring stack"
	@echo "  monitor-stop    - Stop Prometheus and Grafana monitoring stack"
	@echo "  monitor-logs    - Show monitoring stack logs"
	@echo "  monitor-status  - Show monitoring stack status"
	@echo "  help            - Show this help message" 