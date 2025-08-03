# Gollama.cpp Documentation Generator Example

This example is a documentation generator that automatically analyzes all examples in the gollama.cpp project and creates comprehensive Markdown documentation. It parses Go source files to extract command-line flag definitions and generates structured documentation similar to the llama.cpp gen-docs utility.

## Overview

The gen-docs example demonstrates:

1. **Source Code Analysis**: Parsing Go source files to extract flag definitions
2. **Documentation Generation**: Creating comprehensive Markdown documentation
3. **Parameter Categorization**: Organizing flags into logical categories
4. **Multi-format Output**: Generating different documentation formats
5. **Automation**: Automated documentation generation workflow

## Features

### Source Code Analysis
- **Flag Detection**: Automatically finds `flag.String()`, `flag.Int()`, `flag.Bool()` calls
- **Parameter Extraction**: Extracts flag names, types, defaults, and descriptions
- **README Parsing**: Analyzes README.md files for example descriptions and features
- **Smart Categorization**: Automatically categorizes flags into logical groups

### Documentation Generation
- **Comprehensive Reference**: Complete documentation with all examples and parameters
- **Per-Example Docs**: Individual documentation files for each example
- **Quick Usage Guide**: Concise usage summary for rapid reference
- **Structured Output**: Well-organized Markdown with tables and sections

### Parameter Categories
- **Common**: Basic application control flags (help, verbose, output)
- **Model**: Model loading and configuration (model path, vocab, memory mapping)
- **Context**: Context size, batching, threading parameters
- **Generation**: Text generation and prediction parameters
- **Sampling**: Sampling strategies and randomness control
- **Embedding**: Embedding generation and processing
- **Retrieval**: Information retrieval and search parameters
- **Example-specific**: Parameters unique to individual examples

## Usage

### Basic Usage
```bash
# Generate documentation in default 'docs' directory
go run main.go

# Generate documentation in custom directory
go run main.go output-dir

# Specify custom base directory for examples
go run main.go docs ../..
```

### Advanced Usage
```bash
# Generate docs with specific output structure
go run main.go documentation-output ../../

# Run from different location
cd examples/gen-docs
go run main.go docs ../../
```

## Command Line Options

| Option | Description |
|--------|-------------|
| `output-dir` | Output directory for generated documentation (default: "docs") |
| `base-dir` | Base directory containing examples to scan (default: "../..") |

## Output Files

The generator creates several documentation files:

### 1. Comprehensive Reference (`examples-reference.md`)
- Complete documentation for all examples
- Table of contents with navigation links
- Detailed parameter tables organized by category
- Build and usage instructions for each example
- Common parameters reference
- Statistics and category information

### 2. Quick Usage Guide (`quick-usage.md`)
- Concise usage summary for all examples
- Essential command-line examples
- Quick reference format
- Minimal but complete information

### 3. Per-Example Documentation (`examples/`)
- Individual `.md` file for each example
- Focused documentation for single examples
- Detailed parameter tables
- Usage examples and build instructions

## Example Output

### Comprehensive Reference Structure
```markdown
# Gollama.cpp Examples Documentation

## Table of Contents
- [batched](#batched)
- [diffusion](#diffusion)
- [embedding](#embedding)
...

## Examples Overview
The gollama.cpp library includes 7 example programs...

## batched
Demonstrates batched generation concepts...

### Parameters
#### Common Parameters
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-model` | String | `path/to/model.gguf` | Path to the GGUF model file |
...
```

### Parameter Table Format
```markdown
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-model` | String | `../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf` | Path to the GGUF model file |
| `-prompt` | String | `"The future of AI is"` | Prompt text to generate from |
| `-threads` | Int | `4` | Number of threads to use |
| `-ctx` | Int | `2048` | Context size |
```

## Implementation Details

### Flag Detection Algorithm
The generator uses regular expressions to detect flag definitions:

```go
flagPatterns := map[string]*regexp.Regexp{
    "String": regexp.MustCompile(`flag\.String\("([^"]+)",\s*"([^"]*)",\s*"([^"]*)"\)`),
    "Int":    regexp.MustCompile(`flag\.Int\("([^"]+)",\s*(\d+),\s*"([^"]*)"\)`),
    "Bool":   regexp.MustCompile(`flag\.Bool\("([^"]+)",\s*(true|false),\s*"([^"]*)"\)`),
}
```

### Categorization Logic
Flags are automatically categorized based on naming patterns:

```go
commonFlags := []string{"help", "version", "verbose", "quiet"}
modelFlags := []string{"model", "vocab", "mmap", "mlock", "gpu"}
contextFlags := []string{"ctx", "context", "batch", "threads"}
generationFlags := []string{"prompt", "predict", "temp", "top-k"}
```

### README Analysis
The generator parses README.md files to extract:
- Example descriptions from the first substantial paragraph
- Purpose statements from "Overview" sections
- Feature lists from bulleted sections
- Additional context for documentation

## Generated Documentation Features

### Navigation and Structure
- **Table of Contents**: Links to all examples and sections
- **Category Organization**: Parameters grouped by functionality
- **Cross-references**: Links between related sections
- **Consistent Formatting**: Standardized markdown structure

### Comprehensive Information
- **Complete Parameter Lists**: All flags with types, defaults, descriptions
- **Usage Examples**: Command-line examples for each scenario
- **Build Instructions**: Step-by-step compilation and execution
- **Feature Summaries**: Key capabilities of each example

### Statistics and Analysis
- **Example Counts**: Total number of examples and parameters
- **Category Breakdown**: Parameters organized by functional area
- **Usage Patterns**: Common parameter patterns across examples
- **Complexity Metrics**: Average parameters per example

## Integration with Development Workflow

### Automated Documentation
- **CI/CD Integration**: Can be run as part of build process
- **Documentation Updates**: Automatically reflects code changes
- **Version Control**: Generated docs can be committed or ignored
- **Quality Assurance**: Ensures documentation stays current

### Development Benefits
- **Consistency**: Standardized documentation format
- **Completeness**: No missing parameter documentation
- **Accuracy**: Documentation matches actual code
- **Maintenance**: Reduces manual documentation effort

## Build and Run

```bash
# Install dependencies
go mod tidy

# Build
go build -o gen-docs

# Run with default settings
./gen-docs

# Run with custom output directory
./gen-docs output-directory

# Run with custom base directory
./gen-docs docs ../../
```

## Output Example

When run, the generator produces output like:

```
Gollama.cpp Documentation Generator v1.0.0-llamacpp.b6076
Scanning examples in: ../..
Output directory: docs

Scanning examples... found 7 examples
Generating comprehensive documentation... done
Generating per-example documentation... done
Generating usage summary... done

Documentation Generation Complete!
Generated documentation for 7 examples:
  - batched (8 parameters)
  - diffusion (12 parameters)
  - embedding (8 parameters)
  - eval-callback (10 parameters)
  - gen-docs (0 parameters)
  - retrieval (10 parameters)
  - simple-chat (5 parameters)

Output files:
  - docs/examples-reference.md (comprehensive reference)
  - docs/quick-usage.md (quick usage guide)  
  - docs/examples/ (individual example docs)
```

## Use Cases

### Documentation Maintenance
- **Regular Updates**: Run after adding new examples or flags
- **Release Documentation**: Generate docs for releases
- **API Documentation**: Include in API documentation builds
- **Website Integration**: Use generated docs for project websites

### Development Support
- **New Developer Onboarding**: Comprehensive example documentation
- **Parameter Discovery**: Find available options across examples
- **Usage Patterns**: Understand common parameter combinations
- **Feature Exploration**: Discover example capabilities

### Quality Assurance
- **Documentation Completeness**: Ensure all parameters are documented
- **Consistency Checking**: Verify consistent parameter descriptions
- **Change Detection**: Identify when documentation needs updates
- **Standard Compliance**: Enforce documentation standards

## Extending the Generator

### Adding New Flag Types
Support for additional flag types can be added by extending the pattern matching:

```go
flagPatterns["Duration"] = regexp.MustCompile(`flag\.Duration\("([^"]+)",\s*([^,]+),\s*"([^"]*)"\)`)
```

### Custom Categories
New parameter categories can be added to the categorization logic:

```go
customFlags := []string{"custom", "special", "advanced"}
```

### Output Formats
Additional output formats can be implemented:
- **JSON**: Machine-readable parameter definitions
- **HTML**: Web-friendly documentation
- **PDF**: Printable documentation
- **Man Pages**: Unix manual page format

## Related Examples

This example complements other examples in the suite:
- **Simple Chat**: Basic text generation with standard parameters
- **Embedding**: Embedding generation with specialized parameters
- **Retrieval**: Information retrieval with search parameters
- **Eval Callback**: Debugging parameters and monitoring options

## Performance Considerations

### Parsing Performance
- **Regex Efficiency**: Optimized patterns for fast matching
- **File Reading**: Efficient file processing for large codebases
- **Memory Usage**: Minimal memory footprint during analysis
- **Scalability**: Handles projects with many examples

### Output Generation
- **Template Efficiency**: Fast markdown generation
- **File I/O**: Efficient file writing for multiple outputs
- **String Processing**: Optimized text manipulation
- **Batch Processing**: Processes multiple examples efficiently

## Troubleshooting

### Common Issues

1. **No Examples Found**: Check base directory path
2. **Missing Parameters**: Verify flag definition patterns
3. **Parse Errors**: Check for unsupported flag syntax
4. **Output Errors**: Ensure write permissions for output directory

### Debugging Tips

- Use verbose output to see scanning progress
- Check example directory structure
- Verify main.go files exist in example directories
- Ensure proper flag definition syntax

## Future Enhancements

### Planned Features
- **Template System**: Customizable documentation templates
- **Plugin Architecture**: Extensible analysis plugins
- **Configuration Files**: YAML/JSON configuration support
- **Interactive Mode**: Command-line interface for customization

### Integration Possibilities
- **IDE Integration**: Plugin for development environments
- **Web Interface**: Browser-based documentation generator
- **API Documentation**: Integration with API doc generators
- **Testing Integration**: Automated documentation testing
