package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-dir_creator-rtl
func TestPkgDirCreatorRtl(t *testing.T) {
	r := w3ctest.Load(t, "pkg-dir_creator-rtl")

	if !w3ctest.Contains(r.Identifier(), "pkg-dir_creator-rtl") {
		t.Errorf("expected identifier %q, got %v", "pkg-dir_creator-rtl", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pkg-dir_creator-rtl") {
		t.Errorf("expected title %q, got %v", "pkg-dir_creator-rtl", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
