package ocf

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

// MimeType is the media type of an OCF ZIP container as recorded in the
// mimetype file.
const MimeType = "application/epub+zip"

// EPUBContainerMime is the media type of an EPUB package document.
const EPUBContainerMime = "application/oebps-package+xml"

// OCFZipContainer is an in-memory representation of an OCF ZIP container,
// holding both the packaged files and the parsed META-INF metadata.
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

		cleanPath := path.Clean(f.Name)
		if strings.Contains(f.Name, `\`) || !filepath.IsLocal(f.Name) || cleanPath == ".." || strings.HasPrefix(cleanPath, "../") {
			return fmt.Errorf("invalid path in zip: %s", f.Name)
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		// Prevent zip bomb by limiting the read size to the uncompressed size
		// and enforcing a maximum file size (e.g., 1GB)
		const maxFileSize = 1024 * 1024 * 1024
		if f.UncompressedSize64 > maxFileSize {
			rc.Close()
			return fmt.Errorf("file %s is too large: %d bytes", f.Name, f.UncompressedSize64)
		}

		content, err := io.ReadAll(io.LimitReader(rc, int64(f.UncompressedSize64)))
		if err != nil {
			rc.Close()
			return err
		}

		z.files[cleanPath] = content
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
		if err := parseMap[reversedFileName](data); err != nil {
			return err
		}
	}

	return nil
}

// MimeType returns the contents of the mimetype file with surrounding
// whitespace trimmed.
func (z *OCFZipContainer) MimeType() string {
	return strings.TrimSpace(string(z.files["mimetype"]))
}

// Container returns the parsed container.xml document.
func (z *OCFZipContainer) Container() *Container {
	return &z.metaInf.container
}

// Signatures returns the parsed signatures.xml document.
func (z *OCFZipContainer) Signatures() *Signatures {
	return &z.metaInf.signatures
}

// Encryption returns the parsed encryption.xml document.
func (z *OCFZipContainer) Encryption() *Encryption {
	return &z.metaInf.encryption
}

// Metadata returns the parsed metadata.xml document.
func (z *OCFZipContainer) Metadata() *Metadata {
	return &z.metaInf.metadata
}

// Rights returns the parsed rights.xml document.
func (z *OCFZipContainer) Rights() *Rights {
	return &z.metaInf.rights
}

// Manifest returns the parsed manifest.xml document.
func (z *OCFZipContainer) Manifest() *Manifest {
	return &z.metaInf.manifest
}

// AllFiles returns the raw files stored in the container, keyed by their path
// relative to the container root.
func (z *OCFZipContainer) AllFiles() map[string][]byte {
	return z.files
}

// SelectFile returns the content of the file at the given path, or an error if
// the file is not present in the container.
func (z *OCFZipContainer) SelectFile(name string) (data []byte, err error) {
	data, ok := z.files[name]
	if !ok {
		return nil, fmt.Errorf("No file found with name %s", name)
	}

	return data, nil
}

// NonMetaInfFiles returns all container files except those in the META-INF
// directory, keyed by their path relative to the container root.
func (z *OCFZipContainer) NonMetaInfFiles() map[string][]byte {
	files := map[string][]byte{}
	for filePath, data := range z.files {
		if getRootDirectory(filePath) != metaInfDirectoryName {
			files[filePath] = data
		}
	}
	return files
}
