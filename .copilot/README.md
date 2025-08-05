# Copilot Configuration

This directory contains configuration files that help GitHub Copilot automatically maintain documentation and CI configuration in sync with code changes.

## Files

### `instructions.md`
Primary instructions for GitHub Copilot. This file contains:
- Automatic documentation update rules
- CI configuration update guidelines  
- Code quality standards
- File-specific update rules

### `rules.md`
Comprehensive rules document explaining:
- When to update different types of documentation
- Priority guidelines for different types of changes
- Platform-specific update requirements
- **Code validation rules** for automatic lint and security checks

### `templates.md`
Ready-to-use templates for:
- README.md updates
- CHANGELOG.md entries
- CI configuration changes
- Go documentation comments
- Example documentation

### `validation.md`
Code validation guidelines including:
- Pre-commit validation workflow
- Lint and security check requirements
- VS Code task integration
- Troubleshooting guidance

## Automatic Behavior

When GitHub Copilot detects changes to:

### API Files (`gollama.go`, `platform_*.go`, etc.)
- **Runs automatic validation** with `make lint sec`
- Updates README.md examples
- Updates Go doc comments
- Adds CHANGELOG.md entries
- Updates CI if dependencies change

### Example Files (`examples/*/`)
- **Validates code changes** before completion
- Updates corresponding README.md files
- Ensures demo scripts work
- Updates main examples documentation

### Dependencies (`go.mod`, `libs/`)
- **Runs security checks** on new dependencies
- Updates CI configuration
- Updates installation instructions
- Updates version references

### Platform Support
- Updates CI matrix
- Updates supported platforms documentation
- Updates build instructions

## Manual Tools

### Documentation Check Script
Run locally before committing:
```bash
./scripts/check-docs.sh
```

This script:
- Analyzes your changes
- Suggests documentation updates
- Tests example compilation
- Checks for TODOs and formatting issues

### Code Validation
Run validation before committing:
```bash
make lint sec
# Or use VS Code task: "Validate Changes (lint + sec)"
```

This validation:
- Checks code formatting and style
- Performs security analysis
- Ensures code quality standards
- Required for all code changes

### CI Workflow
The `doc-sync-check.yml` workflow runs on PRs to:
- Detect when documentation might be out of sync
- Validate examples still compile
- Check CHANGELOG.md format
- Suggest improvements

## Usage Tips

1. **Let Copilot Help**: When making code changes, Copilot will automatically suggest documentation updates based on these rules.

2. **Review Suggestions**: Always review Copilot's documentation suggestions to ensure accuracy.

3. **Use Templates**: Reference `templates.md` for consistent formatting.

4. **Run Checks**: Use `./scripts/check-docs.sh` before committing to catch issues early.

5. **Update Rules**: Modify these files as the project evolves to keep Copilot's suggestions relevant.

## Integration with Development Workflow

1. **During Development**: 
   - Copilot suggests documentation updates as you code
   - **Automatic validation** runs for code changes (lint + security)
2. **Before Committing**: 
   - Run `./scripts/check-docs.sh` to verify completeness
   - **Run `make lint sec`** to validate code quality  
3. **In Pull Requests**: CI checks validate documentation sync and code quality
4. **After Merging**: Documentation stays current with code changes

This configuration ensures that documentation never falls behind code changes and that all code meets quality and security standards, improving the developer experience and project quality.
