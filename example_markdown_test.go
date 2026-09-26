package epub_test

import (
	"fmt"

	"github.com/raitucarp/epub"
)

// ExampleWriter_AddMarkdown demonstrates converting a Markdown string into a
// content document.
func ExampleWriter_AddMarkdown() {
	w := epub.New("urn:isbn:9780000000001")
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

// ExampleWriter_AddMarkdownFile demonstrates converting a Markdown file read
// from disk into a content document.
func ExampleWriter_AddMarkdownFile() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Markdown Book")
	w.Languages("en")

	if _, err := w.AddMarkdownFile("chapter-1.md"); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}

// ExampleWriter_AddMarkdownDirectory demonstrates building an EPUB from a
// directory of Markdown files, with the table of contents derived from the
// headings in each file.
func ExampleWriter_AddMarkdownDirectory() {
	w := epub.New("urn:isbn:9780000000001")
	w.Title("A Markdown Book")
	w.Author("Jane Doe")
	w.Languages("en")

	if err := w.AddMarkdownDirectory("manuscript"); err != nil {
		panic(err)
	}

	if err := w.Write("book.epub"); err != nil {
		panic(err)
	}
}
