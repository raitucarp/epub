package main

import (
	"bytes"
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
		epubPath = filepath.Join("..", "..", "tests", "data", "arthur-conan-doyle_the-white-company.epub")
	}

	fmt.Printf("=== Edit Mode Example ===\nOpening: %s\n\n", epubPath)

	reader, err := epub.OpenReader(epubPath)
	if err != nil {
		log.Fatalf("Failed to open EPUB: %v", err)
	}

	fmt.Println("Original Metadata:")
	fmt.Printf("  Title : %s\n", strings.Join(reader.Title(), ", "))
	fmt.Printf("  Author: %s\n\n", strings.Join(reader.Author(), ", "))

	// 1. Enter edit mode
	editor, err := reader.Edit()
	if err != nil {
		log.Fatalf("Failed to enter edit mode: %v", err)
	}

	// 2. Modify metadata concisely with method chaining
	fmt.Println("Applying metadata modifications...")
	editor.Title("The White Company (Annotated Edition)").
		Author("Sir Arthur Conan Doyle", "Annotated by Raitucarp").
		Description("Special annotated edition of The White Company with reader/editor demo.").
		Subject("Historical Fiction", "Middle Ages", "Adventure", "Archery").
		Publisher("Custom Publishing").
		Language("en")

	// 3. Add an extra introductory chapter
	fmt.Println("Adding a new introduction chapter...")
	introHTML := []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
  <title>Editor's Note</title>
</head>
<body>
  <section epub:type="preface">
    <h2>Editor's Note</h2>
    <p>This chapter was injected dynamically using the new <code>epub.Editor</code> API!</p>
    <p>All existing files, images, stylesheets, and fonts in the EPUB remain intact.</p>
  </section>
</body>
</html>`)

	res, err := editor.AddContent("text/editors-note.xhtml", introHTML)
	if err != nil {
		log.Fatalf("AddContent failed: %v", err)
	}
	fmt.Printf("  Added resource ID '%s' at '%s'\n", res.ID, res.Href)

	// 4. Update an existing file (e.g. imprint)
	fmt.Println("Updating existing imprint...")
	updatedImprint := []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Imprint</title></head>
<body>
  <p>Updated Imprint &copy; 2026. Edited dynamically with github.com/raitucarp/epub.</p>
</body>
</html>`)
	if err := editor.UpdateContent("text/imprint.xhtml", updatedImprint); err != nil {
		log.Fatalf("UpdateContent failed: %v", err)
	}

	// 5. Test returning directly to Reader mode without touching disk
	fmt.Println("\nSwitching back to Reader mode in memory (editor.Reader())...")
	readerBack, err := editor.Reader()
	if err != nil {
		log.Fatalf("editor.Reader() failed: %v", err)
	}

	fmt.Printf("  Updated Title  : %s\n", strings.Join(readerBack.Title(), ", "))
	fmt.Printf("  Updated Authors: %s\n", strings.Join(readerBack.Author(), ", "))
	fmt.Printf("  Spine items    : %d (originally %d)\n", len(readerBack.Spine()), len(reader.Spine()))

	// 6. Demonstrate saving in all 3 supported ways
	fmt.Println("\nSaving output:")

	// A. Save to disk with SaveAs
	outDir := filepath.Join("temp", "examples_output")
	_ = os.MkdirAll(outDir, 0o755)
	outPath := filepath.Join(outDir, "white-company-edited.epub")
	if err := editor.SaveAs(outPath); err != nil {
		log.Fatalf("SaveAs failed: %v", err)
	}
	fmt.Printf("  [1] SaveAs: Saved to %s\n", outPath)

	// B. Save to an io.Writer (e.g. bytes.Buffer or HTTP response)
	var buf bytes.Buffer
	if err := editor.Save(&buf); err != nil {
		log.Fatalf("Save to io.Writer failed: %v", err)
	}
	fmt.Printf("  [2] Save(io.Writer): Wrote %d bytes to buffer\n", buf.Len())

	// C. Write directly into a byte pointer with WriteBytes
	var rawBytes []byte
	if err := editor.WriteBytes(&rawBytes); err != nil {
		log.Fatalf("WriteBytes failed: %v", err)
	}
	fmt.Printf("  [3] WriteBytes(&b): Populated byte slice of %d bytes\n", len(rawBytes))

	fmt.Println("\n=== Edit Mode Demonstration Complete ===")
}
