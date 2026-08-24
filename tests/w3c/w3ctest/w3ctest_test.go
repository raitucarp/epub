package w3ctest

import (
	"os"
	"testing"
)

// TestSyncBinaries restores the large binary fixtures that are not committed by
// cloning the epub-tests repository. It is an opt-in maintenance test, enabled
// via the SYNC_EPUB_TESTS environment variable.
func TestSyncBinaries(t *testing.T) {
	if os.Getenv("SYNC_EPUB_TESTS") == "" {
		t.Skip("set SYNC_EPUB_TESTS=1 to sync epub-tests binary resources")
	}
	SyncBinaries(t)
}
