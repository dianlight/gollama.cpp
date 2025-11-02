//go:build !embedallowed_no

package gollama

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// embeddedLibFiles contains all files under the libs directory that should be bundled into the binary.
// The directory is expected to contain platform-specific folders such as darwin_amd64_<version>/,
// each holding the shared libraries for that platform.
//
//go:embed libs/**
var embeddedLibFiles embed.FS

var embeddedCopyMu sync.Mutex

// embeddedPlatformDirName returns the directory name that stores the embedded libraries for the
// given platform and architecture.
func embeddedPlatformDirName(goos, goarch string) string {
	return fmt.Sprintf("%s_%s_%s", goos, goarch, LlamaCppBuild)
}

// embeddedPlatformFSPath returns the filesystem path used inside the embedded FS for the platform.
func embeddedPlatformFSPath(goos, goarch string) string {
	return filepath.ToSlash(filepath.Join("libs", embeddedPlatformDirName(goos, goarch)))
}

// hasEmbeddedLibraryForPlatform returns true if the embedded filesystem contains a library bundle
// for the requested platform/arch pair.
func hasEmbeddedLibraryForPlatform(goos, goarch string) bool {
	path := embeddedPlatformFSPath(goos, goarch)
	if path == "" {
		return false
	}

	if _, err := fs.Stat(embeddedLibFiles, path); err != nil {
		return false
	}
	return true
}

// extractEmbeddedLibrariesTo copies the embedded libraries for the requested platform/arch pair to
// the destination directory, replacing any existing contents.
func extractEmbeddedLibrariesTo(dest, goos, goarch string) error {
	if dest == "" {
		return errors.New("destination path cannot be empty")
	}

	platformPath := embeddedPlatformFSPath(goos, goarch)
	if platformPath == "" {
		return fmt.Errorf("invalid platform %s/%s", goos, goarch)
	}

	// Ensure the platform exists in the embedded filesystem.
	if !hasEmbeddedLibraryForPlatform(goos, goarch) {
		return fmt.Errorf("embedded libraries not available for %s/%s", goos, goarch)
	}

	embeddedCopyMu.Lock()
	defer embeddedCopyMu.Unlock()

	if err := os.RemoveAll(dest); err != nil {
		return fmt.Errorf("failed to clean destination %s: %w", dest, err)
	}
	if err := os.MkdirAll(dest, 0o750); err != nil {
		return fmt.Errorf("failed to create destination %s: %w", dest, err)
	}

	return fs.WalkDir(embeddedLibFiles, platformPath, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel := strings.TrimPrefix(path, platformPath)
		rel = strings.TrimPrefix(rel, "/")
		if rel == "" {
			return nil
		}

		targetPath := filepath.Join(dest, filepath.FromSlash(rel))

		if d.IsDir() {
			if err := os.MkdirAll(targetPath, 0o750); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			return nil
		}

		data, err := fs.ReadFile(embeddedLibFiles, path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		if err := os.WriteFile(targetPath, data, 0o600); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}
		return nil
	})
}
