# Gollama.cpp

[![Go Reference](https://pkg.go.dev/badge/github.com/dianlight/gollama.cpp.svg)](https://pkg.go.dev/github.com/dianlight/gollama.cpp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/dianlight/gollama.cpp.svg)](https://github.com/dianlight/gollama.cpp/releases)

A high-performance Go binding for [llama.cpp](https://github.com/ggml-org/llama.cpp) using [libgoffi](https://github.com/noctarius/libgoffi) for robust FFI support.

> **⚠️ Important Change**: This project now uses libgoffi (CGO-based) instead of purego. This provides better struct support and type mapping, but requires CGO and a C compiler. See [Requirements](#requirements) below.

## Features

- **CGO-based FFI**: Uses libgoffi for robust C interoperability
- **Cross-Platform**: Supports macOS (CPU/Metal), Linux (CPU/NVIDIA/AMD)
- **Performance**: Direct bindings to llama.cpp shared libraries
- **Compatibility**: Version-synchronized with llama.cpp releases
- **Easy Integration**: Simple Go API for LLM inference
- **GPU Acceleration**: Supports Metal, CUDA, HIP, Vulkan, OpenCL, SYCL, and other backends
- **Better Struct Support**: libgoffi provides improved handling of C structs

## Requirements

Since this project now uses libgoffi, you need:

1. **CGO enabled**: `CGO_ENABLED=1` (default on most systems)
2. **C compiler**: 
   - Linux: `gcc` or `clang`
   - macOS: Xcode Command Line Tools (`xcode-select --install`)
   - Windows: Not currently supported by libgoffi
3. **libffi**: The libffi library must be installed
   - Linux: `sudo apt-get install libffi-dev` (Debian/Ubuntu) or `sudo yum install libffi-devel` (RHEL/CentOS)
   - macOS: `brew install libffi` (or use system libffi)

## Supported Platforms

Gollama.cpp uses a **platform-specific architecture** with build tags to ensure optimal compatibility and performance across all operating systems.

### ✅ Fully Supported Platforms

#### macOS
- **CPU**: Intel x64, Apple Silicon (ARM64)
- **GPU**: Metal (Apple Silicon)
- **Status**: Full feature support with libgoffi
- **Build Tags**: Uses `!windows` build tag

#### Linux
- **CPU**: x86_64, ARM64
- **GPU**: NVIDIA (CUDA/Vulkan), AMD (HIP/ROCm/Vulkan), Intel (SYCL/Vulkan)
- **Status**: Full feature support with libgoffi
- **Build Tags**: Uses `!windows` build tag

### ❌ Not Supported

#### Windows
- **Status**: libgoffi does not support Windows
- Windows support would require either:
  - Returning to the previous purego implementation for Windows
  - Implementing a Windows-specific libffi wrapper
  - Using a different FFI solution for Windows

### Platform-Specific Implementation Details

Our platform abstraction layer uses Go build tags to provide:

- **Unix-like systems** (`!windows`): Uses [libgoffi](https://github.com/noctarius/libgoffi) (CGO-based, wraps libffi)
- **Windows** (`windows`): Not supported by libgoffi (would need alternative implementation)
- **Automatic detection**: Runtime platform capability detection

## Installation

### Prerequisites

Before installing, ensure you have the required dependencies:

**Linux:**
```bash
# Debian/Ubuntu
sudo apt-get install build-essential libffi-dev

# RHEL/CentOS/Fedora
sudo yum install gcc libffi-devel

# Arch Linux
sudo pacman -S base-devel libffi
```

**macOS:**
```bash
# Install Xcode Command Line Tools (if not already installed)
xcode-select --install

# libffi is usually included, but you can also install via Homebrew
brew install libffi
```

### Install Package

```bash
CGO_ENABLED=1 go get github.com/dianlight/gollama.cpp
```

The Go module automatically downloads pre-built llama.cpp libraries from the official [ggml-org/llama.cpp](https://github.com/ggml-org/llama.cpp) releases on first use. No manual compilation of llama.cpp required!

## Build Requirements

Since this project uses CGO via libgoffi:

- **CGO must be enabled**: Set `CGO_ENABLED=1` environment variable
- **C compiler required**: gcc, clang, or compatible C compiler
- **libffi required**: The Foreign Function Interface library
- **Cross-compilation**: More complex than pure Go due to CGO requirements

## Cross-Platform Development

### Build Compatibility Matrix

**Note**: With CGO and libgoffi, cross-compilation is more complex than with pure Go:

| Target Platform | Build Support | Runtime Support | Notes |
|------------------|:-------------:|:---------------:|-------|
| Linux (amd64)    | ✅            | ✅              | Full support with libgoffi |
| Linux (arm64)    | ✅            | ✅              | Full support with libgoffi |
| macOS (amd64)    | ✅            | ✅              | Full support with libgoffi |
| macOS (arm64)    | ✅            | ✅              | Full support with libgoffi |
| Windows (amd64)  | ❌            | ❌              | libgoffi does not support Windows |
| Windows (arm64)  | ❌            | ❌              | libgoffi does not support Windows |

### Development Workflow

```bash
# Build for current platform (requires CGO)
CGO_ENABLED=1 go build ./...

# Build for specific platform (requires cross-compilation toolchain)
# Cross-compiling with CGO is more complex - you need the target platform's
# C compiler and libraries. See Go CGO cross-compilation documentation.

# Run platform-specific tests
CGO_ENABLED=1 go test -v -run TestPlatformSpecific ./...
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

// Detect available GPU backend
backend := gollama.DetectGpuBackend()
fmt.Printf("Using GPU backend: %s\n", backend.String())

// Platform-specific optimizations:
// - macOS: Uses Metal when available  
// - Linux: Supports CUDA, HIP, Vulkan, and SYCL
// - Windows: Supports CUDA, HIP, Vulkan, OpenCL, and SYCL
params.split_mode = gollama.LLAMA_SPLIT_MODE_LAYER
```

#### GPU Support Matrix

| Platform | GPU Type | Backend | Status |
|----------|----------|---------|--------|
| macOS | Apple Silicon | Metal | ✅ Supported |
| macOS | Intel/AMD | CPU only | ✅ Supported |
| Linux | NVIDIA | CUDA | ✅ Available in releases |
| Linux | NVIDIA | Vulkan | ✅ Available in releases |
| Linux | AMD | HIP/ROCm | ✅ Available in releases |
| Linux | AMD | Vulkan | ✅ Available in releases |
| Linux | Intel | SYCL | ✅ Available in releases |
| Linux | Intel/Other | Vulkan | ✅ Available in releases |
| Linux | Intel/Other | CPU | ✅ Fallback |
| Windows | NVIDIA | CUDA | ✅ Available in releases |
| Windows | NVIDIA | Vulkan | ✅ Available in releases |
| Windows | AMD | HIP | ✅ Available in releases |
| Windows | AMD | Vulkan | ✅ Available in releases |
| Windows | Intel | SYCL | ✅ Available in releases |
| Windows | Qualcomm Adreno | OpenCL | ✅ Available in releases |
| Windows | Intel/Other | Vulkan | ✅ Available in releases |
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
err := gollama.LoadLibraryWithVersion("b6099")

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

# Test GPU detection and functionality
make test-gpu

# Detect available GPU backends
make detect-gpu

# Clean library cache
make clean-libs
```

#### Available Library Variants

The downloader automatically selects the best variant for your platform:

- **macOS**: Metal-enabled binaries (arm64/x64)
- **Linux**: CPU-optimized binaries (CUDA/HIP/Vulkan/SYCL versions available)
- **Windows**: CPU-optimized binaries (CUDA/HIP/Vulkan/OpenCL/SYCL versions available)

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

For example: `v0.2.0-llamacpp.b6099` uses llama.cpp build b6099.

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

## Funding

If you find this project helpful and would like to support its development, you can:

- ⭐ Star this repository on GitHub
- 🐛 Report bugs and suggest improvements
- 📖 Improve documentation

[![GitHub Sponsors](https://img.shields.io/badge/Sponsor-GitHub-pink?style=for-the-badge&logo=github)](https://github.com/sponsors/dianlight)
[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black)](https://www.buymeacoffee.com/ypKZ2I0)

Your support helps maintain and improve this project for the entire community!

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
