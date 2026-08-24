package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-spine-duplicate-item-rendering
func TestPkgSpineDuplicateItemRendering(t *testing.T) {
	r := w3ctest.Load(t, "pkg-spine-duplicate-item-rendering")

	if !w3ctest.Contains(r.Identifier(), "pkg-spine-duplicate-item-rendering") {
		t.Errorf("expected identifier %q, got %v", "pkg-spine-duplicate-item-rendering", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pkg-spine-duplicate-item-rendering") {
		t.Errorf("expected title %q, got %v", "pkg-spine-duplicate-item-rendering", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
