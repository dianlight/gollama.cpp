package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dianlight/gollama.cpp"
)

func main() {
	fmt.Println("=== gollama.cpp Cache Directory Configuration Demo ===\n")

	// Example 1: Get default cache directory
	fmt.Println("1. Default Cache Directory:")
	cacheDir, err := gollama.GetLibraryCacheDir()
	if err != nil {
		log.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Cache directory: %s\n\n", cacheDir)
	}

	// Example 2: Configure via environment variable
	fmt.Println("2. Configure via Environment Variable:")
	customEnvCache := "/tmp/gollama_env_cache"
	os.Setenv("GOLLAMA_CACHE_DIR", customEnvCache)
	fmt.Printf("  Set GOLLAMA_CACHE_DIR=%s\n", customEnvCache)

	// Create a new downloader to pick up the env var
	downloader, err := gollama.NewLibraryDownloader()
	if err != nil {
		log.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Cache directory: %s\n\n", downloader.GetCacheDir())
	}
	os.Unsetenv("GOLLAMA_CACHE_DIR")

	// Example 3: Configure via Config object
	fmt.Println("3. Configure via Config Object:")
	config := gollama.DefaultConfig()
	customConfigCache := "/tmp/gollama_config_cache"
	config.CacheDir = customConfigCache
	fmt.Printf("  Set config.CacheDir=%s\n", customConfigCache)

	if err := gollama.SetGlobalConfig(config); err != nil {
		log.Printf("  Error: %v\n", err)
	} else {
		cacheDir, _ := gollama.GetLibraryCacheDir()
		fmt.Printf("  Cache directory: %s\n\n", cacheDir)
	}

	// Example 4: Load config from file (simulated)
	fmt.Println("4. Configuration File Support:")
	fmt.Println("  Create a JSON config file:")
	fmt.Println(`  {
    "cache_dir": "/custom/path/to/cache",
    "enable_logging": true,
    "num_threads": 8
  }`)
	fmt.Println("  Then load it:")
	fmt.Println(`  config, err := gollama.LoadConfig("config.json")
  if err != nil {
    log.Fatal(err)
  }
  gollama.SetGlobalConfig(config)`)
	fmt.Println()

	// Example 5: Clean cache
	fmt.Println("5. Clean Cache:")
	fmt.Println("  To force re-download of libraries, clean the cache:")
	fmt.Println(`  err := gollama.CleanLibraryCache()
  if err != nil {
    log.Fatal(err)
  }`)
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
	fmt.Println("\nCache Directory Priority:")
	fmt.Println("  1. Config.CacheDir (highest priority)")
	fmt.Println("  2. GOLLAMA_CACHE_DIR environment variable")
	fmt.Println("  3. Platform default (~/.cache/gollama/libs on Unix)")
}
