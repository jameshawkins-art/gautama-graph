---
trigger: always_on
description: Consult the graphify knowledge graph at graphify-out/ for codebase and architecture questions.
---

## graphify

This project has a graphify knowledge graph at graphify-out/ maintained by `github.com/jameshawkins-art/gautama-graph` v1.7.0.

Rules:
- For codebase or architecture questions, when `graphify-out/graph.json` exists, first run `gautama-graph query "<question>"` (or `make graphify-query Q="<question>"`) or `query_graph` (MCP). Use `gautama-graph path "<A>" "<B>"` (or `make graphify-path A="<A>" B="<B>"`) / `shortest_path` for relationships and `gautama-graph explain "<concept>"` (or `make graphify-explain C="<concept>"`) / `get_node` for focused concepts. These return a scoped subgraph, usually much smaller than `GRAPH_REPORT.md` or raw grep output.
- If graphify-out/wiki/index.md exists, navigate it instead of reading raw files.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code files in this session, run `make graphify-update` or `make audit` to execute the turnkey 4-stage pipeline (binary resolution, graphify update, AST code audit, and doc graph audit).
- If broken markdown links are reported across `.agents/` or `docs/`, run `make audit-remediate` to automatically compute and apply canonical relative path fixes.
