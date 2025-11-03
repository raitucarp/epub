package ocf

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"path"
	"slices"
)

const MimeType = "application/epub+zip"
const EPUBContainerMime = "application/oebps-package+xml"

type OCFZipContainer struct {
	files   map[string][]byte
	metaInf MetaInf
}

func (z *OCFZipContainer) readFiles(zrc *zip.Reader) (err error) {
	z.files = make(map[string][]byte)
	for _, f := range zrc.File {
		info := f.FileInfo()
		if info.IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		content, err := io.ReadAll(rc)
		if err != nil {
			return err
		}

		z.files[f.Name] = content
		rc.Close()
	}
	return
}

func (z *OCFZipContainer) parseAllMetaInfFiles() error {
	reservedFiles := map[metaInfReservedFile][]byte{}
	for filePath, data := range z.files {
		if getRootDirectory(filePath) != metaInfDirectoryName {
			continue
		}

		filename := path.Base(filePath)
		if !slices.Contains(metaInfReservedFiles, metaInfReservedFile(filename)) {
			continue
		}

		reservedFiles[metaInfReservedFile(filename)] = data
	}

	for _, filename := range requiredMetaInfFiles {
		_, ok := reservedFiles[filename]
		if !ok {
			return errors.New("Package does not have required files")
		}
	}

	parseMap := map[metaInfReservedFile]func(data []byte) (err error){
		containerFile:  z.metaInf.parseContainer,
		encryptionFile: z.metaInf.parseEncryption,
		manifestFile:   z.metaInf.parseManifest,
		metadataFile:   z.metaInf.parseMetadata,
		rightsFile:     z.metaInf.parseRights,
		signaturesFile: z.metaInf.parseSignatures,
	}

	for reversedFileName, data := range reservedFiles {
		parseMap[reversedFileName](data)
	}

	return nil
}

func (z *OCFZipContainer) MimeType() string {
	return string(z.files["mimetype"])
}

func (z *OCFZipContainer) Container() *Container {
	return &z.metaInf.container
}

func (z *OCFZipContainer) Signatures() *Signatures {
	return &z.metaInf.signatures
}

func (z *OCFZipContainer) Encryption() *Encryption {
	return &z.metaInf.encryption
}

func (z *OCFZipContainer) Metadata() *Metadata {
	return &z.metaInf.metadata
}

func (z *OCFZipContainer) Rights() *Rights {
	return &z.metaInf.rights
}

func (z *OCFZipContainer) Manifest() *Manifest {
	return &z.metaInf.manifest
}

func (z *OCFZipContainer) AllFiles() map[string][]byte {
	return z.files
}

func (z *OCFZipContainer) SelectFile(name string) (data []byte, err error) {
	data, ok := z.files[name]
	if !ok {
		return nil, fmt.Errorf("No file found with name %s", name)
	}

	return data, nil
}

func (z *OCFZipContainer) NonMetaInfFiles() map[string][]byte {
	files := map[string][]byte{}
	for filePath, data := range z.files {
		if getRootDirectory(filePath) != metaInfDirectoryName {
			files[filePath] = data
		}
	}
	return files
}
