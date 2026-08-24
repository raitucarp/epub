package w3c

import (
	"testing"

	"github.com/raitucarp/epub/tests/w3c/w3ctest"
)

// Port of https://github.com/w3c/epub-tests/tests/sec-untrusted-consent_network
func TestSecUntrustedConsentNetwork(t *testing.T) {
	r := w3ctest.Load(t, "sec-untrusted-consent_network")

	if !w3ctest.Contains(r.Identifier(), "sec-untrusted-consent_network") {
		t.Errorf("expected identifier %q, got %v", "sec-untrusted-consent_network", r.Identifier())
	}

	if !w3ctest.Contains(r.Title(), "sec-untrusted-consent_network") {
		t.Errorf("expected title %q, got %v", "sec-untrusted-consent_network", r.Title())
	}

	if !w3ctest.Contains(r.Language(), "en") {
		t.Errorf("expected language %q, got %v", "en", r.Language())
	}
}
