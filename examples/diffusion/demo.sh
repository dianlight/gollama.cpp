#!/bin/bash

# Gollama.cpp Diffusion Generation Interactive Demo
# This script demonstrates various aspects of diffusion-based text generation

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Helper function to print colored headers
print_header() {
    echo -e "\n${CYAN}================================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}================================================${NC}\n"
}

print_section() {
    echo -e "\n${YELLOW}$1${NC}\n"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if binary exists
if [ ! -f "./diffusion" ]; then
    echo -e "${RED}❌ Binary not found. Building...${NC}"
    make build
fi

print_header "🎯 Gollama.cpp Diffusion Generation Interactive Demo"

echo -e "Welcome to the diffusion generation demo! This interactive script will guide you"
echo -e "through various aspects of diffusion-based text generation.\n"

echo -e "Diffusion generation creates text by iteratively refining masked tokens through"
echo -e "multiple denoising steps, using confidence algorithms to determine which tokens"
echo -e "to unmask at each step.\n"

print_warning "Note: This is a conceptual demonstration working within standard chat model constraints."

# Function to wait for user input
wait_for_user() {
    echo -e "\n${MAGENTA}Press Enter to continue...${NC}"
    read -r
}

# Demo 1: Basic Diffusion Generation
print_section "📝 Demo 1: Basic Diffusion Generation"
print_info "Let's start with a simple example using the confidence-based algorithm."
print_info "This will generate text in 8 diffusion steps, showing the basic process."

echo -e "\nCommand: ./diffusion -prompt \"The future of technology\" -diffusion-steps 8 -max-length 48"
wait_for_user

./diffusion -prompt "The future of technology" -diffusion-steps 8 -max-length 48

# Demo 2: Visual Mode
print_section "🎬 Demo 2: Visual Generation Mode"
print_info "Visual mode shows the generation process in real-time."
print_info "You'll see masked positions (underscores) being filled step by step."

echo -e "\nCommand: ./diffusion -prompt \"Machine learning enables\" -diffusion-visual -diffusion-steps 6 -max-length 40"
print_info "Watch as the text is progressively revealed!"
wait_for_user

./diffusion -prompt "Machine learning enables" -diffusion-visual -diffusion-steps 6 -max-length 40

# Demo 3: Algorithm Comparison
print_section "🧠 Demo 3: Algorithm Comparison"
print_info "Different algorithms use different strategies for selecting tokens to unmask."
print_info "Let's compare confidence-based vs entropy-based algorithms."

echo -e "\nTesting prompt: \"Artificial intelligence will\""

echo -e "\n${BLUE}Confidence-based algorithm (uses token probability):${NC}"
./diffusion -prompt "Artificial intelligence will" -diffusion-algorithm 0 -diffusion-steps 6 -max-length 40

wait_for_user

echo -e "\n${BLUE}Entropy-based algorithm (uses distribution entropy):${NC}"
./diffusion -prompt "Artificial intelligence will" -diffusion-algorithm 1 -diffusion-steps 6 -max-length 40

# Demo 4: Verbose Mode
print_section "🔍 Demo 4: Verbose Analysis")
print_info "Verbose mode shows detailed information about each diffusion step."
print_info "You'll see token counts, step progression, and timing information."

echo -e "\nCommand: ./diffusion -prompt \"Scientific innovation\" -verbose -diffusion-steps 5 -max-length 36"
wait_for_user

./diffusion -prompt "Scientific innovation" -verbose -diffusion-steps 5 -max-length 36

# Demo 5: Step Count Effects
print_section "📊 Demo 5: Effect of Step Count"
print_info "The number of diffusion steps affects the quality and refinement of generation."
print_info "More steps generally produce more coherent results but take longer."

prompt="The evolution of computing"

echo -e "\nTesting prompt: \"$prompt\""

echo -e "\n${BLUE}Few steps (3 - faster, less refined):${NC}"
./diffusion -prompt "$prompt" -diffusion-steps 3 -max-length 44

wait_for_user

echo -e "\n${BLUE}Medium steps (8 - balanced):${NC}"
./diffusion -prompt "$prompt" -diffusion-steps 8 -max-length 44

wait_for_user

echo -e "\n${BLUE}Many steps (15 - slower, more refined):${NC}"
./diffusion -prompt "$prompt" -diffusion-steps 15 -max-length 44

# Demo 6: Temperature Effects
print_section "🌡️  Demo 6: Temperature Effects"
print_info "Temperature affects the randomness of token selection."
print_info "Higher temperature produces more creative but less predictable results."

prompt="Once upon a time in a magical land"

echo -e "\nTesting prompt: \"$prompt\""

echo -e "\n${BLUE}Low temperature (0.5 - more deterministic):${NC}"
./diffusion -prompt "$prompt" -temperature 0.5 -diffusion-steps 6 -max-length 48

wait_for_user

echo -e "\n${BLUE}High temperature (1.2 - more creative):${NC}"
./diffusion -prompt "$prompt" -temperature 1.2 -diffusion-steps 6 -max-length 48

# Demo 7: Deterministic Generation
print_section "🔒 Demo 7: Deterministic Generation"
print_info "Using a fixed seed produces reproducible results."
print_info "This is useful for testing and reproducible research."

prompt="Reproducible AI research"

echo -e "\nTesting prompt: \"$prompt\" with seed 12345"

echo -e "\n${BLUE}First run:${NC}"
./diffusion -prompt "$prompt" -seed 12345 -diffusion-steps 6 -max-length 40

wait_for_user

echo -e "\n${BLUE}Second run (should be identical):${NC}"
./diffusion -prompt "$prompt" -seed 12345 -diffusion-steps 6 -max-length 40

# Demo 8: Interactive Custom Prompt
print_section "🎮 Demo 8: Interactive Custom Prompt"
print_info "Now it's your turn! Enter a custom prompt and see diffusion in action."

while true; do
    echo -e "\n${YELLOW}Enter your prompt (or 'quit' to exit):${NC}"
    read -r user_prompt
    
    if [ "$user_prompt" = "quit" ] || [ "$user_prompt" = "exit" ] || [ "$user_prompt" = "q" ]; then
        break
    fi
    
    if [ -z "$user_prompt" ]; then
        print_warning "Empty prompt. Using default."
        user_prompt="The future is"
    fi
    
    echo -e "\n${BLUE}Generating with visual mode enabled...${NC}"
    ./diffusion -prompt "$user_prompt" -diffusion-visual -diffusion-steps 8 -max-length 56
    
    echo -e "\n${MAGENTA}Try another prompt? (Enter to continue, 'quit' to exit)${NC}"
done

# Conclusion
print_header "🎉 Demo Complete!"

echo -e "You've experienced the key concepts of diffusion-based text generation:"
echo -e ""
echo -e "${GREEN}✅ Basic iterative generation process${NC}"
echo -e "${GREEN}✅ Visual real-time generation display${NC}"
echo -e "${GREEN}✅ Different confidence algorithms${NC}"
echo -e "${GREEN}✅ Effect of step count and temperature${NC}"
echo -e "${GREEN}✅ Deterministic and reproducible generation${NC}"
echo -e "${GREEN}✅ Interactive custom prompt testing${NC}"
echo -e ""

print_section "🚀 Next Steps"
echo -e "Explore more with these commands:"
echo -e ""
echo -e "${CYAN}make explain-algorithms${NC}    - Detailed algorithm comparison"
echo -e "${CYAN}make creative${NC}              - Creative writing generation"
echo -e "${CYAN}make benchmark${NC}             - Performance testing"
echo -e "${CYAN}make help${NC}                  - See all available demos"
echo -e ""
echo -e "${CYAN}./diffusion --help${NC}         - View all command-line options"
echo -e ""

print_section "📚 Understanding Diffusion"
echo -e "This demonstration shows conceptual diffusion principles. In production:"
echo -e ""
echo -e "• True diffusion models use specialized architectures (Dream, LLaDA)"
echo -e "• Non-causal attention allows bidirectional context"
echo -e "• Proper mask tokens enable more sophisticated masking"
echo -e "• Advanced scheduling improves generation quality"
echo -e ""

print_info "For more details, see README.md or visit the llama.cpp diffusion documentation."

print_success "Thank you for trying the diffusion generation demo!"
echo -e ""
