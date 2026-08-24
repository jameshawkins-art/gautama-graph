package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// DefaultBinaryManager manages local Graphify binary lifecycle and caching.
type DefaultBinaryManager struct {
	downloader ReleaseDownloader
	mu         sync.Mutex
}

// NewDefaultBinaryManager initializes a DefaultBinaryManager instance.
func NewDefaultBinaryManager(downloader ReleaseDownloader) *DefaultBinaryManager {
	if downloader == nil {
		downloader = NewDefaultReleaseDownloader("")
	}
	return &DefaultBinaryManager{
		downloader: downloader,
	}
}

// ResolvePlatformTarget returns the current runtime OS and Architecture.
func ResolvePlatformTarget() PlatformTarget {
	return PlatformTarget{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

// ResolveDefaultCacheDir returns the default local cache directory for gautama-graph binaries.
func ResolveDefaultCacheDir(workspaceRoot string) string {
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return filepath.Join(home, ".cache", "gautama-graph", "bin")
	}
	if workspaceRoot != "" {
		return filepath.Join(filepath.Clean(workspaceRoot), ".gautama-graph", "bin")
	}
	return filepath.Join(".", ".gautama-graph", "bin")
}

// EnsureBinary locates an existing cached binary or fetches and verifies the latest release from GitHub.
func (m *DefaultBinaryManager) EnsureBinary(ctx context.Context, cfg RunnerConfig) (string, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cacheDir := cfg.CacheDirectoryPath
	if cacheDir == "" {
		cacheDir = ResolveDefaultCacheDir(cfg.WorkspaceRootPath)
	}
	cacheDir = filepath.Clean(cacheDir)

	platform := ResolvePlatformTarget()
	targetBinaryName := fmt.Sprintf("graphify-%s-%s", platform.OS, platform.Arch)
	if platform.OS == "windows" {
		targetBinaryName += ".exe"
	}
	cachedBinaryPath := filepath.Join(cacheDir, targetBinaryName)

	// 1. Check local cache first (if not forced download)
	if !cfg.ForceDownload {
		if info, err := os.Stat(cachedBinaryPath); err == nil && !info.IsDir() {
			if info.Mode()&0111 == 0 {
				_ = os.Chmod(cachedBinaryPath, 0755)
			}
			return cachedBinaryPath, "cached", nil
		}
	}

	// 2. Check if host already has graphify executable on PATH as zero-dependency fallback
	if !cfg.ForceDownload {
		if hostBin, err := exec.LookPath("graphify"); err == nil {
			return hostBin, "system-path", nil
		}
	}

	// 3. Query GitHub Releases API for download
	meta, err := m.downloader.GetLatestRelease(ctx, cfg.GitHubRepoOwner, cfg.GitHubRepoName)
	if err != nil {
		if !cfg.ForceDownload {
			// If network query fails but cached binary exists, fallback gracefully
			if info, statErr := os.Stat(cachedBinaryPath); statErr == nil && !info.IsDir() {
				return cachedBinaryPath, "offline-cached", nil
			}
			// If system path binary exists, fallback to it
			if hostBin, pathErr := exec.LookPath("graphify"); pathErr == nil {
				return hostBin, "offline-system-path", nil
			}
		}
		return "", "", fmt.Errorf("failed to fetch release metadata: %w", err)
	}

	// 4. Match asset for host platform
	var matchedAsset *ReleaseAsset
	for name, asset := range meta.Assets {
		lowerName := strings.ToLower(name)
		osMatch := strings.Contains(lowerName, platform.OS)
		archMatch := strings.Contains(lowerName, platform.Arch) || (platform.Arch == "amd64" && strings.Contains(lowerName, "x86_64"))

		if osMatch && archMatch {
			matchedAsset = &asset
			break
		}
	}

	// Fallback to generic binary or first asset if platform-specific name not found
	if matchedAsset == nil {
		for _, asset := range meta.Assets {
			if strings.Contains(strings.ToLower(asset.Name), "graphify") {
				matchedAsset = &asset
				break
			}
		}
	}

	if matchedAsset == nil {
		// If no assets matched in GitHub release, check if system has graphify
		if hostBin, err := exec.LookPath("graphify"); err == nil {
			return hostBin, meta.TagName, nil
		}
		return "", "", fmt.Errorf("no compatible release asset found for platform %s/%s in release %s", platform.OS, platform.Arch, meta.TagName)
	}

	// 5. Download binary into cache
	versionedBinaryName := fmt.Sprintf("graphify-%s-%s-%s", meta.TagName, platform.OS, platform.Arch)
	if platform.OS == "windows" {
		versionedBinaryName += ".exe"
	}
	destPath := filepath.Join(cacheDir, versionedBinaryName)

	if err := m.downloader.DownloadBinary(ctx, *matchedAsset, destPath); err != nil {
		return "", "", fmt.Errorf("failed downloading binary asset %s: %w", matchedAsset.Name, err)
	}

	// Verify Checksum if provided
	if matchedAsset.ChecksumSHA256 != "" {
		if _, err := m.downloader.VerifyChecksum(destPath, matchedAsset.ChecksumSHA256); err != nil {
			_ = os.Remove(destPath)
			return "", "", fmt.Errorf("downloaded binary failed checksum verification: %w", err)
		}
	}

	// Ensure executable permissions
	_ = os.Chmod(destPath, 0755)

	// Create symlink or copy to targetBinaryName for stable resolution
	_ = os.Remove(cachedBinaryPath)
	if err := os.Symlink(destPath, cachedBinaryPath); err != nil {
		// Fallback: copy file if symlink fails on Windows
		if data, readErr := os.ReadFile(destPath); readErr == nil {
			_ = os.WriteFile(cachedBinaryPath, data, 0755)
		}
	}

	return destPath, meta.TagName, nil
}
