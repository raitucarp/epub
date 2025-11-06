package epub

import (
	"bytes"
	"encoding/xml"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/ocf"
	"github.com/raitucarp/epub/pkg"
	"golang.org/x/net/html"
)

// Writer provides an interface for constructing, modifying, and writing
// EPUB publications to disk or memory. Writer usage documentation is evolving.
type Writer struct {
	identifier string
	epub       *Epub
	textDir    string
	contentDir string
	imagesDir  string
	direction  string
}

// New creates a new Writer with the given publication identifier.
// The identifier is assigned to the package metadata (dc:identifier).
func New(pubId string) *Writer {
	epubWriter := &Writer{
		identifier: pubId,
		epub: &Epub{
			packagePubs:  make(map[string]*pkg.Package),
			zipContainer: ocf.NewOCFZipContainer(),
		},
		textDir:    "text",
		contentDir: "epub",
		imagesDir:  "images",
		direction:  "ltr",
	}

	epubWriter.epub.rendition = "content"
	epubWriter.epub.packagePubs["content"] = &pkg.Package{
		UniqueIdentifier: "pub-id",
		Version:          "3.0",
		Dir:              epubWriter.direction,
		Metadata:         pkg.Metadata{},
		Spine:            pkg.Spine{TOC: "ncx"},
		Manifest:         pkg.Manifest{},
	}

	primaryIdentifier := pkg.DCIdentifier{ID: "pub-id", Value: epubWriter.identifier}
	epubWriter.epub.SelectedPackage().Metadata.Identifiers = append(
		epubWriter.epub.SelectedPackage().Metadata.Identifiers,
		primaryIdentifier,
	)
	epubWriter.epub.zipContainer.AddMimeType()

	return epubWriter
}

// Direction sets the writing direction (ltr or rtl) used by the spine.
func (w *Writer) Direction(dir string) {
	w.epub.SelectedPackage().Dir = dir
}

// SetContentDir sets the directory used for storing content documents.
func (w *Writer) SetContentDir(dir string) {
	w.contentDir = dir
}

// SetTextDir sets the directory used for text document organization.
func (w *Writer) SetTextDir(dir string) {
	w.textDir = dir
}

// SetImageDir sets the directory used for storing image resources.
func (w *Writer) SetImageDir(dir string) {
	w.imagesDir = dir
}

// Title sets one or more title entries in the metadata.
func (w *Writer) Title(title ...string) {
	if len(title) <= 0 {
		return
	}

	w.epub.SelectedPackage().Metadata.Titles = append(
		w.epub.SelectedPackage().Metadata.Titles,
		pkg.DCTitle{ID: "title", Value: title[0]},
	)

	for _, title := range title[1:] {
		w.epub.SelectedPackage().Metadata.Meta = append(w.epub.SelectedPackage().Metadata.Meta, pkg.Meta{
			Refines: "#title",
			Value:   title,
		})
	}
}

// Description sets a short description or summary for the publication.
func (w *Writer) Description(description string) {
	w.epub.SelectedPackage().Metadata.OptionalDC = append(
		w.epub.SelectedPackage().Metadata.OptionalDC,
		pkg.DCOptional{ID: "description", Value: description, XMLName: xml.Name{Local: "dc:description"}},
	)
}

// Author sets the primary creator/author in the package metadata.
func (w *Writer) Author(creator string) {
	w.epub.SelectedPackage().Metadata.OptionalDC = append(
		w.epub.SelectedPackage().Metadata.OptionalDC,
		pkg.DCOptional{ID: "author", Value: creator, XMLName: xml.Name{Local: "dc:creator"}},
	)
}

// Creator adds a creator with a specific identifier attribute to the metadata.
func (w *Writer) Creator(id string, creator string) {
	w.epub.SelectedPackage().Metadata.OptionalDC = append(
		w.epub.SelectedPackage().Metadata.OptionalDC,
		pkg.DCOptional{ID: id, Value: creator, XMLName: xml.Name{Local: "dc:creator"}},
	)
}

// Contributor adds a contributor entry of the specified role or type.
func (w *Writer) Contributor(kind string, contributor string) {
	w.epub.SelectedPackage().Metadata.OptionalDC = append(
		w.epub.SelectedPackage().Metadata.OptionalDC,
		pkg.DCOptional{ID: kind, Value: contributor, XMLName: xml.Name{Local: "dc:contributor"}},
	)
}

// Subject adds a subject or theme classification to the publication.
func (w *Writer) Subject(id string, subject string) {
	w.epub.SelectedPackage().Metadata.OptionalDC = append(
		w.epub.SelectedPackage().Metadata.OptionalDC,
		pkg.DCOptional{ID: id, Value: subject, XMLName: xml.Name{Local: "dc:subject"}},
	)
}

// LongDescription sets an extended descriptive summary.
func (w *Writer) LongDescription(description string) {
	w.Refines("#description", "long-description", description)
}

// Rights sets the copyright or licensing information for the publication.
func (w *Writer) Rights(rights string) {
	w.DublinCores(map[string]string{"rights": rights})
}

// Date sets the publication date metadata.
func (w *Writer) Date(date time.Time) {
	dateString := date.Format(time.RFC3339)
	w.DublinCores(map[string]string{"date": dateString})
}

// Publisher sets the publication publisher.
func (w *Writer) Publisher(publisher string) {
	w.DublinCores(map[string]string{"publisher": publisher})
}

// DublinCores sets multiple Dublin Core metadata fields at once.
func (w *Writer) DublinCores(keyVal map[string]string) {
	for key, value := range keyVal {
		w.epub.SelectedPackage().Metadata.OptionalDC = append(
			w.epub.SelectedPackage().Metadata.OptionalDC,
			pkg.DCOptional{
				XMLName: xml.Name{
					Local: strings.Join([]string{"dc", key}, ":"),
				},
				ID:    key,
				Value: value,
			},
		)
	}
}

// Meta adds a meta element to the package metadata as-is.
func (w *Writer) Meta(meta pkg.Meta) {
	w.epub.SelectedPackage().Metadata.Meta = append(
		w.epub.SelectedPackage().Metadata.Meta,
		meta,
	)
}

// MetaContent adds metadata key/value entries that do not require refinements.
func (w *Writer) MetaContent(keyVal map[string]string) {
	for key, value := range keyVal {
		w.Meta(
			pkg.Meta{Name: key, Content: value},
		)
	}
}

// MetaProperty adds a property-based metadata refinement entry.
func (w *Writer) MetaProperty(id string, property string, value string) {
	w.Meta(
		pkg.Meta{ID: id, Property: property, Value: value},
	)
}

// Refines applies a metadata refinement to an existing metadata item.
func (w *Writer) Refines(refines string, property string, value string, otherAttributes ...pkg.Meta) {
	meta := pkg.Meta{}
	for _, m := range otherAttributes {
		meta = m
	}

	finalMeta := meta
	finalMeta.ID = property
	finalMeta.Refines = refines
	finalMeta.Property = property
	finalMeta.Value = value

	w.Meta(finalMeta)
}

// Identifiers adds one or more identifiers to the package metadata.
func (w *Writer) Identifiers(identifier ...string) {
	for index, id := range identifier {
		pubId := pkg.DCIdentifier{ID: "pub-id-" + strconv.Itoa(index), Value: id, XMLName: xml.Name{Local: "dc:identifier"}}
		w.epub.SelectedPackage().Metadata.Identifiers = append(w.epub.SelectedPackage().Metadata.Identifiers, pubId)
	}
}

// Languages adds one or more language codes to the publication metadata.
func (w *Writer) Languages(language ...string) {
	if len(language) <= 0 {
		return
	}

	for _, l := range language {
		lang := pkg.DCLanguage{ID: l, Value: l}
		w.epub.SelectedPackage().Metadata.Languages = append(w.epub.SelectedPackage().Metadata.Languages, lang)
	}

}

// AddGuide adds a guide reference entry (e.g., "cover", "toc", "title-page")
// to the package metadata.
func (w *Writer) AddGuide(kind pkg.GuideReferenceType, href string, title string) {
	if w.epub.SelectedPackage().Guide == nil {
		w.epub.SelectedPackage().Guide = &pkg.Guide{}
	}

	w.epub.SelectedPackage().Guide.References = append(
		w.epub.SelectedPackage().Guide.References,
		pkg.GuideReference{Type: kind, Title: title, Href: href},
	)
}

// AddContentFile adds a content file to the publication by reading the file
// from disk. Returns the created resource and any file access error.
func (w *Writer) AddContentFile(name string) (res PublicationResource, err error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return
	}

	res = w.AddContent(name, data)

	return
}

// Cover sets the publication cover from a raw image byte slice.
func (w *Writer) Cover(cover []byte) (err error) {
	name := "cover"
	mime := http.DetectContentType(cover)

	content := name
	switch mime {
	case "image/png":
		content += ".png"

	case "image/jpeg":
	case "image/jpg":
		content += ".jpeg"
	}

	w.addImageCover(content, cover)
	w.MetaContent(map[string]string{name: content})
	return
}

// CoverPNG sets the publication cover image from an image.Image encoded as PNG
func (w *Writer) CoverPNG(cover image.Image) (err error) {
	name := "cover"
	content := name + ".png"

	buf := new(bytes.Buffer)

	err = png.Encode(buf, cover)
	if err != nil {
		return err
	}
	w.addImageCover(content, buf.Bytes())
	w.MetaContent(map[string]string{name: content})
	return
}

// CoverJPG sets the publication cover image from an image.Image encoded as JPEG.
func (w *Writer) CoverJPG(cover image.Image) (err error) {
	name := "cover"
	content := name + ".jpg"

	buf := new(bytes.Buffer)

	err = jpeg.Encode(buf, cover, &jpeg.Options{Quality: 70})
	if err != nil {
		return err
	}
	w.addImageCover(content, buf.Bytes())
	w.MetaContent(map[string]string{name: content})
	return
}

// CoverFile sets the publication cover image by file path.
func (w *Writer) CoverFile(name string) {
	data, err := os.ReadFile(name)
	if err != nil {
		return
	}

	content := name
	w.addImageCover(content, data)
	w.MetaContent(map[string]string{"cover": content})
}

func (w *Writer) addImageCover(name string, content []byte) (res PublicationResource) {
	href := path.Join(w.imagesDir, name)
	filePath := path.Join(w.contentDir, href)
	mimeType := http.DetectContentType(content)
	base := filepath.Base(href)
	res = w.addResource(
		base,
		filePath,
		href,
		pkg.PropertyCoverImage,
		mimeType,
		content,
	)

	return res
}

// AddImage adds an image resource from raw bytes to the publication.
func (w *Writer) AddImage(name string, content []byte) (res PublicationResource) {
	href := path.Join(w.imagesDir, name)
	filePath := path.Join(w.contentDir, href)
	mimeType := http.DetectContentType(content)
	base := filepath.Base(href)
	res = w.addResource(
		base,
		filePath,
		href,
		pkg.NotProperty,
		mimeType,
		content,
	)

	return res
}

// AddImageFile adds an image resource to the publication by reading from disk.
func (w *Writer) AddImageFile(name string) (res PublicationResource) {
	data, err := os.ReadFile(name)
	if err != nil {
		return
	}

	res = w.AddImage(name, data)

	return
}

// AddContent adds a content file (such as XHTML or SVG) to the publication
// using the provided filename and raw bytes. Returns the created resource.
func (w *Writer) AddContent(filename string, content []byte) (res PublicationResource) {
	href := filename
	filePath := path.Join(w.contentDir, href)
	mimeType := pkg.MediaTypeXHTML
	base := filepath.Base(href)
	res = w.addResource(
		base,
		filePath,
		href,
		pkg.NotProperty,
		mimeType,
		content,
	)

	w.AddSpineItem(res)
	return
}

func (w *Writer) addResource(
	id string,
	filePath string,
	href string,
	properties pkg.ManifestProperty,
	mimeType string,
	content []byte,
) (pubRes PublicationResource) {
	pubRes = PublicationResource{
		ID:         id,
		Filepath:   filePath,
		Href:       href,
		Properties: properties,
		MIMEType:   mimeType,
		Content:    content,
	}

	w.epub.resources = append(w.epub.resources, pubRes)
	w.epub.zipContainer.AddFile(pubRes.Filepath, content)
	w.epub.SelectedPackage().Manifest.Items = append(
		w.epub.SelectedPackage().Manifest.Items,
		pkg.Item{
			ID:         pubRes.ID,
			Href:       pubRes.Href,
			MediaType:  pubRes.MIMEType,
			Properties: pubRes.Properties,
		},
	)

	return pubRes
}

// AddSpineItem appends the given resource to the spine reading order.
func (w *Writer) AddSpineItem(res PublicationResource) {
	itemRef := pkg.ItemRef{IDRef: res.ID}
	w.epub.SelectedPackage().Spine.ItemRefs = append(
		w.epub.SelectedPackage().Spine.ItemRefs,
		itemRef,
	)
}

func (w *Writer) TableOfContents(name string, toc TOC) (err error) {
	navigation := ncx.NCX{
		NavMap: ncx.NavMap{
			ID:        "navmap",
			NavPoints: []ncx.NavPoint{},
		},
	}
	currentNavPoint := 0
	visitTOC(&toc, func(t *TOC, depth int) {
		if t.Href == "" && depth == 0 {
			navigation.DocTitle.Text = t.Title
			return
		}

		label := ncx.NavLabel{Text: t.Title}
		playOrder := strconv.Itoa(currentNavPoint + 1)
		navPoint := ncx.NavPoint{
			ID:        "nav-point-" + playOrder,
			PlayOrder: playOrder,
			NavLabel:  label,
			Content:   ncx.Content{Src: t.Href},
			NavPoints: []ncx.NavPoint{},
		}

		if depth <= 1 {
			navigation.NavMap.NavPoints = append(navigation.NavMap.NavPoints, navPoint)
		} else {
			navigation.NavMap.NavPoints[depth+1].NavPoints = append(
				navigation.NavMap.NavPoints[depth+1].NavPoints,
				navPoint,
			)
		}

		currentNavPoint++

	})

	w.epub.navigationCenterEXtended = &navigation
	ncxContent, err := xml.MarshalIndent(navigation, "", " ")
	ncxBase := name + ".ncx"
	ncxFilePath := path.Join(w.contentDir, ncxBase)
	w.addResource(
		name,
		ncxFilePath,
		ncxBase,
		pkg.NotProperty,
		pkg.MediaTypeNCX,
		ncxContent,
	)

	filePath := path.Join(w.contentDir, name+".xhtml")
	base := filepath.Base(filePath)
	href := base
	languages := []string{}
	for _, lang := range w.epub.SelectedPackage().Metadata.Languages {
		languages = append(languages, lang.Value)
	}
	xhtmlContent, err := tocToHTMLNode(toc, languages)
	if err != nil {
		return
	}

	var content bytes.Buffer
	html.Render(&content, xhtmlContent)
	w.addResource(
		base,
		filePath,
		href,
		pkg.NavProperty,
		pkg.MediaTypeXHTML,
		content.Bytes(),
	)
	return
}

func (w *Writer) guardCheck() (err error) {
	for _, p := range w.epub.packagePubs {

		if len(p.Metadata.Identifiers) <= 0 {
			return errors.New("Package should have identifiers.")
		}

		if len(p.Metadata.Titles) <= 0 {
			return errors.New("Package should have titles.")
		}

		if len(p.Metadata.Languages) <= 0 {
			return errors.New("Package should have languages.")
		}

		if len(p.Manifest.Items) <= 0 {
			return errors.New("No content insides.")
		} else {
			var content int
			var cover int
			for _, item := range p.Manifest.Items {
				if item.MediaType == pkg.MediaTypeXHTML {
					content++
				}

				if item.Properties == "cover-image" {
					cover++
				}
			}

			if content <= 0 {
				return errors.New("No text content insides.")
			}

			if cover <= 0 {
				return errors.New("No cover images.")
			}
		}

	}

	if w.epub.navigationCenterEXtended == nil {
		return errors.New("No table of contents.")
	}

	return
}

// Write finalizes the EPUB structure and writes it to the specified filename.
func (w *Writer) Write(filename string) (err error) {
	err = w.guardCheck()
	if err != nil {
		return err
	}

	rootFiles := []string{}
	for name, p := range w.epub.packagePubs {
		containerFilePath := path.Join(w.contentDir, name+".opf")
		w.epub.zipContainer.AddPackage(containerFilePath, *p)
		rootFiles = append(rootFiles, containerFilePath)
	}

	w.epub.zipContainer.AddContainerXML(rootFiles...)

	err = w.epub.zipContainer.Write(filename)
	if err != nil {
		return
	}

	return
}
