package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pp-embedded-images-svg
func TestLayPpEmbeddedImagesSvg(t *testing.T) {
	r := w3ctest.Load(t, "lay-pp-embedded-images-svg")

	if !w3ctest.Contains(r.Identifier(), "lay-pp-embedded-images-svg") {
		t.Errorf("expected identifier %q, got %v", "lay-pp-embedded-images-svg", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "A Apple Pie — embedded images in SVG") {
		t.Errorf("expected title %q, got %v", "A Apple Pie — embedded images in SVG", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
