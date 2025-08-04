# GritLM Dual-Purpose Model Example Makefile

# Go settings
GO_BIN := go
BINARY_NAME := gritlm
MAIN_FILE := main.go

# Model settings
MODEL_NAME := gritlm-7b_q4_1.gguf
MODEL_REPO := cohesionet/GritLM-7B_gguf
MODEL_DIR := ../../models
MODEL_PATH := $(MODEL_DIR)/$(MODEL_NAME)
HF_SCRIPT := ../../build/llama.cpp/scripts/hf.sh

# Build settings
BUILD_FLAGS := -ldflags="-s -w"
VERBOSE := false

# Default target
.PHONY: all
all: build

# Help target
.PHONY: help
help:
	@echo "GritLM Dual-Purpose Model Example"
	@echo "=================================="
	@echo ""
	@echo "Available targets:"
	@echo "  build              - Build the gritlm example binary"
	@echo "  run                - Run the example with default model"
	@echo "  run-verbose        - Run with verbose output"
	@echo "  model_download     - Download the GritLM-7B model (~4.2GB)"
	@echo "  model_info         - Show model information"
	@echo "  test               - Run basic functionality tests"
	@echo "  demo               - Run interactive demo"
	@echo "  clean              - Clean build artifacts"
	@echo "  clean-model        - Remove downloaded model"
	@echo "  distclean          - Clean everything including models"
	@echo ""
	@echo "Model Information:"
	@echo "  Model: $(MODEL_REPO)/$(MODEL_NAME)"
	@echo "  Size:  ~4.2GB (Q4_1 quantized)"
	@echo "  Path:  $(MODEL_PATH)"
	@echo ""
	@echo "Usage Examples:"
	@echo "  make build && make model_download && make run"
	@echo "  make demo"
	@echo "  ./$(BINARY_NAME) /path/to/custom/model.gguf"

# Build the binary
.PHONY: build
build:
	@echo "Building GritLM example..."
	$(GO_BIN) build $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Built: ./$(BINARY_NAME)"

# Download the GritLM model using hf.sh script
.PHONY: model_download
model_download:
	@echo "Downloading GritLM-7B model..."
	@echo "Model: $(MODEL_REPO)/$(MODEL_NAME)"
	@echo "Size: ~4.2GB (this may take several minutes)"
	@echo ""
	@if [ ! -f "$(HF_SCRIPT)" ]; then \
		echo "Error: hf.sh script not found at $(HF_SCRIPT)"; \
		echo "Please ensure llama.cpp is properly built in ../../build/llama.cpp/"; \
		exit 1; \
	fi
	@mkdir -p $(MODEL_DIR)
	@if [ -f "$(MODEL_PATH)" ]; then \
		echo "Model already exists at $(MODEL_PATH)"; \
		echo "Use 'make clean-model' to remove and re-download"; \
	else \
		echo "Downloading $(MODEL_NAME)..."; \
		$(HF_SCRIPT) --repo $(MODEL_REPO) --file $(MODEL_NAME) --outdir $(MODEL_DIR); \
		if [ -f "$(MODEL_PATH)" ]; then \
			echo "Successfully downloaded $(MODEL_NAME)"; \
			echo "Model path: $(MODEL_PATH)"; \
		else \
			echo "Failed to download model"; \
			exit 1; \
		fi \
	fi

# Check if model exists and show info
.PHONY: model_info
model_info:
	@echo "GritLM Model Information"
	@echo "========================"
	@echo "Repository: $(MODEL_REPO)"
	@echo "Filename:   $(MODEL_NAME)"
	@echo "Local path: $(MODEL_PATH)"
	@echo "Expected size: ~4.2GB"
	@echo ""
	@if [ -f "$(MODEL_PATH)" ]; then \
		echo "Status: ✓ Model is available"; \
		echo "Size: $$(du -h "$(MODEL_PATH)" | cut -f1)"; \
		echo "Modified: $$(stat -f "%Sm" "$(MODEL_PATH)")"; \
	else \
		echo "Status: ✗ Model not found"; \
		echo "Run 'make model_download' to download"; \
	fi

# Run the example
.PHONY: run
run: build
	@if [ ! -f "$(MODEL_PATH)" ]; then \
		echo "Model not found. Downloading..."; \
		$(MAKE) model_download; \
	fi
	@echo "Running GritLM example..."
	./$(BINARY_NAME) $(MODEL_PATH)

# Run with verbose output
.PHONY: run-verbose
run-verbose: build
	@if [ ! -f "$(MODEL_PATH)" ]; then \
		echo "Model not found. Downloading..."; \
		$(MAKE) model_download; \
	fi
	@echo "Running GritLM example (verbose)..."
	GOLLAMA_VERBOSE=1 ./$(BINARY_NAME) $(MODEL_PATH)

# Interactive demo
.PHONY: demo
demo: build
	@echo "GritLM Interactive Demo"
	@echo "======================="
	@echo ""
	@echo "This demo showcases GritLM's dual-purpose capabilities:"
	@echo "1. Text embedding generation for semantic search"
	@echo "2. Text generation setup for conversational AI"
	@echo "3. RAG (Retrieval-Augmented Generation) pipeline"
	@echo ""
	@if [ ! -f "$(MODEL_PATH)" ]; then \
		echo "Model not found. Would you like to download it? (y/N)"; \
		read -r answer; \
		if [ "$$answer" = "y" ] || [ "$$answer" = "Y" ]; then \
			$(MAKE) model_download; \
		else \
			echo "Cannot run demo without model. Exiting."; \
			exit 1; \
		fi \
	fi
	@echo "Starting demo with model: $(MODEL_PATH)"
	@echo "Press Ctrl+C to exit during execution"
	@echo ""
	./$(BINARY_NAME) $(MODEL_PATH)

# Basic functionality test
.PHONY: test
test: build
	@echo "Testing GritLM example..."
	@if [ ! -f "$(MODEL_PATH)" ]; then \
		echo "Model required for testing. Downloading..."; \
		$(MAKE) model_download; \
	fi
	@echo "Running basic functionality test..."
	@timeout 120 ./$(BINARY_NAME) $(MODEL_PATH) > test_output.log 2>&1 || true
	@if grep -q "Generated.*embedding" test_output.log; then \
		echo "✓ Embedding generation test passed"; \
	else \
		echo "✗ Embedding generation test failed"; \
	fi
	@if grep -q "Best match:" test_output.log; then \
		echo "✓ Similarity search test passed"; \
	else \
		echo "✗ Similarity search test failed"; \
	fi
	@if grep -q "Generation prompt" test_output.log; then \
		echo "✓ Generation setup test passed"; \
	else \
		echo "✗ Generation setup test failed"; \
	fi
	@rm -f test_output.log
	@echo "Basic tests completed"

# Development targets
.PHONY: dev-build
dev-build:
	@echo "Development build (with race detection)..."
	$(GO_BIN) build -race -o $(BINARY_NAME)-dev $(MAIN_FILE)

.PHONY: dev-run
dev-run: dev-build
	@if [ ! -f "$(MODEL_PATH)" ]; then \
		$(MAKE) model_download; \
	fi
	./$(BINARY_NAME)-dev $(MODEL_PATH)

# Format and lint
.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	$(GO_BIN) fmt ./...

.PHONY: vet
vet:
	@echo "Running Go vet..."
	$(GO_BIN) vet ./...

.PHONY: lint
lint: fmt vet
	@echo "Linting completed"

# Documentation
.PHONY: docs
docs:
	@echo "GritLM Example Documentation"
	@echo "============================="
	@echo ""
	@cat README.md

# Benchmark (if model is available)
.PHONY: benchmark
benchmark: build
	@if [ ! -f "$(MODEL_PATH)" ]; then \
		echo "Model required for benchmarking. Run 'make model_download' first."; \
		exit 1; \
	fi
	@echo "Running embedding benchmark..."
	@time ./$(BINARY_NAME) $(MODEL_PATH) > /dev/null
	@echo "Benchmark completed"

# Performance profiling
.PHONY: profile
profile: build
	@if [ ! -f "$(MODEL_PATH)" ]; then \
		echo "Model required for profiling. Run 'make model_download' first."; \
		exit 1; \
	fi
	@echo "Running with CPU profiling..."
	@echo "Profile data will be saved to cpu.prof"
	$(GO_BIN) build -o $(BINARY_NAME)-profile $(MAIN_FILE)
	./$(BINARY_NAME)-profile -cpuprofile=cpu.prof $(MODEL_PATH)
	@echo "View profile with: go tool pprof cpu.prof"

# Clean targets
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-dev $(BINARY_NAME)-profile
	@rm -f test_output.log cpu.prof mem.prof
	@echo "Clean completed"

.PHONY: clean-model
clean-model:
	@echo "Removing downloaded model..."
	@rm -f $(MODEL_PATH)
	@echo "Model removed: $(MODEL_PATH)"

.PHONY: distclean
distclean: clean clean-model
	@echo "Deep clean completed"

# Utility targets
.PHONY: size
size:
	@if [ -f "$(BINARY_NAME)" ]; then \
		echo "Binary size: $$(du -h $(BINARY_NAME) | cut -f1)"; \
	else \
		echo "Binary not found. Run 'make build' first."; \
	fi
	@if [ -f "$(MODEL_PATH)" ]; then \
		echo "Model size:  $$(du -h $(MODEL_PATH) | cut -f1)"; \
	fi

.PHONY: deps
deps:
	@echo "Checking dependencies..."
	@$(GO_BIN) mod tidy
	@$(GO_BIN) mod verify
	@echo "Dependencies verified"

# Check if tools are available
.PHONY: check-env
check-env:
	@echo "Environment Check"
	@echo "=================="
	@which $(GO_BIN) > /dev/null && echo "✓ Go found: $$($(GO_BIN) version)" || echo "✗ Go not found"
	@[ -f "$(HF_SCRIPT)" ] && echo "✓ hf.sh script found" || echo "✗ hf.sh script not found at $(HF_SCRIPT)"
	@[ -d "$(MODEL_DIR)" ] && echo "✓ Model directory exists" || echo "ℹ Model directory will be created"
	@echo ""

# Show current status
.PHONY: status
status: check-env
	@echo "Build Status"
	@echo "============"
	@[ -f "$(BINARY_NAME)" ] && echo "✓ Binary built" || echo "✗ Binary not built"
	@[ -f "$(MODEL_PATH)" ] && echo "✓ Model downloaded" || echo "✗ Model not downloaded"
	@echo ""
	@echo "Ready to run: $$([ -f "$(BINARY_NAME)" ] && [ -f "$(MODEL_PATH)" ] && echo "Yes" || echo "No")"

# Quick start (build + download + run)
.PHONY: quickstart
quickstart:
	@echo "GritLM Quick Start"
	@echo "=================="
	@echo "This will build the example, download the model, and run a demo"
	@echo ""
	$(MAKE) build
	$(MAKE) model_download  
	$(MAKE) run

# Show version info
.PHONY: version
version:
	@echo "GritLM Example Version Information"
	@echo "=================================="
	@$(GO_BIN) version
	@echo "Model: $(MODEL_REPO)/$(MODEL_NAME)"
	@if [ -f "$(BINARY_NAME)" ]; then \
		echo "Binary: $(BINARY_NAME) (built)"; \
	else \
		echo "Binary: $(BINARY_NAME) (not built)"; \
	fi
