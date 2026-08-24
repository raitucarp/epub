package epub

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/url"
	"regexp"
	"slices"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/raitucarp/epub/pkg"
	"golang.org/x/net/html"
	"golang.org/x/text/runes"
)

// UID returns the unique identifier of the publication. It resolves the
// package's unique-identifier attribute to the matching dc:identifier entry.
func (r *Reader) UID() (identifier string) {
	pkgMetadata := r.CurrentSelectedPackage().Metadata
	identifiers := pkgMetadata.Identifiers

	for _, uid := range identifiers {
		if uid.ID == r.CurrentSelectedPackage().UniqueIdentifier {
			return uid.Value
		}
	}

	for _, uid := range identifiers {
		if uid.Value != "" {
			identifier = uid.Value
		}
	}

	return
}

// Version returns the EPUB specification version of the publication.
func (r *Reader) Version() (version string) {
	return r.CurrentSelectedPackage().Version
}

var coverImagePattern = regexp.MustCompile("cover")

func (r *Reader) getCoverInMetadata() (cover *image.Image) {
	metadata := r.Metadata()
	if meta, ok := metadata["meta"]; ok {
		metaMap, ok := meta.(map[string]any)

		if !ok {
			return
		}

		for key, value := range metaMap {
			if coverImagePattern.MatchString(key) {
				if resId, ok := value.(string); ok {
					cover = r.ReadImageById(resId)
				}
			}
		}
		return
	}

	return
}

func (r *Reader) getCoverInResources() (cover *image.Image) {
	resources := r.Resources()
	for _, res := range resources {
		if res.Properties == pkg.CoverImageProperty {
			cover = r.ReadImageByHref(res.Href)
		}

		if cover == nil && coverImagePattern.MatchString(res.ID) {
			cover = r.ReadImageByHref(res.Href)
		}
	}

	return
}

func (r *Reader) getCoverInSpine() (cover *image.Image) {
	for _, item := range r.Spine() {
		if !coverImagePattern.MatchString(item.ID) {
			continue
		}
		res := r.SelectResourceById(item.ID)
		if res != nil {
			cover = r.ReadImageByHref(res.Href)
		}
	}

	return
}

func findFirstImg(n *html.Node) *html.Node {
	return FindNode(n, func(node *html.Node) bool {
		return node.Type == html.ElementNode && node.Data == "img"
	})
}

func getImageSrc(imgNode *html.Node) (href string) {
	srcIndex := slices.IndexFunc(imgNode.Attr, func(attr html.Attribute) bool { return attr.Key == "src" })
	if srcIndex == -1 {
		return
	}
	href = imgNode.Attr[srcIndex].Val
	return
}

func (r *Reader) getCoverFromToc() (cover *image.Image) {
	toc, err := r.TableOfContents()
	if err != nil {
		return
	}

	for _, item := range toc.Items {
		if !coverImagePattern.MatchString(item.Href) {
			continue
		}

		htmlNode := r.ReadContentHTMLByHref(item.Href)
		if htmlNode == nil {
			continue
		}
		firstImg := findFirstImg(htmlNode)
		href := getImageSrc(firstImg)
		cover = r.ReadImageByHref(href)

	}

	return
}

// Cover returns the publication's cover image if present.
func (r *Reader) Cover() (cover *image.Image) {
	cover = r.getCoverInMetadata()

	if cover == nil {
		cover = r.getCoverInResources()
	}

	if cover == nil {
		cover = r.getCoverInSpine()
	}

	if cover == nil {
		cover = r.getCoverFromToc()
	}

	return
}

// CoverBytes returns the raw byte representation of the cover image.
// An error is returned if the publication does not define a cover.
func (r *Reader) CoverBytes() (cover []byte, err error) {
	coverImage := r.Cover()

	buf := new(bytes.Buffer)

	err = png.Encode(buf, *coverImage)
	if err == nil {
		return buf.Bytes(), err
	}

	err = jpeg.Encode(buf, *coverImage, &jpeg.Options{Quality: 70})
	if err == nil {
		return buf.Bytes(), err
	}

	return nil, err
}

var titlePattern = regexp.MustCompile("title")

// Title returns the publication's title metadata.
func (r *Reader) Title() []string {
	var titles []string
	for key, value := range r.Metadata() {
		if key == "title" {
			titles = append(titles, value.([]string)...)
			return titles
		}
	}

	guide := r.CurrentSelectedPackage().Guide

	if guide != nil {
		for _, ref := range guide.References {
			if ref.Type == pkg.GuideRefTitlePage {
				res := r.SelectResourceByHref(ref.Href)
				if res == nil {
					continue
				}

				htmlNode, err := r.parseHTML(res.Content)
				if err != nil {
					continue
				}
				t := getTextByEpubType(htmlNode, "title")
				if t != "" {
					return []string{t}
				}
			}
		}
	}

	for _, ref := range r.epub.resources {
		if titlePattern.MatchString(ref.ID) || titlePattern.MatchString(ref.Href) {
			htmlNode, _ := r.parseHTML(ref.Content)
			t := getTextByEpubType(htmlNode, "title")
			if t == "" {
				t = getTextByEpubType(htmlNode, "fulltitle")
			}

			if t != "" {
				return []string{t}
			}
		}
	}

	return titles
}

// Author returns the author (creator) metadata of the publication.
func (r *Reader) Author() []string {
	var authors []string
	for key, value := range r.Metadata() {
		if key == "creator" {
			authors = append(authors, value.([]string)...)
			return authors
		}
	}

	guide := r.CurrentSelectedPackage().Guide

	if guide != nil {
		for _, ref := range guide.References {
			if ref.Type == pkg.GuideRefTitlePage {
				res := r.SelectResourceByHref(ref.Href)
				if res == nil {
					continue
				}

				htmlNode, err := r.parseHTML(res.Content)
				if err != nil {
					continue
				}
				a := getTextByEpubType(htmlNode, "author")
				if a != "" {
					return []string{a}
				}
			}
		}
	}

	for _, ref := range r.epub.resources {
		if titlePattern.MatchString(ref.ID) || titlePattern.MatchString(ref.Href) {
			htmlNode, _ := r.parseHTML(ref.Content)
			a := getTextByEpubType(htmlNode, "author")

			if a != "" {
				return []string{a}
			}
		}
	}

	if len(authors) == 0 {
		return []string{"Unknown"}
	}
	return authors
}

// Language returns the primary language of the publication, as declared
// in the package metadata (dc:language).
func (r *Reader) Language() []string {
	desc, descriptionExists := r.epub.metadata["language"]
	if descriptionExists {
		return desc.([]string)
	}
	return nil
}

// Identifier returns the primary identifier of the publication as declared
// in the package metadata (often equivalent to UID).
func (r *Reader) Identifier() []string {
	desc, descriptionExists := r.epub.metadata["identifiers"]
	if descriptionExists {
		return desc.([]string)
	}
	return nil
}

var descriptionPattern = regexp.MustCompile("description")

func extractDescriptionFromMetadata(metadata map[string]any) []string {
	desc, descriptionExists := metadata["description"]
	if descriptionExists {
		return desc.([]string)
	}
	return nil
}

func extractDescriptionFromOptionalMeta(metadata map[string]any) []string {
	meta, metaExists := metadata["meta"]
	if !metaExists {
		return nil
	}

	metaMap, isMetaMap := meta.(map[string]any)
	if !isMetaMap {
		return nil
	}

	var allDescriptions []string
	for key, value := range metaMap {
		if !descriptionPattern.MatchString(key) {
			continue
		}

		desc, isSlice := value.([]any)
		if !isSlice {
			continue
		}

		for _, d := range desc {
			if v, ok := d.(string); ok {
				allDescriptions = append(allDescriptions, v)
			}
		}
	}

	return allDescriptions
}

func extractDescriptionFromSummaryMeta(metadata map[string]any) []string {
	meta, metaExists := metadata["meta"]
	if !metaExists {
		return nil
	}

	metaMap, isMetaMap := meta.(map[string]any)
	if !isMetaMap {
		return nil
	}

	for key, value := range metaMap {
		if key != "summary" {
			continue
		}

		desc, isSlice := value.([]any)
		if !isSlice {
			continue
		}

		for _, d := range desc {
			if v, ok := d.(string); ok {
				return []string{v}
			}
		}
	}
	return nil
}

func extractDescriptionFromEpubType(epubType string, htmlNode *html.Node) (description string) {
	for desc := range htmlNode.Descendants() {
		abstractIndex := slices.IndexFunc(desc.Attr, func(attr html.Attribute) bool {
			return attr.Key == "epub:type" && attr.Val == epubType
		})

		if abstractIndex > -1 {
			descByte, err := htmltomarkdown.ConvertNode(desc)
			if err != nil {
				continue
			}
			description = string(descByte)
		}
	}
	return
}

var coverPagePattern = regexp.MustCompile("cover")

func (r *Reader) extractDescriptionFromSpine() []string {
	spine := r.Spine()
	introTypes := []string{"abstract", "foreword", "introduction", "preamble", "preface", "prologue"}
	for _, res := range spine {
		htmlNode := r.ReadContentHTMLByHref(res.Href)
		if htmlNode == nil {
			continue
		}

		for desc := range htmlNode.Descendants() {
			abstractIndex := slices.IndexFunc(desc.Attr, func(attr html.Attribute) bool {
				return attr.Key == "epub:type" && slices.Contains(introTypes, attr.Val)
			})

			if abstractIndex > -1 {
				descByte, err := htmltomarkdown.ConvertNode(desc)
				if err == nil {
					return []string{string(descByte)}
				}
			}
		}
	}

	return nil
}

func getBody(doc *html.Node) *html.Node {
	return FindNode(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "body"
	})
}

func (r *Reader) extractDescriptionFromReferences() []string {
	refs := r.References()

	candidates := []pkg.GuideReferenceType{pkg.GuideRefText, pkg.GuideRefPreface, pkg.GuideRefForeword}

	var description string
	for _, candidate := range candidates {
		if description != "" {
			continue
		}

		if doc, ok := refs[candidate]; ok {
			body := getBody(doc)

			markdownBody, _ := htmltomarkdown.ConvertNode(body)
			description = string(markdownBody)
		}
	}

	if description != "" {
		return []string{description}
	}
	return nil
}

func (r *Reader) extractDescriptionFromFirstFullContentTOCItem() []string {
	toc, err := r.TableOfContents()
	if err != nil {
		return nil
	}

	var firstContentItem TOC
	for _, item := range toc.Items {
		if coverPagePattern.MatchString(item.Href) {
			continue
		}

		if firstContentItem.Title == "" {
			firstContentItem = item
			break
		}
	}

	if firstContentItem.Title == "" {
		return nil
	}

	content := r.ReadContentHTMLByHref(firstContentItem.Href)

	body := getBody(content)
	markdownBody, _ := htmltomarkdown.ConvertNode(body)

	if len(markdownBody) > 0 {
		return []string{string(markdownBody)}
	}
	return nil
}

func convertDescriptionToMd(description string) string {
	descriptionTemp := description
	newDesc, err := htmltomarkdown.ConvertString(description)
	if err != nil {
		description = descriptionTemp
	}
	description = newDesc

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	description, _, _ = transform.String(t, description)

	return description
}

// Description returns the publication's description metadata if defined.
func (r *Reader) Description() []string {
	metadata := r.epub.metadata
	descriptions := extractDescriptionFromMetadata(metadata)

	if len(descriptions) == 0 {
		descriptions = extractDescriptionFromOptionalMeta(metadata)
	}

	if len(descriptions) == 0 {
		descriptions = extractDescriptionFromSummaryMeta(metadata)
		if len(descriptions) > 0 {
			descByte, err := htmltomarkdown.ConvertString(descriptions[0])
			if err == nil {
				descriptions[0] = string(descByte)
			}
		}
	}

	if len(descriptions) == 0 {
		descriptions = r.extractDescriptionFromSpine()
	}

	if len(descriptions) == 0 {
		descriptions = r.extractDescriptionFromReferences()
	}

	if len(descriptions) == 0 {
		descriptions = r.extractDescriptionFromFirstFullContentTOCItem()
	}

	var finalDesc []string
	for _, desc := range descriptions {
		if desc != "" {
			finalDesc = append(finalDesc, convertDescriptionToMd(desc))
		}
	}

	return finalDesc
}

// References returns the structural guide references defined in the package,
// such as "cover", "title-page", "toc", etc. The returned map is keyed by
// reference type and mapped to corresponding HTML content.
func (r *Reader) References() (references map[pkg.GuideReferenceType]*html.Node) {
	references = make(map[pkg.GuideReferenceType]*html.Node)

	guides := r.epub.SelectedPackage().Guide
	if guides == nil {
		return
	}

	for _, ref := range guides.References {

		u, err := url.Parse(ref.Href)
		if err != nil {
			continue
		}

		// Remove the fragment by setting it to empty string
		u.Fragment = ""

		// Get the string representation without fragment
		cleanHref := u.String()
		content := r.ReadContentHTMLByHref(cleanHref)
		references[ref.Type] = content
	}

	return
}
