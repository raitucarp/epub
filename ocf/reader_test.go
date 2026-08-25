package ocf

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// buildZip creates an in-memory ZIP archive from a map of file names to
// contents. The mimetype entry is always stored first and uncompressed.
func buildZip(t *testing.T, files map[string][]byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	write := func(name string, data []byte, store bool) {
		header := &zip.FileHeader{Name: name}
		if store {
			header.Method = zip.Store
		} else {
			header.Method = zip.Deflate
		}
		w, err := zw.CreateHeader(header)
		if err != nil {
			t.Fatalf("create zip entry %s: %v", name, err)
		}
		if _, err := w.Write(data); err != nil {
			t.Fatalf("write zip entry %s: %v", name, err)
		}
	}

	if data, ok := files["mimetype"]; ok {
		write("mimetype", data, true)
	}
	for name, data := range files {
		if name == "mimetype" {
			continue
		}
		write(name, data, false)
	}

	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	return buf.Bytes()
}

const validContainerXML = `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="EPUB/package.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>`

func validEPUBFiles() map[string][]byte {
	return map[string][]byte{
		"mimetype":               []byte(MimeType),
		"META-INF/container.xml": []byte(validContainerXML),
	}
}

func TestNewReader(t *testing.T) {
	t.Run("valid epub", func(t *testing.T) {
		data := buildZip(t, validEPUBFiles())

		container, err := NewReader(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if container.MimeType() != MimeType {
			t.Errorf("expected mimetype %q, got %q", MimeType, container.MimeType())
		}
		if len(container.Container().RootFiles.RootFile) != 1 {
			t.Fatalf("expected 1 rootfile, got %d", len(container.Container().RootFiles.RootFile))
		}
		if got := container.Container().RootFiles.RootFile[0].FullPath; got != "EPUB/package.opf" {
			t.Errorf("unexpected rootfile full-path %q", got)
		}
	})

	t.Run("invalid bytes", func(t *testing.T) {
		if _, err := NewReader([]byte("not a zip")); err == nil {
			t.Error("expected error for invalid bytes, got nil")
		}
	})

	t.Run("mimetype mismatch", func(t *testing.T) {
		files := validEPUBFiles()
		files["mimetype"] = []byte("application/zip")
		data := buildZip(t, files)

		if _, err := NewReader(data); err == nil {
			t.Error("expected error for mimetype mismatch, got nil")
		}
	})

	t.Run("missing container", func(t *testing.T) {
		files := map[string][]byte{
			"mimetype": []byte(MimeType),
		}
		data := buildZip(t, files)

		if _, err := NewReader(data); err == nil {
			t.Error("expected error for missing container.xml, got nil")
		}
	})
}

func TestOpenReader(t *testing.T) {
	t.Run("valid epub file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "valid.epub")
		if err := os.WriteFile(path, buildZip(t, validEPUBFiles()), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		container, err := OpenReader(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if container.MimeType() != MimeType {
			t.Errorf("unexpected mimetype %q", container.MimeType())
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		if _, err := OpenReader(filepath.Join(t.TempDir(), "missing.epub")); err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})

	t.Run("invalid file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "invalid.epub")
		if err := os.WriteFile(path, []byte("garbage"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if _, err := OpenReader(path); err == nil {
			t.Error("expected error for invalid file, got nil")
		}
	})
}

func TestNewContainerAndParse_TrimsMimetypeWhitespace(t *testing.T) {
	files := validEPUBFiles()
	files["mimetype"] = []byte("application/epub+zip\r\n")
	data := buildZip(t, files)

	container, err := NewReader(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if container.MimeType() != MimeType {
		t.Errorf("expected mimetype to be trimmed, got %q", container.MimeType())
	}
}

func TestNewReader_InvalidPathInZip(t *testing.T) {
	files := map[string][]byte{
		"mimetype":               []byte(MimeType),
		"META-INF/container.xml": []byte(validContainerXML),
		"../etc/passwd":          []byte("evil"),
	}
	data := buildZip(t, files)

	if _, err := NewReader(data); err == nil {
		t.Error("expected error for zip path traversal, got nil")
	}
}

func TestReadFiles_SkipsDirectories(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	dirHeader := &zip.FileHeader{Name: "EPUB/"}
	dirHeader.SetMode(0o755 | os.ModeDir)
	if _, err := zw.CreateHeader(dirHeader); err != nil {
		t.Fatalf("create dir entry: %v", err)
	}

	fileHeader := &zip.FileHeader{Name: "mimetype", Method: zip.Store}
	w, err := zw.CreateHeader(fileHeader)
	if err != nil {
		t.Fatalf("create file entry: %v", err)
	}
	if _, err := w.Write([]byte(MimeType)); err != nil {
		t.Fatalf("write file entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}

	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("new reader: %v", err)
	}

	container := &OCFZipContainer{}
	if err := container.readFiles(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := container.files["EPUB/"]; ok {
		t.Error("expected directory entry to be skipped")
	}
	if _, ok := container.files["mimetype"]; !ok {
		t.Error("expected mimetype file to be read")
	}
}

func TestReadFiles_OpenError(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	fh := &zip.FileHeader{Name: "file.txt", Method: 99}
	f, err := zw.CreateRaw(fh)
	if err != nil {
		t.Fatalf("create raw entry: %v", err)
	}
	if _, err := f.Write([]byte("data")); err != nil {
		t.Fatalf("write entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}

	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("new reader: %v", err)
	}

	container := &OCFZipContainer{}
	if err := container.readFiles(r); err == nil {
		t.Error("expected error for unsupported compression algorithm, got nil")
	}
}

func TestReadFiles_ReadError(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	fh := &zip.FileHeader{Name: "file.txt", Method: zip.Deflate}
	f, err := zw.CreateRaw(fh)
	if err != nil {
		t.Fatalf("create raw entry: %v", err)
	}
	if _, err := f.Write([]byte("this is not deflate-compressed data")); err != nil {
		t.Fatalf("write entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}

	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("new reader: %v", err)
	}

	// The raw entry has an unknown uncompressed size; force a positive limit so
	// the corrupt deflate data is actually read and fails to decompress.
	for _, file := range r.File {
		file.UncompressedSize64 = 1024
	}

	container := &OCFZipContainer{}
	if err := container.readFiles(r); err == nil {
		t.Error("expected error for corrupt deflate data, got nil")
	}
}
