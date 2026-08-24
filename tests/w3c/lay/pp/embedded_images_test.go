package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pp-embedded-images
func TestLayPpEmbeddedImages(t *testing.T) {
	r := w3ctest.Load(t, "lay-pp-embedded-images")

	if !w3ctest.Contains(r.Identifier(), "lay-pp-embedded-images") {
		t.Errorf("expected identifier %q, got %v", "lay-pp-embedded-images", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "A Apple Pie — embedded images") {
		t.Errorf("expected title %q, got %v", "A Apple Pie — embedded images", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
