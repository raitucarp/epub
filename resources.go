package epub

import (
	"path/filepath"

	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/pkg"
)

type PublicationResource struct {
	ID         string
	Href       string
	MIMEType   string
	Content    []byte
	Filepath   string
	Properties string
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

func (r *Reader) Resources() (resources []PublicationResource) {
	return r.epub.resources
}

func (r *Reader) SelectResourceById(id string) (resource *PublicationResource) {

	for _, res := range r.epub.resources {
		if res.ID == id {
			return &res
		}
	}
	return nil
}

func (r *Reader) SelectResourceByHref(href string) (resource *PublicationResource) {

	for _, res := range r.epub.resources {
		if res.Href == href {
			return &res
		}
	}
	return nil
}
