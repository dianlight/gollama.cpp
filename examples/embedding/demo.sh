#!/bin/bash

# Demo script for the Gollama.cpp Embedding Example
# This script demonstrates various features of the embedding example

set -e

MODEL_PATH="../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"
EMBEDDING_BIN="./embedding"

echo "=== Gollama.cpp Embedding Example Demo ==="
echo ""

# Check if model exists
if [ ! -f "$MODEL_PATH" ]; then
    echo "❌ Model file not found: $MODEL_PATH"
    echo "Please ensure you have a GGUF model file in the models directory."
    exit 1
fi

# Build the example
echo "🔨 Building embedding example..."
go build -o embedding main.go
echo "✅ Build complete!"
echo ""

# Demo 1: Single text embedding
echo "📝 Demo 1: Single text embedding"
echo "Command: $EMBEDDING_BIN -prompt \"Hello World!\" -verbose"
echo ""
$EMBEDDING_BIN -prompt "Hello World!" -verbose
echo ""

# Demo 2: Multiple texts with similarity matrix
echo "📊 Demo 2: Multiple texts with similarity matrix"
echo "Command: $EMBEDDING_BIN -prompt \"dog|cat|animal|car|vehicle\""
echo ""
$EMBEDDING_BIN -prompt "dog|cat|animal|car|vehicle"
echo ""

# Demo 3: JSON output format
echo "🔧 Demo 3: JSON output format"
echo "Command: $EMBEDDING_BIN -prompt \"Artificial Intelligence\" -output-format json"
echo ""
$EMBEDDING_BIN -prompt "Artificial Intelligence" -output-format json
echo ""

# Demo 4: Array output format
echo "📋 Demo 4: Array output format"
echo "Command: $EMBEDDING_BIN -prompt \"Machine Learning\" -output-format array"
echo ""
$EMBEDDING_BIN -prompt "Machine Learning" -output-format array
echo ""

# Demo 5: Semantic similarity test
echo "🧠 Demo 5: Semantic similarity test"
echo "Command: $EMBEDDING_BIN -prompt \"good|great|excellent|bad|terrible\""
echo "Notice how positive words cluster together and negative words cluster together:"
echo ""
$EMBEDDING_BIN -prompt "good|great|excellent|bad|terrible"
echo ""

echo "🎉 Demo complete!"
echo ""
echo "💡 Try your own examples:"
echo "   $EMBEDDING_BIN -prompt \"your text here\""
echo "   $EMBEDDING_BIN -prompt \"text1|text2|text3\" -output-format json"
echo "   $EMBEDDING_BIN -prompt \"compare these texts\" -verbose"
