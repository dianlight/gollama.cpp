#!/bin/bash

# Test script for the changelog update process
# This simulates the tag-release workflow without git operations

set -e

echo "=== Testing CHANGELOG.md update workflow ==="

# Get current version from Makefile
CURRENT_VERSION=$(grep "^VERSION ?=" Makefile | sed 's/VERSION ?= //')
echo "Current version: $CURRENT_VERSION"

# Step 1: Update CHANGELOG.md for release
echo ""
echo "Step 1: Converting [Unreleased] to [$CURRENT_VERSION] with current date"
bash scripts/update-changelog.sh "$CURRENT_VERSION" "release"

# Show the result
echo ""
echo "=== CHANGELOG.md after release update ==="
head -20 CHANGELOG.md

# Step 2: Simulate version increment
echo ""
echo "Step 2: Simulating version increment..."
NEW_VERSION=$(echo $CURRENT_VERSION | awk -F. '{print $1"."$2"."$3+1}')
echo "Next version would be: $NEW_VERSION"

# Step 3: Add new [Unreleased] section
echo ""
echo "Step 3: Adding new [Unreleased] section for version $NEW_VERSION"
bash scripts/update-changelog.sh "$NEW_VERSION" "unreleased"

# Show the final result
echo ""
echo "=== CHANGELOG.md after adding [Unreleased] ==="
head -25 CHANGELOG.md

echo ""
echo "=== Test completed successfully! ==="
echo "The CHANGELOG.md has been updated following the release workflow."
echo "Remember to restore the original state if this was just a test."
