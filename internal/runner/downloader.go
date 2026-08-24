package runner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DefaultReleaseDownloader implements ReleaseDownloader using net/http.
type DefaultReleaseDownloader struct {
	httpClient *http.Client
	apiBaseURL string
}

// NewDefaultReleaseDownloader constructs a DefaultReleaseDownloader instance.
func NewDefaultReleaseDownloader(apiBaseURL string) *DefaultReleaseDownloader {
	if apiBaseURL == "" {
		apiBaseURL = "https://api.github.com"
	}
	return &DefaultReleaseDownloader{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiBaseURL: strings.TrimRight(apiBaseURL, "/"),
	}
}

type gitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
	ChecksumSHA256     string `json:"checksum_sha256,omitempty"`
}

type gitHubReleaseResponse struct {
	TagName     string        `json:"tag_name"`
	PublishedAt time.Time     `json:"published_at"`
	Assets      []gitHubAsset `json:"assets"`
}

// GetLatestRelease queries GitHub releases API for the latest release metadata.
func (d *DefaultReleaseDownloader) GetLatestRelease(ctx context.Context, owner, repo string) (*ReleaseMetadata, error) {
	if owner == "" {
		owner = "Graphify-Labs"
	}
	if repo == "" {
		repo = "graphify"
	}

	endpoint := fmt.Sprintf("%s/repos/%s/%s/releases/latest", d.apiBaseURL, owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create release request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "gautama-graph-runner/1.1.0")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release from %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("github api returned non-200 status %d: %s", resp.StatusCode, string(body))
	}

	var ghResp gitHubReleaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, fmt.Errorf("failed to decode release json: %w", err)
	}

	meta := &ReleaseMetadata{
		TagName:     ghResp.TagName,
		PublishedAt: ghResp.PublishedAt,
		Assets:      make(map[string]ReleaseAsset),
	}

	for _, a := range ghResp.Assets {
		meta.Assets[a.Name] = ReleaseAsset{
			Name:               a.Name,
			BrowserDownloadURL: a.BrowserDownloadURL,
			Size:               a.Size,
			ChecksumSHA256:     a.ChecksumSHA256,
		}
	}

	return meta, nil
}

// DownloadBinary downloads the target asset to destinationPath with atomic .tmp staging.
func (d *DefaultReleaseDownloader) DownloadBinary(ctx context.Context, asset ReleaseAsset, destinationPath string) error {
	cleanDest := filepath.Clean(destinationPath)
	destDir := filepath.Dir(cleanDest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed creating cache directory at %s: %w", destDir, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.BrowserDownloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed creating download request: %w", err)
	}
	req.Header.Set("User-Agent", "gautama-graph-runner/1.1.0")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed downloading asset from %s: %w", asset.BrowserDownloadURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("asset download failed with HTTP %d", resp.StatusCode)
	}

	tmpFile := cleanDest + ".tmp"
	out, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed creating tmp destination file at %s: %w", tmpFile, err)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed streaming asset body: %w", err)
	}
	out.Close()

	if err := os.Rename(tmpFile, cleanDest); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed committing downloaded binary to %s: %w", cleanDest, err)
	}

	return nil
}

// VerifyChecksum validates that filePath matches expectedSHA256.
func (d *DefaultReleaseDownloader) VerifyChecksum(filePath, expectedSHA256 string) (bool, error) {
	if expectedSHA256 == "" {
		return true, nil // No checksum specified
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return false, fmt.Errorf("failed opening file for checksum verification: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return false, fmt.Errorf("failed hashing file contents: %w", err)
	}

	calculated := hex.EncodeToString(hasher.Sum(nil))
	if !strings.EqualFold(calculated, strings.TrimSpace(expectedSHA256)) {
		return false, fmt.Errorf("checksum mismatch: expected %s, calculated %s", expectedSHA256, calculated)
	}

	return true, nil
}
