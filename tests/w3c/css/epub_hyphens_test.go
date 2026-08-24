package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/css-epub-hyphens
func TestCssEpubHyphens(t *testing.T) {
	r := w3ctest.Load(t, "css-epub-hyphens")

	if !w3ctest.Contains(r.Identifier(), "css-epub-hyphens") {
		t.Errorf("expected identifier %q, got %v", "css-epub-hyphens", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "css-epub-hyphens") {
		t.Errorf("expected title %q, got %v", "css-epub-hyphens", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
