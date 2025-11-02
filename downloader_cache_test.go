package gollama

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCacheDirectoryConfiguration(t *testing.T) {
	t.Run("Default cache directory", func(t *testing.T) {
		downloader, err := NewLibraryDownloader()
		if err != nil {
			t.Fatalf("Failed to create downloader: %v", err)
		}

		cacheDir := downloader.GetCacheDir()
		if cacheDir == "" {
			t.Error("Cache directory should not be empty")
		}

		// Should contain "gollama" in the path
		if !strings.Contains(cacheDir, "gollama") {
			t.Errorf("Cache directory should contain 'gollama': %s", cacheDir)
		}
	})

	t.Run("Custom cache directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		customCache := filepath.Join(tmpDir, "custom_cache")

		downloader, err := NewLibraryDownloaderWithCacheDir(customCache)
		if err != nil {
			t.Fatalf("Failed to create downloader with custom cache: %v", err)
		}

		cacheDir := downloader.GetCacheDir()
		if cacheDir != customCache {
			t.Errorf("Expected cache dir %s, got %s", customCache, cacheDir)
		}

		// Verify directory was created
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			t.Errorf("Cache directory was not created: %s", cacheDir)
		}
	})

	t.Run("Environment variable cache directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		envCache := filepath.Join(tmpDir, "env_cache")

		// Set environment variable
		oldEnv := os.Getenv("GOLLAMA_CACHE_DIR")
		if err := os.Setenv("GOLLAMA_CACHE_DIR", envCache); err != nil {
			t.Fatalf("Failed to set environment variable: %v", err)
		}
		defer func() {
			if err := os.Setenv("GOLLAMA_CACHE_DIR", oldEnv); err != nil {
				t.Errorf("Failed to restore environment variable: %v", err)
			}
		}()

		downloader, err := NewLibraryDownloader()
		if err != nil {
			t.Fatalf("Failed to create downloader: %v", err)
		}

		cacheDir := downloader.GetCacheDir()
		expectedPath := filepath.Join(envCache, "libs")
		if cacheDir != expectedPath {
			t.Errorf("Expected cache dir %s, got %s", expectedPath, cacheDir)
		}
	})

	t.Run("Config cache directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		configCache := filepath.Join(tmpDir, "config_cache")

		config := DefaultConfig()
		config.CacheDir = configCache
		_ = SetGlobalConfig(config)
		defer func() { _ = SetGlobalConfig(LoadDefaultConfig()) }() // Reset to default

		cacheDir, err := GetLibraryCacheDir()
		if err != nil {
			t.Fatalf("Failed to get cache directory: %v", err)
		}

		if cacheDir != configCache {
			t.Errorf("Expected cache dir %s, got %s", configCache, cacheDir)
		}
	})
}

func TestCacheDirValidation(t *testing.T) {
	t.Run("Valid cache directory in config", func(t *testing.T) {
		tmpDir := t.TempDir()
		config := DefaultConfig()
		config.CacheDir = tmpDir

		err := config.Validate()
		if err != nil {
			t.Errorf("Expected no error for valid cache dir, got: %v", err)
		}
	})

	t.Run("Path traversal in cache directory", func(t *testing.T) {
		config := DefaultConfig()
		config.CacheDir = "../../../etc/passwd"

		err := config.Validate()
		if err == nil {
			t.Error("Expected error for path traversal in cache_dir")
		}
		if !strings.Contains(err.Error(), "path traversal") {
			t.Errorf("Expected 'path traversal' error, got: %v", err)
		}
	})
}
