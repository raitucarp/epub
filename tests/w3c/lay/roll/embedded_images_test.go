package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-roll-embedded-images
func TestLayRollEmbeddedImages(t *testing.T) {
	r := w3ctest.Load(t, "lay-roll-embedded-images")

	if !w3ctest.Contains(r.Identifier(), "lay-roll-embedded-images") {
		t.Errorf("expected identifier %q, got %v", "lay-roll-embedded-images", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "A Apple Pie — embedded images in roll") {
		t.Errorf("expected title %q, got %v", "A Apple Pie — embedded images in roll", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
