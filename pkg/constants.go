package pkg

// XML namespace URIs used across the EPUB package document and its
// metadata sections. These are defined by the EPUB 3.3 specification and
// its referenced vocabularies.
const (
	// NamespaceOPF is the namespace of the Package Document elements.
	NamespaceOPF = "http://www.idpf.org/2007/opf"
	// NamespaceDC is the namespace of the Dublin Core metadata elements.
	NamespaceDC = "http://purl.org/dc/elements/1.1/"
	// NamespaceOPS is the namespace of the EPUB content and reserved
	// prefixes (epub:type and friends).
	NamespaceOPS = "http://www.idpf.org/2007/ops"
)

// Constants for common property values
const (
	PropertyNav                = "nav"
	PropertyCoverImage         = "cover-image"
	PropertyMathML             = "mathml"
	PropertyRemoteResources    = "remote-resources"
	PropertyScripted           = "scripted"
	PropertySVG                = "svg"
	PropertySwitch             = "switch"
	PropertyLayoutPrePaginated = "layout-pre-paginated"

	// Media types
	MediaTypeXHTML = "application/xhtml+xml"
	MediaTypeSVG   = "image/svg+xml"
	MediaTypeJPEG  = "image/jpeg"
	MediaTypeGIF   = "image/gif"
	MediaTypeWebP  = "image/webp"
	MediaTypePNG   = "image/png"
	MediaTypeCSS   = "text/css"
	MediaTypeNCX   = "application/x-dtbncx+xml"

	// Spine directions
	SpineDirectionLTR     = "ltr"
	SpineDirectionRTL     = "rtl"
	SpineDirectionDefault = "default"

	// Linear values
	LinearYes = "yes"
	LinearNo  = "no"
)

var ImageMediaTypes = []string{
	MediaTypeSVG,
	MediaTypeJPEG,
	MediaTypeGIF,
	MediaTypeWebP,
	MediaTypePNG,
}
