#!/usr/bin/env bash
# ==============================================================================
# Gautama Graph - Master Graphify Knowledge Graph Synchronization Runner
# Pipeline: 1. Base Graphify Extraction -> 2. AST Code Audit -> 3. Doc Graph Audit
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKSPACE_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
export GOWORK=off

# ANSI Color Definitions
BOLD="\033[1m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
CYAN="\033[0;36m"
YELLOW="\033[1;33m"
RESET="\033[0m"

echo -e "${BOLD}${BLUE}==================================================${RESET}"
echo -e "${BOLD}${CYAN}🚀 Gautama Graph Graphify Knowledge Graph Sync${RESET}"
echo -e "${BOLD}${BLUE}==================================================${RESET}"
echo -e "  ${BOLD}Workspace Root${RESET} : ${WORKSPACE_ROOT}"
echo -e "  ${BOLD}Timestamp     ${RESET} : $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
echo -e "${BOLD}${BLUE}==================================================${RESET}"

# Step 1: Base Graphify Extraction & Community Detection
echo -e "\n${BOLD}${CYAN}[1/3] Running Base Graphify Extraction...${RESET}"
(cd "${WORKSPACE_ROOT}" && graphify update .)

# Step 2: Go & Python AST Code Relationship Audit (Prunes Phantoms)
echo -e "\n${BOLD}${CYAN}[2/3] Running In-Repo AST Code Relationship Audit...${RESET}"
(cd "${WORKSPACE_ROOT}" && go run cmd/graphify-ast-audit/main.go --workspace "${WORKSPACE_ROOT}")

# Step 3: Markdown Documentation Graph Audit (Validates Links)
echo -e "\n${BOLD}${CYAN}[3/3] Running In-Repo Markdown Doc Graph Audit...${RESET}"
(cd "${WORKSPACE_ROOT}" && go run cmd/graphify-doc-audit/main.go --workspace "${WORKSPACE_ROOT}")

echo -e "\n${BOLD}${GREEN}==================================================${RESET}"
echo -e "${BOLD}${GREEN}✅ Master Graphify Synchronization Complete (0 Errors)${RESET}"
echo -e "${BOLD}${GREEN}==================================================${RESET}"
