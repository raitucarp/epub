package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-data-urls_top-level-content
func TestPubDataUrlsTopLevelContent(t *testing.T) {
	r := w3ctest.Load(t, "pub-data-urls_top-level-content")

	if !w3ctest.Contains(r.Identifier(), "pub-data-urls_top-level-content") {
		t.Errorf("expected identifier %q, got %v", "pub-data-urls_top-level-content", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pub-data-urls_top-level-content") {
		t.Errorf("expected title %q, got %v", "pub-data-urls_top-level-content", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
