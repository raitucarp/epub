package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/ocf-font_obfuscation
// An obfuscated font must be de-obfuscated using the identifier-based algorithm
// (SHA-1 of the unique identifier), restoring the TrueType header.
func TestOcfFontObfuscation(t *testing.T) {
	r := w3ctest.Load(t, "ocf-font_obfuscation")

	res := r.SelectResourceById("font_truetype")
	if res == nil {
		t.Fatalf("expected font resource")
	}
	if len(res.Content) < 4 {
		t.Fatalf("expected font content")
	}

	// A TrueType font starts with the sfnt version magic 0x00010000.
	if !(res.Content[0] == 0 && res.Content[1] == 1 && res.Content[2] == 0 && res.Content[3] == 0) {
		t.Errorf("expected de-obfuscated TTF magic, got % x", res.Content[:4])
	}
}
