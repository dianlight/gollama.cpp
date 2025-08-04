package gollama

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestLibraryLoader_GetLibraryName(t *testing.T) {
	loader := &LibraryLoader{}

	// Test current platform
	result, err := loader.getLibraryName()
	if err != nil {
		// If current platform is unsupported, that's fine
		if runtime.GOOS != "darwin" && runtime.GOOS != "linux" && runtime.GOOS != "windows" {
			t.Logf("Current platform %s is unsupported, which is expected", runtime.GOOS)
			return
		}
		t.Errorf("Unexpected error for supported OS %s: %v", runtime.GOOS, err)
		return
	}

	// Verify the result matches expected pattern for current OS
	switch runtime.GOOS {
	case "darwin":
		if result != "libllama.dylib" {
			t.Errorf("Expected libllama.dylib for darwin, got %s", result)
		}
	case "linux":
		if result != "libllama.so" {
			t.Errorf("Expected libllama.so for linux, got %s", result)
		}
	case "windows":
		if result != "llama.dll" {
			t.Errorf("Expected llama.dll for windows, got %s", result)
		}
	default:
		t.Logf("Platform %s returned %s", runtime.GOOS, result)
	}
}

func TestLibraryLoader_LoadSharedLibrary(t *testing.T) {
	loader := &LibraryLoader{}

	t.Run("Windows behavior", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping Windows test on non-Windows platform")
		}

		_, err := loader.loadSharedLibrary("test.dll")
		if err == nil {
			t.Error("Expected error for Windows, but got none")
		}
		if err.Error() != "support for windows platform not yet implemented" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("Unix-like systems", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping Unix test on Windows")
		}

		// Test with invalid path (should fail)
		_, err := loader.loadSharedLibrary("/invalid/path/libtest.so")
		if err == nil {
			t.Error("Expected error for invalid path, but got none")
		}
	})
}

func TestLibraryLoader_ExtractEmbeddedLibrary(t *testing.T) {
	loader := &LibraryLoader{}

	t.Run("Non-existent embedded library", func(t *testing.T) {
		_, err := loader.extractEmbeddedLibrary("nonexistent.so")
		if err == nil {
			t.Error("Expected error for non-existent library, but got none")
		}
	})
}

func TestLibraryLoader_GetHandle(t *testing.T) {
	loader := &LibraryLoader{}

	t.Run("Initially zero", func(t *testing.T) {
		handle := loader.GetHandle()
		if handle != 0 {
			t.Errorf("Expected handle to be 0, got %d", handle)
		}
	})

	t.Run("After setting handle", func(t *testing.T) {
		expectedHandle := uintptr(12345)
		loader.handle = expectedHandle
		loader.loaded = true

		handle := loader.GetHandle()
		if handle != expectedHandle {
			t.Errorf("Expected handle to be %d, got %d", expectedHandle, handle)
		}

		// Reset for other tests
		loader.handle = 0
		loader.loaded = false
	})
}

func TestLibraryLoader_IsLoaded(t *testing.T) {
	loader := &LibraryLoader{}

	t.Run("Initially false", func(t *testing.T) {
		if loader.IsLoaded() {
			t.Error("Expected IsLoaded to be false initially")
		}
	})

	t.Run("After setting loaded", func(t *testing.T) {
		loader.loaded = true

		if !loader.IsLoaded() {
			t.Error("Expected IsLoaded to be true after setting loaded")
		}

		// Reset for other tests
		loader.loaded = false
	})
}

func TestLibraryLoader_UnloadLibrary(t *testing.T) {
	loader := &LibraryLoader{}

	t.Run("Unload when not loaded", func(t *testing.T) {
		err := loader.UnloadLibrary()
		if err != nil {
			t.Errorf("Unexpected error when unloading unloaded library: %v", err)
		}
	})

	t.Run("Unload with temporary directory", func(t *testing.T) {
		// Create a temporary directory
		tempDir, err := os.MkdirTemp("", "gollama-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}

		loader.loaded = true
		loader.handle = uintptr(12345)
		loader.tempDir = tempDir
		loader.libPath = filepath.Join(tempDir, "test.so")

		err = loader.UnloadLibrary()
		if err != nil {
			t.Errorf("Unexpected error when unloading library: %v", err)
		}

		// Check that everything is reset
		if loader.loaded {
			t.Error("Expected loaded to be false after unload")
		}
		if loader.handle != 0 {
			t.Error("Expected handle to be 0 after unload")
		}
		if loader.tempDir != "" {
			t.Error("Expected tempDir to be empty after unload")
		}
		if loader.libPath != "" {
			t.Error("Expected libPath to be empty after unload")
		}

		// Check that temp directory was removed
		if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
			t.Error("Expected temp directory to be removed")
		}
	})
}

func TestLibraryLoader_LoadLibrary(t *testing.T) {
	loader := &LibraryLoader{}

	t.Run("Load when already loaded", func(t *testing.T) {
		loader.loaded = true
		defer func() { loader.loaded = false }()

		err := loader.LoadLibrary()
		if err != nil {
			t.Errorf("Unexpected error when loading already loaded library: %v", err)
		}
	})

	t.Run("Load library - expected failure on missing libs", func(t *testing.T) {
		// Since we don't have actual libraries embedded, this should fail
		// but not panic or cause other issues
		err := loader.LoadLibrary()
		if err == nil {
			// This is unexpected but not necessarily wrong
			t.Log("LoadLibrary succeeded unexpectedly - may have found system library")
		} else {
			t.Logf("LoadLibrary failed as expected: %v", err)
		}

		// Clean up if somehow it succeeded
		if loader.loaded {
			loader.UnloadLibrary()
		}
	})
}

func TestLibraryLoader_ThreadSafety(t *testing.T) {
	loader := &LibraryLoader{}

	t.Run("Concurrent access to GetHandle", func(t *testing.T) {
		const numGoroutines = 100
		var wg sync.WaitGroup

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				_ = loader.GetHandle()
			}()
		}

		wg.Wait()
	})

	t.Run("Concurrent access to IsLoaded", func(t *testing.T) {
		const numGoroutines = 100
		var wg sync.WaitGroup

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				_ = loader.IsLoaded()
			}()
		}

		wg.Wait()
	})

	t.Run("Concurrent LoadLibrary calls", func(t *testing.T) {
		const numGoroutines = 10
		var wg sync.WaitGroup

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				_ = loader.LoadLibrary()
			}()
		}

		wg.Wait()

		// Clean up if any succeeded
		if loader.loaded {
			loader.UnloadLibrary()
		}
	})
}

func TestGlobalFunctions(t *testing.T) {
	t.Run("getLibHandle", func(t *testing.T) {
		handle := getLibHandle()
		expectedHandle := globalLoader.GetHandle()
		if handle != expectedHandle {
			t.Errorf("Expected handle %d, got %d", expectedHandle, handle)
		}
	})

	t.Run("isLibraryLoaded", func(t *testing.T) {
		loaded := isLibraryLoaded()
		expectedLoaded := globalLoader.IsLoaded()
		if loaded != expectedLoaded {
			t.Errorf("Expected loaded %t, got %t", expectedLoaded, loaded)
		}
	})

	t.Run("RegisterFunction with no library", func(t *testing.T) {
		var testFunc func()
		err := RegisterFunction(&testFunc, "test_function")
		if err == nil {
			t.Error("Expected error when registering function with no library loaded")
		}
		if err.Error() != "library not loaded" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("Cleanup", func(t *testing.T) {
		// This should not panic
		Cleanup()
	})
}

func TestLibraryLoader_ExtractEmbeddedLibraryWriteFailure(t *testing.T) {
	loader := &LibraryLoader{}

	t.Run("Write to read-only directory", func(t *testing.T) {
		// Create a temporary directory and make it read-only
		tempDir, err := os.MkdirTemp("", "gollama-readonly-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Make directory read-only
		err = os.Chmod(tempDir, 0444)
		if err != nil {
			t.Skipf("Cannot change directory permissions: %v", err)
		}
		defer os.Chmod(tempDir, 0755) // Restore permissions for cleanup

		// This test is OS-dependent and may not work reliably
		// We'll just verify the function doesn't panic
		_, err = loader.extractEmbeddedLibrary("test.so")
		if err == nil {
			t.Log("extractEmbeddedLibrary succeeded unexpectedly")
		} else {
			t.Logf("extractEmbeddedLibrary failed as expected: %v", err)
		}
	})
}

// Benchmark tests
func BenchmarkLibraryLoader_GetHandle(b *testing.B) {
	loader := &LibraryLoader{}
	loader.handle = uintptr(12345)
	loader.loaded = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = loader.GetHandle()
	}
}

func BenchmarkLibraryLoader_IsLoaded(b *testing.B) {
	loader := &LibraryLoader{}
	loader.loaded = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = loader.IsLoaded()
	}
}

func BenchmarkLibraryLoader_GetLibraryName(b *testing.B) {
	loader := &LibraryLoader{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = loader.getLibraryName()
	}
}

func BenchmarkGlobalFunctions(b *testing.B) {
	b.Run("getLibHandle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = getLibHandle()
		}
	})

	b.Run("isLibraryLoaded", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = isLibraryLoaded()
		}
	})
}

// Test race conditions
func TestLibraryLoader_RaceConditions(t *testing.T) {
	loader := &LibraryLoader{}

	t.Run("Load and Unload race", func(t *testing.T) {
		const iterations = 50
		var wg sync.WaitGroup

		wg.Add(2)

		// Goroutine 1: Try to load
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_ = loader.LoadLibrary()
				time.Sleep(time.Microsecond)
			}
		}()

		// Goroutine 2: Try to unload
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_ = loader.UnloadLibrary()
				time.Sleep(time.Microsecond)
			}
		}()

		wg.Wait()

		// Final cleanup
		if loader.loaded {
			loader.UnloadLibrary()
		}
	})
}

// Test initialization and state
func TestLibraryLoader_InitialState(t *testing.T) {
	loader := &LibraryLoader{}

	if loader.handle != 0 {
		t.Error("Expected initial handle to be 0")
	}
	if loader.loaded {
		t.Error("Expected initial loaded to be false")
	}
	if loader.tempDir != "" {
		t.Error("Expected initial tempDir to be empty")
	}
	if loader.libPath != "" {
		t.Error("Expected initial libPath to be empty")
	}
}

// Test global loader initialization
func TestGlobalLoader(t *testing.T) {
	if globalLoader == nil {
		t.Error("Expected globalLoader to be initialized")
	}

	// Test that global functions work with uninitialized state
	handle := getLibHandle()
	if handle != 0 {
		t.Errorf("Expected global handle to be 0 initially, got %d", handle)
	}

	loaded := isLibraryLoaded()
	if loaded {
		t.Error("Expected global library to not be loaded initially")
	}
}
