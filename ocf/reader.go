package ocf

import (
	"archive/zip"
	"bytes"
	"errors"
)

func newContainerAndParse(file *zip.Reader) (container *OCFZipContainer, err error) {
	container = &OCFZipContainer{}
	err = container.readFiles(file)
	if err != nil {
		return
	}

	err = container.parseAllMetaInfFiles()
	if err != nil {
		return
	}

	if container.MimeType() != MimeType {
		return nil, errors.New("Mimetype mismatch")
	}

	return container, nil
}

// OpenReader opens the OCF ZIP container stored at name.
func OpenReader(name string) (container *OCFZipContainer, err error) {
	z, err := zip.OpenReader(name)
	if err != nil {
		return
	}
	defer z.Close()

	return newContainerAndParse(&z.Reader)
}

// NewReader parses an OCF ZIP container from the raw bytes b.
func NewReader(b []byte) (container *OCFZipContainer, err error) {
	byteReader := bytes.NewReader(b)

	z, err := zip.NewReader(byteReader, int64(byteReader.Len()))
	if err != nil {
		return
	}

	return newContainerAndParse(z)
}
