package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/mol-timing-synchronization_multiple_audio
func TestMolTimingSynchronizationMultipleAudio(t *testing.T) {
	r := w3ctest.Load(t, "mol-timing-synchronization_multiple_audio")

	if !w3ctest.Contains(r.Identifier(), "mol-timing-synchronization_multiple_audio") {
		t.Errorf("expected identifier %q, got %v", "mol-timing-synchronization_multiple_audio", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "mol-timing-synchronization_multiple_audio") {
		t.Errorf("expected title %q, got %v", "mol-timing-synchronization_multiple_audio", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
