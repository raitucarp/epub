package ocf

import (
	"testing"
)

func TestOCFZipContainer_MimeType(t *testing.T) {
	container := &OCFZipContainer{
		files: map[string][]byte{
			"mimetype": []byte(MimeType),
		},
	}

	mime := container.MimeType()
	if mime != MimeType {
		t.Errorf("Expected MimeType to be %s, but got %s", MimeType, mime)
	}
}

func TestOCFZipContainer_MimeType_Empty(t *testing.T) {
	container := &OCFZipContainer{
		files: map[string][]byte{},
	}

	mime := container.MimeType()
	if mime != "" {
		t.Errorf("Expected empty MimeType, but got %s", mime)
	}
}
