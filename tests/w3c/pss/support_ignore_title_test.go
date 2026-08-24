package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pss-support_ignore-title
func TestPssSupportIgnoreTitle(t *testing.T) {
	r := w3ctest.Load(t, "pss-support_ignore-title")

	if !w3ctest.Contains(r.Identifier(), "pss-support_ignore-title") {
		t.Errorf("expected identifier %q, got %v", "pss-support_ignore-title", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pss-support_ignore-title") {
		t.Errorf("expected title %q, got %v", "pss-support_ignore-title", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
