package ocf

import "encoding/xml"

// Signatures is the root element of the signatures.xml document, which holds
// the digital signatures for the publication's resources.
type Signatures struct {
	XMLName   xml.Name    `xml:"urn:oasis:names:tc:opendocument:xmlns:container signatures"`
	XMLNS     string      `xml:"xmlns,attr,omitempty"`
	Signature []Signature `xml:"http://www.w3.org/2000/09/xmldsig# Signature"`
}

// Signature represents an XML Signature as defined by xmldsig-core1
type Signature struct {
	ID             string     `xml:"Id,attr,omitempty"`
	SignedInfo     SignedInfo `xml:"http://www.w3.org/2000/09/xmldsig# SignedInfo"`
	SignatureValue string     `xml:"http://www.w3.org/2000/09/xmldsig# SignatureValue"`
	KeyInfo        KeyInfo    `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo,omitempty"`
	Object         []Object   `xml:"http://www.w3.org/2000/09/xmldsig# Object,omitempty"`
}

// SignedInfo contains the canonicalization method, signature method, and references
type SignedInfo struct {
	CanonicalizationMethod Method      `xml:"http://www.w3.org/2000/09/xmldsig# CanonicalizationMethod"`
	SignatureMethod        Method      `xml:"http://www.w3.org/2000/09/xmldsig# SignatureMethod"`
	Reference              []Reference `xml:"http://www.w3.org/2000/09/xmldsig# Reference"`
}

// Method represents algorithm methods
type Method struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// Reference represents a reference to signed data
type Reference struct {
	URI          string      `xml:"URI,attr,omitempty"`
	Transforms   *Transforms `xml:"http://www.w3.org/2000/09/xmldsig# Transforms,omitempty"`
	DigestMethod Method      `xml:"http://www.w3.org/2000/09/xmldsig# DigestMethod"`
	DigestValue  string      `xml:"http://www.w3.org/2000/09/xmldsig# DigestValue"`
}

// Object contains additional data like Manifest
type Object struct {
	Manifest SignatureManifest `xml:"http://www.w3.org/2000/09/xmldsig# Manifest,omitempty"`
	// Other object types can be added here
}

// SignatureManifest contains references to the signed resources
type SignatureManifest struct {
	ID        string      `xml:"Id,attr,omitempty"`
	Reference []Reference `xml:"http://www.w3.org/2000/09/xmldsig# Reference"`
}
