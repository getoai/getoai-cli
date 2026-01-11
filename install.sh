#!/bin/bash
set -e

# getoai-cli installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/getoai/getoai-cli/master/install.sh | bash

REPO="getoai/getoai-cli"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="getoai"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        darwin) OS="darwin" ;;
        linux) OS="linux" ;;
        mingw*|msys*|cygwin*) OS="windows" ;;
        *) error "Unsupported OS: $OS" ;;
    esac

    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *) error "Unsupported architecture: $ARCH" ;;
    esac

    echo "${OS}-${ARCH}"
}

# Get the latest release version
get_latest_version() {
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# Download and install
install() {
    PLATFORM=$(detect_platform)
    info "Detected platform: $PLATFORM"

    VERSION=$(get_latest_version)
    if [ -z "$VERSION" ]; then
        error "Failed to get latest version"
    fi
    info "Latest version: $VERSION"

    # Construct download URL
    if [ "$OS" = "windows" ]; then
        FILENAME="${BINARY_NAME}-${PLATFORM}.zip"
    else
        FILENAME="${BINARY_NAME}-${PLATFORM}.tar.gz"
    fi
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

    info "Downloading from: $DOWNLOAD_URL"

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    # Download
    if command -v curl &> /dev/null; then
        curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/$FILENAME"
    elif command -v wget &> /dev/null; then
        wget -q "$DOWNLOAD_URL" -O "$TMP_DIR/$FILENAME"
    else
        error "Neither curl nor wget found. Please install one of them."
    fi

    # Extract
    cd "$TMP_DIR"
    if [ "$OS" = "windows" ]; then
        unzip -q "$FILENAME"
    else
        tar -xzf "$FILENAME"
    fi

    # Install
    BINARY="${BINARY_NAME}-${PLATFORM}"
    if [ "$OS" = "windows" ]; then
        BINARY="${BINARY}.exe"
    fi

    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY" "$INSTALL_DIR/$BINARY_NAME"
    else
        info "Requesting sudo access to install to $INSTALL_DIR"
        sudo mv "$BINARY" "$INSTALL_DIR/$BINARY_NAME"
    fi

    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    info "Successfully installed $BINARY_NAME to $INSTALL_DIR"
    info "Run 'getoai --help' to get started"
}

# Check if already installed
if command -v getoai &> /dev/null; then
    CURRENT_VERSION=$(getoai --version 2>/dev/null | head -1 || echo "unknown")
    warn "getoai is already installed (version: $CURRENT_VERSION)"
    read -p "Do you want to reinstall/upgrade? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        info "Installation cancelled"
        exit 0
    fi
fi

install
