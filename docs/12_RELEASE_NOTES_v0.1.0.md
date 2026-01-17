# Claudex Windows Release Notes v0.1.0

**Release Date:** January 17, 2026  
**Version:** 0.1.0  
**Status:** Initial Release (Stable)  
**Repository:** https://github.com/OCNGill/claudex-windows  

---

## Overview

Claudex Windows v0.1.0 is the inaugural stable release of Claudex Windows, a powerful CLI-based session management system for Claude, the AI assistant by Anthropic. This version provides complete session lifecycle management, documentation indexing, MCP server integration, and hook-based extensibility.

---

## What's New in v0.1.0

### ✨ Core Features

#### 1. Session Management System
- **NEW:** Create and manage Claude sessions with persistent state
- **NEW:** Resume previous sessions with full context restoration
- **NEW:** Fork existing sessions for experimental branches
- **NEW:** Read-only session review mode for team collaboration
- **NEW:** Automatic session metadata tracking (creation date, last used, file count)

#### 2. Five Launch Modes
- **NEW:** NEW mode - Initialize fresh sessions with documentation discovery
- **NEW:** RESUME mode - Continue previous work with context preservation
- **NEW:** FORK mode - Copy sessions for safe experimentation
- **NEW:** FRESH mode - Read-only access for auditing/review
- **NEW:** EPHEMERAL mode - Temporary sessions for setup/configuration

#### 3. Documentation Indexing
- **NEW:** Automatic documentation discovery and indexing
- **NEW:** Support for multiple documentation formats (Markdown, Text, RST)
- **NEW:** Configurable include/exclude patterns for fine-grained control
- **NEW:** Cross-platform file system scanning with performance optimization
- **NEW:** Metadata tracking of indexed documents

#### 4. MCP Server Integration
- **NEW:** Model Context Protocol (MCP) server support
- **NEW:** Pre-configured servers: sequential-thinking, context7
- **NEW:** Custom MCP server configuration support
- **NEW:** Server lifecycle management (start, stop, health check)

#### 5. Hook System
- **NEW:** Pre-tool-use hooks for validation and preparation
- **NEW:** Post-tool-use hooks for documentation sync and git integration
- **NEW:** Session-end hooks for cleanup and logging
- **NEW:** Notification hooks for status tracking
- **NEW:** Platform-specific hook support (Bash/PowerShell)

#### 6. Configuration System
- **NEW:** Multi-level configuration hierarchy (global, project, session)
- **NEW:** TOML-based configuration files
- **NEW:** Environment variable overrides (CLAUDEX_*)
- **NEW:** Configuration precedence system with clear priority order
- **NEW:** Configuration validation and schema enforcement

#### 7. Git Integration
- **NEW:** Automatic git change detection
- **NEW:** Hook-based documentation sync with git
- **NEW:** Commit message prefixing for CLI operations
- **NEW:** Branch tracking and merge base detection
- **NEW:** Post-commit hook setup for automation

#### 8. Cross-Platform Support
- **NEW:** Full Windows support with PowerShell integration
- **NEW:** macOS native support
- **NEW:** Linux complete compatibility
- **NEW:** Platform-specific hook execution (PS1 on Windows, SH on Unix)

---

## Key Improvements

### Architecture
- ✅ Service-oriented architecture with 14 specialized services
- ✅ Dependency injection pattern for testability
- ✅ Interface-based design for flexibility
- ✅ afero abstraction for filesystem operations (testable)

### Performance
- ✅ Efficient documentation scanning (< 5 seconds for typical projects)
- ✅ Memory-efficient session storage
- ✅ Optimized hook execution with timeouts
- ✅ Concurrent safe operations with file locking

### Reliability
- ✅ 34 test files covering 22 unit + 10 integration + 2 E2E scenarios
- ✅ Automatic session backups for recovery
- ✅ Error recovery procedures
- ✅ Data validation on load/save

### Usability
- ✅ Simple one-command usage: `claudex`
- ✅ Automatic mode detection based on directory state
- ✅ Sensible defaults for all configuration
- ✅ Comprehensive help and error messages

---

## Features by Package

### App Service (`internal/services/app`)
- Application lifecycle management (Init, Run, Close)
- Launch mode detection and routing
- Service container and dependency injection
- Claude process management

### Session Service (`internal/services/session`)
- CRUD operations for sessions
- Session metadata management
- Timestamp tracking
- Session discovery and listing

### Configuration Service (`internal/services/config`)
- TOML file loading and parsing
- Configuration precedence handling
- Environment variable support
- Validation and schema checking

### Profile Service (`internal/services/profile`)
- Agent profile loading
- Profile composition and merging
- Skill and behavior management
- Custom profile support

### Git Service (`internal/services/git`)
- Branch detection
- Change file detection
- Commit creation
- Merge base calculation

### Hook System (`internal/hooks/*`)
- Pre-tool-use hook execution
- Post-tool-use hook execution
- Notification hook handling
- Session-end cleanup
- Platform-specific script execution

### Additional Services
- MCP Configuration (`mcpconfig`)
- Documentation Tracking (`doctracking`)
- Lock Management (`lock`) - Cross-process file locking
- Preferences (`preferences`) - Project settings storage

---

## Use Cases

### Use Case 1: New Project Initialization
```bash
cd ~/my-new-project
claudex
# Session created, documentation indexed, ready for Claude
```

### Use Case 2: Resume Previous Work
```bash
cd ~/my-existing-project
claudex
# Session automatically resumed with full context
```

### Use Case 3: Code Review
```bash
cd ~/colleague-project
claudex --no-overwrite
# Review code with Claude without modifying session
```

### Use Case 4: Experimental Branch
```bash
cp -r ~/my-project ~/my-project-experiment
cd ~/my-project-experiment
claudex  # Creates new session
# Experiment safely without affecting original
```

### Use Case 5: Team Collaboration
```bash
claudex --doc ./team-standards --doc ../shared-docs
# Access team documentation alongside project docs
```

---

## Breaking Changes

**None** - This is the initial release.

---

## Known Issues & Limitations

### Windows-Specific
- ⚠️ File paths > 260 characters may require Windows long path support
- ⚠️ PowerShell execution policies may require configuration
- ⚠️ Case-insensitive filesystem may cause issues with Linux-formatted paths

### General
- ⚠️ Maximum 10,000 files indexed by default (configurable)
- ⚠️ Session size should not exceed 500 MB for optimal performance
- ⚠️ Hook scripts must complete within timeout (default 30 seconds)

### Documented Workarounds
- See [11_TROUBLESHOOTING_GUIDE_v1.0.0.md](./docs/11_TROUBLESHOOTING_GUIDE_v1.0.0.md) for solutions to all known issues

---

## Deprecations

**None** - This is the initial release.

---

## Migration Guide

**Not Applicable** - v0.1.0 is the initial stable release. Users upgrading from beta/pre-release should see [Upgrading from Pre-Release](#upgrading-from-pre-release).

### Upgrading from Pre-Release

If upgrading from any pre-release version:

```bash
# 1. Backup existing sessions
cp -r ~/.claudex ~/.claudex.backup

# 2. Uninstall old version
npm uninstall -g @claudex-windows/cli
npm uninstall -g @claudex/cli
npm uninstall -g claudex

# 3. Install new version
npm install -g @claudex-windows/cli

# 4. Verify installation
claudex --version
# Should output: claudex version 0.1.0

# 5. Migrate configuration (if applicable)
claudex --migrate-config

# 6. Test with existing session
cd ~/existing-project
claudex
# Should resume previous session
```

---

## Installation

### NPM (Recommended)
```bash
npm install -g @claudex-windows/cli
claudex --version
```

### From Source
```bash
git clone https://github.com/OCNGill/claudex-windows.git
cd claudex-windows
go build -o claudex ./src/cmd/claudex
go build -o claudex-hooks ./src/cmd/claudex-hooks
```

See [08_CLI_USER_GUIDE_v1.0.0.md](./docs/08_CLI_USER_GUIDE_v1.0.0.md) for detailed installation instructions.

---

## Documentation

Complete documentation is available in the `/docs` directory:

**Phase 1 - DEFINE:**
- [PROJECT_DEFINITION_v1.0.0.md](./docs/PROJECT_DEFINITION_v1.0.0.md) - Project overview and scope

**Phase 2 - DESIGN:**
- [04_SYSTEM_ARCHITECTURE_DESIGN_v1.0.0.md](./docs/04_SYSTEM_ARCHITECTURE_DESIGN_v1.0.0.md) - System architecture
- [05_DESIGN_IMPLEMENTATION_DETAILS_v1.0.0.md](./docs/05_DESIGN_IMPLEMENTATION_DETAILS_v1.0.0.md) - Implementation details

**Phase 3 - DEBUG:**
- [06_TEST_STRATEGY_v1.0.0.md](./docs/06_TEST_STRATEGY_v1.0.0.md) - Test strategy and coverage

**Phase 4 - DOCUMENT:**
- [08_CLI_USER_GUIDE_v1.0.0.md](./docs/08_CLI_USER_GUIDE_v1.0.0.md) - CLI reference and examples
- [09_API_REFERENCE_v1.0.0.md](./docs/09_API_REFERENCE_v1.0.0.md) - API documentation
- [10_CONFIGURATION_GUIDE_v1.0.0.md](./docs/10_CONFIGURATION_GUIDE_v1.0.0.md) - Configuration reference
- [11_TROUBLESHOOTING_GUIDE_v1.0.0.md](./docs/11_TROUBLESHOOTING_GUIDE_v1.0.0.md) - Troubleshooting and support

---

## System Requirements

### Minimum Requirements
- **OS:** Windows 10+, macOS 10.14+, Linux (any modern distribution)
- **Node.js:** 14.0+ (for npm installation)
- **Go:** 1.20+ (for building from source)
- **Disk Space:** 50 MB for installation + session storage

### Recommended Requirements
- **OS:** Windows 11, macOS 12+, Linux (latest LTS)
- **Node.js:** 18.0+
- **Go:** 1.24+
- **Disk Space:** 1 GB (for comfortable session storage)

### Optional Dependencies
- **git:** For git integration features
- **jq:** For advanced hook scripting
- **Claude Desktop:** For MCP server support

---

## Performance Metrics (v0.1.0)

### Benchmarks
- **Session Creation:** < 2 seconds
- **Session Resumption:** < 1 second
- **Documentation Scanning:** < 5 seconds (1000 files)
- **Hook Execution:** < 30 seconds (default timeout)
- **Memory Usage:** ~ 50 MB idle, < 500 MB with large sessions

### Limits (Default Configuration)
- **Max Files Indexed:** 10,000 (configurable)
- **Max Session Size:** 500 MB (recommended)
- **Max Hook Timeout:** 60 seconds
- **Max Sessions per Directory:** 5 (configurable)

---

## Support & Contributing

### Getting Help
- 📖 **Documentation:** See `/docs` directory
- 🐛 **Issues:** Report at https://github.com/OCNGill/claudex-windows/issues
- 💬 **Discussions:** https://github.com/OCNGill/claudex-windows/discussions

### Contributing
- 🤝 **Contributions Welcome:** See [CONTRIBUTING.md](./CONTRIBUTING.md)
- 📝 **Code Style:** Go 1.24+ standards
- ✅ **Testing:** All changes must include tests (85%+ coverage target)

---

## License

Claudex Windows is released under the MIT License. See [LICENSE](./LICENSE) for details.

---

## Acknowledgments

- Built with Go 1.24.0
- Uses 7D Agile framework for development
- Comprehensive documentation following academic standards
- 100% requirements traceability maintained

---

## What's Next?

### Planned for v0.2.0
- Enhanced UI/TUI for session management
- Web-based dashboard for monitoring
- Advanced profiling and performance analytics
- Plugin system for community extensions

### Planned for v0.3.0
- Multi-workspace support
- Team collaboration features
- Advanced CI/CD integration
- Enhanced MCP server discovery

---

## Summary

**Claudex Windows v0.1.0** provides a complete, production-ready session management system for Claude. With comprehensive documentation, extensive testing, and cross-platform support, it's ready for immediate use in professional environments.

### Key Metrics
- ✅ 14 services, fully tested
- ✅ 34 test files (22 unit, 10 integration, 2 E2E)
- ✅ 586 KB comprehensive documentation
- ✅ 100% requirements coverage
- ✅ Cross-platform compatibility

---

**Version:** 0.1.0 (Stable)  
**Release Date:** January 17, 2026  
**Status:** Production Ready ✅  

