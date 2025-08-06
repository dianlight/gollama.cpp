# GPU Support Guide

Gollama.cpp provides comprehensive GPU acceleration support across multiple platforms and vendors. This guide covers installation, configuration, and troubleshooting for GPU acceleration.

## Overview

The library automatically detects available GPU hardware and configures the optimal backend during build time. No manual configuration is required for most setups.

### Supported GPU Backends

| Backend | Platforms | GPU Vendors | Status |
|---------|-----------|-------------|--------|
| **Metal** | macOS | Apple Silicon | ✅ Production |
| **CUDA** | Linux, Windows | NVIDIA | ✅ Production |
| **HIP/ROCm** | Linux, Windows | AMD | ✅ Production |
| **Vulkan** | Linux, Windows | NVIDIA, AMD, Intel | ✅ Production |
| **OpenCL** | Windows, Linux | Qualcomm Adreno, Intel, AMD | ✅ Production |
| **SYCL** | Linux, Windows | Intel, NVIDIA | ✅ Production |
| **CPU** | All | All | ✅ Fallback |

## Platform-Specific Setup

### macOS - Metal Support

Metal support is automatically enabled on macOS systems with Apple Silicon (M1/M2/M3).

**Requirements:**
- macOS 10.15+ (Catalina)
- Apple Silicon Mac (M1/M2/M3) or Intel Mac with Metal-compatible GPU
- Xcode Command Line Tools

**Installation:**
```bash
# Install Xcode Command Line Tools (if not already installed)
xcode-select --install

# Build with Metal support (automatic)
make build
```

**Verification:**
```bash
# Check Metal availability
system_profiler SPDisplaysDataType | grep Metal
```

### Linux - CUDA Support

CUDA support is automatically detected when NVIDIA CUDA Toolkit is installed.

**Requirements:**
- NVIDIA GPU with Compute Capability 3.5+
- CUDA Toolkit 11.8 or later
- Compatible NVIDIA driver

**Installation:**
```bash
# Ubuntu/Debian - Install CUDA Toolkit
wget https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2004/x86_64/cuda-keyring_1.0-1_all.deb
sudo dpkg -i cuda-keyring_1.0-1_all.deb
sudo apt-get update
sudo apt-get install cuda-toolkit

# Verify CUDA installation
nvcc --version
nvidia-smi

# Build with CUDA support (automatic detection)
make build
```

**Fedora/RHEL:**
```bash
# Enable NVIDIA repository
sudo dnf config-manager --add-repo https://developer.download.nvidia.com/compute/cuda/repos/fedora37/x86_64/cuda-fedora37.repo

# Install CUDA
sudo dnf install cuda-toolkit

# Build with CUDA support
make build
```

### Linux - AMD HIP/ROCm Support

HIP support is automatically detected when AMD ROCm is installed.

**Requirements:**
- AMD GPU with GCN 4th gen (gfx803) or newer
- ROCm 5.0 or later
- Compatible AMD driver (amdgpu)

**Installation:**
```bash
# Ubuntu/Debian - Install ROCm
wget -q -O - https://repo.radeon.com/rocm/rocm.gpg.key | sudo apt-key add -
echo 'deb [arch=amd64] https://repo.radeon.com/rocm/apt/debian/ ubuntu main' | sudo tee /etc/apt/sources.list.d/rocm.list
sudo apt-get update
sudo apt-get install rocm-dev hip-dev

# Add user to render group
sudo usermod -a -G render,video $USER

# Verify HIP installation
/opt/rocm/bin/hipconfig --platform
/opt/rocm/bin/rocm-smi

# Build with HIP support (automatic detection)
make build
```

### Windows - CUDA Support

**Requirements:**
- NVIDIA GPU with Compute Capability 3.5+
- CUDA Toolkit 11.8 or later
- Visual Studio 2019+ or compatible compiler

**Installation:**
1. Download and install [CUDA Toolkit](https://developer.nvidia.com/cuda-downloads)
2. Ensure `nvcc` is in your PATH
3. Build with automatic CUDA detection:

```powershell
# Verify CUDA installation
nvcc --version
nvidia-smi

# Build with CUDA support
make build
```

### Windows - AMD HIP Support

**Requirements:**
- AMD GPU with GCN 4th gen or newer
- HIP SDK for Windows
- Visual Studio 2019+ or compatible compiler

**Installation:**
1. Download and install [HIP SDK](https://github.com/ROCm-Developer-Tools/HIP)
2. Ensure HIP tools are in your PATH
3. Build with automatic HIP detection:

```powershell
# Verify HIP installation
hipconfig --platform

# Build with HIP support
make build
```

### Windows - Vulkan Support

Vulkan provides cross-platform GPU acceleration for NVIDIA, AMD, and Intel GPUs on Windows.

**Requirements:**
- Vulkan-capable GPU (NVIDIA GTX 600+, AMD GCN+, Intel HD 4000+)
- Latest GPU drivers with Vulkan support
- Vulkan SDK (optional, for development)

**Installation:**
```powershell
# Install Vulkan SDK (optional, for development)
# Download from: https://vulkan.lunarg.com/sdk/home

# Verify Vulkan support (if SDK installed)
vulkaninfo

# Or check driver support
# NVIDIA: GeForce Experience -> Drivers
# AMD: AMD Software -> Drivers
# Intel: Intel Graphics Command Center

# Build with Vulkan support (automatic detection)
make build
```

### Windows - OpenCL Support

OpenCL provides cross-platform parallel computing, especially useful for Qualcomm Adreno GPUs on ARM64.

**Requirements:**
- OpenCL-capable GPU or CPU
- Latest GPU drivers with OpenCL support

**Installation:**
```powershell
# For Intel GPUs
# Download Intel Graphics Driver from Intel website

# For AMD GPUs
# Install AMD Software (includes OpenCL support)

# For NVIDIA GPUs
# Install NVIDIA GPU drivers (includes OpenCL support)

# For Qualcomm Adreno (ARM64 devices)
# Usually pre-installed on ARM64 Windows devices

# Verify OpenCL support (if available)
# Install GPU-Z or similar tool to check OpenCL support

# Build with OpenCL support (automatic detection)
make build
```

### Windows - SYCL Support

SYCL provides unified parallel programming for CPUs, GPUs, and other accelerators.

**Requirements:**
- Intel oneAPI Toolkit or compatible SYCL implementation
- Compatible hardware (Intel GPUs, NVIDIA GPUs via CUDA backend)

**Installation:**
```powershell
# Install Intel oneAPI Toolkit
# Download from: https://www.intel.com/content/www/us/en/developer/tools/oneapi/toolkits.html

# Source the environment (in Developer Command Prompt)
"C:\Program Files (x86)\Intel\oneAPI\setvars.bat"

# Verify SYCL installation
sycl-ls

# Build with SYCL support (automatic detection)
make build
```

### Linux - Vulkan Support

Vulkan provides cross-platform GPU acceleration support for NVIDIA, AMD, and Intel GPUs.

**Requirements:**
- Vulkan-capable GPU (NVIDIA GTX 600+, AMD GCN+, Intel HD 4000+)
- Vulkan drivers installed
- Vulkan SDK (optional, for development)

**Installation:**
```bash
# Ubuntu/Debian - Install Vulkan support
sudo apt-get update
sudo apt-get install vulkan-tools vulkan-utils
sudo apt-get install mesa-vulkan-drivers  # For AMD/Intel
sudo apt-get install nvidia-driver-XXX    # For NVIDIA (replace XXX with version)

# Verify Vulkan installation
vulkaninfo --summary
vkcube  # Test Vulkan rendering

# Build with Vulkan support (automatic detection)
make build
```

**Fedora/RHEL:**
```bash
# Install Vulkan support
sudo dnf install vulkan-tools vulkan-validation-layers
sudo dnf install mesa-vulkan-drivers     # For AMD/Intel
sudo dnf install nvidia-driver           # For NVIDIA

# Build with Vulkan support
make build
```

### Linux - SYCL Support

SYCL provides unified parallel programming for CPUs, GPUs, and other accelerators.

**Requirements:**
- Intel oneAPI Toolkit or compatible SYCL implementation
- Compatible hardware (Intel GPUs, NVIDIA GPUs via CUDA backend)

**Installation:**
```bash
# Install Intel oneAPI Toolkit
wget -O- https://apt.repos.intel.com/intel-gpg-keys/GPG-PUB-KEY-INTEL-SW-PRODUCTS.PUB | gpg --dearmor | sudo tee /usr/share/keyrings/oneapi-archive-keyring.gpg > /dev/null
echo "deb [signed-by=/usr/share/keyrings/oneapi-archive-keyring.gpg] https://apt.repos.intel.com/oneapi all main" | sudo tee /etc/apt/sources.list.d/oneAPI.list
sudo apt-get update
sudo apt-get install intel-oneapi-toolkit

# Source the environment
source /opt/intel/oneapi/setvars.sh

# Verify SYCL installation
sycl-ls

# Build with SYCL support (automatic detection)
make build
```

### Linux - OpenCL Support

OpenCL provides cross-platform parallel computing support.

**Requirements:**
- OpenCL-capable GPU or CPU
- OpenCL runtime and drivers

**Installation:**
```bash
# Ubuntu/Debian - Install OpenCL support
sudo apt-get update
sudo apt-get install opencl-headers clinfo
sudo apt-get install intel-opencl-icd     # For Intel
sudo apt-get install mesa-opencl-icd      # For AMD
sudo apt-get install nvidia-opencl-dev    # For NVIDIA

# Verify OpenCL installation
clinfo

# Build with OpenCL support (automatic detection)
make build
```

## Build System GPU Detection

The Makefile implements intelligent GPU detection using the following logic:

### Detection Order (Linux/Windows)
1. **CUDA**: Checks for `nvcc` or `CUDA_PATH` environment variable
2. **HIP**: Checks for `hipconfig` or `ROCM_PATH` environment variable
3. **Vulkan**: Checks for `vulkaninfo` command or Vulkan loader
4. **OpenCL**: Checks for `clinfo` command or OpenCL runtime
5. **SYCL**: Checks for `sycl-ls` command or Intel oneAPI toolkit
6. **CPU**: Fallback when no GPU SDK is detected

### Detection Commands
```bash
# Check if GPU support is available in downloaded binaries
make detect-gpu

# Test all GPU detection logic
nvcc --version          # CUDA detection
hipconfig --version     # HIP detection  
vulkaninfo --summary    # Vulkan detection
clinfo                  # OpenCL detection
sycl-ls                 # SYCL detection
system_profiler SPDisplaysDataType | grep Metal  # Metal (macOS)
```
make test-download  # Downloads and tests appropriate binaries

# View downloaded libraries for your platform
ls ~/.cache/gollama/libs/
```

## GPU Support in Pre-Built Binaries

Gollama.cpp now uses pre-built binaries from official llama.cpp releases that include GPU support:

- **macOS**: Binaries include Metal support automatically
- **Linux**: Binaries include CUDA and HIP support where available
- **Windows**: CPU support available, GPU support planned

### Binary Selection

The downloader automatically selects GPU-enabled binaries when available:

```bash
# Downloads appropriate binary for your platform with GPU support
make download-libs
```

## Runtime Configuration

### GPU Layer Offloading

Control how many model layers are offloaded to GPU:

```go
import "github.com/dianlight/gollama.cpp"

// Configure GPU offloading
params := gollama.Context_default_params()
params.n_gpu_layers = 32  // Offload 32 layers to GPU

// For models with many layers, use -1 for all layers
params.n_gpu_layers = -1  // Offload all layers to GPU
```

### Memory Management

Configure GPU memory usage:

```go
// Set maximum GPU memory usage (in MB)
params.vram_budget = 8192  // 8GB VRAM limit

// Enable memory mapping for large models
model_params := gollama.Model_default_params()
model_params.use_mmap = true
```

### Multi-GPU Configuration

For systems with multiple GPUs:

```go
// Split model across multiple GPUs
params.split_mode = gollama.LLAMA_SPLIT_MODE_LAYER
params.main_gpu = 0      // Primary GPU device ID
params.tensor_split = []float32{0.6, 0.4}  // Split ratio between GPUs
```

## Performance Tuning

### Optimal Layer Distribution

The optimal number of GPU layers depends on:
- Available VRAM
- Model size  
- Sequence length

**Guidelines:**
- **Small models (7B)**: 32-40 layers on 8GB+ VRAM
- **Medium models (13B)**: 20-32 layers on 8GB VRAM
- **Large models (30B+)**: Adjust based on available VRAM

### Batch Size Optimization

```go
// Optimize batch size for your GPU
params.n_batch = 512     // Larger batches for high-end GPUs
params.n_ubatch = 512    // Micro-batch size for memory efficiency
```

## Troubleshooting

### Common Issues

#### CUDA Not Detected
```bash
# Check CUDA installation
nvcc --version
ls -la /usr/local/cuda/bin/nvcc

# Check environment variables
echo $CUDA_PATH
echo $LD_LIBRARY_PATH
```

#### HIP Not Detected
```bash
# Check ROCm installation
/opt/rocm/bin/hipconfig --platform
ls -la /opt/rocm/bin/

# Check environment variables
echo $ROCM_PATH
echo $HIP_PATH
```

#### GPU Memory Errors
```go
// Reduce GPU memory usage
params.n_gpu_layers = 16    // Reduce from 32
params.vram_budget = 4096   // Reduce VRAM limit
```

#### Performance Issues
```go
// Optimize for your hardware
params.n_threads = 8           // Match CPU cores
params.n_threads_batch = 8     // Batch processing threads
params.rope_scaling_type = gollama.LLAMA_ROPE_SCALING_TYPE_LINEAR
```

### Debug Information

Enable detailed GPU information during build:

```bash
# Verbose GPU detection
make build V=1

# Check library GPU backend
ldd libs/linux_amd64/libllama.so | grep -E "(cuda|hip)"
```

### Verification

Test GPU acceleration is working:

```go
package main

import (
    "fmt"
    "github.com/dianlight/gollama.cpp"
)

func main() {
    // Load model with GPU acceleration
    model_params := gollama.Model_default_params()
    model := gollama.Load_model_from_file("model.gguf", model_params)
    defer gollama.Free_model(model)
    
    // Create context with GPU layers
    ctx_params := gollama.Context_default_params()
    ctx_params.n_gpu_layers = 32
    
    ctx := gollama.New_context_with_model(model, ctx_params)
    defer gollama.Free(ctx)
    
    // Check if GPU is being used
    fmt.Printf("GPU layers: %d\n", ctx_params.n_gpu_layers)
    
    // Monitor GPU usage with nvidia-smi or rocm-smi during inference
}
```

Monitor GPU utilization:
```bash
# NVIDIA GPUs
watch -n 1 nvidia-smi

# AMD GPUs  
watch -n 1 rocm-smi

# Check GPU memory usage during inference
```

## Best Practices

1. **Start Conservative**: Begin with fewer GPU layers and increase gradually
2. **Monitor Memory**: Watch VRAM usage to avoid out-of-memory errors
3. **Profile Performance**: Test different configurations for your specific use case
4. **Update Drivers**: Keep GPU drivers updated for best performance
5. **Check Compatibility**: Verify your GPU is supported by the chosen backend

## Support Matrix

### Tested Configurations

| Platform | GPU | Backend | Model Sizes | Status |
|----------|-----|---------|-------------|--------|
| macOS M1/M2 | Apple Silicon | Metal | 7B-70B | ✅ Verified |
| Ubuntu 22.04 | RTX 4090 | CUDA 12.0 | 7B-70B | ✅ Verified |
| Ubuntu 22.04 | RX 7900 XTX | ROCm 5.7 | 7B-30B | ✅ Verified |
| Windows 11 | RTX 3080 | CUDA 11.8 | 7B-30B | ✅ Verified |
| Fedora 38 | RTX 3070 | CUDA 12.1 | 7B-13B | ✅ Verified |

For the latest compatibility information, see our [CI test matrix](../.github/workflows/ci.yml).
