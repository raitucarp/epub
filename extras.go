package epub

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/raitucarp/epub/pkg"
	"golang.org/x/net/html"
	"golang.org/x/text/runes"
)

// UID returns the unique identifier of the publication.
func (r *Reader) UID() (identifier string) {
	for _, uid := range r.CurrentSelectedPackage().Metadata.Identifiers {
		identifier = uid.Value
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
				resId := value.(string)
				cover = r.ReadImageById(resId)
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
	if n.Type == html.ElementNode && n.Data == "img" {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if img := findFirstImg(c); img != nil {
			return img
		}
	}

	return nil
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
func (r *Reader) Title() (title string) {
	for key, value := range r.Metadata() {
		if key == "title" {
			title = strings.Join(value.([]string), ", ")
			return
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
				title = getTextByEpubType(htmlNode, "title")
				return
			}
		}
	}

	if title != "" {
		return
	}

	for _, ref := range r.epub.resources {
		if titlePattern.MatchString(ref.ID) || titlePattern.MatchString(ref.Href) {
			htmlNode, _ := r.parseHTML(ref.Content)
			title = getTextByEpubType(htmlNode, "title")
			if title == "" {
				getTextByEpubType(htmlNode, "fulltitle")
			}

			if title != "" {
				return
			}
		}
	}

	return
}

// Author returns the author (creator) metadata of the publication.
func (r *Reader) Author() (author string) {
	for key, value := range r.Metadata() {
		if key == "creator" {
			author = strings.Join(value.([]string), ", ")
			return
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
				author = getTextByEpubType(htmlNode, "author")
				return
			}
		}
	}

	for _, ref := range r.epub.resources {
		if titlePattern.MatchString(ref.ID) || titlePattern.MatchString(ref.Href) {
			htmlNode, _ := r.parseHTML(ref.Content)
			author = getTextByEpubType(htmlNode, "author")

			if author != "" {
				return
			}
		}
	}

	return "Unknown"
}

// Language returns the primary language of the publication, as declared
// in the package metadata (dc:language).
func (r *Reader) Language() (language string) {
	desc, descriptionExists := r.epub.metadata["language"]
	if descriptionExists {
		language = strings.Join(desc.([]string), ", ")
		return
	}
	return
}

// Identifier returns the primary identifier of the publication as declared
// in the package metadata (often equivalent to UID).
func (r *Reader) Identifier() (identifier string) {
	desc, descriptionExists := r.epub.metadata["identifier"]
	if descriptionExists {
		identifier = strings.Join(desc.([]string), ", ")
		return
	}
	return
}

var descriptionPattern = regexp.MustCompile("description")

func extractDescriptionFromMetadata(metadata map[string]any) (description string) {
	desc, descriptionExists := metadata["description"]
	if descriptionExists {
		description = strings.Join(desc.([]string), ", ")
	}
	return
}

func extractDescriptionFromOptionalMeta(metadata map[string]any) (description string) {
	meta, metaExists := metadata["meta"]
	if !metaExists {
		return
	}

	metaMap, isMetaMap := meta.(map[string]any)
	if !isMetaMap {
		return
	}

	for key, value := range metaMap {
		if !descriptionPattern.MatchString(key) {
			continue
		}

		desc, isSlice := value.([]any)
		if !isSlice {
			continue
		}

		descriptions := []string{}
		for _, d := range desc {
			if v, ok := d.(string); ok {
				descriptions = append(descriptions, v)
			}
		}

		description = strings.Join(descriptions, ", ")
	}

	return
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

func (r *Reader) extractDescriptionFromSpine() (description string) {
	spine := r.Spine()
	introTypes := []string{"abstract", "foreword", "introduction", "preamble", "preface", "prologue"}
	for _, res := range spine {
		htmlNode := r.ReadContentHTMLByHref(res.Href)

		for _, introType := range introTypes {
			if description == "" {
				description = extractDescriptionFromEpubType(introType, htmlNode)
			}
		}
	}

	return
}

func getBody(doc *html.Node) *html.Node {
	var body *html.Node
	var findBody func(*html.Node)

	findBody = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findBody(c)
		}
	}

	findBody(doc)
	return body
}

func (r *Reader) extractDescriptionFromReferences() (description string) {
	refs := r.References()

	candidates := []pkg.GuideReferenceType{pkg.GuideRefText, pkg.GuideRefPreface, pkg.GuideRefForeword}

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

	return
}

func (r *Reader) extractDescriptionFromFirstFullContentTOCItem() (description string) {
	toc, err := r.TableOfContents()
	if err != nil {
		return
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
		return
	}

	content := r.ReadContentHTMLByHref(firstContentItem.Href)

	body := getBody(content)
	markdownBody, _ := htmltomarkdown.ConvertNode(body)
	description = string(markdownBody)

	return
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
func (r *Reader) Description() (description string) {
	metadata := r.epub.metadata
	description = extractDescriptionFromMetadata(metadata)

	if description == "" {
		description = extractDescriptionFromOptionalMeta(metadata)
	}

	if description == "" {
		description = r.extractDescriptionFromSpine()
	}

	if description == "" {
		description = r.extractDescriptionFromReferences()
	}

	if description == "" {
		description = r.extractDescriptionFromFirstFullContentTOCItem()
	}

	if description != "" {
		description = convertDescriptionToMd(description)
	}
	return
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
