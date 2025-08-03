# Gollama.cpp Evaluation Callback Example

This example demonstrates how evaluation callbacks work in llama.cpp and simulates what they would show during model inference. Evaluation callbacks allow you to inspect tensor operations, monitor graph computation, and debug model execution in real-time.

## Overview

The eval-callback example shows:

1. **Tensor Operation Monitoring**: Track each operation during graph execution
2. **Memory Location Tracking**: Monitor whether tensors are in CPU or GPU memory  
3. **Performance Profiling**: Measure operation timing and throughput
4. **Data Inspection**: View tensor shapes, sizes, and data values
5. **Progress Reporting**: Real-time updates during inference

## Features

### Simulation Mode
- **Operation Logging**: Simulates ggml_debug callback output showing tensor operations
- **Tensor Information**: Displays tensor names, types, dimensions, and memory usage
- **Performance Metrics**: Tracks operation counts, data throughput, and timing
- **Memory Monitoring**: Shows CPU vs GPU memory location for each tensor

### Real Model Mode
- **Actual Inference**: Runs real model evaluation alongside simulation
- **Performance Timing**: Measures actual evaluation speed and token processing
- **Memory Usage**: Shows real memory consumption and processing statistics

## Usage

### Basic Usage
```bash
# Run with simulation only (no model loading)
go run main.go -simulate-only

# Run with actual model evaluation
go run main.go -model ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf

# Custom prompt
go run main.go -prompt "Explain neural networks" -simulate-only
```

### Advanced Options
```bash
# Enable verbose tensor data printing
go run main.go -print-tensor-data -max-logged-ops 20

# Disable operation logging, show only progress
go run main.go -enable-logging=false -enable-progress

# Custom context and threading
go run main.go -ctx 1024 -threads 8
```

## Command Line Options

| Option | Default | Description |
|--------|---------|-------------|
| `-model` | `../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf` | Path to GGUF model file |
| `-prompt` | `"The future of AI is"` | Text prompt to evaluate |
| `-threads` | `4` | Number of threads for computation |
| `-ctx` | `512` | Context size for evaluation |
| `-enable-logging` | `true` | Enable operation logging (simulates ggml_debug) |
| `-enable-progress` | `true` | Enable progress updates during evaluation |
| `-print-tensor-data` | `false` | Print tensor data values (verbose output) |
| `-max-logged-ops` | `50` | Maximum operations to log (0 = unlimited) |
| `-simulate-only` | `false` | Only run simulation without model inference |

## Example Output

### Operation Logging
```
=== Starting Evaluation with Callbacks ===
Tokens to process: 5
Logging enabled: true

ggml_debug:                attn_q_layer_0 = (f32)    MUL_MAT(inp_layer_0{5, 2048}, wq_0{2048, 2048}) = {5, 2048}
                              └─ Op #1, 40960 bytes, GPU memory, 1.234ms since start, 1.234ms since last

ggml_debug:                attn_k_layer_0 = (f32)    MUL_MAT(inp_layer_0{5, 2048}, wk_0{2048, 2048}) = {5, 2048}
                              └─ Op #2, 40960 bytes, GPU memory, 2.456ms since start, 1.222ms since last

ggml_debug:            attn_scores_layer_0 = (f32)    MUL_MAT(attn_q_layer_0{5, 2048}, attn_k_layer_0{5, 2048}) = {5, 5}
                              └─ Op #4, 100 bytes, CPU memory, 3.789ms since start, 1.333ms since last
```

### Progress Updates
```
Progress Update:
  Operations processed: 134
  Tensors processed: 134
  Data processed: 45.2 MB
  Elapsed time: 0.89s
  Average ops/sec: 150.6
  Average throughput: 50.8 MB/sec
```

### Performance Summary
```
=== Performance Information ===
Evaluation time: 234.56 ms
Tokens processed: 5
Processing speed: 21.3 tokens/s
```

## Implementation Details

### Callback Simulation
The example simulates what real eval callbacks would show:

```go
type EvalCallbackData struct {
    OperationCount  int
    TensorCount     int
    BytesProcessed  int64
    StartTime       time.Time
    EnableLogging   bool
    // ... other fields
}
```

### Tensor Information
Each operation logs detailed tensor information:

```go
type SimulatedTensorInfo struct {
    Name       string    // Tensor name (e.g., "attn_q_layer_0")
    Type       string    // Tensor type (e.g., "attn_q")
    Operation  string    // GGML operation (e.g., "MUL_MAT")
    Dimensions []int64   // Tensor dimensions
    SizeBytes  int64     // Memory size in bytes
    IsHost     bool      // CPU vs GPU memory
    DataType   string    // Data type (f32, f16, etc.)
}
```

### Real Implementation Notes
In a real callback implementation, you would:

1. **Set Callback Pointers**: 
   ```go
   ctxParams.CbEval = callbackFunctionPointer
   ctxParams.CbEvalUserData = callbackDataPointer
   ```

2. **Implement C Callback**: Use CGO to create a C callback function that can access tensor data

3. **Access Tensor Data**: Extract actual tensor values, names, and operations during graph execution

4. **Memory Management**: Handle GPU memory transfers and data copying

## Model Operations

The simulation demonstrates typical transformer operations:

1. **Attention Layers**: Query, Key, Value matrix multiplications
2. **Attention Computation**: Score calculation and softmax
3. **Feed-Forward Networks**: Gate and up projections
4. **Layer Normalization**: Input normalization operations
5. **Output Layer**: Final logits computation

## Educational Value

This example helps understand:

- **Graph Execution Flow**: How operations are executed in sequence
- **Memory Management**: CPU vs GPU tensor placement
- **Performance Characteristics**: Operation timing and throughput
- **Debugging Techniques**: How to monitor inference in real-time
- **Tensor Operations**: The low-level operations that make up inference

## Build and Run

```bash
# Install dependencies
go mod tidy

# Build
go build -o eval-callback

# Run simulation
./eval-callback -simulate-only

# Run with model
./eval-callback -model path/to/your/model.gguf
```

## Integration with Other Examples

This example complements other examples in the suite:

- **Simple Chat**: Shows basic inference without debugging
- **Embedding**: Demonstrates embedding extraction
- **Batched**: Shows parallel processing concepts
- **Speculative**: Demonstrates advanced generation techniques

## Performance Considerations

### Callback Overhead
Real eval callbacks introduce performance overhead:

- **Frequent Calls**: Called for every tensor operation
- **Memory Transfers**: GPU to CPU data copying
- **Logging Overhead**: String formatting and output

### Optimization Tips
- **Selective Logging**: Only log important operations
- **Batch Processing**: Group operations before logging
- **Async Logging**: Use separate goroutines for output
- **Memory Pooling**: Reuse callback data structures

## Troubleshooting

### Common Issues

1. **No Callback Output**: Callbacks require proper C function pointers
2. **Memory Errors**: GPU memory access needs proper synchronization  
3. **Performance Impact**: Excessive logging can slow inference
4. **Type Mismatches**: Ensure correct tensor type handling

### Debugging Tips

- Use simulation mode first to understand expected output
- Start with limited operation logging
- Check memory location (CPU vs GPU) for data access
- Monitor performance impact of callback overhead

## Related Documentation

- [llama.cpp eval-callback example](../../build/llama.cpp/examples/eval-callback/)
- [GGML tensor operations documentation](https://github.com/ggerganov/ggml)
- [Callback implementation guide](../README.md)
