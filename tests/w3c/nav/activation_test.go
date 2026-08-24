package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/nav-activation
func TestNavActivation(t *testing.T) {
	r := w3ctest.Load(t, "nav-activation")

	if !w3ctest.Contains(r.Identifier(), "nav-activation") {
		t.Errorf("expected identifier %q, got %v", "nav-activation", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "nav-activation") {
		t.Errorf("expected title %q, got %v", "nav-activation", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
