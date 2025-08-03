# Gollama.cpp Evaluation Callback Example Makefile

.PHONY: all build clean run test help
.PHONY: run-basic run-simulation run-verbose run-minimal run-progress
.PHONY: demo-operations demo-performance demo-memory demo-debugging
.PHONY: test-cpu test-gpu test-threading test-context-sizes
.PHONY: benchmark-ops benchmark-throughput benchmark-memory
.DEFAULT_GOAL := help

# Configuration
MODEL_PATH := ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf
BINARY_NAME := eval-callback

# Build configuration
GO_FLAGS := -v
BUILD_FLAGS := $(GO_FLAGS)

## Build Commands

all: build ## Build the eval-callback example

build: ## Build the eval-callback executable
	@echo "Building eval-callback example..."
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	go clean
	@echo "Clean complete"

test: ## Run Go tests
	@echo "Running tests..."
	go test -v ./...

## Basic Run Commands

run: build ## Run with default settings
	@echo "Running eval-callback with default settings..."
	./$(BINARY_NAME) -model $(MODEL_PATH)

run-simulation: build ## Run simulation-only mode (no model loading)
	@echo "Running eval-callback in simulation mode..."
	./$(BINARY_NAME) -simulate-only

run-basic: build ## Run with basic model evaluation
	@echo "Running basic model evaluation..."
	./$(BINARY_NAME) -model $(MODEL_PATH) -prompt "Hello world"

run-verbose: build ## Run with verbose tensor data printing
	@echo "Running with verbose tensor data output..."
	./$(BINARY_NAME) -simulate-only -print-tensor-data -max-logged-ops 20

run-minimal: build ## Run with minimal logging
	@echo "Running with minimal logging..."
	./$(BINARY_NAME) -simulate-only -enable-logging=false -enable-progress

run-progress: build ## Run showing only progress updates
	@echo "Running with progress updates only..."
	./$(BINARY_NAME) -simulate-only -enable-logging=false -enable-progress

## Demo Commands

demo-operations: build ## Demonstrate operation logging
	@echo "=== Demonstration: Operation Logging ==="
	@echo "This shows detailed tensor operation information..."
	./$(BINARY_NAME) -simulate-only -prompt "Neural networks are" -max-logged-ops 30

demo-performance: build ## Demonstrate performance monitoring
	@echo "=== Demonstration: Performance Monitoring ==="
	@echo "This shows operation timing and throughput metrics..."
	./$(BINARY_NAME) -simulate-only -prompt "The future of artificial intelligence" -enable-progress

demo-memory: build ## Demonstrate memory location tracking  
	@echo "=== Demonstration: Memory Location Tracking ==="
	@echo "This shows CPU vs GPU memory placement..."
	./$(BINARY_NAME) -simulate-only -prompt "Machine learning models" -max-logged-ops 25

demo-debugging: build ## Demonstrate debugging capabilities
	@echo "=== Demonstration: Debugging Capabilities ==="
	@echo "This shows how to debug model inference..."
	./$(BINARY_NAME) -simulate-only -prompt "Deep learning architectures" -print-tensor-data -max-logged-ops 15

## Test Commands

test-cpu: build ## Test CPU-only mode
	@echo "Testing CPU-only evaluation..."
	./$(BINARY_NAME) -simulate-only -threads 1 -prompt "CPU evaluation test"

test-gpu: build ## Test with GPU acceleration simulation
	@echo "Testing GPU acceleration simulation..."
	./$(BINARY_NAME) -simulate-only -threads 8 -prompt "GPU acceleration test"

test-threading: build ## Test different thread counts
	@echo "Testing various thread configurations..."
	@for threads in 1 2 4 8; do \
		echo "Testing with $$threads threads..."; \
		./$(BINARY_NAME) -simulate-only -threads $$threads -prompt "Threading test $$threads" -enable-logging=false; \
	done

test-context-sizes: build ## Test different context sizes
	@echo "Testing various context sizes..."
	@for ctx in 256 512 1024 2048; do \
		echo "Testing with context size $$ctx..."; \
		./$(BINARY_NAME) -simulate-only -ctx $$ctx -prompt "Context test $$ctx" -enable-logging=false; \
	done

## Benchmark Commands

benchmark-ops: build ## Benchmark operation processing
	@echo "=== Benchmark: Operation Processing ==="
	@echo "Testing operation processing speed..."
	./$(BINARY_NAME) -simulate-only -prompt "Benchmark operation processing with longer text" -max-logged-ops 100

benchmark-throughput: build ## Benchmark data throughput
	@echo "=== Benchmark: Data Throughput ==="
	@echo "Testing data processing throughput..."
	./$(BINARY_NAME) -simulate-only -prompt "Throughput benchmark with extended prompt text for evaluation" -enable-progress

benchmark-memory: build ## Benchmark memory operations
	@echo "=== Benchmark: Memory Operations ==="
	@echo "Testing memory operation performance..."
	./$(BINARY_NAME) -simulate-only -prompt "Memory benchmark evaluation" -print-tensor-data -max-logged-ops 50

## Interactive Demos

demo-interactive: build ## Run interactive demonstration
	@echo "=== Interactive Evaluation Callback Demo ==="
	@echo "Enter prompts to see callback simulation (Ctrl+C to exit)"
	@while true; do \
		read -p "Prompt: " prompt; \
		if [ -n "$$prompt" ]; then \
			./$(BINARY_NAME) -simulate-only -prompt "$$prompt" -max-logged-ops 20; \
		fi; \
	done

demo-comparison: build ## Compare different configurations
	@echo "=== Configuration Comparison Demo ==="
	@echo "1. Minimal logging:"
	./$(BINARY_NAME) -simulate-only -prompt "Comparison test" -enable-logging=false
	@echo ""
	@echo "2. Full logging:"
	./$(BINARY_NAME) -simulate-only -prompt "Comparison test" -max-logged-ops 15
	@echo ""
	@echo "3. With tensor data:"
	./$(BINARY_NAME) -simulate-only -prompt "Comparison test" -print-tensor-data -max-logged-ops 10

## Model Testing (requires actual model file)

test-model: build ## Test with actual model (requires model file)
	@if [ -f "$(MODEL_PATH)" ]; then \
		echo "Testing with actual model..."; \
		./$(BINARY_NAME) -model $(MODEL_PATH) -prompt "Hello world" -ctx 512; \
	else \
		echo "Model file not found: $(MODEL_PATH)"; \
		echo "Please ensure the model file exists or run 'make run-simulation'"; \
	fi

test-model-performance: build ## Test model performance (requires model file)
	@if [ -f "$(MODEL_PATH)" ]; then \
		echo "Testing model performance..."; \
		./$(BINARY_NAME) -model $(MODEL_PATH) -prompt "Performance test with actual model evaluation" -ctx 1024; \
	else \
		echo "Model file not found: $(MODEL_PATH)"; \
		echo "Please ensure the model file exists or run 'make benchmark-ops'"; \
	fi

## Documentation

docs: ## Generate documentation
	@echo "Generating documentation..."
	@echo "README.md - Main documentation"
	@echo "Makefile - This makefile with all commands"
	@echo "main.go - Source code with detailed comments"

## Help

help: ## Show this help message
	@echo "Gollama.cpp Evaluation Callback Example"
	@echo "======================================="
	@echo ""
	@echo "This example demonstrates evaluation callbacks for monitoring"
	@echo "tensor operations and debugging model inference."
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Quick start:"
	@echo "  make run-simulation  # Run without loading model"
	@echo "  make demo-operations # See operation logging demo"
	@echo "  make test-model      # Test with actual model (if available)"
	@echo ""
	@echo "Model file expected at: $(MODEL_PATH)"

## Advanced Usage

advanced-demo: build ## Run advanced demonstration suite
	@echo "=== Advanced Evaluation Callback Demonstration ==="
	@echo ""
	@echo "1. Basic operation logging:"
	./$(BINARY_NAME) -simulate-only -prompt "Advanced demo" -max-logged-ops 10
	@echo ""
	@echo "2. Performance monitoring:"
	./$(BINARY_NAME) -simulate-only -prompt "Performance monitoring" -enable-progress -enable-logging=false
	@echo ""
	@echo "3. Memory location tracking:"
	./$(BINARY_NAME) -simulate-only -prompt "Memory tracking" -max-logged-ops 15
	@echo ""
	@echo "4. Tensor data inspection:"
	./$(BINARY_NAME) -simulate-only -prompt "Tensor data" -print-tensor-data -max-logged-ops 5

stress-test: build ## Run stress test with various configurations
	@echo "=== Stress Testing Evaluation Callbacks ==="
	@for ops in 10 50 100 200; do \
		echo "Testing with $$ops operations..."; \
		./$(BINARY_NAME) -simulate-only -prompt "Stress test $$ops operations" -max-logged-ops $$ops -enable-progress=false; \
	done

validate: build test ## Validate the example works correctly
	@echo "=== Validating Evaluation Callback Example ==="
	@echo "1. Testing simulation mode..."
	./$(BINARY_NAME) -simulate-only -prompt "Validation test" -max-logged-ops 5 > /dev/null
	@echo "✓ Simulation mode working"
	@echo "2. Testing progress reporting..."
	./$(BINARY_NAME) -simulate-only -prompt "Progress test" -enable-logging=false > /dev/null
	@echo "✓ Progress reporting working"
	@echo "3. Testing tensor data printing..."
	./$(BINARY_NAME) -simulate-only -prompt "Tensor test" -print-tensor-data -max-logged-ops 3 > /dev/null
	@echo "✓ Tensor data printing working"
	@echo "All validation tests passed!"
