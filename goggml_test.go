package gollama

import (
	"testing"
)

// TestGgmlTypeSize tests the GGML type size function
func TestGgmlTypeSize(t *testing.T) {
	// Initialize backend
	if err := Backend_init(); err != nil {
		t.Fatalf("Failed to initialize backend: %v", err)
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
		t.Run(tt.name, func(t *testing.T) {
			size, err := Ggml_type_size(tt.typ)
			if err != nil {
				t.Errorf("Ggml_type_size() error = %v", err)
				return
			}
			if size != tt.wantSize {
				t.Errorf("Ggml_type_size() = %v, want %v", size, tt.wantSize)
			}
		})
	}
}

// TestGgmlTypeIsQuantized tests whether types are correctly identified as quantized
func TestGgmlTypeIsQuantized(t *testing.T) {
	// Initialize backend
	if err := Backend_init(); err != nil {
		t.Fatalf("Failed to initialize backend: %v", err)
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
		t.Run(tt.name, func(t *testing.T) {
			isQuantized, err := Ggml_type_is_quantized(tt.typ)
			if err != nil {
				// Function may not be available in all builds
				t.Skipf("Ggml_type_is_quantized() not available: %v", err)
				return
			}
			if isQuantized != tt.wantQuantized {
				t.Errorf("Ggml_type_is_quantized() = %v, want %v", isQuantized, tt.wantQuantized)
			}
		})
	}
}

// TestGgmlTypeString tests the String method for GgmlType
func TestGgmlTypeString(t *testing.T) {
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
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.typ.String(); got != tt.want {
				t.Errorf("GgmlType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGgmlBackendDevCount tests the backend device count function
func TestGgmlBackendDevCount(t *testing.T) {
	// Initialize backend
	if err := Backend_init(); err != nil {
		t.Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	count, err := Ggml_backend_dev_count()
	if err != nil {
		t.Skipf("Ggml_backend_dev_count() not available: %v", err)
	}

	// Note: Function may return 0 if GGML functions are not exported
	// This is expected in some llama.cpp builds
	if count == 0 {
		t.Skip("GGML backend device functions not available in this build")
	}

	t.Logf("Found %d backend device(s)", count)
}

// TestGgmlBackendDevInfo tests getting backend device information
func TestGgmlBackendDevInfo(t *testing.T) {
	// Initialize backend
	if err := Backend_init(); err != nil {
		t.Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	count, err := Ggml_backend_dev_count()
	if err != nil {
		t.Fatalf("Ggml_backend_dev_count() error = %v", err)
	}

	if count == 0 {
		t.Skip("No backend devices available")
	}

	// Get info for the first device (usually CPU)
	device, err := Ggml_backend_dev_get(0)
	if err != nil {
		t.Fatalf("Ggml_backend_dev_get(0) error = %v", err)
	}

	name, err := Ggml_backend_dev_name(device)
	if err != nil {
		t.Fatalf("Ggml_backend_dev_name() error = %v", err)
	}

	if name == "" {
		t.Error("Ggml_backend_dev_name() returned empty string")
	}

	t.Logf("Device 0: %s", name)

	desc, err := Ggml_backend_dev_description(device)
	if err != nil {
		t.Fatalf("Ggml_backend_dev_description() error = %v", err)
	}

	if desc != "" {
		t.Logf("Description: %s", desc)
	}

	// Try to get memory info (may not be supported on all devices)
	free, total, err := Ggml_backend_dev_memory(device)
	if err == nil {
		t.Logf("Memory: %d bytes free / %d bytes total", free, total)
	}
}

// TestGgmlBackendCpuBufferType tests getting the CPU buffer type
func TestGgmlBackendCpuBufferType(t *testing.T) {
	// Initialize backend
	if err := Backend_init(); err != nil {
		t.Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	bufType, err := Ggml_backend_cpu_buffer_type()
	if err != nil {
		t.Fatalf("Ggml_backend_cpu_buffer_type() error = %v", err)
	}

	if bufType == 0 {
		t.Error("Ggml_backend_cpu_buffer_type() returned null buffer type")
	}
}

// TestGgmlTypeName tests getting type names via GGML
func TestGgmlTypeName(t *testing.T) {
	// Initialize backend
	if err := Backend_init(); err != nil {
		t.Fatalf("Failed to initialize backend: %v", err)
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
		t.Run(tt.wantName, func(t *testing.T) {
			name, err := Ggml_type_name(tt.typ)
			if err != nil {
				t.Errorf("Ggml_type_name() error = %v", err)
				return
			}
			if name != tt.wantName {
				t.Errorf("Ggml_type_name() = %v, want %v", name, tt.wantName)
			}
		})
	}
}

// TestGgmlBackendLoad tests backend loading by name
func TestGgmlBackendLoad(t *testing.T) {
	// Initialize backend
	if err := Backend_init(); err != nil {
		t.Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	// Try loading CPU backend (this may not work in all builds)
	backend, err := Ggml_backend_load("CPU", "")
	if err != nil {
		t.Logf("ggml_backend_load not available or failed: %v", err)
		return
	}

	if backend != 0 {
		// Successfully loaded a backend
		name, err := Ggml_backend_name(backend)
		if err == nil {
			t.Logf("Loaded backend: %s", name)
		}
		// Note: Don't free the backend here as it may be managed by the library
	}
}

// TestGgmlBackendLoadAll tests loading all available backends
func TestGgmlBackendLoadAll(t *testing.T) {
	// Initialize backend
	if err := Backend_init(); err != nil {
		t.Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	// Try loading all backends
	err := Ggml_backend_load_all()
	if err != nil {
		t.Logf("ggml_backend_load_all not available: %v", err)
		return
	}

	// If successful, try to enumerate devices
	count, err := Ggml_backend_dev_count()
	if err == nil {
		t.Logf("Backend device count after load_all: %d", count)
	}
}

// TestGgmlBackendLoadAllFromPath tests loading all backends from a specific path
func TestGgmlBackendLoadAllFromPath(t *testing.T) {
	// Initialize backend
	if err := Backend_init(); err != nil {
		t.Fatalf("Failed to initialize backend: %v", err)
	}
	defer Backend_free()

	// Try loading all backends from current directory
	err := Ggml_backend_load_all_from_path(".")
	if err != nil {
		t.Logf("ggml_backend_load_all_from_path not available: %v", err)
		return
	}

	// If successful, try to enumerate devices
	count, err := Ggml_backend_dev_count()
	if err == nil {
		t.Logf("Backend device count after load_all_from_path: %d", count)
	}
}

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
