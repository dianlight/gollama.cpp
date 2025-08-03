# Gollama.cpp Speculative Decoding Example

This example demonstrates speculative decoding, an advanced technique for accelerating text generation by using a smaller, faster "draft" model to predict multiple tokens ahead, then verifying these predictions with a larger, more accurate "target" model.

## What is Speculative Decoding?

Speculative decoding is an optimization technique that can significantly speed up text generation:

1. **Draft Phase**: A smaller, faster model generates multiple tokens ahead (speculative predictions)
2. **Verify Phase**: The larger, target model checks these predictions in parallel
3. **Accept/Reject**: Matching predictions are accepted; mismatches trigger new sampling from the target model

This approach can provide 2-4x speedup while maintaining the same quality as the target model alone.

## Quick Start

```bash
cd examples/speculative
go run main.go -prompt "The future of AI is"
```

For real acceleration, use two different models:
```bash
go run main.go -model large_model.gguf -draft-model small_model.gguf -prompt "Your prompt"
```

## Features

- **Dual-model speculative decoding** with separate target and draft models
- **Same-model demonstration mode** for understanding the algorithm
- **Configurable draft length** for performance tuning
- **Temperature sampling support** (with fallback to greedy)
- **Detailed statistics** showing acceptance rates and speedup
- **Verbose mode** for observing the draft/verify process
- **Automatic model compatibility checking** (simplified)

## Command Line Options

- `-model string`: Path to the target (main) GGUF model file (default: "../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf")
- `-draft-model string`: Path to the draft (faster) GGUF model file (if empty, uses same as target)
- `-prompt string`: Prompt text to generate from (default: "The future of AI is")
- `-n-predict int`: Number of tokens to predict (default: 100)
- `-n-draft int`: Number of tokens to draft ahead (default: 5)
- `-threads int`: Number of threads to use (default: 4)
- `-ctx int`: Context size (default: 2048)
- `-temperature float`: Sampling temperature, 0.0 = greedy (default: 0.1)
- `-seed int`: Random seed, -1 for random (default: -1)
- `-verbose`: Enable verbose output showing draft/verify process

## Examples

### Basic Speculative Decoding

```bash
# Using same model (demonstration mode)
go run main.go -prompt "Once upon a time" -n-draft 8 -verbose

# With different models for real speedup
go run main.go -model large.gguf -draft-model small.gguf -prompt "Explain quantum computing"
```

### Tuning Draft Length

```bash
# Short drafts (safer, less speedup)
go run main.go -prompt "The benefits of AI" -n-draft 3

# Medium drafts (balanced)
go run main.go -prompt "The benefits of AI" -n-draft 8

# Long drafts (more speedup if accepted, risky)
go run main.go -prompt "The benefits of AI" -n-draft 15
```

### Temperature Sampling

```bash
# Greedy sampling (deterministic)
go run main.go -prompt "Write a story" -temperature 0.0 -n-draft 6

# Low temperature (focused)
go run main.go -prompt "Write a story" -temperature 0.3 -n-draft 6

# Higher temperature (more creative)
go run main.go -prompt "Write a story" -temperature 0.7 -n-draft 6
```

### Performance Analysis

```bash
# Verbose mode to see acceptance rates
go run main.go -prompt "Machine learning" -n-draft 10 -verbose -n-predict 150
```

## Using the Makefile

The included Makefile provides convenient shortcuts:

```bash
# Build the example
make build

# Run with default settings
make run

# Demonstration with same model
make same-model

# Compare different draft lengths
make draft-comparison

# Temperature sampling demo
make temperature-demo

# Creative writing demo
make creative

# Technical explanation demo
make technical

# Full demonstration
make demo

# Test compilation
make test
```

## Understanding the Output

### Normal Output
The example shows generated text in real-time, followed by statistics:

```
Total tokens generated: 100
Draft tokens created: 523
Draft tokens accepted: 67
Acceptance rate: 12.81%
Generation time: 2.3s
Tokens per second: 43.5
```

### Verbose Output
With `-verbose`, you can see the internal process:

```
[DRAFT] Token 0: 284 ('The')
[DRAFT] Token 1: 1108 (' key')
[VERIFY] Draft: 284 ('The'), Target: 284 ('The')
[ACCEPT] Token 0 accepted
[VERIFY] Draft: 1108 (' key'), Target: 2191 (' main')
[REJECT] Token 1 rejected, using target token
```

## Performance Tuning

### Draft Length (`-n-draft`)
- **Small (3-5)**: Conservative, higher acceptance rate, moderate speedup
- **Medium (6-10)**: Balanced approach, good for most use cases
- **Large (12-20)**: Aggressive, lower acceptance rate, higher potential speedup

### Model Selection
For real speedup, use models with significant size differences:
- **Target model**: Large, high-quality model (13B, 30B, 70B parameters)
- **Draft model**: Small, fast model (1B, 3B, 7B parameters)

### Temperature Settings
- **0.0 (Greedy)**: Best for speculative decoding, deterministic
- **0.1-0.3**: Low temperature, good acceptance rates
- **0.5+**: Higher temperature, lower acceptance rates but more creativity

## Model Requirements

### Target Model
- Any high-quality GGUF model
- Larger models benefit more from speculative decoding
- Examples: LLaMA-13B, Mistral-7B, CodeLlama-34B

### Draft Model
- Smaller, faster GGUF model with similar training
- Should have compatible vocabulary with target model
- Examples: TinyLlama-1.1B, Phi-2, smaller Mistral variants

### Compatibility Requirements
- Similar vocabulary and tokenization
- Compatible special tokens (BOS, EOS, etc.)
- Similar training data/domain (for better acceptance rates)

## How It Works

### Algorithm Overview

1. **Initialization**: Load both target and draft models
2. **Prompt Processing**: Both models process the initial prompt
3. **Main Loop**:
   - **Draft Phase**: Draft model generates N tokens ahead
   - **Verify Phase**: Target model verifies each draft token
   - **Accept/Reject**: Matching tokens are accepted, mismatches trigger resampling
4. **Statistics**: Track acceptance rates and performance metrics

### Key Benefits

- **Maintains Quality**: Output is identical to target model alone
- **Parallel Verification**: Multiple draft tokens verified simultaneously
- **Early Termination**: Stop verification at first mismatch
- **Adaptive**: Works with any combination of compatible models

## Troubleshooting

### Low Acceptance Rates
- Use models with more similar training
- Reduce draft length (`-n-draft`)
- Lower temperature (`-temperature`)
- Check model compatibility

### No Speedup Observed
- Ensure you're using different model sizes
- Draft model should be significantly smaller/faster
- Increase draft length if acceptance rate is high
- Check that both models are properly loaded

### Model Compatibility Issues
- Verify both models use the same tokenizer
- Check vocabulary sizes are similar
- Ensure special tokens match
- Use models from the same family/training

### Memory Issues
- Reduce context size (`-ctx`)
- Use smaller models
- Close other applications
- Monitor memory usage

## Advanced Usage

### Custom Model Combinations

```bash
# Large target with tiny draft
./speculative -model llama-30b.gguf -draft-model tinyllama-1b.gguf

# Code model with general draft
./speculative -model codellama-13b.gguf -draft-model llama-7b.gguf

# Chat model with base draft
./speculative -model vicuna-13b.gguf -draft-model llama-7b.gguf
```

### Performance Benchmarking

```bash
# Benchmark different configurations
./speculative -prompt "Long text generation prompt" -n-predict 500 -n-draft 3
./speculative -prompt "Long text generation prompt" -n-predict 500 -n-draft 8
./speculative -prompt "Long text generation prompt" -n-predict 500 -n-draft 15
```

## Limitations

- **Model Compatibility**: Requires compatible vocabularies
- **Memory Usage**: Loads two models simultaneously
- **Implementation Simplification**: This example uses simplified verification
- **Temperature Sampling**: Limited temperature sampling implementation

## Next Steps

After trying this example:

1. **Experiment with Model Pairs**: Find the best target/draft combinations
2. **Optimize Parameters**: Tune draft length and temperature for your use case
3. **Benchmark Performance**: Measure actual speedup with your models
4. **Advanced Sampling**: Implement more sophisticated sampling methods
5. **Tree-based Speculation**: Explore tree-based speculative decoding

## Requirements

- Go 1.21 or later
- Two GGUF model files (or one for demonstration mode)
- Sufficient RAM to load both models
- CPU with multiple cores for optimal performance
