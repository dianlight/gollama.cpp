You are working on Gollama.cpp, a Go binding for llama.cpp. Follow these specific rules for all code changes:

## Automatic Documentation Updates

### Always Update When:
1. **API Changes**: Update README.md examples and Go doc comments
2. **New Features**: Add to CHANGELOG.md and update feature lists in README.md
3. **Platform Changes**: Update supported platforms section and CI matrix
4. **Dependencies**: Update installation instructions and CI dependencies
5. **Examples**: Update corresponding README.md files and ensure demo scripts work
6. **Roadmap Items**: Update ROADMAP.md when implementing planned features or changing priorities

### Documentation Files to Consider:
- `README.md` (main project documentation)
- `CHANGELOG.md` (version history)
- `docs/ROADMAP.md` (development roadmap and future plans)
- `examples/*/README.md` (example-specific documentation)
- `docs/*.md` (technical documentation)
- Go doc comments in source files

### ROADMAP.md Update Triggers:
- **Feature Implementation**: Move items from planned to completed sections
- **New Enhancement Issues**: Add to appropriate roadmap sections
- **Priority Changes**: Reorganize short/medium/long-term goals
- **Timeline Updates**: Adjust target dates based on progress
- **Platform Support**: Update Windows, GPU, or cross-platform roadmap items
- **Performance Work**: Update optimization and benchmarking roadmap sections

## Automatic CI Updates

### Update `.github/workflows/ci.yml` When:
1. **Go Version Changes**: Update `GO_VERSION` and matrix strategy
2. **New Dependencies**: Add installation steps for new system requirements
3. **Platform Support**: Add new OS to test matrix
4. **Library Updates**: Update `LLAMA_CPP_BUILD` version
5. **Build Process Changes**: Update compilation flags or build steps

### CI Configuration Patterns:
```yaml
env:
  GO_VERSION: '1.21'           # Update when minimum Go version changes
  LLAMA_CPP_BUILD: 'b6099'     # Update when llama.cpp version changes

strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]  # Add new platforms here
    go-version: ['1.21', '1.22', '1.24']              # Update supported versions
```

## Code Quality Standards

### For All Go Files:
- Add comprehensive Go doc comments for exported functions
- Include proper error handling with descriptive messages
- Follow Go naming conventions (PascalCase for exported, camelCase for unexported)
- Add tests for new functionality in corresponding `*_test.go` files

### For Platform-Specific Files (`platform_*.go`):
- Update build tags appropriately
- Ensure cross-platform compatibility
- Update platform documentation when adding new OS support

### For Examples:
- Ensure all examples compile and run
- Update `demo.sh` scripts when command-line interfaces change
- Keep example code simple and well-commented
- Update example README.md with any changes

## Version Management

### When Updating Dependencies:
1. Update `go.mod` files (main and examples)
2. Update `LLAMA_CPP_BUILD` in CI configuration
3. Update version references in documentation
4. Add changelog entry with migration notes if breaking

### When Adding Features:
1. Add to CHANGELOG.md under "Unreleased" section
2. Update feature list in main README.md
3. Add or update relevant examples
4. Update API documentation

## Specific File Rules

### `gollama.go` (main API):
- Update README.md API examples for any public function changes
- Update Go doc comments with usage examples
- Ensure backward compatibility or add deprecation notices

### `platform_*.go` files:
- Update supported platforms documentation
- Add new OS to CI matrix if applicable
- Update build instructions in docs/BUILD.md

### `examples/*/main.go`:
- Update corresponding README.md
- Ensure demo.sh script works
- Keep examples simple and focused

### `.github/workflows/ci.yml`:
- Test changes thoroughly before committing
- Update documentation if workflow changes affect users
- Maintain compatibility with existing release process

## Error Handling Patterns

Always use descriptive error messages:
```go
if err != nil {
    return fmt.Errorf("failed to load model %s: %w", modelPath, err)
}
```

## Testing Requirements

- Add unit tests for new functions in `*_test.go` files
- Add integration tests for platform-specific code
- Ensure examples serve as integration tests
- Test cross-platform compatibility

## Priority Order for Updates

1. **Security and Bug Fixes**: Immediate documentation and CI updates
2. **Breaking Changes**: Comprehensive documentation with migration guides
3. **New Features**: Full documentation and example updates
4. **Performance Improvements**: Update benchmarks and performance notes
5. **Documentation Only**: Ensure accuracy and consistency

Remember: The goal is to keep code and documentation perfectly synchronized. Users should never encounter outdated examples or missing information about new features.
