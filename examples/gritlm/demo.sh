#!/bin/bash

# GritLM Interactive Demo Script
# This script provides an interactive demonstration of GritLM's dual-purpose capabilities

set -e

# Colors for better output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_NAME="gritlm"
MODEL_NAME="gritlm-7b_q4_1.gguf"
MODEL_PATH="../../models/${MODEL_NAME}"

# Print colored output
print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"
    
    # Check if Go is installed
    if command -v go &> /dev/null; then
        print_success "Go is installed: $(go version)"
    else
        print_error "Go is not installed. Please install Go 1.21 or later."
        exit 1
    fi
    
    # Check if we're in the right directory
    if [ ! -f "main.go" ]; then
        print_error "main.go not found. Please run this script from the gritlm example directory."
        exit 1
    fi
    
    # Check if binary exists
    if [ -f "$BINARY_NAME" ]; then
        print_success "Binary found: $BINARY_NAME"
    else
        print_warning "Binary not found. Will build it."
    fi
    
    # Check if model exists
    if [ -f "$MODEL_PATH" ]; then
        print_success "Model found: $MODEL_PATH"
        MODEL_SIZE=$(du -h "$MODEL_PATH" | cut -f1)
        print_info "Model size: $MODEL_SIZE"
    else
        print_warning "Model not found. Will download it."
    fi
    
    echo ""
}

# Build the binary
build_binary() {
    if [ ! -f "$BINARY_NAME" ]; then
        print_header "Building GritLM Example"
        echo "Building the GritLM example binary..."
        
        if make build; then
            print_success "Binary built successfully"
        else
            print_error "Failed to build binary"
            exit 1
        fi
        echo ""
    fi
}

# Download model if needed
download_model() {
    if [ ! -f "$MODEL_PATH" ]; then
        print_header "Downloading GritLM Model"
        echo "The GritLM-7B model (~4.2GB) will be downloaded."
        echo "This may take several minutes depending on your internet connection."
        echo ""
        
        read -p "Do you want to download the model now? (y/N): " -n 1 -r
        echo ""
        
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "Downloading model..."
            if make model_download; then
                print_success "Model downloaded successfully"
            else
                print_error "Failed to download model"
                exit 1
            fi
        else
            print_error "Model is required to run the demo. Exiting."
            exit 1
        fi
        echo ""
    fi
}

# Show demo menu
show_menu() {
    print_header "GritLM Demo Options"
    echo "Choose a demonstration:"
    echo ""
    echo "1) Full Demo          - Complete GritLM dual-purpose demonstration"
    echo "2) Quick Demo         - Abbreviated version for quick overview" 
    echo "3) Embedding Focus    - Focus on embedding generation capabilities"
    echo "4) Custom Documents   - Run with your own documents"
    echo "5) Model Information  - Show detailed model information"
    echo "6) Performance Test   - Run basic performance benchmark"
    echo "7) Exit"
    echo ""
}

# Run full demo
run_full_demo() {
    print_header "GritLM Full Demonstration"
    echo "This demonstration will show:"
    echo "• Document embedding generation"
    echo "• Query processing and semantic search"
    echo "• Similarity ranking and best match identification"
    echo "• Generation setup for RAG pipeline"
    echo ""
    
    read -p "Press Enter to start the demo..." 
    echo ""
    
    echo "Starting GritLM demonstration..."
    echo "================================="
    echo ""
    
    ./"$BINARY_NAME" "$MODEL_PATH"
}

# Run quick demo with timeout
run_quick_demo() {
    print_header "GritLM Quick Demo"
    echo "Running a quick demonstration (60 seconds max)..."
    echo ""
    
    timeout 60 ./"$BINARY_NAME" "$MODEL_PATH" || {
        echo ""
        print_info "Demo completed (may have been truncated for brevity)"
    }
}

# Focus on embeddings
run_embedding_focus() {
    print_header "Embedding Generation Focus"
    echo "This demo emphasizes the embedding generation capabilities:"
    echo "• High-quality vector representations"
    echo "• Semantic similarity computation"
    echo "• Document-query matching"
    echo ""
    
    read -p "Press Enter to start..." 
    echo ""
    
    # Run with focus on embedding output
    ./"$BINARY_NAME" "$MODEL_PATH" | grep -A 20 -B 5 "embedding\|similarity\|Generated.*dimensional"
}

# Custom documents demo
run_custom_demo() {
    print_header "Custom Documents Demo"
    echo "Enter your own documents for embedding and search demonstration."
    echo "Enter up to 3 documents (press Enter twice to finish each document):"
    echo ""
    
    CUSTOM_DOCS=""
    for i in {1..3}; do
        echo "Document $i:"
        read -r DOC
        if [ -n "$DOC" ]; then
            CUSTOM_DOCS="${CUSTOM_DOCS}\"$DOC\","
        fi
    done
    
    if [ -n "$CUSTOM_DOCS" ]; then
        echo ""
        echo "Enter a search query:"
        read -r QUERY
        
        echo ""
        print_info "Running demo with your custom content..."
        echo "Documents: ${CUSTOM_DOCS%,}"
        echo "Query: $QUERY"
        echo ""
        
        # Note: This would require modifying the main program to accept custom input
        # For now, we'll run the standard demo and mention the customization capability
        print_info "Custom input feature would require program modification."
        print_info "Running standard demo to show the concept..."
        ./"$BINARY_NAME" "$MODEL_PATH"
    else
        print_warning "No documents provided. Running standard demo."
        ./"$BINARY_NAME" "$MODEL_PATH"
    fi
}

# Show model information
show_model_info() {
    print_header "GritLM Model Information"
    
    if make model_info; then
        echo ""
        print_info "Model Capabilities:"
        echo "• Embedding generation with 4096 dimensions"
        echo "• Text generation with conversational AI"
        echo "• Unified architecture for RAG applications"
        echo "• Instruction-based mode switching"
        echo ""
        
        print_info "Technical Details:"
        echo "• Architecture: Based on Mistral-7B"
        echo "• Quantization: Q4_1 (4-bit) for efficiency"
        echo "• Context length: 2048 tokens (configurable)"
        echo "• Instruction format: <|user|>...<|embed|> or <|assistant|>"
    else
        print_error "Failed to get model information"
    fi
    echo ""
}

# Run performance test
run_performance_test() {
    print_header "Performance Benchmark"
    echo "Running basic performance test..."
    echo ""
    
    if make benchmark; then
        print_success "Performance test completed"
    else
        print_error "Performance test failed"
    fi
    echo ""
}

# Main menu loop
main_menu() {
    while true; do
        show_menu
        read -p "Enter your choice (1-7): " -n 1 -r
        echo ""
        echo ""
        
        case $REPLY in
            1)
                run_full_demo
                ;;
            2)
                run_quick_demo
                ;;
            3)
                run_embedding_focus
                ;;
            4)
                run_custom_demo
                ;;
            5)
                show_model_info
                ;;
            6)
                run_performance_test
                ;;
            7)
                print_info "Exiting GritLM demo. Thank you!"
                exit 0
                ;;
            *)
                print_error "Invalid option. Please choose 1-7."
                ;;
        esac
        
        echo ""
        read -p "Press Enter to return to menu..." 
        echo ""
    done
}

# Introduction
show_introduction() {
    clear
    print_header "Welcome to GritLM Demo"
    echo "GritLM (Generative Representational Instruction Tuning) is a"
    echo "unified language model that can perform both:"
    echo ""
    echo "🔍 Embedding Generation - Create vector representations for semantic search"
    echo "💬 Text Generation - Generate human-like text responses"
    echo "🔗 RAG Pipeline - Combine retrieval and generation in one model"
    echo ""
    echo "This demo will guide you through GritLM's capabilities and show"
    echo "how it can replace separate embedding and generation models with"
    echo "a single, efficient solution."
    echo ""
    
    read -p "Press Enter to continue..." 
    echo ""
}

# Cleanup function
cleanup() {
    echo ""
    print_info "Demo session ended"
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Main execution
main() {
    # Change to script directory
    cd "$SCRIPT_DIR"
    
    # Show introduction
    show_introduction
    
    # Check prerequisites
    check_prerequisites
    
    # Build binary if needed
    build_binary
    
    # Download model if needed
    download_model
    
    # Start main menu
    main_menu
}

# Run main function
main "$@"
