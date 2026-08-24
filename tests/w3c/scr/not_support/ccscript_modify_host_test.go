package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/scr-not-support_ccscript-modify-host
func TestScrNotSupportCcscriptModifyHost(t *testing.T) {
	r := w3ctest.Load(t, "scr-not-support_ccscript-modify-host")

	if !w3ctest.Contains(r.Identifier(), "scr-not-support_ccscript-modify-host") {
		t.Errorf("expected identifier %q, got %v", "scr-not-support_ccscript-modify-host", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "scr-not-support_ccscript-modify-host") {
		t.Errorf("expected title %q, got %v", "scr-not-support_ccscript-modify-host", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
