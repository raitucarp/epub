package pkg

import (
	"encoding/xml"
	"fmt"
)

// Package represents the root package element
type Package struct {
	XMLName          xml.Name     `xml:"http://www.idpf.org/2007/opf package"`
	Dir              string       `xml:"dir,attr,omitempty"`
	ID               string       `xml:"id,attr,omitempty"`
	Prefix           string       `xml:"prefix,attr,omitempty"`
	Lang             string       `xml:"xml lang,attr,omitempty"`
	UniqueIdentifier string       `xml:"unique-identifier,attr"`
	Version          string       `xml:"version,attr"`
	Metadata         Metadata     `xml:"metadata"`
	Manifest         Manifest     `xml:"manifest"`
	Spine            Spine        `xml:"spine"`
	Guide            *Guide       `xml:"guide,omitempty"`
	Bindings         *Bindings    `xml:"bindings,omitempty"`
	Collections      []Collection `xml:"collection,omitempty"`
}

// Metadata represents the metadata section
type Metadata struct {
	XMLName     xml.Name       `xml:"metadata"`
	Identifiers []DCIdentifier `xml:"http://purl.org/dc/elements/1.1/ identifier"`
	Titles      []DCTitle      `xml:"http://purl.org/dc/elements/1.1/ title"`
	Languages   []DCLanguage   `xml:"http://purl.org/dc/elements/1.1/ language"`
	OptionalDC  []DCOptional   `xml:",any"`
	Meta        []Meta         `xml:"meta"`
	Links       []Link         `xml:"link,omitempty"`
}

// DCIdentifier represents dc:identifier element
type DCIdentifier struct {
	XMLName xml.Name `xml:"http://purl.org/dc/elements/1.1/ identifier"`
	ID      string   `xml:"id,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

// DCTitle represents dc:title element
type DCTitle struct {
	XMLName xml.Name `xml:"http://purl.org/dc/elements/1.1/ title"`
	Dir     string   `xml:"dir,attr,omitempty"`
	ID      string   `xml:"id,attr,omitempty"`
	Lang    string   `xml:"xml lang,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

// DCLanguage represents dc:language element
type DCLanguage struct {
	XMLName xml.Name `xml:"http://purl.org/dc/elements/1.1/ language"`
	ID      string   `xml:"id,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

// DCOptional represents optional Dublin Core elements
type DCOptional struct {
	XMLName xml.Name
	Dir     string `xml:"dir,attr,omitempty"`
	ID      string `xml:"id,attr,omitempty"`
	Lang    string `xml:"xml lang,attr,omitempty"`
	Value   string `xml:",chardata"`
}

// Meta represents the meta element for generic metadata
type Meta struct {
	XMLName  xml.Name `xml:"meta"`
	Name     string   `xml:"name,attr,omitempty"`
	Content  string   `xml:"content,attr,omitempty"`
	Dir      string   `xml:"dir,attr,omitempty"`
	ID       string   `xml:"id,attr,omitempty"`
	Property string   `xml:"property,attr"`
	Refines  string   `xml:"refines,attr,omitempty"`
	Scheme   string   `xml:"scheme,attr,omitempty"`
	Lang     string   `xml:"xml lang,attr,omitempty"`
	Value    string   `xml:",chardata"`
}

// Link represents the link element for associating resources
type Link struct {
	XMLName    xml.Name `xml:"link"`
	Href       string   `xml:"href,attr"`
	HrefLang   string   `xml:"hreflang,attr,omitempty"`
	ID         string   `xml:"id,attr,omitempty"`
	MediaType  string   `xml:"media-type,attr,omitempty"`
	Properties string   `xml:"properties,attr,omitempty"`
	Refines    string   `xml:"refines,attr,omitempty"`
	Rel        string   `xml:"rel,attr"`
}

// Manifest represents the manifest section
type Manifest struct {
	XMLName xml.Name `xml:"manifest"`
	ID      string   `xml:"id,attr,omitempty"`
	Items   []Item   `xml:"item"`
}

// Item represents a publication resource in the manifest
type Item struct {
	XMLName      xml.Name `xml:"item"`
	Fallback     string   `xml:"fallback,attr,omitempty"`
	Href         string   `xml:"href,attr"`
	ID           string   `xml:"id,attr"`
	MediaOverlay string   `xml:"media-overlay,attr,omitempty"`
	MediaType    string   `xml:"media-type,attr"`
	Properties   string   `xml:"properties,attr,omitempty"`
}

// Spine represents the spine section
type Spine struct {
	XMLName                  xml.Name  `xml:"spine"`
	ID                       string    `xml:"id,attr,omitempty"`
	PageProgressionDirection string    `xml:"page-progression-direction,attr,omitempty"`
	TOC                      string    `xml:"toc,attr,omitempty"`
	ItemRefs                 []ItemRef `xml:"itemref"`
}

// ItemRef represents a reference to an item in the spine
type ItemRef struct {
	XMLName    xml.Name `xml:"itemref"`
	ID         string   `xml:"id,attr,omitempty"`
	IDRef      string   `xml:"idref,attr"`
	Linear     string   `xml:"linear,attr,omitempty"`
	Properties string   `xml:"properties,attr,omitempty"`
}

// Guide represents the legacy guide element
type Guide struct {
	XMLName    xml.Name         `xml:"http://www.idpf.org/2007/opf guide"`
	References []GuideReference `xml:"reference"`
}

type GuideReferenceType string

const (
	GuideRefCover            GuideReferenceType = "cover"
	GuideRefTitlePage        GuideReferenceType = "title-page"
	GuideRefToc              GuideReferenceType = "toc"
	GuideRefIndex            GuideReferenceType = "index"
	GuideRefGlossary         GuideReferenceType = "glossary"
	GuideRefAcknowledgements GuideReferenceType = "acknowledgements"
	GuideRefBibliography     GuideReferenceType = "bibliography"
	GuideRefColophon         GuideReferenceType = "colophon"
	GuideRefCopyrightPage    GuideReferenceType = "copyright-page"
	GuideRefDedication       GuideReferenceType = "dedication"
	GuideRefEpigraph         GuideReferenceType = "epigraph"
	GuideRefForeword         GuideReferenceType = "foreword"
	GuideRefLoi              GuideReferenceType = "loi"
	GuideRefLot              GuideReferenceType = "lot"
	GuideRefNotes            GuideReferenceType = "notes"
	GuideRefPreface          GuideReferenceType = "preface"
	GuideRefText             GuideReferenceType = "text"
)

func (e GuideReferenceType) MarshalText() ([]byte, error) {
	return []byte(e), nil
}

func (e *GuideReferenceType) UnmarshalText(text []byte) error {
	s := string(text)
	switch s {
	case string(GuideRefCover):
		*e = GuideRefCover
	case string(GuideRefTitlePage):
		*e = GuideRefTitlePage
	case string(GuideRefToc):
		*e = GuideRefToc
	case string(GuideRefIndex):
		*e = GuideRefIndex
	case string(GuideRefGlossary):
		*e = GuideRefGlossary
	case string(GuideRefAcknowledgements):
		*e = GuideRefAcknowledgements
	case string(GuideRefBibliography):
		*e = GuideRefBibliography
	case string(GuideRefColophon):
		*e = GuideRefColophon
	case string(GuideRefCopyrightPage):
		*e = GuideRefCopyrightPage
	case string(GuideRefDedication):
		*e = GuideRefDedication
	case string(GuideRefEpigraph):
		*e = GuideRefEpigraph
	case string(GuideRefForeword):
		*e = GuideRefForeword
	case string(GuideRefLoi):
		*e = GuideRefLoi
	case string(GuideRefLot):
		*e = GuideRefLot
	case string(GuideRefNotes):
		*e = GuideRefNotes
	case string(GuideRefPreface):
		*e = GuideRefPreface
	case string(GuideRefText):
		*e = GuideRefText

	default:
		return fmt.Errorf("invalid enum value: %s", s)
	}
	return nil
}

// GuideReference represents a reference in the guide
type GuideReference struct {
	XMLName xml.Name           `xml:"reference"`
	Type    GuideReferenceType `xml:"type,attr"`
	Title   string             `xml:"title,attr,omitempty"`
	Href    string             `xml:"href,attr"`
}

// Bindings represents deprecated bindings element
type Bindings struct {
	XMLName    xml.Name    `xml:"http://www.idpf.org/2007/opf bindings"`
	MediaTypes []MediaType `xml:"mediaType"`
}

// MediaType represents a media type in bindings
type MediaType struct {
	XMLName   xml.Name `xml:"mediaType"`
	Handler   string   `xml:"handler,attr"`
	MediaType string   `xml:"media-type,attr"`
}

// Collection represents a collection element
type Collection struct {
	XMLName     xml.Name            `xml:"collection"`
	Dir         string              `xml:"dir,attr,omitempty"`
	ID          string              `xml:"id,attr,omitempty"`
	Role        string              `xml:"role,attr"`
	Lang        string              `xml:"xml lang,attr,omitempty"`
	Metadata    *CollectionMetadata `xml:"metadata,omitempty"`
	Collections []Collection        `xml:"collection,omitempty"`
	Links       []Link              `xml:"link,omitempty"`
}

// CollectionMetadata represents metadata within a collection
type CollectionMetadata struct {
	XMLName xml.Name `xml:"metadata"`
	// Can contain similar content to main metadata but scoped to collection
	Meta  []Meta `xml:"meta,omitempty"`
	Links []Link `xml:"link,omitempty"`
}
