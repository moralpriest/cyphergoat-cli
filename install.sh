#!/bin/bash
set -e

BIN_DIR="/usr/local/bin"
BINARY_NAME="cyphergoat"
REMOTE_URL="https://github.com/moralpriest/cyphergoat-cli/releases/download/v1"
TMP_DIR=$(mktemp -d)
TMP_FILE="${TMP_DIR}/${BINARY_NAME}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
esac

case "$OS" in
    cygwin*|mingw*|msys*) OS="windows" ;;
esac

if [ "$OS" = "windows" ]; then
    ASSET_NAME="${BINARY_NAME}-${OS}-${ARCH}.exe"
else
    ASSET_NAME="${BINARY_NAME}-${OS}-${ARCH}"
fi

echo "Installing CypherGoat CLI..."
echo "Platform: ${OS}-${ARCH}"

echo "Downloading CypherGoat CLI..."
curl -fsSL "${REMOTE_URL}/${ASSET_NAME}" -o "${TMP_FILE}"

chmod +x "${TMP_FILE}"

echo "Moving binary to ${BIN_DIR}..."
if [ -w "${BIN_DIR}" ]; then
    mv "${TMP_FILE}" "${BIN_DIR}/${BINARY_NAME}"

    rm -rf "${TMP_DIR}"

    echo "CypherGoat CLI has been successfully installed!"
    echo "Run '${BINARY_NAME}' to get started."
else
    echo "Elevated permissions required to install to ${BIN_DIR}"
    sudo mv "${TMP_FILE}" "${BIN_DIR}/${BINARY_NAME}"

    rm -rf "${TMP_DIR}"

    echo "CypherGoat CLI has been successfully installed!"
    echo "Run '${BINARY_NAME}' to get started."
fi
