package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-dir-auto_root-rtl
func TestPkgDirAutoRootRtl(t *testing.T) {
	r := w3ctest.Load(t, "pkg-dir-auto_root-rtl")

	if !w3ctest.Contains(r.Identifier(), "pkg-dir-auto_root-rtl") {
		t.Errorf("expected identifier %q, got %v", "pkg-dir-auto_root-rtl", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "CSS: مغامرة جديدة!") {
		t.Errorf("expected title %q, got %v", "CSS: مغامرة جديدة!", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
