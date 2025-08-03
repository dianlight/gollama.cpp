# Gollama.cpp Examples

This directory contains various examples demonstrating how to use gollama.cpp for different use cases.

## Available Examples

### 1. Simple Chat (`simple-chat/`)
A basic example showing how to generate text using a GGUF model.

**Features:**
- Text generation with customizable parameters
- Basic model loading and context creation
- Token prediction with configurable limits

**Usage:**
```bash
cd simple-chat
go run main.go -prompt "The future of AI is"
```

### 2. Embedding Generation (`embedding/`)
Demonstrates how to generate high-dimensional embedding vectors from text.

**Features:**
- Generate embeddings for single or multiple texts
- Support for different output formats (default, JSON, array)
- Automatic embedding normalization using L2 norm
- Cosine similarity matrix computation for multiple texts
- Configurable context size and thread count

**Usage:**
```bash
cd embedding
go run main.go -prompt "Hello World!"

# Multiple texts with similarity matrix
go run main.go -prompt "dog|cat|animal|car|vehicle"

# JSON output format
go run main.go -prompt "Hello World!" -output-format json

# Run the interactive demo
./demo.sh
```

## Getting Started

### Prerequisites
- Go 1.21 or later
- A GGUF model file (included: `tinyllama-1.1b-chat-v1.0.Q2_K.gguf`)

### Building and Running Examples

Each example can be built and run independently:

```bash
# Navigate to any example directory
cd simple-chat  # or embedding

# Build the example
go build

# Run with default parameters
./simple-chat   # or ./embedding

# Or run directly with go
go run main.go [options]
```

### Common Options

Most examples support these common command-line options:

- `-model string`: Path to the GGUF model file
- `-prompt string`: Input text or prompt
- `-threads int`: Number of threads to use (default: 4)
- `-ctx int`: Context size (default: 2048)
- `-verbose`: Enable verbose output

## Model Requirements

### Text Generation Examples
- Any GGUF model that supports text generation
- Models like LLaMA, Mistral, CodeLlama, etc.

### Embedding Examples
- GGUF models that support embedding generation
- Some models are text-generation only and don't provide embeddings
- Verify your model supports embeddings before using the embedding example

## Troubleshooting

### "Failed to load model"
- Check that the model path is correct
- Ensure the model file is a valid GGUF file
- Make sure you have enough RAM to load the model

### "Permission denied" when running examples
- Make sure the example binary is executable: `chmod +x example-name`
- Or use `go run main.go` instead

### Out of memory errors
- Reduce context size with `-ctx 1024` or smaller
- Use a smaller model
- Close other applications to free RAM

## Contributing

When adding new examples:

1. Create a new directory under `examples/`
2. Include a comprehensive `README.md` explaining the example
3. Add a `Makefile` with common targets (build, run, clean)
4. Provide example usage commands
5. Update this main examples README

## Related Documentation

- [Main README](../README.md) - Project overview and installation
- [Build Documentation](../docs/BUILD.md) - Building from source
- [Contributing Guidelines](../CONTRIBUTING.md) - How to contribute
