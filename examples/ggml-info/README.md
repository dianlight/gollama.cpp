# GGML Info Example

This example demonstrates how to use the GGML low-level API bindings to query type information and enumerate backend devices.

## What it Does

- Queries information about various GGML tensor types (F32, F16, Q4_0, Q8_0, etc.)
- Shows type sizes, whether types are quantized, and type names
- Enumerates available backend devices (CPU, GPU, etc.)
- Displays device descriptions and memory information when available

## Running the Example

```bash
# From the examples/ggml-info directory
go run main.go

# Or build and run
go build -o ggml-info
./ggml-info
```

## Expected Output

```
GGML Low-Level API Demo
========================

=== GGML Type Information ===
Type         | Size       | Quantized  | Name      
-------------|------------|------------|------------
f32          | 4 bytes    | false      | f32       
f16          | 2 bytes    | false      | f16       
bf16         | 2 bytes    | false      | bf16      
q4_0         | 2 bytes    | true       | q4_0      
q8_0         | 1 bytes    | true       | q8_0      
iq2_xxs      | N/A        | true       | iq2_xxs   
iq4_xs       | N/A        | true       | iq4_xs    
i32          | 4 bytes    | false      | i32       

=== Backend Devices ===
Backend device enumeration not available in this build
(GGML functions may not be exported)
```

**Note:** GGML functions may not be exported in all llama.cpp builds. If backend enumeration is not available, the example will gracefully handle this and continue.

## Related Documentation

- [GGML API Documentation](../../docs/GGML_API.md)
- [Main README](../../README.md)
