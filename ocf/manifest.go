package ocf

import "encoding/xml"

const manifestNamespace = "urn:oasis:names:tc:opendocument:xmlns:manifest:1.0"

// Manifest is the root element of the manifest.xml document, which lists the
// files stored in the container and their media types.
type Manifest struct {
	XMLName   xml.Name    `xml:"urn:oasis:names:tc:opendocument:xmlns:manifest:1.0 manifest"`
	XMLNS     string      `xml:"xmlns:manifest,attr"`
	Version   string      `xml:"urn:oasis:names:tc:opendocument:xmlns:manifest:1.0 version,attr"`
	FileEntry []FileEntry `xml:"urn:oasis:names:tc:opendocument:xmlns:manifest:1.0 file-entry"`
}

// FileEntry represents a single file listed in manifest.xml.
type FileEntry struct {
	MediaType string `xml:"urn:oasis:names:tc:opendocument:xmlns:manifest:1.0 media-type,attr"`
	Version   string `xml:"urn:oasis:names:tc:opendocument:xmlns:manifest:1.0 version,attr,omitempty"`
	FullPath  string `xml:"urn:oasis:names:tc:opendocument:xmlns:manifest:1.0 full-path,attr"`
}
