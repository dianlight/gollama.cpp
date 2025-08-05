# Makefile for gollama.cpp
# Cross-platform Go bindings for llama.cpp using purego

# Version information
VERSION ?= 0.1.0
LLAMA_CPP_BUILD ?= b6089
FULL_VERSION = v$(VERSION)-llamacpp.$(LLAMA_CPP_BUILD)

# Check everything
.PHONY: check
check: fmt vet lint sec test

# Package releases

# Go configuration
GO ?= go
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Build directories
BUILD_DIR = build
DIST_DIR = dist
EXAMPLES_DIR = examples

# llama.cpp configuration
LLAMA_CPP_DIR = $(BUILD_DIR)/llama.cpp
LLAMA_CPP_REPO = https://github.com/ggerganov/llama.cpp.git

# Platform-specific configurations
PLATFORMS = darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

# Default target
.PHONY: all
all: build

# Clean everything
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	$(GO) clean -cache

# Clean libraries only
.PHONY: clean-libs
clean-libs:
	@echo "Cleaning library cache..."
	$(GO) run ./cmd/gollama-download -clean-cache

# Initialize/update dependencies
.PHONY: deps
deps:
	$(GO) mod download
	$(GO) mod tidy

# Build for current platform
.PHONY: build
build: deps
	@echo "Building gollama.cpp for $(GOOS)/$(GOARCH)"
	mkdir -p $(BUILD_DIR)/$(GOOS)_$(GOARCH)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(BUILD_DIR)/$(GOOS)_$(GOARCH)/ ./...

# Build for all platforms
.PHONY: build-all
build-all: deps
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		echo "Building for $$os/$$arch"; \
		mkdir -p $(BUILD_DIR)/$$os\_$$arch; \
		GOOS=$$os GOARCH=$$arch $(GO) build -o $(BUILD_DIR)/$$os\_$$arch/ ./...; \
	done

# Build examples
.PHONY: build-examples
build-examples: build
	@echo "Building examples"
	cd $(EXAMPLES_DIR) && $(GO) build ./...

# Test with library download
.PHONY: test
test: deps
	@echo "Running tests (libraries will be downloaded automatically)"
	$(GO) test -v ./...

# Test with race detection
.PHONY: test-race
test-race: deps
	@echo "Running tests with race detection"
	$(GO) test -race -v ./...

# Test library download functionality
.PHONY: test-download
test-download: deps
	@echo "Testing library download functionality"
	$(GO) run ./cmd/gollama-download -test-download

# Run platform-specific tests
.PHONY: test-platform
test-platform:
	@echo "Running platform-specific tests"
	$(GO) test -v -run TestPlatformSpecific ./...

# Test cross-compilation for all platforms
.PHONY: test-cross-compile
test-cross-compile:
	@echo "Testing cross-compilation for all platforms..."
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		echo "Building for $$GOOS/$$GOARCH..."; \
		env GOOS=$$GOOS GOARCH=$$GOARCH $(GO) build -v ./... || exit 1; \
	done
	@echo "All cross-compilation tests passed!"

# Test library download for specific platforms
.PHONY: test-download-platforms
test-download-platforms:
	@echo "Testing library download for different platforms..."
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		echo "Testing download for $$GOOS/$$GOARCH..."; \
		env GOOS=$$GOOS GOARCH=$$GOARCH $(GO) run ./cmd/gollama-download -test-download || echo "Download test for $$GOOS/$$GOARCH completed"; \
	done

# Download and verify libraries for current platform
.PHONY: download-libs
download-libs: deps
	@echo "Downloading llama.cpp libraries for $(GOOS)/$(GOARCH)"
	$(GO) run ./cmd/gollama-download -download -version $(LLAMA_CPP_BUILD)

# Download libraries for all platforms (for testing)
.PHONY: download-libs-all
download-libs-all: deps
	@echo "Downloading llama.cpp libraries for all platforms"
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		echo "Downloading for $$GOOS/$$GOARCH..."; \
		env GOOS=$$GOOS GOARCH=$$GOARCH $(GO) run ./cmd/gollama-download -download -version $(LLAMA_CPP_BUILD) || echo "Download for $$GOOS/$$GOARCH completed"; \
	done

# Test compilation for specific platform  
.PHONY: test-compile-windows
test-compile-windows:
	@echo "Testing Windows compilation"
	GOOS=windows GOARCH=amd64 $(GO) build -v ./...
	GOOS=windows GOARCH=arm64 $(GO) build -v ./...

.PHONY: test-compile-linux  
test-compile-linux:
	@echo "Testing Linux compilation"
	GOOS=linux GOARCH=amd64 $(GO) build -v ./...
	GOOS=linux GOARCH=arm64 $(GO) build -v ./...

.PHONY: test-compile-darwin
test-compile-darwin:
	@echo "Testing macOS compilation" 
	GOOS=darwin GOARCH=amd64 $(GO) build -v ./...
	GOOS=darwin GOARCH=arm64 $(GO) build -v ./...

# Benchmark
.PHONY: bench
bench: deps
	@echo "Running benchmarks (libraries will be downloaded automatically)"
	$(GO) test -bench=. -benchmem ./...

# Lint
.PHONY: lint
lint:
	@echo "Running linter"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping linting"; \
	fi

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code"
	$(GO) fmt ./...

# Vet code
.PHONY: vet
vet:
	@echo "Vetting code"
	$(GO) vet ./...

# Security check
.PHONY: sec
sec:
	@echo "Running security check"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec -exclude=G103,G104,G115,G304 -severity=medium ./...; \
	else \
		echo "gosec not found, skipping security check"; \
	fi

# Check everything
.PHONY: check
check: fmt vet lint sec test

# Clone llama.cpp repository for cross-reference checks
.PHONY: clone-llamacpp
clone-llamacpp:
	@if [ ! -d "$(LLAMA_CPP_DIR)" ]; then \
		echo "Cloning llama.cpp repository for cross-reference"; \
		mkdir -p $(BUILD_DIR); \
		git clone $(LLAMA_CPP_REPO) $(LLAMA_CPP_DIR); \
	fi
	@echo "Checking out build $(LLAMA_CPP_BUILD)"
	cd $(LLAMA_CPP_DIR) && git fetch && git checkout $(LLAMA_CPP_BUILD)


# Package releases
.PHONY: release
release: clean build-all
	@echo "Creating release packages"
	mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		echo "Packaging $$os/$$arch"; \
		pkg_name="gollama.cpp-$(FULL_VERSION)-$$os-$$arch"; \
		mkdir -p $(DIST_DIR)/$$pkg_name; \
		cp -r $(BUILD_DIR)/$$os\_$$arch/* $(DIST_DIR)/$$pkg_name/ 2>/dev/null || true; \
		cp README.md LICENSE CHANGELOG.md $(DIST_DIR)/$$pkg_name/; \
		cd $(DIST_DIR) && zip -r $$pkg_name.zip $$pkg_name && rm -rf $$pkg_name; \
	done

# Quick release for current platform
.PHONY: release-current
release-current: clean build
	@echo "Creating release package for $(GOOS)/$(GOARCH)"
	mkdir -p $(DIST_DIR)
	pkg_name="gollama.cpp-$(FULL_VERSION)-$(GOOS)-$(GOARCH)"
	mkdir -p $(DIST_DIR)/$$pkg_name
	cp -r $(BUILD_DIR)/$(GOOS)_$(GOARCH)/* $(DIST_DIR)/$$pkg_name/ 2>/dev/null || true
	cp README.md LICENSE CHANGELOG.md $(DIST_DIR)/$$pkg_name/
	cd $(DIST_DIR) && zip -r $$pkg_name.zip $$pkg_name && rm -rf $$pkg_name
	@echo "Release package created: $(DIST_DIR)/$$pkg_name.zip"

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools"
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest

# Download model file
.PHONY: model_download
model_download:
	@echo "Downloading models"
	@mkdir -p models
	@if [ ! -f "models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf" ]; then \
		echo "Downloading TinyLlama model..."; \
		curl -L -o models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf \
			"https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"; \
		echo "TinyLlama model downloaded successfully"; \
	else \
		echo "TinyLlama model already exists in models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"; \
	fi
	@if [ ! -f "models/gritlm-7b_q4_1.gguf" ]; then \
		echo "Downloading GritLM model..."; \
		curl -L -o models/gritlm-7b_q4_1.gguf \
			"https://huggingface.co/cohesionet/GritLM-7B_gguf/resolve/main/gritlm-7b_q4_1.gguf"; \
		echo "GritLM model downloaded successfully"; \
	else \
		echo "GritLM model already exists in models/gritlm-7b_q4_1.gguf"; \
	fi

# Show version information
.PHONY: version
version:
	@echo "Gollama.cpp Version: $(VERSION)"
	@echo "llama.cpp Build: $(LLAMA_CPP_BUILD)"
	@echo "Full Version: $(FULL_VERSION)"

# Help
.PHONY: help
help:
	@echo "Gollama.cpp Makefile"
	@echo ""
	@echo "Main targets:"
	@echo "  build              Build for current platform"
	@echo "  build-all          Build for all platforms"
	@echo "  build-examples     Build examples"
	@echo "  test               Run tests (downloads libraries automatically)"
	@echo "  test-race          Run tests with race detection"
	@echo "  bench              Run benchmarks"
	@echo "  clean              Clean all build artifacts"
	@echo ""
	@echo "Library management:"
	@echo "  download-libs      Download llama.cpp libraries for current platform"
	@echo "  download-libs-all  Download llama.cpp libraries for all platforms"
	@echo "  test-download      Test library download functionality"
	@echo "  test-download-platforms  Test downloads for all platforms"
	@echo "  clean-libs         Clean library cache (forces re-download)"
	@echo ""
	@echo "Quality assurance:"
	@echo "  check              Run all checks (fmt, vet, lint, sec, test)"
	@echo "  fmt                Format code"
	@echo "  vet                Vet code"
	@echo "  lint               Run linter"
	@echo "  sec                Run security check"
	@echo ""
	@echo "Release:"
	@echo "  release            Create release packages for all platforms"
	@echo "  release-current    Create release package for current platform"
	@echo ""
	@echo "Utilities:"
	@echo "  deps               Update dependencies"
	@echo "  clone-llamacpp     Clone llama.cpp repository for cross-reference"
	@echo "  model_download     Download example models"
	@echo "  install-tools      Install development tools"
	@echo "  version            Show version information"
	@echo "  help               Show this help"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION=$(VERSION)"
	@echo "  LLAMA_CPP_BUILD=$(LLAMA_CPP_BUILD)"
	@echo "  GOOS=$(GOOS)"
	@echo "  GOARCH=$(GOARCH)"
