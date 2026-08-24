package ocf

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
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

func TestOCFZipContainer_WriteMimetypeFirst(t *testing.T) {
	container := NewOCFZipContainer()
	container.AddMimeType()
	container.AddFile("META-INF/container.xml", []byte("<container/>"))
	container.AddFile("EPUB/content.opf", []byte("<package/>"))

	path := filepath.Join(t.TempDir(), "test.epub")
	err := container.Write(path)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	z, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("Expected to open zip, got %v", err)
	}
	defer z.Close()

	if len(z.File) == 0 {
		t.Fatalf("Expected at least one file in zip")
	}

	if z.File[0].Name != "mimetype" {
		t.Errorf("Expected mimetype to be first entry, got %q", z.File[0].Name)
	}

	if z.File[0].Method != zip.Store {
		t.Errorf("Expected mimetype to be stored uncompressed, got method %v", z.File[0].Method)
	}

	if len(z.File[0].Extra) != 0 {
		t.Errorf("Expected mimetype entry to have no extra field, got %x", z.File[0].Extra)
	}
}

func TestOCFZipContainer_WriteIsReadable(t *testing.T) {
	container := NewOCFZipContainer()
	container.AddMimeType()
	container.AddFile("META-INF/container.xml", []byte("<container/>"))

	path := filepath.Join(t.TempDir(), "test.epub")
	if err := container.Write(path); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("Expected written file to exist, got %v", err)
	}
}
