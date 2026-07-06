package tests

import (
	"path"
	"strings"
	"testing"

	"github.com/raitucarp/epub"
)

func TestReadContentMarkdownById(t *testing.T) {
	reader, err := epub.OpenReader(path.Join(epubPath, "arthur-conan-doyle_the-white-company.epub"))
	if err != nil {
		t.Fatalf("Failed to open epub: %v", err)
	}

	t.Run("Existing ID", func(t *testing.T) {
		// "chapter-1.xhtml" is an ID in the epub
		md := reader.ReadContentMarkdownById("chapter-1.xhtml")

		if md == "" {
			t.Error("Expected non-empty markdown for existing ID, got empty string")
		}

		if !strings.Contains(md, "How the Black Sheep Came Forth from the Fold") {
			t.Errorf("Expected markdown to contain chapter title, got: %v", md)
		}
	})

	t.Run("Non-existent ID", func(t *testing.T) {
		md := reader.ReadContentMarkdownById("non-existent-id")

		if md != "" {
			t.Errorf("Expected empty markdown for non-existent ID, got: %v", md)
		}
	})
}
