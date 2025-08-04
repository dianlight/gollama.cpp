# Platform Migration Guide

This document explains the platform-specific architecture changes made to gollama.cpp to resolve Windows compilation issues and improve cross-platform compatibility.

## Overview

We've migrated from a single-platform implementation to a **build tag-based platform-specific architecture** that provides:

- ✅ **Windows compatibility**: Fixed compilation errors on Windows CI
- ✅ **Cross-platform builds**: All platforms can build from any host OS  
- ✅ **Better abstraction**: Platform-specific code is cleanly separated
- ✅ **Future extensibility**: Easy to add new platform support

## Architecture Changes

### Before (Single Implementation)

```go
// gollama.go - Direct purego usage
import "github.com/ebitengine/purego"

func loadLibrary() error {
    if runtime.GOOS == "windows" {
        return errors.New("not implemented") // ❌ Compilation error
    }
    handle, err := purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
    // ...
}
```

**Problems:**
- ❌ Windows compilation failed due to undefined `purego` symbols
- ❌ Platform logic mixed with business logic
- ❌ Hard to extend for new platforms

### After (Platform-Specific Architecture)

```go
// platform_unix.go
//go:build !windows

package gollama
import "github.com/ebitengine/purego"

func loadLibraryPlatform(libPath string) (uintptr, error) {
    return purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}
```

```go
// platform_windows.go  
//go:build windows

package gollama
import "syscall"

func loadLibraryPlatform(libPath string) (uintptr, error) {
    // Windows-specific implementation using LoadLibraryW
}
```

```go
// gollama.go - Platform-agnostic
func loadLibrary() error {
    if !isPlatformSupported() {
        return getPlatformError()
    }
    handle, err := loadLibraryPlatform(libPath) // ✅ Clean abstraction
    // ...
}
```

## File Structure

### New Platform-Specific Files

| File | Build Tag | Purpose |
|------|-----------|---------|
| `platform_unix.go` | `!windows` | Unix-like systems (Linux, macOS) |
| `platform_windows.go` | `windows` | Windows systems |
| `platform_test.go` | none | Cross-platform tests |

### Platform Function Interface

All platforms must implement:

```go
// Core platform functions
func loadLibraryPlatform(libPath string) (uintptr, error)
func closeLibraryPlatform(handle uintptr) error  
func registerLibFunc(fptr interface{}, handle uintptr, fname string)

// Platform capability detection
func isPlatformSupported() bool
func getPlatformError() error
```

## Implementation Details

### Unix-like Platforms (Linux, macOS)

- **Build tag**: `//go:build !windows`
- **Implementation**: Uses [purego](https://github.com/ebitengine/purego) for FFI
- **Status**: ✅ Fully supported
- **Features**: Complete llama.cpp binding support

### Windows Platform  

- **Build tag**: `//go:build windows`
- **Implementation**: Uses native Windows syscalls (`LoadLibraryW`, `FreeLibrary`)
- **Status**: 🚧 Build compatibility implemented, runtime support in development
- **Current capabilities**:
  - ✅ Compiles without errors
  - ✅ Cross-compilation works
  - 🚧 Function registration placeholder
  - 🚧 Full runtime support coming soon

## Testing Strategy

### Cross-Platform Compilation Tests

```bash
# Test all platform builds
make test-cross-compile

# Test specific platforms
make test-compile-windows
make test-compile-linux  
make test-compile-darwin
```

### Platform-Specific Runtime Tests

```bash
# Run platform capability tests
make test-platform

# Full test suite (builds native libraries)
make test
```

### CI Integration

Our CI now tests:

1. **Native compilation** on Ubuntu, macOS, and Windows
2. **Cross-compilation** from Linux to all platforms
3. **Platform-specific tests** for capability detection
4. **Race detection** on supported platforms

## Migration Impact

### For Users

- ✅ **No breaking changes**: Public API remains identical
- ✅ **Better reliability**: Windows builds no longer fail
- ✅ **Same performance**: No runtime overhead added

### For Contributors

- 📝 **New guidelines**: See [CONTRIBUTING.md](../CONTRIBUTING.md) for platform-specific development
- 🧪 **Enhanced testing**: Use `make test-cross-compile` before submitting
- 🏗️ **Build tags**: Understand when to use platform-specific files

## Future Roadmap

### Windows Support Completion

1. **Phase 1** ✅ - Build compatibility (completed)
2. **Phase 2** 🚧 - Function registration via `GetProcAddress`  
3. **Phase 3** 📋 - Full runtime testing and GPU support
4. **Phase 4** 📋 - Windows-specific optimizations

### Additional Platforms

The architecture supports easy extension:

```go
// platform_freebsd.go
//go:build freebsd

// platform_wasm.go  
//go:build js,wasm
```

## Performance Impact

- **Compilation**: ✅ No impact (build tags eliminate unused code)
- **Runtime**: ✅ Zero overhead (function calls are identical)  
- **Binary size**: ✅ Smaller (only relevant platform code included)

## Troubleshooting

### Build Issues

```bash
# Verify platform detection
go test -v -run TestPlatformSpecific

# Test cross-compilation
GOOS=windows GOARCH=amd64 go build -v ./...

# Clean and rebuild  
make clean && make build
```

### Platform Detection

The library automatically detects platform capabilities:

```go
if gollama.IsPlatformSupported() {
    // Platform has full support
} else {
    // Platform has limited or no support
    fmt.Println("Error:", gollama.GetPlatformError())
}
```

## Conclusion

This migration provides a robust foundation for multi-platform support while maintaining backward compatibility and preparing for future platform additions. The clean separation of platform-specific code makes the library more maintainable and extensible.
