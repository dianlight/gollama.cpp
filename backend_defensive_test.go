package gollama

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// BackendDefensiveSuite tests that backend wrappers are defensive when library or symbols are missing
type BackendDefensiveSuite struct {
	BaseSuite
}

// TestBackendInitWithoutLibrary tests that Backend_init returns an error when library is not loaded
func (s *BackendDefensiveSuite) TestBackendInitWithoutLibrary() {
	// Ensure library is not loaded
	Cleanup()

	// Save the current state
	savedInit := llamaBackendInit
	savedLoaded := isLoaded
	savedHandle := libHandle

	// Simulate library loaded but symbol missing
	isLoaded = true
	libHandle = 1 // Non-zero to indicate "loaded"
	llamaBackendInit = nil

	// Restore after test
	defer func() {
		llamaBackendInit = savedInit
		isLoaded = savedLoaded
		libHandle = savedHandle
	}()

	// Backend_init should return an error when the symbol is missing
	err := Backend_init()
	s.Require().Error(err, "Backend_init should return error when llamaBackendInit is nil")
	s.Contains(err.Error(), "not available", "Error should mention function not available")
}

// TestBackendFreeWithoutLibrary tests that Backend_free is a safe no-op when library is not loaded
func (s *BackendDefensiveSuite) TestBackendFreeWithoutLibrary() {
	// Ensure library is not loaded
	Cleanup()

	// Save the current state
	savedFree := llamaBackendFree
	savedLoaded := isLoaded

	// Simulate missing symbol by setting to nil
	llamaBackendFree = nil
	isLoaded = false

	// Restore after test
	defer func() {
		llamaBackendFree = savedFree
		isLoaded = savedLoaded
	}()

	// Backend_free should not panic even when the symbol is missing
	s.Require().NotPanics(func() {
		Backend_free()
	}, "Backend_free should not panic when llamaBackendFree is nil")
}

// TestBackendFreeWithNilFunction tests that Backend_free is safe when function pointer is nil but library is loaded
func (s *BackendDefensiveSuite) TestBackendFreeWithNilFunction() {
	// Save the current state
	savedFree := llamaBackendFree
	savedLoaded := isLoaded

	// Simulate loaded library but missing symbol
	llamaBackendFree = nil
	isLoaded = true

	// Restore after test
	defer func() {
		llamaBackendFree = savedFree
		isLoaded = savedLoaded
	}()

	// Backend_free should not panic even when the symbol is missing
	s.Require().NotPanics(func() {
		Backend_free()
	}, "Backend_free should not panic when llamaBackendFree is nil even if isLoaded is true")
}

// TestGgmlBackendFreeDefensive tests that Ggml_backend_free is a safe no-op when function is missing
func (s *BackendDefensiveSuite) TestGgmlBackendFreeDefensive() {
	// Ensure library is not loaded
	Cleanup()

	// Save the current function pointer
	savedFree := ggmlBackendFree

	// Simulate missing symbol by setting to nil
	ggmlBackendFree = nil

	// Restore after test
	defer func() {
		ggmlBackendFree = savedFree
	}()

	// Ggml_backend_free should not return an error when the symbol is missing (safe no-op)
	err := Ggml_backend_free(0)
	s.Require().EqualError(err, "ggml_backend_free function not available", "Ggml_backend_free should be a safe no-op when ggmlBackendFree is nil")
}

// TestGgmlBackendBufferFreeDefensive tests that Ggml_backend_buffer_free is a safe no-op when function is missing
func (s *BackendDefensiveSuite) TestGgmlBackendBufferFreeDefensive() {
	// Ensure library is not loaded
	Cleanup()

	// Save the current function pointer
	savedFree := ggmlBackendBufferFree

	// Simulate missing symbol by setting to nil
	ggmlBackendBufferFree = nil

	// Restore after test
	defer func() {
		ggmlBackendBufferFree = savedFree
	}()

	// Ggml_backend_buffer_free should not return an error when the symbol is missing (safe no-op)
	err := Ggml_backend_buffer_free(0)
	s.Require().EqualError(err, "ggml_backend_buffer_free function not available", "Ggml_backend_buffer_free should be a safe no-op when ggmlBackendBufferFree is nil")
}

// TestGgmlBackendUnloadDefensive tests that Ggml_backend_unload is a safe no-op when function is missing
func (s *BackendDefensiveSuite) TestGgmlBackendUnloadDefensive() {
	// Ensure library is not loaded
	Cleanup()

	// Save the current function pointer
	savedUnload := ggmlBackendUnload

	// Simulate missing symbol by setting to nil
	ggmlBackendUnload = nil

	// Restore after test
	defer func() {
		ggmlBackendUnload = savedUnload
	}()

	// Ggml_backend_unload should not return an error when the symbol is missing (safe no-op)
	err := Ggml_backend_unload(0)
	s.Require().EqualError(err, "ggml_backend_unload function not available", "Ggml_backend_unload should be a safe no-op when ggmlBackendUnload is nil")
}

// TestGgmlBackendInitFunctionsReturnErrors tests that init functions return proper errors when symbols are missing
func (s *BackendDefensiveSuite) TestGgmlBackendInitFunctionsReturnErrors() {
	// Ensure library is not loaded
	Cleanup()

	// Save the current function pointers
	savedInitBest := ggmlBackendInitBest
	savedInitByName := ggmlBackendInitByName
	savedInitByType := ggmlBackendInitByType
	// Simulate missing symbols by setting to nil
	ggmlBackendInitBest = nil
	ggmlBackendInitByName = nil
	ggmlBackendInitByType = nil

	// Restore after test
	defer func() {
		ggmlBackendInitBest = savedInitBest
		ggmlBackendInitByName = savedInitByName
		ggmlBackendInitByType = savedInitByType
	}()

	// All init functions should return errors when symbols are missing
	_, err := Ggml_backend_init_best()
	s.Require().Error(err, "Ggml_backend_init_best should return error when ggmlBackendInitBest is nil")
	s.Contains(err.Error(), "not available", "Error should mention function not available")

	_, err = Ggml_backend_init_by_name("cpu", "")
	s.Require().Error(err, "Ggml_backend_init_by_name should return error when ggmlBackendInitByName is nil")
	s.Contains(err.Error(), "not available", "Error should mention function not available")

	_, err = Ggml_backend_init_by_type(GGML_BACKEND_DEVICE_TYPE_CPU, "")
	s.Require().Error(err, "Ggml_backend_init_by_type should return error when ggmlBackendInitByType is nil")
	s.Contains(err.Error(), "not available", "Error should mention function not available")
}

// TestCompleteLibraryLoadFailureScenario tests defensive behavior when symbols are missing but library loads
func (s *BackendDefensiveSuite) TestCompleteLibraryLoadFailureScenario() {
	// This test simulates the scenario where the library loads successfully
	// but certain backend symbols are not available (e.g., older version of llama.cpp)

	// Ensure library is loaded
	_ = ensureLoaded()

	// Save current state
	savedState := struct {
		llamaBackendInit      func()
		llamaBackendFree      func()
		ggmlBackendFree       func(GgmlBackend)
		ggmlBackendBufferFree func(GgmlBackendBuffer)
		ggmlBackendUnload     func(GgmlBackendReg)
	}{
		llamaBackendInit:      llamaBackendInit,
		llamaBackendFree:      llamaBackendFree,
		ggmlBackendFree:       ggmlBackendFree,
		ggmlBackendBufferFree: ggmlBackendBufferFree,
		ggmlBackendUnload:     ggmlBackendUnload,
	}

	// Simulate missing symbols
	llamaBackendInit = nil
	llamaBackendFree = nil
	ggmlBackendFree = nil
	ggmlBackendBufferFree = nil
	ggmlBackendUnload = nil

	// Restore state
	defer func() {
		llamaBackendInit = savedState.llamaBackendInit
		llamaBackendFree = savedState.llamaBackendFree
		ggmlBackendFree = savedState.ggmlBackendFree
		ggmlBackendBufferFree = savedState.ggmlBackendBufferFree
		ggmlBackendUnload = savedState.ggmlBackendUnload
	}()

	// Backend_init should return error when symbol is missing
	err := Backend_init()
	s.Require().Error(err, "Backend_init should return error when llamaBackendInit is nil")

	// Backend_free should not panic even when symbol is missing
	s.Require().NotPanics(func() {
		Backend_free()
	}, "Backend_free should not panic when llamaBackendFree is nil")

	// Free functions should be safe no-ops when symbols are missing
	err = Ggml_backend_free(0)
	s.Require().Error(err, "Ggml_backend_free should be safe no-op when ggmlBackendFree is nil")

	err = Ggml_backend_buffer_free(0)
	s.Require().Error(err, "Ggml_backend_buffer_free should be safe no-op when ggmlBackendBufferFree is nil")

	err = Ggml_backend_unload(0)
	s.Require().Error(err, "Ggml_backend_unload should be safe no-op when ggmlBackendUnload is nil")
}

// TestBackendFunctionsWithPartiallyLoadedLibrary tests behavior when library is loaded but some symbols are missing
func (s *BackendDefensiveSuite) TestBackendFunctionsWithPartiallyLoadedLibrary() {
	// This simulates a scenario where the library loads but some symbols are not available
	// (e.g., older version of llama.cpp without certain functions)

	// Save current state
	savedInit := llamaBackendInit
	savedFree := llamaBackendFree
	savedLoaded := isLoaded
	savedHandle := libHandle

	// Simulate partially loaded library
	isLoaded = true
	libHandle = 1 // Non-zero to indicate "loaded"
	llamaBackendInit = nil
	llamaBackendFree = nil

	// Restore after test
	defer func() {
		llamaBackendInit = savedInit
		llamaBackendFree = savedFree
		isLoaded = savedLoaded
		libHandle = savedHandle
	}()

	// Backend_init should return error when symbol is missing
	err := Backend_init()
	s.Require().Error(err, "Backend_init should return error when llamaBackendInit symbol is missing")

	// Backend_free should be safe no-op
	s.Require().NotPanics(func() {
		Backend_free()
	}, "Backend_free should not panic when llamaBackendFree symbol is missing")
}

// TestNoPanicOnNilFunctionPointers is a comprehensive test ensuring no nil pointer panics
func (s *BackendDefensiveSuite) TestNoPanicOnNilFunctionPointers() {
	// Ensure library is not loaded
	Cleanup()

	// Save all backend-related function pointers
	savedBackendFuncs := struct {
		init            func()
		free            func()
		ggmlFree        func(GgmlBackend)
		ggmlBufferFree  func(GgmlBackendBuffer)
		ggmlUnload      func(GgmlBackendReg)
		ggmlInitBest    func() GgmlBackend
		ggmlInitByName  func(*byte, *byte) GgmlBackend
		ggmlInitByType  func(int32, *byte) GgmlBackend
		ggmlLoad        func(*byte) GgmlBackendReg
		ggmlLoadAll     func()
		ggmlLoadAllPath func(*byte)
	}{
		init:            llamaBackendInit,
		free:            llamaBackendFree,
		ggmlFree:        ggmlBackendFree,
		ggmlBufferFree:  ggmlBackendBufferFree,
		ggmlUnload:      ggmlBackendUnload,
		ggmlInitBest:    ggmlBackendInitBest,
		ggmlInitByName:  ggmlBackendInitByName,
		ggmlInitByType:  ggmlBackendInitByType,
		ggmlLoad:        ggmlBackendLoad,
		ggmlLoadAll:     ggmlBackendLoadAll,
		ggmlLoadAllPath: ggmlBackendLoadAllFromPath,
	}

	// Set all to nil
	llamaBackendInit = nil
	llamaBackendFree = nil
	ggmlBackendFree = nil
	ggmlBackendBufferFree = nil
	ggmlBackendUnload = nil
	ggmlBackendInitBest = nil
	ggmlBackendInitByName = nil
	ggmlBackendInitByType = nil
	ggmlBackendLoad = nil
	ggmlBackendLoadAll = nil
	ggmlBackendLoadAllFromPath = nil

	// Restore after test
	defer func() {
		llamaBackendInit = savedBackendFuncs.init
		llamaBackendFree = savedBackendFuncs.free
		ggmlBackendFree = savedBackendFuncs.ggmlFree
		ggmlBackendBufferFree = savedBackendFuncs.ggmlBufferFree
		ggmlBackendUnload = savedBackendFuncs.ggmlUnload
		ggmlBackendInitBest = savedBackendFuncs.ggmlInitBest
		ggmlBackendInitByName = savedBackendFuncs.ggmlInitByName
		ggmlBackendInitByType = savedBackendFuncs.ggmlInitByType
		ggmlBackendLoad = savedBackendFuncs.ggmlLoad
		ggmlBackendLoadAll = savedBackendFuncs.ggmlLoadAll
		ggmlBackendLoadAllFromPath = savedBackendFuncs.ggmlLoadAllPath
	}()

	// Test that none of these cause panics
	s.Require().NotPanics(func() {
		Backend_free()
	}, "Backend_free should not panic with nil llamaBackendFree")

	s.Require().NotPanics(func() {
		_ = Ggml_backend_free(0)
	}, "Ggml_backend_free should not panic with nil ggmlBackendFree")

	s.Require().NotPanics(func() {
		_ = Ggml_backend_buffer_free(0)
	}, "Ggml_backend_buffer_free should not panic with nil ggmlBackendBufferFree")

	s.Require().NotPanics(func() {
		_ = Ggml_backend_unload(0)
	}, "Ggml_backend_unload should not panic with nil ggmlBackendUnload")

	// Init functions should return errors, not panic
	s.Require().NotPanics(func() {
		err := Backend_init()
		s.Error(err)
	}, "Backend_init should return error, not panic with nil llamaBackendInit")

	s.Require().NotPanics(func() {
		_, err := Ggml_backend_init_best()
		s.Error(err)
	}, "Ggml_backend_init_best should return error, not panic with nil ggmlBackendInitBest")

	s.Require().NotPanics(func() {
		_, err := Ggml_backend_init_by_name("cpu", "")
		s.Error(err)
	}, "Ggml_backend_init_by_name should return error, not panic with nil ggmlBackendInitByName")

	s.Require().NotPanics(func() {
		_, err := Ggml_backend_init_by_type(GGML_BACKEND_DEVICE_TYPE_CPU, "")
		s.Error(err)
	}, "Ggml_backend_init_by_type should return error, not panic with nil ggmlBackendInitByType")

	s.Require().NotPanics(func() {
		_, err := Ggml_backend_load("/some/path")
		s.Error(err)
	}, "Ggml_backend_load should return error, not panic with nil ggmlBackendLoad")

	s.Require().NotPanics(func() {
		err := Ggml_backend_load_all()
		s.Error(err)
	}, "Ggml_backend_load_all should return error, not panic with nil ggmlBackendLoadAll")

	s.Require().NotPanics(func() {
		err := Ggml_backend_load_all_from_path("/some/path")
		s.Error(err)
	}, "Ggml_backend_load_all_from_path should return error, not panic with nil ggmlBackendLoadAllFromPath")
}

func TestBackendDefensiveSuite(t *testing.T) {
	suite.Run(t, new(BackendDefensiveSuite))
}
