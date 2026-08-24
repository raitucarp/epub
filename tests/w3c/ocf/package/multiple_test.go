package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/ocf-package_multiple
// Multiple package documents (renditions) must all be listed, with the first
// one selectable as the default rendition.
func TestOcfPackageMultiple(t *testing.T) {
	r := w3ctest.Load(t, "ocf-package_multiple")

	renditions := r.ListRenditions()
	if len(renditions) != 3 {
		t.Fatalf("expected 3 renditions, got %v", renditions)
	}
	if renditions[0] != "default" {
		t.Errorf("expected default rendition first, got %v", renditions)
	}

	for _, rendition := range renditions {
		r.SelectPackageRendition(rendition)
		if r.CurrentSelectedPackage() == nil {
			t.Errorf("expected package for rendition %q", rendition)
		}
	}
}
