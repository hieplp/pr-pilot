#!/usr/bin/env bash
set -e

BINARY="pr-pilot"
INSTALL_DIR="/usr/local/bin"

echo "Building $BINARY..."
go build -o "$BINARY" .

echo "Installing to $INSTALL_DIR/$BINARY..."
if [[ -w "$INSTALL_DIR" ]]; then
    mv "$BINARY" "$INSTALL_DIR/$BINARY"
else
    sudo mv "$BINARY" "$INSTALL_DIR/$BINARY"
fi

echo "Done. Run '$BINARY --help' to get started."
