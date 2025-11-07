package gollama

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Additional tests to increase coverage for gollama.go
type GollamaMoreSuite struct{ BaseSuite }

func (s *GollamaMoreSuite) SetupTest() {
	s.BaseSuite.SetupTest()
	// Initialize backend once per test to exercise ensureLoaded and registration paths
	err := Backend_init()
	if err != nil {
		s.T().Fatalf("Backend_init failed: %v", err)
	}
}

func (s *GollamaMoreSuite) TearDownTest() {
	Backend_free()
	s.BaseSuite.TearDownTest()
}

// Exercises Token_data_array_init happy path without requiring a model
func (s *GollamaMoreSuite) TestTokenDataArrayInit() {
	arr := Token_data_array_init(0)
	require.NotNil(s.T(), arr)
	assert.Equal(s.T(), uint64(256), arr.Size)
	require.NotNil(s.T(), arr.Data)
}

// Load a tiny model, check a few simple APIs that were previously uncovered
func (s *GollamaMoreSuite) TestModelAndContextBasics() {
	// Load model
	modelPath := "./models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"
	params := Model_default_params()
	params.NGpuLayers = 0
	model, err := Model_load_from_file(modelPath, params)
	if err != nil {
		s.T().Errorf("Model not available at %s: %v", modelPath, err)
		return
	}
	defer Model_free(model)

	// Basic model query
	nEmb := Model_n_embd(model)
	assert.Greater(s.T(), nEmb, int32(0))

	// Create context from model
	ctxParams := Context_default_params()
	ctx, err := Init_from_model(model, ctxParams)
	require.NoError(s.T(), err)
	defer Free(ctx)

	// Exercise simple getters/setters that previously had no coverage
	Set_causal_attn(ctx, true)
	Set_embeddings(ctx, true)
	_ = Get_memory(ctx)
	_ = Memory_clear(ctx, true)
	_ = Get_logits(ctx)
	_ = Get_logits_ith(ctx, 0)
}

// Test Batch_get_one path and Token_to_piece using the model's vocab
func (s *GollamaMoreSuite) TestBatchAndTokenPiece() {
	modelPath := "./models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf"
	params := Model_default_params()
	params.NGpuLayers = 0
	model, err := Model_load_from_file(modelPath, params)
	if err != nil {
		s.T().Errorf("Model not available at %s: %v", modelPath, err)
		return
	}
	defer Model_free(model)

	// Create a trivial batch
	tokens := []LlamaToken{1, 2, 3}
	batch := Batch_get_one(tokens)
	// We only assert that the call succeeds without panicking; content may be implementation-defined
	_ = batch

	// Use internal vocab helpers to fetch a known token and convert it to string
	vocab := llamaModelGetVocab(model)
	if vocab != 0 && llamaVocabBos != nil {
		bos := llamaVocabBos(vocab)
		piece := Token_to_piece(model, bos, false)
		// Some builds may return empty BOS piece; just ensure the call path is exercised
		_ = piece
	}
}

// Quick coverage for helpers that return immediately
func (s *GollamaMoreSuite) TestHelpersAndDetect() {
	// These should not error or panic and exercise return paths
	_ = Print_system_info()
	_ = Supports_mmap()
	_ = Supports_mlock()
	_ = Supports_gpu_offload()
	_ = Max_devices()
	// DetectGpuBackend should return a valid enum
	b := DetectGpuBackend()
	assert.GreaterOrEqual(s.T(), int(b), int(LLAMA_GPU_BACKEND_NONE))
}

// Cover the alternate default helpers
func (s *GollamaMoreSuite) TestAlternateDefaultHelpers() {
	md := ModelDefaultParams()
	_ = md
	cd := ContextDefaultParams()
	_ = cd
	sd := SamplerChainDefaultParams()
	_ = sd
}

func TestGollamaMoreSuite(t *testing.T) { suite.Run(t, new(GollamaMoreSuite)) }
