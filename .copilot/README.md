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

### `templates.md`
Ready-to-use templates for:
- README.md updates
- CHANGELOG.md entries
- CI configuration changes
- Go documentation comments
- Example documentation

## Automatic Behavior

When GitHub Copilot detects changes to:

### API Files (`gollama.go`, `platform_*.go`, etc.)
- Updates README.md examples
- Updates Go doc comments
- Adds CHANGELOG.md entries
- Updates CI if dependencies change

### Example Files (`examples/*/`)
- Updates corresponding README.md files
- Ensures demo scripts work
- Updates main examples documentation

### Dependencies (`go.mod`, `libs/`)
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

1. **During Development**: Copilot suggests documentation updates as you code
2. **Before Committing**: Run `./scripts/check-docs.sh` to verify completeness  
3. **In Pull Requests**: CI checks validate documentation sync
4. **After Merging**: Documentation stays current with code changes

This configuration ensures that documentation never falls behind code changes, improving the developer experience and project quality.
