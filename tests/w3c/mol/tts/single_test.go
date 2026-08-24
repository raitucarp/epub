package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/mol-tts_single
func TestMolTtsSingle(t *testing.T) {
	r := w3ctest.Load(t, "mol-tts_single")

	if !w3ctest.Contains(r.Identifier(), "mol-tts_single") {
		t.Errorf("expected identifier %q, got %v", "mol-tts_single", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "mol-tts_single") {
		t.Errorf("expected title %q, got %v", "mol-tts_single", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
