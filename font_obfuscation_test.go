package epub

import (
	"bytes"
	"testing"

	"github.com/raitucarp/epub/pkg"
)

func TestDeobfuscate_RoundTrip(t *testing.T) {
	content := make([]byte, 2048)
	for i := range content {
		content[i] = byte(i % 256)
	}

	key := []byte{0xde, 0xad, 0xbe, 0xef}

	obfuscated := deobfuscate(content, key, obfuscationKeyIDPFMax)
	if bytes.Equal(obfuscated, content) {
		t.Fatalf("expected obfuscated content to differ from original")
	}

	restored := deobfuscate(obfuscated, key, obfuscationKeyIDPFMax)
	if !bytes.Equal(restored, content) {
		t.Errorf("expected restored content to equal original")
	}
}

func TestDeobfuscate_Limit(t *testing.T) {
	content := bytes.Repeat([]byte{0x55}, 3000)
	key := []byte{0xaa}

	obfuscated := deobfuscate(content, key, obfuscationKeyIDPFMax)

	// First 1040 bytes must differ.
	if !bytes.Equal(obfuscated[:1040], bytes.Repeat([]byte{0xff}, 1040)) {
		t.Errorf("expected first 1040 bytes to be XORed")
	}

	// Bytes after 1040 must be unchanged.
	if !bytes.Equal(obfuscated[1040:], content[1040:]) {
		t.Errorf("expected bytes beyond the obfuscation limit to be unchanged")
	}
}

func TestReader_ObfuscationKeyIDPF(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": {
					UniqueIdentifier: "pub-id",
					Metadata: pkg.Metadata{
						Identifiers: []pkg.DCIdentifier{
							{ID: "other", Value: "urn:other"},
							{ID: "pub-id", Value: "ocf-font_obfuscation"},
						},
					},
				},
			},
		},
	}

	key := r.obfuscationKeyIDPF()
	if key == nil || len(key) != 20 {
		t.Fatalf("expected 20-byte SHA-1 key, got %v", key)
	}
}

func TestReader_ObfuscationKeyEmpty(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": {},
			},
		},
	}

	if key := r.obfuscationKeyIDPF(); key != nil {
		t.Errorf("expected nil key for empty identifier, got %v", key)
	}
}
