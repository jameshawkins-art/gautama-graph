---
name: debugger-remediation
description: >-
  Investigates, triages, reproduces, and resolves runtime panics, IPC bridge crashes, deadlocks, and AST parsing lockups across Go and Python subsystems in gautama-graph.
  Use when analyzing stack traces, debugging python/ast_auditor_bridge.py failures, resolving broken link cycle errors, fixing memory leaks, or emitting surgical remediation patches.
---

# Debugger & Remediation Skill (`gautama-mechanic`)

This skill provides step-by-step diagnostic workflows, fault isolation heuristics, safety protocols, and structured artifact delivery specifications for triaging and repairing runtime failures in the `gautama-graph` engine.

---

## Section 1: Core Behavioral Rules & Guardrails

When triaging incidents and engineering bug fixes, you MUST adhere strictly to the following remediation principles:

### 1.1 Root-Cause Isolation & Minimal Surgical Patches
* **Zero Collateral Modification**: Fixes MUST be surgical and strictly confined to the faulting execution path. Do not refactor unrelated code, change public method signatures, or modify formatting of undisturbed files.
* **Deterministic Reproduction First**: Never commit a code patch without first isolating the failure into an executable reproduction case (e.g., in `internal/auditor/auditor_test.go` or standalone scratch test `testdata/repro/`).
* **Preserve Semantic Error Contracts**: Do not suppress errors or convert explicit errors into silent no-ops. Always preserve or augment error chains with `fmt.Errorf("...: %w", err)`.

### 1.2 IPC Bridge & Subprocess Crash Defense (`python/ast_auditor_bridge.py`)
* **JSON Payload Sanitization**: Guard against malformed JSON or empty stdin buffers in `python/ast_auditor_bridge.py`. The bridge MUST return a valid JSON payload `{"status": "error", "error": "..."}` rather than exiting with unhandled Python tracebacks.
* **Subprocess Deadlock Prevention**: When reading stdout and stderr from `os/exec.Cmd`, ensure buffers are read concurrently or via `cmd.CombinedOutput()` / bounded pipes to prevent process deadlocks when pipe buffers fill up.
* **Timeout & Context Cancellation**: Subprocess calls MUST honor `ctx.Done()`. If a Python process hangs on a deeply nested or circular AST, `exec.CommandContext` MUST cleanly kill the subprocess group and release file handles.

### 1.3 Concurrency & Mutex Safety
* **Defer Mutex Unlock**: In `JSONGraphStore` and all stateful structures, `mu.Lock()` MUST immediately be followed by `defer mu.Unlock()`. Never place return statements between lock acquisition and deferral.
* **Avoid Nested Lock Deadlocks**: Never invoke an external interface method or callbacks while holding an internal `sync.Mutex` lock.

### 1.4 AST Walkers & Link Graph Recursion Guardrails
* **Cycle Detection in Doc Links**: When resolving markdown links in `DocGraphAuditor`, maintain a visited set `map[string]bool` to prevent infinite loops on cyclic link graphs.
* **AST Depth Limiter**: Respect `Config.MaxASTDepth` during recursive AST node inspections to prevent stack exhaustion on adversarial or pathological source inputs.

---

## Section 2: Expected Antigravity Artifact Deliverables

The Debugger & Remediation agent MUST generate concrete diagnostic artifacts stored in `<appDataDir>/brain/<conversation-id>/`:

### 2.1 Root Cause Analysis & Remediation Proposal (`remediation_proposal.md`)
Generated upon investigating a bug, detailing the fault timeline, reproduction steps, and candidate patch.

```markdown
# Remediation Proposal: [Incident / Bug Title]

## Incident Summary
- **Symptom**: [e.g., Subprocess deadlock during Python AST bridge batch audit]
- **Component**: `internal/auditor/python_bridge.go` & `python/ast_auditor_bridge.py`
- **Severity**: High (Engine halts on large Python candidate lists)

## Root Cause Analysis
Explain the technical failure mechanism (e.g., unbuffered stdout write exceeding 64KB OS pipe buffer).

## Minimal Reproduction Case
```go
func TestRepro_PythonBridgeDeadlock(t *testing.T) {
    // Isolated reproduction code
}
```

## Proposed Code Patch (Diff)
```diff
--- a/internal/auditor/python_bridge.go
+++ b/internal/auditor/python_bridge.go
@@ -45,3 +45,5 @@
-   out, err := cmd.Output()
+   out, err := cmd.CombinedOutput()
```

## Side-Effect & Regression Risk Analysis
Assess potential impacts on downstream callers or memory footprints.
```

### 2.2 Patch Feasibility Proposal (`patch_feasibility_proposal.md`)
Required when the fix requires structural adjustments or architectural alterations to prevent recurring failures.
* **Metadata**: `RequestFeedback: true`, `UserFacing: true`
* Contains alternative resolution approaches, trade-offs, and rollout considerations.

### 2.3 JSON Remediation Meta-Artifact (`remediation_meta.json`)
Emitted to hand off the verified patch to `gautama-guard` for regression testing.

```json
{
  "$schema": "https://antigravity.internal/schemas/v2/remediation-meta.json",
  "incident_id": "BUG-2026-08-01",
  "timestamp": "2026-08-24T08:05:00Z",
  "mechanic_agent": "gautama-remediation-mechanic",
  "fault_category": "IPC_BRIDGE_DEADLOCK",
  "affected_files": [
    "internal/auditor/python_bridge.go",
    "python/ast_auditor_bridge.py"
  ],
  "reproduction_test_name": "TestRepro_PythonBridgeDeadlock",
  "patch_applied": true,
  "requires_guard_regression_run": true
}
```
