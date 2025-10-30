package ocf

import "encoding/xml"

type Manifest struct {
	XMLName   xml.Name    `xml:"manifest:manifest"`
	XMLNS     string      `xml:"xmlns:manifest,attr"`
	Version   string      `xml:"manifest:version,attr"`
	FileEntry []FileEntry `xml:"manifest:file-entry"`
}

type FileEntry struct {
	MediaType string `xml:"manifest:media-type,attr"`
	Version   string `xml:"manifest:version,attr,omitempty"`
	FullPath  string `xml:"manifest:full-path,attr"`
}
