package ocf

import "encoding/xml"

// Rights represents the rights management file. Its structure is not defined
// by the OCF specification, so the raw inner XML is preserved.
type Rights struct {
	XMLName xml.Name `xml:"rights"`
	// Content stores the raw XML of the rights expression.
	Content string `xml:",innerxml"`
}
