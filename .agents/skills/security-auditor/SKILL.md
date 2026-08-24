---
name: security-auditor
description: >-
  Scans, audits, and enforces code safety, path traversal constraints, memory security, subprocess sanitization, and API compliance across gautama-graph.
  Use when conducting pre-release audits, vetting dependency vulnerabilities with govulncheck, verifying filesystem boundary confinement, checking for unsafe memory blocks, or certifying release readiness.
---

# Security & Compliance Auditor Skill (`gautama-gatekeeper`)

This skill defines the zero-trust security audit protocols, static analysis gates, memory safety checks, and compliance artifact schemas required to certify releases in the `gautama-graph` repository.

---

## Section 1: Core Behavioral Rules & Guardrails

When auditing code, dependency graphs, or filesystem interactions, you MUST enforce the following security invariants:

### 1.1 Path Traversal Defense & Zero-Trust Filesystem Confinement
* **Mandatory Path Containment Verification**: Inspect all call sites of `os.Open`, `os.ReadFile`, `os.WriteFile`, `os.Create`, and `filepath.Walk`. Every file path MUST be sanitized and verified against the workspace root:
  ```go
  cleanRoot := filepath.Clean(workspaceRoot)
  cleanTarget := filepath.Clean(targetPath)
  if !strings.HasPrefix(cleanTarget, cleanRoot+string(filepath.Separator)) && cleanTarget != cleanRoot {
      return fmt.Errorf("security violation: path %s escapes workspace root %s", cleanTarget, cleanRoot)
  }
  ```
* **No Unsanitized Symlink Resolution**: Never blindly follow arbitrary symlinks targeting directories outside the repository root.

### 1.2 Unsafe Memory & CGo Ban
* **Zero `unsafe` Package Usage**: The import and usage of `unsafe.Pointer`, `unsafe.Sizeof`, and unsafe slice conversions are strictly forbidden across the entire `gautama-graph` codebase.
* **No Unvetted CGo**: The codebase MUST remain pure Go 1.26+ (with external Python AST delegation limited strictly to the verified bridge script). Direct `import "C"` invocations are prohibited.

### 1.3 Subprocess Isolation & Command Injection Defense
* **No Shell Interpolation**: Subprocess execution in `python_bridge.go` MUST invoke binaries directly using `exec.CommandContext(ctx, "python3", scriptPath, targetFile)` with discrete string arguments. Never invoke commands via shell interpreters (`sh -c`, `bash -c`, `cmd.exe /c`).
* **Environment Variable Sanitization**: Ensure sub-processes do not inherit unconstrained or hazardous environment variables.
* **Bounded Stdout/Stderr Consumption**: Guard all subprocess readers with bounded size limits (e.g., `io.LimitReader(stdout, 10*1024*1024)`) to mitigate resource exhaustion or zip-bomb style JSON payloads.

### 1.4 Dependency Vulnerability Auditing & Public API Tracking
* **Automated `govulncheck` Scanning**: Run vulnerability assessments on all transitive dependencies in `go.mod` / `go.sum` before release certification.
* **Public API SemVer Enforcement**: Check `internal/auditor/types.go` for breaking signature changes. If exported types are modified, verify that git tagging and release notes reflect appropriate Semantic Versioning increments.

---

## Section 2: Expected Antigravity Artifact Deliverables

The Security & Compliance Auditor agent MUST produce structured verification artifacts stored in `<appDataDir>/brain/<conversation-id>/`:

### 2.1 Security Audit Report (`security_audit_report.md`)
Generated for pre-merge and pre-release audits to summarize security findings and compliance status.

```markdown
# Security & Compliance Audit Report: [Audit ID / Release Version]

## Audit Verdict
- **Certification Status**: APPROVED / REJECTED
- **Security Score**: 100/100
- **Path Traversal Vulnerabilities**: 0
- **Unsafe Code Blocks**: 0
- **Dependency Vulnerabilities (govulncheck)**: 0

## Filesystem Boundary Verification
| Component | File Path | Sanitization Method | Confinement Status |
| :--- | :--- | :--- | :--- |
| AST Parser | `internal/auditor/parser.go` | `filepath.Clean` + Root Prefix | VERIFIED |
| Doc Link Resolver | `internal/auditor/doc_auditor.go` | `filepath.Abs` + Prefix Check | VERIFIED |
| Atomic Store | `internal/auditor/store.go` | `.tmp` Staging + Clean Path | VERIFIED |

## Static Analysis & govulncheck Summary
- `govulncheck ./...`: Clean (0 known vulnerabilities detected).
- Unsafe memory scan: 0 occurrences of `import "unsafe"`.
- CGo scan: 0 occurrences of `import "C"`.
```

### 2.2 Security Compliance Checklist (`compliance_checklist.md`)
Interactive gate checklist verified before production merge:
- [x] All file access calls enforce strict root prefix containment.
- [x] No `unsafe.Pointer` or unvetted cgo usage present.
- [x] Subprocess commands use parameter arrays without shell expansion.
- [x] Atomic write protocols with temporary buffers `.tmp` and `os.Rename` validated.
- [x] Exported API symbols follow strict PascalCase formatting with godoc comments.
- [x] Context timeouts enforced across all AST traversals and IPC operations.

### 2.3 JSON Security Verification Meta-Artifact (`security_verification_meta.json`)
Machine-readable certificate signed by `gautama-gatekeeper` approving deployment.

```json
{
  "$schema": "https://antigravity.internal/schemas/v2/security-verification.json",
  "audit_id": "SEC-AUDIT-2026-08-24-001",
  "timestamp": "2026-08-24T08:15:00Z",
  "auditor_agent": "gautama-security-gatekeeper",
  "status": "APPROVED",
  "checks": {
    "path_traversal_confinement": "PASSED",
    "unsafe_memory_prohibition": "PASSED",
    "subprocess_sanitization": "PASSED",
    "dependency_vulnerabilities": "PASSED",
    "public_api_pascal_case": "PASSED",
    "atomic_persistence": "PASSED"
  },
  "govulncheck_vulnerabilities_count": 0,
  "release_ready": true
}
```
