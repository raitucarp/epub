package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/css-epub-line-break
func TestCssEpubLineBreak(t *testing.T) {
	r := w3ctest.Load(t, "css-epub-line-break")

	if !w3ctest.Contains(r.Identifier(), "css-epub-line-break") {
		t.Errorf("expected identifier %q, got %v", "css-epub-line-break", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "css-epub-line-break") {
		t.Errorf("expected title %q, got %v", "css-epub-line-break", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
