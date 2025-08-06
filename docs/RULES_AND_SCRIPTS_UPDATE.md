# Updated Rules and Scripts Summary

This document summarizes the updates made to rules and scripts to support the reorganized roadmap with the new "Long-term Vision (wait for llama.cpp)" section.

## Updated Files

### 1. `scripts/roadmap-update.sh`

#### New Features Added:
- **Enhanced Section Validation**: Added validation for the new "Long-term Vision (wait for llama.cpp)" section
- **Missing Function Scanner**: New `scan-missing` command to detect code that depends on missing llama.cpp functions
- **Blocked Items Counter**: Enhanced validation to count items that are blocked by missing llama.cpp functions

#### New Commands:
```bash
# Scan for code depending on missing llama.cpp functions
./scripts/roadmap-update.sh scan-missing

# Enhanced validation including new sections
./scripts/roadmap-update.sh validate
```

#### Detection Patterns:
The script now scans for these patterns that indicate missing functionality:
- `not implemented`
- `not yet implemented` 
- `missing.*function`
- `requires.*llama.cpp`
- `// Function doesn't exist`
- `runtime.GOOS != "darwin"`
- `Skip.*non-Darwin`

### 2. `Makefile`

#### New Rules Added:
```makefile
# New rule to scan for missing llama.cpp dependencies
.PHONY: roadmap-scan-missing
roadmap-scan-missing:
	@echo "Scanning for code that depends on missing llama.cpp functions"
	@bash scripts/roadmap-update.sh scan-missing
```

#### Updated Help Documentation:
Added the new `roadmap-scan-missing` command to the help output.

### 3. `.github/workflows/doc-sync-check.yml`

#### New Workflow Step:
Added automatic roadmap validation to the CI pipeline:

```yaml
- name: Validate ROADMAP.md format
  run: |
    if [ -f "docs/ROADMAP.md" ]; then
      echo "Validating ROADMAP.md format and structure..."
      
      # Run the roadmap validation script
      make roadmap-validate
      
      # Check for items that might need to be moved to "wait for llama.cpp" section
      echo ""
      echo "Checking for potential missing llama.cpp dependencies..."
      make roadmap-scan-missing
      
      echo "✅ ROADMAP.md validation completed"
    else
      echo "⚠️  No ROADMAP.md found"
    fi
```

This ensures that:
- Roadmap format is validated on every PR
- Missing llama.cpp dependencies are automatically detected
- Documentation stays in sync with code changes

### 4. Source Code Comments

#### Enhanced TODO Comments:
Updated TODO comments in source code to reference roadmap priorities:

```go
// platform_windows.go
// TODO: Implement proper function registration for Windows - blocks ROADMAP Priority 1 (Windows Runtime Completion)

// config.go  
// TODO: Implement logging configuration once we have the actual logging functions - moved to ROADMAP "wait for llama.cpp" section

// gollama.go
// Performance functions - These may not exist in this llama.cpp version - moved to ROADMAP "wait for llama.cpp" section

// Error messages now include roadmap context
return errors.New("Decode not yet implemented for non-Darwin platforms - blocks ROADMAP Priority 1 (Windows Runtime Completion)")
```

## Usage Examples

### Check Roadmap Status
```bash
# Validate overall roadmap format
make roadmap-validate

# Scan for TODO/FIXME items 
make roadmap-scan

# Scan for missing llama.cpp dependencies
make roadmap-scan-missing

# Update last modified date
make roadmap-update
```

### CI Integration
The updated workflow will automatically:
1. Validate roadmap format on every PR
2. Check for new TODO/FIXME comments
3. Scan for code depending on missing llama.cpp functions
4. Provide recommendations for documentation updates

## Benefits

### 1. **Automated Dependency Tracking**
- Automatically detects when code depends on missing llama.cpp functions
- Helps prioritize roadmap items based on actual code dependencies
- Prevents features from being moved to active development prematurely

### 2. **Improved CI/CD Pipeline** 
- Roadmap validation is now part of the CI process
- Documentation sync is automatically checked
- Reduces manual review overhead

### 3. **Better Code Documentation**
- TODO comments now reference specific roadmap items
- Clear indication of what blocks each feature
- Easier to understand project priorities from code

### 4. **Enhanced Project Management**
- Automatic categorization of features by dependency type
- Clear separation between achievable and blocked goals
- Better communication about what requires upstream changes

## Future Enhancements

These scripts and rules can be extended to:

1. **Automatic Issue Creation**: Generate GitHub issues for items moved to "wait for llama.cpp"
2. **llama.cpp Version Tracking**: Automatically check when blocking functions become available
3. **Progress Metrics**: Generate statistics on roadmap completion rates
4. **Integration Testing**: Automatically test features when dependencies become available

This infrastructure provides a solid foundation for maintaining the roadmap as the project evolves and llama.cpp adds new functionality.
