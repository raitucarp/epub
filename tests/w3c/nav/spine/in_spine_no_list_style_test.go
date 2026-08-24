package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/nav-spine_in-spine-no-list-style
func TestNavSpineInSpineNoListStyle(t *testing.T) {
	r := w3ctest.Load(t, "nav-spine_in-spine-no-list-style")

	if !w3ctest.Contains(r.Identifier(), "nav-spine_in-spine-no-list-style") {
		t.Errorf("expected identifier %q, got %v", "nav-spine_in-spine-no-list-style", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "nav-spine_in-spine-no-list-style") {
		t.Errorf("expected title %q, got %v", "nav-spine_in-spine-no-list-style", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
