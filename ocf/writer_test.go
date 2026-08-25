package ocf

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/raitucarp/epub/pkg"
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

func TestOCFZipContainer_AddPackage(t *testing.T) {
	container := NewOCFZipContainer()

	p := pkg.Package{
		Version:          "3.0",
		UniqueIdentifier: "pub-id",
	}
	p.Metadata.Identifiers = append(p.Metadata.Identifiers, pkg.DCIdentifier{ID: "pub-id", Value: "urn:test"})

	if err := container.AddPackage("EPUB/package.opf", p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, ok := container.files["EPUB/package.opf"]
	if !ok {
		t.Fatal("expected package file to be added")
	}
	if !strings.HasPrefix(string(content), xml.Header) {
		t.Errorf("expected package content to start with xml header, got %q", string(content[:min(len(content), 40)]))
	}
	if !strings.Contains(string(content), "urn:test") {
		t.Errorf("expected package to contain identifier, got %q", string(content))
	}
}

func TestOCFZipContainer_AddContainerXML(t *testing.T) {
	container := NewOCFZipContainer()

	if err := container.AddContainerXML("EPUB/package.opf", "EPUB/pre.opf"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, ok := container.files["META-INF/container.xml"]
	if !ok {
		t.Fatal("expected container.xml to be added")
	}
	s := string(content)
	if !strings.Contains(s, "EPUB/package.opf") || !strings.Contains(s, "EPUB/pre.opf") {
		t.Errorf("expected both rootfiles in container.xml, got %q", s)
	}
	if !strings.Contains(s, EPUBContainerMime) {
		t.Errorf("expected media type in container.xml, got %q", s)
	}
}

func TestOCFZipContainer_Write_InvalidPath(t *testing.T) {
	container := NewOCFZipContainer()
	container.AddMimeType()

	path := filepath.Join(t.TempDir(), "missing-dir", "out.epub")
	if err := container.Write(path); err == nil {
		t.Error("expected error writing to non-existent directory, got nil")
	}
}
