package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/mol-support_xhtml-load-fxl
func TestMolSupportXhtmlLoadFxl(t *testing.T) {
	r := w3ctest.Load(t, "mol-support_xhtml-load-fxl")

	if !w3ctest.Contains(r.Identifier(), "mol-support_xhtml-load-fxl") {
		t.Errorf("expected identifier %q, got %v", "mol-support_xhtml-load-fxl", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "mol-support_xhtml-load-fxl") {
		t.Errorf("expected title %q, got %v", "mol-support_xhtml-load-fxl", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
