package ocf

import "encoding/xml"

type Container struct {
	XMLName   xml.Name  `xml:"urn:oasis:names:tc:opendocument:xmlns:container container"`
	Version   string    `xml:"version,attr"`
	RootFiles RootFiles `xml:"rootfiles"`
	Links     *Links    `xml:"links,omitempty"`
}

// RootFiles represents the rootfiles element
type RootFiles struct {
	RootFile []RootFile `xml:"rootfile"`
}

// RootFile represents a rootfile element
type RootFile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`

	Media      string `xml:"http://www.idpf.org/2013/rendition media,attr,omitempty"`
	Layout     string `xml:"http://www.idpf.org/2013/rendition layout,attr,omitempty"`
	Language   string `xml:"http://www.idpf.org/2013/rendition language,attr,omitempty"`
	AccessMode string `xml:"http://www.idpf.org/2013/rendition accessMode,attr,omitempty"`
	Label      string `xml:"http://www.idpf.org/2013/rendition label,attr,omitempty"`
}

// Links represents the links element
type Links struct {
	Link []Link `xml:"link"`
}

// Link represents a link element
type Link struct {
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr,omitempty"`
	Rel       string `xml:"rel,attr"`
}
