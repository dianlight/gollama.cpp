# Gollama.cpp Tools Documentation

This document describes the various tools and scripts included with gollama.cpp.

## Scripts Directory (`scripts/`)

### hf.sh - Hugging Face Model Downloader

The `hf.sh` script is a powerful utility for downloading models from Hugging Face. It's automatically copied from the llama.cpp repository to ensure compatibility with the current llama.cpp version.

#### Installation

The script is automatically installed when you run:

```bash
make clone-llamacpp
```

To force an update of the script:

```bash
make update-hf-script
```

#### Usage

##### Basic Usage

```bash
# Download by direct URL
./scripts/hf.sh https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/tinyllama-1.1b-chat-v1.0.Q2_K.gguf

# Download using repository and file name
./scripts/hf.sh --repo TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF --file tinyllama-1.1b-chat-v1.0.Q2_K.gguf

# Download to specific directory
./scripts/hf.sh --repo TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF --file tinyllama-1.1b-chat-v1.0.Q2_K.gguf --outdir models
```

##### Options

- `--url <url>`: Direct URL to the model file
- `--repo <repo>`: Hugging Face repository (format: `owner/repo-name`)
- `--file <filename>`: Specific file to download from the repository
- `--outdir <directory>`: Output directory (default: current directory)
- `-h, --help`: Show help message

#### Integration with Examples

The script is integrated with various example projects:

```bash
# Download models for all examples
make model_download

# Example-specific downloads (from example directories)
cd examples/gritlm
make model_download
```

### Other Scripts

#### check-docs.sh
Validates documentation consistency and formatting.

#### increment-version.sh
Utility for version management and release preparation.

## Command Line Tools (`cmd/`)

### gollama-download

The `gollama-download` tool manages llama.cpp library downloads and caching.

#### Installation

Built automatically with the project:

```bash
make build
```

#### Usage

```bash
# Download libraries for current platform
go run ./cmd/gollama-download -download -version b6089

# Test download functionality
go run ./cmd/gollama-download -test-download

# Clean library cache
go run ./cmd/gollama-download -clean-cache
```

#### Options

- `-download`: Download libraries for the current platform
- `-version <build>`: Specify llama.cpp build version (e.g., "b6089")
- `-test-download`: Test the download functionality
- `-clean-cache`: Remove all cached libraries

## Makefile Targets

### Model Management

```bash
# Download example models using hf.sh
make model_download

# Clean model cache
make clean-libs
```

### Script Management

```bash
# Clone llama.cpp and install/update hf.sh script
make clone-llamacpp

# Force update hf.sh script to latest version
make update-hf-script
```

### Development Tools

```bash
# Install development tools (linters, security scanners)
make install-tools

# Run all quality checks
make check

# Format, vet, lint, and test
make fmt vet lint sec test
```

## Integration Examples

### Using hf.sh in Your Projects

```bash
#!/bin/bash
# Example: Download a model for your application

MODEL_DIR="./models"
HF_SCRIPT="./scripts/hf.sh"

# Ensure the script exists
if [ ! -f "$HF_SCRIPT" ]; then
    echo "Installing hf.sh script..."
    make clone-llamacpp
fi

# Download a specific model
echo "Downloading model..."
$HF_SCRIPT --repo microsoft/DialoGPT-medium --file pytorch_model.bin --outdir $MODEL_DIR

echo "Model downloaded to $MODEL_DIR"
```

### Automated Model Management

The project includes automated model management through Make targets that handle:

1. **Dependency Checking**: Ensures `hf.sh` script is available
2. **Conditional Downloads**: Only downloads if model doesn't exist
3. **Error Handling**: Provides clear error messages and recovery instructions
4. **Cross-Platform Support**: Works on all supported platforms

## Troubleshooting

### Script Not Found

If you get an error about `hf.sh` not being found:

```bash
# Install the script
make clone-llamacpp

# Or force update
make update-hf-script
```

### Permission Issues

Ensure the script has execute permissions:

```bash
chmod +x scripts/hf.sh
```

### Download Failures

Common issues and solutions:

1. **Network connectivity**: Check internet connection
2. **Repository access**: Verify the repository exists and is public
3. **File name**: Ensure the file name is correct (case-sensitive)
4. **Disk space**: Check available disk space for large models

### Library Download Issues

For `gollama-download` issues:

```bash
# Clean cache and retry
make clean-libs
make test-download

# Test with verbose output
go run ./cmd/gollama-download -test-download -v
```

## Contributing

When adding new tools or modifying existing ones:

1. Update this documentation
2. Add appropriate Make targets
3. Include usage examples
4. Test on all supported platforms
5. Update integration tests

For questions or issues with tools, please check the main project documentation or open an issue on GitHub.
