package gollama

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// PopulateLibDirectoryFromResults copies downloaded library artifacts into the local libs directory so they
// can be embedded in future builds. Only the llama.cpp build defined by LlamaCppBuild is supported.
func PopulateLibDirectoryFromResults(results []DownloadResult, version, libsDir string) error {
	effectiveVersion := version
	if effectiveVersion == "" {
		effectiveVersion = LlamaCppBuild
	}

	if effectiveVersion != LlamaCppBuild {
		return fmt.Errorf("only llama.cpp build %s can be embedded (requested %s)", LlamaCppBuild, effectiveVersion)
	}

	if libsDir == "" {
		libsDir = "libs"
	}

	if err := os.MkdirAll(libsDir, 0o750); err != nil {
		return fmt.Errorf("failed to ensure libs directory: %w", err)
	}

	// Clean up any old versions to enforce single-version policy.
	if err := pruneLegacyLibVersions(libsDir, effectiveVersion); err != nil {
		return err
	}

	for _, res := range results {
		if !res.Success {
			continue
		}

		goos, goarch, err := splitPlatform(res.Platform)
		if err != nil {
			return err
		}

		srcDir := res.ExtractedDir
		if srcDir == "" && res.LibraryPath != "" {
			srcDir = filepath.Dir(res.LibraryPath)
		}
		if srcDir == "" {
			return fmt.Errorf("could not determine source directory for platform %s", res.Platform)
		}

		if err := copyPlatformLibraries(srcDir, libsDir, goos, goarch, effectiveVersion); err != nil {
			return err
		}
	}

	return nil
}

func pruneLegacyLibVersions(libsDir, version string) error {
	entries, err := os.ReadDir(libsDir)
	if errors.Is(err, fs.ErrNotExist) {
		return os.MkdirAll(libsDir, 0o750)
	}
	if err != nil {
		return fmt.Errorf("failed to read libs directory: %w", err)
	}

	suffix := "_" + version

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, suffix) {
			continue
		}

		if err := os.RemoveAll(filepath.Join(libsDir, name)); err != nil {
			return fmt.Errorf("failed to remove legacy libs directory %s: %w", name, err)
		}
	}

	return nil
}

func copyPlatformLibraries(srcDir, libsDir, goos, goarch, version string) error {
	targetDir := filepath.Join(libsDir, fmt.Sprintf("%s_%s_%s", goos, goarch, version))

	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("failed to clean target directory %s: %w", targetDir, err)
	}
	if err := os.MkdirAll(targetDir, 0o750); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	var copied bool
	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		lower := strings.ToLower(d.Name())
		switch {
		case strings.HasSuffix(lower, ".dylib"), strings.HasSuffix(lower, ".so"), strings.HasSuffix(lower, ".dll"):
		default:
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read library %s: %w", path, err)
		}

		destPath := filepath.Join(targetDir, d.Name())
		if err := os.WriteFile(destPath, data, 0o600); err != nil {
			return fmt.Errorf("failed to write library %s: %w", destPath, err)
		}
		copied = true
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to copy libraries from %s: %w", srcDir, err)
	}

	if !copied {
		return fmt.Errorf("no libraries found in %s for %s/%s", srcDir, goos, goarch)
	}

	return nil
}

func splitPlatform(platform string) (string, string, error) {
	parts := strings.Split(platform, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid platform string: %s", platform)
	}
	return parts[0], parts[1], nil
}
