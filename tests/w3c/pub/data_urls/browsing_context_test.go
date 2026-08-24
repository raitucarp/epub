package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-data-urls_browsing-context
func TestPubDataUrlsBrowsingContext(t *testing.T) {
	r := w3ctest.Load(t, "pub-data-urls_browsing-context")

	if !w3ctest.Contains(r.Identifier(), "pub-data-urls_browsing-context") {
		t.Errorf("expected identifier %q, got %v", "pub-data-urls_browsing-context", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pub-data-urls_browsing-context") {
		t.Errorf("expected title %q, got %v", "pub-data-urls_browsing-context", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
