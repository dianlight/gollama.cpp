package gollama

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CacheDirSuite struct {
	suite.Suite
}

func (s *CacheDirSuite) TestDefaultCacheDirectory() {
	downloader, err := NewLibraryDownloader()
	assert.NoError(s.T(), err, "Failed to create downloader")

	cacheDir := downloader.GetCacheDir()
	assert.NotEmpty(s.T(), cacheDir, "Cache directory should not be empty")
	assert.Contains(s.T(), cacheDir, "gollama", "Cache directory should contain 'gollama'")
}

func (s *CacheDirSuite) TestCustomCacheDirectory() {
	tmpDir := s.T().TempDir()
	customCache := filepath.Join(tmpDir, "custom_cache")

	downloader, err := NewLibraryDownloaderWithCacheDir(customCache)
	assert.NoError(s.T(), err, "Failed to create downloader with custom cache")

	cacheDir := downloader.GetCacheDir()
	assert.Equal(s.T(), customCache, cacheDir)

	// Verify directory was created
	_, statErr := os.Stat(cacheDir)
	assert.False(s.T(), os.IsNotExist(statErr), "Cache directory was not created: %s", cacheDir)
}

func (s *CacheDirSuite) TestEnvironmentVariableCacheDirectory() {
	tmpDir := s.T().TempDir()
	envCache := filepath.Join(tmpDir, "env_cache")

	// Set environment variable
	oldEnv := os.Getenv("GOLLAMA_CACHE_DIR")
	setErr := os.Setenv("GOLLAMA_CACHE_DIR", envCache)
	assert.NoError(s.T(), setErr, "Failed to set environment variable")
	s.T().Cleanup(func() { _ = os.Setenv("GOLLAMA_CACHE_DIR", oldEnv) })

	downloader, err := NewLibraryDownloader()
	assert.NoError(s.T(), err, "Failed to create downloader")

	cacheDir := downloader.GetCacheDir()
	expectedPath := filepath.Join(envCache, "libs")
	assert.Equal(s.T(), expectedPath, cacheDir)
}

func (s *CacheDirSuite) TestConfigCacheDirectory() {
	tmpDir := s.T().TempDir()
	configCache := filepath.Join(tmpDir, "config_cache")

	config := DefaultConfig()
	config.CacheDir = configCache
	_ = SetGlobalConfig(config)
	s.T().Cleanup(func() { _ = SetGlobalConfig(LoadDefaultConfig()) })

	cacheDir, err := GetLibraryCacheDir()
	assert.NoError(s.T(), err, "Failed to get cache directory")
	assert.Equal(s.T(), configCache, cacheDir)
}

func (s *CacheDirSuite) TestValidCacheDirectoryInConfig() {
	tmpDir := s.T().TempDir()
	config := DefaultConfig()
	config.CacheDir = tmpDir

	err := config.Validate()
	assert.NoError(s.T(), err, "Expected no error for valid cache dir")
}

func (s *CacheDirSuite) TestPathTraversalInCacheDirectory() {
	config := DefaultConfig()
	config.CacheDir = "../../../etc/passwd"

	err := config.Validate()
	assert.Error(s.T(), err, "Expected error for path traversal in cache_dir")
	if err != nil {
		assert.Contains(s.T(), err.Error(), "path traversal")
	}
}

func TestCacheDirSuite(t *testing.T) {
	suite.Run(t, new(CacheDirSuite))
}
