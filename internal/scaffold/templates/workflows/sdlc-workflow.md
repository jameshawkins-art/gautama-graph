# SDLC Feature Engineering Workflow Guide

This workflow governs feature engineering and capability expansion across 4 sequential gates with integrated security auditing and strict Antigravity 2.0 context isolation.

---

## Step-by-Step Context Isolation Matrix

To prevent **context rot** and token bloat, each workflow step strictly isolates the active persona, skill, and input/output artifacts. Never load inactive personas or broad file trees into a single step's context.

| Step | Gate Name | Active Persona | Injected Skill | Injected Context / Input | Output Meta-Artifact | Excluded from Context |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **`sdlc-step1`** | Requirements & API Spec | `@nexus.md`, Feature Engineer | Feature Engineering | `docs/specs/*-requirements.md` template | `feature_spec.md` | Test runners, security scanners, mechanic |
| **`sdlc-step2`** | Implementation & Code Gen | Feature Engineer | Feature Engineering | `feature_spec.md` + targeted source files | `feature_delivery.json` | Guard, Gatekeeper, fuzzing harnesses |
| **`sdlc-step3`** | Regression & Coverage Gate | Regression Tester | Regression Testing | `feature_delivery.json` + test files | `test_verification_meta.json` | Production logic modification, mechanic |
| **`sdlc-step4`** | Security Audit & Release | Security Auditor, `@nexus.md` | Security Auditor | `test_verification_meta.json` | `security_verification_meta.json` | Test runners, unvetted source files |

---

## Step Execution Protocol

### Step 1: Feature Inception (`sdlc-step1`)
- **Action**: `@nexus.md` invokes the feature engineer persona using `graphify query` to inspect existing symbol contracts.
- **Output**: Generates `docs/specs/<NNN>-<feature-name>-requirements.md` with interface specifications and acceptance criteria.

### Step 2: Implementation (`sdlc-step2`)
- **Action**: Feature engineer writes code, enforcing atomic file writes, context propagation, and discrete subprocess execution.
- **Output**: Emits `feature_delivery.json`.

### Step 3: Regression & Fuzzing Gate (`sdlc-step3`)
- **Action**: Regression tester runs test suites with race detection, executes test harnesses, and asserts statement coverage thresholds.
- **Output**: Emits `test_verification_meta.json`.

### Step 4: Security Audit & Release Gate (`sdlc-step4`)
- **Action**: Security auditor verifies zero-trust path containment, validates dependencies, executes graph sync, and signs off with `@nexus.md`.
- **Output**: Emits `security_verification_meta.json` and closes the feature.
