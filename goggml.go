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
type GgmlGuid [16]byte // ggml_guid_t

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

// GGML backend device types
type GgmlBackendDevType int32

const (
	GGML_BACKEND_DEVICE_TYPE_CPU   GgmlBackendDevType = 0
	GGML_BACKEND_DEVICE_TYPE_GPU   GgmlBackendDevType = 1
	GGML_BACKEND_DEVICE_TYPE_IGPU  GgmlBackendDevType = 2
	GGML_BACKEND_DEVICE_TYPE_ACCEL GgmlBackendDevType = 3
)

// GGML backend device capabilities
type GgmlBackendDevCaps struct {
	Async             bool // asynchronous operations
	HostBuffer        bool // pinned host buffer
	BufferFromHostPtr bool // creating buffers from host ptr
	Events            bool // event synchronization
}

// GGML backend device properties
type GgmlBackendDevProps struct {
	Name        string             // device name
	Description string             // device description
	MemoryFree  uint64             // device free memory in bytes
	MemoryTotal uint64             // device total memory in bytes
	Type        GgmlBackendDevType // device type
	DeviceID    string             // device id (e.g., PCI bus id)
	Caps        GgmlBackendDevCaps // device capabilities
}

// Function pointers for GGML functions
var (
	// Type size functions
	ggmlTypeSize    func(typ GgmlType) uint64
	ggmlTypeSizeof  func(typ GgmlType) uint64
	ggmlBlckSize    func(typ GgmlType) int32
	ggmlIsQuantized func(typ GgmlType) bool

	// Backend device functions
	ggmlBackendDevCount             func() uint64
	ggmlBackendDevGet               func(index uint64) GgmlBackendDevice
	ggmlBackendDevByName            func(name *byte) GgmlBackendDevice
	ggmlBackendDevByType            func(typ int32) GgmlBackendDevice
	ggmlBackendDevInit              func(device GgmlBackendDevice, params *byte) GgmlBackend
	ggmlBackendDevName              func(device GgmlBackendDevice) *byte
	ggmlBackendDevDescription       func(device GgmlBackendDevice) *byte
	ggmlBackendDevMemory            func(device GgmlBackendDevice, free *uint64, total *uint64)
	ggmlBackendDevType              func(device GgmlBackendDevice) int32
	ggmlBackendDevGetProps          func(device GgmlBackendDevice, props unsafe.Pointer)
	ggmlBackendDevBackendReg        func(device GgmlBackendDevice) GgmlBackendReg
	ggmlBackendDevBufferFromHostPtr func(device GgmlBackendDevice, ptr unsafe.Pointer, size uint64, maxTensorSize uint64) GgmlBackendBuffer
	ggmlBackendDevSupportsOp        func(device GgmlBackendDevice, op GgmlTensor) bool
	ggmlBackendDevSupportsBuft      func(device GgmlBackendDevice, buft GgmlBackendBufferType) bool
	ggmlBackendDevOffloadOp         func(device GgmlBackendDevice, op GgmlTensor) bool

	// Backend buffer type functions
	ggmlBackendDevBufferType     func(device GgmlBackendDevice) GgmlBackendBufferType
	ggmlBackendDevHostBufferType func(device GgmlBackendDevice) GgmlBackendBufferType
	ggmlBackendCpuBufferType     func() GgmlBackendBufferType
	ggmlBackendBuftName          func(buft GgmlBackendBufferType) *byte
	ggmlBackendBuftAllocBuffer   func(buft GgmlBackendBufferType, size uint64) GgmlBackendBuffer
	ggmlBackendBuftGetAlignment  func(buft GgmlBackendBufferType) uint64
	ggmlBackendBuftGetMaxSize    func(buft GgmlBackendBufferType) uint64
	ggmlBackendBuftGetAllocSize  func(buft GgmlBackendBufferType, tensor GgmlTensor) uint64
	ggmlBackendBuftIsHost        func(buft GgmlBackendBufferType) bool
	ggmlBackendBuftGetDevice     func(buft GgmlBackendBufferType) GgmlBackendDevice

	// Backend buffer functions
	ggmlBackendBufferFree         func(buffer GgmlBackendBuffer)
	ggmlBackendBufferGetBase      func(buffer GgmlBackendBuffer) unsafe.Pointer
	ggmlBackendBufferGetSize      func(buffer GgmlBackendBuffer) uint64
	ggmlBackendBufferInitTensor   func(buffer GgmlBackendBuffer, tensor GgmlTensor) int32 // enum ggml_status
	ggmlBackendBufferGetAlignment func(buffer GgmlBackendBuffer) uint64
	ggmlBackendBufferGetMaxSize   func(buffer GgmlBackendBuffer) uint64
	ggmlBackendBufferGetAllocSize func(buffer GgmlBackendBuffer, tensor GgmlTensor) uint64
	ggmlBackendBufferClear        func(buffer GgmlBackendBuffer, value uint8)
	ggmlBackendBufferIsHost       func(buffer GgmlBackendBuffer) bool
	ggmlBackendBufferSetUsage     func(buffer GgmlBackendBuffer, usage int32)
	ggmlBackendBufferGetUsage     func(buffer GgmlBackendBuffer) int32
	ggmlBackendBufferGetType      func(buffer GgmlBackendBuffer) GgmlBackendBufferType
	ggmlBackendBufferName         func(buffer GgmlBackendBuffer) *byte
	ggmlBackendBufferReset        func(buffer GgmlBackendBuffer)

	// Backend functions
	ggmlBackendGuid                 func(backend GgmlBackend, guid *GgmlGuid)
	ggmlBackendFree                 func(backend GgmlBackend)
	ggmlBackendName                 func(backend GgmlBackend) *byte
	ggmlBackendGetDefaultBufferType func(backend GgmlBackend) GgmlBackendBufferType
	ggmlBackendAllocBuffer          func(backend GgmlBackend, size uint64) GgmlBackendBuffer
	ggmlBackendGetAlignment         func(backend GgmlBackend) uint64
	ggmlBackendGetMaxSize           func(backend GgmlBackend) uint64
	ggmlBackendSupports             func(backend GgmlBackend, buft GgmlBackendBufferType) bool
	ggmlBackendGetDevice            func(backend GgmlBackend) GgmlBackendDevice
	ggmlBackendInitBest             func() GgmlBackend
	ggmlBackendInitByName           func(name *byte, params *byte) GgmlBackend
	ggmlBackendInitByType           func(typ int32, params *byte) GgmlBackend

	// Backend registry functions
	ggmlBackendRegName           func(reg GgmlBackendReg) *byte
	ggmlBackendRegDevCount       func(reg GgmlBackendReg) uint64
	ggmlBackendRegDevGet         func(reg GgmlBackendReg, index uint64) GgmlBackendDevice
	ggmlBackendRegGetProcAddress func(reg GgmlBackendReg, name *byte) unsafe.Pointer
	ggmlBackendRegister          func(reg GgmlBackendReg)
	ggmlBackendDeviceRegister    func(device GgmlBackendDevice)
	ggmlBackendRegCount          func() uint64
	ggmlBackendRegGet            func(index uint64) GgmlBackendReg
	ggmlBackendRegByName         func(name *byte) GgmlBackendReg

	// Backend loading functions
	ggmlBackendLoad            func(path *byte) GgmlBackendReg
	ggmlBackendUnload          func(reg GgmlBackendReg)
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

	// Backend device functions (extended)
	_ = tryRegisterLibFunc(&ggmlBackendDevByName, libHandle, "ggml_backend_dev_by_name")
	_ = tryRegisterLibFunc(&ggmlBackendDevType, libHandle, "ggml_backend_dev_type")
	_ = tryRegisterLibFunc(&ggmlBackendDevGetProps, libHandle, "ggml_backend_dev_get_props")
	_ = tryRegisterLibFunc(&ggmlBackendDevBackendReg, libHandle, "ggml_backend_dev_backend_reg")
	_ = tryRegisterLibFunc(&ggmlBackendDevBufferFromHostPtr, libHandle, "ggml_backend_dev_buffer_from_host_ptr")
	_ = tryRegisterLibFunc(&ggmlBackendDevSupportsOp, libHandle, "ggml_backend_dev_supports_op")
	_ = tryRegisterLibFunc(&ggmlBackendDevSupportsBuft, libHandle, "ggml_backend_dev_supports_buft")
	_ = tryRegisterLibFunc(&ggmlBackendDevOffloadOp, libHandle, "ggml_backend_dev_offload_op")

	// Backend buffer type functions (extended)
	_ = tryRegisterLibFunc(&ggmlBackendBuftGetAlignment, libHandle, "ggml_backend_buft_get_alignment")
	_ = tryRegisterLibFunc(&ggmlBackendBuftGetMaxSize, libHandle, "ggml_backend_buft_get_max_size")
	_ = tryRegisterLibFunc(&ggmlBackendBuftGetAllocSize, libHandle, "ggml_backend_buft_get_alloc_size")
	_ = tryRegisterLibFunc(&ggmlBackendBuftGetDevice, libHandle, "ggml_backend_buft_get_device")

	// Backend buffer functions (extended)
	_ = tryRegisterLibFunc(&ggmlBackendBufferInitTensor, libHandle, "ggml_backend_buffer_init_tensor")
	_ = tryRegisterLibFunc(&ggmlBackendBufferGetAlignment, libHandle, "ggml_backend_buffer_get_alignment")
	_ = tryRegisterLibFunc(&ggmlBackendBufferGetMaxSize, libHandle, "ggml_backend_buffer_get_max_size")
	_ = tryRegisterLibFunc(&ggmlBackendBufferGetAllocSize, libHandle, "ggml_backend_buffer_get_alloc_size")
	_ = tryRegisterLibFunc(&ggmlBackendBufferGetUsage, libHandle, "ggml_backend_buffer_get_usage")
	_ = tryRegisterLibFunc(&ggmlBackendBufferReset, libHandle, "ggml_backend_buffer_reset")

	// Backend functions
	_ = tryRegisterLibFunc(&ggmlBackendGuid, libHandle, "ggml_backend_guid")
	_ = tryRegisterLibFunc(&ggmlBackendFree, libHandle, "ggml_backend_free")
	_ = tryRegisterLibFunc(&ggmlBackendName, libHandle, "ggml_backend_name")
	_ = tryRegisterLibFunc(&ggmlBackendSupports, libHandle, "ggml_backend_supports_buft")
	_ = tryRegisterLibFunc(&ggmlBackendGetDefaultBufferType, libHandle, "ggml_backend_get_default_buffer_type")
	_ = tryRegisterLibFunc(&ggmlBackendAllocBuffer, libHandle, "ggml_backend_alloc_buffer")
	_ = tryRegisterLibFunc(&ggmlBackendGetAlignment, libHandle, "ggml_backend_get_alignment")
	_ = tryRegisterLibFunc(&ggmlBackendGetMaxSize, libHandle, "ggml_backend_get_max_size")
	_ = tryRegisterLibFunc(&ggmlBackendGetDevice, libHandle, "ggml_backend_get_device")
	_ = tryRegisterLibFunc(&ggmlBackendInitBest, libHandle, "ggml_backend_init_best")
	_ = tryRegisterLibFunc(&ggmlBackendInitByName, libHandle, "ggml_backend_init_by_name")
	_ = tryRegisterLibFunc(&ggmlBackendInitByType, libHandle, "ggml_backend_init_by_type")
	_ = tryRegisterLibFunc(&ggmlBackendLoad, libHandle, "ggml_backend_load")
	_ = tryRegisterLibFunc(&ggmlBackendUnload, libHandle, "ggml_backend_unload")
	_ = tryRegisterLibFunc(&ggmlBackendLoadAll, libHandle, "ggml_backend_load_all")
	_ = tryRegisterLibFunc(&ggmlBackendLoadAllFromPath, libHandle, "ggml_backend_load_all_from_path")

	// Backend registry functions
	_ = tryRegisterLibFunc(&ggmlBackendRegName, libHandle, "ggml_backend_reg_name")
	_ = tryRegisterLibFunc(&ggmlBackendRegDevCount, libHandle, "ggml_backend_reg_dev_count")
	_ = tryRegisterLibFunc(&ggmlBackendRegDevGet, libHandle, "ggml_backend_reg_dev_get")
	_ = tryRegisterLibFunc(&ggmlBackendRegGetProcAddress, libHandle, "ggml_backend_reg_get_proc_address")
	_ = tryRegisterLibFunc(&ggmlBackendRegister, libHandle, "ggml_backend_register")
	_ = tryRegisterLibFunc(&ggmlBackendDeviceRegister, libHandle, "ggml_backend_device_register")
	_ = tryRegisterLibFunc(&ggmlBackendRegCount, libHandle, "ggml_backend_reg_count")
	_ = tryRegisterLibFunc(&ggmlBackendRegGet, libHandle, "ggml_backend_reg_get")
	_ = tryRegisterLibFunc(&ggmlBackendRegByName, libHandle, "ggml_backend_reg_by_name")

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

// Ggml_backend_init_best initializes the best available backend (GPU or CPU)
func Ggml_backend_init_best() (GgmlBackend, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlBackendInitBest == nil {
		return 0, fmt.Errorf("ggml_backend_init_best function not available")
	}
	backend := ggmlBackendInitBest()
	if backend == 0 {
		return 0, fmt.Errorf("failed to initialize best backend")
	}
	return backend, nil
}

// Ggml_backend_init_by_name initializes a backend by name with optional parameters
func Ggml_backend_init_by_name(name string, params string) (GgmlBackend, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlBackendInitByName == nil {
		return 0, fmt.Errorf("ggml_backend_init_by_name function not available")
	}

	nameBytes := append([]byte(name), 0)
	var paramsPtr *byte
	if params != "" {
		paramsBytes := append([]byte(params), 0)
		paramsPtr = &paramsBytes[0]
	}

	backend := ggmlBackendInitByName(&nameBytes[0], paramsPtr)
	if backend == 0 {
		return 0, fmt.Errorf("failed to initialize backend by name: %s", name)
	}
	return backend, nil
}

// Ggml_backend_init_by_type initializes a backend by device type with optional parameters
func Ggml_backend_init_by_type(deviceType GgmlBackendDevType, params string) (GgmlBackend, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlBackendInitByType == nil {
		return 0, fmt.Errorf("ggml_backend_init_by_type function not available")
	}

	var paramsPtr *byte
	if params != "" {
		paramsBytes := append([]byte(params), 0)
		paramsPtr = &paramsBytes[0]
	}

	backend := ggmlBackendInitByType(int32(deviceType), paramsPtr)
	if backend == 0 {
		return 0, fmt.Errorf("failed to initialize backend by type: %d", deviceType)
	}
	return backend, nil
}

// Ggml_backend_load dynamically loads a backend from a library path and returns a backend registry
func Ggml_backend_load(path string) (GgmlBackendReg, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if ggmlBackendLoad == nil {
		return 0, fmt.Errorf("ggml_backend_load function not available")
	}

	if globalLoader.rootLibPath == "" {
		err := globalLoader.LoadLibrary()
		if err != nil {
			return 0, fmt.Errorf("failed to load library for backend loading: %v", err)
		}
	}

	pathBytes := append([]byte(path), 0)
	reg := ggmlBackendLoad(&pathBytes[0])
	if reg == 0 {
		return 0, fmt.Errorf("failed to load backend from path: %s", path)
	}
	return reg, nil
}

// Ggml_backend_unload unloads a dynamically loaded backend and unregisters it
func Ggml_backend_unload(reg GgmlBackendReg) error {
	if err := ensureLoaded(); err != nil {
		return err
	}
	if ggmlBackendUnload == nil {
		return fmt.Errorf("ggml_backend_unload function not available")
	}

	ggmlBackendUnload(reg)
	return nil
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
