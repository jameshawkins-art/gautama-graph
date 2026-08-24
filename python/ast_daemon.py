#!/usr/bin/env python3
"""
python/ast_daemon.py
Persistent Python AST Relationship Auditor Daemon for Gautama Graph.
Communicates with Go parent process over stdin/stdout via Newline-Delimited JSON (NDJSON).
"""

import ast
import json
import os
import sys
import time
from typing import Any, Dict, List


def audit_python_file(workspace_root: str, file_path: str, candidates: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    results = []
    
    if not os.path.isabs(file_path):
        clean_path = os.path.abspath(os.path.join(workspace_root, file_path))
    else:
        clean_path = os.path.abspath(file_path)

    if not os.path.exists(clean_path):
        for cand in candidates:
            cand_id = cand.get("id", f"{file_path}->{cand.get('target_symbol', '')}")
            results.append({
                "candidate": cand,
                "id": cand_id,
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
            cand_id = cand.get("id", f"{file_path}->{cand.get('target_symbol', '')}")
            results.append({
                "candidate": cand,
                "id": cand_id,
                "provenance_status": "PRUNED_PHANTOM",
                "confidence": 0.0,
                "error": f"AST parse exception: {str(e)}"
            })
        return results

    # Extract all function calls and selector references from Python AST
    calls = set()
    for node in ast.walk(tree):
        if isinstance(node, ast.Call):
            if isinstance(node.func, ast.Name):
                calls.add(node.func.id)
            elif isinstance(node.func, ast.Attribute):
                calls.add(node.func.attr)
        elif isinstance(node, ast.Attribute):
            calls.add(node.attr)

    for cand in candidates:
        target = cand.get("target_symbol", "")
        cand_id = cand.get("id", f"{file_path}->{target}")
        if target and target in calls:
            results.append({
                "candidate": cand,
                "id": cand_id,
                "provenance_status": "EXTRACTED_AST",
                "confidence": 1.0,
                "ast_node_pattern": f"ast.Call({target})"
            })
        else:
            results.append({
                "candidate": cand,
                "id": cand_id,
                "provenance_status": "PRUNED_PHANTOM",
                "confidence": 0.0
            })

    return results


def main() -> None:
    while True:
        try:
            line = sys.stdin.readline()
            if not line:
                break

            line = line.strip()
            if not line:
                continue

            start_time = time.perf_counter()
            req = json.loads(line)
            req_id = req.get("id", "")
            cmd = req.get("command", "")
            workspace_root = req.get("workspace_root", ".")

            if cmd == "PING":
                res = {
                    "id": req_id,
                    "success": True,
                    "duration_ms": (time.perf_counter() - start_time) * 1000.0
                }
            elif cmd == "AUDIT_CANDIDATES":
                source_file = req.get("source_file", "")
                candidates = req.get("candidates", [])
                audited = audit_python_file(workspace_root, source_file, candidates)
                res = {
                    "id": req_id,
                    "success": True,
                    "audited_edges": audited,
                    "duration_ms": (time.perf_counter() - start_time) * 1000.0
                }
            elif cmd == "SHUTDOWN":
                res = {
                    "id": req_id,
                    "success": True,
                    "duration_ms": (time.perf_counter() - start_time) * 1000.0
                }
                sys.stdout.write(json.dumps(res) + "\n")
                sys.stdout.flush()
                break
            else:
                res = {
                    "id": req_id,
                    "success": False,
                    "error": f"unrecognized command: {cmd}",
                    "duration_ms": (time.perf_counter() - start_time) * 1000.0
                }

            sys.stdout.write(json.dumps(res) + "\n")
            sys.stdout.flush()

        except json.JSONDecodeError as jde:
            err_res = {
                "id": "",
                "success": False,
                "error": f"malformed JSON: {str(jde)}",
                "duration_ms": 0.0
            }
            sys.stdout.write(json.dumps(err_res) + "\n")
            sys.stdout.flush()
        except Exception as e:
            err_res = {
                "id": req.get("id", "") if "req" in locals() and isinstance(req, dict) else "",
                "success": False,
                "error": f"unhandled exception: {str(e)}",
                "duration_ms": 0.0
            }
            sys.stdout.write(json.dumps(err_res) + "\n")
            sys.stdout.flush()


if __name__ == "__main__":
    main()
