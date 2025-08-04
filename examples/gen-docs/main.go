package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/dianlight/gollama.cpp"
)

// FlagInfo represents a command-line flag
type FlagInfo struct {
	Name        string
	Type        string
	Default     string
	Description string
	Category    string
	Example     string
}

// ExampleInfo represents an example program
type ExampleInfo struct {
	Name        string
	Description string
	Purpose     string
	Features    []string
	Flags       []FlagInfo
	FilePath    string
}

// DocumentationGenerator generates documentation from Go source files
type DocumentationGenerator struct {
	Examples []ExampleInfo
	BaseDir  string
}

// NewDocumentationGenerator creates a new documentation generator
func NewDocumentationGenerator(baseDir string) *DocumentationGenerator {
	return &DocumentationGenerator{
		BaseDir: baseDir,
	}
}

// parseFlags extracts flag definitions from Go source code
func (dg *DocumentationGenerator) parseFlags(filePath string) ([]FlagInfo, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var flags []FlagInfo
	lines := strings.Split(string(content), "\n")

	// Regex patterns for different flag types
	flagPatterns := map[string]*regexp.Regexp{
		"String": regexp.MustCompile(`(\w+)\s*=\s*flag\.String\("([^"]+)",\s*"([^"]*)",\s*"([^"]*)"\)`),
		"Int":    regexp.MustCompile(`(\w+)\s*=\s*flag\.Int\("([^"]+)",\s*(\d+),\s*"([^"]*)"\)`),
		"Bool":   regexp.MustCompile(`(\w+)\s*=\s*flag\.Bool\("([^"]+)",\s*(true|false),\s*"([^"]*)"\)`),
		"Float":  regexp.MustCompile(`(\w+)\s*=\s*flag\.Float64\("([^"]+)",\s*([0-9.]+),\s*"([^"]*)"\)`),
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		for flagType, pattern := range flagPatterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) > 0 {
				flag := FlagInfo{
					Type: flagType,
				}

				switch flagType {
				case "String":
					flag.Name = matches[2]
					flag.Default = matches[3]
					flag.Description = matches[4]
				case "Int":
					flag.Name = matches[2]
					flag.Default = matches[3]
					flag.Description = matches[4]
				case "Bool":
					flag.Name = matches[2]
					flag.Default = matches[3]
					flag.Description = matches[4]
				case "Float":
					flag.Name = matches[2]
					flag.Default = matches[3]
					flag.Description = matches[4]
				}

				flag.Category = dg.categorizeFlag(flag.Name)
				flags = append(flags, flag)
			}
		}
	}

	return flags, nil
}

// categorizeFlag categorizes flags into common, model, context, generation, etc.
func (dg *DocumentationGenerator) categorizeFlag(flagName string) string {
	commonFlags := []string{"help", "version", "verbose", "quiet", "output", "format"}
	modelFlags := []string{"model", "vocab", "mmap", "mlock", "gpu", "offload"}
	contextFlags := []string{"ctx", "context", "batch", "ubatch", "threads", "seq"}
	generationFlags := []string{"prompt", "predict", "n-predict", "temp", "temperature", "top-k", "top-p", "repeat", "seed"}
	samplingFlags := []string{"temp", "temperature", "top-k", "top-p", "repeat", "penalty", "presence", "frequency", "mirostat", "typical"}
	embedFlags := []string{"embedding", "embed", "normalize", "similarity", "chunk"}
	retrievalFlags := []string{"query", "retrieval", "search", "index", "rank", "score"}

	flagLower := strings.ToLower(flagName)

	for _, flag := range commonFlags {
		if strings.Contains(flagLower, flag) {
			return "Common"
		}
	}

	for _, flag := range modelFlags {
		if strings.Contains(flagLower, flag) {
			return "Model"
		}
	}

	for _, flag := range contextFlags {
		if strings.Contains(flagLower, flag) {
			return "Context"
		}
	}

	for _, flag := range generationFlags {
		if strings.Contains(flagLower, flag) {
			return "Generation"
		}
	}

	for _, flag := range samplingFlags {
		if strings.Contains(flagLower, flag) {
			return "Sampling"
		}
	}

	for _, flag := range embedFlags {
		if strings.Contains(flagLower, flag) {
			return "Embedding"
		}
	}

	for _, flag := range retrievalFlags {
		if strings.Contains(flagLower, flag) {
			return "Retrieval"
		}
	}

	return "Example-specific"
}

// parseExampleInfo extracts information about an example from its source and README
func (dg *DocumentationGenerator) parseExampleInfo(exampleDir string) (ExampleInfo, error) {
	example := ExampleInfo{
		Name:     filepath.Base(exampleDir),
		FilePath: exampleDir,
	}

	// Parse main.go for flags
	mainGoPath := filepath.Join(exampleDir, "main.go")
	if _, err := os.Stat(mainGoPath); err == nil {
		flags, err := dg.parseFlags(mainGoPath)
		if err == nil {
			example.Flags = flags
		}
	}

	// Parse README.md for description and features
	readmePath := filepath.Join(exampleDir, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		content, err := os.ReadFile(readmePath)
		if err == nil {
			example.Description, example.Purpose, example.Features = dg.parseReadme(string(content))
		}
	}

	return example, nil
}

// parseReadme extracts description, purpose, and features from README content
func (dg *DocumentationGenerator) parseReadme(content string) (string, string, []string) {
	lines := strings.Split(content, "\n")

	var description, purpose string
	var features []string
	var inFeatures bool

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Extract description from first paragraph after title
		if description == "" && line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "```") {
			// Look for the first substantial paragraph
			if len(line) > 20 {
				description = line
			}
		}

		// Look for "Overview" or "Features" section
		if strings.Contains(strings.ToLower(line), "overview") ||
			strings.Contains(strings.ToLower(line), "purpose") ||
			strings.Contains(strings.ToLower(line), "demonstrates") {
			if i+1 < len(lines) {
				purpose = strings.TrimSpace(lines[i+1])
			}
		}

		// Extract features from bulleted lists
		if strings.Contains(strings.ToLower(line), "feature") && strings.HasPrefix(line, "#") {
			inFeatures = true
			continue
		}

		if inFeatures {
			if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
				feature := strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")
				if len(feature) > 0 {
					features = append(features, feature)
				}
			} else if strings.HasPrefix(line, "#") {
				inFeatures = false
			}
		}
	}

	return description, purpose, features
}

// ScanExamples scans the examples directory and extracts information
func (dg *DocumentationGenerator) ScanExamples() error {
	examplesDir := filepath.Join(dg.BaseDir, "examples")

	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			examplePath := filepath.Join(examplesDir, entry.Name())

			// Skip if no main.go exists
			mainGoPath := filepath.Join(examplePath, "main.go")
			if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
				continue
			}

			example, err := dg.parseExampleInfo(examplePath)
			if err != nil {
				log.Printf("Warning: failed to parse example %s: %v", entry.Name(), err)
				continue
			}

			dg.Examples = append(dg.Examples, example)
		}
	}

	// Sort examples by name
	sort.Slice(dg.Examples, func(i, j int) bool {
		return dg.Examples[i].Name < dg.Examples[j].Name
	})

	return nil
}

// GenerateMarkdown generates comprehensive markdown documentation
func (dg *DocumentationGenerator) GenerateMarkdown(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	// Header
	fmt.Fprintf(w, "# Gollama.cpp Examples Documentation\n\n")
	fmt.Fprintf(w, "*Generated on %s*\n\n", time.Now().Format("January 2, 2006 at 15:04:05"))
	fmt.Fprintf(w, "This documentation provides comprehensive information about all available examples in the gollama.cpp library.\n\n")

	// Table of Contents
	fmt.Fprintf(w, "## Table of Contents\n\n")
	for _, example := range dg.Examples {
		fmt.Fprintf(w, "- [%s](#%s)\n", example.Name, strings.ToLower(strings.ReplaceAll(example.Name, "-", "")))
	}
	fmt.Fprintf(w, "- [Common Parameters Reference](#common-parameters-reference)\n")
	fmt.Fprintf(w, "- [Flag Categories](#flag-categories)\n\n")

	// Examples overview
	fmt.Fprintf(w, "## Examples Overview\n\n")
	fmt.Fprintf(w, "The gollama.cpp library includes %d example programs demonstrating various capabilities:\n\n", len(dg.Examples))

	for _, example := range dg.Examples {
		description := example.Description
		if description == "" {
			description = "Example demonstrating " + example.Name + " functionality"
		}
		fmt.Fprintf(w, "- **%s**: %s\n", example.Name, description)
	}
	fmt.Fprintf(w, "\n")

	// Individual example documentation
	for _, example := range dg.Examples {
		dg.writeExampleSection(w, example)
	}

	// Common parameters reference
	dg.writeCommonParametersReference(w)

	// Flag categories
	dg.writeFlagCategories(w)

	return nil
}

// writeExampleSection writes documentation for a single example
func (dg *DocumentationGenerator) writeExampleSection(w *bufio.Writer, example ExampleInfo) {
	fmt.Fprintf(w, "## %s\n\n", example.Name)

	// Description
	if example.Description != "" {
		fmt.Fprintf(w, "%s\n\n", example.Description)
	}

	// Purpose
	if example.Purpose != "" {
		fmt.Fprintf(w, "**Purpose**: %s\n\n", example.Purpose)
	}

	// Features
	if len(example.Features) > 0 {
		fmt.Fprintf(w, "**Features**:\n")
		for _, feature := range example.Features {
			fmt.Fprintf(w, "- %s\n", feature)
		}
		fmt.Fprintf(w, "\n")
	}

	// Usage
	fmt.Fprintf(w, "### Usage\n\n")
	fmt.Fprintf(w, "```bash\n")
	fmt.Fprintf(w, "cd examples/%s\n", example.Name)
	fmt.Fprintf(w, "go run main.go [OPTIONS]\n")
	fmt.Fprintf(w, "```\n\n")

	// Parameters
	if len(example.Flags) > 0 {
		fmt.Fprintf(w, "### Parameters\n\n")

		// Group flags by category
		categories := make(map[string][]FlagInfo)
		for _, flag := range example.Flags {
			categories[flag.Category] = append(categories[flag.Category], flag)
		}

		// Sort categories
		categoryOrder := []string{"Common", "Model", "Context", "Generation", "Sampling", "Embedding", "Retrieval", "Example-specific"}
		for _, category := range categoryOrder {
			if flags, exists := categories[category]; exists {
				fmt.Fprintf(w, "#### %s Parameters\n\n", category)
				dg.writeFlagsTable(w, flags)
				fmt.Fprintf(w, "\n")
			}
		}
	}

	// Build and run instructions
	fmt.Fprintf(w, "### Build and Run\n\n")
	fmt.Fprintf(w, "```bash\n")
	fmt.Fprintf(w, "# Build\n")
	fmt.Fprintf(w, "cd examples/%s\n", example.Name)
	fmt.Fprintf(w, "go build -o %s .\n", example.Name)
	fmt.Fprintf(w, "\n# Run with default parameters\n")
	fmt.Fprintf(w, "./%s\n", example.Name)
	fmt.Fprintf(w, "\n# Run with custom parameters\n")
	if len(example.Flags) > 0 {
		// Show example usage with some common flags
		exampleFlags := []string{}
		for _, flag := range example.Flags {
			if flag.Name == "model" || flag.Name == "prompt" || flag.Name == "ctx" {
				switch flag.Type {
				case "String":
					exampleFlags = append(exampleFlags, fmt.Sprintf("-%s \"value\"", flag.Name))
				case "Int":
					exampleFlags = append(exampleFlags, fmt.Sprintf("-%s 512", flag.Name))
				case "Bool":
					exampleFlags = append(exampleFlags, fmt.Sprintf("-%s", flag.Name))
				}
				if len(exampleFlags) >= 2 {
					break
				}
			}
		}
		if len(exampleFlags) > 0 {
			fmt.Fprintf(w, "./%s %s\n", example.Name, strings.Join(exampleFlags, " "))
		}
	}
	fmt.Fprintf(w, "```\n\n")

	fmt.Fprintf(w, "---\n\n")
}

// writeFlagsTable writes a markdown table for flags
func (dg *DocumentationGenerator) writeFlagsTable(w *bufio.Writer, flags []FlagInfo) {
	fmt.Fprintf(w, "| Flag | Type | Default | Description |\n")
	fmt.Fprintf(w, "|------|------|---------|-------------|\n")

	for _, flag := range flags {
		// Escape markdown special characters
		description := strings.ReplaceAll(flag.Description, "|", "\\|")
		default_ := strings.ReplaceAll(flag.Default, "|", "\\|")

		fmt.Fprintf(w, "| `-%s` | %s | `%s` | %s |\n",
			flag.Name, flag.Type, default_, description)
	}
}

// writeCommonParametersReference writes a reference of common parameters
func (dg *DocumentationGenerator) writeCommonParametersReference(w *bufio.Writer) {
	fmt.Fprintf(w, "## Common Parameters Reference\n\n")
	fmt.Fprintf(w, "The following parameters are commonly used across multiple examples:\n\n")

	// Collect all unique flags across examples
	allFlags := make(map[string]FlagInfo)
	for _, example := range dg.Examples {
		for _, flag := range example.Flags {
			if existing, exists := allFlags[flag.Name]; !exists {
				allFlags[flag.Name] = flag
			} else {
				// Merge information if descriptions differ
				if existing.Description != flag.Description && len(flag.Description) > len(existing.Description) {
					allFlags[flag.Name] = flag
				}
			}
		}
	}

	// Group by category and write
	categories := make(map[string][]FlagInfo)
	for _, flag := range allFlags {
		categories[flag.Category] = append(categories[flag.Category], flag)
	}

	categoryOrder := []string{"Common", "Model", "Context", "Generation", "Sampling", "Embedding", "Retrieval"}
	for _, category := range categoryOrder {
		if flags, exists := categories[category]; exists {
			sort.Slice(flags, func(i, j int) bool {
				return flags[i].Name < flags[j].Name
			})

			fmt.Fprintf(w, "### %s Parameters\n\n", category)
			dg.writeFlagsTable(w, flags)
			fmt.Fprintf(w, "\n")
		}
	}
}

// writeFlagCategories writes information about flag categories
func (dg *DocumentationGenerator) writeFlagCategories(w *bufio.Writer) {
	fmt.Fprintf(w, "## Flag Categories\n\n")
	fmt.Fprintf(w, "Parameters are organized into the following categories:\n\n")

	categories := map[string]string{
		"Common":           "Basic application control flags used across multiple examples",
		"Model":            "Model loading and configuration parameters",
		"Context":          "Context size, batching, and threading parameters",
		"Generation":       "Text generation and prediction parameters",
		"Sampling":         "Sampling strategy and randomness control parameters",
		"Embedding":        "Embedding generation and processing parameters",
		"Retrieval":        "Information retrieval and search parameters",
		"Example-specific": "Parameters specific to individual examples",
	}

	for category, description := range categories {
		fmt.Fprintf(w, "### %s\n\n", category)
		fmt.Fprintf(w, "%s\n\n", description)
	}

	// Statistics
	fmt.Fprintf(w, "## Statistics\n\n")
	fmt.Fprintf(w, "- **Total Examples**: %d\n", len(dg.Examples))

	totalFlags := 0
	categoryStats := make(map[string]int)

	for _, example := range dg.Examples {
		totalFlags += len(example.Flags)
		for _, flag := range example.Flags {
			categoryStats[flag.Category]++
		}
	}

	fmt.Fprintf(w, "- **Total Parameters**: %d\n", totalFlags)
	fmt.Fprintf(w, "- **Average Parameters per Example**: %.1f\n\n", float64(totalFlags)/float64(len(dg.Examples)))

	fmt.Fprintf(w, "**Parameters by Category**:\n")
	for category, count := range categoryStats {
		fmt.Fprintf(w, "- %s: %d parameters\n", category, count)
	}
	fmt.Fprintf(w, "\n")
}

// GeneratePerExampleDocs generates individual documentation files for each example
func (dg *DocumentationGenerator) GeneratePerExampleDocs(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0750); err != nil {
		return err
	}

	for _, example := range dg.Examples {
		outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.md", example.Name))

		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}

		w := bufio.NewWriter(file)

		fmt.Fprintf(w, "# %s Example\n\n", example.Name)
		fmt.Fprintf(w, "*Generated on %s*\n\n", time.Now().Format("January 2, 2006"))

		dg.writeExampleSection(w, example)

		w.Flush()
		file.Close()
	}

	return nil
}

// GenerateUsageSummary generates a concise usage summary
func (dg *DocumentationGenerator) GenerateUsageSummary(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	fmt.Fprintf(w, "# Gollama.cpp Examples - Quick Usage Guide\n\n")
	fmt.Fprintf(w, "*Generated on %s*\n\n", time.Now().Format("January 2, 2006"))

	for _, example := range dg.Examples {
		fmt.Fprintf(w, "## %s\n\n", example.Name)
		if example.Description != "" {
			fmt.Fprintf(w, "%s\n\n", example.Description)
		}

		fmt.Fprintf(w, "```bash\n")
		fmt.Fprintf(w, "cd examples/%s && go run main.go", example.Name)

		// Add common flags if they exist
		for _, flag := range example.Flags {
			if flag.Name == "model" {
				fmt.Fprintf(w, " -model path/to/model.gguf")
				break
			}
		}
		fmt.Fprintf(w, "\n```\n\n")
	}

	return nil
}

func main() {
	// Parse command line arguments
	outputDir := "docs"
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	baseDir := "../.."
	if len(os.Args) > 2 {
		baseDir = os.Args[2]
	}

	fmt.Printf("Gollama.cpp Documentation Generator %s\n", gollama.FullVersion)
	fmt.Printf("Scanning examples in: %s\n", baseDir)
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Println()

	// Create documentation generator
	generator := NewDocumentationGenerator(baseDir)

	// Scan examples
	fmt.Print("Scanning examples... ")
	if err := generator.ScanExamples(); err != nil {
		log.Fatalf("Failed to scan examples: %v", err)
	}
	fmt.Printf("found %d examples\n", len(generator.Examples))

	// Create output directory
	if err := os.MkdirAll(outputDir, 0750); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Generate comprehensive documentation
	fmt.Print("Generating comprehensive documentation... ")
	comprehensiveDoc := filepath.Join(outputDir, "examples-reference.md")
	if err := generator.GenerateMarkdown(comprehensiveDoc); err != nil {
		log.Fatalf("Failed to generate comprehensive documentation: %v", err)
	}
	fmt.Printf("done\n")

	// Generate per-example documentation
	fmt.Print("Generating per-example documentation... ")
	exampleDocsDir := filepath.Join(outputDir, "examples")
	if err := generator.GeneratePerExampleDocs(exampleDocsDir); err != nil {
		log.Fatalf("Failed to generate per-example documentation: %v", err)
	}
	fmt.Printf("done\n")

	// Generate usage summary
	fmt.Print("Generating usage summary... ")
	usageSummary := filepath.Join(outputDir, "quick-usage.md")
	if err := generator.GenerateUsageSummary(usageSummary); err != nil {
		log.Fatalf("Failed to generate usage summary: %v", err)
	}
	fmt.Printf("done\n")

	// Print summary
	fmt.Println("\nDocumentation Generation Complete!")
	fmt.Printf("Generated documentation for %d examples:\n", len(generator.Examples))
	for _, example := range generator.Examples {
		fmt.Printf("  - %s (%d parameters)\n", example.Name, len(example.Flags))
	}

	fmt.Printf("\nOutput files:\n")
	fmt.Printf("  - %s (comprehensive reference)\n", comprehensiveDoc)
	fmt.Printf("  - %s (quick usage guide)\n", usageSummary)
	fmt.Printf("  - %s/ (individual example docs)\n", exampleDocsDir)

	fmt.Println("\nNext steps:")
	fmt.Println("  - Review generated documentation")
	fmt.Println("  - Integrate into project documentation")
	fmt.Println("  - Update as examples are added or modified")
}
