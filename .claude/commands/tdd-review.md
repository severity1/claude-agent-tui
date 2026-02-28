---
description: Comprehensive code review using team-based specialized agents with self-coordination
allowed-tools: Read, Grep, Glob, Bash, Edit, Write, Task, WebFetch, TeamCreate, TeamDelete, TaskCreate, TaskUpdate, TaskList, TaskGet, SendMessage, AskUserQuestion
---

# TDD Code Review (Team-Based)

Comprehensive code review using team-based specialized agents with self-coordination. Orchestrates multiple reviewers to analyze code quality, TUI patterns, error handling, and test coverage.

## Quick Mode Exception

For focused reviews, use a single bare Agent call (model: "sonnet") with NO team:
- `--quick` - Only grumpy-gopher
- `--component` - Only tui-component-reviewer
- `--errors` - Only silent-failure-hunter

If `$ARGUMENTS` contains any of these flags, extract the flag, launch the single appropriate agent with the file list, present results, and stop. Do not create a team.

---

## Review Scope Detection

Determine what code to review:

1. **User-specified files**: If user provided specific files or packages after the flag/arguments, use those
2. **Unstaged changes**: `git diff --name-only`
3. **Staged changes**: `git diff --staged --name-only`
4. **Default**: Current package (`./...`)

```bash
git diff --name-only
git diff --staged --name-only
```

Store the file list. If no files found, report and exit.

---

## Reviewer Selection

Based on files identified, select reviewers:

| File Pattern | Reviewer |
|--------------|----------|
| `component/output/*` | tui-component-reviewer |
| `component/input/*` | tui-component-reviewer |
| `layout/*` | tui-component-reviewer |
| `adapter/*` | sdk-adapter-reviewer |
| `style/*` | tui-component-reviewer |
| `*_test.go` | pr-test-analyzer |
| Any Go file | grumpy-gopher |
| Any Go file | silent-failure-hunter |

**Always include**: grumpy-gopher, silent-failure-hunter (for any Go files).

Deduplicate: if tui-component-reviewer would be selected multiple times, launch once with all matching files.

---

## Phase 1: Team Creation & Spawn

### Create Team

Generate a short ID from timestamp or random suffix:
```
TeamCreate(team_name="review-{short-id}", description="Code review for [file summary]")
```

### Create Review Tasks

One task per selected reviewer. No dependencies - all run in parallel.

Example:
```
Task 1: "Review: TUI component patterns" (files: component/output/streamtext/...)
Task 2: "Review: Go idioms and best practices" (files: [all Go files])
Task 3: "Review: Error handling and silent failures" (files: [all Go files])
Task 4: "Review: Test coverage analysis" (files: [test files])
```

Each task description includes: reviewer type, file list, focus areas.

### Spawn Reviewers

All reviewers use `model: "sonnet"`. Each gets this inline prompt:

```
You are a reviewer on team review-{short-id}.

Self-coordination loop:
1. Call TaskList for unblocked, unassigned review tasks
2. Claim lowest-ID matching task via TaskUpdate (set owner to your name)
3. Read task description via TaskGet for file list and focus areas
4. Execute review per your methodology
5. Send findings to lead via SendMessage in this format:
   - CRITICAL: [issues that must fix - bugs, security, data loss]
   - MAJOR: [should fix - significant quality issues]
   - MINOR: [nice to fix - style, minor improvements]
   - POSITIVE: [good patterns observed]
6. Mark task completed via TaskUpdate
7. Check TaskList for more work. If none: message lead "Review complete. Standing by."

Restrictions:
- NEVER use TeamCreate, TeamDelete, or broadcast
- NEVER modify source code - review only
- Report ALL findings with file:line references
```

---

## Phase 2: Collect Results

Wait for all reviewer messages. As each arrives:

1. **Collect findings** from each reviewer message
2. **Deduplicate** - same issue reported by multiple reviewers (keep the most detailed)
3. **Categorize by severity**:
   - **Critical**: Must fix (bugs, security, data loss)
   - **Major**: Should fix (significant quality issues)
   - **Minor**: Nice to fix (style, minor improvements)

---

## Phase 3: Filter MVP-Scoped Issues

Before presenting, filter out issues tracked elsewhere:

### Check GitHub Issues

```bash
gh issue list --label MVP --state open --json number,title,body --limit 50
```

### Filtering Criteria

**FILTER OUT** if:
1. **Tracked in GitHub Issues** (MVP label) - already planned. Note: "Tracked in #123"
2. **Planned in docs/** - features documented as future work
3. **Example/demo code** (`example/` directory) - unless misleading for library users
4. **Documentation-only** - missing docs for planned features
5. **Test infrastructure** - test helpers that ship with their feature

**KEEP** if:
- Blocks current branch functionality
- Code quality issue in core library code
- Security or correctness issue regardless of location

### Present Filtered List

```
FILTERED ISSUES (addressed by other planned work):
- [Issue description] - Reason: Tracked in #123
- [Issue description] - Reason: Planned in docs/ROADMAP.md

Proceed with remaining [N] issues?
```

Use AskUserQuestion to confirm.

---

## Phase 4: Present Summary

```
CODE REVIEW SUMMARY
===================

Files Reviewed: X
Reviewers: [list of reviewer types used]

CRITICAL ISSUES (Must Fix)
--------------------------
| Issue | Location | Reviewer | Recommended Fix |
|-------|----------|----------|-----------------|

MAJOR ISSUES (Should Fix)
-------------------------
| Issue | Location | Reviewer | Recommended Fix |
|-------|----------|----------|-----------------|

MINOR ISSUES (Nice to Fix)
--------------------------
| Issue | Location | Reviewer | Recommended Fix |
|-------|----------|----------|-----------------|

POSITIVE OBSERVATIONS
---------------------
- [Good patterns observed by reviewers]

OVERALL ASSESSMENT
------------------
- Component Quality: X/10
- Error Handling: X/10
- Test Coverage: X/10
- Go Idioms: X/10

Production Ready: [Yes/No - with explanation]
```

---

## Phase 5: Offer Fixes

If issues were found, present options via AskUserQuestion:

1. **Fix all issues** - Apply fixes for all severities
2. **Fix critical only** - Only address critical issues
3. **Fix selected** - User chooses which to fix
4. **Skip fixes** - Report only, no changes

If user chooses to fix:
1. Apply fixes using Edit tool
2. Run `make check` to verify fixes
3. Report what was fixed and verify no regressions

---

## Phase 6: Re-review (Optional)

After fixes applied, offer via AskUserQuestion:

1. **Yes** - Re-run full review (loop back to Phase 1 with same team if still active)
2. **No** - Trust the fixes, proceed to cleanup

If re-reviewing with existing team:
- Create new review tasks for fixed files only
- Message existing reviewers (if still active) or spawn new ones
- Collect and present delta results

---

## Team Cleanup

1. Send `shutdown_request` to each reviewer, wait for confirmations
2. `TeamDelete` to clean up team and task directories
3. Verify directories are removed

---

## Integration with /tdd

When called from `/tdd` command (Phase 4a):
- Review is BLOCKING - must pass before PR creation
- Critical/major issues trigger fix loop automatically
- No user interaction for severity selection - all issues must be addressed
- Findings sent back to lead via SendMessage (not presented as summary table)

---

## Error Handling

- **Agent timeout**: Report which reviewer timed out, continue with others' results
- **No files to review**: Report and exit gracefully, no team created
- **Git errors**: Report git state issue, suggest fix
- **Fix application fails**: Report error, show manual fix instructions
- **Team exists**: If `review-{short-id}` already exists, generate new ID and retry
