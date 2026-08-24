package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pkg-flow-scrolled-continuous
// The rendition:flow package rendering property must be parsed.
func TestLayPkgFlowScrolledContinuous(t *testing.T) {
	r := w3ctest.Load(t, "lay-pkg-flow-scrolled-continuous")

	var flow string
	for _, meta := range r.CurrentSelectedPackage().Metadata.Meta {
		if meta.Property == "rendition:flow" {
			flow = meta.Value
			break
		}
	}

	if flow != "scrolled-continuous" {
		t.Errorf("expected rendition:flow %q, got %q", "scrolled-continuous", flow)
	}
}
