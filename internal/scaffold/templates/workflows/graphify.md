---
description: Turn any folder of files into a navigable knowledge graph
globs: *
alwaysApply: false
---

# Graphify Workflow Playbook

Run this workflow whenever turning a codebase, research folder, or documentation tree into a persistent, navigable knowledge graph.

## Step 1: Run Gautama Graph Pipeline
Run the unified multi-stage extraction and auditing pipeline:
```bash
GOWORK=off go run github.com/jameshawkins-art/gautama-graph/cmd/gautama-graph
```
Or execute the local shell script:
```bash
./scripts/graphify_sync.sh
```

## Step 2: Query Knowledge Graph
Use targeted knowledge queries to inspect architecture and call chains without token bloat:
```bash
# Query concept or relationship
graphify query "<concept or question>"

# Find shortest call path between two symbols
graphify path "<source_symbol>" "<target_symbol>"

# Explain specific symbol or type definition
graphify explain "<type_or_function_name>"
```

## Step 3: Run Deterministic Relationship & Doc Audits
```bash
# Audit Go/Python AST relationships and prune phantoms
GOWORK=off go run github.com/jameshawkins-art/gautama-graph/cmd/graphify-ast-audit

# Audit Markdown doc link topology, orphans, and broken links
GOWORK=off go run github.com/jameshawkins-art/gautama-graph/cmd/graphify-doc-audit

# Auto-remediate broken Markdown relative links in-place
GOWORK=off go run github.com/jameshawkins-art/gautama-graph/cmd/graphify-doc-audit --fix
```
