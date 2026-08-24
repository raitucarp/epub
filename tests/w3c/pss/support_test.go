package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pss-support
func TestPssSupport(t *testing.T) {
	r := w3ctest.Load(t, "pss-support")

	if !w3ctest.Contains(r.Identifier(), "pss-support") {
		t.Errorf("expected identifier %q, got %v", "pss-support", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pss-support") {
		t.Errorf("expected title %q, got %v", "pss-support", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
