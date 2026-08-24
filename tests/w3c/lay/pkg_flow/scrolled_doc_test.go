package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pkg-flow-scrolled-doc
func TestLayPkgFlowScrolledDoc(t *testing.T) {
	r := w3ctest.Load(t, "lay-pkg-flow-scrolled-doc")

	if !w3ctest.Contains(r.Identifier(), "lay-pkg-flow-scrolled-doc") {
		t.Errorf("expected identifier %q, got %v", "lay-pkg-flow-scrolled-doc", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-pkg-flow-scrolled-doc") {
		t.Errorf("expected title %q, got %v", "lay-pkg-flow-scrolled-doc", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
