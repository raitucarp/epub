package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-page-layout-both-spread
func TestLayPageLayoutBothSpread(t *testing.T) {
	r := w3ctest.Load(t, "lay-page-layout-both-spread")

	if !w3ctest.Contains(r.Identifier(), "lay-page-layout-both-spread") {
		t.Errorf("expected identifier %q, got %v", "lay-page-layout-both-spread", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-page-layout-both-spread") {
		t.Errorf("expected title %q, got %v", "lay-page-layout-both-spread", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
