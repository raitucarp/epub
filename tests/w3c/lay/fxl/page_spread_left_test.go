package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-fxl-page-spread-left
func TestLayFxlPageSpreadLeft(t *testing.T) {
	r := w3ctest.Load(t, "lay-fxl-page-spread-left")

	if !w3ctest.Contains(r.Identifier(), "lay-fxl-page-spread-left") {
		t.Errorf("expected identifier %q, got %v", "lay-fxl-page-spread-left", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-fxl-page-spread-left") {
		t.Errorf("expected title %q, got %v", "lay-fxl-page-spread-left", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
