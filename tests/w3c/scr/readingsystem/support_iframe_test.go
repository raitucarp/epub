package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/scr-readingsystem-support_iframe
func TestScrReadingsystemSupportIframe(t *testing.T) {
	r := w3ctest.Load(t, "scr-readingsystem-support_iframe")

	if !w3ctest.Contains(r.Identifier(), "scr-readingsystem-support_iframe") {
		t.Errorf("expected identifier %q, got %v", "scr-readingsystem-support_iframe", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "scr-readingsystem-support_iframe") {
		t.Errorf("expected title %q, got %v", "scr-readingsystem-support_iframe", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
