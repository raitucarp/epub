package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/lay-pp-layout-pre-paginated-spreads
func TestLayPpLayoutPrePaginatedSpreads(t *testing.T) {
	r := w3ctest.Load(t, "lay-pp-layout-pre-paginated-spreads")

	if !w3ctest.Contains(r.Identifier(), "lay-pp-layout-pre-paginated-spreads") {
		t.Errorf("expected identifier %q, got %v", "lay-pp-layout-pre-paginated-spreads", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "lay-pp-layout-pre-paginated-spreads") {
		t.Errorf("expected title %q, got %v", "lay-pp-layout-pre-paginated-spreads", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
