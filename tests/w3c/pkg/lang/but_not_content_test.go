package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-lang_but_not_content
func TestPkgLangButNotContent(t *testing.T) {
	r := w3ctest.Load(t, "pkg-lang_but_not_content")

	if !w3ctest.Contains(r.Identifier(), "pkg-lang_but_not_content") {
		t.Errorf("expected identifier %q, got %v", "pkg-lang_but_not_content", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "Le contenu n'est pas automatiquement français") {
		t.Errorf("expected title %q, got %v", "Le contenu n'est pas automatiquement français", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "fr") {
		t.Errorf("expected language %q, got %v", "fr", r.Language())
	}
}
