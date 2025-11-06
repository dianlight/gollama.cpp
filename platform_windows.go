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
		ret, _, callErr := procSetDefaultDllDirectories.Call(
			uintptr(loadLibrarySearchDefaultDirs | loadLibrarySearchUserDirs | loadLibrarySearchSystem32),
		)
		if ret == 0 {
			fmt.Printf("warning: SetDefaultDllDirectories failed: %v\n", callErr)
		}
	}

	if procAddDllDirectory.Find() == nil {
		pathPtr, err := syscall.UTF16PtrFromString(dir)
		if err == nil {
			ret, _, callErr := procAddDllDirectory.Call(uintptr(unsafe.Pointer(pathPtr)))
			if ret != 0 {
				cookie = ret
				addedDir = true
				fmt.Printf("debug: Added DLL directory: %s\n", dir)
			} else {
				fmt.Printf("warning: AddDllDirectory failed for %s: %v\n", dir, callErr)
			}
		}
	}

	// Fallback for older systems: SetDllDirectoryW (process-wide)
	if !addedDir && procSetDllDirectoryW.Find() == nil {
		pathPtr, err := syscall.UTF16PtrFromString(dir)
		if err == nil {
			ret, _, callErr := procSetDllDirectoryW.Call(uintptr(unsafe.Pointer(pathPtr)))
			if ret != 0 {
				fmt.Printf("debug: Set DLL directory (fallback): %s\n", dir)
			} else {
				fmt.Printf("warning: SetDllDirectoryW failed for %s: %v\n", dir, callErr)
			}
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

	fmt.Printf("debug: Attempting to load library: %s\n", libPath)

	// Prefer LoadLibraryExW with explicit search flags to ensure dependencies
	// in the DLL's directory are discovered reliably.
	var loadErr error
	if procLoadLibraryExW.Find() == nil {
		ret, _, callErr := procLoadLibraryExW.Call(
			uintptr(unsafe.Pointer(pathPtr)),
			0,
			uintptr(loadLibrarySearchDllLoadDir|loadLibrarySearchDefaultDirs|loadLibrarySearchSystem32|loadLibrarySearchUserDirs),
		)
		if ret != 0 {
			fmt.Printf("debug: Successfully loaded library with LoadLibraryExW: %s (handle: 0x%x)\n", libPath, ret)
			// Cleanup any directory we added
			if addedDir && procRemoveDllDirectory.Find() == nil {
				_, _, _ = procRemoveDllDirectory.Call(cookie)
			}
			return ret, nil
		}
		loadErr = fmt.Errorf("LoadLibraryExW failed for %s: %w (GetLastError: %d)", libPath, callErr, callErr.(syscall.Errno))
		fmt.Printf("debug: %v, trying LoadLibraryW...\n", loadErr)
	}

	ret, _, callErr := procLoadLibraryW.Call(uintptr(unsafe.Pointer(pathPtr)))
	if ret == 0 {
		// Cleanup any directory we added before returning
		if addedDir && procRemoveDllDirectory.Find() == nil {
			_, _, _ = procRemoveDllDirectory.Call(cookie)
		}

		// Build detailed error message
		errno := callErr.(syscall.Errno)
		var errMsg string
		switch errno {
		case 126: // ERROR_MOD_NOT_FOUND
			errMsg = fmt.Sprintf("The specified module could not be found (ERROR_MOD_NOT_FOUND). "+
				"This usually means a dependency DLL is missing. "+
				"Library path: %s, Directory: %s", libPath, dir)
		case 193: // ERROR_BAD_EXE_FORMAT
			errMsg = fmt.Sprintf("The library is not a valid Win32 application (ERROR_BAD_EXE_FORMAT). "+
				"This may indicate an architecture mismatch (e.g., trying to load 64-bit DLL in 32-bit process or vice versa). "+
				"Library path: %s", libPath)
		case 2: // ERROR_FILE_NOT_FOUND
			errMsg = fmt.Sprintf("The system cannot find the file specified (ERROR_FILE_NOT_FOUND). "+
				"Library path: %s", libPath)
		default:
			errMsg = fmt.Sprintf("LoadLibraryW failed for %s: %v (GetLastError: %d)", libPath, callErr, errno)
		}

		if loadErr != nil {
			return 0, fmt.Errorf("%s; Previous attempt: %v", errMsg, loadErr)
		}
		return 0, fmt.Errorf("%s", errMsg)
	}

	fmt.Printf("debug: Successfully loaded library with LoadLibraryW: %s (handle: 0x%x)\n", libPath, ret)

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
		// Log the error with detailed information
		fmt.Printf("warning: failed to register %s: %v (handle: 0x%x)\n", fname, err, handle)
		return
	}

	// Cast the function pointer interface to a *uintptr and store the resolved address
	// This works because purego uses *uintptr to store function addresses
	if ptr, ok := fptr.(*uintptr); ok {
		*ptr = procAddr
		fmt.Printf("debug: Successfully registered function %s at address 0x%x\n", fname, procAddr)
	} else {
		fmt.Printf("warning: failed to cast function pointer for %s (type: %T)\n", fname, fptr)
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
	if handle == 0 {
		return 0, fmt.Errorf("invalid library handle (0) when looking up %s", name)
	}

	namePtr, err := syscall.BytePtrFromString(name)
	if err != nil {
		return 0, fmt.Errorf("failed to convert name to byte pointer: %w", err)
	}

	ret, _, err := procGetProcAddress.Call(handle, uintptr(unsafe.Pointer(namePtr)))
	if ret == 0 {
		errno := err.(syscall.Errno)
		return 0, fmt.Errorf("GetProcAddress failed for %s in library handle 0x%x: %w (GetLastError: %d). "+
			"This means the function is not exported by the DLL, or the DLL may not be the correct version",
			name, handle, err, errno)
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
