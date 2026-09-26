package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/raitucarp/epub"
)

func main() {
	fmt.Println("=== Creating EPUB from Markdown Content ===")

	w := epub.New("urn:uuid:markdown-epub")
	w.Title("Markdown Powered Book")
	w.Author("Markdown Enthusiast")
	w.Languages("en")
	w.Description("An EPUB automatically compiled from GitHub Flavored Markdown.")

	// Chapter 1 in Markdown (GFM tables, bold, lists)
	ch1MD := []byte(`# Chapter 1: Introduction to Markdown in EPUB

Markdown makes writing long-form books enjoyable and readable.

## Features

- **Portability**: Plain text source files.
- **Tables**: GFM syntax supported out of the box.
- **Code Blocks**: Formatted cleanly in XHTML.

| Syntax | Description |
| :--- | :--- |
| Header | Title text |
| Paragraph | Body text |

> "Simplicity is prerequisite for reliability." — Edsger W. Dijkstra
`)

	// Chapter 2 in Markdown
	ch2MD := []byte(`# Chapter 2: Syntax and Structure

Writing code snippets and lists:

1. First step
2. Second step
3. Third step

~~~go
package main

import "fmt"

func main() {
    fmt.Println("Hello, EPUB!")
}
~~~
`)

	res1, err := w.AddMarkdown("chapter-1.md", ch1MD)
	if err != nil {
		log.Fatalf("AddMarkdown ch1 failed: %v", err)
	}

	res2, err := w.AddMarkdown("chapter-2.md", ch2MD)
	if err != nil {
		log.Fatalf("AddMarkdown ch2 failed: %v", err)
	}

	// Build Table of Contents
	toc := epub.TOC{
		Title: "Table of Contents",
		Items: []epub.TOC{
			{Title: "Chapter 1: Introduction", Href: res1.Href},
			{Title: "Chapter 2: Syntax and Structure", Href: res2.Href},
		},
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		log.Fatalf("Failed to generate TOC: %v", err)
	}

	// Save
	outDir := filepath.Join("temp", "examples_output")
	_ = os.MkdirAll(outDir, 0o755)
	outPath := filepath.Join(outDir, "markdown-book.epub")

	if err := w.Write(outPath); err != nil {
		log.Fatalf("Write failed: %v", err)
	}

	fmt.Printf("Successfully created: %s\n", outPath)
}
