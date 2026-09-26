package epub

import (
	"encoding/xml"
	"errors"
	"slices"
	"strconv"
	"strings"

	"github.com/raitucarp/epub/ocf"
	"github.com/raitucarp/epub/pkg"
)

// Reader provides an interface for reading and accessing EPUB publication
// data. It offers methods for retrieving metadata, navigation structures,
// content documents, resources, and images.
type Reader struct {
	epub *Epub
}

func newReaderFromZip(zipContainer *ocf.OCFZipContainer) (reader Reader, err error) {
	reader.epub = &Epub{
		packagePaths: make(map[string]string),
		packagePubs:  make(map[string]*pkg.Package),
		metadata:     make(map[string]any),
		zipContainer: zipContainer,
	}

	err = reader.parseRootFiles(zipContainer)
	if err != nil {
		return
	}

	reader.SelectPackageRendition("default")
	return

}

func (r *Reader) parseRootFiles(z *ocf.OCFZipContainer) (err error) {
	for index, rootFile := range z.Container().RootFiles.RootFile {
		packageFullPath := rootFile.FullPath

		data, err := z.SelectFile(packageFullPath)
		if err != nil {
			return err
		}

		var packagePub pkg.Package
		err = xml.Unmarshal(data, &packagePub)
		if err != nil {
			return err
		}

		renditionVars := renditionKey(rootFile, index)
		r.epub.packagePaths[renditionVars] = packageFullPath
		r.epub.packagePubs[renditionVars] = &packagePub
	}

	return nil
}

// renditionKey derives the key used to select a package rendition. The first
// rootfile is always the default rendition, as defined by the EPUB
// multiple-rendition specification. Subsequent renditions are keyed by their
// rendition selection attributes.
func renditionKey(rootFile ocf.RootFile, index int) string {
	if index == 0 {
		return "default"
	}

	parts := []string{}
	if rootFile.Media != "" {
		parts = append(parts, rootFile.Media)
	}
	if rootFile.Layout != "" {
		parts = append(parts, rootFile.Layout)
	}
	if rootFile.Language != "" {
		parts = append(parts, rootFile.Language)
	}
	if rootFile.AccessMode != "" {
		parts = append(parts, rootFile.AccessMode)
	}
	if rootFile.Label != "" {
		parts = append(parts, rootFile.Label)
	}

	if len(parts) > 0 {
		return strings.Join(parts, "_")
	}

	return "rendition-" + strconv.Itoa(index)
}

// OpenReader opens an EPUB file from the provided file path and returns
// a Reader instance. The file must exist and be a valid EPUB container.
func OpenReader(name string) (reader Reader, err error) {
	zipContainer, err := ocf.OpenReader(name)
	if err != nil {
		return
	}

	return newReaderFromZip(zipContainer)
}

// ListRenditions returns the identifiers of all package renditions available
// in the publication. The first rendition is always named "default".
func (r *Reader) ListRenditions() []string {
	renditions := make([]string, 0, len(r.epub.packagePubs))
	for rendition := range r.epub.packagePubs {
		renditions = append(renditions, rendition)
	}

	slices.Sort(renditions)

	// Ensure the default rendition is always listed first.
	if index := slices.Index(renditions, "default"); index > 0 {
		renditions = slices.Delete(renditions, index, index+1)
		renditions = slices.Insert(renditions, 0, "default")
	}

	return renditions
}

// NewReader creates a new Reader instance from a raw EPUB byte slice.
// The byte slice must represent a valid EPUB container.
func NewReader(b []byte) (reader Reader, err error) {
	zipContainer, err := ocf.NewReader(b)
	if err != nil {
		return
	}

	return newReaderFromZip(zipContainer)
}

// Edit creates an Editor instance from this Reader, allowing modification of
// metadata, resources, spine order, and content before writing back to disk or memory.
func (r *Reader) Edit() (*Editor, error) {
	if r == nil || r.epub == nil || r.epub.zipContainer == nil {
		return nil, errors.New("epub: reader is nil or uninitialized")
	}

	return newEditor(r)
}

