package epub

import (
	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/ocf"
	"github.com/raitucarp/epub/pkg"
)

// Epub represents a full EPUB publication, including
// metadata, package information, resources, and navigation.
type Epub struct {
	packagePubs              map[string]*pkg.Package
	packagePaths             map[string]string
	zipContainer             *ocf.OCFZipContainer
	rendition                string
	resources                []PublicationResource
	metadata                 map[string]any
	navigationCenterEXtended *ncx.NCX
}

func (epub *Epub) SelectPackage(name string) *pkg.Package {
	return epub.packagePubs[name]
}

func (epub *Epub) SelectedPackage() *pkg.Package {
	return epub.packagePubs[epub.rendition]
}

func (epub *Epub) DefaultPackage() *pkg.Package {
	return epub.packagePubs["content"]
}
