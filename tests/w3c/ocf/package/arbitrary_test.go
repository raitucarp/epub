package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/ocf-package_arbitrary
// The package document must be located in an arbitrary directory, as referenced
// by container.xml, while unreferenced package documents are ignored.
func TestOcfPackageArbitrary(t *testing.T) {
	r := w3ctest.Load(t, "ocf-package_arbitrary")

	if got := r.CurrentSelectedPackagePath(); got != "FOO/BAR/package.opf" {
		t.Errorf("expected package path %q, got %q", "FOO/BAR/package.opf", got)
	}

	if !w3ctest.Contains(r.Title(), "ocf-package_arbitrary") {
		t.Errorf("expected title %q, got %v", "ocf-package_arbitrary", r.Title())
	}
}
