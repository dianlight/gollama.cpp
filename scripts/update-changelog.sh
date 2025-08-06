#!/bin/bash

# Script to update CHANGELOG.md during release process
# Usage: ./scripts/update-changelog.sh [version] [action]
# Actions: release, unreleased
# - release: converts [Unreleased] to versioned entry with current date
# - unreleased: adds new [Unreleased] section at top

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CHANGELOG_FILE="${CHANGELOG_FILE:-$PROJECT_ROOT/CHANGELOG.md}"

# Parse arguments
VERSION="$1"
ACTION="$2"

if [[ -z "$VERSION" ]] || [[ -z "$ACTION" ]]; then
    echo "Usage: $0 [version] [action]"
    echo "Version format: v0.2.0-llamacpp.b6099 or 0.2.0-llamacpp.b6099"
    echo "Actions:"
    echo "  release    - Convert [Unreleased] to versioned entry"
    echo "  unreleased - Add new [Unreleased] section"
    exit 1
fi

if [[ ! -f "$CHANGELOG_FILE" ]]; then
    echo "Error: CHANGELOG.md not found at $CHANGELOG_FILE"
    exit 1
fi

# Get current date in ISO format
CURRENT_DATE=$(date '+%Y-%m-%d')

# Backup original file
cp "$CHANGELOG_FILE" "$CHANGELOG_FILE.bak"

case "$ACTION" in
    release)
        echo "Converting [Unreleased] to [$VERSION] - $CURRENT_DATE in CHANGELOG.md"
        
        # Check if [Unreleased] section exists
        if ! grep -q "## \[Unreleased\]" "$CHANGELOG_FILE"; then
            echo "Warning: [Unreleased] section not found in CHANGELOG.md"
            rm -f "$CHANGELOG_FILE.bak"
            exit 0
        fi
        
        # Use a simpler approach with awk to replace the content
        awk -v version="$VERSION" -v date="$CURRENT_DATE" '
        /^## \[Unreleased\]$/ {
            print "## [" version "] - " date
            next
        }
        /^## \['"$VERSION"'\] - / {
            print "## [" version "] - " date
            next
        }
        { print }
        ' "$CHANGELOG_FILE" > "$CHANGELOG_FILE.tmp"
        
        mv "$CHANGELOG_FILE.tmp" "$CHANGELOG_FILE"
        echo "Successfully updated CHANGELOG.md for release $VERSION"
        ;;
        
    unreleased)
        echo "Adding new [Unreleased] section to CHANGELOG.md"
        
        # Check if [Unreleased] section already exists
        if grep -q "## \[Unreleased\]" "$CHANGELOG_FILE"; then
            echo "[Unreleased] section already exists in CHANGELOG.md"
            rm -f "$CHANGELOG_FILE.bak"
            exit 0
        fi
        
        # Use awk to insert [Unreleased] section
        awk '
        BEGIN { unreleased_added = 0 }
        /^## \[.*\] - [0-9]{4}-[0-9]{2}-[0-9]{2}$/ && !unreleased_added {
            print "## [Unreleased]"
            print ""
            print "### Added"
            print ""
            print "### Changed"
            print ""
            print "### Fixed"
            print ""
            print "### Removed"
            print ""
            unreleased_added = 1
        }
        /^and this project adheres to/ && !unreleased_added {
            print $0
            print ""
            print "## [Unreleased]"
            print ""
            print "### Added"
            print ""
            print "### Changed"
            print ""
            print "### Fixed"
            print ""
            print "### Removed"
            print ""
            unreleased_added = 1
            next
        }
        { print }
        ' "$CHANGELOG_FILE" > "$CHANGELOG_FILE.tmp"
        
        mv "$CHANGELOG_FILE.tmp" "$CHANGELOG_FILE"
        echo "Successfully added [Unreleased] section to CHANGELOG.md"
        ;;
        
    *)
        echo "Error: Invalid action: $ACTION"
        echo "Valid actions: release, unreleased"
        rm -f "$CHANGELOG_FILE.bak"
        exit 1
        ;;
esac

# Verify the changes using a simpler grep pattern
if [[ "$ACTION" == "release" ]]; then
    # Use a more flexible pattern that doesn't rely on complex escaping
    if grep -F "## [$VERSION] - $CURRENT_DATE" "$CHANGELOG_FILE" > /dev/null; then
        echo "✓ Version entry verified in CHANGELOG.md"
    else
        echo "✗ Failed to verify version entry in CHANGELOG.md"
        echo "Looking for: ## [$VERSION] - $CURRENT_DATE"
        echo "Found in file:"
        grep "## \[$VERSION\]" "$CHANGELOG_FILE" 2>/dev/null || echo "No matching version entries found"
        # Restore backup
        mv "$CHANGELOG_FILE.bak" "$CHANGELOG_FILE"
        exit 1
    fi
elif [[ "$ACTION" == "unreleased" ]]; then
    if grep -F "## [Unreleased]" "$CHANGELOG_FILE" > /dev/null; then
        echo "✓ [Unreleased] section verified in CHANGELOG.md"
    else
        echo "✗ Failed to verify [Unreleased] section in CHANGELOG.md"
        # Restore backup
        mv "$CHANGELOG_FILE.bak" "$CHANGELOG_FILE"
        exit 1
    fi
fi

# Remove backup
rm -f "$CHANGELOG_FILE.bak"

echo "CHANGELOG.md update completed successfully"
