---
name: feature-engineer
description: >-
  Architects, implements, and extends Go 1.26+ features, AST parsers, selector evaluators, and CLI commands for the gautama-graph ecosystem.
  Use when designing new library capabilities, refactoring internal/auditor services, exporting public APIs, adding CLI subcommands to cmd/,
  or updating knowledge graph persistence formats.
---

# Feature Engineer Skill (`gautama-builder`)

This skill provides direct operational instructions, architectural guardrails, and artifact generation templates for implementing production features across the `gautama-graph` AST verification and documentation auditing engine.

---

## Section 1: Core Behavioral Rules & Guardrails

When developing or modifying features within the `gautama-graph` codebase, you MUST adhere strictly to the following core engineering rules:

### 1.1 Public API Export & Go Naming Conventions
* **Strict PascalCase for Exported Identifiers**: All exported structs, interfaces, methods, functions, and constants (e.g., `CandidateEdge`, `ASTAuditReport`, `NewJSONGraphStore`, `AuditGraphFile`) MUST begin with an uppercase letter. Unexported helper functions and fields MUST begin with a lowercase letter (camelCase).
* **Comprehensive Godoc Documentation**: Every exported symbol MUST have a godoc comment beginning with the exact name of the symbol.
* **Interface-First Architecture**: Expose behavior through Go interfaces defined in `internal/auditor/types.go` (e.g., `ASTParser`, `SelectorEvaluator`, `PythonASTBridge`, `GraphStore`, `DocGraphAuditorService`). Struct implementations must accept interfaces for dependencies to preserve testability.

### 1.2 Deterministic Path & Traversal Protection
* **Absolute Path Canonicalization**: Never use raw relative path strings directly for disk operations. Always resolve target paths using `filepath.Clean(rawPath)` and convert to absolute paths via `filepath.Abs(cleanPath)`.
* **Workspace Root Confinement**: Validate that all parsed files and audit targets reside strictly within `Config.WorkspaceRootPath`. Any attempt to resolve a path escaping the root (e.g., `../../etc/passwd`) MUST immediately return an explicit domain error: `fmt.Errorf("path %s escapes workspace root %s: %w", target, root, ErrPathOutOfBounds)`.

### 1.3 Atomic File Persistence via Temporary Buffers
* **Safe Buffer Staging**: Under no circumstances should `graphify-out/graph.json` or `graphify-out/doc_graph_audit.json` be overwritten in-place with partial streams.
* **Two-Phase Commit Protocol**:
  1. Write serialized JSON payload to `targetPath + ".tmp"` with permissions `0644`.
  2. Commit the change using atomic filesystem replacement: `os.Rename(tmpFile, targetPath)`.
  3. Clean up the temporary buffer on any failure branch: `defer func() { _ = os.Remove(tmpFile) }()`.
* **Locking**: Guard all write operations across `JSONGraphStore` and `DocGraphStore` using a dedicated `sync.Mutex`.

### 1.4 Context Propagation & Subprocess Discipline
* **Context Propagation**: All synchronous and asynchronous APIs traversing IO, AST walks, or subprocess executions MUST accept `ctx context.Context` as their first parameter and actively check `ctx.Err()` during long-running traversals.
* **Subprocess Execution**: When dispatching Python AST scripts via `exec.CommandContext`, always set an execution deadline (default 10s), sanitize environment variables, and stream through explicit pipes with bounded memory allocation.

### 1.5 Versioning & Git Tag Consistency
* **Semantic Versioning**: When modifying public interfaces in `internal/auditor/types.go`, ensure backwards compatibility. If breaking changes are unavoidable, update the module version path and prepare git tag release notes (`vX.Y.Z`).

### 1.6 Production Call-Site Invariant & TDD Alignment
* **Production Call-Site Invariant**: Whenever authoring a new helper function, parser method, or utility, immediately wire it into active production call paths (e.g., within `internal/auditor/` engine routines or `cmd/` CLI entrypoints).
* **Zero Orphaned Functions**: Implementing a function tested only within `*_test.go` while leaving production codepaths using duplicate inline logic is strictly classified as a **CRITICAL INCOMPLETE INTEGRATION DEFECT**.
* **DRY De-duplication**: Search and refactor any existing duplicate inline logic across packages to call the centralized function. Verify $C_{prod} > 0$ with `go run cmd/graphify-ast-audit/main.go`.

---

## Section 2: Expected Antigravity Artifact Deliverables

The Feature Engineer agent MUST NOT conclude feature workflows with conversational logs. It must generate structured Antigravity artifacts stored in `<appDataDir>/brain/<conversation-id>/`:

### 2.1 Feature Implementation Plan (`feature_implementation_plan.md`)
Generated prior to code modifications to define scope, interfaces, and migration steps.

```markdown
# Feature Implementation Plan: [Feature Title]

## Summary
Brief description of the feature addition or API modification for gautama-graph.

## Interface Specifications
```go
// Proposed exported types or interface modifications
type NewFeatureService interface {
    ExecuteFeature(ctx context.Context, param string) (*FeatureResult, error)
}
```

## Proposed File Changes
- [NEW] `internal/auditor/new_feature.go`
- [MODIFY] `internal/auditor/types.go`
- [MODIFY] `cmd/graphify-ast-audit/main.go`

## Verification Plan
- Unit tests: `go test -v -race ./internal/auditor -run TestNewFeature`
- Atomic write validation: Verify `.tmp` rename integrity under concurrent writes.
```

### 2.2 Patch Feasibility Proposal (`patch_feasibility_proposal.md`)
Generated when modifying core engine algorithms (e.g., AST visitor recursion, confidence scoring heuristics).

* **Path**: `<appDataDir>/brain/<conversation-id>/patch_feasibility_proposal.md`
* **Artifact Metadata**: `RequestFeedback: true`, `UserFacing: true`
* **Contents**:
  * Root Cause / Feature Driver
  * Architectural Trade-offs (Memory allocation vs AST traversal speed)
  * Backward Compatibility Assessment
  * Diff Preview of `internal/auditor/` changes

### 2.3 Feature Delivery Meta-Artifact (`feature_delivery.json`)
Machine-readable artifact delivered upon feature completion for consumption by `gautama-guard` and `gautama-gatekeeper`.

```json
{
  "$schema": "https://antigravity.internal/schemas/v2/feature-delivery.json",
  "feature_id": "feat-ast-selector-depth",
  "timestamp": "2026-08-24T08:00:00Z",
  "author_agent": "gautama-feature-builder",
  "modified_packages": [
    "internal/auditor",
    "cmd/graphify-ast-audit"
  ],
  "exported_symbols_added": [
    "auditor.SelectorEvaluatorConfig",
    "auditor.NewDeepSelectorEvaluator"
  ],
  "atomic_store_verified": true,
  "path_confinement_verified": true,
  "ready_for_guard_testing": true
}
```
