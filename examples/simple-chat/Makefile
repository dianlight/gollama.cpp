# Makefile for Gollama.cpp Simple Chat Example

.PHONY: build run clean help test demo interactive

# Default target
all: build

# Build the simple chat example
build:
	@echo "Building simple chat example..."
	go build -o simple-chat main.go
	@echo "Built successfully: ./simple-chat"

# Run with default parameters
run: build
	@echo "Running simple chat example with default parameters..."
	./simple-chat

# Run with custom prompt
prompt: build
	@echo "Running simple chat with custom prompt..."
	./simple-chat -prompt "Once upon a time" -n-predict 100

# Run with long text generation
long: build
	@echo "Running simple chat with longer text generation..."
	./simple-chat -prompt "The history of artificial intelligence" -n-predict 200 -ctx 4096

# Run interactive demo with various prompts
demo: build
	@echo "Running simple chat demo with various prompts..."
	./simple-chat -prompt "Explain quantum computing" -n-predict 75
	@echo ""
	@echo "--- Next prompt ---"
	@echo ""
	./simple-chat -prompt "Write a short story about robots" -n-predict 100

# Run creative writing example
creative: build
	@echo "Running creative writing example..."
	./simple-chat -prompt "In a world where time travel is possible," -n-predict 150

# Run technical explanation example
technical: build
	@echo "Running technical explanation example..."
	./simple-chat -prompt "How does machine learning work?" -n-predict 100

# Run conversation starter
conversation: build
	@echo "Running conversation starter example..."
	./simple-chat -prompt "Hello! I'm an AI assistant. How can I help you today?" -n-predict 80

# Test compilation
test:
	@echo "Testing compilation..."
	go build -o /tmp/simple_chat_test main.go
	@rm -f /tmp/simple_chat_test
	@echo "Compilation test passed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f simple-chat
	@echo "Clean complete!"

# Show help
help:
	@echo "Gollama.cpp Simple Chat Example Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build the simple chat example binary"
	@echo "  run          - Build and run with default parameters"
	@echo "  prompt       - Build and run with custom prompt"
	@echo "  long         - Build and run with longer text generation"
	@echo "  demo         - Build and run demonstration with multiple prompts"
	@echo "  creative     - Build and run creative writing example"
	@echo "  technical    - Build and run technical explanation example"
	@echo "  conversation - Build and run conversation starter example"
	@echo "  test         - Test that the code compiles successfully"
	@echo "  clean        - Remove build artifacts"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "Manual usage examples:"
	@echo "  ./simple-chat -prompt \"Your prompt here\""
	@echo "  ./simple-chat -prompt \"Story beginning\" -n-predict 150"
	@echo "  ./simple-chat -model /path/to/model.gguf -prompt \"Text\" -threads 8"
