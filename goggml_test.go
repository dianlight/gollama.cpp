package gollama

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GgmlSuite struct{ BaseSuite }

// Tests the GGML type size function
func (s *GgmlSuite) TestGgmlTypeSize() {
	if err := Backend_init(); err != nil {
		s.T().Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	tests := []struct {
		name     string
		typ      GgmlType
		wantSize uint64
	}{
		{"F32", GGML_TYPE_F32, 4},
		{"F16", GGML_TYPE_F16, 2},
		{"I8", GGML_TYPE_I8, 1},
		{"I16", GGML_TYPE_I16, 2},
		{"I32", GGML_TYPE_I32, 4},
		{"I64", GGML_TYPE_I64, 8},
		{"F64", GGML_TYPE_F64, 8},
	}

	for _, tt := range tests {
		size, err := Ggml_type_size(tt.typ)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), tt.wantSize, size)
	}
}

// Tests whether types are correctly identified as quantized
func (s *GgmlSuite) TestGgmlTypeIsQuantized() {
	if err := Backend_init(); err != nil {
		s.T().Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	tests := []struct {
		name          string
		typ           GgmlType
		wantQuantized bool
	}{
		{"F32", GGML_TYPE_F32, false},
		{"F16", GGML_TYPE_F16, false},
		{"Q4_0", GGML_TYPE_Q4_0, true},
		{"Q4_1", GGML_TYPE_Q4_1, true},
		{"Q5_0", GGML_TYPE_Q5_0, true},
		{"Q8_0", GGML_TYPE_Q8_0, true},
		{"Q2_K", GGML_TYPE_Q2_K, true},
		{"I32", GGML_TYPE_I32, false},
	}

	for _, tt := range tests {
		isQuantized, err := Ggml_type_is_quantized(tt.typ)
		if err != nil {
			s.T().Skipf("Ggml_type_is_quantized() not available: %v", err)
			continue
		}
		assert.Equal(s.T(), tt.wantQuantized, isQuantized)
	}
}

// Tests the String method for GgmlType
func (s *GgmlSuite) TestGgmlTypeString() {
	tests := []struct {
		typ  GgmlType
		want string
	}{
		{GGML_TYPE_F32, "f32"},
		{GGML_TYPE_F16, "f16"},
		{GGML_TYPE_Q4_0, "q4_0"},
		{GGML_TYPE_Q4_1, "q4_1"},
		{GGML_TYPE_Q8_0, "q8_0"},
		{GGML_TYPE_Q2_K, "q2_K"},
		{GGML_TYPE_I32, "i32"},
		{GGML_TYPE_BF16, "bf16"},
		{GgmlType(999), "unknown"},
	}

	for _, tt := range tests {
		assert.Equal(s.T(), tt.want, tt.typ.String())
	}
}

// Tests the backend device count function
func (s *GgmlSuite) TestGgmlBackendDevCount() {
	if err := Backend_init(); err != nil {
		s.T().Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	if err := Ggml_backend_load_all(); err != nil {
		s.T().Logf("ggml_backend_load not available or failed: %v", err)
		return
	}

	count, err := Ggml_backend_dev_count()
	if err != nil {
		s.T().Skipf("Ggml_backend_dev_count() not available: %v", err)
		return
	}
	assert.NotZero(s.T(), count, "GGML no backend device functions available in this build")
	s.T().Logf("Found %d backend device(s)", count)
}

// Tests getting backend device information
func (s *GgmlSuite) TestGgmlBackendDevInfo() {
	if err := Backend_init(); err != nil {
		s.T().Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	if err := Ggml_backend_load_all(); err != nil {
		s.T().Logf("ggml_backend_load not available or failed: %v", err)
		return
	}

	count, err := Ggml_backend_dev_count()
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), count, "No backend devices available")

	device, err := Ggml_backend_dev_get(0)
	assert.NoError(s.T(), err)

	name, err := Ggml_backend_dev_name(device)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), name)
	s.T().Logf("Device 0: %s", name)

	desc, err := Ggml_backend_dev_description(device)
	assert.NoError(s.T(), err)
	if desc != "" {
		s.T().Logf("Description: %s", desc)
	}

	if free, total, err := Ggml_backend_dev_memory(device); err == nil {
		s.T().Logf("Memory: %d bytes free / %d bytes total", free, total)
	}
}

// Tests getting the CPU buffer type
func (s *GgmlSuite) TestGgmlBackendCpuBufferType() {
	if err := Backend_init(); err != nil {
		s.T().Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	bufType, err := Ggml_backend_cpu_buffer_type()
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), bufType, "Ggml_backend_cpu_buffer_type() returned null buffer type")
}

// Tests getting type names via GGML
func (s *GgmlSuite) TestGgmlTypeName() {
	if err := Backend_init(); err != nil {
		s.T().Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	tests := []struct {
		typ      GgmlType
		wantName string
	}{
		{GGML_TYPE_F32, "f32"},
		{GGML_TYPE_F16, "f16"},
		{GGML_TYPE_Q4_0, "q4_0"},
		{GGML_TYPE_Q8_0, "q8_0"},
		{GGML_TYPE_I32, "i32"},
	}

	for _, tt := range tests {
		name, err := Ggml_type_name(tt.typ)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), tt.wantName, name)
	}
}

// Tests backend loading by name
func (s *GgmlSuite) TestGgmlBackendLoad() {
	if err := Backend_init(); err != nil {
		s.T().Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	backend, err := Ggml_backend_load("CPU", "")
	if err != nil {
		s.T().Logf("ggml_backend_load not available or failed: %v", err)
		return
	}
	if backend != 0 {
		if name, err := Ggml_backend_name(backend); err == nil {
			s.T().Logf("Loaded backend: %s", name)
		}
	}
}

// Tests loading all available backends
func (s *GgmlSuite) TestGgmlBackendLoadAll() {
	if err := Backend_init(); err != nil {
		s.T().Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	if err := Ggml_backend_load_all(); err != nil {
		s.T().Logf("ggml_backend_load_all not available: %v", err)
		return
	}
	if count, err := Ggml_backend_dev_count(); err == nil {
		s.T().Logf("Backend device count after load_all: %d", count)
	}
}

// Tests loading all backends from a specific path
func (s *GgmlSuite) TestGgmlBackendLoadAllFromPath() {
	if err := Backend_init(); err != nil {
		s.T().Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	if err := Ggml_backend_load_all_from_path("."); err != nil {
		s.T().Logf("ggml_backend_load_all_from_path not available: %v", err)
		return
	}
	if count, err := Ggml_backend_dev_count(); err == nil {
		s.T().Logf("Backend device count after load_all_from_path: %d", count)
	}
}

func TestGgmlSuite(t *testing.T) { suite.Run(t, new(GgmlSuite)) }

// BenchmarkGgmlTypeSize benchmarks the type size function
func BenchmarkGgmlTypeSize(b *testing.B) {
	// Initialize backend
	if err := Backend_init(); err != nil {
		b.Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Ggml_type_size(GGML_TYPE_F32)
	}
}

// BenchmarkGgmlTypeIsQuantized benchmarks the type quantization check
func BenchmarkGgmlTypeIsQuantized(b *testing.B) {
	// Initialize backend
	if err := Backend_init(); err != nil {
		b.Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Ggml_type_is_quantized(GGML_TYPE_Q4_0)
	}
}
