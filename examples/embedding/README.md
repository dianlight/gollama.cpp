# Gollama.cpp Embedding Example

This example demonstrates how to generate high-dimensional embedding vectors from text using gollama.cpp.

## Quick Start

To get started right away, run the following command, making sure to use the correct path for the model you have:

```bash
cd examples/embedding
go run main.go -model ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf -prompt "Hello World!"
```

## Features

- Generate embeddings for single or multiple texts
- Support for different output formats (default, JSON, array)
- Automatic embedding normalization using L2 norm
- Cosine similarity matrix computation for multiple texts
- Configurable context size and thread count

## Command Line Options

- `-model string`: Path to the GGUF model file (default: "../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf")
- `-prompt string`: Text to generate embeddings for, use `|` to separate multiple texts (default: "Hello World!")
- `-separator string`: Separator for multiple prompts (default: "|")
- `-normalize`: Normalize embeddings using L2 norm (default: true)
- `-threads int`: Number of threads to use (default: 4)
- `-ctx int`: Context size (default: 2048)
- `-verbose`: Enable verbose output (default: false)
- `-output-format string`: Output format: default, json, array (default: "default")

## Examples

### Single Text Embedding

```bash
go run main.go -model ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf -prompt "The quick brown fox"
```

### Multiple Text Embeddings with Similarity Matrix

```bash
go run main.go -model ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf -prompt "dog|cat|animal|car"
```

### JSON Output Format

```bash
go run main.go -model ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf -prompt "Hello World!" -output-format json
```

### Array Output Format

```bash
go run main.go -model ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf -prompt "Hello World!" -output-format array
```

### Verbose Output

```bash
go run main.go -model ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf -prompt "Hello World!" -verbose
```

## Output Formats

### Default Format
Shows the first 5 and last 5 dimensions of each embedding, along with a cosine similarity matrix for multiple prompts.

### JSON Format
Outputs embeddings in JSON format compatible with OpenAI's embedding API response format.

### Array Format
Outputs embeddings as simple arrays of floating-point numbers.

## Understanding the Output

- **Embedding dimensions**: The number of dimensions depends on your model (typically 768, 1024, 4096, etc.)
- **Normalization**: When enabled (default), all embedding vectors have unit length (L2 norm = 1)
- **Cosine similarity**: Values range from -1 to 1, where 1 means identical, 0 means orthogonal, and -1 means opposite
- **Similarity matrix**: Shows how similar each text is to every other text

## Technical Details

This example:

1. Loads a GGUF model with embedding support enabled
2. Tokenizes input text(s) using the model's tokenizer
3. Processes tokens through the model to generate embeddings
4. Optionally normalizes embeddings using L2 normalization
5. Computes cosine similarity between multiple embeddings
6. Outputs results in the requested format

## Requirements

- A GGUF model file that supports embeddings
- Go 1.21 or later
- Sufficient RAM to load the model (varies by model size)

## Troubleshooting

### "Failed to load model"
- Check that the model path is correct
- Ensure the model file is a valid GGUF file
- Make sure you have enough RAM to load the model

### "Failed to get embeddings"
- Verify that your model supports embeddings
- Some models are text generation only and don't provide embeddings
- Try with a different model known to support embeddings

### Empty or Invalid Output
- Check that your input text is not empty
- Verify that the model's tokenizer can process your text
- Enable verbose mode (`-verbose`) for more diagnostic information
