// Package epub provides a comprehensive, high-performance, and idiomatic Go library
// for reading, writing, and in-place editing of EPUB 2 and EPUB 3 publications.
//
// The package strictly implements modern specifications, including EPUB 3.3,
// the Open Container Format (OCF), the Open Packaging Format (OPF), the EPUB
// Navigation Document (XHTML), and legacy NCX tables of contents. It is verified
// against the W3C EPUB conformance test suite and real-world high-quality books
// from Standard Ebooks.
//
// For in-depth tutorials, guides, and an interactive live web book showcase,
// visit https://epub.raitucarp.name/.
//
// # Architecture & Workflows
//
// The library is organized around three primary paradigms:
//
//  1. Reader: Inspect and extract data from existing EPUB publications (files or memory).
//  2. Writer: Create compliant EPUB publications from scratch or from Markdown documents.
//  3. Editor: Mutate existing EPUB files with fluent method chaining and bidirectional Reader transitions.
//
// # 1. Reading Publications
//
// Open a publication using OpenReader (from disk) or NewReader (from in-memory bytes).
// Use the methods on Reader to access metadata, resources, spine order, and contents:
//
//	book, err := epub.OpenReader("path/to/book.epub")
//	if err != nil {
//		log.Fatalf("failed to open epub: %v", err)
//	}
//
//	// Access metadata
//	fmt.Println("Titles:     ", book.Title())
//	fmt.Println("Authors:    ", book.Author())
//	fmt.Println("Languages:  ", book.Language())
//	fmt.Println("Identifier: ", book.UID())
//	fmt.Println("Version:    ", book.Version())
//
//	// Extract table of contents hierarchy
//	toc, err := book.TableOfContents()
//	if err == nil {
//		for _, item := range toc.Items {
//			fmt.Printf("- %s (%s)\n", item.Title, item.Href)
//		}
//	}
//
//	// Read chapter content directly as clean Markdown
//	md := book.ReadContentMarkdownById("chapter-1")
//	fmt.Println(md)
//
//	// Extract cover image
//	coverImg := book.Cover()
//	if coverImg != nil {
//		// process *image.Image
//	}
//
// # 2. Creating Publications
//
// Initialize a Writer with a unique publication identifier using New.
// Add metadata, XHTML documents, images, styles, and a table of contents,
// then serialize to disk or an in-memory buffer:
//
//	w := epub.New("urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6")
//	w.Title("My Adventure")
//	w.Author("Jane Doe")
//	w.Languages("en")
//	w.Publisher("Open Press")
//	w.Date(time.Now())
//
//	// Add content documents
//	w.AddContent("chapter-1.xhtml", []byte(`<?xml version="1.0" encoding="utf-8"?>
//	<!DOCTYPE html>
//	<html xmlns="http://www.w3.org/1999/xhtml">
//	<head><title>Chapter 1</title></head>
//	<body><h1>Chapter 1</h1><p>The journey begins.</p></body>
//	</html>`))
//
//	// Define navigation
//	w.TableOfContents("toc", epub.TOC{
//		Title: "Table of Contents",
//		Items: []epub.TOC{
//			{Title: "Chapter 1", Href: "chapter-1.xhtml"},
//		},
//	})
//
//	// Write directly to file
//	if err := w.Write("output.epub"); err != nil {
//		log.Fatalf("failed to write epub: %v", err)
//	}
//
//	// Or obtain bytes directly without touching the filesystem:
//	// data, err := w.WriteBytes()
//
// # 3. In-Place Editing
//
// The Editor allows you to load an existing EPUB, modify its metadata, replace or
// append chapters and assets, and write the result back to disk, stream to an
// io.Writer, or write into a byte buffer:
//
//	reader, err := epub.OpenReader("original.epub")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Transition from Reader into fluent Editor mode
//	editor, err := reader.Edit()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Chained metadata mutations
//	editor.Title("Updated Title", "Subtitle").
//		Author("Updated Author").
//		Publisher("Revised Edition Publishing").
//		Subject("Fiction", "Adventure").
//		Language("en").
//		Modified(time.Now())
//
//	// Add or update content
//	editor.AddContent("epilogue.xhtml", []byte(`<h1>Epilogue</h1><p>The end.</p>`))
//
//	// Multiple persistence targets:
//	// Target A: Save directly to disk
//	if err := editor.SaveAs("updated.epub"); err != nil {
//		log.Fatal(err)
//	}
//
//	// Target B: Stream to io.Writer (HTTP response, cloud storage, pipe)
//	// err := editor.Save(writer)
//
//	// Target C: Write into an existing byte buffer
//	// var b []byte
//	// err := editor.WriteBytes(&b)
//
//	// Target D: Seamlessly switch back to Reader mode in memory
//	// r, err := editor.Reader()
//
// # 4. Directory & Markdown Publishing
//
// You can compile entire books directly from directories without manual Go declarations
// by providing a metadata.yml file.
//
// For Markdown manuscripts:
//
//	w := epub.New("")
//	// Ingests all *.md files, static assets, cover image, and metadata.yml
//	if err := w.AddMarkdownDirectory("manuscript"); err != nil {
//		log.Fatal(err)
//	}
//	w.Write("book.epub")
//
// For XHTML/HTML content directories:
//
//	w := epub.New("")
//	// Ingests all *.xhtml files, static assets (css, images, fonts), cover, and metadata.yml
//	if err := w.AddDirectory("book_sources"); err != nil {
//		log.Fatal(err)
//	}
//	w.Write("book.epub")
//
// Both Writer and Editor support AddDirectory and AddMarkdownDirectory.
//
// # 5. Advanced Capabilities
//
//   - Font Obfuscation: Transparently de-obfuscates and decrypts embedded fonts using
//     both the IDPF (http://www.idpf.org/2008/embedding) and Adobe
//     (http://ns.adobe.com/pdf/enc/2006/type) obfuscation algorithms.
//   - Multi-Rendition Publications: Support for publications declaring multiple
//     renditions (e.g. reflowable text vs. fixed-layout comics) via ListRenditions
//     and SelectPackageRendition.
//   - Accessibility & Landmarks: Reading systems can access structural landmarks
//     (bodymatter, cover, toc) and page lists via Landmarks() and PageList().
//
// See the runnable Example functions in this package and the programs under the
// examples/ directory for comprehensive recipes.
package epub
