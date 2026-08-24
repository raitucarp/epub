package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pp-layout-default
func TestLayPpLayoutDefault(t *testing.T) {
	r := w3ctest.Load(t, "lay-pp-layout-default")

	if !w3ctest.Contains(r.Identifier(), "lay-pp-layout-default") {
		t.Errorf("expected identifier %q, got %v", "lay-pp-layout-default", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-pp-layout-default") {
		t.Errorf("expected title %q, got %v", "lay-pp-layout-default", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
