# CHANGELOG.md Management Script

## Overview

The `scripts/update-changelog.sh` script automates the management of the `CHANGELOG.md` file during the release process. It is integrated into the `tag-release` target in the Makefile to ensure consistent changelog updates.

## Features

- **Release Action**: Converts `[Unreleased]` section to a versioned entry with current date
- **Unreleased Action**: Adds a new `[Unreleased]` section at the top of the changelog
- **Version Update**: Updates existing version entries with new dates when tags are moved
- **Automatic Verification**: Verifies changes were applied correctly
- **Backup & Restore**: Creates backups and restores on failure

## Usage

### Manual Usage

```bash
# Convert [Unreleased] to versioned entry
bash scripts/update-changelog.sh "1.0.0" "release"

# Add new [Unreleased] section
bash scripts/update-changelog.sh "1.0.1" "unreleased"
```

### Automatic Usage (via Makefile)

The script is automatically called during the `tag-release` process:

```bash
make tag-release
```

## Integration with tag-release

The `tag-release` target in the Makefile now includes these steps:

1. **Before tagging**: Updates CHANGELOG.md to convert `[Unreleased]` to the current version with today's date
2. **After tagging**: Increments the version number for next development cycle
3. **After version increment**: Adds a new `[Unreleased]` section for future changes

## Workflow Example

Starting state:
```markdown
## [Unreleased]
### Added
- New feature X
- New feature Y
```

After `make tag-release` (version 1.0.0):
```markdown
## [Unreleased]

### Added

### Changed

### Fixed

### Removed

## [1.0.0] - 2025-08-06
### Added
- New feature X
- New feature Y
```

## Script Actions

### Release Action
- Searches for `## [Unreleased]` heading
- Replaces it with `## [VERSION] - YYYY-MM-DD`
- If the version already exists, updates the date
- Verifies the change was successful

### Unreleased Action
- Checks if `## [Unreleased]` already exists
- If not, adds a new section with standard subsections:
  - `### Added`
  - `### Changed`
  - `### Fixed`
  - `### Removed`
- Inserts the section before the first versioned entry
- If no versioned entries exist, adds after the header

## Error Handling

- Creates backup files (`.bak`) before making changes
- Restores backups if verification fails
- Exits with error codes on failure
- Provides detailed error messages

## Cross-Platform Compatibility

- Uses portable bash constructs
- Avoids complex regex patterns that differ between platforms
- Uses `awk` for reliable text processing
- Uses `grep -F` for literal string matching to avoid escaping issues

## Testing

Use the test script to verify the workflow:

```bash
bash scripts/test-changelog-workflow.sh
```

This script simulates the entire process without making permanent changes, allowing you to see how the changelog would be updated.

## File Structure

```
scripts/
├── update-changelog.sh          # Main changelog management script
├── test-changelog-workflow.sh   # Test script for the entire workflow
└── increment-version.sh         # Version increment script (existing)
```

## Dependencies

- bash
- awk
- grep
- sed (basic usage)
- date (for timestamp generation)

All dependencies are standard POSIX utilities available on Linux, macOS, and Windows with Git Bash or WSL.
