package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-creator-order
// The reading order of dc:creator elements must be preserved.
func TestPkgCreatorOrder(t *testing.T) {
	r := w3ctest.Load(t, "pkg-creator-order")

	expected := []string{"Dave Cramer", "Wendy Reid", "Dan Lazin", "Ivan Herman", "Brady Duga"}
	if got := r.Author(); !w3ctest.EqualStrings(got, expected) {
		t.Errorf("expected creators %v, got %v", expected, got)
	}
}
