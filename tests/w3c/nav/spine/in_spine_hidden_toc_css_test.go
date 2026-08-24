package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/nav-spine_in-spine-hidden-toc-css
func TestNavSpineInSpineHiddenTocCss(t *testing.T) {
	r := w3ctest.Load(t, "nav-spine_in-spine-hidden-toc-css")

	if !w3ctest.Contains(r.Identifier(), "nav-spine_in-spine-hidden-toc-css") {
		t.Errorf("expected identifier %q, got %v", "nav-spine_in-spine-hidden-toc-css", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "nav-spine_in-spine-hidden-toc-css") {
		t.Errorf("expected title %q, got %v", "nav-spine_in-spine-hidden-toc-css", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
