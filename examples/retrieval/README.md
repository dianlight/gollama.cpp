# Gollama.cpp Retrieval Example

This example demonstrates a document retrieval system using cosine similarity and embedding-based search. It chunks documents into smaller pieces, generates embeddings for each chunk, and allows you to query for the most semantically similar content.

## What is Retrieval-Augmented Generation (RAG)?

This example implements the core retrieval component of RAG systems:

1. **Document Processing**: Split documents into manageable chunks
2. **Embedding Generation**: Create vector representations of text chunks
3. **Similarity Search**: Find chunks most similar to user queries using cosine similarity
4. **Ranked Results**: Return top-K most relevant chunks with similarity scores

This forms the foundation for more advanced systems that combine retrieval with generation.

## Quick Start

```bash
cd examples/retrieval
go run main.go -context-files "sample_ai.txt,sample_programming.txt"
```

Then enter queries like:
- "machine learning"
- "web development"
- "neural networks"

## Features

- **Multi-file document processing** with automatic chunking
- **Configurable chunk sizes** for optimal retrieval performance
- **Interactive query interface** with real-time search
- **Batch processing mode** for automated queries
- **Cosine similarity ranking** with normalized embeddings
- **Verbose mode** for debugging and understanding the process
- **Multiple output formats** with detailed metadata

## Command Line Options

- `-model string`: Path to the GGUF model file that supports embeddings (default: "../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf")
- `-context-files string`: Comma-separated list of files to embed for retrieval
- `-chunk-size int`: Minimum size of each text chunk to be embedded (default: 200)
- `-chunk-separator string`: String to divide chunks by (default: "\n")
- `-top-k int`: Number of top similar chunks to return (default: 3)
- `-threads int`: Number of threads to use (default: 4)
- `-ctx int`: Context size (default: 2048)
- `-verbose`: Enable verbose output showing internal process
- `-interactive`: Enable interactive query mode (default: true)
- `-query string`: Single query to process in non-interactive mode

## Examples

### Basic Interactive Retrieval

```bash
# Single file
go run main.go -context-files document.txt

# Multiple files
go run main.go -context-files "file1.txt,file2.txt,file3.txt"
```

### Single Query Mode

```bash
# Non-interactive single query
go run main.go -context-files document.txt -query "artificial intelligence" -interactive=false

# With custom parameters
go run main.go -context-files data.txt -query "machine learning" -top-k 5 -chunk-size 150 -interactive=false
```

### Advanced Configuration

```bash
# Custom chunk size and separator
go run main.go -context-files document.txt -chunk-size 300 -chunk-separator "."

# Verbose mode to see internal process
go run main.go -context-files document.txt -query "search term" -verbose -interactive=false

# High precision retrieval
go run main.go -context-files document.txt -top-k 10 -chunk-size 100
```

## Using the Makefile

The included Makefile provides convenient demonstrations:

```bash
# Build the example
make build

# Interactive demo with AI concepts
make ai-demo

# Interactive demo with programming languages
make programming-demo

# Combined demo with both sample files
make combined-demo

# Single query demonstrations
make single-query

# Verbose demo showing internal process
make verbose-demo

# Compare different chunk sizes
make chunk-size-demo

# Full demonstration
make demo
```

## Sample Files

The example includes two sample files for demonstration:

### `sample_ai.txt`
- Artificial Intelligence concepts
- Machine Learning definitions
- Deep Learning explanations
- Neural Networks information
- Computer Vision overview

### `sample_programming.txt`
- Programming language descriptions
- Language characteristics
- Development use cases
- Platform information

Try queries like:
- **AI file**: "machine learning", "neural networks", "computer vision"
- **Programming file**: "web development", "systems programming", "mobile apps"
- **Combined**: "intelligent software", "AI programming", "data science"

## Understanding the Output

### Query Results Format

```
Top 3 similar chunks:
filename: sample_ai.txt
filepos: 156
similarity: 0.892341
textdata:
Machine Learning is a subset of AI that provides systems the ability to automatically learn and improve from experience without being explicitly programmed.
--------------------
filename: sample_ai.txt
filepos: 312
similarity: 0.785623
textdata:
Deep Learning is a subset of machine learning that uses artificial neural networks with multiple layers to model and understand complex patterns in data.
--------------------
```

### Output Fields
- **filename**: Source file containing the chunk
- **filepos**: Character position in the original file
- **similarity**: Cosine similarity score (0.0 to 1.0, higher is more similar)
- **textdata**: The actual text content of the chunk

## Performance Tuning

### Chunk Size (`-chunk-size`)
- **Small (50-100)**: More granular matching, more chunks to search
- **Medium (150-300)**: Balanced context and specificity
- **Large (400+)**: More context but less specific matching

### Chunk Separator (`-chunk-separator`)
- **"\n"** (default): Split by lines, good for structured text
- **"."**: Split by sentences, good for narrative text
- **"\n\n"**: Split by paragraphs, good for document sections
- **Custom**: Use domain-specific separators

### Top-K Results (`-top-k`)
- **1-3**: Most relevant results only
- **5-10**: Broader range of relevant content
- **10+**: Comprehensive results for analysis

## Model Requirements

### Embedding Models
This example requires models that support embedding generation:
- **Sentence-BERT models**: BGE, E5, instructor models
- **General embedding models**: Models trained for semantic similarity
- **Domain-specific models**: Models fine-tuned for your content domain

### Model Compatibility
- Must support the embedding flag in context parameters
- Should have pooling enabled (not NONE pooling type)
- Vocabulary should be appropriate for your text domain

### Recommended Models
- **BGE models**: `bge-base-en-v1.5`, `bge-large-en-v1.5`
- **E5 models**: `e5-base-v2`, `e5-large-v2`
- **Instructor models**: `instructor-base`, `instructor-large`

## Use Cases

### Document Search
```bash
# Search through documentation
./retrieval -context-files "docs/api.md,docs/tutorial.md" -query "authentication"
```

### Knowledge Base Query
```bash
# Query a knowledge base
./retrieval -context-files "knowledge_base.txt" -top-k 5 -chunk-size 200
```

### Research Paper Analysis
```bash
# Analyze research papers
./retrieval -context-files "paper1.txt,paper2.txt" -query "methodology" -chunk-separator "."
```

### Code Documentation Search
```bash
# Search code documentation
./retrieval -context-files "readme.md,api_docs.txt" -query "installation" -top-k 3
```

## Troubleshooting

### "Failed to load model"
- Ensure your model supports embeddings
- Check that the model file path is correct
- Verify the model is compatible with gollama.cpp

### "No chunks were created"
- Check that input files exist and are readable
- Verify chunk size isn't larger than your document content
- Ensure chunk separator exists in your documents

### Low Similarity Scores
- Try different chunk sizes for better granularity
- Use models trained on similar domains
- Consider preprocessing text (cleaning, formatting)

### Poor Retrieval Quality
- Experiment with different chunk separators
- Adjust chunk size for your content type
- Try domain-specific embedding models
- Increase top-k to see more results

### Memory Issues
- Reduce context size (`-ctx`)
- Process fewer files at once
- Use smaller embedding models
- Reduce chunk sizes

## Technical Details

### Algorithm Overview

1. **Document Chunking**: Split input files into overlapping or sequential chunks
2. **Tokenization**: Convert text chunks to model tokens
3. **Embedding Generation**: Generate vector representations for each chunk
4. **Normalization**: L2-normalize embeddings for cosine similarity
5. **Query Processing**: Generate embedding for user query
6. **Similarity Calculation**: Compute cosine similarity between query and all chunks
7. **Ranking**: Sort results by similarity score and return top-K

### Implementation Features

- **Batch Processing**: Efficient embedding generation for multiple chunks
- **Memory Management**: Automatic cleanup of intermediate data
- **Error Handling**: Robust handling of file and model errors
- **Streaming Interface**: Real-time query processing

## Advanced Usage

### Custom Text Processing

```bash
# Scientific papers with section-based chunking
./retrieval -context-files "papers/*.txt" -chunk-separator "## " -chunk-size 500

# News articles with paragraph chunking
./retrieval -context-files "news/*.txt" -chunk-separator "\n\n" -chunk-size 250

# Code documentation with function-based chunking
./retrieval -context-files "docs/*.md" -chunk-separator "###" -chunk-size 300
```

### Performance Benchmarking

```bash
# Time different configurations
time ./retrieval -context-files large_document.txt -query "test" -interactive=false

# Compare chunk sizes
for size in 100 200 400; do
  echo "Testing chunk size: $size"
  time ./retrieval -context-files document.txt -chunk-size $size -query "test" -interactive=false
done
```

### Integration with Other Systems

This retrieval system can be integrated with:
- **Text generation models** for RAG systems
- **Question answering systems** for context retrieval
- **Chatbots** for knowledge base queries
- **Search engines** for semantic search
- **Content management systems** for document discovery

## Limitations

- **Model Dependency**: Requires embedding-capable models
- **Memory Usage**: Stores all embeddings in memory
- **Single Query**: Processes one query at a time
- **Basic Chunking**: Simple separator-based chunking strategy
- **No Persistence**: Embeddings are recalculated each run

## Next Steps

After trying this example:

1. **Experiment with Models**: Try different embedding models for your domain
2. **Optimize Chunking**: Develop domain-specific chunking strategies
3. **Add Generation**: Combine with text generation for full RAG
4. **Improve Ranking**: Implement advanced similarity measures
5. **Add Persistence**: Save embeddings to disk for faster startup
6. **Scale Up**: Process larger document collections

## Requirements

- Go 1.21 or later
- A GGUF model that supports embeddings
- Text files to search through
- Sufficient RAM for model and embeddings (varies by content size)
