package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-cmt-mp4
func TestPubCmtMp4(t *testing.T) {
	r := w3ctest.Load(t, "pub-cmt-mp4")

	if !w3ctest.Contains(r.Identifier(), "pub-cmt-mp4") {
		t.Errorf("expected identifier %q, got %v", "pub-cmt-mp4", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pub-cmt-mp4") {
		t.Errorf("expected title %q, got %v", "pub-cmt-mp4", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
