# Gollama.cpp Simple Chat with Library Loader Example

This example demonstrates how to combine library loading capabilities with basic text generation using gollama.cpp. It shows how to dynamically load embedded libraries and then use them for text generation, providing a complete end-to-end demonstration of both the library loader and chat functionality.

## Quick Start

To get started right away, run the following command, making sure to use the correct path for the model you have:

```bash
cd examples/simple-chat-with-loader
go run main.go -prompt "The future of AI is"
```

## Features

- **Dynamic Library Loading**: Automatically extracts and loads embedded llama.cpp libraries
- **Library Lifecycle Management**: Proper loading and unloading of shared libraries
- **Basic Text Generation**: Text completion and conversation capabilities
- **Configurable Model Parameters**: Adjustable context sizes, threads, and prediction length
- **Greedy Sampling**: Consistent text generation with deterministic outputs
- **Real-time Token Generation**: Streaming output with detailed progress logging
- **System Information Display**: GPU support detection, memory mapping capabilities
- **Cross-platform Support**: Works on macOS, Linux, and Windows
- **Resource Cleanup**: Automatic cleanup of temporary files and library handles

## What This Example Demonstrates

1. **Library Loader Integration**: Shows how to use the `LibraryLoader` to dynamically load llama.cpp libraries
2. **Combined Functionality**: Integrates both library management and text generation in a single application
3. **Error Handling**: Proper error handling for both library loading and model operations
4. **Resource Management**: Demonstrates proper cleanup of both library handles and model resources

## Command Line Options

- `-model string`: Path to the GGUF model file (default: "../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf")
- `-prompt string`: Prompt text to generate from (default: "The future of AI is")
- `-n-predict int`: Number of tokens to predict/generate (default: 50)
- `-threads int`: Number of threads to use (default: 4)
- `-ctx int`: Context size - maximum number of tokens in memory (default: 2048)

## Examples

### Basic Text Completion with Library Loading

```bash
go run main.go -prompt "Once upon a time"
```

### Longer Text Generation

```bash
go run main.go -prompt "Explain quantum computing" -n-predict 200
```

### Creative Writing with Custom Context

```bash
go run main.go -prompt "In a world where time travel is possible," -n-predict 150 -ctx 4096
```

## Output Structure

The example produces output in two main sections:

1. **Library Loader Demo**: Shows the process of loading and managing the llama.cpp library
2. **Simple Chat Demo**: Demonstrates text generation using the loaded library

Example output:
```
Gollama.cpp Simple Chat with Library Loader Example v1.0.0
Model: ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf
Prompt: The future of AI is
Threads: 4
Context: 2048

=== Library Loader Demo ===
Testing library extraction and loading... done
Library loaded successfully
Handle: 123456789
IsLoaded: true

=== Simple Chat Demo ===
Initializing backend... done
GPU offload: supported
Memory mapping: true
Memory locking: true
Max devices: 1

Loading model... done
Creating context... done
Tokenizing prompt... done (5 tokens)
Processing prompt... done

Generated text:
The future of AI is bright and full of possibilities...
```

## Technical Details

This example demonstrates:

- **Embedded Library Extraction**: The `LibraryLoader` extracts pre-compiled libraries from embedded files
- **Platform Detection**: Automatically detects the current platform and loads the appropriate library
- **Memory Management**: Proper allocation and deallocation of both library handles and model resources
- **Thread Safety**: Uses mutex locks to ensure thread-safe library operations
- **Temporary File Cleanup**: Automatically removes extracted files when the library is unloaded

## Requirements

- Go 1.21 or later
- Compatible GGUF model file
- Sufficient system memory for the model and context size
- Platform-specific llama.cpp libraries (automatically handled by the loader)

## Building

```bash
cd examples/simple-chat-with-loader
go build -o simple-chat-with-loader main.go
```

## Running

```bash
./simple-chat-with-loader -model path/to/your/model.gguf -prompt "Your prompt here"
```
