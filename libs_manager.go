package gollama

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
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

// MergeVariantLibraries merges library files from multiple variant directories into the
// single target libs directory for the given platform and version. If two variants contain
// a library with the same file name but different content, it returns an error to prevent
// ambiguous or conflicting embeddings.
//
// Target layout: <libsDir>/<goos>_<goarch>_<version>/
// Copied files: *.dylib, *.so, *.dll (base name only; subdir structure is not preserved)
func MergeVariantLibraries(goos, goarch, version, libsDir string, variantDirs []string) error {
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

	targetDir := filepath.Join(libsDir, fmt.Sprintf("%s_%s_%s", goos, goarch, effectiveVersion))

	// Start clean for target directory
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("failed to clean target directory %s: %w", targetDir, err)
	}
	if err := os.MkdirAll(targetDir, 0o750); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	// Walk each variant directory and merge libraries
	var copied bool
	for _, srcDir := range variantDirs {
		if srcDir == "" {
			continue
		}
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

			destPath := filepath.Join(targetDir, d.Name())

			// If destination exists, compare contents; if equal, skip; if different, error.
			if _, err := os.Stat(destPath); err == nil {
				same, cmpErr := filesHaveSameSHA256(path, destPath)
				if cmpErr != nil {
					return fmt.Errorf("failed to compare %s and %s: %w", path, destPath, cmpErr)
				}
				if !same {
					return fmt.Errorf("conflicting library file detected: %s (different contents across variants)", d.Name())
				}
				// identical; skip copy
				return nil
			}

			// Copy file
			in, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open %s: %w", path, err)
			}
			defer func() { _ = in.Close() }()

			out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
			if err != nil {
				return fmt.Errorf("failed to create %s: %w", destPath, err)
			}
			if _, err := io.Copy(out, in); err != nil {
				_ = out.Close()
				return fmt.Errorf("failed to write %s: %w", destPath, err)
			}
			if err := out.Close(); err != nil {
				return fmt.Errorf("failed to close %s: %w", destPath, err)
			}
			copied = true
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to merge libraries from %s: %w", srcDir, err)
		}
	}

	if !copied {
		return fmt.Errorf("no libraries found to copy for %s/%s", goos, goarch)
	}

	return nil
}

func filesHaveSameSHA256(p1, p2 string) (bool, error) {
	h1, err := sha256ForFile(p1)
	if err != nil {
		return false, err
	}
	h2, err := sha256ForFile(p2)
	if err != nil {
		return false, err
	}
	return h1 == h2, nil
}

func sha256ForFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
