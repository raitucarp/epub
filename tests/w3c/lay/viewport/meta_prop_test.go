package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-viewport-meta-prop
func TestLayViewportMetaProp(t *testing.T) {
	r := w3ctest.Load(t, "lay-viewport-meta-prop")

	if !w3ctest.Contains(r.Identifier(), "lay-viewport-meta-prop") {
		t.Errorf("expected identifier %q, got %v", "lay-viewport-meta-prop", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-viewport-meta-prop") {
		t.Errorf("expected title %q, got %v", "lay-viewport-meta-prop", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
