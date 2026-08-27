---
trigger: always_on
description: Enforces the mandatory Red-Green-Refactor TDD cycle and Production Call-Site Invariant to prevent test-only function omissions and phantom test passes.
---

## Test-Driven Development (TDD) Red-Green-Refactor & Production Integration Invariant

All feature development (`sdlc-step3.md`) and bug remediation (`bug-step3.md`) in this project MUST strictly adhere to the 3-stage Test-Driven Development (TDD) cycle and the Production Call-Site Invariant.

---

### 1. 🔴 RED STAGE: Failing Test First
- **Mandate**: Author the automated test in an adjacent test suite file (e.g., `*_test.go`, `test_*.py`, `*.spec.ts`, `*.test.ts`) **BEFORE** writing or modifying production implementation code.
- **Integration Boundary Rule**: The test MUST target the **actual production entrypoint / call site / HTTP endpoint / handler / service / pipeline stage**, in addition to any granular helper unit tests.
  - *Example*: When introducing a new utility function (e.g. atomic persister, cryptographic validator, data parser), do not merely write isolated unit tests for the utility. You MUST write/update tests asserting that the consumer (handler, service, exporter, or controller) uses the utility in its nominal and error workflows.
- **Failure Verification**: Run the project's test command (e.g. `go test -v -run <TestName> ...`, `pytest -k <test_name>`, `npm test -- -t <TestName>`) and verify that the test **FAILS** for the expected reason prior to writing production code.

---

### 2. 🟢 GREEN STAGE: Minimal Production Code & Production Call-Site Wiring
- **Mandate**: Write the minimal, idiomatic code necessary to pass the failing tests.
- **Production Call-Site Invariant**:
  - Whenever a new helper function, struct/class method, or utility is created in production source directories (e.g. `internal/`, `src/`, `lib/`, `pkg/`), the agent **MUST IMMEDIATELY WIRE IT INTO ITS PRODUCTION CALL PATH** (e.g. within handlers, services, repositories, or pipeline stages).
  - **PROHIBITED ANTI-PATTERN**: Implementing a new utility function and only calling it inside test files while leaving production code using duplicate inline logic or obsolete primitives is strictly classified as a **CRITICAL INCOMPLETE INTEGRATION DEFECT**.
- **Pass Verification**: Execute the test command to confirm that the test transitions from RED to GREEN.

---

### 3. 🔵 REFACTOR & INTEGRATION AUDIT STAGE: DRY & AST Verification
- **Mandate**: Clean up, optimize, and eliminate code duplication across the entire subsystem without breaking existing behavior.
- **De-duplication (DRY)**: Search the workspace for any duplicate inline implementations of the new utility and refactor them to call the centralized function.
- **AST Production Caller Verification**:
  - Verify that every newly declared function has **$\ge 1$ active call site in production non-test code** ($C_{prod} > 0$).
  - Use `graphify query`, `graphify path`, or AST auditing tools (`make audit-ast`, `go run cmd/graphify-ast-audit/main.go`) to confirm verified AST caller edges.
- **Full Suite Regression & Graph Sync**:
  - Run the full workspace test suite (e.g. `go test -timeout 30s -v ./...`, `npm test`, `pytest`) to certify zero regressions or race conditions.
  - Run `make graphify-update` (or `./scripts/graphify_sync.sh`) to synchronize AST symbols and eliminate phantom edges.
