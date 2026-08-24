package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/ocf-zip-comp
func TestOcfZipComp(t *testing.T) {
	r := w3ctest.Load(t, "ocf-zip-comp")

	if !w3ctest.Contains(r.Identifier(), "ocf-zip-comp") {
		t.Errorf("expected identifier %q, got %v", "ocf-zip-comp", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "ocf-zip-comp") {
		t.Errorf("expected title %q, got %v", "ocf-zip-comp", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
