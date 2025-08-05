package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dianlight/gollama.cpp"
)

func main() {
	fmt.Println("=== gollama.cpp Parallel Download & Checksum Demo ===\n")

	// Demo 1: Download libraries for specific platforms
	fmt.Println("1. Downloading libraries for specific platforms...")
	platforms := []string{"linux/amd64", "darwin/arm64", "windows/amd64"}

	results, err := gollama.DownloadLibrariesForPlatforms(platforms, "")
	if err != nil {
		log.Fatalf("Failed to download libraries: %v", err)
	}

	fmt.Printf("\nDownload Results:\n")
	fmt.Printf("================\n")
	for _, result := range results {
		if result.Success {
			fmt.Printf("✅ %s: SUCCESS\n", result.Platform)
			if result.LibraryPath != "" {
				fmt.Printf("   Library: %s\n", result.LibraryPath)
			}
			if result.SHA256Sum != "" {
				fmt.Printf("   SHA256: %s\n", result.SHA256Sum)
			}
		} else {
			fmt.Printf("❌ %s: FAILED", result.Platform)
			if result.Error != nil {
				fmt.Printf(" - %s", result.Error.Error())
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Demo 2: Calculate checksum of a downloaded library
	fmt.Println("2. Calculating checksum of downloaded library...")
	for _, result := range results {
		if result.Success && result.LibraryPath != "" {
			if _, err := os.Stat(result.LibraryPath); err == nil {
				checksum, err := gollama.GetSHA256ForFile(result.LibraryPath)
				if err != nil {
					fmt.Printf("Failed to calculate checksum for %s: %v\n", result.LibraryPath, err)
				} else {
					fmt.Printf("File: %s\n", result.LibraryPath)
					fmt.Printf("SHA256: %s\n", checksum)
				}
				break
			}
		}
	}

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("Features demonstrated:")
	fmt.Println("✅ Parallel downloads for multiple platforms")
	fmt.Println("✅ Automatic SHA256 checksum calculation")
	fmt.Println("✅ Platform-specific library detection")
	fmt.Println("✅ Concurrent download processing")
	fmt.Println("✅ Error handling and reporting")
}
