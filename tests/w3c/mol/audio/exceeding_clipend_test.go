package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/mol-audio-exceeding-clipend
func TestMolAudioExceedingClipend(t *testing.T) {
	r := w3ctest.Load(t, "mol-audio-exceeding-clipend")

	if !w3ctest.Contains(r.Identifier(), "mol-audio-exceeding-clipend") {
		t.Errorf("expected identifier %q, got %v", "mol-audio-exceeding-clipend", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "mol-audio-exceeding-clipend") {
		t.Errorf("expected title %q, got %v", "mol-audio-exceeding-clipend", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
