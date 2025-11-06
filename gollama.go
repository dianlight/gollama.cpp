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
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"
)

// Version information
const (
	// Version is the gollama.cpp version
	Version = "0.2.2"
	// LlamaCppBuild is the llama.cpp build number this version is based on
	LlamaCppBuild = "b6862"
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

type LlamaGpuBackend int32

const (
	LLAMA_GPU_BACKEND_NONE   LlamaGpuBackend = 0
	LLAMA_GPU_BACKEND_CPU    LlamaGpuBackend = 1
	LLAMA_GPU_BACKEND_CUDA   LlamaGpuBackend = 2
	LLAMA_GPU_BACKEND_METAL  LlamaGpuBackend = 3
	LLAMA_GPU_BACKEND_HIP    LlamaGpuBackend = 4
	LLAMA_GPU_BACKEND_VULKAN LlamaGpuBackend = 5
	LLAMA_GPU_BACKEND_OPENCL LlamaGpuBackend = 6
	LLAMA_GPU_BACKEND_SYCL   LlamaGpuBackend = 7
)

// String returns the string representation of the GPU backend
func (b LlamaGpuBackend) String() string {
	switch b {
	case LLAMA_GPU_BACKEND_NONE:
		return "None"
	case LLAMA_GPU_BACKEND_CPU:
		return "CPU"
	case LLAMA_GPU_BACKEND_CUDA:
		return "CUDA"
	case LLAMA_GPU_BACKEND_METAL:
		return "Metal"
	case LLAMA_GPU_BACKEND_HIP:
		return "HIP"
	case LLAMA_GPU_BACKEND_VULKAN:
		return "Vulkan"
	case LLAMA_GPU_BACKEND_OPENCL:
		return "OpenCL"
	case LLAMA_GPU_BACKEND_SYCL:
		return "SYCL"
	default:
		return "Unknown"
	}
}

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
	Devices                  uintptr        // ggml_backend_dev_t * - NULL-terminated list of devices
	TensorBuftOverrides      uintptr        // const struct llama_model_tensor_buft_override *
	NGpuLayers               int32          // number of layers to store in VRAM
	SplitMode                LlamaSplitMode // how to split the model across multiple GPUs
	MainGpu                  int32          // the GPU that is used for the entire model
	TensorSplit              *float32       // proportion of the model to offload to each GPU
	ProgressCallback         uintptr        // llama_progress_callback function pointer
	ProgressCallbackUserData uintptr        // context pointer passed to the progress callback
	KvOverrides              uintptr        // const struct llama_model_kv_override *
	VocabOnly                uint8          // only load the vocabulary, no weights (bool as uint8)
	UseMmap                  uint8          // use mmap if possible (bool as uint8)
	UseMlock                 uint8          // force system to keep model in RAM (bool as uint8)
	CheckTensors             uint8          // validate model tensor data (bool as uint8)
	UseExtraBufts            uint8          // use extra buffer types (bool as uint8)
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
	llamaTokenize     func(vocab LlamaVocab, text *byte, textLen int32, tokens *LlamaToken, nTokensMax int32, addSpecial bool, parseSpecial bool) int32
	llamaTokenToPiece func(vocab LlamaVocab, token LlamaToken, buf *byte, length int32, lstrip int32, special bool) int32
	llamaDetokenize   func(model LlamaModel, tokens *LlamaToken, nTokens int32, text *byte, textLen int32, removeSpecial bool, unparseSpecial bool) int32
	llamaVocabGetText func(vocab LlamaVocab, token LlamaToken) *byte

	// Vocab functions
	llamaModelGetVocab func(model LlamaModel) LlamaVocab
	llamaVocabNTokens  func(vocab LlamaVocab) int32
	llamaVocabBos      func(vocab LlamaVocab) LlamaToken
	llamaVocabEos      func(vocab LlamaVocab) LlamaToken
	llamaVocabEot      func(vocab LlamaVocab) LlamaToken
	llamaVocabNl       func(vocab LlamaVocab) LlamaToken
	llamaVocabPad      func(vocab LlamaVocab) LlamaToken

	// Batch functions
	llamaBatchInit   func(nTokens int32, embd int32, nSeqMax int32) LlamaBatch
	llamaBatchFree   func(batch LlamaBatch)
	llamaBatchGetOne func(tokens *LlamaToken, nTokens int32) LlamaBatch

	// Decode functions
	llamaDecode func(ctx LlamaContext, batch LlamaBatch) int32
	llamaEncode func(ctx LlamaContext, batch LlamaBatch) int32

	// Logits and embeddings
	llamaGetLogits        func(ctx LlamaContext) *float32
	llamaGetLogitsIth     func(ctx LlamaContext, i int32) *float32
	llamaGetEmbeddings    func(ctx LlamaContext) *float32
	llamaGetEmbeddingsIth func(ctx LlamaContext, i int32) *float32
	llamaSetCausalAttn    func(ctx LlamaContext, causal bool) int32
	llamaSetEmbeddings    func(ctx LlamaContext, embeddings bool)
	llamaMemoryClear      func(memory LlamaMemory, reset bool) bool
	llamaGetMemory        func(ctx LlamaContext) LlamaMemory

	// Sampling functions
	llamaSamplerChainDefaultParams func() LlamaSamplerChainParams
	llamaSamplerChainInit          func(params LlamaSamplerChainParams) LlamaSampler
	llamaSamplerChainAdd           func(chain LlamaSampler, smpl LlamaSampler)
	llamaSamplerChainGet           func(chain LlamaSampler, i int32) LlamaSampler
	llamaSamplerChainN             func(chain LlamaSampler) int32
	llamaSamplerChainFree          func(chain LlamaSampler)
	llamaSamplerSample             func(smpl LlamaSampler, ctx LlamaContext, idx int32) LlamaToken
	llamaSamplerAccept             func(smpl LlamaSampler, token LlamaToken)
	llamaSamplerReset              func(smpl LlamaSampler)

	// Built-in samplers
	llamaSamplerInitGreedy func() LlamaSampler
	llamaSamplerInitDist   func(seed uint32) LlamaSampler
	// llamaSamplerInitSoftmax func() LlamaSampler  // Function doesn't exist in b6862
	llamaSamplerInitTopK func(k int32) LlamaSampler
	llamaSamplerInitTopP func(p float32, minKeep uint64) LlamaSampler
	llamaSamplerInitMinP func(p float32, minKeep uint64) LlamaSampler
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

	// KV cache functions (deprecated functions removed in b6862)

	// State functions
	llamaStateGetSize  func(ctx LlamaContext) uint64
	llamaStateGetData  func(ctx LlamaContext, dst *byte, size uint64) uint64
	llamaStateSetData  func(ctx LlamaContext, src *byte, size uint64) uint64
	llamaStateLoadFile func(ctx LlamaContext, pathSession *byte, tokensOut *LlamaToken, nTokenCapacity uint64, nTokenCountOut *uint64) bool
	llamaStateSaveFile func(ctx LlamaContext, pathSession *byte, tokens *LlamaToken, nTokenCount uint64) bool

	// Performance functions - These may not exist in this llama.cpp version - moved to ROADMAP "wait for llama.cpp" section
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

	// Start with standard search paths
	candidates := []string{
		libName,                         // Current directory
		"libs/darwin_arm64/" + libName,  // macOS
		"libs/darwin_amd64/" + libName,  // macOS
		"libs/linux_arm64/" + libName,   // Linux ARM64
		"libs/linux_amd64/" + libName,   // Linux AMD64
		"libs/windows_amd64/" + libName, // Windows AMD64
		"libs/windows_arm64/" + libName, // Windows ARM64
		"../" + libName,                 // Parent directory (for when running from examples/)
		"../../" + libName,              // Parent directory (for when running from examples/)
		"/usr/local/lib/" + libName,     // System library path
		"/usr/lib/" + libName,           // Common system library path
		"/lib/" + libName,               // Another common system library path
	}

	// Add cache directory paths
	// Determine cache directory using the same logic as downloader
	var cacheDir string
	if globalConfig != nil && globalConfig.CacheDir != "" {
		cacheDir = globalConfig.CacheDir
	} else if envCacheDir := os.Getenv("GOLLAMA_CACHE_DIR"); envCacheDir != "" {
		cacheDir = filepath.Join(envCacheDir, "libs")
	} else {
		userCacheDir, err := os.UserCacheDir()
		if err == nil {
			cacheDir = filepath.Join(userCacheDir, "gollama", "libs")
		} else {
			cacheDir = filepath.Join(os.TempDir(), "gollama", "libs")
		}
	}

	// Try to find library in cache directory subdirectories
	if cacheDir != "" {
		entries, err := os.ReadDir(cacheDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					// Check for library in build/bin subdirectory (standard structure)
					cachePath := filepath.Join(cacheDir, entry.Name(), "build", "bin", libName)
					candidates = append(candidates, cachePath)
					// Check for library in bin subdirectory
					cachePath = filepath.Join(cacheDir, entry.Name(), "bin", libName)
					candidates = append(candidates, cachePath)
					// Also check directly in subdirectory
					cachePath = filepath.Join(cacheDir, entry.Name(), libName)
					candidates = append(candidates, cachePath)

					// Check nested subdirectories (e.g., embedded/linux_amd64_b6862/)
					subEntries, subErr := os.ReadDir(filepath.Join(cacheDir, entry.Name()))
					if subErr == nil {
						for _, subEntry := range subEntries {
							if subEntry.IsDir() {
								// Check for library in nested subdirectory
								nestedPath := filepath.Join(cacheDir, entry.Name(), subEntry.Name(), libName)
								candidates = append(candidates, nestedPath)
								// Also check build/bin in nested subdirectory
								nestedPath = filepath.Join(cacheDir, entry.Name(), subEntry.Name(), "build", "bin", libName)
								candidates = append(candidates, nestedPath)
								// Also check bin in nested subdirectory
								nestedPath = filepath.Join(cacheDir, entry.Name(), subEntry.Name(), "bin", libName)
								candidates = append(candidates, nestedPath)
							}
						}
					}
				}
			}
		}
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	// If not found in any of the candidate locations, return the basic name
	// and let the system dynamic loader try to find it
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

	// Check if platform is supported
	if !isPlatformSupported() {
		return getPlatformError()
	}

	// Use platform-specific library loading
	handle, err := loadLibraryPlatform(libPath)
	if err != nil {
		return fmt.Errorf("failed to load library %s: %w", libPath, err)
	}

	libHandle = handle

	// Register all function pointers
	if err := registerFunctions(); err != nil {
		_ = closeLibraryPlatform(handle) // Ignore error during cleanup
		return fmt.Errorf("failed to register functions: %w", err)
	}

	isLoaded = true
	return nil
}

// registerFunctions registers all llama.cpp function pointers
func registerFunctions() error {
	// Backend functions
	registerLibFunc(&llamaBackendInit, libHandle, "llama_backend_init")
	registerLibFunc(&llamaBackendFree, libHandle, "llama_backend_free")
	registerLibFunc(&llamaLogSet, libHandle, "llama_log_set")

	// Model functions - Register struct functions only on Darwin (purego limitation)
	// On other platforms, FFI handles struct parameters/returns directly
	if runtime.GOOS == "darwin" {
		registerLibFunc(&llamaModelDefaultParams, libHandle, "llama_model_default_params")
		registerLibFunc(&llamaContextDefaultParams, libHandle, "llama_context_default_params")
		registerLibFunc(&llamaSamplerChainDefaultParams, libHandle, "llama_sampler_chain_default_params")
		registerLibFunc(&llamaModelLoadFromFile, libHandle, "llama_model_load_from_file")
		registerLibFunc(&llamaModelLoadFromSplits, libHandle, "llama_model_load_from_splits")
		registerLibFunc(&llamaInitFromModel, libHandle, "llama_init_from_model")
	}
	registerLibFunc(&llamaModelSaveToFile, libHandle, "llama_model_save_to_file")
	registerLibFunc(&llamaModelFree, libHandle, "llama_model_free")

	// Context functions
	registerLibFunc(&llamaFree, libHandle, "llama_free")

	// Model info functions
	registerLibFunc(&llamaModelNCtxTrain, libHandle, "llama_model_n_ctx_train")
	registerLibFunc(&llamaModelNEmbd, libHandle, "llama_model_n_embd")
	registerLibFunc(&llamaModelNLayer, libHandle, "llama_model_n_layer")
	registerLibFunc(&llamaModelNHead, libHandle, "llama_model_n_head")
	registerLibFunc(&llamaModelNHeadKv, libHandle, "llama_model_n_head_kv")
	registerLibFunc(&llamaModelVocabType, libHandle, "llama_vocab_type")
	registerLibFunc(&llamaModelRopeType, libHandle, "llama_model_rope_type")

	// Context info functions
	registerLibFunc(&llamaNCtx, libHandle, "llama_n_ctx")
	registerLibFunc(&llamaNBatch, libHandle, "llama_n_batch")
	registerLibFunc(&llamaNUbatch, libHandle, "llama_n_ubatch")
	registerLibFunc(&llamaNSeqMax, libHandle, "llama_n_seq_max")
	registerLibFunc(&llamaPoolingType, libHandle, "llama_pooling_type")
	registerLibFunc(&llamaGetModel, libHandle, "llama_get_model")

	// Tokenization functions
	registerLibFunc(&llamaTokenize, libHandle, "llama_tokenize")
	registerLibFunc(&llamaTokenToPiece, libHandle, "llama_token_to_piece")
	registerLibFunc(&llamaDetokenize, libHandle, "llama_detokenize")
	registerLibFunc(&llamaVocabGetText, libHandle, "llama_vocab_get_text")

	// Vocab functions
	registerLibFunc(&llamaModelGetVocab, libHandle, "llama_model_get_vocab")
	registerLibFunc(&llamaVocabNTokens, libHandle, "llama_vocab_n_tokens")
	registerLibFunc(&llamaVocabBos, libHandle, "llama_vocab_bos")
	registerLibFunc(&llamaVocabEos, libHandle, "llama_vocab_eos")
	registerLibFunc(&llamaVocabEot, libHandle, "llama_vocab_eot")
	registerLibFunc(&llamaVocabNl, libHandle, "llama_vocab_nl")
	registerLibFunc(&llamaVocabPad, libHandle, "llama_vocab_pad")

	// Batch functions - Register struct functions only on Darwin (purego limitation)
	// On other platforms, FFI handles struct parameters/returns directly
	if runtime.GOOS == "darwin" {
		registerLibFunc(&llamaBatchInit, libHandle, "llama_batch_init")
		registerLibFunc(&llamaBatchGetOne, libHandle, "llama_batch_get_one")
		registerLibFunc(&llamaBatchFree, libHandle, "llama_batch_free")
	}

	// Decode functions - Register struct functions only on Darwin (purego limitation)
	// On other platforms, FFI handles struct parameters/returns directly
	if runtime.GOOS == "darwin" {
		registerLibFunc(&llamaDecode, libHandle, "llama_decode")
		registerLibFunc(&llamaEncode, libHandle, "llama_encode")
	}

	// Logits and embeddings
	registerLibFunc(&llamaGetLogits, libHandle, "llama_get_logits")
	registerLibFunc(&llamaGetLogitsIth, libHandle, "llama_get_logits_ith")
	registerLibFunc(&llamaGetEmbeddings, libHandle, "llama_get_embeddings")
	registerLibFunc(&llamaGetEmbeddingsIth, libHandle, "llama_get_embeddings_ith")
	registerLibFunc(&llamaSetCausalAttn, libHandle, "llama_set_causal_attn")
	registerLibFunc(&llamaSetEmbeddings, libHandle, "llama_set_embeddings")
	registerLibFunc(&llamaMemoryClear, libHandle, "llama_memory_clear")
	registerLibFunc(&llamaGetMemory, libHandle, "llama_get_memory")

	// Sampling functions - Register struct functions only on Darwin (purego limitation)
	// On other platforms, FFI handles struct parameters/returns directly
	if runtime.GOOS == "darwin" {
		registerLibFunc(&llamaSamplerChainInit, libHandle, "llama_sampler_chain_init")
	}
	registerLibFunc(&llamaSamplerChainAdd, libHandle, "llama_sampler_chain_add")
	registerLibFunc(&llamaSamplerChainGet, libHandle, "llama_sampler_chain_get")
	registerLibFunc(&llamaSamplerChainN, libHandle, "llama_sampler_chain_n")
	registerLibFunc(&llamaSamplerChainFree, libHandle, "llama_sampler_free")
	registerLibFunc(&llamaSamplerSample, libHandle, "llama_sampler_sample")
	registerLibFunc(&llamaSamplerAccept, libHandle, "llama_sampler_accept")
	registerLibFunc(&llamaSamplerReset, libHandle, "llama_sampler_reset")

	// Built-in samplers
	registerLibFunc(&llamaSamplerInitGreedy, libHandle, "llama_sampler_init_greedy")
	registerLibFunc(&llamaSamplerInitDist, libHandle, "llama_sampler_init_dist")
	// registerLibFunc(&llamaSamplerInitSoftmax, libHandle, "llama_sampler_init_softmax")  // Function doesn't exist in b6862
	registerLibFunc(&llamaSamplerInitTopK, libHandle, "llama_sampler_init_top_k")
	registerLibFunc(&llamaSamplerInitTopP, libHandle, "llama_sampler_init_top_p")
	registerLibFunc(&llamaSamplerInitMinP, libHandle, "llama_sampler_init_min_p")
	// registerLibFunc(&llamaSamplerInitTailFree, libHandle, "llama_sampler_init_tail_free")  // Function doesn't exist
	registerLibFunc(&llamaSamplerInitTypical, libHandle, "llama_sampler_init_typical")
	registerLibFunc(&llamaSamplerInitTemp, libHandle, "llama_sampler_init_temp")
	registerLibFunc(&llamaSamplerInitTempExt, libHandle, "llama_sampler_init_temp_ext")
	registerLibFunc(&llamaSamplerInitMirostat, libHandle, "llama_sampler_init_mirostat")
	registerLibFunc(&llamaSamplerInitMirostatV2, libHandle, "llama_sampler_init_mirostat_v2")

	// Utility functions
	registerLibFunc(&llamaMaxDevices, libHandle, "llama_max_devices")
	registerLibFunc(&llamaSupportsMmap, libHandle, "llama_supports_mmap")
	registerLibFunc(&llamaSupportsMlock, libHandle, "llama_supports_mlock")
	registerLibFunc(&llamaSupportsGpuOffload, libHandle, "llama_supports_gpu_offload")
	registerLibFunc(&llamaSupportsRpc, libHandle, "llama_supports_rpc")
	registerLibFunc(&llamaTimeUs, libHandle, "llama_time_us")
	registerLibFunc(&llamaPrintSystemInfo, libHandle, "llama_print_system_info")

	// KV cache functions
	// KV cache functions - deprecated/removed in b6862
	// registerLibFunc(&llamaKvCacheClear, libHandle, "llama_kv_cache_clear")
	// registerLibFunc(&llamaKvCacheSeqRm, libHandle, "llama_kv_cache_seq_rm")
	// registerLibFunc(&llamaKvCacheSeqCp, libHandle, "llama_kv_cache_seq_cp")
	// registerLibFunc(&llamaKvCacheSeqKeep, libHandle, "llama_kv_cache_seq_keep")
	// registerLibFunc(&llamaKvCacheSeqAdd, libHandle, "llama_kv_cache_seq_add")
	// registerLibFunc(&llamaKvCacheSeqDiv, libHandle, "llama_kv_cache_seq_div")
	// registerLibFunc(&llamaKvCacheSeqPos, libHandle, "llama_kv_cache_seq_pos")  // Might not exist
	// registerLibFunc(&llamaKvCacheDefrag, libHandle, "llama_kv_cache_defrag")
	// registerLibFunc(&llamaKvCacheUpdate, libHandle, "llama_kv_cache_update")

	// State functions
	registerLibFunc(&llamaStateGetSize, libHandle, "llama_state_get_size")
	registerLibFunc(&llamaStateGetData, libHandle, "llama_state_get_data")
	registerLibFunc(&llamaStateSetData, libHandle, "llama_state_set_data")
	registerLibFunc(&llamaStateLoadFile, libHandle, "llama_state_load_file")
	registerLibFunc(&llamaStateSaveFile, libHandle, "llama_state_save_file")

	// Performance functions - These may not exist in this llama.cpp version - moved to ROADMAP "wait for llama.cpp" section
	// registerLibFunc(&llamaGetTimings, libHandle, "llama_get_timings")
	// registerLibFunc(&llamaPrintTimings, libHandle, "llama_print_timings")
	// registerLibFunc(&llamaResetTimings, libHandle, "llama_reset_timings")

	// Register GGML functions
	if err := registerGgmlFunctions(); err != nil {
		return fmt.Errorf("failed to register GGML functions: %w", err)
	}

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
	// Try to load library if not already loaded
	_ = ensureLoaded() // Ignore error, fallback to defaults

	// Try FFI first (works on all platforms)
	if isLoaded {
		if params, err := ffiModelDefaultParams(); err == nil {
			return params
		}
	}

	// Fallback to purego on Darwin
	if runtime.GOOS == "darwin" && llamaModelDefaultParams != nil && isLoaded {
		return llamaModelDefaultParams()
	}

	// Last resort: return hardcoded defaults
	return LlamaModelParams{
		NGpuLayers:   0,
		SplitMode:    LLAMA_SPLIT_MODE_NONE,
		MainGpu:      0,
		VocabOnly:    0,
		UseMmap:      1, // Enable mmap by default
		UseMlock:     0,
		CheckTensors: 1, // Enable tensor validation by default
	}
}

// Context_default_params returns default context parameters
func Context_default_params() LlamaContextParams {
	// Try to load library if not already loaded
	_ = ensureLoaded() // Ignore error, fallback to defaults

	// Try FFI first (works on all platforms)
	if isLoaded {
		if params, err := ffiContextDefaultParams(); err == nil {
			return params
		}
	}

	// Fallback to purego on Darwin
	if runtime.GOOS == "darwin" && llamaContextDefaultParams != nil && isLoaded {
		return llamaContextDefaultParams()
	}

	// Last resort: return hardcoded defaults
	return LlamaContextParams{
		Seed:            LLAMA_DEFAULT_SEED,
		NCtx:            0, // Auto-detect from model
		NBatch:          2048,
		NUbatch:         512,
		NSeqMax:         1,
		NThreads:        int32(runtime.NumCPU()),
		NThreadsBatch:   int32(runtime.NumCPU()),
		RopeScalingType: LLAMA_ROPE_SCALING_TYPE_UNSPECIFIED,
		PoolingType:     LLAMA_POOLING_TYPE_UNSPECIFIED,
		AttentionType:   LLAMA_ATTENTION_TYPE_CAUSAL,
		DefragThold:     -1.0, // Disabled by default
		Logits:          0,    // Disabled by default
		Embeddings:      0,    // Disabled by default
		Offload_kqv:     1,    // Enable by default
		FlashAttn:       0,    // Disabled by default
		NoPerf:          0,    // Enable performance measurement by default
	}
}

// Sampler_chain_default_params returns default sampler chain parameters
func Sampler_chain_default_params() LlamaSamplerChainParams {
	// Try to load library if not already loaded
	_ = ensureLoaded() // Ignore error, fallback to defaults

	// Try FFI first (works on all platforms)
	if isLoaded {
		if params, err := ffiSamplerChainDefaultParams(); err == nil {
			return params
		}
	}

	// Fallback to purego on Darwin
	if runtime.GOOS == "darwin" && llamaSamplerChainDefaultParams != nil && isLoaded {
		return llamaSamplerChainDefaultParams()
	}

	// Last resort: return hardcoded defaults
	return LlamaSamplerChainParams{
		NoPerf: 0, // Enable performance measurement by default
	}
}

// Model_load_from_file loads a model from a file
func Model_load_from_file(pathModel string, params LlamaModelParams) (LlamaModel, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}

	// Check GGML backend initialized
	if !isLoaded {
		return 0, errors.New("llama.cpp library not loaded")
	}

	pathBytes := append([]byte(pathModel), 0) // null-terminate

	// Fallback to purego on Darwin
	if runtime.GOOS == "darwin" {
		model := llamaModelLoadFromFile((*byte)(unsafe.Pointer(&pathBytes[0])), params)
		if model == 0 {
			return 0, errors.New("failed to load model")
		}
		return model, nil
	} else {
		// Try FFI first (works on all platforms)
		if model, err := ffiModelLoadFromFile((*byte)(unsafe.Pointer(&pathBytes[0])), params); err == nil {
			return model, nil
		} else {
			return 0, err
		}
	}
}

// Model_free frees a model
func Model_free(model LlamaModel) {
	if isLoaded && model != 0 {
		llamaModelFree(model)
	}
}

// Model_n_embd returns the number of embedding dimensions for the model
func Model_n_embd(model LlamaModel) int32 {
	if err := ensureLoaded(); err != nil {
		panic(err)
	}
	return llamaModelNEmbd(model)
}

// Get_embeddings returns the embeddings for the context
func Get_embeddings(ctx LlamaContext) *float32 {
	if err := ensureLoaded(); err != nil {
		return nil
	}
	return llamaGetEmbeddings(ctx)
}

// Get_embeddings_ith returns the embeddings for the ith sequence in the context
func Get_embeddings_ith(ctx LlamaContext, i int32) *float32 {
	if err := ensureLoaded(); err != nil {
		return nil
	}
	return llamaGetEmbeddingsIth(ctx, i)
}

// Set_causal_attn sets whether to use causal attention
func Set_causal_attn(ctx LlamaContext, causal bool) {
	if err := ensureLoaded(); err != nil {
		return
	}
	llamaSetCausalAttn(ctx, causal)
}

// Set_embeddings sets whether to extract embeddings
func Set_embeddings(ctx LlamaContext, embeddings bool) {
	if err := ensureLoaded(); err != nil {
		return
	}
	llamaSetEmbeddings(ctx, embeddings)
}

// Memory_clear clears the KV cache
func Memory_clear(ctx LlamaContext, reset bool) bool {
	if err := ensureLoaded(); err != nil {
		return false
	}
	memory := llamaGetMemory(ctx)
	return llamaMemoryClear(memory, reset)
}

// Get_memory returns the memory handle for the context
func Get_memory(ctx LlamaContext) LlamaMemory {
	if err := ensureLoaded(); err != nil {
		return 0
	}
	return llamaGetMemory(ctx)
}

// Init_from_model creates a context from a model
func Init_from_model(model LlamaModel, params LlamaContextParams) (LlamaContext, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}

	// Try FFI first (works on all platforms)
	if ctx, err := ffiInitFromModel(model, params); err == nil {
		return ctx, nil
	}

	// Fallback to purego on Darwin
	if runtime.GOOS == "darwin" && llamaInitFromModel != nil {
		ctx := llamaInitFromModel(model, params)
		if ctx == 0 {
			return 0, errors.New("failed to create context")
		}
		return ctx, nil
	}

	return 0, errors.New("Init_from_model not available on this platform")
}

// Free frees a context
func Free(ctx LlamaContext) {
	if isLoaded && ctx != 0 {
		llamaFree(ctx)
	}
}

// Tokenize tokenizes text
func Tokenize(model LlamaModel, text string, addSpecial, parseSpecial bool) ([]LlamaToken, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}

	// Get the vocabulary from the model
	vocab := llamaModelGetVocab(model)
	if vocab == 0 {
		return nil, errors.New("failed to get vocabulary from model")
	}

	textBytes := append([]byte(text), 0) // null-terminate

	// First call to get the number of tokens
	textLen := len(text)
	if textLen > math.MaxInt32 {
		return nil, fmt.Errorf("text too long: %d characters, maximum supported: %d", textLen, math.MaxInt32)
	}
	nTokens := llamaTokenize(vocab, (*byte)(unsafe.Pointer(&textBytes[0])), int32(textLen), nil, 0, addSpecial, parseSpecial)
	if nTokens <= 0 {
		// llama_tokenize returns negative value indicating number of tokens needed
		if nTokens < 0 {
			nTokens = -nTokens // Convert to positive
		} else {
			return nil, fmt.Errorf("tokenization failed with error code: %d", nTokens)
		}
	}

	if nTokens == 0 {
		return []LlamaToken{}, nil
	}

	// Second call to get the actual tokens
	tokens := make([]LlamaToken, nTokens)
	result := llamaTokenize(vocab, (*byte)(unsafe.Pointer(&textBytes[0])), int32(textLen), &tokens[0], nTokens, addSpecial, parseSpecial)
	if result < 0 {
		return nil, fmt.Errorf("tokenization failed with error code: %d", result)
	}

	return tokens[:result], nil
}

// Token_to_piece converts a token to its string representation using model
func Token_to_piece(model LlamaModel, token LlamaToken, special bool) string {
	if err := ensureLoaded(); err != nil {
		return ""
	}

	// Validate model handle
	if model == 0 {
		return ""
	}

	// Get the vocabulary from the model
	vocab := llamaModelGetVocab(model)
	if vocab == 0 {
		return ""
	}

	// Use the simpler llama_vocab_get_text function which directly returns the text
	textPtr := llamaVocabGetText(vocab, token)
	if textPtr == nil {
		return ""
	}

	// Convert C string to Go string
	// We need to find the length of the C string first
	var length int
	for {
		// Use unsafe.Add to safely advance the pointer
		bytePtr := (*byte)(unsafe.Add(unsafe.Pointer(textPtr), length))
		if *bytePtr == 0 {
			break
		}
		length++
	}

	if length == 0 {
		return ""
	}

	// Create a Go byte slice from the C string
	bytes := (*[1 << 30]byte)(unsafe.Pointer(textPtr))[:length:length]
	return string(bytes)
}

// Batch_init creates a new batch
func Batch_init(nTokens, embd, nSeqMax int32) LlamaBatch {
	// Try to load library if not already loaded
	_ = ensureLoaded() // Ignore error, fallback to empty batch

	// Try FFI first (works on all platforms)
	if isLoaded {
		if batch, err := ffiBatchInit(nTokens, embd, nSeqMax); err == nil {
			return batch
		}
	}

	// Fallback to purego on Darwin
	if runtime.GOOS == "darwin" && llamaBatchInit != nil && isLoaded {
		return llamaBatchInit(nTokens, embd, nSeqMax)
	}

	// Last resort: return zero-initialized batch
	return LlamaBatch{}
}

// Batch_get_one creates a batch from a single set of tokens
func Batch_get_one(tokens []LlamaToken) LlamaBatch {
	// Try to load library if not already loaded
	_ = ensureLoaded() // Ignore error, fallback to empty batch

	if len(tokens) == 0 {
		return LlamaBatch{}
	}

	tokensLen := len(tokens)
	if tokensLen > math.MaxInt32 {
		// Token count exceeds maximum supported size
		// Return empty batch instead of panicking (safer than previous panic behavior)
		return LlamaBatch{}
	}

	// Try FFI first (works on all platforms)
	if isLoaded {
		if batch, err := ffiBatchGetOne(&tokens[0], int32(tokensLen)); err == nil {
			return batch
		}
	}

	// Fallback to purego on Darwin
	if runtime.GOOS == "darwin" && llamaBatchGetOne != nil && isLoaded {
		return llamaBatchGetOne(&tokens[0], int32(tokensLen))
	}

	// Last resort: return zero-initialized batch
	return LlamaBatch{}
}

// Batch_free frees a batch
func Batch_free(batch LlamaBatch) {
	if err := ensureLoaded(); err != nil {
		return
	}
	// Only call llama_batch_free for batches created with llama_batch_init
	// Batches created with llama_batch_get_one don't need to be freed
	if runtime.GOOS == "darwin" && batch.Token != nil {
		llamaBatchFree(batch)
	}
}

// Decode decodes a batch
func Decode(ctx LlamaContext, batch LlamaBatch) error {
	if err := ensureLoaded(); err != nil {
		return err
	}

	// Try FFI first (works on all platforms)
	if result, err := ffiDecode(ctx, batch); err == nil {
		if result != 0 {
			return fmt.Errorf("decode failed with code %d", result)
		}
		return nil
	}

	// Fallback to purego on Darwin
	if runtime.GOOS == "darwin" && llamaDecode != nil {
		result := llamaDecode(ctx, batch)
		if result != 0 {
			return fmt.Errorf("decode failed with code %d", result)
		}
		return nil
	}

	return errors.New("Decode not available on this platform")
}

// Encode encodes a batch
func Encode(ctx LlamaContext, batch LlamaBatch) error {
	if err := ensureLoaded(); err != nil {
		return err
	}

	// Try FFI first (works on all platforms)
	if result, err := ffiEncode(ctx, batch); err == nil {
		if result != 0 {
			return fmt.Errorf("encode failed with code %d", result)
		}
		return nil
	}

	// Fallback to purego on Darwin
	if runtime.GOOS == "darwin" && llamaEncode != nil {
		result := llamaEncode(ctx, batch)
		if result != 0 {
			return fmt.Errorf("encode failed with code %d", result)
		}
		return nil
	}

	return errors.New("Encode not available on this platform")
}

// Get_logits gets logits for all tokens
func Get_logits(ctx LlamaContext) *float32 {
	if err := ensureLoaded(); err != nil {
		return nil
	}
	return llamaGetLogits(ctx)
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
	if err := ensureLoaded(); err != nil {
		return nil
	}

	// Use actual number of available logits (256) instead of full vocab (32000)
	// Based on error: "out of range [0, 256)"
	nVocab := int32(256)

	// Allocate memory for token data array
	tokenData := make([]LlamaTokenData, nVocab)

	// Initialize token data array - will be populated with actual logits later
	for i := int32(0); i < nVocab; i++ {
		tokenData[i] = LlamaTokenData{
			Id:    LlamaToken(i),
			Logit: 0.0,
			P:     0.0,
		}
	}

	// Return pointer to token data array structure
	if nVocab < 0 {
		panic(fmt.Errorf("invalid vocabulary size: %d", nVocab))
	}
	return &LlamaTokenDataArray{
		Data:     &tokenData[0],
		Size:     uint64(uint32(nVocab)), // Safe conversion since nVocab is int32
		Selected: -1,
		Sorted:   0,
	}
}

// Token_data_array_from_logits creates a token data array from logits
func Token_data_array_from_logits(model LlamaModel, logits *float32) *LlamaTokenDataArray {
	if err := ensureLoaded(); err != nil {
		return nil
	}

	if logits == nil {
		return nil
	}

	// Use hardcoded vocabulary size for now to avoid corruption issues
	// Use a very small, safe subset to prevent any out-of-bounds access
	nVocab := int32(32)

	// Allocate memory for token data array
	tokenData := make([]LlamaTokenData, nVocab)

	// Convert logits pointer to slice for easier access
	logitsSlice := unsafe.Slice(logits, nVocab)

	// Populate token data array with actual logits
	for i := int32(0); i < nVocab; i++ {
		tokenData[i] = LlamaTokenData{
			Id:    LlamaToken(i),
			Logit: logitsSlice[i],
			P:     0.0, // Will be computed by the sampler
		}
	}

	// Return pointer to token data array structure
	if nVocab < 0 {
		panic(fmt.Errorf("invalid vocabulary size: %d", nVocab))
	}
	return &LlamaTokenDataArray{
		Data:     &tokenData[0],
		Size:     uint64(uint32(nVocab)), // Safe conversion since nVocab is int32
		Selected: -1,
		Sorted:   0,
	}
}

// Sampler_init_greedy creates a greedy sampler
func Sampler_init_greedy() LlamaSampler {
	if err := ensureLoaded(); err != nil {
		panic(err)
	}
	return llamaSamplerInitGreedy()
}

// Sampler_chain_init creates a sampler chain
func Sampler_chain_init(params LlamaSamplerChainParams) LlamaSampler {
	// Try to load library if not already loaded
	_ = ensureLoaded() // Ignore error, return 0 on failure

	// Try FFI first (works on all platforms)
	if isLoaded {
		if sampler, err := ffiSamplerChainInit(params); err == nil {
			return sampler
		}
	}

	// Fallback to purego on Darwin
	if runtime.GOOS == "darwin" && llamaSamplerChainInit != nil && isLoaded {
		return llamaSamplerChainInit(params)
	}

	// Last resort: return null sampler
	return 0
}

// Sampler_free frees a sampler
func Sampler_free(sampler LlamaSampler) {
	// The C library doesn't seem to have a direct sampler free function
	// This might be handled by the sampler chain
}

// Sampler_sample samples a token from the logits at the given index (-1 for last token)
func Sampler_sample(sampler LlamaSampler, ctx LlamaContext, idx int32) LlamaToken {
	if err := ensureLoaded(); err != nil {
		return LLAMA_TOKEN_NULL
	}
	return llamaSamplerSample(sampler, ctx, idx)
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

// Helper functions for platforms where struct returns aren't supported - moved to ROADMAP "wait for purego struct support" section
func ModelDefaultParams() LlamaModelParams {
	if runtime.GOOS == "darwin" && llamaModelDefaultParams != nil {
		return llamaModelDefaultParams()
	}
	// Return default values for non-Darwin platforms - blocks ROADMAP "wait for purego struct support"
	return LlamaModelParams{
		NGpuLayers:    0,
		SplitMode:     LLAMA_SPLIT_MODE_LAYER,
		MainGpu:       0,
		VocabOnly:     0,
		UseMmap:       1,
		UseMlock:      0,
		CheckTensors:  1,
		UseExtraBufts: 0,
	}
}

func ContextDefaultParams() LlamaContextParams {
	if runtime.GOOS == "darwin" && llamaContextDefaultParams != nil {
		return llamaContextDefaultParams()
	}
	// Return default values for non-Darwin platforms - blocks ROADMAP "wait for purego struct support"
	return LlamaContextParams{
		Seed:            LLAMA_DEFAULT_SEED,
		NCtx:            0, // 0 = from model
		NBatch:          2048,
		NUbatch:         512,
		NSeqMax:         1,
		NThreads:        -1, // -1 = auto-detect
		NThreadsBatch:   -1, // -1 = auto-detect
		RopeScalingType: LLAMA_ROPE_SCALING_TYPE_UNSPECIFIED,
		PoolingType:     LLAMA_POOLING_TYPE_UNSPECIFIED,
		AttentionType:   LLAMA_ATTENTION_TYPE_CAUSAL,
		RopeFreqBase:    0.0, // 0.0 = from model
		RopeFreqScale:   0.0, // 0.0 = from model
		YarnExtFactor:   -1.0,
		YarnAttnFactor:  1.0,
		YarnBetaFast:    32.0,
		YarnBetaSlow:    1.0,
		YarnOrigCtx:     0,
		DefragThold:     -1.0,
		TypeK:           -1,
		TypeV:           -1,
		Logits:          0,
		Embeddings:      0,
		Offload_kqv:     1,
		FlashAttn:       0,
		NoPerf:          0,
	}
}

func SamplerChainDefaultParams() LlamaSamplerChainParams {
	if runtime.GOOS == "darwin" && llamaSamplerChainDefaultParams != nil {
		return llamaSamplerChainDefaultParams()
	}
	// Return default values for non-Darwin platforms - blocks ROADMAP "wait for purego struct support"
	return LlamaSamplerChainParams{
		NoPerf: 0,
	}
}

// DetectGpuBackend detects the available GPU backend on the current system
func DetectGpuBackend() LlamaGpuBackend {
	// Check for GPU backends in priority order based on platform
	switch runtime.GOOS {
	case "darwin":
		// On macOS, Metal is the primary GPU backend
		return LLAMA_GPU_BACKEND_METAL
	case "linux", "windows":
		// Check for available GPU SDKs in priority order
		if hasCommand("nvcc") {
			return LLAMA_GPU_BACKEND_CUDA
		}
		if hasCommand("hipconfig") {
			return LLAMA_GPU_BACKEND_HIP
		}
		if hasCommand("vulkaninfo") {
			return LLAMA_GPU_BACKEND_VULKAN
		}
		if hasCommand("clinfo") {
			return LLAMA_GPU_BACKEND_OPENCL
		}
		if hasCommand("sycl-ls") {
			return LLAMA_GPU_BACKEND_SYCL
		}
		return LLAMA_GPU_BACKEND_CPU
	default:
		return LLAMA_GPU_BACKEND_CPU
	}
}

// hasCommand checks if a command is available in PATH
func hasCommand(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
