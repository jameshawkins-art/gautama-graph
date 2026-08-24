[.agents/personas/nexus.md](../../.agents/personas/nexus.md) [.agents/personas/feature-engineer.md](../../.agents/personas/feature-engineer.md) [.agents/personas/security-auditor.md](../../.agents/personas/security-auditor.md)

## CONTEXT & OBJECTIVE
Execute Phase 1 (Requirements, Go Interface Specifications & Security Gate) of the Gautama Graph Software Development Lifecycle for the target feature or roadmap task.

You are tasked with producing a formal Requirements Specification, Go Interface Contracts, Path Confinement Plans, Cyber Security Threat Model & Dependency Audit, and Acceptance Criteria document in `docs/specs/` for the feature (e.g., `docs/specs/<NNN>-<feature-name>-requirements.md` where `<NNN>` is the 3-digit sequence code of the target roadmap item in `docs/roadmap/roadmap.md`).

---

## 🔒 STEP 1 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)), `feature-engineer` ([@feature-engineer.md](../../.agents/personas/feature-engineer.md)), `security-auditor` ([@security-auditor.md](../../.agents/personas/security-auditor.md)).
- **Injected Skills**: `skills/feature-engineer` ([SKILL.md](../../.agents/skills/feature-engineer/SKILL.md)), `skills/security-auditor` ([SKILL.md](../../.agents/skills/security-auditor/SKILL.md)).
- **Excluded Context**: Bug remediation test runners, chaos fuzzing harnesses, runtime mechanic tools.

---

## 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH DISCOVERY (TOKEN OPTIMIZATION)
Before performing raw file reads or broad greps across the repository, the Phase 1 persona team MUST query `graphify query "<feature concept>"`, `graphify path "<A>" "<B>"`, `graphify explain "<concept>"`, or navigate `graphify-out/wiki/index.md` to extract existing package boundaries, AST structures, interfaces in `internal/auditor/types.go`, and storage paths with minimal token consumption.

---

## 🛑 PHASE 1 EXECUTION CONSTRAINTS
1. **No Implementation Code**: Writing Go production code, editing AST parsers, modifying Python scripts (`python/*.py`), creating `*_test.go` files, or drafting Phase 2 architecture blueprints is strictly forbidden in Phase 1.
2. **No Release Matrix Updates**: Updating feature release status matrices in roadmap documents is strictly restricted to Step 4 / Phase 5 & 6 (`@nexus.md`).
3. **Artifact Output**: The specification must be saved directly as a dedicated markdown file inside `docs/specs/` using the 3-digit roadmap sequence code prefix (e.g., `docs/specs/<NNN>-<feature-name>-requirements.md`).
4. **Phase Boundary Rule**: User approval of a Phase 1 document signifies requirements, Go interface definitions, path boundary constraints, and security parameters sign-off ONLY. Agents MUST stop and wait for explicit user invocation of Step 2 / Phase 2 (`execute docs/prompts/sdlc-step2.md with docs/specs/<NNN>-<feature-name>-requirements.md`) in a subsequent prompt.

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Requirements Scope & Go Interface Specifications (`@nexus.md`, `@feature-engineer.md`)
- Define feature functional scope and public API contracts for `internal/auditor/types.go`.
- Specify exported Go structs and interfaces in strict PascalCase with comprehensive godoc comments.
- Map symbol relationships, AST traversal depths (`Config.MaxASTDepth`), and edge confidence models (`EXTRACTED_AST`, `INFERRED_HEURISTIC`, `PRUNED_PHANTOM`).
- Map Markdown doc-graph link parsing, anchor stripping, and orphan metrics requirements (`InDegree == 0`).

### 2. Filesystem Confinement & Two-Phase Persistence Plan (`@feature-engineer.md`, `@security-auditor.md`)
- Formulate zero-trust path containment specifications: all file resolutions MUST use `filepath.Clean` and assert `strings.HasPrefix(target, cleanRoot)`.
- Define atomic write staging protocols: all mutations to `graphify-out/graph.json` and `graphify-out/doc_graph_audit.json` must stage to `<target>.tmp` buffers before committing via `os.Rename`.
- Guard all stateful graph stores with dedicated `sync.Mutex` locks.

### 3. Cyber Security Threat Modeling & Subprocess Safety (`@security-auditor.md`)
- **Path Traversal Defense**: Validate 100% boundary check coverage on all file access call sites.
- **Unsafe Code & CGo Prohibition**: Enforce strict zero-tolerance for `import "unsafe"` and unvetted CGo.
- **Subprocess Argument Sanitization**: Specify discrete argument slice passing for `exec.CommandContext` without shell expansion (`sh -c`, `bash -c`).
- **Dependency Audit**: Run `govulncheck` assessment on `go.mod` dependencies to ensure zero known vulnerabilities.

### 4. Definition of Done (DoD) & Acceptance Criteria (`@nexus.md`, `@feature-engineer.md`)
- Establish explicit, verifiable Acceptance Criteria for nominal and boundary execution paths.
- Define test coverage expectations ($\ge 85\%$ statement coverage on modified packages, race detector cleanliness with `-race`).
- Include an **Edge Case & Failure Mode Matrix** detailing behavior on malformed syntax, cyclic doc links, subprocess timeout cancellation, and missing disk targets.

---

## 📄 OUTPUT FILE REQUIREMENT
Save the finalized specification to `docs/specs/<NNN>-<feature-name>-requirements.md` (e.g. `docs/specs/001-ast-selector-depth-requirements.md`) and provide the exact clickable file link in your response.
