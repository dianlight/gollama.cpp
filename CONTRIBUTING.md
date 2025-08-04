# Contributing to Gollama.cpp

We welcome contributions to gollama.cpp! This document provides guidelines for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Documentation](#documentation)
- [Release Process](#release-process)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

### Our Standards

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

Before contributing, ensure you have:

- Go 1.21 or later
- Git
- Make
- CMake 3.14+
- Platform-specific build tools (see [BUILD.md](BUILD.md))

### Setting up the Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/gollama.cpp
   cd gollama.cpp
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/dianlight/gollama.cpp
   ```

4. Install dependencies:
   ```bash
   make deps
   make install-tools
   ```

5. Build the project:
   ```bash
   make build
   ```

6. Run tests to ensure everything works:
   ```bash
   make test
   ```

## Development Setup

### Recommended IDE Setup

#### VS Code
Install the following extensions:
- Go (official Go extension)
- C/C++ (for llama.cpp development)
- GitLens
- Code Spell Checker

#### Recommended VS Code settings (`.vscode/settings.json`):
```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "workspace",
    "go.formatTool": "goimports",
    "go.generateTestsFlags": ["-parallel"]
}
```

### Environment Variables

Set up your development environment:

```bash
# Optional: Specify custom llama.cpp build
export LLAMA_CPP_BUILD=b6076

# Optional: Enable verbose testing
export GOLLAMA_TEST_VERBOSE=1

# Optional: Specify test model path
export GOLLAMA_TEST_MODEL=/path/to/test/model.gguf
```

## Making Changes

### Platform-Specific Development

Gollama.cpp uses a **platform-specific architecture** with Go build tags. When contributing:

#### Build Tags and Platform Support

We use the following build tag strategy:
- **`!windows`**: Unix-like systems (Linux, macOS) using purego
- **`windows`**: Windows systems using native syscalls

#### Platform-Specific Files

When working on platform support:

1. **Unix-like platforms** (`platform_unix.go`):
   ```go
   //go:build !windows
   
   package gollama
   
   import "github.com/ebitengine/purego"
   
   func loadLibraryPlatform(libPath string) (uintptr, error) {
       return purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
   }
   ```

2. **Windows platforms** (`platform_windows.go`):
   ```go
   //go:build windows
   
   package gollama
   
   import "syscall"
   
   func loadLibraryPlatform(libPath string) (uintptr, error) {
       // Windows-specific implementation using LoadLibraryW
   }
   ```

#### Testing Platform Changes

Always test platform-specific code:

```bash
# Test current platform
go test -v -run TestPlatformSpecific ./...

# Test cross-compilation for all platforms
make test-cross-compile

# Test specific platform (without running)
GOOS=windows GOARCH=amd64 go test -c ./...
GOOS=linux GOARCH=arm64 go test -c ./...
GOOS=darwin GOARCH=arm64 go test -c ./...
```

#### Windows Development Guidelines

When contributing Windows support:

1. Use native Windows APIs via `syscall` package
2. Implement proper error handling for Windows-specific errors
3. Test both compilation and runtime on Windows when possible
4. Ensure cross-compilation works from other platforms

### Branch Naming

Use descriptive branch names:
- `feature/add-sampling-method` - for new features
- `fix/memory-leak-context` - for bug fixes
- `docs/update-api-reference` - for documentation
- `refactor/simplify-tokenization` - for refactoring

### Workflow

1. Create a new branch from `main`:
   ```bash
   git checkout main
   git pull upstream main
   git checkout -b feature/your-feature-name
   ```

2. Make your changes in logical commits
3. Write or update tests
4. Update documentation if needed
5. Ensure all tests pass:
   ```bash
   make check
   ```

### Commit Messages

Follow the conventional commit format:

```
type(scope): short description

Longer description if needed.

Fixes #123
```

Types:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
- `feat(sampling): add temperature sampling method`
- `fix(memory): resolve context memory leak`
- `docs(api): update tokenization examples`

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Run benchmarks
make bench

# Run tests for specific package
go test -v ./pkg/sampling/
```

### Writing Tests

#### Unit Tests
- Test files should end with `_test.go`
- Test functions should start with `Test`
- Use table-driven tests for multiple test cases

Example:
```go
func TestTokenize(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected []LlamaToken
        wantErr  bool
    }{
        {
            name:     "simple text",
            input:    "hello world",
            expected: []LlamaToken{123, 456},
            wantErr:  false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Tokenize(model, tt.input, true, false)
            if (err != nil) != tt.wantErr {
                t.Errorf("Tokenize() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("Tokenize() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

#### Benchmarks
```go
func BenchmarkTokenize(b *testing.B) {
    text := "The quick brown fox jumps over the lazy dog"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := Tokenize(model, text, true, false)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

#### Integration Tests
- Place in `tests/` directory
- Require actual model files
- May be skipped in CI if models not available

### Test Requirements

- All new code must have tests
- Tests must pass on all supported platforms
- Code coverage should not decrease significantly
- Benchmarks should not show performance regressions

## Submitting Changes

### Pull Request Process

1. Ensure your branch is up to date:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. Push your branch:
   ```bash
   git push origin feature/your-feature-name
   ```

3. Create a pull request on GitHub

4. Fill out the pull request template completely

5. Respond to code review feedback

### Pull Request Requirements

- [ ] Code follows project style guidelines
- [ ] Tests are included and passing
- [ ] Documentation is updated (if applicable)
- [ ] CHANGELOG.md is updated (for significant changes)
- [ ] No breaking changes (unless approved)
- [ ] Performance impact is assessed

### Review Process

1. Automated checks must pass (CI/CD)
2. At least one maintainer review required
3. Address all feedback before merging
4. Maintainer will merge when ready

## Code Style

### Go Code Style

Follow standard Go conventions:

- Use `gofmt` for formatting
- Use `goimports` for import organization
- Follow effective Go guidelines
- Use meaningful variable and function names
- Add comments for exported functions and types

### Linting

We use `golangci-lint` with the following configuration:

```bash
# Run linter
make lint

# Fix auto-fixable issues
golangci-lint run --fix
```

### Documentation Comments

```go
// Tokenize converts text into a sequence of tokens using the specified model.
// The addSpecial parameter determines whether to add special tokens (BOS/EOS).
// The parseSpecial parameter determines whether to parse special token sequences.
//
// Returns a slice of tokens and an error if tokenization fails.
func Tokenize(model LlamaModel, text string, addSpecial, parseSpecial bool) ([]LlamaToken, error) {
    // Implementation...
}
```

## Documentation

### API Documentation

- All exported functions must have documentation comments
- Include examples in documentation when helpful
- Update README.md for significant API changes

### User Documentation

- Update relevant documentation in `docs/`
- Include examples for new features
- Update CHANGELOG.md for user-visible changes

### Example Code

- Include working examples for new features
- Test all example code
- Keep examples simple and focused

## Release Process

### Version Numbering

We follow semantic versioning with llama.cpp build numbers:

- Format: `vX.Y.Z-llamacpp.BUILD`
- Example: `v1.2.0-llamacpp.b6076`

### Release Checklist

For maintainers:

1. Update CHANGELOG.md
2. Update version constants in code
3. Create and push git tag
4. GitHub Actions will build and release automatically
5. Update documentation if needed
6. Announce release

## Types of Contributions

### Bug Reports

Use the issue template and include:
- Go version
- Operating system and architecture
- Steps to reproduce
- Expected vs actual behavior
- Minimal code example

### Feature Requests

- Describe the use case
- Explain why the feature is needed
- Provide examples of how it would be used
- Consider implementation complexity

### Code Contributions

- New features
- Bug fixes
- Performance improvements
- Documentation improvements
- Test improvements

### Documentation Contributions

- API documentation
- Tutorials and guides
- Example code
- README improvements

## Platform-Specific Contributions

### Cross-Platform Testing

When contributing:
- Test on multiple platforms if possible
- Note any platform-specific behavior
- Update platform-specific documentation

### GPU Support

For GPU-related contributions:
- Test on relevant hardware when possible
- Document hardware requirements
- Include fallback behavior for unsupported systems

## Getting Help

- Create an issue for bugs or feature requests
- Start a discussion for questions
- Check existing issues and discussions first
- Be patient and respectful when asking for help

## Recognition

Contributors will be:
- Listed in CHANGELOG.md for significant contributions
- Credited in release notes
- Mentioned in project documentation

Thank you for contributing to gollama.cpp!
