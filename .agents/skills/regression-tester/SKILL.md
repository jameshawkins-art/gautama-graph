---
name: regression-tester
description: >-
  Executes, designs, and validates unit, integration, and fuzz test suites across the gautama-graph AST auditor and doc link graph subsystems.
  Use when verifying new features, checking patches against regressions, validating code coverage thresholds, running race condition detectors, or mocking IPC bridges and filesystem fixtures.
---

# Regression & Test Automation Skill (`gautama-guard`)

This skill defines the testing standards, regression verification workflows, mocking protocols, and automated reporting formats for certifying code changes in the `gautama-graph` repository.

---

## Section 1: Core Behavioral Rules & Guardrails

When executing test suites or writing regression tests, you MUST adhere strictly to the following testing laws:

### 1.1 Test Isolation & Hermetic Filesystem Usage
* **No Leaked Disk Fixtures**: Never write directly to the real repository workspace during tests. Always use `t.TempDir()` to construct sandboxed test directories for `graphify-out/graph.json` and `graphify-out/doc_graph_audit.json`.
* **Subprocess Mocking**: For unit tests in `internal/auditor/`, mock the `PythonASTBridge` interface rather than executing live Python subprocesses unless specifically running end-to-end integration tests.
* **Deterministic Fixtures**: Avoid non-deterministic values (e.g., `time.Now()` without freezing, random seeds without static values) in test assertions.

### 1.2 Table-Driven Test Patterns & Boundary Testing
* **Table-Driven Structure**: All new Go unit tests MUST follow standard table-driven test patterns:
  ```go
  tests := []struct {
      name      string
      input     CandidateEdge
      want      AuditedEdge
      wantErr   bool
  }{
      // test cases
  }
  ```
* **Compulsory Edge Case Coverage**: Every test suite for AST auditing and link resolution MUST test:
  1. Empty input / empty file / zero nodes.
  2. Malformed AST syntax (unparseable Go/Python source).
  3. Circular references (cyclic markdown links).
  4. Path traversal attempts (`../../shadow_file.go`).
  5. Cancelled contexts (`ctx, cancel := context.WithCancel(...)`).

### 1.3 Race Detection & Concurrency Verification
* **Mandatory `-race` Flag**: All unit test executions MUST run with Go's race detector enabled:
  `go test -v -race ./internal/auditor/... ./cmd/...`
* **Parallel Test Safety**: Use `t.Parallel()` on independent sub-tests, ensuring shared mocks or variables are properly synchronized.

### 1.4 Coverage Gate & Quality Thresholds
* **Strict 85% Minimum Coverage**: Any pull request or feature branch modifying `internal/auditor/` MUST maintain >= 85% statement coverage. If coverage drops below 85%, the Guard agent MUST flag a regression and fail the release gate.
* **Fuzz Testing for Parsers**: Implement Go native fuzzing (`func FuzzParseDocLinks(f *testing.F)`) for user-input parsers (e.g., Markdown regex extraction in `doc_auditor.go`).

---

## Section 2: Expected Antigravity Artifact Deliverables

The Regression & Test Automation agent MUST produce the following structured Antigravity artifacts stored in `<appDataDir>/brain/<conversation-id>/`:

### 2.1 Unit Execution Report (`unit_execution_report.md`)
Generated after running test suites to document pass/fail statuses, run durations, and coverage metrics.

```markdown
# Unit Execution Report: [Test Run Identifier]

## Executive Summary
- **Overall Status**: PASSED / FAILED
- **Total Tests Executed**: 42
- **Passed**: 42 | **Failed**: 0 | **Skipped**: 0
- **Duration**: 1.84s
- **Race Detector Clean**: YES

## Package Coverage Matrix
| Package | Statement Coverage | Minimum Required | Status |
| :--- | :--- | :--- | :--- |
| `internal/auditor` | 89.4% | 85.0% | PASS |
| `cmd/graphify-ast-audit` | 91.2% | 85.0% | PASS |
| `cmd/graphify-doc-audit` | 88.0% | 85.0% | PASS |

## Test Suite Breakdown
- `TestEngine_AuditGraphFile_GoAST`: PASS (0.02s)
- `TestEngine_AuditGraphFile_PythonBridge`: PASS (0.15s)
- `TestDocGraphAuditor_OrphanAndBrokenLinks`: PASS (0.08s)
- `TestJSONGraphStore_AtomicWriteRename`: PASS (0.04s)
```

### 2.2 Regression Matrix Checklist (`regression_matrix.md`)
Matrix verifying all legacy capabilities remain functional after recent modifications.

| Feature Area | Test Target | Pre-Condition | Verification Method | Result |
| :--- | :--- | :--- | :--- | :--- |
| Go AST Parsing | `parser.go` | Go 1.26 syntax | `TestASTParser_ValidGoFile` | Verified |
| Python Bridge IPC | `python_bridge.go` | Subprocess pipe | `TestPythonBridge_MockAndE2E` | Verified |
| Markdown Link Scan | `doc_auditor.go` | Relative & absolute links | `TestDocAuditor_PathResolution` | Verified |
| Atomic Store | `store.go` | Concurrent writers | `TestStore_AtomicRenameSafety` | Verified |

### 2.3 Test Verification Meta-Artifact (`test_verification_meta.json`)
Machine-readable output consumed by `gautama-gatekeeper` to approve or reject code changes.

```json
{
  "$schema": "https://antigravity.internal/schemas/v2/test-verification.json",
  "run_id": "TEST-2026-08-24-001",
  "timestamp": "2026-08-24T08:10:00Z",
  "guard_agent": "gautama-regression-guard",
  "all_passed": true,
  "race_detected": false,
  "coverage_percentage": 89.4,
  "coverage_gate_met": true,
  "fuzz_cycles_completed": 10000,
  "ready_for_security_audit": true
}
```
