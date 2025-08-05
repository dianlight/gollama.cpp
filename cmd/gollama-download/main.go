package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dianlight/gollama.cpp"
)

func main() {
	var (
		download     = flag.Bool("download", false, "Download llama.cpp library for current platform")
		version      = flag.String("version", "", "Specific version to download (default: latest)")
		testDownload = flag.Bool("test-download", false, "Test download functionality without loading library")
		cleanCache   = flag.Bool("clean-cache", false, "Clean library cache")
		showVersion  = flag.Bool("v", false, "Show version information")
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
		return
	}

	// Default behavior: show help
	fmt.Printf("gollama.cpp library downloader\n\n")
	fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	fmt.Printf("\nExamples:\n")
	fmt.Printf("  %s -download                     # Download latest version\n", os.Args[0])
	fmt.Printf("  %s -download -version b6089      # Download specific version\n", os.Args[0])
	fmt.Printf("  %s -test-download               # Test download without loading\n", os.Args[0])
	fmt.Printf("  %s -clean-cache                 # Clean cache directory\n", os.Args[0])
}
