package ocf

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestZipSlipPrevention(t *testing.T) {
	tests := []struct {
		name        string
		zipFilename string
		wantErr     bool
	}{
		{"valid path", "valid.txt", false},
		{"valid directory path", "folder/valid.txt", false},
		{"absolute path", "/etc/passwd", true},
		{"parent directory traversal", "../../../etc/passwd", true},
		{"backslash path traversal", "..\\..\\etc\\passwd", true},
		{"backslash absolute path", "C:\\Windows\\system.ini", true},
		{"hidden parent traversal", "folder/../../etc/passwd", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			w := zip.NewWriter(buf)

			fh := &zip.FileHeader{
				Name:   tt.zipFilename,
				Method: zip.Store,
			}
			f, _ := w.CreateHeader(fh)
			f.Write([]byte("some data"))
			w.Close()

			r, _ := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))

			container := &OCFZipContainer{}
			err := container.readFiles(r)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for %q, got nil", tt.zipFilename)
				} else if !strings.Contains(err.Error(), "invalid path in zip") {
					t.Errorf("Expected 'invalid path in zip' error for %q, got: %v", tt.zipFilename, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for %q, got: %v", tt.zipFilename, err)
				}
			}
		})
	}
}
