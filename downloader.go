package gollama

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	llamaCppRepo      = "ggml-org/llama.cpp"
	githubReleasesAPI = "https://api.github.com/repos"
	githubReleasesURL = "https://github.com/ggml-org/llama.cpp/releases/download"
	downloadTimeout   = 10 * time.Minute
	userAgent         = "gollama.cpp/1.0.0"
)

// isValidPath checks if a file path is safe for extraction
func isValidPath(dest, filename string) error {
	// Clean the filename to resolve any .. components
	cleanName := filepath.Clean(filename)

	// Check for absolute paths or paths that start with ..
	if filepath.IsAbs(cleanName) || strings.HasPrefix(cleanName, "..") {
		return fmt.Errorf("unsafe path: %s", filename)
	}

	// Join with destination and check final path
	finalPath := filepath.Join(dest, cleanName)
	cleanDest := filepath.Clean(dest) + string(os.PathSeparator)

	if !strings.HasPrefix(finalPath, cleanDest) {
		return fmt.Errorf("path traversal attempt: %s", filename)
	}

	return nil
}

// ReleaseInfo represents GitHub release information
type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
}

// DownloadTask represents a single download task for parallel processing
type DownloadTask struct {
	Platform     string
	AssetName    string
	DownloadURL  string
	TargetDir    string
	ExpectedSHA2 string
}

// DownloadResult represents the result of a download task
type DownloadResult struct {
	Platform    string
	Success     bool
	Error       error
	LibraryPath string
	SHA256Sum   string
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

	if err := os.MkdirAll(cacheDir, 0750); err != nil {
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

// DownloadAndExtractWithChecksum downloads and extracts the library archive with checksum verification
func (d *LibraryDownloader) DownloadAndExtractWithChecksum(downloadURL, filename, expectedChecksum string) (string, string, error) {
	// Create target directory for this release
	targetDir := filepath.Join(d.cacheDir, strings.TrimSuffix(filename, ".zip"))

	// Check if already extracted
	if d.isLibraryReady(targetDir) {
		// Calculate checksum of existing file if available
		archivePath := filepath.Join(d.cacheDir, filename)
		if _, err := os.Stat(archivePath); err == nil {
			checksum, _ := d.calculateSHA256(archivePath)
			return targetDir, checksum, nil
		}
		return targetDir, "", nil
	}

	// Download the archive with checksum calculation
	archivePath := filepath.Join(d.cacheDir, filename)
	checksum, err := d.downloadFileWithChecksum(downloadURL, archivePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to download %s: %w", filename, err)
	}

	// Verify checksum if provided
	if err := d.verifySHA256(archivePath, expectedChecksum); err != nil {
		// Remove corrupted file
		_ = os.Remove(archivePath)
		return "", "", fmt.Errorf("checksum verification failed for %s: %w", filename, err)
	}

	// Extract the archive
	if err := d.extractZip(archivePath, targetDir); err != nil {
		return "", "", fmt.Errorf("failed to extract %s: %w", filename, err)
	}

	// Clean up the archive file
	_ = os.Remove(archivePath)

	return targetDir, checksum, nil
}

// GetPlatformAssetPatternForPlatform returns the asset name pattern for a specific platform
func (d *LibraryDownloader) GetPlatformAssetPatternForPlatform(goos, goarch string) (string, error) {
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

// DownloadMultiplePlatforms downloads libraries for multiple platforms in parallel
func (d *LibraryDownloader) DownloadMultiplePlatforms(platforms []string, version string) ([]DownloadResult, error) {
	var release *ReleaseInfo
	var err error

	// Get release information
	if version != "" {
		release, err = d.GetReleaseByTag(version)
	} else {
		release, err = d.GetLatestRelease()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get release information: %w", err)
	}

	// Create download tasks
	var tasks []DownloadTask
	for _, platform := range platforms {
		parts := strings.Split(platform, "/")
		if len(parts) != 2 {
			continue // Skip invalid platform specifications
		}
		goos, goarch := parts[0], parts[1]

		pattern, err := d.GetPlatformAssetPatternForPlatform(goos, goarch)
		if err != nil {
			continue // Skip unsupported platforms
		}

		assetName, downloadURL, err := d.FindAssetByPattern(release, pattern)
		if err != nil {
			continue // Skip platforms without available assets
		}

		targetDir := filepath.Join(d.cacheDir, strings.TrimSuffix(assetName, ".zip"))
		tasks = append(tasks, DownloadTask{
			Platform:    platform,
			AssetName:   assetName,
			DownloadURL: downloadURL,
			TargetDir:   targetDir,
			// No expected checksum since llama.cpp doesn't provide them
			ExpectedSHA2: "",
		})
	}

	// Execute downloads in parallel
	return d.executeParallelDownloads(tasks)
}

// executeParallelDownloads executes multiple download tasks concurrently
func (d *LibraryDownloader) executeParallelDownloads(tasks []DownloadTask) ([]DownloadResult, error) {
	results := make([]DownloadResult, len(tasks))
	var wg sync.WaitGroup

	// Use a semaphore to limit concurrent downloads (max 4 concurrent)
	semaphore := make(chan struct{}, 4)

	for i, task := range tasks {
		wg.Add(1)
		go func(index int, t DownloadTask) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := DownloadResult{
				Platform: t.Platform,
				Success:  false,
			}

			// Check if already exists and ready
			if d.isLibraryReady(t.TargetDir) {
				result.Success = true
				// Extract platform info from task
				parts := strings.Split(t.Platform, "/")
				if len(parts) == 2 {
					libPath, err := d.FindLibraryPathForPlatform(t.TargetDir, parts[0])
					if err == nil {
						result.LibraryPath = libPath
					}
				}
				// Try to calculate checksum of existing archive if available
				archivePath := filepath.Join(d.cacheDir, t.AssetName)
				if checksum, err := d.calculateSHA256(archivePath); err == nil {
					result.SHA256Sum = checksum
				}
				results[index] = result
				return
			}

			// Download and extract with checksum
			extractedDir, checksum, err := d.DownloadAndExtractWithChecksum(t.DownloadURL, t.AssetName, t.ExpectedSHA2)
			if err != nil {
				result.Error = err
				results[index] = result
				return
			}

			// Find library path for the specific platform
			parts := strings.Split(t.Platform, "/")
			if len(parts) != 2 {
				result.Error = fmt.Errorf("invalid platform format: %s", t.Platform)
				results[index] = result
				return
			}

			libPath, err := d.FindLibraryPathForPlatform(extractedDir, parts[0])
			if err != nil {
				result.Error = fmt.Errorf("library not found after extraction: %w", err)
				results[index] = result
				return
			}

			result.Success = true
			result.LibraryPath = libPath
			result.SHA256Sum = checksum
			results[index] = result
		}(i, task)
	}

	wg.Wait()
	return results, nil
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

// downloadFileWithChecksum downloads a file and calculates its SHA256 checksum
func (d *LibraryDownloader) downloadFileWithChecksum(url, filepath string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", d.userAgent)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Create a hash writer that computes SHA256 while writing
	hash := sha256.New()
	multiWriter := io.MultiWriter(out, hash)

	_, err = io.Copy(multiWriter, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return the hexadecimal representation of the hash
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// calculateSHA256 calculates the SHA256 checksum of a file
func (d *LibraryDownloader) calculateSHA256(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// verifySHA256 verifies that a file matches the expected SHA256 checksum
func (d *LibraryDownloader) verifySHA256(filepath, expectedChecksum string) error {
	if expectedChecksum == "" {
		// No checksum to verify
		return nil
	}

	actualChecksum, err := d.calculateSHA256(filepath)
	if err != nil {
		return err
	}

	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
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
	if err := os.MkdirAll(dest, 0750); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Extract files
	for _, file := range reader.File {
		// Validate path security
		if err := isValidPath(dest, file.Name); err != nil {
			return err
		}

		// #nosec G305 - Path is validated by isValidPath function above
		path := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
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

		// Limit extraction to prevent decompression bombs (max 1GB per file)
		const maxFileSize = 1 << 30 // 1GB
		limitedReader := io.LimitReader(fileReader, maxFileSize)

		_, err = io.Copy(targetFile, limitedReader)
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

// getExpectedLibraryNameForPlatform returns the expected library filename for a specific platform
func getExpectedLibraryNameForPlatform(goos string) (string, error) {
	switch goos {
	case "darwin":
		return "libllama.dylib", nil
	case "linux":
		return "libllama.so", nil
	case "windows":
		return "llama.dll", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", goos)
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

// FindLibraryPathForPlatform finds the main library file for a specific platform
func (d *LibraryDownloader) FindLibraryPathForPlatform(extractedDir, goos string) (string, error) {
	expectedLib, err := getExpectedLibraryNameForPlatform(goos)
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
