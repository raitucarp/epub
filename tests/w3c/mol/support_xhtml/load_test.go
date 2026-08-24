package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/mol-support_xhtml-load
func TestMolSupportXhtmlLoad(t *testing.T) {
	r := w3ctest.Load(t, "mol-support_xhtml-load")

	if !w3ctest.Contains(r.Identifier(), "mol-support_xhtml-load") {
		t.Errorf("expected identifier %q, got %v", "mol-support_xhtml-load", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "mol-support_xhtml-load") {
		t.Errorf("expected title %q, got %v", "mol-support_xhtml-load", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
