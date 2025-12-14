#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "Assembling npm packages..."

# Copy binaries to platform packages
for platform in darwin-arm64 darwin-x64 linux-x64 linux-arm64; do
    src_dir="$ROOT_DIR/dist/$platform"
    dest_dir="$ROOT_DIR/npm/@claudex/$platform/bin"

    if [ -d "$src_dir" ]; then
        mkdir -p "$dest_dir"
        cp "$src_dir/claudex" "$dest_dir/"
        cp "$src_dir/claudex-hooks" "$dest_dir/"
        chmod +x "$dest_dir/claudex" "$dest_dir/claudex-hooks"
        echo "✓ Assembled $platform"
    else
        echo "⚠ Skipped $platform (not built)"
    fi
done

echo "Assembly complete!"
