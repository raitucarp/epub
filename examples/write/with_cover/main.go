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

func main() {
	fmt.Println("=== Creating EPUB with Cover ===")

	w := epub.New("urn:uuid:book-with-cover")
	w.Title("Illustrated Guide")
	w.Author("Jane Smith")
	w.Languages("en")
	w.Description("Demonstrating cover image embedding with PNG and JPEG encoders.")

	// 1. Generate an in-memory PNG cover image (or load from disk with w.CoverFile)
	img := image.NewRGBA(image.Rect(0, 0, 400, 600))
	for y := 0; y < 600; y++ {
		for x := 0; x < 400; x++ {
			// Gradient background
			img.Set(x, y, color.RGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: 180,
				A: 255,
			})
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatalf("Failed to encode PNG: %v", err)
	}

	// Set cover via raw byte slice (detects image/png automatically)
	if err := w.Cover(buf.Bytes()); err != nil {
		log.Fatalf("Failed to set cover: %v", err)
	}

	// 2. Add content
	w.AddContent("intro.xhtml", []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Introduction</title></head>
<body>
  <h1>Introduction</h1>
  <p>This EPUB includes a dedicated cover image registered in the manifest with property <code>cover-image</code>.</p>
</body>
</html>`))

	// 3. Table of Contents
	w.TableOfContents("toc", epub.TOC{
		Title: "Table of Contents",
		Items: []epub.TOC{
			{Title: "Introduction", Href: "intro.xhtml"},
		},
	})

	// 4. Save
	outDir := filepath.Join("temp", "examples_output")
	_ = os.MkdirAll(outDir, 0o755)
	outPath := filepath.Join(outDir, "with-cover.epub")

	if err := w.Write(outPath); err != nil {
		log.Fatalf("Write failed: %v", err)
	}

	fmt.Printf("Successfully created: %s\n", outPath)
}
