# Renovate Configuration for llama.cpp Updates

This document describes how Renovate is configured to automatically track and update llama.cpp dependencies in this project.

## Overview

The project uses [Renovate](https://docs.renovatebot.com/) to automatically track updates to the llama.cpp library and manage version bumping. When a new llama.cpp release is detected, Renovate will:

1. Update `LLAMA_CPP_BUILD` references across multiple files
2. Automatically increment the project's minor version
3. Update the CHANGELOG.md with the dependency changes
4. Create a pull request with all changes

## Files Monitored

Renovate monitors the following files for llama.cpp version references:

- **Makefile**: `LLAMA_CPP_BUILD ?= b6862`
- **gollama.go**: `LlamaCppBuild = "b6862"`
- **.github/workflows/ci.yml**: `LLAMA_CPP_BUILD: 'b6862'`
- **CONTRIBUTING.md**: `export LLAMA_CPP_BUILD=b6862`
- **docs/BUILD.md**: `export LLAMA_CPP_BUILD=b6862`

## Version Bumping

When llama.cpp is updated, the project version is automatically incremented using the following logic:

- **Current version**: `0.2.0`
- **After llama.cpp update**: `0.3.0` (minor version incremented, patch reset to 0)

The version is updated in both:
- `Makefile`: `VERSION ?= 0.2.0`
- `gollama.go`: `Version = "0.2.0"`
- `CHANGELOG.md`: Dependency update information

## Configuration Details

### Custom Managers

Renovate uses custom regex managers to detect llama.cpp version references:

1. **GitHub Workflows**: Matches `LLAMA_CPP_BUILD: 'b6862'` patterns
2. **Makefile**: Matches `LLAMA_CPP_BUILD ?= b6862` patterns  
3. **Go Source**: Matches `LlamaCppBuild = "b6862"` patterns
4. **Documentation**: Matches `export LLAMA_CPP_BUILD=b6862` patterns

### Post-Upgrade Tasks

After updating llama.cpp references, Renovate runs:

```bash
chmod +x scripts/increment-version.sh
./scripts/increment-version.sh minor
```

This automatically increments the project version in both `Makefile` and `gollama.go`.

### Package Rules

llama.cpp updates are grouped together with the following configuration:

- **Group Name**: "llama.cpp"
- **Labels**: `llama.cpp`, `dependencies`
- **Schedule**: Before 4am on Monday
- **Minimum Release Age**: 1 day
- **Commit Type**: `feat(deps)`
- **PR Title**: `feat(deps): update llama.cpp to {{newVersion}}`

## Manual Version Management

For manual version management, use the provided script:

```bash
# Increment patch version (0.2.0 -> 0.2.1)
./scripts/increment-version.sh patch

# Increment minor version (0.2.0 -> 0.3.0)  
./scripts/increment-version.sh minor

# Increment major version (0.2.0 -> 1.0.0)
./scripts/increment-version.sh major
```

The script updates both `Makefile` and `gollama.go` consistently.

**Note**: When using manual version management, remember to also update `CHANGELOG.md` with the appropriate version entry and change information.

## Backup Workflow

A GitHub Action workflow (`auto-version-bump.yml`) provides a backup mechanism that detects Renovate PRs for llama.cpp and automatically increments the version if the post-upgrade task fails.

## Testing

To test the Renovate configuration locally:

1. **Test version increment script**:
   ```bash
   ./scripts/increment-version.sh patch
   git checkout -- Makefile gollama.go  # Reset changes
   ```

2. **Check current versions**:
   ```bash
   grep "^VERSION" Makefile
   grep -E "^\s*Version = " gollama.go
   ```

3. **Validate Renovate config**:
   ```bash
   npx renovate-config-validator .github/renovate.json
   ```

## Troubleshooting

### Common Issues

1. **Version format mismatch**: Ensure versions follow `MAJOR.MINOR.PATCH` format
2. **Regex not matching**: Check that file patterns exactly match the expected format
3. **Post-upgrade tasks failing**: Verify script has execute permissions and handles file updates correctly

### Manual Intervention

If automatic updates fail, manually update:

1. Update all `LLAMA_CPP_BUILD` references to the new llama.cpp version
2. Run `./scripts/increment-version.sh minor` to bump the project version
3. Update `CHANGELOG.md` with documenting the llama.cpp update
4. Commit changes with format: `feat(deps): update llama.cpp to <version>`

## Configuration Files

- **Main Config**: `.github/renovate.json`
- **Version Script**: `scripts/increment-version.sh`
- **Backup Workflow**: `.github/workflows/auto-version-bump.yml`
