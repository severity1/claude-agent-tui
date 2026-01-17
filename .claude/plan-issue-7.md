# TDD Implementation Plan: Issue #7 - Initialize Project with Charmbracelet Dependencies

## Issue Summary

**Issue:** #7 - [MVP-1] Initialize project with Charmbracelet dependencies
**State:** OPEN
**Labels:** mvp, phase-1, setup

## Requirements

1. Install dependencies:
   - github.com/charmbracelet/bubbletea
   - github.com/charmbracelet/lipgloss
   - github.com/charmbracelet/bubbles/textarea
   - github.com/charmbracelet/bubbles/viewport
   - github.com/severity1/claude-agent-sdk-go

2. Create directory structure:
   ```
   claude-agent-tui/
   ├── go.mod          (already exists)
   ├── go.sum
   ├── example/
   │   └── chat/
   ├── component/
   │   ├── output/
   │   │   └── streamtext/
   │   └── input/
   │       ├── chatinput/
   │       └── permission/
   ├── adapter/
   └── layout/
       └── chat/
   ```

3. Acceptance criteria:
   - All Charmbracelet dependencies in go.sum
   - claude-agent-sdk-go dependency in go.sum
   - `go build ./...` passes without errors
   - Directory structure created

---

## RED Phase - Tests First

Since this is a project initialization issue, tests verify build/import success rather than business logic.

### Test 1: Import Verification Test

**File:** `imports_test.go`

Verifies all required dependencies are importable:
- `tea "github.com/charmbracelet/bubbletea"` - tea.Model, tea.Cmd accessible
- `"github.com/charmbracelet/lipgloss"` - lipgloss.NewStyle() works
- `"github.com/charmbracelet/bubbles/textarea"` - textarea.New() works
- `"github.com/charmbracelet/bubbles/viewport"` - viewport.New() works
- `claudecode "github.com/severity1/claude-agent-sdk-go"` - package accessible

### Test 2: Directory Structure Test

**File:** `structure_test.go`

Verifies required directories exist:
- example/chat
- component/output/streamtext
- component/input/chatinput
- component/input/permission
- adapter
- layout/chat

### Test 3: Package Placeholder Tests

Each package gets a minimal `*_test.go` to verify it compiles:
- `component/output/streamtext/streamtext_test.go`
- `component/input/chatinput/chatinput_test.go`
- `component/input/permission/permission_test.go`
- `adapter/adapter_test.go`
- `layout/chat/chat_test.go`

---

## GREEN Phase - Implementation

### Step 1: Install Dependencies

```bash
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/severity1/claude-agent-sdk-go@latest
go mod tidy
```

### Step 2: Create Directory Structure

```bash
mkdir -p example/chat
mkdir -p component/output/streamtext
mkdir -p component/input/chatinput
mkdir -p component/input/permission
mkdir -p adapter
mkdir -p layout/chat
```

### Step 3: Create Package Files

Each package needs a `doc.go` for documentation:
- `component/output/streamtext/doc.go` - Package streamtext
- `component/input/chatinput/doc.go` - Package chatinput
- `component/input/permission/doc.go` - Package permission
- `adapter/doc.go` - Package adapter
- `layout/chat/doc.go` - Package chat
- `example/chat/main.go` - Placeholder main
- `tui.go` - Root package documentation

### Step 4: Verify

```bash
go test ./...
go build ./...
```

---

## BLUE Phase - Refactoring

1. `go fmt ./...` - Format all code
2. `go vet ./...` - Static analysis
3. Verify package naming conventions (lowercase, no underscores)
4. Verify import grouping (stdlib, third-party, local)

---

## Final Directory Structure

```
claude-agent-tui/
├── go.mod                    # Updated with dependencies
├── go.sum                    # Generated
├── tui.go                    # Root package doc
├── imports_test.go           # Import verification test
├── structure_test.go         # Directory structure test
├── example/
│   └── chat/
│       └── main.go           # Placeholder main
├── component/
│   ├── output/
│   │   └── streamtext/
│   │       ├── doc.go
│   │       └── streamtext_test.go
│   └── input/
│       ├── chatinput/
│       │   ├── doc.go
│       │   └── chatinput_test.go
│       └── permission/
│           ├── doc.go
│           └── permission_test.go
├── adapter/
│   ├── doc.go
│   └── adapter_test.go
└── layout/
    └── chat/
        ├── doc.go
        └── chat_test.go
```

---

## Acceptance Criteria Mapping

| Criterion | Verification |
|-----------|--------------|
| Charmbracelet dependencies in go.sum | `TestDependencyImports` passes |
| claude-agent-sdk-go in go.sum | `TestDependencyImports` claudecode subtest passes |
| `go build ./...` passes | Build command exits 0 |
| Directory structure created | `TestDirectoryStructure` passes |

---

## Potential Risks

1. **SDK availability** - `github.com/severity1/claude-agent-sdk-go` must be published on GitHub. If not accessible, import test will fail.

2. **Version compatibility** - Using latest Charmbracelet packages. Pin versions if breaking changes occur.
