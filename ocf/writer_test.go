package ocf

import (
	"bytes"
	"testing"
)

func TestOCFZipContainer_AddMimeType(t *testing.T) {
	container := NewOCFZipContainer()

	// Initial check
	if len(container.files) != 0 {
		t.Fatalf("Expected newly created container to have no files, got %d", len(container.files))
	}

	// Call the method under test
	container.AddMimeType()

	// Verify the result
	if len(container.files) != 1 {
		t.Fatalf("Expected container to have 1 file, got %d", len(container.files))
	}

	content, exists := container.files["mimetype"]
	if !exists {
		t.Fatalf("Expected file 'mimetype' to exist in container files map")
	}

	expectedContent := []byte(MimeType)
	if !bytes.Equal(content, expectedContent) {
		t.Errorf("Expected mimetype file content to be %s, but got %s", expectedContent, content)
	}
}
