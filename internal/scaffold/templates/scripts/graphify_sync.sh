#!/usr/bin/env bash
# ==============================================================================
# Gautama Graph Master Knowledge Graph Synchronization Pipeline
# ==============================================================================
set -euo pipefail

WORKSPACE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export GOWORK=off

echo "=================================================="
echo "🚀 Gautama Graph Knowledge Graph Synchronization"
echo "=================================================="
echo "  Workspace Root : ${WORKSPACE_ROOT}"
echo "  Timestamp      : $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
echo "=================================================="

# Stage 1: Base Graphify Extraction & Encapsulated Binary Execution
echo ""
echo "[1/3] Running Base Knowledge Graph Extraction..."
go run github.com/jameshawkins-art/gautama-graph/cmd/gautama-graph --workspace="${WORKSPACE_ROOT}"

# Stage 2: In-Repo Deterministic AST Code Relationship Audit
echo ""
echo "[2/3] Running In-Repo AST Code Relationship Audit..."
go run github.com/jameshawkins-art/gautama-graph/cmd/graphify-ast-audit -graph="${WORKSPACE_ROOT}/graphify-out/graph.json"

# Stage 3: Markdown Documentation Graph Link & Orphan Audit
echo ""
echo "[3/3] Running In-Repo Markdown Doc Graph Audit..."
go run github.com/jameshawkins-art/gautama-graph/cmd/graphify-doc-audit -workspace="${WORKSPACE_ROOT}"

echo ""
echo "=================================================="
echo "✅ Knowledge Graph Synchronization Complete"
echo "=================================================="
