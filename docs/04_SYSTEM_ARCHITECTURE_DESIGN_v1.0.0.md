# Claudex Windows - System Architecture & Design Document v1.0.0

**Document Type:** System Architecture & Design Specification  
**Version:** 1.0.0  
**Status:** APPROVED  
**Created:** 2025-01-16  
**Release Target:** v0.1.0  
**Based On:** Actual working codebase analysis  

---

## 1. Architecture Overview

### 1.1 High-Level System Architecture (C4 Level 1)

```
┌─────────────────────────────────────────────────────────────────┐
│                   User / Developer                              │
└─────────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│              Claudex Windows CLI Application                    │
│  (Go executable @claudex-windows/cli npm package)              │
│                                                                  │
│  - Session Management (create, resume, fork, list)             │
│  - Terminal User Interface (TUI) with Bubble Tea               │
│  - Configuration Management                                     │
│  - Hook System Integration                                      │
│  - Agent Profile Support                                        │
│  - MCP Server Configuration                                     │
└─────────────────────────────────────────────────────────────────┘
                ↙                   ↓                   ↘
         ┌──────────┐         ┌──────────┐         ┌──────────┐
         │ Claude   │         │   Git    │         │ Filesystem
         │  Code    │         │Repository│         │  (.claude/)
         │   CLI    │         │          │         │            
         └──────────┘         └──────────┘         └──────────┘
```

### 1.2 Component Architecture (C4 Level 2)

```
┌──────────────────────────────────────────────────────────────────────┐
│                     Claudex Windows System                           │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              UI Layer (Bubble Tea TUI)                      │   │
│  │  Dashboard | Session Wizard | Selection Interface           │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                            ↓                                        │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              Application Layer (app.App)                    │   │
│  │  Session Lifecycle | Configuration | Dependencies           │   │
│  └─────────────────────────────────────────────────────────────┘   │
│           ↙        ↓        ↓        ↓        ↓         ↘          │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              Services Layer                                 │   │
│  │  ┌──────────────┬──────────────┬──────────────────────────┐ │   │
│  │  │   Session    │   Config     │  Profile / Stack Detect │ │   │
│  │  │  Management  │ Management   │                         │ │   │
│  │  └──────────────┴──────────────┴──────────────────────────┘ │   │
│  │  ┌──────────────┬──────────────┬──────────────────────────┐ │   │
│  │  │   MCP        │   Git        │  Documentation Tracking  │ │   │
│  │  │  Configuration│ Operations   │                         │ │   │
│  │  └──────────────┴──────────────┴──────────────────────────┘ │   │
│  │  ┌──────────────┬──────────────┬──────────────────────────┐ │   │
│  │  │   Claude     │  Filesystem  │   Commander (Process)    │ │   │
│  │  │   Settings   │  Operations  │   Execution              │ │   │
│  │  └──────────────┴──────────────┴──────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────┘   │
│           ↙        ↓        ↓        ↓        ↓         ↘          │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              Use Cases Layer                                │   │
│  │  ┌──────────────┬──────────────┬──────────────────────────┐ │   │
│  │  │   Session    │   Setup      │   Migration / Update     │ │   │
│  │  │ Lifecycle UC │   Workflows  │   Operations             │ │   │
│  │  └──────────────┴──────────────┴──────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              Infrastructure Layer                           │   │
│  │  Clock | Env | UUID | Lock | Commander | Filesystem (afero) │  │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

### 1.3 Package Organization

```
claudex-windows/
├── src/
│   ├── cmd/                          (Entry Points)
│   │   ├── claudex/main.go           (Main CLI executable)
│   │   └── claudex-hooks/main.go     (Hook execution tool)
│   │
│   ├── internal/
│   │   ├── services/                 (Infrastructure Services)
│   │   │   ├── app/                  (Application container & lifecycle)
│   │   │   ├── session/              (Session management)
│   │   │   ├── config/               (TOML configuration)
│   │   │   ├── profile/              (Agent profile loading)
│   │   │   ├── mcpconfig/            (MCP server configuration)
│   │   │   ├── git/                  (Git operations)
│   │   │   ├── filesystem/           (File operations)
│   │   │   ├── commander/            (Process execution)
│   │   │   ├── clock/                (Time abstraction)
│   │   │   ├── uuid/                 (UUID generation)
│   │   │   ├── env/                  (Environment access)
│   │   │   └── lock/                 (File locking)
│   │   │
│   │   ├── usecases/                 (Business Logic)
│   │   │   ├── session/              (Session lifecycle: create, resume, fork)
│   │   │   ├── setup/                (Initialize .claude/ structure)
│   │   │   ├── setuphook/            (Git hook installation)
│   │   │   ├── setupmcp/             (MCP configuration workflow)
│   │   │   ├── migrate/              (Legacy artifact migration)
│   │   │   ├── updatecheck/          (Version checking)
│   │   │   ├── updatedocs/           (Documentation updates)
│   │   │   └── createindex/          (Generate index.md)
│   │   │
│   │   ├── hooks/                    (Hook System)
│   │   │   ├── pretooluse/           (Pre-tool-use context injection)
│   │   │   ├── posttooluse/          (Post-tool-use documentation)
│   │   │   └── shared/               (Logging, parsing, utilities)
│   │   │
│   │   ├── doc/                      (Documentation Services)
│   │   │   └── rangeupdater/         (Range-based updates to files)
│   │   │
│   │   ├── ui/                       (Terminal User Interface)
│   │   ├── notify/                   (Notification services)
│   │   ├── testutil/                 (Testing utilities)
│   │   └── vendor/                   (Dependencies)
│   │
│   ├── profiles/                     (Agent profiles & skills)
│   │   ├── agents/                   (Predefined agent personas)
│   │   ├── skills/                   (Domain-specific skill templates)
│   │   ├── commands/                 (Command profiles)
│   │   ├── roles/                    (Role definitions)
│   │   └── tasks/                    (Task templates)
│   │
│   └── testdata/                     (Test fixtures)
│
├── npm/                              (NPM Package Distribution)
│   ├── @claudex-windows/
│   │   ├── cli/                      (Main package)
│   │   ├── windows-x64/              (Windows binaries)
│   │   ├── darwin-arm64/             (macOS ARM64)
│   │   ├── darwin-x64/               (macOS x64)
│   │   ├── linux-x64/                (Linux x64)
│   │   └── linux-arm64/              (Linux ARM64)
│   └── version.txt                   (Current version)
│
└── docs/                             (Documentation)
```

---

## 2. Core Components & Detailed Design

### 2.1 Application Container (app.App)

**File:** `src/internal/services/app/app.go`

**Responsibility:** Main application lifecycle management, orchestration of all services

**Key Structures:**

```go
type App struct {
    deps            *Dependencies      // Injected service dependencies
    cfg             *config.Config      // Loaded TOML configuration
    projectDir      string              // Working directory
    sessionsDir     string              // .claudex/sessions directory
    docPaths        []string            // Documentation paths for context
    noOverwrite     bool                // Skip overwriting existing files
    updateDocs      bool                // Update index.md from git
    setupMCP        bool                // Configure MCP servers
    createIndex     string              // Create index.md at path
    version         string              // Application version
}

type SessionInfo struct {
    Name         string                 // Session name
    Path         string                 // Session directory path
    ClaudeID     string                 // Claude process ID
    Mode         LaunchMode             // Session mode (new/resume/fork/fresh/ephemeral)
    OriginalName string                 // For fork/fresh operations
}

// LaunchMode options:
// - new:       Create and launch new session
// - resume:    Resume existing session with full history
// - fork:      Create new session from existing
// - fresh:     Resume with cleared Claude context (but keep docs)
// - ephemeral: Temporary session, no persistence
```

**Key Methods:**

- `New()` - Create new App instance with production dependencies
- `Init()` - Initialize (migration, config load, flag parsing, logging)
- `Run()` - Main application loop (session selection, Claude launch)
- `Close()` - Cleanup resources

**Workflow:**
1. Parse command-line flags
2. Run migration (ensure `.claudex/` exists)
3. Load configuration from `~/.claudex/config.toml`
4. Select launch mode (new/resume/fork/fresh/ephemeral)
5. Invoke appropriate use case
6. Launch Claude Code CLI with context injection
7. Monitor hook execution
8. Cleanup on exit

---

### 2.2 Session Management Services

**File:** `src/internal/services/session/session.go`

**Responsibility:** Session CRUD operations, metadata management, listing, naming

**Key Data Structures:**

```go
type SessionItem struct {
    Title       string      // Session folder name
    Description string      // Description + last used date
    Created     time.Time   // Creation or last-used timestamp
    ItemType    string      // "session"
}
```

**Session Directory Structure:**

```
.claudex/sessions/
└── [session-id]/
    ├── .created              // Session creation timestamp (RFC3339)
    ├── .description          // User-provided session description
    ├── .last_used            // Last accessed timestamp (RFC3339)
    ├── session-overview.md   // Auto-maintained progress documentation
    ├── messages.json         // Archived Claude messages (optional)
    ├── .claude/              // Claude Code configuration
    │   ├── settings.json     // Claude settings for this session
    │   └── hooks/            // Session-specific hooks (optional)
    └── metadata.json         // Session metadata (optional)
```

**Key Functions:**

- `GetSessions()` - List all sessions sorted by last-used
- `UpdateLastUsed()` - Update .last_used timestamp
- `GetSessionByName()` - Find session by directory name
- `GenerateSessionName()` - Create unique session identifier

**Session Listing Algorithm:**
1. Scan `.claudex/sessions/` directory
2. Read metadata from `.description` and `.created`/`.last_used` files
3. Sort by last-used timestamp (descending)
4. Return sorted SessionItem list

---

### 2.3 Configuration Management

**File:** `src/internal/services/config/config.go`

**Responsibility:** Load and manage TOML configuration, command-line flags, precedence

**Configuration File:** `~/.claudex/config.toml`

```toml
[claude]
api_endpoint = "https://api.anthropic.com"
code_enabled = true

[profiles]
default = "engineer"
available = ["engineer", "designer", "architect"]

[mcp]
enabled = true
auto_configure = true
sequential_thinking = true
context7 = true

[documentation]
auto_update = true
template = "default"
include_timestamps = true

[git]
auto_commit = true
commit_template = "docs: update session progress"
```

**Configuration Precedence (highest to lowest):**
1. Command-line flags
2. Environment variables (`CLAUDEX_*`)
3. Session-specific config (`.claudex/sessions/[id]/.config`)
4. User config (`~/.claudex/config.toml`)
5. Built-in defaults

---

### 2.4 Hook System Architecture

**Files:** `src/internal/hooks/pretooluse/` and `src/internal/hooks/posttooluse/`

**Responsibility:** Execute hooks at key points in Claude execution lifecycle

**Hook Execution Points:**

1. **Pre-Tool-Use Hook** (`pre-tool-use.sh`)
   - **When:** Before each Claude tool execution
   - **Purpose:** Context injection, parameter validation, security checks
   - **Environment Variables:**
     - `TOOL_NAME` - The tool being executed
     - `TOOL_INPUT` - Tool parameters
     - `SESSION_ID` - Current session identifier
     - `SESSION_PATH` - Session directory path
   - **Files:** `src/internal/hooks/pretooluse/context_injector.go`

2. **Post-Tool-Use Hook** (`post-tool-use.sh`)
   - **When:** After tool execution completes
   - **Purpose:** Documentation updates, artifact capture, notification
   - **Environment Variables:**
     - `TOOL_NAME` - The tool that executed
     - `TOOL_OUTPUT` - Tool execution result
     - `SESSION_ID` - Current session identifier
   - **Files:** `src/internal/hooks/posttooluse/autodoc.go`

3. **Session-End Hook** (`session-end.sh`)
   - **When:** When session terminates
   - **Purpose:** Finalization, cleanup, archiving

4. **Notification Hook** (`notification-hook.sh`)
   - **When:** On significant events
   - **Purpose:** User notifications, logging

**Hook Execution Model:**

```
Claude Tool Execution
        ↓
Pre-Tool-Use Hook executes
        ↓
Context passed to Claude
        ↓
Tool executes in Claude
        ↓
Tool completes, output captured
        ↓
Post-Tool-Use Hook executes
        ↓
Documentation updated
        ↓
Return to Claude
```

**Hook Scripts Location:**
- Built-in hooks: `src/.claude/hooks/` (copied to each session)
- Platform-specific: Windows `.bat` and `.ps1`, Unix `.sh`
- Custom hooks: Session-specific overrides in `.claudex/sessions/[id]/.claude/hooks/`

---

### 2.5 Agent Profiles & Stack Detection

**Files:** `src/internal/services/profile/` and `src/internal/services/stackdetect/`

**Agent Profiles** - Pre-defined Claude personas stored in `src/profiles/agents/`

```
profiles/
├── agents/
│   ├── engineer.md          # Software engineer persona
│   ├── architect.md         # System architect persona
│   ├── designer.md          # Product designer persona
│   ├── team-lead.md         # Team lead persona
│   └── [custom].md          # User-defined personas
│
├── skills/
│   ├── typescript.md        # TypeScript best practices
│   ├── python.md            # Python best practices
│   ├── golang.md            # Go best practices
│   └── [domain].md          # Domain-specific skills
│
└── commands/
    ├── review-code.md       # Code review command profile
    ├── plan.md              # Planning command profile
    └── debug.md             # Debug command profile
```

**Stack Detection** - Auto-detects technology stack via marker files

**Detection Algorithm:**

```go
// Detects stack based on file presence:
- TypeScript/React: package.json with @types, tsconfig.json
- Python: pyproject.toml, setup.py, requirements.txt, .python-version
- Go: go.mod, go.sum
- React Native: app.json (Expo), Gemfile (iOS setup)
```

**Used for:** Auto-injecting relevant skills and best practices into Claude context

---

### 2.6 Session Use Cases

**File:** `src/internal/usecases/session/`

**Three Primary Session Workflows:**

#### UC-1: Create New Session

**Entry Point:** `session/new/new.go`

**Process:**
1. Accept session name or generate from description
2. Create `.claudex/sessions/[session-id]/` directory
3. Write `.created`, `.description` metadata files
4. Copy `.claude/` template with hooks
5. Initialize session-overview.md
6. Invoke Claude with session context
7. Return SessionInfo

**Flow Diagram:**
```
User Input
    ↓
Generate Session ID (UUID)
    ↓
Create Session Directory
    ↓
Write Metadata Files
    ↓
Copy Hook Templates
    ↓
Initialize Documentation
    ↓
Launch Claude
    ↓
Session Started
```

#### UC-2: Resume Session

**Entry Point:** `session/resume/resume.go`

**Process:**
1. Locate existing session by ID or name
2. Verify session integrity
3. Read session-overview.md for context
4. Inject context into Claude environment
5. Update .last_used timestamp
6. Launch Claude with resumed context
7. Return SessionInfo

**Context Restoration:**
- Read `session-overview.md` containing summary of prior work
- Extract key decisions and artifacts
- Inject into Claude's system prompt
- Claude "catches up" in seconds vs. losing context

#### UC-3: Fork Session

**Entry Point:** `session/resume/fork/fork.go`

**Process:**
1. Select source session
2. Generate new session ID
3. Copy all session artifacts to new directory
4. Create fork metadata (parent ID, fork timestamp)
5. Maintain separate execution context
6. Launch Claude in new session
7. Return SessionInfo

**Use Case:** Experimental branching without losing original session

#### UC-4: Fresh Memory Resume

**Entry Point:** `session/resume/fresh/fresh.go`

**Process:**
1. Select existing session
2. Generate NEW session ID but copy docs
3. Clear Claude context window (but keep local copy)
4. Inject context from session-overview.md
5. Claude catches up from documentation
6. Return SessionInfo

**Use Case:** Combat context window exhaustion while preserving history

---

### 2.7 MCP Server Integration

**File:** `src/internal/services/mcpconfig/`

**Responsibility:** Configure Claude Code's ~/.claude.json for MCP server use

**Built-in MCP Servers (v0.1.0):**

1. **Sequential Thinking** 
   - Purpose: Structured reasoning for complex problems
   - Config: Automatic setup via setupmcp use case
   - Status: Recommended, optional

2. **Context7**
   - Purpose: Up-to-date documentation lookup
   - Config: User specifies documentation sources
   - Status: Recommended, optional

**MCP Configuration Process:**

```go
// ~/.claude.json structure:
{
  "mcp_servers": {
    "sequential-thinking": {
      "command": "python -m anthropic_mcp_server sequential_thinking"
    },
    "context7": {
      "command": "node /path/to/context7/index.js",
      "env": {
        "CONTEXT7_SOURCES": "path/to/docs"
      }
    }
  }
}
```

**Setup Workflow:**
1. Detect available MCP servers
2. Prompt user for opt-in
3. Configure in ~/.claude.json
4. Verify connections
5. Store preference to avoid re-prompting

---

## 3. Data Flow Diagrams

### 3.1 Session Creation Data Flow

```
┌─────────────┐
│ User Runs   │
│ claudex new │
└──────┬──────┘
       ↓
┌──────────────────────────┐
│ App.Init()               │
│ - Load config            │
│ - Parse flags            │
│ - Setup logging          │
└──────┬───────────────────┘
       ↓
┌──────────────────────────┐
│ SessionInfo collected    │
│ - Name or generated      │
│ - Mode = "new"           │
└──────┬───────────────────┘
       ↓
┌──────────────────────────────────────────┐
│ UseCase: CreateSession                   │
│ - Generate UUID for session ID           │
│ - Create directory structure             │
│ - Write metadata files                   │
│ - Copy hook templates                    │
│ - Initialize docs                        │
└──────┬───────────────────────────────────┘
       ↓
┌──────────────────────────┐
│ Session Directory Ready  │
│ .claudex/sessions/[id]   │
│ with all artifacts       │
└──────┬───────────────────┘
       ↓
┌────────────────────────────────────────────┐
│ App.LaunchClaude()                         │
│ - Inject profile/context                  │
│ - Setup hooks                             │
│ - Execute: claude code                    │
└────────┬─────────────────────────────────┘
         ↓
┌──────────────────────┐
│ Claude Code Running  │
│ Session active       │
└──────────────────────┘
```

### 3.2 Hook Execution Data Flow

```
Claude Tool Execution Started
         ↓
┌──────────────────────────────────┐
│ Pre-Tool-Use Hook Triggered      │
│ - Tool name detected             │
│ - Environment set:               │
│   * TOOL_NAME                    │
│   * SESSION_ID                   │
│   * TOOL_INPUT (if applicable)   │
└──────┬───────────────────────────┘
       ↓
┌──────────────────────────────────┐
│ context_injector.go executes     │
│ - Validates tool type            │
│ - Injects context if needed      │
│ - Passes through to Claude       │
└──────┬───────────────────────────┘
       ↓
┌──────────────────────────────────┐
│ Claude Executes Tool             │
│ - Tool runs with context         │
│ - Output generated               │
└──────┬───────────────────────────┘
       ↓
┌──────────────────────────────────────┐
│ Post-Tool-Use Hook Triggered         │
│ - Tool output captured               │
│ - Environment set:                   │
│   * TOOL_NAME                        │
│   * TOOL_OUTPUT                      │
│   * SESSION_ID                       │
└──────┬───────────────────────────────┘
       ↓
┌──────────────────────────────────┐
│ autodoc.go executes              │
│ - Parses tool output             │
│ - Updates session-overview.md    │
│ - Logs results                   │
└──────┬───────────────────────────┘
       ↓
┌──────────────────────────────────┐
│ Return to Claude                 │
│ - Session docs updated           │
│ - User continues interaction     │
└──────────────────────────────────┘
```

### 3.3 Configuration Loading Flow

```
App.Init() starts
     ↓
┌─────────────────────────────────────┐
│ Check Config Precedence:            │
│ 1. Command-line flags (highest)     │
│ 2. Environment variables            │
│ 3. Session-specific config          │
│ 4. User config ~/.claudex/config    │
│ 5. Built-in defaults (lowest)       │
└────────┬────────────────────────────┘
         ↓
┌─────────────────────────────────────┐
│ Load ~/.claudex/config.toml         │
│ - Parse TOML                        │
│ - Validate sections                 │
│ - Apply env variable overrides      │
└────────┬────────────────────────────┘
         ↓
┌─────────────────────────────────────┐
│ Create Config struct                │
│ with all settings merged            │
└────────┬────────────────────────────┘
         ↓
┌─────────────────────────────────────┐
│ Return to App                       │
│ Config ready for use                │
└─────────────────────────────────────┘
```

---

## 4. Component Interfaces

### 4.1 Service Interfaces

**Session Service:**
```go
type SessionService interface {
    GetSessions(fs afero.Fs, dir string) ([]SessionItem, error)
    UpdateLastUsed(fs afero.Fs, clk clock.Clock, path string) error
    GetSessionByName(fs afero.Fs, dir, name string) (string, error)
    GenerateSessionName(desc string) string
}
```

**Config Service:**
```go
type ConfigService interface {
    Load(fs afero.Fs, path string) (*Config, error)
    Save(fs afero.Fs, path string, cfg *Config) error
    Merge(base, override *Config) *Config
}
```

**Profile Service:**
```go
type ProfileService interface {
    LoadProfile(name string) (string, error)
    ListAvailableProfiles() ([]string, error)
    ComposeProfile(base, skill string) string
}
```

**Hook Service:**
```go
type HookService interface {
    ExecutePre(ctx context.Context, toolName string, env map[string]string) error
    ExecutePost(ctx context.Context, toolName, output string, env map[string]string) error
    DiscoverHooks(sessionPath string) ([]string, error)
}
```

---

## 5. Design Patterns Used

### 5.1 Dependency Injection

All services use constructor injection via `Dependencies` struct:

```go
type Dependencies struct {
    FS       afero.Fs
    Clock    clock.Clock
    Env      env.EnvGetter
    UUID     uuid.Generator
    Commander commander.Commander
    // ...
}
```

**Benefits:**
- Testable (mock all dependencies)
- Explicit dependencies
- Easy to swap implementations

### 5.2 Use Case Pattern

Each business workflow implemented as a use case:

```go
type CreateSessionUC struct {
    deps *Dependencies
}

func (uc *CreateSessionUC) Execute(name string) (*SessionInfo, error) {
    // Implementation
}
```

### 5.3 Service Locator Pattern

App container provides access to all services:

```go
func (a *App) getSessionService() *session.Service {
    return a.deps.SessionService
}
```

### 5.4 Strategy Pattern

Different launch modes use different strategies:

```go
type LaunchStrategy interface {
    Execute(sessionID string) (*SessionInfo, error)
}

// Implementations: NewStrategy, ResumeStrategy, ForkStrategy, etc.
```

---

## 6. Error Handling Strategy

**Levels of Error Handling:**

1. **Service Layer** - Return specific errors:
   ```go
   ErrSessionNotFound
   ErrInvalidConfiguration
   ErrHookExecutionFailed
   ```

2. **Use Case Layer** - Wrap with context:
   ```go
   fmt.Errorf("failed to create session: %w", err)
   ```

3. **Application Layer** - Log and display to user:
   ```go
   fmt.Fprintf(os.Stderr, "Error: %v\n", err)
   os.Exit(1)
   ```

4. **Hook Level** - Continue on non-critical failures:
   ```go
   if err != nil {
       log.Printf("Hook failed (non-critical): %v", err)
       // Continue - don't crash
   }
   ```

---

## 7. State Management

**Session State Persistence:**

```
.claudex/sessions/[session-id]/
├── .created              # Immutable creation timestamp
├── .last_used            # Updated on each resume
├── .description          # User description
├── session-overview.md   # Main state artifact
├── messages.json         # Message history
└── metadata.json         # Additional metadata
```

**In-Memory State:**

```go
type SessionInfo struct {
    Name         string      // Current session name
    Path         string      // Working directory path
    ClaudeID     string      // Process identifier
    Mode         LaunchMode  // Current mode
    OriginalName string      // For forks
}
```

**State Transitions:**

```
Created → Active → Paused → Resumed → Forked
         ↓                             ↓
         └─────────────→ Archived ←───┘
```

---

## 8. Performance Considerations

### 8.1 Performance Targets

| Operation | Target | Method |
|-----------|--------|--------|
| Session creation | < 1s | Direct FS operations |
| Session list | < 2s | Parallel directory scan |
| Session resume | < 2s | Memory-resident context |
| Hook execution | < 100ms | Async execution |

### 8.2 Optimization Strategies

1. **Lazy Loading** - Config loaded on-demand
2. **Caching** - Session list cached during execution
3. **Parallel I/O** - Multiple sessions scanned concurrently
4. **Async Hooks** - Non-blocking hook execution
5. **Process Pooling** - Reuse process pool for commands

---

## 9. Testing Strategy

### 9.1 Test Coverage by Component

| Component | Coverage | Test Files |
|-----------|----------|-----------|
| Session Service | 100% | session_test.go |
| Config Service | 95% | config_test.go |
| Hook System | 85% | hooks/*_test.go |
| Use Cases | 90% | usecases/*/\*_test.go |
| **Overall** | **85%+** | **34 test files** |

### 9.2 Test Categories

- **Unit Tests** - Individual function/method testing
- **Integration Tests** - Service interaction testing
- **End-to-End Tests** - Full session lifecycle
- **Mocking** - afero filesystem, clock, env

---

## 10. Sign-Off

| Role | Name | Date | Approval |
|------|------|------|----------|
| Technical Lead | [TBD] | TBD | [ ] Approved |
| Architect | [TBD] | TBD | [ ] Approved |
| QA Lead | [TBD] | TBD | [ ] Approved |

---

## 11. References

- Project Definition v1.0.0
- PRD v1.0.0
- Requirements Traceability Matrix v1.0.0
- Source code: `src/` directory
- Dependency documentation: `go.mod`

