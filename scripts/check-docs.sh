#!/bin/bash

# Documentation Sync Checker
# Run this script before committing to check if your changes need documentation updates

set -e

echo "🔍 Checking if your changes need documentation updates..."
echo ""

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ This script must be run from within a git repository"
    exit 1
fi

# Get staged and unstaged changes
STAGED_FILES=$(git diff --cached --name-only 2>/dev/null || true)
UNSTAGED_FILES=$(git diff --name-only 2>/dev/null || true)
ALL_CHANGED_FILES="$STAGED_FILES"$'\n'"$UNSTAGED_FILES"

if [ -z "$ALL_CHANGED_FILES" ]; then
    echo "ℹ️  No changes detected. Make some changes first!"
    exit 0
fi

echo "Changed files:"
echo "$ALL_CHANGED_FILES" | grep -v '^$' | sed 's/^/  - /'
echo ""

# Analyze changes
GO_FILES_CHANGED=$(echo "$ALL_CHANGED_FILES" | grep -E '\.go$' || true)
API_FILES_CHANGED=$(echo "$ALL_CHANGED_FILES" | grep -E '^(gollama\.go|platform_.*\.go|config\.go|loader\.go)$' || true)
EXAMPLE_FILES_CHANGED=$(echo "$ALL_CHANGED_FILES" | grep -E '^examples/' || true)
CI_FILES_CHANGED=$(echo "$ALL_CHANGED_FILES" | grep -E '\.github/workflows/' || true)
DOC_FILES_CHANGED=$(echo "$ALL_CHANGED_FILES" | grep -E '\.(md|txt)$' || true)

echo "📊 Change Analysis:"
echo "  - Go files: $([ -n "$GO_FILES_CHANGED" ] && echo "✓ Modified" || echo "○ Unchanged")"
echo "  - API files: $([ -n "$API_FILES_CHANGED" ] && echo "✓ Modified" || echo "○ Unchanged")"
echo "  - Examples: $([ -n "$EXAMPLE_FILES_CHANGED" ] && echo "✓ Modified" || echo "○ Unchanged")"
echo "  - CI config: $([ -n "$CI_FILES_CHANGED" ] && echo "✓ Modified" || echo "○ Unchanged")"
echo "  - Documentation: $([ -n "$DOC_FILES_CHANGED" ] && echo "✓ Modified" || echo "○ Unchanged")"
echo ""

# Generate recommendations
NEEDS_ATTENTION=false

echo "💡 Recommendations:"

if [ -n "$API_FILES_CHANGED" ]; then
    echo "  🔧 API files changed - Consider updating:"
    echo "     • README.md (public API examples)"
    echo "     • Go doc comments in changed files"
    echo "     • CHANGELOG.md with new features"
    
    if [ -z "$DOC_FILES_CHANGED" ]; then
        echo "     ⚠️  No documentation files changed yet!"
        NEEDS_ATTENTION=true
    fi
    echo ""
fi

if [ -n "$EXAMPLE_FILES_CHANGED" ]; then
    echo "  📚 Example files changed - Check:"
    
    # Check if corresponding README files were updated
    for example_file in $EXAMPLE_FILES_CHANGED; do
        if [[ "$example_file" =~ ^examples/([^/]+)/ ]]; then
            example_dir="${BASH_REMATCH[1]}"
            example_readme="examples/$example_dir/README.md"
            
            if [ -f "$example_readme" ]; then
                if ! echo "$ALL_CHANGED_FILES" | grep -q "$example_readme"; then
                    echo "     ⚠️  $example_readme may need updates"
                    NEEDS_ATTENTION=true
                fi
            fi
        fi
    done
    echo ""
fi

if [ -n "$GO_FILES_CHANGED" ]; then
    # Check if go.mod changed
    if echo "$ALL_CHANGED_FILES" | grep -q 'go\.mod$'; then
        echo "  📦 go.mod changed - Consider:"
        echo "     • Updating CI dependencies if new system packages needed"
        echo "     • Updating Go version in CI if minimum version changed"
        echo "     • Updating installation instructions in README.md"
        
        if [ -z "$CI_FILES_CHANGED" ]; then
            echo "     ⚠️  CI configuration unchanged - check if updates needed"
            NEEDS_ATTENTION=true
        fi
        echo ""
    fi
fi

# Check if examples still compile
if [ -n "$GO_FILES_CHANGED" ] || [ -n "$EXAMPLE_FILES_CHANGED" ]; then
    echo "  🏗️  Testing example compilation..."
    
    COMPILE_ERRORS=false
    for example_dir in examples/*/; do
        if [ -f "$example_dir/go.mod" ]; then
            example_name=$(basename "$example_dir")
            echo -n "     • $example_name: "
            
            if (cd "$example_dir" && go build -o /tmp/example-test . > /dev/null 2>&1); then
                echo "✓"
            else
                echo "❌ Failed to compile"
                COMPILE_ERRORS=true
                NEEDS_ATTENTION=true
            fi
        fi
    done
    
    if [ "$COMPILE_ERRORS" = true ]; then
        echo "     ⚠️  Some examples failed to compile - fix before committing"
    fi
    echo ""
fi

# Check for TODO/FIXME in staged files
if [ -n "$STAGED_FILES" ]; then
    TODOS=$(git diff --cached | grep -E '^\+.*\b(TODO|FIXME|XXX|HACK)\b' || true)
    if [ -n "$TODOS" ]; then
        echo "  📝 Found TODO/FIXME in staged changes:"
        echo "$TODOS" | sed 's/^/     /'
        echo "     Consider creating GitHub issues to track these"
        echo ""
    fi
fi

# Check CHANGELOG.md
if [ -f "CHANGELOG.md" ] && [ -n "$GO_FILES_CHANGED" ]; then
    if ! grep -q "## \[Unreleased\]" CHANGELOG.md; then
        echo "  📋 CHANGELOG.md missing [Unreleased] section"
        echo "     Add section for tracking new changes"
        NEEDS_ATTENTION=true
        echo ""
    elif ! echo "$ALL_CHANGED_FILES" | grep -q "CHANGELOG.md"; then
        echo "  📋 Consider updating CHANGELOG.md with your changes"
        echo ""
    fi
fi

# Summary
echo "🎯 Summary:"
if [ "$NEEDS_ATTENTION" = true ]; then
    echo "  ⚠️  Some recommendations need attention before committing"
    echo "  📖 See .copilot/templates.md for update templates"
    echo "  🤖 GitHub Copilot will help with these updates automatically"
else
    echo "  ✅ All checks passed! Your changes look well-documented."
fi

echo ""
echo "💡 Tips:"
echo "  • Use 'git add .' to stage documentation updates"
echo "  • GitHub Copilot can help generate documentation automatically"
echo "  • Run this script again after making documentation updates"
