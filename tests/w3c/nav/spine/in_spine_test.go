package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/nav-spine_in-spine
// The toc nav must be parsed from the EPUB navigation document.
func TestNavSpineInSpine(t *testing.T) {
	r := w3ctest.Load(t, "nav-spine_in-spine")

	toc, err := r.TableOfContents()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(toc.Items) != 2 {
		t.Errorf("expected 2 toc items, got %d", len(toc.Items))
	}
}
