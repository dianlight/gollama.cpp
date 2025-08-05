# Copilot Validation Rules

## Pre-commit Validation
Before any code changes are finalized, run these validation commands:

1. **Lint Check**: `make lint`
   - Validates Go code formatting
   - Checks for style violations
   - Ensures code follows project conventions

2. **Security Check**: `make sec`
   - Runs security analysis tools
   - Identifies potential vulnerabilities
   - Validates secure coding practices

3. **Combined Check**: `make lint sec`
   - Runs both lint and security checks
   - Available as VS Code task: "Validate Changes (lint + sec)"

## Validation Workflow

### For any Go code changes:
```bash
# Before completing changes
make lint sec

# Or use the VS Code task
# Terminal > Run Task > "Validate Changes (lint + sec)"
```

### Expected Results:
- All lint checks must pass (no violations)
- All security checks must pass (no vulnerabilities)
- Fix any issues before considering changes complete

## Integration Points

### VS Code Task
A pre-configured task is available:
- **Task Name**: "Validate Changes (lint + sec)"
- **Command**: `make lint sec`
- **Usage**: Terminal > Run Task > Select task

### GitHub Copilot Rules
This validation is automatically enforced by Copilot rules defined in `.copilot/rules.md`:
- Code changes trigger automatic validation
- Lint and security issues are addressed before completion
- Ensures consistent code quality across the project

## Tools Used

### Lint Tools
- `go fmt` - Code formatting
- `go vet` - Static analysis
- `golangci-lint` - Comprehensive linting (if configured)

### Security Tools  
- `gosec` - Security analysis for Go code
- Static security analysis
- Vulnerability detection

## Troubleshooting

### Common Lint Issues
- Run `go fmt` to fix formatting
- Check for unused variables/imports
- Follow Go naming conventions

### Common Security Issues
- Review error handling patterns
- Check for hardcoded credentials
- Validate input sanitization
- Review file permissions and paths

### Getting Help
- Check `make help` for available targets
- Review tool documentation for specific errors
- Consult project documentation in `docs/` directory
