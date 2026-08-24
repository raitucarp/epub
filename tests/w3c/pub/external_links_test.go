package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-external-links
func TestPubExternalLinks(t *testing.T) {
	r := w3ctest.Load(t, "pub-external-links")

	if !w3ctest.Contains(r.Identifier(), "pub-external-links") {
		t.Errorf("expected identifier %q, got %v", "pub-external-links", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pub-external-links") {
		t.Errorf("expected title %q, got %v", "pub-external-links", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
