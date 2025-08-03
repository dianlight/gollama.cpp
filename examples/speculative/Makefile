# Makefile for Gollama.cpp Speculative Decoding Example

.PHONY: build run clean help test demo fast-demo same-model different-models

# Default target
all: build

# Build the speculative decoding example
build:
	@echo "Building speculative decoding example..."
	go build -o speculative main.go
	@echo "Built successfully: ./speculative"

# Run with default parameters (same model for both target and draft)
run: build
	@echo "Running speculative decoding example with default parameters..."
	./speculative

# Run with same model for both target and draft (no acceleration but shows algorithm)
same-model: build
	@echo "Running speculative decoding with same model (demonstration mode)..."
	./speculative -prompt "The benefits of artificial intelligence include" -n-predict 150 -n-draft 8 -verbose

# Run demo with different draft lengths
draft-comparison: build
	@echo "Comparing different draft lengths..."
	@echo "=== Draft length: 3 ==="
	./speculative -prompt "Machine learning is" -n-predict 50 -n-draft 3
	@echo ""
	@echo "=== Draft length: 8 ==="
	./speculative -prompt "Machine learning is" -n-predict 50 -n-draft 8
	@echo ""
	@echo "=== Draft length: 15 ==="
	./speculative -prompt "Machine learning is" -n-predict 50 -n-draft 15

# Run with temperature sampling
temperature-demo: build
	@echo "Running with temperature sampling..."
	./speculative -prompt "Once upon a time in a distant galaxy" -n-predict 100 -temperature 0.7 -n-draft 6 -verbose

# Run creative writing demo
creative: build
	@echo "Running creative writing with speculative decoding..."
	./speculative -prompt "In the year 2050, robots and humans" -n-predict 200 -n-draft 10 -temperature 0.3

# Run technical explanation demo
technical: build
	@echo "Running technical explanation with speculative decoding..."
	./speculative -prompt "Quantum computing works by utilizing" -n-predict 150 -n-draft 8 -temperature 0.1

# Run conversation demo
conversation: build
	@echo "Running conversation with speculative decoding..."
	./speculative -prompt "Hello! I'm an AI assistant and I can help you with" -n-predict 120 -n-draft 6

# Run performance comparison (if you have two different models)
performance-demo: build
	@echo "Running performance demonstration..."
	@echo "Note: For real speedup, use a smaller draft model with -draft-model flag"
	./speculative -prompt "The history of artificial intelligence" -n-predict 200 -n-draft 12 -verbose

# Full demonstration with various scenarios
demo: build
	@echo "=== Speculative Decoding Full Demonstration ==="
	@echo ""
	@echo "1. Basic speculative decoding:"
	./speculative -prompt "The future of technology" -n-predict 80 -n-draft 5
	@echo ""
	@echo "2. Verbose mode (shows draft/verify process):"
	./speculative -prompt "Artificial intelligence" -n-predict 40 -n-draft 4 -verbose
	@echo ""
	@echo "3. Longer generation with more drafts:"
	./speculative -prompt "Climate change solutions include" -n-predict 120 -n-draft 10

# Test compilation
test:
	@echo "Testing compilation..."
	go build -o /tmp/speculative_test main.go
	@rm -f /tmp/speculative_test
	@echo "Compilation test passed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f speculative
	@echo "Clean complete!"

# Show help
help:
	@echo "Gollama.cpp Speculative Decoding Example Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build             - Build the speculative decoding example binary"
	@echo "  run               - Build and run with default parameters"
	@echo "  same-model        - Run with same model (demonstration mode)"
	@echo "  draft-comparison  - Compare different draft lengths"
	@echo "  temperature-demo  - Run with temperature sampling"
	@echo "  creative          - Run creative writing demo"
	@echo "  technical         - Run technical explanation demo"
	@echo "  conversation      - Run conversation demo"
	@echo "  performance-demo  - Run performance demonstration"
	@echo "  demo              - Run full demonstration with multiple scenarios"
	@echo "  test              - Test that the code compiles successfully"
	@echo "  clean             - Remove build artifacts"
	@echo "  help              - Show this help message"
	@echo ""
	@echo "Manual usage examples:"
	@echo "  ./speculative -prompt \"Your prompt here\""
	@echo "  ./speculative -draft-model /path/to/fast_model.gguf -prompt \"Text\""
	@echo "  ./speculative -prompt \"Story\" -n-draft 10 -temperature 0.5 -verbose"
	@echo ""
	@echo "For real speedup, use a smaller/faster model as the draft model:"
	@echo "  ./speculative -model large_model.gguf -draft-model small_model.gguf"
