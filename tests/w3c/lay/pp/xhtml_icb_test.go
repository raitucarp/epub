package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pp-xhtml-icb
func TestLayPpXhtmlIcb(t *testing.T) {
	r := w3ctest.Load(t, "lay-pp-xhtml-icb")

	if !w3ctest.Contains(r.Identifier(), "lay-pp-xhtml-icb") {
		t.Errorf("expected identifier %q, got %v", "lay-pp-xhtml-icb", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-pp-xhtml-icb") {
		t.Errorf("expected title %q, got %v", "lay-pp-xhtml-icb", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
