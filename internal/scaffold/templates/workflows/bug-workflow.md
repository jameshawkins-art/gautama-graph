# Bug Remediation Workflow Guide

This workflow governs incident triage and defect resolution across 3 sequential gates with reproduction-first enforcement, test-driven remediation, and strict context isolation.

---

## Step-by-Step Context Isolation Matrix

To prevent **context rot**, each bug remediation step isolates active agents and inputs to the minimum required for the task.

| Step | Gate Name | Active Persona | Injected Skill | Injected Context / Input | Output Meta-Artifact | Excluded from Context |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **`bug-step1`** | Defect Triage & Isolation | `@nexus.md`, Debugger/Remediation | Debugger/Remediation | Incident report / stack trace | Repro test + `docs/bugs/bug-*.md` | Production code refactoring, security scanner |
| **`bug-step2`** | Surgical Remediation | Debugger/Remediation | Debugger/Remediation | Repro test + faulting source file | `remediation_meta.json` | Feature engineering, unrelated source files |
| **`bug-step3`** | Regression & Security Gate | Regression Tester, Security Auditor, `@nexus.md` | Regression Tester, Security Auditor | `remediation_meta.json` + test suites | `security_verification_meta.json` | Mechanic refactoring, unconstrained source reads |

---

## Step Execution Protocol

### Step 1: Defect Capture & Isolation (`bug-step1`)
- **Action**: `@nexus.md` and the remediation persona trace faulting paths using `graphify path` and create a standalone minimal reproduction test case.
- **Output**: Generates `docs/bugs/bug-<description>-<id>.md` and failing repro test.

### Step 2: Surgical Remediation (`bug-step2`)
- **Action**: The remediation persona delivers minimal, non-breaking code patches confined strictly to the faulting execution path, ensuring clean resource management.
- **Output**: Emits `remediation_meta.json`.

### Step 3: Regression & Security Gate (`bug-step3`)
- **Action**: Regression tester verifies the repro test passes and runs the full test suite. Security auditor verifies path containment and dependency stability. Runs graph sync.
- **Output**: Marks defect `CLOSED (🟢 RESOLVED)`.
