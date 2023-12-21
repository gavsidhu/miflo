#!/bin/bash
set -e

# Define the GitHub repository
GITHUB_REPO="gavsidhu/miflo"

# Fetch the latest release tag (version) and remove the 'v' prefix
LATEST_VERSION=$(curl -s https://api.github.com/repos/$GITHUB_REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')

# Construct the binary URL using the latest version without the 'v'
BINARY_URL="https://github.com/$GITHUB_REPO/releases/download/v$LATEST_VERSION/miflo_${LATEST_VERSION}_darwin_arm64.tar.gz"

# Debug: Print the URL
echo "Downloading from: $BINARY_URL"

# Download the binary
curl -L $BINARY_URL -o miflo.tar.gz

# Debug: Check the downloaded file type
file miflo.tar.gz

# Extract and move the binary to /usr/local/bin
tar -xzf miflo.tar.gz
sudo mv miflo /usr/local/bin/miflo
sudo chmod +x /usr/local/bin/miflo

echo "miflo installed successfully"

