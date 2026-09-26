package epub_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/raitucarp/epub"
	"github.com/raitucarp/epub/pkg"
)

// Example_roundTrip demonstrates reading an existing publication, adjusting its
// metadata, and writing a new EPUB while preserving the content and table of
// contents.
func Example_roundTrip() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	w := epub.New(book.UID())
	w.Title(append([]string{"Second Edition"}, book.Title()...)...)
	w.Author(book.Author()...)
	w.Languages(book.Language()...)

	for _, res := range book.Spine() {
		w.AddContent(res.Href, res.Content)
	}

	toc, err := book.TableOfContents()
	if err != nil {
		panic(err)
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		panic(err)
	}

	if err := w.Write("second-edition.epub"); err != nil {
		panic(err)
	}
}

// Example_epubToText demonstrates extracting the plain text of every spine
// document in reading order.
func Example_epubToText() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	for _, res := range book.Spine() {
		node := book.ReadContentHTMLById(res.ID)
		fmt.Println("==", res.Href, "==")
		fmt.Println(epub.GetTextContent(node))
	}
}

// Example_epubToMarkdown demonstrates converting every content document to a
// Markdown file on disk.
func Example_epubToMarkdown() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	dir, err := os.MkdirTemp("", "epub-markdown")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	for _, id := range book.ListContentDocumentIds() {
		markdown := book.ReadContentMarkdownById(id)
		path := filepath.Join(dir, id+".md")
		if err := os.WriteFile(path, []byte(markdown), 0o644); err != nil {
			panic(err)
		}
		fmt.Println(path)
	}
}

// Example_bookPipeline demonstrates assembling a complete publication in one
// pass: full metadata, a cover, multiple chapters, guide references, and a
// nested table of contents.
func Example_bookPipeline() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Complete Book")
	w.Author("Jane Doe")
	w.Languages("en")
	w.Description("An example of a complete book.")
	w.Publisher("Example Press")
	w.Rights("All rights reserved")

	if err := w.Cover(samplePNG()); err != nil {
		panic(err)
	}

	w.AddContent("cover.xhtml", sampleChapter("Cover"))
	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	w.AddContent("chapter-2.xhtml", sampleChapter("Chapter 2"))

	w.AddGuide(pkg.GuideRefCover, "cover.xhtml", "Cover")

	toc := epub.TOC{
		Title: "Contents",
		Items: []epub.TOC{
			{Title: "Chapter 1", Href: "chapter-1.xhtml", Items: []epub.TOC{
				{Title: "Section 1.1", Href: "chapter-1.xhtml#section-1-1"},
			}},
			{Title: "Chapter 2", Href: "chapter-2.xhtml"},
		},
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// Example_reproduceStandardEbook demonstrates reproducing a publication from
// the repository's Standard Ebooks sample. The source path is relative to the
// repository root.
func Example_reproduceStandardEbook() {
	src := "tests/data/arthur-conan-doyle_the-white-company.epub"

	book, err := epub.OpenReader(src)
	if err != nil {
		panic(err)
	}

	w := epub.New(strings.Join(book.Identifier(), ", "))
	w.Title(book.Title()...)
	w.Author(book.Author()...)
	w.Languages(book.Language()...)

	if cover, err := book.CoverBytes(); err == nil {
		if err := w.Cover(cover); err != nil {
			panic(err)
		}
	}

	for name, data := range book.ImageResources() {
		w.AddImage(name, data)
	}

	for _, res := range book.Spine() {
		w.AddContent(res.Href, res.Content)
	}

	toc, err := book.TableOfContents()
	if err != nil {
		panic(err)
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		panic(err)
	}

	if err := w.Write("reproduced.epub"); err != nil {
		panic(err)
	}
}
