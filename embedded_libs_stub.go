//go:build embedallowed_no

package gollama

import "fmt"

func embeddedPlatformDirName(goos, goarch string) string {
	return fmt.Sprintf("%s_%s_%s", goos, goarch, LlamaCppBuild)
}

func embeddedPlatformFSPath(goos, goarch string) string {
	return ""
}

func hasEmbeddedLibraryForPlatform(goos, goarch string) bool {
	return false
}

func listEmbeddedPlatformDirs() ([]string, error) {
	return nil, nil
}

func extractEmbeddedLibrariesTo(dest, goos, goarch string) error {
	return fmt.Errorf("embedded libraries disabled by build tag")
}
