package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-cmt-png
// A PNG core media type resource must be decodable as an image.
func TestPubCmtPng(t *testing.T) {
	r := w3ctest.Load(t, "pub-cmt-png")

	img := r.ReadImageById("img001")
	if img == nil {
		t.Fatalf("expected to decode image resource img001")
	}
	if (*img).Bounds().Dx() == 0 || (*img).Bounds().Dy() == 0 {
		t.Errorf("expected non-zero image dimensions, got %v", (*img).Bounds())
	}
}
