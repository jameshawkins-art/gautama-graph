---
trigger: always_on
description: Consult and synchronize the Graphify knowledge graph at graphify-out/ for codebase navigation, AST pruning, and doc link integrity.
---

# Graphify Knowledge Graph & Topology Rules

This project maintains a deterministic knowledge graph and documentation topology at `graphify-out/`.

## 1. Graph Discovery First (Token Optimization)
- For codebase, relationship, or architecture questions when `graphify-out/graph.json` exists, **first** run `graphify query "<question>"` (CLI) or `query_graph` (MCP).
- Use `graphify path "<A>" "<B>"` / `shortest_path` for call chains and `graphify explain "<concept>"` / `get_node` for focused symbol definitions. These return a scoped subgraph, minimizing token consumption compared to raw file reads or greps.
- If `graphify-out/wiki/index.md` exists, navigate it instead of reading raw files.
- Read `graphify-out/GRAPH_REPORT.md` only for broad architecture review or when query/path/explain do not surface sufficient context.

## 2. Post-Modification Update & AST Pruning
- After modifying code files in a session, run `graphify update .` to update graph nodes and candidate relationships.
- Run `go run cmd/graphify-ast-audit/main.go` to prune heuristic phantom edges against actual Go/Python ASTs and assign provenance metadata (`EXTRACTED_AST`, `INFERRED_HEURISTIC`, `PRUNED_PHANTOM`).

## 3. Documentation Link Standards
- Ensure all workspace Markdown links use clean relative paths (`[Label](./target.md)` or `[Label](../path/to/target.md)`).
- Avoid `file://` URIs and code backticks inside brackets (`[`file.md`](...)`) to guarantee graph edge extraction.

## 4. Full Pipeline Synchronization
- Use `./scripts/graphify_sync.sh` to run the master three-stage pipeline sequentially:
  1. Base extraction (`graphify update .`)
  2. Deterministic AST code relationship auditing (`go run cmd/graphify-ast-audit/main.go`)
  3. Markdown documentation link graph validation (`go run cmd/graphify-doc-audit/main.go`)
