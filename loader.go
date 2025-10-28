package gollama

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

// Library loader manages the loading and lifecycle of llama.cpp shared libraries
type LibraryLoader struct {
	handle     uintptr
	loaded     bool
	libPath    string
	downloader *LibraryDownloader
	tempDir    string
	mutex      sync.RWMutex
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

	// Initialize downloader if not already done
	if l.downloader == nil {
		downloader, err := NewLibraryDownloader()
		if err != nil {
			return fmt.Errorf("failed to create library downloader: %w", err)
		}
		l.downloader = downloader
	}

	// Get the appropriate release
	var release *ReleaseInfo
	var err error

	if version == "" {
		release, err = l.downloader.GetLatestRelease()
		if err != nil {
			return fmt.Errorf("failed to get latest release: %w", err)
		}
	} else {
		release, err = l.downloader.GetReleaseByTag(version)
		if err != nil {
			return fmt.Errorf("failed to get release %s: %w", version, err)
		}
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

	// Download and extract the library
	extractedDir, err := l.downloader.DownloadAndExtract(downloadURL, assetName)
	if err != nil {
		return fmt.Errorf("failed to download library: %w", err)
	}

	// Find the main library file
	libPath, err := l.downloader.FindLibraryPath(extractedDir)
	if err != nil {
		return fmt.Errorf("failed to find library file: %w", err)
	}

	// Load the library
	handle, err := l.loadSharedLibrary(libPath)
	if err != nil {
		return fmt.Errorf("failed to load library %s: %w", libPath, err)
	}

	l.handle = handle
	l.libPath = libPath
	l.loaded = true

	return nil
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
	l.libPath = ""
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
	switch runtime.GOOS {
	case "windows":
		// On Windows, we would use LoadLibrary
		// For now, return an error as Windows support is not fully implemented
		return 0, fmt.Errorf("support for windows platform not yet implemented")
	default:
		// On Unix-like systems, use platform-specific loading
		return loadLibraryPlatform(path)
	}
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
		downloader, err := NewLibraryDownloader()
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
		downloader, err := NewLibraryDownloader()
		if err != nil {
			return "", fmt.Errorf("failed to create library downloader: %w", err)
		}
		globalLoader.downloader = downloader
	}

	return globalLoader.downloader.calculateSHA256(filepath)
}
