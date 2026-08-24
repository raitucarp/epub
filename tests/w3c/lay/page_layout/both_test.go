package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-page-layout-both
func TestLayPageLayoutBoth(t *testing.T) {
	r := w3ctest.Load(t, "lay-page-layout-both")

	if !w3ctest.Contains(r.Identifier(), "lay-page-layout-both") {
		t.Errorf("expected identifier %q, got %v", "lay-page-layout-both", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-page-layout-both") {
		t.Errorf("expected title %q, got %v", "lay-page-layout-both", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
