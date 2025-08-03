# Batched Generation Example

This example demonstrates **batched text generation concepts** - showing how multiple independent text sequences could be generated efficiently. This implementation uses a simplified approach to illustrate the principles of batched processing without the complex parallel execution that would require advanced batch management.

**Note:** This is a simplified demonstration of batched generation concepts. A full implementation would use true parallel batch processing that isn't readily available in the current Go bindings.

## What is Batched Generation?

Batched generation allows you to:
- **Generate Multiple Sequences**: Create several different continuations from the same prompt
- **Improve Efficiency**: Process multiple sequences more efficiently than sequential generation  
- **Explore Variations**: See different creative outputs from the same starting point
- **Parallel Processing**: Utilize modern hardware capabilities for better throughput

This is particularly useful for:
- **Creative Writing**: Generate multiple story continuations
- **Code Generation**: Explore different implementation approaches
- **Research & Analysis**: Compare model behavior across multiple runs
- **Content Creation**: Generate variations for A/B testing

## Implementation Notes

This Go implementation demonstrates the batched generation concept using a simplified approach:

1. **Sequential Generation**: For simplicity, sequences are generated one after another
2. **Shared Context**: All sequences use the same initial prompt evaluation
3. **Independent Sampling**: Each sequence uses independent token sampling
4. **Configurable Parameters**: Adjust generation parameters per sequence

The original llama.cpp example uses true parallel batching with complex KV cache management. A full parallel implementation would require:
- Advanced batch manipulation APIs
- KV cache sequence management  
- Parallel logits computation
- Complex state tracking

## Usage Examples

### Basic Multiple Sequences
```bash
# Generate 4 sequences with default settings
./batched -prompt "The future of AI is" -n-parallel 4

# Generate 2 longer sequences
./batched -prompt "Once upon a time" -n-parallel 2 -n-predict 64
```

### Creative Writing
```bash
# Generate story variations
./batched -prompt "In a world where magic exists," -n-parallel 3 -n-predict 50 -temperature 0.9

# Generate different character descriptions  
./batched -prompt "The mysterious stranger was" -n-parallel 5 -n-predict 30
```

### Technical Content
```bash
# Generate code explanations
./batched -prompt "// Function to calculate" -n-parallel 3 -n-predict 40 -temperature 0.3

# Generate documentation variations
./batched -prompt "This API endpoint" -n-parallel 4 -n-predict 25
```

### Parameter Exploration
```bash
# High creativity (high temperature)
./batched -prompt "Innovation means" -temperature 1.2 -n-parallel 4

# Low creativity (low temperature)  
./batched -prompt "Innovation means" -temperature 0.2 -n-parallel 4

# Compare different sampling settings
./batched -prompt "The solution is" -top-k 10 -top-p 0.8 -n-parallel 3
```

### Verbose Mode
```bash
# See detailed processing information
./batched -prompt "Machine learning" -verbose -n-parallel 2 -n-predict 20
```

## Command Line Options

| Option | Default | Description |
|--------|---------|-------------|
| `-model` | `../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf` | Path to GGUF model file |
| `-prompt` | `"Hello my name is"` | Input prompt for all sequences |
| `-n-parallel` | `4` | Number of parallel sequences to generate |
| `-n-predict` | `32` | Number of tokens to predict per sequence |
| `-ctx` | `2048` | Context size for the model |
| `-threads` | `4` | Number of CPU threads to use |
| `-temperature` | `0.8` | Temperature for sampling (0.0-2.0) |
| `-top-k` | `40` | Top-k sampling parameter |
| `-top-p` | `0.9` | Top-p (nucleus) sampling parameter |
| `-verbose` | `false` | Enable detailed output |

## Sampling Parameters

### Temperature
- **0.1-0.3**: Very focused, predictable output
- **0.5-0.8**: Balanced creativity and coherence  
- **0.9-1.2**: More creative and varied output
- **1.3+**: Highly creative but potentially chaotic

### Top-K
- **1-10**: Very constrained vocabulary
- **20-40**: Moderate vocabulary diversity
- **50-100**: High vocabulary diversity

### Top-P  
- **0.1-0.3**: Very focused word selection
- **0.5-0.8**: Balanced word selection
- **0.9-0.95**: Diverse word selection

## Understanding the Output

The example shows:

1. **Individual Sequences**: Each sequence numbered and displayed separately
2. **Performance Statistics**: Token generation speed and throughput
3. **Configuration Summary**: Parameters used for generation
4. **Processing Details**: (in verbose mode) Internal steps and timing

Example output:
```
Sequence 0:

Hello my name is Sarah and I work as a software engineer at a tech startup.

Sequence 1:

Hello my name is Michael, I'm a photographer who loves capturing natural landscapes.

Performance Statistics:
  Decoded 64 tokens in 2.30 seconds
  Speed: 27.83 tokens/second
  Parallel sequences: 2
```

## Performance Considerations

### Memory Usage
- **Model Size**: Larger models require more RAM
- **Context Size**: Larger context requires more memory
- **Batch Size**: More sequences need additional memory

### Processing Speed
- **CPU Threads**: More threads can improve performance
- **Model Complexity**: Simpler models generate faster
- **Sequence Length**: Longer sequences take more time

### Optimization Tips
- Start with smaller models for testing
- Use appropriate context size for your content
- Adjust thread count based on your CPU
- Balance creativity (temperature) with speed

## Comparison with Other Examples

| Example | Purpose | Key Feature |
|---------|---------|-------------|
| **simple-chat** | Single sequence generation | Interactive conversation |
| **batched** | Multiple sequence generation | Parallel processing concept |
| **speculative** | Accelerated generation | Draft-verify algorithm |
| **embedding** | Text representation | Vector embeddings |
| **retrieval** | Document search | Semantic similarity |

## Advanced Usage

### Scripting Multiple Runs
```bash
# Generate multiple batches with different prompts
for prompt in "The weather today" "Technology will" "Art is about"; do
    echo "=== Prompt: $prompt ==="
    ./batched -prompt "$prompt" -n-parallel 3 -n-predict 25
    echo ""
done
```

### Comparing Temperatures
```bash
# Compare low vs high temperature
echo "=== Low Temperature ==="
./batched -prompt "Science is" -temperature 0.2 -n-parallel 3

echo "=== High Temperature ==="  
./batched -prompt "Science is" -temperature 1.0 -n-parallel 3
```

### Performance Testing
```bash
# Test different sequence counts
for count in 1 2 4 8; do
    echo "=== $count sequences ==="
    time ./batched -prompt "Test prompt" -n-parallel $count -n-predict 20
done
```

## Implementation Details

### Sequence Generation Process
1. **Load Model**: Initialize the language model
2. **Tokenize Prompt**: Convert text to token sequence
3. **Setup Context**: Configure generation parameters
4. **Generate Sequences**: Create multiple independent continuations
5. **Collect Results**: Gather and display all sequences

### Batch Management
- Uses simplified sequential processing for clarity
- Each sequence maintains independent state
- Context is reset between sequences for isolation
- Memory management handled automatically

### Sampling Strategy
- Greedy sampling with temperature adjustment
- Top-k and top-p parameters available but simplified
- Could be enhanced with more sophisticated sampling methods

## Troubleshooting

### Model Loading Issues
```bash
# Verify model file exists
ls -la ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf

# Test with verbose output
./batched -verbose -n-parallel 1 -n-predict 5
```

### Memory Problems
```bash
# Reduce context size
./batched -ctx 1024 -n-parallel 2

# Use fewer sequences
./batched -n-parallel 1 -n-predict 10
```

### Performance Issues
```bash
# Adjust thread count
./batched -threads 2

# Use smaller predictions
./batched -n-predict 15 -n-parallel 2
```

### Output Quality
```bash
# Adjust sampling parameters
./batched -temperature 0.5 -top-k 20 -top-p 0.8

# Try different prompts
./batched -prompt "A clear and simple explanation:"
```

## Building and Running

```bash
# Build the example
go build -o batched main.go

# Run with default settings
./batched

# Run with custom settings
./batched -prompt "Your prompt here" -n-parallel 3 -verbose

# Show all options
./batched -help
```

## Model Requirements

- **Text Generation Models**: Any GGUF model that supports text generation
- **Model Types**: LLaMA, Mistral, CodeLlama, etc.
- **Size Considerations**: Larger models produce better quality but require more resources
- **Memory**: Ensure sufficient RAM for model + context + sequences

## Future Enhancements

This example could be extended with:
- **True Parallel Batching**: Implement proper parallel processing
- **Advanced Sampling**: Better top-k/top-p implementation
- **KV Cache Management**: Efficient sequence state handling
- **Streaming Output**: Real-time sequence display
- **Comparison Tools**: Side-by-side sequence analysis
- **Export Options**: Save sequences to files
- **Interactive Mode**: Real-time parameter adjustment

## See Also

- [Simple Chat Example](../simple-chat/) - Single sequence generation
- [Speculative Example](../speculative/) - Accelerated generation
- [Main Examples README](../README.md) - Overview of all examples
- [llama.cpp batched example](../../build/llama.cpp/examples/batched/) - Original C++ implementation
