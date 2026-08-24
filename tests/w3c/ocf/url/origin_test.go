package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/ocf-url_origin
func TestOcfUrlOrigin(t *testing.T) {
	r := w3ctest.Load(t, "ocf-url_origin")

	if !w3ctest.Contains(r.Identifier(), "ocf-url_origin") {
		t.Errorf("expected identifier %q, got %v", "ocf-url_origin", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "ocf-url_origin") {
		t.Errorf("expected title %q, got %v", "ocf-url_origin", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
