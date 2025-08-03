package gollama

import (
	"errors"
	"fmt"
)

// Error types for different categories of errors
var (
	// Library errors
	ErrLibraryNotLoaded   = errors.New("llama.cpp library not loaded")
	ErrLibraryLoadFailed  = errors.New("failed to load llama.cpp library")
	ErrFunctionNotFound   = errors.New("function not found in library")
	ErrInvalidLibraryPath = errors.New("invalid library path")

	// Model errors
	ErrModelNotLoaded       = errors.New("model not loaded")
	ErrModelLoadFailed      = errors.New("failed to load model")
	ErrModelSaveFailed      = errors.New("failed to save model")
	ErrInvalidModelPath     = errors.New("invalid model path")
	ErrModelCorrupted       = errors.New("model file corrupted")
	ErrUnsupportedModelType = errors.New("unsupported model type")

	// Context errors
	ErrContextNotCreated     = errors.New("context not created")
	ErrContextCreationFailed = errors.New("failed to create context")
	ErrInvalidContextSize    = errors.New("invalid context size")
	ErrContextFull           = errors.New("context is full")

	// Token errors
	ErrTokenizationFailed = errors.New("tokenization failed")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenOutOfRange    = errors.New("token out of vocabulary range")

	// Generation errors
	ErrGenerationFailed      = errors.New("text generation failed")
	ErrSamplingFailed        = errors.New("token sampling failed")
	ErrInvalidSamplingParams = errors.New("invalid sampling parameters")

	// Memory errors
	ErrOutOfMemory            = errors.New("out of memory")
	ErrMemoryAllocationFailed = errors.New("memory allocation failed")
	ErrInvalidMemorySize      = errors.New("invalid memory size")

	// Configuration errors
	ErrInvalidConfig          = errors.New("invalid configuration")
	ErrConfigValidationFailed = errors.New("configuration validation failed")
	ErrUnsupportedPlatform    = errors.New("unsupported platform")

	// Backend errors
	ErrBackendNotAvailable = errors.New("backend not available")
	ErrBackendInitFailed   = errors.New("backend initialization failed")
	ErrGPUNotAvailable     = errors.New("GPU not available")
	ErrCUDANotAvailable    = errors.New("CUDA not available")
	ErrMetalNotAvailable   = errors.New("metal backend not available")
	ErrVulkanNotAvailable  = errors.New("vulkan backend not available")

	// File I/O errors
	ErrFileNotFound      = errors.New("file not found")
	ErrFileReadFailed    = errors.New("failed to read file")
	ErrFileWriteFailed   = errors.New("failed to write file")
	ErrInvalidFileFormat = errors.New("invalid file format")

	// Parameter errors
	ErrInvalidParameter    = errors.New("invalid parameter")
	ErrParameterOutOfRange = errors.New("parameter out of range")
	ErrMissingParameter    = errors.New("missing required parameter")

	// Thread/concurrency errors
	ErrThreadingFailed      = errors.New("threading operation failed")
	ErrConcurrencyViolation = errors.New("concurrency violation")
	ErrDeadlock             = errors.New("deadlock detected")
)

// LlamaError represents a structured error from the llama.cpp library
type LlamaError struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Function string `json:"function,omitempty"`
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Cause    error  `json:"cause,omitempty"`
}

// Error implements the error interface
func (e *LlamaError) Error() string {
	if e.Function != "" {
		return fmt.Sprintf("llama error [%d]: %s (in %s)", e.Code, e.Message, e.Function)
	}
	return fmt.Sprintf("llama error [%d]: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause
func (e *LlamaError) Unwrap() error {
	return e.Cause
}

// Is implements error matching
func (e *LlamaError) Is(target error) bool {
	if t, ok := target.(*LlamaError); ok {
		return e.Code == t.Code
	}
	return errors.Is(e.Cause, target)
}

// NewLlamaError creates a new LlamaError
func NewLlamaError(code int, message string) *LlamaError {
	return &LlamaError{
		Code:    code,
		Message: message,
	}
}

// NewLlamaErrorWithCause creates a new LlamaError with an underlying cause
func NewLlamaErrorWithCause(code int, message string, cause error) *LlamaError {
	return &LlamaError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewLlamaErrorWithContext creates a new LlamaError with function context
func NewLlamaErrorWithContext(code int, message, function string) *LlamaError {
	return &LlamaError{
		Code:     code,
		Message:  message,
		Function: function,
	}
}

// Error codes matching llama.cpp return values
const (
	LLAMA_ERR_SUCCESS         = 0
	LLAMA_ERR_FAIL            = -1
	LLAMA_ERR_INVALID_PARAM   = -2
	LLAMA_ERR_OUT_OF_MEMORY   = -3
	LLAMA_ERR_FILE_NOT_FOUND  = -4
	LLAMA_ERR_FILE_READ       = -5
	LLAMA_ERR_FILE_WRITE      = -6
	LLAMA_ERR_INVALID_FORMAT  = -7
	LLAMA_ERR_UNSUPPORTED     = -8
	LLAMA_ERR_BACKEND_INIT    = -9
	LLAMA_ERR_CONTEXT_FULL    = -10
	LLAMA_ERR_TOKEN_INVALID   = -11
	LLAMA_ERR_MODEL_CORRUPTED = -12
	LLAMA_ERR_GPU_UNAVAILABLE = -13
)

// ErrorFromCode converts a llama.cpp error code to a Go error
func ErrorFromCode(code int) error {
	switch code {
	case LLAMA_ERR_SUCCESS:
		return nil
	case LLAMA_ERR_FAIL:
		return NewLlamaError(code, "operation failed")
	case LLAMA_ERR_INVALID_PARAM:
		return NewLlamaErrorWithCause(code, "invalid parameter", ErrInvalidParameter)
	case LLAMA_ERR_OUT_OF_MEMORY:
		return NewLlamaErrorWithCause(code, "out of memory", ErrOutOfMemory)
	case LLAMA_ERR_FILE_NOT_FOUND:
		return NewLlamaErrorWithCause(code, "file not found", ErrFileNotFound)
	case LLAMA_ERR_FILE_READ:
		return NewLlamaErrorWithCause(code, "file read error", ErrFileReadFailed)
	case LLAMA_ERR_FILE_WRITE:
		return NewLlamaErrorWithCause(code, "file write error", ErrFileWriteFailed)
	case LLAMA_ERR_INVALID_FORMAT:
		return NewLlamaErrorWithCause(code, "invalid file format", ErrInvalidFileFormat)
	case LLAMA_ERR_UNSUPPORTED:
		return NewLlamaErrorWithCause(code, "unsupported operation", ErrUnsupportedPlatform)
	case LLAMA_ERR_BACKEND_INIT:
		return NewLlamaErrorWithCause(code, "backend initialization failed", ErrBackendInitFailed)
	case LLAMA_ERR_CONTEXT_FULL:
		return NewLlamaErrorWithCause(code, "context is full", ErrContextFull)
	case LLAMA_ERR_TOKEN_INVALID:
		return NewLlamaErrorWithCause(code, "invalid token", ErrInvalidToken)
	case LLAMA_ERR_MODEL_CORRUPTED:
		return NewLlamaErrorWithCause(code, "model file corrupted", ErrModelCorrupted)
	case LLAMA_ERR_GPU_UNAVAILABLE:
		return NewLlamaErrorWithCause(code, "GPU not available", ErrGPUNotAvailable)
	default:
		return NewLlamaError(code, fmt.Sprintf("unknown error code: %d", code))
	}
}

// ErrorHandler provides centralized error handling and logging
type ErrorHandler struct {
	enableLogging bool
	logCallback   func(level int, message string)
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(enableLogging bool) *ErrorHandler {
	return &ErrorHandler{
		enableLogging: enableLogging,
	}
}

// SetLogCallback sets the log callback function
func (eh *ErrorHandler) SetLogCallback(callback func(level int, message string)) {
	eh.logCallback = callback
}

// HandleError processes and logs an error
func (eh *ErrorHandler) HandleError(err error, context string) error {
	if err == nil {
		return nil
	}

	if eh.enableLogging && eh.logCallback != nil {
		message := fmt.Sprintf("Error in %s: %v", context, err)
		eh.logCallback(3, message) // LLAMA_LOG_LEVEL_ERROR
	}

	// Wrap with context if it's not already a LlamaError
	if _, ok := err.(*LlamaError); !ok {
		return NewLlamaErrorWithContext(-1, err.Error(), context)
	}

	return err
}

// WrapError wraps an error with additional context
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// WrapErrorf wraps an error with formatted additional context
func WrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific retryable errors
	return errors.Is(err, ErrOutOfMemory) ||
		errors.Is(err, ErrFileReadFailed) ||
		errors.Is(err, ErrBackendInitFailed)
}

// IsFatalError checks if an error is fatal and should stop execution
func IsFatalError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific fatal errors
	return errors.Is(err, ErrModelCorrupted) ||
		errors.Is(err, ErrInvalidFileFormat) ||
		errors.Is(err, ErrUnsupportedPlatform) ||
		errors.Is(err, ErrLibraryNotLoaded)
}

// ErrorCategory represents different categories of errors
type ErrorCategory int

const (
	CategoryLibrary ErrorCategory = iota
	CategoryModel
	CategoryContext
	CategoryToken
	CategoryGeneration
	CategoryMemory
	CategoryConfig
	CategoryBackend
	CategoryFile
	CategoryParameter
	CategoryThread
)

// String returns the string representation of an error category
func (ec ErrorCategory) String() string {
	switch ec {
	case CategoryLibrary:
		return "Library"
	case CategoryModel:
		return "Model"
	case CategoryContext:
		return "Context"
	case CategoryToken:
		return "Token"
	case CategoryGeneration:
		return "Generation"
	case CategoryMemory:
		return "Memory"
	case CategoryConfig:
		return "Config"
	case CategoryBackend:
		return "Backend"
	case CategoryFile:
		return "File"
	case CategoryParameter:
		return "Parameter"
	case CategoryThread:
		return "Thread"
	default:
		return "Unknown"
	}
}

// CategorizeError determines the category of an error
func CategorizeError(err error) ErrorCategory {
	if err == nil {
		return CategoryLibrary // Default category
	}

	switch {
	case errors.Is(err, ErrLibraryNotLoaded), errors.Is(err, ErrLibraryLoadFailed):
		return CategoryLibrary
	case errors.Is(err, ErrModelNotLoaded), errors.Is(err, ErrModelLoadFailed):
		return CategoryModel
	case errors.Is(err, ErrContextNotCreated), errors.Is(err, ErrContextCreationFailed):
		return CategoryContext
	case errors.Is(err, ErrTokenizationFailed), errors.Is(err, ErrInvalidToken):
		return CategoryToken
	case errors.Is(err, ErrGenerationFailed), errors.Is(err, ErrSamplingFailed):
		return CategoryGeneration
	case errors.Is(err, ErrOutOfMemory), errors.Is(err, ErrMemoryAllocationFailed):
		return CategoryMemory
	case errors.Is(err, ErrInvalidConfig), errors.Is(err, ErrConfigValidationFailed):
		return CategoryConfig
	case errors.Is(err, ErrBackendNotAvailable), errors.Is(err, ErrGPUNotAvailable):
		return CategoryBackend
	case errors.Is(err, ErrFileNotFound), errors.Is(err, ErrFileReadFailed):
		return CategoryFile
	case errors.Is(err, ErrInvalidParameter), errors.Is(err, ErrParameterOutOfRange):
		return CategoryParameter
	case errors.Is(err, ErrThreadingFailed), errors.Is(err, ErrConcurrencyViolation):
		return CategoryThread
	default:
		return CategoryLibrary
	}
}

// Global error handler instance
var globalErrorHandler = NewErrorHandler(true)

// SetGlobalErrorHandler sets the global error handler
func SetGlobalErrorHandler(handler *ErrorHandler) {
	globalErrorHandler = handler
}

// GetGlobalErrorHandler returns the global error handler
func GetGlobalErrorHandler() *ErrorHandler {
	return globalErrorHandler
}

// HandleError is a convenience function that uses the global error handler
func HandleError(err error, context string) error {
	return globalErrorHandler.HandleError(err, context)
}

// CheckResult checks a result code and returns an appropriate error
func CheckResult(result int, function string) error {
	if result == LLAMA_ERR_SUCCESS {
		return nil
	}

	err := ErrorFromCode(result)
	if llamaErr, ok := err.(*LlamaError); ok {
		llamaErr.Function = function
	}

	return HandleError(err, function)
}

// Must panics if the error is not nil - use only for initialization
func Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("gollama initialization failed: %v", err))
	}
}

// Try executes a function and handles any panics as errors
func Try(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case error:
				err = x
			case string:
				err = errors.New(x)
			default:
				err = fmt.Errorf("unknown panic: %v", r)
			}
		}
	}()

	return fn()
}
