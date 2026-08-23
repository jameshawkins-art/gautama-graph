#!/usr/bin/env python3
"""
python/ast_auditor_bridge.py
Python AST Auditor Bridge for Graphify Knowledge Graph Parsing Verification
"""

import ast
import json
import os
import sys
from typing import Any, Dict, List


def audit_python_file(file_path: str, candidates: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    results = []
    clean_path = os.path.abspath(file_path)

    if not os.path.exists(clean_path):
        for cand in candidates:
            results.append({
                "id": cand.get("id", ""),
                "provenance_status": "PRUNED_PHANTOM",
                "confidence": 0.0,
                "error": f"file {clean_path} does not exist"
            })
        return results

    try:
        with open(clean_path, "r", encoding="utf-8") as f:
            tree = ast.parse(f.read(), filename=clean_path)
    except Exception as e:
        for cand in candidates:
            results.append({
                "id": cand.get("id", ""),
                "provenance_status": "PRUNED_PHANTOM",
                "confidence": 0.0,
                "error": f"AST parse exception: {str(e)}"
            })
        return results

    # Extract all function calls from Python AST
    calls = set()
    for node in ast.walk(tree):
        if isinstance(node, ast.Call):
            if isinstance(node.func, ast.Name):
                calls.add(node.func.id)
            elif isinstance(node.func, ast.Attribute):
                calls.add(node.func.attr)

    for cand in candidates:
        target = cand.get("target_symbol", "")
        cand_id = cand.get("id", f"{file_path}->{target}")
        if target and target in calls:
            results.append({
                "id": cand_id,
                "provenance_status": "EXTRACTED_AST",
                "confidence": 1.0,
                "ast_node_pattern": "ast.Call"
            })
        else:
            results.append({
                "id": cand_id,
                "provenance_status": "PRUNED_PHANTOM",
                "confidence": 0.0
            })

    return results


def main() -> None:
    if len(sys.argv) < 2:
        sys.stderr.write("Usage: ast_auditor_bridge.py <target_file_path>\n")
        sys.exit(1)

    target_file = sys.argv[1]
    if sys.stdin.isatty():
        sys.stdout.write(json.dumps({"status": "success", "audited_edges": []}, indent=2) + "\n")
        sys.exit(0)

    raw_input = sys.stdin.read()

    if not raw_input.strip():
        sys.stdout.write(json.dumps({"status": "success", "audited_edges": []}, indent=2) + "\n")
        sys.exit(0)

    try:
        payload = json.loads(raw_input)
        candidates = payload.get("candidates", [])
        audited = audit_python_file(target_file, candidates)
        sys.stdout.write(json.dumps({"status": "success", "audited_edges": audited}, indent=2))
        sys.stdout.write("\n")
    except Exception as e:
        sys.stderr.write(f"AST Auditor Bridge Execution Failure: {str(e)}\n")
        sys.exit(1)


if __name__ == "__main__":
    main()
