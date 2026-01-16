#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "Assembling npm packages..."

# Copy binaries to platform packages
for platform in darwin-arm64 darwin-x64 linux-x64 linux-arm64 windows-x64; do
    src_dir="$ROOT_DIR/dist/$platform"
    dest_dir="$ROOT_DIR/npm/@claudex-windows/$platform/bin"

    if [ -d "$src_dir" ]; then
        mkdir -p "$dest_dir"
        if [ "$platform" = "windows-x64" ]; then
            cp "$src_dir/claudex-windows.exe" "$dest_dir/"
            cp "$src_dir/claudex-windows-hooks.exe" "$dest_dir/"
        else
            cp "$src_dir/claudex-windows" "$dest_dir/"
            cp "$src_dir/claudex-windows-hooks" "$dest_dir/"
            chmod +x "$dest_dir/claudex-windows" "$dest_dir/claudex-windows-hooks"
        fi
        echo "✓ Assembled $platform"
    else
        echo "⚠ Skipped $platform (not built)"
    fi
done

echo "Assembly complete!"
