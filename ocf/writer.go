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

// NewOCFZipContainer creates an empty OCF ZIP container ready to be populated
// with files and META-INF data.
func NewOCFZipContainer() *OCFZipContainer {
	return &OCFZipContainer{
		files:   make(map[string][]byte),
		metaInf: MetaInf{},
	}
}

// AddFile stores the given content in the container at filePath.
func (z *OCFZipContainer) AddFile(filePath string, content []byte) {
	z.files[filePath] = content
}

// RemoveFile removes the file at filePath from the container.
func (z *OCFZipContainer) RemoveFile(filePath string) {
	delete(z.files, filePath)
}

// Clone creates an independent copy of the OCF ZIP container.
func (z *OCFZipContainer) Clone() *OCFZipContainer {
	if z == nil {
		return nil
	}
	clonedFiles := make(map[string][]byte, len(z.files))
	for k, v := range z.files {
		data := make([]byte, len(v))
		copy(data, v)
		clonedFiles[k] = data
	}
	cloned := &OCFZipContainer{
		files:   clonedFiles,
		metaInf: z.metaInf,
	}
	if len(z.metaInf.container.RootFiles.RootFile) > 0 {
		cloned.metaInf.container.RootFiles.RootFile = make([]RootFile, len(z.metaInf.container.RootFiles.RootFile))
		copy(cloned.metaInf.container.RootFiles.RootFile, z.metaInf.container.RootFiles.RootFile)
	}
	return cloned
}

// AddMimeType adds the required mimetype file to the container.
func (z *OCFZipContainer) AddMimeType() {
	z.AddFile("mimetype", []byte(MimeType))
}

// AddPackage marshals packageData and stores it in the container at filename.
func (z *OCFZipContainer) AddPackage(filename string, packageData pkg.Package) (err error) {
	content, err := xml.MarshalIndent(packageData, "", "  ")
	if err != nil {
		return
	}

	finalXml := append([]byte(xml.Header), content...)
	z.AddFile(filename, finalXml)
	return
}

// AddContainerXML writes a container.xml file listing the given root file
// paths, each registered with the EPUB package media type.
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

// Write serializes the container to a ZIP archive at filename. The mimetype
// file is written first, uncompressed, as required by the OCF specification.
func (z *OCFZipContainer) Write(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return z.WriteArchive(file)
}

// WriteArchive serializes the container as a ZIP archive to w. The mimetype
// file is written first, uncompressed, as required by the OCF specification.
func (z *OCFZipContainer) WriteArchive(w io.Writer) (err error) {
	zipWriter := zip.NewWriter(w)
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
