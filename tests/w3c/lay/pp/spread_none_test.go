package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pp-spread-none
func TestLayPpSpreadNone(t *testing.T) {
	r := w3ctest.Load(t, "lay-pp-spread-none")

	if !w3ctest.Contains(r.Identifier(), "lay-pp-spread-none") {
		t.Errorf("expected identifier %q, got %v", "lay-pp-spread-none", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-pp-spread-none") {
		t.Errorf("expected title %q, got %v", "lay-pp-spread-none", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
