package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pkg-flow-paginated
func TestLayPkgFlowPaginated(t *testing.T) {
	r := w3ctest.Load(t, "lay-pkg-flow-paginated")

	if !w3ctest.Contains(r.Identifier(), "lay-pkg-flow-paginated") {
		t.Errorf("expected identifier %q, got %v", "lay-pkg-flow-paginated", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-pkg-flow-paginated") {
		t.Errorf("expected title %q, got %v", "lay-pkg-flow-paginated", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
