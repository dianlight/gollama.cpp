# Copilot Rules for Gollama.cpp

> **Note**: This file has been migrated to `.github/copilot-instructions.md` (the official standard location).
> Please refer to that file for the current rules. This file is kept for reference only.

---

**👉 See [.github/copilot-instructions.md](../.github/copilot-instructions.md) for the active Copilot rules.**

---

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

## Git Operations

### Rule: NO Automatic Git Operations
**CRITICAL**: Never perform git operations automatically without explicit user request:

1. **Prohibited Actions**:
   - Do NOT run `git add` without explicit user request
   - Do NOT run `git commit` without explicit user request
   - Do NOT run `git push` without explicit user request
   - Do NOT run `git pull` without explicit user request
   - Do NOT create branches without explicit user request
   - Do NOT merge branches without explicit user request

2. **When to Ask**:
   - Always ask the user if they want to commit changes
   - Always ask the user if they want to push changes
   - Always ask the user which files should be staged
   - Always ask the user for commit messages

3. **Allowed Git Operations**:
   - Read-only operations like `git status` or `git diff` are acceptable
   - Informing the user about uncommitted changes is acceptable
   - Suggesting git commands for the user to run manually is acceptable

## Code Quality Rules

### Rule: Automatic Code Validation
Before completing any code changes, always run validation tools:

1. **Lint Validation**:
   - Run `make lint` to check code formatting and style
   - Fix any linting issues before submitting changes
   - Ensure code follows Go best practices and project conventions

2. **Security Validation**:
   - Run `make sec` to perform security analysis
   - Address any security vulnerabilities or warnings
   - Verify that new code doesn't introduce security risks

3. **Combined Validation**:
   - Use the available VS Code task "Validate Changes (lint + sec)" to run both checks
   - Alternatively run `make lint sec` to execute both validations
   - Ensure all validation passes before considering code changes complete

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

1. **Code changes in Go files**:
   - Automatically run `make lint` to validate code style
   - Automatically run `make sec` to check for security issues
   - Fix any linting or security issues before completing changes

2. **New Go module dependencies**: 
   - Check if CI needs updated system dependencies
   - Update README.md installation instructions if needed
   - Run validation tools to ensure new dependencies don't introduce issues

3. **API signature changes**:
   - Update all example code
   - Update documentation with new signatures
   - Add deprecation notices if needed
   - Validate changes with lint and security tools

4. **New platform support**:
   - Add platform to CI matrix
   - Update README.md supported platforms section
   - Update build documentation
   - Run validation to ensure cross-platform compatibility

5. **Version bumps**:
   - Update CHANGELOG.md
   - Update version references in documentation
   - Update CI configuration if needed
   - Validate all changes before completing version update

## Roadmap Management

### Rule: Auto-update ROADMAP.md
When implementing features or making significant changes that affect the project roadmap:

1. **Feature Completion**:
   - Move completed items from "In Progress" to "Completed Major Features"
   - Update progress indicators (✅, 🚧, [ ]) appropriately
   - Add completion dates for major milestones
   - Update the "Last Updated" date at the bottom

2. **New Feature Planning**:
   - Add new planned features to appropriate priority sections
   - Update timeline estimates based on current progress
   - Adjust dependencies and technical requirements
   - Update success metrics if applicable

3. **Priority Changes**:
   - Reassess priorities based on user feedback and technical constraints
   - Move items between Short-term, Medium-term, and Long-term sections
   - Update implementation priorities (High/Medium/Low)
   - Adjust target dates accordingly

4. **Technical Dependencies**:
   - Update external dependencies when new requirements are discovered
   - Modify internal dependencies based on architectural changes
   - Update risk assessments for changing technical landscape
   - Revise contribution opportunities

### Roadmap Triggers for Updates:
- **Platform Support Changes**: Update platform-specific roadmap items
- **GPU Backend Updates**: Modify GPU support roadmap sections
- **API Changes**: Update advanced features and developer experience items
- **Performance Improvements**: Adjust performance optimization roadmap
- **Community Growth**: Update contribution and success metrics
- **External Dependencies**: Modify dependency-related roadmap items

## Priority Guidelines

1. **High Priority**: Security fixes, breaking changes, new platform support, roadmap-critical items
2. **Medium Priority**: Feature additions, performance improvements, roadmap enhancements
3. **Low Priority**: Documentation improvements, example updates, roadmap refinements

Always prioritize keeping documentation in sync with code changes to maintain project quality and user experience. The roadmap should reflect the current state and realistic future plans.
