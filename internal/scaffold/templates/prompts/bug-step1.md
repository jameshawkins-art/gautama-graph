[.agents/personas/nexus.md](../../.agents/personas/nexus.md)

## CONTEXT & OBJECTIVE
Execute Phase B1 (Defect Capture, Triage & Minimal Repro Gate) of the Bug Remediation Lifecycle for a reported issue, runtime panic, crash, broken dependency, or lockup.

You are tasked with analyzing the defect symptoms, tracing execution paths using Graphify, constructing an isolated minimal reproduction test case, and producing a formal Bug Specification Document in `docs/bugs/` (e.g. `docs/bugs/bug-<description>-<id>.md`).

---

## 🔒 STEP B1 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)) and the debugger/remediation persona.
- **Injected Skills**: Debugging, remediation, and error isolation skills.
- **Excluded Context**: Feature roadmap specs, public API extensions, full release test suites.

---

## 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH DISCOVERY (TOKEN OPTIMIZATION)
The triage team MUST query `graphify path "<Caller>" "<Callee>"`, `graphify explain "<type or function>"`, or `graphify query "<component>"` to locate exact error origin sites and upstream/downstream callers with minimal token overhead instead of broad greps.

---

## 🛑 PHASE B1 EXECUTION CONSTRAINTS
1. **Reproduction Before Fix**: Writing production code fixes or broad refactorings is strictly forbidden in Phase B1. The goal is reproduction and triage only.
2. **Deterministic Repro Harness**: The bug MUST be captured in an executable, isolated unit/integration test case that consistently reproduces the failure.
3. **Artifact Output**: The defect report must be saved directly to `docs/bugs/bug-<description>-<id>.md`.
4. **Phase Boundary Rule**: Agents MUST stop after Phase B1 and wait for user invocation of Bug Step 2 (`execute docs/prompts/bug-step2.md with docs/bugs/bug-<description>-<id>.md`).

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Defect Classification & Impact Analysis
- Capture incident symptoms: stack trace, error logs, and affected components.
- Classify fault category: syntax fault, IPC deadlock, circular link loop, path boundary escape attempt, concurrency violation, or resource leak.
- Determine defect severity (Critical, High, Medium, Low) and system blast radius.

### 2. Deterministic Minimal Reproduction Case
- Author an isolated reproduction test case in the test suite.
- Execute the test to prove failure (Red state) and capture empirical failure logs.

### 3. Bug Specification Document (`docs/bugs/bug-<description>-<id>.md`)
- **Metadata**: Title, Bug ID, Category, Severity, Status `(🔴 OPEN / TRIAGED)`.
- **Symptom & Stack Trace**: Full runtime logs and error signatures.
- **Isolated Repro Code**: Snippet of the failing test case.
- **Hypothesized Root Cause**: Technical analysis of the failure mechanism.

---

## 📄 OUTPUT FILE REQUIREMENT
Save the bug specification to `docs/bugs/bug-<description>-<id>.md` (e.g. `docs/bugs/bug-resource-leak-001.md`) and provide the exact clickable file link in your response.
