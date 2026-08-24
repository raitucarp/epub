package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/ocf-url_manifest
func TestOcfUrlManifest(t *testing.T) {
	r := w3ctest.Load(t, "ocf-url_manifest")

	if !w3ctest.Contains(r.Identifier(), "ocf-url_manifest") {
		t.Errorf("expected identifier %q, got %v", "ocf-url_manifest", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "ocf-url_manifest") {
		t.Errorf("expected title %q, got %v", "ocf-url_manifest", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
