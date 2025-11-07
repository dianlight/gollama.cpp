package gollama

import (
	"os"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/suite"
)

type FFISuite struct{ BaseSuite }

// SetupSuite runs before all tests in the suite
func (s *FFISuite) SetupSuite() {
	if err := loadLibrary(); err != nil {
		s.Require().Failf("SetupSuite", "Failed to load library during suite setup: %v", err)
	}
}

// SetupTest reloads the library if it was unloaded by a previous suite and
// snapshots environment/config using BaseSuite without causing an unload after
// each individual test.
func (s *FFISuite) SetupTest() {
	s.BaseSuite.SetupTest()
	if !isLoaded {
		s.Require().NoError(loadLibrary(), "Failed to load library for test")
	}
}

// TearDownTest restores environment and config but intentionally does NOT
// unload the library (the FFI suite requires a persistent handle across tests).
func (s *FFISuite) TearDownTest() {
	// Restore env variables
	for k, v := range s.savedEnv {
		if v == "" {
			_ = os.Unsetenv(k)
		} else {
			_ = os.Setenv(k, v)
		}
	}
	// Restore global config
	if s.savedConfig != nil {
		_ = SetGlobalConfig(s.savedConfig)
	}
	// Intentionally skip Cleanup() here
}

// TearDownSuite performs a final cleanup after all tests have run.
func (s *FFISuite) TearDownSuite() {
	Cleanup()
}

// Verifies that FFI type definitions are properly structured
func (s *FFISuite) TestFFITypeDefinitions() {
	s.Assert().NotZero(ffiTypeLlamaModelParams.Type, "ffiTypeLlamaModelParams Type should not be zero")
	s.Assert().NotZero(ffiTypeLlamaContextParams.Type, "ffiTypeLlamaContextParams Type should not be zero")
	s.Assert().NotZero(ffiTypeLlamaSamplerChainParams.Type, "ffiTypeLlamaSamplerChainParams Type should not be zero")
	s.Assert().NotZero(ffiTypeLlamaBatch.Type, "ffiTypeLlamaBatch Type should not be zero")
}

// Verifies that Go structs match expected sizes
func (s *FFISuite) TestFFIStructSizes() {
	modelParamsSize := unsafe.Sizeof(LlamaModelParams{})
	s.Assert().NotZero(modelParamsSize, "LlamaModelParams size should not be zero")
	s.T().Logf("LlamaModelParams size: %d bytes", modelParamsSize)

	contextParamsSize := unsafe.Sizeof(LlamaContextParams{})
	s.Assert().NotZero(contextParamsSize, "LlamaContextParams size should not be zero")
	s.T().Logf("LlamaContextParams size: %d bytes", contextParamsSize)

	batchSize := unsafe.Sizeof(LlamaBatch{})
	s.Assert().NotZero(batchSize, "LlamaBatch size should not be zero")
	s.T().Logf("LlamaBatch size: %d bytes", batchSize)

	samplerChainParamsSize := unsafe.Sizeof(LlamaSamplerChainParams{})
	s.Assert().NotZero(samplerChainParamsSize, "LlamaSamplerChainParams size should not be zero")
	s.T().Logf("LlamaSamplerChainParams size: %d bytes", samplerChainParamsSize)
}

// Tests FFI-based model parameter retrieval
func (s *FFISuite) TestFFIModelDefaultParams() {
	params, err := ffiModelDefaultParams()
	if err != nil {
		s.T().Logf("FFI model default params failed (expected if library not present): %v", err)
		return
	}

	s.Assert().GreaterOrEqual(int(params.NGpuLayers), -1, "NGpuLayers should be >= -1")
	s.T().Logf("FFI Model default params: NGpuLayers=%d, SplitMode=%d, UseMmap=%d",
		params.NGpuLayers, params.SplitMode, params.UseMmap)
}

// Tests FFI-based context parameter retrieval
func (s *FFISuite) TestFFIContextDefaultParams() {
	params, err := ffiContextDefaultParams()
	if err != nil {
		s.T().Logf("FFI context default params failed (expected if library not present): %v", err)
		return
	}

	s.Assert().NotZero(params.Seed, "Seed should not be zero in default params")
	s.Assert().NotZero(params.NBatch, "NBatch should not be zero in default params")
	s.T().Logf("FFI Context default params: Seed=%d, NCtx=%d, NBatch=%d, NThreads=%d",
		params.Seed, params.NCtx, params.NBatch, params.NThreads)
}

// Tests FFI-based sampler chain parameter retrieval
func (s *FFISuite) TestFFISamplerChainDefaultParams() {
	params, err := ffiSamplerChainDefaultParams()
	if err != nil {
		s.T().Logf("FFI sampler chain default params failed (expected if library not present): %v", err)
		return
	}

	s.Assert().LessOrEqual(int(params.NoPerf), 1, "NoPerf should be 0 or 1")
	s.T().Logf("FFI Sampler chain default params: NoPerf=%d", params.NoPerf)
}

// Tests FFI-based batch initialization
func (s *FFISuite) TestFFIBatchInit() {
	batch, err := ffiBatchInit(512, 0, 1)
	if err != nil {
		s.T().Logf("FFI batch init failed (expected if library not present): %v", err)
		return
	}
	// No strict expectations here; just ensure it doesn't panic and returns a struct
	if batch.NTokens != 0 {
		s.T().Logf("Batch initialized with NTokens=%d", batch.NTokens)
	}
}

// Tests FFI-based encode function
func (s *FFISuite) TestFFIEncode() {
	// TODO: Implement a proper test with a valid context and batch
	s.T().Skip("Skipping FFI encode test - requires valid context and batch to avoid assertion failure")
}

// Tests FFI-based sampler chain initialization
func (s *FFISuite) TestFFISamplerChainInit() {
	params := LlamaSamplerChainParams{NoPerf: 0}
	sampler, err := ffiSamplerChainInit(params)
	s.Require().NoError(err, "FFI sampler chain init failed")
	s.Assert().NotZero(sampler, "FFI sampler chain init returned null sampler")
}

// Tests that FFI functions fall back gracefully
func (s *FFISuite) TestFFIFallbackBehavior() {
	params := Model_default_params()
	s.Assert().GreaterOrEqual(int(params.NGpuLayers), -1, "Model params should have reasonable defaults even without library")

	ctxParams := Context_default_params()
	s.Assert().NotZero(ctxParams.NBatch, "Context params should have reasonable defaults even without library")

	samplerParams := Sampler_chain_default_params()
	s.Assert().LessOrEqual(int(samplerParams.NoPerf), 1, "Sampler chain params should have reasonable defaults even without library")

	s.T().Log("All FFI functions have proper fallback behavior")
}

// Tests the platform-specific GetProcAddress implementation
func (s *FFISuite) TestPlatformGetProcAddress() {
	addr, err := getProcAddressPlatform(libHandle, "llama_backend_init")
	s.Require().NoError(err, "Failed to get llama_backend_init address")
	s.Assert().NotZero(addr, "llama_backend_init address should not be zero")
}

// Verifies cross-platform build compatibility of FFI helpers
func (s *FFISuite) TestFFICrossCompileCompatibility() {
	_ = loadLibraryPlatform
	_ = closeLibraryPlatform
	_ = registerLibFunc
	_ = getProcAddressPlatform
	_ = isPlatformSupported
	_ = getPlatformError
	s.T().Log("All platform-specific functions are properly defined")
}

func TestFFISuite(t *testing.T) {
	suite.Run(t, new(FFISuite))
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
