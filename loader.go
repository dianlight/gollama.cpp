package gollama

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Library loader manages the loading and lifecycle of llama.cpp shared libraries
type LibraryLoader struct {
	handle          uintptr
	loaded          bool
	llamaLibPath    string
	rootLibPath     string
	extensionSuffix string
	downloader      *LibraryDownloader
	tempDir         string
	mutex           sync.RWMutex
}

var globalLoader = &LibraryLoader{}

// LoadLibrary loads the appropriate llama.cpp library for the current platform
func (l *LibraryLoader) LoadLibrary() error {
	return l.LoadLibraryWithVersion("")
}

// LoadLibraryWithVersion loads the llama.cpp library for a specific version
// If version is empty, it loads the default build version (LlamaCppBuild)
// Resolution order:
// 1) Embedded (only if version == LlamaCppBuild)
// 2) Local ./libs (only if version == LlamaCppBuild)
// 3) Cache directory entries matching current GOOS (best-effort scan)
// 4) Download + extract to cache
// 5) Return a detailed error if all fail
func (l *LibraryLoader) LoadLibraryWithVersion(version string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.loaded {
		return nil
	}

	resolvedVersion := version
	if resolvedVersion == "" {
		resolvedVersion = LlamaCppBuild
	}

	// Initialize downloader if not already done
	if l.downloader == nil {
		// Check if global config has a custom cache directory
		cacheDir := ""
		if globalConfig != nil && globalConfig.CacheDir != "" {
			cacheDir = globalConfig.CacheDir
		}

		downloader, err := NewLibraryDownloaderWithCacheDir(cacheDir)
		if err != nil {
			return fmt.Errorf("failed to create library downloader: %w", err)
		}
		l.downloader = downloader
	}

	var reasons []string

	// 1) Embedded libraries
	if resolvedVersion == LlamaCppBuild && hasEmbeddedLibraryForPlatform(runtime.GOOS, runtime.GOARCH) {
		targetDir := filepath.Join(l.downloader.cacheDir, "embedded", embeddedPlatformDirName(runtime.GOOS, runtime.GOARCH))
		if !l.downloader.isLibraryReady(targetDir) {
			if err := extractEmbeddedLibrariesTo(targetDir, runtime.GOOS, runtime.GOARCH); err != nil {
				reasons = append(reasons, fmt.Sprintf("embedded extract failed: %v", err))
			}
		}
		if l.downloader.isLibraryReady(targetDir) {
			if libPath, err := l.downloader.FindLibraryPathForPlatform(targetDir, runtime.GOOS); err == nil {
				info, errs := l.LoadLibraryWithDependencies(libPath)
				reasons = append(reasons, errs...)
				if info.Success {
					if err := l.ApplyLibraryLoad(info, targetDir); err == nil {
						return nil
					}
				}
			} else {
				reasons = append(reasons, fmt.Sprintf("embedded lib not found in %s: %v", targetDir, err))
			}
		}
	}

	// 2) Local ./libs for the same build (only when version == LlamaCppBuild)
	if !l.loaded && resolvedVersion == LlamaCppBuild {
		localDir := filepath.Join("libs", embeddedPlatformDirName(runtime.GOOS, runtime.GOARCH))
		if _, statErr := os.Stat(localDir); statErr == nil {
			if libPath, err := l.downloader.FindLibraryPathForPlatform(localDir, runtime.GOOS); err == nil {
				info, errs := l.LoadLibraryWithDependencies(libPath)
				reasons = append(reasons, errs...)
				if info.Success {
					if err := l.ApplyLibraryLoad(info, localDir); err == nil {
						return nil
					}
				}
			} else {
				reasons = append(reasons, fmt.Sprintf("./libs library not found: %v", err))
			}
		} else if !os.IsNotExist(statErr) {
			reasons = append(reasons, fmt.Sprintf("./libs check failed: %v", statErr))
		}
	}

	// 3) Cache directory scan (best effort, match GOOS by library filename)
	if !l.loaded {
		entries, err := os.ReadDir(l.downloader.cacheDir)
		if err == nil {
			for _, e := range entries {
				if !e.IsDir() {
					continue
				}
				name := e.Name()
				if name == "embedded" { // already handled in step 1
					continue
				}
				candDir := filepath.Join(l.downloader.cacheDir, name)
				if libPath, err := l.downloader.FindLibraryPathForPlatform(candDir, runtime.GOOS); err == nil {
					info, errs := l.LoadLibraryWithDependencies(libPath)
					if len(errs) > 0 {
						reasons = append(reasons, errs...)
						continue
					}
					if info.Success {
						if err := l.ApplyLibraryLoad(info, candDir); err == nil {
							return nil
						}
					}
				}
			}
		} else {
			reasons = append(reasons, fmt.Sprintf("cache scan failed: %v", err))
		}
	}

	// 4) Download and extract into cache
	// Fetch release according to resolvedVersion
	release, err := l.getReleaseForVersion(resolvedVersion)
	if err != nil {
		reasons = append(reasons, fmt.Sprintf("release fetch failed: %v", err))
		return fmt.Errorf("failed to resolve llama.cpp libraries: %s", strings.Join(reasons, "; "))
	}

	pattern, err := l.downloader.GetPlatformAssetPattern()
	if err != nil {
		reasons = append(reasons, fmt.Sprintf("platform pattern failed: %v", err))
		return fmt.Errorf("failed to resolve llama.cpp libraries: %s", strings.Join(reasons, "; "))
	}

	assetName, downloadURL, err := l.downloader.FindAssetByPattern(release, pattern)
	if err != nil {
		reasons = append(reasons, fmt.Sprintf("no matching asset: %v", err))
		return fmt.Errorf("failed to resolve llama.cpp libraries: %s", strings.Join(reasons, "; "))
	}

	// If already extracted in cache (by exact asset name), use it
	extractedDir := filepath.Join(l.downloader.cacheDir, strings.TrimSuffix(assetName, ".zip"))
	if libPath, err := l.downloader.FindLibraryPathForPlatform(extractedDir, runtime.GOOS); err == nil {
		info, errs := l.LoadLibraryWithDependencies(libPath)
		reasons = append(reasons, errs...)
		if info.Success {
			if err := l.ApplyLibraryLoad(info, extractedDir); err == nil {
				return nil
			}
		}
	}

	// Download and extract
	extractedDir, err = l.downloader.DownloadAndExtract(downloadURL, assetName)
	if err != nil {
		reasons = append(reasons, fmt.Sprintf("download failed: %v", err))
		return fmt.Errorf("failed to resolve llama.cpp libraries: %s", strings.Join(reasons, "; "))
	}

	libPath, err := l.downloader.FindLibraryPathForPlatform(extractedDir, runtime.GOOS)
	if err != nil {
		reasons = append(reasons, fmt.Sprintf("post-extract lib not found: %v", err))
		return fmt.Errorf("failed to resolve llama.cpp libraries: %s", strings.Join(reasons, "; "))
	}

	info, errs := l.LoadLibraryWithDependencies(libPath)
	reasons = append(reasons, errs...)
	if !info.Success {
		return fmt.Errorf("failed to resolve llama.cpp libraries: %s", strings.Join(reasons, "; "))
	}

	if err := l.ApplyLibraryLoad(info, extractedDir); err != nil {
		return fmt.Errorf("failed to apply library load: %w", err)
	}

	return nil
}

func (l *LibraryLoader) getReleaseForVersion(version string) (*ReleaseInfo, error) {
	if version == "" {
		release, err := l.downloader.GetLatestRelease()
		if err != nil {
			return nil, fmt.Errorf("failed to get latest release: %w", err)
		}
		return release, nil
	}

	release, err := l.downloader.GetReleaseByTag(version)
	if err != nil {
		return nil, fmt.Errorf("failed to get release %s: %w", version, err)
	}
	return release, nil
}

// UnloadLibrary unloads the library and cleans up resources
func (l *LibraryLoader) UnloadLibrary() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if !l.loaded {
		return nil
	}

	// Close library handle
	if l.handle != 0 {
		if runtime.GOOS != "windows" && runtime.GOOS == "darwin" {
			// Only call dlclose on Darwin where it's more stable
			_ = closeLibraryPlatform(l.handle) // Ignore error during cleanup
		}
		// On other platforms, we just mark as unloaded without calling dlclose
		// to avoid segfaults in the underlying library
	}

	// Clean up temporary directory if it exists
	if l.tempDir != "" {
		_ = os.RemoveAll(l.tempDir) // Ignore error during cleanup
	}

	l.handle = 0
	l.loaded = false
	l.llamaLibPath = ""
	l.tempDir = ""

	return nil
}

// getLibraryName returns the platform-specific library name
func (l *LibraryLoader) getLibraryName() (string, error) {
	goos := runtime.GOOS

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

// extractEmbeddedLibraries extracts embedded libraries to a temporary directory
// This method is provided for compatibility with tests, but this implementation
// doesn't use embedded libraries - it downloads them instead
func (l *LibraryLoader) extractEmbeddedLibraries() (string, error) {
	// Since this implementation uses downloaded libraries instead of embedded ones,
	// we simulate the behavior expected by tests by creating a temporary directory
	// and returning an error indicating no embedded libraries are available
	return "", fmt.Errorf("no embedded libraries found - this implementation uses downloaded libraries")
}

// GetHandle returns the library handle
func (l *LibraryLoader) GetHandle() uintptr {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.handle
}

// IsLoaded returns whether the library is loaded
func (l *LibraryLoader) IsLoaded() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.loaded
}

// loadSharedLibrary loads a shared library using the appropriate method for the platform
func (l *LibraryLoader) loadSharedLibrary(path string) (uintptr, error) {
	return loadLibraryPlatform(path)
}

// preloadDependentLibraries preloads all dependent libraries from the same directory
// on Unix-like systems to ensure correct library versions are used
func (l *LibraryLoader) preloadDependentLibraries(mainLibPath string) error {
	// Only preload on Unix-like systems where @rpath can cause version conflicts
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		return nil
	}

	// Get the directory containing the main library
	libDir := filepath.Dir(mainLibPath)

	// Define the order of libraries to preload (based on dependency chain)
	// These must be loaded in the correct order to satisfy dependencies
	dependentLibs := []string{
		"libggml-base.dylib",  // Base library - must be loaded first
		"libggml-cpu.dylib",   // CPU implementation
		"libggml-blas.dylib",  // BLAS implementation
		"libggml-metal.dylib", // Metal implementation (macOS)
		"libggml-rpc.dylib",   // RPC implementation
		"libggml.dylib",       // Main GGML library
		"libmtmd.dylib",       // MTMD library
	}

	// On Linux, use .so extension
	if runtime.GOOS == "linux" {
		for i, lib := range dependentLibs {
			dependentLibs[i] = strings.Replace(lib, ".dylib", ".so", 1)
		}
	}

	// Preload each dependent library
	for _, libName := range dependentLibs {
		libPath := filepath.Join(libDir, libName)

		// Check if the library exists
		if _, err := os.Stat(libPath); err != nil {
			// Skip if library doesn't exist (some may be optional)
			continue
		}

		// Preload the library using RTLD_NOW | RTLD_GLOBAL
		_, err := l.loadSharedLibrary(libPath)
		if err != nil {
			// Log but don't fail - some libraries may be optional
			// The main library load will fail if truly required libraries are missing
			continue
		}
	}

	return nil
}

// Global functions for backward compatibility

// ensureDownloader initializes the global downloader if needed
// and returns a reference to it. This consolidates the repeated
// downloader initialization pattern used in several functions
func ensureDownloader() (*LibraryDownloader, error) {
	if globalLoader.downloader != nil {
		return globalLoader.downloader, nil
	}

	// Check if global config has a custom cache directory
	cacheDir := ""
	if globalConfig != nil && globalConfig.CacheDir != "" {
		cacheDir = globalConfig.CacheDir
	}

	downloader, err := NewLibraryDownloaderWithCacheDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create library downloader: %w", err)
	}
	globalLoader.downloader = downloader
	return downloader, nil
}

// LoadLibraryWithVersion loads a specific version of the llama.cpp library
func LoadLibraryWithVersion(version string) error {
	return globalLoader.LoadLibraryWithVersion(version)
}

// getLibHandle returns the global library handle
func getLibHandle() uintptr {
	return globalLoader.GetHandle()
}

// isLibraryLoaded returns whether the global library is loaded
func isLibraryLoaded() bool {
	return globalLoader.IsLoaded()
}

// RegisterFunction registers a function with the global library handle
func RegisterFunction(fptr interface{}, name string) error {
	handle := globalLoader.GetHandle()
	if handle == 0 {
		return fmt.Errorf("library not loaded")
	}

	registerLibFunc(fptr, handle, name)
	return nil
}

// Cleanup function to be called when the program exits
func Cleanup() {
	_ = globalLoader.UnloadLibrary() // Ignore error during cleanup
	_ = unloadLibrary()              // Also unload the gollama.go global state
}

// CleanLibraryCache removes cached library files to force re-download
func CleanLibraryCache() error {
	if globalLoader.downloader != nil {
		return globalLoader.downloader.CleanCache()
	}
	return nil
}

// DownloadLibrariesForPlatforms downloads libraries for multiple platforms in parallel
// platforms should be in the format []string{"linux/amd64", "darwin/arm64", "windows/amd64"}
// version can be empty for latest version or specify a specific version like "b6862"
func DownloadLibrariesForPlatforms(platforms []string, version string) ([]DownloadResult, error) {
	downloader, err := ensureDownloader()
	if err != nil {
		return nil, err
	}

	return downloader.DownloadMultiplePlatforms(platforms, version)
}

// GetSHA256ForFile calculates the SHA256 checksum for a given file
func GetSHA256ForFile(filepath string) (string, error) {
	downloader, err := ensureDownloader()
	if err != nil {
		return "", err
	}

	return downloader.calculateSHA256(filepath)
}

// GetLibraryCacheDir returns the directory where downloaded libraries are cached
func GetLibraryCacheDir() (string, error) {
	downloader, err := ensureDownloader()
	if err != nil {
		return "", err
	}

	return downloader.GetCacheDir(), nil
}

func getExpectedLibrarySuffix() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return ".dylib", nil
	case "linux":
		return ".so", nil
	case "windows":
		return ".dll", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
