package pkg

import (
	"encoding/xml"
	"strings"
	"testing"
)

const packageFixture = `<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" xmlns:dc="http://purl.org/dc/elements/1.1/" version="3.0" unique-identifier="pub-id">
  <metadata>
    <dc:identifier id="pub-id">urn:isbn:1234567890</dc:identifier>
    <dc:title id="title">My Book</dc:title>
    <dc:language>en</dc:language>
    <dc:creator id="creator">Jane Doe</dc:creator>
    <meta refines="#title" property="title-type">main</meta>
  </metadata>
  <manifest>
    <item id="nav" properties="nav" href="nav.xhtml" media-type="application/xhtml+xml"/>
    <item id="content" href="content.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
  <spine>
    <itemref idref="content"/>
  </spine>
</package>`

func TestPackage_UnmarshalNamespaces(t *testing.T) {
	var p Package
	if err := xml.Unmarshal([]byte(packageFixture), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(p.Metadata.Identifiers) != 1 {
		t.Fatalf("expected 1 identifier, got %d", len(p.Metadata.Identifiers))
	}
	if p.Metadata.Identifiers[0].Value != "urn:isbn:1234567890" {
		t.Errorf("unexpected identifier value %q", p.Metadata.Identifiers[0].Value)
	}
	if p.Metadata.Identifiers[0].ID != "pub-id" {
		t.Errorf("unexpected identifier id %q", p.Metadata.Identifiers[0].ID)
	}

	if len(p.Metadata.Titles) != 1 || p.Metadata.Titles[0].Value != "My Book" {
		t.Errorf("expected title to be parsed, got %+v", p.Metadata.Titles)
	}

	if len(p.Metadata.Languages) != 1 || p.Metadata.Languages[0].Value != "en" {
		t.Errorf("expected language to be parsed, got %+v", p.Metadata.Languages)
	}

	if len(p.Metadata.OptionalDC) != 1 {
		t.Fatalf("expected 1 optional dc element (dc:creator), got %d", len(p.Metadata.OptionalDC))
	}
	creator := p.Metadata.OptionalDC[0]
	if creator.XMLName.Local != "creator" || creator.Value != "Jane Doe" {
		t.Errorf("unexpected creator %+v", creator)
	}

	if len(p.Metadata.Meta) != 1 || p.Metadata.Meta[0].Property != "title-type" {
		t.Errorf("unexpected meta %+v", p.Metadata.Meta)
	}

	if len(p.Manifest.Items) != 2 {
		t.Errorf("expected 2 manifest items, got %d", len(p.Manifest.Items))
	}

	if len(p.Spine.ItemRefs) != 1 || p.Spine.ItemRefs[0].IDRef != "content" {
		t.Errorf("unexpected spine %+v", p.Spine.ItemRefs)
	}
}

func TestPackage_MarshalNamespacesRoundTrip(t *testing.T) {
	p := Package{
		Version:          "3.0",
		UniqueIdentifier: "pub-id",
	}
	p.Metadata.Identifiers = append(p.Metadata.Identifiers, DCIdentifier{ID: "pub-id", Value: "urn:isbn:1234567890"})
	p.Metadata.Titles = append(p.Metadata.Titles, DCTitle{ID: "title", Value: "My Book"})
	p.Metadata.Languages = append(p.Metadata.Languages, DCLanguage{Value: "en"})
	p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, DCOptional{
		XMLName: xml.Name{Space: NamespaceDC, Local: "creator"},
		Value:   "Jane Doe",
	})

	out, err := xml.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(out)
	if !strings.Contains(s, "xmlns=\"http://www.idpf.org/2007/opf\"") {
		t.Errorf("expected package namespace declaration, got:\n%s", s)
	}
	if !strings.Contains(s, "http://purl.org/dc/elements/1.1/") {
		t.Errorf("expected Dublin Core namespace declaration, got:\n%s", s)
	}

	var parsed Package
	if err := xml.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("unexpected error unmarshalling output: %v", err)
	}
	if len(parsed.Metadata.Identifiers) != 1 || parsed.Metadata.Identifiers[0].Value != "urn:isbn:1234567890" {
		t.Errorf("identifier did not round-trip, got %+v", parsed.Metadata.Identifiers)
	}
	if len(parsed.Metadata.Titles) != 1 || parsed.Metadata.Titles[0].Value != "My Book" {
		t.Errorf("title did not round-trip, got %+v", parsed.Metadata.Titles)
	}
	if len(parsed.Metadata.OptionalDC) != 1 || parsed.Metadata.OptionalDC[0].Value != "Jane Doe" {
		t.Errorf("creator did not round-trip, got %+v", parsed.Metadata.OptionalDC)
	}
}

func TestCollection_Unmarshal(t *testing.T) {
	fixture := `<collection xmlns="http://www.idpf.org/2007/opf" role="series" xmlns:dc="http://purl.org/dc/elements/1.1/">
  <metadata>
    <dc:title>My Series</dc:title>
  </metadata>
  <link href="content.xhtml"/>
</collection>`

	var c Collection
	if err := xml.Unmarshal([]byte(fixture), &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Role != "series" {
		t.Errorf("expected role series, got %q", c.Role)
	}
	if c.Metadata == nil {
		t.Fatalf("expected collection metadata, got nil")
	}
	if len(c.Metadata.OptionalDC) != 1 || c.Metadata.OptionalDC[0].Value != "My Series" {
		t.Errorf("expected dc:title to be captured, got %+v", c.Metadata.OptionalDC)
	}
	if len(c.Links) != 1 || c.Links[0].Href != "content.xhtml" {
		t.Errorf("expected link, got %+v", c.Links)
	}
}

func TestPackage_XmlLangRoundTrip(t *testing.T) {
	fixture := `<package xmlns="http://www.idpf.org/2007/opf" version="3.0" xml:lang="en" unique-identifier="pub-id"></package>`

	var p Package
	if err := xml.Unmarshal([]byte(fixture), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Lang != "en" {
		t.Errorf("expected xml:lang to be parsed as 'en', got %q", p.Lang)
	}

	out, err := xml.Marshal(p)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, `xml:lang="en"`) {
		t.Errorf("expected xml:lang attribute in output, got %s", s)
	}
	if strings.Contains(s, "_xml") {
		t.Errorf("output must not contain a mangled xml prefix, got %s", s)
	}
}

func TestGuideReferenceType_MarshalText(t *testing.T) {
	values := []GuideReferenceType{
		GuideRefCover, GuideRefTitlePage, GuideRefToc, GuideRefIndex,
		GuideRefGlossary, GuideRefAcknowledgements, GuideRefBibliography,
		GuideRefColophon, GuideRefCopyrightPage, GuideRefDedication,
		GuideRefEpigraph, GuideRefForeword, GuideRefLoi, GuideRefLot,
		GuideRefNotes, GuideRefPreface, GuideRefText,
	}

	for _, v := range values {
		t.Run(string(v), func(t *testing.T) {
			out, err := v.MarshalText()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(out) != string(v) {
				t.Errorf("expected %q, got %q", string(v), string(out))
			}
		})
	}
}

func TestGuideReferenceType_UnmarshalText(t *testing.T) {
	known := []GuideReferenceType{
		GuideRefCover, GuideRefTitlePage, GuideRefToc, GuideRefIndex,
		GuideRefGlossary, GuideRefAcknowledgements, GuideRefBibliography,
		GuideRefColophon, GuideRefCopyrightPage, GuideRefDedication,
		GuideRefEpigraph, GuideRefForeword, GuideRefLoi, GuideRefLot,
		GuideRefNotes, GuideRefPreface, GuideRefText,
	}

	for _, v := range known {
		t.Run(string(v), func(t *testing.T) {
			var got GuideReferenceType
			if err := got.UnmarshalText([]byte(v)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != v {
				t.Errorf("expected %q, got %q", v, got)
			}
		})
	}

	t.Run("unknown value preserved", func(t *testing.T) {
		var got GuideReferenceType
		if err := got.UnmarshalText([]byte("custom-type")); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != GuideReferenceType("custom-type") {
			t.Errorf("expected custom-type to be preserved, got %q", got)
		}
	})
}
