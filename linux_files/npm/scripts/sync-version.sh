#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
NPM_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

VERSION=$(cat "$NPM_DIR/version.txt")

if [ -z "$VERSION" ]; then
  echo "Error: version.txt is empty or missing" >&2
  exit 1
fi

echo "Syncing version $VERSION to all packages..."

# Update main package (@claudex-windows/cli)
jq ".version = \"$VERSION\" | .optionalDependencies[\"@claudex-windows/darwin-arm64\"] = \"$VERSION\" | .optionalDependencies[\"@claudex-windows/darwin-x64\"] = \"$VERSION\" | .optionalDependencies[\"@claudex-windows/linux-x64\"] = \"$VERSION\" | .optionalDependencies[\"@claudex-windows/linux-arm64\"] = \"$VERSION\" | .optionalDependencies[\"@claudex-windows/windows-x64\"] = \"$VERSION\"" \
  "$NPM_DIR/@claudex-windows/cli/package.json" > "$NPM_DIR/@claudex-windows/cli/package.json.tmp" && \
  mv "$NPM_DIR/@claudex-windows/cli/package.json.tmp" "$NPM_DIR/@claudex-windows/cli/package.json"
echo "✓ Updated @claudex-windows/cli"

# Update platform packages
for platform in darwin-arm64 darwin-x64 linux-x64 linux-arm64 windows-x64; do
  pkg_json="$NPM_DIR/@claudex-windows/$platform/package.json"
  jq ".version = \"$VERSION\"" "$pkg_json" > "$pkg_json.tmp" && mv "$pkg_json.tmp" "$pkg_json"
  echo "✓ Updated @claudex-windows/$platform"
done

echo "Version sync complete: $VERSION"
