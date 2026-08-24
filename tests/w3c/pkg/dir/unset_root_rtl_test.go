package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-dir_unset-root-rtl
func TestPkgDirUnsetRootRtl(t *testing.T) {
	r := w3ctest.Load(t, "pkg-dir_unset-root-rtl")

	if !w3ctest.Contains(r.Identifier(), "pkg-dir_unset-root-rtl") {
		t.Errorf("expected identifier %q, got %v", "pkg-dir_unset-root-rtl", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "CSS: مغامرة جديدة!") {
		t.Errorf("expected title %q, got %v", "CSS: مغامرة جديدة!", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "ar") {
		t.Errorf("expected language %q, got %v", "ar", r.Language())
	}
}
