package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pp-spine-overrides_image-spine-pp
func TestLayPpSpineOverridesImageSpinePp(t *testing.T) {
	r := w3ctest.Load(t, "lay-pp-spine-overrides_image-spine-pp")

	if !w3ctest.Contains(r.Identifier(), "lay-pp-spine-overrides_image-spine-pp") {
		t.Errorf("expected identifier %q, got %v", "lay-pp-spine-overrides_image-spine-pp", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-pp-spine-overrides_image-spine-pp") {
		t.Errorf("expected title %q, got %v", "lay-pp-spine-overrides_image-spine-pp", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
