# Purego Struct Support Roadmap Update Summary

This document summarizes the updates made to reorganize the roadmap by moving goals that depend on purego struct support limitations to a new "Long-term Vision (wait for purego struct support)" section.

## Understanding Purego Struct Limitations

### Current State
The project uses [purego](https://github.com/ebitengine/purego) v0.9.0-alpha.10 for cross-platform C interoperability without CGO. However, purego currently has limitations with struct parameters and return values on non-Darwin platforms.

### What Works vs What Doesn't

#### ✅ **Currently Supported (All Platforms)**
- Simple function calls with primitive parameters
- Pointer parameters and return values
- Basic C function registration
- Library loading and symbol resolution

#### ❌ **Limited to Darwin (macOS) Only**
- Functions that take struct parameters
- Functions that return structs
- Complex type marshaling
- Batch processing with LlamaBatch structs
- Context initialization with LlamaContextParams
- Model loading with LlamaModelParams

### Technical Impact

This limitation significantly affects Windows and Linux support because:

1. **Core Runtime Functions**: Model loading, context creation, and batch processing all require struct support
2. **Cross-platform Compatibility**: Many essential features only work on macOS
3. **Feature Parity**: Windows/Linux users have limited functionality compared to macOS users

## Roadmap Reorganization

### Goals Moved to "Long-term Vision (wait for purego struct support)"

#### 🏗️ **Core Runtime Functions**
- **Complete Windows runtime library loading implementation** - Requires struct parameter/return support
- **Windows GPU acceleration support** - Requires struct parameter/return support  
- **Windows-specific examples and testing** - Requires core runtime functions
- **Performance optimization for Windows platform** - Requires struct parameter/return support

#### 📦 **Batch Processing and Context Management**
- **Cross-platform batch processing** - Requires LlamaBatch struct support
- **Advanced context management** - Requires LlamaContextParams struct support
- **Model parameter configuration** - Requires LlamaModelParams struct support

#### 🎯 **Sampling and Generation**
- **Advanced sampling configurations** - Requires sampler struct support
- **Custom sampling strategies** - Requires struct-based sampling API

### Priority Reorganization

#### **Before** (Original Priorities):
1. ~~Windows Runtime Completion~~ (moved to blocked)
2. Enhanced GPU Support
3. Advanced Model Management

#### **After** (Updated Priorities):
1. **Automated Dependency Management** - Can be implemented without struct support
2. **Enhanced GPU Support** - Core value proposition that works on macOS
3. **Advanced Model Management** - Improves user experience on supported platforms

## Updated Infrastructure

### 🔧 **Enhanced Scripts**

#### New `roadmap-update.sh` Commands:
```bash
# Scan for purego struct limitations
./scripts/roadmap-update.sh scan-purego

# Updated validation with purego section
./scripts/roadmap-update.sh validate
```

#### Detection Patterns:
The script now scans for these purego limitation indicators:
- `Skip.*struct.*non-Darwin`
- `struct.*return.*not.*supported`  
- `Helper functions.*struct returns`
- `runtime.GOOS.*darwin.*struct`
- `non-Darwin platforms.*struct`
- `Return default.*non-Darwin`

### 📋 **Enhanced Makefile**

#### New Rules:
```makefile
# Scan for purego struct support dependencies
roadmap-scan-purego:
	@echo "Scanning for code that depends on purego struct support"
	@bash scripts/roadmap-update.sh scan-purego
```

### 🚀 **Updated CI Pipeline**

#### Enhanced `doc-sync-check.yml`:
```yaml
- name: Validate ROADMAP.md format
  run: |
    # Check for purego struct support dependencies
    echo "Checking for potential purego struct support dependencies..."
    make roadmap-scan-purego
```

### 💻 **Source Code Improvements**

#### Enhanced Comments and Error Messages:
```go
// Model functions - Skip functions that use structs on non-Darwin platforms - moved to ROADMAP "wait for purego struct support" section
if runtime.GOOS == "darwin" {

// Return default values for non-Darwin platforms - blocks ROADMAP "wait for purego struct support"
return LlamaModelParams{

// Error messages with roadmap context
return errors.New("Model_load_from_file not yet implemented for non-Darwin platforms - blocks ROADMAP Priority 1 (wait for purego struct support)")
```

## Current Detection Results

The scanning functionality successfully detected **10+ struct-related limitations**:

### Key Findings:
- **5 function registration blocks** - Functions only registered on Darwin
- **3 helper function fallbacks** - Default structs returned on non-Darwin
- **2+ batch processing limitations** - LlamaBatch functions unavailable
- **Multiple error conditions** - Clear indication of struct dependency

### Example Detection Output:
```
Found Darwin-only function registrations (likely due to struct limitations):
- gollama.go:565: Model functions (Darwin only)
- gollama.go:612: Batch functions (Darwin only) 
- gollama.go:619: Decode functions (Darwin only)
- gollama.go:635: Sampling functions (Darwin only)
```

## Benefits of This Reorganization

### 🎯 **Realistic Timeline Management**
- **Clear Separation**: Distinguishes between achievable goals and blocked features
- **Dependency Tracking**: Explicitly tracks external library limitations
- **Priority Clarity**: Focuses development effort on implementable features

### 🔍 **Improved Transparency**  
- **User Expectations**: Clear communication about platform limitations
- **Developer Guidance**: Obvious indication of what requires upstream changes
- **Community Awareness**: Understanding of project constraints

### ⚡ **Development Focus**
- **Achievable Goals**: Prioritizes features that can be implemented now
- **Resource Allocation**: Avoids wasted effort on blocked functionality
- **Strategic Planning**: Aligns development with external library roadmaps

## Future Path Forward

### 📋 **Monitoring purego Development**
1. **Track purego Releases**: Monitor struct support progress
2. **Test Compatibility**: Validate new versions against gollama.cpp needs
3. **Gradual Migration**: Move features back to active development as support becomes available

### 🔧 **Alternative Approaches**
While waiting for purego struct support, potential alternatives include:
1. **CGO Fallback**: Optional CGO support for struct-heavy operations
2. **Wrapper Functions**: C wrapper functions that decompose structs into primitives
3. **Manual Marshaling**: Custom struct serialization/deserialization

### 📈 **Success Metrics**
- **purego struct support**: Track upstream library progress
- **Cross-platform parity**: Measure feature availability across platforms
- **User adoption**: Monitor usage on Windows/Linux vs macOS

## Testing the Updates

### ✅ **Validation Results**
```bash
# Roadmap structure validation
make roadmap-validate
# ✅ All required sections found
# ℹ️ Found 63 planned items, 0 completed items
# ℹ️ Blocked items: 10 (llama.cpp), 10 (purego struct support)

# Purego limitation detection  
make roadmap-scan-purego
# ⚠️ Found Darwin-only function registrations (likely due to struct limitations)
```

This comprehensive reorganization provides a realistic, well-documented approach to managing the project roadmap while acknowledging and properly categorizing external dependency limitations.

---

**Impact Summary**: This update moves 4+ major goals from active development to the "wait for purego struct support" section, providing clarity about what can realistically be achieved with current tooling versus what requires upstream library improvements.
