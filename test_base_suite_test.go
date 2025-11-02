package gollama

import (
	"os"

	"github.com/stretchr/testify/suite"
)

// BaseSuite provides shared setup/teardown for all test suites.
// Responsibilities:
// - Snapshot and restore global configuration between tests
// - Snapshot and restore key environment variables used by tests
// - Ensure the llama library is unloaded after each test to avoid cross-test state
type BaseSuite struct {
	suite.Suite

	savedConfig *Config
	savedEnv    map[string]string
}

// envKeys are the environment variables we preserve across tests
var envKeys = []string{
	"GOLLAMA_CACHE_DIR",
	"GOLLAMA_LIBRARY_PATH",
	"GOLLAMA_USE_EMBEDDED",
	"GOLLAMA_ENABLE_LOGGING",
	"GOLLAMA_LOG_LEVEL",
}

// SetupTest runs before each test
func (s *BaseSuite) SetupTest() {
	// Snapshot global config (make a copy of the struct)
	if cfg := GetGlobalConfig(); cfg != nil {
		copied := *cfg
		s.savedConfig = &copied
	} else {
		s.savedConfig = DefaultConfig()
	}

	// Snapshot selected environment variables
	s.savedEnv = make(map[string]string, len(envKeys))
	for _, k := range envKeys {
		s.savedEnv[k] = os.Getenv(k)
	}
}

// TearDownTest runs after each test
func (s *BaseSuite) TearDownTest() {
	// Restore environment variables
	for k, v := range s.savedEnv {
		if v == "" {
			_ = os.Unsetenv(k)
		} else {
			_ = os.Setenv(k, v)
		}
	}

	// Restore global config
	if s.savedConfig != nil {
		_ = SetGlobalConfig(s.savedConfig)
	}

	// Ensure the library is unloaded to prevent cross-test contamination
	Cleanup()
}
