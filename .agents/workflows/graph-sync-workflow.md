# Knowledge Graph Synchronization & Integrity Workflow Guide

This workflow governs the deterministic validation, phantom pruning, and topological health auditing of the Gautama knowledge graph (`graphify-out/graph.json` and `graphify-out/doc_graph_audit.json`).

## Pipeline Execution Stages
1. **Stage 1: Heuristic Extraction & Community Detection**
   - Command: `graphify update .`
   - Role: Generates base knowledge graph nodes and candidate code/document relationships.
2. **Stage 2: Deterministic AST Code Relationship Auditing**
   - Command: `go run cmd/graphify-ast-audit/main.go`
   - Role: Parses Go AST (`go/ast`) and Python AST (`python/ast_auditor_bridge.py`), evaluates selector expressions, applies provenance tags (`EXTRACTED_AST`, `INFERRED_HEURISTIC`, `PRUNED_PHANTOM`), and prunes unverified phantom edges.
3. **Stage 3: Markdown Documentation Graph Auditing**
   - Command: `go run cmd/graphify-doc-audit/main.go`
   - Role: Parses all workspace Markdown files, strips code blocks, validates relative link targets against disk, computes in/out degrees, and reports orphan documents (`InDegree == 0`).
4. **Stage 4: Atomic Store Persistence & Verification**
   - Enforces two-phase atomic write protocol (`.tmp` buffer + `os.Rename`) on all graph output files.

## Master Script Invocation
Execute all stages sequentially using the master script:
```bash
./scripts/graphify_sync.sh
```

Under `--strict` mode, any pruned phantom edges, broken links, or orphan documents will halt the pipeline with a non-zero exit code.
