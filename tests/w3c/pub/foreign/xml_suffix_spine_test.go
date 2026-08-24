package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-foreign_xml-suffix-spine
func TestPubForeignXmlSuffixSpine(t *testing.T) {
	r := w3ctest.Load(t, "pub-foreign_xml-suffix-spine")

	if !w3ctest.Contains(r.Identifier(), "pub-foreign_xml-suffix-spine") {
		t.Errorf("expected identifier %q, got %v", "pub-foreign_xml-suffix-spine", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pub-foreign_xml-suffix-spine") {
		t.Errorf("expected title %q, got %v", "pub-foreign_xml-suffix-spine", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
