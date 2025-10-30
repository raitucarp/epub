package ocf

import "encoding/xml"

// Rights represents the rights management file
// The structure is not defined in this specification version
type Rights struct {
	XMLName xml.Name `xml:"rights"`
	// Since the structure is reserved but not defined, use flexible content
	Content any `xml:",innerxml"` // Store raw XML for rights expressions
}
