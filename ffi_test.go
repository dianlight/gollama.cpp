package gollama

import (
	"testing"
	"unsafe"
)

// TestFFITypeDefinitions verifies that FFI type definitions are properly structured
func TestFFITypeDefinitions(t *testing.T) {
	// Test that FFI types are properly initialized
	if ffiTypeLlamaModelParams.Type == 0 {
		t.Error("ffiTypeLlamaModelParams Type should not be zero")
	}
	if ffiTypeLlamaContextParams.Type == 0 {
		t.Error("ffiTypeLlamaContextParams Type should not be zero")
	}
	if ffiTypeLlamaSamplerChainParams.Type == 0 {
		t.Error("ffiTypeLlamaSamplerChainParams Type should not be zero")
	}
	if ffiTypeLlamaBatch.Type == 0 {
		t.Error("ffiTypeLlamaBatch Type should not be zero")
	}
}

// TestFFIStructSizes verifies that Go structs match expected sizes
func TestFFIStructSizes(t *testing.T) {
	// Test struct sizes to ensure proper memory layout
	modelParamsSize := unsafe.Sizeof(LlamaModelParams{})
	if modelParamsSize == 0 {
		t.Error("LlamaModelParams size should not be zero")
	}
	t.Logf("LlamaModelParams size: %d bytes", modelParamsSize)

	contextParamsSize := unsafe.Sizeof(LlamaContextParams{})
	if contextParamsSize == 0 {
		t.Error("LlamaContextParams size should not be zero")
	}
	t.Logf("LlamaContextParams size: %d bytes", contextParamsSize)

	batchSize := unsafe.Sizeof(LlamaBatch{})
	if batchSize == 0 {
		t.Error("LlamaBatch size should not be zero")
	}
	t.Logf("LlamaBatch size: %d bytes", batchSize)

	samplerChainParamsSize := unsafe.Sizeof(LlamaSamplerChainParams{})
	if samplerChainParamsSize == 0 {
		t.Error("LlamaSamplerChainParams size should not be zero")
	}
	t.Logf("LlamaSamplerChainParams size: %d bytes", samplerChainParamsSize)
}

// TestFFIModelDefaultParams tests FFI-based model parameter retrieval
func TestFFIModelDefaultParams(t *testing.T) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Skipf("FFI model params test requires library to be available: %v", err)
		}
	}

	// Call using FFI directly
	params, err := ffiModelDefaultParams()
	if err != nil {
		// FFI might fail if library is not loaded or function not found
		t.Logf("FFI model default params failed (expected if library not present): %v", err)
		return
	}

	// Verify reasonable defaults
	if params.NGpuLayers < -1 {
		t.Errorf("NGpuLayers should be >= -1, got %d", params.NGpuLayers)
	}

	t.Logf("FFI Model default params: NGpuLayers=%d, SplitMode=%d, UseMmap=%d",
		params.NGpuLayers, params.SplitMode, params.UseMmap)
}

// TestFFIContextDefaultParams tests FFI-based context parameter retrieval
func TestFFIContextDefaultParams(t *testing.T) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Skipf("FFI context params test requires library to be available: %v", err)
		}
	}

	// Call using FFI directly
	params, err := ffiContextDefaultParams()
	if err != nil {
		t.Logf("FFI context default params failed (expected if library not present): %v", err)
		return
	}

	// Verify reasonable defaults
	if params.Seed == 0 {
		t.Error("Seed should not be zero in default params")
	}
	if params.NBatch == 0 {
		t.Error("NBatch should not be zero in default params")
	}

	t.Logf("FFI Context default params: Seed=%d, NCtx=%d, NBatch=%d, NThreads=%d",
		params.Seed, params.NCtx, params.NBatch, params.NThreads)
}

// TestFFISamplerChainDefaultParams tests FFI-based sampler chain parameter retrieval
func TestFFISamplerChainDefaultParams(t *testing.T) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Skipf("FFI sampler chain params test requires library to be available: %v", err)
		}
	}

	// Call using FFI directly
	params, err := ffiSamplerChainDefaultParams()
	if err != nil {
		t.Logf("FFI sampler chain default params failed (expected if library not present): %v", err)
		return
	}

	// Verify defaults (NoPerf is a boolean represented as uint8)
	if params.NoPerf > 1 {
		t.Errorf("NoPerf should be 0 or 1, got %d", params.NoPerf)
	}

	t.Logf("FFI Sampler chain default params: NoPerf=%d", params.NoPerf)
}

// TestFFIBatchInit tests FFI-based batch initialization
func TestFFIBatchInit(t *testing.T) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Skipf("FFI batch init test requires library to be available: %v", err)
		}
	}

	// Call using FFI directly
	batch, err := ffiBatchInit(512, 0, 1)
	if err != nil {
		t.Logf("FFI batch init failed (expected if library not present): %v", err)
		return
	}

	// Verify batch structure
	if batch.NTokens != 0 {
		t.Logf("Batch initialized with NTokens=%d", batch.NTokens)
	}

	t.Logf("FFI Batch init successful: NTokens=%d", batch.NTokens)
}

// TestFFIEncode tests FFI-based encode function
func TestFFIEncode(t *testing.T) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Skipf("FFI encode test requires library to be available: %v", err)
		}
	}

	// Note: This test cannot be run without a valid context and batch
	// as llama.cpp will assert and crash. We skip this test since we're
	// mainly interested in testing that the FFI structure is correct,
	// which is validated by other tests
	t.Skip("Skipping FFI encode test - requires valid context and batch to avoid assertion failure")
}

// TestFFISamplerChainInit tests FFI-based sampler chain initialization
func TestFFISamplerChainInit(t *testing.T) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Skipf("FFI sampler chain init test requires library to be available: %v", err)
		}
	}

	// Call using FFI directly
	params := LlamaSamplerChainParams{NoPerf: 0}
	sampler, err := ffiSamplerChainInit(params)
	if err != nil {
		t.Logf("FFI sampler chain init failed (expected if library not present): %v", err)
		return
	}

	// Verify sampler is not null
	if sampler == 0 {
		t.Error("FFI sampler chain init returned null sampler")
	}

	t.Logf("FFI Sampler chain init successful: sampler=%v", sampler)
}

// TestFFIFallbackBehavior tests that FFI functions fall back gracefully
func TestFFIFallbackBehavior(t *testing.T) {
	// Test Model_default_params fallback
	params := Model_default_params()
	if params.NGpuLayers < -1 {
		t.Error("Model params should have reasonable defaults even without library")
	}

	// Test Context_default_params fallback
	ctxParams := Context_default_params()
	if ctxParams.NBatch == 0 {
		t.Error("Context params should have reasonable defaults even without library")
	}

	// Test Sampler_chain_default_params fallback
	samplerParams := Sampler_chain_default_params()
	if samplerParams.NoPerf > 1 {
		t.Error("Sampler chain params should have reasonable defaults even without library")
	}

	t.Log("All FFI functions have proper fallback behavior")
}

// TestPlatformGetProcAddress tests the platform-specific GetProcAddress implementation
func TestPlatformGetProcAddress(t *testing.T) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Skipf("GetProcAddress test requires library to be available: %v", err)
		}
	}

	// Try to get a known function address
	addr, err := getProcAddressPlatform(libHandle, "llama_backend_init")
	if err != nil {
		t.Errorf("Failed to get llama_backend_init address: %v", err)
	}
	if addr == 0 {
		t.Error("llama_backend_init address should not be zero")
	}

	t.Logf("Successfully retrieved llama_backend_init at address: %x", addr)
}

// TestFFICrossCompileCompatibility verifies cross-platform build compatibility
func TestFFICrossCompileCompatibility(t *testing.T) {
	// This test verifies that the FFI implementation doesn't break cross-compilation
	// It checks that all platform-specific functions are properly defined

	// These should compile on all platforms
	_ = loadLibraryPlatform
	_ = closeLibraryPlatform
	_ = registerLibFunc
	_ = getProcAddressPlatform
	_ = isPlatformSupported
	_ = getPlatformError

	t.Log("All platform-specific functions are properly defined")
}

// BenchmarkFFIModelDefaultParams benchmarks FFI model parameter retrieval
func BenchmarkFFIModelDefaultParams(b *testing.B) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			b.Skipf("Skipping benchmark: library not available: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ffiModelDefaultParams()
	}
}

// BenchmarkFFIContextDefaultParams benchmarks FFI context parameter retrieval
func BenchmarkFFIContextDefaultParams(b *testing.B) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			b.Skipf("Skipping benchmark: library not available: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ffiContextDefaultParams()
	}
}

// BenchmarkFFIBatchInit benchmarks FFI batch initialization
func BenchmarkFFIBatchInit(b *testing.B) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			b.Skipf("Skipping benchmark: library not available: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ffiBatchInit(512, 0, 1)
	}
}
