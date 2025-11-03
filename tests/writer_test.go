package tests

import (
	"testing"

	"github.com/raitucarp/epub"
)

func TestCreateEpubWithoutRequiredFieldShouldFail(t *testing.T) {
	epubData, err := epub.OpenReader("./data/arthur-conan-doyle_the-white-company.epub")
	if err != nil {
		t.Errorf("Something error %s", err)
		return
	}

	epubWriter := epub.New("https://standardebooks.org/ebooks/arthur-conan-doyle/the-white-company")
	epubWriter.Author(epubData.Author())

	err = epubWriter.Write("./temp/new-book.epub")
	if err == nil {
		t.Errorf("Write should fail because missing required")
	}
}

func TestCreateEpubSuccess(t *testing.T) {
	epubData, err := epub.OpenReader("./data/arthur-conan-doyle_the-white-company.epub")
	if err != nil {
		t.Errorf("Something error %s", err)
		return
	}

	identifier := epubData.Identifier()
	epubWriter := epub.New(identifier)
	epubWriter.Title(epubData.Title())
	epubWriter.Languages(epubData.Language())
	epubWriter.Subject("subject1", "Historical Fiction")

	// contents := epubData.ContentDocumentXHTMLString()
	spines := epubData.Spine()
	toc, err := epubData.TableOfContents()
	if err != nil {
		t.Errorf("Sommething error %s", err)
	}

	coverImage, err := epubData.CoverBytes()
	if err != nil {
		t.Errorf("%s", err)
	}
	epubWriter.Cover(coverImage)

	images := epubData.ImageResources()
	for href, imageData := range images {
		epubWriter.AddImage(href, imageData)
	}

	for _, itemRef := range spines {
		epubWriter.AddContent(itemRef.Href, itemRef.Content)
	}

	err = epubWriter.TableOfContents("table_of_contents", toc)
	if err != nil {
		t.Errorf("Sommething error %s", err)
	}

	err = epubWriter.Write("./temp/new-book.epub")
	if err != nil {
		t.Errorf("Sommething error %s", err)
	}
}
