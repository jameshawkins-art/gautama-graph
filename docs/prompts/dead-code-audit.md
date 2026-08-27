[.agents/personas/nexus.md](../../.agents/personas/nexus.md) [.agents/personas/feature-engineer.md](../../.agents/personas/feature-engineer.md) [.agents/personas/regression-tester.md](../../.agents/personas/regression-tester.md) [.agents/personas/security-auditor.md](../../.agents/personas/security-auditor.md)

# Unused Symbol, Dead Code & Test-Only Function Audit Directive

## Context & Operational Mandate
Execute a comprehensive AST static analysis and knowledge graph audit across the entire Gautama Graph codebase (`internal/auditor/`, `internal/runner/`, `internal/scaffold/`, `cmd/`, `python/`) to identify:
1. **Test-Only Production Functions (Incomplete Integration Defects)**: Functions, helpers, or structs implemented in production source files (`*.go`, `*.py`) that are **ONLY called from test files** (`*_test.go`, `python/tests/*.py`) with zero callers in actual production application codepaths.
2. **Orphaned / Dead Symbols**: Functions, types, constants, or variables with **0 callers anywhere** in the workspace.
3. **Duplicate Inline Re-implementations**: Production codepaths that manually implement inline logic instead of consuming existing dedicated utility functions.
4. **Unreferenced / Dead Interfaces & Exports**: Interface definitions or exported abstractions that are never implemented or referenced.

---

## 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
- **Lead Personas**:
  - **Lead AI Workflow Architect ([@nexus.md](../../.agents/personas/nexus.md))**: System gatekeeper, rule compliance, and audit synthesis.
  - **Feature Engineer ([@feature-engineer.md](../../.agents/personas/feature-engineer.md))**: Go AST symbol analysis, call graph tracing, and call-site refactoring.
  - **Regression Tester ([@regression-tester.md](../../.agents/personas/regression-tester.md))**: Test call site isolation and TDD cycle verification.
  - **Security Auditor ([@security-auditor.md](../../.agents/personas/security-auditor.md))**: Zero-trust path containment, memory safety, and dependency verification.

---

## 🕸️ Mandatory Graphify Knowledge Graph Scoping
Before scanning raw files:
1. Query `graphify query "<subsystem>"` to trace package dependency trees and symbol callers.
2. Run `graphify path "<Caller>" "<Callee>"` to verify whether intended production routes actually reach newly declared utility functions.
3. Inspect `graphify-out/graph.json` and run `go run cmd/graphify-ast-audit/main.go` to check AST verified relationships vs heuristic references.

---

## 📋 Audit Scope & Methodology

### 1. Static AST Symbol & Caller Graph Extraction
- Enumerate all declared symbols in production packages:
  - Go: `internal/**/*.go` (excluding `*_test.go`), `cmd/**/*.go`
  - Python: `python/*.py` (excluding `python/tests/*.py`)
- For each declared symbol, calculate its **Caller Matrix**:
  - **Production Caller Count ($C_{prod}$)**: Invocations from non-test production code.
  - **Test Caller Count ($C_{test}$)**: Invocations from test files (`*_test.go`, `python/tests/*`).

### 2. Defect Classification Taxonomy

| Severity | Category | Definition | Remediation Action |
| :--- | :--- | :--- | :--- |
| **CRITICAL** | **Incomplete Integration ($C_{prod} = 0, C_{test} > 0$)** | Function was authored for a feature/spec but never wired into production caller paths (only exercised in synthetic unit tests, while production handlers use duplicate inline code). | Refactor production call sites to immediately consume the centralized utility. |
| **HIGH** | **Dead Code ($C_{prod} = 0, C_{test} = 0$)** | Symbol has zero callers in both production and test code (stale remnant from previous refactorings). | Decommission and prune the dead code; run AST audit to synchronize graph. |
| **MEDIUM** | **Duplicate Inline Logic** | A production function exists, but another component re-implements the same logic inline instead of calling it. | Consolidate inline logic to call the existing function (DRY invariant). |
| **LOW / INFO** | **Exported Library API** | Public package API intentionally exported for external module consumption, adhering to ISP ($\le 3$ methods). | Verify documentation comments and confirm intended public export contract. |

---

## 🧪 Verification Protocol
1. For every **Incomplete Integration** finding:
   - Identify the production caller that should be invoking the function.
   - Refactor the production caller to consume the function.
   - Re-run `GOWORK=off go test -timeout 30s -v -race ./...`.
   - Verify that the production caller test exercises the function transitively.
2. Run `go run cmd/graphify-ast-audit/main.go` to certify that the AST relationship edge is recorded in `graphify-out/graph.json`.
3. Run `go run cmd/graphify-doc-audit/main.go` to guarantee doc link graph integrity.

---

## 📄 OUTPUT FILE REQUIREMENT
Save the finalized audit report to `docs/audits/unused-symbol-audit-report-<date>.md` containing:
- Executive Summary & Metrics (Total symbols scanned, orphaned count, test-only count).
- Itemized Finding Table with exact file links and line numbers.
- Concrete Remediation Plan and verified diffs.
