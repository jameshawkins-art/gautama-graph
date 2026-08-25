[.agents/personas/nexus.md](../../.agents/personas/nexus.md)

# Pre-SDLC Roadmap Item Formulation Directive

## Context & Operational Mandate
You are executing the **Pre-SDLC Roadmap Item Decomposition** for a specific capability, module extension, or architectural enhancement in this workspace.

Reference the active project personas:
- **Lead AI Workflow Architect & Governor ([@nexus.md](../../.agents/personas/nexus.md))**: Master orchestrator, Prompt Ops director, Dynamic Scaffolding governor, and release status gatekeeper.
- **Specialized Engineering Personas** (as defined in `.agents/personas/`): Feature engineers, test guards, security auditors, and system mechanics.

---

## 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)) and relevant engineering personas.
- **Excluded Context**: Production code edits, unit test executions, release builds.

---

## 🕸️ Mandatory Graphify Knowledge Graph Discovery (Token Optimization)
Before performing raw file reads, the persona team MUST query `graphify query "<topic>"`, `graphify path "<A>" "<B>"`, or `graphify explain "<concept>"` to map target packages, existing interfaces, and storage paths with minimal token consumption.

---

## 🛑 STRICT PRE-SDLC EXECUTION CONSTRAINTS
1. **No SDLC Kick-Off**: This prompt is strictly PRE-SDLC. Do **NOT** generate requirements specification files in `docs/specs/` (e.g. `docs/specs/<NNN>-<feature>-requirements.md`).
2. **No Implementation Code or Blueprints**: Writing source code, scripts, test files, or Phase 2 technical blueprints is strictly forbidden during roadmap item generation.
3. **Target Files ONLY**:
   - Master Roadmap: Update [docs/roadmap/roadmap.md](../roadmap/roadmap.md) (Master Feature Table + Detailed Specifications section).
   - Dedicated Item Spec: Create [docs/roadmap/<topic>-roadmap-<NNN>.md](../roadmap/) (e.g. `docs/roadmap/core-module-roadmap-001.md`).

---

## 📋 REQUIRED DELIVERABLES & OUTPUT FORMAT

1. **Master Roadmap Update ([docs/roadmap/roadmap.md](../roadmap/roadmap.md))**:
   - Query `docs/roadmap/` for the highest 3-digit zero-padded sequence code `NNN` and increment by 1 (`NNN+1`).
   - Append new row to the Master Feature Roadmap Table with status `(🔴 NOT STARTED)`.
   - Append detailed item summary under Detailed Item Specifications section.

2. **Dedicated Roadmap Item Document ([docs/roadmap/<topic>-roadmap-<NNN>.md](../roadmap/))**:
   - **Header & Metadata**: Document Title, Sequence Code `NNN`, Persona Drivers & Gatekeepers, Status `(🔴 NOT STARTED)`.
   - **1. Executive Summary & Strategic Objective**: Problem statement, target architecture, knowledge graph value proposition.
   - **2. Subsystem / Component Matrix**: Table mapping packages, modules, interface types, scripts, and storage paths using Graphify queries.
   - **3. Phased Master Task Matrix**: Task decomposition (`NNN.1`, `NNN.2`...) with driver personas, priority, estimated effort, and target SDLC phase.
   - **4. Definition of Done (DoD)**: Explicit, verifiable acceptance criteria for future SDLC execution.

---

## 📄 OUTPUT FILE REQUIREMENT
Provide clickable file links to both updated/created roadmap documents:
- `docs/roadmap/roadmap.md`
- `docs/roadmap/<topic>-roadmap-<NNN>.md`
