# Releasing claudex

This document describes how to release new versions of claudex to npm.

## Prerequisites

- npm account with publish access to `claudex` and `@claudex/*` packages
- `NPM_TOKEN` secret configured in GitHub repository (for automated releases)
- Go 1.24+ installed locally (for manual releases)

## Version Format

claudex uses semantic versioning: `MAJOR.MINOR.PATCH`

- **MAJOR**: Breaking changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible

## Automated Release (Recommended)

1. Update the version in `npm/version.txt`:
   ```bash
   echo "1.0.0" > npm/version.txt
   ```

2. Commit the version change:
   ```bash
   git add npm/version.txt
   git commit -m "chore: bump version to 1.0.0"
   ```

3. Create and push a version tag:
   ```bash
   git tag v1.0.0
   git push origin main --tags
   ```

4. GitHub Actions will automatically:
   - Build binaries for all platforms (darwin-arm64, darwin-x64, linux-x64, linux-arm64)
   - Sync versions across all package.json files
   - Publish platform packages first, then the main package

## Manual Release

If you need to release manually (e.g., GitHub Actions is down):

1. Ensure you're logged into npm:
   ```bash
   npm login
   npm whoami  # Verify login
   ```

2. Update version and build:
   ```bash
   echo "1.0.0" > npm/version.txt
   make npm-package
   ```

3. Publish:
   ```bash
   make npm-publish
   ```

## Package Structure

The npm distribution consists of 5 packages:

| Package | Description |
|---------|-------------|
| `claudex` | Main package with bin wrappers and postinstall |
| `@claudex/darwin-arm64` | macOS ARM64 (Apple Silicon) binary |
| `@claudex/darwin-x64` | macOS x64 (Intel) binary |
| `@claudex/linux-x64` | Linux x64 binary |
| `@claudex/linux-arm64` | Linux ARM64 binary |

The main package uses `optionalDependencies` to automatically download only the platform-specific binary needed.

## Troubleshooting

### "Not logged in to npm"

Run `npm login` and authenticate with your npm credentials.

### "Package @claudex/darwin-arm64 not found"

Platform packages must be published before the main package. The publish script handles this automatically, but if publishing manually, ensure you publish in the correct order.

### Build fails for a specific platform

Check the Go cross-compilation setup. All builds use `CGO_ENABLED=0` for static linking.

## Testing a Release Locally

Before publishing, test the packages locally:

```bash
# Build and package
make npm-package

# Create test directory
mkdir /tmp/claudex-npm-test
cd /tmp/claudex-npm-test

# Pack packages
cd /path/to/claudex/npm/claudex
npm pack

cd /path/to/claudex/npm/@claudex/darwin-arm64  # or your platform
npm pack

# Test installation
npm install -g ./claudex-*.tgz
claudex --version
```
