package epub_test

import (
	"fmt"
	"strings"

	"github.com/raitucarp/epub"
)

// Example demonstrates the core write-then-read workflow. A publication is
// built in memory, serialized to bytes, and read back to verify its metadata.
func Example() {
	w := epub.New("urn:example:book")
	w.Title("The Example Book")
	w.Author("Jane Doe")
	w.Languages("en")

	w.AddContent("chapter-1.xhtml", []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body><h1>Chapter 1</h1><p>Hello, world.</p></body></html>`))
	if err := w.TableOfContents("toc", epub.TOC{
		Title: "Contents",
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

	fmt.Println(strings.Join(book.Title(), ", "))
	fmt.Println(strings.Join(book.Author(), ", "))
	fmt.Println(book.UID())

	// Output:
	// The Example Book
	// Jane Doe
	// urn:example:book
}

// ExampleOpenReader demonstrates opening an EPUB file from disk and reading its
// metadata and table of contents.
func ExampleOpenReader() {
	book, err := epub.OpenReader("book.epub")
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(book.Title(), ", "))
	fmt.Println(book.Version())

	toc, err := book.TableOfContents()
	if err != nil {
		panic(err)
	}
	for _, item := range toc.Items {
		fmt.Println(item.Title, item.Href)
	}
}

// ExampleWriter demonstrates building a publication with a cover image and
// multiple chapters before writing it to disk.
func ExampleWriter() {
	w := epub.New("urn:example:writer")
	w.Title("A Novel")
	w.Author("Jane Doe")
	w.Languages("en")

	w.AddContent("chapter-1.xhtml", []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body><h1>Chapter 1</h1></body></html>`))
	w.AddContent("chapter-2.xhtml", []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body><h1>Chapter 2</h1></body></html>`))

	toc := epub.TOC{
		Title: "Contents",
		Items: []epub.TOC{
			{Title: "Chapter 1", Href: "chapter-1.xhtml"},
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

// ExampleWriter_AddMarkdown demonstrates converting a single Markdown document
// into an EPUB content document.
func ExampleWriter_AddMarkdown() {
	w := epub.New("urn:example:markdown")
	w.Title("A Markdown Book")
	w.Languages("en")

	res, err := w.AddMarkdown("chapter-1.md", []byte("# Chapter 1\n\nOnce upon a time...\n"))
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Href)
	fmt.Println(res.MIMEType)

	// Output:
	// chapter-1.xhtml
	// application/xhtml+xml
}

// ExampleWriter_AddMarkdownDirectory demonstrates building an entire EPUB from
// a directory of Markdown files. The table of contents is derived from the
// headings in each file.
func ExampleWriter_AddMarkdownDirectory() {
	w := epub.New("urn:example:markdown-dir")
	w.Title("A Markdown Directory Book")
	w.Author("Jane Doe")
	w.Languages("en")

	if err := w.AddMarkdownDirectory("manuscript"); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// Example_roundTrip demonstrates reading an existing publication, adjusting its
// metadata, and writing a new EPUB.
func Example_roundTrip() {
	book, err := epub.OpenReader("original.epub")
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
