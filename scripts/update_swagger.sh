#!/bin/bash

# update_swagger.sh: Script to update OpenAPI specifications

set -euo pipefail

# Logging function
log() {
    local level="$1"
    local message="$2"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $message"
}

# Error handling function
handle_error() {
    log "ERROR" "$1"
    exit 1
}

# Check if a directory is writable
check_writable() {
    local dir="$1"
    if [ -d "$dir" ] && [ ! -w "$dir" ]; then
        handle_error "Directory $dir is not writable. Please check permissions."
    fi
}

# Update OpenAPI specifications
update_openapi() {
    log "INFO" "Updating OpenAPI specifications..."
    
    # Check if swe-protobuf-shared directory already exists
    local repo_dir="swe-protobuf-shared"
    if [ -d "$repo_dir" ]; then
        log "INFO" "Removing existing $repo_dir directory..."
        check_writable "$repo_dir"
        rm -rf "$repo_dir" || handle_error "Failed to remove $repo_dir"
    fi

    # Clone repository
    git clone https://github.com/ngdangkietswe/swe-protobuf-shared.git "$repo_dir" || handle_error "Failed to clone swe-protobuf-shared repository"

    # Copy OpenAPI files
    if [ -d "$repo_dir/openapiv2" ]; then
        mkdir -p swagger || handle_error "Failed to create swagger directory"
        cp -r "$repo_dir/openapiv2/"* swagger/ || handle_error "Failed to copy OpenAPI files"
    else
        handle_error "OpenAPI directory not found in $repo_dir"
    fi

    # Clean up
    log "INFO" "Cleaning up $repo_dir directory..."
    check_writable "$repo_dir"
    rm -rf "$repo_dir" || handle_error "Failed to remove $repo_dir"
    
    log "INFO" "OpenAPI specifications updated successfully!"
}

# Main execution
main() {
    log "INFO" "Starting OpenAPI update script..."

    # Check for required commands
    command -v git >/dev/null 2>&1 || handle_error "Git is not installed"

    update_openapi

    log "INFO" "OpenAPI update script completed successfully!"
}

main