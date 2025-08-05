# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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
- **Breaking**: Migrated from compilation-based to download-based architecture
- **Simplified build process**: No longer requires CMake, compilers, or GPU SDKs
- Library loading now uses automatic download instead of local compilation
- Updated documentation to reflect new download-based workflow
- **Model download system**: Now uses hf.sh script from llama.cpp instead of direct curl commands
- **Example projects**: Updated to use local hf.sh script from `scripts/` directory
- **Documentation**: Updated to reflect hf.sh script integration and usage

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
- Metal support for macOS (Apple Silicon and Intel)
- CUDA support for NVIDIA GPUs  
- HIP support for AMD GPUs
- Memory mapping and locking options
- Batch processing capabilities
- Multiple sampling strategies
- Model quantization support
- Context state management
- Token manipulation utilities

### Platform Support
- **macOS**: ✅ Intel x64, Apple Silicon (ARM64) with Metal - **Fully supported**
- **Linux**: ✅ x86_64, ARM64 with CUDA/HIP - **Fully supported**  
- **Windows**: 🚧 x86_64, ARM64 with CUDA - **Build compatibility implemented, runtime support in development**

### Technical Details
- **Unix-like platforms** (Linux, macOS): Use purego for dynamic library loading
- **Windows platform**: Use native Windows syscalls for library management
- **Build tags**: `!windows` for Unix-like, `windows` for Windows-specific code
- **Zero runtime overhead**: Platform abstraction has no performance impact

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
