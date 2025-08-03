// Package gollama provides Go bindings for llama.cpp using purego.
// This package allows you to use llama.cpp functionality from Go without CGO.
//
// The bindings are designed to be as close to the original llama.cpp C API as possible,
// while providing Go-friendly interfaces where appropriate.
//
// Example usage:
//
//	// Initialize the library
//	gollama.Backend_init()
//	defer gollama.Backend_free()
//
//	// Load a model
//	params := gollama.Model_default_params()
//	model, err := gollama.Model_load_from_file("model.gguf", params)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer gollama.Model_free(model)
//
//	// Create context and generate text
//	ctxParams := gollama.Context_default_params()
//	ctx, err := gollama.Init_from_model(model, ctxParams)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer gollama.Free(ctx)
package gollama

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// Version information
const (
	// Version is the gollama.cpp version
	Version = "1.0.0"
	// LlamaCppBuild is the llama.cpp build number this version is based on
	LlamaCppBuild = "b6076"
	// FullVersion combines both version numbers
	FullVersion = "v" + Version + "-llamacpp." + LlamaCppBuild
)

// Platform-specific library names
var libNames = map[string]map[string]string{
	"darwin": {
		"amd64": "libllama.dylib",
		"arm64": "libllama.dylib",
	},
	"linux": {
		"amd64": "libllama.so",
		"arm64": "libllama.so",
	},
	"windows": {
		"amd64": "llama.dll",
		"arm64": "llama.dll",
	},
}

// Global library handle
var (
	libHandle uintptr
	libMutex  sync.RWMutex
	isLoaded  bool
)

// Common types matching llama.cpp
type (
	LlamaToken  int32
	LlamaPos    int32
	LlamaSeqId  int32
	LlamaMemory uintptr
)

// Constants from llama.h
const (
	LLAMA_DEFAULT_SEED = 0xFFFFFFFF
	LLAMA_TOKEN_NULL   = -1

	// File magic numbers
	LLAMA_FILE_MAGIC_GGLA = 0x67676c61
	LLAMA_FILE_MAGIC_GGSN = 0x6767736e
	LLAMA_FILE_MAGIC_GGSQ = 0x67677371

	// Session constants
	LLAMA_SESSION_MAGIC   = LLAMA_FILE_MAGIC_GGSN
	LLAMA_SESSION_VERSION = 9

	LLAMA_STATE_SEQ_MAGIC   = LLAMA_FILE_MAGIC_GGSQ
	LLAMA_STATE_SEQ_VERSION = 2
)

// Enums
type LlamaVocabType int32

const (
	LLAMA_VOCAB_TYPE_NONE LlamaVocabType = iota
	LLAMA_VOCAB_TYPE_SPM
	LLAMA_VOCAB_TYPE_BPE
	LLAMA_VOCAB_TYPE_WPM
	LLAMA_VOCAB_TYPE_UGM
	LLAMA_VOCAB_TYPE_RWKV
)

type LlamaTokenType int32

const (
	LLAMA_TOKEN_TYPE_UNDEFINED LlamaTokenType = iota
	LLAMA_TOKEN_TYPE_NORMAL
	LLAMA_TOKEN_TYPE_UNKNOWN
	LLAMA_TOKEN_TYPE_CONTROL
	LLAMA_TOKEN_TYPE_USER_DEFINED
	LLAMA_TOKEN_TYPE_UNUSED
	LLAMA_TOKEN_TYPE_BYTE
)

type LlamaTokenAttr int32

const (
	LLAMA_TOKEN_ATTR_UNDEFINED   LlamaTokenAttr = 0
	LLAMA_TOKEN_ATTR_UNKNOWN     LlamaTokenAttr = 1 << 0
	LLAMA_TOKEN_ATTR_UNUSED      LlamaTokenAttr = 1 << 1
	LLAMA_TOKEN_ATTR_NORMAL      LlamaTokenAttr = 1 << 2
	LLAMA_TOKEN_ATTR_CONTROL     LlamaTokenAttr = 1 << 3
	LLAMA_TOKEN_ATTR_USER_DEF    LlamaTokenAttr = 1 << 4
	LLAMA_TOKEN_ATTR_BYTE        LlamaTokenAttr = 1 << 5
	LLAMA_TOKEN_ATTR_LSTRIP      LlamaTokenAttr = 1 << 6
	LLAMA_TOKEN_ATTR_RSTRIP      LlamaTokenAttr = 1 << 7
	LLAMA_TOKEN_ATTR_SINGLE_WORD LlamaTokenAttr = 1 << 8
)

type LlamaFtype int32

const (
	LLAMA_FTYPE_ALL_F32        LlamaFtype = 0
	LLAMA_FTYPE_MOSTLY_F16     LlamaFtype = 1
	LLAMA_FTYPE_MOSTLY_Q4_0    LlamaFtype = 2
	LLAMA_FTYPE_MOSTLY_Q4_1    LlamaFtype = 3
	LLAMA_FTYPE_MOSTLY_Q8_0    LlamaFtype = 7
	LLAMA_FTYPE_MOSTLY_Q5_0    LlamaFtype = 8
	LLAMA_FTYPE_MOSTLY_Q5_1    LlamaFtype = 9
	LLAMA_FTYPE_MOSTLY_Q2_K    LlamaFtype = 10
	LLAMA_FTYPE_MOSTLY_Q3_K_S  LlamaFtype = 11
	LLAMA_FTYPE_MOSTLY_Q3_K_M  LlamaFtype = 12
	LLAMA_FTYPE_MOSTLY_Q3_K_L  LlamaFtype = 13
	LLAMA_FTYPE_MOSTLY_Q4_K_S  LlamaFtype = 14
	LLAMA_FTYPE_MOSTLY_Q4_K_M  LlamaFtype = 15
	LLAMA_FTYPE_MOSTLY_Q5_K_S  LlamaFtype = 16
	LLAMA_FTYPE_MOSTLY_Q5_K_M  LlamaFtype = 17
	LLAMA_FTYPE_MOSTLY_Q6_K    LlamaFtype = 18
	LLAMA_FTYPE_MOSTLY_IQ2_XXS LlamaFtype = 19
	LLAMA_FTYPE_MOSTLY_IQ2_XS  LlamaFtype = 20
	LLAMA_FTYPE_MOSTLY_Q2_K_S  LlamaFtype = 21
	LLAMA_FTYPE_MOSTLY_IQ3_XS  LlamaFtype = 22
)

type LlamaRopeScalingType int32

const (
	LLAMA_ROPE_SCALING_TYPE_UNSPECIFIED LlamaRopeScalingType = -1
	LLAMA_ROPE_SCALING_TYPE_NONE        LlamaRopeScalingType = 0
	LLAMA_ROPE_SCALING_TYPE_LINEAR      LlamaRopeScalingType = 1
	LLAMA_ROPE_SCALING_TYPE_YARN        LlamaRopeScalingType = 2
)

type LlamaPoolingType int32

const (
	LLAMA_POOLING_TYPE_UNSPECIFIED LlamaPoolingType = -1
	LLAMA_POOLING_TYPE_NONE        LlamaPoolingType = 0
	LLAMA_POOLING_TYPE_MEAN        LlamaPoolingType = 1
	LLAMA_POOLING_TYPE_CLS         LlamaPoolingType = 2
	LLAMA_POOLING_TYPE_LAST        LlamaPoolingType = 3
	LLAMA_POOLING_TYPE_RANK        LlamaPoolingType = 4
)

type LlamaAttentionType int32

const (
	LLAMA_ATTENTION_TYPE_CAUSAL     LlamaAttentionType = 0
	LLAMA_ATTENTION_TYPE_NON_CAUSAL LlamaAttentionType = 1
)

type LlamaSplitMode int32

const (
	LLAMA_SPLIT_MODE_NONE  LlamaSplitMode = 0
	LLAMA_SPLIT_MODE_LAYER LlamaSplitMode = 1
	LLAMA_SPLIT_MODE_ROW   LlamaSplitMode = 2
)

// Opaque types (represented as pointers)
type LlamaModel uintptr
type LlamaContext uintptr
type LlamaVocab uintptr
type LlamaSampler uintptr
type LlamaAdapterLora uintptr

// Structs
type LlamaTokenData struct {
	Id    LlamaToken // token id
	Logit float32    // log-odds of the token
	P     float32    // probability of the token
}

type LlamaTokenDataArray struct {
	Data     *LlamaTokenData // pointer to token data array
	Size     uint64          // number of tokens
	Selected int64           // index of selected token (-1 if none)
	Sorted   uint8           // whether the array is sorted by probability (bool as uint8)
}

type LlamaBatch struct {
	NTokens int32        // number of tokens
	Token   *LlamaToken  // tokens
	Embd    *float32     // embeddings (if using embeddings instead of tokens)
	Pos     *LlamaPos    // positions
	NSeqId  *int32       // number of sequence IDs per token
	SeqId   **LlamaSeqId // sequence IDs
	Logits  *int8        // whether to compute logits for each token
}

// Model parameters
type LlamaModelParams struct {
	NGpuLayers               int32          // number of layers to store in VRAM
	SplitMode                LlamaSplitMode // how to split the model across multiple GPUs
	MainGpu                  int32          // the GPU that is used for the entire model
	TensorSplit              *float32       // proportion of the model (layers or rows) to offload to each GPU
	RpcServers               *byte          // comma separated list of RPC servers
	ProgressCallback         uintptr        // progress callback function pointer
	ProgressCallbackUserData uintptr        // user data for progress callback
	KvOverrides              uintptr        // model key-value overrides
	VocabOnly                uint8          // only load the vocabulary, no weights (bool as uint8)
	UseMmap                  uint8          // use mmap if possible (bool as uint8)
	UseMlock                 uint8          // force system to keep model in RAM (bool as uint8)
	CheckTensors             uint8          // validate model tensor data (bool as uint8)
}

// Context parameters
type LlamaContextParams struct {
	Seed              uint32               // RNG seed, -1 for random
	NCtx              uint32               // text context, 0 = from model
	NBatch            uint32               // logical maximum batch size
	NUbatch           uint32               // physical maximum batch size
	NSeqMax           uint32               // max number of sequences
	NThreads          int32                // number of threads to use for generation
	NThreadsBatch     int32                // number of threads to use for batch processing
	RopeScalingType   LlamaRopeScalingType // RoPE scaling type
	PoolingType       LlamaPoolingType     // pooling type for embeddings
	AttentionType     LlamaAttentionType   // attention type
	RopeFreqBase      float32              // RoPE base frequency
	RopeFreqScale     float32              // RoPE frequency scaling factor
	YarnExtFactor     float32              // YaRN extrapolation mix factor
	YarnAttnFactor    float32              // YaRN magnitude scaling factor
	YarnBetaFast      float32              // YaRN low correction dim
	YarnBetaSlow      float32              // YaRN high correction dim
	YarnOrigCtx       uint32               // YaRN original context size
	DefragThold       float32              // defragment the KV cache if holes/size > thold
	CbEval            uintptr              // evaluation callback
	CbEvalUserData    uintptr              // user data for evaluation callback
	TypeK             int32                // data type for K cache
	TypeV             int32                // data type for V cache
	AbortCallback     uintptr              // abort callback
	AbortCallbackData uintptr              // user data for abort callback
	Logits            uint8                // whether to compute and return logits (bool as uint8)
	Embeddings        uint8                // whether to compute and return embeddings (bool as uint8)
	Offload_kqv       uint8                // whether to offload K, Q, V to GPU (bool as uint8)
	FlashAttn         uint8                // whether to use flash attention (bool as uint8)
	NoPerf            uint8                // whether to measure performance (bool as uint8)
}

// Model quantize parameters
type LlamaModelQuantizeParams struct {
	NThread              int32      // number of threads to use for quantizing
	Ftype                LlamaFtype // quantize to this llama_ftype
	OutputTensorType     int32      // output tensor type
	TokenEmbeddingType   int32      // itoken embeddings tensor type
	AllowRequantize      uint8      // allow quantizing non-f32/f16 tensors (bool as uint8)
	QuantizeOutputTensor uint8      // quantize output.weight (bool as uint8)
	OnlyF32              uint8      // quantize only f32 tensors (bool as uint8)
	PureF16              uint8      // disable k-quant mixtures and quantize all tensors to the same type (bool as uint8)
	KeepSplit            uint8      // keep split tensors (bool as uint8)
	IMatrix              *byte      // importance matrix data
	KqsWarning           uint8      // warning for quantization quality loss (bool as uint8)
}

// Chat message
type LlamaChatMessage struct {
	Role    *byte // role string
	Content *byte // content string
}

// Sampler chain parameters
type LlamaSamplerChainParams struct {
	NoPerf uint8 // whether to measure performance timings (bool as uint8)
}

// Logit bias
type LlamaLogitBias struct {
	Token LlamaToken
	Bias  float32
}

// Function pointers for C functions
var (
	// Backend functions
	llamaBackendInit func()
	llamaBackendFree func()
	llamaLogSet      func(logCallback uintptr, userData uintptr)

	// Model functions
	llamaModelDefaultParams  func() LlamaModelParams
	llamaModelLoadFromFile   func(pathModel *byte, params LlamaModelParams) LlamaModel
	llamaModelLoadFromSplits func(paths **byte, nPaths uint64, params LlamaModelParams) LlamaModel
	llamaModelSaveToFile     func(model LlamaModel, pathModel *byte)
	llamaModelFree           func(model LlamaModel)

	// Context functions
	llamaContextDefaultParams func() LlamaContextParams
	llamaInitFromModel        func(model LlamaModel, params LlamaContextParams) LlamaContext
	llamaFree                 func(ctx LlamaContext)

	// Model info functions
	llamaModelNCtxTrain func(model LlamaModel) int32
	llamaModelNEmbd     func(model LlamaModel) int32
	llamaModelNLayer    func(model LlamaModel) int32
	llamaModelNHead     func(model LlamaModel) int32
	llamaModelNHeadKv   func(model LlamaModel) int32
	llamaModelVocabType func(model LlamaModel) LlamaVocabType
	llamaModelRopeType  func(model LlamaModel) int32

	// Context info functions
	llamaNCtx        func(ctx LlamaContext) uint32
	llamaNBatch      func(ctx LlamaContext) uint32
	llamaNUbatch     func(ctx LlamaContext) uint32
	llamaNSeqMax     func(ctx LlamaContext) uint32
	llamaPoolingType func(ctx LlamaContext) LlamaPoolingType
	llamaGetModel    func(ctx LlamaContext) LlamaModel

	// Tokenization functions
	llamaTokenize     func(model LlamaModel, text *byte, textLen int32, tokens *LlamaToken, nTokensMax int32, addSpecial uint8, parseSpecial uint8) int32
	llamaTokenToPiece func(model LlamaModel, token LlamaToken, buf *byte, length int32, lstrip uint8, special uint8) int32
	llamaDetokenize   func(model LlamaModel, tokens *LlamaToken, nTokens int32, text *byte, textLen int32, removeSpecial uint8, unparseSpecial uint8) int32

	// Vocab functions
	llamaModelGetVocab func(model LlamaModel) LlamaVocab
	llamaVocabNTokens  func(vocab LlamaVocab) int32
	llamaVocabBos      func(vocab LlamaVocab) LlamaToken
	llamaVocabEos      func(vocab LlamaVocab) LlamaToken
	llamaVocabEot      func(vocab LlamaVocab) LlamaToken
	llamaVocabNl       func(vocab LlamaVocab) LlamaToken
	llamaVocabPad      func(vocab LlamaVocab) LlamaToken

	// Batch functions
	llamaBatchInit func(nTokens int32, embd int32, nSeqMax int32) LlamaBatch
	llamaBatchFree func(batch LlamaBatch)
	llamaBatchGet1 func(tokens *LlamaToken, nTokens int32, pos0 LlamaPos, seqId LlamaSeqId) LlamaBatch

	// Decode functions
	llamaDecode func(ctx LlamaContext, batch LlamaBatch) int32
	llamaEncode func(ctx LlamaContext, batch LlamaBatch) int32

	// Logits and embeddings
	llamaGetLogits     func(ctx LlamaContext) *float32
	llamaGetLogitsIth  func(ctx LlamaContext, i int32) *float32
	llamaGetEmbeddings func(ctx LlamaContext) *float32

	// Sampling functions
	llamaSamplerChainDefaultParams func() LlamaSamplerChainParams
	llamaSamplerChainInit          func(params LlamaSamplerChainParams) LlamaSampler
	llamaSamplerChainAdd           func(chain LlamaSampler, smpl LlamaSampler)
	llamaSamplerChainGet           func(chain LlamaSampler, i int32) LlamaSampler
	llamaSamplerChainN             func(chain LlamaSampler) int32
	llamaSamplerChainFree          func(chain LlamaSampler)
	llamaSamplerSample             func(smpl LlamaSampler, ctx LlamaContext, candidates *LlamaTokenDataArray) LlamaToken
	llamaSamplerAccept             func(smpl LlamaSampler, token LlamaToken)
	llamaSamplerReset              func(smpl LlamaSampler)

	// Built-in samplers
	llamaSamplerInitGreedy  func() LlamaSampler
	llamaSamplerInitDist    func(seed uint32) LlamaSampler
	llamaSamplerInitSoftmax func() LlamaSampler
	llamaSamplerInitTopK    func(k int32) LlamaSampler
	llamaSamplerInitTopP    func(p float32, minKeep uint64) LlamaSampler
	llamaSamplerInitMinP    func(p float32, minKeep uint64) LlamaSampler
	// llamaSamplerInitTailFree   func(z float32, minKeep uint64) LlamaSampler  // Function doesn't exist
	llamaSamplerInitTypical    func(p float32, minKeep uint64) LlamaSampler
	llamaSamplerInitTemp       func(temp float32) LlamaSampler
	llamaSamplerInitTempExt    func(temp float32, delta float32, exponent float32) LlamaSampler
	llamaSamplerInitMirostat   func(tau float32, eta float32, m int32, seed uint32) LlamaSampler
	llamaSamplerInitMirostatV2 func(tau float32, eta float32, seed uint32) LlamaSampler

	// Utility functions
	llamaMaxDevices         func() uint64
	llamaSupportsMmap       func() bool
	llamaSupportsMlock      func() bool
	llamaSupportsGpuOffload func() bool
	llamaSupportsRpc        func() bool
	llamaTimeUs             func() int64
	llamaPrintSystemInfo    func() *byte

	// KV cache functions
	llamaKvCacheClear   func(ctx LlamaContext)
	llamaKvCacheSeqRm   func(ctx LlamaContext, seqId LlamaSeqId, p0 LlamaPos, p1 LlamaPos) bool
	llamaKvCacheSeqCp   func(ctx LlamaContext, seqIdSrc LlamaSeqId, seqIdDst LlamaSeqId, p0 LlamaPos, p1 LlamaPos)
	llamaKvCacheSeqKeep func(ctx LlamaContext, seqId LlamaSeqId)
	llamaKvCacheSeqAdd  func(ctx LlamaContext, seqId LlamaSeqId, p0 LlamaPos, p1 LlamaPos, delta LlamaPos)
	llamaKvCacheSeqDiv  func(ctx LlamaContext, seqId LlamaSeqId, p0 LlamaPos, p1 LlamaPos, d int32)
	// llamaKvCacheSeqPos  func(ctx LlamaContext, seqId LlamaSeqId, p0 LlamaPos, p1 LlamaPos, delta LlamaPos)  // Function doesn't exist
	llamaKvCacheDefrag func(ctx LlamaContext)
	llamaKvCacheUpdate func(ctx LlamaContext)

	// State functions
	llamaStateGetSize  func(ctx LlamaContext) uint64
	llamaStateGetData  func(ctx LlamaContext, dst *byte, size uint64) uint64
	llamaStateSetData  func(ctx LlamaContext, src *byte, size uint64) uint64
	llamaStateLoadFile func(ctx LlamaContext, pathSession *byte, tokensOut *LlamaToken, nTokenCapacity uint64, nTokenCountOut *uint64) bool
	llamaStateSaveFile func(ctx LlamaContext, pathSession *byte, tokens *LlamaToken, nTokenCount uint64) bool

	// Performance functions - These may not exist in this llama.cpp version
	// llamaGetTimings   func(ctx LlamaContext) uintptr
	// llamaPrintTimings func(ctx LlamaContext)
	// llamaResetTimings func(ctx LlamaContext)
)

// Library loading and initialization
func getLibraryPath() (string, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	archMap, ok := libNames[goos]
	if !ok {
		return "", fmt.Errorf("unsupported OS: %s", goos)
	}

	libName, ok := archMap[goarch]
	if !ok {
		return "", fmt.Errorf("unsupported architecture: %s on %s", goarch, goos)
	}

	// For now, assume the library is in the same directory or system path
	// In a real implementation, you would embed or package the libraries
	return libName, nil
}

// loadLibrary loads the llama.cpp shared library
func loadLibrary() error {
	libMutex.Lock()
	defer libMutex.Unlock()

	if isLoaded {
		return nil
	}

	libPath, err := getLibraryPath()
	if err != nil {
		return fmt.Errorf("failed to get library path: %w", err)
	}

	var handle uintptr
	if runtime.GOOS == "windows" {
		// On Windows, use LoadLibrary via syscall
		// This would need to be implemented with proper Windows API calls
		return errors.New("support for windows platform not yet implemented")
	} else {
		// On Unix-like systems, use purego's Dlopen
		handle, err = purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			return fmt.Errorf("failed to load library %s: %w", libPath, err)
		}
	}

	libHandle = handle

	// Register all function pointers
	if err := registerFunctions(); err != nil {
		purego.Dlclose(handle)
		return fmt.Errorf("failed to register functions: %w", err)
	}

	isLoaded = true
	return nil
}

// registerFunctions registers all llama.cpp function pointers
func registerFunctions() error {
	// Backend functions
	purego.RegisterLibFunc(&llamaBackendInit, libHandle, "llama_backend_init")
	purego.RegisterLibFunc(&llamaBackendFree, libHandle, "llama_backend_free")
	purego.RegisterLibFunc(&llamaLogSet, libHandle, "llama_log_set")

	// Model functions
	purego.RegisterLibFunc(&llamaModelDefaultParams, libHandle, "llama_model_default_params")
	purego.RegisterLibFunc(&llamaModelLoadFromFile, libHandle, "llama_model_load_from_file")
	purego.RegisterLibFunc(&llamaModelLoadFromSplits, libHandle, "llama_model_load_from_splits")
	purego.RegisterLibFunc(&llamaModelSaveToFile, libHandle, "llama_model_save_to_file")
	purego.RegisterLibFunc(&llamaModelFree, libHandle, "llama_model_free")

	// Context functions
	purego.RegisterLibFunc(&llamaContextDefaultParams, libHandle, "llama_context_default_params")
	purego.RegisterLibFunc(&llamaInitFromModel, libHandle, "llama_init_from_model")
	purego.RegisterLibFunc(&llamaFree, libHandle, "llama_free")

	// Model info functions
	purego.RegisterLibFunc(&llamaModelNCtxTrain, libHandle, "llama_model_n_ctx_train")
	purego.RegisterLibFunc(&llamaModelNEmbd, libHandle, "llama_model_n_embd")
	purego.RegisterLibFunc(&llamaModelNLayer, libHandle, "llama_model_n_layer")
	purego.RegisterLibFunc(&llamaModelNHead, libHandle, "llama_model_n_head")
	purego.RegisterLibFunc(&llamaModelNHeadKv, libHandle, "llama_model_n_head_kv")
	purego.RegisterLibFunc(&llamaModelVocabType, libHandle, "llama_vocab_type")
	purego.RegisterLibFunc(&llamaModelRopeType, libHandle, "llama_model_rope_type")

	// Context info functions
	purego.RegisterLibFunc(&llamaNCtx, libHandle, "llama_n_ctx")
	purego.RegisterLibFunc(&llamaNBatch, libHandle, "llama_n_batch")
	purego.RegisterLibFunc(&llamaNUbatch, libHandle, "llama_n_ubatch")
	purego.RegisterLibFunc(&llamaNSeqMax, libHandle, "llama_n_seq_max")
	purego.RegisterLibFunc(&llamaPoolingType, libHandle, "llama_pooling_type")
	purego.RegisterLibFunc(&llamaGetModel, libHandle, "llama_get_model")

	// Tokenization functions
	purego.RegisterLibFunc(&llamaTokenize, libHandle, "llama_tokenize")
	purego.RegisterLibFunc(&llamaTokenToPiece, libHandle, "llama_token_to_piece")
	purego.RegisterLibFunc(&llamaDetokenize, libHandle, "llama_detokenize")

	// Vocab functions
	purego.RegisterLibFunc(&llamaModelGetVocab, libHandle, "llama_model_get_vocab")
	purego.RegisterLibFunc(&llamaVocabNTokens, libHandle, "llama_vocab_n_tokens")
	purego.RegisterLibFunc(&llamaVocabBos, libHandle, "llama_vocab_bos")
	purego.RegisterLibFunc(&llamaVocabEos, libHandle, "llama_vocab_eos")
	purego.RegisterLibFunc(&llamaVocabEot, libHandle, "llama_vocab_eot")
	purego.RegisterLibFunc(&llamaVocabNl, libHandle, "llama_vocab_nl")
	purego.RegisterLibFunc(&llamaVocabPad, libHandle, "llama_vocab_pad")

	// Batch functions
	purego.RegisterLibFunc(&llamaBatchInit, libHandle, "llama_batch_init")
	purego.RegisterLibFunc(&llamaBatchFree, libHandle, "llama_batch_free")
	purego.RegisterLibFunc(&llamaBatchGet1, libHandle, "llama_batch_get_one")

	// Decode functions
	purego.RegisterLibFunc(&llamaDecode, libHandle, "llama_decode")
	purego.RegisterLibFunc(&llamaEncode, libHandle, "llama_encode")

	// Logits and embeddings
	purego.RegisterLibFunc(&llamaGetLogits, libHandle, "llama_get_logits")
	purego.RegisterLibFunc(&llamaGetLogitsIth, libHandle, "llama_get_logits_ith")
	purego.RegisterLibFunc(&llamaGetEmbeddings, libHandle, "llama_get_embeddings")

	// Sampling functions
	purego.RegisterLibFunc(&llamaSamplerChainDefaultParams, libHandle, "llama_sampler_chain_default_params")
	purego.RegisterLibFunc(&llamaSamplerChainInit, libHandle, "llama_sampler_chain_init")
	purego.RegisterLibFunc(&llamaSamplerChainAdd, libHandle, "llama_sampler_chain_add")
	purego.RegisterLibFunc(&llamaSamplerChainGet, libHandle, "llama_sampler_chain_get")
	purego.RegisterLibFunc(&llamaSamplerChainN, libHandle, "llama_sampler_chain_n")
	purego.RegisterLibFunc(&llamaSamplerChainFree, libHandle, "llama_sampler_free")
	purego.RegisterLibFunc(&llamaSamplerSample, libHandle, "llama_sampler_sample")
	purego.RegisterLibFunc(&llamaSamplerAccept, libHandle, "llama_sampler_accept")
	purego.RegisterLibFunc(&llamaSamplerReset, libHandle, "llama_sampler_reset")

	// Built-in samplers
	purego.RegisterLibFunc(&llamaSamplerInitGreedy, libHandle, "llama_sampler_init_greedy")
	purego.RegisterLibFunc(&llamaSamplerInitDist, libHandle, "llama_sampler_init_dist")
	purego.RegisterLibFunc(&llamaSamplerInitSoftmax, libHandle, "llama_sampler_init_softmax")
	purego.RegisterLibFunc(&llamaSamplerInitTopK, libHandle, "llama_sampler_init_top_k")
	purego.RegisterLibFunc(&llamaSamplerInitTopP, libHandle, "llama_sampler_init_top_p")
	purego.RegisterLibFunc(&llamaSamplerInitMinP, libHandle, "llama_sampler_init_min_p")
	// purego.RegisterLibFunc(&llamaSamplerInitTailFree, libHandle, "llama_sampler_init_tail_free")  // Function doesn't exist
	purego.RegisterLibFunc(&llamaSamplerInitTypical, libHandle, "llama_sampler_init_typical")
	purego.RegisterLibFunc(&llamaSamplerInitTemp, libHandle, "llama_sampler_init_temp")
	purego.RegisterLibFunc(&llamaSamplerInitTempExt, libHandle, "llama_sampler_init_temp_ext")
	purego.RegisterLibFunc(&llamaSamplerInitMirostat, libHandle, "llama_sampler_init_mirostat")
	purego.RegisterLibFunc(&llamaSamplerInitMirostatV2, libHandle, "llama_sampler_init_mirostat_v2")

	// Utility functions
	purego.RegisterLibFunc(&llamaMaxDevices, libHandle, "llama_max_devices")
	purego.RegisterLibFunc(&llamaSupportsMmap, libHandle, "llama_supports_mmap")
	purego.RegisterLibFunc(&llamaSupportsMlock, libHandle, "llama_supports_mlock")
	purego.RegisterLibFunc(&llamaSupportsGpuOffload, libHandle, "llama_supports_gpu_offload")
	purego.RegisterLibFunc(&llamaSupportsRpc, libHandle, "llama_supports_rpc")
	purego.RegisterLibFunc(&llamaTimeUs, libHandle, "llama_time_us")
	purego.RegisterLibFunc(&llamaPrintSystemInfo, libHandle, "llama_print_system_info")

	// KV cache functions
	purego.RegisterLibFunc(&llamaKvCacheClear, libHandle, "llama_kv_self_clear")
	purego.RegisterLibFunc(&llamaKvCacheSeqRm, libHandle, "llama_kv_self_seq_rm")
	purego.RegisterLibFunc(&llamaKvCacheSeqCp, libHandle, "llama_kv_self_seq_cp")
	purego.RegisterLibFunc(&llamaKvCacheSeqKeep, libHandle, "llama_kv_self_seq_keep")
	purego.RegisterLibFunc(&llamaKvCacheSeqAdd, libHandle, "llama_kv_self_seq_add")
	purego.RegisterLibFunc(&llamaKvCacheSeqDiv, libHandle, "llama_kv_self_seq_div")
	// purego.RegisterLibFunc(&llamaKvCacheSeqPos, libHandle, "llama_kv_self_seq_pos")  // Might not exist
	purego.RegisterLibFunc(&llamaKvCacheDefrag, libHandle, "llama_kv_self_defrag")
	purego.RegisterLibFunc(&llamaKvCacheUpdate, libHandle, "llama_kv_self_update")

	// State functions
	purego.RegisterLibFunc(&llamaStateGetSize, libHandle, "llama_state_get_size")
	purego.RegisterLibFunc(&llamaStateGetData, libHandle, "llama_state_get_data")
	purego.RegisterLibFunc(&llamaStateSetData, libHandle, "llama_state_set_data")
	purego.RegisterLibFunc(&llamaStateLoadFile, libHandle, "llama_state_load_file")
	purego.RegisterLibFunc(&llamaStateSaveFile, libHandle, "llama_state_save_file")

	// Performance functions - These may not exist in this llama.cpp version
	// purego.RegisterLibFunc(&llamaGetTimings, libHandle, "llama_get_timings")
	// purego.RegisterLibFunc(&llamaPrintTimings, libHandle, "llama_print_timings")
	// purego.RegisterLibFunc(&llamaResetTimings, libHandle, "llama_reset_timings")

	return nil
}

// ensureLoaded ensures the library is loaded before calling any functions
func ensureLoaded() error {
	libMutex.RLock()
	if isLoaded {
		libMutex.RUnlock()
		return nil
	}
	libMutex.RUnlock()

	return loadLibrary()
}

// Public API functions

// Backend_init initializes the llama + ggml backend
func Backend_init() error {
	if err := ensureLoaded(); err != nil {
		return err
	}
	llamaBackendInit()
	return nil
}

// Backend_free frees the llama + ggml backend
func Backend_free() {
	if isLoaded {
		llamaBackendFree()
	}
}

// Model_default_params returns default model parameters
func Model_default_params() LlamaModelParams {
	if err := ensureLoaded(); err != nil {
		panic(err) // In a real implementation, handle this better
	}
	return llamaModelDefaultParams()
}

// Model_load_from_file loads a model from a file
func Model_load_from_file(pathModel string, params LlamaModelParams) (LlamaModel, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}

	pathBytes := append([]byte(pathModel), 0) // null-terminate
	model := llamaModelLoadFromFile((*byte)(unsafe.Pointer(&pathBytes[0])), params)
	if model == 0 {
		return 0, errors.New("failed to load model")
	}
	return model, nil
}

// Model_free frees a model
func Model_free(model LlamaModel) {
	if isLoaded && model != 0 {
		llamaModelFree(model)
	}
}

// Context_default_params returns default context parameters
func Context_default_params() LlamaContextParams {
	if err := ensureLoaded(); err != nil {
		panic(err)
	}
	return llamaContextDefaultParams()
}

// Init_from_model creates a context from a model
func Init_from_model(model LlamaModel, params LlamaContextParams) (LlamaContext, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}

	ctx := llamaInitFromModel(model, params)
	if ctx == 0 {
		return 0, errors.New("failed to create context")
	}
	return ctx, nil
}

// Free frees a context
func Free(ctx LlamaContext) {
	if isLoaded && ctx != 0 {
		llamaFree(ctx)
	}
}

// Helper function to convert bool to uint8 for C interop
func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// Tokenize tokenizes text
func Tokenize(model LlamaModel, text string, addSpecial, parseSpecial bool) ([]LlamaToken, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}

	textBytes := append([]byte(text), 0) // null-terminate

	// First call to get the number of tokens
	nTokens := llamaTokenize(model, (*byte)(unsafe.Pointer(&textBytes[0])), int32(len(text)), nil, 0, boolToUint8(addSpecial), boolToUint8(parseSpecial))
	if nTokens < 0 {
		return nil, errors.New("tokenization failed")
	}

	if nTokens == 0 {
		return []LlamaToken{}, nil
	}

	// Second call to get the actual tokens
	tokens := make([]LlamaToken, nTokens)
	result := llamaTokenize(model, (*byte)(unsafe.Pointer(&textBytes[0])), int32(len(text)), &tokens[0], nTokens, boolToUint8(addSpecial), boolToUint8(parseSpecial))
	if result < 0 {
		return nil, errors.New("tokenization failed")
	}

	return tokens[:result], nil
}

// Token_to_piece converts a token to its string representation
func Token_to_piece(model LlamaModel, token LlamaToken, special bool) string {
	if err := ensureLoaded(); err != nil {
		return ""
	}

	// First call to get buffer size
	bufSize := llamaTokenToPiece(model, token, nil, 0, boolToUint8(false), boolToUint8(special))
	if bufSize <= 0 {
		return ""
	}

	// Second call to get actual string
	buf := make([]byte, bufSize)
	result := llamaTokenToPiece(model, token, &buf[0], bufSize, boolToUint8(false), boolToUint8(special))
	if result <= 0 {
		return ""
	}

	// Find the null terminator and return the string
	for i, b := range buf {
		if b == 0 {
			return string(buf[:i])
		}
	}
	return string(buf)
}

// Batch_init creates a new batch
func Batch_init(nTokens, embd, nSeqMax int32) LlamaBatch {
	if err := ensureLoaded(); err != nil {
		panic(err)
	}
	return llamaBatchInit(nTokens, embd, nSeqMax)
}

// Batch_free frees a batch
func Batch_free(batch LlamaBatch) {
	if isLoaded {
		llamaBatchFree(batch)
	}
}

// Batch_add adds a token to a batch (helper function)
func Batch_add(batch LlamaBatch, token LlamaToken, pos LlamaPos, seqIds []LlamaSeqId, logits bool) {
	// This is a helper function - actual implementation would manipulate the batch struct
	// For now, this is a placeholder
}

// Decode decodes a batch
func Decode(ctx LlamaContext, batch LlamaBatch) error {
	if err := ensureLoaded(); err != nil {
		return err
	}

	result := llamaDecode(ctx, batch)
	if result != 0 {
		return fmt.Errorf("decode failed with code %d", result)
	}
	return nil
}

// Get_logits_ith gets logits for a specific token
func Get_logits_ith(ctx LlamaContext, i int32) *float32 {
	if err := ensureLoaded(); err != nil {
		return nil
	}
	return llamaGetLogitsIth(ctx, i)
}

// Token_data_array_init creates a token data array (helper function)
func Token_data_array_init(model LlamaModel) *LlamaTokenDataArray {
	// This would be implemented to create a proper token data array
	// For now, this is a placeholder
	return nil
}

// Sampler_init_greedy creates a greedy sampler
func Sampler_init_greedy() LlamaSampler {
	if err := ensureLoaded(); err != nil {
		panic(err)
	}
	return llamaSamplerInitGreedy()
}

// Sampler_free frees a sampler
func Sampler_free(sampler LlamaSampler) {
	// The C library doesn't seem to have a direct sampler free function
	// This might be handled by the sampler chain
}

// Sampler_sample samples a token
func Sampler_sample(sampler LlamaSampler, ctx LlamaContext, candidates *LlamaTokenDataArray) LlamaToken {
	if err := ensureLoaded(); err != nil {
		return LLAMA_TOKEN_NULL
	}
	return llamaSamplerSample(sampler, ctx, candidates)
}

// Additional utility functions

// Print_system_info prints system information
func Print_system_info() string {
	if err := ensureLoaded(); err != nil {
		return ""
	}

	ptr := llamaPrintSystemInfo()
	if ptr == nil {
		return ""
	}

	// Convert C string to Go string
	// This is unsafe and needs proper implementation
	return ""
}

// Supports_mmap returns whether mmap is supported
func Supports_mmap() bool {
	if err := ensureLoaded(); err != nil {
		return false
	}
	return llamaSupportsMmap()
}

// Supports_mlock returns whether mlock is supported
func Supports_mlock() bool {
	if err := ensureLoaded(); err != nil {
		return false
	}
	return llamaSupportsMlock()
}

// Supports_gpu_offload returns whether GPU offload is supported
func Supports_gpu_offload() bool {
	if err := ensureLoaded(); err != nil {
		return false
	}
	return llamaSupportsGpuOffload()
}

// Max_devices returns the maximum number of devices
func Max_devices() uint64 {
	if err := ensureLoaded(); err != nil {
		return 0
	}
	return llamaMaxDevices()
}
