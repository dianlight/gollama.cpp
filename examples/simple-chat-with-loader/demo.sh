#!/bin/bash

# Demo script for the Gollama.cpp Simple Chat with Library Loader Example
# This script demonstrates both library loading and text generation features

set -e

MODEL_PATH="../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"
CHAT_BIN="./simple-chat-with-loader"

echo "=== Gollama.cpp Simple Chat with Library Loader Example Demo ==="
echo ""

# Check if model exists
if [ ! -f "$MODEL_PATH" ]; then
    echo "❌ Model file not found: $MODEL_PATH"
    echo "Please ensure you have a GGUF model file in the models directory."
    exit 1
fi

# Build the example
echo "🔨 Building simple chat with library loader example..."
go build -o simple-chat-with-loader main.go
echo "✅ Build complete!"
echo ""

echo "🔄 This demo showcases:"
echo "   • Dynamic library loading and management"
echo "   • Embedded library extraction" 
echo "   • Text generation using loaded libraries"
echo "   • Proper resource cleanup"
echo ""

# Demo 1: Basic functionality with library loading
echo "📚 Demo 1: Library Loading + Basic Text Completion"
echo "Command: $CHAT_BIN -prompt \"Once upon a time\" -n-predict 80"
echo ""
$CHAT_BIN -prompt "Once upon a time" -n-predict 80
echo ""
echo "---"
echo ""

# Demo 2: Technical explanation with library management
echo "🔬 Demo 2: Library Management + Technical Explanation"
echo "Command: $CHAT_BIN -prompt \"How does artificial intelligence work?\" -n-predict 100"
echo ""
$CHAT_BIN -prompt "How does artificial intelligence work?" -n-predict 100
echo ""
echo "---"
echo ""

# Demo 3: Creative writing demonstrating full workflow
echo "✨ Demo 3: Full Workflow + Creative Writing"
echo "Command: $CHAT_BIN -prompt \"In the year 2050, robots and humans\" -n-predict 120"
echo ""
$CHAT_BIN -prompt "In the year 2050, robots and humans" -n-predict 120
echo ""
echo "---"
echo ""

# Demo 4: Testing library loading with minimal output
echo "⚙️ Demo 4: Library Loading Focus (Minimal Generation)"
echo "Command: $CHAT_BIN -prompt \"Testing library loader\" -n-predict 20"
echo ""
$CHAT_BIN -prompt "Testing library loader" -n-predict 20
echo ""
echo "---"
echo ""

# Demo 5: Conversation starter with resource management
echo "💬 Demo 5: Resource Management + Conversation"
echo "Command: $CHAT_BIN -prompt \"Hello! I'm an AI assistant. I can help you with\" -n-predict 60"
echo ""
$CHAT_BIN -prompt "Hello! I'm an AI assistant. I can help you with" -n-predict 60
echo ""
echo "---"
echo ""

# Demo 6: Longer generation testing library stability
echo "📚 Demo 6: Library Stability + Longer Generation"
echo "Command: $CHAT_BIN -prompt \"The benefits of renewable energy include\" -n-predict 150 -ctx 4096"
echo ""
$CHAT_BIN -prompt "The benefits of renewable energy include" -n-predict 150 -ctx 4096
echo ""
echo "---"
echo ""

# Demo 7: Performance testing with library overhead
echo "⚡ Demo 7: Performance Testing (Library Loading Overhead)"
echo "Testing generation time including library loading:"
echo "Command: $CHAT_BIN -prompt \"Machine learning is\" -n-predict 50 -threads 4"
echo ""
echo "Timing full execution (including library load/unload):"
time $CHAT_BIN -prompt "Machine learning is" -n-predict 50 -threads 4 >/dev/null 2>&1
echo ""
echo "---"
echo ""

# Interactive section
echo "🎮 Interactive Mode with Library Loader"
echo ""
echo "Try your own prompts with dynamic library loading! Suggestions:"
echo ""
echo "Library testing prompts:"
echo "  $CHAT_BIN -prompt \"Testing embedded libraries\" -n-predict 50"
echo "  $CHAT_BIN -prompt \"Dynamic loading works\" -n-predict 30"
echo ""
echo "Story starters:"
echo "  $CHAT_BIN -prompt \"The last person on Earth\" -n-predict 100"
echo "  $CHAT_BIN -prompt \"A mysterious letter arrived\" -n-predict 150"
echo ""
echo "Educational prompts:"
echo "  $CHAT_BIN -prompt \"The solar system consists of\" -n-predict 120"
echo "  $CHAT_BIN -prompt \"Climate change is caused by\" -n-predict 100"
echo ""
echo "Creative prompts:"
echo "  $CHAT_BIN -prompt \"If I could travel anywhere\" -n-predict 80"
echo "  $CHAT_BIN -prompt \"The recipe for happiness\" -n-predict 100"
echo ""
echo "Technical prompts:"
echo "  $CHAT_BIN -prompt \"To build a website, you need\" -n-predict 120"
echo "  $CHAT_BIN -prompt \"The difference between AI and ML\" -n-predict 100"
echo ""

# Ask if user wants to try interactive mode
read -p "Would you like to try an interactive prompt? (y/N): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    read -p "Enter your prompt: " user_prompt
    read -p "Number of tokens to generate (default 100): " user_tokens
    user_tokens=${user_tokens:-100}
    
    echo ""
    echo "🚀 Generating text with library loader for: \"$user_prompt\""
    echo "Command: $CHAT_BIN -prompt \"$user_prompt\" -n-predict $user_tokens"
    echo ""
    $CHAT_BIN -prompt "$user_prompt" -n-predict $user_tokens
    echo ""
fi

echo "🎉 Library Loader Demo complete!"
echo ""
echo "🔧 What this demo showed:"
echo "   • Automatic library extraction from embedded files"
echo "   • Dynamic loading of platform-specific libraries"
echo "   • Seamless integration with text generation"
echo "   • Proper cleanup of library handles and temporary files"
echo "   • Cross-platform compatibility"
echo ""
echo "💡 Tips for better results:"
echo "   • Library loading adds minimal overhead"
echo "   • Embedded libraries ensure portability"
echo "   • Use clear, specific prompts for better generation"
echo "   • Adjust -n-predict based on desired response length"
echo "   • Increase -ctx for longer conversations"
echo "   • Use more -threads for faster generation"
echo ""
echo "📖 For more information, see the README.md file"
echo "🔧 Use 'make help' to see all available Makefile targets"
echo "🎯 Compare with regular simple-chat to see the loader benefits"
