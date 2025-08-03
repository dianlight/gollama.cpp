#!/bin/bash

# Demo script for the Gollama.cpp Retrieval Example
# This script demonstrates document retrieval and semantic search capabilities

set -e

MODEL_PATH="../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"
RETRIEVAL_BIN="./retrieval"

echo "=== Gollama.cpp Retrieval Example Demo ==="
echo ""

# Check if model exists
if [ ! -f "$MODEL_PATH" ]; then
    echo "❌ Model file not found: $MODEL_PATH"
    echo "Please ensure you have a GGUF model file that supports embeddings."
    exit 1
fi

# Check if sample files exist
if [ ! -f "sample_ai.txt" ] || [ ! -f "sample_programming.txt" ]; then
    echo "❌ Sample files not found. They should be created automatically."
    echo "Please run 'make sample-files' or check the repository."
    exit 1
fi

# Build the example
echo "🔨 Building retrieval example..."
go build -o retrieval main.go
echo "✅ Build complete!"
echo ""

# Introduction
echo "📚 What is Document Retrieval?"
echo ""
echo "Document retrieval is a core component of modern AI systems that allows you to:"
echo "• Search through large document collections using natural language"
echo "• Find semantically similar content (not just keyword matching)"
echo "• Rank results by relevance using embedding similarity"
echo "• Build the foundation for Retrieval-Augmented Generation (RAG)"
echo ""
echo "Process:"
echo "1. Split documents into chunks"
echo "2. Generate embeddings for each chunk"
echo "3. Generate embedding for user query"
echo "4. Compute similarity scores and rank results"
echo ""
echo "---"
echo ""

# Demo 1: Basic retrieval with AI concepts
echo "🔍 Demo 1: AI Concepts Retrieval"
echo "Searching through AI and machine learning concepts..."
echo ""
echo "Sample queries to try:"
echo "• machine learning"
echo "• neural networks" 
echo "• computer vision"
echo ""
echo "Command: $RETRIEVAL_BIN -context-files sample_ai.txt -query \"machine learning\" -interactive=false -top-k 3"
echo ""
$RETRIEVAL_BIN -context-files sample_ai.txt -query "machine learning" -interactive=false -top-k 3
echo ""
echo "---"
echo ""

# Demo 2: Programming languages retrieval
echo "💻 Demo 2: Programming Languages Retrieval"
echo "Searching through programming language descriptions..."
echo ""
echo "Command: $RETRIEVAL_BIN -context-files sample_programming.txt -query \"web development\" -interactive=false -top-k 3"
echo ""
$RETRIEVAL_BIN -context-files sample_programming.txt -query "web development" -interactive=false -top-k 3
echo ""
echo "---"
echo ""

# Demo 3: Cross-domain search
echo "🌐 Demo 3: Cross-Domain Search"
echo "Searching across both AI and programming files..."
echo ""
echo "Command: $RETRIEVAL_BIN -context-files \"sample_ai.txt,sample_programming.txt\" -query \"intelligent software\" -interactive=false -top-k 4"
echo ""
$RETRIEVAL_BIN -context-files "sample_ai.txt,sample_programming.txt" -query "intelligent software" -interactive=false -top-k 4
echo ""
echo "---"
echo ""

# Demo 4: Different chunk sizes
echo "📏 Demo 4: Impact of Chunk Size"
echo "Comparing different chunk sizes for the same query..."
echo ""

echo "🔸 Small chunks (chunk-size=80):"
echo "Command: $RETRIEVAL_BIN -context-files sample_ai.txt -query \"neural networks\" -chunk-size 80 -interactive=false -top-k 2"
echo ""
$RETRIEVAL_BIN -context-files sample_ai.txt -query "neural networks" -chunk-size 80 -interactive=false -top-k 2
echo ""

echo "🔸 Large chunks (chunk-size=250):"
echo "Command: $RETRIEVAL_BIN -context-files sample_ai.txt -query \"neural networks\" -chunk-size 250 -interactive=false -top-k 2"
echo ""
$RETRIEVAL_BIN -context-files sample_ai.txt -query "neural networks" -chunk-size 250 -interactive=false -top-k 2
echo ""
echo "---"
echo ""

# Demo 5: Verbose mode
echo "🔍 Demo 5: Verbose Mode (Internal Process)"
echo "Showing the internal processing steps..."
echo ""
echo "Command: $RETRIEVAL_BIN -context-files sample_ai.txt -query \"deep learning\" -interactive=false -verbose -top-k 2"
echo ""
$RETRIEVAL_BIN -context-files sample_ai.txt -query "deep learning" -interactive=false -verbose -top-k 2
echo ""
echo "---"
echo ""

# Demo 6: Different separators
echo "📄 Demo 6: Different Chunk Separators"
echo "Using sentence-based chunking instead of line-based..."
echo ""
echo "Command: $RETRIEVAL_BIN -context-files sample_programming.txt -query \"systems programming\" -chunk-separator \".\" -interactive=false -top-k 2"
echo ""
$RETRIEVAL_BIN -context-files sample_programming.txt -query "systems programming" -chunk-separator "." -interactive=false -top-k 2
echo ""
echo "---"
echo ""

# Demo 7: High-precision search
echo "🎯 Demo 7: High-Precision Search"
echo "Using more results to find comprehensive matches..."
echo ""
echo "Command: $RETRIEVAL_BIN -context-files \"sample_ai.txt,sample_programming.txt\" -query \"algorithms\" -top-k 6 -interactive=false"
echo ""
$RETRIEVAL_BIN -context-files "sample_ai.txt,sample_programming.txt" -query "algorithms" -top-k 6 -interactive=false
echo ""
echo "---"
echo ""

# Demo 8: Similarity score analysis
echo "📊 Demo 8: Understanding Similarity Scores"
echo "Comparing different queries to show similarity ranges..."
echo ""

echo "🔸 Exact match query:"
echo "Command: $RETRIEVAL_BIN -context-files sample_ai.txt -query \"Artificial Intelligence\" -interactive=false -top-k 1"
echo ""
$RETRIEVAL_BIN -context-files sample_ai.txt -query "Artificial Intelligence" -interactive=false -top-k 1
echo ""

echo "🔸 Related concept query:"
echo "Command: $RETRIEVAL_BIN -context-files sample_ai.txt -query \"AI technology\" -interactive=false -top-k 1"
echo ""
$RETRIEVAL_BIN -context-files sample_ai.txt -query "AI technology" -interactive=false -top-k 1
echo ""

echo "🔸 Distant concept query:"
echo "Command: $RETRIEVAL_BIN -context-files sample_ai.txt -query \"cooking recipes\" -interactive=false -top-k 1"
echo ""
$RETRIEVAL_BIN -context-files sample_ai.txt -query "cooking recipes" -interactive=false -top-k 1
echo ""
echo "---"
echo ""

# Performance comparison
echo "⚡ Performance Analysis"
echo ""
echo "The retrieval system processes documents in these steps:"
echo "1. Document chunking and tokenization"
echo "2. Embedding generation for all chunks"
echo "3. Query embedding generation"
echo "4. Similarity calculation and ranking"
echo ""
echo "Performance factors:"
echo "• Number of chunks (affects search time)"
echo "• Chunk size (affects context quality)"
echo "• Model size (affects embedding quality and speed)"
echo "• Top-K value (affects result comprehensiveness)"
echo ""

# Interactive section
echo "🎮 Try Interactive Mode!"
echo ""
echo "The demos above show automated queries, but the real power comes from"
echo "interactive exploration. Here are some commands to try:"
echo ""
echo "AI and Machine Learning queries:"
echo "  make ai-demo"
echo "  # Then try: 'supervised learning', 'computer vision', 'data science'"
echo ""
echo "Programming and Technology queries:"
echo "  make programming-demo"
echo "  # Then try: 'mobile development', 'web frameworks', 'functional programming'"
echo ""
echo "Combined domain queries:"
echo "  make combined-demo"
echo "  # Then try: 'AI programming', 'intelligent systems', 'automation'"
echo ""
echo "Custom file queries:"
echo "  $RETRIEVAL_BIN -context-files \"your_file.txt\""
echo "  # Then enter any query related to your content"
echo ""

# Ask if user wants to try interactive mode
read -p "Would you like to try interactive mode now? (y/N): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "🚀 Starting Interactive Retrieval Session"
    echo ""
    echo "Choose a demo:"
    echo "1. AI concepts (sample_ai.txt)"
    echo "2. Programming languages (sample_programming.txt)"
    echo "3. Both files combined"
    echo ""
    read -p "Enter your choice (1-3): " -n 1 -r choice
    echo ""
    
    case $choice in
        1)
            echo "Loading AI concepts file..."
            $RETRIEVAL_BIN -context-files sample_ai.txt -top-k 3
            ;;
        2)
            echo "Loading programming languages file..."
            $RETRIEVAL_BIN -context-files sample_programming.txt -top-k 3
            ;;
        3)
            echo "Loading both files..."
            $RETRIEVAL_BIN -context-files "sample_ai.txt,sample_programming.txt" -top-k 4
            ;;
        *)
            echo "Loading AI concepts file (default)..."
            $RETRIEVAL_BIN -context-files sample_ai.txt -top-k 3
            ;;
    esac
    echo ""
fi

echo "🎉 Demo complete!"
echo ""
echo "🧠 Key Takeaways:"
echo "   • Retrieval systems enable semantic search beyond keyword matching"
echo "   • Chunk size affects the granularity and context of results"
echo "   • Similarity scores help rank relevance of retrieved content"
echo "   • Cross-domain search can find unexpected connections"
echo "   • Interactive mode is powerful for exploratory research"
echo ""
echo "💡 Advanced Use Cases:"
echo "   • Document Q&A systems"
echo "   • Knowledge base search"
echo "   • Research paper analysis"
echo "   • Code documentation search"
echo "   • Customer support knowledge retrieval"
echo "   • Content recommendation systems"
echo ""
echo "🔧 Optimization Tips:"
echo "   • Use domain-specific embedding models for better results"
echo "   • Experiment with chunk sizes for your content type"
echo "   • Consider preprocessing text for better chunking"
echo "   • Combine retrieval with generation for full RAG systems"
echo ""
echo "📖 For more information, see the README.md file"
echo "🛠️  Use 'make help' to see all available Makefile targets"
