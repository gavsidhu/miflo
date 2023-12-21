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
        # For Intel Macs
        BINARY_FILE="miflo_${LATEST_VERSION}_darwin_amd64.tar.gz"
        ;;
    "arm64")
        # For Apple Silicon Macs
        BINARY_FILE="miflo_${LATEST_VERSION}_darwin_arm64.tar.gz"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

BINARY_URL="https://github.com/$GITHUB_REPO/releases/download/v$LATEST_VERSION/$BINARY_FILE"

# Download the binary
curl -L $BINARY_URL -o miflo.tar.gz


# Extract and move the binary to /usr/local/bin
tar -xzf miflo.tar.gz
sudo mv miflo /usr/local/bin/miflo
sudo chmod +x /usr/local/bin/miflo

echo "miflo installed successfully"

