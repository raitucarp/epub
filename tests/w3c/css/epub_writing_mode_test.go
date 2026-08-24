package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/css-epub-writing-mode
func TestCssEpubWritingMode(t *testing.T) {
	r := w3ctest.Load(t, "css-epub-writing-mode")

	if !w3ctest.Contains(r.Identifier(), "css-epub-writing-mode") {
		t.Errorf("expected identifier %q, got %v", "css-epub-writing-mode", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "css-epub-writing-mode") {
		t.Errorf("expected title %q, got %v", "css-epub-writing-mode", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
