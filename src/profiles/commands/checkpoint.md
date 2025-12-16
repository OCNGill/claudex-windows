# Checkpoint Command

Update the current session overview file with all findings, progress, and current status. This command helps maintain context across conversation resets and allows you to resume work seamlessly.

## Step 1: Gather Current State

Run these commands to understand the current project state:

```bash
# Check git status
git status --short

# Show current branch
git branch --show-current

# Show uncommitted changes
git diff --stat

# Check for stashed changes
git stash list

# Show recent commits
git log --oneline -5
```

## Step 2: Identify Session File Location

The session overview file is located at:
```
.claudex/sessions/{session-id}/session-overview.md
```

Determine the current session ID from the session context provided at conversation start.

## Step 3: Update Session Overview

Update the session overview file with the following sections:

### Session Summary
- Keep the original goal/purpose of the session
- Add any clarifications or scope changes that occurred

### Current Focus
- Update with what you're currently working on
- If task is complete, describe what was just finished
- If blocked, describe the blocker clearly

### Key Documents
- Add any new documentation files created
- Add any modified configuration files
- Include execution plans, architecture docs, test reports
- Use absolute paths for all file references

### Progress Timeline
- Add a new timestamped entry describing:
  - What was accomplished
  - What decisions were made
  - What files were modified/created
  - Any blockers encountered
  - Next planned steps
- Format: `- **YYYY-MM-DDTHH:MM:SSZ** - Description of progress`

### Status Field
Update the status to one of:
- `In Progress` - Active work ongoing
- `Blocked` - Waiting on external input or resolution
- `Completed` - All tasks finished
- `Paused` - Work temporarily suspended

### Last Updated Timestamp
- Update the footer timestamp
- Include a brief description of the update type
- Format: `*Last updated: YYYY-MM-DDTHH:MM:SSZ (Checkpoint: brief description)*`

## Step 4: Update Todo List

If a todo list exists in the conversation:
- Mark completed tasks with [x]
- Add any new tasks discovered during work
- Update task descriptions if scope changed
- Remove tasks that are no longer relevant

## Step 5: Capture Key Context

Include in the progress timeline entry:

**Git State:**
- Current branch
- Modified files count
- Uncommitted changes summary
- Stashed changes (if any)

**Work Completed:**
- Files created or modified
- Tests written or fixed
- Documentation updated
- Features implemented or bugs fixed

**Current Blockers:**
- Any issues preventing progress
- Questions that need answering
- Dependencies waiting on external input

**Next Steps:**
- Immediate next actions
- Planned tasks in priority order
- Any follow-up items

## Step 6: Confirm Update

After updating the session overview, confirm with a summary:

```
Checkpoint saved to session overview!

Updated:
- Status: [current status]
- Current Focus: [brief description]
- Progress Timeline: Added entry at [timestamp]
- Key Documents: [number] files tracked
- Todo List: [X] completed, [Y] in progress, [Z] pending

Git State:
- Branch: [branch name]
- Modified: [number] files
- Uncommitted changes: [yes/no]

Next Steps:
- [step 1]
- [step 2]
```

## Important Notes

- This checkpoint captures point-in-time state for session resumption
- The session overview is the single source of truth for session context
- Always use absolute paths when referencing files
- Timestamps should use ISO 8601 format (UTC)
- Be specific in progress descriptions to help future resumption
- Include both successes and blockers for complete context

## When to Use This Command

Use `/checkpoint` at these key moments:
- After completing a significant milestone
- Before taking a break or ending a work session
- When switching focus to a different task
- After encountering a blocker
- Before making major architectural decisions
- When tests pass after bug fixes
- After updating documentation

Regular checkpoints ensure no context is lost and make session forking more reliable.
