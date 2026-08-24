package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-unique-id
// The dc:identifier referenced by unique-identifier must be resolved as the
// publication's unique identifier.
func TestPkgUniqueID(t *testing.T) {
	r := w3ctest.Load(t, "pkg-unique-id")

	if got := r.UID(); got != "pkg-unique-id" {
		t.Errorf("expected unique id %q, got %q", "pkg-unique-id", got)
	}
}
