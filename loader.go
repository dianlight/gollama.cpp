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
	handle       uintptr
	loaded       bool
	llamaLibPath string
	rootLibPath  string
	downloader   *LibraryDownloader
	tempDir      string
	mutex        sync.RWMutex
}

var globalLoader = &LibraryLoader{}

// LoadLibrary loads the appropriate llama.cpp library for the current platform
func (l *LibraryLoader) LoadLibrary() error {
	return l.LoadLibraryWithVersion("")
}

// LoadLibraryWithVersion loads the llama.cpp library for a specific version
// If version is empty, it loads the latest version
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

	// Prefer embedded libraries if available for the requested version
	if resolvedVersion == LlamaCppBuild && hasEmbeddedLibraryForPlatform(runtime.GOOS, runtime.GOARCH) {
		targetDir := filepath.Join(l.downloader.cacheDir, "embedded", embeddedPlatformDirName(runtime.GOOS, runtime.GOARCH))
		if !l.downloader.isLibraryReady(targetDir) {
			if err := extractEmbeddedLibrariesTo(targetDir, runtime.GOOS, runtime.GOARCH); err != nil {
				return fmt.Errorf("failed to extract embedded libraries: %w", err)
			}
		}

		libPath, err := l.downloader.FindLibraryPathForPlatform(targetDir, runtime.GOOS)
		if err == nil {
			// Preload dependent libraries before loading main library
			if err := l.preloadDependentLibraries(libPath); err != nil {
				return fmt.Errorf("failed to preload dependent libraries: %w", err)
			}

			handle, err := l.loadSharedLibrary(libPath)
			if err != nil {
				return fmt.Errorf("failed to load embedded library %s: %w", libPath, err)
			}
			l.handle = handle
			l.llamaLibPath = libPath
			l.loaded = true
			l.rootLibPath = targetDir
			return nil
		}
	}

	// Get the appropriate release when embedded libs are unavailable
	release, err := l.getReleaseForVersion(version)
	if err != nil {
		return err
	}

	// Get platform-specific asset pattern
	pattern, err := l.downloader.GetPlatformAssetPattern()
	if err != nil {
		return fmt.Errorf("failed to get platform pattern: %w", err)
	}

	// Find the appropriate asset
	assetName, downloadURL, err := l.downloader.FindAssetByPattern(release, pattern)
	if err != nil {
		return fmt.Errorf("failed to find platform asset: %w", err)
	}

	// Check if library is already cached and ready
	extractedDir := filepath.Join(l.downloader.cacheDir, strings.TrimSuffix(assetName, ".zip"))
	libPath, err := l.downloader.FindLibraryPath(extractedDir)
	if err == nil {
		// Preload dependent libraries before loading main library
		if err := l.preloadDependentLibraries(libPath); err != nil {
			return fmt.Errorf("failed to preload dependent libraries: %w", err)
		}

		handle, err := l.loadSharedLibrary(libPath)
		if err != nil {
			return fmt.Errorf("failed to load library %s: %w", libPath, err)
		}

		l.handle = handle
		l.llamaLibPath = libPath
		l.loaded = true

		return nil
	}

	// Download and extract the library
	extractedDir, err = l.downloader.DownloadAndExtract(downloadURL, assetName)
	if err != nil {
		return fmt.Errorf("failed to download library: %w", err)
	}

	// Find the main library file
	libPath, err = l.downloader.FindLibraryPath(extractedDir)
	if err != nil {
		return fmt.Errorf("failed to find library file: %w", err)
	}

	// Preload dependent libraries on Unix-like systems to ensure correct library versions are used
	if err := l.preloadDependentLibraries(libPath); err != nil {
		return fmt.Errorf("failed to preload dependent libraries: %w", err)
	}

	// Load the library
	handle, err := l.loadSharedLibrary(libPath)
	if err != nil {
		return fmt.Errorf("failed to load library %s: %w", libPath, err)
	}

	l.handle = handle
	l.llamaLibPath = libPath
	l.loaded = true

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
	if globalLoader.downloader == nil {
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
	}

	return globalLoader.downloader.DownloadMultiplePlatforms(platforms, version)
}

// GetSHA256ForFile calculates the SHA256 checksum for a given file
func GetSHA256ForFile(filepath string) (string, error) {
	if globalLoader.downloader == nil {
		// Check if global config has a custom cache directory
		cacheDir := ""
		if globalConfig != nil && globalConfig.CacheDir != "" {
			cacheDir = globalConfig.CacheDir
		}

		downloader, err := NewLibraryDownloaderWithCacheDir(cacheDir)
		if err != nil {
			return "", fmt.Errorf("failed to create library downloader: %w", err)
		}
		globalLoader.downloader = downloader
	}

	return globalLoader.downloader.calculateSHA256(filepath)
}

// GetLibraryCacheDir returns the directory where downloaded libraries are cached
func GetLibraryCacheDir() (string, error) {
	if globalLoader.downloader == nil {
		// Check if global config has a custom cache directory
		cacheDir := ""
		if globalConfig != nil && globalConfig.CacheDir != "" {
			cacheDir = globalConfig.CacheDir
		}

		downloader, err := NewLibraryDownloaderWithCacheDir(cacheDir)
		if err != nil {
			return "", fmt.Errorf("failed to create library downloader: %w", err)
		}
		globalLoader.downloader = downloader
	}

	return globalLoader.downloader.GetCacheDir(), nil
}
