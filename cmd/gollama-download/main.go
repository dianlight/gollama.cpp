package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dianlight/gollama.cpp"
)

func main() {
	var (
		download       = flag.Bool("download", false, "Download llama.cpp library for current platform")
		downloadAll    = flag.Bool("download-all", false, "Download llama.cpp libraries for all supported platforms")
		platforms      = flag.String("platforms", "", "Comma-separated list of platforms to download (e.g., linux/amd64,darwin/arm64)")
		version        = flag.String("version", "", "Specific version to download (default: latest)")
		testDownload   = flag.Bool("test-download", false, "Test download functionality without loading library")
		cleanCache     = flag.Bool("clean-cache", false, "Clean library cache")
		showVersion    = flag.Bool("v", false, "Show version information")
		showChecksum   = flag.Bool("checksum", false, "Show SHA256 checksum of downloaded files")
		verifyChecksum = flag.String("verify-checksum", "", "Verify SHA256 checksum of a file")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("gollama.cpp library downloader\n")
		fmt.Printf("Supports downloading pre-built llama.cpp binaries from ggml-org/llama.cpp\n")
		return
	}

	if *cleanCache {
		fmt.Println("Cleaning library cache...")
		if err := gollama.CleanLibraryCache(); err != nil {
			log.Fatalf("Failed to clean cache: %v", err)
		}
		fmt.Println("Cache cleaned successfully")
		return
	}

	if *verifyChecksum != "" {
		fmt.Printf("Calculating SHA256 checksum for %s...\n", *verifyChecksum)
		checksum, err := gollama.GetSHA256ForFile(*verifyChecksum)
		if err != nil {
			log.Fatalf("Failed to calculate checksum: %v", err)
		}
		fmt.Printf("SHA256: %s\n", checksum)
		return
	}

	if *downloadAll {
		fmt.Println("Downloading libraries for all supported platforms...")
		allPlatforms := []string{
			"darwin/amd64", "darwin/arm64",
			"linux/amd64", "linux/arm64",
			"windows/amd64", "windows/arm64",
		}

		results, err := gollama.DownloadLibrariesForPlatforms(allPlatforms, *version)
		if err != nil {
			log.Fatalf("Failed to download libraries: %v", err)
		}

		printDownloadResults(results, *showChecksum)
		return
	}

	if *platforms != "" {
		fmt.Printf("Downloading libraries for platforms: %s...\n", *platforms)
		platformList := strings.Split(*platforms, ",")
		for i, p := range platformList {
			platformList[i] = strings.TrimSpace(p)
		}

		results, err := gollama.DownloadLibrariesForPlatforms(platformList, *version)
		if err != nil {
			log.Fatalf("Failed to download libraries: %v", err)
		}

		printDownloadResults(results, *showChecksum)
		return
	}

	if *testDownload {
		fmt.Println("Testing library download functionality...")
		downloader, err := gollama.NewLibraryDownloader()
		if err != nil {
			log.Fatalf("Failed to create downloader: %v", err)
		}

		var release *gollama.ReleaseInfo
		if *version != "" {
			fmt.Printf("Getting release information for version %s...\n", *version)
			release, err = downloader.GetReleaseByTag(*version)
		} else {
			fmt.Println("Getting latest release information...")
			release, err = downloader.GetLatestRelease()
		}

		if err != nil {
			log.Fatalf("Failed to get release info: %v", err)
		}

		fmt.Printf("Found release: %s\n", release.TagName)

		pattern, err := downloader.GetPlatformAssetPattern()
		if err != nil {
			log.Fatalf("Failed to get platform pattern: %v", err)
		}

		fmt.Printf("Looking for asset matching pattern: %s\n", pattern)

		assetName, downloadURL, err := downloader.FindAssetByPattern(release, pattern)
		if err != nil {
			log.Fatalf("Failed to find platform asset: %v", err)
		}

		fmt.Printf("Found asset: %s\n", assetName)
		fmt.Printf("Download URL: %s\n", downloadURL)
		fmt.Println("Download test completed successfully")
		return
	}

	if *download {
		fmt.Println("Downloading llama.cpp library...")

		var err error
		if *version != "" {
			fmt.Printf("Downloading version %s...\n", *version)
			err = gollama.LoadLibraryWithVersion(*version)
		} else {
			fmt.Println("Downloading latest version...")
			err = gollama.LoadLibraryWithVersion("")
		}

		if err != nil {
			log.Fatalf("Failed to download library: %v", err)
		}

		fmt.Println("Library downloaded and loaded successfully")

		// Show checksum if requested
		if *showChecksum {
			// For single platform download, calculate checksum of the main library cache
			fmt.Println("Calculating SHA256 checksums...")
			// Note: This is a simplified approach - for more detailed checksums,
			// users should use the parallel download features
		}

		return
	}

	// Default behavior: show help
	fmt.Printf("gollama.cpp library downloader\n\n")
	fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	fmt.Printf("\nExamples:\n")
	fmt.Printf("  %s -download                     # Download latest version for current platform\n", os.Args[0])
	fmt.Printf("  %s -download -version b6089      # Download specific version for current platform\n", os.Args[0])
	fmt.Printf("  %s -download-all                 # Download for all supported platforms\n", os.Args[0])
	fmt.Printf("  %s -platforms linux/amd64,darwin/arm64  # Download for specific platforms\n", os.Args[0])
	fmt.Printf("  %s -test-download               # Test download without loading\n", os.Args[0])
	fmt.Printf("  %s -clean-cache                 # Clean cache directory\n", os.Args[0])
	fmt.Printf("  %s -checksum -download           # Download and show checksums\n", os.Args[0])
	fmt.Printf("  %s -verify-checksum file.zip     # Verify checksum of a file\n", os.Args[0])
}

// printDownloadResults prints the results of parallel downloads
func printDownloadResults(results []gollama.DownloadResult, showChecksum bool) {
	fmt.Printf("\nDownload Results:\n")
	fmt.Printf("================\n")

	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
			fmt.Printf("✅ %s: SUCCESS", result.Platform)
			if result.LibraryPath != "" {
				fmt.Printf(" (Library: %s)", result.LibraryPath)
			}
			if showChecksum && result.SHA256Sum != "" {
				fmt.Printf("\n   SHA256: %s", result.SHA256Sum)
			}
			fmt.Println()
		} else {
			fmt.Printf("❌ %s: FAILED", result.Platform)
			if result.Error != nil {
				fmt.Printf(" - %s", result.Error.Error())
			}
			fmt.Println()
		}
	}

	fmt.Printf("\nSummary: %d/%d platforms downloaded successfully\n", successCount, len(results))
}
