package main

import (
	"fmt"
	"log"

	"github.com/dianlight/gollama.cpp"
)

func main() {
	fmt.Println("GGML Low-Level API Demo")
	fmt.Println("========================\n")

	// Initialize the library
	if err := gollama.Backend_init(); err != nil {
		log.Fatal(err)
	}
	defer gollama.Backend_free()

	// Demonstrate type information queries
	demonstrateTypeInfo()

	// Demonstrate backend device enumeration
	demonstrateBackendDevices()
}

func demonstrateTypeInfo() {
	fmt.Println("=== GGML Type Information ===")

	types := []gollama.GgmlType{
		gollama.GGML_TYPE_F32,
		gollama.GGML_TYPE_F16,
		gollama.GGML_TYPE_BF16,
		gollama.GGML_TYPE_Q4_0,
		gollama.GGML_TYPE_Q8_0,
		gollama.GGML_TYPE_IQ2_XXS,
		gollama.GGML_TYPE_IQ4_XS,
		gollama.GGML_TYPE_I32,
	}

	fmt.Printf("%-12s | %-10s | %-10s | %-10s\n", "Type", "Size", "Quantized", "Name")
	fmt.Println("-------------|------------|------------|------------")

	for _, typ := range types {
		// Get type size
		size, sizeErr := gollama.Ggml_type_size(typ)
		sizeStr := "N/A"
		if sizeErr == nil {
			sizeStr = fmt.Sprintf("%d bytes", size)
		}

		// Check if quantized
		isQuant, quantErr := gollama.Ggml_type_is_quantized(typ)
		quantStr := "N/A"
		if quantErr == nil {
			quantStr = fmt.Sprintf("%v", isQuant)
		}

		// Get type name
		name, nameErr := gollama.Ggml_type_name(typ)
		if nameErr != nil {
			name = "N/A"
		}

		fmt.Printf("%-12s | %-10s | %-10s | %-10s\n",
			typ.String(), sizeStr, quantStr, name)
	}

	fmt.Println()
}

func demonstrateBackendDevices() {
	fmt.Println("=== Backend Devices ===")

	count, err := gollama.Ggml_backend_dev_count()
	if err != nil {
		fmt.Println("Backend device enumeration not available in this build")
		fmt.Println("(GGML functions may not be exported)")
		return
	}

	if count == 0 {
		fmt.Println("No backend devices available")
		return
	}

	fmt.Printf("Found %d backend device(s):\n\n", count)

	for i := uint64(0); i < count; i++ {
		dev, err := gollama.Ggml_backend_dev_get(i)
		if err != nil {
			fmt.Printf("Device %d: Error getting device\n", i)
			continue
		}

		name, err := gollama.Ggml_backend_dev_name(dev)
		if err != nil {
			fmt.Printf("Device %d: Error getting name\n", i)
			continue
		}

		fmt.Printf("Device %d: %s\n", i, name)

		// Try to get description
		desc, err := gollama.Ggml_backend_dev_description(dev)
		if err == nil && desc != "" {
			fmt.Printf("  Description: %s\n", desc)
		}

		// Try to get memory info (may not be supported on all devices)
		free, total, err := gollama.Ggml_backend_dev_memory(dev)
		if err == nil && total > 0 {
			fmt.Printf("  Memory: %.2f MB free / %.2f MB total (%.1f%% used)\n",
				float64(free)/(1024*1024),
				float64(total)/(1024*1024),
				float64(total-free)/float64(total)*100)
		}

		fmt.Println()
	}
}
