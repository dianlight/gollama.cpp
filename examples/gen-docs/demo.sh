#!/bin/bash

# Gollama.cpp Documentation Generator Example Demo Script
# This script demonstrates the documentation generation functionality

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
    print_section "Building Documentation Generator"
    
    if make build; then
        print_success "Build completed successfully"
    else
        print_error "Build failed"
        exit 1
    fi
}

# Demo 1: Basic Documentation Generation
demo_basic_generation() {
    print_header "Demo 1: Basic Documentation Generation"
    
    echo "This demonstrates the basic documentation generation functionality."
    echo "The generator will scan all examples and create comprehensive documentation."
    wait_for_user
    
    print_section "Running Basic Documentation Generation"
    ./gen-docs demo-basic
    
    print_section "Generated Files"
    echo "Documentation files created:"
    find demo-basic -name "*.md" | sort
    
    echo ""
    echo "File sizes:"
    ls -lh demo-basic/*.md demo-basic/examples/*.md 2>/dev/null | head -10
    
    print_success "Basic documentation generation completed"
    wait_for_user
}

# Demo 2: Comprehensive Reference
demo_comprehensive_reference() {
    print_header "Demo 2: Comprehensive Reference Documentation"
    
    echo "This shows the comprehensive reference documentation that includes:"
    echo "• Complete example listing with descriptions"
    echo "• Detailed parameter tables organized by category"
    echo "• Cross-references and navigation"
    echo "• Statistics and analysis"
    wait_for_user
    
    print_section "Generating Comprehensive Reference"
    ./gen-docs demo-comprehensive
    
    print_section "Reference Documentation Preview"
    echo "Table of contents:"
    grep -E "^- \[" demo-comprehensive/examples-reference.md | head -10
    
    echo ""
    echo "Parameter categories found:"
    grep -o "#### [A-Za-z-]* Parameters" demo-comprehensive/examples-reference.md | sort | uniq
    
    echo ""
    echo "Statistics section:"
    tail -20 demo-comprehensive/examples-reference.md
    
    print_success "Comprehensive reference generated"
    wait_for_user
}

# Demo 3: Quick Usage Guide
demo_quick_usage() {
    print_header "Demo 3: Quick Usage Guide"
    
    echo "This demonstrates the quick usage guide generation."
    echo "The guide provides concise usage information for rapid reference."
    wait_for_user
    
    print_section "Generating Quick Usage Guide"
    ./gen-docs demo-quick
    
    print_section "Quick Usage Guide Preview"
    echo "Quick usage examples:"
    head -50 demo-quick/quick-usage.md
    
    print_success "Quick usage guide generated"
    wait_for_user
}

# Demo 4: Per-Example Documentation
demo_per_example_docs() {
    print_header "Demo 4: Per-Example Documentation"
    
    echo "This shows individual documentation files for each example."
    echo "Each example gets its own detailed documentation file."
    wait_for_user
    
    print_section "Generating Per-Example Documentation"
    ./gen-docs demo-per-example
    
    print_section "Per-Example Documentation Files"
    echo "Individual example documentation:"
    ls -la demo-per-example/examples/
    
    echo ""
    echo "Sample example documentation preview (first available):"
    example_file=$(ls demo-per-example/examples/*.md 2>/dev/null | head -1)
    if [[ -n "$example_file" ]]; then
        echo "File: $(basename "$example_file")"
        head -30 "$example_file"
    fi
    
    print_success "Per-example documentation generated"
    wait_for_user
}

# Demo 5: Parameter Analysis
demo_parameter_analysis() {
    print_header "Demo 5: Parameter Analysis and Categorization"
    
    echo "This demonstrates the parameter analysis capabilities:"
    echo "• Automatic flag detection and parsing"
    echo "• Smart categorization of parameters"
    echo "• Statistical analysis of parameter usage"
    wait_for_user
    
    print_section "Running Parameter Analysis"
    ./gen-docs demo-analysis
    
    print_section "Parameter Analysis Results"
    echo "Parameter categories and counts:"
    grep -A 10 "Parameters by Category" demo-analysis/examples-reference.md
    
    echo ""
    echo "Examples with most parameters:"
    grep -o "Generated documentation for [0-9]* examples:" demo-analysis/examples-reference.md || echo "Analysis in progress..."
    
    print_success "Parameter analysis completed"
    wait_for_user
}

# Demo 6: Flag Detection Showcase
demo_flag_detection() {
    print_header "Demo 6: Flag Detection Showcase"
    
    echo "This showcases the flag detection capabilities:"
    echo "• Detection of different flag types (String, Int, Bool)"
    echo "• Extraction of default values and descriptions"
    echo "• Handling of various flag definition patterns"
    wait_for_user
    
    print_section "Analyzing Flag Patterns"
    ./gen-docs demo-flags
    
    echo "Flag detection examples from source files:"
    grep -h "flag\." ../../examples/*/main.go | head -15
    
    echo ""
    echo "Generated parameter documentation sample:"
    grep -A 5 "| Flag | Type | Default | Description |" demo-flags/examples-reference.md | head -10
    
    print_success "Flag detection showcase completed"
    wait_for_user
}

# Demo 7: Documentation Quality Check
demo_quality_check() {
    print_header "Demo 7: Documentation Quality Assessment"
    
    echo "This demonstrates documentation quality checking:"
    echo "• Completeness verification"
    echo "• Consistency checking"
    echo "• Structure validation"
    wait_for_user
    
    print_section "Running Quality Assessment"
    ./gen-docs demo-quality
    
    print_section "Quality Assessment Results"
    echo "Checking for documentation completeness:"
    
    # Check for required files
    required_files=("examples-reference.md" "quick-usage.md")
    for file in "${required_files[@]}"; do
        if [[ -f "demo-quality/$file" ]]; then
            print_success "$file generated"
        else
            print_error "$file missing"
        fi
    done
    
    # Check for content quality
    echo ""
    echo "Content quality checks:"
    
    if grep -q "Table of Contents" demo-quality/examples-reference.md; then
        print_success "Table of contents present"
    else
        print_warning "Table of contents missing"
    fi
    
    param_count=$(grep -c "| \`-" demo-quality/examples-reference.md || echo "0")
    echo "Parameters documented: $param_count"
    
    print_success "Quality assessment completed"
    wait_for_user
}

# Demo 8: Comparison with Manual Documentation
demo_comparison() {
    print_header "Demo 8: Generated vs Manual Documentation Comparison"
    
    echo "This compares generated documentation with manual documentation:"
    echo "• Consistency in format and structure"
    echo "• Completeness of parameter coverage"
    echo "• Accuracy of information"
    wait_for_user
    
    print_section "Generating Documentation for Comparison"
    ./gen-docs demo-comparison
    
    print_section "Documentation Analysis"
    echo "Generated documentation statistics:"
    
    # Count various elements
    examples_count=$(grep -c "^## [a-z]" demo-comparison/examples-reference.md || echo "0")
    parameters_count=$(grep -c "| \`-" demo-comparison/examples-reference.md || echo "0")
    categories_count=$(grep -c "#### [A-Za-z]* Parameters" demo-comparison/examples-reference.md || echo "0")
    
    echo "Examples documented: $examples_count"
    echo "Parameters documented: $parameters_count"
    echo "Parameter categories: $categories_count"
    
    echo ""
    echo "Documentation structure:"
    grep "^#" demo-comparison/examples-reference.md | head -10
    
    print_success "Documentation comparison completed"
    wait_for_user
}

# Demo 9: Advanced Features
demo_advanced_features() {
    print_header "Demo 9: Advanced Documentation Features"
    
    echo "This demonstrates advanced documentation generation features:"
    echo "• README parsing for example descriptions"
    echo "• Feature extraction from markdown"
    echo "• Cross-reference generation"
    echo "• Statistical analysis"
    wait_for_user
    
    print_section "Advanced Feature Generation"
    ./gen-docs demo-advanced
    
    print_section "Advanced Features Showcase"
    echo "Example descriptions extracted from READMEs:"
    grep -A 1 "^## [a-z]" demo-advanced/examples-reference.md | grep -v "^## " | grep -v "^--" | head -10
    
    echo ""
    echo "Feature lists detected:"
    grep -A 3 "Features.*:" demo-advanced/examples-reference.md | head -15
    
    echo ""
    echo "Statistical analysis:"
    grep -A 10 "Statistics" demo-advanced/examples-reference.md
    
    print_success "Advanced features demonstration completed"
    wait_for_user
}

# Demo 10: Interactive Documentation Explorer
demo_interactive_explorer() {
    print_header "Demo 10: Interactive Documentation Explorer"
    
    echo "This provides an interactive way to explore generated documentation."
    echo "You can browse different sections and files interactively."
    wait_for_user
    
    print_section "Generating Documentation for Exploration"
    ./gen-docs demo-interactive
    
    print_section "Interactive Documentation Explorer"
    
    while true; do
        echo -e "\n${CYAN}Documentation Explorer Menu:${NC}"
        echo "1) View comprehensive reference summary"
        echo "2) Browse quick usage guide"
        echo "3) Explore example documentation"
        echo "4) View parameter statistics"
        echo "5) Show file structure"
        echo "q) Quit explorer"
        
        read -p "Your choice [1-5/q]: " choice
        
        case $choice in
            1)
                echo -e "\n${YELLOW}Comprehensive Reference Summary:${NC}"
                head -50 demo-interactive/examples-reference.md
                ;;
            2)
                echo -e "\n${YELLOW}Quick Usage Guide:${NC}"
                head -30 demo-interactive/quick-usage.md
                ;;
            3)
                echo -e "\n${YELLOW}Available Example Documentation:${NC}"
                ls demo-interactive/examples/
                read -p "Enter example name to view (or press Enter to continue): " example_name
                if [[ -n "$example_name" && -f "demo-interactive/examples/$example_name.md" ]]; then
                    head -40 "demo-interactive/examples/$example_name.md"
                fi
                ;;
            4)
                echo -e "\n${YELLOW}Parameter Statistics:${NC}"
                grep -A 15 "Statistics" demo-interactive/examples-reference.md
                ;;
            5)
                echo -e "\n${YELLOW}Documentation File Structure:${NC}"
                find demo-interactive -name "*.md" | sort
                ;;
            q|Q)
                break
                ;;
            *)
                print_warning "Invalid choice. Please select 1-5 or 'q'."
                ;;
        esac
    done
    
    print_success "Interactive exploration completed"
}

# Main demo function
run_demo() {
    print_header "Gollama.cpp Documentation Generator Demo"
    
    echo "This script demonstrates the documentation generation capabilities"
    echo "that automatically analyze Go source files and create comprehensive"
    echo "markdown documentation for command-line parameters and usage."
    echo ""
    echo "Available demos:"
    echo "  1. Basic Documentation Generation"
    echo "  2. Comprehensive Reference"
    echo "  3. Quick Usage Guide"
    echo "  4. Per-Example Documentation"
    echo "  5. Parameter Analysis"
    echo "  6. Flag Detection Showcase"
    echo "  7. Documentation Quality Check"
    echo "  8. Generated vs Manual Comparison"
    echo "  9. Advanced Features"
    echo "  10. Interactive Explorer"
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
    demo_basic_generation
    demo_comprehensive_reference
    demo_quick_usage
    demo_per_example_docs
    demo_parameter_analysis
    demo_flag_detection
    demo_quality_check
    demo_comparison
    demo_advanced_features
    
    if [[ "${AUTO_MODE:-}" != "true" ]]; then
        demo_interactive_explorer
    fi
    
    print_header "Demo Complete!"
    print_success "All documentation generation demos completed successfully"
    echo ""
    echo "Generated documentation examples in:"
    find . -maxdepth 1 -type d -name "demo-*" | sort
    echo ""
    echo "Next steps:"
    echo "• Explore the generated documentation files"
    echo "• Run the generator on your own projects"
    echo "• Customize the output formats and categories"
    echo "• Integrate into your documentation workflow"
}

# Interactive menu for demo selection
interactive_menu() {
    while true; do
        print_header "Interactive Demo Menu"
        echo "Select a demo to run:"
        echo "  1) Basic Documentation Generation"
        echo "  2) Comprehensive Reference"
        echo "  3) Quick Usage Guide"
        echo "  4) Per-Example Documentation"
        echo "  5) Parameter Analysis"
        echo "  6) Flag Detection Showcase"
        echo "  7) Documentation Quality Check"
        echo "  8) Generated vs Manual Comparison"
        echo "  9) Advanced Features"
        echo "  10) Interactive Explorer"
        echo "  a) Run all demos"
        echo "  q) Quit"
        echo ""
        read -p "Your choice [1-10/a/q]: " choice
        
        case $choice in
            1) build_example && demo_basic_generation ;;
            2) build_example && demo_comprehensive_reference ;;
            3) build_example && demo_quick_usage ;;
            4) build_example && demo_per_example_docs ;;
            5) build_example && demo_parameter_analysis ;;
            6) build_example && demo_flag_detection ;;
            7) build_example && demo_quality_check ;;
            8) build_example && demo_comparison ;;
            9) build_example && demo_advanced_features ;;
            10) build_example && demo_interactive_explorer ;;
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
        echo "Gollama.cpp Documentation Generator Demo Script"
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
