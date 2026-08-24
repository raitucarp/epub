package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-foreign_json-spine
func TestPubForeignJsonSpine(t *testing.T) {
	r := w3ctest.Load(t, "pub-foreign_json-spine")

	if !w3ctest.Contains(r.Identifier(), "pub-foreign_json-spine") {
		t.Errorf("expected identifier %q, got %v", "pub-foreign_json-spine", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pub-foreign_json-spine") {
		t.Errorf("expected title %q, got %v", "pub-foreign_json-spine", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
