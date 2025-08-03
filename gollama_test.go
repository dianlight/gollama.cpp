package gollama

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if LlamaCppBuild == "" {
		t.Error("LlamaCppBuild should not be empty")
	}
	if FullVersion == "" {
		t.Error("FullVersion should not be empty")
	}

	expectedFull := "v" + Version + "-llamacpp." + LlamaCppBuild
	if FullVersion != expectedFull {
		t.Errorf("FullVersion mismatch: got %s, want %s", FullVersion, expectedFull)
	}
}

func TestLibraryPath(t *testing.T) {
	path, err := getLibraryPath()
	if err != nil {
		t.Skipf("Skipping library path test on unsupported platform: %v", err)
	}
	if path == "" {
		t.Error("Library path should not be empty")
	}
}

func TestConstants(t *testing.T) {
	// Test that constants have expected values
	if LLAMA_DEFAULT_SEED != 0xFFFFFFFF {
		t.Errorf("LLAMA_DEFAULT_SEED mismatch: got %x, want %x", LLAMA_DEFAULT_SEED, 0xFFFFFFFF)
	}
	if LLAMA_TOKEN_NULL != -1 {
		t.Errorf("LLAMA_TOKEN_NULL mismatch: got %d, want %d", LLAMA_TOKEN_NULL, -1)
	}
}

// Test that we can call functions that don't require a loaded library
func TestUtilityFunctions(t *testing.T) {
	// These tests will only run if the library can be loaded
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Skipf("Skipping utility function tests: library not available: %v", err)
		}
	}

	// Test system capability queries
	_ = Supports_mmap()
	_ = Supports_mlock()
	_ = Supports_gpu_offload()
	_ = Max_devices()

	// These should not panic
	t.Log("Utility functions executed successfully")
}

func TestBackendInitialization(t *testing.T) {
	// Test backend initialization
	err := Backend_init()
	if err != nil {
		t.Skipf("Skipping backend test: %v", err)
	}

	// Clean up
	Backend_free()
}

func TestModelParams(t *testing.T) {
	// This test will only run if the library can be loaded
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Skipf("Skipping model params test: library not available: %v", err)
		}
	}

	params := Model_default_params()

	// Check that we got some reasonable defaults
	// The exact values depend on the llama.cpp implementation
	if params.NGpuLayers < 0 {
		t.Error("NGpuLayers should not be negative")
	}
}

func TestContextParams(t *testing.T) {
	// This test will only run if the library can be loaded
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Skipf("Skipping context params test: library not available: %v", err)
		}
	}

	params := Context_default_params()

	// Check that we got some reasonable defaults
	if params.NCtx == 0 {
		t.Error("NCtx should not be zero")
	}
	if params.NBatch == 0 {
		t.Error("NBatch should not be zero")
	}
}

// Benchmark basic operations
func BenchmarkGetLibraryPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getLibraryPath()
	}
}

func BenchmarkModelDefaultParams(b *testing.B) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			b.Skipf("Skipping benchmark: library not available: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Model_default_params()
	}
}

func BenchmarkContextDefaultParams(b *testing.B) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			b.Skipf("Skipping benchmark: library not available: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Context_default_params()
	}
}
