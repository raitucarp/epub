package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/nav-spine_not-in-spine
func TestNavSpineNotInSpine(t *testing.T) {
	r := w3ctest.Load(t, "nav-spine_not-in-spine")

	if !w3ctest.Contains(r.Identifier(), "nav-spine_not-in-spine") {
		t.Errorf("expected identifier %q, got %v", "nav-spine_not-in-spine", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "nav-spine_not-in-spine") {
		t.Errorf("expected title %q, got %v", "nav-spine_not-in-spine", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
