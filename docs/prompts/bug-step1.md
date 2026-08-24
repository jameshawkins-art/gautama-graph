[.agents/personas/nexus.md](../../.agents/personas/nexus.md) [.agents/personas/debugger-remediation.md](../../.agents/personas/debugger-remediation.md)

## CONTEXT & OBJECTIVE
Execute Phase B1 (Defect Capture, Triage & Minimal Repro Gate) of the Gautama Graph Bug Remediation Lifecycle for a reported issue, runtime panic, IPC bridge crash, doc link audit failure, or parser lockup.

You are tasked with analyzing the defect symptoms, tracing execution paths using Graphify, constructing an isolated minimal reproduction test case (in `internal/auditor/auditor_test.go` or `testdata/repro/`), and producing a formal Bug Specification Document in `docs/bugs/` (e.g. `docs/bugs/bug-<description>-<id>.md`).

---

## 🔒 STEP B1 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)), `debugger-remediation` ([@debugger-remediation.md](../../.agents/personas/debugger-remediation.md)).
- **Injected Skills**: `skills/debugger-remediation` ([SKILL.md](../../.agents/skills/debugger-remediation/SKILL.md)).
- **Excluded Context**: Feature roadmap specs, public API extensions, full release test suites.

---

## 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH DISCOVERY (TOKEN OPTIMIZATION)
The triage team MUST query `graphify path "<Caller>" "<Callee>"`, `graphify explain "<type or function>"`, or `graphify query "<component>"` to locate exact error origin sites and upstream/downstream callers with minimal token overhead instead of broad greps.

---

## 🛑 PHASE B1 EXECUTION CONSTRAINTS
1. **Reproduction Before Fix**: Writing production code fixes or refactoring Go packages is strictly forbidden in Phase B1. The goal is reproduction and triage only.
2. **Deterministic Repro Harness**: The bug MUST be captured in an executable, isolated unit/integration test case that consistently reproduces the failure.
3. **Artifact Output**: The defect report must be saved directly to `docs/bugs/bug-<description>-<id>.md`.
4. **Phase Boundary Rule**: Agents MUST stop after Phase B1 and wait for user invocation of Bug Step 2 / Phase B2 (`execute docs/prompts/bug-step2.md with docs/bugs/bug-<description>-<id>.md`).

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Defect Classification & Impact Analysis (`@debugger-remediation.md`, `@nexus.md`)
- Capture incident symptoms: stack trace, error logs, affected packages (`internal/auditor`, `python/ast_auditor_bridge.py`, `cmd/`).
- Classify fault category: AST parsing syntax fault, IPC bridge deadlock, circular doc link loop, path boundary escape attempt, concurrent mutex violation, or memory leak.
- Determine defect severity (Critical, High, Medium, Low) and system blast radius.

### 2. Deterministic Minimal Reproduction Case (`@debugger-remediation.md`)
- Author an isolated reproduction test case in `internal/auditor/auditor_test.go` or `testdata/repro/` (e.g. `TestRepro_FaultName(t *testing.T)`).
- Execute the test to prove failure (Red state) and capture empirical failure logs.

### 3. Bug Specification Document (`docs/bugs/bug-<description>-<id>.md`)
- **Metadata**: Title, Bug ID, Category, Severity, Status `(🔴 OPEN / TRIAGED)`.
- **Symptom & Stack Trace**: Full runtime logs and error signatures.
- **Isolated Repro Code**: Snippet of the failing test case.
- **Hypothesized Root Cause**: Technical analysis of the failure mechanism.

---

## 📄 OUTPUT FILE REQUIREMENT
Save the bug specification to `docs/bugs/bug-<description>-<id>.md` (e.g. `docs/bugs/bug-ipc-bridge-deadlock-001.md`) and provide the exact clickable file link in your response.
