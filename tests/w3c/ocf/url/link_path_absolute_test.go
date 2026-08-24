package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/ocf-url_link-path-absolute
func TestOcfUrlLinkPathAbsolute(t *testing.T) {
	r := w3ctest.Load(t, "ocf-url_link-path-absolute")

	if !w3ctest.Contains(r.Identifier(), "ocf-url_link-path-absolute") {
		t.Errorf("expected identifier %q, got %v", "ocf-url_link-path-absolute", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "ocf-url_link-path-absolute") {
		t.Errorf("expected title %q, got %v", "ocf-url_link-path-absolute", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
