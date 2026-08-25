[.agents/personas/nexus.md](../../.agents/personas/nexus.md)

## CONTEXT & OBJECTIVE
Execute Phase 2 (Architecture Blueprint & Technical Design Gate) of the Software Development Lifecycle using the approved Phase 1 requirements specification (e.g., `docs/specs/<NNN>-<feature-name>-requirements.md`).

You are tasked with producing a formal Technical Blueprint & Threat Model document inside `docs/specs/` that defines clean interface abstractions, struct/class definitions, algorithmic flows, IPC/API contracts, atomic persistence protocols, and comprehensive security controls.

---

## 🔒 STEP 2 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)) and the responsible engineering persona.
- **Injected Skills**: Relevant domain skills (e.g. feature engineering, architecture design, security auditing).
- **Excluded Context**: Active bug reports, test runner outputs, irrelevant subprocess scripts.

---

## 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH MAPPING (TOKEN OPTIMIZATION)
The Phase 2 blueprint team MUST query `graphify path "<A>" "<B>"`, `graphify explain "<concept>"`, or `graphify query "<package/interface>"` to construct clean interface contracts, caller hierarchies, and storage paths with minimal token consumption instead of reading raw files.

---

## 🛑 PHASE 2 EXECUTION CONSTRAINTS
1. **Forbidden Actions**: Writing production implementation code, editing existing sources, creating test suites, running build binaries, or updating project release status matrices is strictly forbidden in Phase 2.
2. **Artifact Output**: The blueprint must be saved directly as a dedicated markdown file inside `docs/specs/` using the 3-digit roadmap sequence code prefix (e.g., `docs/specs/<NNN>-<feature-name>-architecture-blueprint.md`).
3. **Phase Boundary Rule**: Receiving user feedback/approval on a Phase 2 blueprint document signifies architecture and security design sign-off ONLY. Agents MUST stop execution after Phase 2 and wait for explicit invocation of Step 3 (`execute docs/prompts/sdlc-step3.md with docs/specs/<NNN>-<feature-name>-architecture-blueprint.md`) in a subsequent prompt.

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Engine & Interface Architecture Blueprint
- Draft granular interface contracts adhering to the Interface Segregation Principle.
- Specify traversal architecture, receiver bindings, and algorithmic limits.
- Detail the two-phase atomic persistence protocol (`.tmp` staging + atomic rename) and concurrency guards.

### 2. Subprocess & API Contracts
- Blueprint subprocess and IPC contracts with bounded execution timeouts, stream separation, and robust error payloads.
- Define communication lifecycle parameters and memory bounds.

### 3. Cyber Security Architecture & Hardening Strategy
- **Zero-Trust Path Confinement**: Blueprint boundary verification using path cleaning and root prefix containment assertions.
- **Memory & Concurrency Safety**: Enforce memory safety and clean thread/goroutine lifecycle invariants.
- **Command Injection Prevention**: Verify all external commands pass discrete argument parameters without shell evaluation.

---

## 📄 OUTPUT FILE REQUIREMENT
Save the completed technical blueprint to `docs/specs/<NNN>-<feature-name>-architecture-blueprint.md` (e.g., `docs/specs/001-core-feature-architecture-blueprint.md`) and provide the exact clickable file link in your response.
