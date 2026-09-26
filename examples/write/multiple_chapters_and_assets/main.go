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
	fmt.Println("=== Creating EPUB with Multiple Chapters & Assets ===")

	w := epub.New("urn:uuid:multi-chapter-assets")
	w.Title("Modern Web & Graphics")
	w.Author("Alice Developer")
	w.Languages("en")
	w.Description("Comprehensive example with multiple chapters, embedded PNG images, and SVG illustrations.")

	// 1. Add an embedded PNG image
	icon := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			icon.Set(x, y, color.RGBA{R: 30, G: 144, B: 255, A: 255})
		}
	}
	var imgBuf bytes.Buffer
	_ = png.Encode(&imgBuf, icon)
	w.AddImage("icon.png", imgBuf.Bytes())

	// 2. Add SVG vector graphic content
	svgContent := []byte(`<?xml version="1.0" encoding="utf-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100" height="100">
  <circle cx="50" cy="50" r="40" stroke="green" stroke-width="4" fill="yellow" />
</svg>`)
	w.AddContent("diagram.svg", svgContent)

	// 3. Add Multiple Chapters
	w.AddContent("chapter-1.xhtml", []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Chapter 1: The Landscape</title></head>
<body>
  <h1>Chapter 1: The Landscape</h1>
  <p>In this chapter, we explore how assets are referenced inside EPUB packages.</p>
  <p><img src="images/icon.png" alt="Sample Icon" /></p>
</body>
</html>`))

	w.AddContent("chapter-2.xhtml", []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Chapter 2: Vector Art</title></head>
<body>
  <h1>Chapter 2: Vector Art</h1>
  <p>EPUB 3 natively supports embedded SVG diagrams as content documents or inline markup.</p>
</body>
</html>`))

	// 4. Hierarchical Table of Contents
	toc := epub.TOC{
		Title: "Table of Contents",
		Items: []epub.TOC{
			{Title: "Chapter 1: The Landscape", Href: "chapter-1.xhtml"},
			{Title: "Chapter 2: Vector Art", Href: "chapter-2.xhtml"},
			{Title: "Appendix: Diagram", Href: "diagram.svg"},
		},
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		log.Fatalf("Failed to generate TOC: %v", err)
	}

	// 5. Save
	outDir := filepath.Join("temp", "examples_output")
	_ = os.MkdirAll(outDir, 0o755)
	outPath := filepath.Join(outDir, "multi-chapter-assets.epub")

	if err := w.Write(outPath); err != nil {
		log.Fatalf("Write failed: %v", err)
	}

	fmt.Printf("Successfully created: %s\n", outPath)
}
