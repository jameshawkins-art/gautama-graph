package runner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

// DefaultSubprocessRunner executes Graphify commands via Go standard library exec.CommandContext.
type DefaultSubprocessRunner struct {
	maxOutputBytes int64
}

// NewDefaultSubprocessRunner initializes a DefaultSubprocessRunner instance.
func NewDefaultSubprocessRunner() *DefaultSubprocessRunner {
	return &DefaultSubprocessRunner{
		maxOutputBytes: 10 * 1024 * 1024, // 10MB memory limit
	}
}

// ExecuteCommand executes a binary command with discrete argument passing and stream separation.
func (r *DefaultSubprocessRunner) ExecuteCommand(ctx context.Context, binaryPath, workspaceRoot string, args ...string) ([]byte, []byte, error) {
	cleanBinary := filepath.Clean(binaryPath)
	cleanRoot := filepath.Clean(workspaceRoot)

	cmd := exec.CommandContext(ctx, cleanBinary, args...)
	cmd.Dir = cleanRoot
	cmd.Stdin = nil // Headless execution

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	stdoutBytes := stdoutBuf.Bytes()
	if int64(len(stdoutBytes)) > r.maxOutputBytes {
		stdoutBytes = stdoutBytes[:r.maxOutputBytes]
	}

	stderrBytes := stderrBuf.Bytes()
	if int64(len(stderrBytes)) > r.maxOutputBytes {
		stderrBytes = stderrBytes[:r.maxOutputBytes]
	}

	if err != nil {
		return stdoutBytes, stderrBytes, fmt.Errorf("subprocess execution failed: %w (stderr: %s)", err, string(stderrBytes))
	}

	return stdoutBytes, stderrBytes, nil
}
