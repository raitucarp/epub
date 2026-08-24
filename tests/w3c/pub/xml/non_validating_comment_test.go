package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-xml-non-validating_comment
func TestPubXmlNonValidatingComment(t *testing.T) {
	r := w3ctest.Load(t, "pub-xml-non-validating_comment")

	if !w3ctest.Contains(r.Identifier(), "pub-xml-non-validating_comment") {
		t.Errorf("expected identifier %q, got %v", "pub-xml-non-validating_comment", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pub-xml-non-validating_comment") {
		t.Errorf("expected title %q, got %v", "pub-xml-non-validating_comment", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
