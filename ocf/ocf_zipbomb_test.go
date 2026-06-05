package ocf

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestZipBombPrevention(t *testing.T) {
	// Create a zip with a file whose uncompressed size is large, but actual size is 0
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	fh := &zip.FileHeader{
		Name:               "large_file.txt",
		Method:             zip.Store,
		UncompressedSize64: 2000000000, // 2GB (more than maxFileSize)
	}
	f, _ := w.CreateRaw(fh)
	f.Write([]byte("some data"))
	w.Close()

	r, _ := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))

	container := &OCFZipContainer{}
	err := container.readFiles(r)

	if err == nil {
		t.Fatalf("Expected error for file too large, got nil")
	}

	if !strings.Contains(err.Error(), "too large") {
		t.Errorf("Expected 'too large' error, got: %v", err)
	}
}

func TestZipBombPrevention_LimitRead(t *testing.T) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	fh := &zip.FileHeader{
		Name:   "limited_file.txt",
		Method: zip.Deflate,
	}
	f, _ := w.CreateHeader(fh)
	f.Write([]byte("some large data that should be truncated"))
	w.Close()

	r, _ := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	// manually change uncompressed size to trigger truncation test correctly for Deflate
	for _, file := range r.File {
		file.UncompressedSize64 = 5
	}

	container := &OCFZipContainer{}
	err := container.readFiles(r)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	content := string(container.files["limited_file.txt"])
	if len(content) > 5 {
		t.Errorf("Content read is larger than uncompressed size: %d", len(content))
	}
	if content != "some " {
		t.Errorf("Content is not truncated correctly: %s", content)
	}
}
