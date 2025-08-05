# Copilot Rules for Gollama.cpp

This file defines automatic rules for GitHub Copilot to follow when making changes to the codebase.

## Documentation Updates

### Rule: Auto-update Documentation
When making changes to code, always consider updating relevant documentation files:

1. **README.md Updates**: 
   - Update feature lists when adding/removing functionality
   - Update code examples when APIs change
   - Update supported platforms when adding new platform support
   - Update version compatibility information
   - Update installation instructions if dependencies change

2. **CHANGELOG.md Updates**:
   - Add entries for new features, bug fixes, and breaking changes
   - Follow semantic versioning principles
   - Include migration guides for breaking changes
   - Reference issue/PR numbers

3. **API Documentation**:
   - Update Go doc comments when function signatures change
   - Update example code in documentation
   - Update any inline documentation

4. **Example Documentation**:
   - Update `examples/*/README.md` when example code changes
   - Ensure all examples compile and run correctly
   - Update demo scripts if command-line interfaces change

## CI/CD Updates

### Rule: Auto-update CI Configuration
When making changes that affect the build process, testing, or deployment:

1. **Go Version Updates**:
   - Update `GO_VERSION` in `.github/workflows/ci.yml`
   - Update go.mod files if minimum Go version changes
   - Update matrix strategy in CI if supporting new Go versions

2. **Dependency Changes**:
   - Update CI dependencies when adding new system requirements
   - Update cache keys when dependency structure changes
   - Add new test steps for new functionality

3. **Platform Support Changes**:
   - Update CI matrix when adding/removing platform support
   - Add new OS runners when extending platform compatibility
   - Update build tags and compilation flags

4. **Library Updates**:
   - Update `LLAMA_CPP_BUILD` version when upgrading llama.cpp
   - Update library paths and download URLs
   - Update build scripts and Makefiles

## Code Quality Rules

### Rule: Maintain Code Standards
When writing or modifying code:

1. **Go Standards**:
   - Follow Go naming conventions
   - Add proper error handling
   - Include comprehensive tests
   - Add Go doc comments for exported functions

2. **Test Coverage**:
   - Add tests for new functionality
   - Update existing tests when APIs change
   - Ensure examples have corresponding tests

3. **Version Consistency**:
   - Keep version numbers synchronized across files
   - Update version references in documentation
   - Update download URLs and checksums

## File-Specific Rules

### For `gollama.go` changes:
- Update main README.md with API changes
- Update examples if public API changes
- Update CI tests if new dependencies are added

### For `platform_*.go` changes:
- Update platform support documentation in README.md
- Update CI matrix if new platforms are added
- Update build documentation

### For `examples/` changes:
- Update corresponding README.md in the example directory
- Update main examples documentation
- Ensure demo.sh scripts work correctly

### For `libs/` or dependency changes:
- Update CI.yml with new library versions
- Update download scripts
- Update platform-specific documentation

### For `.github/workflows/` changes:
- Update CI documentation if workflow changes affect users
- Test changes thoroughly as they affect release process
- Update any references to workflow names or steps

## Automatic Actions

When Copilot detects:

1. **New Go module dependencies**: 
   - Check if CI needs updated system dependencies
   - Update README.md installation instructions if needed

2. **API signature changes**:
   - Update all example code
   - Update documentation with new signatures
   - Add deprecation notices if needed

3. **New platform support**:
   - Add platform to CI matrix
   - Update README.md supported platforms section
   - Update build documentation

4. **Version bumps**:
   - Update CHANGELOG.md
   - Update version references in documentation
   - Update CI configuration if needed

## Priority Guidelines

1. **High Priority**: Security fixes, breaking changes, new platform support
2. **Medium Priority**: Feature additions, performance improvements
3. **Low Priority**: Documentation improvements, example updates

Always prioritize keeping documentation in sync with code changes to maintain project quality and user experience.
