package auditor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// workerSession represents an individual persistent Python worker daemon subprocess.
type workerSession struct {
	id            int
	cmd           *exec.Cmd
	stdin         io.WriteCloser
	stdoutScanner *bufio.Scanner
	mu            sync.Mutex
	state         WorkerState
	stats         WorkerStats
	workspaceRoot string
	daemonPath    string
}

func newWorkerSession(id int, workspaceRoot, daemonPath string) *workerSession {
	return &workerSession{
		id:            id,
		workspaceRoot: workspaceRoot,
		daemonPath:    daemonPath,
		state:         WorkerStateStarting,
		stats: WorkerStats{
			WorkerID: id,
			State:    WorkerStateStarting,
		},
	}
}

func (w *workerSession) start(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	cleanDaemon := filepath.Clean(w.daemonPath)
	if !filepath.IsAbs(cleanDaemon) && w.workspaceRoot != "" {
		cleanDaemon = filepath.Join(w.workspaceRoot, w.daemonPath)
	}

	cmd := exec.CommandContext(ctx, "python3", cleanDaemon)
	cmd.Dir = w.workspaceRoot
	cmd.Env = append(os.Environ(), "PYTHONUNBUFFERED=1")

	// Enforce process group isolation on POSIX
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed creating stdin pipe for worker %d: %w", w.id, err)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdinPipe.Close()
		return fmt.Errorf("failed creating stdout pipe for worker %d: %w", w.id, err)
	}

	if err := cmd.Start(); err != nil {
		_ = stdinPipe.Close()
		_ = stdoutPipe.Close()
		return fmt.Errorf("failed starting worker %d (cmd: python3 %s): %w", w.id, cleanDaemon, err)
	}

	w.cmd = cmd
	w.stdin = stdinPipe
	w.stdoutScanner = bufio.NewScanner(stdoutPipe)
	// Allow large responses up to 10MB
	buf := make([]byte, 64*1024)
	w.stdoutScanner.Buffer(buf, 10*1024*1024)

	w.state = WorkerStateIdle
	w.stats.PID = cmd.Process.Pid
	w.stats.State = WorkerStateIdle
	w.stats.LastActive = time.Now()

	return nil
}

// Send dispatches an IPCRequest and reads the correlated IPCResponse from the worker.
func (w *workerSession) Send(ctx context.Context, req *IPCRequest) (*IPCResponse, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state == WorkerStateCrashed || w.state == WorkerStateTerminated || w.cmd == nil || w.cmd.Process == nil {
		return nil, fmt.Errorf("worker %d is not running (state: %s)", w.id, w.state)
	}

	w.state = WorkerStateBusy
	w.stats.State = WorkerStateBusy
	w.stats.RequestsTotal++
	w.stats.LastActive = time.Now()

	defer func() {
		if w.state == WorkerStateBusy {
			w.state = WorkerStateIdle
			w.stats.State = WorkerStateIdle
		}
	}()

	payload, err := json.Marshal(req)
	if err != nil {
		w.stats.ErrorsTotal++
		return nil, fmt.Errorf("failed marshaling IPCRequest: %w", err)
	}

	payload = append(payload, '\n')

	// Write request to child stdin
	if _, err := w.stdin.Write(payload); err != nil {
		w.state = WorkerStateCrashed
		w.stats.State = WorkerStateCrashed
		w.stats.ErrorsTotal++
		w.killProcess()
		return nil, fmt.Errorf("failed writing to worker %d stdin: %w", w.id, err)
	}

	// Read response from child stdout
	type scanResult struct {
		line []byte
		err  error
	}

	resCh := make(chan scanResult, 1)
	go func() {
		if w.stdoutScanner.Scan() {
			text := w.stdoutScanner.Bytes()
			resCh <- scanResult{line: append([]byte(nil), text...)}
		} else {
			scanErr := w.stdoutScanner.Err()
			if scanErr == nil {
				scanErr = io.EOF
			}
			resCh <- scanResult{err: scanErr}
		}
	}()

	select {
	case <-ctx.Done():
		w.state = WorkerStateCrashed
		w.stats.State = WorkerStateCrashed
		w.stats.ErrorsTotal++
		w.killProcess()
		return nil, fmt.Errorf("worker %d timed out awaiting response: %w", w.id, ctx.Err())

	case res := <-resCh:
		if res.err != nil {
			w.state = WorkerStateCrashed
			w.stats.State = WorkerStateCrashed
			w.stats.ErrorsTotal++
			w.killProcess()
			return nil, fmt.Errorf("worker %d stdout scan failure: %w", w.id, res.err)
		}

		var response IPCResponse
		if err := json.Unmarshal(res.line, &response); err != nil {
			w.stats.ErrorsTotal++
			return nil, fmt.Errorf("worker %d malformed JSON response: %w", w.id, err)
		}

		return &response, nil
	}
}

func (w *workerSession) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	req := &IPCRequest{
		ID:            fmt.Sprintf("ping-%d-%d", w.id, start.UnixNano()),
		Command:       CmdPing,
		WorkspaceRoot: w.workspaceRoot,
		Timestamp:     start,
	}

	resp, err := w.Send(ctx, req)
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, fmt.Errorf("ping returned unsuccessful: %s", resp.Error)
	}

	return time.Since(start), nil
}

func (w *workerSession) Status() WorkerState {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.state
}

func (w *workerSession) PID() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.stats.PID
}

func (w *workerSession) killProcess() {
	if w.stdin != nil {
		_ = w.stdin.Close()
	}
	if w.cmd != nil && w.cmd.Process != nil {
		pid := w.cmd.Process.Pid
		_ = syscall.Kill(-pid, syscall.SIGKILL)
		_ = w.cmd.Wait()
	}
}

func (w *workerSession) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state == WorkerStateTerminated {
		return nil
	}

	w.state = WorkerStateTerminated
	w.stats.State = WorkerStateTerminated

	// Send shutdown command
	if w.stdin != nil {
		shutdownReq := IPCRequest{
			ID:        fmt.Sprintf("shutdown-%d", w.id),
			Command:   CmdShutdown,
			Timestamp: time.Now(),
		}
		data, _ := json.Marshal(shutdownReq)
		_, _ = w.stdin.Write(append(data, '\n'))
		_ = w.stdin.Close()
	}

	// Give worker brief grace period before SIGKILL
	if w.cmd != nil && w.cmd.Process != nil {
		done := make(chan error, 1)
		go func() {
			done <- w.cmd.Wait()
		}()

		select {
		case <-time.After(500 * time.Millisecond):
			_ = syscall.Kill(-w.cmd.Process.Pid, syscall.SIGKILL)
			<-done
		case <-done:
		}
	}

	return nil
}

// DefaultIPCWorkerPool manages a persistent pool of Python worker daemon subprocesses.
type DefaultIPCWorkerPool struct {
	mu            sync.RWMutex
	workers       []*workerSession
	idleQueue     chan *workerSession
	workspaceRoot string
	daemonPath    string
	poolSize      int
	closed        bool
	reqCounter    int64
}

// NewDefaultIPCWorkerPool initializes a DefaultIPCWorkerPool instance.
func NewDefaultIPCWorkerPool(workspaceRoot, daemonPath string, poolSize int) *DefaultIPCWorkerPool {
	if poolSize <= 0 {
		poolSize = runtime.NumCPU()
		if poolSize < 2 {
			poolSize = 2
		}
		if poolSize > 8 {
			poolSize = 8
		}
	}

	if daemonPath == "" {
		daemonPath = "python/ast_daemon.py"
	}

	return &DefaultIPCWorkerPool{
		workspaceRoot: workspaceRoot,
		daemonPath:    daemonPath,
		poolSize:      poolSize,
		idleQueue:     make(chan *workerSession, poolSize*2),
	}
}

// SpawnWorkers starts the pool workers up to poolSize.
func (p *DefaultIPCWorkerPool) SpawnWorkers(ctx context.Context, poolSize int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if poolSize > 0 {
		p.poolSize = poolSize
	}

	for i := 0; i < p.poolSize; i++ {
		w := newWorkerSession(i+1, p.workspaceRoot, p.daemonPath)
		if err := w.start(ctx); err != nil {
			return fmt.Errorf("failed spawning worker %d: %w", i+1, err)
		}
		p.workers = append(p.workers, w)
		p.idleQueue <- w
	}

	return nil
}

// AuditPython dispatches candidate edges to an idle worker daemon over streaming NDJSON pipes.
func (p *DefaultIPCWorkerPool) AuditPython(ctx context.Context, sourceFile string, candidates []CandidateEdge) ([]AuditedEdge, error) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil, fmt.Errorf("IPC worker pool is closed")
	}
	if len(p.workers) == 0 {
		p.mu.RUnlock()
		// Auto-initialize workers on first call if not spawned
		if err := p.SpawnWorkers(ctx, p.poolSize); err != nil {
			return nil, fmt.Errorf("failed auto-spawning IPC workers: %w", err)
		}
	} else {
		p.mu.RUnlock()
	}

	cleanFile := filepath.Clean(sourceFile)
	if p.workspaceRoot != "" {
		if _, err := ValidatePathBoundary(p.workspaceRoot, cleanFile); err != nil {
			return nil, err
		}
	}

	reqID := fmt.Sprintf("req-%d-%d", atomic.AddInt64(&p.reqCounter, 1), time.Now().UnixNano())
	req := &IPCRequest{
		ID:            reqID,
		Command:       CmdAuditPythonCandidates,
		WorkspaceRoot: p.workspaceRoot,
		SourceFile:    cleanFile,
		Candidates:    candidates,
		Timestamp:     time.Now(),
	}

	// Retry loop for transparent recovery if a worker crashes
	maxRetries := 2
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case worker := <-p.idleQueue:
			resp, err := worker.Send(ctx, req)
			if err != nil {
				lastErr = err
				// Worker crashed or timed out: respawn replacement
				go p.recycleWorker(worker)
				continue
			}

			// Success: return worker to idle queue
			p.idleQueue <- worker

			if !resp.Success {
				return nil, fmt.Errorf("worker execution error: %s", resp.Error)
			}

			return resp.AuditedEdges, nil
		}
	}

	return nil, fmt.Errorf("all %d worker dispatch attempts failed: %w", maxRetries+1, lastErr)
}

func (p *DefaultIPCWorkerPool) recycleWorker(crashed *workerSession) {
	crashed.killProcess()
	crashed.stats.RestartsTotal++

	// Spawn replacement
	replacement := newWorkerSession(crashed.id, p.workspaceRoot, p.daemonPath)
	replacement.stats.RestartsTotal = crashed.stats.RestartsTotal
	if err := replacement.start(context.Background()); err == nil {
		p.mu.Lock()
		for i, w := range p.workers {
			if w.id == crashed.id {
				p.workers[i] = replacement
				break
			}
		}
		p.mu.Unlock()
		p.idleQueue <- replacement
	}
}

// Stats returns diagnostic metrics across all daemon workers in the pool.
func (p *DefaultIPCWorkerPool) Stats() PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := PoolStats{
		TotalWorkers: len(p.workers),
		Workers:      make([]WorkerStats, 0, len(p.workers)),
	}

	for _, w := range p.workers {
		w.mu.Lock()
		s := w.stats
		w.mu.Unlock()

		switch s.State {
		case WorkerStateIdle:
			stats.IdleWorkers++
		case WorkerStateBusy:
			stats.BusyWorkers++
		}
		stats.Workers = append(stats.Workers, s)
	}

	return stats
}

// Close terminates all Python daemon subprocesses cleanly.
func (p *DefaultIPCWorkerPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}
	p.closed = true

	var errs []error
	for _, w := range p.workers {
		if err := w.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	p.workers = nil

	if len(errs) > 0 {
		return fmt.Errorf("errors closing worker pool: %v", errs)
	}
	return nil
}
