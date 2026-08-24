package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/mol-timing-synchronization
func TestMolTimingSynchronization(t *testing.T) {
	r := w3ctest.Load(t, "mol-timing-synchronization")

	if !w3ctest.Contains(r.Identifier(), "mol-timing-synchronization") {
		t.Errorf("expected identifier %q, got %v", "mol-timing-synchronization", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "mol-timing-synchronization") {
		t.Errorf("expected title %q, got %v", "mol-timing-synchronization", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
