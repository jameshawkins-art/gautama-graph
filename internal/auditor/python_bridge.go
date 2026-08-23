package auditor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// DefaultPythonASTBridge implements PythonASTBridge using python/ast_auditor_bridge.py via subprocess.
type DefaultPythonASTBridge struct {
	workspaceRoot string
	pythonBinPath string
	bridgeScript  string
}

// NewDefaultPythonASTBridge constructs a DefaultPythonASTBridge locating the best available Python interpreter.
func NewDefaultPythonASTBridge(workspaceRoot string) *DefaultPythonASTBridge {
	cleanRoot := filepath.Clean(workspaceRoot)
	if cleanRoot == "" || cleanRoot == "." {
		cwd, err := os.Getwd()
		if err == nil {
			cleanRoot = cwd
		}
	}

	// 1. Check virtual environment python
	venvPython := filepath.Join(cleanRoot, ".venv", "bin", "python3")
	pythonBin := venvPython
	if _, err := os.Stat(venvPython); err != nil {
		// Fallback: check PATH for python3
		if pathLook, err := exec.LookPath("python3"); err == nil {
			pythonBin = pathLook
		} else {
			pythonBin = "python3"
		}
	}

	bridgeScript := filepath.Join(cleanRoot, "python", "ast_auditor_bridge.py")

	return &DefaultPythonASTBridge{
		workspaceRoot: cleanRoot,
		pythonBinPath: pythonBin,
		bridgeScript:  bridgeScript,
	}
}

type pyBridgePayload struct {
	TargetFile string          `json:"target_file"`
	Candidates []CandidateEdge `json:"candidates"`
}

type pyBridgeResponse struct {
	Status       string        `json:"status"`
	AuditedEdges []AuditedEdge `json:"audited_edges"`
	Error        string        `json:"error,omitempty"`
}

// AuditPythonCandidates executes python/ast_auditor_bridge.py via subprocess.
func (b *DefaultPythonASTBridge) AuditPythonCandidates(ctx context.Context, targetFile string, candidates []CandidateEdge) ([]AuditedEdge, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if _, err := os.Stat(b.bridgeScript); err != nil {
		// Script missing; return graceful heuristic fallback
		results := make([]AuditedEdge, len(candidates))
		for i, c := range candidates {
			results[i] = AuditedEdge{
				CandidateEdge:    c,
				ProvenanceStatus: "INFERRED_HEURISTIC",
				Confidence:       0.5,
			}
		}
		return results, fmt.Errorf("python bridge script not found at %s: %w", b.bridgeScript, err)
	}

	payload := pyBridgePayload{
		TargetFile: targetFile,
		Candidates: candidates,
	}

	inputBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal python bridge payload: %w", err)
	}

	cmd := exec.CommandContext(ctx, b.pythonBinPath, b.bridgeScript, targetFile)
	cmd.Stdin = bytes.NewReader(inputBytes)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("python bridge subprocess failed: %w (stderr: %s)", err, stderr.String())
	}

	var resp pyBridgeResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal python bridge response (%s): %w", stdout.String(), err)
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("python bridge execution reported failure: %s", resp.Error)
	}

	return resp.AuditedEdges, nil
}
