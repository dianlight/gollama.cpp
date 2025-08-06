# Documentation Update Templates

Use these templates when automatically updating documentation files.

## README.md Update Templates

### Adding a New Feature
```markdown
## Features

- **Pure Go**: No CGO required, uses purego for C interop
- **Cross-Platform**: Supports macOS (CPU/Metal), Linux (CPU/NVIDIA/AMD), Windows (CPU/NVIDIA/AMD)
- **[NEW FEATURE NAME]**: [Brief description of what it does]
```

### API Change Example
```markdown
## Quick Start

```go
package main

import (
    "fmt"
    "github.com/dianlight/gollama.cpp"
)

func main() {
    // [Update this section when API changes]
    model, err := gollama.LoadModel("path/to/model.gguf")
    if err != nil {
        panic(err)
    }
    defer model.Close()
    
    // [Keep examples current with actual API]
}
```

### Platform Support Update
```markdown
### ✅ Fully Supported Platforms

#### [Platform Name]
- **CPU**: [Architectures]
- **GPU**: [GPU Types]
- **Status**: [Current status]
- **Build Tags**: [Relevant build tags]
```

## CHANGELOG.md Update Templates

### New Release Section
```markdown
## [Unreleased]

### Added
- [New feature description]

### Changed
- [Changed functionality]

### Deprecated
- [Deprecated features with migration path]

### Removed
- [Removed features]

### Fixed
- [Bug fixes]

### Security
- [Security improvements]
```

### Breaking Change Entry
```markdown
### Changed
- **BREAKING**: [Description of breaking change]
  - **Migration**: [How to update existing code]
  - **Reason**: [Why the change was necessary]
  - **Example**: 
    ```go
    // Old way (deprecated)
    model.OldMethod()
    
    // New way
    model.NewMethod()
    ```
```

## CI Configuration Update Templates

### Adding New Go Version
```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]
    go-version: ['1.21', '1.22', '1.24', '1.25']  # Add new version here
```

### Adding System Dependencies
```yaml
- name: Install dependencies (Ubuntu)
  if: matrix.os == 'ubuntu-latest'
  run: |
    sudo apt-get update
    sudo apt-get install -y cmake build-essential [new-dependency]
```

### Updating llama.cpp Version
```yaml
env:
  GO_VERSION: '1.21'
  LLAMA_CPP_BUILD: 'b6099'  # Update this when upgrading llama.cpp
```

### Adding New Platform
```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest, [new-platform]]
    include:
      - os: [new-platform]
        go-version: '1.24'  # Specify if different Go version needed
```

## Example README.md Update Template

```markdown
# [Example Name]

[Brief description of what this example demonstrates]

## Prerequisites

- Go 1.21 or later
- [Any specific requirements]

## Building

```bash
make build
# or
go build -o [binary-name] main.go
```

## Running

```bash
./demo.sh
# or
./[binary-name] [arguments]
```

## Expected Output

```
[Sample output that users should expect]
```

## Notes

- [Any important notes about the example]
- [Performance considerations]
- [Platform-specific behavior]
```

## Go Doc Comment Templates

### Function Documentation
```go
// [FunctionName] [brief description of what it does].
//
// [Detailed description if needed, including:]
// - Parameter explanations
// - Return value descriptions  
// - Usage examples
// - Error conditions
//
// Example:
//
//	result, err := FunctionName(param1, param2)
//	if err != nil {
//		// handle error
//	}
//	// use result
//
// Returns an error if [condition that causes error].
func FunctionName(param1 Type1, param2 Type2) (ReturnType, error) {
```

### Type Documentation
```go
// [TypeName] represents [what this type is for].
//
// [Detailed description including:]
// - What fields are used for
// - Initialization requirements
// - Thread safety considerations
//
// Example usage:
//
//	config := &TypeName{
//		Field1: value1,
//		Field2: value2,
//	}
type TypeName struct {
    // Field1 [description of what this field does]
    Field1 Type1
    
    // Field2 [description of what this field does]
    Field2 Type2
}
```

## Automated Update Checklist

When making changes, ensure these files are updated:

### Code Changes:
- [ ] Go doc comments updated
- [ ] Examples updated if API changed
- [ ] Tests added/updated
- [ ] Error messages are descriptive

### Documentation:
- [ ] README.md (if public API changed)
- [ ] CHANGELOG.md (with proper section)
- [ ] Example README files (if examples changed)
- [ ] Migration guides (for breaking changes)

### CI/CD:
- [ ] Go version matrix (if minimum version changed)
- [ ] System dependencies (if new deps added)
- [ ] Build environment (if compilation changed)
- [ ] Platform matrix (if new platforms added)

### Version Management:
- [ ] go.mod files updated
- [ ] Version references in docs
- [ ] Library version numbers
- [ ] Download URLs and checksums
