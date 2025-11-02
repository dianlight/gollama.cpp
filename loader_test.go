package gollama

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ BaseSuite }

func (s *LoaderSuite) TestGetLibraryName() {
	loader := &LibraryLoader{}
	result, err := loader.getLibraryName()
	if err != nil {
		if runtime.GOOS != "darwin" && runtime.GOOS != "linux" && runtime.GOOS != "windows" {
			s.T().Logf("Current platform %s is unsupported, which is expected", runtime.GOOS)
			return
		}
		s.T().Fatalf("Unexpected error for supported OS %s: %v", runtime.GOOS, err)
	}
	switch runtime.GOOS {
	case "darwin":
		assert.Equal(s.T(), "libllama.dylib", result)
	case "linux":
		assert.Equal(s.T(), "libllama.so", result)
	case "windows":
		assert.Equal(s.T(), "llama.dll", result)
	default:
		s.T().Logf("Platform %s returned %s", runtime.GOOS, result)
	}
}

func (s *LoaderSuite) TestLoadSharedLibrary_WindowsBehavior() {
	if runtime.GOOS != "windows" {
		s.T().Skip("Skipping Windows test on non-Windows platform")
	}
	loader := &LibraryLoader{}
	_, err := loader.loadSharedLibrary("test.dll")
	assert.Error(s.T(), err, "Expected error for non-existent DLL")
}

func (s *LoaderSuite) TestLoadSharedLibrary_Unix() {
	if runtime.GOOS == "windows" {
		s.T().Skip("Skipping Unix test on Windows")
	}
	loader := &LibraryLoader{}
	_, err := loader.loadSharedLibrary("/invalid/path/libtest.so")
	assert.Error(s.T(), err, "Expected error for invalid path")
}

func (s *LoaderSuite) TestExtractEmbeddedLibraries_NonExistent() {
	loader := &LibraryLoader{}
	_, err := loader.extractEmbeddedLibraries()
	if err == nil {
		s.T().Log("extractEmbeddedLibraries succeeded - embedded libraries are present")
		if loader.tempDir != "" {
			_ = os.RemoveAll(loader.tempDir)
			loader.tempDir = ""
		}
	} else {
		s.T().Logf("extractEmbeddedLibraries failed as expected when no libraries present: %v", err)
	}
}

func (s *LoaderSuite) TestExtractEmbeddedLibraries_PlatformSpecificSets() {
	loader := &LibraryLoader{}
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	expectedExtensions := map[string]string{
		"darwin":  ".dylib",
		"linux":   ".so",
		"windows": ".dll",
	}
	if expectedExt, exists := expectedExtensions[goos]; exists {
		mainLibPath, err := loader.extractEmbeddedLibraries()
		if err == nil {
			s.T().Log("extractEmbeddedLibraries succeeded - libraries may be present")
			expectedMainLib, _ := loader.getLibraryName()
			assert.Equal(s.T(), expectedMainLib, filepath.Base(mainLibPath))
			if loader.tempDir != "" {
				files, err := os.ReadDir(loader.tempDir)
				if err == nil {
					s.T().Logf("Extracted %d files to %s", len(files), loader.tempDir)
					for _, file := range files {
						fileName := file.Name()
						assert.Equalf(s.T(), expectedExt, filepath.Ext(fileName), "Unexpected file extension for %s", fileName)
					}
				}
				_ = os.RemoveAll(loader.tempDir)
				loader.tempDir = ""
			}
		} else {
			if err.Error() == fmt.Sprintf("unsupported OS: %s", goos) {
				s.T().Fatalf("Platform should be supported: %s", goos)
			} else {
				s.T().Logf("Expected failure for missing libraries: %v", err)
			}
		}
		s.T().Logf("Platform %s_%s expects files with extension: %s", goos, goarch, expectedExt)
	} else {
		s.T().Logf("Platform %s not in expected patterns", goos)
	}
}

func (s *LoaderSuite) TestGetHandle_InitiallyZero() {
	loader := &LibraryLoader{}
	handle := loader.GetHandle()
	assert.Equal(s.T(), uintptr(0), handle)
}

func (s *LoaderSuite) TestGetHandle_AfterSettingHandle() {
	loader := &LibraryLoader{}
	expectedHandle := uintptr(12345)
	loader.handle = expectedHandle
	loader.loaded = true
	handle := loader.GetHandle()
	assert.Equal(s.T(), expectedHandle, handle)
	loader.handle = 0
	loader.loaded = false
}

func (s *LoaderSuite) TestIsLoaded_InitiallyFalse() {
	loader := &LibraryLoader{}
	assert.False(s.T(), loader.IsLoaded())
}

func (s *LoaderSuite) TestIsLoaded_AfterSettingLoaded() {
	loader := &LibraryLoader{}
	loader.loaded = true
	assert.True(s.T(), loader.IsLoaded())
	loader.loaded = false
}

func (s *LoaderSuite) TestUnloadLibrary_WhenNotLoaded() {
	loader := &LibraryLoader{}
	err := loader.UnloadLibrary()
	assert.NoError(s.T(), err)
}

func (s *LoaderSuite) TestUnloadLibrary_WithTemporaryDirectory() {
	loader := &LibraryLoader{}
	tempDir, err := os.MkdirTemp("", "gollama-test-*")
	if err != nil {
		s.T().Fatalf("Failed to create temp dir: %v", err)
	}
	loader.loaded = true
	loader.handle = uintptr(12345)
	loader.tempDir = tempDir
	loader.llamaLibPath = filepath.Join(tempDir, "test.so")
	err = loader.UnloadLibrary()
	assert.NoError(s.T(), err)
	assert.False(s.T(), loader.loaded)
	assert.Equal(s.T(), uintptr(0), loader.handle)
	assert.Empty(s.T(), loader.tempDir)
	assert.Empty(s.T(), loader.llamaLibPath)
	_, statErr := os.Stat(tempDir)
	assert.True(s.T(), os.IsNotExist(statErr), "Expected temp directory to be removed")
}

func (s *LoaderSuite) TestLoadLibrary_WhenAlreadyLoaded() {
	loader := &LibraryLoader{}
	loader.loaded = true
	defer func() { loader.loaded = false }()
	err := loader.LoadLibrary()
	assert.NoError(s.T(), err)
}

func (s *LoaderSuite) TestLoadLibrary_DownloadsWhenEmbeddedDisabled() {
	loader := &LibraryLoader{}
	err := loader.LoadLibrary()
	if err != nil {
		s.T().Errorf("Expected LoadLibrary to succeed by downloading libraries, but got error: %v", err)
	} else {
		s.T().Log("LoadLibrary succeeded by downloading libraries as expected")
	}
	if loader.loaded {
		_ = loader.UnloadLibrary()
	}
}

func (s *LoaderSuite) TestThreadSafety_GetHandle() {
	loader := &LibraryLoader{}
	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_ = loader.GetHandle()
		}()
	}
	wg.Wait()
}

func (s *LoaderSuite) TestThreadSafety_IsLoaded() {
	loader := &LibraryLoader{}
	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_ = loader.IsLoaded()
		}()
	}
	wg.Wait()
}

func (s *LoaderSuite) TestThreadSafety_ConcurrentLoadLibraryCalls() {
	loader := &LibraryLoader{}
	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_ = loader.LoadLibrary()
		}()
	}
	wg.Wait()
	if loader.loaded {
		_ = loader.UnloadLibrary()
	}
}

func (s *LoaderSuite) TestGlobalFunctions_GetLibHandle() {
	handle := getLibHandle()
	expectedHandle := globalLoader.GetHandle()
	assert.Equal(s.T(), expectedHandle, handle)
}

func (s *LoaderSuite) TestGlobalFunctions_IsLibraryLoaded() {
	loaded := isLibraryLoaded()
	expectedLoaded := globalLoader.IsLoaded()
	assert.Equal(s.T(), expectedLoaded, loaded)
}

func (s *LoaderSuite) TestGlobalFunctions_RegisterFunctionNoLibrary() {
	var testFunc func()
	err := RegisterFunction(&testFunc, "test_function")
	assert.Error(s.T(), err, "Expected error when registering function with no library loaded")
	if err != nil {
		assert.Equal(s.T(), "library not loaded", err.Error())
	}
}

func (s *LoaderSuite) TestGlobalFunctions_Cleanup() { Cleanup() }

func (s *LoaderSuite) TestExtractEmbeddedLibrariesWriteFailure_WriteToReadOnlyDirectory() {
	loader := &LibraryLoader{}
	tempDir, err := os.MkdirTemp("", "gollama-readonly-test-*")
	if err != nil {
		s.T().Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()
	if err = os.Chmod(tempDir, 0444); err != nil {
		s.T().Skipf("Cannot change directory permissions: %v", err)
	}
	defer func() { _ = os.Chmod(tempDir, 0755) }()
	_, err = loader.extractEmbeddedLibraries()
	if err == nil {
		s.T().Log("extractEmbeddedLibraries succeeded unexpectedly")
	} else {
		s.T().Logf("extractEmbeddedLibraries failed as expected: %v", err)
	}
}

// Benchmark tests
func BenchmarkLibraryLoader_GetHandle(b *testing.B) {
	loader := &LibraryLoader{}
	loader.handle = uintptr(12345)
	loader.loaded = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = loader.GetHandle()
	}
}

func BenchmarkLibraryLoader_IsLoaded(b *testing.B) {
	loader := &LibraryLoader{}
	loader.loaded = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = loader.IsLoaded()
	}
}

func BenchmarkLibraryLoader_GetLibraryName(b *testing.B) {
	loader := &LibraryLoader{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = loader.getLibraryName()
	}
}

func BenchmarkGlobalFunctions(b *testing.B) {
	b.Run("getLibHandle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = getLibHandle()
		}
	})

	b.Run("isLibraryLoaded", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = isLibraryLoaded()
		}
	})
}

// Test race conditions
func (s *LoaderSuite) TestRaceConditions_LoadAndUnloadRace() {
	loader := &LibraryLoader{}
	const iterations = 50
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = loader.LoadLibrary()
			time.Sleep(time.Microsecond)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = loader.UnloadLibrary()
			time.Sleep(time.Microsecond)
		}
	}()
	wg.Wait()
	if loader.loaded {
		_ = loader.UnloadLibrary()
	}
}

// Test initialization and state
func (s *LoaderSuite) TestInitialState() {
	loader := &LibraryLoader{}
	assert.Equal(s.T(), uintptr(0), loader.handle)
	assert.False(s.T(), loader.loaded)
	assert.Empty(s.T(), loader.tempDir)
	assert.Empty(s.T(), loader.llamaLibPath)
}

// Test global loader initialization
func (s *LoaderSuite) TestGlobalLoader() {
	assert.NotNil(s.T(), globalLoader, "Expected globalLoader to be initialized")
	handle := getLibHandle()
	assert.Equal(s.T(), uintptr(0), handle, "Expected global handle to be 0 initially")
	loaded := isLibraryLoaded()
	assert.False(s.T(), loaded, "Expected global library to not be loaded initially")
}

func TestLoaderSuite(t *testing.T) { suite.Run(t, new(LoaderSuite)) }
