package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/fxl-spine-overrides_behave-as-global
func TestFxlSpineOverridesBehaveAsGlobal(t *testing.T) {
	r := w3ctest.Load(t, "fxl-spine-overrides_behave-as-global")

	if !w3ctest.Contains(r.Identifier(), "fxl-spine-overrides_behave-as-global") {
		t.Errorf("expected identifier %q, got %v", "fxl-spine-overrides_behave-as-global", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "fxl-spine-overrides_behave-as-global") {
		t.Errorf("expected title %q, got %v", "fxl-spine-overrides_behave-as-global", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
