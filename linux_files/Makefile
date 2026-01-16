.PHONY: all build clean install install-project uninstall help test run deps fmt vet check

# Configuration
SRC_DIR = src
CONFIG_DIR = $(HOME)/.config/claudex-windows
BIN_DIR = $(HOME)/.local/bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Default target
all: help

# Build the executable
build:
	@echo "Building claudex-windows $(VERSION)..."
	@cd $(SRC_DIR) && go build -ldflags "-X main.Version=$(VERSION)" -o ../claudex-windows ./cmd/claudex
	@echo "✓ Built: claudex-windows $(VERSION)"

# Build hooks binary
build-hooks:
	@echo "Building claudex-windows-hooks..."
	@cd $(SRC_DIR) && go build -o ../bin/claudex-windows-hooks ./cmd/claudex-hooks
	@echo "✓ Built: bin/claudex-windows-hooks"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@cd $(SRC_DIR) && go mod tidy
	@echo "✓ Dependencies installed"

# Run tests
test:
	@echo "Running tests..."
	@cd $(SRC_DIR) && go test -v ./...
	@echo "✓ Tests complete"

# Format code
fmt:
	@echo "Formatting code..."
	@cd $(SRC_DIR) && go fmt ./...
	@echo "✓ Formatting complete"

# Vet code
vet:
	@echo "Vetting code..."
	@cd $(SRC_DIR) && go vet ./...
	@echo "✓ Vet complete"

# Run all checks (fmt, vet, test) - use before submitting PRs
check: fmt vet test
	@echo "✓ All checks passed"

# Install hooks (binary + proxies)
install-hooks: build-hooks
	@echo "Installing claudex-windows-hooks..."
	@mkdir -p $(BIN_DIR) $(CONFIG_DIR)/hooks
	@install -m 755 bin/claudex-windows-hooks $(BIN_DIR)/claudex-windows-hooks
	@install -m 755 $(SRC_DIR)/scripts/proxies/*.sh $(CONFIG_DIR)/hooks/
	@echo "✓ Installed claudex-windows-hooks to $(BIN_DIR)"
	@echo "✓ Installed hook proxies to $(CONFIG_DIR)/hooks"

# Configure recommended MCPs (opt-in)
install-mcp: build
	@echo "Configuring recommended MCPs..."
	@./claudex-windows --setup-mcp

# Install claudex for current user
install: build build-hooks
	@echo "Installing claudex-windows..."
	@mkdir -p $(CONFIG_DIR) $(BIN_DIR)
	@cp -r $(SRC_DIR)/profiles $(CONFIG_DIR)/
	@ln -sf $(CURDIR)/claudex-windows $(BIN_DIR)/claudex-windows
	@install -m 755 bin/claudex-windows-hooks $(BIN_DIR)/claudex-windows-hooks
	@mkdir -p $(CONFIG_DIR)/hooks
	@install -m 755 $(SRC_DIR)/scripts/proxies/*.sh $(CONFIG_DIR)/hooks/
	@echo "✓ Installed to $(CONFIG_DIR)"
	@echo "✓ Binary linked at $(BIN_DIR)/claudex-windows"
	@echo "✓ Hooks installed to $(CONFIG_DIR)/hooks"
	@if ! echo "$$PATH" | grep -q "$(BIN_DIR)"; then \
		echo "⚠ Add to your shell config: export PATH=\"\$$HOME/.local/bin:\$$PATH\""; \
	fi

# Uninstall claudex
uninstall:
	@echo "Uninstalling claudex-windows..."
	@rm -rf $(CONFIG_DIR)
	@rm -f $(BIN_DIR)/claudex-windows
	@echo "✓ Uninstalled"

# Install to current project directory
install-project:
	@echo "Installing claudex-windows to current project..."
	@mkdir -p .claude
	@if [ -d "$(CONFIG_DIR)/profiles" ]; then \
		cp -r $(CONFIG_DIR)/profiles .claude/; \
		echo "✓ Copied profiles from $(CONFIG_DIR)"; \
	elif [ -d "profiles" ]; then \
		cp -r profiles .claude/; \
		echo "✓ Copied profiles from local directory"; \
	else \
		echo "✗ No profiles directory found"; \
		exit 1; \
	fi
	@if [ -d "$(CONFIG_DIR)/hooks" ]; then \
		cp -r $(CONFIG_DIR)/hooks .claude/; \
		echo "✓ Copied hooks from $(CONFIG_DIR)"; \
	elif [ -d ".claude/hooks" ]; then \
		echo "✓ Hooks already exist in .claude/"; \
	else \
		echo "⚠ No hooks directory found"; \
	fi
	@echo "✓ Project installation complete"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f claudex-windows
	@rm -rf bin/
	@echo "✓ Cleaned"

# Run the program
run: build
	@./claudex-windows

# npm packaging targets
NPM_PLATFORMS := darwin-arm64 darwin-amd64 linux-amd64 linux-arm64 windows-x64

.PHONY: npm-build npm-package npm-clean npm-sync-version npm-publish

npm-build: $(addprefix npm-build-,$(NPM_PLATFORMS))

npm-build-darwin-arm64:
	@echo "Building for darwin-arm64..."
	@mkdir -p dist/darwin-arm64
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -C src -ldflags "-X main.Version=$(VERSION)" -o ../dist/darwin-arm64/claudex ./cmd/claudex
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -C src -o ../dist/darwin-arm64/claudex-hooks ./cmd/claudex-hooks

npm-build-darwin-amd64:
	@echo "Building for darwin-amd64..."
	@mkdir -p dist/darwin-x64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -C src -ldflags "-X main.Version=$(VERSION)" -o ../dist/darwin-x64/claudex ./cmd/claudex
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -C src -o ../dist/darwin-x64/claudex-hooks ./cmd/claudex-hooks

npm-build-linux-amd64:
	@echo "Building for linux-amd64..."
	@mkdir -p dist/linux-x64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -C src -ldflags "-X main.Version=$(VERSION)" -o ../dist/linux-x64/claudex ./cmd/claudex
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -C src -o ../dist/linux-x64/claudex-hooks ./cmd/claudex-hooks

npm-build-linux-arm64:
	@echo "Building for linux-arm64..."
	@mkdir -p dist/linux-arm64
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C src -ldflags "-X main.Version=$(VERSION)" -o ../dist/linux-arm64/claudex-windows ./cmd/claudex
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C src -o ../dist/linux-arm64/claudex-windows-hooks ./cmd/claudex-hooks

npm-build-windows-x64:
	@echo "Building for windows-x64..."
	@mkdir -p dist/windows-x64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -C src -ldflags "-X main.Version=$(VERSION)" -o ../dist/windows-x64/claudex-windows.exe ./cmd/claudex
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -C src -o ../dist/windows-x64/claudex-windows-hooks.exe ./cmd/claudex-hooks

npm-package: npm-build
	@echo "Assembling npm packages..."
	./npm/scripts/assemble.sh

npm-sync-version:
	./npm/scripts/sync-version.sh

npm-publish: npm-package npm-sync-version
	./npm/scripts/publish.sh

npm-clean:
	rm -rf dist/
	rm -rf npm/@claudex/*/bin/*

# Show help
help:
	@echo "Available targets:"
	@echo "  make all             - Show this help (default)"
	@echo "  make build           - Build claudex-windows"
	@echo "  make build-hooks     - Build claudex-windows-hooks binary"
	@echo "  make deps            - Install/update dependencies"
	@echo "  make test            - Run tests"
	@echo "  make fmt             - Format code with go fmt"
	@echo "  make vet             - Run go vet"
	@echo "  make check           - Run fmt, vet, and test (pre-PR validation)"
	@echo "  make install         - Install claudex-windows and hooks to ~/.local/bin and ~/.config/claudex-windows"
	@echo "  make install-hooks   - Install only hooks binary and proxies"
	@echo "  make install-mcp     - Configure recommended MCP servers (sequential-thinking, context7)"
	@echo "  make uninstall       - Remove claudex-windows installation"
	@echo "  make install-project - Install profiles/hooks to current project .claude/"
	@echo "  make clean           - Remove build artifacts"
	@echo "  make run             - Build and run"
	@echo "  make npm-build       - Cross-compile binaries for all npm platforms"
	@echo "  make npm-package     - Build and assemble npm packages"
	@echo "  make npm-sync-version - Sync version from version.txt to all package.json files"
	@echo "  make npm-publish     - Build, sync versions, and publish to npm"
	@echo "  make npm-clean       - Remove npm build artifacts"
	@echo "  make help            - Show this help"
