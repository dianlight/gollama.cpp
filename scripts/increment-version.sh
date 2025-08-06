#!/bin/bash

# Script to increment the version in Makefile and gollama.go
# Usage: ./scripts/increment-version.sh [major|minor|patch]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
MAKEFILE="$PROJECT_ROOT/Makefile"
GOLLAMA_GO="$PROJECT_ROOT/gollama.go"

# Default increment type
INCREMENT_TYPE="${1:-minor}"

if [[ ! -f "$MAKEFILE" ]]; then
    echo "Error: Makefile not found at $MAKEFILE"
    exit 1
fi

if [[ ! -f "$GOLLAMA_GO" ]]; then
    echo "Error: gollama.go not found at $GOLLAMA_GO"
    exit 1
fi

# Extract current version from Makefile
CURRENT_VERSION=$(grep "^VERSION ?=" "$MAKEFILE" | sed 's/VERSION ?= //')

if [[ -z "$CURRENT_VERSION" ]]; then
    echo "Error: Could not find VERSION in Makefile"
    exit 1
fi

echo "Current version: $CURRENT_VERSION"

# Parse version parts
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# Validate version parts
if [[ ! "$MAJOR" =~ ^[0-9]+$ ]] || [[ ! "$MINOR" =~ ^[0-9]+$ ]] || [[ ! "$PATCH" =~ ^[0-9]+$ ]]; then
    echo "Error: Invalid version format: $CURRENT_VERSION"
    echo "Expected format: MAJOR.MINOR.PATCH (e.g., 1.0.0)"
    exit 1
fi

# Increment version based on type
case "$INCREMENT_TYPE" in
    major)
        NEW_MAJOR=$((MAJOR + 1))
        NEW_MINOR=0
        NEW_PATCH=0
        ;;
    minor)
        NEW_MAJOR=$MAJOR
        NEW_MINOR=$((MINOR + 1))
        NEW_PATCH=0
        ;;
    patch)
        NEW_MAJOR=$MAJOR
        NEW_MINOR=$MINOR
        NEW_PATCH=$((PATCH + 1))
        ;;
    *)
        echo "Error: Invalid increment type: $INCREMENT_TYPE"
        echo "Valid options: major, minor, patch"
        exit 1
        ;;
esac

NEW_VERSION="$NEW_MAJOR.$NEW_MINOR.$NEW_PATCH"
echo "New version: $NEW_VERSION"

# Backup original files
cp "$MAKEFILE" "$MAKEFILE.bak"
cp "$GOLLAMA_GO" "$GOLLAMA_GO.bak"

# Update VERSION in Makefile
sed -i.tmp "s/^VERSION ?= .*/VERSION ?= $NEW_VERSION/" "$MAKEFILE"
rm -f "$MAKEFILE.tmp"

# Update Version in gollama.go (handle tab character, avoid FullVersion)
sed -i.tmp "s/^\([[:space:]]*\)Version = \"[^\"]*\"/\1Version = \"$NEW_VERSION\"/" "$GOLLAMA_GO"
rm -f "$GOLLAMA_GO.tmp"

# Verify the changes
NEW_VERSION_MAKEFILE=$(grep "^VERSION ?=" "$MAKEFILE" | sed 's/VERSION ?= //')
NEW_VERSION_GO=$(grep -E "^\s*Version = " "$GOLLAMA_GO" | sed 's/.*Version = "\([^"]*\)".*/\1/')

if [[ "$NEW_VERSION_MAKEFILE" != "$NEW_VERSION" ]]; then
    echo "Error: Makefile version update failed"
    echo "Expected: $NEW_VERSION"
    echo "Got: $NEW_VERSION_MAKEFILE"
    # Restore backups
    mv "$MAKEFILE.bak" "$MAKEFILE"
    mv "$GOLLAMA_GO.bak" "$GOLLAMA_GO"
    exit 1
fi

if [[ "$NEW_VERSION_GO" != "$NEW_VERSION" ]]; then
    echo "Error: gollama.go version update failed"
    echo "Expected: $NEW_VERSION"
    echo "Got: $NEW_VERSION_GO"
    # Restore backups
    mv "$MAKEFILE.bak" "$MAKEFILE"
    mv "$GOLLAMA_GO.bak" "$GOLLAMA_GO"
    exit 1
fi

# Remove backups
rm -f "$MAKEFILE.bak"
rm -f "$GOLLAMA_GO.bak"

echo "Successfully updated version from $CURRENT_VERSION to $NEW_VERSION"
echo "Updated files:"
echo "  - $MAKEFILE"
echo "  - $GOLLAMA_GO"
