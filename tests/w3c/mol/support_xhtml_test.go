package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/mol-support_xhtml
// Media overlay metadata (media-overlay manifest attribute and media:duration
// meta property) must be parsed.
func TestMolSupportXhtml(t *testing.T) {
	r := w3ctest.Load(t, "mol-support_xhtml")

	var hasOverlay bool
	for _, item := range r.CurrentSelectedPackage().Manifest.Items {
		if item.MediaOverlay != "" {
			hasOverlay = true
			break
		}
	}
	if !hasOverlay {
		t.Errorf("expected a manifest item with a media-overlay attribute")
	}

	var hasDuration bool
	for _, meta := range r.CurrentSelectedPackage().Metadata.Meta {
		if meta.Property == "media:duration" {
			hasDuration = true
			break
		}
	}
	if !hasDuration {
		t.Errorf("expected a media:duration meta property")
	}
}
