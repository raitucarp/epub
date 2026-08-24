package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/cnt-svg-css-inclusion
func TestCntSvgCssInclusion(t *testing.T) {
	r := w3ctest.Load(t, "cnt-svg-css-inclusion")

	if !w3ctest.Contains(r.Identifier(), "cnt-svg-css-inclusion") {
		t.Errorf("expected identifier %q, got %v", "cnt-svg-css-inclusion", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "cnt-svg-css-inclusion") {
		t.Errorf("expected title %q, got %v", "cnt-svg-css-inclusion", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
