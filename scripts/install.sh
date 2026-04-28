#!/bin/sh
set -e

REPO="github.com/suapapa/gox"
BINARY="gox"

# Check if go is installed
if ! command -v go >/dev/null 2>&1; then
    echo "Error: 'go' tool is not installed."
    echo "Please install the Go toolchain from https://go.dev/doc/install and try again."
    exit 1
fi

echo "Installing $BINARY from $REPO..."
go install "$REPO@latest"

# Check if the binary is in PATH
GOPATH_BIN="$(go env GOPATH)/bin"
if [ -n "$GOBIN" ]; then
    GOPATH_BIN="$GOBIN"
fi

if ! echo "$PATH" | grep -q "$GOPATH_BIN"; then
    echo ""
    echo "Warning: $GOPATH_BIN is not in your PATH."
    echo "To run '$BINARY', please add it to your shell configuration (e.g., .zshrc or .bashrc):"
    echo ""
    echo "    export PATH=\"$GOPATH_BIN:\$PATH\""
    echo ""
else
    echo ""
    echo "$BINARY installed successfully!"
fi
