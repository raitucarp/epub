package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-meta-whitespace
// Leading, trailing, and repeated internal whitespace in metadata values must
// be collapsed to a single space.
func TestPkgMetaWhitespace(t *testing.T) {
	r := w3ctest.Load(t, "pkg-meta-whitespace")

	creators := r.Author()
	if len(creators) != 1 || creators[0] != "Dave Cramer" {
		t.Errorf("expected normalized creator %q, got %v", "Dave Cramer", creators)
	}
}
