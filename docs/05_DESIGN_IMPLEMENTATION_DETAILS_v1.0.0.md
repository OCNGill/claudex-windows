# Claudex Windows - Design Implementation Details v1.0.0

**Document Type:** Detailed Component Design & Module Specifications  
**Version:** 1.0.0  
**Status:** APPROVED  
**Created:** 2025-01-16  
**Release Target:** v0.1.0  

---

## 1. Detailed Service Specifications

### 1.1 App Service (app.App)

**Location:** `src/internal/services/app/app.go`

**Core Responsibility:** Application lifecycle orchestration and main container

**Data Structures:**

```go
// Main application container
type App struct {
    // Configuration & State
    cfg             *config.Config
    projectDir      string
    sessionsDir     string
    docPaths        []string
    
    // Behavior Flags
    noOverwrite     bool
    updateDocs      bool
    setupMCP        bool
    createIndex     string
    version         string
    
    // Dependency Container
    deps            *Dependencies
}

// All injected dependencies
type Dependencies struct {
    FS              afero.Fs
    Clock           clock.Clock
    Env             env.EnvGetter
    UUID            uuid.Generator
    
    // Services
    SessionService  *session.Service
    ConfigService   *config.Service
    ProfileService  *profile.Service
    HookService     *hooks.Service
    CommandService  *commander.Service
    FilesystemSvc   *filesystem.Service
    GitService      *git.Service
    
    // Infrastructure
    Logger          *log.Logger
    Lock            *lock.FileLock
}
```

**Key Methods:**

| Method | Input | Output | Purpose |
|--------|-------|--------|---------|
| `New()` | - | `*App, error` | Create production App instance |
| `Init()` | - | `error` | Initialize (migration, config load, setup) |
| `Run()` | - | `error` | Main application loop |
| `Close()` | - | `error` | Cleanup & shutdown |
| `SelectLaunchMode()` | - | `LaunchMode, error` | Determine new/resume/fork/fresh/ephemeral |
| `LaunchClaude()` | `SessionInfo` | `error` | Execute Claude with context |

**Initialization Sequence:**

1. **New()** - Construct with production dependencies
2. **Init()** - Sequential initialization:
   - `a.runMigration()` - Ensure `.claudex/` structure
   - `a.loadConfig()` - Load from `~/.claudex/config.toml`
   - `a.parseFlags()` - Command-line flag override
   - `a.setupLogging()` - Initialize logger
3. **Run()** - Main loop:
   - `a.SelectLaunchMode()` - Determine mode (new/resume/fork/fresh)
   - `a.ExecuteUseCase()` - Run appropriate use case
   - `a.LaunchClaude()` - Execute Claude Code CLI
   - `a.MonitorExecution()` - Hook execution
4. **Close()** - Cleanup:
   - Flush logs
   - Release locks
   - Close processes

**Usage Example:**

```go
// Create app
app, err := app.New()
if err != nil {
    log.Fatal(err)
}
defer app.Close()

// Initialize
if err := app.Init(); err != nil {
    log.Fatal(err)
}

// Run main loop
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

---

### 1.2 Session Service (session.Service)

**Location:** `src/internal/services/session/session.go`

**Core Responsibility:** Session CRUD operations, listing, metadata management

**Data Structures:**

```go
// Public API
type SessionItem struct {
    Title       string      // Session name (.claudex/sessions/[title])
    Description string      // Full description + last used
    Created     time.Time   // Creation or last-used time
    ItemType    string      // Always "session"
}

// Internal representation
type sessionMetadata struct {
    Name        string
    Description string
    Created     time.Time
    LastUsed    time.Time
    ClaudeID    string      // Process ID if active
    LaunchMode  string      // how session was created
}
```

**Session Metadata Files:**

```
.claudex/sessions/[session-id]/
├── .created              → RFC3339 formatted creation timestamp
├── .last_used            → RFC3339 formatted last access timestamp
├── .description          → User-provided or auto-generated description
├── .claude_id            → Process ID if currently active
├── session-overview.md   → Main session documentation
└── metadata.json         → Structured metadata object
```

**Key Methods:**

| Method | Input | Output | Logic |
|--------|-------|--------|-------|
| `GetSessions()` | `fs, dir` | `[]SessionItem, error` | Scan `.claudex/sessions/`, read metadata, sort by last-used DESC |
| `UpdateLastUsed()` | `fs, clk, path` | `error` | Write current time to `.last_used` file |
| `GetSessionByName()` | `fs, dir, name` | `string, error` | Find session by partial name match |
| `GenerateSessionName()` | `desc` | `string` | Create UUID-based name, optionally append first words of description |
| `GetSessionMetadata()` | `fs, path` | `*sessionMetadata, error` | Read all metadata files into struct |

**Session Directory Creation:**

```go
func (s *Service) CreateSessionDir(fs afero.Fs, baseDir string, name string) (string, error) {
    // 1. Create base directory: baseDir/[name]
    sessionPath := filepath.Join(baseDir, name)
    if err := fs.MkdirAll(sessionPath, 0755); err != nil {
        return "", err
    }
    
    // 2. Write metadata files
    now := time.Now()
    if err := afero.WriteFile(fs, filepath.Join(sessionPath, ".created"), 
        []byte(now.Format(time.RFC3339)), 0644); err != nil {
        return "", err
    }
    
    // 3. Write description
    desc := fmt.Sprintf("Created: %s", now.Format("2006-01-02 15:04:05"))
    if err := afero.WriteFile(fs, filepath.Join(sessionPath, ".description"), 
        []byte(desc), 0644); err != nil {
        return "", err
    }
    
    // 4. Initialize session-overview.md
    overview := fmt.Sprintf("# Session: %s\n\nCreated: %s\n\n## Activities\n\n", 
        name, now.Format(time.RFC3339))
    if err := afero.WriteFile(fs, filepath.Join(sessionPath, "session-overview.md"), 
        []byte(overview), 0644); err != nil {
        return "", err
    }
    
    return sessionPath, nil
}
```

**Session Listing Algorithm:**

```go
func (s *Service) GetSessions(fs afero.Fs, dir string) ([]SessionItem, error) {
    items := []SessionItem{}
    
    // 1. Read sessions directory
    entries, err := afero.ReadDir(fs, dir)
    if err != nil {
        return nil, err
    }
    
    // 2. For each entry
    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }
        
        sessionPath := filepath.Join(dir, entry.Name())
        
        // 3. Read metadata
        descBytes, _ := afero.ReadFile(fs, filepath.Join(sessionPath, ".description"))
        createdBytes, _ := afero.ReadFile(fs, filepath.Join(sessionPath, ".created"))
        lastUsedBytes, _ := afero.ReadFile(fs, filepath.Join(sessionPath, ".last_used"))
        
        // 4. Parse timestamps
        created := parseTime(string(createdBytes))
        lastUsed := parseTime(string(lastUsedBytes))
        if lastUsed.IsZero() {
            lastUsed = created
        }
        
        // 5. Create SessionItem
        items = append(items, SessionItem{
            Title:       entry.Name(),
            Description: string(descBytes),
            Created:     lastUsed,  // Sort by last-used
            ItemType:    "session",
        })
    }
    
    // 6. Sort by Created (last-used DESC)
    sort.Slice(items, func(i, j int) bool {
        return items[i].Created.After(items[j].Created)
    })
    
    return items, nil
}
```

---

### 1.3 Configuration Service (config.Service)

**Location:** `src/internal/services/config/config.go`

**Core Responsibility:** TOML configuration loading and management

**Configuration Structure:**

```toml
# ~/.claudex/config.toml

[claude]
api_endpoint = "https://api.anthropic.com"
code_enabled = true
model = "claude-3-5-sonnet-20241022"

[profiles]
default = "engineer"
available = ["engineer", "designer", "architect", "team-lead"]
custom_profiles_dir = "~/.claudex/profiles"

[mcp]
enabled = true
auto_configure = true
sequential_thinking = true
context7 = true
default_servers = ["sequential-thinking"]

[documentation]
auto_update = true
include_timestamps = true
template = "default"
track_changes = true

[git]
auto_commit = true
commit_template = "docs: update session progress"
hooks_enabled = true

[sessions]
default_launch_mode = "resume"
auto_fork_on_context_exhaustion = true
retention_days = 90

[hooks]
pre_tool_use_enabled = true
post_tool_use_enabled = true
custom_hooks_dir = ".claude/hooks"
```

**Go Data Structures:**

```go
type Config struct {
    Claude          claudeConfig
    Profiles        profilesConfig
    MCP             mcpConfig
    Documentation   docConfig
    Git             gitConfig
    Sessions        sessionsConfig
    Hooks           hooksConfig
    Raw             *toml.Tree     // For extensions
}

type claudeConfig struct {
    APIEndpoint string   `toml:"api_endpoint"`
    CodeEnabled bool     `toml:"code_enabled"`
    Model       string   `toml:"model"`
}

type profilesConfig struct {
    Default            string   `toml:"default"`
    Available          []string `toml:"available"`
    CustomProfilesDir  string   `toml:"custom_profiles_dir"`
}

type mcpConfig struct {
    Enabled              bool     `toml:"enabled"`
    AutoConfigure        bool     `toml:"auto_configure"`
    SequentialThinking   bool     `toml:"sequential_thinking"`
    Context7             bool     `toml:"context7"`
    DefaultServers       []string `toml:"default_servers"`
}

// ... other config structs
```

**Configuration Loading Precedence:**

```go
func (s *Service) Load(projectDir string) (*Config, error) {
    cfg := &Config{}
    
    // 1. Load defaults (built-in)
    cfg = s.loadDefaults()
    
    // 2. Load user config (~/.claudex/config.toml)
    userCfgPath := filepath.Join(os.ExpandEnv("$HOME"), ".claudex", "config.toml")
    if userCfg, err := s.loadFile(userCfgPath); err == nil {
        cfg = s.merge(cfg, userCfg)  // Override defaults
    }
    
    // 3. Load environment variable overrides
    if os.Getenv("CLAUDEX_NO_OVERWRITE") != "" {
        cfg.CLI.NoOverwrite = true
    }
    // ... other env vars
    
    // 4. Load session-specific config (highest priority)
    sessionCfgPath := filepath.Join(projectDir, ".claudex", "sessions", 
        "[sessionID]", ".config.toml")
    if sessionCfg, err := s.loadFile(sessionCfgPath); err == nil {
        cfg = s.merge(cfg, sessionCfg)
    }
    
    return cfg, nil
}

// Merge strategy: override.Field ?? base.Field
func (s *Service) merge(base, override *Config) *Config {
    result := *base
    if override.Claude.APIEndpoint != "" {
        result.Claude.APIEndpoint = override.Claude.APIEndpoint
    }
    // ... merge all fields with non-empty override values
    return &result
}
```

---

### 1.4 Profile Service (profile.Service)

**Location:** `src/internal/services/profile/profile.go`

**Core Responsibility:** Load and compose agent profiles from files

**Profile Structure:**

```markdown
# profiles/agents/engineer.md

## Engineer Profile

You are an expert software engineer with deep knowledge of modern software development practices.

### Strengths
- Architecture design
- Code optimization
- System design

### Approach
1. Always consider scalability
2. Follow SOLID principles
3. Write well-documented code

### Tools
- Code analysis and refactoring
- Architecture design
- Performance optimization
```

**Available Profiles:**

```
src/profiles/
├── agents/
│   ├── engineer.md          # Software engineer (default)
│   ├── architect.md         # System architect
│   ├── designer.md          # Product/UX designer
│   ├── team-lead.md         # Team lead
│   └── README.md            # Profile documentation
│
├── skills/
│   ├── typescript.md        # TypeScript best practices
│   ├── golang.md            # Go best practices
│   ├── python.md            # Python best practices
│   ├── react.md             # React patterns
│   └── README.md            # Skill documentation
│
├── commands/
│   ├── review-code.md       # Code review instructions
│   ├── plan.md              # Planning instructions
│   ├── debug.md             # Debugging instructions
│   └── README.md
│
└── README.md                # Profile guide
```

**Key Methods:**

```go
func (s *Service) LoadProfile(name string) (string, error) {
    // 1. Find profile file
    path := filepath.Join(s.profilesDir, "agents", name+".md")
    if _, err := s.fs.Stat(path); err != nil {
        return "", fmt.Errorf("profile not found: %s", name)
    }
    
    // 2. Read content
    content, err := afero.ReadFile(s.fs, path)
    if err != nil {
        return "", err
    }
    
    return string(content), nil
}

func (s *Service) ListProfiles() ([]string, error) {
    profiles := []string{}
    
    // List all .md files in agents/
    entries, err := afero.ReadDir(s.fs, filepath.Join(s.profilesDir, "agents"))
    if err != nil {
        return nil, err
    }
    
    for _, entry := range entries {
        if strings.HasSuffix(entry.Name(), ".md") {
            name := strings.TrimSuffix(entry.Name(), ".md")
            profiles = append(profiles, name)
        }
    }
    
    sort.Strings(profiles)
    return profiles, nil
}

func (s *Service) ComposeProfile(agentName, skillName string) (string, error) {
    // 1. Load agent profile
    agent, err := s.LoadProfile(agentName)
    if err != nil {
        return "", err
    }
    
    // 2. Load skill profile
    skillPath := filepath.Join(s.profilesDir, "skills", skillName+".md")
    skillContent, err := afero.ReadFile(s.fs, skillPath)
    if err != nil {
        return "", fmt.Errorf("skill not found: %s", skillName)
    }
    
    // 3. Compose (agent + skill)
    composed := fmt.Sprintf("%s\n\n## Additional Guidance\n\n%s", 
        agent, string(skillContent))
    
    return composed, nil
}
```

**Profile Injection into Claude:**

```bash
# In pre-tool-use hook:
export CLAUDE_SYSTEM_PROMPT=$(cat ~/.claudex/profiles/agents/engineer.md)
export CLAUDE_SYSTEM_PROMPT="${CLAUDE_SYSTEM_PROMPT}

## Technology Stack
$(cat ~/.claudex/profiles/skills/typescript.md)"

# Then execute Claude with expanded environment
```

---

### 1.5 Hook Service (hooks.Service)

**Location:** `src/internal/hooks/`

**Core Responsibility:** Hook discovery, execution, and logging

**Hook Execution Points:**

```
.claude/hooks/
├── pre-tool-use.sh              # Before each tool execution
├── pre-tool-use.ps1             # Windows version
├── post-tool-use.sh             # After each tool execution
├── post-tool-use.ps1            # Windows version
├── session-end.sh               # When session terminates
├── session-end.ps1              # Windows version
├── notification.sh              # On significant events
└── notification.ps1             # Windows version
```

**Hook Execution Model:**

```go
func (h *Service) ExecutePre(ctx context.Context, toolName string, input string) error {
    // 1. Discover hook script (platform-specific)
    hookPath := h.discoverHook("pre-tool-use")
    if hookPath == "" {
        return nil  // No hook, continue
    }
    
    // 2. Setup environment
    env := os.Environ()
    env = append(env, "TOOL_NAME="+toolName)
    env = append(env, "TOOL_INPUT="+input)
    env = append(env, "SESSION_ID="+h.sessionID)
    env = append(env, "SESSION_PATH="+h.sessionPath)
    
    // 3. Execute hook with timeout
    cmd := exec.CommandContext(ctx, hookPath)
    cmd.Env = env
    cmd.Stdout = h.logFile
    cmd.Stderr = h.logFile
    
    if err := cmd.Run(); err != nil {
        // Log but don't crash - hooks are non-critical
        h.logger.Printf("Pre-tool hook failed: %v", err)
        return nil
    }
    
    return nil
}

func (h *Service) ExecutePost(ctx context.Context, toolName, output string) error {
    // 1. Discover hook script
    hookPath := h.discoverHook("post-tool-use")
    if hookPath == "" {
        return nil
    }
    
    // 2. Setup environment (similar to pre-tool-use)
    env := os.Environ()
    env = append(env, "TOOL_NAME="+toolName)
    env = append(env, "TOOL_OUTPUT="+output)
    env = append(env, "SESSION_ID="+h.sessionID)
    
    // 3. Execute hook
    cmd := exec.CommandContext(ctx, hookPath)
    cmd.Env = env
    cmd.Stdout = h.logFile
    cmd.Stderr = h.logFile
    
    if err := cmd.Run(); err != nil {
        h.logger.Printf("Post-tool hook failed: %v", err)
        return nil
    }
    
    return nil
}

func (h *Service) discoverHook(name string) string {
    // 1. Check for custom hook in session
    sessionHook := filepath.Join(h.sessionPath, ".claude", "hooks", name+h.scriptExt)
    if _, err := os.Stat(sessionHook); err == nil {
        return sessionHook
    }
    
    // 2. Check for user hook in ~/.claude
    userHook := filepath.Join(os.ExpandEnv("$HOME"), ".claude", "hooks", name+h.scriptExt)
    if _, err := os.Stat(userHook); err == nil {
        return userHook
    }
    
    // 3. Check for built-in hook
    builtinHook := filepath.Join(h.builtinDir, name+h.scriptExt)
    if _, err := os.Stat(builtinHook); err == nil {
        return builtinHook
    }
    
    return ""
}
```

**Platform-Specific Hook Scripts:**

**Windows (pre-tool-use.ps1):**
```powershell
# Pre-tool-use hook for Windows
param(
    [string]$ToolName = $env:TOOL_NAME,
    [string]$ToolInput = $env:TOOL_INPUT,
    [string]$SessionId = $env:SESSION_ID,
    [string]$SessionPath = $env:SESSION_PATH
)

# Log the tool call
Add-Content -Path "$SessionPath/hook.log" -Value "$(Get-Date): Executing $ToolName with input: $ToolInput"

# Validate tool
if ($ToolName -eq "file_editor") {
    # Validate file paths in input
    Write-Host "Validating file editor input..."
}

exit 0
```

**Unix (pre-tool-use.sh):**
```bash
#!/bin/bash
# Pre-tool-use hook for Unix

TOOL_NAME="${TOOL_NAME:-}"
TOOL_INPUT="${TOOL_INPUT:-}"
SESSION_ID="${SESSION_ID:-}"
SESSION_PATH="${SESSION_PATH:-}"

# Log the tool call
echo "$(date): Executing $TOOL_NAME with input: $TOOL_INPUT" >> "$SESSION_PATH/hook.log"

# Validate tool
if [ "$TOOL_NAME" = "file_editor" ]; then
    echo "Validating file editor input..." >&2
fi

exit 0
```

---

### 1.6 Git Service (git.Service)

**Location:** `src/internal/services/git/git.go`

**Core Responsibility:** Git operations (commit SHA, file changes, merge base)

**Key Methods:**

```go
func (g *Service) GetCurrentCommitSHA(projectDir string) (string, error) {
    // Execute: git rev-parse HEAD
    cmd := exec.Command("git", "rev-parse", "HEAD")
    cmd.Dir = projectDir
    
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("failed to get commit SHA: %w", err)
    }
    
    return strings.TrimSpace(string(output)), nil
}

func (g *Service) GetChangedFiles(projectDir string) ([]string, error) {
    // Execute: git diff --name-only HEAD
    cmd := exec.Command("git", "diff", "--name-only", "HEAD")
    cmd.Dir = projectDir
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    files := strings.Split(strings.TrimSpace(string(output)), "\n")
    return files, nil
}

func (g *Service) GetMergeBase(projectDir, branch string) (string, error) {
    // Execute: git merge-base HEAD [branch]
    cmd := exec.Command("git", "merge-base", "HEAD", branch)
    cmd.Dir = projectDir
    
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    
    return strings.TrimSpace(string(output)), nil
}

func (g *Service) GetCommitLog(projectDir, ref string, limit int) ([]Commit, error) {
    // Execute: git log --format=%H%n%an%n%ae%n%ai%n%B [ref] -n [limit]
    format := "%H%n%an%n%ae%n%ai%n%B%n---END_COMMIT---%n"
    cmd := exec.Command("git", "log", fmt.Sprintf("--format=%s", format), 
        ref, "-n", fmt.Sprint(limit))
    cmd.Dir = projectDir
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    // Parse output into Commit structs
    commits := parseCommitOutput(string(output))
    return commits, nil
}

type Commit struct {
    SHA     string
    Author  string
    Email   string
    Date    time.Time
    Message string
}
```

---

### 1.7 Filesystem Service (filesystem.Service)

**Location:** `src/internal/services/filesystem/filesystem.go`

**Core Responsibility:** Abstract filesystem operations using afero interface

**Key Methods:**

```go
func (f *Service) EnsureDir(path string) error {
    return f.fs.MkdirAll(path, 0755)
}

func (f *Service) WriteFile(path, content string) error {
    return afero.WriteFile(f.fs, path, []byte(content), 0644)
}

func (f *Service) ReadFile(path string) (string, error) {
    content, err := afero.ReadFile(f.fs, path)
    return string(content), err
}

func (f *Service) SearchFiles(dir, pattern string) ([]string, error) {
    // Walk directory and match pattern
    matches := []string{}
    
    err := afero.Walk(f.fs, dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if match, _ := filepath.Match(pattern, filepath.Base(path)); match {
            matches = append(matches, path)
        }
        
        return nil
    })
    
    return matches, err
}

func (f *Service) ListDir(dir string) ([]os.FileInfo, error) {
    return afero.ReadDir(f.fs, dir)
}
```

**Testability Benefit:** All filesystem operations go through `afero.Fs` interface, enabling mock filesystem in tests.

---

## 2. Use Case Implementations

### 2.1 Create New Session Use Case

**Location:** `src/internal/usecases/session/new/`

**Class Diagram:**

```
CreateSessionUseCase
├── deps: *Dependencies
├── Execute(name string): (*SessionInfo, error)
├── validateName(name string): error
├── createSessionDir(name string): (string, error)
└── setupHooks(sessionPath string): error
```

**Implementation:**

```go
type CreateSessionUC struct {
    deps *Dependencies
}

func (uc *CreateSessionUC) Execute(name string) (*SessionInfo, error) {
    // 1. Validate name
    if err := uc.validateName(name); err != nil {
        return nil, err
    }
    
    // 2. Generate session ID
    sessionID := uc.deps.UUID.Generate()
    sessionPath := filepath.Join(uc.deps.SessionsDir, sessionID)
    
    // 3. Create session directory structure
    if err := uc.createSessionDir(sessionPath); err != nil {
        return nil, fmt.Errorf("failed to create session dir: %w", err)
    }
    
    // 4. Setup hooks
    if err := uc.setupHooks(sessionPath); err != nil {
        return nil, fmt.Errorf("failed to setup hooks: %w", err)
    }
    
    // 5. Create and return SessionInfo
    return &SessionInfo{
        Name:     name,
        Path:     sessionPath,
        Mode:     LaunchModeNew,
        ClaudeID: "",  // Will be set after Claude launches
    }, nil
}

func (uc *CreateSessionUC) createSessionDir(sessionPath string) error {
    // Create: .claudex/sessions/[id]/
    if err := uc.deps.FS.MkdirAll(sessionPath, 0755); err != nil {
        return err
    }
    
    // Create metadata files
    now := uc.deps.Clock.Now()
    metadata := map[string][]byte{
        ".created":  []byte(now.Format(time.RFC3339)),
        ".description": []byte("New session"),
        ".last_used": []byte(now.Format(time.RFC3339)),
    }
    
    for name, content := range metadata {
        path := filepath.Join(sessionPath, name)
        if err := afero.WriteFile(uc.deps.FS, path, content, 0644); err != nil {
            return err
        }
    }
    
    // Create session-overview.md
    overview := fmt.Sprintf("# Session: %s\n\nCreated: %s\n\n## Activities\n\n", 
        filepath.Base(sessionPath), now.Format("2006-01-02 15:04:05"))
    
    overviewPath := filepath.Join(sessionPath, "session-overview.md")
    return afero.WriteFile(uc.deps.FS, overviewPath, []byte(overview), 0644)
}

func (uc *CreateSessionUC) setupHooks(sessionPath string) error {
    // Copy built-in hooks to .claude/hooks/
    hooksDir := filepath.Join(sessionPath, ".claude", "hooks")
    if err := uc.deps.FS.MkdirAll(hooksDir, 0755); err != nil {
        return err
    }
    
    hookFiles := []string{"pre-tool-use", "post-tool-use", "session-end"}
    ext := ".sh"
    if runtime.GOOS == "windows" {
        ext = ".ps1"
    }
    
    for _, hook := range hookFiles {
        src := filepath.Join(uc.deps.BuiltinHooksDir, hook+ext)
        dst := filepath.Join(hooksDir, hook+ext)
        
        content, err := afero.ReadFile(uc.deps.FS, src)
        if err != nil {
            continue  // Skip missing built-in hooks
        }
        
        if err := afero.WriteFile(uc.deps.FS, dst, content, 0755); err != nil {
            return err
        }
    }
    
    return nil
}
```

---

### 2.2 Resume Session Use Case

**Location:** `src/internal/usecases/session/resume/`

**Implementation:**

```go
type ResumeSessionUC struct {
    deps *Dependencies
}

func (uc *ResumeSessionUC) Execute(sessionID string) (*SessionInfo, error) {
    sessionPath := filepath.Join(uc.deps.SessionsDir, sessionID)
    
    // 1. Verify session exists
    if _, err := uc.deps.FS.Stat(sessionPath); err != nil {
        return nil, fmt.Errorf("session not found: %s", sessionID)
    }
    
    // 2. Load session metadata
    info, err := uc.loadSessionInfo(sessionPath)
    if err != nil {
        return nil, err
    }
    
    // 3. Update last-used timestamp
    if err := uc.deps.SessionService.UpdateLastUsed(uc.deps.FS, uc.deps.Clock, sessionPath); err != nil {
        return nil, fmt.Errorf("failed to update last-used: %w", err)
    }
    
    // 4. Load session context
    context, err := uc.loadSessionContext(sessionPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load context: %w", err)
    }
    
    info.Context = context
    return info, nil
}

func (uc *ResumeSessionUC) loadSessionContext(sessionPath string) (string, error) {
    // Read session-overview.md for context
    overviewPath := filepath.Join(sessionPath, "session-overview.md")
    
    content, err := afero.ReadFile(uc.deps.FS, overviewPath)
    if err != nil {
        return "", err
    }
    
    return string(content), nil
}
```

**Context Injection:**

```bash
# In main Claude launcher:
export CLAUDE_CONTEXT=$(cat $SESSION_PATH/session-overview.md)
export CLAUDE_SESSION_ID=$SESSION_ID

exec claude code \
    --project-dir $PROJECT_DIR \
    --inject-context "$CLAUDE_CONTEXT" \
    --session-id "$CLAUDE_SESSION_ID"
```

---

### 2.3 Fork Session Use Case

**Location:** `src/internal/usecases/session/fork/`

**Implementation:**

```go
type ForkSessionUC struct {
    deps *Dependencies
}

func (uc *ForkSessionUC) Execute(sourceSessionID, newName string) (*SessionInfo, error) {
    sourceSessionPath := filepath.Join(uc.deps.SessionsDir, sourceSessionID)
    
    // 1. Verify source session exists
    if _, err := uc.deps.FS.Stat(sourceSessionPath); err != nil {
        return nil, fmt.Errorf("source session not found: %s", sourceSessionID)
    }
    
    // 2. Generate new session ID
    newSessionID := uc.deps.UUID.Generate()
    newSessionPath := filepath.Join(uc.deps.SessionsDir, newSessionID)
    
    // 3. Copy entire session directory
    if err := uc.copySessionDir(sourceSessionPath, newSessionPath); err != nil {
        return nil, fmt.Errorf("failed to copy session: %w", err)
    }
    
    // 4. Update metadata to indicate fork
    forkMetadata := fmt.Sprintf("Forked from: %s\nForked at: %s\n", 
        sourceSessionID, uc.deps.Clock.Now().Format(time.RFC3339))
    
    forkPath := filepath.Join(newSessionPath, ".fork_source")
    if err := afero.WriteFile(uc.deps.FS, forkPath, []byte(forkMetadata), 0644); err != nil {
        return nil, err
    }
    
    // 5. Create new SessionInfo
    return &SessionInfo{
        Name:         newName,
        Path:         newSessionPath,
        Mode:         LaunchModeFork,
        OriginalName: sourceSessionID,
    }, nil
}

func (uc *ForkSessionUC) copySessionDir(src, dst string) error {
    // Recursively copy directory using afero
    return afero.Walk(uc.deps.FS, src, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        relPath, _ := filepath.Rel(src, path)
        dstPath := filepath.Join(dst, relPath)
        
        if info.IsDir() {
            return uc.deps.FS.MkdirAll(dstPath, info.Mode())
        }
        
        content, _ := afero.ReadFile(uc.deps.FS, path)
        return afero.WriteFile(uc.deps.FS, dstPath, content, info.Mode())
    })
}
```

---

### 2.4 Setup MCP Use Case

**Location:** `src/internal/usecases/setupmcp/`

**Implementation:**

```go
type SetupMCPUC struct {
    deps *Dependencies
}

func (uc *SetupMCPUC) Execute() error {
    // 1. Detect available MCP servers
    available := uc.detectAvailableMCPServers()
    
    // 2. Prompt user for each
    selected := uc.promptUserForMCPSelection(available)
    
    // 3. Generate ~/.claude.json config
    config := uc.generateMCPConfig(selected)
    
    // 4. Write to ~/.claude.json
    claudePath := filepath.Join(os.ExpandEnv("$HOME"), ".claude", "claude.json")
    configJSON, _ := json.MarshalIndent(config, "", "  ")
    
    if err := afero.WriteFile(uc.deps.FS, claudePath, configJSON, 0600); err != nil {
        return err
    }
    
    // 5. Verify connection
    return uc.verifyMCPConnection(selected)
}

func (uc *SetupMCPUC) detectAvailableMCPServers() []string {
    servers := []string{}
    
    // Check for sequential-thinking (Python)
    if _, err := exec.LookPath("python3"); err == nil {
        servers = append(servers, "sequential-thinking")
    }
    
    // Check for context7 (Node.js)
    if _, err := exec.LookPath("node"); err == nil {
        servers = append(servers, "context7")
    }
    
    return servers
}

func (uc *SetupMCPUC) generateMCPConfig(selected []string) map[string]interface{} {
    config := map[string]interface{}{
        "mcp_servers": map[string]interface{}{},
    }
    
    servers := config["mcp_servers"].(map[string]interface{})
    
    for _, name := range selected {
        switch name {
        case "sequential-thinking":
            servers["sequential-thinking"] = map[string]interface{}{
                "command": "python3 -m anthropic_mcp_server sequential_thinking",
            }
        case "context7":
            servers["context7"] = map[string]interface{}{
                "command": "node /opt/context7/index.js",
                "env": map[string]string{
                    "CONTEXT7_PORT": "8000",
                },
            }
        }
    }
    
    return config
}
```

---

## 3. Sequence Diagrams

### 3.1 Session Creation Workflow

```
User              App           SessionUC       SessionService    FS
  │                 │                 │                │           │
  ├─ claudex new ──>│                 │                │           │
  │                 │                 │                │           │
  │                 ├─ Init() ───────>│                │           │
  │                 │                 │                │           │
  │                 ├─ Execute() ────>│                │           │
  │                 │                 │                           │
  │                 │                 ├─ createSessionDir() ─────>│
  │                 │                 │                   │       │
  │                 │                 │                   │<─ OK ─┤
  │                 │                 │                           │
  │                 │                 ├─ setupHooks() ──────────>│
  │                 │                 │                   │       │
  │                 │                 │                   │<─ OK ─┤
  │                 │                 │                           │
  │                 │<─ SessionInfo ──┤                          │
  │                 │                 │                           │
  │                 ├─ LaunchClaude() │                           │
  │                 │  (exec: claude code ...) → Claude Running
  │<─ Session Started
  │
```

### 3.2 Hook Execution Workflow

```
Claude              Pre-Hook       Tool             Post-Hook        Session
  │                    │            │                  │              │
  ├─ Pre-Tool ────────>│            │                  │              │
  │  (environment)     │            │                  │              │
  │                    │            │                  │              │
  │                    ├─ Validate ─┤                  │              │
  │                    │            │                  │              │
  │                    └─ Pass ─────────────────────>  │              │
  │                                 │                  │              │
  │<────── Context ────────────────┤                  │              │
  │                                 │                  │              │
  │ Tool Executes                   │                  │              │
  │                                 │                  │              │
  ├─ Post-Tool ────────────────────────────────────> │              │
  │  (output)                                         │              │
  │                                                   │              │
  │                                  ├─ Parse ───────>│              │
  │                                  │                 │              │
  │                                  ├─ Update ──────────────────> │
  │                                  │    (session-overview.md)     │
  │                                  │                 │<── OK ─────┤
  │                                  │                 │              │
  │<──────────── Continue ──────────┤                 │
  │
```

---

## 4. Error Handling Strategy

**Error Hierarchy:**

```
error
├── ErrSessionNotFound
│   └── Used when session directory doesn't exist
├── ErrInvalidSessionName
│   └── Used for invalid session naming
├── ErrConfigurationError
│   └── Used for TOML parsing failures
├── ErrHookExecutionFailed
│   └── Non-blocking, logged but continues
├── ErrGitOperationFailed
│   └── Used for git command failures
└── ErrFileSystemError
    └── Used for afero operations
```

**Error Handling in Use Cases:**

```go
func (uc *CreateSessionUC) Execute(name string) (*SessionInfo, error) {
    // Validation errors: return immediately
    if err := uc.validateName(name); err != nil {
        return nil, fmt.Errorf("invalid session name: %w", err)
    }
    
    // Critical errors: return immediately
    if err := uc.createSessionDir(sessionPath); err != nil {
        return nil, fmt.Errorf("failed to create session directory: %w", err)
    }
    
    // Non-critical errors: log but continue
    if err := uc.setupHooks(sessionPath); err != nil {
        uc.deps.Logger.Printf("Warning: failed to setup hooks: %v", err)
        // Continue - hooks are optional
    }
    
    return info, nil
}
```

---

## 5. Testing Strategy

### 5.1 Unit Test Examples

**Session Service Test:**

```go
func TestGetSessions(t *testing.T) {
    // 1. Create mock filesystem
    fs := afero.NewMemMapFs()
    
    // 2. Create test session directories
    fs.MkdirAll(".claudex/sessions/session-1", 0755)
    afero.WriteFile(fs, ".claudex/sessions/session-1/.created", 
        []byte(time.Now().Format(time.RFC3339)), 0644)
    afero.WriteFile(fs, ".claudex/sessions/session-1/.description", 
        []byte("Test session"), 0644)
    
    // 3. Create service with mock FS
    svc := session.New(fs)
    
    // 4. Test GetSessions
    sessions, err := svc.GetSessions(fs, ".claudex/sessions")
    
    // 5. Assert
    if err != nil {
        t.Fatalf("GetSessions failed: %v", err)
    }
    if len(sessions) != 1 {
        t.Errorf("Expected 1 session, got %d", len(sessions))
    }
}
```

**Hook Execution Test:**

```go
func TestExecutePre(t *testing.T) {
    // 1. Create mock hook file
    fs := afero.NewMemMapFs()
    fs.MkdirAll(".claude/hooks", 0755)
    
    hookContent := `#!/bin/bash
echo "TOOL_NAME=$TOOL_NAME" > /tmp/hook_output.txt
`
    afero.WriteFile(fs, ".claude/hooks/pre-tool-use.sh", []byte(hookContent), 0755)
    
    // 2. Create hook service
    svc := hooks.New(fs, ".claude/hooks")
    
    // 3. Execute pre-hook
    err := svc.ExecutePre(context.Background(), "file_editor", "input.txt")
    
    // 4. Assert
    if err != nil {
        t.Fatalf("ExecutePre failed: %v", err)
    }
}
```

---

## 6. Sign-Off

| Role | Name | Date | Approval |
|------|------|------|----------|
| Architecture | [TBD] | TBD | [ ] Approved |
| Development Lead | [TBD] | TBD | [ ] Approved |

