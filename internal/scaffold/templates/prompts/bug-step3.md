[.agents/personas/nexus.md](../../.agents/personas/nexus.md)

## CONTEXT & OBJECTIVE
Execute Phase B3 (Test-Driven Bug Remediation & Regression Verification) and Phase B4 (SQA Certification & Defect Closure) of the Bug Remediation Lifecycle for the approved bug blueprint in `docs/bugs/` (e.g. `docs/bugs/bug-<description>-<id>.md`).

You are tasked with applying the minimal surgical code fix, proving that the failing reproduction test turns green, running the full test suite with concurrency/race detection, executing knowledge graph sync (`make graphify-update` or `./scripts/graphify_sync.sh`), updating the bug specification status to `CLOSED (🟢 RESOLVED)`, and creating or updating `walkthrough.md`.

---

## 🔒 STEP B3 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)), debugger/remediation, test guard, and security auditor personas.
- **Injected Skills**: Debugging, remediation, and regression testing skills.
- **Excluded Context**: New feature specifications, pre-SDLC roadmaps, unimpacted subsystem implementations.

---

## 🕸️ MANDATORY GRAPHIFY DISCOVERY & POST-FIX GRAPH SYNC
1. **Graphify Discovery (Token Optimization)**: Query `graphify path` or `graphify explain` to review call chains before editing code.
2. **Post-Fix Graph Sync**: After applying the fix, the agent team MUST run `make graphify-update` or `./scripts/graphify_sync.sh` (or `graphify update .`) to update graph nodes and verify graph health.

---

## 🛑 PHASE B3 & B4 EXECUTION CONSTRAINTS
1. **Test-Driven Discipline**: The minimal reproduction test written in Phase B1 MUST fail before the fix and pass cleanly after applying the patch.
2. **Full Test Suite Verification**: The full test suite MUST pass 100% with zero race conditions or unhandled errors.
3. **Artifact Output**: Update the target bug specification document (`docs/bugs/bug-<description>-<id>.md`) status to `CLOSED (🟢 RESOLVED)`, emit `remediation_meta.json`, and update `walkthrough.md`.
4. **Phase Boundary Rule**: Completing Phase B3 & B4 certifies bug resolution, regression immunity, and defect closure.

---

## 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### 1. Targeted Code Remediation
- Apply the minimal, surgical patch strictly confined to the faulting execution path.
- Ensure error wrapping with context, immediate resource cleanup, and stream hygiene.
- Emit machine-readable meta-artifact `remediation_meta.json`.

### 2. Regression & SQA Verification
- Execute the isolated reproduction test to confirm green pass.
- Run full regression suite with concurrency/race detection.
- Validate that statement coverage on modified packages meets project thresholds.

### 3. Security Check & Defect Closure
- Verify zero-trust path containment, memory safety, and sanitized subprocess calls.
- Run knowledge graph sync to update `graphify-out/graph.json`.
- Update bug specification status in `docs/bugs/bug-<description>-<id>.md` to `CLOSED (🟢 RESOLVED)`.
- Create or update `walkthrough.md` summarizing the defect, root cause, surgical patch diff, and verification test outputs.

---

## 📄 OUTPUT FILE REQUIREMENT
Provide clickable file links to the updated bug specification document (`docs/bugs/bug-<description>-<id>.md`), `remediation_meta.json`, and `walkthrough.md`.
