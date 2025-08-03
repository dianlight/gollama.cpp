#!/bin/bash

# Demo script for the Gollama.cpp Simple Chat Example
# This script demonstrates various features and use cases

set -e

MODEL_PATH="../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"
CHAT_BIN="./simple-chat"

echo "=== Gollama.cpp Simple Chat Example Demo ==="
echo ""

# Check if model exists
if [ ! -f "$MODEL_PATH" ]; then
    echo "❌ Model file not found: $MODEL_PATH"
    echo "Please ensure you have a GGUF model file in the models directory."
    exit 1
fi

# Build the example
echo "🔨 Building simple chat example..."
go build -o simple-chat main.go
echo "✅ Build complete!"
echo ""

# Demo 1: Basic text completion
echo "📝 Demo 1: Basic Text Completion"
echo "Command: $CHAT_BIN -prompt \"Once upon a time\" -n-predict 80"
echo ""
$CHAT_BIN -prompt "Once upon a time" -n-predict 80
echo ""
echo "---"
echo ""

# Demo 2: Technical explanation
echo "🔬 Demo 2: Technical Explanation"
echo "Command: $CHAT_BIN -prompt \"How does artificial intelligence work?\" -n-predict 100"
echo ""
$CHAT_BIN -prompt "How does artificial intelligence work?" -n-predict 100
echo ""
echo "---"
echo ""

# Demo 3: Creative writing
echo "✨ Demo 3: Creative Writing"
echo "Command: $CHAT_BIN -prompt \"In the year 2050, robots and humans\" -n-predict 120"
echo ""
$CHAT_BIN -prompt "In the year 2050, robots and humans" -n-predict 120
echo ""
echo "---"
echo ""

# Demo 4: Conversation starter
echo "💬 Demo 4: Conversation Starter"
echo "Command: $CHAT_BIN -prompt \"Hello! I'm an AI assistant. I can help you with\" -n-predict 60"
echo ""
$CHAT_BIN -prompt "Hello! I'm an AI assistant. I can help you with" -n-predict 60
echo ""
echo "---"
echo ""

# Demo 5: Code explanation
echo "💻 Demo 5: Code/Programming Context"
echo "Command: $CHAT_BIN -prompt \"Python is a programming language that\" -n-predict 80"
echo ""
$CHAT_BIN -prompt "Python is a programming language that" -n-predict 80
echo ""
echo "---"
echo ""

# Demo 6: Longer generation with higher context
echo "📚 Demo 6: Longer Text Generation"
echo "Command: $CHAT_BIN -prompt \"The benefits of renewable energy include\" -n-predict 150 -ctx 4096"
echo ""
$CHAT_BIN -prompt "The benefits of renewable energy include" -n-predict 150 -ctx 4096
echo ""
echo "---"
echo ""

# Demo 7: Performance comparison
echo "⚡ Demo 7: Performance with Different Thread Counts"
echo "Testing with 1 thread:"
echo "Command: $CHAT_BIN -prompt \"Machine learning is\" -n-predict 50 -threads 1"
echo ""
time $CHAT_BIN -prompt "Machine learning is" -n-predict 50 -threads 1 >/dev/null 2>&1
echo ""

echo "Testing with 4 threads:"
echo "Command: $CHAT_BIN -prompt \"Machine learning is\" -n-predict 50 -threads 4"
echo ""
time $CHAT_BIN -prompt "Machine learning is" -n-predict 50 -threads 4 >/dev/null 2>&1
echo ""
echo "---"
echo ""

# Interactive section
echo "🎮 Interactive Mode"
echo ""
echo "Now you can try your own prompts! Here are some suggestions:"
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
    echo "🚀 Generating text for: \"$user_prompt\""
    echo "Command: $CHAT_BIN -prompt \"$user_prompt\" -n-predict $user_tokens"
    echo ""
    $CHAT_BIN -prompt "$user_prompt" -n-predict $user_tokens
    echo ""
fi

echo "🎉 Demo complete!"
echo ""
echo "💡 Tips for better results:"
echo "   • Use clear, specific prompts"
echo "   • Adjust -n-predict based on desired response length"
echo "   • Increase -ctx for longer conversations"
echo "   • Use more -threads for faster generation"
echo "   • Try different types of prompts (creative, technical, conversational)"
echo ""
echo "📖 For more information, see the README.md file"
echo "🔧 Use 'make help' to see all available Makefile targets"
