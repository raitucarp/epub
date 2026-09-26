package epub_test

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"time"

	"github.com/raitucarp/epub"
	"github.com/raitucarp/epub/pkg"
)

// ExampleNew demonstrates creating a Writer with complete metadata and writing
// a minimal publication.
func ExampleNew() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Author("Jane Doe")
	w.Languages("en")
	w.Description("A short description")
	w.Publisher("Example Press")
	w.Rights("All rights reserved")
	w.Date(time.Now())
	w.Modified(time.Now())

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter demonstrates creating an EPUB publication, adding chapters,
// configuring metadata and navigation, and producing the binary archive.
func ExampleWriter() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("Adventures in Go")
	w.Author("Jane Doe")
	w.Languages("en")
	w.Description("A complete guide to building software with Go.")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1: The Beginning"))
	w.AddContent("chapter-2.xhtml", sampleChapter("Chapter 2: Concurrency"))

	if err := w.TableOfContents("toc", epub.TOC{
		Title: "Table of Contents",
		Items: []epub.TOC{
			{Title: "Chapter 1", Href: "chapter-1.xhtml"},
			{Title: "Chapter 2", Href: "chapter-2.xhtml"},
		},
	}); err != nil {
		panic(err)
	}

	data, err := w.WriteBytes()
	if err != nil {
		panic(err)
	}

	// Verify publication was generated
	book, err := epub.NewReader(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(book.Title(), ", "))
	fmt.Println(strings.Join(book.Author(), ", "))
	fmt.Println("Chapters in Spine:", len(book.Spine()))

	// Output:
	// Adventures in Go
	// Jane Doe
	// Chapters in Spine: 2
}

// ExampleWriter_Title demonstrates setting a title with an additional subtitle.
func ExampleWriter_Title() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book", "The Subtitled Edition")
	w.Languages("en")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Author demonstrates setting the primary author.
func ExampleWriter_Author() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Author("Jane Doe")
	w.Languages("en")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Creator demonstrates adding a creator with an explicit
// identifier attribute.
func ExampleWriter_Creator() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Creator("creator-1", "Jane Doe")
	w.Languages("en")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Contributor demonstrates adding a contributor of a given role.
func ExampleWriter_Contributor() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Author("Jane Doe")
	w.Contributor("editor", "John Editor")
	w.Languages("en")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Subject demonstrates adding subject classifications.
func ExampleWriter_Subject() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Author("Jane Doe")
	w.Subject("subject-1", "Fiction", "Adventure")
	w.Languages("en")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Description demonstrates setting a short and a long
// description.
func ExampleWriter_Description() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")
	w.Description("A short summary.")
	w.LongDescription("A much longer and more detailed description of the work.")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Publisher demonstrates setting the publisher.
func ExampleWriter_Publisher() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")
	w.Publisher("Example Press")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Rights demonstrates setting copyright information.
func ExampleWriter_Rights() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")
	w.Rights("Copyright 2026 Jane Doe. All rights reserved.")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Date demonstrates setting the publication and modification
// dates.
func ExampleWriter_Date() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")
	w.Date(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	w.Modified(time.Now())

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Languages demonstrates setting multiple language codes.
func ExampleWriter_Languages() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en", "fr")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Identifiers demonstrates adding multiple identifiers.
func ExampleWriter_Identifiers() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")
	w.Identifiers("urn:isbn:9780000000002", "urn:uuid:0f1f1e7e-0000-0000-0000-000000000000")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_AddContent demonstrates adding multiple chapter documents.
func ExampleWriter_AddContent() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	w.AddContent("chapter-2.xhtml", sampleChapter("Chapter 2"))

	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{
			{Title: "Chapter 1", Href: "chapter-1.xhtml"},
			{Title: "Chapter 2", Href: "chapter-2.xhtml"},
		},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_AddContentFile demonstrates adding a content document read
// from disk.
func ExampleWriter_AddContentFile() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")

	if _, err := w.AddContentFile("chapter-1.xhtml"); err != nil {
		panic(err)
	}

	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_AddImage demonstrates embedding an image resource.
func ExampleWriter_AddImage() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")

	w.AddImage("figure-1.png", samplePNG())

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Cover demonstrates setting the cover from raw image bytes.
func ExampleWriter_Cover() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")

	if err := w.Cover(samplePNG()); err != nil {
		panic(err)
	}

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_CoverPNG demonstrates setting the cover from an image.Image.
func ExampleWriter_CoverPNG() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")

	cover := image.NewRGBA(image.Rect(0, 0, 100, 150))
	cover.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := w.CoverPNG(cover); err != nil {
		panic(err)
	}

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_AddGuide demonstrates adding structural guide references.
func ExampleWriter_AddGuide() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")

	w.AddContent("cover.xhtml", sampleChapter("Cover"))
	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))

	w.AddGuide(pkg.GuideRefCover, "cover.xhtml", "Cover")
	w.AddGuide(pkg.GuideRefText, "chapter-1.xhtml", "Chapter 1")

	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_TableOfContents demonstrates writing a nested table of
// contents.
func ExampleWriter_TableOfContents() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))

	toc := epub.TOC{
		Title: "Contents",
		Items: []epub.TOC{
			{
				Title: "Chapter 1",
				Href:  "chapter-1.xhtml",
				Items: []epub.TOC{
					{Title: "Section 1.1", Href: "chapter-1.xhtml#section-1-1"},
					{Title: "Section 1.2", Href: "chapter-1.xhtml#section-1-2"},
				},
			},
		},
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_Meta demonstrates adding custom metadata and refinements.
func ExampleWriter_Meta() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")

	w.Meta(pkg.Meta{Property: "schema:accessMode", Value: "textual"})
	w.Refines("#title", "title-type", "main")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_DublinCores demonstrates setting Dublin Core fields in bulk.
func ExampleWriter_DublinCores() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")
	w.DublinCores(map[string]string{
		"publisher": "Example Press",
		"rights":    "All rights reserved",
	})

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_WriteBytes demonstrates serializing to an in-memory byte slice
// and reading it back.
func ExampleWriter_WriteBytes() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Book")
	w.Languages("en")
	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	data, err := w.WriteBytes()
	if err != nil {
		panic(err)
	}

	book, err := epub.NewReader(data)
	if err != nil {
		panic(err)
	}
	_ = book
}
