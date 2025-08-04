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

### 2. Simple Chat with Library Loader (`simple-chat-with-loader/`)
A comprehensive example that combines dynamic library loading with text generation capabilities.

**Features:**
- **Dynamic Library Loading**: Automatically extracts and loads embedded llama.cpp libraries
- **Library Lifecycle Management**: Proper loading and unloading of shared libraries with cleanup
- **Cross-platform Support**: Works on macOS, Linux, and Windows with platform detection
- **Embedded Library Extraction**: Pre-compiled libraries embedded in the binary for portability
- **Combined Functionality**: Full integration of library management and text generation
- **Resource Management**: Automatic cleanup of temporary files and library handles
- **Error Handling**: Comprehensive error handling for both library loading and model operations

**Usage:**
```bash
cd simple-chat-with-loader
go run main.go -prompt "The future of AI is"

# Test library loading specifically
go run main.go -prompt "Testing library loader" -n-predict 20

# Creative writing with loader
go run main.go -prompt "Once upon a time" -n-predict 150

# Run the comprehensive demo
./demo.sh

# Use Makefile shortcuts
make test-loader  # Test library loading functionality
make creative     # Creative writing demo with loader
make technical    # Technical explanation demo with loader
```

### 3. Embedding Generation (`embedding/`)
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

### 4. Speculative Decoding (`speculative/`)
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

### 5. Document Retrieval (`retrieval/`)
Comprehensive document retrieval system using embedding-based semantic search.

**Features:**
- Multi-document semantic search using embeddings
- Configurable text chunking with custom separators and sizes
- Cosine similarity ranking for relevance scoring
- Interactive query mode for exploratory research
- Cross-domain search across multiple document types
- Verbose mode showing internal processing steps
- Support for various document formats and chunking strategies

**Usage:**
```bash
cd retrieval
go run main.go -context-files "document.txt" -query "search terms"

# Multiple files with interactive mode
go run main.go -context-files "file1.txt,file2.txt"

# Custom chunking and ranking
go run main.go -context-files "doc.txt" -chunk-size 150 -top-k 5

# Run comprehensive demo
./demo.sh

# Use Makefile shortcuts
make ai-demo          # AI concepts retrieval
make programming-demo # Programming languages retrieval
make combined-demo    # Cross-domain search
```

### 6. Batched Generation (`batched/`)
Simplified demonstration of batched text generation concepts for multiple sequence generation.

**Features:**
- Conceptual demonstration of batched generation principles
- Multiple sequence generation from the same prompt
- Configurable parallel sequence count and generation parameters
- Performance statistics and timing analysis
- Educational example showing what full batched processing would entail
- Clear documentation of implementation limitations and future improvements

**Usage:**
```bash
cd batched
go run main.go -prompt "The future of technology" -n-parallel 4

# Creative writing with higher temperature
go run main.go -prompt "Once upon a time" -n-parallel 3 -temperature 1.0

# Quick batch demo
go run main.go -prompt "Hello world" -n-parallel 2 -n-predict 20

# Run the interactive demo
./demo.sh

# Use Makefile shortcuts
make demo         # Comprehensive demo
make creative     # Creative writing demo
make tech-demo    # Technology concepts demo
```

**Note:** This is a simplified demonstration that generates sequences sequentially to illustrate batched generation concepts. A full implementation would require true parallel batch processing with advanced batch management.

### 7. Diffusion Generation (`diffusion/`)
Conceptual demonstration of diffusion-based text generation principles using iterative token refinement.

**Features:**
- Iterative token refinement through multiple diffusion steps
- Multiple confidence algorithms (confidence-based, entropy-based, margin-based, random)
- Visual mode showing real-time generation progress
- Configurable transfer scheduling and step parameters
- Educational demonstration of diffusion model concepts
- Performance analysis and step-by-step verbose output

**Usage:**
```bash
cd diffusion
go run main.go -prompt "The future of AI" -diffusion-steps 10

# Visual mode with real-time progress
go run main.go -prompt "Machine learning" -diffusion-visual -diffusion-steps 8

# Compare different algorithms
go run main.go -prompt "Technology advances" -diffusion-algorithm 1 -verbose

# Run the interactive demo
./demo.sh

# Use Makefile shortcuts
make demo              # Comprehensive demo
make visual-demo       # Interactive visual generation
make confidence-demo   # Confidence-based algorithm
make entropy-demo      # Entropy-based algorithm
make interactive       # Interactive prompt input
```

**Note:** This is a conceptual demonstration of diffusion principles. A full implementation would require specialized diffusion model architectures and non-causal attention mechanisms not available in standard chat models.

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

### Text Generation Examples (simple-chat, speculative)
- Any GGUF model that supports text generation
- Models like LLaMA, Mistral, CodeLlama, etc.

### Embedding Examples (embedding, retrieval)
- GGUF models that support embedding generation
- Some models are text-generation only and don't provide embeddings
- Verify your model supports embeddings before using these examples
- Popular embedding-capable models: all-MiniLM, sentence-transformers models
- Note: Even some chat models like TinyLlama support embeddings

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
