#!/bin/bash

# Demo script for the Gollama.cpp Batched Generation Example
# This script demonstrates multiple sequence generation capabilities

set -e

MODEL_PATH="../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"
BATCHED_BIN="./batched"

echo "=== Gollama.cpp Batched Generation Example Demo ==="
echo ""

# Check if model exists
if [ ! -f "$MODEL_PATH" ]; then
    echo "❌ Model file not found: $MODEL_PATH"
    echo "Please ensure you have a GGUF model file in the models directory."
    exit 1
fi

# Build the example
echo "🔨 Building batched generation example..."
go build -o batched main.go
echo "✅ Build complete!"
echo ""

# Introduction
echo "📚 What is Batched Generation?"
echo ""
echo "Batched generation allows you to create multiple independent text sequences"
echo "from the same prompt, which is useful for:"
echo "• Exploring different creative directions"
echo "• Comparing model outputs with different parameters"
echo "• Generating content variations for A/B testing"
echo "• Improving throughput for multiple sequences"
echo ""
echo "This implementation demonstrates the concept with sequential processing"
echo "for simplicity, while the original llama.cpp uses true parallel batching."
echo ""
echo "---"
echo ""

# Demo 1: Basic batched generation
echo "🚀 Demo 1: Basic Batched Generation"
echo "Generating 4 different sequences from the same prompt..."
echo ""
echo "Command: $BATCHED_BIN -prompt \"Hello my name is\" -n-parallel 4 -n-predict 25"
echo ""
$BATCHED_BIN -prompt "Hello my name is" -n-parallel 4 -n-predict 25
echo ""
echo "---"
echo ""

# Demo 2: Creative writing with high temperature
echo "✍️  Demo 2: Creative Writing Variations"
echo "Using high temperature for creative diversity..."
echo ""
echo "Command: $BATCHED_BIN -prompt \"In a world where magic exists,\" -n-parallel 3 -n-predict 40 -temperature 0.9"
echo ""
$BATCHED_BIN -prompt "In a world where magic exists," -n-parallel 3 -n-predict 40 -temperature 0.9
echo ""
echo "---"
echo ""

# Demo 3: Technical content with low temperature
echo "💻 Demo 3: Technical Content Generation"
echo "Using low temperature for focused, consistent output..."
echo ""
echo "Command: $BATCHED_BIN -prompt \"// Function to calculate\" -n-parallel 3 -n-predict 30 -temperature 0.3"
echo ""
$BATCHED_BIN -prompt "// Function to calculate" -n-parallel 3 -n-predict 30 -temperature 0.3
echo ""
echo "---"
echo ""

# Demo 4: Parameter comparison
echo "⚙️  Demo 4: Temperature Comparison"
echo "Comparing low vs high temperature with the same prompt..."
echo ""

echo "🔸 Low Temperature (Focused):"
echo "Command: $BATCHED_BIN -prompt \"Innovation means\" -temperature 0.2 -n-parallel 3 -n-predict 20"
echo ""
$BATCHED_BIN -prompt "Innovation means" -temperature 0.2 -n-parallel 3 -n-predict 20
echo ""

echo "🔸 High Temperature (Creative):"
echo "Command: $BATCHED_BIN -prompt \"Innovation means\" -temperature 1.1 -n-parallel 3 -n-predict 20"
echo ""
$BATCHED_BIN -prompt "Innovation means" -temperature 1.1 -n-parallel 3 -n-predict 20
echo ""
echo "---"
echo ""

# Demo 5: Performance analysis
echo "📊 Demo 5: Performance Analysis"
echo "Comparing generation speed with different sequence counts..."
echo ""

echo "🔸 Single Sequence:"
echo "Command: time $BATCHED_BIN -prompt \"Performance test\" -n-parallel 1 -n-predict 20"
echo ""
time $BATCHED_BIN -prompt "Performance test" -n-parallel 1 -n-predict 20
echo ""

echo "🔸 Multiple Sequences (4):"
echo "Command: time $BATCHED_BIN -prompt \"Performance test\" -n-parallel 4 -n-predict 20"
echo ""
time $BATCHED_BIN -prompt "Performance test" -n-parallel 4 -n-predict 20
echo ""
echo "---"
echo ""

# Demo 6: Sampling parameter effects
echo "🎯 Demo 6: Sampling Parameter Effects"
echo "Comparing different top-k and top-p settings..."
echo ""

echo "🔸 Conservative Sampling (top-k=10, top-p=0.5):"
echo "Command: $BATCHED_BIN -prompt \"Artificial intelligence\" -top-k 10 -top-p 0.5 -n-parallel 2 -n-predict 25"
echo ""
$BATCHED_BIN -prompt "Artificial intelligence" -top-k 10 -top-p 0.5 -n-parallel 2 -n-predict 25
echo ""

echo "🔸 Diverse Sampling (top-k=50, top-p=0.9):"
echo "Command: $BATCHED_BIN -prompt \"Artificial intelligence\" -top-k 50 -top-p 0.9 -n-parallel 2 -n-predict 25"
echo ""
$BATCHED_BIN -prompt "Artificial intelligence" -top-k 50 -top-p 0.9 -n-parallel 2 -n-predict 25
echo ""
echo "---"
echo ""

# Demo 7: Verbose mode
echo "🔍 Demo 7: Verbose Mode (Internal Processing)"
echo "Showing detailed processing information..."
echo ""
echo "Command: $BATCHED_BIN -prompt \"Machine learning is\" -n-parallel 2 -n-predict 15 -verbose"
echo ""
$BATCHED_BIN -prompt "Machine learning is" -n-parallel 2 -n-predict 15 -verbose
echo ""
echo "---"
echo ""

# Demo 8: Different content types
echo "🌟 Demo 8: Content Type Variations"
echo "Showing batched generation across different content types..."
echo ""

echo "🔸 Storytelling:"
echo "Command: $BATCHED_BIN -prompt \"Once upon a time in a distant kingdom\" -n-parallel 2 -n-predict 30 -temperature 0.8"
echo ""
$BATCHED_BIN -prompt "Once upon a time in a distant kingdom" -n-parallel 2 -n-predict 30 -temperature 0.8
echo ""

echo "🔸 Scientific Explanation:"
echo "Command: $BATCHED_BIN -prompt \"The theory of relativity explains\" -n-parallel 2 -n-predict 25 -temperature 0.4"
echo ""
$BATCHED_BIN -prompt "The theory of relativity explains" -n-parallel 2 -n-predict 25 -temperature 0.4
echo ""

echo "🔸 Code Documentation:"
echo "Command: $BATCHED_BIN -prompt \"This API endpoint allows users to\" -n-parallel 2 -n-predict 20 -temperature 0.3"
echo ""
$BATCHED_BIN -prompt "This API endpoint allows users to" -n-parallel 2 -n-predict 20 -temperature 0.3
echo ""
echo "---"
echo ""

# Demo 9: Conversation starters
echo "💬 Demo 9: Conversation Starters"
echo "Generating multiple conversation openers..."
echo ""
echo "Command: $BATCHED_BIN -prompt \"What I find most interesting about\" -n-parallel 4 -n-predict 15 -temperature 0.7"
echo ""
$BATCHED_BIN -prompt "What I find most interesting about" -n-parallel 4 -n-predict 15 -temperature 0.7
echo ""
echo "---"
echo ""

# Demo 10: Brainstorming
echo "💡 Demo 10: Brainstorming Session"
echo "Using batched generation for idea exploration..."
echo ""
echo "Command: $BATCHED_BIN -prompt \"Creative solutions to reduce plastic waste include\" -n-parallel 5 -n-predict 20 -temperature 0.8"
echo ""
$BATCHED_BIN -prompt "Creative solutions to reduce plastic waste include" -n-parallel 5 -n-predict 20 -temperature 0.8
echo ""
echo "---"
echo ""

# Performance comparison with other methods
echo "🔄 Performance Comparison"
echo ""
echo "Batched generation provides several advantages:"
echo "• Multiple outputs from single prompt evaluation"
echo "• Variety in responses for comparison"
echo "• Efficient use of model context"
echo "• Parallel processing capabilities (in full implementations)"
echo ""

echo "Sequential vs Batched comparison:"
echo "Sequential: Generate one sequence at a time"
echo "Batched: Generate multiple sequences together"
echo ""

# Interactive section
echo "🎮 Try Interactive Mode!"
echo ""
echo "The demos above show automated examples, but you can also run"
echo "custom batched generation sessions:"
echo ""
echo "Basic Examples:"
echo "  $BATCHED_BIN -prompt \"Your prompt here\" -n-parallel 3"
echo ""
echo "Creative Writing:"
echo "  $BATCHED_BIN -prompt \"In a world where...\" -temperature 0.9 -n-parallel 4"
echo ""
echo "Technical Content:"
echo "  $BATCHED_BIN -prompt \"This algorithm...\" -temperature 0.3 -n-parallel 2"
echo ""
echo "Brainstorming:"
echo "  $BATCHED_BIN -prompt \"Solutions to...\" -temperature 0.8 -n-parallel 5"
echo ""
echo "Performance Testing:"
echo "  $BATCHED_BIN -prompt \"Test\" -n-parallel 8 -n-predict 10"
echo ""

# Ask if user wants to try interactive mode
read -p "Would you like to try a custom prompt now? (y/N): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "🚀 Interactive Batched Generation"
    echo ""
    read -p "Enter your prompt: " user_prompt
    read -p "Number of sequences (1-8, default 3): " num_sequences
    read -p "Tokens per sequence (10-100, default 25): " num_tokens
    read -p "Temperature (0.1-1.5, default 0.8): " temperature
    
    # Set defaults
    num_sequences=${num_sequences:-3}
    num_tokens=${num_tokens:-25}
    temperature=${temperature:-0.8}
    
    echo ""
    echo "Generating $num_sequences sequences with prompt: \"$user_prompt\""
    echo "Parameters: tokens=$num_tokens, temperature=$temperature"
    echo ""
    
    $BATCHED_BIN -prompt "$user_prompt" -n-parallel "$num_sequences" -n-predict "$num_tokens" -temperature "$temperature" -verbose
    echo ""
fi

echo "🎉 Demo complete!"
echo ""
echo "🧠 Key Takeaways:"
echo "   • Batched generation creates multiple sequences from one prompt"
echo "   • Temperature controls creativity vs consistency"
echo "   • Top-k/top-p parameters affect vocabulary diversity"
echo "   • Different content types benefit from different parameters"
echo "   • Verbose mode shows internal processing details"
echo ""
echo "💡 Use Cases:"
echo "   • Creative writing variations"
echo "   • A/B testing content"
echo "   • Brainstorming sessions"
echo "   • Parameter comparison"
echo "   • Quality vs diversity analysis"
echo ""
echo "🔧 Optimization Tips:"
echo "   • Adjust sequence count based on your needs"
echo "   • Use appropriate temperature for content type"
echo "   • Monitor memory usage with large sequence counts"
echo "   • Experiment with sampling parameters"
echo ""
echo "📖 For more information, see the README.md file"
echo "🛠️  Use 'make help' to see all available Makefile targets"
