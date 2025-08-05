# Parallel Download & Checksum Demo

This example demonstrates the new parallel download and checksum verification features of gollama.cpp.

## Features Demonstrated

- **Parallel Downloads**: Download libraries for multiple platforms simultaneously
- **Checksum Verification**: Automatic SHA256 calculation and verification
- **Platform Detection**: Automatic detection of correct library files (.so, .dylib, .dll)
- **Error Handling**: Graceful handling of download failures

## Usage

```bash
go run main.go
```

## Expected Output

```
=== gollama.cpp Parallel Download & Checksum Demo ===

1. Downloading libraries for specific platforms...

Download Results:
================
✅ linux/amd64: SUCCESS
   Library: /home/user/.cache/gollama/libs/llama-b6089-bin-ubuntu-x64/build/bin/libllama.so
   SHA256: d3f76db17295aaebe984db4edd5b08a3a4da1106d32123d9dad1d640ae607622

✅ darwin/arm64: SUCCESS
   Library: /home/user/.cache/gollama/libs/llama-b6089-bin-macos-arm64/build/bin/libllama.dylib
   SHA256: e5ec9a20b0e77ba87ed5d8938e846ab5f03c3e11faeea23c38941508f3008ff8

✅ windows/amd64: SUCCESS
   Library: /home/user/.cache/gollama/libs/llama-b6089-bin-win-cpu-x64/llama.dll
   SHA256: 7e7d3de87806f0b780ecd9458da3afe0fe11bf8edb5e042aafec1d71ff9eb9e8

2. Calculating checksum of downloaded library...
File: /home/user/.cache/gollama/libs/llama-b6089-bin-ubuntu-x64/build/bin/libllama.so
SHA256: d3f76db17295aaebe984db4edd5b08a3a4da1106d32123d9dad1d640ae607622

=== Demo Complete ===
Features demonstrated:
✅ Parallel downloads for multiple platforms
✅ Automatic SHA256 checksum calculation
✅ Platform-specific library detection
✅ Concurrent download processing
✅ Error handling and reporting
```

## Command Line Usage

For more control, use the command-line tool:

```bash
# Download for all platforms with checksums
go run ../../cmd/gollama-download -download-all -checksum

# Download for specific platforms
go run ../../cmd/gollama-download -platforms "linux/amd64,darwin/arm64,windows/amd64" -checksum

# Verify checksum of a file
go run ../../cmd/gollama-download -verify-checksum /path/to/file

# Use Makefile targets
make -C ../.. download-libs-parallel
make -C ../.. download-libs-platforms
```

## Performance Benefits

- **Parallel Processing**: Downloads multiple platforms simultaneously (up to 4 concurrent)
- **Reduced Wait Time**: Total download time is reduced from sequential to maximum single download time
- **Integrity Verification**: SHA256 checksums ensure downloaded files are not corrupted
- **Caching**: Downloaded libraries are cached and reused across downloads
