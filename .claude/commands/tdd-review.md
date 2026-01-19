---
description: Comprehensive code review using parallel specialized agents for TUI components and SDK integration
allowed-tools: Read, Grep, Glob, Bash(go vet:*), Bash(golangci-lint run:*), Bash(go test:*), Bash(git diff:*), Task
---

# TDD Code Review

Comprehensive code review using parallel specialized agents. This command orchestrates multiple reviewers to analyze code quality, TUI patterns, error handling, and test coverage.

## Review Scope Detection

Determine what code to review:

1. **User-specified files**: If user provided specific files or packages, use those
2. **Recent changes**: Otherwise, check `git diff` for uncommitted changes
3. **Staged changes**: Check `git diff --staged` for staged changes
4. **Default**: Review current package if no changes detected

```bash
# Check for changes
git status --porcelain
git diff --name-only
git diff --staged --name-only
```

Store the list of files to review.

---

## Reviewer Selection

Based on files identified, determine which reviewers to launch:

| File Pattern | Reviewer |
|--------------|----------|
| `component/output/*` | tui-component-reviewer |
| `component/input/*` | tui-component-reviewer |
| `layout/*` | tui-component-reviewer |
| `adapter/*` | sdk-adapter-reviewer |
| `style/*` | tui-component-reviewer |
| `*_test.go` | pr-test-analyzer |
| Any Go file | grumpy-gopher, silent-failure-hunter |

**Always include**: grumpy-gopher, silent-failure-hunter

---

## Phase 1: Launch Parallel Review Agents

Launch up to 3 review agents IN PARALLEL using the Task tool:

### Agent 1: Component/Adapter Reviewer

If TUI component files detected, launch **tui-component-reviewer**:
```
Review the following TUI component files for Bubble Tea patterns:

Files to review:
[List of component files]

Focus on:
- Model struct correctness
- Init/Update/View implementation
- Message type design
- Focus management (for input components)
- Lip Gloss styling consistency
- Functional options pattern

Provide findings in severity-categorized table format.
```

If adapter files detected, launch **sdk-adapter-reviewer**:
```
Review the following SDK adapter files:

Files to review:
[List of adapter files]

Focus on:
- Stream event conversion
- Channel and goroutine safety
- Error propagation
- Client lifecycle management

Provide findings in severity-categorized table format.
```

### Agent 2: Go Expert Reviewer (grumpy-gopher)

```
Review the following Go files for best practices and idioms:

Files to review:
[List of all Go files]

Focus on:
- Go idioms and best practices
- Error handling patterns
- Interface design
- Resource management
- Naming conventions

Be thorough and technical. Apply your grumpy personality after analysis.
```

### Agent 3: Error Handling Reviewer (silent-failure-hunter)

```
Analyze error handling in the following files:

Files to review:
[List of all Go files]

Hunt for:
- Silent failures (errors swallowed)
- Empty catch blocks or ignored errors
- Inappropriate fallback behavior hiding problems
- Errors logged but not propagated
- Missing error checks on function returns
- Update handlers ignoring error messages

Report all findings with file:line references and severity.
```

**WAIT for all agents to complete.**

---

## Phase 2: Compile Results

Aggregate findings from all reviewers:

1. **Collect all issues** from each agent
2. **Deduplicate** - Same issue reported by multiple agents
3. **Categorize by severity**:
   - **Critical**: Must fix before merge (bugs, security, data loss)
   - **Major**: Should fix (significant quality issues)
   - **Minor**: Nice to fix (style, minor improvements)

---

## Phase 2.5: Filter MVP-Scoped Issues

Before presenting results, filter out issues that will be addressed by other MVP GitHub issues:

### Check GitHub Issues First

Query open MVP issues to understand planned work:

```bash
# List MVP-labeled issues
gh issue list --label MVP --state open --json number,title,body --limit 50

# Search for specific feature mentions
gh issue list --search "adapter" --state open --json number,title
```

Cross-reference review findings against open issues. If a finding is already tracked, note the issue number and filter it.

### Filtering Criteria

**FILTER OUT** issues that match:

1. **Tracked in GitHub Issues (MVP label)** - Already planned work
   - Run `gh issue list --label MVP` to get current MVP scope
   - If finding matches an open issue, note: "Tracked in #123"
   - Example: "Missing ToolResultMsg" - if Issue #15 covers this, filter out

2. **Planned features documented in `docs/`** - Features explicitly marked as planned/future work
   - Check `docs/ADAPTERS.md`, `docs/ROADMAP.md` for planned components
   - Example: "Missing TaskStartedMsg type" - if documented as planned, filter out

3. **Example/demo code issues** - Issues in `example/` directory that don't affect core library
   - Example app is for validation, not production use
   - Filter unless the issue would mislead users about library usage

4. **Documentation-only changes** - Issues about missing docs for planned features

5. **Test infrastructure** - Issues about test helpers that will be added with their features

### How to Filter

For each issue, check:
```
1. Is this tracked in OR directly supports a GitHub MVP issue? -> FILTER (note issue #)
2. Is this in docs/ as planned work? -> FILTER
3. Is this in example/ and non-critical? -> FILTER
4. Does this block current branch functionality? -> KEEP
5. Is this a code quality issue in core library? -> KEEP
```

**"Directly supports"** means the finding relates to functionality that an MVP issue depends on or enables, even if not explicitly mentioned in that issue. Use judgment - if fixing this finding would be part of completing an MVP issue, filter it.

### Present Filtered List

After filtering, ask user to confirm:
```
FILTERED ISSUES (will be addressed by other MVP work):
- [Issue 1] - Reason: Documented in docs/ADAPTERS.md as Phase 2
- [Issue 2] - Reason: Example app only, not core library

Proceed with remaining [N] issues? [Y/n]
```

---

## Phase 3: Present Summary

Display aggregated results in table format:

```
CODE REVIEW SUMMARY
===================

Files Reviewed: X
Reviewers: tui-component-reviewer, grumpy-gopher, silent-failure-hunter

CRITICAL ISSUES (Must Fix)
--------------------------
| Issue | Location | Reviewer | Recommended Fix |
|-------|----------|----------|-----------------|
| ...   | file:line | ...     | ...             |

MAJOR ISSUES (Should Fix)
-------------------------
| Issue | Location | Reviewer | Recommended Fix |
|-------|----------|----------|-----------------|
| ...   | file:line | ...     | ...             |

MINOR ISSUES (Nice to Fix)
--------------------------
| Issue | Location | Reviewer | Recommended Fix |
|-------|----------|----------|-----------------|
| ...   | file:line | ...     | ...             |

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

## Phase 4: Offer Fixes

If issues were found, ask user:

1. **Fix all issues** - Automatically apply fixes for all issues
2. **Fix critical only** - Only fix critical issues
3. **Fix selected** - Let user choose which to fix
4. **Skip fixes** - Don't fix, just report

If user chooses to fix:

1. Apply fixes using Edit tool
2. Run `go fmt ./...` and `go vet ./...`
3. Run tests to verify fixes don't break anything
4. Report what was fixed

---

## Phase 5: Re-review (Optional)

After fixes are applied, offer to re-run review:

```
Fixes applied. Would you like to re-run the review to verify all issues are resolved?

1. Yes - Run full review again
2. No - Trust the fixes, we're done
```

If yes, loop back to Phase 1 with the fixed files.

---

## Quick Review Mode

For fast reviews (less thorough), user can specify:
- `/tdd-review --quick` - Only run grumpy-gopher
- `/tdd-review --component` - Only run tui-component-reviewer
- `/tdd-review --errors` - Only run silent-failure-hunter

Default is full review with all agents.

---

## Integration with /tdd

The `/tdd` command uses this review workflow in Phase 5. When called from `/tdd`:

1. Review is BLOCKING - must pass before PR
2. Critical/Major issues trigger fix loop
3. No user interaction for severity selection - all issues must be fixed

---

## Error Handling

If review fails:

1. **Agent timeout**: Report which agent timed out, continue with others
2. **No files to review**: Report and exit gracefully
3. **Git errors**: Report git state issue, suggest fix
4. **Fix application fails**: Report error, show manual fix instructions
