# Global Workspace Execution & Gatekeeping Rules (rules.md)

## 1. System Architecture & Technical Stack Boundaries

### Language & Runtime Boundaries
- **Core Engine & CLI Layer**: All active AST parsing, selector evaluation, graph persistence, doc auditing, and CLI tools MUST be implemented strictly in **Go 1.26+**.
- **Python AST Subprocess Bridge**: Python usage is strictly confined to `python/ast_auditor_bridge.py` for Python source AST extraction, invoked exclusively via Go's `exec.CommandContext` with bounded context timeouts.
- **Pure Go Ecosystem**: Zero `unsafe` package usage and zero unvetted CGo bindings across the entire codebase.

---

## 2. Module Gatekeeper & Lifecycle Protocol

- **Verification Gate Requirement**: No pipeline phase, subagent handoff, or release candidate may proceed without an explicit, validated meta-artifact (`feature_delivery.json`, `remediation_meta.json`, `test_verification_meta.json`, `security_verification_meta.json`).
- **Failure Escalation**: Any failure or test regression automatically routes control to the Debugger & Remediation Agent (`@debugger-remediation.md`) for minimal reproduction and root cause isolation.

---

## 3. Filesystem Safety & Atomic Persistence

- **Zero-Trust Path Confinement**: Every file access call (`os.Open`, `os.ReadFile`, `os.WriteFile`, `os.Create`, `filepath.Walk`) MUST validate path containment using `filepath.Clean` and assert that the target path begins with the workspace root.
- **Two-Phase Atomic File Commit**: Under no circumstances should `graphify-out/graph.json` or `graphify-out/doc_graph_audit.json` be overwritten in-place with partial streams. Writes must stage to `<target>.tmp` and commit via `os.Rename`.

---

## 4. Concurrency, Stream Hygiene & Code Quality

- **`bufio.Scanner` Stream Iteration Check**: Whenever scanning stream inputs or subprocess stdout/stderr pipes, developers and agents MUST check `scanner.Err()` immediately after the `for scanner.Scan()` loop terminates (e.g., `if err := scanner.Err(); err != nil`).
- **Mutex Safety**: Every `mu.Lock()` call on stateful stores MUST be immediately followed by `defer mu.Unlock()`.
- **Zero Suppression**: Never swallow errors or use silent fallback defaults without returning an explicit, wrapped error (`fmt.Errorf("...: %w", err)`).
- **Public API Export**: All exported Go structs, interfaces, methods, and constants must follow PascalCase naming with comprehensive godoc comments.
- **AST Knowledge Graph Sync**: Following any modification to codebase files, developers and agents MUST execute `graphify update .` and `go run cmd/graphify-ast-audit/main.go` to maintain graph accuracy.
