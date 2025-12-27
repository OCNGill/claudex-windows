# Review Code Command

Perform a comprehensive code review of implementation changes. This command analyzes code quality, identifies optimization opportunities, suggests maintainability improvements, and highlights potential bugs or edge cases.

## Step 1: Gather Context

Determine what code to review by running these commands:

```bash
# Check git status for staged and unstaged changes
git status --short

# Show staged changes with stats
git diff --cached --stat

# Show unstaged changes with stats
git diff --stat

# Show recent commits on current branch
git log --oneline -10

# Show changes in last commit
git diff HEAD~1 --stat
```

**Ask the user to clarify the review scope:**
- Review staged changes?
- Review specific files or directories?
- Review recent commits (how many)?
- Review all uncommitted changes?
- Review specific pull request changes?

## Step 2: Read the Code

Based on the scope identified:

1. Read all relevant files using the Read tool
2. Read the actual diff content to understand what changed:

```bash
# For staged changes
git diff --cached

# For unstaged changes
git diff

# For specific files
git diff path/to/file

# For recent commits
git diff HEAD~N..HEAD
```

3. Load related test files if they exist
4. Check for documentation files related to the code

## Step 3: Code Quality & Readability Review

Analyze the code for quality and readability issues:

**Naming & Conventions:**
- Are variable, function, and type names clear and descriptive?
- Do names follow project conventions and language idioms?
- Are abbreviations avoided or well-established?

**Code Structure:**
- Is the code well-organized with clear separation of concerns?
- Are functions/methods focused on a single responsibility?
- Is the code DRY (Don't Repeat Yourself)?
- Are there overly complex functions that should be split?

**Readability:**
- Is the code easy to understand on first reading?
- Are complex logic sections commented appropriately?
- Is the formatting consistent?
- Are magic numbers/strings extracted to named constants?

**Documentation:**
- Are public APIs documented?
- Are complex algorithms explained?
- Are assumptions and edge cases documented?

## Step 4: Optimization Opportunities Review

Look for performance and efficiency improvements:

**Algorithmic Efficiency:**
- Are there O(n²) or worse algorithms that could be optimized?
- Can any loops be eliminated or reduced?
- Are there redundant operations?

**Resource Usage:**
- Are there memory leaks or unnecessary allocations?
- Are resources (files, connections, locks) properly released?
- Can caching reduce repeated expensive operations?

**Language-Specific:**
- Are language-specific performance best practices followed?
- Are appropriate data structures used?
- Are there opportunities to use built-in optimized functions?

**Premature Optimization:**
- Note if optimizations are premature and not justified by profiling
- Balance readability vs performance

## Step 5: Maintainability Improvements Review

Evaluate long-term code health:

**Testability:**
- Is the code easily testable?
- Are dependencies injected or hardcoded?
- Can functions be tested in isolation?
- Are there sufficient tests for the changes?

**Modularity:**
- Are components loosely coupled?
- Are interfaces well-defined?
- Can modules be reused or replaced easily?

**Error Handling:**
- Are errors handled appropriately?
- Are error messages helpful for debugging?
- Are edge cases handled gracefully?
- Is error recovery logic sound?

**Dependencies:**
- Are new dependencies justified?
- Are dependencies up-to-date and maintained?
- Could standard library features replace dependencies?

**Technical Debt:**
- Are there TODOs or FIXMEs that should be addressed?
- Are workarounds or hacks properly documented?
- Should any code be refactored before proceeding?

## Step 6: Potential Bugs & Edge Cases Review

Identify correctness issues and edge cases:

**Logic Errors:**
- Are conditionals correct (&&, ||, negations)?
- Are boundary conditions handled (off-by-one errors)?
- Are null/undefined/nil values handled?
- Are type conversions safe?

**Concurrency Issues:**
- Are race conditions possible?
- Is shared state properly synchronized?
- Are async operations handled correctly?

**Security Concerns:**
- Are inputs validated and sanitized?
- Are there SQL injection or XSS vulnerabilities?
- Are secrets or credentials exposed?
- Are authentication/authorization checks present?

**Edge Cases:**
- What happens with empty inputs?
- What happens with very large inputs?
- What happens with invalid/malformed data?
- Are timezone/locale issues considered?
- Are integer overflow/underflow possible?

**Resource Limits:**
- What happens when disk/memory is full?
- Are timeouts configured appropriately?
- Are retry mechanisms in place?

## Step 7: Produce Review Summary

Generate a structured review report:

```markdown
# Code Review Summary

**Reviewed:** [file paths or commit range]
**Reviewer:** Claude
**Date:** [current date]

## Overall Assessment

[Brief 2-3 sentence summary of code quality and readiness]

## Critical Issues 🔴

[Issues that MUST be fixed before merging - bugs, security issues, breaking changes]

- **[File:Line]** [Issue description]
  - Impact: [what could go wrong]
  - Recommendation: [how to fix]

## Important Improvements 🟡

[Issues that SHOULD be addressed - quality, maintainability, performance concerns]

- **[File:Line]** [Issue description]
  - Impact: [why this matters]
  - Recommendation: [suggested improvement]

## Minor Suggestions 🟢

[Optional improvements - style, documentation, minor optimizations]

- **[File:Line]** [Suggestion]
  - Benefit: [why this helps]

## Positive Highlights ✅

[What was done well - good patterns, clever solutions, excellent tests]

- [Positive observation]

## Test Coverage

- [Assessment of test completeness]
- [Missing test scenarios if any]

## Next Steps

1. [Prioritized list of actions to take]
2. [...]

## Review Statistics

- Files reviewed: [N]
- Lines changed: [+X -Y]
- Critical issues: [N]
- Important improvements: [N]
- Minor suggestions: [N]
```

## Step 8: Interactive Clarification

If you need more context to complete the review:

**Ask about:**
- Project conventions or standards
- Performance requirements or constraints
- Browser/platform support requirements
- Test coverage expectations
- Deployment environment considerations
- Known technical debt or future plans

**Offer to:**
- Deep dive into specific files
- Review related code for consistency
- Check integration points
- Review test coverage
- Suggest refactoring approaches

## Important Notes

- Reviews should be constructive and specific
- Provide examples of improvements when possible
- Consider the project's tech stack and conventions
- Balance perfectionism with pragmatism
- Highlight both problems AND good practices
- Use file:line references for precise feedback
- Categorize by severity (critical/important/minor)
- Consider the context: is this a prototype, production code, or refactoring?

## When to Use This Command

Use `/review-code` in these situations:

- Before committing significant changes
- Before creating a pull request
- After implementing a new feature
- When refactoring existing code
- After a merge to review integration
- When investigating bugs in recent changes
- Before a release to review critical paths
- When learning a new codebase (review recent quality work)
- After team feedback to validate improvements

Regular code reviews catch issues early, improve code quality, and maintain project health.
