# Gollama.cpp

[![Go Reference](https://pkg.go.dev/badge/github.com/dianlight/gollama.cpp.svg)](https://pkg.go.dev/github.com/dianlight/gollama.cpp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/ltarantino/gollama.cpp.svg)](https://github.com/dianlight/gollama.cpp/releases)

A high-performance Go binding for [llama.cpp](https://github.com/ggml-org/llama.cpp) using [purego](https://github.com/ebitengine/purego) for cross-platform compatibility without CGO.

## Features

- **Pure Go**: No CGO required, uses purego for C interop
- **Cross-Platform**: Supports macOS (CPU/Metal), Linux (CPU/NVIDIA/AMD), Windows (CPU/NVIDIA/AMD)
- **Performance**: Direct bindings to llama.cpp shared libraries
- **Compatibility**: Version-synchronized with llama.cpp releases
- **Easy Integration**: Simple Go API for LLM inference
- **GPU Acceleration**: Supports Metal, CUDA, HIP, and other backends

## Supported Platforms

Gollama.cpp uses a **platform-specific architecture** with build tags to ensure optimal compatibility and performance across all operating systems.

### ✅ Fully Supported Platforms

#### macOS
- **CPU**: Intel x64, Apple Silicon (ARM64)
- **GPU**: Metal (Apple Silicon)
- **Status**: Full feature support with purego
- **Build Tags**: Uses `!windows` build tag

#### Linux
- **CPU**: x86_64, ARM64
- **GPU**: NVIDIA (CUDA), AMD (HIP/ROCm)
- **Status**: Full feature support with purego
- **Build Tags**: Uses `!windows` build tag

### 🚧 In Development

#### Windows
- **CPU**: x86_64, ARM64 
- **GPU**: NVIDIA (CUDA), AMD (HIP) - planned
- **Status**: **Build compatibility implemented**, runtime support in development
- **Build Tags**: Uses `windows` build tag with syscall-based library loading
- **Current State**: 
  - ✅ Compiles without errors on Windows
  - ✅ Cross-compilation from other platforms works
  - 🚧 Runtime functionality being implemented
  - 🚧 GPU acceleration being added

### Platform-Specific Implementation Details

Our platform abstraction layer uses Go build tags to provide:

- **Unix-like systems** (`!windows`): Uses [purego](https://github.com/ebitengine/purego) for dynamic library loading
- **Windows** (`windows`): Uses native Windows syscalls (`LoadLibraryW`, `FreeLibrary`, `GetProcAddress`)
- **Cross-compilation**: Supports building for any platform from any platform
- **Automatic detection**: Runtime platform capability detection

## Installation

```bash
go get github.com/dianlight/gollama.cpp
```

The Go module automatically downloads pre-built llama.cpp libraries from the official [ggml-org/llama.cpp](https://github.com/ggml-org/llama.cpp) releases on first use. No manual compilation required!

## Cross-Platform Development

### Build Compatibility Matrix

Our CI system tests compilation across all platforms:

| Target Platform | Build From Linux | Build From macOS | Build From Windows |
|------------------|:----------------:|:----------------:|:------------------:|
| Linux (amd64)    | ✅               | ✅               | ✅                 |
| Linux (arm64)    | ✅               | ✅               | ✅                 |
| macOS (amd64)    | ✅               | ✅               | ✅                 |
| macOS (arm64)    | ✅               | ✅               | ✅                 |
| Windows (amd64)  | ✅               | ✅               | ✅                 |
| Windows (arm64)  | ✅               | ✅               | ✅                 |

### Development Workflow

```bash
# Test cross-compilation for all platforms
make test-cross-compile

# Build for specific platform
GOOS=windows GOARCH=amd64 go build ./...
GOOS=linux GOARCH=arm64 go build ./...
GOOS=darwin GOARCH=arm64 go build ./...

# Run platform-specific tests
go test -v -run TestPlatformSpecific ./...
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/dianlight/gollama.cpp"
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

Gollama.cpp automatically downloads the appropriate pre-built binaries with GPU support and configures the optimal backend:

```go
// Automatic GPU detection and configuration
params := gollama.Context_default_params()
params.n_gpu_layers = 32 // Offload layers to GPU (if available)

// Platform-specific optimizations:
// - macOS: Uses Metal when available  
// - Linux: Supports CUDA and HIP
// - Windows: Supports CUDA and HIP
params.split_mode = gollama.LLAMA_SPLIT_MODE_LAYER
```

#### GPU Support Matrix

| Platform | GPU Type | Backend | Status |
|----------|----------|---------|--------|
| macOS | Apple Silicon | Metal | ✅ Supported |
| macOS | Intel/AMD | CPU only | ✅ Supported |
| Linux | NVIDIA | CUDA | ✅ Available in releases |
| Linux | AMD | HIP/ROCm | ✅ Available in releases |
| Linux | Intel/Other | CPU | ✅ Fallback |
| Windows | NVIDIA | CUDA | ✅ Available in releases |
| Windows | AMD | HIP | ✅ Available in releases |
| Windows | Intel/Other | CPU | ✅ Fallback |

The library automatically downloads pre-built binaries from the official llama.cpp releases with the appropriate GPU support for your platform. The download happens automatically on first use!

### Model Loading Options

```go
params := gollama.Model_default_params()
params.n_ctx = 4096           // Context size
params.use_mmap = true        // Memory mapping
params.use_mlock = true       // Memory locking
params.vocab_only = false     // Load full model
```

### Library Management

Gollama.cpp automatically downloads pre-built binaries from the official llama.cpp releases. You can also manage libraries manually:

```go
// Load a specific version
err := gollama.LoadLibraryWithVersion("b6089")

// Clean cache to force re-download
err := gollama.CleanLibraryCache()
```

#### Command Line Tools

```bash
# Download libraries for current platform
make download-libs

# Download libraries for all platforms  
make download-libs-all

# Test download functionality
make test-download

# Clean library cache
make clean-libs
```

#### Available Library Variants

The downloader automatically selects the best variant for your platform:

- **macOS**: Metal-enabled binaries (arm64/x64)
- **Linux**: CPU-optimized binaries (CUDA/HIP versions available)
- **Windows**: CPU-optimized binaries (CUDA/HIP versions available)

#### Cache Location

Downloaded libraries are cached in:
- **Linux/macOS**: `~/.cache/gollama/libs/`
- **Windows**: `%LOCALAPPDATA%/gollama/libs/`

## Building from Source

### Prerequisites

- Go 1.21 or later
- Make

### Quick Start

```bash
# Clone and build
git clone https://github.com/dianlight/gollama.cpp
cd gollama.cpp

# Build for current platform
make build

# Run tests (downloads libraries automatically)
make test

# Build examples
make build-examples

# Run tests
make test

# Generate release packages
make release
```

### GPU Detection Logic

The Makefile implements intelligent GPU detection:

1. **CUDA Detection**: Checks for `nvcc` compiler and CUDA toolkit
2. **HIP Detection**: Checks for `hipconfig` and ROCm installation  
3. **Priority Order**: CUDA > HIP > CPU (on Linux/Windows)
4. **Metal**: Always enabled on macOS when Xcode is available

No manual configuration or environment variables required!

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

- [API Reference](https://pkg.go.dev/github.com/dianlight/gollama.cpp)
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

- [Issues](https://github.com/dianlight/gollama.cpp/issues) - Bug reports and feature requests
- [Discussions](https://github.com/dianlight/gollama.cpp/discussions) - Questions and community support
