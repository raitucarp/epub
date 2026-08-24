package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/fxl-spine-overrides_behave-as-global-bis
func TestFxlSpineOverridesBehaveAsGlobalBis(t *testing.T) {
	r := w3ctest.Load(t, "fxl-spine-overrides_behave-as-global-bis")

	if !w3ctest.Contains(r.Identifier(), "fxl-spine-overrides_behave-as-global-bis") {
		t.Errorf("expected identifier %q, got %v", "fxl-spine-overrides_behave-as-global-bis", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "fxl-spine-overrides_behave-as-global-bis") {
		t.Errorf("expected title %q, got %v", "fxl-spine-overrides_behave-as-global-bis", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
