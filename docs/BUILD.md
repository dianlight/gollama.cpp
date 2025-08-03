# Build Guide

This guide explains how to build gollama.cpp from source on different platforms.

## Prerequisites

### Common Requirements

- Go 1.21 or later
- Git
- Make
- CMake 3.14 or later

### Platform-Specific Requirements

#### macOS
```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install Homebrew (if not already installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install dependencies
brew install cmake git
```

#### Linux (Ubuntu/Debian)
```bash
# Install build tools
sudo apt-get update
sudo apt-get install -y build-essential cmake git

# For NVIDIA GPU support (optional)
# Install CUDA Toolkit from https://developer.nvidia.com/cuda-toolkit

# For AMD GPU support (optional)
# Install ROCm from https://rocm.docs.amd.com/
```

#### Linux (CentOS/RHEL/Fedora)
```bash
# CentOS/RHEL
sudo yum groupinstall "Development Tools"
sudo yum install cmake git

# Fedora
sudo dnf groupinstall "Development Tools"
sudo dnf install cmake git
```

#### Windows
```bash
# Install using Chocolatey
choco install cmake git mingw

# Or install Visual Studio 2019/2022 with C++ support
# Download from https://visualstudio.microsoft.com/
```

## Quick Build

The fastest way to build gollama.cpp:

```bash
# Clone the repository
git clone https://github.com/ltarantino/gollama.cpp
cd gollama.cpp

# Build everything (Go bindings + llama.cpp libraries)
make build build-llamacpp-current

# Run tests
make test
```

## Detailed Build Process

### 1. Clone and Setup

```bash
git clone https://github.com/ltarantino/gollama.cpp
cd gollama.cpp

# Download Go dependencies
make deps
```

### 2. Build llama.cpp Libraries

#### Current Platform Only
```bash
# Build for your current platform
make build-llamacpp-current
```

#### All Platforms
```bash
# Build for all supported platforms (requires cross-compilation setup)
make build-llamacpp-all
```

#### Specific Platforms
```bash
# macOS (with Metal support)
make build-llamacpp-darwin-amd64
make build-llamacpp-darwin-arm64

# Linux (with CUDA support where available)
make build-llamacpp-linux-amd64
make build-llamacpp-linux-arm64

# Windows (with CUDA support where available)
make build-llamacpp-windows-amd64
```

#### GPU-Specific Builds
```bash
# Build with GPU acceleration (CUDA, Metal, etc.)
make build-libs-gpu

# Build CPU-only versions
make build-libs-cpu
```

### 3. Build Go Bindings

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build specific platform
GOOS=linux GOARCH=amd64 make build
```

### 4. Build Examples

```bash
make build-examples
```

## Custom Build Configuration

### Environment Variables

You can customize the build with environment variables:

```bash
# Specify llama.cpp version
export LLAMA_CPP_BUILD=b6076

# Specify Go version
export GO_VERSION=1.21

# Platform targeting
export GOOS=linux
export GOARCH=amd64

make build
```

### CMake Options for llama.cpp

When building llama.cpp, you can pass custom CMake flags:

#### GPU Acceleration
```bash
# NVIDIA CUDA
cmake -DGGML_CUDA=ON

# AMD HIP/ROCm
cmake -DGGML_HIP=ON

# Apple Metal
cmake -DGGML_METAL=ON

# Intel oneAPI
cmake -DGGML_SYCL=ON

# Vulkan
cmake -DGGML_VULKAN=ON
```

#### CPU Optimizations
```bash
# AVX support
cmake -DGGML_AVX=ON -DGGML_AVX2=ON -DGGML_AVX512=ON

# ARM NEON
cmake -DGGML_NEON=ON

# Disable optimizations for compatibility
cmake -DGGML_NATIVE=OFF
```

#### Other Options
```bash
# Build shared libraries (required for gollama.cpp)
cmake -DBUILD_SHARED_LIBS=ON

# Build tools
cmake -DBUILD_TOOLS=ON

# Enable debug info
cmake -DCMAKE_BUILD_TYPE=Debug
```

## Cross-Compilation

### Linux ARM64 from x86_64

```bash
# Install cross-compilation tools
sudo apt-get install gcc-aarch64-linux-gnu g++-aarch64-linux-gnu

# Build llama.cpp
cd build/llama.cpp
cmake -B build-linux-arm64 \
  -DCMAKE_SYSTEM_NAME=Linux \
  -DCMAKE_SYSTEM_PROCESSOR=aarch64 \
  -DCMAKE_C_COMPILER=aarch64-linux-gnu-gcc \
  -DCMAKE_CXX_COMPILER=aarch64-linux-gnu-g++ \
  -DBUILD_SHARED_LIBS=ON
cmake --build build-linux-arm64

# Build Go bindings
GOOS=linux GOARCH=arm64 go build
```

### Windows from Linux (MinGW)

```bash
# Install MinGW cross-compiler
sudo apt-get install mingw-w64

# Build llama.cpp
cd build/llama.cpp
cmake -B build-windows-amd64 \
  -DCMAKE_SYSTEM_NAME=Windows \
  -DCMAKE_C_COMPILER=x86_64-w64-mingw32-gcc \
  -DCMAKE_CXX_COMPILER=x86_64-w64-mingw32-g++ \
  -DBUILD_SHARED_LIBS=ON
cmake --build build-windows-amd64

# Build Go bindings
GOOS=windows GOARCH=amd64 go build
```

## Troubleshooting

### Common Issues

#### 1. CMake Not Found
```bash
# Solution: Install CMake
# macOS: brew install cmake
# Ubuntu: sudo apt-get install cmake
# Windows: choco install cmake
```

#### 2. CUDA Not Found
```bash
# Solution: Install CUDA Toolkit and set environment variables
export CUDA_PATH=/usr/local/cuda
export PATH=$CUDA_PATH/bin:$PATH
export LD_LIBRARY_PATH=$CUDA_PATH/lib64:$LD_LIBRARY_PATH
```

#### 3. Metal Framework Not Found (macOS)
```bash
# Solution: Install Xcode Command Line Tools
xcode-select --install
```

#### 4. Go Module Issues
```bash
# Solution: Clean and reinstall modules
go clean -modcache
go mod download
go mod tidy
```

#### 5. Library Loading Issues
```bash
# Linux: Add library path
export LD_LIBRARY_PATH=$PWD/libs/linux_amd64:$LD_LIBRARY_PATH

# macOS: Add library path
export DYLD_LIBRARY_PATH=$PWD/libs/darwin_amd64:$DYLD_LIBRARY_PATH

# Windows: Add library path
set PATH=%CD%\libs\windows_amd64;%PATH%
```

### Build Logs

Enable verbose build output for debugging:

```bash
# CMake verbose
cmake --build . --verbose

# Make verbose
make V=1

# Go verbose
go build -v
```

### Platform-Specific Notes

#### macOS
- Metal support requires macOS 10.13+ and Xcode 9+
- Universal binaries can be built with `-DCMAKE_OSX_ARCHITECTURES="x86_64;arm64"`
- Code signing may be required for distribution

#### Linux
- CUDA support requires NVIDIA driver 450+ and CUDA 11.0+
- ROCm support requires AMD drivers and ROCm 4.0+
- Some distributions may need additional development packages

#### Windows
- Visual Studio 2019+ or MinGW-w64 is required
- CUDA support requires Visual Studio integration
- PATH environment variable must include library directory

## Performance Optimization

### Build Flags

For maximum performance:

```bash
# Go build flags
go build -ldflags="-s -w"

# CMake release build
cmake -DCMAKE_BUILD_TYPE=Release -DGGML_NATIVE=ON
```

### CPU Architecture Specific

```bash
# For modern Intel/AMD CPUs
cmake -DGGML_AVX=ON -DGGML_AVX2=ON -DGGML_F16C=ON -DGGML_FMA=ON

# For ARM64
cmake -DGGML_NEON=ON

# For compatibility (slower but works everywhere)
cmake -DGGML_NATIVE=OFF
```

## Contributing

When contributing builds:

1. Test on multiple platforms if possible
2. Document any platform-specific requirements
3. Update this guide with new build instructions
4. Ensure CI/CD pipeline passes

## Additional Resources

- [llama.cpp Build Documentation](https://github.com/ggml-org/llama.cpp#build)
- [Go Cross Compilation](https://golang.org/doc/install/cross)
- [CMake Documentation](https://cmake.org/documentation/)
- [CUDA Installation Guide](https://docs.nvidia.com/cuda/cuda-installation-guide-linux/)
- [ROCm Installation Guide](https://rocm.docs.amd.com/)
