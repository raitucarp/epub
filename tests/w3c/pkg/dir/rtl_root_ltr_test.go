package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-dir_rtl-root-ltr
func TestPkgDirRtlRootLtr(t *testing.T) {
	r := w3ctest.Load(t, "pkg-dir_rtl-root-ltr")

	if !w3ctest.Contains(r.Identifier(), "pkg-dir_rtl-root-ltr") {
		t.Errorf("expected identifier %q, got %v", "pkg-dir_rtl-root-ltr", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "CSS: הרפתקה חדשה!") {
		t.Errorf("expected title %q, got %v", "CSS: הרפתקה חדשה!", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
