package epub

import (
	"bytes"
	"fmt"
	"image"
	"regexp"
	"slices"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/pkg"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// SelectPackageRendition changes the active package rendition by its
// rendition identifier. Useful when multiple reading layouts are available.
func (r *Reader) SelectPackageRendition(rendition string) {
	r.epub.rendition = rendition
	r.parseResources()
}

// CurrentSelectedPackage returns the currently active package rendition.
// EPUB publications may have multiple renditions.
func (r *Reader) CurrentSelectedPackage() *pkg.Package {
	return r.epub.packagePubs[r.epub.rendition]
}

// CurrentSelectedPackagePath returns the internal path to the currently
// selected package document.
func (r *Reader) CurrentSelectedPackagePath() string {
	return r.epub.packagePaths[r.epub.rendition]
}

func (r *Reader) parseHTML(htmlByte []byte) (node *html.Node, err error) {
	content := bytes.NewReader(htmlByte)
	node, err = html.Parse(content)

	return
}

// ListContentDocumentIds returns the IDs of all content documents
// (XHTML/SVG) registered in the publication manifest.
func (r *Reader) ListContentDocumentIds() (ids []string) {
	for _, res := range r.Resources() {
		if res.MIMEType == pkg.MediaTypeXHTML {
			ids = append(ids, res.ID)
		}
	}

	return
}

// ListImageIds returns the IDs of all image-based resources
// (e.g., PNG, JPEG, SVG) in the publication manifest.
func (r *Reader) ListImageIds() (ids []string) {
	mimes := []string{pkg.MediaTypeJPEG, pkg.MediaTypePNG, pkg.MediaTypeSVG, pkg.MediaTypeGIF, pkg.MediaTypeWebP}
	for _, res := range r.Resources() {
		if slices.Contains(mimes, res.MIMEType) {
			ids = append(ids, res.ID)
		}
	}

	return
}

// ContentDocumentXHTML returns XHTML content documents parsed into html.Node trees.
// The returned map is keyed by EPUB manifest item ID.
func (r *Reader) ContentDocumentXHTML() (documents map[string]*html.Node) {
	documents = make(map[string]*html.Node)

	for _, res := range r.epub.resources {
		if res.MIMEType == pkg.MediaTypeXHTML {
			node, err := r.parseHTML(res.Content)

			if err != nil {
				continue
			}
			documents[res.ID] = node
		}
	}
	return
}

// ContentDocumentXHTMLString returns XHTML content documents as raw strings.
// The returned map is keyed by EPUB manifest item ID.
func (r *Reader) ContentDocumentXHTMLString() (documents map[string]string) {
	resourcesHtml := r.ContentDocumentXHTML()
	documents = make(map[string]string)

	for resId, res := range resourcesHtml {
		var buffer bytes.Buffer
		html.Render(&buffer, res)
		documents[resId] = buffer.String()
	}
	return
}

func cleanupHTML(node *html.Node) (newNode *html.Node) {
	for desc := range node.Descendants() {
		if desc.DataAtom == atom.Title {
			desc.Parent.RemoveChild(desc)
			continue
		}
		if desc.Type == html.ElementNode && desc.Data == "a" {
			hrefAttributeIndex := slices.IndexFunc(desc.Attr,
				func(attr html.Attribute) bool {
					return attr.Key == "href"
				})
			if hrefAttributeIndex < 0 {
				desc.Data = "div"
			}
		}
	}
	return node
}

func extractTitle(node *html.Node) string {
	var title string

	for desc := range node.Descendants() {
		if desc.DataAtom == atom.Title {
			title = desc.FirstChild.Data
		}
	}
	return title
}

func getTextByEpubType(node *html.Node, attributeValue string) (text string) {
	for desc := range node.Descendants() {
		attrIndex := slices.IndexFunc(desc.Attr, func(attr html.Attribute) bool {
			matched, _ := regexp.MatchString(attributeValue, attr.Val)
			return attr.Key == "epub:type" && matched
		})

		if attrIndex > -1 {
			text = desc.FirstChild.Data
		}
	}
	return
}

// ContentDocumentMarkdown returns content documents converted into Markdown
// form. The returned map is keyed by EPUB manifest item ID.
func (r *Reader) ContentDocumentMarkdown() (documents map[string]string) {
	resourcesHtml := r.ContentDocumentXHTML()
	documents = make(map[string]string)

	for resId, res := range resourcesHtml {
		frontMatters := ""
		title := extractTitle(res)
		cleanedHTML := cleanupHTML(res)
		if title != "" {
			frontMatters = fmt.Sprintf(`---
title: %#v
---`, title)
		}
		md, err := htmltomarkdown.ConvertNode(cleanedHTML)
		if err != nil {
			continue
		}

		markdownString := string(md)
		if frontMatters != "" {
			markdownString = frontMatters + "\n" + string(markdownString)
		}
		documents[resId] = markdownString
	}
	return
}

// ReadContentHTMLById returns the XHTML/HTML content document associated
// with the given manifest ID, parsed into an html.Node tree.
func (r *Reader) ReadContentHTMLById(id string) (doc *html.Node) {
	resourcesHtml := r.ContentDocumentXHTML()
	for resId, res := range resourcesHtml {
		if resId == id {
			return res
		}
	}
	return
}

// ReadContentHTMLByHref returns the content document associated with the given
// manifest href. The returned document is parsed into an html.Node tree.
func (r *Reader) ReadContentHTMLByHref(href string) (doc *html.Node) {
	contentIndex := slices.IndexFunc(r.epub.resources, func(r PublicationResource) bool {
		return r.Href == href
	})

	if contentIndex > -1 {
		res := r.epub.resources[contentIndex]
		return r.ReadContentHTMLById(res.ID)
	}

	return
}

// ReadContentMarkdownById returns a Markdown string representation of the
// content document associated with the given manifest ID.
func (r *Reader) ReadContentMarkdownById(id string) (md string) {
	resourcesMd := r.ContentDocumentMarkdown()
	for resId, res := range resourcesMd {
		if resId == id {
			return res
		}
	}
	return
}

// ReadImageById returns the image resource associated with the given
// manifest ID.
func (r *Reader) ReadImageById(id string) (img *image.Image) {

	for _, res := range r.epub.resources {
		if res.ID == id {

			reader := bytes.NewReader(res.Content)
			img, _, _ := image.Decode(reader)
			return &img
		}
	}
	return
}

// ReadImageByHref returns the image resource referenced by the given href,
// if present in the manifest.
func (r *Reader) ReadImageByHref(href string) (img *image.Image) {

	for _, res := range r.epub.resources {
		if res.Href == href {

			reader := bytes.NewReader(res.Content)
			img, _, _ := image.Decode(reader)
			return &img
		}
	}
	return
}

// ContentDocumentSVG returns SVG content documents parsed into html.Node trees.
// The returned map is keyed by EPUB manifest item ID.
func (r *Reader) ContentDocumentSVG() (documents map[string]*html.Node) {
	documents = make(map[string]*html.Node)

	for _, res := range r.epub.resources {
		if res.MIMEType == pkg.MediaTypeSVG {
			content := bytes.NewReader(res.Content)
			node, err := html.Parse(content)
			if err != nil {
				continue
			}
			documents[res.ID] = node
		}
	}
	return
}

// Images returns all image resources in the publication, keyed by manifest ID.
func (r *Reader) Images() (images map[string]image.Image) {
	images = make(map[string]image.Image)

	for _, res := range r.epub.resources {
		if res.MIMEType == pkg.MediaTypeJPEG ||
			res.MIMEType == pkg.MediaTypePNG ||
			res.MIMEType == pkg.MediaTypeGIF ||
			res.MIMEType == pkg.MediaTypeWebP {
			reader := bytes.NewReader(res.Content)

			img, _, err := image.Decode(reader)
			if err != nil {
				continue
			}

			images[res.ID] = img
		}
	}
	return
}

// ImageBytes returns all image resources in the publication, keyed by manifest ID.
func (r *Reader) ImageResources() (images map[string][]byte) {
	images = make(map[string][]byte)

	for _, res := range r.epub.resources {
		if res.MIMEType == pkg.MediaTypeJPEG ||
			res.MIMEType == pkg.MediaTypePNG ||
			res.MIMEType == pkg.MediaTypeGIF ||
			res.MIMEType == pkg.MediaTypeWebP {

			images[res.ID] = res.Content
		}
	}
	return
}

// Spine returns publication's spines, ordered resources like table of contents.
func (r *Reader) Spine() (orderedResources []PublicationResource) {
	for _, item := range r.CurrentSelectedPackage().Spine.ItemRefs {
		for _, res := range r.epub.resources {
			if item.IDRef == res.ID {
				orderedResources = append(orderedResources, res)
			}
		}
	}

	return
}

func (r *Reader) parseMetadata() {
	r.epub.metadata = make(map[string]any)
	packageMetadata := r.CurrentSelectedPackage().Metadata

	identifiers := []string{}
	for _, dcIdentifiers := range packageMetadata.Identifiers {
		identifiers = append(identifiers, dcIdentifiers.Value)
	}
	r.epub.metadata["identifiers"] = identifiers

	titles := []string{}
	for _, title := range packageMetadata.Titles {
		titles = append(titles, title.Value)
	}
	r.epub.metadata["title"] = titles

	languages := []string{}
	for _, language := range packageMetadata.Languages {
		r.epub.metadata["language"] = append(languages, language.Value)
	}
	r.epub.metadata["language"] = languages

	for _, optional := range packageMetadata.OptionalDC {
		name := optional.XMLName.Local
		if r.epub.metadata[name] == nil {
			r.epub.metadata[name] = []string{}
		}

		r.epub.metadata[name] = append(r.epub.metadata[name].([]string), optional.Value)
	}

	r.epub.metadata["meta"] = map[string]any{}
	for _, meta := range packageMetadata.Meta {
		name := meta.Property

		if name == "" {
			name = meta.Name
			r.epub.metadata["meta"].(map[string]any)[name] = meta.Content
			continue
		}

		if r.epub.metadata["meta"].(map[string]any)[name] == nil {
			r.epub.metadata["meta"].(map[string]any)[name] = []any{}
		}

		r.epub.metadata["meta"].(map[string]any)[name] = append(
			r.epub.metadata["meta"].(map[string]any)[name].([]any), meta.Value)

		r.epub.metadata["meta"].(map[string]any)[name] = slices.Compact(
			r.epub.metadata["meta"].(map[string]any)[name].([]any),
		)
	}

}

// Metadata returns the complete metadata block of the publication.
// The returned map may include standard as well as extended metadata fields.
func (r *Reader) Metadata() (metadata map[string]any) {
	return r.epub.metadata
}

// Refines returns metadata refinement relationships. The returned map is
// keyed by subject identifier, mapping to properties and their assigned values.
func (r *Reader) Refines() (refines map[string]map[string][]string) {
	refines = make(map[string]map[string][]string)
	packageMetadata := r.CurrentSelectedPackage().Metadata

	for _, dcIdentifiers := range packageMetadata.Identifiers {
		if dcIdentifiers.ID != "" {
			refineName := dcIdentifiers.ID
			if refines[refineName] == nil {
				refines[refineName] = make(map[string][]string)
			}

			refines[refineName][dcIdentifiers.XMLName.Local] = append(refines[refineName][dcIdentifiers.XMLName.Local], dcIdentifiers.Value)
		}
	}

	for _, title := range packageMetadata.Titles {
		if title.ID != "" {
			refineName := title.ID
			if refines[refineName] == nil {
				refines[refineName] = make(map[string][]string)
			}

			refines[refineName][title.XMLName.Local] = append(refines[refineName][title.XMLName.Local], title.Value)
		}
	}

	for _, language := range packageMetadata.Languages {
		if language.ID != "" {
			refineName := language.ID
			if refines[refineName] == nil {
				refines[refineName] = make(map[string][]string)
			}

			refines[refineName][language.XMLName.Local] = append(refines[refineName][language.XMLName.Local], language.Value)
		}
	}

	for _, optional := range packageMetadata.OptionalDC {
		if optional.ID != "" {
			refineName := optional.ID
			if refines[refineName] == nil {
				refines[refineName] = make(map[string][]string)
			}

			refines[refineName][optional.XMLName.Local] = append(refines[refineName][optional.XMLName.Local], optional.Value)
		}
	}

	refineCounter := map[string]int{}

	for _, meta := range packageMetadata.Meta {
		if meta.Refines != "" {
			refineName := strings.Trim(meta.Refines, "#")
			refineCounter[refineName]++
		}
	}

	for _, meta := range packageMetadata.Meta {
		if meta.ID != "" {
			refineName := meta.ID
			if refines[refineName] == nil {
				refines[refineName] = make(map[string][]string)
			}

			if refineCounter[refineName] > 0 {
				refines[refineName][meta.Property] = append(refines[refineName][meta.Property], meta.Value)
			} else {
				delete(refines, refineName)
			}
		}
	}

	for _, meta := range packageMetadata.Meta {
		if meta.Refines == "" {
			continue
		}
		refineName := strings.Trim(meta.Refines, "#")

		if refineCounter[refineName] > 0 {
			refines[refineName][meta.Property] = append(refines[refineName][meta.Property], meta.Value)
		}
	}

	// packageMetadata
	return
}

// NavigationCenterExtended returns the NCX navigation document (if available).
// This is primarily used for EPUB 2.x backward compatibility.
func (r *Reader) NavigationCenterExtended() *ncx.NCX {
	return r.epub.navigationCenterEXtended
}
