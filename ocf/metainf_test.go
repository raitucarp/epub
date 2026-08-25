package ocf

import (
	"strings"
	"testing"
)

const encryptionXML = `<?xml version="1.0" encoding="UTF-8"?>
<encryption xmlns="urn:oasis:names:tc:opendocument:xmlns:container" xmlns:enc="http://www.w3.org/2001/04/xmlenc#">
  <enc:EncryptedData>
    <enc:EncryptionMethod Algorithm="http://www.idpf.org/2008/embedding"/>
    <enc:CipherData>
      <enc:CipherReference URI="EPUB/fonts/font.ttf"/>
    </enc:CipherData>
  </enc:EncryptedData>
</encryption>`

const manifestXML = `<?xml version="1.0" encoding="UTF-8"?>
<manifest xmlns="urn:oasis:names:tc:opendocument:xmlns:manifest:1.0" version="1.0">
  <file-entry full-path="EPUB/package.opf" media-type="application/oebps-package+xml"/>
</manifest>`

const metadataXML = `<?xml version="1.0" encoding="UTF-8"?>
<metadata xmlns="http://www.idpf.org/2013/metadata">
  <rendition>default</rendition>
</metadata>`

const rightsXML = `<?xml version="1.0" encoding="UTF-8"?>
<rights>
  <license>all rights reserved</license>
</rights>`

const signaturesXML = `<?xml version="1.0" encoding="UTF-8"?>
<signatures xmlns="urn:oasis:names:tc:opendocument:xmlns:container" xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
  <ds:Signature Id="sig1">
    <ds:SignedInfo>
      <ds:CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
      <ds:SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/>
      <ds:Reference URI="EPUB/package.opf">
        <ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
        <ds:DigestValue>abc123</ds:DigestValue>
      </ds:Reference>
    </ds:SignedInfo>
    <ds:SignatureValue>signed-bytes</ds:SignatureValue>
  </ds:Signature>
</signatures>`

func allReservedFiles() map[string][]byte {
	return map[string][]byte{
		"mimetype":                []byte(MimeType),
		"META-INF/container.xml":  []byte(validContainerXML),
		"META-INF/encryption.xml": []byte(encryptionXML),
		"META-INF/manifest.xml":   []byte(manifestXML),
		"META-INF/metadata.xml":   []byte(metadataXML),
		"META-INF/rights.xml":     []byte(rightsXML),
		"META-INF/signatures.xml": []byte(signaturesXML),
	}
}

func TestParseAllMetaInfFiles(t *testing.T) {
	container, err := NewReader(buildZip(t, allReservedFiles()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(container.Container().RootFiles.RootFile) != 1 {
		t.Errorf("expected 1 rootfile, got %d", len(container.Container().RootFiles.RootFile))
	}

	enc := container.Encryption()
	if len(enc.EncryptedData) != 1 {
		t.Errorf("expected 1 encrypted data, got %d", len(enc.EncryptedData))
	} else if enc.EncryptedData[0].EncryptionMethod == nil {
		t.Error("expected encryption method to be parsed")
	}

	manifest := container.Manifest()
	if len(manifest.FileEntry) != 1 {
		t.Errorf("expected 1 manifest file-entry, got %d", len(manifest.FileEntry))
	}

	metadata := container.Metadata()
	if metadata.Content == "" {
		t.Error("expected metadata content to be captured")
	}

	rights := container.Rights()
	if rights.Content == "" {
		t.Error("expected rights content to be captured")
	}

	sigs := container.Signatures()
	if len(sigs.Signature) != 1 {
		t.Errorf("expected 1 signature, got %d", len(sigs.Signature))
	} else {
		sig := sigs.Signature[0]
		if sig.ID != "sig1" {
			t.Errorf("expected signature id %q, got %q", "sig1", sig.ID)
		}
		if sig.SignatureValue != "signed-bytes" {
			t.Errorf("expected signature value, got %q", sig.SignatureValue)
		}
		if len(sig.SignedInfo.Reference) != 1 {
			t.Errorf("expected 1 reference, got %d", len(sig.SignedInfo.Reference))
		} else if sig.SignedInfo.Reference[0].DigestValue != "abc123" {
			t.Errorf("expected digest value, got %q", sig.SignedInfo.Reference[0].DigestValue)
		}
	}
}

func TestParseAllMetaInfFiles_MissingRequired(t *testing.T) {
	z := NewOCFZipContainer()
	z.files = map[string][]byte{
		"mimetype": []byte(MimeType),
	}

	if err := z.parseAllMetaInfFiles(); err == nil {
		t.Error("expected error for missing container.xml, got nil")
	}
}

func TestParseAllMetaInfFiles_MalformedContainer(t *testing.T) {
	z := NewOCFZipContainer()
	z.files = map[string][]byte{
		"mimetype":               []byte(MimeType),
		"META-INF/container.xml": []byte("<container><unclosed>"),
	}

	if err := z.parseAllMetaInfFiles(); err == nil {
		t.Error("expected error for malformed container.xml, got nil")
	}
}

func TestMetaInfReservedFile_Unmarshal_Error(t *testing.T) {
	var f metaInfReservedFile = containerFile
	var v Container
	if err := f.Unmarshal([]byte("<container><broken>"), &v); err == nil {
		t.Error("expected error for malformed xml, got nil")
	}
}

func TestParseAllMetaInfFiles_IgnoresNonReservedFiles(t *testing.T) {
	z := NewOCFZipContainer()
	z.files = map[string][]byte{
		"mimetype":                  []byte(MimeType),
		"META-INF/container.xml":    []byte(validContainerXML),
		"META-INF/not-reserved.xml": []byte("<foo/>"),
		"EPUB/package.opf":          []byte("<package/>"),
	}

	if err := z.parseAllMetaInfFiles(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(z.Container().RootFiles.RootFile) != 1 {
		t.Errorf("expected container.xml to be parsed, got %d rootfiles", len(z.Container().RootFiles.RootFile))
	}
}

func TestParseMetadata_InnerXML(t *testing.T) {
	z := NewOCFZipContainer()
	if err := z.metaInf.parseMetadata([]byte(metadataXML)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(z.metaInf.metadata.Content, "rendition") {
		t.Errorf("expected inner xml to contain 'rendition', got %v", z.metaInf.metadata.Content)
	}
}

func TestParseRights_InnerXML(t *testing.T) {
	z := NewOCFZipContainer()
	if err := z.metaInf.parseRights([]byte(rightsXML)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(z.metaInf.rights.Content, "license") {
		t.Errorf("expected inner xml to contain 'license', got %v", z.metaInf.rights.Content)
	}
}
