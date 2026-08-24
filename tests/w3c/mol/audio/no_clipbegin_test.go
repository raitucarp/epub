package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/mol-audio-no-clipbegin
func TestMolAudioNoClipbegin(t *testing.T) {
	r := w3ctest.Load(t, "mol-audio-no-clipbegin")

	if !w3ctest.Contains(r.Identifier(), "mol-audio-no-clipbegin") {
		t.Errorf("expected identifier %q, got %v", "mol-audio-no-clipbegin", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "mol-audio-no-clipbegin") {
		t.Errorf("expected title %q, got %v", "mol-audio-no-clipbegin", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
