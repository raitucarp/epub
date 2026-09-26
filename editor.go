package epub

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/ocf"
	"github.com/raitucarp/epub/pkg"
	"golang.org/x/net/html"
)

// Editor provides an interface for mutating and updating an existing EPUB publication,
// including its metadata, resources, spine reading order, and table of contents, and
// saving the edited publication back to disk, an io.Writer, or a byte buffer.
type Editor struct {
	zipContainer *ocf.OCFZipContainer
	packagePubs  map[string]*pkg.Package
	packagePaths map[string]string
	rendition    string
	modifiedSet  bool
}

// newEditor clones the underlying container and package data from Reader to create
// an independent Editor.
func newEditor(r *Reader) (*Editor, error) {
	clonedZip := r.epub.zipContainer.Clone()
	if clonedZip == nil {
		return nil, errors.New("epub: failed to clone zip container")
	}

	clonedPubs := make(map[string]*pkg.Package, len(r.epub.packagePubs))
	for k, packageItem := range r.epub.packagePubs {
		b, err := xml.Marshal(packageItem)
		if err != nil {
			return nil, fmt.Errorf("failed to clone package %s: %w", k, err)
		}
		var cloned pkg.Package
		if err := xml.Unmarshal(b, &cloned); err != nil {
			return nil, fmt.Errorf("failed to unmarshal package %s: %w", k, err)
		}
		clonedPubs[k] = &cloned
	}

	if len(clonedPubs) == 0 {
		return nil, errors.New("epub: no package document found")
	}

	clonedPaths := make(map[string]string, len(r.epub.packagePaths))
	for k, v := range r.epub.packagePaths {
		clonedPaths[k] = v
	}

	rendition := r.epub.rendition
	if rendition == "" {
		rendition = "default"
	}
	if _, ok := clonedPubs[rendition]; !ok {
		for k := range clonedPubs {
			rendition = k
			break
		}
	}

	return &Editor{
		zipContainer: clonedZip,
		packagePubs:  clonedPubs,
		packagePaths: clonedPaths,
		rendition:    rendition,
	}, nil
}

// Package returns the currently active package document for direct inspection or modification.
func (e *Editor) Package() *pkg.Package {
	return e.CurrentPackage()
}

// CurrentPackage returns the active package publication.
func (e *Editor) CurrentPackage() *pkg.Package {
	return e.packagePubs[e.rendition]
}

// CurrentPackagePath returns the container file path of the active package document.
func (e *Editor) CurrentPackagePath() string {
	return e.packagePaths[e.rendition]
}

// SelectPackageRendition switches the active package rendition.
func (e *Editor) SelectPackageRendition(rendition string) error {
	if _, ok := e.packagePubs[rendition]; !ok {
		return fmt.Errorf("rendition '%s' not found", rendition)
	}
	e.rendition = rendition
	return nil
}

// ListRenditions returns the identifiers of all package renditions available.
func (e *Editor) ListRenditions() []string {
	renditions := make([]string, 0, len(e.packagePubs))
	for r := range e.packagePubs {
		renditions = append(renditions, r)
	}
	slices.Sort(renditions)
	return renditions
}

// Title sets (replaces) the publication title(s).
func (e *Editor) Title(titles ...string) *Editor {
	p := e.CurrentPackage()
	if p == nil || len(titles) == 0 {
		return e
	}
	p.Metadata.Titles = nil
	for i, t := range titles {
		id := "title"
		if i > 0 {
			id = fmt.Sprintf("title-%d", i)
		}
		p.Metadata.Titles = append(p.Metadata.Titles, pkg.DCTitle{
			ID:    id,
			Value: t,
		})
	}
	return e
}

// AddTitle appends an additional title entry to the metadata.
func (e *Editor) AddTitle(title string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	id := fmt.Sprintf("title-%d", len(p.Metadata.Titles))
	p.Metadata.Titles = append(p.Metadata.Titles, pkg.DCTitle{
		ID:    id,
		Value: title,
	})
	return e
}

// Author sets (replaces) the publication creator/author(s).
func (e *Editor) Author(authors ...string) *Editor {
	e.RemoveMetadata("creator")
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	for i, a := range authors {
		id := "author"
		if i > 0 {
			id = fmt.Sprintf("author-%d", i)
		}
		p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
			XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "creator"},
			ID:      id,
			Value:   a,
		})
	}
	return e
}

// AddAuthor appends an additional author/creator to the metadata.
func (e *Editor) AddAuthor(author string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	id := "author"
	p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
		XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "creator"},
		ID:      id,
		Value:   author,
	})
	return e
}

// Creator adds a creator with a specific identifier attribute to the metadata.
func (e *Editor) Creator(id string, creators ...string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	for _, c := range creators {
		p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
			XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "creator"},
			ID:      id,
			Value:   c,
		})
	}
	return e
}

// Description sets (replaces) the publication description(s).
func (e *Editor) Description(descriptions ...string) *Editor {
	e.RemoveMetadata("description")
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	for i, d := range descriptions {
		id := "description"
		if i > 0 {
			id = fmt.Sprintf("description-%d", i)
		}
		p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
			XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "description"},
			ID:      id,
			Value:   d,
		})
	}
	return e
}

// AddDescription appends a description to the metadata.
func (e *Editor) AddDescription(description string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
		XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "description"},
		ID:      "description",
		Value:   description,
	})
	return e
}

// Publisher sets (replaces) the publication publisher(s).
func (e *Editor) Publisher(publishers ...string) *Editor {
	e.RemoveMetadata("publisher")
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	for _, pub := range publishers {
		p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
			XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "publisher"},
			ID:      "publisher",
			Value:   pub,
		})
	}
	return e
}

// AddPublisher appends a publisher to the metadata.
func (e *Editor) AddPublisher(publisher string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
		XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "publisher"},
		ID:      "publisher",
		Value:   publisher,
	})
	return e
}

// Contributor adds a contributor entry of the specified role or type.
func (e *Editor) Contributor(kind string, contributors ...string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	for _, c := range contributors {
		p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
			XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "contributor"},
			ID:      kind,
			Value:   c,
		})
	}
	return e
}

// Subject sets (replaces) the publication subjects/themes.
func (e *Editor) Subject(subjects ...string) *Editor {
	e.RemoveMetadata("subject")
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	for i, s := range subjects {
		id := "subject"
		if i > 0 {
			id = fmt.Sprintf("subject-%d", i)
		}
		p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
			XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "subject"},
			ID:      id,
			Value:   s,
		})
	}
	return e
}

// AddSubject appends a subject to the metadata.
func (e *Editor) AddSubject(subject string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
		XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "subject"},
		ID:      "subject",
		Value:   subject,
	})
	return e
}

// Rights sets the copyright or licensing information.
func (e *Editor) Rights(rights string) *Editor {
	e.RemoveMetadata("rights")
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
		XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "rights"},
		ID:      "rights",
		Value:   rights,
	})
	return e
}

// Date sets the publication date metadata.
func (e *Editor) Date(date time.Time) *Editor {
	e.RemoveMetadata("date")
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
		XMLName: xml.Name{Space: pkg.NamespaceDC, Local: "date"},
		ID:      "date",
		Value:   date.Format(time.RFC3339),
	})
	return e
}

// Modified sets the last modification date (dcterms:modified) metadata.
func (e *Editor) Modified(date time.Time) *Editor {
	e.modifiedSet = true
	val := date.UTC().Format("2006-01-02T15:04:05Z")
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	found := false
	for i, m := range p.Metadata.Meta {
		if m.Property == "dcterms:modified" && m.Refines == "" {
			p.Metadata.Meta[i].Value = val
			found = true
			break
		}
	}
	if !found {
		p.Metadata.Meta = append(p.Metadata.Meta, pkg.Meta{
			Property: "dcterms:modified",
			Value:    val,
		})
	}
	return e
}

func (e *Editor) ensureModifiedDate() {
	if !e.modifiedSet {
		e.Modified(time.Now())
	}
}

// Language sets (replaces) the publication language codes.
func (e *Editor) Language(languages ...string) *Editor {
	p := e.CurrentPackage()
	if p == nil || len(languages) == 0 {
		return e
	}
	p.Metadata.Languages = nil
	for _, l := range languages {
		p.Metadata.Languages = append(p.Metadata.Languages, pkg.DCLanguage{
			ID:    l,
			Value: l,
		})
	}
	return e
}

// AddLanguage appends a language code to the publication metadata.
func (e *Editor) AddLanguage(language string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	p.Metadata.Languages = append(p.Metadata.Languages, pkg.DCLanguage{
		ID:    language,
		Value: language,
	})
	return e
}

// Identifier sets (replaces) the publication identifiers.
func (e *Editor) Identifier(identifiers ...string) *Editor {
	p := e.CurrentPackage()
	if p == nil || len(identifiers) == 0 {
		return e
	}
	p.Metadata.Identifiers = nil
	for i, id := range identifiers {
		pubId := fmt.Sprintf("pub-id-%d", i)
		if i == 0 && p.UniqueIdentifier != "" {
			pubId = p.UniqueIdentifier
		}
		p.Metadata.Identifiers = append(p.Metadata.Identifiers, pkg.DCIdentifier{
			ID:    pubId,
			Value: id,
		})
	}
	return e
}

// AddIdentifier appends an identifier to the package metadata.
func (e *Editor) AddIdentifier(id, value string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	p.Metadata.Identifiers = append(p.Metadata.Identifiers, pkg.DCIdentifier{
		ID:    id,
		Value: value,
	})
	return e
}

// UniqueIdentifier sets the package unique-identifier attribute.
func (e *Editor) UniqueIdentifier(id string) *Editor {
	p := e.CurrentPackage()
	if p != nil {
		p.UniqueIdentifier = id
	}
	return e
}

// Version sets the EPUB publication version (e.g. "3.0" or "2.0").
func (e *Editor) Version(v string) *Editor {
	p := e.CurrentPackage()
	if p != nil {
		p.Version = v
	}
	return e
}

// Direction sets the reading direction ("ltr" or "rtl").
func (e *Editor) Direction(dir string) *Editor {
	p := e.CurrentPackage()
	if p != nil {
		p.Dir = dir
		p.Spine.PageProgressionDirection = dir
	}
	return e
}

// DublinCore adds an arbitrary Dublin Core element to the metadata.
func (e *Editor) DublinCore(localName, value string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	p.Metadata.OptionalDC = append(p.Metadata.OptionalDC, pkg.DCOptional{
		XMLName: xml.Name{Space: pkg.NamespaceDC, Local: localName},
		ID:      localName,
		Value:   value,
	})
	return e
}

// RemoveMetadata removes optional Dublin Core elements matching the given XML local name.
func (e *Editor) RemoveMetadata(xmlLocalName string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	p.Metadata.OptionalDC = slices.DeleteFunc(p.Metadata.OptionalDC, func(dc pkg.DCOptional) bool {
		return dc.XMLName.Local == xmlLocalName
	})
	return e
}

// Meta adds a meta element to the package metadata.
func (e *Editor) Meta(meta pkg.Meta) *Editor {
	p := e.CurrentPackage()
	if p != nil {
		p.Metadata.Meta = append(p.Metadata.Meta, meta)
	}
	return e
}

// SetMeta sets or updates a meta element matching the specified property.
func (e *Editor) SetMeta(property, value string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	found := false
	for i, m := range p.Metadata.Meta {
		if m.Property == property {
			p.Metadata.Meta[i].Value = value
			found = true
			break
		}
	}
	if !found {
		p.Metadata.Meta = append(p.Metadata.Meta, pkg.Meta{
			Property: property,
			Value:    value,
		})
	}
	return e
}

// MetaContent adds metadata key/value entries that do not require refinements.
func (e *Editor) MetaContent(keyVal map[string]string) *Editor {
	for k, v := range keyVal {
		e.Meta(pkg.Meta{Name: k, Content: v})
	}
	return e
}

// MetaProperty adds a property-based metadata refinement entry.
func (e *Editor) MetaProperty(id, property, value string) *Editor {
	return e.Meta(pkg.Meta{ID: id, Property: property, Value: value})
}

// Refines applies a metadata refinement to an existing metadata item.
func (e *Editor) Refines(refines, property, value string, otherAttributes ...pkg.Meta) *Editor {
	meta := pkg.Meta{}
	for _, m := range otherAttributes {
		meta = m
	}
	meta.ID = property
	meta.Refines = refines
	meta.Property = property
	meta.Value = value
	return e.Meta(meta)
}

// RemoveMeta removes meta elements whose Property or Name matches propertyOrName.
func (e *Editor) RemoveMeta(propertyOrName string) *Editor {
	p := e.CurrentPackage()
	if p == nil {
		return e
	}
	p.Metadata.Meta = slices.DeleteFunc(p.Metadata.Meta, func(m pkg.Meta) bool {
		return m.Property == propertyOrName || m.Name == propertyOrName
	})
	return e
}

func (e *Editor) resolveContainerPath(href string) string {
	packageDir := path.Dir(e.CurrentPackagePath())
	if packageDir == "." || packageDir == "" {
		return href
	}
	return path.Clean(path.Join(packageDir, href))
}

func (e *Editor) uniqueID(id string) string {
	p := e.CurrentPackage()
	used := false
	for _, item := range p.Manifest.Items {
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
		for _, item := range p.Manifest.Items {
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

// AddFile directly adds or updates any arbitrary file path inside the container.
func (e *Editor) AddFile(filePath string, content []byte) *Editor {
	e.zipContainer.AddFile(filePath, content)
	return e
}

// AddResource adds a resource to the publication manifest and container.
func (e *Editor) AddResource(id, href, mimeType string, properties pkg.ManifestProperty, content []byte) (PublicationResource, error) {
	p := e.CurrentPackage()
	if p == nil {
		return PublicationResource{}, errors.New("no active package selected")
	}

	if id == "" {
		id = strings.TrimSuffix(filepath.Base(href), filepath.Ext(href))
	}
	id = e.uniqueID(id)

	filePath := e.resolveContainerPath(href)
	res := PublicationResource{
		ID:         id,
		Href:       href,
		MIMEType:   mimeType,
		Properties: properties,
		Filepath:   filePath,
		Content:    content,
	}

	e.zipContainer.AddFile(filePath, content)
	p.Manifest.Items = append(p.Manifest.Items, pkg.Item{
		ID:         id,
		Href:       href,
		MediaType:  mimeType,
		Properties: properties,
	})

	return res, nil
}

// AddContent adds a content document (e.g. XHTML) to the publication manifest, spine, and container.
func (e *Editor) AddContent(filename string, content []byte) (PublicationResource, error) {
	mimeType := detectContentMediaType(filename, content)
	res, err := e.AddResource("", filename, mimeType, pkg.NotProperty, content)
	if err != nil {
		return res, err
	}

	e.CurrentPackage().Spine.ItemRefs = append(e.CurrentPackage().Spine.ItemRefs, pkg.ItemRef{IDRef: res.ID})
	return res, nil
}

// AddContentFile adds a content document by reading it from the local filesystem.
func (e *Editor) AddContentFile(name string) (PublicationResource, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return PublicationResource{}, err
	}
	return e.AddContent(filepath.Base(name), data)
}

// AddMarkdown converts content from Markdown to compliant XHTML and adds it to the
// publication as a spine content document. The file extension is automatically converted to .xhtml.
func (e *Editor) AddMarkdown(filename string, content []byte) (PublicationResource, error) {
	doc, err := markdownToXHTML(content)
	if err != nil {
		return PublicationResource{}, err
	}
	return e.AddContent(replaceExt(filename, ".xhtml"), []byte(doc.body))
}

// AddMarkdownFile reads the Markdown file at name from disk, converts it to compliant XHTML,
// and adds it to the publication as a spine content document.
func (e *Editor) AddMarkdownFile(name string) (PublicationResource, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return PublicationResource{}, err
	}
	return e.AddMarkdown(filepath.Base(name), data)
}

// AddImage adds an image resource to the publication manifest and container.
func (e *Editor) AddImage(name string, content []byte) (PublicationResource, error) {
	mimeType := http.DetectContentType(content)
	return e.AddResource("", name, mimeType, pkg.NotProperty, content)
}

// AddImageFile adds an image resource by reading it from the local filesystem.
func (e *Editor) AddImageFile(name string) (PublicationResource, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return PublicationResource{}, err
	}
	return e.AddImage(filepath.Base(name), data)
}

// UpdateResource updates the byte content of an existing resource identified by ID or Href.
func (e *Editor) UpdateResource(idOrHref string, content []byte) error {
	p := e.CurrentPackage()
	if p == nil {
		return errors.New("no active package selected")
	}

	var foundItem *pkg.Item
	for i, item := range p.Manifest.Items {
		if item.ID == idOrHref || item.Href == idOrHref {
			foundItem = &p.Manifest.Items[i]
			break
		}
	}
	if foundItem == nil {
		return fmt.Errorf("resource '%s' not found in manifest", idOrHref)
	}

	filePath := e.resolveContainerPath(foundItem.Href)
	e.zipContainer.AddFile(filePath, content)
	return nil
}

// UpdateContent updates the byte content of an existing content document identified by ID or Href.
func (e *Editor) UpdateContent(idOrHref string, content []byte) error {
	return e.UpdateResource(idOrHref, content)
}

// RemoveResource removes a resource from the manifest, spine (if referenced), and container.
func (e *Editor) RemoveResource(idOrHref string) error {
	p := e.CurrentPackage()
	if p == nil {
		return errors.New("no active package selected")
	}

	var foundItem *pkg.Item
	var itemIndex = -1
	for i, item := range p.Manifest.Items {
		if item.ID == idOrHref || item.Href == idOrHref {
			foundItem = &p.Manifest.Items[i]
			itemIndex = i
			break
		}
	}
	if foundItem == nil {
		return fmt.Errorf("resource '%s' not found in manifest", idOrHref)
	}

	// Remove from manifest
	p.Manifest.Items = append(p.Manifest.Items[:itemIndex], p.Manifest.Items[itemIndex+1:]...)

	// Remove from spine if present
	p.Spine.ItemRefs = slices.DeleteFunc(p.Spine.ItemRefs, func(ref pkg.ItemRef) bool {
		return ref.IDRef == foundItem.ID
	})

	// Remove from zip container
	filePath := e.resolveContainerPath(foundItem.Href)
	e.zipContainer.RemoveFile(filePath)

	return nil
}

// AddSpineItem appends a manifest item to the spine reading order.
func (e *Editor) AddSpineItem(idOrHref string) error {
	p := e.CurrentPackage()
	if p == nil {
		return errors.New("no active package selected")
	}

	id := idOrHref
	for _, item := range p.Manifest.Items {
		if item.Href == idOrHref {
			id = item.ID
			break
		}
	}

	exists := false
	for _, item := range p.Manifest.Items {
		if item.ID == id {
			exists = true
			break
		}
	}
	if !exists {
		return fmt.Errorf("item with ID or Href '%s' not found in manifest", idOrHref)
	}

	p.Spine.ItemRefs = append(p.Spine.ItemRefs, pkg.ItemRef{IDRef: id})
	return nil
}

// RemoveSpineItem removes an item from the spine reading order without removing the resource.
func (e *Editor) RemoveSpineItem(idOrHref string) error {
	p := e.CurrentPackage()
	if p == nil {
		return errors.New("no active package selected")
	}

	id := idOrHref
	for _, item := range p.Manifest.Items {
		if item.Href == idOrHref {
			id = item.ID
			break
		}
	}

	beforeLen := len(p.Spine.ItemRefs)
	p.Spine.ItemRefs = slices.DeleteFunc(p.Spine.ItemRefs, func(ref pkg.ItemRef) bool {
		return ref.IDRef == id
	})
	if len(p.Spine.ItemRefs) == beforeLen {
		return fmt.Errorf("spine item '%s' not found", idOrHref)
	}

	return nil
}

// Resources returns all publication resources declared in the manifest.
func (e *Editor) Resources() []PublicationResource {
	p := e.CurrentPackage()
	if p == nil {
		return nil
	}

	allFiles := e.zipContainer.AllFiles()
	res := make([]PublicationResource, 0, len(p.Manifest.Items))
	for _, item := range p.Manifest.Items {
		filePath := e.resolveContainerPath(item.Href)
		res = append(res, PublicationResource{
			ID:         item.ID,
			Href:       item.Href,
			MIMEType:   item.MediaType,
			Content:    allFiles[filePath],
			Filepath:   filePath,
			Properties: item.Properties,
		})
	}
	return res
}

// SelectResourceById retrieves a publication resource by its manifest ID.
func (e *Editor) SelectResourceById(id string) *PublicationResource {
	for _, res := range e.Resources() {
		if res.ID == id {
			return &res
		}
	}
	return nil
}

// SelectResourceByHref retrieves a publication resource by its manifest Href.
func (e *Editor) SelectResourceByHref(href string) *PublicationResource {
	for _, res := range e.Resources() {
		if res.Href == href {
			return &res
		}
	}
	return nil
}

// Cover sets the publication cover image from raw bytes.
func (e *Editor) Cover(cover []byte) error {
	p := e.CurrentPackage()
	if p == nil {
		return errors.New("no active package selected")
	}

	name := "cover"
	mime := http.DetectContentType(cover)
	contentName := name
	switch mime {
	case "image/png":
		contentName += ".png"
	case "image/jpeg", "image/jpg":
		contentName += ".jpeg"
	default:
		contentName += ".jpeg"
	}

	filePath := e.resolveContainerPath(contentName)
	id := "cover-image"

	// Remove existing cover image manifest items
	p.Manifest.Items = slices.DeleteFunc(p.Manifest.Items, func(it pkg.Item) bool {
		return it.ID == id || it.Properties == pkg.CoverImageProperty
	})

	e.zipContainer.AddFile(filePath, cover)
	p.Manifest.Items = append(p.Manifest.Items, pkg.Item{
		ID:         id,
		Href:       contentName,
		MediaType:  mime,
		Properties: pkg.CoverImageProperty,
	})

	foundMeta := false
	for i, m := range p.Metadata.Meta {
		if m.Name == "cover" {
			p.Metadata.Meta[i].Content = id
			foundMeta = true
			break
		}
	}
	if !foundMeta {
		p.Metadata.Meta = append(p.Metadata.Meta, pkg.Meta{
			Name:    "cover",
			Content: id,
		})
	}

	return nil
}

// CoverPNG sets the cover image from an image.Image encoded as PNG.
func (e *Editor) CoverPNG(cover image.Image) error {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, cover); err != nil {
		return err
	}
	return e.Cover(buf.Bytes())
}

// CoverJPG sets the cover image from an image.Image encoded as JPEG.
func (e *Editor) CoverJPG(cover image.Image) error {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, cover, &jpeg.Options{Quality: 70}); err != nil {
		return err
	}
	return e.Cover(buf.Bytes())
}

// CoverFile sets the cover image by reading from a file path.
func (e *Editor) CoverFile(name string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	return e.Cover(data)
}

// TableOfContents generates and updates the EPUB navigation document (XHTML) and legacy NCX.
func (e *Editor) TableOfContents(name string, toc TOC) error {
	p := e.CurrentPackage()
	if p == nil {
		return errors.New("no active package selected")
	}

	// Remove existing NAV and NCX items from manifest and zip container
	p.Manifest.Items = slices.DeleteFunc(p.Manifest.Items, func(it pkg.Item) bool {
		if it.Properties == pkg.NavProperty || it.MediaType == pkg.MediaTypeNCX || it.ID == p.Spine.TOC {
			filePath := e.resolveContainerPath(it.Href)
			e.zipContainer.RemoveFile(filePath)
			return true
		}
		return false
	})

	navigation := ncx.NCX{
		NavMap: ncx.NavMap{
			ID:        "navmap",
			NavPoints: []ncx.NavPoint{},
		},
	}
	navigation.DocTitle.Text = toc.Title
	navigation.NavMap.NavPoints = buildNavPoints(toc.Items, 0)

	ncxContent, err := xml.MarshalIndent(navigation, "", " ")
	if err != nil {
		return err
	}
	ncxBase := name + ".ncx"
	ncxFilePath := e.resolveContainerPath(ncxBase)
	ncxID := e.uniqueID(name + "-ncx")

	e.zipContainer.AddFile(ncxFilePath, ncxContent)
	p.Manifest.Items = append(p.Manifest.Items, pkg.Item{
		ID:        ncxID,
		Href:      ncxBase,
		MediaType: pkg.MediaTypeNCX,
	})
	p.Spine.TOC = ncxID

	languages := []string{}
	for _, lang := range p.Metadata.Languages {
		languages = append(languages, lang.Value)
	}
	xhtmlNode, err := tocToHTMLNode(toc, languages)
	if err != nil {
		return err
	}

	var content bytes.Buffer
	if err := html.Render(&content, xhtmlNode); err != nil {
		return err
	}
	xhtmlBase := name + ".xhtml"
	xhtmlFilePath := e.resolveContainerPath(xhtmlBase)
	xhtmlID := e.uniqueID(name + "-nav")

	e.zipContainer.AddFile(xhtmlFilePath, content.Bytes())
	p.Manifest.Items = append(p.Manifest.Items, pkg.Item{
		ID:         xhtmlID,
		Href:       xhtmlBase,
		MediaType:  pkg.MediaTypeXHTML,
		Properties: pkg.NavProperty,
	})

	return nil
}

// Build serializes the modified publication package(s) and creates an in-memory EPUB zip archive.
func (e *Editor) Build() ([]byte, error) {
	e.ensureModifiedDate()

	for name, packageItem := range e.packagePubs {
		pkgPath, ok := e.packagePaths[name]
		if !ok {
			pkgPath = path.Join("epub", name+".opf")
			e.packagePaths[name] = pkgPath
		}
		if err := e.zipContainer.AddPackage(pkgPath, *packageItem); err != nil {
			return nil, fmt.Errorf("failed to marshal package %s: %w", name, err)
		}
	}

	var buf bytes.Buffer
	if err := e.zipContainer.WriteArchive(&buf); err != nil {
		return nil, fmt.Errorf("failed to write epub archive: %w", err)
	}

	return buf.Bytes(), nil
}

// Save writes the edited EPUB archive to an io.Writer.
func (e *Editor) Save(writer io.Writer) error {
	if writer == nil {
		return errors.New("epub: writer cannot be nil")
	}
	data, err := e.Build()
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

// SaveAs writes the edited EPUB archive to the specified file on disk.
func (e *Editor) SaveAs(filename string) error {
	data, err := e.Build()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0o644)
}

// Write writes the edited EPUB archive to the specified file on disk (alias for SaveAs).
func (e *Editor) Write(filename string) error {
	return e.SaveAs(filename)
}

// WriteBytes writes the edited EPUB bytes into the byte slice pointed to by b.
func (e *Editor) WriteBytes(b *[]byte) error {
	if b == nil {
		return errors.New("epub: byte pointer cannot be nil")
	}
	data, err := e.Build()
	if err != nil {
		return err
	}
	*b = data
	return nil
}

// SaveByte is an alias for WriteBytes for backwards compatibility.
func (e *Editor) SaveByte(b *[]byte) error {
	return e.WriteBytes(b)
}

// SaveBytes returns the edited EPUB archive as a byte slice.
func (e *Editor) SaveBytes() ([]byte, error) {
	return e.Build()
}

// Reader finalizes all pending edits and converts the publication back into a Reader instance.
func (e *Editor) Reader() (Reader, error) {
	data, err := e.Build()
	if err != nil {
		return Reader{}, err
	}
	return NewReader(data)
}


