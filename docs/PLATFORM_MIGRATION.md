# Platform Migration Guide

This document explains the platform-specific architecture changes made to gollama.cpp and the migration from compilation to pre-built binary downloads to improve cross-platform compatibility and eliminate build complexity.

## Overview

We've migrated from a **compilation-based architecture** to a **download-based architecture** with platform-specific code that provides:

- ✅ **No compilation required**: Downloads pre-built binaries from official llama.cpp releases
- ✅ **Windows compatibility**: Full support with runtime library loading
- ✅ **Cross-platform builds**: All platforms work from any host OS  
- ✅ **Better user experience**: No need for CMake, compilers, or GPU SDKs
- ✅ **Always up-to-date**: Uses latest official llama.cpp releases
- ✅ **GPU support**: Automatically selects GPU-enabled binaries when available

## Architecture Changes

### Before (Compilation-Based)

```go
// Required complex build process
make clone-llamacpp
make build-llamacpp-current  // Requires CMake, compilers, GPU SDKs
```

**Problems:**
- ❌ Required CMake, compilers, and build tools
- ❌ Windows compilation failed frequently
- ❌ GPU SDK detection was complex and error-prone
- ❌ Long build times and large repository sizes
- ❌ Dependency hell with different CUDA/HIP versions

### After (Download-Based Architecture)

```go
// Automatic download on first use
import "github.com/dianlight/gollama.cpp"

func main() {
    // Library downloads automatically on first LoadLibrary() call
    gollama.Backend_init() // ✅ Just works!
}
```

**Benefits:**
- ✅ No build dependencies required
- ✅ Uses official pre-built binaries from ggml-org/llama.cpp
- ✅ Automatic platform and GPU variant detection
- ✅ Fast startup and small repository size
- ✅ Always compatible with latest llama.cpp releases

## File Structure

### New Download-Based Files

| File | Purpose |
|------|---------|
| `downloader.go` | Handles downloading pre-built binaries from GitHub releases |
| `loader.go` | Platform-agnostic library loading with download integration |
| `platform_unix.go` | Unix-like systems (Linux, macOS) library loading |
| `platform_windows.go` | Windows systems library loading |
| `cmd/gollama-download/main.go` | Command-line tool for manual library management |

### Library Download Interface

All platforms use the same download interface:

```go
// Core download functions
func NewLibraryDownloader() (*LibraryDownloader, error)
func (d *LibraryDownloader) GetLatestRelease() (*ReleaseInfo, error)
func (d *LibraryDownloader) DownloadAndExtract(url, filename string) (string, error)

// Platform-specific loading
func loadLibraryPlatform(libPath string) (uintptr, error)

// High-level API
func LoadLibraryWithVersion(version string) error
```

## Implementation Details

### Binary Download System

The new architecture downloads appropriate binaries from [ggml-org/llama.cpp releases](https://github.com/ggml-org/llama.cpp/releases):

- **macOS**: `llama-{version}-bin-macos-{arch}.zip` (Metal-enabled)
- **Linux**: `llama-{version}-bin-ubuntu-{arch}.zip` (CPU/CUDA/HIP variants)
- **Windows**: `llama-{version}-bin-win-{variant}-{arch}.zip` (CPU/CUDA/HIP variants)

### Platform Support Matrix

| Platform | Architecture | Binary Variant | Status |
|----------|--------------|----------------|--------|
| macOS | x64 | Metal-enabled | ✅ Fully supported |
| macOS | ARM64 | Metal-enabled | ✅ Fully supported |
| Linux | x64 | CPU/CUDA/HIP | ✅ Fully supported |
| Linux | ARM64 | CPU | ✅ Fully supported |
| Windows | x64 | CPU/CUDA/HIP | ✅ Fully supported |
| Windows | ARM64 | CPU | ✅ Fully supported |

### Automatic Platform Detection

```go
// Automatic platform and variant selection
func (d *LibraryDownloader) GetPlatformAssetPattern() (string, error) {
    switch runtime.GOOS {
    case "darwin":
        return fmt.Sprintf("llama-.*-bin-macos-%s.zip", arch), nil
    case "linux":
        return fmt.Sprintf("llama-.*-bin-ubuntu-%s.zip", arch), nil
    case "windows":
        return fmt.Sprintf("llama-.*-bin-win-cpu-%s.zip", arch), nil
    }
}
```

## Testing Strategy

### Download and Library Tests

```bash
# Test library download functionality
make test-download

# Test downloads for all platforms
make test-download-platforms

# Test cross-compilation
make test-cross-compile

# Full test suite (downloads libraries automatically)
make test
```

### Platform-Specific Runtime Tests

```bash
# Test specific platform downloads
GOOS=windows GOARCH=amd64 go run ./cmd/gollama-download -test-download
GOOS=linux GOARCH=arm64 go run ./cmd/gollama-download -test-download
GOOS=darwin GOARCH=arm64 go run ./cmd/gollama-download -test-download

# Manual library management
go run ./cmd/gollama-download -download -version b6089
go run ./cmd/gollama-download -clean-cache
```

### CI Integration

Our CI now tests:

1. **Cross-compilation** on Ubuntu, macOS, and Windows
2. **Download functionality** for all supported platforms
3. **Platform-specific tests** for library loading
4. **Integration tests** with actual llama.cpp binaries

## Migration Impact

### For Users

- ✅ **Dramatically improved experience**: No build dependencies required
- ✅ **No breaking changes**: Public API remains identical
- ✅ **Faster setup**: Downloads happen automatically and quickly
- ✅ **Better reliability**: Uses official tested binaries
- ✅ **Always up-to-date**: Automatically gets latest llama.cpp versions

### For Contributors

- 📝 **Simplified development**: No need for CMake, compilers, or GPU SDKs
- 🧪 **Faster testing**: Use `make test-download` for quick validation
- 🏗️ **Cleaner codebase**: Removed complex build logic
- 📦 **Smaller repository**: No embedded binaries or build artifacts

### Migration Timeline

1. **Phase 1** ✅ - Download infrastructure (completed)
2. **Phase 2** ✅ - Platform-specific loading (completed)  
3. **Phase 3** ✅ - Command-line tools (completed)
4. **Phase 4** ✅ - Documentation updates (completed)
5. **Phase 5** 📋 - Automated version tracking with Renovate

## Future Roadmap

### Automated Version Tracking

The next phase involves implementing automated tracking of llama.cpp releases:

1. **Renovate Integration** 📋 - Automatic PR creation for new llama.cpp releases
2. **Version Validation** � - Automated testing of new binary releases
3. **Release Automation** 📋 - Automated gollama.cpp releases when llama.cpp updates
4. **GPU Variant Selection** 📋 - Intelligent GPU binary selection based on system capabilities

### Enhanced Library Management

Future improvements to the download system:

```go
// Future API enhancements
gollama.LoadLibraryWithOptions(&gollama.LoadOptions{
    Version:    "b6089",
    GPUVariant: "cuda",  // cuda, hip, cpu, auto
    Cache:      true,
    Verify:     true,
})
```

### Additional Platforms

The architecture supports easy extension to new platforms as they become available in llama.cpp releases.

## Performance Impact

- **Compilation**: ✅ Eliminated (no build time)
- **Download**: ✅ Fast initial setup (caching prevents re-download)
- **Runtime**: ✅ Zero overhead (same performance as compiled libraries)  
- **Binary size**: ✅ Smaller (downloads only needed platforms)
- **Memory usage**: ✅ Unchanged (same binaries, different delivery method)

## Troubleshooting

### Download Issues

```bash
# Test connectivity and platform detection
go run ./cmd/gollama-download -test-download

# Force cache refresh
go run ./cmd/gollama-download -clean-cache
go run ./cmd/gollama-download -download

# Check specific version
go run ./cmd/gollama-download -download -version b6089
```

### Platform Detection

The library automatically detects platform capabilities:

```go
if gollama.IsPlatformSupported() {
    // Platform has full support
} else {
    // Platform may have limited support
    fmt.Println("Error:", gollama.GetPlatformError())
}
```

## Conclusion

This migration from compilation-based to download-based architecture provides a dramatically improved developer experience while maintaining full compatibility and performance. The elimination of build dependencies and the use of official pre-built binaries ensures reliability and reduces maintenance overhead for both users and contributors.
