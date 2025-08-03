#!/bin/bash

# Gollama.cpp Evaluation Callback Example Demo Script
# This script demonstrates the evaluation callback functionality

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Utility functions
print_header() {
    echo -e "\n${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}\n"
}

print_section() {
    echo -e "\n${CYAN}--- $1 ---${NC}\n"
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

wait_for_user() {
    if [[ "${AUTO_MODE:-}" != "true" ]]; then
        echo -e "\n${YELLOW}Press Enter to continue...${NC}"
        read -r
    else
        sleep 2
    fi
}

# Build the example
build_example() {
    print_section "Building Evaluation Callback Example"
    
    if make build; then
        print_success "Build completed successfully"
    else
        print_error "Build failed"
        exit 1
    fi
}

# Demo 1: Basic Simulation
demo_basic_simulation() {
    print_header "Demo 1: Basic Evaluation Callback Simulation"
    
    echo "This demonstrates the basic evaluation callback simulation."
    echo "It shows how tensor operations would be logged during inference."
    wait_for_user
    
    print_section "Running Basic Simulation"
    ./eval-callback -simulate-only -prompt "Hello world" -max-logged-ops 15
    
    print_success "Basic simulation completed"
    wait_for_user
}

# Demo 2: Operation Logging
demo_operation_logging() {
    print_header "Demo 2: Detailed Operation Logging"
    
    echo "This shows detailed tensor operation information including:"
    echo "• Operation types (MUL_MAT, SOFT_MAX, etc.)"
    echo "• Tensor names and dimensions"
    echo "• Memory locations (CPU vs GPU)"
    echo "• Timing information"
    wait_for_user
    
    print_section "Detailed Operation Logging"
    ./eval-callback -simulate-only -prompt "Neural networks are powerful" -max-logged-ops 25
    
    print_success "Operation logging demo completed"
    wait_for_user
}

# Demo 3: Performance Monitoring
demo_performance_monitoring() {
    print_header "Demo 3: Performance Monitoring"
    
    echo "This demonstrates performance monitoring capabilities:"
    echo "• Operation counting"
    echo "• Data throughput measurement"
    echo "• Timing analysis"
    echo "• Progress reporting"
    wait_for_user
    
    print_section "Performance Monitoring with Progress Updates"
    ./eval-callback -simulate-only -prompt "The future of artificial intelligence technology" -enable-progress
    
    print_success "Performance monitoring demo completed"
    wait_for_user
}

# Demo 4: Memory Location Tracking
demo_memory_tracking() {
    print_header "Demo 4: Memory Location Tracking"
    
    echo "This shows how eval callbacks can track memory usage:"
    echo "• CPU vs GPU memory placement"
    echo "• Memory transfer operations"
    echo "• Tensor size monitoring"
    wait_for_user
    
    print_section "Memory Location Tracking"
    ./eval-callback -simulate-only -prompt "Memory management in machine learning" -max-logged-ops 20
    
    print_success "Memory tracking demo completed"
    wait_for_user
}

# Demo 5: Tensor Data Inspection
demo_tensor_data() {
    print_header "Demo 5: Tensor Data Inspection"
    
    echo "This demonstrates tensor data value inspection:"
    echo "• Viewing actual tensor values"
    echo "• Data type information"
    echo "• Tensor shape analysis"
    echo ""
    print_warning "Note: This produces verbose output"
    wait_for_user
    
    print_section "Tensor Data Inspection (Limited Output)"
    ./eval-callback -simulate-only -prompt "Data inspection" -print-tensor-data -max-logged-ops 8
    
    print_success "Tensor data inspection demo completed"
    wait_for_user
}

# Demo 6: Configuration Comparison
demo_configuration_comparison() {
    print_header "Demo 6: Configuration Comparison"
    
    echo "This compares different callback configurations:"
    echo "• Minimal logging vs full logging"
    echo "• Progress-only vs detailed operations"
    echo "• Different logging limits"
    wait_for_user
    
    print_section "Minimal Logging"
    ./eval-callback -simulate-only -prompt "Configuration test" -enable-logging=false -enable-progress
    
    print_section "Full Operation Logging"
    ./eval-callback -simulate-only -prompt "Configuration test" -max-logged-ops 12
    
    print_section "Tensor Data with Limited Operations"
    ./eval-callback -simulate-only -prompt "Configuration test" -print-tensor-data -max-logged-ops 5
    
    print_success "Configuration comparison completed"
    wait_for_user
}

# Demo 7: Threading and Context Tests
demo_performance_tests() {
    print_header "Demo 7: Performance Configuration Tests"
    
    echo "This tests different performance configurations:"
    echo "• Various thread counts"
    echo "• Different context sizes"
    echo "• Performance impact analysis"
    wait_for_user
    
    print_section "Testing Different Thread Counts"
    for threads in 1 2 4 8; do
        echo "Testing with $threads threads..."
        ./eval-callback -simulate-only -threads $threads -prompt "Threading test" -enable-logging=false -enable-progress
    done
    
    print_section "Testing Different Context Sizes"
    for ctx in 256 512 1024; do
        echo "Testing with context size $ctx..."
        ./eval-callback -simulate-only -ctx $ctx -prompt "Context test" -enable-logging=false -enable-progress
    done
    
    print_success "Performance tests completed"
    wait_for_user
}

# Demo 8: Real Model Evaluation (if available)
demo_real_model() {
    print_header "Demo 8: Real Model Evaluation"
    
    local model_path="../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"
    
    if [[ -f "$model_path" ]]; then
        echo "Model file found! This will run actual model inference"
        echo "alongside the callback simulation."
        print_warning "This requires more memory and time"
        wait_for_user
        
        print_section "Real Model Evaluation with Callbacks"
        ./eval-callback -model "$model_path" -prompt "Hello world" -ctx 512
        
        print_success "Real model evaluation completed"
    else
        print_warning "Model file not found: $model_path"
        echo "Skipping real model evaluation."
        echo "You can download a model and place it at the expected location"
        echo "to test with real model inference."
    fi
    wait_for_user
}

# Demo 9: Advanced Debugging Scenarios
demo_advanced_debugging() {
    print_header "Demo 9: Advanced Debugging Scenarios"
    
    echo "This demonstrates advanced debugging use cases:"
    echo "• Complex prompt processing"
    echo "• Long sequence handling"
    echo "• Performance bottleneck identification"
    wait_for_user
    
    print_section "Complex Prompt Processing"
    ./eval-callback -simulate-only -prompt "Explain how transformer attention mechanisms work in detail" -max-logged-ops 30
    
    print_section "Performance Bottleneck Analysis"
    ./eval-callback -simulate-only -prompt "Analyze performance patterns" -enable-progress -enable-logging=false
    
    print_success "Advanced debugging demo completed"
    wait_for_user
}

# Demo 10: Interactive Mode
demo_interactive() {
    print_header "Demo 10: Interactive Evaluation Callback"
    
    echo "This provides an interactive demonstration where you can"
    echo "enter custom prompts and see the callback simulation."
    echo ""
    print_warning "Type 'quit' or press Ctrl+C to exit"
    wait_for_user
    
    print_section "Interactive Mode"
    
    while true; do
        echo -e "\n${CYAN}Enter a prompt (or 'quit' to exit):${NC}"
        read -r prompt
        
        if [[ "$prompt" == "quit" ]] || [[ "$prompt" == "exit" ]]; then
            break
        fi
        
        if [[ -n "$prompt" ]]; then
            echo -e "\n${YELLOW}Processing: $prompt${NC}"
            ./eval-callback -simulate-only -prompt "$prompt" -max-logged-ops 15
            print_success "Prompt processed"
        fi
    done
    
    print_success "Interactive demo completed"
}

# Main demo function
run_demo() {
    print_header "Gollama.cpp Evaluation Callback Example Demo"
    
    echo "This script demonstrates the evaluation callback functionality"
    echo "that allows monitoring tensor operations during model inference."
    echo ""
    echo "Available demos:"
    echo "  1. Basic Simulation"
    echo "  2. Operation Logging"
    echo "  3. Performance Monitoring"
    echo "  4. Memory Tracking"
    echo "  5. Tensor Data Inspection"
    echo "  6. Configuration Comparison"
    echo "  7. Performance Tests"
    echo "  8. Real Model Evaluation"
    echo "  9. Advanced Debugging"
    echo "  10. Interactive Mode"
    echo ""
    
    if [[ "${AUTO_MODE:-}" == "true" ]]; then
        echo "Running in automatic mode..."
        echo "Set AUTO_MODE=false to run interactively"
        wait_for_user
    else
        echo "Choose demo mode:"
        echo "  a) Run all demos automatically"
        echo "  i) Interactive demo selection"
        echo "  q) Quit"
        echo ""
        read -p "Your choice [a/i/q]: " choice
        
        case $choice in
            q|Q) exit 0 ;;
            i|I) 
                interactive_menu
                return
                ;;
            *) 
                echo "Running all demos automatically..."
                ;;
        esac
    fi
    
    # Build first
    build_example
    
    # Run all demos
    demo_basic_simulation
    demo_operation_logging
    demo_performance_monitoring
    demo_memory_tracking
    demo_tensor_data
    demo_configuration_comparison
    demo_performance_tests
    demo_real_model
    demo_advanced_debugging
    
    if [[ "${AUTO_MODE:-}" != "true" ]]; then
        demo_interactive
    fi
    
    print_header "Demo Complete!"
    print_success "All evaluation callback demos completed successfully"
    echo ""
    echo "Next steps:"
    echo "• Explore the source code in main.go"
    echo "• Read the detailed README.md"
    echo "• Try the Makefile commands for specific scenarios"
    echo "• Experiment with different prompts and configurations"
}

# Interactive menu for demo selection
interactive_menu() {
    while true; do
        print_header "Interactive Demo Menu"
        echo "Select a demo to run:"
        echo "  1) Basic Simulation"
        echo "  2) Operation Logging"
        echo "  3) Performance Monitoring"
        echo "  4) Memory Tracking"
        echo "  5) Tensor Data Inspection"
        echo "  6) Configuration Comparison"
        echo "  7) Performance Tests"
        echo "  8) Real Model Evaluation"
        echo "  9) Advanced Debugging"
        echo "  10) Interactive Mode"
        echo "  a) Run all demos"
        echo "  q) Quit"
        echo ""
        read -p "Your choice [1-10/a/q]: " choice
        
        case $choice in
            1) build_example && demo_basic_simulation ;;
            2) build_example && demo_operation_logging ;;
            3) build_example && demo_performance_monitoring ;;
            4) build_example && demo_memory_tracking ;;
            5) build_example && demo_tensor_data ;;
            6) build_example && demo_configuration_comparison ;;
            7) build_example && demo_performance_tests ;;
            8) build_example && demo_real_model ;;
            9) build_example && demo_advanced_debugging ;;
            10) build_example && demo_interactive ;;
            a|A) 
                export AUTO_MODE=true
                run_demo
                return
                ;;
            q|Q) exit 0 ;;
            *) print_error "Invalid choice. Please select 1-10, 'a', or 'q'." ;;
        esac
    done
}

# Handle command line arguments
case "${1:-}" in
    --auto)
        export AUTO_MODE=true
        run_demo
        ;;
    --interactive)
        build_example
        interactive_menu
        ;;
    --help|-h)
        echo "Gollama.cpp Evaluation Callback Demo Script"
        echo ""
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  --auto         Run all demos automatically"
        echo "  --interactive  Run interactive demo menu"
        echo "  --help, -h     Show this help message"
        echo ""
        echo "Without options, prompts for demo mode selection."
        ;;
    *)
        run_demo
        ;;
esac
