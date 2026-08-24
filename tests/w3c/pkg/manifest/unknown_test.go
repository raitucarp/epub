package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-manifest-unknown
func TestPkgManifestUnknown(t *testing.T) {
	r := w3ctest.Load(t, "pkg-manifest-unknown")

	if !w3ctest.Contains(r.Identifier(), "pkg-manifest-unknown") {
		t.Errorf("expected identifier %q, got %v", "pkg-manifest-unknown", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pkg-manifest-unknown") {
		t.Errorf("expected title %q, got %v", "pkg-manifest-unknown", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
