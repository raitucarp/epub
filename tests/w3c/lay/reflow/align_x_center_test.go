package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-reflow-align-x-center
func TestLayReflowAlignXCenter(t *testing.T) {
	r := w3ctest.Load(t, "lay-reflow-align-x-center")

	if !w3ctest.Contains(r.Identifier(), "lay-reflow-align-x-center") {
		t.Errorf("expected identifier %q, got %v", "lay-reflow-align-x-center", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-reflow-align-x-center") {
		t.Errorf("expected title %q, got %v", "lay-reflow-align-x-center", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
