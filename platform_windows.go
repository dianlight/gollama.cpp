//go:build windows

package gollama

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procLoadLibraryW = kernel32.NewProc("LoadLibraryW")
	procFreeLibrary  = kernel32.NewProc("FreeLibrary")
	//procGetProcAddress = kernel32.NewProc("GetProcAddress")
)

// loadLibraryPlatform loads a shared library using platform-specific methods
func loadLibraryPlatform(libPath string) (uintptr, error) {
	pathPtr, err := syscall.UTF16PtrFromString(libPath)
	if err != nil {
		return 0, fmt.Errorf("failed to convert path to UTF16: %w", err)
	}

	ret, _, err := procLoadLibraryW.Call(uintptr(unsafe.Pointer(pathPtr)))
	if ret == 0 {
		return 0, fmt.Errorf("LoadLibraryW failed: %w", err)
	}

	return ret, nil
}

// closeLibraryPlatform closes a shared library using platform-specific methods
func closeLibraryPlatform(handle uintptr) error {
	ret, _, err := procFreeLibrary.Call(handle)
	if ret == 0 {
		return fmt.Errorf("FreeLibrary failed: %w", err)
	}
	return nil
}

// registerLibFunc registers a library function using platform-specific methods
// For Windows, this is a placeholder implementation
func registerLibFunc(fptr interface{}, handle uintptr, fname string) {
	// TODO: Implement proper function registration for Windows - blocks ROADMAP Priority 1 (Windows Runtime Completion)
	// This would need to use GetProcAddress and set the function pointer
	// For now, this is a no-op to prevent build failures
}

// isPlatformSupported returns whether the current platform is supported
func isPlatformSupported() bool {
	// For now, return false to indicate Windows support is not complete
	// This can be changed to true once full Windows support is implemented
	return false
}

// getPlatformError returns a platform-specific error message
func getPlatformError() error {
	return fmt.Errorf("support for windows platform not yet implemented")
}
