package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/css-epub-word-break
func TestCssEpubWordBreak(t *testing.T) {
	r := w3ctest.Load(t, "css-epub-word-break")

	if !w3ctest.Contains(r.Identifier(), "css-epub-word-break") {
		t.Errorf("expected identifier %q, got %v", "css-epub-word-break", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "css-epub-word-break") {
		t.Errorf("expected title %q, got %v", "css-epub-word-break", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
