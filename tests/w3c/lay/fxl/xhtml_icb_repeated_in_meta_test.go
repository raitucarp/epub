package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-fxl-xhtml-icb_repeated-in-meta
func TestLayFxlXhtmlIcbRepeatedInMeta(t *testing.T) {
	r := w3ctest.Load(t, "lay-fxl-xhtml-icb_repeated-in-meta")

	if !w3ctest.Contains(r.Identifier(), "lay-fxl-xhtml-icb_repeated-in-meta") {
		t.Errorf("expected identifier %q, got %v", "lay-fxl-xhtml-icb_repeated-in-meta", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-fxl-xhtml-icb_repeated-in-meta") {
		t.Errorf("expected title %q, got %v", "lay-fxl-xhtml-icb_repeated-in-meta", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
