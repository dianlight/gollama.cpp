package gollama

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v68/github"
)

const (
	downloadTimeout = 10 * time.Minute
	userAgent       = "gollama.cpp/1.0.0"
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
// Using go-github's Release type directly
type ReleaseInfo = github.RepositoryRelease

// DownloadTask represents a single download task for parallel processing
type DownloadTask struct {
	Platform     string
	AssetName    string
	DownloadURL  string
	TargetDir    string
	ExpectedSHA2 string
	ResultIndex  int
}

// DownloadResult represents the result of a download task
type DownloadResult struct {
	Platform     string
	Success      bool
	Error        error
	LibraryPath  string
	SHA256Sum    string
	ExtractedDir string
	Embedded     bool
}

// VariantAsset represents a single variant asset for a platform
type VariantAsset struct {
	AssetName   string
	DownloadURL string
	Variant     string // e.g., "cpu", "cuda-12.6.0", "vulkan", "hip-6.2"
}

// VariantDownloadResult represents the result of downloading all variants for a platform
type VariantDownloadResult struct {
	Platform      string
	Variants      []VariantInfo
	Success       bool
	Error         error
	CommonLibPath string // Path to verified common library files
}

// VariantInfo contains information about a downloaded variant
type VariantInfo struct {
	Variant      string
	ExtractedDir string
	SHA256Sum    string
	Success      bool
	Error        error
}

// LibraryDownloader handles downloading pre-built llama.cpp binaries
type LibraryDownloader struct {
	cacheDir  string
	userAgent string
	client    *github.Client
}

// NewLibraryDownloader creates a new library downloader instance
func NewLibraryDownloader() (*LibraryDownloader, error) {
	return NewLibraryDownloaderWithCacheDir("")
}

// NewLibraryDownloaderWithCacheDir creates a new library downloader instance with a custom cache directory
func NewLibraryDownloaderWithCacheDir(customCacheDir string) (*LibraryDownloader, error) {
	var cacheDir string

	// Use custom cache directory if provided
	if customCacheDir != "" {
		cacheDir = customCacheDir
	} else {
		// Check for environment variable
		if envCacheDir := os.Getenv("GOLLAMA_CACHE_DIR"); envCacheDir != "" {
			cacheDir = filepath.Join(envCacheDir, "libs")
		} else {
			// Try user cache directory first
			userCacheDir, err := os.UserCacheDir()
			if err == nil {
				cacheDir = filepath.Join(userCacheDir, "gollama", "libs")
			} else {
				// Fallback to temp directory
				cacheDir = filepath.Join(os.TempDir(), "gollama", "libs")
			}
		}
	}

	if err := os.MkdirAll(cacheDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Create go-github client with optional authentication
	var client *github.Client
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		// Authenticated client using GITHUB_TOKEN
		httpClient := &http.Client{Timeout: downloadTimeout}
		client = github.NewClient(httpClient).WithAuthToken(token)
	} else {
		// Unauthenticated client
		httpClient := &http.Client{Timeout: downloadTimeout}
		client = github.NewClient(httpClient)
	}

	return &LibraryDownloader{
		cacheDir:  cacheDir,
		userAgent: userAgent,
		client:    client,
	}, nil
}

// GetLatestRelease fetches the latest release information from GitHub
func (d *LibraryDownloader) GetLatestRelease() (*ReleaseInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	release, _, err := d.client.Repositories.GetLatestRelease(ctx, "ggml-org", "llama.cpp")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}

	return release, nil
}

// GetReleaseByTag fetches release information for a specific tag
func (d *LibraryDownloader) GetReleaseByTag(tag string) (*ReleaseInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	release, _, err := d.client.Repositories.GetReleaseByTag(ctx, "ggml-org", "llama.cpp", tag)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w (%s)", err, tag)
	}

	return release, nil
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
		// Auto-detect available GPU backends
		return d.getLinuxVariantPattern(arch), nil
	case "windows":
		// Auto-detect available GPU backends
		return d.getWindowsVariantPattern(arch), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", goos)
	}
}

// getLinuxVariantPattern detects and returns the best GPU variant pattern for Linux
func (d *LibraryDownloader) getLinuxVariantPattern(arch string) string {
	// Priority order: CUDA > HIP > Vulkan > SYCL > CPU

	// Check for CUDA
	if d.hasCommand("nvcc") {
		// Try CUDA variant first
		return fmt.Sprintf("llama-.*-bin-ubuntu-cuda-.*-%s.zip", arch)
	}

	// Check for HIP/ROCm
	if d.hasCommand("hipconfig") {
		return fmt.Sprintf("llama-.*-bin-ubuntu-hip-.*-%s.zip", arch)
	}

	// Check for Vulkan
	if d.hasCommand("vulkaninfo") {
		return fmt.Sprintf("llama-.*-bin-ubuntu-vulkan-%s.zip", arch)
	}

	// Check for SYCL (Intel oneAPI)
	if d.hasCommand("sycl-ls") {
		return fmt.Sprintf("llama-.*-bin-ubuntu-sycl-%s.zip", arch)
	}

	// Fallback to CPU
	return fmt.Sprintf("llama-.*-bin-ubuntu-%s.zip", arch)
}

// getWindowsVariantPattern detects and returns the best GPU variant pattern for Windows
func (d *LibraryDownloader) getWindowsVariantPattern(arch string) string {
	// Priority order: CUDA > HIP > Vulkan > OpenCL > SYCL > CPU

	// Check for CUDA
	if d.hasCommand("nvcc") {
		return fmt.Sprintf("llama-.*-bin-win-cuda-.*-%s.zip", arch)
	}

	// Check for HIP (Windows)
	if d.hasCommand("hipconfig") {
		return fmt.Sprintf("llama-.*-bin-win-hip-.*-%s.zip", arch)
	}

	// Check for Vulkan
	if d.hasCommand("vulkaninfo") {
		return fmt.Sprintf("llama-.*-bin-win-vulkan-%s.zip", arch)
	}

	// Check for OpenCL (especially for ARM64/Adreno)
	if d.hasCommand("clinfo") || arch == "arm64" {
		return fmt.Sprintf("llama-.*-bin-win-opencl-.*-%s.zip", arch)
	}

	// Check for SYCL (Intel oneAPI)
	if d.hasCommand("sycl-ls") {
		return fmt.Sprintf("llama-.*-bin-win-sycl-%s.zip", arch)
	}

	// Fallback to CPU
	return fmt.Sprintf("llama-.*-bin-win-cpu-%s.zip", arch)
}

// hasCommand checks if a command is available in PATH
func (d *LibraryDownloader) hasCommand(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// FindAssetByPattern finds an asset that matches the given pattern
func (d *LibraryDownloader) FindAssetByPattern(release *ReleaseInfo, pattern string) (string, string, error) {
	// Compile the pattern as a regular expression
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return "", "", fmt.Errorf("invalid pattern: %w", err)
	}

	for _, asset := range release.Assets {
		// go-github returns pointers to strings, so we need to dereference them
		assetName := ""
		downloadURL := ""
		if asset.Name != nil {
			assetName = *asset.Name
		}
		if asset.BrowserDownloadURL != nil {
			downloadURL = *asset.BrowserDownloadURL
		}

		if regex.MatchString(assetName) {
			return assetName, downloadURL, nil
		}
	}
	return "", "", fmt.Errorf("no asset found matching pattern: %s", pattern)
}

// FindAllVariantAssets finds all variant assets for a specific platform and architecture
// Pattern: llama-<version>-bin-<os>-<variant>[-<variant-version>]-<arch>.zip
func (d *LibraryDownloader) FindAllVariantAssets(release *ReleaseInfo, goos, goarch string) ([]VariantAsset, error) {
	// Convert Go arch to llama.cpp naming convention
	var arch string
	switch goarch {
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "arm64"
	default:
		return nil, fmt.Errorf("unsupported architecture: %s", goarch)
	}

	// Build platform-specific base pattern
	var osPrefix string
	switch goos {
	case "darwin":
		osPrefix = "macos"
	case "linux":
		osPrefix = "ubuntu"
	case "windows":
		osPrefix = "win"
	default:
		return nil, fmt.Errorf("unsupported OS: %s", goos)
	}

	// Pattern to match all variants: llama-<version>-bin-<os>-<variant>[-<variant-version>]-<arch>.zip
	// Examples:
	//   llama-b1234-bin-ubuntu-x64.zip (CPU)
	//   llama-b1234-bin-ubuntu-cuda-12.6.0-x64.zip
	//   llama-b1234-bin-ubuntu-vulkan-x64.zip
	//   llama-b1234-bin-macos-arm64.zip
	// The pattern captures everything between the OS prefix and the arch
	patternStr := fmt.Sprintf(`^llama-[^-]+-bin-%s-(.+-)%s\.zip$`, osPrefix, arch)
	cpuPatternStr := fmt.Sprintf(`^llama-[^-]+-bin-%s-%s\.zip$`, osPrefix, arch)

	regex, err := regexp.Compile(patternStr)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	cpuRegex, err := regexp.Compile(cpuPatternStr)
	if err != nil {
		return nil, fmt.Errorf("invalid CPU pattern: %w", err)
	}

	var variants []VariantAsset
	for _, asset := range release.Assets {
		assetName := ""
		downloadURL := ""
		if asset.Name != nil {
			assetName = *asset.Name
		}
		if asset.BrowserDownloadURL != nil {
			downloadURL = *asset.BrowserDownloadURL
		}

		// Check CPU-only variant first (simplest pattern)
		if cpuRegex.MatchString(assetName) {
			variants = append(variants, VariantAsset{
				AssetName:   assetName,
				DownloadURL: downloadURL,
				Variant:     "cpu",
			})
			continue
		}

		// Check for GPU variants
		if matches := regex.FindStringSubmatch(assetName); matches != nil {
			// Extract variant string (everything between os prefix and arch)
			// e.g., "cuda-12.6.0-" -> "cuda-12.6.0"
			variantStr := strings.TrimSuffix(matches[1], "-")

			variants = append(variants, VariantAsset{
				AssetName:   assetName,
				DownloadURL: downloadURL,
				Variant:     variantStr,
			})
		}
	}

	if len(variants) == 0 {
		return nil, fmt.Errorf("no variants found for %s/%s", goos, goarch)
	}

	return variants, nil
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
		// Auto-detect available GPU backends for Linux
		return d.getLinuxVariantPattern(arch), nil
	case "windows":
		// Auto-detect available GPU backends for Windows
		return d.getWindowsVariantPattern(arch), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", goos)
	}
}

// DownloadMultiplePlatforms downloads libraries for multiple platforms in parallel
func (d *LibraryDownloader) DownloadMultiplePlatforms(platforms []string, version string) ([]DownloadResult, error) {
	preferEmbedded := version == "" || version == LlamaCppBuild
	effectiveVersion := version
	if effectiveVersion == "" {
		effectiveVersion = LlamaCppBuild
	}

	results := make([]DownloadResult, 0, len(platforms))
	var tasks []DownloadTask
	var release *ReleaseInfo

	fetchRelease := func() error {
		if release != nil {
			return nil
		}
		var err error
		if version != "" {
			release, err = d.GetReleaseByTag(version)
		} else {
			release, err = d.GetLatestRelease()
			if err == nil && release != nil && release.TagName != nil {
				effectiveVersion = release.GetTagName()
			}
		}
		if err != nil {
			return fmt.Errorf("failed to get release information: %w (%s)", err, version)
		}
		return nil
	}

	for _, platform := range platforms {
		parts := strings.Split(platform, "/")
		if len(parts) != 2 {
			continue
		}
		goos, goarch := parts[0], parts[1]

		// Attempt to use embedded libraries when allowed and available.
		if preferEmbedded && effectiveVersion == LlamaCppBuild && hasEmbeddedLibraryForPlatform(goos, goarch) {
			targetDir := filepath.Join(d.cacheDir, "embedded", embeddedPlatformDirName(goos, goarch))
			if !d.isLibraryReady(targetDir) {
				if err := extractEmbeddedLibrariesTo(targetDir, goos, goarch); err == nil {
					// extracted successfully
				} else {
					targetDir = ""
				}
			}

			if targetDir != "" {
				if libPath, err := d.FindLibraryPathForPlatform(targetDir, goos); err == nil {
					results = append(results, DownloadResult{
						Platform:     platform,
						Success:      true,
						LibraryPath:  libPath,
						ExtractedDir: targetDir,
						Embedded:     true,
					})
					continue
				}
			}
		}

		if err := fetchRelease(); err != nil {
			return nil, err
		}

		pattern, err := d.GetPlatformAssetPatternForPlatform(goos, goarch)
		if err != nil {
			results = append(results, DownloadResult{
				Platform: platform,
				Success:  false,
				Error:    err,
			})
			continue
		}

		assetName, downloadURL, err := d.FindAssetByPattern(release, pattern)
		if err != nil {
			results = append(results, DownloadResult{
				Platform: platform,
				Success:  false,
				Error:    err,
			})
			continue
		}

		targetDir := filepath.Join(d.cacheDir, strings.TrimSuffix(assetName, ".zip"))
		idx := len(results)
		results = append(results, DownloadResult{Platform: platform})
		tasks = append(tasks, DownloadTask{
			Platform:     platform,
			AssetName:    assetName,
			DownloadURL:  downloadURL,
			TargetDir:    targetDir,
			ExpectedSHA2: "",
			ResultIndex:  idx,
		})
	}

	if len(tasks) == 0 {
		return results, nil
	}

	taskResults, err := d.executeParallelDownloads(tasks)
	if err != nil {
		return nil, err
	}

	for i, taskResult := range taskResults {
		idx := tasks[i].ResultIndex
		if idx < len(results) {
			results[idx] = taskResult
		} else {
			results = append(results, taskResult)
		}
	}

	return results, nil
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
				Platform:     t.Platform,
				Success:      false,
				ExtractedDir: t.TargetDir,
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
				result.ExtractedDir = t.TargetDir
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
			result.ExtractedDir = extractedDir
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

	// Use the HTTP client from go-github
	httpClient := &http.Client{Timeout: downloadTimeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer func() {
		_ = resp.Body.Close() // Ignore error in defer
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		_ = out.Close() // Ignore error in defer
	}()

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

	// Use a fresh HTTP client for file downloads
	httpClient := &http.Client{Timeout: downloadTimeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer func() {
		_ = resp.Body.Close() // Ignore error in defer
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		_ = out.Close() // Ignore error in defer
	}()

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
	defer func() {
		_ = file.Close() // Ignore error in defer
	}()

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
	defer func() {
		_ = reader.Close() // Ignore error in defer
	}()

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
		defer func(fr io.ReadCloser) {
			_ = fr.Close() // Ignore error in defer
		}(fileReader)

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			return fmt.Errorf("failed to create target file: %w", err)
		}
		defer func(tf *os.File) {
			_ = tf.Close() // Ignore error in defer
		}(targetFile)

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

// DownloadAllVariants downloads all variants for a platform and verifies common files are identical
func (d *LibraryDownloader) DownloadAllVariants(release *ReleaseInfo, goos, goarch string) (*VariantDownloadResult, error) {
	result := &VariantDownloadResult{
		Platform: fmt.Sprintf("%s/%s", goos, goarch),
	}

	// Find all variant assets for this platform
	variants, err := d.FindAllVariantAssets(release, goos, goarch)
	if err != nil {
		result.Error = err
		return result, err
	}

	// Download and extract all variants in parallel
	result.Variants = make([]VariantInfo, len(variants))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 4) // Limit concurrent downloads

	for i, variant := range variants {
		wg.Add(1)
		go func(index int, v VariantAsset) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			variantInfo := VariantInfo{
				Variant: v.Variant,
			}

			// Create target directory for this variant
			targetDir := filepath.Join(d.cacheDir, strings.TrimSuffix(v.AssetName, ".zip"))

			// Check if already extracted
			if d.isLibraryReady(targetDir) {
				variantInfo.Success = true
				variantInfo.ExtractedDir = targetDir

				// Try to calculate checksum if archive still exists
				archivePath := filepath.Join(d.cacheDir, v.AssetName)
				if checksum, err := d.calculateSHA256(archivePath); err == nil {
					variantInfo.SHA256Sum = checksum
				}
				result.Variants[index] = variantInfo
				return
			}

			// Download and extract with checksum
			extractedDir, checksum, err := d.DownloadAndExtractWithChecksum(v.DownloadURL, v.AssetName, "")
			if err != nil {
				variantInfo.Error = err
				variantInfo.Success = false
				result.Variants[index] = variantInfo
				return
			}

			variantInfo.Success = true
			variantInfo.ExtractedDir = extractedDir
			variantInfo.SHA256Sum = checksum
			result.Variants[index] = variantInfo
		}(i, variant)
	}

	wg.Wait()

	// Check if all downloads succeeded
	allSuccess := true
	for _, v := range result.Variants {
		if !v.Success {
			allSuccess = false
			if result.Error == nil && v.Error != nil {
				result.Error = fmt.Errorf("variant %s failed: %w", v.Variant, v.Error)
			}
		}
	}

	if !allSuccess {
		result.Success = false
		return result, result.Error
	}

	// Verify common files are identical across all variants
	if err := d.verifyCommonFiles(result.Variants, goos); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("common file verification failed: %w", err)
		return result, result.Error
	}

	// All variants downloaded successfully and common files verified
	result.Success = true

	// Use the first variant's directory as the common lib path
	if len(result.Variants) > 0 && result.Variants[0].ExtractedDir != "" {
		result.CommonLibPath = result.Variants[0].ExtractedDir
	}

	return result, nil
}

// verifyCommonFiles checks that common files are identical across all variant directories
func (d *LibraryDownloader) verifyCommonFiles(variants []VariantInfo, goos string) error {
	if len(variants) < 2 {
		// Nothing to compare
		return nil
	}

	baseDir := variants[0].ExtractedDir
	if baseDir == "" {
		return fmt.Errorf("base variant has no extracted directory")
	}

	// Get expected library name for platform
	expectedLib, err := getExpectedLibraryNameForPlatform(goos)
	if err != nil {
		return err
	}

	// Common files to check (relative paths within extracted directory)
	commonFiles := []string{
		filepath.Join("build", "bin", expectedLib),
		filepath.Join("bin", expectedLib),
		expectedLib,
	}

	// Find which common files actually exist in base directory
	var existingCommonFiles []string
	for _, relPath := range commonFiles {
		fullPath := filepath.Join(baseDir, relPath)
		if _, err := os.Stat(fullPath); err == nil {
			existingCommonFiles = append(existingCommonFiles, relPath)
			break // Only check the first existing library file
		}
	}

	if len(existingCommonFiles) == 0 {
		// No common library file found, this is expected behavior
		// Different variants may have different structures
		return nil
	}

	// Calculate checksums of base files
	baseChecksums := make(map[string]string)
	for _, relPath := range existingCommonFiles {
		fullPath := filepath.Join(baseDir, relPath)
		checksum, err := d.calculateSHA256(fullPath)
		if err != nil {
			return fmt.Errorf("failed to calculate checksum for %s: %w", relPath, err)
		}
		baseChecksums[relPath] = checksum
	}

	// Compare with other variants
	for i := 1; i < len(variants); i++ {
		variantDir := variants[i].ExtractedDir
		if variantDir == "" {
			continue
		}

		for relPath, baseChecksum := range baseChecksums {
			variantPath := filepath.Join(variantDir, relPath)

			// Check if file exists in variant
			if _, err := os.Stat(variantPath); os.IsNotExist(err) {
				// File doesn't exist in this variant - this is OK as different variants
				// may have different file structures
				continue
			}

			// Calculate checksum
			variantChecksum, err := d.calculateSHA256(variantPath)
			if err != nil {
				return fmt.Errorf("failed to calculate checksum for %s in variant %s: %w",
					relPath, variants[i].Variant, err)
			}

			// Compare checksums
			if baseChecksum != variantChecksum {
				return fmt.Errorf("file %s differs between variant %s and %s",
					relPath, variants[0].Variant, variants[i].Variant)
			}
		}
	}

	return nil
}

// CleanCache removes old cached library files
func (d *LibraryDownloader) CleanCache() error {
	return os.RemoveAll(d.cacheDir)
}

// GetCacheDir returns the cache directory being used
func (d *LibraryDownloader) GetCacheDir() string {
	return d.cacheDir
}
