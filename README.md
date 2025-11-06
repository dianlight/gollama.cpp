# Gollama.cpp

[![Go Reference](https://pkg.go.dev/badge/github.com/dianlight/gollama.cpp.svg)](https://pkg.go.dev/github.com/dianlight/gollama.cpp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/dianlight/gollama.cpp.svg)](https://github.com/dianlight/gollama.cpp/releases)

A high-performance Go binding for [llama.cpp](https://github.com/ggml-org/llama.cpp) using [purego](https://github.com/ebitengine/purego) and [libffi](https://github.com/jupiterrider/ffi) for cross-platform compatibility without CGO.

## Features

- **Pure Go**: No CGO required, uses purego and libffi for C interop
- **Cross-Platform**: Supports macOS (CPU/Metal), Linux (CPU/NVIDIA/AMD), Windows (CPU/NVIDIA/AMD)
- **Struct Support**: Uses libffi for calling C functions with struct parameters/returns on all platforms
- **Performance**: Direct bindings to llama.cpp shared libraries
- **Compatibility**: Version-synchronized with llama.cpp releases
- **Easy Integration**: Simple Go API for LLM inference
- **GPU Acceleration**: Supports Metal, CUDA, HIP, Vulkan, OpenCL, SYCL, and other backends
- **Embedded Runtime Libraries**: Optional go:embed bundle for all supported platforms
- **GGML Bindings**: Low-level GGML tensor library bindings for advanced use cases

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
- **GPU**: NVIDIA (CUDA/Vulkan), AMD (HIP/ROCm/Vulkan), Intel (SYCL/Vulkan)
- **Status**: Full feature support with purego and libffi
- **Build Tags**: Uses `!windows` build tag

#### Windows
- **CPU**: x86_64, ARM64 
- **GPU**: NVIDIA (CUDA/Vulkan), AMD (HIP/Vulkan), Intel (SYCL/Vulkan), Qualcomm Adreno (OpenCL)
- **Status**: **Full feature support with libffi**
- **Build Tags**: Uses `windows` build tag with syscall-based library loading
- **Current State**: 
  - ✅ Compiles without errors on Windows
  - ✅ Cross-compilation from other platforms works
  - ✅ Runtime functionality fully enabled via libffi and GetProcAddress
  - ✅ Full struct parameter/return support through function registration
  - 🚧 GPU acceleration being tested

> Windows runtime notes
>
> - The loader now adds the DLL's directory to the Windows DLL search path and uses `LoadLibraryExW` with safe search flags to reliably resolve sibling dependencies (ggml, libomp, libcurl, etc.).
> - When a symbol isn't found in `llama.dll`, resolution automatically searches sibling DLLs from the same directory (e.g., `ggml*.dll`). This matches how upstream splits exports on Windows and fixes missing `llama_backend_*` on some builds.
> - If you see “The specified module could not be found.” while loading `llama.dll`, it often indicates a missing system runtime (e.g., Microsoft Visual C++ Redistributable 2015–2022). Installing the latest x64/x86 redistributable typically resolves it.
> - CI runners set PATH for later steps, but the downloader verifies loading immediately after download; the improved loader handles dependency resolution without relying on PATH.

### Platform-Specific Implementation Details

Our platform abstraction layer uses Go build tags to provide:

- **Unix-like systems** (`!windows`): Uses [purego](https://github.com/ebitengine/purego) for dynamic library loading
- **Windows** (`windows`): Uses native Windows syscalls (`LoadLibraryW`, `FreeLibrary`, `GetProcAddress`)
- **All platforms**: Uses [libffi](https://github.com/jupiterrider/ffi) for calling C functions with struct parameters/returns
- **Cross-compilation**: Supports building for any platform from any platform
- **Automatic detection**: Runtime platform capability detection

## Installation

```bash
go get github.com/dianlight/gollama.cpp
```

The Go module automatically downloads pre-built llama.cpp libraries from the official [ggml-org/llama.cpp](https://github.com/ggml-org/llama.cpp) releases on first use. No manual compilation required!

### Embedding Libraries

For reproducible builds you can embed the pre-built libraries directly into the Go module. A helper Makefile target downloads the configured llama.cpp build (`LLAMA_CPP_BUILD`) for every supported platform and synchronises the `./libs` directory which is picked up by `go:embed`:

```bash
# Download all platform builds for the configured llama.cpp version and populate ./libs
make populate-libs

# Alternatively, use the CLI directly
go run ./cmd/gollama-download -download-all -version b6862 -copy-libs
```

Only a single llama.cpp version is stored in `./libs` at a time. Running `populate-libs` removes outdated directories automatically. Subsequent `go build` invocations embed the freshly synchronised libraries and `LoadLibraryWithVersion("")` will prefer the embedded bundle.

## Cross-Platform Development

### Build Compatibility Matrix

Our CI system tests compilation across all platforms:

| Target Platform | Build From Linux | Build From macOS | Build From Windows |
| --------------- | :--------------: | :--------------: | :----------------: |
| Linux (amd64)   |        ✅         |        ✅         |         ✅          |
| Linux (arm64)   |        ✅         |        ✅         |         ✅          |
| macOS (amd64)   |        ✅         |        ✅         |         ✅          |
| macOS (arm64)   |        ✅         |        ✅         |         ✅          |
| Windows (amd64) |        ✅         |        ✅         |         ✅          |
| Windows (arm64) |        ✅         |        ✅         |         ✅          |

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

### GGML Low-Level API

For advanced use cases, gollama.cpp provides direct access to GGML (the tensor library powering llama.cpp):

```go
// Check GGML type information
typeSize, err := gollama.Ggml_type_size(gollama.GGML_TYPE_F32)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("F32 type size: %d bytes\n", typeSize)

// Check if a type is quantized
isQuantized, err := gollama.Ggml_type_is_quantized(gollama.GGML_TYPE_Q4_0)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Q4_0 is quantized: %v\n", isQuantized)

// Enumerate backend devices
devCount, err := gollama.Ggml_backend_dev_count()
if err == nil && devCount > 0 {
    for i := uint64(0); i < devCount; i++ {
        dev, _ := gollama.Ggml_backend_dev_get(i)
        name, _ := gollama.Ggml_backend_dev_name(dev)
        fmt.Printf("Device %d: %s\n", i, name)
    }
}
```

**Supported GGML Features:**
- 31 tensor type definitions (F32, F16, Q4_0, Q8_0, BF16, etc.)
- Type size and quantization utilities
- Backend device enumeration and management
- Buffer allocation and management
- Type information queries

**Note:** GGML functions may not be exported in all llama.cpp builds. The library gracefully handles missing functions without errors.

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

| Platform | GPU Type        | Backend  | Status                  |
| -------- | --------------- | -------- | ----------------------- |
| macOS    | Apple Silicon   | Metal    | ✅ Supported             |
| macOS    | Intel/AMD       | CPU only | ✅ Supported             |
| Linux    | NVIDIA          | CUDA     | ✅ Available in releases |
| Linux    | NVIDIA          | Vulkan   | ✅ Available in releases |
| Linux    | AMD             | HIP/ROCm | ✅ Available in releases |
| Linux    | AMD             | Vulkan   | ✅ Available in releases |
| Linux    | Intel           | SYCL     | ✅ Available in releases |
| Linux    | Intel/Other     | Vulkan   | ✅ Available in releases |
| Linux    | Intel/Other     | CPU      | ✅ Fallback              |
| Windows  | NVIDIA          | CUDA     | ✅ Available in releases |
| Windows  | NVIDIA          | Vulkan   | ✅ Available in releases |
| Windows  | AMD             | HIP      | ✅ Available in releases |
| Windows  | AMD             | Vulkan   | ✅ Available in releases |
| Windows  | Intel           | SYCL     | ✅ Available in releases |
| Windows  | Qualcomm Adreno | OpenCL   | ✅ Available in releases |
| Windows  | Intel/Other     | Vulkan   | ✅ Available in releases |
| Windows  | Intel/Other     | CPU      | ✅ Fallback              |

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
err := gollama.LoadLibraryWithVersion("b6862")

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

Downloaded libraries are cached in platform-specific locations:
- **Linux/macOS**: `~/.cache/gollama/libs/`
- **Windows**: `%LOCALAPPDATA%/gollama/libs/`

You can customize the cache directory in several ways:

**Environment Variable:**
```bash
export GOLLAMA_CACHE_DIR=/custom/path/to/cache
```

**Configuration File:**
```json
{
  "cache_dir": "/custom/path/to/cache"
}
```

**Programmatically:**
```go
config := gollama.DefaultConfig()
config.CacheDir = "/custom/path/to/cache"
gollama.SetGlobalConfig(config)
```

To get the current cache directory:
```go
cacheDir, err := gollama.GetLibraryCacheDir()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Using cache directory: %s\n", cacheDir)
```

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

## Running Tests

```bash
make test
```

Tests use `github.com/stretchr/testify/suite` along with a shared `BaseSuite` (see `test_base_suite_test.go`) that automatically snapshots/restores configuration and environment variables and unloads the llama library after each test. See the Contributing guide for details.

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

For example: `v0.2.0-llamacpp.b6862` uses llama.cpp build b6862.

## Documentation

- [API Reference](https://pkg.go.dev/github.com/dianlight/gollama.cpp)
  - **High-Level API**: `gollama.go` - Complete llama.cpp bindings
  - **Low-Level API**: `goggml.go` - GGML tensor library bindings
- [GGML Low-Level API Guide](./docs/GGML_API.md) - Detailed GGML bindings documentation
- [Examples](./examples/)
- [Build Guide](./docs/BUILD.md)
- [GPU Setup](./docs/GPU.md)

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
