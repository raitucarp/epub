package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-meta-unknown
func TestPkgMetaUnknown(t *testing.T) {
	r := w3ctest.Load(t, "pkg-meta-unknown")

	if !w3ctest.Contains(r.Identifier(), "pkg-meta-unknown") {
		t.Errorf("expected identifier %q, got %v", "pkg-meta-unknown", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pkg-meta-unknown") {
		t.Errorf("expected title %q, got %v", "pkg-meta-unknown", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
