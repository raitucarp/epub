package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-roll-images-mixed
func TestLayRollImagesMixed(t *testing.T) {
	r := w3ctest.Load(t, "lay-roll-images-mixed")

	if !w3ctest.Contains(r.Identifier(), "lay-roll-images-mixed") {
		t.Errorf("expected identifier %q, got %v", "lay-roll-images-mixed", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "A Apple Pie — embedded images except for two in the spine in a roll") {
		t.Errorf("expected title %q, got %v", "A Apple Pie — embedded images except for two in the spine in a roll", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
