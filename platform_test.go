package gollama

import (
	"runtime"
	"testing"
)

func TestPlatformSpecific(t *testing.T) {
	t.Run("Platform support detection", func(t *testing.T) {
		supported := isPlatformSupported()

		if runtime.GOOS == "windows" {
			// Windows support is now enabled with FFI
			if !supported {
				t.Error("Windows platform should be supported with FFI")
			}

			err := getPlatformError()
			if err != nil {
				t.Errorf("getPlatformError should return nil for Windows with FFI support, got: %v", err)
			}

			t.Log("Windows platform correctly reports as supported with FFI")
		} else {
			// Unix-like platforms should be supported
			if !supported {
				t.Error("Unix-like platforms should be supported")
			}

			err := getPlatformError()
			if err != nil {
				t.Errorf("getPlatformError should return nil for supported platforms, got: %v", err)
			}

			t.Log("Unix-like platform correctly reports as supported")
		}
	})

	t.Run("Platform library functions", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			// Test that Windows functions don't panic
			_, err := loadLibraryPlatform("nonexistent.dll")
			if err == nil {
				t.Error("loadLibraryPlatform should fail for non-existent library")
			}
			t.Logf("Windows loadLibraryPlatform correctly failed: %v", err)

			// Test closeLibraryPlatform with invalid handle
			err = closeLibraryPlatform(0)
			if err == nil {
				t.Error("closeLibraryPlatform should fail for invalid handle")
			}
			t.Logf("Windows closeLibraryPlatform correctly failed: %v", err)

			// Test registerLibFunc doesn't panic
			var dummy uintptr
			registerLibFunc(&dummy, 0, "test_function")
			t.Log("Windows registerLibFunc completed without panic")
			
			// Test getProcAddressPlatform with invalid handle
			_, err = getProcAddressPlatform(0, "test_function")
			if err == nil {
				t.Error("getProcAddressPlatform should fail for invalid handle")
			}
			t.Logf("Windows getProcAddressPlatform correctly failed: %v", err)
		} else {
			// For Unix-like systems, we can test that the functions exist
			// but we don't want to actually load libraries in unit tests
			t.Log("Unix-like platform functions are available through purego")
		}
	})
}
