package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-rendition-flow-pre-pag
func TestLayRenditionFlowPrePag(t *testing.T) {
	r := w3ctest.Load(t, "lay-rendition-flow-pre-pag")

	if !w3ctest.Contains(r.Identifier(), "lay-rendition-flow-pre-pag") {
		t.Errorf("expected identifier %q, got %v", "lay-rendition-flow-pre-pag", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-rendition-flow-pre-pag") {
		t.Errorf("expected title %q, got %v", "lay-rendition-flow-pre-pag", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
