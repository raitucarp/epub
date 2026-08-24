package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/css-epub-text-orientation
func TestCssEpubTextOrientation(t *testing.T) {
	r := w3ctest.Load(t, "css-epub-text-orientation")

	if !w3ctest.Contains(r.Identifier(), "css-epub-text-orientation") {
		t.Errorf("expected identifier %q, got %v", "css-epub-text-orientation", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "css-epub-text-orientation") {
		t.Errorf("expected title %q, got %v", "css-epub-text-orientation", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
