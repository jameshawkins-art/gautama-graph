[.agents/personas/nexus.md](../../.agents/personas/nexus.md) [.agents/personas/feature-engineer.md](../../.agents/personas/feature-engineer.md) [.agents/personas/security-auditor.md](../../.agents/personas/security-auditor.md) [.agents/personas/regression-tester.md](../../.agents/personas/regression-tester.md)

# Master Engineering Roadmap Formulation Directive (Pre-SDLC)

## Context & Operational Mandate
You are executing the **Pre-SDLC Initial Roadmap Formulation** for **Gautama Graph** (`github.com/jameshawkins-art/gautama-graph`), the core AST verification and documentation integrity engine for the Gautama ecosystem.

Reference the active project personas:
- **Lead AI Workflow Architect & Governor ([@nexus.md](../../.agents/personas/nexus.md))**: Master orchestrator, Prompt Ops director, Dynamic Scaffolding governor, and release status gatekeeper.
- **Feature Engineer ([@feature-engineer.md](../../.agents/personas/feature-engineer.md))**: Go 1.26+ AST parsing engine, selector evaluators, Markdown doc-graph topology, and CLI commands.
- **Security & Compliance Auditor ([@security-auditor.md](../../.agents/personas/security-auditor.md))**: Zero-trust path traversal defense, unsafe memory bans, subprocess sanitization, and govulncheck.
- **Regression & Test Guard ([@regression-tester.md](../../.agents/personas/regression-tester.md))**: Table-driven test suites, race detection (`-race`), coverage gates ($\ge 85\%$), and parser fuzzing.

---

## 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)), `feature-engineer` ([@feature-engineer.md](../../.agents/personas/feature-engineer.md)), `security-auditor` ([@security-auditor.md](../../.agents/personas/security-auditor.md)).
- **Excluded Context**: Production code implementations, runtime bug remediation scripts.

---

## 🕸️ Mandatory Graphify Knowledge Graph Discovery (Token Optimization)
Before performing raw file reads, the roadmap team MUST query:
- `graphify query "AST"` to inspect existing engine packages and selector evaluators.
- `graphify query "DocAuditor"` to inspect markdown link graph parsing.
- `graphify path "<A>" "<B>"` to trace dependencies between `cmd/`, `internal/auditor/`, and `python/ast_auditor_bridge.py`.
- Navigate `graphify-out/wiki/index.md` to map subsystem architecture with minimal token overhead.

---

## 🛑 STRICT PRE-SDLC EXECUTION CONSTRAINTS
1. **No SDLC Kick-Off**: This prompt is strictly PRE-SDLC. Do **NOT** generate requirements specification files in `docs/specs/` (e.g. `docs/specs/<NNN>-<feature>-requirements.md`).
2. **No Implementation Code or Blueprints**: Writing Go source code, Python scripts, test files (`*_test.go`), or Phase 2 technical blueprints is strictly forbidden during roadmap formulation.
3. **Target File ONLY**: The output must be written directly to [`docs/roadmap/roadmap.md`](file:///home/slvr/source/gautama-graph/docs/roadmap/roadmap.md).

---

## 📋 REQUIRED DELIVERABLES & OUTPUT FORMAT

Generate a comprehensive, structured Master Product Roadmap document saved to [`docs/roadmap/roadmap.md`](file:///home/slvr/source/gautama-graph/docs/roadmap/roadmap.md) containing:
1. **Executive Summary**: Strategic mission of Gautama Graph (deterministic code AST auditing, phantom edge pruning, Markdown link graph integrity, atomic persistence).
2. **Architecture & Subsystem Alignment**: Overview of `internal/auditor/` (Engine, Parser, Evaluator, PythonBridge, Store, DocAuditor), `cmd/`, and `python/ast_auditor_bridge.py`.
3. **Master Feature Roadmap Table**:
   - Sequential 3-digit zero-padded sequence codes (`001`, `002`...).
   - Feature Title, Lead Personas, Target Milestone, Status `(🔴 NOT STARTED)`.
4. **Phased Milestone Breakdown**:
   - Milestone 1: Core AST Engine & Provenance Model
   - Milestone 2: Markdown Doc-Graph Topology & Orphan Detection
   - Milestone 3: Python IPC Bridge Hardening & Subprocess Lifecycle
   - Milestone 4: CI/CD, Fuzzing & High-Throughput Batch Auditing
5. **Quality & Security Guardrails**: Coverage gates ($\ge 85\%$), race detector enforcement, zero-trust path containment, atomic write guarantees.

---

## 📄 OUTPUT FILE REQUIREMENT
Save the finalized master roadmap to `docs/roadmap/roadmap.md` and provide the exact clickable file link in your response.
