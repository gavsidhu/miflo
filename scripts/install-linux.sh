#!/bin/bash
set -e

GITHUB_REPO="gavsidhu/miflo"

# Fetch the latest release tag (version) and remove the 'v' prefix
LATEST_VERSION=$(curl -s https://api.github.com/repos/$GITHUB_REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')

# Detect the architecture
ARCH=$(uname -m)
BINARY_FILE=""

case "$ARCH" in
    "x86_64")
        BINARY_FILE="miflo_${LATEST_VERSION}_linux_amd64"
        ;;
    "aarch64")
        BINARY_FILE="miflo_${LATEST_VERSION}_linux_arm64"
        ;;
    "armv6l")
        BINARY_FILE="miflo_${LATEST_VERSION}_linux_armv6"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

BINARY_URL="https://github.com/$GITHUB_REPO/releases/download/v$LATEST_VERSION/$BINARY_FILE"

# Download the binary
sudo curl -L $BINARY_URL -o /usr/local/bin/miflo

# Make the binary executable
sudo chmod +x /usr/local/bin/miflo

echo "miflo installed successfully"

