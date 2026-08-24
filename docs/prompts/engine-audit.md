[.agents/personas/nexus.md](../../.agents/personas/nexus.md) [.agents/personas/security-auditor.md](../../.agents/personas/security-auditor.md) [.agents/personas/feature-engineer.md](../../.agents/personas/feature-engineer.md) [.agents/personas/debugger-remediation.md](../../.agents/personas/debugger-remediation.md) [.agents/personas/regression-tester.md](../../.agents/personas/regression-tester.md)

# Lead Engine Architect & Security Code Audit Directive

## Context & Operational Mandate
You are executing a comprehensive, multi-lens Architectural, Subprocess Lifecycle, Doc Topology, and Security Code Audit on the core engine and CLI subsystems of **Gautama Graph** (`github.com/jameshawkins-art/gautama-graph`).

Reference the active project personas:
- **Lead AI Workflow Architect ([@nexus.md](../../.agents/personas/nexus.md))**: System gatekeeper, Graphify knowledge graph enforcement, SDLC boundary compliance, and dynamic scaffolding governance.
- **Feature Engineer ([@feature-engineer.md](../../.agents/personas/feature-engineer.md))**: Go standard library AST parsers (`go/ast`, `go/parser`), `ast.Inspect` selector expression evaluation, interface segregation, and atomic storage persistence (`internal/auditor/store.go`).
- **Debugger & Remediation ([@debugger-remediation.md](../../.agents/personas/debugger-remediation.md))**: IPC bridge deadlocks, `bufio.Scanner` / `scanner.Err()` checks, mutex contention, circular doc link recursion loops, and minimal patch triage.
- **Security & Compliance Auditor ([@security-auditor.md](../../.agents/personas/security-auditor.md))**: Zero-trust workspace path containment (`filepath.Clean`, `filepath.Abs`), unsafe memory prohibition (zero `unsafe`/cgo), subprocess argument injection defense, and `govulncheck`.
- **Regression & Test Guard ([@regression-tester.md](../../.agents/personas/regression-tester.md))**: Table-driven test coverage ($\ge 85\%$), race detector enforcement (`-race`), parser fuzzing, and hermetic filesystem sandboxing (`t.TempDir()`).

---

## 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
To prevent **context rot** and token bloat:
- **Active Personas**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md)), `security-auditor` ([@security-auditor.md](../../.agents/personas/security-auditor.md)), `feature-engineer` ([@feature-engineer.md](../../.agents/personas/feature-engineer.md)).
- **Injected Skills**: `skills/security-auditor` ([SKILL.md](../../.agents/skills/security-auditor/SKILL.md)), `skills/feature-engineer` ([SKILL.md](../../.agents/skills/feature-engineer/SKILL.md)).
- **Excluded Context**: Production deployment credentials, unimpacted external tools.

---

## 🕸️ Mandatory Graphify Knowledge Graph Scoping (Token Optimization)
Before performing raw file reads or broad greps across the repository, the audit team MUST query:
- `graphify query "<subsystem>"` to trace package dependencies and interface contracts.
- `graphify path "<Caller>" "<Callee>"` to map IPC invocation flows between Go engine and Python bridge scripts.
- `graphify explain "<type or function>"` to analyze struct definitions and interface abstractions.
- Navigate [graphify-out/wiki/index.md](../../graphify-out/wiki/index.md) to inspect community clusters and high-level architecture with minimal token overhead.

---

## 1. Audit Scope & Target Codepaths

Audit the **Core AST Engine**, **Doc-Graph Auditor**, **Python IPC Bridge**, **Atomic GraphStore**, and **CLI Tools**:

### Target Codepaths & Subsystems:
1. **AST Parsing & Selector Evaluation**:
   - [internal/auditor/engine.go](../../internal/auditor/engine.go): Pipeline orchestration, candidate edge processing, extension routing (`.go` vs `.py`).
   - [internal/auditor/parser.go](../../internal/auditor/parser.go): Go file AST parsing, workspace path boundary checks.
   - [internal/auditor/evaluator.go](../../internal/auditor/evaluator.go): `ast.Inspect` selector and call expression matching.
2. **Doc-Graph Topology & Orphan Detection**:
   - [internal/auditor/doc_auditor.go](../../internal/auditor/doc_auditor.go): Markdown link parsing, code-block stripping, relative link resolution, `InDegree` calculation, and orphan identification.
3. **Python AST IPC Bridge**:
   - [internal/auditor/python_bridge.go](../../internal/auditor/python_bridge.go): Subprocess execution wrapper, stream buffering, timeout enforcement.
   - [python/ast_auditor_bridge.py](../../python/ast_auditor_bridge.py): Python AST visitor, JSON stdin/stdout serialization, error handling.
4. **Atomic Store Persistence**:
   - [internal/auditor/store.go](../../internal/auditor/store.go): Mutex locking, two-phase staging (`.tmp`), atomic `os.Rename`.
5. **Command-Line Entrypoints**:
   - [cmd/graphify-ast-audit/main.go](../../cmd/graphify-ast-audit/main.go): AST audit CLI, flag parsing (`--strict`, `--verbose`), exit code discipline.
   - [cmd/graphify-doc-audit/main.go](../../cmd/graphify-doc-audit/main.go): Doc audit CLI, exit code gating.
   - [scripts/graphify_sync.sh](../../scripts/graphify_sync.sh): Master synchronization orchestration.

---

## 2. Audit Execution & Deliverables

Conduct the audit across 5 technical pillars:
1. **Pillar 1: Go Interface Abstractions & Export Standards**: PascalCase naming, Godoc coverage, interface segregation in `types.go`.
2. **Pillar 2: Zero-Trust Path Confinement**: Path containment assertion (`filepath.Clean`, `filepath.Abs`, root prefix validation) on all file operations.
3. **Pillar 3: Atomic Persistence & Concurrency Hygiene**: `.tmp` staging buffer commitment via `os.Rename`, deferral of `mu.Unlock()`.
4. **Pillar 4: Subprocess Lifecycle & Stream Safety**: `exec.CommandContext` deadlines, discrete stdout JSON / stderr logs, mandatory `scanner.Err()` checks.
5. **Pillar 5: Test Suite Integrity & Vulnerability Assessment**: `GOWORK=off go test -v -race ./...`, statement coverage $\ge 85\%$, `govulncheck ./...`.

---

## 📄 OUTPUT FILE REQUIREMENT
Save the finalized audit findings to `docs/audits/engine-audit-report-<date>.md` and provide the exact clickable file link in your response.
