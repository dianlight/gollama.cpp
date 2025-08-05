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
build-llamacpp-darwin-amd64: clone-llamacpp
	@echo "Building llama.cpp for macOS x86_64"
	mkdir -p $(LIB_DIR)/darwin_amd64
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-darwin-amd64 \
		-DCMAKE_OSX_ARCHITECTURES=x86_64 \
		-DCMAKE_OSX_DEPLOYMENT_TARGET=10.15 \
		-DGGML_METAL=ON \
		-DGGML_NATIVE=OFF \
		-DBUILD_SHARED_LIBS=ON \
		-DLLAMA_CURL=OFF && \
	cmake --build build-darwin-amd64 --config Release -j$$(sysctl -n hw.ncpu) && \
	cp build-darwin-amd64/bin/*.dylib ../../$(LIB_DIR)/darwin_amd64/ && \
	for lib in ../../$(LIB_DIR)/darwin_amd64/*.dylib; do \
		install_name_tool -id "@rpath/$$(basename $$lib)" "$$lib"; \
		for dep in ../../$(LIB_DIR)/darwin_amd64/*.dylib; do \
			if [ "$$lib" != "$$dep" ]; then \
				install_name_tool -change "$$(otool -L $$lib | grep $$(basename $$dep) | awk '{print $$1}')" "@rpath/$$(basename $$dep)" "$$lib" 2>/dev/null || true; \
			fi; \
		done; \
	done

.PHONY: build-llamacpp-darwin-arm64
build-llamacpp-darwin-arm64: clone-llamacpp
	@echo "Building llama.cpp for macOS ARM64"
	mkdir -p $(LIB_DIR)/darwin_arm64
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-darwin-arm64 -DCMAKE_OSX_ARCHITECTURES=arm64 -DGGML_METAL=ON -DBUILD_SHARED_LIBS=ON -DLLAMA_CURL=OFF && \
	cmake --build build-darwin-arm64 --config Release -j$$(sysctl -n hw.ncpu) && \
	cp build-darwin-arm64/bin/*.dylib ../../$(LIB_DIR)/darwin_arm64/ && \
	for lib in ../../$(LIB_DIR)/darwin_arm64/*.dylib; do \
		install_name_tool -id "@rpath/$$(basename $$lib)" "$$lib"; \
		for dep in ../../$(LIB_DIR)/darwin_arm64/*.dylib; do \
			if [ "$$lib" != "$$dep" ]; then \
				install_name_tool -change "$$(otool -L $$lib | grep $$(basename $$dep) | awk '{print $$1}')" "@rpath/$$(basename $$dep)" "$$lib" 2>/dev/null || true; \
			fi; \
		done; \
	done

# Linux builds
.PHONY: build-llamacpp-linux-amd64
build-llamacpp-linux-amd64: clone-llamacpp
	@echo "Building llama.cpp for Linux x86_64"
	@echo "Checking for GPU SDK availability..."
	@if [ -d "/usr/local/cuda" ] || [ -d "/opt/cuda" ] || command -v nvcc >/dev/null 2>&1; then \
		echo "CUDA SDK found - building with CUDA support"; \
		GPU_FLAGS="-DGGML_CUDA=ON"; \
	elif [ -d "/opt/rocm" ] || [ -d "/usr/local/rocm" ] || command -v hipcc >/dev/null 2>&1; then \
		echo "ROCm/HIP SDK found - building with AMD GPU support"; \
		GPU_FLAGS="-DGGML_HIP=ON"; \
	else \
		echo "No GPU SDK found - building CPU-only version"; \
		GPU_FLAGS="-DGGML_CUDA=OFF -DGGML_HIP=OFF"; \
	fi; \
	mkdir -p $(LIB_DIR)/linux_amd64; \
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-linux-amd64 $$GPU_FLAGS -DBUILD_SHARED_LIBS=ON -DLLAMA_CURL=OFF && \
	cmake --build build-linux-amd64 --config Release -j$$(nproc) && \
	cp build-linux-amd64/bin/lib*.so ../../$(LIB_DIR)/linux_amd64/ && \
	for lib in ../../$(LIB_DIR)/linux_amd64/*.so; do \
		patchelf --set-soname "$$(basename $$lib)" "$$lib"; \
		for dep in ../../$(LIB_DIR)/linux_amd64/*.so; do \
			if [ "$$lib" != "$$dep" ]; then \
				patchelf --replace-needed "$$(basename $$dep)" "$$(basename $$dep)" "$$lib" 2>/dev/null || true; \
			fi; \
		done; \
	done

.PHONY: build-llamacpp-linux-arm64
build-llamacpp-linux-arm64: clone-llamacpp
	@echo "Building llama.cpp for Linux ARM64"
	@echo "Checking for GPU SDK availability..."
	@if [ -d "/usr/local/cuda" ] || [ -d "/opt/cuda" ] || command -v nvcc >/dev/null 2>&1; then \
		echo "CUDA SDK found - building with CUDA support"; \
		GPU_FLAGS="-DGGML_CUDA=ON"; \
	elif [ -d "/opt/rocm" ] || [ -d "/usr/local/rocm" ] || command -v hipcc >/dev/null 2>&1; then \
		echo "ROCm/HIP SDK found - building with AMD GPU support"; \
		GPU_FLAGS="-DGGML_HIP=ON"; \
	else \
		echo "No GPU SDK found - building CPU-only version"; \
		GPU_FLAGS="-DGGML_CUDA=OFF -DGGML_HIP=OFF"; \
	fi; \
	mkdir -p $(LIB_DIR)/linux_arm64; \
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-linux-arm64 -DCMAKE_SYSTEM_PROCESSOR=aarch64 $$GPU_FLAGS -DBUILD_SHARED_LIBS=ON -DLLAMA_CURL=OFF && \
	cmake --build build-linux-arm64 --config Release -j$$(nproc) && \
	cp build-linux-arm64/bin/lib*.so ../../$(LIB_DIR)/linux_arm64/ && \
	for lib in ../../$(LIB_DIR)/linux_arm64/*.so; do \
		patchelf --set-soname "$$(basename $$lib)" "$$lib"; \
		for dep in ../../$(LIB_DIR)/linux_arm64/*.so; do \
			if [ "$$lib" != "$$dep" ]; then \
				patchelf --replace-needed "$$(basename $$dep)" "$$(basename $$dep)" "$$lib" 2>/dev/null || true; \
			fi; \
		done; \
	done

# Windows builds (require cross-compilation setup)
.PHONY: build-llamacpp-windows-amd64
build-llamacpp-windows-amd64: clone-llamacpp
	@echo "Building llama.cpp for Windows x86_64"
	mkdir -p $(LIB_DIR)/windows_amd64
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-windows-amd64 -DGGML_CUDA=ON -DBUILD_SHARED_LIBS=ON -DLLAMA_CURL=OFF -DCMAKE_TOOLCHAIN_FILE=cmake/x64-windows-llvm.cmake && \
	cmake --build build-windows-amd64 --config Release -j$$(nproc) && \
	cp build-windows-amd64/bin/*.dll ../../$(LIB_DIR)/windows_amd64/

.PHONY: build-llamacpp-windows-arm64
build-llamacpp-windows-arm64: clone-llamacpp
	@echo "Building llama.cpp for Windows ARM64"
	mkdir -p $(LIB_DIR)/windows_arm64
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-windows-arm64 -DBUILD_SHARED_LIBS=ON -DLLAMA_CURL=OFF -DCMAKE_TOOLCHAIN_FILE=cmake/arm64-windows-llvm.cmake && \
	cmake --build build-windows-arm64 --config Release -j$$(nproc) && \
	cp build-windows-arm64/bin/*.dll ../../$(LIB_DIR)/windows_arm64/

# Build libraries with GPU support
.PHONY: build-libs-gpu
build-libs-gpu: clone-llamacpp
	@echo "Building llama.cpp libraries with GPU support"
	$(MAKE) build-llamacpp-darwin-amd64  # Metal support
	$(MAKE) build-llamacpp-darwin-arm64  # Metal support
	$(MAKE) build-llamacpp-linux-amd64   # Auto-detect CUDA/HIP support
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
		cmake -B build-$$os-$$arch-cpu -DBUILD_SHARED_LIBS=ON -DGGML_CUDA=OFF -DGGML_HIP=OFF -DGGML_METAL=OFF -DLLAMA_CURL=OFF && \
		cmake --build build-$$os-$$arch-cpu --config Release && \
		cp build-$$os-$$arch-cpu/src/lib* ../../$(LIB_DIR)/$$os\_$$arch/ 2>/dev/null || true; \
		cp build-$$os-$$arch-cpu/src/*.dll ../../$(LIB_DIR)/$$os\_$$arch/ 2>/dev/null || true; \
	done

# Force HIP/AMD GPU build for testing (requires ROCm)
.PHONY: build-llamacpp-linux-amd64-hip
build-llamacpp-linux-amd64-hip: clone-llamacpp
	@echo "Building llama.cpp for Linux x86_64 with forced HIP/AMD GPU support"
	@echo "Note: This requires ROCm SDK to be installed"
	mkdir -p $(LIB_DIR)/linux_amd64
	cd $(LLAMA_CPP_DIR) && \
	cmake -B build-linux-amd64-hip -DGGML_HIP=ON -DGGML_CUDA=OFF -DBUILD_SHARED_LIBS=ON -DLLAMA_CURL=OFF && \
	cmake --build build-linux-amd64-hip --config Release -j$$(nproc) && \
	cp build-linux-amd64-hip/bin/lib*.so ../../$(LIB_DIR)/linux_amd64/ && \
	for lib in ../../$(LIB_DIR)/linux_amd64/*.so; do \
		patchelf --set-soname "$$(basename $$lib)" "$$lib"; \
		for dep in ../../$(LIB_DIR)/linux_amd64/*.so; do \
			if [ "$$lib" != "$$dep" ]; then \
				patchelf --replace-needed "$$(basename $$dep)" "$$(basename $$dep)" "$$lib" 2>/dev/null || true; \
			fi; \
		done; \
	done

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
