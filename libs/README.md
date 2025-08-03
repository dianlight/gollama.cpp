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
├── darwin_amd64/
│   ├── libggml.dylib
│   ├── libggml-base.dylib
│   ├── libggml-blas.dylib
│   ├── libggml-cpu.dylib
│   ├── libggml-metal.dylib
│   ├── libllama.dylib
│   └── libmtmd.dylib
├── darwin_arm64/
│   ├── libggml.dylib
│   ├── libggml-base.dylib
│   ├── libggml-blas.dylib
│   ├── libggml-cpu.dylib
│   ├── libggml-metal.dylib
│   ├── libllama.dylib
│   └── libmtmd.dylib
├── linux_amd64/
│   ├── libggml.so
│   ├── libggml-base.so
│   ├── libggml-blas.so
│   ├── libggml-cpu.so
│   ├── libggml-cuda.so
│   ├── libllama.so
│   └── libmtmd.so
├── linux_arm64/
│   ├── libggml.so
│   ├── libggml-base.so
│   ├── libggml-blas.so
│   ├── libggml-cpu.so
│   ├── libllama.so
│   └── libmtmd.so
├── windows_amd64/
│   ├── ggml.dll
│   ├── ggml-base.dll
│   ├── ggml-blas.dll
│   ├── ggml-cpu.dll
│   ├── ggml-cuda.dll
│   ├── llama.dll
│   └── mtmd.dll
└── windows_arm64/
    ├── ggml.dll
    ├── ggml-base.dll
    ├── ggml-blas.dll
    ├── ggml-cpu.dll
    ├── llama.dll
    └── mtmd.dll
```

These libraries will be embedded into the Go binary and extracted at runtime when needed.
The build system will automatically populate this directory during the release process.
