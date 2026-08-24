package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/css-epub-text-emphasis
func TestCssEpubTextEmphasis(t *testing.T) {
	r := w3ctest.Load(t, "css-epub-text-emphasis")

	if !w3ctest.Contains(r.Identifier(), "css-epub-text-emphasis") {
		t.Errorf("expected identifier %q, got %v", "css-epub-text-emphasis", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "css-epub-text-emphasis") {
		t.Errorf("expected title %q, got %v", "css-epub-text-emphasis", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
