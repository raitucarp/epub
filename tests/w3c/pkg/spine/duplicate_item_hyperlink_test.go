package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-spine-duplicate-item-hyperlink
func TestPkgSpineDuplicateItemHyperlink(t *testing.T) {
	r := w3ctest.Load(t, "pkg-spine-duplicate-item-hyperlink")

	if !w3ctest.Contains(r.Identifier(), "pkg-spine-duplicate-item-hyperlink") {
		t.Errorf("expected identifier %q, got %v", "pkg-spine-duplicate-item-hyperlink", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pkg-spine-duplicate-item-hyperlink") {
		t.Errorf("expected title %q, got %v", "pkg-spine-duplicate-item-hyperlink", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
