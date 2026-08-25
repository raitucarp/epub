// Command write builds a minimal EPUB from in-memory XHTML content and writes
// it to book.epub.
//
// Usage:
//
//	go run ./examples/write
package main

import (
	"log"

	"github.com/raitucarp/epub"
)

func main() {
	w := epub.New("urn:example:write")
	w.Title("A Small Book")
	w.Author("Jane Doe")
	w.Languages("en")

	chapter := `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<body>
<h1>Chapter 1</h1>
<p>Hello, world.</p>
</body>
</html>`

	w.AddContent("chapter-1.xhtml", []byte(chapter))

	toc := epub.TOC{
		Title: "Contents",
		Items: []epub.TOC{
			{Title: "Chapter 1", Href: "chapter-1.xhtml"},
		},
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		log.Fatal(err)
	}

	if err := w.Write("book.epub"); err != nil {
		log.Fatal(err)
	}
}
