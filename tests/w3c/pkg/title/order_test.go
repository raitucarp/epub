package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-title-order
// The first dc:title element must be the primary title.
func TestPkgTitleOrder(t *testing.T) {
	r := w3ctest.Load(t, "pkg-title-order")

	titles := r.Title()
	if len(titles) < 6 {
		t.Fatalf("expected at least 6 titles, got %d", len(titles))
	}
	if titles[0] != "pkg-title-order" {
		t.Errorf("expected first title %q, got %q", "pkg-title-order", titles[0])
	}
}
