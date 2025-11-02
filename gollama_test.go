package gollama

import (
	"runtime"
	"testing"
	"unsafe"
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
		t.Fatalf("Skipping library path test on unsupported platform: %v", err)
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
			t.Fatalf("Skipping utility function tests: library not available: %v", err)
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
		t.Fatalf("Skipping backend test: %v", err)
	}

	// Clean up
	Backend_free()
}

func TestModelParams(t *testing.T) {
	// This test will only run if the library can be loaded
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Fatalf("Skipping model params test: library not available: %v", err)
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
	// FFI now enables struct-returning functions on all platforms
	// This test will only run if the library can be loaded
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Fatalf("Skipping context params test: library not available: %v", err)
		}
	}

	params := Context_default_params()

	// Check that we got some reasonable defaults
	// Note: Some values might be 0 if using fallback defaults
	if params.NBatch == 0 {
		t.Error("NBatch should not be zero")
	}
}

// Benchmark basic operations
func BenchmarkGetLibraryPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = getLibraryPath() // Ignore return values in benchmark
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

// Test default parameters functionality (from debug-params.go)
func TestContextDefaultParams(t *testing.T) {
	// FFI now enables struct-returning functions on all platforms

	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Fatalf("Skipping context default params test: library not available: %v", err)
		}
	}

	params := Context_default_params()

	// Verify that parameters have reasonable default values
	if params.NSeqMax == 0 {
		t.Error("NSeqMax should not be zero")
	}
	if params.NCtx == 0 {
		t.Error("NCtx should not be zero")
	}
	if params.NBatch == 0 {
		t.Error("NBatch should not be zero")
	}
	if params.NUbatch == 0 {
		t.Error("NUbatch should not be zero")
	}

	t.Logf("Default NSeqMax: %d", params.NSeqMax)
	t.Logf("Default NCtx: %d", params.NCtx)
	t.Logf("Default NBatch: %d", params.NBatch)
	t.Logf("Default NUbatch: %d", params.NUbatch)
}

// Test token data array functionality (from token_array_test.go)
func TestTokenDataArrayFromLogits(t *testing.T) {
	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Fatalf("Skipping token data array test: library not available: %v", err)
		}
	}

	// Create dummy logits array
	logits := make([]float32, 256)
	for i := 0; i < 256; i++ {
		logits[i] = float32(i) * 0.1
	}

	// Call our function with the logits
	// We don't need a real model since the function doesn't use it currently
	tokenArray := Token_data_array_from_logits(LlamaModel(0), &logits[0])

	if tokenArray == nil {
		t.Fatal("Token array should not be nil")
	}

	// The actual size may be different from what we expect, so let's just check it's reasonable
	if tokenArray.Size == 0 {
		t.Error("Token array size should not be zero")
	}

	// Check that Selected is initialized correctly
	if tokenArray.Selected != -1 {
		t.Errorf("Expected Selected to be -1, got %d", tokenArray.Selected)
	}

	// Check that Sorted is initialized correctly
	if tokenArray.Sorted != 0 {
		t.Errorf("Expected Sorted to be 0, got %d", tokenArray.Sorted)
	}

	// Check that we can access the first element
	if tokenArray.Data == nil {
		t.Fatal("Data pointer should not be nil")
	}

	firstToken := tokenArray.Data
	if firstToken.Id != 0 {
		t.Errorf("Expected first token ID to be 0, got %d", firstToken.Id)
	}

	if firstToken.Logit != 0.0 {
		t.Errorf("Expected first token logit to be 0.0, got %f", firstToken.Logit)
	}

	t.Logf("SUCCESS: Token array created with size %d", tokenArray.Size)
	t.Logf("Data pointer: %p", tokenArray.Data)
	t.Logf("First token: ID=%d, Logit=%f", firstToken.Id, firstToken.Logit)

	// Only test accessing other elements if size is large enough
	if tokenArray.Size > 1 {
		// Check that we can access the last element
		lastIndex := tokenArray.Size - 1
		lastElement := (*LlamaTokenData)(unsafe.Pointer(uintptr(unsafe.Pointer(tokenArray.Data)) + uintptr(lastIndex)*unsafe.Sizeof(LlamaTokenData{})))
		t.Logf("Last token: ID=%d, Logit=%f", lastElement.Id, lastElement.Logit)
	}
}

// Test tokenization functionality (from test_tokenize.go)
// This test requires a model file, so it's marked as an integration test
func TestTokenization(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if !isLoaded {
		err := loadLibrary()
		if err != nil {
			t.Fatalf("Skipping tokenization test: library not available: %v", err)
		}
	}

	err := Backend_init()
	if err != nil {
		t.Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	// Look for a test model in the models directory
	modelPath := "./models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"

	// Load the model
	params := Model_default_params()
	params.NGpuLayers = 0 // Use CPU for testing

	model, err := Model_load_from_file(modelPath, params)
	if err != nil {
		t.Fatalf("Tokenization test: model not available at %s: %v", modelPath, err)
	}
	defer Model_free(model)

	t.Log("Model loaded successfully")

	// Test simple tokenization
	testText := "Hello world"
	tokens, err := Tokenize(model, testText, false, false)
	if err != nil {
		t.Fatalf("Failed to tokenize: %v", err)
	}

	if len(tokens) == 0 {
		t.Error("Expected at least one token")
	}

	t.Logf("Tokenized '%s' into %d tokens: %v", testText, len(tokens), tokens)

	// Test with different parameters
	tokensWithBos, err := Tokenize(model, testText, true, false)
	if err != nil {
		t.Fatalf("Failed to tokenize with BOS: %v", err)
	}

	if len(tokensWithBos) <= len(tokens) {
		t.Error("Expected more tokens when adding BOS")
	}

	t.Logf("Tokenized with BOS: %d tokens: %v", len(tokensWithBos), tokensWithBos)
}

// TestGpuBackendDetection tests GPU backend detection functionality
func TestGpuBackendDetection(t *testing.T) {
	backend := DetectGpuBackend()

	t.Logf("Detected GPU backend: %s (%d)", backend.String(), int(backend))

	// Verify the backend is valid
	if backend < LLAMA_GPU_BACKEND_NONE || backend > LLAMA_GPU_BACKEND_SYCL {
		t.Errorf("Invalid GPU backend detected: %d", backend)
	}

	// Platform-specific expectations
	switch runtime.GOOS {
	case "darwin":
		// On macOS, we expect Metal (unless running in a container/VM)
		if backend != LLAMA_GPU_BACKEND_METAL && backend != LLAMA_GPU_BACKEND_CPU {
			t.Logf("Note: Expected Metal on macOS, got %s", backend.String())
		}
	case "linux", "windows":
		// On Linux/Windows, we expect any valid backend
		if backend == LLAMA_GPU_BACKEND_NONE {
			t.Error("Expected valid GPU backend detection on Linux/Windows")
		}
	}
}

// TestGpuBackendString tests the String() method of LlamaGpuBackend
func TestGpuBackendString(t *testing.T) {
	tests := []struct {
		backend  LlamaGpuBackend
		expected string
	}{
		{LLAMA_GPU_BACKEND_NONE, "None"},
		{LLAMA_GPU_BACKEND_CPU, "CPU"},
		{LLAMA_GPU_BACKEND_CUDA, "CUDA"},
		{LLAMA_GPU_BACKEND_METAL, "Metal"},
		{LLAMA_GPU_BACKEND_HIP, "HIP"},
		{LLAMA_GPU_BACKEND_VULKAN, "Vulkan"},
		{LLAMA_GPU_BACKEND_OPENCL, "OpenCL"},
		{LLAMA_GPU_BACKEND_SYCL, "SYCL"},
		{LlamaGpuBackend(999), "Unknown"},
	}

	for _, tt := range tests {
		result := tt.backend.String()
		if result != tt.expected {
			t.Errorf("Backend %d String() = %s, want %s", tt.backend, result, tt.expected)
		}
	}
}

// TestCommandDetection tests the hasCommand function
func TestCommandDetection(t *testing.T) {
	// Test with commands that should exist on most systems
	commonCommands := []string{"go", "echo"}

	for _, cmd := range commonCommands {
		if !hasCommand(cmd) {
			t.Logf("Command '%s' not found (this may be expected in some environments)", cmd)
		}
	}

	// Test with a command that definitely shouldn't exist
	if hasCommand("definitely-not-a-real-command-12345") {
		t.Error("hasCommand should return false for non-existent commands")
	}

	// Test GPU-related commands (these may or may not be available)
	gpuCommands := []string{"nvcc", "hipconfig", "vulkaninfo", "clinfo", "sycl-ls"}
	for _, cmd := range gpuCommands {
		found := hasCommand(cmd)
		t.Logf("GPU command '%s' found: %t", cmd, found)
	}
}
