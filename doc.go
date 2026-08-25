// Package epub provides a reader and writer for EPUB publications.
//
// The package implements the EPUB 3.3 specification: it parses the Open
// Container Format (OCF), the package document, the EPUB navigation document,
// and the legacy NCX format, and exposes a high-level API for inspecting and
// producing publications.
//
// # Reading
//
// Open a publication from a file or from memory with OpenReader or NewReader,
// then use the accessor methods on Reader to inspect its metadata, resources,
// table of contents, and content documents:
//
//	book, err := epub.OpenReader("book.epub")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(book.Title(), book.Author(), book.Language())
//
// # Writing
//
// Build a publication with a Writer, add content documents and metadata, then
// serialize it to a file or to memory with Write or WriteBytes:
//
//	w := epub.New("urn:isbn:9780000000001")
//	w.Title("A Book")
//	w.Author("Jane Doe")
//	w.Languages("en")
//	w.AddContent("chapter-1.xhtml", []byte("<html><body><h1>Chapter 1</h1></body></html>"))
//	w.TableOfContents("toc", epub.TOC{Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}}})
//	if err := w.Write("book.epub"); err != nil {
//		log.Fatal(err)
//	}
//
// # Markdown
//
// Content documents can also be produced from Markdown. AddMarkdown converts
// a single document, while AddMarkdownDirectory converts every Markdown file
// in a directory and derives a table of contents from the document headings:
//
//	w := epub.New("urn:isbn:9780000000001")
//	w.Title("A Markdown Book")
//	w.Languages("en")
//	if err := w.AddMarkdownDirectory("manuscript"); err != nil {
//		log.Fatal(err)
//	}
//	w.Write("book.epub")
//
// See the example functions in this package and the programs under the
// examples directory for further usage.
package epub
