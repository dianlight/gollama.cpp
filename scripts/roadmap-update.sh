#!/bin/bash
# roadmap-update.sh - Script to help update ROADMAP.md
# Usage: ./scripts/roadmap-update.sh [command] [options]

set -e

ROADMAP_FILE="docs/ROADMAP.md"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if roadmap file exists
check_roadmap_exists() {
    if [ ! -f "$ROOT_DIR/$ROADMAP_FILE" ]; then
        log_error "ROADMAP.md not found at $ROOT_DIR/$ROADMAP_FILE"
        exit 1
    fi
}

# Update the "Last Updated" date
update_last_updated() {
    local current_date=$(date "+%B %d, %Y")
    log_info "Updating last updated date to $current_date"
    
    if sed -i.bak "s/\*\*Last Updated\*\*: .*/\*\*Last Updated\*\*: $current_date/" "$ROOT_DIR/$ROADMAP_FILE"; then
        rm -f "$ROOT_DIR/$ROADMAP_FILE.bak"
        log_success "Updated last updated date"
    else
        log_error "Failed to update last updated date"
        exit 1
    fi
}

# Mark a feature as completed
complete_feature() {
    local feature_name="$1"
    if [ -z "$feature_name" ]; then
        log_error "Feature name required for completion"
        echo "Usage: $0 complete-feature \"Feature Name\""
        exit 1
    fi
    
    log_info "Marking feature '$feature_name' as completed"
    
    # This would need more sophisticated logic to move items between sections
    log_warning "Feature completion requires manual editing of ROADMAP.md"
    log_info "Please manually move '$feature_name' from its current section to '✅ Completed Major Features'"
    
    update_last_updated
}

# Add a new planned feature
add_feature() {
    local feature_name="$1"
    local priority="$2"
    local target_date="$3"
    
    if [ -z "$feature_name" ] || [ -z "$priority" ]; then
        log_error "Feature name and priority required"
        echo "Usage: $0 add-feature \"Feature Name\" \"Priority 1|2|3\" [target_date]"
        exit 1
    fi
    
    log_info "Adding new planned feature: '$feature_name' with priority $priority"
    log_warning "Feature addition requires manual editing of ROADMAP.md"
    log_info "Please add the feature to the appropriate priority section"
    
    update_last_updated
}

# Check for TODO/FIXME comments that might indicate roadmap items
scan_todos() {
    log_info "Scanning for TODO/FIXME comments that might indicate roadmap items..."
    
    # Scan for TODO/FIXME in Go files
    if find "$ROOT_DIR" -name "*.go" -type f -exec grep -Hn "TODO\|FIXME" {} \; > /tmp/roadmap_todos.txt 2>/dev/null; then
        if [ -s /tmp/roadmap_todos.txt ]; then
            log_warning "Found TODO/FIXME comments that might need roadmap entries:"
            cat /tmp/roadmap_todos.txt
        else
            log_success "No TODO/FIXME comments found in Go files"
        fi
    fi
    
    rm -f /tmp/roadmap_todos.txt
}

# Validate roadmap format
validate_roadmap() {
    log_info "Validating ROADMAP.md format..."
    
    # Check for required sections
    local required_sections=(
        "## Current Status"
        "## Short-term Goals"
        "## Medium-term Goals"
        "## Long-term Vision"
        "## Implementation Priorities"
        "## Success Metrics"
        "## Getting Involved"
    )
    
    local missing_sections=()
    for section in "${required_sections[@]}"; do
        if ! grep -q "$section" "$ROOT_DIR/$ROADMAP_FILE"; then
            missing_sections+=("$section")
        fi
    done
    
    if [ ${#missing_sections[@]} -eq 0 ]; then
        log_success "All required sections found"
    else
        log_warning "Missing sections:"
        printf '%s\n' "${missing_sections[@]}"
    fi
    
    # Check for proper checkbox format
    local checkbox_count=$(grep -c "^- \[ \]" "$ROOT_DIR/$ROADMAP_FILE" || true)
    local completed_count=$(grep -c "^- \[x\]" "$ROOT_DIR/$ROADMAP_FILE" || true)
    
    log_info "Found $checkbox_count planned items and $completed_count completed items"
    
    # Check if last updated date is recent (within 30 days)
    local last_updated_line=$(grep "Last Updated" "$ROOT_DIR/$ROADMAP_FILE" || echo "")
    if [ -n "$last_updated_line" ]; then
        log_info "Last updated info: $last_updated_line"
    else
        log_warning "No 'Last Updated' line found"
    fi
}

# Show help
show_help() {
    echo "ROADMAP.md Update Script"
    echo ""
    echo "Usage: $0 [command] [options]"
    echo ""
    echo "Commands:"
    echo "  update-date                    Update the 'Last Updated' date to today"
    echo "  complete-feature \"name\"       Mark a feature as completed (manual editing required)"
    echo "  add-feature \"name\" priority   Add a new planned feature (manual editing required)"
    echo "  scan-todos                     Scan for TODO/FIXME comments"
    echo "  validate                       Validate roadmap format and content"
    echo "  help                          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 update-date"
    echo "  $0 complete-feature \"Windows Runtime Support\""
    echo "  $0 add-feature \"Multi-GPU Support\" \"Priority 2\" \"Q1 2026\""
    echo "  $0 scan-todos"
    echo "  $0 validate"
    echo ""
    echo "Note: Most operations require manual editing of ROADMAP.md for safety."
    echo "This script primarily helps with validation and date updates."
}

# Main script logic
main() {
    check_roadmap_exists
    
    case "${1:-help}" in
        "update-date")
            update_last_updated
            ;;
        "complete-feature")
            complete_feature "$2"
            ;;
        "add-feature")
            add_feature "$2" "$3" "$4"
            ;;
        "scan-todos")
            scan_todos
            ;;
        "validate")
            validate_roadmap
            ;;
        "help"|"--help"|"-h")
            show_help
            ;;
        *)
            log_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
