// Package w3ctest provides shared helpers for the W3C EPUB test suite port.
// The actual test cases live in nested directories under tests/w3c and import
// this package to load fixtures and build OCF containers.
package w3ctest

import (
	"archive/zip"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/raitucarp/epub"
)

// largeFileThreshold is the size above which a fixture is treated as a "large"
// binary resource and excluded from the committed testdata. Such files (audio,
// layout images, and fonts) are restored by SyncBinaries.
const largeFileThreshold = 500 * 1024

// dataDir resolves the location of the committed epub-tests fixtures relative
// to this source file, independent of the current working directory.
func dataDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("w3ctest: cannot resolve source location")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "testdata", "epub-tests"))
}

// Contains reports whether a slice of strings contains the given value.
func Contains(values []string, want string) bool {
	return slices.Contains(values, want)
}

// EqualStrings reports whether two string slices contain the same values in
// the same order.
func EqualStrings(a, b []string) bool {
	return slices.Equal(a, b)
}

// BuildEpubFromDir packages a loose test directory (mimetype, META-INF, and
// content) into an OCF ZIP container in memory, mirroring the epub-tests
// generateEpubs.sh layout rules: mimetype first and stored uncompressed.
func BuildEpubFromDir(t *testing.T, dir string) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	addEntry := func(name, path string, store bool) {
		t.Helper()
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read %s: %v", path, err)
		}

		header := &zip.FileHeader{Name: name}
		if store {
			header.Method = zip.Store
		} else {
			header.Method = zip.Deflate
		}

		w, err := zw.CreateHeader(header)
		if err != nil {
			t.Fatalf("failed to create zip entry %s: %v", name, err)
		}
		if _, err := w.Write(data); err != nil {
			t.Fatalf("failed to write zip entry %s: %v", name, err)
		}
	}

	addEntry("mimetype", filepath.Join(dir, "mimetype"), true)

	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "mimetype" {
			return nil
		}

		paths = append(paths, rel)
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk %s: %v", dir, err)
	}

	sort.Strings(paths)
	for _, rel := range paths {
		addEntry(rel, filepath.Join(dir, filepath.FromSlash(rel)), false)
	}

	if err := zw.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}

	return buf.Bytes()
}

// Load builds and opens a single test publication by its directory name.
func Load(t *testing.T, name string) epub.Reader {
	t.Helper()

	dir := filepath.Join(dataDir(), name)
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("epub-tests fixture %q not found at %s", name, dir)
	}

	reader, err := epub.NewReader(BuildEpubFromDir(t, dir))
	if err != nil {
		t.Fatalf("failed to open %s: %v", name, err)
	}
	return reader
}

// binaryFile reports whether the given file (relative to the epub-tests
// repository) is a large binary resource that is excluded from the committed
// testdata and therefore restored by SyncBinaries.
func binaryFile(rel string) bool {
	ext := strings.ToLower(filepath.Ext(rel))

	switch ext {
	case ".mp4", ".m4a", ".aac", ".wav", ".dmg", ".z01":
		return true
	case ".otf", ".ttf", ".woff", ".woff2":
		// Keep the obfuscation fonts committed; they are small and required by
		// the font obfuscation tests.
		return !strings.Contains(rel, "ocf-font_obfuscation")
	}

	return false
}

// SyncBinaries clones the W3C epub-tests repository (when absent) and copies
// the large binary resources that are excluded from the committed testdata back
// into the testdata directory. It is idempotent and a no-op once synced.
func SyncBinaries(t *testing.T) {
	t.Helper()

	marker := filepath.Join(dataDir(), ".binaries-synced")
	if _, err := os.Stat(marker); err == nil {
		return
	}

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available; skipping epub-tests binary sync")
	}

	repo := filepath.Join(repoRoot(), "temp", "epub-tests")
	if _, err := os.Stat(filepath.Join(repo, "tests")); err != nil {
		if err := clone(repo); err != nil {
			t.Skipf("failed to clone epub-tests: %v", err)
		}
	}

	src := filepath.Join(repo, "tests")
	var synced int
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		full := filepath.Join(src, filepath.FromSlash(rel))
		if !binaryFile(rel) {
			if info.Size() > largeFileThreshold {
				copyFile(t, full, filepath.Join(dataDir(), filepath.FromSlash(rel)))
				synced++
			}
			return nil
		}

		copyFile(t, full, filepath.Join(dataDir(), filepath.FromSlash(rel)))
		synced++
		return nil
	})
	if err != nil {
		t.Fatalf("failed to sync binaries: %v", err)
	}

	if err := os.WriteFile(marker, []byte("synced"), 0644); err != nil {
		t.Fatalf("failed to write sync marker: %v", err)
	}
	t.Logf("synced %d binary resources", synced)
}

// repoRoot returns the module root directory.
func repoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("w3ctest: cannot resolve source location")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
}

func clone(dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	cmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/w3c/epub-tests.git", dst)
	return cmd.Run()
}

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("failed to read %s: %v", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		t.Fatalf("failed to create dir for %s: %v", dst, err)
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		t.Fatalf("failed to write %s: %v", dst, err)
	}
}
