# Makefile for gollama.cpp
# Cross-platform Go bindings for llama.cpp using purego

# Version information
VERSION ?= 0.2.1
LLAMA_CPP_BUILD ?= b6862
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
	env GOOS= GOARCH= $(GO) run ./cmd/gollama-download -clean-cache

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
test: deps download-libs
	@echo "Running tests (libraries will be downloaded automatically)"
	$(GO) test -p 1 -failfast -timeout 120s -tags embedallowed_no -coverprofile=coverage.out -cover ./... && \
	$(GO) tool cover -func=coverage.out | grep total: | awk '{print "Total coverage: " $$3}'


# Test with race detection
.PHONY: test-race
test-race: deps
	@echo "Running tests with race detection"
	$(GO) test -race -v ./...

# Test library download functionality
.PHONY: test-download
test-download: deps
	@echo "Testing library download functionality"
	env GOOS= GOARCH= $(GO) run ./cmd/gollama-download -test-download

# Test GPU detection and backend functionality
.PHONY: test-gpu
test-gpu: deps
	@echo "Testing GPU detection and backend functionality"
	$(GO) test -v -run TestGpu ./...
	make detect-gpu

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
	env GOOS= GOARCH= $(GO) run ./cmd/gollama-download -download -version $(LLAMA_CPP_BUILD)

# Download libraries for all platforms (for testing)
.PHONY: download-libs-all
download-libs-all: deps
	@echo "Downloading llama.cpp libraries for all platforms"
	env GOOS= GOARCH= $(GO) run ./cmd/gollama-download -download-all -version $(LLAMA_CPP_BUILD)

# Download libraries for all platforms with parallel processing
.PHONY: download-libs-parallel
download-libs-parallel: deps
	@echo "Downloading llama.cpp libraries for all platforms (parallel)"
	env GOOS= GOARCH= $(GO) run ./cmd/gollama-download -download-all -version $(LLAMA_CPP_BUILD) -checksum

# Download libraries for specific platforms
.PHONY: download-libs-platforms
download-libs-platforms: deps
	@echo "Downloading llama.cpp libraries for specific platforms"
	env GOOS= GOARCH= $(GO) run ./cmd/gollama-download -platforms "linux/amd64,darwin/arm64,windows/amd64" -version $(LLAMA_CPP_BUILD) -checksum

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
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

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
	$(GO) run github.com/securego/gosec/v2/cmd/gosec@latest -exclude=G103,G104,G115,G304 -severity=medium ./...

# Check everything
.PHONY: check
check: fmt vet lint sec test

# GPU Detection Logic
.PHONY: detect-gpu
detect-gpu:
	@echo "Detecting available GPU backends..."
	@echo "Platform: $(GOOS)/$(GOARCH)"
	@if command -v nvcc >/dev/null 2>&1; then \
		echo "✅ CUDA detected (nvcc found)"; \
		nvcc --version 2>/dev/null | head -1 || echo "   Version info not available"; \
	else \
		echo "❌ CUDA not detected"; \
	fi
	@if command -v hipconfig >/dev/null 2>&1; then \
		echo "✅ HIP/ROCm detected (hipconfig found)"; \
		hipconfig --version 2>/dev/null || echo "   Version info not available"; \
	else \
		echo "❌ HIP/ROCm not detected"; \
	fi
	@if command -v vulkaninfo >/dev/null 2>&1; then \
		echo "✅ Vulkan detected (vulkaninfo found)"; \
		vulkaninfo --summary 2>/dev/null | head -5 || echo "   Summary not available"; \
	else \
		echo "❌ Vulkan not detected"; \
	fi
	@if command -v clinfo >/dev/null 2>&1; then \
		echo "✅ OpenCL detected (clinfo found)"; \
		clinfo --list 2>/dev/null | head -10 || echo "   Device list not available"; \
	else \
		echo "❌ OpenCL not detected"; \
	fi
	@if command -v sycl-ls >/dev/null 2>&1; then \
		echo "✅ SYCL detected (sycl-ls found)"; \
		sycl-ls 2>/dev/null || echo "   Device list not available"; \
	else \
		echo "❌ SYCL not detected"; \
	fi
	@if [ "$(GOOS)" = "darwin" ]; then \
		if system_profiler SPDisplaysDataType 2>/dev/null | grep -q "Metal"; then \
			echo "✅ Metal detected"; \
		else \
			echo "❌ Metal not detected"; \
		fi \
	fi

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
	@echo "Copying hf.sh script if different or missing"
	@if [ ! -f "scripts/hf.sh" ] || ! cmp -s "$(LLAMA_CPP_DIR)/scripts/hf.sh" "scripts/hf.sh"; then \
		echo "Copying hf.sh from llama.cpp/scripts/"; \
		cp "$(LLAMA_CPP_DIR)/scripts/hf.sh" "scripts/hf.sh"; \
		chmod +x "scripts/hf.sh"; \
		echo "hf.sh script updated"; \
	else \
		echo "hf.sh script is already up to date"; \
	fi

# Update hf.sh script from llama.cpp repository
.PHONY: update-hf-script
update-hf-script: clone-llamacpp
	@echo "Forcing update of hf.sh script from llama.cpp"
	@if [ -f "$(LLAMA_CPP_DIR)/scripts/hf.sh" ]; then \
		cp "$(LLAMA_CPP_DIR)/scripts/hf.sh" "scripts/hf.sh"; \
		chmod +x "scripts/hf.sh"; \
		echo "hf.sh script updated from llama.cpp build $(LLAMA_CPP_BUILD)"; \
	else \
		echo "Error: hf.sh script not found in llama.cpp repository"; \
		exit 1; \
	fi


# Automated tag and release
.PHONY: tag-release
tag-release:
	@echo "Starting automated tag and release process..."
	
	# Check if we're on the main branch
	@current_branch=$$(git rev-parse --abbrev-ref HEAD); \
	if [ "$$current_branch" != "main" ]; then \
		echo "Error: tag-release can only be run from the main branch. Current branch: $$current_branch"; \
		exit 1; \
	fi
	
	# Check if current branch is up to date with origin
	@echo "Checking if main branch is up to date with origin..."
	@git fetch origin main >/dev/null 2>&1; \
	local_commit=$$(git rev-parse HEAD); \
	remote_commit=$$(git rev-parse origin/main); \
	if [ "$$local_commit" != "$$remote_commit" ]; then \
		echo "Error: Local main branch is not up to date with origin/main"; \
		echo "Local:  $$local_commit"; \
		echo "Remote: $$remote_commit"; \
		echo "Please pull latest changes: git pull origin main"; \
		exit 1; \
	fi
	
	# Update CHANGELOG.md for release
	@echo "Updating CHANGELOG.md for release $(FULL_VERSION)..."
	@bash scripts/update-changelog.sh "$(FULL_VERSION)" "release"
	@echo "CHANGELOG.md updated successfully"
	
	# Commit the changelog update
	@echo "Committing CHANGELOG.md update..."
	@git add CHANGELOG.md
	@git commit -m ":rocket: chore(release): update CHANGELOG.md for $(FULL_VERSION)"
	@git push origin main
	@echo "CHANGELOG.md committed successfully"
	
	# Check if current version tag exists
	@tag_name="$(FULL_VERSION)"; \
	echo "Checking if tag $$tag_name exists..."; \
	if git tag -l | grep -q "^$$tag_name$$"; then \
		echo "Tag $$tag_name already exists"; \
		\
		# Check if GitHub release exists for this tag \
		echo "Checking if GitHub release exists for tag $$tag_name..."; \
		if command -v gh >/dev/null 2>&1; then \
			if gh release view $$tag_name >/dev/null 2>&1; then \
				echo "GitHub release already exists for tag $$tag_name"; \
				echo "Moving tag to current HEAD and updating release..."; \
				git tag -d $$tag_name; \
				git push origin :refs/tags/$$tag_name; \
				git tag $$tag_name HEAD; \
				git push origin $$tag_name; \
				echo "Tag $$tag_name moved to current HEAD and pushed"; \
			else \
				echo "Tag exists but no GitHub release found. Creating release..."; \
				git push origin $$tag_name; \
				echo "Tag $$tag_name pushed to origin"; \
			fi; \
		else \
			echo "Warning: GitHub CLI (gh) not found. Cannot check for existing releases."; \
			echo "Moving tag to current HEAD..."; \
			git tag -d $$tag_name; \
			git push origin :refs/tags/$$tag_name; \
			git tag $$tag_name HEAD; \
			git push origin $$tag_name; \
			echo "Tag $$tag_name moved to current HEAD and pushed"; \
		fi; \
	else \
		echo "Tag $$tag_name does not exist. Creating and pushing..."; \
		git tag $$tag_name HEAD; \
		git push origin $$tag_name; \
		echo "Tag $$tag_name created and pushed"; \
	fi
	
	# Wait for GitHub Actions to process the tag and create the release
	@echo "Waiting for GitHub Actions to process the tag and create release artifacts..."
	@echo "You can monitor the progress at: https://github.com/$$(git config --get remote.origin.url | sed 's/.*github.com[:/]\([^/]*\/[^/]*\)\.git/\1/')/actions"
	
	# Increment patch version for next development cycle
	@echo "Incrementing patch version for next development cycle..."
	@bash scripts/increment-version.sh patch
	@echo "Version incremented successfully (no git commit performed)"
	
	# Add new [Unreleased] section to CHANGELOG.md
	@echo "Adding new [Unreleased] section to CHANGELOG.md..."
	@new_version=$$(echo $(VERSION) | awk -F. '{print $$1"."$$2"."$$3+1}'); \
	new_full_version="v$$new_version-llamacpp.$(LLAMA_CPP_BUILD)"; \
	bash scripts/update-changelog.sh "$$new_full_version" "unreleased"
	@echo "[Unreleased] section added successfully"
	
	@echo "Tag and release process completed successfully!"
	@echo "Released version: $(FULL_VERSION)"
	@echo "Next development version: v$$(echo $(VERSION) | awk -F. '{print $$1"."$$2"."$$3+1}')-llamacpp.$(LLAMA_CPP_BUILD)"

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

# Download model file using hf.sh script
.PHONY: model_download
model_download:
	@echo "Downloading models using hf.sh script"
	@mkdir -p models
	@if [ ! -f "scripts/hf.sh" ]; then \
		echo "Error: hf.sh script not found. Run 'make clone-llamacpp' first."; \
		exit 1; \
	fi
	@if [ ! -f "models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf" ]; then \
		echo "Downloading TinyLlama model using hf.sh..."; \
		bash scripts/hf.sh --repo TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF --file tinyllama-1.1b-chat-v1.0.Q2_K.gguf --outdir models; \
		echo "TinyLlama model downloaded successfully"; \
	else \
		echo "TinyLlama model already exists in models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"; \
	fi
	@if [ ! -f "models/gritlm-7b_q4_1.gguf" ]; then \
		echo "Downloading GritLM model using hf.sh..."; \
		bash scripts/hf.sh --repo cohesionet/GritLM-7B_gguf --file gritlm-7b_q4_1.gguf --outdir models; \
		echo "GritLM model downloaded successfully"; \
	else \
		echo "GritLM model already exists in models/gritlm-7b_q4_1.gguf"; \
	fi

# Roadmap management
.PHONY: roadmap-update
roadmap-update:
	@echo "Updating ROADMAP.md last updated date"
	@bash scripts/roadmap-update.sh update-date

.PHONY: roadmap-validate
roadmap-validate:
	@echo "Validating ROADMAP.md format and content"
	@bash scripts/roadmap-update.sh validate

.PHONY: roadmap-scan
roadmap-scan:
	@echo "Scanning for potential roadmap items in code"
	@bash scripts/roadmap-update.sh scan-todos

.PHONY: roadmap-scan-missing
roadmap-scan-missing:
	@echo "Scanning for code that depends on missing llama.cpp functions"
	@bash scripts/roadmap-update.sh scan-missing

.PHONY: roadmap-scan-purego
roadmap-scan-purego:
	@echo "Scanning for code that depends on purego struct support"
	@bash scripts/roadmap-update.sh scan-purego

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
	@echo "  tag-release        Automated tag and release (main branch only)"
	@echo ""
	@echo "Utilities:"
	@echo "  deps               Update dependencies"
	@echo "  clone-llamacpp     Clone llama.cpp repository for cross-reference"
	@echo "  update-hf-script   Update hf.sh script from llama.cpp repository"
	@echo "  model_download     Download example models using hf.sh script"
	@echo "  install-tools      Install development tools"
	@echo "  roadmap-update     Update ROADMAP.md last updated date"
	@echo "  roadmap-validate   Validate ROADMAP.md format and content"
	@echo "  roadmap-scan       Scan for potential roadmap items in code"
	@echo "  roadmap-scan-missing  Scan for code depending on missing llama.cpp functions"
	@echo "  roadmap-scan-purego   Scan for code depending on purego struct support"
	@echo "  version            Show version information"
	@echo "  help               Show this help"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION=$(VERSION)"
	@echo "  LLAMA_CPP_BUILD=$(LLAMA_CPP_BUILD)"
	@echo "  GOOS=$(GOOS)"
	@echo "  GOARCH=$(GOARCH)"
