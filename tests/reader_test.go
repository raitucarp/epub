package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/raitucarp/epub"
)

func TestOpenReaderError(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		_, err := epub.OpenReader("nonexistent_file.epub")
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})

	t.Run("invalid dummy file", func(t *testing.T) {
		dummyPath := filepath.Join(os.TempDir(), "dummy.epub")
		err := os.WriteFile(dummyPath, []byte("not an epub file"), 0644)
		if err != nil {
			t.Fatalf("failed to create dummy file: %v", err)
		}
		defer os.Remove(dummyPath)

		_, err = epub.OpenReader(dummyPath)
		if err == nil {
			t.Error("expected error for invalid dummy file, got nil")
		}
	})
}
