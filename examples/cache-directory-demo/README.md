# Cache Directory Configuration Demo

This example demonstrates how to configure and manage the cache directory used by gollama.cpp for downloading and storing llama.cpp library binaries.

## Features

- Get the default cache directory location
- Configure cache directory via environment variable
- Configure cache directory via Config object
- Configure cache directory via JSON configuration file
- Clean cache to force re-download

## Usage

```bash
# Run the demo
go run main.go

# With custom environment variable
GOLLAMA_CACHE_DIR=/tmp/my_cache go run main.go
```

## Output Example

```
=== gollama.cpp Cache Directory Configuration Demo ===

1. Default Cache Directory:
  Cache directory: /home/user/.cache/gollama/libs

2. Configure via Environment Variable:
  Set GOLLAMA_CACHE_DIR=/tmp/gollama_env_cache
  Cache directory: /tmp/gollama_env_cache/libs

3. Configure via Config Object:
  Set config.CacheDir=/tmp/gollama_config_cache
  Cache directory: /tmp/gollama_config_cache

4. Configuration File Support:
  Create a JSON config file:
  {
    "cache_dir": "/custom/path/to/cache",
    "enable_logging": true,
    "num_threads": 8
  }
  Then load it:
  config, err := gollama.LoadConfig("config.json")
  if err != nil {
    log.Fatal(err)
  }
  gollama.SetGlobalConfig(config)

5. Clean Cache:
  To force re-download of libraries, clean the cache:
  err := gollama.CleanLibraryCache()
  if err != nil {
    log.Fatal(err)
  }

=== Demo Complete ===

Cache Directory Priority:
  1. Config.CacheDir (highest priority)
  2. GOLLAMA_CACHE_DIR environment variable
  3. Platform default (~/.cache/gollama/libs on Unix)
```

## Configuration Methods

### 1. Environment Variable

The simplest way to configure the cache directory:

```bash
export GOLLAMA_CACHE_DIR=/path/to/cache
```

### 2. Config Object

For programmatic configuration:

```go
config := gollama.DefaultConfig()
config.CacheDir = "/path/to/cache"
gollama.SetGlobalConfig(config)
```

### 3. Configuration File

Create a `config.json` file:

```json
{
  "cache_dir": "/path/to/cache",
  "enable_logging": true,
  "num_threads": 8,
  "context_size": 4096
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

## Cache Directory Structure

The cache directory contains downloaded and extracted library binaries:

```
~/.cache/gollama/libs/
├── llama-b6862-bin-macos-arm64/
│   └── build/
│       └── bin/
│           └── libllama.dylib
├── llama-b6862-bin-ubuntu-x64/
│   └── build/
│       └── bin/
│           └── libllama.so
└── llama-b6862-bin-win-cpu-x64/
    └── llama.dll
```

## Cleaning the Cache

To force a re-download of libraries:

```go
err := gollama.CleanLibraryCache()
if err != nil {
    log.Fatal(err)
}
```

Or from the command line:

```bash
make clean-libs
```

## Platform-Specific Defaults

| Platform | Default Cache Directory          |
| -------- | -------------------------------- |
| Linux    | `~/.cache/gollama/libs/`         |
| macOS    | `~/Library/Caches/gollama/libs/` |
| Windows  | `%LOCALAPPDATA%\gollama\libs\`   |

## Security

- The cache directory validation prevents path traversal attacks
- Downloaded files are verified with SHA256 checksums
- File permissions are set to `0750` for directories and preserve original permissions for extracted files
