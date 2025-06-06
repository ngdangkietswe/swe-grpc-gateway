#!/bin/bash

# build.sh: Script to update dependencies and vendor them properly

set -euo pipefail
SECONDS=0  # Track total build time

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

# Update Go dependencies with forced update
update_dependencies() {
    log "INFO" "Checking and updating dependencies..."

    # Clean Go mod cache to avoid stale modules
    log "INFO" "Cleaning module cache..."
    go clean -modcache

    # Fetch the latest version forcefully
    log "INFO" "Fetching latest versions of required modules..."
    go get -u github.com/ngdangkietswe/swe-protobuf-shared@latest &
    go get -u github.com/ngdangkietswe/swe-go-common-shared@latest &
    wait || handle_error "Failed to update dependencies"

    # Ensure the module is actually used
    log "INFO" "Checking if module is used..."
    go mod why github.com/ngdangkietswe/swe-go-common-shared || log "WARNING" "Module common not used directly in code."
    go mod why github.com/ngdangkietswe/swe-protobuf-shared || log "WARNING" "Module protobuf not used directly in code."

    # Tidy modules
    log "INFO" "Tidying Go modules..."
    go mod tidy || handle_error "Failed to tidy Go modules"

    # Clear vendor directory if it exists
    if [ -d "vendor" ]; then
        log "INFO" "Clearing existing vendor directory..."
        check_writable "vendor"
        rm -rf vendor || handle_error "Failed to remove vendor directory"
    fi

    # Vendor dependencies
    log "INFO" "Vendoring dependencies..."
    go mod vendor -v || handle_error "Failed to vendor Go modules"

    # Confirm vendor content
    # Check for swe-go-common-shared
    if [ -d "vendor/github.com/ngdangkietswe/swe-go-common-shared" ]; then
        log "INFO" "✅ Module 'swe-go-common-shared' vendored successfully!"
    else
        log "WARNING" "⚠ Module 'swe-go-common-shared' not found in vendor! Ensure it is imported in the code."
    fi

    # Check for swe-protobuf-shared
    if [ -d "vendor/github.com/ngdangkietswe/swe-protobuf-shared" ]; then
        log "INFO" "✅ Module 'swe-protobuf-shared' vendored successfully!"
    else
        log "WARNING" "⚠ Module 'swe-protobuf-shared' not found in vendor! Ensure it is imported in the code."
    fi
}

# Main execution
main() {
    log "INFO" "Starting build script..."

    # Check for required commands
    command -v go >/dev/null 2>&1 || handle_error "Go is not installed"
    command -v git >/dev/null 2>&1 || handle_error "Git is not installed"

    update_dependencies

    duration=$SECONDS
    log "INFO" "Build completed in $((duration / 60)) minute(s) and $((duration % 60)) second(s)"
}

main
