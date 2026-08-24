package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/ocf-font_obfuscation_bis
// The obfuscated font used a different publication identifier, so de-obfuscation
// with the actual unique identifier must NOT restore a valid TrueType header.
func TestOcfFontObfuscationBis(t *testing.T) {
	r := w3ctest.Load(t, "ocf-font_obfuscation_bis")

	res := r.SelectResourceById("font_truetype")
	if res == nil {
		t.Fatalf("expected font resource")
	}
	if len(res.Content) < 4 {
		t.Fatalf("expected font content")
	}

	// The unique identifier differs from the one used to obfuscate the font, so
	// the sfnt magic must not be restored.
	if res.Content[0] == 0 && res.Content[1] == 1 && res.Content[2] == 0 && res.Content[3] == 0 {
		t.Errorf("expected the font to remain obfuscated (mismatched identifier), got % x", res.Content[:4])
	}
}
