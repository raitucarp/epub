package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-cmt-svg
// An SVG core media type resource must be recognized as an SVG content document.
func TestPubCmtSvg(t *testing.T) {
	r := w3ctest.Load(t, "pub-cmt-svg")

	svgs := r.ContentDocumentSVG()
	if len(svgs) != 1 {
		t.Errorf("expected 1 SVG content document, got %d", len(svgs))
	}
}
