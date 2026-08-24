package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/cnt-css-fonts_ot
func TestCntCssFontsOt(t *testing.T) {
	r := w3ctest.Load(t, "cnt-css-fonts_ot")

	if !w3ctest.Contains(r.Identifier(), "cnt-css-fonts_ot") {
		t.Errorf("expected identifier %q, got %v", "cnt-css-fonts_ot", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "cnt-css-fonts_ot") {
		t.Errorf("expected title %q, got %v", "cnt-css-fonts_ot", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
