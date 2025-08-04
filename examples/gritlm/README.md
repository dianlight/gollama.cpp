# GritLM Dual-Purpose Model Example

This example demonstrates the capabilities of **GritLM** (Generative Representational Instruction Tuning for Large Language Models), a unified model that can perform both embedding generation and text generation tasks using a single model architecture.

## Overview

GritLM represents a breakthrough in unified language model design, enabling:

1. **Text Embedding Generation** - High-quality vector representations for semantic search and retrieval
2. **Text Generation** - Conversational AI and content generation capabilities  
3. **Retrieval-Augmented Generation (RAG)** - End-to-end pipeline using a single model

### Key Advantages

- **Single Model Architecture**: Eliminates the need for separate embedding and generation models
- **Instruction-Based Interface**: Uses special tokens to switch between embedding and generation modes
- **High Performance**: Competitive quality on both embedding and generation benchmarks
- **Efficiency**: Reduced model management and deployment complexity

## Model Information

- **Default Model**: `cohesionet/GritLM-7B_gguf` (gritlm-7b_q4_1.gguf)
- **Architecture**: Based on Mistral-7B with specialized training
- **Embedding Dimension**: 4096
- **Quantization**: Q4_1 (4-bit quantized for efficiency)

## Features Demonstrated

### 1. Embedding Generation
- Document encoding using `<|user|>...<|embed|>` instruction format
- Query encoding for semantic search
- Cosine similarity computation for relevance ranking
- Vector normalization for optimal similarity computation

### 2. Text Generation Setup
- Generation prompt formatting using `<|user|>...<|assistant|>` instruction format
- Context preparation for response generation
- Integration with retrieval results for RAG pipeline

### 3. Semantic Search
- Document corpus indexing with embeddings
- Query-document similarity scoring
- Best match identification and ranking

## Usage

### Basic Usage

```bash
# Using default model (will be downloaded if not present)
./gritlm ../../models/gritlm-7b_q4_1.gguf

# Using alternative model
./gritlm /path/to/your/gritlm-model.gguf
```

### Expected Output

The example demonstrates:

1. **Document Processing**: Encoding 5 sample documents about AI/ML topics
2. **Query Processing**: Processing 3 sample queries about machine learning concepts  
3. **Similarity Ranking**: Finding the best matching documents for each query
4. **Generation Setup**: Preparing prompts for text generation (setup demonstration)

Sample output:
```
=== Part 1: Text Embedding Generation ===
Encoding document 1: Machine learning is a subset of artificial intelligence...
  → Generated 4096-dimensional embedding

Query 1: What is machine learning?
  Similarity scores:
    Doc 1: 0.8234 - Machine learning is a subset of artificial intelligence...
    Doc 2: 0.6543 - The transformer architecture revolutionized natural...
  ✓ Best match: Document 1 (similarity: 0.8234)
```

## Technical Implementation

### Instruction Format

GritLM uses specific instruction templates to control model behavior:

**For Embeddings:**
```
<|user|>
{text content}
<|embed|>
```

**For Generation:**
```
<|user|>
{prompt content}
<|assistant|>
```

### Key Functions

- `EncodeInstruction()` - Formats text for embedding extraction
- `GenerateInstruction()` - Formats text for generation
- `cosineSimilarity()` - Computes semantic similarity between vectors
- `normalizeVector()` - L2 normalizes embedding vectors
- `addSequenceToBatch()` - Manages token batching for processing

## Build and Run

### Prerequisites

- Go 1.21 or later
- GritLM model file (automatically downloaded via Makefile)

### Building

```bash
# Build the example
make build

# Download the default GritLM model
make model_download

# Run with downloaded model
make run

# Clean build artifacts
make clean
```

### Manual Build

```bash
go build -o gritlm main.go
```

## Model Download

The example uses the `cohesionet/GritLM-7B_gguf` model. The Makefile includes an automated download:

```bash
make model_download
```

This downloads `gritlm-7b_q4_1.gguf` (~4.2GB) from Hugging Face using the project's `hf.sh` script.

## Architecture Details

### Dual-Purpose Design

GritLM achieves dual functionality through:

1. **Specialized Training**: Joint training on embedding and generation tasks
2. **Instruction Following**: Different instruction templates trigger different behaviors
3. **Shared Representations**: Common foundation enables knowledge transfer between tasks

### Context Modes

- **Embedding Mode**: `Embeddings = 1` in context parameters
- **Generation Mode**: `Embeddings = 0` in context parameters

## Applications

### Document Search
```go
// Index documents
for _, doc := range documents {
    embedding := generateEmbedding(model, EncodeInstruction(doc))
    index.Add(doc, embedding)
}

// Search
query := "What is machine learning?"
queryEmbedding := generateEmbedding(model, EncodeInstruction(query))
results := index.Search(queryEmbedding, topK=5)
```

### RAG Pipeline
```go
// 1. Retrieve relevant documents
docs := searchIndex(query)

// 2. Generate response using retrieved context
context := buildContext(docs)
prompt := GenerateInstruction(context + "\n\nAnswer: " + query)
response := generateText(model, prompt)
```

## Performance Characteristics

- **Embedding Quality**: Competitive with specialized embedding models
- **Generation Quality**: Maintains strong text generation capabilities
- **Efficiency**: Single model reduces memory footprint and complexity
- **Latency**: Optimized for both embedding extraction and text generation

## Comparison with Traditional Approaches

| Aspect | Traditional (Separate Models) | GritLM (Unified) |
|--------|------------------------------|------------------|
| Models Required | 2+ (embedder + generator) | 1 (dual-purpose) |
| Memory Usage | High (multiple models) | Lower (single model) |
| Deployment | Complex (model coordination) | Simple (single endpoint) |
| Consistency | Variable (different training) | High (shared foundation) |
| Maintenance | Multiple model updates | Single model management |

## Research Background

GritLM is based on the research paper:
- **"GritLM: Generative Representational Instruction Tuning for Large Language Models"**
- Demonstrates that unified models can achieve competitive performance on both embedding and generation tasks
- Introduces instruction-based mode switching for dual functionality

## Limitations and Considerations

1. **Model Size**: 7B parameters require significant computational resources
2. **Quantization Trade-offs**: Q4_1 quantization reduces precision for size efficiency
3. **Context Switching**: Separate contexts needed for embedding vs generation modes
4. **Specialization**: May not match highly specialized single-purpose models in specific domains

## Future Enhancements

Potential improvements to this example:

1. **Full Generation Loop**: Complete text generation with sampling strategies
2. **Advanced RAG**: Multi-document retrieval and context management
3. **Performance Optimization**: Batch processing and caching strategies
4. **Model Variants**: Support for different GritLM model sizes and quantizations

## References

- [GritLM Paper](https://arxiv.org/abs/2402.09906)
- [GritLM Hugging Face](https://huggingface.co/cohesionet/GritLM-7B_gguf)
- [llama.cpp Documentation](https://github.com/ggerganov/llama.cpp)

## Contributing

This example is part of the gollama.cpp project. For contributions:

1. Follow the project's coding standards
2. Add appropriate tests for new functionality
3. Update documentation for API changes
4. Consider performance implications of modifications

---

**Note**: This example demonstrates the core concepts of GritLM usage. For production applications, consider additional error handling, performance optimization, and resource management strategies.
