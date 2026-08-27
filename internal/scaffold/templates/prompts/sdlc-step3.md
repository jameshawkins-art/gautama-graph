[.agents/personas/nexus.md](../../.agents/personas/nexus.md) [.agents/rules/tdd-cycle.md](../../.agents/rules/tdd-cycle.md)

## CONTEXT & OBJECTIVE
Execute Phase 3 (Deterministic Code Implementation) and Phase 4 (Regression & SQA Verification Gate) of the Software Development Lifecycle for the approved technical blueprint in `docs/specs/` (e.g. `docs/specs/<NNN>-<feature-name>-architecture-blueprint.md`).

You are tasked with writing clean, robust code, creating comprehensive test suites, executing race/concurrency detection, verifying test coverage thresholds, running graph integrity audits, and generating the `walkthrough.md` verification artifact and `feature_delivery.json` meta-artifact.

---

## 🔒 STEP 3 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)), feature engineers, and test guards.
- **Injected Skills**: Feature engineering, test automation, and regression verification skills.
- **Excluded Context**: Pre-SDLC roadmap authoring, external documentation drafting, release status matrix updating.

---

## 🕸️ MANDATORY GRAPHIFY DISCOVERY & POST-IMPLEMENTATION SYNC
1. **Graphify Discovery (Token Optimization)**: Query `graphify query "<feature>"` or `graphify explain` to locate exact symbol definitions and package boundaries before editing code.
2. **Post-Implementation Sync**: After modifying code, the agent team MUST run `make graphify-update` or `graphify update .` to update graph nodes and verify graph integrity (`graphify-out/graph.json`).

---

## 🛑 PHASE 3 & 4 EXECUTION CONSTRAINTS
1. **Test-Driven Development (TDD) & Production Call-Site Invariant ([.agents/rules/tdd-cycle.md](../../.agents/rules/tdd-cycle.md))**:
   - **Red Stage**: Author failing automated tests targeting the true entrypoint/caller before writing production code.
   - **Green Stage**: Write minimal production code and **wire all newly declared utility functions directly into production caller paths**.
   - **Refactor Stage**: Eliminate duplicate inline logic across packages, verify AST caller graph connections ($C_{prod} > 0$), and ensure zero test-only orphaned utilities.
2. **Forbidden Release Actions**: Modifying project release status matrices in roadmap documents (`docs/roadmap/roadmap.md`) or claiming release completion is strictly forbidden in Step 3.
3. **Mandatory Test Suite Execution**: The project's automated test suite MUST be executed and pass 100% with zero race conditions or uncaught errors.
4. **Artifact Output**: Create or update the `walkthrough.md` artifact summarizing code changes, test results, and coverage metrics, and emit `feature_delivery.json`.
5. **Phase Boundary Rule**: Completing Step 3 certifies local implementation and regression verification ONLY. Agents MUST stop after Phase 4 and wait for explicit user invocation of Step 4 (`execute docs/prompts/sdlc-step4.md with docs/specs/<NNN>-<feature-name>-architecture-blueprint.md`) in a subsequent prompt.

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Hardened Implementation & Production Call-Site Wiring
- Implement the requested modules, interfaces, and CLI commands according to the architecture blueprint.
- Wire all newly authored utilities directly into their production caller paths (handlers, services, repositories, pipelines).
- Eliminate duplicate inline logic across packages in favor of centralized utilities (DRY invariant).
- Enforce standard naming, documentation, and error handling conventions.
- Enforce the two-phase atomic write protocol (`.tmp` buffer + atomic rename) on state mutations.
- Enforce zero-trust path containment checks.
- Enforce subprocess hygiene: bounded timeouts, discrete argument slices, and stream separation.

### 2. Table-Driven Test Suites & Regression Verification
- Build matching unit and integration test suites.
- Test boundary cases: empty inputs, malformed syntax, edge dependencies, path traversal attempts, and context cancellations.
- Execute test suites with race detection and verify test coverage thresholds.

### 3. SQA Certification & Deliverable Handoff
- Run project audit scripts or doc/code verification tools.
- Summarize code diffs, security controls verified, and test outputs in `walkthrough.md`.
- Emit machine-readable meta-artifact `feature_delivery.json` for subsequent release gate consumption.

---

## 📄 OUTPUT FILE REQUIREMENT
Provide clickable file links to the modified source files, unit test files, `walkthrough.md`, and `feature_delivery.json`.
