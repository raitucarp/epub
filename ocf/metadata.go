package ocf

import "encoding/xml"

// Metadata represents the container-level metadata file
// Root element: metadata in namespace http://www.idpf.org/2013/metadata
type Metadata struct {
	XMLName xml.Name `xml:"http://www.idpf.org/2013/metadata metadata"`
	// The content model of this file is defined by the EPUB
	// multiple-rendition specification; preserve the raw inner XML so that
	// callers can process it without data loss.
	Content any `xml:",innerxml"`
}
