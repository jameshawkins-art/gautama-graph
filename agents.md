# Gautama Graph: Antigravity 2.0 Multi-Agent Team Manifest

This directory manifest configures the declarative orchestration routes for the `gautama-graph` subsystem.

## Team Layout
* **nexus**: Master Orchestrator, Prompt Ops Director & Dynamic Scaffolding Governor.
* **gautama-builder**: Feature Engineer Agent optimized for Go AST structures.
* **gautama-mechanic**: Debugger & Remediation Agent handling IPC bridges.
* **gautama-guard**: Regression & Test Automation Agent asserting code metrics.
* **gautama-gatekeeper**: Security & Compliance Auditor enforcing isolation boundaries.

## Routing Manifest
```yaml
---
orchestrator: nexus
version: 2.0.0
registry:
  - persona: .agents/personas/nexus.md
  - persona: .agents/personas/feature-engineer.md
  - persona: .agents/personas/debugger-remediation.md
  - persona: .agents/personas/regression-tester.md
  - persona: .agents/personas/security-auditor.md
---
```

## Lifecycle Handoff Schema
| Stage | Origin Agent | Target Agent | Handoff Artifact | Validation Gate |
| :--- | :--- | :--- | :--- | :--- |
| **1. Inception** | `nexus` | `gautama-builder` | `feature_spec.md` | Interface definition & Godoc check |
| **2. Code Gen** | `gautama-builder` | `gautama-guard` | `feature_delivery.json` | PascalCase & Atomic persistence check |
| **3. Test Gate** | `gautama-guard` | `gautama-gatekeeper` | `test_verification_meta.json` | Coverage >= 85% & zero test failures |
| **4. Fault Triage**| `gautama-guard` | `gautama-mechanic` | `remediation_meta.json` | Fault isolation & Repro test case |
| **5. Security Gate**| `gautama-gatekeeper` | `nexus` | `security_verification_meta.json` | Path containment, zero vuln, clean lint |
| **6. Release Gate**| `nexus` | `Production / Main` | `release_manifest.json` | Final sign-off & sync completion |
