[.agents/personas/feature-engineer.md](../../.agents/personas/feature-engineer.md) [.agents/personas/regression-tester.md](../../.agents/personas/regression-tester.md) [.agents/personas/security-auditor.md](../../.agents/personas/security-auditor.md) [.agents/personas/nexus.md](../../.agents/personas/nexus.md)

## CONTEXT & OBJECTIVE
Execute Phase 3 (Deterministic Code Implementation) and Phase 4 (Regression & SQA Verification Gate) of the Gautama Graph Software Development Lifecycle for the approved technical blueprint in `docs/specs/` (e.g. `docs/specs/<NNN>-<feature-name>-architecture-blueprint.md`).

You are tasked with writing clean, robust Go 1.26+ code, creating adjacent table-driven `*_test.go` suites, executing race detection (`-race`), ensuring test coverage $\ge 85\%$, running AST relationship audits, and generating the `walkthrough.md` verification artifact and `feature_delivery.json` meta-artifact.

---

## 🔒 STEP 3 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `feature-engineer` ([@feature-engineer.md](../../.agents/personas/feature-engineer.md)), `regression-tester` ([@regression-tester.md](../../.agents/personas/regression-tester.md)), `security-auditor` ([@security-auditor.md](../../.agents/personas/security-auditor.md)), `nexus` ([@nexus.md](../../.agents/personas/nexus.md)).
- **Injected Skills**: `skills/feature-engineer` ([SKILL.md](../../.agents/skills/feature-engineer/SKILL.md)), `skills/regression-tester` ([SKILL.md](../../.agents/skills/regression-tester/SKILL.md)).
- **Excluded Context**: Pre-SDLC roadmap authoring, external documentation drafting, release status matrix updating.

---

## 🕸️ MANDATORY GRAPHIFY DISCOVERY & POST-IMPLEMENTATION SYNC
1. **Graphify Discovery (Token Optimization)**: Query `graphify query "<feature>"` or `graphify explain` to locate exact symbol definitions and package boundaries before editing code.
2. **Post-Implementation Sync**: After modifying code, the agent team MUST run `graphify update .` followed by `go run cmd/graphify-ast-audit/main.go` to prune phantom edges and verify graph integrity (`graphify-out/graph.json`).

---

## 🛑 PHASE 3 & 4 EXECUTION CONSTRAINTS
1. **Forbidden Release Actions**: Modifying project release status matrices in roadmap documents (`docs/roadmap/roadmap.md`) or claiming Phase 5/6 release completion is strictly forbidden in Step 3.
2. **Mandatory Test Suite Execution**: `GOWORK=off go test -timeout 30s -v -race ./internal/auditor/...` MUST be executed and pass 100% with zero race conditions.
3. **Artifact Output**: Create or update the `walkthrough.md` artifact summarizing code changes, unit test results, and coverage metrics, and emit `feature_delivery.json`.
4. **Phase Boundary Rule**: Completing Step 3 / Phase 3 & 4 certifies local implementation and regression verification ONLY. Agents MUST stop after Phase 4 and wait for explicit user invocation of Step 4 / Phase 5 & 6 (`execute docs/prompts/sdlc-step4.md with docs/specs/<NNN>-<feature-name>-architecture-blueprint.md`) in a subsequent prompt.

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Hardened Go Engine Implementation (`@feature-engineer.md`, `@nexus.md`)
- Implement Go 1.26+ AST parsing, selector evaluation, doc-graph analysis, and CLI commands in `internal/auditor/` and `cmd/`.
- Enforce PascalCase naming with godoc comments on all exported identifiers.
- Enforce the two-phase atomic write protocol (`.tmp` buffer + `os.Rename`) on all graph store mutations.
- Enforce zero-trust path containment checks (`filepath.Clean`, `filepath.Abs`, root prefix validation).
- Enforce subprocess hygiene: `exec.CommandContext` with bounded timeouts, discrete argument slices, stream separation, and mandatory `scanner.Err()` checks.

### 2. Table-Driven Test Suites & Regression Verification (`@regression-tester.md`)
- Build matching table-driven unit and integration tests in `internal/auditor/*_test.go`.
- Test boundary cases: empty graphs, malformed AST syntax, circular markdown links, path traversal attacks, and context cancellations.
- Execute test suites with race detection: `GOWORK=off go test -timeout 30s -v -race ./internal/auditor/...`.
- Assert statement coverage $\ge 85\%$ on modified packages.
- Implement parser fuzzing harnesses (`func Fuzz...`) where applicable.

### 3. SQA Certification & Deliverable Handoff (`@regression-tester.md`, `@security-auditor.md`, `@nexus.md`)
- Run `go run cmd/graphify-ast-audit/main.go` and `go run cmd/graphify-doc-audit/main.go` to verify CLI utility functionality.
- Summarize code diffs, security controls verified, and test outputs in `walkthrough.md`.
- Emit machine-readable meta-artifact `feature_delivery.json` for Phase 5 & 6 consumption.

---

## 📄 OUTPUT FILE REQUIREMENT
Provide clickable file links to the modified Go source files, unit test files, `walkthrough.md`, and `feature_delivery.json`.
