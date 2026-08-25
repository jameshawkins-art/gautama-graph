[.agents/personas/nexus.md](../../.agents/personas/nexus.md)

## CONTEXT & OBJECTIVE
Execute Phase 1 (Requirements, Interface Specifications & Security Gate) of the Software Development Lifecycle for the target feature or roadmap task.

You are tasked with producing a formal Requirements Specification, Interface Contracts, Path Confinement Plans, Cyber Security Threat Model & Dependency Audit, and Acceptance Criteria document in `docs/specs/` for the feature (e.g., `docs/specs/<NNN>-<feature-name>-requirements.md` where `<NNN>` is the 3-digit sequence code of the target roadmap item in `docs/roadmap/roadmap.md`).

---

## 🔒 STEP 1 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)) and the responsible engineering persona.
- **Injected Skills**: Relevant domain skills (e.g. feature engineering, security auditing).
- **Excluded Context**: Bug remediation test runners, chaos fuzzing harnesses, runtime mechanic tools.

---

## 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH DISCOVERY (TOKEN OPTIMIZATION)
Before performing raw file reads or broad greps across the repository, the Phase 1 persona team MUST query `graphify query "<feature concept>"`, `graphify path "<A>" "<B>"`, `graphify explain "<concept>"`, or navigate `graphify-out/wiki/index.md` to extract existing package boundaries, AST structures, interfaces, and storage paths with minimal token consumption.

---

## 🛑 PHASE 1 EXECUTION CONSTRAINTS
1. **No Implementation Code**: Writing production code, creating test suites, or drafting Phase 2 architecture blueprints is strictly forbidden in Phase 1.
2. **No Release Matrix Updates**: Updating feature release status matrices in roadmap documents is strictly restricted to Step 4 / Phase 5 & 6 (`@nexus.md`).
3. **Artifact Output**: The specification must be saved directly as a dedicated markdown file inside `docs/specs/` using the 3-digit roadmap sequence code prefix (e.g., `docs/specs/<NNN>-<feature-name>-requirements.md`).
4. **Phase Boundary Rule**: User approval of a Phase 1 document signifies requirements, interface definitions, path boundary constraints, and security parameters sign-off ONLY. Agents MUST stop and wait for explicit user invocation of Step 2 (`execute docs/prompts/sdlc-step2.md with docs/specs/<NNN>-<feature-name>-requirements.md`) in a subsequent prompt.

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Requirements Scope & Interface Specifications
- Define feature functional scope and public API contracts.
- Specify exported structs, types, and interfaces in clean, idiomatic style with comprehensive documentation.
- Map symbol relationships and edge dependencies.

### 2. Filesystem Confinement & Atomic Persistence Plan
- Formulate zero-trust path containment specifications: all file resolutions MUST use clean relative paths and prevent workspace directory traversal.
- Define atomic write staging protocols: state mutations must stage to temporary buffers before committing via atomic rename.
- Guard all shared/stateful resources with proper synchronization primitives.

### 3. Cyber Security Threat Modeling & Dependency Safety
- **Path Traversal Defense**: Validate boundary check coverage on all file access call sites.
- **Unsafe Code & External Subprocess Safety**: Enforce memory safety and sanitize subprocess invocations.
- **Dependency Audit**: Audit project dependencies to ensure zero known vulnerabilities.

### 4. Definition of Done (DoD) & Acceptance Criteria
- Establish explicit, verifiable Acceptance Criteria for nominal and boundary execution paths.
- Define test coverage expectations and race/concurrency cleanliness.
- Include an **Edge Case & Failure Mode Matrix** detailing behavior on malformed syntax, timeout cancellation, and missing disk targets.

---

## 📄 OUTPUT FILE REQUIREMENT
Save the finalized specification to `docs/specs/<NNN>-<feature-name>-requirements.md` (e.g. `docs/specs/001-core-feature-requirements.md`) and provide the exact clickable file link in your response.
