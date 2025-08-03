#!/bin/bash

# Demo script for the Gollama.cpp Speculative Decoding Example
# This script demonstrates the concepts and benefits of speculative decoding

set -e

MODEL_PATH="../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"
SPEC_BIN="./speculative"

echo "=== Gollama.cpp Speculative Decoding Example Demo ==="
echo ""

# Check if model exists
if [ ! -f "$MODEL_PATH" ]; then
    echo "❌ Model file not found: $MODEL_PATH"
    echo "Please ensure you have a GGUF model file in the models directory."
    exit 1
fi

# Build the example
echo "🔨 Building speculative decoding example..."
go build -o speculative main.go
echo "✅ Build complete!"
echo ""

# Introduction
echo "📚 What is Speculative Decoding?"
echo ""
echo "Speculative decoding is an optimization technique that can speed up text generation"
echo "by using a smaller 'draft' model to predict tokens ahead, then verifying them"
echo "with a larger 'target' model. This can provide 2-4x speedup while maintaining"
echo "the same quality as using the target model alone."
echo ""
echo "Process:"
echo "1. Draft model generates N tokens ahead"
echo "2. Target model verifies these predictions in parallel" 
echo "3. Matching predictions are accepted; mismatches trigger new sampling"
echo ""
echo "---"
echo ""

# Demo 1: Basic demonstration with verbose output
echo "🔍 Demo 1: Basic Speculative Decoding (Verbose Mode)"
echo "This shows the internal draft/verify process using the same model for both target and draft."
echo ""
echo "Command: $SPEC_BIN -prompt \"The future of artificial intelligence is\" -n-predict 60 -n-draft 5 -verbose"
echo ""
$SPEC_BIN -prompt "The future of artificial intelligence is" -n-predict 60 -n-draft 5 -verbose
echo ""
echo "---"
echo ""

# Demo 2: Different draft lengths comparison
echo "📊 Demo 2: Impact of Draft Length"
echo "Comparing different numbers of draft tokens to show the trade-off between"
echo "potential speedup and acceptance rate."
echo ""

echo "🔸 Short drafts (n-draft=3):"
echo "Command: $SPEC_BIN -prompt \"Machine learning algorithms\" -n-predict 50 -n-draft 3"
echo ""
$SPEC_BIN -prompt "Machine learning algorithms" -n-predict 50 -n-draft 3
echo ""

echo "🔸 Medium drafts (n-draft=8):"
echo "Command: $SPEC_BIN -prompt \"Machine learning algorithms\" -n-predict 50 -n-draft 8"
echo ""
$SPEC_BIN -prompt "Machine learning algorithms" -n-predict 50 -n-draft 8
echo ""

echo "🔸 Long drafts (n-draft=15):"
echo "Command: $SPEC_BIN -prompt \"Machine learning algorithms\" -n-predict 50 -n-draft 15"
echo ""
$SPEC_BIN -prompt "Machine learning algorithms" -n-predict 50 -n-draft 15
echo ""
echo "---"
echo ""

# Demo 3: Creative writing with temperature
echo "✨ Demo 3: Creative Writing with Temperature"
echo "Using temperature sampling for more creative output."
echo ""
echo "Command: $SPEC_BIN -prompt \"Once upon a time in a distant galaxy\" -n-predict 120 -n-draft 6 -temperature 0.3"
echo ""
$SPEC_BIN -prompt "Once upon a time in a distant galaxy" -n-predict 120 -n-draft 6 -temperature 0.3
echo ""
echo "---"
echo ""

# Demo 4: Technical explanation
echo "🔬 Demo 4: Technical Explanation Generation"
echo "Generating technical content with speculative decoding."
echo ""
echo "Command: $SPEC_BIN -prompt \"Quantum computing works by utilizing\" -n-predict 100 -n-draft 8 -temperature 0.1"
echo ""
$SPEC_BIN -prompt "Quantum computing works by utilizing" -n-predict 100 -n-draft 8 -temperature 0.1
echo ""
echo "---"
echo ""

# Demo 5: Performance analysis
echo "⚡ Demo 5: Performance Analysis"
echo "Longer generation to better show performance characteristics."
echo ""
echo "Command: $SPEC_BIN -prompt \"The benefits of renewable energy include\" -n-predict 200 -n-draft 10 -verbose"
echo ""
$SPEC_BIN -prompt "The benefits of renewable energy include" -n-predict 200 -n-draft 10 -verbose
echo ""
echo "---"
echo ""

# Demo 6: Conversation generation
echo "💬 Demo 6: Conversation Generation"
echo "Generating conversational responses with speculative decoding."
echo ""
echo "Command: $SPEC_BIN -prompt \"Hello! I'm an AI assistant and I can help you with\" -n-predict 80 -n-draft 6"
echo ""
$SPEC_BIN -prompt "Hello! I'm an AI assistant and I can help you with" -n-predict 80 -n-draft 6
echo ""
echo "---"
echo ""

# Performance comparison section
echo "📈 Performance Analysis Summary"
echo ""
echo "In this demonstration using the same model for both target and draft,"
echo "we can observe the algorithm behavior but won't see real speedup."
echo ""
echo "For actual acceleration, you would use two different models:"
echo "• Target model: Large, high-quality model (13B, 30B+ parameters)"
echo "• Draft model: Small, fast model (1B, 3B parameters)"
echo ""
echo "Example with real speedup:"
echo "  $SPEC_BIN -model large_model.gguf -draft-model small_model.gguf"
echo ""

# Interactive section
echo "🎮 Try It Yourself!"
echo ""
echo "Here are some commands you can try:"
echo ""
echo "Different topics:"
echo "  $SPEC_BIN -prompt \"Climate change solutions\" -n-draft 8 -n-predict 150"
echo "  $SPEC_BIN -prompt \"The history of space exploration\" -n-draft 10 -n-predict 200"
echo "  $SPEC_BIN -prompt \"How to learn programming\" -n-draft 6 -n-predict 120"
echo ""
echo "Different parameters:"
echo "  $SPEC_BIN -prompt \"Your prompt\" -n-draft 3 -temperature 0.0    # Conservative"
echo "  $SPEC_BIN -prompt \"Your prompt\" -n-draft 12 -temperature 0.5   # Aggressive"
echo "  $SPEC_BIN -prompt \"Your prompt\" -n-draft 8 -verbose            # See process"
echo ""
echo "Performance testing:"
echo "  $SPEC_BIN -prompt \"Long story prompt\" -n-predict 500 -n-draft 15"
echo ""

# Ask if user wants to try interactive mode
read -p "Would you like to try a custom prompt? (y/N): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    read -p "Enter your prompt: " user_prompt
    read -p "Number of tokens to generate (default 100): " user_tokens
    user_tokens=${user_tokens:-100}
    read -p "Number of draft tokens (default 8): " user_draft
    user_draft=${user_draft:-8}
    read -p "Enable verbose mode? (y/N): " -n 1 -r verbose_flag
    echo ""
    
    verbose_arg=""
    if [[ $verbose_flag =~ ^[Yy]$ ]]; then
        verbose_arg="-verbose"
    fi
    
    echo ""
    echo "🚀 Generating text with speculative decoding:"
    echo "Command: $SPEC_BIN -prompt \"$user_prompt\" -n-predict $user_tokens -n-draft $user_draft $verbose_arg"
    echo ""
    $SPEC_BIN -prompt "$user_prompt" -n-predict $user_tokens -n-draft $user_draft $verbose_arg
    echo ""
fi

echo "🎉 Demo complete!"
echo ""
echo "🧠 Key Takeaways:"
echo "   • Speculative decoding can significantly speed up text generation"
echo "   • Draft length affects the trade-off between speedup and acceptance rate"
echo "   • Real speedup requires using different model sizes"
echo "   • The algorithm maintains the same quality as the target model alone"
echo "   • Temperature affects both quality and acceptance rates"
echo ""
echo "💡 Next Steps:"
echo "   • Try with different model combinations for real speedup"
echo "   • Experiment with various draft lengths and temperatures"
echo "   • Benchmark performance with longer generations"
echo "   • Explore advanced speculative decoding techniques"
echo ""
echo "📖 For more information, see the README.md file"
echo "🔧 Use 'make help' to see all available Makefile targets"
