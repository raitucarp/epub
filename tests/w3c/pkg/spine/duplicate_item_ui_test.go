package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-spine-duplicate-item-ui
func TestPkgSpineDuplicateItemUi(t *testing.T) {
	r := w3ctest.Load(t, "pkg-spine-duplicate-item-ui")

	if !w3ctest.Contains(r.Identifier(), "pkg-spine-duplicate-item-ui") {
		t.Errorf("expected identifier %q, got %v", "pkg-spine-duplicate-item-ui", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pkg-spine-duplicate-item-ui") {
		t.Errorf("expected title %q, got %v", "pkg-spine-duplicate-item-ui", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
