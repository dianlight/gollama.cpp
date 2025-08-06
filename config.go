package gollama

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Config holds configuration options for gollama
type Config struct {
	// Library settings
	LibraryPath   string `json:"library_path,omitempty"`
	UseEmbedded   bool   `json:"use_embedded"`
	EnableLogging bool   `json:"enable_logging"`
	LogLevel      int    `json:"log_level"`

	// Performance settings
	NumThreads    int  `json:"num_threads"`
	EnableGPU     bool `json:"enable_gpu"`
	GPULayers     int  `json:"gpu_layers"`
	MetalEnabled  bool `json:"metal_enabled"`
	CUDAEnabled   bool `json:"cuda_enabled"`
	VulkanEnabled bool `json:"vulkan_enabled"`

	// Memory settings
	ContextSize       int  `json:"context_size"`
	BatchSize         int  `json:"batch_size"`
	UbatchSize        int  `json:"ubatch_size"`
	MemoryMapEnabled  bool `json:"memory_map_enabled"`
	MemoryLockEnabled bool `json:"memory_lock_enabled"`

	// Model settings
	ModelPath        string `json:"model_path,omitempty"`
	VocabOnly        bool   `json:"vocab_only"`
	UseQuantization  bool   `json:"use_quantization"`
	QuantizationType string `json:"quantization_type,omitempty"`

	// Backend settings
	BackendType string `json:"backend_type,omitempty"`
	DeviceID    int    `json:"device_id"`

	// Debug settings
	VerboseLogging bool `json:"verbose_logging"`
	DebugMode      bool `json:"debug_mode"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	numCPU := runtime.NumCPU()

	return &Config{
		// Library settings
		UseEmbedded:   true,
		EnableLogging: true,
		LogLevel:      1, // LLAMA_LOG_LEVEL_INFO

		// Performance settings
		NumThreads:    numCPU,
		EnableGPU:     detectGPU(),
		GPULayers:     -1, // Use all layers on GPU if available
		MetalEnabled:  runtime.GOOS == "darwin",
		CUDAEnabled:   runtime.GOOS == "linux" || runtime.GOOS == "windows",
		VulkanEnabled: false,

		// Memory settings
		ContextSize:       2048,
		BatchSize:         512,
		UbatchSize:        512,
		MemoryMapEnabled:  true,
		MemoryLockEnabled: false,

		// Model settings
		VocabOnly:       false,
		UseQuantization: false,

		// Backend settings
		BackendType: "auto",
		DeviceID:    0,

		// Debug settings
		VerboseLogging: false,
		DebugMode:      false,
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	// Validate path to prevent directory traversal attacks
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("invalid path: path traversal detected")
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	// Library settings
	if path := os.Getenv("GOLLAMA_LIBRARY_PATH"); path != "" {
		config.LibraryPath = path
	}
	if embedded := os.Getenv("GOLLAMA_USE_EMBEDDED"); embedded != "" {
		config.UseEmbedded = parseEnvBool(embedded, config.UseEmbedded)
	}
	if logging := os.Getenv("GOLLAMA_ENABLE_LOGGING"); logging != "" {
		config.EnableLogging = parseEnvBool(logging, config.EnableLogging)
	}
	if level := os.Getenv("GOLLAMA_LOG_LEVEL"); level != "" {
		if val, err := strconv.Atoi(level); err == nil {
			config.LogLevel = val
		}
	}

	// Performance settings
	if threads := os.Getenv("GOLLAMA_NUM_THREADS"); threads != "" {
		if val, err := strconv.Atoi(threads); err == nil && val > 0 {
			config.NumThreads = val
		}
	}
	if gpu := os.Getenv("GOLLAMA_ENABLE_GPU"); gpu != "" {
		config.EnableGPU = parseEnvBool(gpu, config.EnableGPU)
	}
	if layers := os.Getenv("GOLLAMA_GPU_LAYERS"); layers != "" {
		if val, err := strconv.Atoi(layers); err == nil {
			config.GPULayers = val
		}
	}
	if metal := os.Getenv("GOLLAMA_METAL_ENABLED"); metal != "" {
		config.MetalEnabled = parseEnvBool(metal, config.MetalEnabled)
	}
	if cuda := os.Getenv("GOLLAMA_CUDA_ENABLED"); cuda != "" {
		config.CUDAEnabled = parseEnvBool(cuda, config.CUDAEnabled)
	}
	if vulkan := os.Getenv("GOLLAMA_VULKAN_ENABLED"); vulkan != "" {
		config.VulkanEnabled = parseEnvBool(vulkan, config.VulkanEnabled)
	}

	// Memory settings
	if ctx := os.Getenv("GOLLAMA_CONTEXT_SIZE"); ctx != "" {
		if val, err := strconv.Atoi(ctx); err == nil && val > 0 {
			config.ContextSize = val
		}
	}
	if batch := os.Getenv("GOLLAMA_BATCH_SIZE"); batch != "" {
		if val, err := strconv.Atoi(batch); err == nil && val > 0 {
			config.BatchSize = val
		}
	}
	if ubatch := os.Getenv("GOLLAMA_UBATCH_SIZE"); ubatch != "" {
		if val, err := strconv.Atoi(ubatch); err == nil && val > 0 {
			config.UbatchSize = val
		}
	}
	if mmap := os.Getenv("GOLLAMA_MEMORY_MAP_ENABLED"); mmap != "" {
		config.MemoryMapEnabled = parseEnvBool(mmap, config.MemoryMapEnabled)
	}
	if mlock := os.Getenv("GOLLAMA_MEMORY_LOCK_ENABLED"); mlock != "" {
		config.MemoryLockEnabled = parseEnvBool(mlock, config.MemoryLockEnabled)
	}

	// Model settings
	if path := os.Getenv("GOLLAMA_MODEL_PATH"); path != "" {
		config.ModelPath = path
	}
	if vocab := os.Getenv("GOLLAMA_VOCAB_ONLY"); vocab != "" {
		config.VocabOnly = parseEnvBool(vocab, config.VocabOnly)
	}
	if quant := os.Getenv("GOLLAMA_USE_QUANTIZATION"); quant != "" {
		config.UseQuantization = parseEnvBool(quant, config.UseQuantization)
	}
	if quantType := os.Getenv("GOLLAMA_QUANTIZATION_TYPE"); quantType != "" {
		config.QuantizationType = quantType
	}

	// Backend settings
	if backend := os.Getenv("GOLLAMA_BACKEND_TYPE"); backend != "" {
		config.BackendType = backend
	}
	if device := os.Getenv("GOLLAMA_DEVICE_ID"); device != "" {
		if val, err := strconv.Atoi(device); err == nil && val >= 0 {
			config.DeviceID = val
		}
	}

	// Debug settings
	if verbose := os.Getenv("GOLLAMA_VERBOSE_LOGGING"); verbose != "" {
		config.VerboseLogging = parseEnvBool(verbose, config.VerboseLogging)
	}
	if debug := os.Getenv("GOLLAMA_DEBUG_MODE"); debug != "" {
		config.DebugMode = parseEnvBool(debug, config.DebugMode)
	}

	return config
}

// SaveConfig saves configuration to a JSON file
func (c *Config) SaveConfig(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.NumThreads <= 0 {
		return fmt.Errorf("num_threads must be positive, got %d", c.NumThreads)
	}

	if c.ContextSize <= 0 {
		return fmt.Errorf("context_size must be positive, got %d", c.ContextSize)
	}

	if c.BatchSize <= 0 {
		return fmt.Errorf("batch_size must be positive, got %d", c.BatchSize)
	}

	if c.UbatchSize <= 0 {
		return fmt.Errorf("ubatch_size must be positive, got %d", c.UbatchSize)
	}

	if c.DeviceID < 0 {
		return fmt.Errorf("device_id must be non-negative, got %d", c.DeviceID)
	}

	// Validate library path if specified
	if c.LibraryPath != "" {
		if _, err := os.Stat(c.LibraryPath); os.IsNotExist(err) {
			return fmt.Errorf("library_path does not exist: %s", c.LibraryPath)
		}
	}

	// Validate model path if specified
	if c.ModelPath != "" {
		if _, err := os.Stat(c.ModelPath); os.IsNotExist(err) {
			return fmt.Errorf("model_path does not exist: %s", c.ModelPath)
		}
	}

	// Validate backend type
	validBackends := []string{"auto", "cpu", "gpu", "metal", "cuda", "vulkan", "opencl"}
	if c.BackendType != "" {
		found := false
		for _, backend := range validBackends {
			if c.BackendType == backend {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid backend_type: %s, must be one of: %v", c.BackendType, validBackends)
		}
	}

	// Validate quantization type if specified
	if c.QuantizationType != "" {
		validQuantTypes := []string{"q4_0", "q4_1", "q5_0", "q5_1", "q8_0", "q2_k", "q3_k", "q4_k", "q5_k", "q6_k", "q8_k"}
		found := false
		for _, quantType := range validQuantTypes {
			if c.QuantizationType == quantType {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid quantization_type: %s, must be one of: %v", c.QuantizationType, validQuantTypes)
		}
	}

	return nil
}

// GetConfigPath returns the default configuration file path
func GetConfigPath() string {
	if configDir := os.Getenv("XDG_CONFIG_HOME"); configDir != "" {
		return filepath.Join(configDir, "gollama", "config.json")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "gollama-config.json"
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "gollama", "config.json")
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "gollama", "config.json")
		}
		return filepath.Join(homeDir, "AppData", "Roaming", "gollama", "config.json")
	default:
		return filepath.Join(homeDir, ".config", "gollama", "config.json")
	}
}

// LoadDefaultConfig loads configuration from the default locations
func LoadDefaultConfig() *Config {
	// Start with environment variables
	config := LoadConfigFromEnv()

	// Try to load from config file
	configPath := GetConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		if fileConfig, err := LoadConfig(configPath); err == nil {
			// Merge file config with env config (env takes precedence)
			mergeConfigs(config, fileConfig)
		}
	}

	return config
}

// mergeConfigs merges source config into target config (target takes precedence)
func mergeConfigs(target, source *Config) {
	if target.LibraryPath == "" && source.LibraryPath != "" {
		target.LibraryPath = source.LibraryPath
	}
	if target.ModelPath == "" && source.ModelPath != "" {
		target.ModelPath = source.ModelPath
	}
	if target.BackendType == "auto" && source.BackendType != "" && source.BackendType != "auto" {
		target.BackendType = source.BackendType
	}
	if target.QuantizationType == "" && source.QuantizationType != "" {
		target.QuantizationType = source.QuantizationType
	}
}

// detectGPU attempts to detect if GPU acceleration is available
func detectGPU() bool {
	switch runtime.GOOS {
	case "darwin":
		// On macOS, Metal is usually available on modern systems
		return true
	case "linux", "windows":
		// On Linux/Windows, check for common GPU vendors
		// This is a simplified check - in a real implementation you might
		// check for CUDA/OpenCL/Vulkan libraries
		return checkForGPULibraries()
	default:
		return false
	}
}

// checkForGPULibraries checks for the presence of GPU acceleration libraries
func checkForGPULibraries() bool {
	// Check for CUDA
	cudaPaths := []string{
		"/usr/local/cuda/lib64/libcuda.so",
		"/usr/lib/x86_64-linux-gnu/libcuda.so",
		"C:\\Program Files\\NVIDIA GPU Computing Toolkit\\CUDA\\v11.0\\bin\\nvcuda.dll",
		"C:\\Windows\\System32\\nvcuda.dll",
	}

	for _, path := range cudaPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	// Check for OpenCL
	openclPaths := []string{
		"/usr/lib/x86_64-linux-gnu/libOpenCL.so",
		"/usr/lib/libOpenCL.so",
		"C:\\Windows\\System32\\OpenCL.dll",
	}

	for _, path := range openclPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// parseEnvBool parses a boolean environment variable with a default fallback
func parseEnvBool(value string, defaultValue bool) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "true", "1", "yes", "on", "enable", "enabled":
		return true
	case "false", "0", "no", "off", "disable", "disabled":
		return false
	default:
		return defaultValue
	}
}

// ApplyConfig applies the configuration to the library
func ApplyConfig(config *Config) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Set library path if specified
	if config.LibraryPath != "" {
		globalLoader.mutex.Lock()
		// Force reload with new path
		if globalLoader.loaded {
			_ = globalLoader.UnloadLibrary() // Ignore error during configuration
		}
		globalLoader.mutex.Unlock()
	}

	// Apply logging configuration
	// TODO: Implement logging configuration once we have the actual logging functions - moved to ROADMAP "wait for llama.cpp" section
	// if config.EnableLogging {
	//     llamaLogSet(logCallback, nil)
	// }

	return nil
}

// Global configuration instance
var globalConfig = LoadDefaultConfig()

// SetGlobalConfig sets the global configuration
func SetGlobalConfig(config *Config) error {
	if err := ApplyConfig(config); err != nil {
		return err
	}
	globalConfig = config
	return nil
}

// GetGlobalConfig returns the global configuration
func GetGlobalConfig() *Config {
	return globalConfig
}
