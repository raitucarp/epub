package epub

import (
	"bytes"
	"testing"

	"github.com/raitucarp/epub/ocf"
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

	if !bytes.Equal(obfuscated[:1040], bytes.Repeat([]byte{0xff}, 1040)) {
		t.Errorf("expected first 1040 bytes to be XORed")
	}
	if !bytes.Equal(obfuscated[1040:], content[1040:]) {
		t.Errorf("expected bytes beyond the obfuscation limit to be unchanged")
	}
}

func TestDeobfuscate_EmptyKeyOrContent(t *testing.T) {
	content := []byte("hello")
	if got := deobfuscate(content, nil, obfuscationKeyIDPFMax); !bytes.Equal(got, content) {
		t.Errorf("expected unchanged content for nil key")
	}
	if got := deobfuscate(nil, []byte("key"), obfuscationKeyIDPFMax); got != nil {
		t.Errorf("expected nil for empty content")
	}
}

func TestReader_ObfuscationKeyIDPF(t *testing.T) {
	r := &Reader{epub: &Epub{
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
	}}

	key := r.obfuscationKeyIDPF()
	if key == nil || len(key) != 20 {
		t.Fatalf("expected 20-byte SHA-1 key, got %v", key)
	}
}

func TestReader_ObfuscationKeyAdobe(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {
				UniqueIdentifier: "pub-id",
				Metadata: pkg.Metadata{
					Identifiers: []pkg.DCIdentifier{{ID: "pub-id", Value: "unique-id-value"}},
				},
			},
		},
	}}

	if key := r.obfuscationKeyAdobe(); string(key) != "unique-id-value" {
		t.Errorf("expected raw identifier key, got %q", key)
	}
}

func TestReader_ObfuscationKeyEmpty(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {},
		},
	}}

	if key := r.obfuscationKeyIDPF(); key != nil {
		t.Errorf("expected nil key for empty identifier, got %v", key)
	}
	if key := r.obfuscationKeyAdobe(); key != nil {
		t.Errorf("expected nil key for empty identifier, got %v", key)
	}
}

func TestReader_UniqueIdentifierValue(t *testing.T) {
	t.Run("resolves unique-identifier", func(t *testing.T) {
		r := &Reader{epub: &Epub{
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": {
					UniqueIdentifier: "pub-id",
					Metadata: pkg.Metadata{
						Identifiers: []pkg.DCIdentifier{{ID: "pub-id", Value: "  urn:x  "}},
					},
				},
			},
		}}

		if got := r.uniqueIdentifierValue(); got != "urn:x" {
			t.Errorf("expected %q, got %q", "urn:x", got)
		}
	})

	t.Run("falls back to first non-empty", func(t *testing.T) {
		r := &Reader{epub: &Epub{
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": {
					Metadata: pkg.Metadata{
						Identifiers: []pkg.DCIdentifier{{ID: "id-1", Value: "urn:first"}},
					},
				},
			},
		}}

		if got := r.uniqueIdentifierValue(); got != "urn:first" {
			t.Errorf("expected %q, got %q", "urn:first", got)
		}
	})

	t.Run("no identifiers", func(t *testing.T) {
		r := &Reader{epub: &Epub{
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": {},
			},
		}}

		if got := r.uniqueIdentifierValue(); got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})
}

func TestReader_DeobfuscateFonts(t *testing.T) {
	original := []byte("font-bytes-here")
	key := make([]byte, 20)
	key[0] = 0x01

	idpfKey := key
	obfuscated := deobfuscate(original, idpfKey, obfuscationKeyIDPFMax)

	z := ocf.NewOCFZipContainer()
	z.Encryption().EncryptedData = []ocf.EncryptedData{
		{
			EncryptionMethod: &ocf.EncryptionMethod{Algorithm: fontObfuscationIDPF},
			CipherData:       ocf.CipherData{CipherReference: &ocf.CipherReference{URI: "fonts/font.ttf"}},
		},
	}

	r := &Reader{epub: &Epub{
		rendition:    "default",
		zipContainer: z,
		packagePubs: map[string]*pkg.Package{
			"default": {
				UniqueIdentifier: "pub-id",
				Metadata: pkg.Metadata{
					Identifiers: []pkg.DCIdentifier{{ID: "pub-id", Value: "identifier"}},
				},
			},
		},
		resources: []PublicationResource{
			{ID: "font", Filepath: "fonts/font.ttf", MIMEType: "font/ttf", Content: obfuscated},
		},
	}}

	// The obfuscation key is derived from the identifier at deobfuscation time,
	// so obfuscate with the same derived key to verify round-trip.
	expectedKey := r.obfuscationKeyIDPF()
	reobfuscated := deobfuscate(original, expectedKey, obfuscationKeyIDPFMax)

	r.epub.resources[0].Content = reobfuscated
	r.deobfuscateFonts()

	if !bytes.Equal(r.epub.resources[0].Content, original) {
		t.Errorf("expected font to be de-obfuscated, got %q", r.epub.resources[0].Content)
	}
}

func TestReader_DeobfuscateFonts_NoEncryption(t *testing.T) {
	z := ocf.NewOCFZipContainer()
	r := &Reader{epub: &Epub{
		rendition:    "default",
		zipContainer: z,
		packagePubs:  map[string]*pkg.Package{"default": {}},
		resources:    []PublicationResource{{ID: "font", Filepath: "f.ttf", Content: []byte("data")}},
	}}

	r.deobfuscateFonts()
	if string(r.epub.resources[0].Content) != "data" {
		t.Error("expected content unchanged without encryption")
	}
}
