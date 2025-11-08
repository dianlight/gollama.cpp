# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.2.2-llamacpp.b6862] - 2025-11-08

### Added

- **GGML Low-Level Bindings** (`goggml.go`): Direct access to GGML tensor library functions
  - 31 GGML type definitions (F32, F16, Q4_0, Q8_0, BF16, IQ2_XXS, IQ4_XS, etc.)
  - Backend device enumeration and management functions
  - Buffer allocation and memory management
  - Type utility functions (size, quantization checks, name lookups)
  - Quantization helper functions
  - Graceful handling of optional GGML functions (may not be exported in all builds)
  - Helper function `tryRegisterLibFunc()` for optional function registration
  - Comprehensive test suite (`goggml_test.go`) with type tests, backend tests, and benchmarks
- **Comprehensive GGML Backend API** synchronized with `ggml-backend.h`:
  - **Backend Device Functions** (8 new): `ggml_backend_dev_by_name`, `ggml_backend_dev_type`, `ggml_backend_dev_get_props`, `ggml_backend_dev_backend_reg`, `ggml_backend_dev_buffer_from_host_ptr`, `ggml_backend_dev_supports_op`, `ggml_backend_dev_supports_buft`, `ggml_backend_dev_offload_op`
  - **Backend Registry Functions** (9 new): `ggml_backend_reg_name`, `ggml_backend_reg_dev_count`, `ggml_backend_reg_dev_get`, `ggml_backend_reg_get_proc_address`, `ggml_backend_register`, `ggml_backend_device_register`, `ggml_backend_reg_count`, `ggml_backend_reg_get`, `ggml_backend_reg_by_name`
  - **Backend Buffer Type Functions** (4 new): `ggml_backend_buft_get_alignment`, `ggml_backend_buft_get_max_size`, `ggml_backend_buft_get_alloc_size`, `ggml_backend_buft_get_device`
  - **Backend Buffer Functions** (6 new): `ggml_backend_buffer_init_tensor`, `ggml_backend_buffer_get_alignment`, `ggml_backend_buffer_get_max_size`, `ggml_backend_buffer_get_alloc_size`, `ggml_backend_buffer_get_usage`, `ggml_backend_buffer_reset`
  - **Backend Functions** (6 new): `ggml_backend_guid`, `ggml_backend_get_default_buffer_type`, `ggml_backend_alloc_buffer`, `ggml_backend_get_alignment`, `ggml_backend_get_max_size`, `ggml_backend_get_device`
  - **New Types**: `GgmlGuid` (16-byte backend identifier), `GgmlBackendDevCaps` (device capabilities), `GgmlBackendDevProps` (device properties)
  - **Device Type Constants**: `GGML_BACKEND_DEVICE_TYPE_CPU`, `GGML_BACKEND_DEVICE_TYPE_GPU`, `GGML_BACKEND_DEVICE_TYPE_IGPU`, `GGML_BACKEND_DEVICE_TYPE_ACCEL`

### Changed

### Fixed

- **GGML Backend Signatures**: Fixed `Ggml_backend_load()` to accept only one parameter (path) instead of two, matching the actual GGML C API. Function now returns `GgmlBackendReg` instead of `GgmlBackend`.
- **Function Parameter Validation**: Verified all ~50 GGML backend function signatures against official `ggml-backend.h` to ensure correct parameter types and counts.

### Removed


## [v0.2.1-llamacpp.b6862] - 2025-11-02

### Added
- **Windows Function Registration**: Implemented proper `registerLibFunc` for Windows using `GetProcAddress` to resolve function addresses from loaded DLLs. This enables cross-platform struct parameter/return support on Windows through libffi.
- **libffi Support**: Cross-platform struct handling for C function calls on all platforms (Windows, Linux, macOS)
  - Added `github.com/jupiterrider/ffi v0.5.1` dependency for FFI support
  - Implemented FFI wrapper layer (`ffi.go`) with 10 wrapper functions for struct-based operations
  - Added `Encode()` and `Sampler_chain_init()` public wrapper functions
  - Platform-agnostic GetProcAddress helper for symbol resolution
- **Embedded Library Packaging**: `go:embed` support for pre-built llama.cpp binaries across all platforms
  - New `make populate-libs` target and `-copy-libs` CLI flag to sync the `./libs` directory
  - Runtime prefers embedded libraries when requesting the bundled llama.cpp build
  - Release workflow auto-populates and commits the embedded library bundle

### Changed
- **llama.cpp Version**: Updated from b6099 to b6862
  - Updated Makefile, CI workflows, and all documentation
  - Improved compatibility with latest llama.cpp features
  - Deprecated KV cache functions (removed from b6862 API): `llama_kv_cache_clear`, `llama_kv_cache_seq_*`, `llama_kv_cache_defrag`, `llama_kv_cache_update`
  - Removed non-existent functions: `llama_sampler_init_softmax`
- **CI/CD Improvements**: 
  - Added automatic library download step before running tests
  - Configured platform-specific library paths (LD_LIBRARY_PATH, DYLD_LIBRARY_PATH, PATH)
  - Updated GO_VERSION to 1.25
- **Windows Support**: Enabled full runtime support with FFI (previously build-only)
- **Test Behavior**: FFI tests now fail instead of skip when library is unavailable

### Fixed

- Windows DLL loading: Improve reliability by adding the library directory to the DLL search path and preferring `LoadLibraryExW` with safe search flags. This resolves "The specified module could not be found" when loading embedded `llama.dll` on CI.

### Removed


## [v0.2.0-llamacpp.b6099] - 2025-08-06

### 🚀 Major GPU Backend Enhancements

This release significantly expands GPU support with **three new GPU backends** and comprehensive detection systems:

- **🔥 Vulkan Support**: Cross-platform GPU acceleration for NVIDIA, AMD, and Intel GPUs
- **⚡ OpenCL Support**: Broad compatibility including Qualcomm Adreno GPUs on ARM64 devices  
- **🧠 SYCL Support**: Intel oneAPI unified parallel programming for CPUs and GPUs
- **🔍 Smart Detection**: Automatic GPU backend detection with optimal library selection
- **📦 Intelligent Downloads**: GPU-aware library downloader selects best variants automatically

### Added
- **Enhanced GPU Backend Support** with comprehensive detection and automatic selection
  - **Vulkan support** for cross-platform GPU acceleration (NVIDIA, AMD, Intel GPUs)
  - **OpenCL support** for diverse GPU vendors including Qualcomm Adreno on ARM64
  - **SYCL support** for Intel oneAPI toolkit and unified parallel programming
  - **GPU backend detection** with `DetectGpuBackend()` function and `LlamaGpuBackend` enum
  - **Automatic GPU variant selection** in library downloader based on available SDKs
  - **GPU detection tools** with new `detect-gpu` Makefile target
  - **GPU testing suite** with `test-gpu` target and comprehensive backend tests
- **Enhanced Documentation** for GPU setup and configuration
  - Comprehensive GPU installation guides for Linux (Vulkan, OpenCL, SYCL)
  - Windows GPU setup instructions for all supported backends
  - Platform-specific GPU requirements and verification steps
  - Updated GPU Support Matrix with all available backends
- **Improved CI/CD Pipeline** with GPU detection testing
  - Added GPU detection tools installation in CI workflows
  - New `gpu-detection` job for comprehensive GPU backend testing
  - GPU-aware library downloading validation
- **Download-based architecture** using pre-built binaries from official llama.cpp releases
- Automatic library download system with platform detection
- Library cache management with `clean-libs` target
- Cross-platform download testing (`test-download-platforms`)
- Command-line download tool (`cmd/gollama-download`)
- `clone-llamacpp` target for developers needing source code cross-reference
- **Platform-specific architecture** with Go build tags for improved cross-platform support
- Windows compilation compatibility using native syscalls (`LoadLibraryW`, `FreeLibrary`)
- Cross-platform compilation testing in CI pipeline
- Platform capability detection functions (`isPlatformSupported`, `getPlatformError`)
- **Integrated hf.sh script management** for Hugging Face model downloads
- `update-hf-script` target for updating hf.sh from llama.cpp repository
- Enhanced model download system using hf.sh instead of direct curl
- Comprehensive tools documentation (`docs/TOOLS.md`)
- Dedicated platform-specific test suite (`TestPlatformSpecific`)
- Enhanced Makefile with cross-compilation targets (`test-cross-compile`, `test-compile-*`)
- Comprehensive platform migration documentation
- Initial Go binding for llama.cpp using purego
- Cross-platform support (macOS, Linux, Windows)
- CPU and GPU acceleration support

### Changed
- **Enhanced GPU Backend Priority**: Updated detection order to CUDA → HIP → Vulkan → OpenCL → SYCL → CPU
- **Intelligent Library Selection**: Downloader now automatically selects optimal GPU variant based on available tools
- **Expanded Platform Support**: Added comprehensive GPU backend support across Linux and Windows
- **Updated GPU Configuration**: Enhanced context parameters with automatic GPU backend detection
- **Dependencies**: Updated llama.cpp from build b6076 to b6099 (automated via Renovate)
- **Breaking**: Migrated from compilation-based to download-based architecture
- **Simplified build process**: No longer requires CMake, compilers, or GPU SDKs
- Library loading now uses automatic download instead of local compilation
- Updated documentation to reflect new download-based workflow
- **Model download system**: Now uses hf.sh script from llama.cpp instead of direct curl commands
- **Example projects**: Updated to use local hf.sh script from `scripts/` directory
- **Documentation**: Updated to reflect hf.sh script integration and usage

### GPU Backend Support Matrix
| Backend      | Platforms      | GPU Vendors                 | Detection Command | Status       |
| ------------ | -------------- | --------------------------- | ----------------- | ------------ |
| **Metal**    | macOS          | Apple Silicon               | `system_profiler` | ✅ Production |
| **CUDA**     | Linux, Windows | NVIDIA                      | `nvcc`            | ✅ Production |
| **HIP/ROCm** | Linux, Windows | AMD                         | `hipconfig`       | ✅ Production |
| **Vulkan**   | Linux, Windows | NVIDIA, AMD, Intel          | `vulkaninfo`      | ✅ Production |
| **OpenCL**   | Windows, Linux | Qualcomm Adreno, Intel, AMD | `clinfo`          | ✅ Production |
| **SYCL**     | Linux, Windows | Intel, NVIDIA               | `sycl-ls`         | ✅ Production |
| **CPU**      | All            | All                         | N/A               | ✅ Fallback   |

### Dependencies
- llama.cpp: Updated from b6076 to b6099 (managed by Renovate)

### Removed
- All `build-llamacpp-*` compilation targets (no longer needed)
- CMake and compiler dependencies for regular builds
- Complex GPU SDK detection at build time
- `build-libs-gpu` and `build-libs-cpu` targets
- Complete API coverage for llama.cpp functions
- Pre-built llama.cpp libraries for all platforms
- Comprehensive examples and documentation
- GitHub Actions CI/CD pipeline
- Automated release system

### Changed
- **Breaking internal change**: Migrated from direct purego imports to platform-specific abstraction layer
- Separated platform-specific code into `platform_unix.go` and `platform_windows.go` with appropriate build tags
- Updated CI to test cross-compilation for all platforms (Windows, Linux, macOS on both amd64 and arm64)
- Enhanced documentation to reflect platform-specific implementation details

### Fixed
- **Windows CI compilation errors**: Fixed undefined `purego.Dlopen`, `purego.RTLD_NOW`, `purego.RTLD_GLOBAL`, and `purego.Dlclose` symbols
- Cross-compilation now works from any platform to any platform
- Platform detection properly handles unsupported/incomplete platforms

### Features
- Pure Go implementation (no CGO required)
- **Enhanced GPU Support**:
  - Metal support for macOS (Apple Silicon and Intel)
  - CUDA support for NVIDIA GPUs
  - HIP support for AMD GPUs
  - **NEW**: Vulkan support for cross-platform GPU acceleration
  - **NEW**: OpenCL support for diverse GPU vendors (including Qualcomm Adreno)
  - **NEW**: SYCL support for Intel oneAPI and unified parallel programming
  - **NEW**: Automatic GPU backend detection and selection
  - **NEW**: GPU-aware library downloading with optimal variant selection
- Memory mapping and locking options
- Batch processing capabilities
- Multiple sampling strategies
- Model quantization support
- Context state management
- Token manipulation utilities

### Platform Support
- **macOS**: ✅ Intel x64, Apple Silicon (ARM64) with Metal - **Fully supported**
- **Linux**: ✅ x86_64, ARM64 with CUDA/HIP/Vulkan/SYCL - **Fully supported**  
- **Windows**: 🚧 x86_64, ARM64 with CUDA/HIP/Vulkan/OpenCL/SYCL - **Build compatibility implemented, runtime support in development**

### Technical Details
- **Unix-like platforms** (Linux, macOS): Use purego for dynamic library loading
- **Windows platform**: Use native Windows syscalls for library management
- **Build tags**: `!windows` for Unix-like, `windows` for Windows-specific code
- **Zero runtime overhead**: Platform abstraction has no performance impact
- **GPU Detection Priority**: CUDA → HIP → Vulkan → OpenCL → SYCL → CPU (Linux/Windows)
- **Automatic Fallback**: Graceful degradation to CPU when GPU backends unavailable
- **Command Detection**: Uses `exec.LookPath()` for cross-platform command availability
- **Pattern Matching**: Regex-based asset selection for optimal GPU variant downloads

## [0.0.0-llamacpp.b6076] - 2025-01-XX

### Added
- Initial release based on llama.cpp build b6076
- Core llama.cpp API bindings
- Model loading and management
- Context creation and management  
- Text generation and sampling
- Tokenization utilities
- Batch processing
- Memory management
- Error handling
- Cross-platform library loading

### Dependencies
- llama.cpp: b6076
- purego: v0.8.1
- Go: 1.21+

### Platforms
- darwin/amd64 (macOS Intel)
- darwin/arm64 (macOS Apple Silicon)
- linux/amd64 (Linux x86_64)
- linux/arm64 (Linux ARM64)
- windows/amd64 (Windows x86_64)
- windows/arm64 (Windows ARM64)

### Known Issues
- None at initial release

### Breaking Changes
- N/A (initial release)

### Migration Guide
- N/A (initial release)

---

## Version Naming Convention

This project follows the version naming convention:
```
vX.Y.Z-llamacpp.BUILD
```

Where:
- `X.Y.Z` follows semantic versioning for the Go binding
- `BUILD` corresponds to the llama.cpp build number being used

For example: `v1.0.0-llamacpp.b6076` means:
- Go binding version 1.0.0
- Using llama.cpp build b6076

## Release Process

1. Update CHANGELOG.md with new version information
2. Tag the release: `git tag v1.0.0-llamacpp.b6076`
3. Push the tag: `git push origin v1.0.0-llamacpp.b6076`
4. GitHub Actions will automatically build and release binaries
5. Update documentation if needed

## Breaking Changes Policy

- Major version bumps (X.0.0) may include breaking changes
- Minor version bumps (X.Y.0) should be backward compatible
- Patch version bumps (X.Y.Z) are for bug fixes only
- llama.cpp build updates may introduce breaking changes and will be noted

## Support Policy

- Latest version receives full support
- Previous major version receives security updates for 6 months
- Older versions are community-supported only
