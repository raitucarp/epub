package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pp-images-mixed
func TestLayPpImagesMixed(t *testing.T) {
	r := w3ctest.Load(t, "lay-pp-images-mixed")

	if !w3ctest.Contains(r.Identifier(), "lay-pp-images-mixed") {
		t.Errorf("expected identifier %q, got %v", "lay-pp-images-mixed", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "A Apple Pie — embedded images except for two in the spine") {
		t.Errorf("expected title %q, got %v", "A Apple Pie — embedded images except for two in the spine", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
