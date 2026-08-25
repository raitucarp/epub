package epub

import (
	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/ocf"
	"github.com/raitucarp/epub/pkg"
)

// Epub represents a full EPUB publication, including metadata, package
// information, resources, and navigation.
type Epub struct {
	packagePubs              map[string]*pkg.Package
	packagePaths             map[string]string
	zipContainer             *ocf.OCFZipContainer
	rendition                string
	resources                []PublicationResource
	metadata                 map[string]any
	navigationCenterEXtended *ncx.NCX
}

// SelectPackage returns the package document registered under the given name,
// or nil if no such package exists. Use Reader.SelectPackageRendition to switch
// the active package instead.
func (epub *Epub) SelectPackage(name string) *pkg.Package {
	return epub.packagePubs[name]
}

// SelectedPackage returns the currently active package document.
func (epub *Epub) SelectedPackage() *pkg.Package {
	return epub.packagePubs[epub.rendition]
}

// DefaultPackage returns the package document registered under the conventional
// "content" key.
func (epub *Epub) DefaultPackage() *pkg.Package {
	return epub.packagePubs["content"]
}
