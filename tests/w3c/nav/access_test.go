package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/nav-access
func TestNavAccess(t *testing.T) {
	r := w3ctest.Load(t, "nav-access")

	if !w3ctest.Contains(r.Identifier(), "nav-access") {
		t.Errorf("expected identifier %q, got %v", "nav-access", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "nav-access") {
		t.Errorf("expected title %q, got %v", "nav-access", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
