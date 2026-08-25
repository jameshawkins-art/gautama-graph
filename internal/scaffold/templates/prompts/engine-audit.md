[.agents/personas/nexus.md](../../.agents/personas/nexus.md)

# Architecture, Code Quality & Security Audit Directive

## Context & Operational Mandate
You are executing a comprehensive, multi-lens Architectural, Subprocess Lifecycle, Documentation Topology, and Security Code Audit on this workspace.

Reference the active project personas:
- **Lead AI Workflow Architect ([@nexus.md](../../.agents/personas/nexus.md))**: System gatekeeper, Graphify knowledge graph enforcement, SDLC boundary compliance, and dynamic scaffolding governance.
- **Specialized Engineering Personas** (as defined in `.agents/personas/`): Feature engineers, test guards, security auditors, and system mechanics.

---

## 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)) and relevant audit personas.
- **Injected Skills**: Security auditing, code inspection, and architecture analysis skills.
- **Excluded Context**: Production deployment credentials, unimpacted external tools.

---

## 🕸️ Mandatory Graphify Knowledge Graph Scoping (Token Optimization)
Before performing raw file reads or broad greps across the repository, the audit team MUST query:
- `graphify query "<subsystem>"` to trace package dependencies and interface contracts.
- `graphify path "<Caller>" "<Callee>"` to map invocation flows.
- `graphify explain "<type or function>"` to analyze struct definitions and interface abstractions.
- Navigate [graphify-out/wiki/index.md](../../graphify-out/wiki/index.md) to inspect community clusters and high-level architecture with minimal token overhead.

---

## 1. Audit Scope & Target Codepaths

Audit the **Core Engine Packages**, **Documentation Topology**, **Subprocess/IPC Layers**, **State Stores**, and **CLI Entrypoints**:

### Target Codepaths & Subsystems:
1. **Core Processing & Business Logic**: Interface segregation, package boundaries, caller hierarchies.
2. **Documentation Topology & Relative Link Graph**: Link resolution, relative path normalization, orphan document identification.
3. **Subprocess & IPC Lifecycle**: Process lifecycle management, timeout enforcement, stream buffering.
4. **State Storage & Persistence**: Concurrency locks, two-phase atomic write staging (`.tmp` buffer + atomic rename).
5. **Command-Line Entrypoints**: Flag parsing, exit code discipline, error reporting.

---

## 2. Audit Execution & Deliverables

Conduct the audit across 5 technical pillars:
1. **Pillar 1: Interface Abstractions & Export Standards**: Idiomatic naming, documentation coverage, interface segregation.
2. **Pillar 2: Zero-Trust Path Confinement**: Path containment assertion on all file operations.
3. **Pillar 3: Atomic Persistence & Concurrency Hygiene**: Atomic staging buffers, clean lock management.
4. **Pillar 4: Subprocess Lifecycle & Stream Safety**: Process timeout deadlines, stream hygiene, error checking.
5. **Pillar 5: Test Suite Integrity & Vulnerability Assessment**: Full test suite execution, coverage checks, dependency vulnerability scanning.

---

## 📄 OUTPUT FILE REQUIREMENT
Save the finalized audit findings to `docs/audits/engine-audit-report-<date>.md` and provide the exact clickable file link in your response.
