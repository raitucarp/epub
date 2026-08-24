package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/pkg-spine-nonlinear-activation
func TestPkgSpineNonlinearActivation(t *testing.T) {
	r := w3ctest.Load(t, "pkg-spine-nonlinear-activation")

	if !w3ctest.Contains(r.Identifier(), "pkg-spine-nonlinear-activation") {
		t.Errorf("expected identifier %q, got %v", "pkg-spine-nonlinear-activation", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "pkg-spine-nonlinear-activation") {
		t.Errorf("expected title %q, got %v", "pkg-spine-nonlinear-activation", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
