package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-linked-records
// A metadata record link (rel="record") must be parsed with its properties.
func TestPkgLinkedRecords(t *testing.T) {
	r := w3ctest.Load(t, "pkg-linked-records")

	var foundRecord bool
	for _, link := range r.CurrentSelectedPackage().Metadata.Links {
		if link.Rel == "record" {
			foundRecord = true
			if link.Href == "" {
				t.Errorf("expected record link to have an href")
			}
			if link.Properties == "" {
				t.Errorf("expected record link to carry properties")
			}
		}
	}

	if !foundRecord {
		t.Errorf("expected a record link in metadata")
	}
}
