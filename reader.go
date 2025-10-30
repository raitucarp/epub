package epub

import (
	"encoding/xml"
	"strings"

	"github.com/raitucarp/epub/ocf"
	"github.com/raitucarp/epub/pkg"
)

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
	reader.parseMetadata()
	return

}

func (r *Reader) parseRootFiles(z *ocf.OCFZipContainer) (err error) {
	for _, rootFile := range z.Container().RootFiles.RootFile {
		packageFullPath := rootFile.FullPath

		rendition := []string{"default"}
		if rootFile.Media != "" {
			rendition = append(rendition, rootFile.Media)
		}

		if rootFile.Layout != "" {
			rendition = append(rendition, rootFile.Layout)
		}

		if rootFile.Language != "" {
			rendition = append(rendition, rootFile.Language)
		}

		if rootFile.AccessMode != "" {
			rendition = append(rendition, rootFile.AccessMode)
		}

		if rootFile.Label != "" {
			rendition = append(rendition, rootFile.Label)
		}

		data, err := z.SelectFile(packageFullPath)
		if err != nil {
			return err
		}

		var packagePub pkg.Package
		err = xml.Unmarshal(data, &packagePub)
		if err != nil {
			return err
		}

		renditionVars := strings.Join(rendition, "_")
		r.epub.packagePaths[renditionVars] = packageFullPath
		r.epub.packagePubs[renditionVars] = &packagePub
	}

	return nil
}

func OpenReader(name string) (reader Reader, err error) {
	zipContainer, err := ocf.OpenReader(name)
	if err != nil {
		return
	}

	return newReaderFromZip(zipContainer)
}

func NewReader(b []byte) (reader Reader, err error) {
	zipContainer, err := ocf.NewReader(b)
	if err != nil {
		return
	}

	return newReaderFromZip(zipContainer)
}
