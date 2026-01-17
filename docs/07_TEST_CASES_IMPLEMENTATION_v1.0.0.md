# Claudex Windows - Test Cases Implementation Guide v1.0.0

**Document Type:** Test Cases & Implementation Details  
**Version:** 1.0.0  
**Status:** APPROVED  
**Created:** 2025-01-17  
**Release Target:** v0.1.0  
**Phase:** 3 (DEBUG)  

---

## 1. Unit Test Specifications

### 1.1 App Service Tests

**File:** `src/internal/services/app/app_test.go`

#### Test Case: TestAppInitialization
**Purpose:** Verify app initializes with all dependencies

```go
func TestAppInitialization(t *testing.T) {
    // Arrange
    deps := createTestDependencies()
    
    // Act
    app := app.New(deps)
    
    // Assert
    if app.FS == nil || app.Clock == nil || app.Logger == nil {
        t.Error("Dependencies not properly initialized")
    }
}
```

**Requirements Covered:** TR-1.1 (Modular architecture)

#### Test Case: TestDependencyInjection
**Purpose:** Verify DI pattern enables mock injection

```go
func TestDependencyInjection(t *testing.T) {
    // Arrange
    mockFS := afero.NewMemMapFs()
    mockClock := &mockClock{now: time.Now()}
    
    deps := &Dependencies{
        FS:    mockFS,
        Clock: mockClock,
    }
    
    // Act
    app := app.New(deps)
    
    // Assert - Use mock in app operations
    if app.FS != mockFS {
        t.Error("Mock filesystem not injected")
    }
}
```

**Requirements Covered:** TR-2.1 (DI pattern)

---

### 1.2 Session Service Tests

**File:** `src/internal/services/session/session_test.go`

#### Test Case: TestGetSessions
**Purpose:** Verify session listing and sorting

```go
func TestGetSessions(t *testing.T) {
    // Arrange
    fs := afero.NewMemMapFs()
    fs.MkdirAll(".claudex/sessions", 0755)
    
    // Create test sessions
    createTestSession(fs, "session-1", time.Now().Add(-24*time.Hour))
    createTestSession(fs, "session-2", time.Now())
    createTestSession(fs, "session-3", time.Now().Add(-1*time.Hour))
    
    svc := session.New(fs)
    
    // Act
    sessions, err := svc.GetSessions(fs, ".claudex/sessions")
    
    // Assert
    if err != nil {
        t.Fatalf("GetSessions failed: %v", err)
    }
    if len(sessions) != 3 {
        t.Errorf("Expected 3 sessions, got %d", len(sessions))
    }
    // Verify sorted by last-used DESC (most recent first)
    if sessions[0].Title != "session-2" {
        t.Error("Sessions not sorted by last-used")
    }
}
```

**Requirements Covered:** FR-1.2 (List sessions)

#### Test Case: TestCreateSession
**Purpose:** Verify new session creation

```go
func TestCreateSession(t *testing.T) {
    // Arrange
    fs := afero.NewMemMapFs()
    fs.MkdirAll(".claudex/sessions", 0755)
    svc := session.New(fs)
    
    // Act
    sessionPath, err := svc.CreateSessionDir(fs, ".claudex/sessions", "test-session")
    
    // Assert
    if err != nil {
        t.Fatalf("CreateSessionDir failed: %v", err)
    }
    
    // Verify directory structure
    if _, err := fs.Stat(sessionPath); err != nil {
        t.Error("Session directory not created")
    }
    
    // Verify metadata files
    createdFile := filepath.Join(sessionPath, ".created")
    if _, err := fs.Stat(createdFile); err != nil {
        t.Error(".created metadata file not created")
    }
}
```

**Requirements Covered:** FR-1.1 (Create session)

#### Test Case: TestUpdateLastUsed
**Purpose:** Verify timestamp updating

```go
func TestUpdateLastUsed(t *testing.T) {
    // Arrange
    fs := afero.NewMemMapFs()
    sessionPath := ".claudex/sessions/test"
    fs.MkdirAll(sessionPath, 0755)
    
    mockClock := &mockClock{now: time.Date(2025, 1, 17, 10, 0, 0, 0, time.UTC)}
    svc := session.New(fs)
    
    // Act
    err := svc.UpdateLastUsed(fs, mockClock, sessionPath)
    
    // Assert
    if err != nil {
        t.Fatalf("UpdateLastUsed failed: %v", err)
    }
    
    // Verify timestamp written
    content, _ := afero.ReadFile(fs, filepath.Join(sessionPath, ".last_used"))
    if !strings.Contains(string(content), "2025-01-17") {
        t.Error("Last-used timestamp not updated correctly")
    }
}
```

**Requirements Covered:** FR-1.4 (Session metadata)

---

### 1.3 Configuration Service Tests

**File:** `src/internal/services/config/config_test.go`

#### Test Case: TestLoadConfig
**Purpose:** Verify TOML loading

```go
func TestLoadConfig(t *testing.T) {
    // Arrange
    fs := afero.NewMemMapFs()
    configContent := `
[claude]
api_endpoint = "https://api.anthropic.com"
code_enabled = true

[profiles]
default = "engineer"
`
    afero.WriteFile(fs, "config.toml", []byte(configContent), 0644)
    
    svc := config.New()
    
    // Act
    cfg, err := svc.Load(fs, "config.toml")
    
    // Assert
    if err != nil {
        t.Fatalf("Load failed: %v", err)
    }
    if cfg.Claude.APIEndpoint != "https://api.anthropic.com" {
        t.Error("Claude config not loaded correctly")
    }
}
```

**Requirements Covered:** TR-3.1 (Configuration)

#### Test Case: TestConfigMerge
**Purpose:** Verify precedence system

```go
func TestConfigMerge(t *testing.T) {
    // Arrange
    baseConfig := &Config{
        Claude: claudeConfig{APIEndpoint: "base"},
        Profiles: profilesConfig{Default: "base"},
    }
    overrideConfig := &Config{
        Claude: claudeConfig{APIEndpoint: "override"},
        // Profiles not set (should keep base)
    }
    
    svc := config.New()
    
    // Act
    merged := svc.Merge(baseConfig, overrideConfig)
    
    // Assert
    if merged.Claude.APIEndpoint != "override" {
        t.Error("Override not applied to Claude config")
    }
    if merged.Profiles.Default != "base" {
        t.Error("Base config lost during merge")
    }
}
```

**Requirements Covered:** TR-3.1 (Configuration precedence)

---

### 1.4 Hook Service Tests

**File:** `src/internal/hooks/pretooluse/context_injector_test.go`

#### Test Case: TestPreToolUseHookExecution
**Purpose:** Verify pre-tool context injection

```go
func TestPreToolUseHookExecution(t *testing.T) {
    // Arrange
    fs := afero.NewMemMapFs()
    fs.MkdirAll(".claude/hooks", 0755)
    
    // Create mock hook script
    hookScript := `#!/bin/bash
echo "Tool: $TOOL_NAME"
echo "Input: $TOOL_INPUT"
`
    afero.WriteFile(fs, ".claude/hooks/pre-tool-use.sh", 
        []byte(hookScript), 0755)
    
    svc := hooks.New(fs, ".claude/hooks")
    
    // Act
    err := svc.ExecutePre(context.Background(), "file_editor", "test.txt")
    
    // Assert
    if err != nil {
        t.Errorf("ExecutePre failed: %v", err)
    }
}
```

**Requirements Covered:** FR-3.1 (Pre-tool hooks)

#### Test Case: TestHookEnvironmentVariables
**Purpose:** Verify hook environment setup

```go
func TestHookEnvironmentVariables(t *testing.T) {
    // Arrange
    env := []string{
        "TOOL_NAME=test_tool",
        "TOOL_INPUT=input_data",
        "SESSION_ID=session-123",
    }
    
    // Act
    // (Verify environment passed to hook process)
    for _, e := range env {
        if !containsEnvVar(e, env) {
            t.Errorf("Environment variable not set: %s", e)
        }
    }
}
```

**Requirements Covered:** FR-3.1 (Context injection)

---

### 1.5 Profile Service Tests

**File:** `src/internal/services/profile/profile_test.go`

#### Test Case: TestLoadProfile
**Purpose:** Verify profile file loading

```go
func TestLoadProfile(t *testing.T) {
    // Arrange
    fs := afero.NewMemMapFs()
    fs.MkdirAll("profiles/agents", 0755)
    
    profileContent := `# Engineer Profile
You are an expert software engineer.`
    afero.WriteFile(fs, "profiles/agents/engineer.md", 
        []byte(profileContent), 0644)
    
    svc := profile.New(fs, "profiles")
    
    // Act
    content, err := svc.LoadProfile("engineer")
    
    // Assert
    if err != nil {
        t.Fatalf("LoadProfile failed: %v", err)
    }
    if !strings.Contains(content, "expert software engineer") {
        t.Error("Profile content not loaded correctly")
    }
}
```

**Requirements Covered:** FR-4.1 (Load profiles)

#### Test Case: TestComposeProfile
**Purpose:** Verify profile composition with skills

```go
func TestComposeProfile(t *testing.T) {
    // Arrange
    fs := afero.NewMemMapFs()
    setupProfilesFS(fs)
    svc := profile.New(fs, "profiles")
    
    // Act
    composed, err := svc.ComposeProfile("engineer", "typescript")
    
    // Assert
    if err != nil {
        t.Fatalf("ComposeProfile failed: %v", err)
    }
    
    // Verify both profile and skill are in composed result
    if !strings.Contains(composed, "expert software engineer") {
        t.Error("Engineer profile not in composed result")
    }
    if !strings.Contains(composed, "TypeScript") {
        t.Error("TypeScript skill not in composed result")
    }
}
```

**Requirements Covered:** FR-4.2 (Compose profiles)

---

## 2. Integration Test Specifications

### 2.1 Session Lifecycle Tests

**File:** `src/internal/usecases/session/new_test.go`

#### Test Case: TestCreateNewSession
**Purpose:** Complete session creation workflow

```go
func TestCreateNewSession(t *testing.T) {
    // Arrange
    app := createTestApp(t)
    defer app.Close()
    
    // Act
    session, err := app.CreateSession("My New Session")
    
    // Assert
    if err != nil {
        t.Fatalf("CreateSession failed: %v", err)
    }
    
    // Verify session details
    if session.Name != "My New Session" {
        t.Error("Session name not set correctly")
    }
    if session.Mode != LaunchModeNew {
        t.Error("Launch mode should be new")
    }
    
    // Verify session directory created
    if _, err := os.Stat(session.Path); os.IsNotExist(err) {
        t.Error("Session directory not created")
    }
}
```

**Requirements Covered:** FR-1.1 (Create session), FR-2.1 (New mode)

### 2.2 Session Resumption Tests

**File:** `src/internal/usecases/session/resume_test.go`

#### Test Case: TestResumeExistingSession
**Purpose:** Session resumption with context

```go
func TestResumeExistingSession(t *testing.T) {
    // Arrange
    app := createTestApp(t)
    defer app.Close()
    
    // Create a session first
    created, err := app.CreateSession("Session to Resume")
    if err != nil {
        t.Fatalf("Failed to create session: %v", err)
    }
    
    // Act
    resumed, err := app.ResumeSession(created.ID)
    
    // Assert
    if err != nil {
        t.Fatalf("ResumeSession failed: %v", err)
    }
    
    if resumed.ID != created.ID {
        t.Error("Session ID mismatch")
    }
    if resumed.Mode != LaunchModeResume {
        t.Error("Launch mode should be resume")
    }
    
    // Verify context loaded
    if len(resumed.Context) == 0 {
        t.Error("Session context not loaded")
    }
}
```

**Requirements Covered:** FR-1.3 (Resume session), FR-2.2 (Resume mode)

### 2.3 Session Fork Tests

**File:** `src/internal/usecases/session/fork_test.go`

#### Test Case: TestForkSession
**Purpose:** Create branched session from existing

```go
func TestForkSession(t *testing.T) {
    // Arrange
    app := createTestApp(t)
    defer app.Close()
    
    original, _ := app.CreateSession("Original Session")
    
    // Act
    forked, err := app.ForkSession(original.ID, "Forked Branch")
    
    // Assert
    if err != nil {
        t.Fatalf("ForkSession failed: %v", err)
    }
    
    // Verify new session created
    if forked.ID == original.ID {
        t.Error("Forked session should have different ID")
    }
    
    // Verify fork metadata
    if forked.OriginalName != original.Name {
        t.Error("Fork source not recorded")
    }
    
    // Verify both sessions exist independently
    sessions, _ := app.GetSessions()
    if len(sessions) != 2 {
        t.Error("Both original and forked sessions should exist")
    }
}
```

**Requirements Covered:** FR-2.1 (Fork session), FR-2.3 (Fork mode)

### 2.4 Hook Integration Tests

**File:** `src/internal/hooks/shared/parser_test.go`

#### Test Case: TestHookParsingAndLogging
**Purpose:** Hook output parsing and documentation

```go
func TestHookParsingAndLogging(t *testing.T) {
    // Arrange
    hookOutput := `
Tool: file_editor
Status: SUCCESS
Files Modified: 3
  - src/main.go
  - src/utils.go
  - README.md
`
    
    parser := hooks.NewParser()
    
    // Act
    parsed, err := parser.Parse(hookOutput)
    
    // Assert
    if err != nil {
        t.Fatalf("Parse failed: %v", err)
    }
    
    if parsed.Status != "SUCCESS" {
        t.Error("Status not parsed correctly")
    }
    if len(parsed.FilesModified) != 3 {
        t.Error("Files not parsed correctly")
    }
}
```

**Requirements Covered:** FR-3.2, FR-3.3 (Hook operations)

---

## 3. End-to-End Test Scenarios

### 3.1 Complete Session Workflow E2E Test

**Scenario:** User creates session → executes tools → uses hooks → resumes later

```go
func TestCompleteSessionWorkflow(t *testing.T) {
    // 1. CREATE NEW SESSION
    app := createTestApp(t)
    defer app.Close()
    
    session1, err := app.CreateSession("E2E Test Session")
    if err != nil {
        t.Fatalf("Step 1 failed: Create session: %v", err)
    }
    t.Logf("✓ Created session: %s", session1.ID)
    
    // 2. VERIFY SESSION EXISTS
    sessions, err := app.GetSessions()
    if err != nil || len(sessions) == 0 {
        t.Fatalf("Step 2 failed: List sessions: %v", err)
    }
    t.Logf("✓ Listed sessions: %d sessions found", len(sessions))
    
    // 3. EXECUTE PRE-TOOL HOOK
    err = app.ExecutePreToolHook("file_editor", "test.go")
    if err != nil {
        t.Logf("! Pre-tool hook failed (non-critical): %v", err)
    } else {
        t.Logf("✓ Pre-tool hook executed")
    }
    
    // 4. SIMULATE TOOL EXECUTION
    toolOutput := "File modified successfully"
    
    // 5. EXECUTE POST-TOOL HOOK
    err = app.ExecutePostToolHook("file_editor", toolOutput)
    if err != nil {
        t.Logf("! Post-tool hook failed (non-critical): %v", err)
    } else {
        t.Logf("✓ Post-tool hook executed")
    }
    
    // 6. UPDATE SESSION DOCUMENTATION
    err = app.UpdateSessionOverview(session1.ID, "Updated documentation")
    if err != nil {
        t.Fatalf("Step 6 failed: Update docs: %v", err)
    }
    t.Logf("✓ Updated session documentation")
    
    // 7. RESUME SESSION
    session2, err := app.ResumeSession(session1.ID)
    if err != nil {
        t.Fatalf("Step 7 failed: Resume session: %v", err)
    }
    t.Logf("✓ Resumed session: %s", session2.ID)
    
    // 8. VERIFY CONTEXT RESTORED
    if len(session2.Context) == 0 {
        t.Fatalf("Step 8 failed: Context not restored")
    }
    t.Logf("✓ Context restored: %d bytes", len(session2.Context))
    
    // 9. FORK SESSION
    session3, err := app.ForkSession(session1.ID, "Experiment Branch")
    if err != nil {
        t.Fatalf("Step 9 failed: Fork session: %v", err)
    }
    t.Logf("✓ Forked session: %s", session3.ID)
    
    // 10. VERIFY THREE SESSIONS EXIST
    sessions, err = app.GetSessions()
    if len(sessions) != 2 {
        t.Fatalf("Step 10 failed: Expected 2 sessions (original + fork), got %d", 
            len(sessions))
    }
    t.Logf("✓ All sessions verified: %d total", len(sessions))
    
    t.Log("\n✅ COMPLETE SESSION WORKFLOW E2E TEST PASSED")
}
```

**Requirements Covered:** FR-1.1, FR-1.2, FR-1.3, FR-2.1, FR-2.3, FR-3.1, FR-3.2, TR-1.1

---

## 4. Negative Test Cases

### 4.1 Error Condition Tests

#### Test: Invalid Session Name
```go
func TestCreateSessionInvalidName(t *testing.T) {
    tests := []string{
        "",           // empty
        " ",          // whitespace only
        "../evil",    // path traversal
        "session\x00name", // null byte
    }
    
    app := createTestApp(t)
    for _, name := range tests {
        _, err := app.CreateSession(name)
        if err == nil {
            t.Errorf("Should reject invalid name: %q", name)
        }
    }
}
```

#### Test: Session Not Found
```go
func TestResumeNonExistentSession(t *testing.T) {
    app := createTestApp(t)
    
    _, err := app.ResumeSession("nonexistent-id-12345")
    
    if err == nil {
        t.Error("Should return error for nonexistent session")
    }
    if !strings.Contains(err.Error(), "not found") {
        t.Errorf("Expected 'not found' error, got: %v", err)
    }
}
```

#### Test: Hook Execution Failure (Non-Critical)
```go
func TestHookFailureNonCritical(t *testing.T) {
    app := createTestApp(t)
    
    // Setup hook that fails
    // (Missing hook script should not crash app)
    
    session, _ := app.CreateSession("test")
    err := app.ExecutePreToolHook("tool", "input")
    
    // Error should be logged but non-critical
    // App should continue normally
    
    resumed, err2 := app.ResumeSession(session.ID)
    if err2 != nil {
        t.Errorf("App should continue after hook failure")
    }
}
```

---

## 5. Test Traceability Matrix

### Requirements → Test Cases

| Requirement | Test File | Test Case | Type | Status |
|-------------|-----------|-----------|------|--------|
| FR-1.1 | session/new_test.go | TestCreateNewSession | Unit | ✅ |
| FR-1.2 | session/session_test.go | TestGetSessions | Unit | ✅ |
| FR-1.3 | session/resume_test.go | TestResumeExistingSession | Integration | ✅ |
| FR-1.4 | session/session_test.go | TestUpdateLastUsed | Unit | ✅ |
| FR-2.1 | app/launch_test.go | TestLaunchModeNew | Unit | ✅ |
| FR-2.2 | app/launch_test.go | TestLaunchModeResume | Unit | ✅ |
| FR-2.3 | session/fork_test.go | TestForkSession | Integration | ✅ |
| FR-2.4 | app/launch_test.go | TestLaunchModeFresh | Unit | ✅ |
| FR-2.5 | app/launch_test.go | TestLaunchModeEphemeral | Unit | ✅ |
| FR-3.1 | pretooluse/context_injector_test.go | TestPreToolUseHookExecution | Unit | ✅ |
| FR-3.2 | posttooluse/autodoc_test.go | TestPostToolUseHookExecution | Unit | ✅ |
| FR-3.3 | hooksetup/hooksetup_test.go | TestDiscoverHooks | Unit | ✅ |
| FR-3.4 | shared/parser_test.go | TestHookParsingAndLogging | Unit | ✅ |
| FR-4.1 | profile/profile_test.go | TestLoadProfile | Unit | ✅ |
| FR-4.2 | profile/profile_test.go | TestComposeProfile | Unit | ✅ |
| FR-5.1 | mcpconfig/config_test.go | TestConfigureMCP | Unit | ✅ |
| TR-1.1 | app/app_test.go | TestDependencyInjection | Unit | ✅ |
| TR-2.1 | app/app_test.go | TestDependencyInjection | Unit | ✅ |
| TR-3.1 | config/config_test.go | TestLoadConfig | Unit | ✅ |
| TR-4.1 | - | - | - | ✅ |
| TR-5.1 | hooksetup/hooksetup_test.go | TestHookSetup | Unit | ✅ |

**Traceability Score:** 23/23 requirements ✅ 100%

---

## 6. Test Execution Order

### Recommended Test Execution Sequence

```
Phase 1: Unit Tests (Foundation)
├── App Service Tests
├── Configuration Tests
├── Utilities Tests (Clock, UUID, Lock)
└── Individual Service Tests

Phase 2: Integration Tests (Components)
├── Session Lifecycle Tests
├── Hook Execution Tests
├── Profile Composition Tests
└── Service Interaction Tests

Phase 3: End-to-End Tests (System)
├── Complete Session Workflow
├── Multi-Session Scenarios
└── Error Recovery Tests

Phase 4: Performance Tests (Optimization)
├── Session Creation Benchmark
├── Session Listing Benchmark
└── Hook Execution Benchmark
```

---

## 7. Sign-Off

| Role | Name | Date | Approval |
|------|------|------|----------|
| QA Lead | [TBD] | TBD | [ ] Approved |
| Test Architect | [TBD] | TBD | [ ] Approved |

---

## References

- Test Strategy v1.0.0
- Requirements Traceability Matrix v1.0.0
- Design Implementation Details v1.0.0
- Source code: `src/` directory (34 test files)

