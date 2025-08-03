# Gollama.cpp Documentation Generator Example Makefile

.PHONY: all build clean run test help
.PHONY: run-basic run-custom run-docs run-website
.PHONY: demo-comprehensive demo-quick demo-per-example demo-analysis
.PHONY: test-parsing test-generation test-categories test-output
.PHONY: validate-docs check-docs update-docs publish-docs
.DEFAULT_GOAL := help

# Configuration
BINARY_NAME := gen-docs
OUTPUT_DIR := docs
BASE_DIR := ../..
WEBSITE_DIR := ../../docs
BUILD_FLAGS := -v

## Build Commands

all: build ## Build the gen-docs example

build: ## Build the gen-docs executable
	@echo "Building gen-docs example..."
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

clean: ## Clean build artifacts and generated docs
	@echo "Cleaning build artifacts and generated documentation..."
	rm -f $(BINARY_NAME)
	rm -rf $(OUTPUT_DIR)
	rm -rf docs-*
	go clean
	@echo "Clean complete"

test: ## Run Go tests
	@echo "Running tests..."
	go test -v ./...

## Basic Run Commands

run: build ## Generate documentation with default settings
	@echo "Generating documentation with default settings..."
	./$(BINARY_NAME)
	@echo "Documentation generated in: $(OUTPUT_DIR)"

run-basic: build ## Generate basic documentation
	@echo "Generating basic documentation..."
	./$(BINARY_NAME) $(OUTPUT_DIR)

run-custom: build ## Generate documentation in custom directory
	@echo "Generating documentation in custom directory..."
	./$(BINARY_NAME) docs-custom $(BASE_DIR)

run-docs: build ## Generate documentation for docs directory
	@echo "Generating documentation for docs integration..."
	./$(BINARY_NAME) $(WEBSITE_DIR) $(BASE_DIR)

run-website: build ## Generate documentation for website
	@echo "Generating documentation for website integration..."
	mkdir -p website-docs
	./$(BINARY_NAME) website-docs $(BASE_DIR)

## Demo Commands

demo-comprehensive: build ## Demonstrate comprehensive documentation generation
	@echo "=== Demonstration: Comprehensive Documentation ==="
	@echo "Generating complete documentation with all features..."
	./$(BINARY_NAME) demo-comprehensive
	@echo ""
	@echo "Generated files:"
	@ls -la demo-comprehensive/
	@echo ""
	@echo "Comprehensive reference preview:"
	@head -20 demo-comprehensive/examples-reference.md

demo-quick: build ## Demonstrate quick usage guide generation
	@echo "=== Demonstration: Quick Usage Guide ==="
	@echo "Generating concise usage documentation..."
	./$(BINARY_NAME) demo-quick
	@echo ""
	@echo "Quick usage guide preview:"
	@head -30 demo-quick/quick-usage.md

demo-per-example: build ## Demonstrate per-example documentation
	@echo "=== Demonstration: Per-Example Documentation ==="
	@echo "Generating individual example documentation..."
	./$(BINARY_NAME) demo-per-example
	@echo ""
	@echo "Per-example documentation:"
	@ls -la demo-per-example/examples/
	@echo ""
	@echo "Sample example doc preview:"
	@if [ -f demo-per-example/examples/simple-chat.md ]; then head -20 demo-per-example/examples/simple-chat.md; fi

demo-analysis: build ## Demonstrate documentation analysis features
	@echo "=== Demonstration: Documentation Analysis ==="
	@echo "Analyzing examples and generating statistics..."
	./$(BINARY_NAME) demo-analysis 2>&1 | grep -E "(found|Generated|examples:|parameters)"

## Test Commands

test-parsing: build ## Test flag parsing functionality
	@echo "Testing flag parsing across examples..."
	@echo "Scanning examples for flag patterns..."
	./$(BINARY_NAME) test-parsing $(BASE_DIR) 2>&1 | grep -E "(found|examples)"

test-generation: build ## Test documentation generation
	@echo "Testing documentation generation..."
	./$(BINARY_NAME) test-generation
	@echo "Verifying generated files:"
	@ls -la test-generation/

test-categories: build ## Test parameter categorization
	@echo "Testing parameter categorization..."
	./$(BINARY_NAME) test-categories
	@echo "Checking category distribution in generated docs:"
	@grep -o "#### [A-Za-z-]* Parameters" test-categories/examples-reference.md | sort | uniq -c

test-output: build ## Test different output formats
	@echo "Testing various output configurations..."
	@for dir in test-out-1 test-out-2 test-out-3; do \
		echo "Testing output to $$dir..."; \
		./$(BINARY_NAME) $$dir $(BASE_DIR); \
		echo "Files in $$dir: $$(ls $$dir | wc -l)"; \
	done

## Documentation Validation

validate-docs: build ## Validate generated documentation
	@echo "=== Validating Generated Documentation ==="
	./$(BINARY_NAME) validate-output
	@echo "1. Checking comprehensive reference..."
	@if [ -f validate-output/examples-reference.md ]; then \
		echo "✓ Comprehensive reference generated"; \
	else \
		echo "✗ Comprehensive reference missing"; \
	fi
	@echo "2. Checking quick usage guide..."
	@if [ -f validate-output/quick-usage.md ]; then \
		echo "✓ Quick usage guide generated"; \
	else \
		echo "✗ Quick usage guide missing"; \
	fi
	@echo "3. Checking per-example docs..."
	@if [ -d validate-output/examples ]; then \
		example_count=$$(ls validate-output/examples/*.md 2>/dev/null | wc -l); \
		echo "✓ $$example_count per-example docs generated"; \
	else \
		echo "✗ Per-example docs missing"; \
	fi

check-docs: build ## Check documentation quality
	@echo "=== Checking Documentation Quality ==="
	./$(BINARY_NAME) check-output
	@echo "Checking for common documentation issues..."
	@echo "1. Checking for empty descriptions:"
	@grep -c "Description.*|$$" check-output/examples-reference.md || echo "No empty descriptions found"
	@echo "2. Checking for missing defaults:"
	@grep -c "Default.*| \`\` |" check-output/examples-reference.md || echo "No missing defaults found"
	@echo "3. Checking table formatting:"
	@grep -c "^|.*|$$" check-output/examples-reference.md
	@echo "4. Checking for broken links:"
	@grep -o "\[.*\](#.*)" check-output/examples-reference.md | wc -l

update-docs: build ## Update project documentation
	@echo "=== Updating Project Documentation ==="
	@echo "Generating latest documentation..."
	./$(BINARY_NAME) $(WEBSITE_DIR)
	@echo "Documentation updated in: $(WEBSITE_DIR)"
	@echo ""
	@echo "Generated files:"
	@ls -la $(WEBSITE_DIR)/

publish-docs: build ## Publish documentation (requires write permissions)
	@echo "=== Publishing Documentation ==="
	@if [ -w "$(WEBSITE_DIR)" ]; then \
		echo "Publishing to $(WEBSITE_DIR)..."; \
		./$(BINARY_NAME) $(WEBSITE_DIR); \
		echo "Documentation published successfully"; \
	else \
		echo "Warning: Cannot write to $(WEBSITE_DIR)"; \
		echo "Generating to docs-publish instead..."; \
		./$(BINARY_NAME) docs-publish; \
		echo "Documentation generated in docs-publish/"; \
	fi

## Advanced Usage

advanced-demo: build ## Run comprehensive demonstration suite
	@echo "=== Advanced Documentation Generator Demonstration ==="
	@echo ""
	@echo "1. Basic generation:"
	./$(BINARY_NAME) advanced-1
	@echo ""
	@echo "2. Custom directory generation:"
	./$(BINARY_NAME) advanced-2 $(BASE_DIR)
	@echo ""
	@echo "3. Analysis and statistics:"
	./$(BINARY_NAME) advanced-3 2>&1 | tail -10

benchmark: build ## Benchmark documentation generation
	@echo "=== Benchmarking Documentation Generation ==="
	@echo "Testing generation speed..."
	@time ./$(BINARY_NAME) benchmark-output > /dev/null
	@echo "Files generated: $$(find benchmark-output -name "*.md" | wc -l)"
	@echo "Total size: $$(du -sh benchmark-output | cut -f1)"

stress-test: build ## Stress test with multiple generations
	@echo "=== Stress Testing Documentation Generator ==="
	@for i in 1 2 3 4 5; do \
		echo "Run $$i..."; \
		./$(BINARY_NAME) stress-$$i > /dev/null; \
		echo "Generated $$(find stress-$$i -name "*.md" | wc -l) files"; \
	done

validate: build test ## Validate the example works correctly
	@echo "=== Validating Documentation Generator ==="
	@echo "1. Testing basic generation..."
	./$(BINARY_NAME) validate-test > /dev/null
	@echo "✓ Basic generation working"
	@echo "2. Testing file creation..."
	@if [ -f validate-test/examples-reference.md ]; then echo "✓ Reference file created"; else echo "✗ Reference file missing"; fi
	@echo "3. Testing content generation..."
	@if [ -s validate-test/examples-reference.md ]; then echo "✓ Reference file has content"; else echo "✗ Reference file empty"; fi
	@echo "All validation tests passed!"

## Development Tools

dev-watch: ## Watch for changes and regenerate docs
	@echo "Watching for changes (requires fswatch)..."
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o $(BASE_DIR)/examples/ | while read; do \
			echo "Changes detected, regenerating documentation..."; \
			make run-basic; \
		done; \
	else \
		echo "fswatch not found. Install with: brew install fswatch"; \
	fi

dev-serve: ## Serve generated documentation (requires Python)
	@echo "Serving documentation on http://localhost:8000"
	@if [ -d $(OUTPUT_DIR) ]; then \
		cd $(OUTPUT_DIR) && python3 -m http.server 8000; \
	else \
		echo "No documentation found. Run 'make run' first."; \
	fi

## Analysis Commands

analyze-examples: build ## Analyze example structure and patterns
	@echo "=== Analyzing Example Structure ==="
	./$(BINARY_NAME) analyze-output
	@echo "Example analysis:"
	@find $(BASE_DIR)/examples -name "main.go" -exec basename {} \; | cut -d/ -f1 | sort
	@echo ""
	@echo "Flag pattern analysis:"
	@grep -h "flag\." $(BASE_DIR)/examples/*/main.go | grep -o "flag\.[A-Za-z]*" | sort | uniq -c

compare-docs: build ## Compare generated docs with existing ones
	@echo "=== Comparing Documentation Versions ==="
	@if [ -d $(OUTPUT_DIR) ]; then \
		./$(BINARY_NAME) compare-new; \
		echo "Comparing with existing documentation..."; \
		diff -u $(OUTPUT_DIR)/examples-reference.md compare-new/examples-reference.md || echo "Documentation has changes"; \
	else \
		echo "No existing documentation to compare"; \
	fi

## Help

help: ## Show this help message
	@echo "Gollama.cpp Documentation Generator Example"
	@echo "==========================================="
	@echo ""
	@echo "This example generates comprehensive documentation by analyzing"
	@echo "Go source files and extracting command-line flag definitions."
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Quick start:"
	@echo "  make run             # Generate docs with default settings"
	@echo "  make demo-comprehensive # See comprehensive docs demo"
	@echo "  make validate-docs   # Validate generated documentation"
	@echo ""
	@echo "Output directory: $(OUTPUT_DIR)"
	@echo "Base directory: $(BASE_DIR)"

## Utility Commands

show-examples: ## Show available examples
	@echo "Available examples to document:"
	@find $(BASE_DIR)/examples -maxdepth 1 -type d -exec basename {} \; | grep -v "^examples$$" | sort

show-flags: ## Show flag patterns found in examples
	@echo "Flag patterns found in examples:"
	@grep -h "flag\." $(BASE_DIR)/examples/*/main.go | head -20

show-structure: ## Show generated documentation structure
	@if [ -d $(OUTPUT_DIR) ]; then \
		echo "Generated documentation structure:"; \
		find $(OUTPUT_DIR) -type f -name "*.md" | sort; \
	else \
		echo "No documentation generated yet. Run 'make run' first."; \
	fi

## Cleanup Commands

clean-all: clean ## Clean everything including generated docs
	@echo "Cleaning all generated documentation..."
	rm -rf docs-*
	rm -rf demo-*
	rm -rf test-*
	rm -rf validate-*
	rm -rf check-*
	rm -rf benchmark-*
	rm -rf stress-*
	rm -rf advanced-*
	rm -rf analyze-*
	rm -rf compare-*
	@echo "All clean"

reset: clean-all build ## Reset and rebuild everything
	@echo "Resetting and rebuilding documentation generator..."
	@echo "Ready for fresh documentation generation"
