# Gautama Graph SDLC Workflow Guide

This workflow governs feature engineering and public API expansion across 4 sequential gates with integrated security auditing and strict Antigravity 2.0 context isolation.

---

## Step-by-Step Context Isolation Matrix

To prevent **context rot** and token bloat, each workflow step strictly isolates the active persona, skill, and input/output artifacts. Never load inactive personas or broad file trees into a single step's context.

| Step | Gate Name | Active Persona | Injected Skill | Injected Context / Input | Output Meta-Artifact | Excluded from Context |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **`sdlc-step1`** | Requirements & API Spec | `@nexus.md`, `@feature-engineer.md` | `skills/feature-engineer` | `docs/specs/*-requirements.md` template | `feature_spec.md` | Test runners, security scanners, mechanic |
| **`sdlc-step2`** | Implementation & Code Gen | `@feature-engineer.md` | `skills/feature-engineer` | `feature_spec.md` + targeted Go files | `feature_delivery.json` | Guard, Gatekeeper, fuzzing harnesses |
| **`sdlc-step3`** | Regression & Coverage Gate | `@regression-tester.md` | `skills/regression-tester` | `feature_delivery.json` + `*_test.go` | `test_verification_meta.json` | Production logic modification, mechanic |
| **`sdlc-step4`** | Security Audit & Release | `@security-auditor.md`, `@nexus.md` | `skills/security-auditor` | `test_verification_meta.json` | `security_verification_meta.json` | Test runners, unvetted source files |

---

## Step Execution Protocol

### Step 1: Feature Inception (`sdlc-step1`)
- **Action**: `@nexus.md` invokes `@feature-engineer.md` using `graphify query` to inspect existing symbol contracts.
- **Output**: Generates `docs/specs/<NNN>-<feature-name>-requirements.md` with PascalCase Go interfaces and Godoc signatures.

### Step 2: Implementation (`sdlc-step2`)
- **Action**: `@feature-engineer.md` writes Go 1.26+ code in `internal/auditor/` and `cmd/`, enforcing `.tmp` staging buffers with `os.Rename`, context propagation, and discrete subprocess execution.
- **Output**: Emits `feature_delivery.json`.

### Step 3: Regression & Fuzzing Gate (`sdlc-step3`)
- **Action**: `@regression-tester.md` runs table-driven test suites with race detection (`go test -v -race ./internal/auditor/...`), executes fuzzing harnesses, and asserts $\ge 85\%$ statement coverage.
- **Output**: Emits `test_verification_meta.json`.

### Step 4: Security Audit & Release Gate (`sdlc-step4`)
- **Action**: `@security-auditor.md` verifies zero-trust path containment, asserts zero `unsafe`/cgo usage, runs `govulncheck ./...`, executes `./scripts/graphify_sync.sh`, and signs off with `@nexus.md`.
- **Output**: Emits `security_verification_meta.json` and closes the feature.
