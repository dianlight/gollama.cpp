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

The Go module includes pre-built llama.cpp libraries for all supported platforms. No additional installation required!

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

Gollama.cpp automatically detects available GPU hardware and configures the optimal backend:

```go
// Automatic GPU detection and configuration
params := gollama.Context_default_params()
params.n_gpu_layers = 32 // Offload layers to GPU (if available)

// Platform-specific optimizations:
// - macOS: Uses Metal when available
// - Linux: Prefers CUDA > HIP > CPU
// - Windows: Prefers CUDA > CPU
params.split_mode = gollama.LLAMA_SPLIT_MODE_LAYER
```

#### GPU Support Matrix

| Platform | GPU Type | Backend | Status |
|----------|----------|---------|--------|
| macOS | Apple Silicon | Metal | ✅ Supported |
| macOS | Intel/AMD | CPU only | ✅ Supported |
| Linux | NVIDIA | CUDA | ✅ Auto-detected |
| Linux | AMD | HIP/ROCm | ✅ Auto-detected |
| Linux | Intel/Other | CPU | ✅ Fallback |
| Windows | NVIDIA | CUDA | ✅ Auto-detected |
| Windows | AMD | HIP | ✅ Auto-detected |
| Windows | Intel/Other | CPU | ✅ Fallback |

The build system automatically detects available GPU SDKs and enables the appropriate backends during compilation. No manual configuration required!

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
- CMake 3.15 or later
- Platform-specific build tools:
  - **macOS**: Xcode Command Line Tools
  - **Linux**: GCC or Clang
  - **Windows**: Visual Studio 2019+ or MinGW-w64

### Optional GPU SDKs (Auto-detected)

The build system automatically detects and enables GPU support when available:

- **NVIDIA CUDA**: CUDA Toolkit 11.8+ (auto-detected via `nvcc` or `CUDA_PATH`)
- **AMD HIP/ROCm**: ROCm 5.0+ (auto-detected via `hipconfig` or `ROCM_PATH`)
- **Apple Metal**: Automatically available on macOS with Xcode

### Build Commands

```bash
# Build for current platform with automatic GPU detection
make build

# Build libraries only (useful for development)
make clone-llamacpp
make build-llamacpp-linux-amd64  # Detects CUDA/HIP automatically

# Build for all platforms
make build-all

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
