package gollama

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GollamaSuite struct{ BaseSuite }

func (s *GollamaSuite) TestVersion() {
	assert.NotEmpty(s.T(), Version)
	assert.NotEmpty(s.T(), LlamaCppBuild)
	assert.NotEmpty(s.T(), FullVersion)
	expectedFull := "v" + Version + "-llamacpp." + LlamaCppBuild
	assert.Equal(s.T(), expectedFull, FullVersion)
}

func (s *GollamaSuite) TestLibraryPath() {
	path, err := getLibraryPath()
	if err != nil {
		s.T().Skipf("Skipping library path test on unsupported platform: %v", err)
	}
	assert.NotEmpty(s.T(), path)
}

func (s *GollamaSuite) TestConstants() {
	assert.Equal(s.T(), int(0xFFFFFFFF), LLAMA_DEFAULT_SEED)
	assert.Equal(s.T(), int(-1), LLAMA_TOKEN_NULL)
}

// Ensure we can call functions that don't require a loaded library
func (s *GollamaSuite) TestUtilityFunctions() {
	if !isLoaded {
		if err := loadLibrary(); err != nil {
			s.T().Skipf("Skipping utility function tests: library not available: %v", err)
		}
	}
	_ = Supports_mmap()
	_ = Supports_mlock()
	_ = Supports_gpu_offload()
	_ = Max_devices()
	s.T().Log("Utility functions executed successfully")
}

func (s *GollamaSuite) TestBackendInitialization() {
	if err := Backend_init(); err != nil {
		s.T().Skipf("Skipping backend test: %v", err)
	}
	Backend_free()
}

func (s *GollamaSuite) TestModelParams() {
	if !isLoaded {
		if err := loadLibrary(); err != nil {
			s.T().Skipf("Skipping model params test: library not available: %v", err)
		}
	}
	params := Model_default_params()
	assert.GreaterOrEqual(s.T(), int(params.NGpuLayers), 0, "NGpuLayers should not be negative")
}

func (s *GollamaSuite) TestContextParams() {
	if !isLoaded {
		if err := loadLibrary(); err != nil {
			s.T().Skipf("Skipping context params test: library not available: %v", err)
		}
	}
	params := Context_default_params()
	assert.NotZero(s.T(), params.NBatch, "NBatch should not be zero")
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
func (s *GollamaSuite) TestContextDefaultParamsDetailed() {
	if !isLoaded {
		if err := loadLibrary(); err != nil {
			s.T().Skipf("Skipping context default params test: library not available: %v", err)
		}
	}
	params := Context_default_params()
	assert.NotZero(s.T(), params.NSeqMax)
	assert.NotZero(s.T(), params.NCtx)
	assert.NotZero(s.T(), params.NBatch)
	assert.NotZero(s.T(), params.NUbatch)
	s.T().Logf("Default NSeqMax: %d", params.NSeqMax)
	s.T().Logf("Default NCtx: %d", params.NCtx)
	s.T().Logf("Default NBatch: %d", params.NBatch)
	s.T().Logf("Default NUbatch: %d", params.NUbatch)
}

// Test token data array functionality (from token_array_test.go)
func (s *GollamaSuite) TestTokenDataArrayFromLogits() {
	if !isLoaded {
		if err := loadLibrary(); err != nil {
			s.T().Skipf("Skipping token data array test: library not available: %v", err)
		}
	}
	logits := make([]float32, 256)
	for i := 0; i < 256; i++ {
		logits[i] = float32(i) * 0.1
	}
	tokenArray := Token_data_array_from_logits(LlamaModel(0), &logits[0])
	if tokenArray == nil {
		s.T().Fatal("Token array should not be nil")
	}
	assert.NotZero(s.T(), tokenArray.Size)
	assert.Equal(s.T(), int64(-1), tokenArray.Selected)
	assert.Equal(s.T(), uint8(0x0), tokenArray.Sorted)
	if tokenArray.Data == nil {
		s.T().Fatal("Data pointer should not be nil")
	}
	firstToken := tokenArray.Data
	assert.Equal(s.T(), LlamaToken(0), firstToken.Id)
	assert.Equal(s.T(), float32(0.0), firstToken.Logit)
	s.T().Logf("SUCCESS: Token array created with size %d", tokenArray.Size)
	s.T().Logf("Data pointer: %p", tokenArray.Data)
	s.T().Logf("First token: ID=%d, Logit=%f", firstToken.Id, firstToken.Logit)
	if tokenArray.Size > 1 {
		lastIndex := tokenArray.Size - 1
		lastElement := (*LlamaTokenData)(unsafe.Pointer(uintptr(unsafe.Pointer(tokenArray.Data)) + uintptr(lastIndex)*unsafe.Sizeof(LlamaTokenData{})))
		s.T().Logf("Last token: ID=%d, Logit=%f", lastElement.Id, lastElement.Logit)
	}
}

// Test tokenization functionality (from test_tokenize.go)
// This test requires a model file, so it's marked as an integration test
func (s *GollamaSuite) TestTokenization() {
	if testing.Short() {
		s.T().Skip("Skipping integration test in short mode")
	}
	if !isLoaded {
		if err := loadLibrary(); err != nil {
			s.T().Skipf("Skipping tokenization test: library not available: %v", err)
		}
	}
	if err := Backend_init(); err != nil {
		s.T().Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()
	if err := Ggml_backend_load_all(); err != nil {
		s.T().Logf("ggml_backend_load not available or failed: %v", err)
		return
	}
	modelPath := "./models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"
	params := Model_default_params()
	params.NGpuLayers = 0
	model, err := Model_load_from_file(modelPath, params)
	if err != nil {
		s.T().Fatalf("Tokenization test: model not available at %s: %v", modelPath, err)
	}
	defer Model_free(model)
	s.T().Log("Model loaded successfully")
	testText := "Hello world"
	tokens, err := Tokenize(model, testText, false, false)
	if err != nil {
		s.T().Fatalf("Failed to tokenize: %v", err)
	}
	assert.NotEmpty(s.T(), tokens)
	s.T().Logf("Tokenized '%s' into %d tokens: %v", testText, len(tokens), tokens)
	tokensWithBos, err := Tokenize(model, testText, true, false)
	if err != nil {
		s.T().Fatalf("Failed to tokenize with BOS: %v", err)
	}
	assert.Greater(s.T(), len(tokensWithBos), len(tokens), "Expected more tokens when adding BOS")
	s.T().Logf("Tokenized with BOS: %d tokens: %v", len(tokensWithBos), tokensWithBos)
}

// TestGpuBackendDetection tests GPU backend detection functionality
func (s *GollamaSuite) TestGpuBackendDetection() {
	backend := DetectGpuBackend()
	s.T().Logf("Detected GPU backend: %s (%d)", backend.String(), int(backend))
	assert.GreaterOrEqual(s.T(), int(backend), int(LLAMA_GPU_BACKEND_NONE))
	assert.LessOrEqual(s.T(), int(backend), int(LLAMA_GPU_BACKEND_SYCL))
	switch runtime.GOOS {
	case "darwin":
		if backend != LLAMA_GPU_BACKEND_METAL && backend != LLAMA_GPU_BACKEND_CPU {
			s.T().Logf("Note: Expected Metal on macOS, got %s", backend.String())
		}
	case "linux", "windows":
		assert.NotEqual(s.T(), LLAMA_GPU_BACKEND_NONE, backend, "Expected valid GPU backend detection on Linux/Windows")
	}
}

// TestGpuBackendString tests the String() method of LlamaGpuBackend
func (s *GollamaSuite) TestGpuBackendString() {
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
		assert.Equal(s.T(), tt.expected, result)
	}
}

// TestCommandDetection tests the hasCommand function
func (s *GollamaSuite) TestCommandDetection() {
	// Test with commands that should exist on most systems
	commonCommands := []string{"go", "echo"}

	for _, cmd := range commonCommands {
		if !hasCommand(cmd) {
			s.T().Logf("Command '%s' not found (this may be expected in some environments)", cmd)
		}
	}

	// Test with a command that definitely shouldn't exist
	assert.False(s.T(), hasCommand("definitely-not-a-real-command-12345"), "hasCommand should return false for non-existent commands")

	// Test GPU-related commands (these may or may not be available)
	gpuCommands := []string{"nvcc", "hipconfig", "vulkaninfo", "clinfo", "sycl-ls"}
	for _, cmd := range gpuCommands {
		found := hasCommand(cmd)
		s.T().Logf("GPU command '%s' found: %t", cmd, found)
	}
}

func TestGollamaSuite(t *testing.T) { suite.Run(t, new(GollamaSuite)) }
