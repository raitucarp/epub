package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/ocf-metainf-inc
func TestOcfMetainfInc(t *testing.T) {
	r := w3ctest.Load(t, "ocf-metainf-inc")

	if !w3ctest.Contains(r.Identifier(), "ocf-metainf-inc") {
		t.Errorf("expected identifier %q, got %v", "ocf-metainf-inc", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "ocf-metainf-inc") {
		t.Errorf("expected title %q, got %v", "ocf-metainf-inc", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
