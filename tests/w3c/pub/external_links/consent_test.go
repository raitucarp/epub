package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-external-links_consent
func TestPubExternalLinksConsent(t *testing.T) {
	r := w3ctest.Load(t, "pub-external-links_consent")

	if !w3ctest.Contains(r.Identifier(), "pub-external-links_consent") {
		t.Errorf("expected identifier %q, got %v", "pub-external-links_consent", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pub-external-links_consent") {
		t.Errorf("expected title %q, got %v", "pub-external-links_consent", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
