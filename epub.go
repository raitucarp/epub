package epub

import (
	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/ocf"
	"github.com/raitucarp/epub/pkg"
)

type Epub struct {
	packagePubs              map[string]*pkg.Package
	packagePaths             map[string]string
	zipContainer             *ocf.OCFZipContainer
	rendition                string
	resources                []PublicationResource
	metadata                 map[string]any
	navigationCenterEXtended *ncx.NCX
}
