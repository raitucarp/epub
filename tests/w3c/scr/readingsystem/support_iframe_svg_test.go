package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/scr-readingsystem-support_iframe_svg
func TestScrReadingsystemSupportIframeSvg(t *testing.T) {
	r := w3ctest.Load(t, "scr-readingsystem-support_iframe_svg")

	if !w3ctest.Contains(r.Identifier(), "scr-readingsystem-support_iframe_svg") {
		t.Errorf("expected identifier %q, got %v", "scr-readingsystem-support_iframe_svg", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "scr-readingsystem-support_iframe_svg") {
		t.Errorf("expected title %q, got %v", "scr-readingsystem-support_iframe_svg", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
