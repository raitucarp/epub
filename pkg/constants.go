package pkg

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
