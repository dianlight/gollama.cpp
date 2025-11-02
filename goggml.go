// Package gollama provides Go bindings for ggml (the tensor library used by llama.cpp).
// This file contains bindings for the core GGML tensor operations and utilities.
//
// GGML (Georgi Gerganov Machine Learning) is the tensor library that powers llama.cpp.
// It provides low-level operations for neural network computation.
//
// # Usage
//
// Most users should use the high-level llama.cpp API in gollama.go. Use the GGML
// bindings when you need direct access to tensor operations, type information,
// backend management, or low-level memory operations.
//
// # Important Note
//
// GGML functions may not be exported in all llama.cpp builds. This package gracefully
// handles missing functions by returning errors instead of panicking, allowing code to
// compile and run even when GGML symbols are not available.
//
// # Example Usage
//
//	// Initialize the library
//	gollama.Backend_init()
//	defer gollama.Backend_free()
//
//	// Query type information
//	size, err := gollama.Ggml_type_size(gollama.GGML_TYPE_F32)
//	if err != nil {
//	    log.Printf("GGML function not available: %v", err)
//	    return
//	}
//	fmt.Printf("F32 size: %d bytes\n", size)
//
//	// Check if a type is quantized
//	isQuantized, err := gollama.Ggml_type_is_quantized(gollama.GGML_TYPE_Q4_0)
//	if err == nil {
//	    fmt.Printf("Q4_0 is quantized: %v\n", isQuantized)
//	}
//
//	// Enumerate backend devices
//	count, err := gollama.Ggml_backend_dev_count()
//	if err == nil && count > 0 {
//	    for i := uint64(0); i < count; i++ {
//	        dev, _ := gollama.Ggml_backend_dev_get(i)
//	        name, _ := gollama.Ggml_backend_dev_name(dev)
//	        fmt.Printf("Device %d: %s\n", i, name)
//	    }
//	}
//
// For more details, see the GGML API documentation at:
// https://github.com/dianlight/gollama.cpp/blob/main/docs/GGML_API.md
package gollama

import (
	"fmt"
	"log/slog"
	"unsafe"
)

// GGML tensor types
type GgmlType int32

const (
	GGML_TYPE_F32     GgmlType = 0
	GGML_TYPE_F16     GgmlType = 1
	GGML_TYPE_Q4_0    GgmlType = 2
	GGML_TYPE_Q4_1    GgmlType = 3
	GGML_TYPE_Q5_0    GgmlType = 6
	GGML_TYPE_Q5_1    GgmlType = 7
	GGML_TYPE_Q8_0    GgmlType = 8
	GGML_TYPE_Q8_1    GgmlType = 9
	GGML_TYPE_Q2_K    GgmlType = 10
	GGML_TYPE_Q3_K    GgmlType = 11
	GGML_TYPE_Q4_K    GgmlType = 12
	GGML_TYPE_Q5_K    GgmlType = 13
	GGML_TYPE_Q6_K    GgmlType = 14
	GGML_TYPE_Q8_K    GgmlType = 15
	GGML_TYPE_IQ2_XXS GgmlType = 16
	GGML_TYPE_IQ2_XS  GgmlType = 17
	GGML_TYPE_IQ3_XXS GgmlType = 18
	GGML_TYPE_IQ1_S   GgmlType = 19
	GGML_TYPE_IQ4_NL  GgmlType = 20
	GGML_TYPE_IQ3_S   GgmlType = 21
	GGML_TYPE_IQ2_S   GgmlType = 22
	GGML_TYPE_IQ4_XS  GgmlType = 23
	GGML_TYPE_I8      GgmlType = 24
	GGML_TYPE_I16     GgmlType = 25
	GGML_TYPE_I32     GgmlType = 26
	GGML_TYPE_I64     GgmlType = 27
	GGML_TYPE_F64     GgmlType = 28
	GGML_TYPE_IQ1_M   GgmlType = 29
	GGML_TYPE_BF16    GgmlType = 30
	GGML_TYPE_COUNT   GgmlType = 31
)

// String returns the string representation of a GGML type
func (t GgmlType) String() string {
	switch t {
	case GGML_TYPE_F32:
		return "f32"
	case GGML_TYPE_F16:
		return "f16"
	case GGML_TYPE_Q4_0:
		return "q4_0"
	case GGML_TYPE_Q4_1:
		return "q4_1"
	case GGML_TYPE_Q5_0:
		return "q5_0"
	case GGML_TYPE_Q5_1:
		return "q5_1"
	case GGML_TYPE_Q8_0:
		return "q8_0"
	case GGML_TYPE_Q8_1:
		return "q8_1"
	case GGML_TYPE_Q2_K:
		return "q2_K"
	case GGML_TYPE_Q3_K:
		return "q3_K"
	case GGML_TYPE_Q4_K:
		return "q4_K"
	case GGML_TYPE_Q5_K:
		return "q5_K"
	case GGML_TYPE_Q6_K:
		return "q6_K"
	case GGML_TYPE_Q8_K:
		return "q8_K"
	case GGML_TYPE_IQ2_XXS:
		return "iq2_xxs"
	case GGML_TYPE_IQ2_XS:
		return "iq2_xs"
	case GGML_TYPE_IQ3_XXS:
		return "iq3_xxs"
	case GGML_TYPE_IQ1_S:
		return "iq1_s"
	case GGML_TYPE_IQ4_NL:
		return "iq4_nl"
	case GGML_TYPE_IQ3_S:
		return "iq3_s"
	case GGML_TYPE_IQ2_S:
		return "iq2_s"
	case GGML_TYPE_IQ4_XS:
		return "iq4_xs"
	case GGML_TYPE_I8:
		return "i8"
	case GGML_TYPE_I16:
		return "i16"
	case GGML_TYPE_I32:
		return "i32"
	case GGML_TYPE_I64:
		return "i64"
	case GGML_TYPE_F64:
		return "f64"
	case GGML_TYPE_IQ1_M:
		return "iq1_m"
	case GGML_TYPE_BF16:
		return "bf16"
	default:
		return "unknown"
	}
}

// GGML backend types
type GgmlBackend uintptr
type GgmlBackendBuffer uintptr
type GgmlBackendBufferType uintptr
type GgmlBackendDevice uintptr
type GgmlBackendReg uintptr

// GGML tensor type
type GgmlTensor uintptr

// GGML context type
type GgmlContext uintptr

// GGML compute plan
type GgmlCplan uintptr

// GGML object type
type GgmlObject int32

const (
	GGML_OBJECT_TENSOR GgmlObject = 0
	GGML_OBJECT_GRAPH  GgmlObject = 1
	GGML_OBJECT_WORK   GgmlObject = 2
)

// GGML operation types
type GgmlOp int32

const (
	GGML_OP_NONE GgmlOp = 0
	GGML_OP_DUP  GgmlOp = 1
	GGML_OP_ADD  GgmlOp = 2
	GGML_OP_SUB  GgmlOp = 3
	GGML_OP_MUL  GgmlOp = 4
	GGML_OP_DIV  GgmlOp = 5
	// Add more operations as needed
)

// Function pointers for GGML functions
var (
	// Type size functions
	ggmlTypeSize    func(typ GgmlType) uint64
	ggmlTypeSizeof  func(typ GgmlType) uint64
	ggmlBlckSize    func(typ GgmlType) int32
	ggmlIsQuantized func(typ GgmlType) bool

	// Backend device functions
	ggmlBackendDevCount       func() uint64
	ggmlBackendDevGet         func(index uint64) GgmlBackendDevice
	ggmlBackendDevByType      func(typ int32) GgmlBackendDevice
	ggmlBackendDevInit        func(device GgmlBackendDevice, params uintptr) GgmlBackend
	ggmlBackendDevName        func(device GgmlBackendDevice) *byte
	ggmlBackendDevDescription func(device GgmlBackendDevice) *byte
	ggmlBackendDevMemory      func(device GgmlBackendDevice, free *uint64, total *uint64)

	// Backend buffer type functions
	ggmlBackendDevBufferType     func(device GgmlBackendDevice) GgmlBackendBufferType
	ggmlBackendDevHostBufferType func(device GgmlBackendDevice) GgmlBackendBufferType
	ggmlBackendCpuBufferType     func() GgmlBackendBufferType
	ggmlBackendBuftName          func(buft GgmlBackendBufferType) *byte
	ggmlBackendBuftAllocBuffer   func(buft GgmlBackendBufferType, size uint64) GgmlBackendBuffer
	ggmlBackendBuftIsHost        func(buft GgmlBackendBufferType) bool

	// Backend buffer functions
	ggmlBackendBufferFree     func(buffer GgmlBackendBuffer)
	ggmlBackendBufferGetBase  func(buffer GgmlBackendBuffer) unsafe.Pointer
	ggmlBackendBufferGetSize  func(buffer GgmlBackendBuffer) uint64
	ggmlBackendBufferClear    func(buffer GgmlBackendBuffer, value uint8)
	ggmlBackendBufferIsHost   func(buffer GgmlBackendBuffer) bool
	ggmlBackendBufferSetUsage func(buffer GgmlBackendBuffer, usage int32)
	ggmlBackendBufferGetType  func(buffer GgmlBackendBuffer) GgmlBackendBufferType
	ggmlBackendBufferName     func(buffer GgmlBackendBuffer) *byte

	// Backend functions
	ggmlBackendFree            func(backend GgmlBackend)
	ggmlBackendName            func(backend GgmlBackend) *byte
	ggmlBackendSupports        func(backend GgmlBackend, buft GgmlBackendBufferType) bool
	ggmlBackendLoad            func(name *byte, search_path *byte) GgmlBackend
	ggmlBackendLoadAll         func()
	ggmlBackendLoadAllFromPath func(path *byte)

	// Tensor utility functions
	ggmlNbytes       func(tensor GgmlTensor) uint64
	ggmlRowSize      func(typ GgmlType, ne int64) uint64
	ggmlTypeToString func(typ GgmlType) *byte
	ggmlElementSize  func(tensor GgmlTensor) uint64

	// Quantization functions
	ggmlQuantizeChunk func(typ GgmlType, src *float32, dst unsafe.Pointer, start int32, nrows int32, ncols int64, hist *int64) uint64
)

// registerGgmlFunctions registers all GGML function pointers
// Note: GGML functions may not be exported in all llama.cpp builds
// This function attempts to register them but doesn't fail if they're not available
func registerGgmlFunctions() error {
	// Try to register functions, but don't fail if they don't exist
	// Most GGML functions are internal to llama.cpp and not exported

	// Type size functions - these are usually available
	_ = tryRegisterLibFunc(&ggmlTypeSize, libHandle, "ggml_type_size")
	_ = tryRegisterLibFunc(&ggmlTypeSizeof, libHandle, "ggml_type_sizef")
	_ = tryRegisterLibFunc(&ggmlBlckSize, libHandle, "ggml_blck_size")
	_ = tryRegisterLibFunc(&ggmlIsQuantized, libHandle, "ggml_is_quantized")

	// Backend device functions
	_ = tryRegisterLibFunc(&ggmlBackendDevCount, libHandle, "ggml_backend_dev_count")
	_ = tryRegisterLibFunc(&ggmlBackendDevGet, libHandle, "ggml_backend_dev_get")
	_ = tryRegisterLibFunc(&ggmlBackendDevByType, libHandle, "ggml_backend_dev_by_type")
	_ = tryRegisterLibFunc(&ggmlBackendDevInit, libHandle, "ggml_backend_dev_init")
	_ = tryRegisterLibFunc(&ggmlBackendDevName, libHandle, "ggml_backend_dev_name")
	_ = tryRegisterLibFunc(&ggmlBackendDevDescription, libHandle, "ggml_backend_dev_description")
	_ = tryRegisterLibFunc(&ggmlBackendDevMemory, libHandle, "ggml_backend_dev_memory")

	// Backend buffer type functions
	_ = tryRegisterLibFunc(&ggmlBackendDevBufferType, libHandle, "ggml_backend_dev_buffer_type")
	_ = tryRegisterLibFunc(&ggmlBackendDevHostBufferType, libHandle, "ggml_backend_dev_host_buffer_type")
	_ = tryRegisterLibFunc(&ggmlBackendCpuBufferType, libHandle, "ggml_backend_cpu_buffer_type")
	_ = tryRegisterLibFunc(&ggmlBackendBuftName, libHandle, "ggml_backend_buft_name")
	_ = tryRegisterLibFunc(&ggmlBackendBuftAllocBuffer, libHandle, "ggml_backend_buft_alloc_buffer")
	_ = tryRegisterLibFunc(&ggmlBackendBuftIsHost, libHandle, "ggml_backend_buft_is_host")

	// Backend buffer functions
	_ = tryRegisterLibFunc(&ggmlBackendBufferFree, libHandle, "ggml_backend_buffer_free")
	_ = tryRegisterLibFunc(&ggmlBackendBufferGetBase, libHandle, "ggml_backend_buffer_get_base")
	_ = tryRegisterLibFunc(&ggmlBackendBufferGetSize, libHandle, "ggml_backend_buffer_get_size")
	_ = tryRegisterLibFunc(&ggmlBackendBufferClear, libHandle, "ggml_backend_buffer_clear")
	_ = tryRegisterLibFunc(&ggmlBackendBufferIsHost, libHandle, "ggml_backend_buffer_is_host")
	_ = tryRegisterLibFunc(&ggmlBackendBufferSetUsage, libHandle, "ggml_backend_buffer_set_usage")
	_ = tryRegisterLibFunc(&ggmlBackendBufferGetType, libHandle, "ggml_backend_buffer_get_type")
	_ = tryRegisterLibFunc(&ggmlBackendBufferName, libHandle, "ggml_backend_buffer_name")

	// Backend functions
	_ = tryRegisterLibFunc(&ggmlBackendFree, libHandle, "ggml_backend_free")
	_ = tryRegisterLibFunc(&ggmlBackendName, libHandle, "ggml_backend_name")
	_ = tryRegisterLibFunc(&ggmlBackendSupports, libHandle, "ggml_backend_supports_buft")
	_ = tryRegisterLibFunc(&ggmlBackendLoad, libHandle, "ggml_backend_load")
	_ = tryRegisterLibFunc(&ggmlBackendLoadAll, libHandle, "ggml_backend_load_all")
	_ = tryRegisterLibFunc(&ggmlBackendLoadAllFromPath, libHandle, "ggml_backend_load_all_from_path")

	// Tensor utility functions
	_ = tryRegisterLibFunc(&ggmlNbytes, libHandle, "ggml_nbytes")
	_ = tryRegisterLibFunc(&ggmlRowSize, libHandle, "ggml_row_size")
	_ = tryRegisterLibFunc(&ggmlTypeToString, libHandle, "ggml_type_name")
	_ = tryRegisterLibFunc(&ggmlElementSize, libHandle, "ggml_element_size")

	// Quantization functions
	_ = tryRegisterLibFunc(&ggmlQuantizeChunk, libHandle, "ggml_quantize_chunk")

	return nil
}

// Public API functions for GGML

// Ggml_type_size returns the size in bytes of a GGML type element
func Ggml_type_size(typ GgmlType) (uint64, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlTypeSize == nil {
		return 0, fmt.Errorf("ggml_type_size function not available")
	}
	return ggmlTypeSize(typ), nil
}

// Ggml_type_sizef returns the size in bytes of a GGML type (float version)
func Ggml_type_sizef(typ GgmlType) (uint64, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlTypeSizeof == nil {
		return 0, fmt.Errorf("ggml_type_sizef function not available")
	}
	return ggmlTypeSizeof(typ), nil
}

// Ggml_blck_size returns the block size of a GGML type
func Ggml_blck_size(typ GgmlType) (int32, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlBlckSize == nil {
		return 0, fmt.Errorf("ggml_blck_size function not available")
	}
	return ggmlBlckSize(typ), nil
}

// Ggml_type_is_quantized returns whether a GGML type is quantized
func Ggml_type_is_quantized(typ GgmlType) (bool, error) {
	if err := ensureLoaded(); err != nil {
		return false, err
	}
	if ggmlIsQuantized == nil {
		return false, fmt.Errorf("ggml_is_quantized function not available")
	}
	return ggmlIsQuantized(typ), nil
}

// Ggml_backend_dev_count returns the number of available backend devices
func Ggml_backend_dev_count() (uint64, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlBackendDevCount == nil {
		return 0, fmt.Errorf("ggml_backend_dev_count function not available")
	}
	return ggmlBackendDevCount(), nil
}

// Ggml_backend_dev_get returns a backend device by index
func Ggml_backend_dev_get(index uint64) (GgmlBackendDevice, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlBackendDevGet == nil {
		return 0, fmt.Errorf("ggml_backend_dev_get function not available")
	}
	return ggmlBackendDevGet(index), nil
}

// Ggml_backend_dev_name returns the name of a backend device
func Ggml_backend_dev_name(device GgmlBackendDevice) (string, error) {
	if err := ensureLoaded(); err != nil {
		return "", err
	}
	if ggmlBackendDevName == nil {
		return "", fmt.Errorf("ggml_backend_dev_name function not available")
	}
	namePtr := ggmlBackendDevName(device)
	if namePtr == nil {
		return "", nil
	}
	return bytePointerToString(namePtr), nil
}

// Ggml_backend_dev_description returns the description of a backend device
func Ggml_backend_dev_description(device GgmlBackendDevice) (string, error) {
	if err := ensureLoaded(); err != nil {
		return "", err
	}
	if ggmlBackendDevDescription == nil {
		return "", fmt.Errorf("ggml_backend_dev_description function not available")
	}
	descPtr := ggmlBackendDevDescription(device)
	if descPtr == nil {
		return "", nil
	}
	return bytePointerToString(descPtr), nil
}

// Ggml_backend_dev_memory returns the memory statistics of a backend device
func Ggml_backend_dev_memory(device GgmlBackendDevice) (free uint64, total uint64, err error) {
	if err := ensureLoaded(); err != nil {
		return 0, 0, err
	}
	if ggmlBackendDevMemory == nil {
		return 0, 0, fmt.Errorf("ggml_backend_dev_memory function not available")
	}
	ggmlBackendDevMemory(device, &free, &total)
	return free, total, nil
}

// Ggml_backend_cpu_buffer_type returns the CPU buffer type
func Ggml_backend_cpu_buffer_type() (GgmlBackendBufferType, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlBackendCpuBufferType == nil {
		return 0, fmt.Errorf("ggml_backend_cpu_buffer_type function not available")
	}
	return ggmlBackendCpuBufferType(), nil
}

// Ggml_backend_buffer_name returns the name of a backend buffer
func Ggml_backend_buffer_name(buffer GgmlBackendBuffer) (string, error) {
	if err := ensureLoaded(); err != nil {
		return "", err
	}
	if ggmlBackendBufferName == nil {
		return "", fmt.Errorf("ggml_backend_buffer_name function not available")
	}
	namePtr := ggmlBackendBufferName(buffer)
	if namePtr == nil {
		return "", nil
	}
	return bytePointerToString(namePtr), nil
}

// Ggml_backend_buffer_free frees a backend buffer
func Ggml_backend_buffer_free(buffer GgmlBackendBuffer) error {
	if err := ensureLoaded(); err != nil {
		return err
	}
	if ggmlBackendBufferFree == nil {
		return fmt.Errorf("ggml_backend_buffer_free function not available")
	}
	ggmlBackendBufferFree(buffer)
	return nil
}

// Ggml_backend_buffer_get_size returns the size of a backend buffer
func Ggml_backend_buffer_get_size(buffer GgmlBackendBuffer) (uint64, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlBackendBufferGetSize == nil {
		return 0, fmt.Errorf("ggml_backend_buffer_get_size function not available")
	}
	return ggmlBackendBufferGetSize(buffer), nil
}

// Ggml_backend_buffer_is_host checks if a buffer is host memory
func Ggml_backend_buffer_is_host(buffer GgmlBackendBuffer) (bool, error) {
	if err := ensureLoaded(); err != nil {
		return false, err
	}
	if ggmlBackendBufferIsHost == nil {
		return false, fmt.Errorf("ggml_backend_buffer_is_host function not available")
	}
	return ggmlBackendBufferIsHost(buffer), nil
}

// Ggml_backend_name returns the name of a backend
func Ggml_backend_name(backend GgmlBackend) (string, error) {
	if err := ensureLoaded(); err != nil {
		return "", err
	}
	if ggmlBackendName == nil {
		return "", fmt.Errorf("ggml_backend_name function not available")
	}
	namePtr := ggmlBackendName(backend)
	if namePtr == nil {
		return "", nil
	}
	return bytePointerToString(namePtr), nil
}

// Ggml_backend_free frees a backend
func Ggml_backend_free(backend GgmlBackend) error {
	if err := ensureLoaded(); err != nil {
		return err
	}
	if ggmlBackendFree == nil {
		return fmt.Errorf("ggml_backend_free function not available")
	}
	ggmlBackendFree(backend)
	return nil
}

// Ggml_backend_is_cpu checks if a backend is CPU-based
// Note: This function is not available in current GGML builds
func Ggml_backend_is_cpu(backend GgmlBackend) (bool, error) {
	if err := ensureLoaded(); err != nil {
		return false, err
	}
	// This function is not exported in GGML, return error
	return false, fmt.Errorf("ggml_backend_is_cpu function not available")
}

// Ggml_type_name returns the string name of a GGML type
func Ggml_type_name(typ GgmlType) (string, error) {
	if err := ensureLoaded(); err != nil {
		return "", err
	}
	if ggmlTypeToString == nil {
		return "", fmt.Errorf("ggml_type_name function not available")
	}
	namePtr := ggmlTypeToString(typ)
	if namePtr == nil {
		return "", nil
	}
	return bytePointerToString(namePtr), nil
}

// Ggml_backend_load dynamically loads a backend by name from a search path
func Ggml_backend_load(name string, searchPath string) (GgmlBackend, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlBackendLoad == nil {
		return 0, fmt.Errorf("ggml_backend_load function not available")
	}

	nameBytes := append([]byte(name), 0)
	var pathPtr *byte
	if searchPath != "" {
		pathBytes := append([]byte(searchPath), 0)
		pathPtr = &pathBytes[0]
	}

	return ggmlBackendLoad(&nameBytes[0], pathPtr), nil
}

// Ggml_backend_load_all loads all available backends
func Ggml_backend_load_all() error {
	if err := ensureLoaded(); err != nil {
		return err
	}
	if ggmlBackendLoadAll == nil {
		return fmt.Errorf("ggml_backend_load_all function not available")
	}

	//	os.Setenv("GGML_BACKEND_PATH", globalLoader.libPath)
	if globalLoader.rootLibPath == "" {

		err := globalLoader.LoadLibrary()
		if err != nil {
			return fmt.Errorf("failed to load library for backend loading: %v", err)
		}
	}
	slog.Info("Loading GGML backends from path", "path", globalLoader.rootLibPath)
	ggmlBackendLoadAllFromPath(&[]byte(globalLoader.rootLibPath + "\x00")[0])
	return nil
}

// Ggml_backend_load_all_from_path loads all available backends from a specific path
func Ggml_backend_load_all_from_path(path string) error {
	if err := ensureLoaded(); err != nil {
		return err
	}
	if ggmlBackendLoadAllFromPath == nil {
		return fmt.Errorf("ggml_backend_load_all_from_path function not available")
	}

	var pathPtr *byte
	if path != "" {
		pathBytes := append([]byte(path), 0)
		pathPtr = &pathBytes[0]
	}

	ggmlBackendLoadAllFromPath(pathPtr)
	return nil
}

// Helper function to convert byte pointer to Go string
func bytePointerToString(ptr *byte) string {
	if ptr == nil {
		return ""
	}
	var length int
	for {
		bytePtr := (*byte)(unsafe.Add(unsafe.Pointer(ptr), length))
		if *bytePtr == 0 {
			break
		}
		length++
	}
	if length == 0 {
		return ""
	}
	bytes := (*[1 << 30]byte)(unsafe.Pointer(ptr))[:length:length]
	return string(bytes)
}
