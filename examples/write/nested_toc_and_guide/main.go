package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/raitucarp/epub"
	"github.com/raitucarp/epub/pkg"
)

func main() {
	fmt.Println("=== Creating EPUB with Nested TOC & Guide References ===")

	w := epub.New("urn:uuid:nested-toc-guide")
	w.Title("Compendium of Knowledge", "Volume 1")
	w.Author("Prof. Robert Langdon")
	w.Languages("en")
	w.Subject("academic", "History", "Science", "Reference")

	// 1. Add Chapters across different parts
	w.AddContent("titlepage.xhtml", []byte(`<html><body><h1>Compendium of Knowledge</h1></body></html>`))
	w.AddContent("part1-ch1.xhtml", []byte(`<html><body><h2>Part 1 - Chapter 1</h2><p>Origins of Science.</p></body></html>`))
	w.AddContent("part1-ch2.xhtml", []byte(`<html><body><h2>Part 1 - Chapter 2</h2><p>Early Methods.</p></body></html>`))
	w.AddContent("part2-ch1.xhtml", []byte(`<html><body><h2>Part 2 - Chapter 1</h2><p>Modern Physics.</p></body></html>`))
	w.AddContent("colophon.xhtml", []byte(`<html><body><h2>Colophon</h2><p>Typeset in Go.</p></body></html>`))

	// 2. Add Guide/Landmarks references (EPUB standard navigation landmarks)
	w.AddGuide(pkg.GuideRefTitlePage, "titlepage.xhtml", "Title Page")
	w.AddGuide(pkg.GuideRefColophon, "colophon.xhtml", "Colophon")

	// 3. Multi-level Nested Table of Contents
	toc := epub.TOC{
		Title: "Table of Contents",
		Items: []epub.TOC{
			{Title: "Title Page", Href: "titlepage.xhtml"},
			{
				Title: "Part I: Classical Era",
				Href:  "part1-ch1.xhtml",
				Items: []epub.TOC{
					{Title: "Chapter 1: Origins of Science", Href: "part1-ch1.xhtml"},
					{Title: "Chapter 2: Early Methods", Href: "part1-ch2.xhtml"},
				},
			},
			{
				Title: "Part II: Modern Era",
				Href:  "part2-ch1.xhtml",
				Items: []epub.TOC{
					{Title: "Chapter 1: Modern Physics", Href: "part2-ch1.xhtml"},
				},
			},
			{Title: "Colophon", Href: "colophon.xhtml"},
		},
	}

	if err := w.TableOfContents("toc", toc); err != nil {
		log.Fatalf("Failed to generate TOC: %v", err)
	}

	// 4. Save
	outDir := filepath.Join("temp", "examples_output")
	_ = os.MkdirAll(outDir, 0o755)
	outPath := filepath.Join(outDir, "nested-toc-guide.epub")

	if err := w.Write(outPath); err != nil {
		log.Fatalf("Write failed: %v", err)
	}

	fmt.Printf("Successfully created: %s\n", outPath)
}
