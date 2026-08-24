package ocf

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/raitucarp/epub/pkg"
)

func NewOCFZipContainer() *OCFZipContainer {
	return &OCFZipContainer{
		files:   make(map[string][]byte),
		metaInf: MetaInf{},
	}
}

func (z *OCFZipContainer) AddFile(filePath string, content []byte) {
	z.files[filePath] = content
}

func (z *OCFZipContainer) AddMimeType() {
	z.AddFile("mimetype", []byte(MimeType))
}

func (z *OCFZipContainer) AddPackage(filename string, packageData pkg.Package) (err error) {
	content, err := xml.MarshalIndent(packageData, "", "  ")
	if err != nil {
		return
	}

	finalXml := append([]byte(xml.Header), content...)
	z.AddFile(filename, finalXml)
	return
}

func (z *OCFZipContainer) AddContainerXML(rootFiles ...string) (err error) {
	container := Container{Version: "1.0"}
	container.XMLName.Space = "urn:oasis:names:tc:opendocument:xmlns:container"
	for _, rootFile := range rootFiles {
		container.RootFiles.RootFile = append(container.RootFiles.RootFile, RootFile{
			FullPath:  rootFile,
			MediaType: EPUBContainerMime,
		})
	}

	content, err := xml.MarshalIndent(container, "", "  ")
	if err != nil {
		return
	}

	finalXml := append([]byte(xml.Header), content...)
	z.AddFile("META-INF/container.xml", finalXml)
	return
}

func addFileToZip(zipWriter *zip.Writer, filename string, content []byte) error {
	// Create a header for the file
	header := &zip.FileHeader{
		Name:   filename,
		Method: zip.Deflate, // Use compression for all files
	}

	// Special handling for mimetype file (must be uncompressed and first).
	// The OCF specification also requires it to have no extra field, so the
	// modification timestamp (which Go encodes as an extra field) is omitted.
	if filename == "mimetype" {
		header.Method = zip.Store // No compression
	} else {
		header.Modified = time.Now()
	}

	// Create the file in the zip
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Write the content
	_, err = io.Copy(writer, bytes.NewReader(content))
	return err
}

func (z *OCFZipContainer) Write(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// The OCF specification requires the mimetype file to be the first
	// entry in the ZIP archive and stored without compression. Write it
	// explicitly before the remaining files to guarantee ordering.
	if content, ok := z.files["mimetype"]; ok {
		err := addFileToZip(zipWriter, "mimetype", content)
		if err != nil {
			return fmt.Errorf("error adding mimetype: %w", err)
		}
	}

	// Write remaining files in a deterministic order.
	names := make([]string, 0, len(z.files))
	for name := range z.files {
		if name != "mimetype" {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	for _, name := range names {
		err := addFileToZip(zipWriter, name, z.files[name])
		if err != nil {
			return fmt.Errorf("error adding %s: %w", name, err)
		}
	}

	return nil
}
