[.agents/personas/nexus.md](../../.agents/personas/nexus.md) [.agents/personas/debugger-remediation.md](../../.agents/personas/debugger-remediation.md) [.agents/personas/security-auditor.md](../../.agents/personas/security-auditor.md)

## CONTEXT & OBJECTIVE
Execute Phase B2 (Root Cause Analysis & Fix Blueprint Gate) of the Gautama Graph Bug Remediation Lifecycle for the triaged defect specification in `docs/bugs/` (e.g. `docs/bugs/bug-<description>-<id>.md`).

You are tasked with identifying the exact failure mechanism, designing a surgical code patch that eliminates the defect without collateral changes or API breakage, analyzing side-effects, and updating the bug specification document with the RCA and Technical Remediation Blueprint.

---

## 🔒 STEP B2 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)), `debugger-remediation` ([@debugger-remediation.md](../../.agents/personas/debugger-remediation.md)), `security-auditor` ([@security-auditor.md](../../.agents/personas/security-auditor.md)).
- **Injected Skills**: `skills/debugger-remediation` ([SKILL.md](../../.agents/skills/debugger-remediation/SKILL.md)), `skills/security-auditor` ([SKILL.md](../../.agents/skills/security-auditor/SKILL.md)).
- **Excluded Context**: Feature roadmap specs, unimpacted package implementations, full test suite harnesses.

---

## 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH MAPPING (TOKEN OPTIMIZATION)
The RCA team MUST query `graphify path` and `graphify explain` to analyze execution paths, error handling chains, and mutex boundaries with minimal token consumption instead of reading raw files.

---

## 🛑 PHASE B2 EXECUTION CONSTRAINTS
1. **No Production Code Edits**: Modifying production Go code (`internal/auditor/*.go`), Python bridge scripts, or running full release builds is strictly forbidden in Phase B2.
2. **Zero Collateral Refactoring**: The proposed patch must be minimal and strictly confined to the faulting execution path.
3. **In-Place Update**: Document RCA and technical remediation steps directly inside the target bug report (`docs/bugs/bug-<description>-<id>.md`).
4. **Phase Boundary Rule**: Receiving user feedback/approval on a Phase B2 blueprint signifies RCA and fix plan sign-off ONLY. Agents MUST stop and wait for explicit invocation of Bug Step 3 (`execute docs/prompts/bug-step3.md with docs/bugs/bug-<description>-<id>.md`).

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Root Cause Analysis (RCA) (`@debugger-remediation.md`, `@nexus.md`)
- Trace code execution paths using `graphify path` to isolate the exact technical failure mechanism (e.g., unbuffered subprocess pipe deadlock, un-deferred mutex unlock, unhandled circular doc link, unhandled JSON syntax exception in Python bridge).
- Document interactions causing memory leaks, hanging goroutines, or data corruption in `graphify-out/`.

### 2. Technical Remediation Blueprint & Surgical Diff Preview (`@debugger-remediation.md`)
- Author a minimal, surgical code patch preview:
  - Preserves public API signatures in `internal/auditor/types.go`.
  - Preserves error wrapping chains (`fmt.Errorf("...: %w", err)`).
  - Ensures proper `scanner.Err()` checks and `defer mu.Unlock()`.
  - Ensures context cancellation kills Python subprocess groups cleanly.

### 3. Security & Side-Effect Assessment (`@security-auditor.md`)
- Verify that the proposed fix does not compromise zero-trust path containment, introduce `unsafe` memory blocks, or create command injection vectors in `exec.CommandContext`.
- Assess potential downstream regressions.

---

## 📄 OUTPUT FILE REQUIREMENT
Update the target bug document (`docs/bugs/bug-<description>-<id>.md`) with the Root Cause Analysis (RCA) and Technical Remediation Blueprint sections, and provide the exact clickable file link in your response.
