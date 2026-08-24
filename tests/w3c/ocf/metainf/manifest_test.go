package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/ocf-metainf-manifest
// An ancillary META-INF/manifest.xml (containing an extra spine item) must not
// affect the spine defined by the package document.
func TestOcfMetainfManifest(t *testing.T) {
	r := w3ctest.Load(t, "ocf-metainf-manifest")

	if got := len(r.Spine()); got != 1 {
		t.Errorf("expected 1 spine item, got %d", got)
	}

	if !w3ctest.Contains(r.Identifier(), "ocf-metainf-manifest") {
		t.Errorf("expected identifier %q, got %v", "ocf-metainf-manifest", r.Identifier())
	}
}
