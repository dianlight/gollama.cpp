# Diffusion Generation Example

A conceptual demonstration of **diffusion-based text generation principles** using gollama.cpp bindings to llama.cpp. This implementation shows the core concepts of Diffusion Language Models (DLLMs) while working within the constraints of standard chat model architectures.

## Overview

This example demonstrates the principles of diffusion text generation - an advanced technique that generates text by iteratively refining masked tokens through multiple denoising steps. While true diffusion models require specialized architectures with non-causal attention, this implementation illustrates the key concepts using available APIs.

**Note:** This is a conceptual demonstration of diffusion principles. A full implementation would require specialized diffusion model architectures (like Dream or LLaDA) and non-causal attention mechanisms not available in standard chat models.

## What is Diffusion Text Generation?

Diffusion text generation is inspired by diffusion models used in image generation. The process works by:

1. **Starting with masked tokens**: Begin with a sequence where some positions are masked/unknown
2. **Iterative refinement**: Through multiple steps, gradually replace masked tokens with predicted tokens
3. **Confidence-based selection**: Use various algorithms to determine which tokens to unmask at each step
4. **Progressive denoising**: Each step reduces uncertainty and improves the overall sequence quality

### Key Concepts

- **Masked Language Modeling**: Working with partially masked sequences
- **Non-causal Attention**: Attending to both past and future tokens (bidirectional)
- **Iterative Refinement**: Multiple denoising steps rather than single-pass generation
- **Confidence Algorithms**: Different strategies for selecting which tokens to unmask
- **Transfer Scheduling**: Controlling how many tokens to unmask per step

## Diffusion Algorithms

This implementation supports several confidence-based algorithms for token selection:

### 1. Confidence-Based (Default)
Uses the probability of the selected token as confidence measure.

### 2. Entropy-Based  
Calculates confidence based on the entropy of the probability distribution.

### 3. Margin-Based
Uses the difference between the top two token probabilities.

### 4. Random
Randomly selects tokens to unmask (for comparison/ablation studies).

## Example Output

```
Gollama.cpp Diffusion Generation Example v1.0.0-llamacpp.b6076
Configuration:
  Model: ../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf
  Prompt: "The future of AI is"
  Diffusion steps: 10
  Max length: 64
  Algorithm: CONFIDENCE_BASED
  Context size: 2048
  Threads: 4
  Temperature: 0.80
  Top-K: 40
  Top-P: 0.90
  Epsilon: 0.010000
  Seed: 1691234567890
  Visual mode: false

Loading model...
Tokenizing prompt...
Prompt tokens: 4

Starting diffusion with 60 masked positions
Step 1: unmasking 6 tokens (remaining: 60)
Step 2: unmasking 7 tokens (remaining: 54)
...
Diffusion step: 10/10 [==================================================] 100%

Diffusion completed in 2.34 seconds

Diffusion Generation Complete!
Generated text:

The future of AI is bright and full of innovative possibilities for humanity.

Generation Summary:
  Diffusion steps: 10
  Algorithm: CONFIDENCE_BASED
  Generated tokens: 11
  Total length: 56 characters

Note: This is a conceptual demonstration of diffusion principles.
A full implementation would require specialized diffusion model architectures
and non-causal attention mechanisms not available in standard chat models.
```

## Visual Mode

Enable visual mode to see the generation process in real-time:

```bash
./diffusion --diffusion-visual --prompt "Once upon a time"
```

Visual mode shows:
- Progress bar for current diffusion step
- Current state of the sequence with masked positions shown as underscores
- Real-time updates as tokens are progressively unmasked

## Usage

### Basic Generation
```bash
# Simple diffusion generation
./diffusion --prompt "The future of technology"

# With custom parameters
./diffusion --prompt "Machine learning" --diffusion-steps 15 --max-length 80
```

### Algorithm Comparison
```bash
# Confidence-based algorithm (default)
./diffusion --prompt "Hello world" --diffusion-algorithm 0

# Entropy-based algorithm
./diffusion --prompt "Hello world" --diffusion-algorithm 1

# Margin-based algorithm  
./diffusion --prompt "Hello world" --diffusion-algorithm 2

# Random algorithm (for comparison)
./diffusion --prompt "Hello world" --diffusion-algorithm 3
```

### Visual and Verbose Modes
```bash
# Visual mode with real-time generation display
./diffusion --prompt "Science fiction" --diffusion-visual

# Verbose mode with detailed step information
./diffusion --prompt "Technology trends" --verbose

# Combined visual and verbose
./diffusion --prompt "Future predictions" --diffusion-visual --verbose
```

### Advanced Parameters
```bash
# Custom diffusion parameters
./diffusion \
  --prompt "Artificial intelligence" \
  --diffusion-steps 20 \
  --diffusion-eps 0.005 \
  --max-length 100 \
  --temperature 0.9

# Deterministic generation with seed
./diffusion --prompt "Deterministic output" --seed 12345
```

## Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `--model` | Path to GGUF model file | `../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf` |
| `--prompt` | Input prompt for generation | `"The future of AI is"` |
| `--diffusion-steps` | Number of diffusion steps | `10` |
| `--diffusion-algorithm` | Algorithm (0-3) | `0` (confidence-based) |
| `--diffusion-eps` | Epsilon parameter for scheduling | `0.01` |
| `--diffusion-visual` | Enable visual mode | `false` |
| `--max-length` | Maximum sequence length | `64` |
| `--temperature` | Sampling temperature | `0.8` |
| `--top-k` | Top-K sampling | `40` |
| `--top-p` | Top-P sampling | `0.9` |
| `--ctx` | Context size | `2048` |
| `--threads` | Number of threads | `4` |
| `--seed` | Random seed (-1 for random) | `-1` |
| `--verbose` | Enable verbose output | `false` |

## Implementation Details

### Scheduling
The implementation uses **timestep-based scheduling** similar to the Dream architecture:
```
t = 1.0 - step/total_steps * (1.0 - eps)
s = 1.0 - (step+1)/total_steps * (1.0 - eps)  
p_transfer = (1.0 - s/t) if step < total_steps-1 else 1.0
```

### Limitations
This conceptual implementation has several limitations compared to true diffusion models:

1. **Causal Attention**: Uses standard causal attention instead of non-causal
2. **Model Architecture**: Works with standard chat models, not specialized diffusion architectures
3. **Mask Tokens**: Simulates masking rather than using actual mask tokens
4. **Simplified Sampling**: Uses basic token sampling rather than full diffusion sampling

### Educational Value
Despite limitations, this example demonstrates:
- Iterative refinement concepts
- Confidence-based token selection algorithms
- Transfer scheduling strategies
- Multi-step generation processes
- Visual progress tracking

## Real Diffusion Models

For production diffusion text generation, you would need:

### Specialized Architectures
- **Dream**: Diffusion model with epsilon-based scheduling
- **LLaDA**: Block-based diffusion with transfer scheduling
- **SUNDAE**: Semi-autoregressive diffusion models

### Key Requirements
- Non-causal attention mechanisms
- Proper mask token support
- Bidirectional context processing
- Specialized training procedures

## Building and Running

```bash
# Build the example
go build

# Run with default settings
./diffusion

# Run with custom prompt
./diffusion --prompt "Your custom prompt here"

# Interactive visual demonstration
./diffusion --diffusion-visual --prompt "Watch this generation"
```

## Related Research

- **Diffusion-LM**: Controllable text generation with diffusion models
- **DiffusionBERT**: BERT-like diffusion for text generation
- **SUNDAE**: Semi-autoregressive diffusion models
- **Dream**: Diffusion rectification and estimation models

## Contributing

This example can be enhanced with:
- More sophisticated confidence algorithms
- Better visual representations
- Performance optimizations
- Integration with specialized diffusion models
- Advanced scheduling strategies

## Technical Notes

The current implementation serves as an educational tool to understand diffusion concepts. For production use, consider:

1. Specialized diffusion model architectures
2. Non-causal attention implementations  
3. Proper mask token handling
4. Advanced sampling techniques
5. GPU-optimized implementations
