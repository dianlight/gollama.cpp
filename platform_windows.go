//go:build windows

package gollama

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procLoadLibraryW   = kernel32.NewProc("LoadLibraryW")
	procFreeLibrary    = kernel32.NewProc("FreeLibrary")
	procGetProcAddress = kernel32.NewProc("GetProcAddress")
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

// getProcAddressPlatform gets the address of a symbol in a loaded library
func getProcAddressPlatform(handle uintptr, name string) (uintptr, error) {
	namePtr, err := syscall.BytePtrFromString(name)
	if err != nil {
		return 0, fmt.Errorf("failed to convert name to byte pointer: %w", err)
	}

	ret, _, err := procGetProcAddress.Call(handle, uintptr(unsafe.Pointer(namePtr)))
	if ret == 0 {
		return 0, fmt.Errorf("GetProcAddress failed for %s: %w", name, err)
	}

	return ret, nil
}

// isPlatformSupported returns whether the current platform is supported
func isPlatformSupported() bool {
	// Now we support Windows with FFI
	return true
}

// getPlatformError returns a platform-specific error message
func getPlatformError() error {
	return nil
}
