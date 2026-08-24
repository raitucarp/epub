package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-dir_but_not_content
func TestPkgDirButNotContent(t *testing.T) {
	r := w3ctest.Load(t, "pkg-dir_but_not_content")

	if !w3ctest.Contains(r.Identifier(), "pkg-dir_but_not_content") {
		t.Errorf("expected identifier %q, got %v", "pkg-dir_but_not_content", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "CSS: הרפתקה חדשה!") {
		t.Errorf("expected title %q, got %v", "CSS: הרפתקה חדשה!", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "he") {
		t.Errorf("expected language %q, got %v", "he", r.Language())
	}
}
