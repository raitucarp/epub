package ocf

import "encoding/xml"

// Metadata represents the container-level metadata file
// Root element: metadata in namespace http://www.idpf.org/2013/metadata
type Metadata struct {
	XMLName xml.Name `xml:"http://www.idpf.org/2013/metadata metadata"`
	// Content is flexible since this version doesn't define specific metadata
	// Using Any to allow any namespace-qualified elements
	Any []xml.Name `xml:",any"`
}
