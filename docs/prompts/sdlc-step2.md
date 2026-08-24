[.agents/personas/nexus.md](../../.agents/personas/nexus.md) [.agents/personas/feature-engineer.md](../../.agents/personas/feature-engineer.md) [.agents/personas/security-auditor.md](../../.agents/personas/security-auditor.md)

## CONTEXT & OBJECTIVE
Execute Phase 2 (Architecture Blueprint & Security Gate) of the Gautama Graph Software Development Lifecycle using the approved Phase 1 requirements specification (e.g., `docs/specs/<NNN>-<feature-name>-requirements.md`).

You are tasked with producing a formal Technical Blueprint & Threat Model document inside `docs/specs/` that defines clean Go interface abstractions, struct definitions, AST traversal algorithms, Python IPC JSON contracts (`python/ast_auditor_bridge.py`), atomic persistence protocols, and comprehensive security controls.

---

## 🔒 STEP 2 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)), `feature-engineer` ([@feature-engineer.md](../../.agents/personas/feature-engineer.md)), `security-auditor` ([@security-auditor.md](../../.agents/personas/security-auditor.md)).
- **Injected Skills**: `skills/feature-engineer` ([SKILL.md](../../.agents/skills/feature-engineer/SKILL.md)), `skills/security-auditor` ([SKILL.md](../../.agents/skills/security-auditor/SKILL.md)).
- **Excluded Context**: Active bug reports, test runner outputs, irrelevant subprocess scripts.

---

## 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH MAPPING (TOKEN OPTIMIZATION)
The Phase 2 blueprint team MUST query `graphify path "<A>" "<B>"`, `graphify explain "<concept>"`, or `graphify query "<package/interface>"` to construct clean Go interface contracts, AST walker hierarchies, Python IPC bridge signatures, and storage paths with minimal token consumption instead of reading raw files.

---

## 🛑 PHASE 2 EXECUTION CONSTRAINTS
1. **Forbidden Actions**: Writing Go implementation code, modifying `internal/auditor/*.go` or `cmd/**/*.go`, editing Python scripts, creating `*_test.go` files, running compilation binaries (`go build`), or updating project release status matrices is strictly forbidden in Phase 2.
2. **Artifact Output**: The blueprint must be saved directly as a dedicated markdown file inside `docs/specs/` using the 3-digit roadmap sequence code prefix (e.g., `docs/specs/<NNN>-<feature-name>-architecture-blueprint.md`).
3. **Phase Boundary Rule**: Receiving user feedback/approval on a Phase 2 blueprint document signifies architecture and security design sign-off ONLY. Agents MUST stop execution after Phase 2 and wait for explicit invocation of Step 3 / Phase 3 & 4 (`execute docs/prompts/sdlc-step3.md with docs/specs/<NNN>-<feature-name>-architecture-blueprint.md`) in a subsequent prompt.

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Go Engine & Interface Architecture Blueprint (`@feature-engineer.md`, `@nexus.md`)
- Draft granular Go interface contracts in `internal/auditor/types.go` adhering to the Interface Segregation Principle (e.g. `ASTParser`, `SelectorEvaluator`, `PythonASTBridge`, `GraphStore`, `DocGraphAuditorService`).
- Specify AST traversal architecture: `ast.Inspect` recursive walkers, receiver type bindings, and depth limiters (`Config.MaxASTDepth`).
- Blueprint Doc Graph topology engine: regex link extraction with code-block stripping, relative path normalization (`filepath.Join`), fragment stripping (`#anchor`), and `InDegree` / `OutDegree` tracking.
- Detail the two-phase atomic persistence protocol (`os.CreateTemp` / `.tmp` staging + `os.Rename`) and mutex concurrency guards.

### 2. Python IPC Bridge & Subprocess Contract (`@feature-engineer.md`, `@security-auditor.md`)
- Blueprint CLI subprocess contracts for `python/ast_auditor_bridge.py` called via Go `exec.CommandContext`.
- Specify JSON streaming schemas for inter-process communication: input file candidates, symbol queries, AST call matches, error payloads (`{"status": "error", "error": "..."}`).
- Define subprocess lifecycle parameters: 15s execution timeout (`context.WithTimeout`), discrete stdout/stderr stream separation, bounded memory limits (`io.LimitReader`), and mandatory `scanner.Err()` checks.

### 3. Cyber Security Architecture & Hardening Strategy (`@security-auditor.md`)
- **Zero-Trust Path Confinement**: Blueprint boundary verification using `filepath.Clean` and `filepath.Abs` prefix matching against `Config.WorkspaceRootPath`.
- **Pure Go & Unsafe Code Ban**: Enforce zero `unsafe` and zero CGo policies across all planned packages.
- **Command Injection Prevention**: Verify all subprocess arguments pass as discrete string arrays without shell evaluation.

---

## 📄 OUTPUT FILE REQUIREMENT
Save the completed technical blueprint to `docs/specs/<NNN>-<feature-name>-architecture-blueprint.md` (e.g., `docs/specs/001-ast-selector-depth-architecture-blueprint.md`) and provide the exact clickable file link in your response.
