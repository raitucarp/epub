package ocf

import (
	"encoding/xml"
)

const metaInfDirectoryName = "META-INF"

type metaInfReservedFile string

const (
	containerFile  metaInfReservedFile = "container.xml"
	encryptionFile metaInfReservedFile = "encryption.xml"
	manifestFile   metaInfReservedFile = "manifest.xml"
	metadataFile   metaInfReservedFile = "metadata.xml"
	rightsFile     metaInfReservedFile = "rights.xml"
	signaturesFile metaInfReservedFile = "signatures.xml"
)

func (f metaInfReservedFile) Unmarshal(data []byte, v any) (err error) {
	err = xml.Unmarshal(data, v)
	if err != nil {
		return
	}
	return
}

var metaInfReservedFiles = []metaInfReservedFile{
	containerFile,
	encryptionFile,
	manifestFile,
	metadataFile,
	rightsFile,
	signaturesFile,
}

var requiredMetaInfFiles = []metaInfReservedFile{containerFile}

type MetaInf struct {
	container  Container
	signatures Signatures
	encryption Encryption
	metadata   Metadata
	rights     Rights
	manifest   Manifest
}

func (metaInf *MetaInf) parseContainer(data []byte) (err error) {
	return containerFile.Unmarshal(data, &metaInf.container)
}

func (metaInf *MetaInf) parseEncryption(data []byte) (err error) {
	return encryptionFile.Unmarshal(data, &metaInf.encryption)
}

func (metaInf *MetaInf) parseManifest(data []byte) (err error) {
	return manifestFile.Unmarshal(data, &metaInf.manifest)
}

func (metaInf *MetaInf) parseMetadata(data []byte) (err error) {
	return manifestFile.Unmarshal(data, &metaInf.metadata)
}

func (metaInf *MetaInf) parseRights(data []byte) (err error) {
	return manifestFile.Unmarshal(data, &metaInf.rights)
}

func (metaInf *MetaInf) parseSignatures(data []byte) (err error) {
	return manifestFile.Unmarshal(data, &metaInf.signatures)
}
