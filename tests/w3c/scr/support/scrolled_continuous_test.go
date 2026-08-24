package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/scr-support_scrolled-continuous
func TestScrSupportScrolledContinuous(t *testing.T) {
	r := w3ctest.Load(t, "scr-support_scrolled-continuous")

	if !w3ctest.Contains(r.Identifier(), "scr-support_scrolled-continuous") {
		t.Errorf("expected identifier %q, got %v", "scr-support_scrolled-continuous", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "scr-support_scrolled-continuous") {
		t.Errorf("expected title %q, got %v", "scr-support_scrolled-continuous", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
