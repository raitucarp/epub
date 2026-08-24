package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pp-xhtml-icb_multi_declarations
func TestLayPpXhtmlIcbMultiDeclarations(t *testing.T) {
	r := w3ctest.Load(t, "lay-pp-xhtml-icb_multi_declarations")

	if !w3ctest.Contains(r.Identifier(), "lay-pp-xhtml-icb_multi_declarations") {
		t.Errorf("expected identifier %q, got %v", "lay-pp-xhtml-icb_multi_declarations", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-pp-xhtml-icb_multi_declarations") {
		t.Errorf("expected title %q, got %v", "lay-pp-xhtml-icb_multi_declarations", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
