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
	sourcePath := filepath.Join("tests", "data", "arthur-conan-doyle_the-white-company.epub")
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		sourcePath = filepath.Join("..", "..", "tests", "data", "arthur-conan-doyle_the-white-company.epub")
	}

	fmt.Printf("=== Reconstructing Standard Ebooks EPUB ===\nSource: %s\n\n", sourcePath)

	// 1. Open the original Standard Ebooks EPUB with Reader
	originalReader, err := epub.OpenReader(sourcePath)
	if err != nil {
		log.Fatalf("Failed to open source EPUB: %v", err)
	}

	origTitles := originalReader.Title()
	origAuthors := originalReader.Author()
	origUID := originalReader.UID()
	origLangs := originalReader.Language()

	fmt.Println("Original Standard Ebooks Metadata:")
	fmt.Printf("  Title : %s\n", strings.Join(origTitles, ", "))
	fmt.Printf("  Author: %s\n", strings.Join(origAuthors, ", "))
	fmt.Printf("  UID   : %s\n", origUID)
	fmt.Printf("  Spine : %d documents\n\n", len(originalReader.Spine()))

	// 2. Initialize a new Writer to reconstruct the publication
	w := epub.New(origUID)
	if len(origTitles) > 0 {
		w.Title(origTitles...)
	}
	if len(origAuthors) > 0 {
		w.Author(origAuthors...)
	}
	if len(origLangs) > 0 {
		w.Languages(origLangs...)
	}
	w.Description("Reconstructed Standard Ebooks edition of The White Company by Arthur Conan Doyle.")

	// 3. Transfer the cover image if available
	coverBytes, err := originalReader.CoverBytes()
	if err == nil && len(coverBytes) > 0 {
		if err := w.Cover(coverBytes); err != nil {
			log.Printf("Warning: failed to transfer cover: %v", err)
		} else {
			fmt.Printf("Transferred cover image (%d bytes)\n", len(coverBytes))
		}
	}

	// 4. Transfer image resources
	images := originalReader.ImageResources()
	fmt.Printf("Transferring %d image assets...\n", len(images))
	for id, content := range images {
		// Ignore the cover which was handled separately
		if strings.Contains(strings.ToLower(id), "cover") {
			continue
		}
		w.AddImage(id+".png", content)
	}

	// 5. Transfer content documents from the original spine
	spine := originalReader.Spine()
	fmt.Printf("Transferring spine documents (%d items)...\n", len(spine))

	var tocItems []epub.TOC
	for i, res := range spine {
		filename := filepath.Base(res.Href)
		w.AddContent(filename, res.Content)

		// Create a corresponding TOC entry for notable chapters (first 10 for demonstration)
		if i < 10 {
			itemTitle := fmt.Sprintf("Document %d", i+1)
			if strings.Contains(filename, "chapter-") {
				num := strings.TrimSuffix(strings.TrimPrefix(filename, "chapter-"), ".xhtml")
				itemTitle = fmt.Sprintf("Chapter %s", num)
			} else {
				itemTitle = strings.TrimSuffix(filename, ".xhtml")
			}
			tocItems = append(tocItems, epub.TOC{
				Title: itemTitle,
				Href:  filename,
			})
		}
	}

	// 6. Generate the Table of Contents
	toc := epub.TOC{
		Title: "Table of Contents",
		Items: tocItems,
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		log.Fatalf("Failed to generate TOC: %v", err)
	}

	// 7. Save reconstructed EPUB
	outDir := filepath.Join("temp", "examples_output")
	_ = os.MkdirAll(outDir, 0o755)
	outPath := filepath.Join(outDir, "white-company-reconstructed.epub")

	if err := w.Write(outPath); err != nil {
		log.Fatalf("Failed to write reconstructed EPUB: %v", err)
	}
	fmt.Printf("\nSaved reconstructed EPUB to: %s\n\n", outPath)

	// 8. Verify by reading back the newly reconstructed EPUB
	fmt.Println("Verifying reconstructed EPUB with Reader:")
	reconstructedReader, err := epub.OpenReader(outPath)
	if err != nil {
		log.Fatalf("Failed to open reconstructed EPUB: %v", err)
	}

	fmt.Printf("  Verified Title : %s\n", strings.Join(reconstructedReader.Title(), ", "))
	fmt.Printf("  Verified Author: %s\n", strings.Join(reconstructedReader.Author(), ", "))
	fmt.Printf("  Verified Spine : %d documents\n", len(reconstructedReader.Spine()))
	if reconstructedReader.Cover() != nil {
		bounds := (*reconstructedReader.Cover()).Bounds()
		fmt.Printf("  Verified Cover : %dx%d pixels\n", bounds.Dx(), bounds.Dy())
	}

	fmt.Println("\n=== Standard Ebooks Reconstruction Successful ===")
}
