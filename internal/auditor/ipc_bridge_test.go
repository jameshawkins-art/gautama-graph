package auditor

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"testing"
	"time"
)

func createSyntheticPythonWorkspace(t *testing.T) (string, string) {
	tmpDir := t.TempDir()

	// Copy/link python/ast_daemon.py into test workspace or point to repository daemon
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed finding repo root: %v", err)
	}
	daemonPath := filepath.Join(repoRoot, "python", "ast_daemon.py")

	pyCode := `
import os

def helper_function():
    return "ok"

class WorkerService:
    def execute(self):
        helper_function()
        os.path.abspath(".")
        return True
`
	srcFile := filepath.Join(tmpDir, "service.py")
	if err := os.WriteFile(srcFile, []byte(pyCode), 0644); err != nil {
		t.Fatalf("failed writing python fixture: %v", err)
	}

	return tmpDir, daemonPath
}

func TestIPCWorkerSession_Lifecycle_And_Ping(t *testing.T) {
	wsDir, daemonPath := createSyntheticPythonWorkspace(t)

	session := newWorkerSession(1, wsDir, daemonPath)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := session.start(ctx); err != nil {
		t.Fatalf("failed starting worker session: %v", err)
	}

	if session.Status() != WorkerStateIdle {
		t.Errorf("expected worker state IDLE, got %s", session.Status())
	}

	if session.PID() <= 0 {
		t.Errorf("expected positive PID, got %d", session.PID())
	}

	// Ping test
	lat, err := session.Ping(ctx)
	if err != nil {
		t.Fatalf("ping failed: %v", err)
	}

	if lat <= 0 {
		t.Errorf("expected positive ping latency, got %v", lat)
	}

	// Close session
	if err := session.Close(); err != nil {
		t.Errorf("failed closing session: %v", err)
	}

	if session.Status() != WorkerStateTerminated {
		t.Errorf("expected worker state TERMINATED, got %s", session.Status())
	}
}

func TestIPCWorkerPool_NominalExecution(t *testing.T) {
	wsDir, daemonPath := createSyntheticPythonWorkspace(t)
	srcFile := filepath.Join(wsDir, "service.py")

	pool := NewDefaultIPCWorkerPool(wsDir, daemonPath, 2)
	defer func() { _ = pool.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pool.SpawnWorkers(ctx, 2); err != nil {
		t.Fatalf("failed spawning workers: %v", err)
	}

	candidates := []CandidateEdge{
		{
			ID:           "c1",
			SourceFile:   "service.py",
			SourceSymbol: "execute",
			TargetSymbol: "helper_function",
		},
		{
			ID:           "c2",
			SourceFile:   "service.py",
			SourceSymbol: "execute",
			TargetSymbol: "non_existent_func",
		},
	}

	results, err := pool.AuditPython(ctx, srcFile, candidates)
	if err != nil {
		t.Fatalf("AuditPython failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 audited edges, got %d", len(results))
	}

	// Verify c1 -> EXTRACTED_AST
	if results[0].ProvenanceStatus != string(ProvenanceExtractedAST) {
		t.Errorf("expected c1 provenance %s, got %s", ProvenanceExtractedAST, results[0].ProvenanceStatus)
	}
	if results[0].Confidence != 1.0 {
		t.Errorf("expected c1 confidence 1.0, got %f", results[0].Confidence)
	}

	// Verify c2 -> PRUNED_PHANTOM
	if results[1].ProvenanceStatus != string(ProvenancePrunedPhantom) {
		t.Errorf("expected c2 provenance %s, got %s", ProvenancePrunedPhantom, results[1].ProvenanceStatus)
	}
	if results[1].Confidence != 0.0 {
		t.Errorf("expected c2 confidence 0.0, got %f", results[1].Confidence)
	}

	// Test PoolStats
	stats := pool.Stats()
	if stats.TotalWorkers != 2 {
		t.Errorf("expected 2 total workers, got %d", stats.TotalWorkers)
	}
	if stats.IdleWorkers != 2 {
		t.Errorf("expected 2 idle workers after completion, got %d", stats.IdleWorkers)
	}
}

func TestIPCWorkerPool_CrashRecovery(t *testing.T) {
	wsDir, daemonPath := createSyntheticPythonWorkspace(t)
	srcFile := filepath.Join(wsDir, "service.py")

	pool := NewDefaultIPCWorkerPool(wsDir, daemonPath, 2)
	defer func() { _ = pool.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pool.SpawnWorkers(ctx, 2); err != nil {
		t.Fatalf("failed spawning workers: %v", err)
	}

	// Deliberately kill the first worker process with SIGKILL to simulate a hard crash
	pool.mu.RLock()
	victimPID := pool.workers[0].PID()
	pool.mu.RUnlock()

	_ = syscall.Kill(victimPID, syscall.SIGKILL)
	time.Sleep(100 * time.Millisecond)

	candidates := []CandidateEdge{
		{
			ID:           "c1",
			SourceFile:   "service.py",
			SourceSymbol: "execute",
			TargetSymbol: "helper_function",
		},
	}

	// AuditPython should auto-detect the crash, recycle the worker, and succeed
	results, err := pool.AuditPython(ctx, srcFile, candidates)
	if err != nil {
		t.Fatalf("expected crash recovery and success, got error: %v", err)
	}

	if len(results) != 1 || results[0].ProvenanceStatus != string(ProvenanceExtractedAST) {
		t.Errorf("expected successful recovery with EXTRACTED_AST edge, got %v", results)
	}
}

func TestIPCWorkerPool_ZeroOrphanCleanup(t *testing.T) {
	wsDir, daemonPath := createSyntheticPythonWorkspace(t)

	pool := NewDefaultIPCWorkerPool(wsDir, daemonPath, 3)
	ctx := context.Background()

	if err := pool.SpawnWorkers(ctx, 3); err != nil {
		t.Fatalf("failed spawning workers: %v", err)
	}

	var pids []int
	pool.mu.RLock()
	for _, w := range pool.workers {
		pids = append(pids, w.PID())
	}
	pool.mu.RUnlock()

	if len(pids) != 3 {
		t.Fatalf("expected 3 worker PIDs, got %d", len(pids))
	}

	// Close pool
	if err := pool.Close(); err != nil {
		t.Fatalf("failed closing worker pool: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	// Verify all child PIDs are no longer active in OS process table
	for _, pid := range pids {
		err := syscall.Kill(pid, 0)
		if err == nil {
			t.Errorf("orphan process detected: PID %d is still alive after pool.Close()", pid)
		}
	}
}

func TestIPCWorkerPool_ConcurrencyLoad(t *testing.T) {
	wsDir, daemonPath := createSyntheticPythonWorkspace(t)
	srcFile := filepath.Join(wsDir, "service.py")

	pool := NewDefaultIPCWorkerPool(wsDir, daemonPath, 4)
	defer func() { _ = pool.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pool.SpawnWorkers(ctx, 4); err != nil {
		t.Fatalf("failed spawning workers: %v", err)
	}

	numGoroutines := 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			cands := []CandidateEdge{
				{
					ID:           "concurrent-c1",
					SourceFile:   "service.py",
					SourceSymbol: "execute",
					TargetSymbol: "helper_function",
				},
			}
			res, err := pool.AuditPython(ctx, srcFile, cands)
			if err != nil || len(res) != 1 || res[0].ProvenanceStatus != string(ProvenanceExtractedAST) {
				t.Errorf("goroutine %d failed audit: %v, res: %v", idx, err, res)
			}
		}(i)
	}

	wg.Wait()
}

func TestIPCWorkerPool_PathBoundary_And_Closed(t *testing.T) {
	wsDir, daemonPath := createSyntheticPythonWorkspace(t)

	pool := NewDefaultIPCWorkerPool(wsDir, daemonPath, 2)
	defer func() { _ = pool.Close() }()

	ctx := context.Background()

	// 1. Path boundary escape test
	_, errBoundary := pool.AuditPython(ctx, "../outside.py", []CandidateEdge{})
	if errBoundary == nil {
		t.Errorf("expected error for path boundary escape, got nil")
	}

	// 2. Closed pool test
	_ = pool.Close()
	_, errClosed := pool.AuditPython(ctx, filepath.Join(wsDir, "service.py"), []CandidateEdge{})
	if errClosed == nil {
		t.Errorf("expected error calling closed pool, got nil")
	}
}

func TestEngine_Python_IPC_Integration(t *testing.T) {
	wsDir, daemonPath := createSyntheticPythonWorkspace(t)
	srcFile := filepath.Join(wsDir, "service.py")

	cfg := Config{
		WorkspaceRootPath: wsDir,
		AuditorTimeout:    10 * time.Second,
	}
	engine := NewDefaultEngine(cfg)
	defer func() { _ = engine.Close() }()

	pool := NewDefaultIPCWorkerPool(wsDir, daemonPath, 2)
	engine.SetIPCPool(pool)

	candidates := []CandidateEdge{
		{
			ID:           "py-c1",
			SourceFile:   srcFile,
			SourceSymbol: "execute",
			TargetSymbol: "helper_function",
		},
		{
			ID:           "py-c2",
			SourceFile:   srcFile,
			SourceSymbol: "execute",
			TargetSymbol: "unknown_target",
		},
	}

	results, err := engine.AuditCandidates(context.Background(), candidates)
	if err != nil {
		t.Fatalf("AuditCandidates with IPC failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0].ProvenanceStatus != string(ProvenanceExtractedAST) {
		t.Errorf("expected py-c1 EXTRACTED_AST, got %s", results[0].ProvenanceStatus)
	}

	if results[1].ProvenanceStatus != string(ProvenancePrunedPhantom) {
		t.Errorf("expected py-c2 PRUNED_PHANTOM, got %s", results[1].ProvenanceStatus)
	}
}
