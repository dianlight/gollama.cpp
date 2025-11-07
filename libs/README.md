# Library Directory

This directory contains the pre-built llama.cpp libraries and their dependencies for different platforms.

The libraries include:
- `libllama` - Main llama.cpp library
- `libggml` - Core GGML library
- `libggml-base` - Base GGML components
- `libggml-blas` - BLAS acceleration support
- `libggml-cpu` - CPU-specific optimizations
- `libggml-metal` - Metal acceleration (macOS only)
- `libggml-cuda` - CUDA acceleration (Linux only)
- `libmtmd` - Multi-threading support

All libraries are configured with proper rpath settings to ensure correct dependency resolution at runtime.

Expected structure:
```
libs/
├── darwin_amd64_<version>/
│   ├── libggml.dylib
│   ├── libggml-base.dylib
│   ├── libggml-blas.dylib
│   ├── libggml-cpu.dylib
│   ├── libggml-metal.dylib
│   ├── libllama.dylib
│   └── libmtmd.dylib
├── darwin_arm64_<version>/
│   ├── libggml.dylib
│   ├── libggml-base.dylib
│   ├── libggml-blas.dylib
│   ├── libggml-cpu.dylib
│   ├── libggml-metal.dylib
│   ├── libllama.dylib
│   └── libmtmd.dylib
├── linux_amd64_<version>/
│   ├── libggml.so
│   ├── libggml-base.so
│   ├── libggml-blas.so
│   ├── libggml-cpu.so
│   ├── libggml-cuda.so
│   ├── libllama.so
│   └── libmtmd.so
├── linux_arm64_<version>/
│   ├── libggml.so
│   ├── libggml-base.so
│   ├── libggml-blas.so
│   ├── libggml-cpu.so
│   ├── libllama.so
│   └── libmtmd.so
├── windows_amd64_<version>/
│   ├── ggml.dll
│   ├── ggml-base.dll
│   ├── ggml-blas.dll
│   ├── ggml-cpu-*.dll       # Multiple CPU architecture variants (x64, sse42, alderlake, etc.)
│   ├── ggml-cuda.dll
│   ├── llama.dll
│   └── mtmd.dll
└── windows_arm64_<version>/
    ├── ggml.dll
    ├── ggml-base.dll
    ├── ggml-blas.dll
    ├── ggml-cpu-*.dll       # Multiple CPU architecture variants
    ├── llama.dll
    └── mtmd.dll
```

These libraries are embedded into the Go binary (via `go:embed`) and extracted to the cache at runtime when needed. Use the tooling below to keep the directory in sync with the configured llama.cpp build:

```bash
# Download all supported platforms for the configured llama.cpp build and refresh ./libs
make populate-libs

# Or run the downloader directly
go run ./cmd/gollama-download -download-all -version b6862 -copy-libs
```

Only one llama.cpp version is stored at a time; older directories are removed automatically during the sync.
