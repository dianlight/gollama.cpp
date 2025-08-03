# Makefile for Gollama.cpp Embedding Example

.PHONY: build run clean help test

# Default target
all: build

# Build the embedding example
build:
	@echo "Building embedding example..."
	go build -o embedding main.go
	@echo "Built successfully: ./embedding"

# Run with default parameters
run: build
	@echo "Running embedding example with default parameters..."
	./embedding

# Run with multiple texts to show similarity matrix
demo: build
	@echo "Running embedding example with multiple texts..."
	./embedding -prompt "dog|cat|animal|car|vehicle" -verbose

# Run with JSON output
json: build
	@echo "Running embedding example with JSON output..."
	./embedding -prompt "Hello World!" -output-format json

# Test compilation
test:
	@echo "Testing compilation..."
	go build -o /tmp/embedding_test main.go
	@rm -f /tmp/embedding_test
	@echo "Compilation test passed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f embedding
	@echo "Clean complete!"

# Show help
help:
	@echo "Gollama.cpp Embedding Example Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build    - Build the embedding example binary"
	@echo "  run      - Build and run with default parameters"
	@echo "  demo     - Build and run with multiple texts for similarity demo"
	@echo "  json     - Build and run with JSON output format"
	@echo "  test     - Test that the code compiles successfully"
	@echo "  clean    - Remove build artifacts"
	@echo "  help     - Show this help message"
	@echo ""
	@echo "Manual usage examples:"
	@echo "  ./embedding -prompt \"Hello World!\""
	@echo "  ./embedding -prompt \"dog|cat|animal\" -output-format json"
	@echo "  ./embedding -model /path/to/model.gguf -prompt \"Text\" -verbose"
