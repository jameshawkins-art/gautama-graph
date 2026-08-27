---
trigger: always_on
description: Enforces the mandatory Red-Green-Refactor TDD cycle and Production Call-Site Invariant to prevent test-only function omissions and phantom test passes.
---

## Test-Driven Development (TDD) Red-Green-Refactor & Production Integration Invariant

All feature development (`sdlc-step3.md`) and bug remediation (`bug-step3.md`) in Gautama Graph MUST strictly adhere to the 3-stage Test-Driven Development (TDD) cycle and the Production Call-Site Invariant.

---

### 1. 🔴 RED STAGE: Failing Test First
- **Mandate**: Author the automated test in an adjacent `*_test.go` file (or Python test runner in `python/tests/`) **BEFORE** writing or modifying production implementation code.
- **Integration Boundary Rule**: The test MUST target the **actual production call site / CLI command / parser / auditor service / IPC bridge**, in addition to any granular helper unit tests.
  - *Example*: When introducing a new AST traversal or selector evaluator utility, do not merely write isolated helper unit tests. You MUST write/update tests asserting that the engine (`internal/auditor/engine.go`), CLI commands (`cmd/graphify-ast-audit/`), or doc remediators (`internal/auditor/doc_remediator.go`) consume the utility in their nominal and error workflows.
- **Failure Verification**: Run the test command (`GOWORK=off go test -v -run <TestName> ./...`) and verify that the test **FAILS** for the expected reason prior to writing production code.

---

### 2. 🟢 GREEN STAGE: Minimal Production Code & Production Call-Site Wiring
- **Mandate**: Write the minimal, idiomatic Go/Python code necessary to pass the failing tests.
- **Production Call-Site Invariant**:
  - Whenever a new helper function, struct method, or utility is created in production packages (`internal/auditor/`, `internal/runner/`, `cmd/`, `python/*.py`), the agent **MUST IMMEDIATELY WIRE IT INTO ITS PRODUCTION CALL PATH** (e.g. within parsing passes, selector evaluators, store committers, or CLI execution routines).
  - **PROHIBITED ANTI-PATTERN**: Implementing a new utility function and only calling it inside `*_test.go` while leaving the production engine using duplicate inline code or obsolete primitives is strictly classified as a **CRITICAL INCOMPLETE INTEGRATION DEFECT**.
- **Pass Verification**: Execute `GOWORK=off go test -v -run <TestName> ./...` to confirm that the test transitions from RED to GREEN.

---

### 3. 🔵 REFACTOR & INTEGRATION AUDIT STAGE: DRY & AST Verification
- **Mandate**: Clean up, optimize, and eliminate code duplication across the entire subsystem without breaking existing behavior.
- **De-duplication (DRY)**: Search the workspace for any duplicate inline implementations of the new utility and refactor them to call the centralized function.
- **AST Production Caller Verification**:
  - Verify that every newly declared function has **$\ge 1$ active call site in production non-test code** ($C_{prod} > 0$).
  - Use `graphify query`, `graphify path`, or `go run cmd/graphify-ast-audit/main.go` to confirm verified AST caller edges.
- **Full Suite Regression & Graph Sync**:
  - Run the full Go test suite with race detection: `GOWORK=off go test -timeout 30s -v -race ./...`
  - Run `./scripts/graphify_sync.sh` (or `go run cmd/graphify-ast-audit/main.go` and `go run cmd/graphify-doc-audit/main.go`) to synchronize AST symbols, eliminate phantom edges, and validate documentation links.
