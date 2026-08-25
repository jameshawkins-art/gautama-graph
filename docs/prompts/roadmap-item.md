[.agents/personas/nexus.md](../../.agents/personas/nexus.md) [.agents/personas/feature-engineer.md](../../.agents/personas/feature-engineer.md) [.agents/personas/security-auditor.md](../../.agents/personas/security-auditor.md)

# Pre-SDLC Roadmap Item Formulation Directive

## Context & Operational Mandate
You are executing the **Pre-SDLC Roadmap Item Decomposition** for a specific capability, module extension, or architectural enhancement in **Gautama Graph**.

Reference the active project personas:
- **Lead AI Workflow Architect & Governor ([@nexus.md](../../.agents/personas/nexus.md))**: Master orchestrator, Prompt Ops director, Dynamic Scaffolding governor, and release status gatekeeper.
- **Feature Engineer ([@feature-engineer.md](../../.agents/personas/feature-engineer.md))**: Go 1.26+ AST parsing engine, selector evaluators, Markdown doc-graph topology, and CLI commands.
- **Security & Compliance Auditor ([@security-auditor.md](../../.agents/personas/security-auditor.md))**: Zero-trust path traversal defense, unsafe memory bans, subprocess sanitization, and govulncheck.

---

## 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)), `feature-engineer` ([@feature-engineer.md](../../.agents/personas/feature-engineer.md)), `security-auditor` ([@security-auditor.md](../../.agents/personas/security-auditor.md)).
- **Excluded Context**: Production code edits, unit test executions, release builds.

---

## 🕸️ Mandatory Graphify Knowledge Graph Discovery (Token Optimization)
Before performing raw file reads, `@feature-engineer.md` and `@nexus.md` MUST query `graphify query "<topic>"`, `graphify path "<A>" "<B>"`, or `graphify explain "<concept>"` to map target packages, existing interfaces, and storage paths with minimal token consumption.

---

## 🛑 STRICT PRE-SDLC EXECUTION CONSTRAINTS
1. **No SDLC Kick-Off**: This prompt is strictly PRE-SDLC. Do **NOT** generate requirements specification files in `docs/specs/` (e.g. `docs/specs/<NNN>-<feature>-requirements.md`).
2. **No Implementation Code or Blueprints**: Writing Go source code, Python scripts, test files (`*_test.go`), or Phase 2 technical blueprints is strictly forbidden during roadmap item generation.
3. **Target Files ONLY**:
   - Master Roadmap: Update [docs/roadmap/roadmap.md](../roadmap/roadmap.md) (Master Feature Table + Detailed Specifications section).
   - Dedicated Item Spec: Create [docs/roadmap/<topic>-roadmap-<NNN>.md](../roadmap/) (e.g. `docs/roadmap/ast-selector-depth-roadmap-001.md`).

---

## 📋 REQUIRED DELIVERABLES & OUTPUT FORMAT

1. **Master Roadmap Update ([docs/roadmap/roadmap.md](../roadmap/roadmap.md))**:
   - Query `docs/roadmap/` for the highest 3-digit zero-padded sequence code `NNN` and increment by 1 (`NNN+1`).
   - Append new row to the Master Feature Roadmap Table with status `(🔴 NOT STARTED)`.
   - Append detailed item summary under Detailed Item Specifications section.

2. **Dedicated Roadmap Item Document ([docs/roadmap/<topic>-roadmap-<NNN>.md](../roadmap/))**:
   - **Header & Metadata**: Document Title, Sequence Code `NNN`, Persona Drivers & Gatekeepers, Status `(🔴 NOT STARTED)`.
   - **1. Executive Summary & Strategic Objective**: Problem statement, target architecture, knowledge graph value proposition.
   - **2. Subsystem / Engine Component Matrix**: Table mapping Go packages, interfaces in `internal/auditor/types.go`, Python scripts, and storage paths using Graphify queries.
   - **3. Phased Master Task Matrix**: Task decomposition (`NNN.1`, `NNN.2`...) with driver personas, priority, estimated effort, and target SDLC phase.
   - **4. Definition of Done (DoD)**: Explicit, verifiable product manager acceptance criteria for future SDLC execution.

---

## 📄 OUTPUT FILE REQUIREMENT
Provide clickable file links to both updated/created roadmap documents:
- `docs/roadmap/roadmap.md`
- `docs/roadmap/<topic>-roadmap-<NNN>.md`
