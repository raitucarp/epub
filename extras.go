package epub

import (
	"image"
	"regexp"
	"strings"

	"github.com/raitucarp/epub/pkg"
)

func (r *Reader) UID() (identifier string) {
	for _, uid := range r.CurrentSelectedPackage().Metadata.Identifiers {
		identifier = uid.Value
	}
	return
}

func (r *Reader) Version() (version string) {
	return r.CurrentSelectedPackage().Version
}

var coverImagePattern = regexp.MustCompile("cover")

func (r *Reader) Cover() (cover *image.Image) {
	metadata := r.Metadata()
	resources := r.Resources()

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

	for _, res := range resources {
		if coverImagePattern.MatchString(res.Properties) || coverImagePattern.MatchString(res.ID) {
			cover = r.ReadImageById(res.ID)
		}
	}

	if cover != nil {
		return
	}

	for _, item := range r.Spine() {
		if coverImagePattern.MatchString(item.Properties) || coverImagePattern.MatchString(item.ID) {
			cover = r.ReadImageById(item.ID)
		}
	}

	if cover != nil {
		return
	}

	for _, guide := range r.CurrentSelectedPackage().Guide.References {
		if guide.Type != pkg.GuideRefCover {
			continue
		}

		cover = r.ReadImageByHref(guide.Href)
	}

	return
}

var titlePattern = regexp.MustCompile("title")

func (r *Reader) Title() (title string) {
	for key, value := range r.Metadata() {
		if key == "title" {
			title = strings.Join(value.([]string), ", ")
			return
		}
	}

	for _, ref := range r.CurrentSelectedPackage().Guide.References {
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

func (r *Reader) Author() (author string) {
	for key, value := range r.Metadata() {
		if key == "creator" {
			author = strings.Join(value.([]string), ", ")
			return
		}
	}

	for _, ref := range r.CurrentSelectedPackage().Guide.References {
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

	for _, ref := range r.epub.resources {
		if titlePattern.MatchString(ref.ID) || titlePattern.MatchString(ref.Href) {
			htmlNode, _ := r.parseHTML(ref.Content)
			author = getTextByEpubType(htmlNode, "author")

			if author != "" {
				return
			}
		}
	}

	return "Anonymous"
}

var descriptionPattern = regexp.MustCompile("description")

func (r *Reader) Description() (description string) {
	desc, descriptionExists := r.epub.metadata["description"]
	if descriptionExists {
		description = strings.Join(desc.([]string), ", ")
		return
	}

	meta, metaExists := r.epub.metadata["meta"]
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
		return
	}

	return
}

func (r *Reader) TableOfContents() (version string) {
	return r.CurrentSelectedPackage().Version
}

// cover	the book cover(s), jacket information, etc.
// title-page	page with possibly title, author, publisher, and other metadata
// toc	table of contents
// index	back-of-book style index
// glossary
// acknowledgements
// bibliography
// colophon
// copyright-page
// dedication
// epigraph
// foreword
// loi	list of illustrations
// lot	list of tables
// notes
// preface
// text	First "real" page of content (e.g. "Chapter 1")
