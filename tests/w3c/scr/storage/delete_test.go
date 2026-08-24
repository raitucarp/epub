package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/scr-storage-delete
func TestScrStorageDelete(t *testing.T) {
	r := w3ctest.Load(t, "scr-storage-delete")

	if !w3ctest.Contains(r.Identifier(), "scr-storage-delete") {
		t.Errorf("expected identifier %q, got %v", "scr-storage-delete", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "scr-storage-delete") {
		t.Errorf("expected title %q, got %v", "scr-storage-delete", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
