package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-roll-images-in-spine
func TestLayRollImagesInSpine(t *testing.T) {
	r := w3ctest.Load(t, "lay-roll-images-in-spine")

	if !w3ctest.Contains(r.Identifier(), "lay-roll-images-in-spine") {
		t.Errorf("expected identifier %q, got %v", "lay-roll-images-in-spine", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "A Apple Pie — images in the spine in roll") {
		t.Errorf("expected title %q, got %v", "A Apple Pie — images in the spine in roll", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
