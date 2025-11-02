//go:build windows

package gollama

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procLoadLibraryW             = kernel32.NewProc("LoadLibraryW")
	procLoadLibraryExW           = kernel32.NewProc("LoadLibraryExW")
	procFreeLibrary              = kernel32.NewProc("FreeLibrary")
	procGetProcAddress           = kernel32.NewProc("GetProcAddress")
	procAddDllDirectory          = kernel32.NewProc("AddDllDirectory")
	procRemoveDllDirectory       = kernel32.NewProc("RemoveDllDirectory")
	procSetDefaultDllDirectories = kernel32.NewProc("SetDefaultDllDirectories")
	procSetDllDirectoryW         = kernel32.NewProc("SetDllDirectoryW")
)

// Flags for LoadLibraryEx and SetDefaultDllDirectories
const (
	loadLibrarySearchDllLoadDir  = 0x00000100
	loadLibrarySearchSystem32    = 0x00000800
	loadLibrarySearchDefaultDirs = 0x00001000
	loadLibrarySearchUserDirs    = 0x00000400
)

// loadLibraryPlatform loads a shared library using platform-specific methods
func loadLibraryPlatform(libPath string) (uintptr, error) {
	// Ensure Windows can find dependencies alongside the target DLL by
	// temporarily adding its directory to the DLL search path.
	dir := filepath.Dir(libPath)

	// Try modern safe APIs first: SetDefaultDllDirectories + AddDllDirectory
	var cookie uintptr
	addedDir := false

	if procSetDefaultDllDirectories.Find() == nil {
		// Set search to default dirs + user dirs (added via AddDllDirectory) + System32
		// This avoids using the current working directory and supports side-by-side loading.
		_, _, _ = procSetDefaultDllDirectories.Call(
			uintptr(loadLibrarySearchDefaultDirs | loadLibrarySearchUserDirs | loadLibrarySearchSystem32),
		)
	}

	if procAddDllDirectory.Find() == nil {
		pathPtr, err := syscall.UTF16PtrFromString(dir)
		if err == nil {
			ret, _, _ := procAddDllDirectory.Call(uintptr(unsafe.Pointer(pathPtr)))
			if ret != 0 {
				cookie = ret
				addedDir = true
			}
		}
	}

	// Fallback for older systems: SetDllDirectoryW (process-wide)
	if !addedDir && procSetDllDirectoryW.Find() == nil {
		pathPtr, err := syscall.UTF16PtrFromString(dir)
		if err == nil {
			_, _, _ = procSetDllDirectoryW.Call(uintptr(unsafe.Pointer(pathPtr)))
		}
	}

	pathPtr, err := syscall.UTF16PtrFromString(libPath)
	if err != nil {
		// Best-effort cleanup
		if addedDir && procRemoveDllDirectory.Find() == nil {
			_, _, _ = procRemoveDllDirectory.Call(cookie)
		}
		return 0, fmt.Errorf("failed to convert path to UTF16: %w", err)
	}

	// Prefer LoadLibraryExW with explicit search flags to ensure dependencies
	// in the DLL's directory are discovered reliably.
	if procLoadLibraryExW.Find() == nil {
		ret, _, callErr := procLoadLibraryExW.Call(
			uintptr(unsafe.Pointer(pathPtr)),
			0,
			uintptr(loadLibrarySearchDllLoadDir|loadLibrarySearchDefaultDirs|loadLibrarySearchSystem32|loadLibrarySearchUserDirs),
		)
		if ret != 0 {
			// Cleanup any directory we added
			if addedDir && procRemoveDllDirectory.Find() == nil {
				_, _, _ = procRemoveDllDirectory.Call(cookie)
			}
			return ret, nil
		}
		// If LoadLibraryExW failed, fall back to LoadLibraryW
		_ = callErr
	}

	ret, _, callErr := procLoadLibraryW.Call(uintptr(unsafe.Pointer(pathPtr)))
	if ret == 0 {
		// Cleanup any directory we added before returning
		if addedDir && procRemoveDllDirectory.Find() == nil {
			_, _, _ = procRemoveDllDirectory.Call(cookie)
		}
		return 0, fmt.Errorf("LoadLibraryW failed: %w", callErr)
	}

	// Cleanup any directory we added
	if addedDir && procRemoveDllDirectory.Find() == nil {
		_, _, _ = procRemoveDllDirectory.Call(cookie)
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
// For Windows, this uses GetProcAddress to resolve the function and stores it in the function pointer
func registerLibFunc(fptr interface{}, handle uintptr, fname string) {
	procAddr, err := getProcAddressPlatform(handle, fname)
	if err != nil {
		// Log the error but don't panic - let the caller handle unresolved functions
		fmt.Printf("warning: failed to register %s: %v\n", fname, err)
		return
	}

	// Cast the function pointer interface to a *uintptr and store the resolved address
	// This works because purego uses *uintptr to store function addresses
	if ptr, ok := fptr.(*uintptr); ok {
		*ptr = procAddr
	}
}

// tryRegisterLibFunc attempts to register a library function, returning an error if it fails
// This is useful for optional functions that may not exist in all library builds
func tryRegisterLibFunc(fptr interface{}, handle uintptr, fname string) error {
	procAddr, err := getProcAddressPlatform(handle, fname)
	if err != nil {
		return err
	}

	// Cast the function pointer interface to a *uintptr and store the resolved address
	if ptr, ok := fptr.(*uintptr); ok {
		*ptr = procAddr
	}
	return nil
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
