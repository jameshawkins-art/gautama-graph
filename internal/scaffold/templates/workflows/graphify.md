---
name: graphify
description: Turn any folder of files into a navigable knowledge graph with automated AST and doc verification
---

# Workflow: graphify

This project utilizes `github.com/jameshawkins-art/gautama-graph` v1.5.0 for turnkey knowledge graph orchestration, deep AST cross-package verification, and markdown doc link auditing.

## Turnkey Execution

1. **Full Turnkey Pipeline & Strict Quality Gate**:
   ```bash
   make audit
   ```
   Executes binary resolution/caching, base graph extraction (`graphify update .`), deep Go AST code relationship audit, and strict markdown doc graph audit.

2. **Update Knowledge Graph Artifacts**:
   ```bash
   make graphify-update
   ```
   Refreshes `graphify-out/graph.json`, `graphify-out/GRAPH_REPORT.md`, `graphify-out/graph.html`, and `graphify-out/doc_graph_audit.json`.

3. **Automated Markdown Doc Link Remediation**:
   ```bash
   make audit-remediate
   ```
   Automatically repairs broken relative markdown links across `.agents/` and `docs/` using canonical path calculation and fuzzy basename matching.

4. **Standalone Audits**:
   - `make audit-docs`: Sub-second doc graph validation (`cmd/graphify-doc-audit -strict`).
   - `make audit-ast`: Go AST code relationship validation (`cmd/graphify-ast-audit`).
