package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-xml-external-id
func TestPubXmlExternalId(t *testing.T) {
	r := w3ctest.Load(t, "pub-xml-external-id")

	if !w3ctest.Contains(r.Identifier(), "pub-xml-external-id") {
		t.Errorf("expected identifier %q, got %v", "pub-xml-external-id", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pub-xml-external-id") {
		t.Errorf("expected title %q, got %v", "pub-xml-external-id", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
