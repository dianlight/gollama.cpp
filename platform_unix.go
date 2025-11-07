//go:build !windows

package gollama

import (
	"github.com/ebitengine/purego"
)

// loadLibraryPlatform loads a shared library using platform-specific methods
func loadLibraryPlatform(libPath string) (uintptr, error) {
	return purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}

// closeLibraryPlatform closes a shared library using platform-specific methods
func closeLibraryPlatform(handle uintptr) error {
	return purego.Dlclose(handle)
}

// registerLibFunc registers a library function using platform-specific methods
func registerLibFunc(fptr interface{}, handle uintptr, fname string) {
	purego.RegisterLibFunc(fptr, handle, fname)
}

// tryRegisterLibFunc attempts to register a library function, returning an error if it fails
// This is useful for optional functions that may not exist in all library builds
func tryRegisterLibFunc(fptr interface{}, handle uintptr, fname string) error {
	// First check if the symbol exists
	_, err := purego.Dlsym(handle, fname)
	if err != nil {
		return err
	}
	// If it exists, register it
	purego.RegisterLibFunc(fptr, handle, fname)
	return nil
}

// getProcAddressPlatform gets the address of a symbol in a loaded library
func getProcAddressPlatform(handle uintptr, name string) (uintptr, error) {
	return purego.Dlsym(handle, name)
}

// isPlatformSupported returns whether the current platform is supported
func isPlatformSupported() bool {
	return true
}

// getPlatformError returns a platform-specific error message
func getPlatformError() error {
	return nil
}

// clearLoadedDllHandles is a no-op on Unix platforms (only used on Windows)
func clearLoadedDllHandles() {
	// No-op: Unix platforms don't maintain a sibling DLL registry
}
