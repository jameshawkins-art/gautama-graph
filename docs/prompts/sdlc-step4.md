[.agents/personas/nexus.md](../../.agents/personas/nexus.md) [.agents/personas/security-auditor.md](../../.agents/personas/security-auditor.md) [.agents/personas/regression-tester.md](../../.agents/personas/regression-tester.md)

## CONTEXT & OBJECTIVE
Execute Phase 5 (Security Audit & Dependency Verification Gate) and Phase 6 (Master Graph Sync & Product Release Sign-off) of the Gautama Graph Software Development Lifecycle following Phase 4 code quality and test certification.

You are tasked with running security vulnerability scans (`govulncheck`), validating CLI compilation binaries, executing master synchronization (`./scripts/graphify_sync.sh`), issuing the security verification meta-artifact, and updating status matrices in roadmap documents (`docs/roadmap/roadmap.md`).

---

## 🔒 STEP 4 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)), `security-auditor` ([@security-auditor.md](../../.agents/personas/security-auditor.md)), `regression-tester` ([@regression-tester.md](../../.agents/personas/regression-tester.md)).
- **Injected Skills**: `skills/security-auditor` ([SKILL.md](../../.agents/skills/security-auditor/SKILL.md)).
- **Excluded Context**: Raw Go implementation files, intermediate compiler debug logs, scratch test scripts.

---

## 🛑 PHASE 5 & 6 EXECUTION AUTHORITY
1. **Exclusive Release Matrix Update Authority**: Updating task/feature status matrices in `./docs/roadmap/` (including both target feature roadmaps and the master product roadmap `docs/roadmap/roadmap.md`) is strictly restricted to Step 4 / Phase 5 & 6 (`@nexus.md`).
2. **Status Title Suffix Requirement**: When updating release status matrices, completed feature or task titles MUST be suffixed with `(status version)` (e.g. `(🟢 COMPLETED V1.0)`).
3. **Master Roadmap Synchronization**: Agents MUST update `docs/roadmap/roadmap.md` to reflect task progress metrics, updated status badges, task table entries, and overall milestone completion states.

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Security Verification & Dependency Audit (`@security-auditor.md`, `@nexus.md`)
- Run vulnerability scanning: `govulncheck ./...` to verify zero known vulnerabilities in transitive dependencies.
- Verify zero `unsafe.Pointer` or unvetted cgo usage across the codebase.
- Verify all file operations strictly assert workspace path confinement.
- Emit machine-readable certificate `security_verification_meta.json`.

### 2. Compilation & Master Knowledge Graph Sync (`@nexus.md`, `@regression-tester.md`)
- Validate compilation: `GOWORK=off go build -o /dev/null ./cmd/graphify-ast-audit/main.go` and `GOWORK=off go build -o /dev/null ./cmd/graphify-doc-audit/main.go`.
- Execute master synchronization script: `./scripts/graphify_sync.sh` (or `graphify update .` -> `graphify-ast-audit` -> `graphify-doc-audit`).
- Verify that `graphify-out/graph.json` and `graphify-out/doc_graph_audit.json` are cleanly generated with atomic persistence.

### 3. Release Matrix Update & Production Synchronization (`@nexus.md`)
- Update target feature release matrices inside `./docs/roadmap/` (e.g., `docs/roadmap/<feature-roadmap-name>.md`), applying the mandatory `(🟢 COMPLETED V1.0)` title suffix.
- Synchronize the Master Product Roadmap document (`docs/roadmap/roadmap.md`), updating completion progress metrics, task progression status tables, and global release matrix badges.

---

## 📄 OUTPUT FILE REQUIREMENT
Provide clickable file links to `security_verification_meta.json`, `docs/roadmap/roadmap.md`, and the updated feature specification document.
