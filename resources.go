package epub

import (
	"path/filepath"

	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/pkg"
)

// PublicationResource represents a single resource inside the EPUB
// container. Resources may include XHTML documents, images, stylesheets,
// SVG files, and auxiliary data referenced by the publication manifest.
type PublicationResource struct {
	ID         string
	Href       string
	MIMEType   string
	Content    []byte
	Filepath   string
	Properties pkg.ManifestProperty
}

func (r *Reader) parseResources() {
	allFiles := r.epub.zipContainer.AllFiles()
	currentPackagePath := r.CurrentSelectedPackagePath()
	for _, item := range r.CurrentSelectedPackage().Manifest.Items {
		itemPath := filepath.ToSlash(
			filepath.Clean(
				filepath.Join(filepath.Dir(currentPackagePath), item.Href),
			),
		)

		content := allFiles[itemPath]

		if item.MediaType == pkg.MediaTypeNCX {
			r.epub.navigationCenterEXtended, _ = ncx.Parse(content)
		}
		r.epub.resources = append(r.epub.resources, PublicationResource{
			ID:         item.ID,
			Href:       item.Href,
			MIMEType:   item.MediaType,
			Content:    content,
			Filepath:   itemPath,
			Properties: item.Properties,
		})
	}
}

// Resources returns all publication resources declared in the manifest.
func (r *Reader) Resources() (resources []PublicationResource) {
	return r.epub.resources
}

// SelectResourceById retrieves a resource referenced by its manifest ID.
func (r *Reader) SelectResourceById(id string) (resource *PublicationResource) {

	for _, res := range r.epub.resources {
		if res.ID == id {
			return &res
		}
	}
	return nil
}

// SelectResourceByHref retrieves a resource referenced by its manifest href.
func (r *Reader) SelectResourceByHref(href string) (resource *PublicationResource) {

	for _, res := range r.epub.resources {
		if res.Href == href {
			return &res
		}
	}
	return nil
}
