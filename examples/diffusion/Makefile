# Makefile for Diffusion Generation Example

# Build configuration
BINARY_NAME=diffusion
MAIN_FILE=main.go
MODEL_PATH=../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf

# Build targets
.PHONY: all build clean run demo help

all: build

build:
	@echo "Building diffusion generation example..."
	@go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Build complete: ./$(BINARY_NAME)"

run: build
	@echo "Running diffusion generation example..."
	@./$(BINARY_NAME)

# Demo targets showing different diffusion algorithms and configurations
demo: build
	@echo "🌟 Comprehensive Diffusion Generation Demo"
	@echo ""
	@echo "This demo shows different aspects of diffusion-based text generation:"
	@echo ""
	@echo "📝 1. Basic Diffusion Generation"
	@echo "   Generating text using confidence-based algorithm..."
	@./$(BINARY_NAME) -prompt "The future of technology" -diffusion-steps 8 -max-length 48
	@echo ""
	@echo "🧠 2. Algorithm Comparison"
	@echo "   Comparing different confidence algorithms..."
	@echo ""
	@echo "   Confidence-based algorithm:"
	@./$(BINARY_NAME) -prompt "Machine learning is" -diffusion-algorithm 0 -diffusion-steps 6 -max-length 32
	@echo ""
	@echo "   Entropy-based algorithm:"  
	@./$(BINARY_NAME) -prompt "Machine learning is" -diffusion-algorithm 1 -diffusion-steps 6 -max-length 32
	@echo ""
	@echo "🎨 3. Creative Generation"
	@echo "   Using higher temperature for creativity..."
	@./$(BINARY_NAME) -prompt "Once upon a time" -temperature 1.0 -diffusion-steps 10 -max-length 56
	@echo ""
	@echo "🔬 4. Verbose Analysis"
	@echo "   Detailed step-by-step generation process..."
	@./$(BINARY_NAME) -prompt "Science and innovation" -verbose -diffusion-steps 5 -max-length 40
	@echo ""
	@echo "✨ Demo complete! Try interactive mode with: make visual-demo"

visual-demo: build
	@echo "🎬 Visual Diffusion Demo"
	@echo "Watch the generation process in real-time..."
	@echo "Press Ctrl+C to stop"
	@./$(BINARY_NAME) -prompt "The evolution of artificial intelligence" -diffusion-visual -diffusion-steps 12 -max-length 72

# Algorithm comparison demos
confidence-demo: build
	@echo "🎯 Confidence-Based Algorithm Demo"
	@./$(BINARY_NAME) -prompt "Artificial intelligence will" -diffusion-algorithm 0 -verbose -diffusion-steps 8

entropy-demo: build
	@echo "📊 Entropy-Based Algorithm Demo"
	@./$(BINARY_NAME) -prompt "Artificial intelligence will" -diffusion-algorithm 1 -verbose -diffusion-steps 8

margin-demo: build
	@echo "📏 Margin-Based Algorithm Demo"
	@./$(BINARY_NAME) -prompt "Artificial intelligence will" -diffusion-algorithm 2 -verbose -diffusion-steps 8

random-demo: build
	@echo "🎲 Random Algorithm Demo (for comparison)"
	@./$(BINARY_NAME) -prompt "Artificial intelligence will" -diffusion-algorithm 3 -verbose -diffusion-steps 8

# Creative writing demos
creative: build
	@echo "🎨 Creative Writing with Diffusion"
	@./$(BINARY_NAME) -prompt "In a world where robots and humans" -temperature 0.9 -diffusion-steps 10 -max-length 64

sci-fi: build
	@echo "🚀 Science Fiction Generation"
	@./$(BINARY_NAME) -prompt "The year is 2150, and humanity has" -temperature 0.8 -diffusion-steps 12 -max-length 80

technical: build
	@echo "🔬 Technical Documentation Style"
	@./$(BINARY_NAME) -prompt "The implementation of neural networks requires" -temperature 0.6 -diffusion-steps 8 -max-length 72

# Parameter exploration demos
long-demo: build
	@echo "📏 Long-form Generation Demo"
	@./$(BINARY_NAME) -prompt "The history of computing began" -diffusion-steps 15 -max-length 96

fast-demo: build
	@echo "⚡ Fast Generation Demo"
	@./$(BINARY_NAME) -prompt "Quick summary:" -diffusion-steps 4 -max-length 32

deterministic: build
	@echo "🔒 Deterministic Generation Demo"
	@echo "Running the same prompt twice with fixed seed:"
	@echo ""
	@echo "First run:"
	@./$(BINARY_NAME) -prompt "Deterministic example" -seed 12345 -diffusion-steps 6
	@echo ""
	@echo "Second run (should be identical):"
	@./$(BINARY_NAME) -prompt "Deterministic example" -seed 12345 -diffusion-steps 6

# Benchmarking and testing
benchmark: build
	@echo "⏱️  Diffusion Generation Benchmark"
	@echo "Testing performance with different configurations..."
	@echo ""
	@echo "Small configuration:"
	@time ./$(BINARY_NAME) -prompt "Benchmark test" -diffusion-steps 5 -max-length 32 > /dev/null
	@echo ""
	@echo "Medium configuration:"
	@time ./$(BINARY_NAME) -prompt "Benchmark test" -diffusion-steps 10 -max-length 64 > /dev/null
	@echo ""
	@echo "Large configuration:"
	@time ./$(BINARY_NAME) -prompt "Benchmark test" -diffusion-steps 15 -max-length 96 > /dev/null

# Testing different models (if available)
test-models: build
	@echo "🧪 Testing Different Models"
	@if [ -f "../../models/another-model.gguf" ]; then \
		echo "Testing with alternative model..."; \
		./$(BINARY_NAME) -model "../../models/another-model.gguf" -prompt "Model comparison test"; \
	else \
		echo "Only testing with default model (no alternative models found)"; \
		./$(BINARY_NAME) -prompt "Single model test"; \
	fi

# Interactive demos
interactive: build
	@echo "🎮 Interactive Diffusion Generation"
	@echo "Enter your own prompts and see diffusion in action!"
	@echo "Default settings: 8 steps, confidence-based algorithm"
	@echo ""
	@read -p "Enter your prompt: " prompt; \
	./$(BINARY_NAME) -prompt "$$prompt" -diffusion-steps 8 -max-length 60

interactive-visual: build
	@echo "🎬 Interactive Visual Diffusion"
	@echo "Watch your prompt being generated step by step!"
	@echo ""
	@read -p "Enter your prompt: " prompt; \
	./$(BINARY_NAME) -prompt "$$prompt" -diffusion-visual -diffusion-steps 10 -max-length 64

# Educational demos
explain-algorithms: build
	@echo "🎓 Diffusion Algorithms Educational Demo"
	@echo ""
	@echo "This demo shows how different algorithms affect token selection:"
	@echo ""
	@prompt="Learning about diffusion models"; \
	echo "Testing prompt: $$prompt"; \
	echo ""; \
	echo "1. Confidence-based (uses token probability):"; \
	./$(BINARY_NAME) -prompt "$$prompt" -diffusion-algorithm 0 -diffusion-steps 6 -verbose; \
	echo ""; \
	echo "2. Entropy-based (uses distribution entropy):"; \
	./$(BINARY_NAME) -prompt "$$prompt" -diffusion-algorithm 1 -diffusion-steps 6 -verbose; \
	echo ""; \
	echo "3. Margin-based (uses probability difference):"; \
	./$(BINARY_NAME) -prompt "$$prompt" -diffusion-algorithm 2 -diffusion-steps 6 -verbose

explain-steps: build
	@echo "📚 Diffusion Steps Educational Demo"
	@echo ""
	@echo "Showing how step count affects generation quality:"
	@echo ""
	@prompt="The principles of machine learning"; \
	echo "Testing prompt: $$prompt"; \
	echo ""; \
	echo "Few steps (faster, less refined):"; \
	./$(BINARY_NAME) -prompt "$$prompt" -diffusion-steps 3 -max-length 48; \
	echo ""; \
	echo "Medium steps (balanced):"; \
	./$(BINARY_NAME) -prompt "$$prompt" -diffusion-steps 8 -max-length 48; \
	echo ""; \
	echo "Many steps (slower, more refined):"; \
	./$(BINARY_NAME) -prompt "$$prompt" -diffusion-steps 15 -max-length 48

# Utility targets
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@echo "✅ Clean complete"

check-model:
	@if [ ! -f "$(MODEL_PATH)" ]; then \
		echo "❌ Model not found: $(MODEL_PATH)"; \
		echo "Please ensure the model file exists before running examples."; \
		exit 1; \
	else \
		echo "✅ Model found: $(MODEL_PATH)"; \
	fi

help:
	@echo "Gollama.cpp Diffusion Generation Example"
	@echo "========================================="
	@echo ""
	@echo "Build targets:"
	@echo "  build         - Build the diffusion example"
	@echo "  clean         - Remove build artifacts"
	@echo "  run           - Build and run with default settings"
	@echo ""
	@echo "Demo targets:"
	@echo "  demo          - Comprehensive demo showing all features"
	@echo "  visual-demo   - Interactive visual generation demo"
	@echo ""
	@echo "Algorithm demos:"
	@echo "  confidence-demo - Confidence-based algorithm"
	@echo "  entropy-demo    - Entropy-based algorithm"
	@echo "  margin-demo     - Margin-based algorithm"
	@echo "  random-demo     - Random algorithm (comparison)"
	@echo ""
	@echo "Content demos:"
	@echo "  creative      - Creative writing generation"
	@echo "  sci-fi        - Science fiction generation"
	@echo "  technical     - Technical documentation style"
	@echo ""
	@echo "Configuration demos:"
	@echo "  long-demo     - Long-form generation"
	@echo "  fast-demo     - Fast generation"
	@echo "  deterministic - Deterministic generation with seed"
	@echo ""
	@echo "Interactive demos:"
	@echo "  interactive        - Interactive prompt input"
	@echo "  interactive-visual - Interactive with visual mode"
	@echo ""
	@echo "Educational demos:"
	@echo "  explain-algorithms - Compare different algorithms"
	@echo "  explain-steps      - Show effect of step count"
	@echo ""
	@echo "Testing:"
	@echo "  benchmark     - Performance benchmark"
	@echo "  test-models   - Test with different models"
	@echo "  check-model   - Verify model file exists"
	@echo ""
	@echo "Usage examples:"
	@echo "  make demo                    # Run comprehensive demo"
	@echo "  make visual-demo             # Watch generation in real-time"
	@echo "  make confidence-demo         # Test confidence algorithm"
	@echo "  make interactive            # Enter your own prompts"
	@echo ""
	@echo "Direct usage:"
	@echo "  ./$(BINARY_NAME) --help     # Show all command-line options"
	@echo "  ./$(BINARY_NAME) --prompt \"Your text\" --diffusion-visual"
