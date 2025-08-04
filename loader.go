package gollama

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// Embedded libraries - in a real implementation, you would embed the pre-built libraries
// For now, we'll use an empty embed so the build doesn't fail
//
//go:embed libs
var embeddedLibs embed.FS

// Library loader manages the loading and lifecycle of llama.cpp shared libraries
type LibraryLoader struct {
	handle  uintptr
	loaded  bool
	tempDir string
	libPath string
	mutex   sync.RWMutex
}

var globalLoader = &LibraryLoader{}

// LoadLibrary loads the appropriate llama.cpp library for the current platform
func (l *LibraryLoader) LoadLibrary() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.loaded {
		return nil
	}

	// Get platform-specific library name
	libName, err := l.getLibraryName()
	if err != nil {
		return fmt.Errorf("failed to get library name: %w", err)
	}

	// Try to load from embedded files first, then fallback to system paths
	libPath, err := l.extractEmbeddedLibraries()
	if err != nil {
		// Fallback to system library
		libPath = libName
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

// UnloadLibrary unloads the library and cleans up temporary files
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

	// Clean up temporary files
	if l.tempDir != "" {
		_ = os.RemoveAll(l.tempDir) // Ignore error during cleanup
	}

	l.handle = 0
	l.loaded = false
	l.tempDir = ""
	l.libPath = ""

	return nil
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

// extractEmbeddedLibraries extracts all embedded libraries for the current platform to a temporary location
func (l *LibraryLoader) extractEmbeddedLibraries() (string, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Define file extensions for each platform
	var libExtension string
	switch goos {
	case "darwin":
		libExtension = ".dylib"
	case "linux":
		libExtension = ".so"
	case "windows":
		libExtension = ".dll"
	default:
		return "", fmt.Errorf("unsupported OS: %s", goos)
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "gollama-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Extract all libraries for this platform
	embeddedDir := fmt.Sprintf("libs/%s_%s", goos, goarch)
	extractedCount := 0

	// Read the embedded directory to find all library files
	dirEntries, err := embeddedLibs.ReadDir(embeddedDir)
	if err != nil {
		_ = os.RemoveAll(tempDir) // Ignore error during cleanup
		return "", fmt.Errorf("no embedded libraries directory found for platform %s_%s: %w", goos, goarch, err)
	}

	// Extract all files with the appropriate extension
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if filepath.Ext(fileName) != libExtension {
			continue
		}

		embeddedPath := fmt.Sprintf("%s/%s", embeddedDir, fileName)

		// Try to read the embedded file
		data, err := embeddedLibs.ReadFile(embeddedPath)
		if err != nil {
			// Skip files that can't be read
			continue
		}

		// Write library to temporary file
		tempLibPath := filepath.Join(tempDir, fileName)
		err = os.WriteFile(tempLibPath, data, 0600)
		if err != nil {
			_ = os.RemoveAll(tempDir) // Ignore error during cleanup
			return "", fmt.Errorf("failed to write temp library %s: %w", fileName, err)
		}
		extractedCount++
	}

	if extractedCount == 0 {
		_ = os.RemoveAll(tempDir) // Ignore error during cleanup
		return "", fmt.Errorf("no embedded libraries found for platform %s_%s", goos, goarch)
	}

	l.tempDir = tempDir

	// Return path to main library
	mainLib, err := l.getLibraryName()
	if err != nil {
		return "", err
	}

	return filepath.Join(tempDir, mainLib), nil
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
