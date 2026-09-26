package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/raitucarp/epub"
)

func main() {
	fmt.Println("=== Creating Basic EPUB ===")

	// 1. Initialize Writer with a unique publication identifier
	w := epub.New("urn:uuid:12345678-1234-5678-1234-567812345678")

	// 2. Set core metadata
	w.Title("My First EPUB")
	w.Author("John Doe")
	w.Languages("en")
	w.Description("A minimal clean EPUB example created with github.com/raitucarp/epub.")

	// 3. Add chapter content (XHTML)
	chapter1 := []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <title>Chapter 1</title>
</head>
<body>
  <h1>Chapter 1: The Beginning</h1>
  <p>This is the first paragraph of our minimal EPUB publication.</p>
  <p>EPUB documents are standard XHTML files organized inside a ZIP archive.</p>
</body>
</html>`)

	w.AddContent("chapter-1.xhtml", chapter1)

	// 4. Create Table of Contents (required by EPUB specification)
	toc := epub.TOC{
		Title: "Table of Contents",
		Items: []epub.TOC{
			{Title: "Chapter 1: The Beginning", Href: "chapter-1.xhtml"},
		},
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		log.Fatalf("Failed to generate TOC: %v", err)
	}

	// 5. Save to disk
	outDir := filepath.Join("temp", "examples_output")
	_ = os.MkdirAll(outDir, 0o755)
	outPath := filepath.Join(outDir, "basic-book.epub")

	if err := w.Write(outPath); err != nil {
		log.Fatalf("Failed to write EPUB: %v", err)
	}

	fmt.Printf("Successfully created: %s\n", outPath)
}
