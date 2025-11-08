//go:build windows

package gollama

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/ebitengine/purego"
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

// keep a small registry of loaded DLL handles from the target directory so we can
// resolve symbols that might be exported by sibling DLLs (e.g., ggml.dll)
var loadedDllHandles []uintptr

// addLoadedHandle saves a successfully loaded DLL handle for later symbol lookup
func addLoadedHandle(h uintptr) {
	// avoid duplicates and nil
	if h == 0 {
		return
	}
	for _, existing := range loadedDllHandles {
		if existing == h {
			return
		}
	}
	loadedDllHandles = append(loadedDllHandles, h)
}

// clearLoadedDllHandles clears the registry of loaded DLL handles
// This should be called when unloading the library to avoid stale handles
func clearLoadedDllHandles() {
	loadedDllHandles = nil
}

// Flags for LoadLibraryEx and SetDefaultDllDirectories
const (
	loadLibrarySearchDllLoadDir  = 0x00000100
	loadLibrarySearchSystem32    = 0x00000800
	loadLibrarySearchDefaultDirs = 0x00001000
	loadLibrarySearchUserDirs    = 0x00000400
)

// loadLibraryPlatform loads a shared library using platform-specific methods
func loadLibraryPlatform(libPath string) (uintptr, error) {
	slog.Debug("loadLibraryPlatform: starting library load", "path", libPath)

	// Ensure Windows can find dependencies alongside the target DLL by
	// temporarily adding its directory to the DLL search path.
	dir := filepath.Dir(libPath)
	slog.Debug("loadLibraryPlatform: DLL directory", "dir", dir)

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
			slog.Warn("loadLibraryPlatform: SetDefaultDllDirectories failed", "error", callErr)
		} else {
			slog.Debug("loadLibraryPlatform: SetDefaultDllDirectories succeeded")
		}
	}

	if procAddDllDirectory.Find() == nil {
		pathPtr, err := syscall.UTF16PtrFromString(dir)
		if err == nil {
			ret, _, callErr := procAddDllDirectory.Call(uintptr(unsafe.Pointer(pathPtr)))
			if ret != 0 {
				cookie = ret
				addedDir = true
				slog.Debug("loadLibraryPlatform: Added DLL directory via AddDllDirectory", "dir", dir, "cookie", fmt.Sprintf("0x%x", cookie))
			} else {
				slog.Warn("loadLibraryPlatform: AddDllDirectory failed", "dir", dir, "error", callErr)
			}
		}
	}

	// Fallback for older systems: SetDllDirectoryW (process-wide)
	if !addedDir && procSetDllDirectoryW.Find() == nil {
		pathPtr, err := syscall.UTF16PtrFromString(dir)
		if err == nil {
			ret, _, callErr := procSetDllDirectoryW.Call(uintptr(unsafe.Pointer(pathPtr)))
			if ret != 0 {
				slog.Debug("loadLibraryPlatform: Set DLL directory (fallback SetDllDirectoryW)", "dir", dir)
			} else {
				slog.Warn("loadLibraryPlatform: SetDllDirectoryW failed", "dir", dir, "error", callErr)
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

	slog.Debug("loadLibraryPlatform: attempting to load library with LoadLibraryExW", "path", libPath)

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
			slog.Debug("loadLibraryPlatform: Successfully loaded library with LoadLibraryExW", "path", libPath, "handle", fmt.Sprintf("0x%x", ret))
			// Cleanup any directory we added
			if addedDir && procRemoveDllDirectory.Find() == nil {
				_, _, _ = procRemoveDllDirectory.Call(cookie)
			}
			// Also try to proactively load sibling DLLs from the same directory to ensure
			// all exports are available (some symbols may live in ggml*.dll on Windows).
			slog.Debug("loadLibraryPlatform: preloading sibling DLLs", "dir", dir)
			preloadSiblingDlls(dir, ret)
			return ret, nil
		}
		loadErr = fmt.Errorf("LoadLibraryExW failed for %s: %w (GetLastError: %d)", libPath, callErr, callErr.(syscall.Errno))
		slog.Debug("loadLibraryPlatform: LoadLibraryExW failed, trying LoadLibraryW", "error", loadErr)
	}

	slog.Debug("loadLibraryPlatform: attempting to load library with LoadLibraryW (fallback)", "path", libPath)

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

	slog.Debug("loadLibraryPlatform: Successfully loaded library with LoadLibraryW", "path", libPath, "handle", fmt.Sprintf("0x%x", ret))

	// Cleanup any directory we added
	if addedDir && procRemoveDllDirectory.Find() == nil {
		_, _, _ = procRemoveDllDirectory.Call(cookie)
	}

	// Proactively load sibling DLLs from the same directory
	slog.Debug("loadLibraryPlatform: preloading sibling DLLs", "dir", dir)
	preloadSiblingDlls(dir, ret)

	return ret, nil
}

// preloadSiblingDlls loads other DLLs from the same directory that commonly contain
// exports used by llama.dll (e.g., ggml*.dll). This improves GetProcAddress success
// on setups where functions are exported by a different module.
// The allowlist ensures critical DLLs like ggml-base.dll are loaded first, before
// searching for symbols, as they may contain core functionality like ggml_backend_cpu_buffer_type.
func preloadSiblingDlls(dir string, mainHandle uintptr) {
	// Track the main handle
	addLoadedHandle(mainHandle)
	slog.Debug("preloadSiblingDlls: starting DLL preload", "directory", dir, "mainHandle", fmt.Sprintf("0x%x", mainHandle))

	// Scan directory for DLLs and load a short allowlist first, then best-effort all *.dll
	// Priority list of likely dependencies - ORDER MATTERS as some must be loaded before others
	// ggml-base.dll is listed early because it exports core functionality like ggml_backend_cpu_buffer_type
	allowlist := []string{
		"ggml-base.dll",    // Core GGML functionality - MUST be loaded before other GGML modules
		"ggml.dll",         // Main GGML library
		"ggml-cpu-x64.dll", // Generic x64 CPU backend (replaces ggml-cpu.dll which doesn't exist)
		"ggml-blas.dll",    // BLAS backend
		"ggml-rpc.dll",     // RPC backend
		"ggml-cuda.dll",    // CUDA backend
		"ggml-metal.dll",   // Metal backend (macOS/iOS)
		"ggml-kompute.dll", // Kompute backend
		"ggml-sycl.dll",    // SYCL backend
	}

	slog.Debug("preloadSiblingDlls: loading allowlisted DLLs", "count", len(allowlist))
	for _, name := range allowlist {
		dllPath := filepath.Join(dir, name)
		if _, err := os.Stat(dllPath); err == nil {
			slog.Debug("preloadSiblingDlls: found allowlisted DLL", "name", name, "path", dllPath)
			if h, err := loadOneDll(dllPath); err == nil {
				addLoadedHandle(h)
				slog.Debug("preloadSiblingDlls: successfully loaded DLL", "name", name, "handle", fmt.Sprintf("0x%x", h))
			} else {
				slog.Warn("preloadSiblingDlls: failed to load allowlisted DLL", "name", name, "error", err)
			}
		}
	}

	// Best-effort: load remaining DLLs in the directory (skip those already loaded)
	slog.Debug("preloadSiblingDlls: scanning directory for additional DLLs", "directory", dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		slog.Warn("preloadSiblingDlls: failed to read directory", "directory", dir, "error", err)
		return
	}

	loadedCount := 0
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".dll" {
			continue
		}
		name := e.Name()
		// Skip main llama.dll; we already have it
		if name == "llama.dll" {
			continue
		}
		// Skip those in allowlist (handled above)
		skip := false
		for _, a := range allowlist {
			if a == name {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		dllPath := filepath.Join(dir, name)
		if h, err := loadOneDll(dllPath); err == nil {
			addLoadedHandle(h)
			loadedCount++
			slog.Debug("preloadSiblingDlls: loaded additional DLL", "name", name, "handle", fmt.Sprintf("0x%x", h))
		}
	}
	slog.Debug("preloadSiblingDlls: completed", "additionalDllsLoaded", loadedCount, "totalLoadedHandles", len(loadedDllHandles))
}

// loadOneDll loads a single DLL by absolute path using LoadLibraryExW with safe flags
func loadOneDll(path string) (uintptr, error) {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		slog.Debug("loadOneDll: failed to convert path", "path", path, "error", err)
		return 0, err
	}
	if procLoadLibraryExW.Find() == nil {
		if ret, _, callErr := procLoadLibraryExW.Call(
			uintptr(unsafe.Pointer(p)),
			0,
			uintptr(loadLibrarySearchDllLoadDir|loadLibrarySearchDefaultDirs|loadLibrarySearchSystem32|loadLibrarySearchUserDirs),
		); ret != 0 {
			slog.Debug("loadOneDll: successfully loaded with LoadLibraryExW", "path", path, "handle", fmt.Sprintf("0x%x", ret))
			return ret, nil
		} else {
			slog.Debug("loadOneDll: LoadLibraryExW failed", "path", path, "error", callErr)
		}
	}
	if ret, _, callErr := procLoadLibraryW.Call(uintptr(unsafe.Pointer(p))); ret != 0 {
		slog.Debug("loadOneDll: successfully loaded with LoadLibraryW (fallback)", "path", path, "handle", fmt.Sprintf("0x%x", ret))
		return ret, nil
	} else {
		slog.Debug("loadOneDll: LoadLibraryW failed", "path", path, "error", callErr)
	}
	return 0, fmt.Errorf("failed to preload dll: %s", path)
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
		slog.Warn("failed to register function", "name", fname, "error", err, "handle", fmt.Sprintf("0x%x", handle))
		return
	}

	if fptr == nil {
		slog.Warn("registerLibFunc nil pointer", "name", fname)
		return
	}

	t := reflect.TypeOf(fptr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Func {
		elemKind := "<nil>"
		elemType := "<nil>"
		if t.Kind() == reflect.Ptr {
			elemKind = t.Elem().Kind().String()
			elemType = t.Elem().String()
		}
		slog.Warn("unexpected pointer type for function",
			"name", fname,
			"type", t.String(),
			"elem_kind", elemKind,
			"elem_type", elemType,
			"addr", fmt.Sprintf("0x%x", procAddr),
		)
		return
	}

	// Safe registration wrapper to avoid panics from purego when symbol signatures are unsupported
	if err := safeRegisterLibFunc(fptr, handle, fname); err != nil {
		slog.Warn("failed binding function",
			"name", fname,
			"error", err,
			"addr", fmt.Sprintf("0x%x", procAddr),
		)
		return
	}
	slog.Debug("registered function",
		"name", fname,
		"addr", fmt.Sprintf("0x%x", procAddr),
	)
}

// findSymbolHandle finds which DLL handle contains the given symbol
// It searches the main handle first, then all loaded sibling DLLs
// Returns the handle that contains the symbol, or an error if not found
func findSymbolHandle(handle uintptr, name string) (uintptr, error) {
	if handle == 0 {
		return 0, fmt.Errorf("invalid library handle (0) when looking up %s", name)
	}

	// Try the main handle first
	namePtr, err := syscall.BytePtrFromString(name)
	if err != nil {
		return 0, fmt.Errorf("failed to convert name to byte pointer: %w", err)
	}

	ret, _, lastErr := procGetProcAddress.Call(handle, uintptr(unsafe.Pointer(namePtr)))
	if ret != 0 {
		slog.Debug("symbol found in main library", "symbol", name, "handle", fmt.Sprintf("0x%x", handle))
		return handle, nil
	}

	// If not found in main handle, search all loaded DLL handles
	if len(loadedDllHandles) > 0 {
		slog.Debug("searching for symbol in sibling DLL handles", "symbol", name, "totalHandles", len(loadedDllHandles))
		for i, h := range loadedDllHandles {
			if h == 0 || h == handle {
				continue
			}
			namePtr, err := syscall.BytePtrFromString(name)
			if err != nil {
				lastErr = err
				continue
			}
			addr, _, _ := procGetProcAddress.Call(h, uintptr(unsafe.Pointer(namePtr)))
			if addr != 0 {
				slog.Debug("symbol found in sibling DLL", "symbol", name, "handle", fmt.Sprintf("0x%x", h), "handleIndex", i)
				return h, nil
			}
		}
	}

	// Not found anywhere
	errno := syscall.Errno(0)
	if err, ok := lastErr.(syscall.Errno); ok {
		errno = err
	}
	return 0, fmt.Errorf("symbol %s not found in library handle 0x%x or %d loaded DLLs: %w (GetLastError: %d)",
		name, handle, len(loadedDllHandles), lastErr, errno)
}

// tryRegisterLibFunc attempts to register a library function, returning an error if it fails
// This is useful for optional functions that may not exist in all library builds
func tryRegisterLibFunc(fptr interface{}, handle uintptr, fname string) error {
	// Find which handle actually has this symbol (might be in a sibling DLL)
	actualHandle, err := findSymbolHandle(handle, fname)
	if err != nil {
		return err
	}
	if fptr == nil {
		return fmt.Errorf("tryRegisterLibFunc: nil function pointer for %s", fname)
	}
	t := reflect.TypeOf(fptr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Func {
		return fmt.Errorf("tryRegisterLibFunc: expected pointer to func for %s, got %s", fname, t.String())
	}
	// Register using the actual handle where the symbol was found
	if err := safeRegisterLibFunc(fptr, actualHandle, fname); err != nil {
		return err
	}
	return nil
}

// safeRegisterLibFunc wraps purego.RegisterLibFunc with panic protection.
// Some optional symbols or exotic signatures on Windows may trigger a panic inside purego.
// We convert those to errors so callers can degrade gracefully.
func safeRegisterLibFunc(fptr interface{}, handle uintptr, fname string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during RegisterLibFunc for %s: %v", fname, r)
		}
	}()
	purego.RegisterLibFunc(fptr, handle, fname)
	return nil
}

// getProcAddressPlatform gets the address of a symbol in a loaded library
func getProcAddressPlatform(handle uintptr, name string) (uintptr, error) {
	if handle == 0 {
		return 0, fmt.Errorf("invalid library handle (0) when looking up %s", name)
	}

	// Log all available DLL handles for debugging
	if len(loadedDllHandles) > 0 {
		slog.Debug("getProcAddressPlatform: searching for symbol", "symbol", name, "mainHandle", fmt.Sprintf("0x%x", handle), "totalSiblingHandles", len(loadedDllHandles))
	}

	// Try the name on the provided handle first
	namePtr, err := syscall.BytePtrFromString(name)
	if err != nil {
		return 0, fmt.Errorf("failed to convert name to byte pointer: %w", err)
	}

	ret, _, lastErr := procGetProcAddress.Call(handle, uintptr(unsafe.Pointer(namePtr)))
	if ret != 0 {
		slog.Debug("symbol resolved from main library", "symbol", name, "handle", fmt.Sprintf("0x%x", handle))
		return ret, nil
	}
	slog.Debug("symbol not found in main handle", "symbol", name, "handle", fmt.Sprintf("0x%x", handle), "error", fmt.Sprintf("%v", lastErr))

	// If not found in main handle, try ALL loaded DLL handles (not just non-main ones)
	// This is important because some symbols may only be in specific dlls like ggml-base.dll
	if len(loadedDllHandles) > 0 {
		slog.Debug("searching in all loaded DLL handles", "symbol", name, "totalHandles", len(loadedDllHandles))
		for i, h := range loadedDllHandles {
			if h == 0 {
				continue
			}
			// Skip the main handle since we already tried it, but be exhaustive with siblings
			if h == handle {
				continue
			}
			namePtr, err := syscall.BytePtrFromString(name)
			if err != nil {
				lastErr = err
				continue
			}
			addr, _, _ := procGetProcAddress.Call(h, uintptr(unsafe.Pointer(namePtr)))
			if addr != 0 {
				slog.Debug("symbol resolved from sibling DLL", "symbol", name, "handle", fmt.Sprintf("0x%x", h), "handleIndex", i)
				return addr, nil
			}
		}
		slog.Debug("symbol not found in any loaded DLL handle", "symbol", name, "totalHandlesSearched", len(loadedDllHandles))
	} else {
		slog.Debug("no sibling DLL handles available for symbol lookup", "symbol", name)
	}

	// Not found anywhere; return the original error context
	errno := syscall.Errno(0)
	if err, ok := lastErr.(syscall.Errno); ok {
		errno = err
	}
	return 0, fmt.Errorf("GetProcAddress failed for %s in library handle 0x%x and %d loaded DLLs: %w (GetLastError: %d). "+
		"The symbol may not be exported by this build.",
		name, handle, len(loadedDllHandles), lastErr, errno)
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
