package runner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jameshawkins-art/gautama-graph/internal/auditor"
)

func TestDefaultReleaseDownloader_GetLatestRelease(t *testing.T) {
	mockRelease := gitHubReleaseResponse{
		TagName:     "v0.9.48",
		PublishedAt: time.Now(),
		Assets: []gitHubAsset{
			{
				Name:               "graphify-linux-amd64",
				BrowserDownloadURL: "http://example.com/download/linux-amd64",
				Size:               10240,
			},
			{
				Name:               "graphify-darwin-arm64",
				BrowserDownloadURL: "http://example.com/download/darwin-arm64",
				Size:               10240,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/Graphify-Labs/graphify/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(mockRelease)
			return
		}
		if r.URL.Path == "/repos/ErrorOwner/ErrorRepo/releases/latest" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found"))
			return
		}
		if r.URL.Path == "/repos/BadJson/BadRepo/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("invalid-json{"))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	downloader := NewDefaultReleaseDownloader(server.URL)
	meta, err := downloader.GetLatestRelease(context.Background(), "", "")
	if err != nil {
		t.Fatalf("unexpected error fetching release: %v", err)
	}

	if meta.TagName != "v0.9.48" {
		t.Errorf("expected tag v0.9.48, got %s", meta.TagName)
	}

	if len(meta.Assets) != 2 {
		t.Errorf("expected 2 assets, got %d", len(meta.Assets))
	}

	// Error path: non-200 status
	_, errNotFound := downloader.GetLatestRelease(context.Background(), "ErrorOwner", "ErrorRepo")
	if errNotFound == nil {
		t.Errorf("expected error for 404 release endpoint, got nil")
	}

	// Error path: malformed JSON
	_, errBadJSON := downloader.GetLatestRelease(context.Background(), "BadJson", "BadRepo")
	if errBadJSON == nil {
		t.Errorf("expected error for invalid json release endpoint, got nil")
	}
}

func TestDefaultReleaseDownloader_DownloadAndVerifyChecksum(t *testing.T) {
	binaryContent := []byte("#!/bin/sh\necho 'mock graphify'\n")
	hasher := sha256.New()
	hasher.Write(binaryContent)
	expectedSHA := hex.EncodeToString(hasher.Sum(nil))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/fail" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(binaryContent)
	}))
	defer server.Close()

	downloader := NewDefaultReleaseDownloader(server.URL)
	tmpDir := t.TempDir()
	destPath := filepath.Join(tmpDir, "bin", "graphify")

	asset := ReleaseAsset{
		Name:               "graphify-linux-amd64",
		BrowserDownloadURL: server.URL + "/graphify-linux-amd64",
		Size:               int64(len(binaryContent)),
	}

	err := downloader.DownloadBinary(context.Background(), asset, destPath)
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}

	// Verify Valid Checksum
	valid, err := downloader.VerifyChecksum(destPath, expectedSHA)
	if err != nil || !valid {
		t.Fatalf("expected checksum match, got valid=%v, err=%v", valid, err)
	}

	// Empty Checksum skips check and returns true
	validEmpty, errEmpty := downloader.VerifyChecksum(destPath, "")
	if errEmpty != nil || !validEmpty {
		t.Fatalf("expected empty checksum to pass, got valid=%v, err=%v", validEmpty, errEmpty)
	}

	// Verify Mismatched Checksum
	invalidSHA := "0000000000000000000000000000000000000000000000000000000000000000"
	validInv, errInv := downloader.VerifyChecksum(destPath, invalidSHA)
	if errInv == nil || validInv {
		t.Errorf("expected checksum mismatch error, got err=%v", errInv)
	}

	// Non-existent file for checksum verification
	_, errMissing := downloader.VerifyChecksum(filepath.Join(tmpDir, "non_existent"), expectedSHA)
	if errMissing == nil {
		t.Errorf("expected error for missing file, got nil")
	}

	// Invalid URL
	errInvalidURL := downloader.DownloadBinary(context.Background(), ReleaseAsset{BrowserDownloadURL: "http://invalid.host:99999/test"}, destPath)
	if errInvalidURL == nil {
		t.Errorf("expected error on invalid URL, got nil")
	}

	// Write error (destination directory exists where tmp file needed)
	invalidDest := filepath.Join(tmpDir, "invalid_dest_dir")
	_ = os.MkdirAll(invalidDest+".tmp", 0755)
	errWrite := downloader.DownloadBinary(context.Background(), asset, invalidDest)
	if errWrite == nil {
		t.Errorf("expected error when tmp file cannot be opened, got nil")
	}
}

func TestResolveDefaultCacheDir(t *testing.T) {
	cache1 := ResolveDefaultCacheDir("")
	if cache1 == "" {
		t.Errorf("expected non-empty default cache dir")
	}

	cache2 := ResolveDefaultCacheDir("/mock/workspace")
	if cache2 == "" {
		t.Errorf("expected non-empty workspace cache dir")
	}
}

func TestDefaultBinaryManager_EnsureBinary_Scenarios(t *testing.T) {
	tmpDir := t.TempDir()
	cacheDir := filepath.Join(tmpDir, "cache")
	_ = os.MkdirAll(cacheDir, 0755)

	binaryContent := []byte("#!/bin/sh\necho test\n")
	hasher := sha256.New()
	hasher.Write(binaryContent)
	expectedSHA := hex.EncodeToString(hasher.Sum(nil))

	platform := ResolvePlatformTarget()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/Graphify-Labs/graphify/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(gitHubReleaseResponse{
				TagName:     "v0.9.48",
				PublishedAt: time.Now(),
				Assets: []gitHubAsset{
					{
						Name:               fmt.Sprintf("graphify-%s-%s", platform.OS, platform.Arch),
						BrowserDownloadURL: "http://" + r.Host + "/download/binary",
						Size:               int64(len(binaryContent)),
						ChecksumSHA256:     expectedSHA,
					},
				},
			})
			return
		}
		if r.URL.Path == "/download/binary" {
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write(binaryContent)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	downloader := NewDefaultReleaseDownloader(server.URL)
	manager := NewDefaultBinaryManager(downloader)

	cfg := RunnerConfig{
		CacheDirectoryPath: cacheDir,
		ForceDownload:      true, // Trigger download pass
		GitHubRepoOwner:    "Graphify-Labs",
		GitHubRepoName:     "graphify",
	}

	// 1. Download scenario with valid ChecksumSHA256
	resolvedPath, version, err := manager.EnsureBinary(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error downloading binary: %v", err)
	}

	if version != "v0.9.48" {
		t.Errorf("expected version v0.9.48, got %s", version)
	}

	if _, err := os.Stat(resolvedPath); err != nil {
		t.Errorf("expected downloaded binary to exist at %s", resolvedPath)
	}

	// 2. Cached scenario (ForceDownload: false) with unexecutable mode check
	_ = os.Chmod(resolvedPath, 0644)
	cfg.ForceDownload = false
	cachedPath, cachedVer, errCached := manager.EnsureBinary(context.Background(), cfg)
	if errCached != nil {
		t.Fatalf("unexpected error resolving cache: %v", errCached)
	}
	if cachedVer != "cached" {
		t.Errorf("expected 'cached', got %s", cachedVer)
	}
	if cachedPath == "" {
		t.Errorf("expected non-empty cached path")
	}

	// 3. Fallback when network fails but cached binary exists
	failDownloader := NewDefaultReleaseDownloader("http://invalid.host:1234")
	failManager := NewDefaultBinaryManager(failDownloader)
	offlinePath, offlineVer, offlineErr := failManager.EnsureBinary(context.Background(), RunnerConfig{
		CacheDirectoryPath: cacheDir,
		ForceDownload:      false,
	})
	if offlineErr != nil {
		t.Fatalf("expected fallback to cached binary, got error: %v", offlineErr)
	}
	if offlineVer != "cached" {
		t.Errorf("expected 'cached', got %s", offlineVer)
	}
	if offlinePath == "" {
		t.Errorf("expected non-empty cached path")
	}

	// 4. Download failure on HTTP error
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer errorServer.Close()
	errDownloader := NewDefaultReleaseDownloader(errorServer.URL)
	errManager := NewDefaultBinaryManager(errDownloader)
	emptyCacheDir := filepath.Join(tmpDir, "empty_cache")
	_, _, expectedErr := errManager.EnsureBinary(context.Background(), RunnerConfig{
		CacheDirectoryPath: emptyCacheDir,
		ForceDownload:      true,
	})
	if expectedErr == nil {
		t.Errorf("expected error when release query fails and no cache exists, got nil")
	}

	// 5. Generic asset matching fallback
	genericServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/Graphify-Labs/graphify/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(gitHubReleaseResponse{
				TagName:     "v0.9.48",
				PublishedAt: time.Now(),
				Assets: []gitHubAsset{
					{
						Name:               "graphify-standalone-binary",
						BrowserDownloadURL: "http://" + r.Host + "/download/generic",
						Size:               int64(len(binaryContent)),
					},
				},
			})
			return
		}
		if r.URL.Path == "/download/generic" {
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write(binaryContent)
			return
		}
		http.NotFound(w, r)
	}))
	defer genericServer.Close()

	genericDownloader := NewDefaultReleaseDownloader(genericServer.URL)
	genericManager := NewDefaultBinaryManager(genericDownloader)
	genericCache := filepath.Join(tmpDir, "generic_cache")
	genPath, genVer, genErr := genericManager.EnsureBinary(context.Background(), RunnerConfig{
		CacheDirectoryPath: genericCache,
		ForceDownload:      true,
	})
	if genErr != nil {
		t.Fatalf("unexpected error downloading generic asset: %v", genErr)
	}
	if genVer != "v0.9.48" {
		t.Errorf("expected version v0.9.48, got %s", genVer)
	}
	if genPath == "" {
		t.Errorf("expected non-empty path")
	}

	// 6. Checksum mismatch during EnsureBinary
	corruptServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/Graphify-Labs/graphify/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(gitHubReleaseResponse{
				TagName:     "v0.9.48",
				PublishedAt: time.Now(),
				Assets: []gitHubAsset{
					{
						Name:               "graphify-linux-amd64",
						BrowserDownloadURL: "http://" + r.Host + "/download/corrupt",
						Size:               int64(len(binaryContent)),
					},
				},
			})
			return
		}
		if r.URL.Path == "/download/corrupt" {
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write([]byte("corrupt binary data"))
			return
		}
		http.NotFound(w, r)
	}))
	defer corruptServer.Close()

	corruptDownloader := NewDefaultReleaseDownloader(corruptServer.URL)
	corruptManager := NewDefaultBinaryManager(corruptDownloader)
	corruptCache := filepath.Join(tmpDir, "corrupt_cache")
	_, _, corruptErr := corruptManager.EnsureBinary(context.Background(), RunnerConfig{
		CacheDirectoryPath: corruptCache,
		ForceDownload:      true,
	})
	_ = corruptErr

	// 7. No assets matched scenario
	emptyAssetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/Graphify-Labs/graphify/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(gitHubReleaseResponse{
				TagName:     "v0.9.48",
				PublishedAt: time.Now(),
				Assets:      []gitHubAsset{},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer emptyAssetServer.Close()

	emptyDownloader := NewDefaultReleaseDownloader(emptyAssetServer.URL)
	emptyManager := NewDefaultBinaryManager(emptyDownloader)
	emptyCache := filepath.Join(tmpDir, "empty_assets_cache")
	_, _, emptyErr := emptyManager.EnsureBinary(context.Background(), RunnerConfig{
		CacheDirectoryPath: emptyCache,
		ForceDownload:      true,
	})
	_ = emptyErr

	_ = expectedSHA
}

func TestDefaultSubprocessRunner_ExecuteCommand(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "mock_runner.sh")
	_ = os.WriteFile(scriptPath, []byte("#!/bin/sh\necho 'runner success'\n"), 0755)

	runner := NewDefaultSubprocessRunner()
	stdout, stderr, err := runner.ExecuteCommand(context.Background(), scriptPath, tmpDir)
	if err != nil {
		t.Fatalf("execution failed: %v (stderr: %s)", err, string(stderr))
	}

	if string(stdout) != "runner success\n" {
		t.Errorf("expected 'runner success\\n', got %q", string(stdout))
	}

	// Error execution
	failScript := filepath.Join(tmpDir, "fail_runner.sh")
	_ = os.WriteFile(failScript, []byte("#!/bin/sh\necho 'error output' >&2\nexit 1\n"), 0755)
	_, _, errFail := runner.ExecuteCommand(context.Background(), failScript, tmpDir)
	if errFail == nil {
		t.Errorf("expected error on exit code 1, got nil")
	}
}

type mockASTAuditor struct {
	shouldFail bool
}

func (m *mockASTAuditor) AuditGraphFile(ctx context.Context, graphPath string, verbose bool) (*auditor.ASTAuditReport, error) {
	if m.shouldFail {
		return nil, errors.New("mock ast audit failure")
	}
	return &auditor.ASTAuditReport{
		TotalEdges:         10,
		VerifiedASTCount:   5,
		PrunedPhantomCount: 2,
		HeuristicCount:     3,
	}, nil
}

type mockDocAuditor struct {
	shouldFail bool
}

func (m *mockDocAuditor) AuditDocGraph(ctx context.Context) (*auditor.DocAuditReport, error) {
	if m.shouldFail {
		return nil, errors.New("mock doc audit failure")
	}
	return &auditor.DocAuditReport{
		TotalDocNodes:   4,
		TotalDocEdges:   8,
		OrphanCount:     0,
		BrokenLinkCount: 0,
	}, nil
}

type mockBinaryManager struct {
	binaryPath string
	version    string
	shouldFail bool
}

func (m *mockBinaryManager) EnsureBinary(ctx context.Context, cfg RunnerConfig) (string, string, error) {
	if m.shouldFail {
		return "", "", errors.New("mock ensure binary failure")
	}
	return m.binaryPath, m.version, nil
}

type mockSubRunner struct {
	shouldFail bool
}

func (m *mockSubRunner) ExecuteCommand(ctx context.Context, binaryPath, workspaceRoot string, args ...string) ([]byte, []byte, error) {
	if m.shouldFail {
		return nil, []byte("mock command failed"), errors.New("mock command error")
	}
	return []byte("mock extraction complete"), nil, nil
}

func TestDefaultOrchestrator_RunPipeline_FullSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	mockBM := &mockBinaryManager{
		binaryPath: "/bin/echo",
		version:    "v0.9.48",
	}
	mockSR := &mockSubRunner{}
	mockAST := &mockASTAuditor{}
	mockDoc := &mockDocAuditor{}

	orchestrator := NewDefaultOrchestrator(mockBM, mockSR, mockAST, mockDoc)

	cfg := RunnerConfig{
		WorkspaceRootPath: tmpDir,
		ExecutionTimeout:  10 * time.Second,
	}

	report, err := orchestrator.RunPipeline(context.Background(), cfg)
	if err != nil {
		t.Fatalf("pipeline run failed: %v", err)
	}

	if len(report.Stages) != 4 {
		t.Errorf("expected 4 stages, got %d", len(report.Stages))
	}

	for _, s := range report.Stages {
		if !s.Success {
			t.Errorf("stage %s failed: %s", s.StageName, s.Error)
		}
	}

	if report.PrunedPhantoms != 2 {
		t.Errorf("expected 2 pruned phantoms, got %d", report.PrunedPhantoms)
	}

	if report.GraphNodeCount != 4 {
		t.Errorf("expected 4 doc nodes, got %d", report.GraphNodeCount)
	}
}

func TestDefaultOrchestrator_RunPipeline_StageFailures(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Stage 1 Failure (Binary Manager)
	orch1 := NewDefaultOrchestrator(
		&mockBinaryManager{shouldFail: true},
		&mockSubRunner{},
		&mockASTAuditor{},
		&mockDocAuditor{},
	)
	_, err1 := orch1.RunPipeline(context.Background(), RunnerConfig{WorkspaceRootPath: tmpDir})
	if err1 == nil {
		t.Errorf("expected error when stage 1 fails, got nil")
	}

	// 2. Stage 2 Failure (Subprocess Runner)
	orch2 := NewDefaultOrchestrator(
		&mockBinaryManager{binaryPath: "/bin/echo"},
		&mockSubRunner{shouldFail: true},
		&mockASTAuditor{},
		&mockDocAuditor{},
	)
	_, err2 := orch2.RunPipeline(context.Background(), RunnerConfig{WorkspaceRootPath: tmpDir})
	if err2 == nil {
		t.Errorf("expected error when stage 2 fails, got nil")
	}

	// 3. Stage 3 Failure (AST Auditor)
	orch3 := NewDefaultOrchestrator(
		&mockBinaryManager{binaryPath: "/bin/echo"},
		&mockSubRunner{},
		&mockASTAuditor{shouldFail: true},
		&mockDocAuditor{},
	)
	_, err3 := orch3.RunPipeline(context.Background(), RunnerConfig{WorkspaceRootPath: tmpDir})
	if err3 == nil {
		t.Errorf("expected error when stage 3 fails, got nil")
	}

	// 4. Stage 4 Failure (Doc Auditor)
	orch4 := NewDefaultOrchestrator(
		&mockBinaryManager{binaryPath: "/bin/echo"},
		&mockSubRunner{},
		&mockASTAuditor{},
		&mockDocAuditor{shouldFail: true},
	)
	_, err4 := orch4.RunPipeline(context.Background(), RunnerConfig{WorkspaceRootPath: tmpDir})
	if err4 == nil {
		t.Errorf("expected error when stage 4 fails, got nil")
	}
}

func TestNewStandardOrchestrator(t *testing.T) {
	tmpDir := t.TempDir()
	orch := NewStandardOrchestrator(tmpDir)
	if orch == nil {
		t.Fatalf("expected non-nil StandardOrchestrator")
	}
	if orch.binaryManager == nil || orch.subprocessRunner == nil || orch.astAuditor == nil || orch.docAuditor == nil {
		t.Errorf("expected all orchestrator subsystems to be initialized")
	}

	dl := NewDefaultReleaseDownloader("")
	if dl.apiBaseURL != "https://api.github.com" {
		t.Errorf("expected default base URL https://api.github.com, got %s", dl.apiBaseURL)
	}

	bm := NewDefaultBinaryManager(nil)
	if bm.downloader == nil {
		t.Errorf("expected default downloader in binary manager")
	}
}
