[.agents/personas/nexus.md](../../.agents/personas/nexus.md)

## CONTEXT & OBJECTIVE
Execute Phase B2 (Root Cause Analysis & Fix Blueprint Gate) of the Bug Remediation Lifecycle for the triaged defect specification in `docs/bugs/` (e.g. `docs/bugs/bug-<description>-<id>.md`).

You are tasked with identifying the exact failure mechanism, designing a surgical code patch that eliminates the defect without collateral changes or API breakage, analyzing side-effects, and updating the bug specification document with the RCA and Technical Remediation Blueprint.

---

## 🔒 STEP B2 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)), debugger/remediation, and security auditor personas.
- **Injected Skills**: Debugging, remediation, and security auditing skills.
- **Excluded Context**: Feature roadmap specs, unimpacted package implementations, full test suite harnesses.

---

## 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH MAPPING (TOKEN OPTIMIZATION)
The RCA team MUST query `graphify path` and `graphify explain` to analyze execution paths, error handling chains, and mutex boundaries with minimal token consumption instead of reading raw files.

---

## 🛑 PHASE B2 EXECUTION CONSTRAINTS
1. **No Production Code Edits**: Modifying production code or running full release builds is strictly forbidden in Phase B2.
2. **Zero Collateral Refactoring**: The proposed patch must be minimal and strictly confined to the faulting execution path.
3. **In-Place Update**: Document RCA and technical remediation steps directly inside the target bug report (`docs/bugs/bug-<description>-<id>.md`).
4. **Phase Boundary Rule**: Receiving user feedback/approval on a Phase B2 blueprint signifies RCA and fix plan sign-off ONLY. Agents MUST stop and wait for explicit invocation of Bug Step 3 (`execute docs/prompts/bug-step3.md with docs/bugs/bug-<description>-<id>.md`).

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Root Cause Analysis (RCA)
- Trace code execution paths using `graphify path` to isolate the exact technical failure mechanism (e.g., pipe deadlock, mutex deadlock, missing error handling, unhandled JSON parsing exception).
- Document interactions causing memory leaks, hanging threads/goroutines, or data corruption.

### 2. Technical Remediation Blueprint & Surgical Diff Preview
- Author a minimal, surgical code patch preview:
  - Preserves public API signatures and contracts.
  - Preserves error wrapping chains.
  - Ensures clean resource cleanup and unlock deferrals.
  - Ensures proper context cancellation handling.

### 3. Security & Side-Effect Assessment
- Verify that the proposed fix does not compromise zero-trust path containment, introduce unsafe memory blocks, or create command injection vectors.
- Assess potential downstream regressions.

---

## 📄 OUTPUT FILE REQUIREMENT
Update the target bug document (`docs/bugs/bug-<description>-<id>.md`) with the Root Cause Analysis (RCA) and Technical Remediation Blueprint sections, and provide the exact clickable file link in your response.
