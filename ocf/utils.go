package ocf

import (
	"path/filepath"
	"strings"
)

func getRootDirectory(path string) string {
	cleanPath := filepath.Clean(path)

	// Split the path and return the first component
	parts := strings.Split(cleanPath, string(filepath.Separator))
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}

	volume := filepath.VolumeName(path)

	if volume != "" {
		// Windows path with drive letter
		return volume + string(filepath.Separator)
	}
	// Unix-like path
	return string(filepath.Separator)
}
