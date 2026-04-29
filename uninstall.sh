#!/usr/bin/env bash
set -e

BINARY="pr-pilot"
INSTALL_DIR="/usr/local/bin"
TARGET="$INSTALL_DIR/$BINARY"

if [[ -f "$TARGET" ]]; then
    echo "Removing $TARGET..."
    if [[ -w "$INSTALL_DIR" ]]; then
        rm "$TARGET"
    else
        sudo rm "$TARGET"
    fi
    echo "Done. $BINARY has been uninstalled."
else
    echo "$TARGET not found. Nothing to uninstall."
fi
