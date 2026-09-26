package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/raitucarp/epub"
)

func main() {
	epubPath := filepath.Join("tests", "data", "arthur-conan-doyle_the-white-company.epub")
	if _, err := os.Stat(epubPath); os.IsNotExist(err) {
		// Fallback for running from within examples/read/ directory
		epubPath = filepath.Join("..", "..", "tests", "data", "arthur-conan-doyle_the-white-company.epub")
	}

	fmt.Printf("=== Reading EPUB: %s ===\n\n", epubPath)

	reader, err := epub.OpenReader(epubPath)
	if err != nil {
		log.Fatalf("Failed to open EPUB: %v", err)
	}

	// 1. Publication Metadata
	fmt.Println("--- Metadata ---")
	fmt.Printf("Title       : %s\n", strings.Join(reader.Title(), ", "))
	fmt.Printf("Author      : %s\n", strings.Join(reader.Author(), ", "))
	fmt.Printf("Language    : %s\n", strings.Join(reader.Language(), ", "))
	fmt.Printf("Identifier  : %s\n", strings.Join(reader.Identifier(), ", "))
	fmt.Printf("UID         : %s\n", reader.UID())
	fmt.Printf("Version     : %s\n", reader.Version())
	if desc, ok := reader.Metadata()["description"]; ok {
		fmt.Printf("Description : %v\n", desc)
	}

	// 2. Cover image
	fmt.Println("\n--- Cover Image ---")
	cover := reader.Cover()
	if cover != nil {
		bounds := (*cover).Bounds()
		fmt.Printf("Cover found: %dx%d pixels\n", bounds.Dx(), bounds.Dy())
	} else {
		fmt.Println("No cover image detected.")
	}

	// 3. Table of Contents
	fmt.Println("\n--- Table of Contents ---")
	toc, err := reader.TableOfContents()
	if err != nil {
		fmt.Printf("Error reading TOC: %v\n", err)
	} else {
		fmt.Printf("TOC Title: %s\n", toc.Title)
		printTOCItems(toc.Items, 1)
	}

	// 4. Reading Order (Spine)
	fmt.Println("\n--- Spine (Reading Order) ---")
	spine := reader.Spine()
	fmt.Printf("Total spine items: %d\n", len(spine))
	for i, item := range spine {
		if i < 5 || i >= len(spine)-2 {
			fmt.Printf("  [%2d] ID: %-25s Href: %s\n", i+1, item.ID, item.Href)
		} else if i == 5 {
			fmt.Println("  ...")
		}
	}

	// 5. Reading Content Documents & Markdown Conversion
	fmt.Println("\n--- Content Sample (First Content Document) ---")
	docIds := reader.ListContentDocumentIds()
	if len(docIds) > 0 {
		firstID := docIds[0]
		md := reader.ReadContentMarkdownById(firstID)
		lines := strings.Split(md, "\n")
		sampleLines := lines
		if len(sampleLines) > 10 {
			sampleLines = sampleLines[:10]
		}
		fmt.Printf("Document ID: %s (first 10 lines converted to Markdown):\n", firstID)
		fmt.Println(strings.Join(sampleLines, "\n"))
	}

	// 6. Resources summary
	fmt.Println("\n--- Resources Summary ---")
	resources := reader.Resources()
	fmt.Printf("Total resources declared: %d\n", len(resources))
	imageIds := reader.ListImageIds()
	fmt.Printf("Total images: %d\n", len(imageIds))
}

func printTOCItems(items []epub.TOC, depth int) {
	indent := strings.Repeat("  ", depth)
	for i, item := range items {
		if i >= 10 && depth == 1 {
			fmt.Printf("%s... (%d more items)\n", indent, len(items)-i)
			break
		}
		fmt.Printf("%s- %s (%s)\n", indent, item.Title, item.Href)
		if len(item.Items) > 0 {
			printTOCItems(item.Items, depth+1)
		}
	}
}
