# Makefile for gollama.cpp
# Cross-platform Go bindings for llama.cpp using purego

# Version information
VERSION ?= 1.0.0
LLAMA_CPP_BUILD ?= b6076
FULL_VERSION = v$(VERSION)-llamacpp.$(LLAMA_CPP_BUILD)

# Go configuration
GO ?= go
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Build directories
BUILD_DIR = build
LIB_DIR = libs
DIST_DIR = dist
EXAMPLES_DIR = examples

# Platform-specific configurations
PLATFORMS = darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

# llama.cpp build configurations
LLAMA_CPP_REPO = https://github.com/ggml-org/llama.cpp
LLAMA_CPP_DIR = $(BUILD_DIR)/llama.cpp

# Library names per platform
LIB_darwin_amd64 = libllama.dylib
LIB_darwin_arm64 = libllama.dylib  
LIB_linux_amd64 = libllama.so
LIB_linux_arm64 = libllama.so
LIB_windows_amd64 = llama.dll
LIB_windows_arm64 = llama.dll

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
	rm -rf $(LIB_DIR)

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

# Test
.PHONY: test
test: deps build-llamacpp-current
	@echo "Running tests"
	$(GO) test -v ./...

# Test with race detection
.PHONY: test-race
test-race: deps build-llamacpp-current
	@echo "Running tests with race detection"
	$(GO) test -race -v ./...

# Benchmark
.PHONY: bench
bench: deps build-llamacpp-current
	@echo "Running benchmarks"
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
		gosec ./...; \
	else \
		echo "gosec not found, skipping security check"; \
	fi

# Check everything
.PHONY: check
check: fmt vet lint sec test

# Clone llama.cpp repository
.PHONY: clone-llamacpp
clone-llamacpp:
	@if [ ! -d "$(LLAMA_CPP_DIR)" ]; then \
		echo "Cloning llama.cpp"; \
		mkdir -p $(BUILD_DIR); \
		git clone $(LLAMA_CPP_REPO) $(LLAMA_CPP_DIR); \
	fi
	@echo "Checking out build $(LLAMA_CPP_BUILD)"
	cd $(LLAMA_CPP_DIR) && git fetch && git checkout $(LLAMA_CPP_BUILD)

# Build llama.cpp for current platform
.PHONY: build-llamacpp-current
build-llamacpp-current: clone-llamacpp
	@echo "Building llama.cpp for $(GOOS)/$(GOARCH)"
	mkdir -p $(LIB_DIR)/$(GOOS)_$(GOARCH)
	$(MAKE) build-llamacpp-$(GOOS)-$(GOARCH)

# Build llama.cpp for all platforms
.PHONY: build-llamacpp-all
build-llamacpp-all: clone-llamacpp
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		echo "Building llama.cpp for $$os/$$arch"; \
		$(MAKE) build-llamacpp-$$os-$$arch; \
	done

# macOS builds
.PHONY: build-llamacpp-darwin-amd64
build-llamacpp-darwin-amd64:
	@echo "Building llama.cpp for macOS x86_64"
	mkdir -p $(LIB_DIR)/darwin_amd64
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-darwin-amd64 -DCMAKE_OSX_ARCHITECTURES=x86_64 -DGGML_METAL=ON -DBUILD_SHARED_LIBS=ON && \
	cmake --build build-darwin-amd64 --config Release -j$$(sysctl -n hw.ncpu) && \
	cp build-darwin-amd64/bin/libllama.dylib ../../$(LIB_DIR)/darwin_amd64/

.PHONY: build-llamacpp-darwin-arm64
build-llamacpp-darwin-arm64:
	@echo "Building llama.cpp for macOS ARM64"
	mkdir -p $(LIB_DIR)/darwin_arm64
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-darwin-arm64 -DCMAKE_OSX_ARCHITECTURES=arm64 -DGGML_METAL=ON -DBUILD_SHARED_LIBS=ON && \
	cmake --build build-darwin-arm64 --config Release -j$$(sysctl -n hw.ncpu) && \
	cp build-darwin-arm64/bin/libllama.dylib ../../$(LIB_DIR)/darwin_arm64/

# Linux builds
.PHONY: build-llamacpp-linux-amd64
build-llamacpp-linux-amd64:
	@echo "Building llama.cpp for Linux x86_64"
	mkdir -p $(LIB_DIR)/linux_amd64
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-linux-amd64 -DGGML_CUDA=ON -DBUILD_SHARED_LIBS=ON && \
	cmake --build build-linux-amd64 --config Release -j$$(nproc) && \
	cp build-linux-amd64/bin/libllama.so ../../$(LIB_DIR)/linux_amd64/

.PHONY: build-llamacpp-linux-arm64
build-llamacpp-linux-arm64:
	@echo "Building llama.cpp for Linux ARM64"
	mkdir -p $(LIB_DIR)/linux_arm64
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-linux-arm64 -DCMAKE_SYSTEM_PROCESSOR=aarch64 -DBUILD_SHARED_LIBS=ON && \
	cmake --build build-linux-arm64 --config Release -j$$(nproc) && \
	cp build-linux-arm64/bin/libllama.so ../../$(LIB_DIR)/linux_arm64/

# Windows builds (require cross-compilation setup)
.PHONY: build-llamacpp-windows-amd64
build-llamacpp-windows-amd64:
	@echo "Building llama.cpp for Windows x86_64"
	mkdir -p $(LIB_DIR)/windows_amd64
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-windows-amd64 -DGGML_CUDA=ON -DBUILD_SHARED_LIBS=ON -DCMAKE_TOOLCHAIN_FILE=cmake/x86_64-w64-mingw32.cmake && \
	cmake --build build-windows-amd64 --config Release -j$$(nproc) && \
	cp build-windows-amd64/bin/llama.dll ../../$(LIB_DIR)/windows_amd64/

.PHONY: build-llamacpp-windows-arm64
build-llamacpp-windows-arm64:
	@echo "Building llama.cpp for Windows ARM64"
	mkdir -p $(LIB_DIR)/windows_arm64
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-windows-arm64 -DBUILD_SHARED_LIBS=ON -DCMAKE_TOOLCHAIN_FILE=cmake/aarch64-w64-mingw32.cmake && \
	cmake --build build-windows-arm64 --config Release -j$$(nproc) && \
	cp build-windows-arm64/bin/llama.dll ../../$(LIB_DIR)/windows_arm64/

# Build libraries with GPU support
.PHONY: build-libs-gpu
build-libs-gpu: clone-llamacpp
	@echo "Building llama.cpp libraries with GPU support"
	$(MAKE) build-llamacpp-darwin-amd64  # Metal support
	$(MAKE) build-llamacpp-darwin-arm64  # Metal support
	$(MAKE) build-llamacpp-linux-amd64   # CUDA support
	$(MAKE) build-llamacpp-windows-amd64 # CUDA support

# Build CPU-only libraries
.PHONY: build-libs-cpu
build-libs-cpu: clone-llamacpp
	@echo "Building llama.cpp libraries (CPU only)"
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		echo "Building CPU-only library for $$os/$$arch"; \
		mkdir -p $(LIB_DIR)/$$os\_$$arch; \
		cd $(LLAMA_CPP_DIR) && \
		cmake -B build-$$os-$$arch-cpu -DBUILD_SHARED_LIBS=ON -DGGML_CUDA=OFF -DGGML_METAL=OFF && \
		cmake --build build-$$os-$$arch-cpu --config Release && \
		cp build-$$os-$$arch-cpu/src/lib* ../../$(LIB_DIR)/$$os\_$$arch/ 2>/dev/null || true; \
		cp build-$$os-$$arch-cpu/src/*.dll ../../$(LIB_DIR)/$$os\_$$arch/ 2>/dev/null || true; \
	done

# Package releases
.PHONY: release
release: clean build-all build-llamacpp-all
	@echo "Creating release packages"
	mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		echo "Packaging $$os/$$arch"; \
		pkg_name="gollama.cpp-$(FULL_VERSION)-$$os-$$arch"; \
		mkdir -p $(DIST_DIR)/$$pkg_name; \
		cp -r $(BUILD_DIR)/$$os\_$$arch/* $(DIST_DIR)/$$pkg_name/ 2>/dev/null || true; \
		cp -r $(LIB_DIR)/$$os\_$$arch/* $(DIST_DIR)/$$pkg_name/ 2>/dev/null || true; \
		cp README.md LICENSE CHANGELOG.md $(DIST_DIR)/$$pkg_name/; \
		cd $(DIST_DIR) && zip -r $$pkg_name.zip $$pkg_name && rm -rf $$pkg_name; \
	done

# Quick release for current platform
.PHONY: release-current
release-current: clean build build-llamacpp-current
	@echo "Creating release package for $(GOOS)/$(GOARCH)"
	mkdir -p $(DIST_DIR)
	pkg_name="gollama.cpp-$(FULL_VERSION)-$(GOOS)-$(GOARCH)"
	mkdir -p $(DIST_DIR)/$$pkg_name
	cp -r $(BUILD_DIR)/$(GOOS)_$(GOARCH)/* $(DIST_DIR)/$$pkg_name/ 2>/dev/null || true
	cp -r $(LIB_DIR)/$(GOOS)_$(GOARCH)/* $(DIST_DIR)/$$pkg_name/ 2>/dev/null || true
	cp README.md LICENSE CHANGELOG.md $(DIST_DIR)/$$pkg_name/
	cd $(DIST_DIR) && zip -r $$pkg_name.zip $$pkg_name && rm -rf $$pkg_name
	@echo "Release package created: $(DIST_DIR)/$$pkg_name.zip"

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools"
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Download model file
.PHONY: model_download
model_download:
	@echo "Downloading tinyllama-1.1b-chat-v1.0.Q2_K.gguf model"
	@mkdir -p models
	@if [ ! -f "models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf" ]; then \
		echo "Downloading model from Hugging Face..."; \
		curl -L -o models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf \
			"https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"; \
		echo "Model downloaded successfully"; \
	else \
		echo "Model already exists in models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"; \
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
	@echo "  test               Run tests"
	@echo "  test-race          Run tests with race detection"
	@echo "  bench              Run benchmarks"
	@echo "  clean              Clean all build artifacts"
	@echo ""
	@echo "Library building:"
	@echo "  build-llamacpp-current    Build llama.cpp for current platform"
	@echo "  build-llamacpp-all        Build llama.cpp for all platforms"
	@echo "  build-libs-gpu            Build libraries with GPU support"
	@echo "  build-libs-cpu            Build CPU-only libraries"
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
	@echo "  model_download     Download tinyllama-1.1b-chat-v1.0.Q2_K.gguf model"
	@echo "  install-tools      Install development tools"
	@echo "  version            Show version information"
	@echo "  help               Show this help"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION=$(VERSION)"
	@echo "  LLAMA_CPP_BUILD=$(LLAMA_CPP_BUILD)"
	@echo "  GOOS=$(GOOS)"
	@echo "  GOARCH=$(GOARCH)"
