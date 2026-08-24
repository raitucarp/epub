package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-collections-unknown
// A collection with an unknown role must be parsed without error.
func TestPkgCollectionsUnknown(t *testing.T) {
	r := w3ctest.Load(t, "pkg-collections-unknown")

	collections := r.CurrentSelectedPackage().Collections
	if len(collections) != 1 {
		t.Fatalf("expected 1 collection, got %d", len(collections))
	}

	if collections[0].Role != "foo" {
		t.Errorf("expected role %q, got %q", "foo", collections[0].Role)
	}

	if collections[0].Metadata == nil || len(collections[0].Metadata.OptionalDC) != 1 {
		t.Errorf("expected collection metadata to capture dc:title, got %+v", collections[0].Metadata)
	}

	if len(collections[0].Links) != 1 {
		t.Errorf("expected 1 collection link, got %d", len(collections[0].Links))
	}
}
