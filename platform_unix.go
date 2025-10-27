//go:build !windows

package gollama

import (
	"fmt"
	"sync"

	goffi "github.com/clevabit/libgoffi"
)

var (
	libgoffiLibrary *goffi.Library
	libgoffiMutex   sync.RWMutex
)

// loadLibraryPlatform loads a shared library using platform-specific methods
func loadLibraryPlatform(libPath string) (uintptr, error) {
	libgoffiMutex.Lock()
	defer libgoffiMutex.Unlock()

	// Load library using libgoffi
	lib, err := goffi.NewLibrary(libPath, goffi.BindNow|goffi.BindGlobal)
	if err != nil {
		return 0, fmt.Errorf("failed to load library with libgoffi: %w", err)
	}

	libgoffiLibrary = lib
	
	// Return a dummy handle (libgoffi manages the library internally)
	// We use 1 as a non-zero handle to indicate success
	return 1, nil
}

// closeLibraryPlatform closes a shared library using platform-specific methods
func closeLibraryPlatform(handle uintptr) error {
	libgoffiMutex.Lock()
	defer libgoffiMutex.Unlock()

	if libgoffiLibrary != nil {
		err := libgoffiLibrary.Close()
		libgoffiLibrary = nil
		return err
	}
	return nil
}

// registerLibFunc registers a library function using platform-specific methods
func registerLibFunc(fptr interface{}, handle uintptr, fname string) {
	libgoffiMutex.RLock()
	defer libgoffiMutex.RUnlock()

	if libgoffiLibrary == nil {
		// Library not loaded, skip registration
		return
	}

	// Use libgoffi's Import method which automatically maps types
	if err := libgoffiLibrary.Import(fname, fptr); err != nil {
		// Note: Original purego implementation doesn't return errors from RegisterLibFunc
		// For compatibility, we silently ignore errors here as well
		// In production, you may want to log or handle this differently
		_ = err
	}
}

// isPlatformSupported returns whether the current platform is supported
func isPlatformSupported() bool {
	return true
}

// getPlatformError returns a platform-specific error message
func getPlatformError() error {
	return nil
}
