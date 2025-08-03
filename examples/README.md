# Gollama.cpp Examples

This directory contains various examples demonstrating how to use gollama.cpp for different use cases.

## Available Examples

### 1. Simple Chat (`simple-chat/`)
A comprehensive example showing how to generate text using a GGUF model.

**Features:**
- Text generation and completion with configurable parameters
- Real-time token generation with streaming output
- System information display (GPU support, memory mapping, etc.)
- Detailed progress logging and error handling
- Support for various text types: creative writing, technical explanations, conversations
- Performance optimization with configurable threading

**Usage:**
```bash
cd simple-chat
go run main.go -prompt "The future of AI is"

# Creative writing
go run main.go -prompt "Once upon a time" -n-predict 150

# Technical explanation
go run main.go -prompt "How does machine learning work?" -n-predict 100

# Run the interactive demo
./demo.sh

# Use Makefile shortcuts
make creative    # Creative writing demo
make technical   # Technical explanation demo
make conversation # Conversation starter demo
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

### 3. Speculative Decoding (`speculative/`)
Advanced example demonstrating speculative decoding for accelerated text generation.

**Features:**
- Dual-model speculative decoding with separate target and draft models
- Same-model demonstration mode for understanding the algorithm
- Configurable draft length for performance tuning
- Temperature sampling support with detailed statistics
- Verbose mode for observing the draft/verify process
- Performance analysis showing acceptance rates and speedup

**Usage:**
```bash
cd speculative
go run main.go -prompt "The future of AI is"

# With different models for real speedup
go run main.go -model large.gguf -draft-model small.gguf -prompt "Your prompt"

# Demonstration with verbose output
go run main.go -prompt "Machine learning" -n-draft 8 -verbose

# Run the interactive demo
./demo.sh

# Use Makefile shortcuts
make demo              # Full demonstration
make draft-comparison  # Compare different draft lengths
make temperature-demo  # Temperature sampling demo
```

### 3. Speculative Decoding (`speculative/`)
Advanced example demonstrating speculative decoding for accelerated text generation.

**Features:**
- Dual-model speculative decoding with separate target and draft models
- Same-model demonstration mode for understanding the algorithm
- Configurable draft length for performance tuning
- Temperature sampling support with detailed statistics
- Verbose mode for observing the draft/verify process
- Performance analysis showing acceptance rates and speedup

**Usage:**
```bash
cd speculative
go run main.go -prompt "The future of AI is"

# With different models for real speedup
go run main.go -model large.gguf -draft-model small.gguf -prompt "Your prompt"

# Demonstration with verbose output
go run main.go -prompt "Machine learning" -n-draft 8 -verbose

# Run the interactive demo
./demo.sh

# Use Makefile shortcuts
make demo              # Full demonstration
make draft-comparison  # Compare different draft lengths
make temperature-demo  # Temperature sampling demo
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
