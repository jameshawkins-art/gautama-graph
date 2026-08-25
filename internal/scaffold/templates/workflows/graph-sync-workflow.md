# Knowledge Graph Synchronization & Integrity Workflow Guide

This workflow governs the deterministic validation, phantom pruning, and topological health auditing of the project knowledge graph (`graphify-out/graph.json`).

## Pipeline Execution Stages
1. **Stage 1: Base Graphify Extraction & Community Detection**
   - Command: `graphify update .` (or `make graphify-update`)
   - Role: Generates base knowledge graph nodes and candidate code/document relationships.
2. **Stage 2: Deterministic AST Code Relationship Auditing**
   - Command: `make audit-ast` (or `graphify-ast-audit`)
   - Role: Parses source ASTs, evaluates selector expressions, applies provenance tags, and prunes unverified phantom edges.
3. **Stage 3: Markdown Documentation Graph Auditing**
   - Command: `make audit-docs` (or `graphify-doc-audit`)
   - Role: Parses all workspace Markdown files, strips code blocks, validates relative link targets against disk, and reports broken links or orphan documents.
4. **Stage 4: Atomic Store Persistence & Verification**
   - Enforces two-phase atomic write protocol (`.tmp` buffer + atomic rename) on all graph output files.

## Master Script Invocation
Execute all stages sequentially using the sync script or Makefile:
```bash
./scripts/graphify_sync.sh
# or
make audit
```
