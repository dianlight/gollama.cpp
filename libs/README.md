# Library Directory

This directory will contain the pre-built llama.cpp libraries for different platforms.

Expected structure:
```
libs/
├── darwin_amd64/
│   └── libllama.dylib
├── darwin_arm64/
│   └── libllama.dylib
├── linux_amd64/
│   ├── libllama.so
│   ├── libllama-cuda.so
│   └── libllama-vulkan.so
├── linux_arm64/
│   └── libllama.so
├── windows_amd64/
│   ├── llama.dll
│   └── llama-cuda.dll
└── windows_arm64/
    └── llama.dll
```

These libraries will be embedded into the Go binary and extracted at runtime when needed.
The build system will automatically populate this directory during the release process.
