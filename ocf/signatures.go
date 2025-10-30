package ocf

import "encoding/xml"

type Signatures struct {
	XMLName   xml.Name    `xml:"urn:oasis:names:tc:opendocument:xmlns:container signatures"`
	XMLNS     string      `xml:"xmlns,attr,omitempty"`
	Signature []Signature `xml:"http://www.w3.org/2000/09/xmldsig# Signature"`
}

// Signature represents an XML Signature as defined by xmldsig-core1
type Signature struct {
	ID             string     `xml:"Id,attr,omitempty"`
	SignedInfo     SignedInfo `xml:"SignedInfo"`
	SignatureValue string     `xml:"SignatureValue"`
	KeyInfo        KeyInfo    `xml:"KeyInfo,omitempty"`
	Object         []Object   `xml:"Object,omitempty"`
}

// SignedInfo contains the canonicalization method, signature method, and references
type SignedInfo struct {
	CanonicalizationMethod Method      `xml:"CanonicalizationMethod"`
	SignatureMethod        Method      `xml:"SignatureMethod"`
	Reference              []Reference `xml:"Reference"`
}

// Method represents algorithm methods
type Method struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// Reference represents a reference to signed data
type Reference struct {
	URI          string      `xml:"URI,attr,omitempty"`
	Transforms   *Transforms `xml:"Transforms,omitempty"`
	DigestMethod Method      `xml:"DigestMethod"`
	DigestValue  string      `xml:"DigestValue"`
}

// DSAKeyValue represents DSA key parameters
type DSAKeyValue struct {
	P string `xml:"P,omitempty"`
	Q string `xml:"Q,omitempty"`
	G string `xml:"G,omitempty"`
	Y string `xml:"Y,omitempty"`
}

// Object contains additional data like Manifest
type Object struct {
	Manifest SignatureManifest `xml:"Manifest,omitempty"`
	// Other object types can be added here
}

// Manifest contains references to the signed resources
type SignatureManifest struct {
	ID        string      `xml:"Id,attr,omitempty"`
	Reference []Reference `xml:"Reference"`
}
