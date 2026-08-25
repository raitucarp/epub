// Command markdown builds an EPUB from a directory of Markdown files and writes
// it to book.epub.
//
// The directory is expected to contain one Markdown file per chapter, named so
// that a lexicographic sort yields the desired reading order, for example:
//
//	manuscript/
//	  chapter_01.md
//	  chapter_02.md
//
// Headings within each file (h1, h2, and so on) are used to build the table of
// contents automatically.
//
// Usage:
//
//	go run ./examples/markdown manuscript
package main

import (
	"log"
	"os"

	"github.com/raitucarp/epub"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: markdown <directory>")
	}

	w := epub.New("urn:example:markdown")
	w.Title("A Markdown Book")
	w.Author("Jane Doe")
	w.Languages("en")

	if err := w.AddMarkdownDirectory(os.Args[1]); err != nil {
		log.Fatal(err)
	}

	if err := w.Write("book.epub"); err != nil {
		log.Fatal(err)
	}
}
