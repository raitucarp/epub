package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/nav-non-text_img_title
func TestNavNonTextImgTitle(t *testing.T) {
	r := w3ctest.Load(t, "nav-non-text_img_title")

	if !w3ctest.Contains(r.Identifier(), "nav-non-text_img_title") {
		t.Errorf("expected identifier %q, got %v", "nav-non-text_img_title", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "nav-non-text_img_title") {
		t.Errorf("expected title %q, got %v", "nav-non-text_img_title", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
