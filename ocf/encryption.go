package ocf

import (
	"encoding/xml"
)

// Encryption represents the root element of encryption.xml
type Encryption struct {
	XMLName       xml.Name        `xml:"urn:oasis:names:tc:opendocument:xmlns:container encryption"`
	EncryptedKey  []EncryptedKey  `xml:"http://www.w3.org/2001/04/xmlenc# EncryptedKey"`
	EncryptedData []EncryptedData `xml:"http://www.w3.org/2001/04/xmlenc# EncryptedData"`
}

// EncryptedKey represents the xenc:EncryptedKey element
type EncryptedKey struct {
	XMLName              xml.Name              `xml:"http://www.w3.org/2001/04/xmlenc# EncryptedKey"`
	Id                   *string               `xml:"Id,attr,omitempty"`
	Type                 *string               `xml:"Type,attr,omitempty"`
	MimeType             *string               `xml:"MimeType,attr,omitempty"`
	Encoding             *string               `xml:"Encoding,attr,omitempty"`
	Recipient            *string               `xml:"Recipient,attr,omitempty"`
	EncryptionMethod     *EncryptionMethod     `xml:"http://www.w3.org/2001/04/xmlenc# EncryptionMethod"`
	KeyInfo              KeyInfo               `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo"`
	CipherData           CipherData            `xml:"http://www.w3.org/2001/04/xmlenc# CipherData"`
	EncryptionProperties *EncryptionProperties `xml:"http://www.w3.org/2001/04/xmlenc# EncryptionProperties,omitempty"`
	ReferenceList        *ReferenceList        `xml:"http://www.w3.org/2001/04/xmlenc# ReferenceList,omitempty"`
	CarriedKeyName       *string               `xml:"http://www.w3.org/2001/04/xmlenc# CarriedKeyName,omitempty"`
}

// EncryptedData represents the xenc:EncryptedData element
type EncryptedData struct {
	XMLName              xml.Name              `xml:"http://www.w3.org/2001/04/xmlenc# EncryptedData"`
	Id                   *string               `xml:"Id,attr,omitempty"`
	Type                 *string               `xml:"Type,attr,omitempty"`
	MimeType             *string               `xml:"MimeType,attr,omitempty"`
	Encoding             *string               `xml:"Encoding,attr,omitempty"`
	EncryptionMethod     *EncryptionMethod     `xml:"http://www.w3.org/2001/04/xmlenc# EncryptionMethod,omitempty"`
	KeyInfo              *KeyInfo              `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo,omitempty"`
	CipherData           CipherData            `xml:"http://www.w3.org/2001/04/xmlenc# CipherData"`
	EncryptionProperties *EncryptionProperties `xml:"http://www.w3.org/2001/04/xmlenc# EncryptionProperties,omitempty"`
}

// EncryptionMethod represents the xenc:EncryptionMethod element
type EncryptionMethod struct {
	XMLName    xml.Name `xml:"http://www.w3.org/2001/04/xmlenc# EncryptionMethod"`
	Algorithm  string   `xml:"Algorithm,attr"`
	KeySize    *int     `xml:"http://www.w3.org/2001/04/xmlenc# KeySize,omitempty"`
	OAEPparams *string  `xml:"http://www.w3.org/2001/04/xmlenc# OAEPparams,omitempty"`
	// Other algorithm-specific parameters could be added here
}

// CipherData represents the xenc:CipherData element
type CipherData struct {
	XMLName         xml.Name         `xml:"http://www.w3.org/2001/04/xmlenc# CipherData"`
	CipherValue     *string          `xml:"http://www.w3.org/2001/04/xmlenc# CipherValue,omitempty"`
	CipherReference *CipherReference `xml:"http://www.w3.org/2001/04/xmlenc# CipherReference,omitempty"`
}

// CipherReference represents the xenc:CipherReference element
type CipherReference struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2001/04/xmlenc# CipherReference"`
	URI        string      `xml:"URI,attr"`
	Transforms *Transforms `xml:"http://www.w3.org/2000/09/xmldsig# Transforms,omitempty"`
}

// Transforms represents the ds:Transforms element
type Transforms struct {
	XMLName   xml.Name    `xml:"http://www.w3.org/2000/09/xmldsig# Transforms"`
	Transform []Transform `xml:"http://www.w3.org/2000/09/xmldsig# Transform"`
}

// Transform represents the ds:Transform element
type Transform struct {
	XMLName   xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# Transform"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// KeyInfo represents the ds:KeyInfo element
type KeyInfo struct {
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo"`
	// KeyInfo can contain various elements - this is a basic implementation
	// You may need to extend this based on your specific needs
	X509Data *X509Data `xml:"http://www.w3.org/2000/09/xmldsig# X509Data,omitempty"`
	KeyName  *string   `xml:"http://www.w3.org/2000/09/xmldsig# KeyName,omitempty"`
	KeyValue *KeyValue `xml:"http://www.w3.org/2000/09/xmldsig# KeyValue,omitempty"`
}

// X509Data represents the ds:X509Data element
type X509Data struct {
	XMLName         xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# X509Data"`
	X509Certificate *string  `xml:"http://www.w3.org/2000/09/xmldsig# X509Certificate,omitempty"`
	X509SubjectName *string  `xml:"http://www.w3.org/2000/09/xmldsig# X509SubjectName,omitempty"`
}

// KeyValue represents the ds:KeyValue element
type KeyValue struct {
	XMLName     xml.Name     `xml:"http://www.w3.org/2000/09/xmldsig# KeyValue"`
	RSAKeyValue *RSAKeyValue `xml:"http://www.w3.org/2000/09/xmldsig# RSAKeyValue,omitempty"`
}

// RSAKeyValue represents the ds:RSAKeyValue element
type RSAKeyValue struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# RSAKeyValue"`
	Modulus  string   `xml:"http://www.w3.org/2000/09/xmldsig# Modulus"`
	Exponent string   `xml:"http://www.w3.org/2000/09/xmldsig# Exponent"`
}

// EncryptionProperties represents the xenc:EncryptionProperties element
type EncryptionProperties struct {
	XMLName            xml.Name             `xml:"http://www.w3.org/2001/04/xmlenc# EncryptionProperties"`
	Id                 *string              `xml:"Id,attr,omitempty"`
	EncryptionProperty []EncryptionProperty `xml:"http://www.w3.org/2001/04/xmlenc# EncryptionProperty"`
}

// EncryptionProperty represents the xenc:EncryptionProperty element
type EncryptionProperty struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/04/xmlenc# EncryptionProperty"`
	Target  *string  `xml:"Target,attr,omitempty"`
	Id      *string  `xml:"Id,attr,omitempty"`
	// Content is ANY, so we might need to handle various content types
}

// ReferenceList represents the xenc:ReferenceList element
type ReferenceList struct {
	XMLName       xml.Name        `xml:"http://www.w3.org/2001/04/xmlenc# ReferenceList"`
	DataReference []DataReference `xml:"http://www.w3.org/2001/04/xmlenc# DataReference,omitempty"`
	KeyReference  []KeyReference  `xml:"http://www.w3.org/2001/04/xmlenc# KeyReference,omitempty"`
}

// DataReference represents the xenc:DataReference element
type DataReference struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/04/xmlenc# DataReference"`
	URI     *string  `xml:"URI,attr,omitempty"`
	Id      *string  `xml:"Id,attr,omitempty"`
}

// KeyReference represents the xenc:KeyReference element
type KeyReference struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/04/xmlenc# KeyReference"`
	URI     *string  `xml:"URI,attr,omitempty"`
	Id      *string  `xml:"Id,attr,omitempty"`
}
