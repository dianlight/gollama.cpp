# Makefile for Gollama.cpp Batched Generation Example

# Configuration
BINARY_NAME = batched
GO_FILES = main.go
MODEL_PATH = ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf

# Default target
.PHONY: all
all: build

# Build the example
.PHONY: build
build:
	@echo "Building batched generation example..."
	go build -o $(BINARY_NAME) $(GO_FILES)
	@echo "✅ Build complete: ./$(BINARY_NAME)"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	@echo "✅ Clean complete"

# Show help
.PHONY: help
help:
	@echo "Gollama.cpp Batched Generation Example"
	@echo ""
	@echo "Available targets:"
	@echo "  build              Build the batched example"
	@echo "  clean              Remove build artifacts"
	@echo "  demo               Run comprehensive demo"
	@echo "  basic-demo         Run basic batched generation"
	@echo "  creative-demo      Run creative writing demo"
	@echo "  technical-demo     Run technical content demo"
	@echo "  parameter-demo     Run parameter comparison demo"
	@echo "  performance-demo   Run performance analysis"
	@echo "  help               Show this help message"
	@echo ""
	@echo "Quick examples:"
	@echo "  make basic-demo    # Generate 4 sequences"
	@echo "  make creative-demo # Creative writing variations"
	@echo "  make verbose-demo  # Detailed processing info"

# Basic demo - 4 sequences with default settings
.PHONY: basic-demo
basic-demo: build
	@echo "🚀 Running Basic Batched Generation Demo"
	@echo ""
	@echo "Generating 4 sequences with default prompt..."
	./$(BINARY_NAME) -prompt "Hello my name is" -n-parallel 4 -n-predict 25

# Creative writing demo
.PHONY: creative-demo
creative-demo: build
	@echo "✍️  Running Creative Writing Demo"
	@echo ""
	@echo "Generating story variations with high temperature..."
	./$(BINARY_NAME) -prompt "In a world where magic exists," -n-parallel 3 -n-predict 40 -temperature 0.9

# Technical content demo
.PHONY: technical-demo
technical-demo: build
	@echo "💻 Running Technical Content Demo"
	@echo ""
	@echo "Generating technical explanations with low temperature..."
	./$(BINARY_NAME) -prompt "// Function to calculate" -n-parallel 3 -n-predict 30 -temperature 0.3

# Parameter comparison demo
.PHONY: parameter-demo
parameter-demo: build
	@echo "⚙️  Running Parameter Comparison Demo"
	@echo ""
	@echo "=== High Temperature (Creative) ==="
	./$(BINARY_NAME) -prompt "Innovation means" -temperature 1.1 -n-parallel 3 -n-predict 20
	@echo ""
	@echo "=== Low Temperature (Focused) ==="
	./$(BINARY_NAME) -prompt "Innovation means" -temperature 0.2 -n-parallel 3 -n-predict 20

# Performance analysis demo
.PHONY: performance-demo
performance-demo: build
	@echo "📊 Running Performance Analysis Demo"
	@echo ""
	@echo "=== 1 Sequence ==="
	time ./$(BINARY_NAME) -prompt "Performance test" -n-parallel 1 -n-predict 20
	@echo ""
	@echo "=== 2 Sequences ==="
	time ./$(BINARY_NAME) -prompt "Performance test" -n-parallel 2 -n-predict 20
	@echo ""
	@echo "=== 4 Sequences ==="
	time ./$(BINARY_NAME) -prompt "Performance test" -n-parallel 4 -n-predict 20

# Verbose demo with detailed output
.PHONY: verbose-demo
verbose-demo: build
	@echo "🔍 Running Verbose Demo"
	@echo ""
	@echo "Generating sequences with detailed processing information..."
	./$(BINARY_NAME) -prompt "Machine learning is" -n-parallel 2 -n-predict 15 -verbose

# Interactive demo
.PHONY: interactive-demo
interactive-demo: build
	@echo "🎮 Interactive Batched Generation Demo"
	@echo ""
	@echo "Try these prompts with different parameters:"
	@echo ""
	@echo "Creative Writing:"
	@echo "  ./$(BINARY_NAME) -prompt \"Once upon a time\" -n-parallel 3 -temperature 0.9"
	@echo ""
	@echo "Technical Content:"
	@echo "  ./$(BINARY_NAME) -prompt \"The algorithm works by\" -n-parallel 2 -temperature 0.3"
	@echo ""
	@echo "Conversation Starters:"
	@echo "  ./$(BINARY_NAME) -prompt \"What I find interesting is\" -n-parallel 4 -temperature 0.7"
	@echo ""
	@echo "Run with -verbose to see detailed processing steps!"

# Comprehensive demo showing all features
.PHONY: demo
demo: build
	@echo "🌟 Comprehensive Batched Generation Demo"
	@echo ""
	@echo "This demo shows multiple aspects of batched text generation:"
	@echo ""
	
	@echo "📝 1. Basic Batched Generation"
	@echo "   Generating multiple sequences from the same prompt..."
	./$(BINARY_NAME) -prompt "The future of technology" -n-parallel 3 -n-predict 25
	@echo ""
	
	@echo "🎨 2. Creative Writing Variations"
	@echo "   Using higher temperature for creative diversity..."
	./$(BINARY_NAME) -prompt "In a distant galaxy" -n-parallel 2 -n-predict 30 -temperature 1.0
	@echo ""
	
	@echo "🔬 3. Technical Content Generation"
	@echo "   Using lower temperature for focused output..."
	./$(BINARY_NAME) -prompt "// This function" -n-parallel 2 -n-predict 20 -temperature 0.3
	@echo ""
	
	@echo "⚡ 4. Performance Comparison"
	@echo "   Comparing different sequence counts..."
	@echo "   Single sequence:"
	time -p ./$(BINARY_NAME) -prompt "Speed test" -n-parallel 1 -n-predict 15 2>&1 | grep real || true
	@echo "   Multiple sequences:"
	time -p ./$(BINARY_NAME) -prompt "Speed test" -n-parallel 3 -n-predict 15 2>&1 | grep real || true
	@echo ""
	
	@echo "🎯 5. Parameter Impact Demo"
	@echo "   Same prompt, different sampling parameters..."
	@echo "   Conservative sampling:"
	./$(BINARY_NAME) -prompt "Artificial intelligence" -top-k 10 -top-p 0.5 -n-parallel 2 -n-predict 15
	@echo "   Diverse sampling:"
	./$(BINARY_NAME) -prompt "Artificial intelligence" -top-k 50 -top-p 0.9 -n-parallel 2 -n-predict 15
	@echo ""
	
	@echo "🎉 Demo Complete!"
	@echo ""
	@echo "💡 Tips for Further Exploration:"
	@echo "   • Try different prompts with 'make creative-demo'"
	@echo "   • Use 'make verbose-demo' to see internal processing"
	@echo "   • Run 'make performance-demo' to analyze throughput"
	@echo "   • Experiment with different temperature values (0.1-1.5)"
	@echo "   • Test various sequence counts (1-8) to find optimal settings"

# Comparison with other examples
.PHONY: comparison-demo
comparison-demo: build
	@echo "🔄 Comparison with Other Generation Methods"
	@echo ""
	@echo "=== Batched Generation (Multiple Sequences) ==="
	./$(BINARY_NAME) -prompt "The best way to learn" -n-parallel 3 -n-predict 20
	@echo ""
	@echo "=== Single Sequence (for comparison) ==="
	./$(BINARY_NAME) -prompt "The best way to learn" -n-parallel 1 -n-predict 20
	@echo ""
	@echo "Notice how batched generation provides multiple perspectives on the same topic!"

# Test with different models (if available)
.PHONY: model-test
model-test: build
	@echo "🧪 Testing with Available Models"
	@echo ""
	@if [ -f "$(MODEL_PATH)" ]; then \
		echo "✅ Using model: $(MODEL_PATH)"; \
		./$(BINARY_NAME) -model "$(MODEL_PATH)" -prompt "Model test" -n-parallel 2 -n-predict 10 -verbose; \
	else \
		echo "❌ Model not found: $(MODEL_PATH)"; \
		echo "Please ensure you have a GGUF model file in the models directory."; \
	fi

# Quick test to verify functionality
.PHONY: test
test: build
	@echo "🧪 Quick Functionality Test"
	@echo ""
	./$(BINARY_NAME) -prompt "Test prompt" -n-parallel 2 -n-predict 5 -verbose

# Install dependencies
.PHONY: deps
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	@echo "✅ Dependencies installed"

# Run with custom prompt (use PROMPT variable)
.PHONY: custom
custom: build
	@echo "🎯 Running with custom prompt..."
	@if [ -z "$(PROMPT)" ]; then \
		echo "Usage: make custom PROMPT=\"Your prompt here\""; \
		echo "Example: make custom PROMPT=\"The secret to happiness\""; \
	else \
		./$(BINARY_NAME) -prompt "$(PROMPT)" -n-parallel 3 -n-predict 30; \
	fi

# Show example usage patterns
.PHONY: examples
examples:
	@echo "📋 Example Usage Patterns"
	@echo ""
	@echo "Creative Writing:"
	@echo "  make creative-demo"
	@echo "  ./$(BINARY_NAME) -prompt \"Once upon a time\" -temperature 0.9 -n-parallel 4"
	@echo ""
	@echo "Technical Documentation:"
	@echo "  make technical-demo"  
	@echo "  ./$(BINARY_NAME) -prompt \"This API endpoint\" -temperature 0.3 -n-parallel 3"
	@echo ""
	@echo "Brainstorming:"
	@echo "  ./$(BINARY_NAME) -prompt \"Solutions to climate change\" -temperature 0.8 -n-parallel 5"
	@echo ""
	@echo "Code Comments:"
	@echo "  ./$(BINARY_NAME) -prompt \"// This algorithm\" -temperature 0.2 -n-parallel 3"
	@echo ""
	@echo "Story Variations:"
	@echo "  ./$(BINARY_NAME) -prompt \"The hero discovered\" -temperature 1.0 -n-parallel 4"

# Check if model file exists
.PHONY: check-model
check-model:
	@echo "🔍 Checking for model file..."
	@if [ -f "$(MODEL_PATH)" ]; then \
		echo "✅ Model found: $(MODEL_PATH)"; \
		ls -lh "$(MODEL_PATH)"; \
	else \
		echo "❌ Model not found: $(MODEL_PATH)"; \
		echo "Please download a GGUF model file to the models directory."; \
		echo "You can find models at: https://huggingface.co/models"; \
	fi
