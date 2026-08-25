package ocf

import (
	"path"
	"strings"
)

// getRootDirectory returns the first path component (the root directory) of a
// slash-separated path. OCF container paths always use forward slashes.
func getRootDirectory(p string) string {
	clean := path.Clean(p)
	root, _, _ := strings.Cut(clean, "/")
	return root
}
