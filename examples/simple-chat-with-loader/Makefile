# Makefile for Gollama.cpp Simple Chat with Library Loader Example

.PHONY: build run clean help test demo interactive

# Default target
all: build

# Build the simple chat with loader example
build:
	@echo "Building simple chat with library loader example..."
	go build -o simple-chat-with-loader main.go
	@echo "Built successfully: ./simple-chat-with-loader"

# Run with default parameters
run: build
	@echo "Running simple chat with library loader example..."
	./simple-chat-with-loader

# Run with custom prompt
prompt: build
	@echo "Running simple chat with loader and custom prompt..."
	./simple-chat-with-loader -prompt "Once upon a time" -n-predict 100

# Run with long text generation
long: build
	@echo "Running simple chat with loader for longer text generation..."
	./simple-chat-with-loader -prompt "The history of artificial intelligence" -n-predict 200 -ctx 4096

# Run interactive demo with various prompts
demo: build
	@echo "Running simple chat with loader demo..."
	./simple-chat-with-loader -prompt "Explain quantum computing" -n-predict 75
	@echo ""
	@echo "--- Next prompt ---"
	@echo ""
	./simple-chat-with-loader -prompt "Write a short story about robots" -n-predict 100

# Run creative writing example
creative: build
	@echo "Running creative writing example with library loader..."
	./simple-chat-with-loader -prompt "In a world where time travel is possible," -n-predict 150

# Run technical explanation example
technical: build
	@echo "Running technical explanation example with library loader..."
	./simple-chat-with-loader -prompt "How does machine learning work?" -n-predict 100

# Run conversation starter
conversation: build
	@echo "Running conversation starter example with library loader..."
	./simple-chat-with-loader -prompt "Hello! I'm an AI assistant. How can I help you today?" -n-predict 80

# Test library loading specifically
test-loader: build
	@echo "Testing library loading functionality..."
	./simple-chat-with-loader -prompt "Testing library loader" -n-predict 10

# Test compilation
test:
	@echo "Testing compilation..."
	go build -o /tmp/simple_chat_loader_test main.go
	@rm -f /tmp/simple_chat_loader_test
	@echo "Compilation test passed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f simple-chat-with-loader
	@echo "Clean complete!"

# Show help
help:
	@echo "Gollama.cpp Simple Chat with Library Loader Example Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build the simple chat with loader example binary"
	@echo "  run          - Build and run with default parameters"
	@echo "  prompt       - Build and run with custom prompt"
	@echo "  long         - Build and run with longer text generation"
	@echo "  demo         - Build and run demonstration with multiple prompts"
	@echo "  creative     - Build and run creative writing example"
	@echo "  technical    - Build and run technical explanation example"
	@echo "  conversation - Build and run conversation starter example"
	@echo "  test-loader  - Build and run with minimal output to test loader"
	@echo "  test         - Test that the code compiles successfully"
	@echo "  clean        - Remove build artifacts"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "Manual usage examples:"
	@echo "  ./simple-chat-with-loader -prompt \"Your prompt here\""
	@echo "  ./simple-chat-with-loader -prompt \"Story beginning\" -n-predict 150"
	@echo "  ./simple-chat-with-loader -model /path/to/model.gguf -prompt \"Text\" -threads 8"
	@echo ""
	@echo "This example demonstrates both:"
	@echo "  1. Dynamic library loading with embedded llama.cpp libraries"
	@echo "  2. Text generation using the loaded libraries"
