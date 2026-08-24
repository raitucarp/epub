package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-version-backward
// A package document with a version attribute lower than 3.0 must still be
// opened, preserving the declared version.
func TestPkgVersionBackward(t *testing.T) {
	r := w3ctest.Load(t, "pkg-version-backward")

	if got := r.Version(); got != "0" {
		t.Errorf("expected version %q, got %q", "0", got)
	}
}
