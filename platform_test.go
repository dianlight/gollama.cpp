package gollama

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PlatformSuite struct{ BaseSuite }

func (s *PlatformSuite) TestPlatformSupportDetection() {
	supported := isPlatformSupported()

	if runtime.GOOS == "windows" {
		assert.True(s.T(), supported, "Windows platform should be supported with FFI")
		err := getPlatformError()
		assert.NoError(s.T(), err, "getPlatformError should return nil for Windows with FFI support")
		s.T().Log("Windows platform correctly reports as supported with FFI")
	} else {
		assert.True(s.T(), supported, "Unix-like platforms should be supported")
		err := getPlatformError()
		assert.NoError(s.T(), err, "getPlatformError should return nil for supported platforms")
		s.T().Log("Unix-like platform correctly reports as supported")
	}
}

func (s *PlatformSuite) TestPlatformLibraryFunctions() {
	if runtime.GOOS == "windows" {
		_, err := loadLibraryPlatform("nonexistent.dll")
		assert.Error(s.T(), err, "loadLibraryPlatform should fail for non-existent library")
		s.T().Logf("Windows loadLibraryPlatform correctly failed: %v", err)

		err = closeLibraryPlatform(0)
		assert.Error(s.T(), err, "closeLibraryPlatform should fail for invalid handle")
		s.T().Logf("Windows closeLibraryPlatform correctly failed: %v", err)

		var dummy uintptr
		registerLibFunc(&dummy, 0, "test_function")
		s.T().Log("Windows registerLibFunc completed without panic")

		_, err = getProcAddressPlatform(0, "test_function")
		assert.Error(s.T(), err, "getProcAddressPlatform should fail for invalid handle")
		s.T().Logf("Windows getProcAddressPlatform correctly failed: %v", err)
	} else {
		s.T().Log("Unix-like platform functions are available through purego")
	}
}

func TestPlatformSuite(t *testing.T) { suite.Run(t, new(PlatformSuite)) }
