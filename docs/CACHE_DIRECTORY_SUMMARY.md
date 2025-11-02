# Cache Directory Configuration - Summary

## Overview

Added comprehensive support for configuring cache directories where downloaded llama.cpp library binaries are stored. Users can now customize the cache location through multiple methods with a clear priority hierarchy.

## Changes Made

### 1. Core Configuration (`config.go`)
- Added `CacheDir` field to the `Config` struct
- Added environment variable support (`GOLLAMA_CACHE_DIR`)
- Added validation for cache directory paths (prevents path traversal attacks)
- Updated `mergeConfigs` to handle cache directory merging

### 2. Downloader (`downloader.go`)
- Added `NewLibraryDownloaderWithCacheDir()` function for custom cache directories
- Updated `NewLibraryDownloader()` to check for `GOLLAMA_CACHE_DIR` environment variable
- Added `GetCacheDir()` method to retrieve the current cache directory
- Cache directory priority:
  1. Custom cache directory passed to `NewLibraryDownloaderWithCacheDir()`
  2. `GOLLAMA_CACHE_DIR` environment variable
  3. Platform-specific defaults

### 3. Loader (`loader.go`)
- Updated all downloader initialization to respect global config's cache directory
- Added `GetLibraryCacheDir()` public function to retrieve current cache directory
- Modified `LoadLibraryWithVersion()`, `DownloadLibrariesForPlatforms()`, and `GetSHA256ForFile()` to use config cache directory

### 4. Tests (`downloader_cache_test.go`)
Added comprehensive test coverage:
- `TestCacheDirectoryConfiguration`: Tests default, custom, environment, and config-based cache directories
- `TestCacheDirValidation`: Tests path traversal security validation

All tests pass successfully.

### 5. Documentation (`README.md`)
Enhanced cache location documentation with:
- Default platform-specific cache directories
- Environment variable configuration example
- Configuration file example
- Programmatic configuration example
- Code example to retrieve current cache directory

### 6. Example (`examples/cache-directory-demo/`)
Created a comprehensive demonstration example showing:
- Default cache directory retrieval
- Environment variable configuration
- Config object configuration
- JSON configuration file usage
- Cache cleaning
- Configuration priority explanation

## Usage Examples

### Get Current Cache Directory
```go
cacheDir, err := gollama.GetLibraryCacheDir()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Cache directory: %s\n", cacheDir)
```

### Configure via Environment Variable
```bash
export GOLLAMA_CACHE_DIR=/custom/path/to/cache
```

### Configure via Config Object
```go
config := gollama.DefaultConfig()
config.CacheDir = "/custom/path/to/cache"
gollama.SetGlobalConfig(config)
```

### Configure via JSON File
```json
{
  "cache_dir": "/custom/path/to/cache",
  "enable_logging": true,
  "num_threads": 8
}
```

Then load it:
```go
config, err := gollama.LoadConfig("config.json")
if err != nil {
    log.Fatal(err)
}
gollama.SetGlobalConfig(config)
```

## Configuration Priority

The cache directory is determined in the following priority order (highest to lowest):

1. **Config.CacheDir** - Set in global config object (highest priority)
2. **GOLLAMA_CACHE_DIR** - Environment variable
3. **Platform Default** - System-specific cache directory:
   - Linux/Unix: `~/.cache/gollama/libs/`
   - macOS: `~/Library/Caches/gollama/libs/`
   - Windows: `%LOCALAPPDATA%\gollama\libs\`
   - Fallback: `<TEMP>/gollama/libs/`

## Security Features

- Path traversal attack prevention through validation
- Clean path normalization
- Proper directory permissions (0750)
- Input validation in `Config.Validate()`

## Backward Compatibility

All changes are backward compatible:
- Default behavior unchanged if no custom cache directory is specified
- Existing code continues to work without modifications
- New functionality is opt-in through configuration

## Testing

All tests pass successfully:
- ✅ Default cache directory detection
- ✅ Custom cache directory creation
- ✅ Environment variable configuration
- ✅ Config object configuration
- ✅ Path traversal validation
- ✅ Existing functionality unchanged

## Files Modified

1. `config.go` - Added cache directory configuration support
2. `downloader.go` - Added cache directory customization
3. `loader.go` - Updated to use configured cache directory
4. `README.md` - Enhanced documentation
5. `examples/README.md` - Added cache-directory-demo section

## Files Created

1. `downloader_cache_test.go` - Comprehensive test suite
2. `examples/cache-directory-demo/main.go` - Demonstration example
3. `examples/cache-directory-demo/go.mod` - Module file
4. `examples/cache-directory-demo/README.md` - Example documentation
5. `CACHE_DIRECTORY_SUMMARY.md` - This summary document

## Benefits

1. **Flexibility**: Users can customize cache location for their needs
2. **Environment Integration**: Works well with CI/CD and containerized environments
3. **Disk Management**: Allows users to control where large files are stored
4. **Multi-tenant Support**: Different configurations for different applications
5. **Security**: Validates paths to prevent security issues
6. **Documentation**: Clear examples and priority hierarchy
