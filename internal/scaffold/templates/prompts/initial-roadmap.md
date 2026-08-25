[.agents/personas/nexus.md](../../.agents/personas/nexus.md)

# Master Engineering Roadmap Formulation Directive (Pre-SDLC)

## Context & Operational Mandate
You are executing the **Pre-SDLC Initial Roadmap Formulation** for this workspace.

Reference the active project personas:
- **Lead AI Workflow Architect & Governor ([@nexus.md](../../.agents/personas/nexus.md))**: Master orchestrator, Prompt Ops director, Dynamic Scaffolding governor, and release status gatekeeper.
- **Specialized Engineering Personas** (as defined in `.agents/personas/`): Feature engineers, test guards, security auditors, and system mechanics.

---

## 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)) and core engineering personas.
- **Excluded Context**: Production code implementations, runtime bug remediation scripts.

---

## 🕸️ Mandatory Graphify Knowledge Graph Discovery (Token Optimization)
Before performing raw file reads, the roadmap team MUST query:
- `graphify query "<subsystem>"` to inspect existing packages, interfaces, and modules.
- `graphify path "<A>" "<B>"` to trace dependencies across components.
- Navigate `graphify-out/wiki/index.md` or query graph concepts to map subsystem architecture with minimal token overhead.

---

## 🛑 STRICT PRE-SDLC EXECUTION CONSTRAINTS
1. **No SDLC Kick-Off**: This prompt is strictly PRE-SDLC. Do **NOT** generate requirements specification files in `docs/specs/` (e.g. `docs/specs/<NNN>-<feature>-requirements.md`).
2. **No Implementation Code or Blueprints**: Writing source code, test files, or Phase 2 technical blueprints is strictly forbidden during roadmap formulation.
3. **Target File ONLY**: The output must be written directly to `docs/roadmap/roadmap.md`.

---

## 📋 REQUIRED DELIVERABLES & OUTPUT FORMAT

Generate a comprehensive, structured Master Product Roadmap document saved to `docs/roadmap/roadmap.md` containing:
1. **Executive Summary**: Strategic mission of the project, core value proposition, and architectural goals.
2. **Architecture & Subsystem Alignment**: Overview of key packages, services, CLI commands, and libraries.
3. **Master Feature Roadmap Table**:
   - Sequential 3-digit zero-padded sequence codes (`001`, `002`...).
   - Feature Title, Lead Personas, Target Milestone, Status `(🔴 NOT STARTED)`.
4. **Phased Milestone Breakdown**: Logical milestone progression across foundational architecture, feature capabilities, test/audit automation, and releases.
5. **Quality & Security Guardrails**: Coverage targets, zero-trust path containment, and automated validation gates.

---

## 📄 OUTPUT FILE REQUIREMENT
Save the finalized master roadmap to `docs/roadmap/roadmap.md` and provide the exact clickable file link in your response.
