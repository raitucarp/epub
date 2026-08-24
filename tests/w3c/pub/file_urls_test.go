package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-file-urls
func TestPubFileUrls(t *testing.T) {
	r := w3ctest.Load(t, "pub-file-urls")

	if !w3ctest.Contains(r.Identifier(), "pub-file-urls") {
		t.Errorf("expected identifier %q, got %v", "pub-file-urls", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pub-file-urls") {
		t.Errorf("expected title %q, got %v", "pub-file-urls", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
