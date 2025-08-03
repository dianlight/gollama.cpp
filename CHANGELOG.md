# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial Go binding for llama.cpp using purego
- Cross-platform support (macOS, Linux, Windows)
- CPU and GPU acceleration support
- Complete API coverage for llama.cpp functions
- Pre-built llama.cpp libraries for all platforms
- Comprehensive examples and documentation
- GitHub Actions CI/CD pipeline
- Automated release system

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
- **macOS**: Intel x64, Apple Silicon (ARM64) with Metal
- **Linux**: x86_64, ARM64 with CUDA/HIP
- **Windows**: x86_64, ARM64 with CUDA

## [1.0.0-llamacpp.b6076] - 2025-01-XX

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
