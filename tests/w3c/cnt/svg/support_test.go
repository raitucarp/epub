package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/cnt-svg-support
func TestCntSvgSupport(t *testing.T) {
	r := w3ctest.Load(t, "cnt-svg-support")

	if !w3ctest.Contains(r.Identifier(), "cnt-svg-support") {
		t.Errorf("expected identifier %q, got %v", "cnt-svg-support", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "cnt-svg-support") {
		t.Errorf("expected title %q, got %v", "cnt-svg-support", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
