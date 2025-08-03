# Gollama.cpp

[![Go Reference](https://pkg.go.dev/badge/github.com/ltarantino/gollama.cpp.svg)](https://pkg.go.dev/github.com/ltarantino/gollama.cpp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/ltarantino/gollama.cpp.svg)](https://github.com/ltarantino/gollama.cpp/releases)

A high-performance Go binding for [llama.cpp](https://github.com/ggml-org/llama.cpp) using [purego](https://github.com/ebitengine/purego) for cross-platform compatibility without CGO.

## Features

- **Pure Go**: No CGO required, uses purego for C interop
- **Cross-Platform**: Supports macOS (CPU/Metal), Linux (CPU/NVIDIA/AMD), Windows (CPU/NVIDIA/AMD)
- **Performance**: Direct bindings to llama.cpp shared libraries
- **Compatibility**: Version-synchronized with llama.cpp releases
- **Easy Integration**: Simple Go API for LLM inference
- **GPU Acceleration**: Supports Metal, CUDA, HIP, and other backends

## Supported Platforms

### macOS
- **CPU**: Intel x64, Apple Silicon (ARM64)
- **GPU**: Metal (Apple Silicon)

### Linux
- **CPU**: x86_64, ARM64
- **GPU**: NVIDIA (CUDA), AMD (HIP/ROCm)

### Windows
- **CPU**: x86_64, ARM64
- **GPU**: NVIDIA (CUDA), AMD (HIP)

## Installation

```bash
go get github.com/ltarantino/gollama.cpp
```

The Go module includes pre-built llama.cpp libraries for all supported platforms. No additional installation required!

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/ltarantino/gollama.cpp"
)

func main() {
    // Initialize the library
    gollama.Backend_init()
    defer gollama.Backend_free()

    // Load model
    params := gollama.Model_default_params()
    model, err := gollama.Model_load_from_file("path/to/model.gguf", params)
    if err != nil {
        log.Fatal(err)
    }
    defer gollama.Model_free(model)

    // Create context
    ctxParams := gollama.Context_default_params()
    ctx, err := gollama.Init_from_model(model, ctxParams)
    if err != nil {
        log.Fatal(err)
    }
    defer gollama.Free(ctx)

    // Tokenize and generate
    prompt := "The future of AI is"
    tokens, err := gollama.Tokenize(model, prompt, true, false)
    if err != nil {
        log.Fatal(err)
    }

    // Create batch and decode
    batch := gollama.Batch_init(len(tokens), 0, 1)
    defer gollama.Batch_free(batch)

    for i, token := range tokens {
        gollama.Batch_add(batch, token, int32(i), []int32{0}, false)
    }

    if err := gollama.Decode(ctx, batch); err != nil {
        log.Fatal(err)
    }

    // Sample next token
    logits := gollama.Get_logits_ith(ctx, -1)
    candidates := gollama.Token_data_array_init(model)
    
    sampler := gollama.Sampler_init_greedy()
    defer gollama.Sampler_free(sampler)
    
    newToken := gollama.Sampler_sample(sampler, ctx, candidates)
    
    // Convert token to text
    text := gollama.Token_to_piece(model, newToken, false)
    fmt.Printf("Generated: %s\n", text)
}
```

## Advanced Usage

### GPU Configuration

```go
// Enable Metal on macOS
params := gollama.Context_default_params()
params.n_gpu_layers = 32 // Offload layers to GPU

// Enable CUDA on Linux/Windows
params.split_mode = gollama.LLAMA_SPLIT_MODE_LAYER
```

### Model Loading Options

```go
params := gollama.Model_default_params()
params.n_ctx = 4096           // Context size
params.use_mmap = true        // Memory mapping
params.use_mlock = true       // Memory locking
params.vocab_only = false     // Load full model
```

### Sampling Configuration

```go
// Temperature sampling
sampler := gollama.Sampler_init_temp(0.8)

// Top-k + Top-p sampling
chain := gollama.Sampler_chain_init(gollama.Sampler_chain_default_params())
gollama.Sampler_chain_add(chain, gollama.Sampler_init_top_k(40))
gollama.Sampler_chain_add(chain, gollama.Sampler_init_top_p(0.9, 1))
gollama.Sampler_chain_add(chain, gollama.Sampler_init_temp(0.8))
```

## Building from Source

### Prerequisites

- Go 1.21 or later
- Make
- Platform-specific build tools:
  - **macOS**: Xcode Command Line Tools
  - **Linux**: GCC, CUDA SDK (for NVIDIA), ROCm (for AMD)
  - **Windows**: Visual Studio or MinGW-w64, CUDA SDK (for NVIDIA)

### Build Commands

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build with GPU support
make build-gpu

# Run tests
make test

# Generate release packages
make release
```

## Version Compatibility

This library tracks llama.cpp versions. The version number format is:

```
vX.Y.Z-llamacpp.ABCD
```

Where:
- `X.Y.Z` is the gollama.cpp semantic version
- `ABCD` is the corresponding llama.cpp build number

For example: `v1.0.0-llamacpp.b6076` uses llama.cpp build b6076.

## Documentation

- [API Reference](https://pkg.go.dev/github.com/ltarantino/gollama.cpp)
- [Examples](./examples/)
- [Build Guide](./docs/BUILD.md)
- [GPU Setup](./docs/GPU.md)
- [Performance Tuning](./docs/PERFORMANCE.md)

## Examples

See the [examples](./examples/) directory for complete examples:

- [Simple Chat](./examples/simple-chat/)
- [Chat with History](./examples/chat-history/)
- [Embedding Generation](./examples/embeddings/)
- [Model Quantization](./examples/quantize/)
- [Batch Processing](./examples/batch/)
- [GPU Acceleration](./examples/gpu/)

## Contributing

Contributions are welcome! Please read our [Contributing Guide](./CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

This license is compatible with llama.cpp's MIT license.

## Acknowledgments

- [llama.cpp](https://github.com/ggml-org/llama.cpp) - The underlying C++ library
- [purego](https://github.com/ebitengine/purego) - Pure Go C interop library
- [ggml](https://github.com/ggml-org/ggml) - Machine learning tensor library

## Support

- [Issues](https://github.com/ltarantino/gollama.cpp/issues) - Bug reports and feature requests
- [Discussions](https://github.com/ltarantino/gollama.cpp/discussions) - Questions and community support
