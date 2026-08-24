package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pub-cmt-jxl
func TestPubCmtJxl(t *testing.T) {
	r := w3ctest.Load(t, "pub-cmt-jxl")

	if !w3ctest.Contains(r.Identifier(), "pub-cmt-jxl") {
		t.Errorf("expected identifier %q, got %v", "pub-cmt-jxl", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pub-cmt-jxl") {
		t.Errorf("expected title %q, got %v", "pub-cmt-jxl", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
