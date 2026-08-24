package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-spine-order-svg
func TestPkgSpineOrderSvg(t *testing.T) {
	r := w3ctest.Load(t, "pkg-spine-order-svg")

	if !w3ctest.Contains(r.Identifier(), "pkg-spine-order-svg") {
		t.Errorf("expected identifier %q, got %v", "pkg-spine-order-svg", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pkg-spine-order-svg") {
		t.Errorf("expected title %q, got %v", "pkg-spine-order-svg", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
