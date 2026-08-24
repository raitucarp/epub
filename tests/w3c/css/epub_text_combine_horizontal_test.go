package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/css-epub-text-combine-horizontal
func TestCssEpubTextCombineHorizontal(t *testing.T) {
	r := w3ctest.Load(t, "css-epub-text-combine-horizontal")

	if !w3ctest.Contains(r.Identifier(), "css-epub-text-combine-horizontal") {
		t.Errorf("expected identifier %q, got %v", "css-epub-text-combine-horizontal", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "css-epub-text-combine-horizontal") {
		t.Errorf("expected title %q, got %v", "css-epub-text-combine-horizontal", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
