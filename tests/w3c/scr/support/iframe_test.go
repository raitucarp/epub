package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/scr-support_iframe
func TestScrSupportIframe(t *testing.T) {
	r := w3ctest.Load(t, "scr-support_iframe")

	if !w3ctest.Contains(r.Identifier(), "scr-support_iframe") {
		t.Errorf("expected identifier %q, got %v", "scr-support_iframe", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "scr-support_iframe") {
		t.Errorf("expected title %q, got %v", "scr-support_iframe", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
