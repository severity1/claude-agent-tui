---
argument-hint: <issue-number>
description: Autonomous TDD team with task-list-driven self-coordination - scales to scope, reviewers as teammates
allowed-tools: Read, Grep, Glob, Bash, Edit, Write, Task, WebFetch, TeamCreate, TeamDelete, TaskCreate, TaskUpdate, TaskList, TaskGet, SendMessage, AskUserQuestion, EnterPlanMode, ExitPlanMode
---

# TDD Development Cycle (Team-Based)

Execute complete TDD workflow for GitHub issue `$ARGUMENTS` using team-based agent orchestration with self-coordinating teammates.

## Input Validation

Parse argument as GitHub issue number:
- Strip leading `#` if present (e.g., `#42` -> `42`)
- Validate it's a positive integer
- **STOP** if invalid - report error and exit

---

## Phase 1: Preflight Checks

Verify the development environment is ready:

1. **Working directory** - `git status` must show clean working tree
2. **Current branch** - Must be on `main`
3. **Pull latest** - `git pull` to sync with remote
4. **Existing PRs** - `gh pr list --search "issue:$ARGUMENTS"` to avoid duplicates
5. **Go environment** - `go version` confirms toolchain available
6. **Blocking issues** - Check issue body for "Depends on" or "Blocked by"

**STOP and report if any check fails.**

---

## Phase 2: Discovery & Planning

### Enter Plan Mode

Use EnterPlanMode to begin planning.

### Fetch Issue

```bash
gh issue view $ARGUMENTS --json title,body,labels,state,comments
```

Validate: issue exists, is open, has clear requirements. Extract title, body, labels, comments.

### Codebase Exploration

Launch 2 Explore agents IN PARALLEL (model: "sonnet"):

**Agent 1 - Build System & Conventions:**
```
Explore claude-agent-tui for build and quality conventions:
- Read Makefile for targets (make check, make test-race, make test-tui)
- Read .claude/skills/charmbracelet/SKILL.md for Bubble Tea patterns
- Identify test framework (teatest, golden files), commit conventions
- Find existing component patterns: Model structs, Init/Update/View, functional options
- Report: build commands, quality gates, test patterns, commit prefixes
```

**Agent 2 - Affected Code & Architecture:**
```
Analyze code relevant to issue #$ARGUMENTS:
- Identify affected files, packages, dependencies
- Map existing patterns in affected areas
- Check adapter/ for SDK integration patterns if relevant
- Check component/ for TUI patterns if relevant
- Read .claude/skills/claude-agent-sdk-go/SKILL.md if SDK work needed
- Report: files to modify/create, dependencies, patterns to follow
```

### Design Task Graph

After exploration completes, design the implementation plan:

**Task graph structure:**
```
RED -> QA RED -> GREEN -> QA GREEN -> Review -> Create PR
```

- RED: Write failing tests (commit prefix: `test:`)
- QA RED: Verify existing gates pass, new tests FAIL, no regressions
- GREEN: Write production-quality code to pass tests (commit prefix: `feat:`)
- QA GREEN: Run ALL quality gates (make check, make test-race)
- Review: Specialized code reviewers (lead dispatches JIT)
- Create PR: Final step (lead-owned)

**Track types based on scope:**
- **Standard**: RED -> QA RED -> GREEN -> QA GREEN (for new features, complex changes)
- **Lightweight**: IMPLEMENT -> QA VERIFY (for simple fixes, config changes)

**Multi-issue/multi-track**: Parallel tracks with cross-deps only when files overlap.

**Team composition (determined during planning):**
- Minimum: 1 implementer + 1 verifier
- Scale implementers to independent work tracks (max 3)
- Reviewers: spawned JIT when QA GREEN passes

### Present Plan

Use ExitPlanMode. Do NOT proceed until user approves.

---

## Phase 3: Team Spawn

After user approval:

### Setup

1. Sanitize issue number for branch/team names (e.g., `42` -> `tdd-42`)
2. Create branch: `feature/issue-$ARGUMENTS-[short-description-kebab-case]`
3. Check for existing team: if `~/.claude/teams/tdd-$ARGUMENTS/` exists, recover (see Error Recovery)

### Create Team

```
TeamCreate(team_name="tdd-$ARGUMENTS", description="TDD for issue #$ARGUMENTS")
```

### Create Task Graph

Use TaskCreate for each task, then TaskUpdate to set `addBlockedBy` dependencies.

Example for standard track:
```
Task 1: "RED: Write failing tests for [feature]"
Task 2: "QA RED: Verify test failures" (blockedBy: [1])
Task 3: "GREEN: Implement [feature]" (blockedBy: [2])
Task 4: "QA GREEN: Full quality gate verification" (blockedBy: [3])
Task 5: "Review: Specialized code review" (blockedBy: [4], owner: lead)
Task 6: "Create PR" (blockedBy: [5], owner: lead)
```

Each task description includes: what to do, files involved, success criteria.

### Spawn Teammates

All teammates use `model: "sonnet"`. Never use `mode: "plan"`.

**Implementer(s)** - one per independent work track:
```
You are an implementer on team tdd-$ARGUMENTS.

Self-coordination loop:
1. Call TaskList to find unblocked, unassigned tasks
2. Filter to implementation tasks (RED, GREEN, IMPLEMENT, fix - NOT QA/verify tasks)
3. Claim lowest-ID matching task via TaskUpdate (set owner to your name)
4. Read task description via TaskGet for full requirements
5. Work to completion, commit with appropriate prefix (test: for RED, feat: for GREEN, fix: for iterations)
6. Mark task completed via TaskUpdate
7. Message verifier with summary of what was done
8. Loop back to step 1

If no matching tasks available, message lead: "No available implementation work. Standing by."

Restrictions:
- NEVER use TeamCreate, TeamDelete, or broadcast
- NEVER create tasks unilaterally - propose to lead via SendMessage
- NEVER skip or modify QA/verify tasks
- Follow project conventions in CLAUDE.md
```

**Verifier** - one per team:
```
You are the verifier on team tdd-$ARGUMENTS.

Self-coordination loop:
1. Call TaskList to find unblocked, unassigned tasks
2. Filter to QA/verification tasks (QA, verify, VERIFY)
3. Claim lowest-ID matching task via TaskUpdate (set owner to your name)
4. Read task description via TaskGet for full criteria
5. Run ALL success criteria - report every failure, not just the first
6. On PASS: mark completed, loop back to step 1
7. On FAIL: message implementer with specific failures, keep task in_progress
   - Track iteration count. After 3 failed iterations, escalate to lead
8. When QA GREEN passes, message lead: "QA GREEN complete for [task description]"

If no matching tasks available, message lead: "No available QA work. Standing by."

QA RED criteria:
- Existing quality gates pass (make check)
- New tests FAIL (they must fail before implementation)
- No regressions in existing tests

QA GREEN criteria:
- make check passes (fmt, vet, lint, test)
- make test-race passes (race condition detection)
- make test-tui passes (golden file tests, if applicable)
- All new tests pass
- No regressions

Restrictions:
- NEVER use TeamCreate, TeamDelete, or broadcast
- NEVER modify source code - only run verification commands
- Report ALL failures, not just the first one found
```

---

## Phase 4: Autonomous Execution

Lead monitoring loop - wait for messages, react:

```
while shutdown criteria NOT met:
  wait for next teammate message

  on "QA GREEN complete":
    -> Dispatch reviews (Phase 4a)

  on review verdicts:
    -> Process review loop (Phase 4b)

  on escalation (QA iteration >= 3):
    -> Read failure details, attempt direct fix or ask user

  on task proposal from implementer:
    -> Evaluate, create task if appropriate, assign back

  on "No available work. Standing by.":
    -> Check TaskList. If unblocked work exists, point teammate to it.
    -> If all work genuinely done, acknowledge.

  on idle notification (teammate stopped):
    -> Normal behavior. Only act if teammate has in_progress tasks.
```

Teammates self-coordinate via TaskList. Lead does NOT proactively assign unclaimed tasks unless asked.

### Phase 4a: Review Dispatch (JIT)

When verifier messages "QA GREEN complete":

1. Assess diff to determine review scope: `git diff main...HEAD --stat`
2. Create review tasks via TaskCreate (no dependencies between reviewers - all parallel)
3. Spawn reviewer teammates (model: "sonnet"), selecting based on modified files:

| Modified Path | Reviewer Agent Type |
|---------------|---------------------|
| `component/` | tui-component-reviewer |
| `adapter/` | sdk-adapter-reviewer |
| Any Go file | grumpy-gopher |
| Any Go file | silent-failure-hunter |
| `*_test.go` | pr-test-analyzer |

Each reviewer gets this inline prompt:
```
You are a reviewer on team tdd-$ARGUMENTS.

1. Call TaskList, claim your review task via TaskUpdate
2. Execute review per your methodology on the files listed in the task
3. Send findings to lead via SendMessage with severity categories
4. Mark task completed via TaskUpdate
5. If no more work, message lead: "Review complete. Standing by."

Restrictions: NEVER use TeamCreate, TeamDelete, or broadcast
```

4. Wait for all reviewer messages

### Phase 4b: Review Loop

When reviewers report findings:

1. If all reviewers report no critical/major issues -> mark "Review" task complete
2. If critical/major issues found:
   a. Message implementer with consolidated findings
   b. Implementer fixes and commits (prefix: `fix:` or `refactor:`)
   c. Create "QA regression" task for verifier (blockedBy: implementer's fix)
   d. After QA passes, reviewer re-checks
   e. If review iteration >= 3: escalate to user with full context
3. Review iteration budgets: 3 per review loop (independent from QA iteration budget)

### Shutdown Criteria

ALL must be true:
1. Every task in TaskList is completed
2. No teammates have in_progress tasks
3. Final verification: `make check && make test-race`

---

## Phase 5: PR Creation & Cleanup

### Create PR

```bash
gh pr create --title "feat: [Issue Title]" --body "$(cat <<'EOF'
## Summary
[1-3 bullet points from implementation]

## Test Plan
- [x] TDD RED: Tests written first, verified failing
- [x] TDD GREEN: Implementation passes all tests
- [x] Quality gates: make check, make test-race
- [x] Code review: [list reviewers that participated]

## Review Summary
[Key findings addressed, reviewer verdicts]

Closes #$ARGUMENTS
EOF
)"
```

### Team Cleanup

1. Send `shutdown_request` to each teammate, wait for confirmations
2. `TeamDelete` to clean up team and task directories
3. Verify `~/.claude/teams/tdd-$ARGUMENTS/` and `~/.claude/tasks/tdd-$ARGUMENTS/` are gone

### Present PR URL to user

---

## Error Recovery

**Preflight failures**: Report to user, provide fix suggestions, stop.

**Issue not found**: Report error, stop.

**Planning rejected**: Revise plan based on user feedback, re-present.

**Unresponsive teammate**: Ping via SendMessage. If no response after 2 pings, escalate to user.

**Repeated QA failures** (iteration >= 3): Include full test output in escalation to user.

**Reviewer disagreement**: Attempt resolution via messages. If deadlocked after 2 rounds, escalate to user.

**Lead fallback**: If implementer stuck after 2 messages, lead writes fix directly and commits.

**Session recovery**: If team directory exists from previous run:
1. Read TaskList for last-known state
2. Reset stale `in_progress` tasks (no owner or owner gone) to `pending`
3. Spawn fresh teammates to continue from current state

**Branch preserved on failure** - user can manually inspect and continue.
