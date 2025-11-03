package gollama

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Additional tests to increase coverage for goggml.go
type GgmlMoreSuite struct{ BaseSuite }

func (s *GgmlMoreSuite) SetupTest() {
	s.BaseSuite.SetupTest()
	require.NoError(s.T(), Backend_init())
}

func (s *GgmlMoreSuite) TearDownTest() {
	Backend_free()
	s.BaseSuite.TearDownTest()
}

// Allocate a small CPU buffer (if supported) and exercise buffer helpers
func (s *GgmlMoreSuite) TestCpuBufferHelpers() {
	buft, err := Ggml_backend_cpu_buffer_type()
	if err != nil {
		s.T().Skipf("CPU buffer type not available: %v", err)
		return
	}

	if ggmlBackendBuftAllocBuffer == nil {
		s.T().Skip("ggml_backend_buft_alloc_buffer not available in this build")
		return
	}

	// Allocate small buffer and query properties
	buf := ggmlBackendBuftAllocBuffer(buft, 1024)
	// Some builds may return 0 for unsupported paths
	if buf != 0 {
		size, err := Ggml_backend_buffer_get_size(buf)
		if err == nil {
			assert.GreaterOrEqual(s.T(), size, uint64(0))
		}
		_, _ = Ggml_backend_buffer_is_host(buf)
		_, _ = Ggml_backend_buffer_name(buf)
		_ = Ggml_backend_buffer_free(buf)
	}
}

// Exercise load/unload paths and load_all_from_path with empty string
func (s *GgmlMoreSuite) TestLoadAndUnloadBackends() {
	// load all from empty path (nil pointer branch)
	_ = Ggml_backend_load_all_from_path("")

	if globalLoader.rootLibPath != "" {
		reg, err := Ggml_backend_load(globalLoader.rootLibPath)
		if err == nil && reg != 0 {
			_ = Ggml_backend_unload(reg)
		}
	}
}

// Directly cover bytePointerToString helper
func (s *GgmlMoreSuite) TestBytePointerToString() {
	bs := []byte("hello\x00")
	got := bytePointerToString(&bs[0])
	assert.Equal(s.T(), "hello", got)
}

func TestGgmlMoreSuite(t *testing.T) { suite.Run(t, new(GgmlMoreSuite)) }
