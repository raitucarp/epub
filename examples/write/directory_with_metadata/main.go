package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/raitucarp/epub"
)

func createSamplePNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{R: 1, G: 73, B: 124, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func main() {
	fmt.Println("=== Creating EPUB from Directory with metadata.yml ===")

	// 1. Set up a sample manuscript folder
	workDir, err := os.MkdirTemp("", "epub_dir_example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(workDir)

	// metadata.yml: all publication information is declared declaratively here!
	metaYAML := `identifier: urn:uuid:dir-publishing-sample
title:
  - Declarative Directory Publishing
  - Zero-Boilerplate EPUB Generation
author:
  - Raitucarp Team
  - Contributor
language: en
description: Demonstrating automatic metadata.yml parsing and directory ingestion for both Markdown and XHTML.
publisher: Raitucarp Open Source
rights: CC-BY-SA 4.0
tags:
  - Go
  - EPUB3
  - Markdown
direction: ltr
cover: cover.png
toc_title: Table of Contents
`
	if err := os.WriteFile(filepath.Join(workDir, "metadata.yml"), []byte(metaYAML), 0o644); err != nil {
		log.Fatal(err)
	}

	// Cover image
	if err := os.WriteFile(filepath.Join(workDir, "cover.png"), createSamplePNG(), 0o644); err != nil {
		log.Fatal(err)
	}

	// Stylesheet
	cssContent := `body { font-family: Merriweather, serif; margin: 1.5em; line-height: 1.6; }
h1, h2 { color: #013a63; }
code { background: #f0f4f8; padding: 2px 4px; border-radius: 3px; }`
	if err := os.WriteFile(filepath.Join(workDir, "style.css"), []byte(cssContent), 0o644); err != nil {
		log.Fatal(err)
	}

	// Markdown chapters
	ch1 := `# Chapter 1: Declarative Setup

With **metadata.yml**, you no longer have to manually call:
- ` + "`w.Title(...)`" + `
- ` + "`w.Author(...)`" + `
- ` + "`w.Languages(...)`" + `
- ` + "`w.Identifier(...)`" + `
- ` + "`w.CoverFile(...)`" + `

All of this is inferred from ` + "`metadata.yml`" + ` automatically!

## Nested Headings

Subheadings like this automatically become nested entries in the Table of Contents.
`
	ch2 := `# Chapter 2: Compiling in Go

Just create a ` + "`Writer`" + ` and call ` + "`w.AddMarkdownDirectory(dir)`" + `:

~~~go
w := epub.New("")
if err := w.AddMarkdownDirectory("manuscript"); err != nil {
    log.Fatal(err)
}
w.Write("book.epub")
~~~
`
	if err := os.WriteFile(filepath.Join(workDir, "01-declarative.md"), []byte(ch1), 0o644); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workDir, "02-compiling.md"), []byte(ch2), 0o644); err != nil {
		log.Fatal(err)
	}

	// 2. Compile directly from directory!
	w := epub.New("")
	if err := w.AddMarkdownDirectory(workDir); err != nil {
		log.Fatalf("AddMarkdownDirectory failed: %v", err)
	}

	outDir := filepath.Join("temp", "examples_output")
	_ = os.MkdirAll(outDir, 0o755)
	outPath := filepath.Join(outDir, "directory-metadata-book.epub")

	if err := w.Write(outPath); err != nil {
		log.Fatalf("Failed to write EPUB: %v", err)
	}

	fmt.Printf("Successfully created EPUB from directory: %s\n", outPath)

	// 3. Inspect the resulting EPUB to verify metadata
	r, err := epub.OpenReader(outPath)
	if err != nil {
		log.Fatalf("Failed to open reader: %v", err)
	}

	fmt.Printf("  Title      : %v\n", r.Title())
	fmt.Printf("  Author     : %v\n", r.Author())
	fmt.Printf("  Language   : %v\n", r.Language())
	fmt.Printf("  Description: %v\n", r.Description())
	fmt.Printf("  Spine items: %d\n", len(r.Spine()))
}
