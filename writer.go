package epub

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
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

// Writer provides an interface for constructing and writing EPUB publications
// to disk or memory.
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

	for _, alt := range title[1:] {
		w.epub.SelectedPackage().Metadata.Titles = append(
			w.epub.SelectedPackage().Metadata.Titles,
			pkg.DCTitle{Value: alt},
		)
	}
}

// Description sets a short description or summary for the publication.
func (w *Writer) Description(description ...string) {
	for _, d := range description {
		w.epub.SelectedPackage().Metadata.OptionalDC = append(
			w.epub.SelectedPackage().Metadata.OptionalDC,
			pkg.DCOptional{ID: "description", Value: d, XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "description"}},
		)
	}
}

// Author sets the primary creator/author in the package metadata.
func (w *Writer) Author(creator ...string) {
	for _, c := range creator {
		w.epub.SelectedPackage().Metadata.OptionalDC = append(
			w.epub.SelectedPackage().Metadata.OptionalDC,
			pkg.DCOptional{ID: "author", Value: c, XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "creator"}},
		)
	}
}

// Creator adds a creator with a specific identifier attribute to the metadata.
func (w *Writer) Creator(id string, creator ...string) {
	for _, c := range creator {
		w.epub.SelectedPackage().Metadata.OptionalDC = append(
			w.epub.SelectedPackage().Metadata.OptionalDC,
			pkg.DCOptional{ID: id, Value: c, XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "creator"}},
		)
	}
}

// Contributor adds a contributor entry of the specified role or type.
func (w *Writer) Contributor(kind string, contributor ...string) {
	for _, c := range contributor {
		w.epub.SelectedPackage().Metadata.OptionalDC = append(
			w.epub.SelectedPackage().Metadata.OptionalDC,
			pkg.DCOptional{ID: kind, Value: c, XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "contributor"}},
		)
	}
}

// Subject adds a subject or theme classification to the publication.
func (w *Writer) Subject(id string, subject ...string) {
	for _, s := range subject {
		w.epub.SelectedPackage().Metadata.OptionalDC = append(
			w.epub.SelectedPackage().Metadata.OptionalDC,
			pkg.DCOptional{ID: id, Value: s, XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "subject"}},
		)
	}
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

// Modified sets the last modification date (dcterms:modified) metadata, which
// the EPUB specification requires in every package document. The value is
// expressed in UTC using the extended ISO 8601 format (YYYY-MM-DDThh:mm:ssZ).
func (w *Writer) Modified(date time.Time) {
	w.Meta(pkg.Meta{Property: "dcterms:modified", Value: date.UTC().Format("2006-01-02T15:04:05Z")})
}

// ensureModifiedDate adds a dcterms:modified property if one is not already
// present, as required by the EPUB specification.
func (w *Writer) ensureModifiedDate() {
	for _, meta := range w.epub.SelectedPackage().Metadata.Meta {
		if meta.Property == "dcterms:modified" && meta.Refines == "" {
			return
		}
	}

	w.Modified(time.Now())
}

// Publisher sets the publication publisher.
func (w *Writer) Publisher(publisher ...string) {
	for _, p := range publisher {
		w.DublinCores(map[string]string{"publisher": p})
	}
}

// DublinCores sets multiple Dublin Core metadata fields at once.
func (w *Writer) DublinCores(keyVal map[string]string) {
	for key, value := range keyVal {
		w.epub.SelectedPackage().Metadata.OptionalDC = append(
			w.epub.SelectedPackage().Metadata.OptionalDC,
			pkg.DCOptional{
				XMLName: xml.Name{
					Space: pkg.NamespaceDC,
					Local: key,
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
		pubId := pkg.DCIdentifier{ID: "pub-id-" + strconv.Itoa(index), Value: id}
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
	if !filepath.IsLocal(name) {
		return res, fmt.Errorf("invalid path: path must be local")
	}

	root, err := os.OpenRoot(".")
	if err != nil {
		return res, err
	}
	defer root.Close()

	data, err := root.ReadFile(name)
	if err != nil {
		return
	}

	res = w.AddContent(filepath.Base(name), data)

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

// CoverPNG sets the publication cover image from an image.Image encoded as PNG.
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
func (w *Writer) CoverFile(name string) error {
	if !filepath.IsLocal(name) {
		return fmt.Errorf("invalid path: path must be local")
	}

	root, err := os.OpenRoot(".")
	if err != nil {
		return err
	}
	defer root.Close()

	data, err := root.ReadFile(name)
	if err != nil {
		return err
	}

	content := filepath.Base(name)
	w.addImageCover(content, data)
	w.MetaContent(map[string]string{"cover": content})
	return nil
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
	if !filepath.IsLocal(name) {
		return
	}

	root, err := os.OpenRoot(".")
	if err != nil {
		return
	}
	defer root.Close()

	data, err := root.ReadFile(name)
	if err != nil {
		return
	}

	res = w.AddImage(filepath.Base(name), data)

	return
}

// AddContent adds a content file (such as XHTML or SVG) to the publication
// using the provided filename and raw bytes. Returns the created resource.
// The media type is inferred from the filename extension; unknown extensions
// default to XHTML.
func (w *Writer) AddContent(filename string, content []byte) (res PublicationResource) {
	href := filename
	filePath := path.Join(w.contentDir, href)
	mimeType := detectContentMediaType(filename, content)
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

// detectContentMediaType returns the media type for a content document based
// on its filename extension, falling back to XHTML for unknown extensions.
func detectContentMediaType(filename string, content []byte) string {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".svg":
		return pkg.MediaTypeSVG
	case ".xhtml", ".html", ".htm", ".xml":
		return pkg.MediaTypeXHTML
	default:
		if mime := http.DetectContentType(content); strings.HasPrefix(mime, "image/svg") {
			return pkg.MediaTypeSVG
		}
		return pkg.MediaTypeXHTML
	}
}

// AddResource adds an arbitrary resource (such as a stylesheet, font, or asset)
// to the publication manifest and container.
func (w *Writer) AddResource(id, href, mimeType string, properties pkg.ManifestProperty, content []byte) (PublicationResource, error) {
	if id == "" {
		id = filepath.Base(href)
	}
	filePath := path.Join(w.contentDir, href)
	res := w.addResource(id, filePath, href, properties, mimeType, content)
	return res, nil
}

func (w *Writer) addResource(
	id string,
	filePath string,
	href string,
	properties pkg.ManifestProperty,
	mimeType string,
	content []byte,
) (pubRes PublicationResource) {
	id = w.uniqueID(id)

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

// uniqueID returns an identifier that is unique within the manifest, appending
// a numeric suffix when the requested identifier is already in use. EPUB
// manifest item ids must be unique.
func (w *Writer) uniqueID(id string) string {
	used := false
	for _, item := range w.epub.SelectedPackage().Manifest.Items {
		if item.ID == id {
			used = true
			break
		}
	}
	if !used {
		return id
	}

	base := id
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		taken := false
		for _, item := range w.epub.SelectedPackage().Manifest.Items {
			if item.ID == candidate {
				taken = true
				break
			}
		}
		if !taken {
			return candidate
		}
	}
}

// AddSpineItem appends the given resource to the spine reading order.
func (w *Writer) AddSpineItem(res PublicationResource) {
	itemRef := pkg.ItemRef{IDRef: res.ID}
	w.epub.SelectedPackage().Spine.ItemRefs = append(
		w.epub.SelectedPackage().Spine.ItemRefs,
		itemRef,
	)
}

// uniqueResourceName returns a name whose derived nav (name.xhtml) and NCX
// (name.ncx) hrefs do not collide with any existing manifest item href. This
// prevents TableOfContents from silently overwriting a content document that
// happens to share the generated navigation document's file name.
func (w *Writer) uniqueResourceName(name string) string {
	collides := func(n string) bool {
		for _, item := range w.epub.SelectedPackage().Manifest.Items {
			if item.Href == n+".xhtml" || item.Href == n+".ncx" {
				return true
			}
		}
		return false
	}
	if !collides(name) {
		return name
	}
	base := name
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if !collides(candidate) {
			return candidate
		}
	}
}

// TableOfContents generates the EPUB navigation document and its legacy NCX
// counterpart from the given TOC structure, using name as the base for the
// generated file names.
func (w *Writer) TableOfContents(name string, toc TOC) (err error) {
	name = w.uniqueResourceName(name)

	navigation := ncx.NCX{
		NavMap: ncx.NavMap{
			ID:        "navmap",
			NavPoints: []ncx.NavPoint{},
		},
	}
	navigation.DocTitle.Text = toc.Title
	navigation.NavMap.NavPoints = buildNavPoints(toc.Items, 0)

	w.epub.navigationCenterEXtended = &navigation
	ncxContent, err := xml.MarshalIndent(navigation, "", " ")
	ncxBase := name + ".ncx"
	ncxFilePath := path.Join(w.contentDir, ncxBase)
	ncxRes := w.addResource(
		name,
		ncxFilePath,
		ncxBase,
		pkg.NotProperty,
		pkg.MediaTypeNCX,
		ncxContent,
	)
	w.epub.SelectedPackage().Spine.TOC = ncxRes.ID

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

// buildNavPoints recursively converts TOC items into NCX navPoints, assigning
// sequential playOrder and id values in document order.
func buildNavPoints(items []TOC, playOrder int) []ncx.NavPoint {
	counter := playOrder
	var build func([]TOC) []ncx.NavPoint
	build = func(list []TOC) []ncx.NavPoint {
		var points []ncx.NavPoint
		for _, item := range list {
			counter++
			order := strconv.Itoa(counter)
			point := ncx.NavPoint{
				ID:        "nav-point-" + order,
				PlayOrder: order,
				NavLabel:  ncx.NavLabel{Text: item.Title},
				Content:   ncx.Content{Src: item.Href},
				NavPoints: []ncx.NavPoint{},
			}
			if len(item.Items) > 0 {
				point.NavPoints = build(item.Items)
			}
			points = append(points, point)
		}
		return points
	}
	return build(items)
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
			for _, item := range p.Manifest.Items {
				if isXHTMLContent(item.MediaType) {
					content++
				}
			}

			if content <= 0 {
				return errors.New("No text content insides.")
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
	data, err := w.Build()
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0o644)
}

// WriteBytes finalizes the EPUB structure and returns the resulting bytes
// without touching the filesystem.
func (w *Writer) WriteBytes() ([]byte, error) {
	return w.Build()
}

// Build assembles the package and container documents and serializes the
// publication to an in-memory ZIP archive.
func (w *Writer) Build() ([]byte, error) {
	if err := w.guardCheck(); err != nil {
		return nil, err
	}

	w.ensureModifiedDate()

	rootFiles := make([]string, 0, len(w.epub.packagePubs))
	for name, p := range w.epub.packagePubs {
		containerFilePath := path.Join(w.contentDir, name+".opf")
		w.epub.zipContainer.AddPackage(containerFilePath, *p)
		rootFiles = append(rootFiles, containerFilePath)
	}

	w.epub.zipContainer.AddContainerXML(rootFiles...)

	var buf bytes.Buffer
	if err := w.epub.zipContainer.WriteArchive(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
