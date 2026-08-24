# Gautama Graph Bug Remediation Workflow Guide

This workflow governs incident triage and defect resolution across 3 sequential gates with reproduction-first enforcement, test-driven remediation, and strict context isolation.

---

## Step-by-Step Context Isolation Matrix

To prevent **context rot**, each bug remediation step isolates active agents and inputs to the minimum required for the task.

| Step | Gate Name | Active Persona | Injected Skill | Injected Context / Input | Output Meta-Artifact | Excluded from Context |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **`bug-step1`** | Defect Triage & Isolation | `@nexus.md`, `@debugger-remediation.md` | `skills/debugger-remediation` | Incident report / stack trace | `repro_case.go` + `docs/bugs/bug-*.md` | Production code refactoring, security scanner |
| **`bug-step2`** | Surgical Remediation | `@debugger-remediation.md` | `skills/debugger-remediation` | `repro_case.go` + faulting source file | `remediation_meta.json` | Feature engineering, unrelated source files |
| **`bug-step3`** | Regression & Security Gate | `@regression-tester.md`, `@security-auditor.md`, `@nexus.md` | `skills/regression-tester`, `skills/security-auditor` | `remediation_meta.json` + test suites | `security_verification_meta.json` | Mechanic refactoring, unconstrained source reads |

---

## Step Execution Protocol

### Step 1: Defect Capture & Isolation (`bug-step1`)
- **Action**: `@nexus.md` and `@debugger-remediation.md` trace faulting paths using `graphify path` and create a standalone minimal reproduction test case (`internal/auditor/auditor_test.go` or `testdata/repro/`).
- **Output**: Generates `docs/bugs/bug-<description>-<id>.md` and failing repro test.

### Step 2: Surgical Remediation (`bug-step2`)
- **Action**: `@debugger-remediation.md` delivers minimal, non-breaking code patches confined strictly to the faulting execution path, resolves pipe buffer deadlocks, ensures `scanner.Err()` and `mu.Unlock()` deferrals.
- **Output**: Emits `remediation_meta.json`.

### Step 3: Regression & Security Gate (`bug-step3`)
- **Action**: `@regression-tester.md` verifies the repro test passes and runs the full test suite under `-race`. `@security-auditor.md` verifies path containment and dependency stability. Runs `./scripts/graphify_sync.sh`.
- **Output**: Marks defect `CLOSED (🟢 RESOLVED)`.
