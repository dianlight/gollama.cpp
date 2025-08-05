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

// isPlatformSupported returns whether the current platform is supported
func isPlatformSupported() bool {
	return true
}

// getPlatformError returns a platform-specific error message
func getPlatformError() error {
	return nil
}
