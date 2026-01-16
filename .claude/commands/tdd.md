---
argument-hint: <issue-number>
description: Autonomous TDD development cycle with orchestrated subagents for TUI components
---

# TDD Development Cycle

Execute complete TDD workflow for GitHub issue `$ARGUMENTS` with built-in quality gates and subagent orchestration.

## Input Validation

Parse argument as GitHub issue number:
- Strip leading `#` if present (e.g., `#42` -> `42`)
- Validate it's a positive integer
- **STOP** if invalid - report error and exit

---

## Phase 1: Pre-flight Checks

Verify the development environment is ready:

1. **Check working directory status** - Run `git status` to ensure no uncommitted changes
2. **Verify current branch** - Should be on `main` branch
3. **Pull latest changes** - Run `git pull` to sync with remote
4. **Check for existing PRs** - Run `gh pr list --search "issue:$ARGUMENTS"` to avoid duplicate work
5. **Verify Go environment** - Run `go version` and ensure toolchain is available
6. **Check for blocking issues** - Look for "Depends on" or "Blocked by" in issue body

**STOP and report to user if any check fails** - Don't proceed until issues are resolved.

---

## Phase 2: Task Validation

Fetch and validate the GitHub issue:

1. **Fetch issue details**:
   ```bash
   gh issue view $ARGUMENTS --json title,body,labels,state,comments
   ```

2. **Validate issue**:
   - Issue must exist and be open
   - Issue must have a clear description or acceptance criteria

3. **Extract requirements**:
   - Title (feature/fix description)
   - Body (detailed requirements)
   - Labels (component type: output, input, adapter, layout)
   - Comments (any additional context or decisions)

**If incomplete:** Report gaps to user and ask if should proceed anyway.
**If complete:** Display issue summary and continue.

---

## Phase 3: Discovery & Planning

### Codebase Exploration (Automated)

Launch 3 Explore agents IN PARALLEL using the Task tool:

**Agent 1 - Local Pattern Discovery:**
```
Explore the claude-agent-tui codebase for patterns relevant to this task:

Project structure:
- component/output/ - Display components (StreamText, MessageBubble, etc.)
- component/input/ - Input components (ChatInput, PermissionPrompt, etc.)
- layout/ - Composite screens
- adapter/ - SDK integration layer
- style/ - Lip Gloss defaults

Search for:
- Similar component implementations in component/
- Bubble Tea model patterns (Init, Update, View)
- Lip Gloss styling patterns in style/
- Functional options patterns (WithStyles, WithXxx)

Use ast-grep for structural search:
- Find Model structs: ast-grep --pattern 'type Model struct { $$$ }' --lang go
- Find Update methods: ast-grep --pattern 'func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { $$$ }' --lang go

Report patterns found that should be followed for this implementation.
```

**Agent 2 - Test Conventions:**
```
Explore testing patterns in claude-agent-tui:

Search for:
- Existing test files (*_test.go)
- Table-driven test patterns
- teatest usage for TUI component testing
- Golden file patterns for ANSI output verification
- Test helpers and fixtures

Use ast-grep:
- Find test functions: ast-grep --pattern 'func Test$NAME(t *testing.T) { $$$ }' --lang go
- Find subtests: ast-grep --pattern 't.Run($NAME, func($$$) { $$$ })' --lang go

Report testing conventions that should be followed.
```

**Agent 3 - SDK Integration Patterns:**
```
Explore SDK integration patterns relevant to this task:

Check adapter/ directory for:
- Stream event to tea.Msg conversion patterns
- Hook callback handling
- Client lifecycle management
- Error propagation patterns

Check ../claude-code-sdk-go/ for:
- Message types to handle
- Client API usage
- Functional options available

Report integration patterns that apply to this implementation.
```

**WAIT for all 3 agents to complete before proceeding.**

### Architecture Planning (Automated)

Launch 1 Plan agent with context from the Explore agents:

```
Design a TDD implementation plan for GitHub issue #$ARGUMENTS:

Issue Details:
[Insert issue title, body, requirements from Phase 2]

Context from Exploration:
- Local Patterns: [Agent 1 findings]
- Test Conventions: [Agent 2 findings]
- SDK Integration: [Agent 3 findings]

Requirements:
1. Follow existing Bubble Tea patterns in the codebase
2. Use teatest with golden files for ANSI verification
3. Match styling patterns in style/
4. Follow functional options pattern

Provide a detailed TDD plan structured as:

## RED Phase - Tests First
- List specific test cases to write
- Include table-driven test structure
- Include edge cases (nil, empty, error conditions)
- Specify golden file tests for visual output

## GREEN Phase - Implementation
- List files to create/modify
- Describe Model struct fields needed
- Describe messages to define
- Describe Init/Update/View implementations

## BLUE Phase - Refactoring
- Quality improvements to apply
- Pattern alignment checks

## Acceptance Criteria Mapping
- Map each requirement to specific test cases
```

### Critical Checkpoint: User Approval

**Use ExitPlanMode tool to present the plan and await user approval.**

Do NOT continue to Phase 4 until user approves the plan.

---

## Phase 4: TDD Implementation

**Delegate to tdd-implementer agent** using the Task tool:

```
Execute TDD implementation for GitHub issue #$ARGUMENTS

## Approved Plan
[Insert the approved plan from Phase 3]

## Issue Details
- Title: [Issue title]
- Requirements: [Issue body/requirements]

## Branch Name
feature/issue-$ARGUMENTS-[short-description-kebab-case]

## Instructions
1. Create feature branch from main
2. Execute RED phase - write failing tests with teatest golden files
3. Execute GREEN phase - implement to pass tests
4. Execute BLUE phase - refactor with code-simplifier agent
5. Push branch to origin
6. Return: branch name, commit SHAs, test results

Follow the tdd-implementer agent workflow exactly.
```

**WAIT for tdd-implementer to complete.**

Capture returned:
- Branch name
- Commit SHAs (RED, GREEN, BLUE)
- Test results

---

## Phase 5: Automated Code Review

Launch 3 review agents IN PARALLEL using the Task tool (BLOCKING):

### Determine Reviewer Type

Based on files modified:
- If `component/` modified -> use `tui-component-reviewer`
- If `adapter/` modified -> use `sdk-adapter-reviewer`
- If both -> use both reviewers (launch 4 agents total)

### tui-component-reviewer OR sdk-adapter-reviewer Agent
```
Review the code changes in branch [branch-name] for TUI component/adapter patterns.
Focus on files modified in this branch.
Report issues in severity-categorized format.
```

### silent-failure-hunter Agent
```
Analyze error handling in the code changes on branch [branch-name]:
- Silent failures where errors are swallowed
- Inadequate error handling
- Inappropriate fallback behavior
- Missing error checks on function returns
- Bubble Tea Update handlers that ignore error messages

Report all findings with file:line references.
```

### pr-test-analyzer Agent
```
Review test coverage for the code changes on branch [branch-name]:
- Table-driven tests used for multiple cases
- Test helpers call t.Helper()
- Edge cases covered (nil, empty, error conditions)
- teatest golden files for ANSI verification
- Tests are deterministic (no flaky tests)
- Coverage adequate for new functionality

Report gaps and recommendations.
```

**WAIT for all review agents to complete.**

**BLOCKING GATE**: If ANY agent finds Critical or Major issues:
1. Compile all findings
2. Launch tdd-implementer agent again with fix instructions
3. Re-run review agents
4. Repeat until all issues resolved

---

## Phase 6: Validation

Run full validation suite:

### Test Suite
```bash
go test -cover -race ./...
```

Verify:
1. All tests pass
2. Coverage report acceptable
3. No race conditions detected

### Build Check
```bash
go build ./...
```

Ensure all packages build without errors.

### Linter Check
```bash
golangci-lint run
```

Address any linting issues (auto-fix with `go fmt` where possible).

**STOP if validation fails** - Report issues to user and await instructions.

---

## Phase 7: PR Creation & Merge

### Create Pull Request

```bash
gh pr create --title "feat: [Issue Title]" --body "$(cat <<'EOF'
## Summary
[1-2 sentence overview of changes]

## Issue Reference
Closes #$ARGUMENTS

## Changes

### Files Created
- `path/to/file.go` - Description

### Files Modified
- `path/to/file.go` - What changed

## Test Plan
- [x] All tests passing
- [x] Coverage maintained/improved
- [x] Race detector clean
- [x] Build succeeds
- [x] Linter clean

## TDD Cycle
- RED: [commit SHA] - Tests written first
- GREEN: [commit SHA] - Implementation added
- BLUE: [commit SHA] - Refactored (if applicable)
EOF
)"
```

### Interactive Checkpoint: PR Review

Display PR URL to user and present options:

1. **Approve** - Proceed to merge
2. **Request changes** - Wait for user edits
3. **Reject** - Close PR and rollback

### After User Approval: Merge PR

```bash
gh pr merge --squash --delete-branch
git checkout main
git pull
```

---

## Completion Summary

Display final summary:

```
TDD Development Cycle Complete for Issue #$ARGUMENTS

Phase 1: Pre-flight Checks - Done
Phase 2: Task Validation - Done
Phase 3: Discovery & Planning (User Approved) - Done
Phase 4: TDD Implementation (via tdd-implementer) - Done
  - RED: Tests written with teatest golden files
  - GREEN: Implementation complete
  - BLUE: Refactored (if applicable)
Phase 5: Automated Code Review - Done
  - Component/Adapter review: Passed
  - Silent failure hunt: Passed
  - Test coverage analysis: Passed
Phase 6: Validation - Done
Phase 7: PR Merged - Done

PR: #[number] (merged and branch deleted)
Issue: #$ARGUMENTS (closed)
Branch: main (updated)

Test Results:
- Tests: X passed
- Coverage: XX%
- Race conditions: None
```

---

## Error Recovery

If any phase fails:

1. **Pre-flight Failures**: Report to user, provide fix suggestions, stop
2. **Issue Not Found**: Report error, stop
3. **Planning Rejected**: Ask user for guidance, revise plan
4. **TDD Failures**: tdd-implementer handles internally
5. **Review Failures**: Loop back to fix and re-review
6. **Validation Failures**: Report issues, offer to fix
7. **PR Failures**: Report error, provide manual commands

**Branch is preserved on failure** - User can manually inspect and continue.
