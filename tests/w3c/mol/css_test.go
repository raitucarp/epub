package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/mol-css
func TestMolCss(t *testing.T) {
	r := w3ctest.Load(t, "mol-css")

	if !w3ctest.Contains(r.Identifier(), "mol-css") {
		t.Errorf("expected identifier %q, got %v", "mol-css", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "mol-css") {
		t.Errorf("expected title %q, got %v", "mol-css", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
