package epub

import (
	"crypto/sha1"
	"strings"
)

// Font obfuscation algorithms defined by the OCF specification.
const (
	fontObfuscationIDPF    = "http://www.idpf.org/2008/embedding"
	fontObfuscationAdobe   = "http://ns.adobe.com/pdf/enc#RC"
	obfuscationKeyIDPFMax  = 1040
	obfuscationKeyAdobeMax = 1024
)

// deobfuscateFonts restores fonts that have been obfuscated as described in
// the OCF specification. The obfuscation key is derived from the publication's
// unique identifier, so the de-obfuscation is scoped to the selected package.
func (r *Reader) deobfuscateFonts() {
	encryption := r.epub.zipContainer.Encryption()
	if encryption == nil || len(encryption.EncryptedData) == 0 {
		return
	}

	idpfKey := r.obfuscationKeyIDPF()
	adobeKey := r.obfuscationKeyAdobe()

	for _, data := range encryption.EncryptedData {
		if data.EncryptionMethod == nil || data.CipherData.CipherReference == nil {
			continue
		}

		var key []byte
		var maxLen int
		switch data.EncryptionMethod.Algorithm {
		case fontObfuscationIDPF:
			key = idpfKey
			maxLen = obfuscationKeyIDPFMax
		case fontObfuscationAdobe:
			key = adobeKey
			maxLen = obfuscationKeyAdobeMax
		default:
			continue
		}

		if len(key) == 0 {
			continue
		}

		uri := data.CipherData.CipherReference.URI
		for i := range r.epub.resources {
			if r.epub.resources[i].Filepath == uri {
				r.epub.resources[i].Content = deobfuscate(r.epub.resources[i].Content, key, maxLen)
			}
		}
	}
}

// obfuscationKeyIDPF returns the SHA-1 digest of the unique identifier, used
// as the key for the IDPF font obfuscation algorithm.
func (r *Reader) obfuscationKeyIDPF() []byte {
	identifier := r.uniqueIdentifierValue()
	if identifier == "" {
		return nil
	}

	sum := sha1.Sum([]byte(identifier))
	return sum[:]
}

// obfuscationKeyAdobe returns the raw unique identifier bytes, used as the
// key for the legacy Adobe font obfuscation algorithm.
func (r *Reader) obfuscationKeyAdobe() []byte {
	identifier := r.uniqueIdentifierValue()
	if identifier == "" {
		return nil
	}

	return []byte(identifier)
}

func (r *Reader) uniqueIdentifierValue() string {
	metadata := r.CurrentSelectedPackage().Metadata

	for _, id := range metadata.Identifiers {
		if id.ID == r.CurrentSelectedPackage().UniqueIdentifier {
			return strings.TrimSpace(id.Value)
		}
	}

	for _, id := range metadata.Identifiers {
		if id.Value != "" {
			return strings.TrimSpace(id.Value)
		}
	}

	return ""
}

func deobfuscate(content []byte, key []byte, maxLen int) []byte {
	if len(key) == 0 || len(content) == 0 {
		return content
	}

	out := make([]byte, len(content))
	copy(out, content)

	limit := len(out)
	if limit > maxLen {
		limit = maxLen
	}

	for i := 0; i < limit; i++ {
		out[i] = content[i] ^ key[i%len(key)]
	}

	return out
}
