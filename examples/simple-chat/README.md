# Gollama.cpp Simple Chat Example

This example demonstrates basic text generation using gollama.cpp. It shows how to load a GGUF model, tokenize input text, and generate new tokens in a conversational or completion style.

## Quick Start

To get started right away, run the following command, making sure to use the correct path for the model you have:

```bash
cd examples/simple-chat
go run main.go -prompt "The future of AI is"
```

**Note:** The example will automatically download the required llama.cpp libraries if they are not found on your system. This ensures the example works out of the box without manual setup.

## Features

- Basic text generation and completion
- Configurable model parameters
- Support for different context sizes and thread counts
- Greedy sampling for consistent outputs
- Real-time token generation with streaming output
- System information display (GPU support, memory mapping, etc.)
- Detailed progress logging

## Command Line Options

- `-model string`: Path to the GGUF model file (default: "../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf")
- `-prompt string`: Prompt text to generate from (default: "The future of AI is")
- `-n-predict int`: Number of tokens to predict/generate (default: 50)
- `-threads int`: Number of threads to use (default: 4)
- `-ctx int`: Context size - maximum number of tokens in memory (default: 2048)

## Examples

### Basic Text Completion

```bash
go run main.go -prompt "Once upon a time"
```

### Longer Text Generation

```bash
go run main.go -prompt "Explain quantum computing" -n-predict 200
```

### Creative Writing

```bash
go run main.go -prompt "In a world where time travel is possible," -n-predict 150
```

### Technical Explanation

```bash
go run main.go -prompt "How does machine learning work?" -n-predict 100
```

### Conversation Starter

```bash
go run main.go -prompt "Hello! I'm an AI assistant. How can I help you today?" -n-predict 80
```

### High-Performance Configuration

```bash
go run main.go -prompt "Your prompt" -threads 8 -ctx 4096 -n-predict 300
```

## Using the Makefile

The included Makefile provides convenient shortcuts:

```bash
# Build the example
make build

# Run with default settings
make run

# Run creative writing demo
make creative

# Run technical explanation demo
make technical

# Run conversation demo
make conversation

# Run full demonstration
make demo

# Test compilation
make test

# Clean build artifacts
make clean
```

## Understanding the Output

The example provides detailed logging of each step:

1. **Initialization**: Backend setup and system information
2. **Model Loading**: Loading the GGUF model file into memory
3. **Context Creation**: Setting up the inference context with specified parameters
4. **Tokenization**: Converting input text into tokens the model understands
5. **Prompt Processing**: Running the prompt through the model
6. **Token Generation**: Generating new tokens one by one with real-time output

### System Information

The example displays important system capabilities:
- **GPU offload support**: Whether GPU acceleration is available
- **Memory mapping**: Whether the model can be memory-mapped for efficiency
- **Memory locking**: Whether memory can be locked to prevent swapping
- **Max devices**: Maximum number of compute devices available

### Generation Process

During generation, you'll see:
- Token sampling information
- Individual token IDs and their text representations
- Real-time text output as it's generated
- Final statistics (number of tokens generated)

## Configuration Tips

### Context Size (`-ctx`)
- **Smaller (512-1024)**: Faster, less memory, shorter conversations
- **Medium (2048-4096)**: Good balance for most use cases
- **Larger (8192+)**: Longer conversations, more memory required

### Token Count (`-n-predict`)
- **Short (20-50)**: Quick responses, sentence completion
- **Medium (100-200)**: Paragraph-length responses
- **Long (300+)**: Detailed explanations, stories

### Thread Count (`-threads`)
- Use the number of CPU cores available for best performance
- More threads = faster generation (up to hardware limits)

## Model Requirements

### Compatible Models
- Any GGUF model that supports text generation
- Popular choices: LLaMA, Mistral, CodeLlama, Alpaca, Vicuna
- Both chat-tuned and base models work

### Model Size Considerations
- **Small models (1-7B parameters)**: Fast, lower quality, less RAM
- **Medium models (13-30B parameters)**: Good balance of speed and quality
- **Large models (65B+ parameters)**: High quality, slower, much more RAM

## Troubleshooting

### "Failed to load model"
- Check that the model path is correct
- Ensure the model file is a valid GGUF file
- Make sure you have enough RAM to load the model

### "Failed to create context"
- Reduce context size with `-ctx 1024` or smaller
- Reduce batch size (this is set automatically)
- Close other applications to free memory

### Slow Generation
- Increase thread count with `-threads 8` (adjust for your CPU)
- Use a smaller model
- Reduce context size if not needed

### Repetitive or Poor Quality Output
- The example uses greedy sampling for consistency
- Try different prompts with more specific instructions
- Consider using a different (potentially larger) model

### Out of Memory Errors
- Use a smaller model
- Reduce context size significantly (`-ctx 512`)
- Close other applications
- Consider using memory mapping if supported

## Technical Details

This example demonstrates:

1. **Backend Initialization**: Setting up the llama.cpp backend
2. **Model Loading**: Loading GGUF models with custom parameters
3. **Context Management**: Creating and configuring inference contexts
4. **Tokenization**: Converting text to tokens and back
5. **Batch Processing**: Using batches for efficient inference
6. **Sampling**: Using greedy sampling for token selection
7. **Memory Management**: Proper cleanup of resources

### Code Structure

- **Initialization**: Backend and system info setup
- **Model Loading**: GGUF model loading with parameters
- **Context Creation**: Inference context with custom settings
- **Tokenization**: Prompt to tokens conversion
- **Inference Loop**: Token-by-token generation with real-time output
- **Cleanup**: Proper resource deallocation

## Next Steps

After trying this example, you might want to explore:

1. **[Embedding Example](../embedding/)**: Generate embeddings for semantic analysis
2. **Advanced Sampling**: Implement temperature, top-k, top-p sampling
3. **Chat Interface**: Build a full conversational interface
4. **Model Comparison**: Compare outputs from different models
5. **Performance Optimization**: Implement GPU offloading, quantization

## Requirements

- Go 1.21 or later
- A GGUF model file
- Sufficient RAM to load the model (varies by model size)
- CPU with multiple cores for optimal performance (optional)
