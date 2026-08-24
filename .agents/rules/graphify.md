---
description: Rules for querying, updating, and maintaining the Graphify knowledge graph and doc topology.
always_on: true
---

# Graphify Knowledge Graph Rules

1. **Discovery First**: Before performing broad file searches or greps, query the graph via `graphify query "<topic>"`, `graphify path "<A>" "<B>"`, or inspect `graphify-out/wiki/index.md`.
2. **Deterministic Pruning**: Always run `go run cmd/graphify-ast-audit/main.go` after code changes to eliminate phantom edges.
3. **Doc Link Standards**: Ensure Markdown links use clean relative paths (`[Label](./target.md)`). Avoid `file://` URIs and backticks in brackets.
4. **Synchronization**: Use `./scripts/graphify_sync.sh` to run the full extraction, AST auditing, and doc validation pipeline.
