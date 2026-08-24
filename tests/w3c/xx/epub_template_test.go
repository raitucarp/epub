package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/xx-epub-template
func TestXxEpubTemplate(t *testing.T) {
	r := w3ctest.Load(t, "xx-epub-template")

	if !w3ctest.Contains(r.Identifier(), "TODO: test-id") {
		t.Errorf("expected identifier %q, got %v", "TODO: test-id", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "TODO: test-id") {
		t.Errorf("expected title %q, got %v", "TODO: test-id", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
