package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-dir-auto_root-unset
func TestPkgDirAutoRootUnset(t *testing.T) {
	r := w3ctest.Load(t, "pkg-dir-auto_root-unset")

	if !w3ctest.Contains(r.Identifier(), "pkg-dir-auto_root-unset") {
		t.Errorf("expected identifier %q, got %v", "pkg-dir-auto_root-unset", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "CSS: הרפתקה חדשה!") {
		t.Errorf("expected title %q, got %v", "CSS: הרפתקה חדשה!", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
