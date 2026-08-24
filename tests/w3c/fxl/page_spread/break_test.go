package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/fxl-page-spread-break
func TestFxlPageSpreadBreak(t *testing.T) {
	r := w3ctest.Load(t, "fxl-page-spread-break")

	if !w3ctest.Contains(r.Identifier(), "fxl-page-spread-break") {
		t.Errorf("expected identifier %q, got %v", "fxl-page-spread-break", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "fxl-page-spread-break") {
		t.Errorf("expected title %q, got %v", "fxl-page-spread-break", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
