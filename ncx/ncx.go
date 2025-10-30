package ncx

import (
	"encoding/xml"
)

type NCX struct {
	XMLName   xml.Name    `xml:"http://www.daisy.org/z3986/2005/ncx/ ncx"`
	Version   string      `xml:"version,attr"`
	Lang      string      `xml:"lang,attr,omitempty"`
	Head      Head        `xml:"head"`
	DocTitle  TextElement `xml:"docTitle"`
	DocAuthor TextElement `xml:"docAuthor,omitempty"`
	NavMap    NavMap      `xml:"navMap"`
	PageList  *PageList   `xml:"pageList,omitempty"`
	NavLists  []NavList   `xml:"navList,omitempty"`
}

func Parse(data []byte) (ncx *NCX, err error) {
	err = xml.Unmarshal(data, &ncx)
	return
}

// Head represents the head element containing metadata
type Head struct {
	Meta []Meta `xml:"meta"`
}

// Meta represents metadata in the head section
type Meta struct {
	Name    string `xml:"name,attr"`
	Content string `xml:"content,attr"`
}

// TextElement represents elements that contain text
type TextElement struct {
	Text string `xml:"text"`
}

// NavMap represents the main navigation map
type NavMap struct {
	NavPoints []NavPoint `xml:"navPoint"`
}

// NavPoint represents a navigation point in the navMap
type NavPoint struct {
	ID        string     `xml:"id,attr"`
	Class     string     `xml:"class,attr,omitempty"`
	PlayOrder string     `xml:"playOrder,attr,omitempty"` // Optional in EPUB
	NavLabel  NavLabel   `xml:"navLabel"`
	Content   Content    `xml:"content"`
	NavPoints []NavPoint `xml:"navPoint,omitempty"` // Nested navPoints
}

// NavLabel represents the label for a navigation point
type NavLabel struct {
	Text string `xml:"text"`
}

// Content represents the content target
type Content struct {
	Src string `xml:"src,attr"`
}

// PageList represents the page list navigation
type PageList struct {
	PageTargets []PageTarget `xml:"pageTarget"`
}

// PageTarget represents a page target in the pageList
type PageTarget struct {
	ID       string   `xml:"id,attr"`
	Type     string   `xml:"type,attr,omitempty"`
	Value    string   `xml:"value,attr"`
	NavLabel NavLabel `xml:"navLabel"`
	Content  Content  `xml:"content"`
}

// NavList represents additional navigation lists (illustrations, tables, etc.)
type NavList struct {
	NavLabel   NavLabel    `xml:"navLabel"`
	NavTargets []NavTarget `xml:"navTarget"`
}

// NavTarget represents a target in a navList
type NavTarget struct {
	ID       string   `xml:"id,attr"`
	NavLabel NavLabel `xml:"navLabel"`
	Content  Content  `xml:"content"`
}

// NCXProcessor provides functionality to work with NCX data
type NCXProcessor struct {
	NCX *NCX
}
