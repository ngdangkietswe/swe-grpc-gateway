#!/bin/bash

# build.sh: Script to update dependencies

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

# Update Go dependencies
update_dependencies() {
    log "INFO" "Checking and updating dependencies..."

    # Fetch specific versions in parallel (replace @main with @vX.Y.Z or @commit-hash if needed)
    GOPROXY=direct go get github.com/ngdangkietswe/swe-protobuf-shared@main &
    GOPROXY=direct go get github.com/ngdangkietswe/swe-go-common-shared@main &
    wait || handle_error "Failed to update dependencies"

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

    log "INFO" "Dependencies updated and vendored successfully!"
}

# Main execution
main() {
    log "INFO" "Starting build script..."

    # Check for required commands
    command -v go >/dev/null 2>&1 || handle_error "Go is not installed"
    command -v git >/dev/null 2>&1 || handle_error "Git is not installed"

    update_dependencies

    log "INFO" "Build script completed successfully!"
}

main