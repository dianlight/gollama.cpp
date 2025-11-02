# GGML Low-Level API Documentation

This document describes the GGML (Georgi Gerganov Machine Learning) low-level tensor library bindings available in gollama.cpp.

## Overview

GGML is the tensor library that powers llama.cpp. While most users will interact with the high-level llama.cpp API (`gollama.go`), the GGML bindings (`goggml.go`) provide direct access to low-level tensor operations and backend management for advanced use cases.

**Important Note:** GGML functions may not be exported in all llama.cpp builds. The library gracefully handles missing functions without errors, allowing code to compile and run even when GGML symbols are not available.

## When to Use GGML Bindings

Use GGML bindings when you need:

- **Type Information**: Query tensor type sizes, block sizes, or quantization status
- **Backend Management**: Enumerate available compute backends (CPU, GPU, etc.)
- **Memory Management**: Direct buffer allocation and management
- **Quantization**: Access to low-level quantization utilities
- **Advanced Integration**: Building custom tensor operations or tools

For most LLM inference tasks, use the high-level llama.cpp API in `gollama.go`.

## Available Features

### Tensor Types (31 types)

GGML supports various data types for tensors:

#### Floating Point Types
- `GGML_TYPE_F32` - 32-bit float (4 bytes)
- `GGML_TYPE_F16` - 16-bit float (2 bytes)
- `GGML_TYPE_F64` - 64-bit float (8 bytes)
- `GGML_TYPE_BF16` - BFloat16 (2 bytes)

#### Integer Types
- `GGML_TYPE_I8` - 8-bit integer (1 byte)
- `GGML_TYPE_I16` - 16-bit integer (2 bytes)
- `GGML_TYPE_I32` - 32-bit integer (4 bytes)
- `GGML_TYPE_I64` - 64-bit integer (8 bytes)

#### Quantized Types (K-quants)
- `GGML_TYPE_Q4_0`, `GGML_TYPE_Q4_1` - 4-bit quantization
- `GGML_TYPE_Q5_0`, `GGML_TYPE_Q5_1` - 5-bit quantization
- `GGML_TYPE_Q8_0`, `GGML_TYPE_Q8_1` - 8-bit quantization
- `GGML_TYPE_Q2_K` - 2-bit K-quant
- `GGML_TYPE_Q3_K` - 3-bit K-quant
- `GGML_TYPE_Q4_K` - 4-bit K-quant
- `GGML_TYPE_Q5_K` - 5-bit K-quant
- `GGML_TYPE_Q6_K` - 6-bit K-quant
- `GGML_TYPE_Q8_K` - 8-bit K-quant

#### Importance Quantization (IQ)
- `GGML_TYPE_IQ1_S`, `GGML_TYPE_IQ1_M` - 1-bit importance quantization
- `GGML_TYPE_IQ2_XXS`, `GGML_TYPE_IQ2_XS`, `GGML_TYPE_IQ2_S` - 2-bit IQ variants
- `GGML_TYPE_IQ3_XXS`, `GGML_TYPE_IQ3_S` - 3-bit IQ variants
- `GGML_TYPE_IQ4_NL`, `GGML_TYPE_IQ4_XS` - 4-bit IQ variants

## API Functions

### Type Utilities

#### Ggml_type_size
```go
func Ggml_type_size(typ GgmlType) (uint64, error)
```
Returns the size in bytes of a GGML type element.

**Example:**
```go
size, err := gollama.Ggml_type_size(gollama.GGML_TYPE_F32)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("F32 size: %d bytes\n", size) // Output: F32 size: 4 bytes
```

#### Ggml_blck_size
```go
func Ggml_blck_size(typ GgmlType) (int32, error)
```
Returns the block size of a GGML type (relevant for quantized types).

#### Ggml_type_is_quantized
```go
func Ggml_type_is_quantized(typ GgmlType) (bool, error)
```
Returns whether a GGML type is quantized.

**Example:**
```go
isQuantized, err := gollama.Ggml_type_is_quantized(gollama.GGML_TYPE_Q4_0)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Q4_0 is quantized: %v\n", isQuantized) // Output: Q4_0 is quantized: true
```

#### Ggml_type_name
```go
func Ggml_type_name(typ GgmlType) (string, error)
```
Returns the string name of a GGML type.

**Example:**
```go
name, err := gollama.Ggml_type_name(gollama.GGML_TYPE_F32)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Type name: %s\n", name) // Output: Type name: f32
```

### Backend Device Management

#### Ggml_backend_dev_count
```go
func Ggml_backend_dev_count() (uint64, error)
```
Returns the number of available backend devices.

**Example:**
```go
count, err := gollama.Ggml_backend_dev_count()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d backend device(s)\n", count)
```

#### Ggml_backend_dev_get
```go
func Ggml_backend_dev_get(index uint64) (GgmlBackendDevice, error)
```
Returns a backend device by index.

#### Ggml_backend_dev_name
```go
func Ggml_backend_dev_name(device GgmlBackendDevice) (string, error)
```
Returns the name of a backend device.

**Example:**
```go
count, _ := gollama.Ggml_backend_dev_count()
for i := uint64(0); i < count; i++ {
    dev, err := gollama.Ggml_backend_dev_get(i)
    if err != nil {
        continue
    }
    name, err := gollama.Ggml_backend_dev_name(dev)
    if err != nil {
        continue
    }
    fmt.Printf("Device %d: %s\n", i, name)
}
```

#### Ggml_backend_dev_description
```go
func Ggml_backend_dev_description(device GgmlBackendDevice) (string, error)
```
Returns the description of a backend device.

#### Ggml_backend_dev_memory
```go
func Ggml_backend_dev_memory(device GgmlBackendDevice) (free uint64, total uint64, err error)
```
Returns the memory statistics of a backend device (free and total memory in bytes).

### Buffer Management

#### Ggml_backend_cpu_buffer_type
```go
func Ggml_backend_cpu_buffer_type() (GgmlBackendBufferType, error)
```
Returns the CPU buffer type.

#### Ggml_backend_buffer_free
```go
func Ggml_backend_buffer_free(buffer GgmlBackendBuffer) error
```
Frees a backend buffer.

#### Ggml_backend_buffer_get_size
```go
func Ggml_backend_buffer_get_size(buffer GgmlBackendBuffer) (uint64, error)
```
Returns the size of a backend buffer in bytes.

#### Ggml_backend_buffer_is_host
```go
func Ggml_backend_buffer_is_host(buffer GgmlBackendBuffer) (bool, error)
```
Checks if a buffer is in host memory (RAM).

#### Ggml_backend_buffer_name
```go
func Ggml_backend_buffer_name(buffer GgmlBackendBuffer) (string, error)
```
Returns the name of a backend buffer.

## Complete Example

Here's a comprehensive example using GGML bindings:

```go
package main

import (
    "fmt"
    "log"

    "github.com/dianlight/gollama.cpp"
)

func main() {
    // Initialize the library
    if err := gollama.Backend_init(); err != nil {
        log.Fatal(err)
    }
    defer gollama.Backend_free()

    // Query type information
    fmt.Println("=== Type Information ===")
    types := []gollama.GgmlType{
        gollama.GGML_TYPE_F32,
        gollama.GGML_TYPE_F16,
        gollama.GGML_TYPE_Q4_0,
        gollama.GGML_TYPE_Q8_0,
    }

    for _, typ := range types {
        // Get type size
        size, err := gollama.Ggml_type_size(typ)
        if err != nil {
            fmt.Printf("Type %s: size unavailable\n", typ.String())
            continue
        }

        // Check if quantized
        isQuant, _ := gollama.Ggml_type_is_quantized(typ)

        // Get type name
        name, _ := gollama.Ggml_type_name(typ)

        fmt.Printf("Type: %-10s | Size: %2d bytes | Quantized: %v | Name: %s\n",
            typ.String(), size, isQuant, name)
    }

    // Enumerate backend devices
    fmt.Println("\n=== Backend Devices ===")
    count, err := gollama.Ggml_backend_dev_count()
    if err != nil {
        fmt.Println("Backend device enumeration not available")
        return
    }

    if count == 0 {
        fmt.Println("No backend devices available")
        return
    }

    for i := uint64(0); i < count; i++ {
        dev, err := gollama.Ggml_backend_dev_get(i)
        if err != nil {
            continue
        }

        name, err := gollama.Ggml_backend_dev_name(dev)
        if err != nil {
            continue
        }

        desc, _ := gollama.Ggml_backend_dev_description(dev)
        fmt.Printf("Device %d: %s\n", i, name)
        if desc != "" {
            fmt.Printf("  Description: %s\n", desc)
        }

        // Try to get memory info (may not be supported)
        free, total, err := gollama.Ggml_backend_dev_memory(dev)
        if err == nil {
            fmt.Printf("  Memory: %.2f MB free / %.2f MB total\n",
                float64(free)/(1024*1024),
                float64(total)/(1024*1024))
        }
    }
}
```

**Expected Output:**
```
=== Type Information ===
Type: f32        | Size:  4 bytes | Quantized: false | Name: f32
Type: f16        | Size:  2 bytes | Quantized: false | Name: f16
Type: q4_0       | Size:  2 bytes | Quantized: true  | Name: q4_0
Type: q8_0       | Size:  1 bytes | Quantized: true  | Name: q8_0

=== Backend Devices ===
Device 0: CPU
  Description: CPU backend
```

## Type Conversions

The `GgmlType` enum provides a `String()` method for easy display:

```go
typ := gollama.GGML_TYPE_Q4_0
fmt.Println(typ.String()) // Output: q4_0
```

## Error Handling

All GGML functions return an error that should be checked:

```go
size, err := gollama.Ggml_type_size(gollama.GGML_TYPE_F32)
if err != nil {
    // Function not available in this build or library not loaded
    log.Printf("Warning: %v", err)
    return
}
// Use size...
```

## Testing

The GGML bindings include comprehensive tests in `goggml_test.go`:

```bash
# Run all GGML tests
go test -v -run TestGgml

# Run specific test
go test -v -run TestGgmlTypeSize

# Run benchmarks
go test -v -bench=BenchmarkGgml
```

## Limitations

1. **Optional Functions**: GGML functions may not be exported in all llama.cpp builds. The library handles this gracefully by returning errors instead of panicking.

2. **Platform Differences**: Some functions may have different behavior or availability across platforms.

3. **Build Variants**: Different llama.cpp builds (CPU-only vs GPU-enabled) may export different GGML symbols.

4. **Version Compatibility**: GGML API may change between llama.cpp versions. Always use the version of gollama.cpp that matches your llama.cpp build.

## Related Documentation

- [Main README](../README.md) - High-level overview and quick start
- [Build Guide](BUILD.md) - Building from source
- [GPU Setup](GPU.md) - GPU acceleration configuration
- [API Reference](https://pkg.go.dev/github.com/dianlight/gollama.cpp) - Full Go API documentation

## Support

If you encounter issues with GGML bindings:

1. Check that your llama.cpp build exports GGML symbols
2. Verify you're using a compatible gollama.cpp version
3. Report issues at: https://github.com/dianlight/gollama.cpp/issues
