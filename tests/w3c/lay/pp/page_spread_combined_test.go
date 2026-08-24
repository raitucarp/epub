package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pp-page-spread-combined
func TestLayPpPageSpreadCombined(t *testing.T) {
	r := w3ctest.Load(t, "lay-pp-page-spread-combined")

	if !w3ctest.Contains(r.Identifier(), "lay-pp-page-spread-combined") {
		t.Errorf("expected identifier %q, got %v", "lay-pp-page-spread-combined", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-pp-page-spread-combined") {
		t.Errorf("expected title %q, got %v", "lay-pp-page-spread-combined", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
