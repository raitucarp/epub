package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/ocf-zip-mult
func TestOcfZipMult(t *testing.T) {
	r := w3ctest.Load(t, "ocf-zip-mult")

	if !w3ctest.Contains(r.Identifier(), "ocf-zip-mult") {
		t.Errorf("expected identifier %q, got %v", "ocf-zip-mult", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "ocf-zip-mult") {
		t.Errorf("expected title %q, got %v", "ocf-zip-mult", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
