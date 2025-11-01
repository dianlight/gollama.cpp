# Gollama.cpp Roadmap

This document outlines the planned development path for gollama.cpp, prioritized by user impact and technical feasibility.

## Current Status (v0.2.1 - August 2025)

### ✅ Completed Major Features
- **Download-based Architecture**: Eliminated build dependencies, uses official llama.cpp binaries
- **Cross-platform Support**: Full support for macOS/Linux, Windows build compatibility
- **GPU Backend Support**: Metal, CUDA, HIP, Vulkan, OpenCL, SYCL with automatic detection
- **Parallel Downloads**: Concurrent library downloads with checksum verification
- **Platform-specific Loading**: Native Windows syscalls, purego for Unix-like systems
- **Comprehensive Examples**: 9 example implementations covering various use cases
- **CI/CD Pipeline**: Cross-platform testing, automated releases, documentation generation

### 🚧 In Progress
- **Renovate Integration**: Automated llama.cpp dependency updates (Phase 5)

### ✅ Recently Completed
- **Windows Runtime Support**: Implemented proper function registration using GetProcAddress for Windows FFI support

## Short-term Goals (Q3-Q4 2025)

### Priority 1: Windows Runtime Completion
**Target: September 2025**
- [ ] Windows installation and setup documentation

**Technical Details:**
- Create Windows-specific test suites
- Optimize memory management for Windows platform

### Priority 2: Enhanced GPU Support
**Target: October 2025**
- [ ] Intelligent GPU variant selection based on hardware detection
- [ ] GPU memory management and optimization features
- [ ] Multi-GPU support and load balancing
- [ ] GPU performance benchmarking and monitoring
- [ ] Advanced GPU configuration options

**New Features:**
```go
// Future GPU configuration API
gollama.LoadLibraryWithOptions(&gollama.LoadOptions{
    GPUVariant:    "auto",        // auto, cuda, hip, vulkan, opencl, sycl, cpu
    GPUDevices:    []int{0, 1},   // Multi-GPU support
    GPUMemory:     8192,          // Memory limit in MB
    PreferGPU:     true,          // GPU preference over CPU
})
```

### Priority 3: Advanced Model Management
**Target: November 2025**
- [ ] Model registry and versioning system
- [ ] Automatic model downloading and caching

**New Components:**
- `ModelRegistry` for centralized model management
- `ModelDownloader` with progress tracking
- Integration with Hugging Face Hub

## Medium-term Goals (Q1-Q2 2026)

### Priority 1: Performance Optimization
**Target: Q1 2026**
- [ ] Memory mapping optimizations
- [ ] Zero-copy operations where possible
- [ ] Batch processing improvements
- [ ] Streaming optimizations
- [ ] Memory usage reduction

### Priority 2: Developer Experience Enhancements
**Target: Q1 2026**
- [ ] VS Code extension for gollama.cpp development
- [ ] Interactive debugging tools
- [ ] Performance profiling dashboard
- [ ] Comprehensive benchmarking suite
- [ ] Developer documentation portal

### Priority 3: Advanced Features
**Target: Q2 2026**
- [ ] Advanced context management
- [ ] Real-time model switching

**Plugin System Design:**
```go
// Future plugin architecture
type Plugin interface {
    Name() string
    Initialize(ctx *Context) error
    ProcessTokens(tokens []Token) []Token
    Cleanup() error
}

// Plugin registration
gollama.RegisterPlugin(&CustomSamplingPlugin{})
```

## Long-term Vision (wait for purego struct support)

### Features Requiring purego Struct Support
These features depend on purego's ability to handle struct parameters and return values, which is currently only available on Darwin (macOS) platforms. They will be moved to active development once purego gains cross-platform struct support.

#### Core Runtime Functions
- [ ] Complete Windows runtime library loading implementation - *Requires struct parameter/return support*
- [ ] Windows GPU acceleration support (CUDA, HIP, Vulkan, OpenCL, SYCL) - *Requires struct parameter/return support*
- [ ] Windows-specific examples and testing - *Requires core runtime functions*
- [ ] Performance optimization for Windows platform - *Requires struct parameter/return support*

#### Batch Processing and Context Management
- [ ] Cross-platform batch processing - *Requires LlamaBatch struct support*
- [ ] Advanced context management - *Requires LlamaContextParams struct support*
- [ ] Model parameter configuration - *Requires LlamaModelParams struct support*

#### Sampling and Generation
- [ ] Advanced sampling configurations - *Requires sampler struct support*
- [ ] Custom sampling strategies - *Requires struct-based sampling API*

**Technical Dependencies:**
- purego struct parameter support on Windows/Linux
- purego struct return value support on Windows/Linux
- Cross-platform function registration for struct-based APIs

## Long-term Vision (wait for llama.cpp)

### Features Requiring llama.cpp Function Implementation
These features depend on specific llama.cpp functions that are either missing or not fully implemented in the current version. They will be moved to active development once the required functions become available.

#### Model Processing and Conversion
- [ ] Model quantization utilities - *Requires full quantization API*
- [ ] Model format conversion tools - *Requires conversion functions*
- [ ] Model performance profiling - *Requires timing/profiling functions*

#### Advanced Sampling and Generation  
- [ ] Plugin architecture for extensibility - *Requires callback/plugin API*
- [ ] Custom sampling strategies - *Requires extended sampling functions*
- [ ] Distributed inference support - *Requires distributed computing API*

#### Performance and Monitoring
- [ ] Performance profiling dashboard - *Requires timing functions*
- [ ] Advanced GPU configuration options - *Requires extended GPU API*

#### Multi-modal and Advanced AI
- [ ] Multi-modal support (text + images) - *Requires multi-modal API*
- [ ] Voice integration capabilities - *Requires audio processing API*
- [ ] Real-time translation services - *Requires translation functions*
- [ ] Advanced reasoning frameworks - *Requires reasoning API*
- [ ] Knowledge graph integration - *Requires graph processing API*

## Long-term Vision (2026+)

### Advanced AI Integration

### Enterprise Features
- [ ] High-availability deployment patterns
- [ ] Monitoring and observability tools
- [ ] Security and compliance features
- [ ] Enterprise-grade authentication
- [ ] Usage analytics and reporting

### Ecosystem Development
- [ ] Community plugin marketplace
- [ ] Third-party integrations (databases, APIs)
- [ ] Cloud-native deployment options
- [ ] Containerization and orchestration
- [ ] Edge computing optimizations

## Implementation Priorities

### High Priority (Critical Path)
1. **Automated Dependency Management** - Essential for maintenance
2. **Enhanced GPU Support** - Core value proposition
3. **Advanced Model Management** - Improves user experience

### Medium Priority (Value-Adding)
1. **Performance Optimizations** - Competitive advantage
2. **Developer Tools** - Community growth

### Low Priority (Future Enhancement)
1. **Enterprise Features** - Commercial applications
2. **Community Features** - Plugin marketplace and ecosystem

### Waiting for purego struct support (Blocked)
1. **Windows Runtime Completion** - Requires struct parameter/return support
2. **Cross-platform Batch Processing** - Requires LlamaBatch struct support
3. **Advanced Sampling** - Requires struct-based sampling configurations

### Waiting for llama.cpp (Blocked)
1. **Model Processing Tools** - Quantization, conversion, profiling utilities
2. **Multi-modal Support** - Image, voice, and advanced AI features

## Technical Dependencies

### External Dependencies
- **llama.cpp releases**: Continued compatibility with upstream changes
- **GPU vendors**: SDK availability and compatibility
- **Go language**: Version compatibility and feature availability
- **Platform vendors**: OS-specific API stability

### Internal Dependencies
- **Build system**: Makefile and CI/CD pipeline maintenance
- **Testing infrastructure**: Comprehensive cross-platform testing
- **Documentation**: Keeping pace with feature development
- **Community**: Feedback and contribution management

## Success Metrics

### Short-term (2025)
- [ ] 100% Windows runtime compatibility
- [ ] Zero manual compilation required for any platform
- [ ] <60 second setup time for new users
- [ ] 95% GPU detection accuracy

### Medium-term (2026)
- [ ] 10+ community contributors
- [ ] 50+ stars on GitHub
- [ ] Integration in 5+ downstream projects
- [ ] Performance within 10% of native llama.cpp

### Long-term (2027+)
- [ ] Industry-standard Go LLM binding
- [ ] Enterprise adoption
- [ ] Plugin ecosystem with 10+ plugins

## Community & Contributions

### Current Contributors
- Core maintainer: [@dianlight](https://github.com/dianlight)
- Contributors welcome for all roadmap items

### How to Contribute
1. **Windows Runtime**: Help with Windows-specific development and testing
2. **GPU Support**: Contribute GPU backend implementations and optimizations
3. **Examples**: Create new use-case examples and tutorials
4. **Documentation**: Improve guides and API documentation
5. **Testing**: Cross-platform testing and bug reports

### Contribution Priorities
1. **Windows Platform**: Critical for project completion
2. **GPU Optimization**: High-value technical contributions
3. **Documentation**: Essential for user adoption
4. **Examples**: Showcase project capabilities

## Risk Management

### Technical Risks
- **llama.cpp API changes**: Mitigated by automated testing and version tracking
- **Platform deprecation**: Diversified platform support reduces impact
- **GPU driver issues**: Fallback to CPU ensures functionality
- **Memory constraints**: Progressive loading and optimization features

### Project Risks
- **Maintainer availability**: Seeking additional core contributors
- **Community adoption**: Focus on documentation and examples
- **Competition**: Differentiate through ease of use and performance
- **Scope creep**: Maintain focus on core value proposition

## Getting Involved

### For Users
- Try the examples and provide feedback
- Report platform-specific issues
- Request missing features through GitHub issues
- Share your use cases and success stories

### For Developers
- Review the [CONTRIBUTING.md](../CONTRIBUTING.md) guide
- Check the [GitHub Issues](https://github.com/dianlight/gollama.cpp/issues) for `enhancement` labels
- Join discussions in GitHub Discussions
- Submit pull requests for roadmap items

### For Organizations
- Evaluate gollama.cpp for your AI/ML projects
- Provide feedback on enterprise requirements
- Consider sponsoring development of specific features
- Share performance benchmarks and use cases

---

**Last Updated**: August 06, 2025
**Next Review**: September 1, 2025

*This roadmap is a living document and will be updated based on community feedback, technical discoveries, and changing requirements.*
