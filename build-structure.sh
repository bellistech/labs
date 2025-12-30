#!/bin/bash

################################################################################
# Labs Repository Structure Builder
#
# This script creates the complete directory structure for the Labs repository
# including all course folders, lesson directories, example folders, reference
# materials, and utility scripts.
#
# Generated with assistance from Claude AI
# 
# Usage:
#   bash build-structure.sh                    # Build in current directory
#   bash build-structure.sh /path/to/labs      # Build in specific directory
#   bash build-structure.sh --verbose          # Show detailed output
#
# Requirements:
#   - bash 4.0+
#   - write permissions in target directory
#   - ~500MB disk space for empty structure
#
################################################################################

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
TARGET_DIR="${1:-.}"
VERBOSE="${VERBOSE:-false}"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

# Statistics
DIR_COUNT=0
FILE_COUNT=0

################################################################################
# Utility Functions
################################################################################

log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*" >&2
}

verbose() {
    if [[ "$VERBOSE" == "true" ]]; then
        echo -e "${BLUE}[VERBOSE]${NC} $*"
    fi
}

# Create directory and track count
create_dir() {
    local dir_path="$1"
    if [[ ! -d "$dir_path" ]]; then
        mkdir -p "$dir_path"
        ((DIR_COUNT++))
        verbose "Created directory: $dir_path"
    fi
}

# Create empty .gitkeep file to preserve empty directories
create_gitkeep() {
    local dir_path="$1"
    if [[ ! -f "$dir_path/.gitkeep" ]]; then
        touch "$dir_path/.gitkeep"
        ((FILE_COUNT++))
    fi
}

# Create directory structure for a course
create_course_structure() {
    local course_name="$1"
    local course_dir="$TARGET_DIR/$course_name"
    
    log_info "Building course structure: $course_name"
    
    # Create main course directory
    create_dir "$course_dir"
    
    # Create section directories
    create_dir "$course_dir/01-fundamentals"
    create_dir "$course_dir/01-fundamentals/examples"
    create_gitkeep "$course_dir/01-fundamentals/examples"
    
    create_dir "$course_dir/02-core-skills"
    create_dir "$course_dir/02-core-skills/examples"
    create_gitkeep "$course_dir/02-core-skills/examples"
    
    create_dir "$course_dir/03-advanced-patterns"
    create_dir "$course_dir/03-advanced-patterns/examples"
    create_gitkeep "$course_dir/03-advanced-patterns/examples"
    
    create_dir "$course_dir/04-labs"
    create_dir "$course_dir/04-labs/solutions"
    create_gitkeep "$course_dir/04-labs/solutions"
    
    create_dir "$course_dir/05-reference"
    
    create_dir "$course_dir/scripts"
    
    # Create placeholder files
    touch "$course_dir/README.md"
    touch "$course_dir/QUICKSTART.md"
    touch "$course_dir/01-fundamentals/lesson-01-basics.md"
    touch "$course_dir/01-fundamentals/lesson-02-core-concepts.md"
    touch "$course_dir/02-core-skills/lesson-01-implementation.md"
    touch "$course_dir/02-core-skills/lesson-02-best-practices.md"
    touch "$course_dir/03-advanced-patterns/lesson-01-architecture.md"
    touch "$course_dir/03-advanced-patterns/lesson-02-patterns.md"
    touch "$course_dir/04-labs/lab-01-scenario.md"
    touch "$course_dir/04-labs/lab-02-scenario.md"
    touch "$course_dir/05-reference/common-mistakes.md"
    touch "$course_dir/05-reference/troubleshooting.md"
    touch "$course_dir/05-reference/faq.md"
    touch "$course_dir/05-reference/additional-resources.md"
    touch "$course_dir/scripts/setup.sh"
    
    FILE_COUNT=$((FILE_COUNT + 15))
    
    verbose "Completed structure for: $course_name"
}

# Create Kubernetes course with extra sections
create_kubernetes_structure() {
    local course_dir="$TARGET_DIR/kubernetes"
    
    log_info "Building Kubernetes course (8-part deep dive)"
    create_dir "$course_dir"
    
    # Kubernetes has more sections
    for i in {01..07}; do
        local section_name=""
        case $i in
            01) section_name="fundamentals" ;;
            02) section_name="core-skills" ;;
            03) section_name="storage-and-config" ;;
            04) section_name="advanced-patterns" ;;
            05) section_name="operations" ;;
            06) section_name="labs" ;;
            07) section_name="reference" ;;
        esac
        
        create_dir "$course_dir/$i-$section_name"
        if [[ $i -lt 6 ]]; then
            create_dir "$course_dir/$i-$section_name/examples"
            create_gitkeep "$course_dir/$i-$section_name/examples"
        fi
    done
    
    # Create placeholder files
    touch "$course_dir/README.md"
    touch "$course_dir/QUICKSTART.md"
    
    # Files for each section
    for dir in "$course_dir"/0*-*/; do
        touch "${dir%/}/lesson-01-basics.md"
        touch "${dir%/}/lesson-02-core-concepts.md"
    done
    
    # Create solutions directory
    create_dir "$course_dir/04-labs/solutions"
    create_gitkeep "$course_dir/04-labs/solutions"
    
    # Create scripts directory
    create_dir "$course_dir/scripts"
    touch "$course_dir/scripts/setup-minikube.sh"
    touch "$course_dir/scripts/cleanup.sh"
    
    verbose "Completed Kubernetes course structure"
}

# Create NixOS course with extra sections
create_nixos_structure() {
    local course_dir="$TARGET_DIR/nixos"
    
    log_info "Building NixOS course (24-file curriculum)"
    create_dir "$course_dir"
    
    # NixOS has extended sections
    for i in {01..07}; do
        local section_name=""
        case $i in
            01) section_name="fundamentals" ;;
            02) section_name="core-configuration" ;;
            03) section_name="advanced-nix" ;;
            04) section_name="reproducibility" ;;
            05) section_name="operations" ;;
            06) section_name="labs" ;;
            07) section_name="reference" ;;
        esac
        
        create_dir "$course_dir/$i-$section_name"
        if [[ $i -lt 6 ]]; then
            create_dir "$course_dir/$i-$section_name/examples"
            create_gitkeep "$course_dir/$i-$section_name/examples"
        fi
    done
    
    # Create placeholder files
    touch "$course_dir/README.md"
    touch "$course_dir/QUICKSTART.md"
    
    # Files for each section
    for dir in "$course_dir"/0*-*/; do
        touch "${dir%/}/lesson-01-basics.md"
        touch "${dir%/}/lesson-02-core-concepts.md"
        touch "${dir%/}/lesson-03-advanced.md"
    done
    
    # Create solutions directory
    create_dir "$course_dir/06-labs/solutions"
    create_gitkeep "$course_dir/06-labs/solutions"
    
    # Create scripts directory
    create_dir "$course_dir/scripts"
    touch "$course_dir/scripts/setup.sh"
    touch "$course_dir/scripts/validate.sh"
    
    verbose "Completed NixOS course structure"
}

# Create Go course with extra sections
create_go_structure() {
    local course_dir="$TARGET_DIR/go"
    
    log_info "Building Go programming course"
    create_dir "$course_dir"
    
    # Go has systems-integration section
    for i in {01..05}; do
        local section_name=""
        case $i in
            01) section_name="fundamentals" ;;
            02) section_name="core-skills" ;;
            03) section_name="advanced-patterns" ;;
            04) section_name="systems-integration" ;;
            05) section_name="reference" ;;
        esac
        
        create_dir "$course_dir/$i-$section_name"
        if [[ $i -le 4 ]]; then
            create_dir "$course_dir/$i-$section_name/examples"
            create_gitkeep "$course_dir/$i-$section_name/examples"
        fi
    done
    
    # Create placeholder files
    touch "$course_dir/README.md"
    touch "$course_dir/QUICKSTART.md"
    
    # Files for each section
    for dir in "$course_dir"/0{1,2,3,4}-*/; do
        touch "${dir%/}/lesson-01-basics.md"
        touch "${dir%/}/lesson-02-core-concepts.md"
    done
    
    # Create labs directory
    create_dir "$course_dir/04-labs"
    create_dir "$course_dir/04-labs/solutions"
    create_gitkeep "$course_dir/04-labs/solutions"
    touch "$course_dir/04-labs/lab-01-scenario.md"
    touch "$course_dir/04-labs/lab-02-scenario.md"
    
    # Create scripts directory
    create_dir "$course_dir/scripts"
    touch "$course_dir/scripts/setup.sh"
    touch "$course_dir/scripts/run-tests.sh"
    
    verbose "Completed Go course structure"
}

################################################################################
# Main Execution
################################################################################

main() {
    # Validate target directory
    if [[ ! -d "$TARGET_DIR" ]]; then
        log_error "Target directory does not exist: $TARGET_DIR"
        exit 1
    fi
    
    if [[ ! -w "$TARGET_DIR" ]]; then
        log_error "No write permission in target directory: $TARGET_DIR"
        exit 1
    fi
    
    echo ""
    log_success "Starting Labs repository structure builder"
    log_info "Target directory: $TARGET_DIR"
    log_info "Timestamp: $TIMESTAMP"
    echo ""
    
    # Create root directory if needed
    create_dir "$TARGET_DIR"
    
    # List of standard courses (using basic structure)
    local standard_courses=(
        "terraform"
        "ansible"
        "bash"
        "monitoring"
        "logging"
        "gitops"
        "networking"
        "llm-deployment"
        "ebpf"
        "databases"
    )
    
    # Create standard courses
    for course in "${standard_courses[@]}"; do
        create_course_structure "$course" || log_warn "Failed to create $course"
    done
    
    # Create special courses with custom structures
    create_kubernetes_structure || log_warn "Failed to create Kubernetes course"
    create_nixos_structure || log_warn "Failed to create NixOS course"
    create_go_structure || log_warn "Failed to create Go course"
    
    # Create root-level documentation structure indicators
    create_dir "$TARGET_DIR/.github"
    create_dir "$TARGET_DIR/.github/workflows"
    create_gitkeep "$TARGET_DIR/.github/workflows"
    
    # Summary
    echo ""
    log_success "Repository structure created successfully!"
    echo ""
    echo "Summary:"
    echo "  Directories created: $DIR_COUNT"
    echo "  Files created: $FILE_COUNT"
    echo "  Total courses: 14"
    echo ""
    echo "Next steps:"
    echo "  1. Review the structure: ls -la $TARGET_DIR"
    echo "  2. Navigate to a course: cd $TARGET_DIR/terraform"
    echo "  3. Start adding content to README.md files"
    echo "  4. Create working examples in examples/ directories"
    echo "  5. Add lesson content and labs"
    echo ""
    log_info "To add content, start with: cat > $TARGET_DIR/terraform/README.md"
    echo ""
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [TARGET_DIR] [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  --verbose    Show detailed output during execution"
        echo "  --help       Display this help message"
        echo ""
        echo "Examples:"
        echo "  $0                          # Build in current directory"
        echo "  $0 /path/to/labs            # Build in specific directory"
        echo "  $0 . --verbose              # Build with verbose output"
        exit 0
        ;;
    --verbose)
        VERBOSE="true"
        ;;
esac

# Run main function
main

exit 0
