package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-unique-id_duplicate
func TestPkgUniqueIdDuplicate(t *testing.T) {
	r := w3ctest.Load(t, "pkg-unique-id_duplicate")

	if !w3ctest.Contains(r.Identifier(), "pkg-unique-id") {
		t.Errorf("expected identifier %q, got %v", "pkg-unique-id", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pkg-unique-id_duplicate") {
		t.Errorf("expected title %q, got %v", "pkg-unique-id_duplicate", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
