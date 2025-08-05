package gollama

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	llamaCppRepo      = "ggml-org/llama.cpp"
	githubReleasesAPI = "https://api.github.com/repos"
	githubReleasesURL = "https://github.com/ggml-org/llama.cpp/releases/download"
	downloadTimeout   = 10 * time.Minute
	userAgent         = "gollama.cpp/1.0.0"
)

// ReleaseInfo represents GitHub release information
type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
}

// LibraryDownloader handles downloading pre-built llama.cpp binaries
type LibraryDownloader struct {
	cacheDir   string
	userAgent  string
	httpClient *http.Client
}

// NewLibraryDownloader creates a new library downloader instance
func NewLibraryDownloader() (*LibraryDownloader, error) {
	// Create cache directory in user's cache or temp directory
	var cacheDir string

	// Try user cache directory first
	userCacheDir, err := os.UserCacheDir()
	if err == nil {
		cacheDir = filepath.Join(userCacheDir, "gollama", "libs")
	} else {
		// Fallback to temp directory
		cacheDir = filepath.Join(os.TempDir(), "gollama", "libs")
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &LibraryDownloader{
		cacheDir:   cacheDir,
		userAgent:  userAgent,
		httpClient: &http.Client{Timeout: downloadTimeout},
	}, nil
}

// GetLatestRelease fetches the latest release information from GitHub
func (d *LibraryDownloader) GetLatestRelease() (*ReleaseInfo, error) {
	url := fmt.Sprintf("%s/%s/releases/latest", githubReleasesAPI, llamaCppRepo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", d.userAgent)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release info: %w", err)
	}

	return &release, nil
}

// GetReleaseByTag fetches release information for a specific tag
func (d *LibraryDownloader) GetReleaseByTag(tag string) (*ReleaseInfo, error) {
	url := fmt.Sprintf("%s/%s/releases/tags/%s", githubReleasesAPI, llamaCppRepo, tag)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", d.userAgent)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("release %s not found", tag)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release info: %w", err)
	}

	return &release, nil
}

// GetPlatformAssetPattern returns the asset name pattern for the current platform
func (d *LibraryDownloader) GetPlatformAssetPattern() (string, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Convert Go arch to llama.cpp naming convention
	var arch string
	switch goarch {
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "arm64"
	default:
		return "", fmt.Errorf("unsupported architecture: %s", goarch)
	}

	switch goos {
	case "darwin":
		return fmt.Sprintf("llama-.*-bin-macos-%s.zip", arch), nil
	case "linux":
		// Prefer CPU version for compatibility, could be enhanced to detect GPU capabilities
		return fmt.Sprintf("llama-.*-bin-ubuntu-%s.zip", arch), nil
	case "windows":
		// Start with CPU version for compatibility
		return fmt.Sprintf("llama-.*-bin-win-cpu-%s.zip", arch), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", goos)
	}
}

// FindAssetByPattern finds an asset that matches the given pattern
func (d *LibraryDownloader) FindAssetByPattern(release *ReleaseInfo, pattern string) (string, string, error) {
	// Compile the pattern as a regular expression
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return "", "", fmt.Errorf("invalid pattern: %w", err)
	}

	for _, asset := range release.Assets {
		if regex.MatchString(asset.Name) {
			return asset.Name, asset.BrowserDownloadURL, nil
		}
	}
	return "", "", fmt.Errorf("no asset found matching pattern: %s", pattern)
}

// DownloadAndExtract downloads and extracts the library archive
func (d *LibraryDownloader) DownloadAndExtract(downloadURL, filename string) (string, error) {
	// Create target directory for this release
	targetDir := filepath.Join(d.cacheDir, strings.TrimSuffix(filename, ".zip"))

	// Check if already extracted
	if d.isLibraryReady(targetDir) {
		return targetDir, nil
	}

	// Download the archive
	archivePath := filepath.Join(d.cacheDir, filename)
	if err := d.downloadFile(downloadURL, archivePath); err != nil {
		return "", fmt.Errorf("failed to download %s: %w", filename, err)
	}

	// Extract the archive
	if err := d.extractZip(archivePath, targetDir); err != nil {
		return "", fmt.Errorf("failed to extract %s: %w", filename, err)
	}

	// Clean up the archive file
	_ = os.Remove(archivePath)

	return targetDir, nil
}

// downloadFile downloads a file from URL to the specified path
func (d *LibraryDownloader) downloadFile(url, filepath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", d.userAgent)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// extractZip extracts a ZIP archive to the specified directory
func (d *LibraryDownloader) extractZip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer reader.Close()

	// Create destination directory
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Extract files
	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)

		// Security check: ensure the file path is within the destination directory
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.FileInfo().Mode())
			continue
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// Extract file
		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in archive: %w", err)
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			return fmt.Errorf("failed to create target file: %w", err)
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return fmt.Errorf("failed to extract file: %w", err)
		}
	}

	return nil
}

// isLibraryReady checks if the library files are already extracted and ready
func (d *LibraryDownloader) isLibraryReady(dir string) bool {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}

	// Check if we have the main library file
	expectedLib, err := getExpectedLibraryName()
	if err != nil {
		return false
	}

	// Check common paths where the library might be located
	searchPaths := []string{
		filepath.Join(dir, "build", "bin", expectedLib),
		filepath.Join(dir, "bin", expectedLib),
		filepath.Join(dir, expectedLib),
		filepath.Join(dir, "lib", expectedLib),
		filepath.Join(dir, "src", expectedLib),
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// getExpectedLibraryName returns the expected library filename for the current platform
func getExpectedLibraryName() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "libllama.dylib", nil
	case "linux":
		return "libllama.so", nil
	case "windows":
		return "llama.dll", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// FindLibraryPath finds the main library file in the extracted directory
func (d *LibraryDownloader) FindLibraryPath(extractedDir string) (string, error) {
	expectedLib, err := getExpectedLibraryName()
	if err != nil {
		return "", err
	}

	// Common paths where the library might be located
	searchPaths := []string{
		filepath.Join(extractedDir, "build", "bin", expectedLib),
		filepath.Join(extractedDir, "bin", expectedLib),
		filepath.Join(extractedDir, expectedLib),
		filepath.Join(extractedDir, "lib", expectedLib),
		filepath.Join(extractedDir, "src", expectedLib),
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("library file %s not found in %s", expectedLib, extractedDir)
}

// CleanCache removes old cached library files
func (d *LibraryDownloader) CleanCache() error {
	return os.RemoveAll(d.cacheDir)
}
