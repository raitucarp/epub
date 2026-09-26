package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/raitucarp/epub"
)

func main() {
	fmt.Println("=== Writing a Basic EPUB Publication ===")

	// 1. Initialize Writer with a unique publication identifier
	w := epub.New("urn:uuid:simple-write-example")

	// 2. Set core metadata
	w.Title("A Simple Book")
	w.Author("Jane Doe")
	w.Languages("en")
	w.Description("A minimal clean EPUB example.")

	// 3. Add chapter content
	w.AddContent("chapter-1.xhtml", []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Chapter 1</title></head>
<body>
  <h1>Chapter 1: The Beginning</h1>
  <p>Hello from github.com/raitucarp/epub!</p>
</body>
</html>`))

	// 4. Create Table of Contents
	toc := epub.TOC{
		Title: "Table of Contents",
		Items: []epub.TOC{
			{Title: "Chapter 1", Href: "chapter-1.xhtml"},
		},
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		log.Fatalf("Failed to generate TOC: %v", err)
	}

	// 5. Save to disk
	outDir := filepath.Join("temp", "examples_output")
	_ = os.MkdirAll(outDir, 0o755)
	outPath := filepath.Join(outDir, "simple-book.epub")

	if err := w.Write(outPath); err != nil {
		log.Fatalf("Failed to write EPUB: %v", err)
	}

	fmt.Printf("Successfully created: %s\n", outPath)
	fmt.Println("\nFor more comprehensive variants, explore:")
	fmt.Println("  - examples/write/with_cover")
	fmt.Println("  - examples/write/multiple_chapters_and_assets")
	fmt.Println("  - examples/write/nested_toc_and_guide")
	fmt.Println("  - examples/write/markdown")
	fmt.Println("  - examples/write/advanced_multilingual_rtl")
	fmt.Println("  - examples/write/reconstruct_standardebooks")
}
