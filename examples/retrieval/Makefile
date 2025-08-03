# Makefile for Gollama.cpp Retrieval Example

.PHONY: build run clean help test demo ai-demo programming-demo sample-files

# Default target
all: build

# Build the retrieval example
build:
	@echo "Building retrieval example..."
	go build -o retrieval main.go
	@echo "Built successfully: ./retrieval"

# Create sample files for demonstration
sample-files:
	@echo "Sample files already exist:"
	@echo "  - sample_ai.txt (AI and machine learning concepts)"
	@echo "  - sample_programming.txt (programming languages)"

# Run with sample AI file
ai-demo: build sample-files
	@echo "Running retrieval demo with AI concepts file..."
	@echo "This will start an interactive session. Try queries like:"
	@echo "  - 'machine learning'"
	@echo "  - 'neural networks'"
	@echo "  - 'computer vision'"
	@echo ""
	./retrieval -context-files sample_ai.txt -chunk-size 150 -top-k 3

# Run with programming languages file
programming-demo: build sample-files
	@echo "Running retrieval demo with programming languages file..."
	@echo "This will start an interactive session. Try queries like:"
	@echo "  - 'web development'"
	@echo "  - 'systems programming'"
	@echo "  - 'mobile development'"
	@echo ""
	./retrieval -context-files sample_programming.txt -chunk-size 100 -top-k 3

# Run with both sample files
combined-demo: build sample-files
	@echo "Running retrieval demo with both sample files..."
	@echo "This will start an interactive session. Try queries like:"
	@echo "  - 'artificial intelligence'"
	@echo "  - 'programming language'"
	@echo "  - 'machine learning vs programming'"
	@echo ""
	./retrieval -context-files "sample_ai.txt,sample_programming.txt" -chunk-size 120 -top-k 4

# Run with default parameters
run: build sample-files
	@echo "Running retrieval example with sample files..."
	./retrieval -context-files "sample_ai.txt,sample_programming.txt"

# Non-interactive single query demo
single-query: build sample-files
	@echo "Running single query demos..."
	@echo ""
	@echo "=== Query: 'neural networks' ==="
	./retrieval -context-files sample_ai.txt -query "neural networks" -interactive=false -top-k 2
	@echo ""
	@echo "=== Query: 'web programming' ==="
	./retrieval -context-files sample_programming.txt -query "web programming" -interactive=false -top-k 2
	@echo ""
	@echo "=== Query: 'machine learning programming' ==="
	./retrieval -context-files "sample_ai.txt,sample_programming.txt" -query "machine learning programming" -interactive=false -top-k 3

# Verbose demo showing internal process
verbose-demo: build sample-files
	@echo "Running verbose demo to show internal process..."
	./retrieval -context-files sample_ai.txt -query "artificial intelligence" -interactive=false -verbose -top-k 3

# Different chunk size comparison
chunk-size-demo: build sample-files
	@echo "Comparing different chunk sizes..."
	@echo ""
	@echo "=== Small chunks (50 chars) ==="
	./retrieval -context-files sample_ai.txt -query "machine learning" -interactive=false -chunk-size 50 -top-k 3
	@echo ""
	@echo "=== Medium chunks (150 chars) ==="
	./retrieval -context-files sample_ai.txt -query "machine learning" -interactive=false -chunk-size 150 -top-k 3
	@echo ""
	@echo "=== Large chunks (300 chars) ==="
	./retrieval -context-files sample_ai.txt -query "machine learning" -interactive=false -chunk-size 300 -top-k 3

# Full demonstration with multiple scenarios
demo: build sample-files
	@echo "=== Gollama.cpp Retrieval Example Full Demo ==="
	@echo ""
	@echo "1. AI concepts retrieval:"
	./retrieval -context-files sample_ai.txt -query "deep learning" -interactive=false -top-k 2
	@echo ""
	@echo "2. Programming languages retrieval:"
	./retrieval -context-files sample_programming.txt -query "systems programming" -interactive=false -top-k 2
	@echo ""
	@echo "3. Cross-domain retrieval:"
	./retrieval -context-files "sample_ai.txt,sample_programming.txt" -query "intelligent software" -interactive=false -top-k 3
	@echo ""
	@echo "4. Interactive mode is also available with: make ai-demo"

# Test compilation
test:
	@echo "Testing compilation..."
	go build -o /tmp/retrieval_test main.go
	@rm -f /tmp/retrieval_test
	@echo "Compilation test passed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f retrieval
	@echo "Clean complete!"

# Show help
help:
	@echo "Gollama.cpp Retrieval Example Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build           - Build the retrieval example binary"
	@echo "  run             - Build and run with sample files (interactive)"
	@echo "  ai-demo         - Interactive demo with AI concepts file"
	@echo "  programming-demo - Interactive demo with programming languages file"
	@echo "  combined-demo   - Interactive demo with both sample files"
	@echo "  single-query    - Non-interactive demo with predefined queries"
	@echo "  verbose-demo    - Verbose demo showing internal process"
	@echo "  chunk-size-demo - Compare different chunk sizes"
	@echo "  demo            - Full demonstration with multiple scenarios"
	@echo "  sample-files    - Information about sample files"
	@echo "  test            - Test that the code compiles successfully"
	@echo "  clean           - Remove build artifacts"
	@echo "  help            - Show this help message"
	@echo ""
	@echo "Manual usage examples:"
	@echo "  ./retrieval -context-files \"file1.txt,file2.txt\""
	@echo "  ./retrieval -context-files document.txt -query \"search term\" -interactive=false"
	@echo "  ./retrieval -context-files data.txt -chunk-size 200 -top-k 5 -verbose"
	@echo ""
	@echo "Interactive mode commands:"
	@echo "  Enter any search query to find similar text chunks"
	@echo "  Type 'quit' or 'exit' to end the session"
