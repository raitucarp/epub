package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pp-page-spread-right
func TestLayPpPageSpreadRight(t *testing.T) {
	r := w3ctest.Load(t, "lay-pp-page-spread-right")

	if !w3ctest.Contains(r.Identifier(), "lay-pp-page-spread-right") {
		t.Errorf("expected identifier %q, got %v", "lay-pp-page-spread-right", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-pp-page-spread-right") {
		t.Errorf("expected title %q, got %v", "lay-pp-page-spread-right", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
