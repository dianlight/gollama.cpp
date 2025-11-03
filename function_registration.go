package gollama

import (
	"fmt"
	"log/slog"
	"runtime"
)

// FunctionRegistration holds metadata for a single function to register
type FunctionRegistration struct {
	FunctionPtr interface{} // Pointer to the function pointer variable
	Name        string      // C function name
	Required    bool        // Whether registration must succeed
	OnlyDarwin  bool        // Register only on Darwin platform
}

// RegisterFunctions registers multiple function pointers with error handling
// If any required function fails to register and we're not in try mode, it returns an error
// This consolidates repeated registerLibFunc calls and reduces code duplication
func RegisterFunctionSet(libHandle uintptr, functions []FunctionRegistration, tryMode bool) error {
	for _, fn := range functions {
		// Skip Darwin-only functions on non-Darwin platforms
		if fn.OnlyDarwin {
			if runtime.GOOS != "darwin" {
				continue
			}
		}

		if tryMode {
			// In try mode, ignore registration errors (for optional functions)
			_ = tryRegisterLibFunc(fn.FunctionPtr, libHandle, fn.Name)
		} else {
			// In strict mode, fail on any registration error if required
			if fn.Required {
				registerLibFunc(fn.FunctionPtr, libHandle, fn.Name)
			} else {
				_ = tryRegisterLibFunc(fn.FunctionPtr, libHandle, fn.Name)
			}
		}
	}
	return nil
}

// BatchRegisterFunctions is a convenience wrapper that registers functions in batch
// and returns true if all required functions were registered successfully
func BatchRegisterFunctions(libHandle uintptr, functions []FunctionRegistration, tryMode bool) error {
	return RegisterFunctionSet(libHandle, functions, tryMode)
}

// LibraryLoadInfo consolidates the result of a library load attempt
type LibraryLoadInfo struct {
	Success bool
	Handle  uintptr
	Path    string
	Error   string
}

// LoadLibraryWithDependencies encapsulates the common pattern of:
// 1. Preloading dependent libraries
// 2. Loading the main library
// 3. Setting up configuration on success
// This is used in multiple places in LoadLibraryWithVersion to reduce duplication
func (l *LibraryLoader) LoadLibraryWithDependencies(libPath string) (*LibraryLoadInfo, []string) {
	var reasons []string

	if err := l.preloadDependentLibraries(libPath); err != nil {
		reasons = append(reasons, fmt.Sprintf("preload failed: %v", err))
		return &LibraryLoadInfo{Success: false}, reasons
	}

	handle, err := l.loadSharedLibrary(libPath)
	if err != nil {
		reasons = append(reasons, fmt.Sprintf("dlopen failed: %v", err))
		return &LibraryLoadInfo{Success: false}, reasons
	}

	return &LibraryLoadInfo{
		Success: true,
		Handle:  handle,
		Path:    libPath,
	}, reasons
}

// ApplyLibraryLoad applies the result of a successful library load to the loader state
func (l *LibraryLoader) ApplyLibraryLoad(info *LibraryLoadInfo, rootPath string) error {
	l.handle = info.Handle
	l.llamaLibPath = info.Path
	l.rootLibPath = rootPath
	l.loaded = true

	suffix, err := getExpectedLibrarySuffix()
	if err != nil {
		slog.Warn("Failed to get expected library suffix", "error", err)
	}
	l.extensionSuffix = suffix
	return nil
}
