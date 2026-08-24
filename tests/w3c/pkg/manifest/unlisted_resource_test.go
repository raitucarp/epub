package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-manifest-unlisted-resource
func TestPkgManifestUnlistedResource(t *testing.T) {
	r := w3ctest.Load(t, "pkg-manifest-unlisted-resource")

	if !w3ctest.Contains(r.Identifier(), "pkg-manifest-unlisted-resource") {
		t.Errorf("expected identifier %q, got %v", "pkg-manifest-unlisted-resource", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pkg-manifest-unlisted-resource") {
		t.Errorf("expected title %q, got %v", "pkg-manifest-unlisted-resource", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
