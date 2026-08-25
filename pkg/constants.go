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
	// NamespaceXML is the reserved XML namespace used by xml:lang and other
	// xml:* attributes.
	NamespaceXML = "http://www.w3.org/XML/1998/namespace"
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
	MediaTypeHTML  = "text/html"
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

// ImageMediaTypes lists the media types that identify image resources in the
// package manifest.
var ImageMediaTypes = []string{
	MediaTypeSVG,
	MediaTypeJPEG,
	MediaTypeGIF,
	MediaTypeWebP,
	MediaTypePNG,
}
