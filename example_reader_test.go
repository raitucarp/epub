package epub_test

import (
	"fmt"
	"strings"

	"github.com/raitucarp/epub"
)

// ExampleOpenReader demonstrates opening an EPUB file from disk.
func ExampleOpenReader() {
	book, err := epub.OpenReader("book.epub")
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(book.Title(), ", "))
	fmt.Println(book.Version())
}

// ExampleNewReader demonstrates reading an EPUB from an in-memory byte slice.
func ExampleNewReader() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(book.Title(), ", "))
	fmt.Println(book.UID())

	// Output:
	// The Sample Book
	// urn:example:sample
}

// ExampleReader_Title demonstrates reading every title declared in the package
// metadata.
func ExampleReader_Title() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	for i, title := range book.Title() {
		fmt.Printf("%d. %s\n", i+1, title)
	}

	// Output:
	// 1. The Sample Book
}

// ExampleReader_Author demonstrates reading the creator metadata.
func ExampleReader_Author() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(book.Author(), ", "))

	// Output:
	// Jane Doe
}

// ExampleReader_Identifier demonstrates reading identifiers and resolving the
// unique identifier.
func ExampleReader_Identifier() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	fmt.Println("Identifiers:", book.Identifier())
	fmt.Println("Unique:", book.UID())

	// Output:
	// Identifiers: [urn:example:sample]
	// Unique: urn:example:sample
}

// ExampleReader_Language demonstrates reading the publication languages.
func ExampleReader_Language() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	fmt.Println(book.Language())

	// Output:
	// [en]
}

// ExampleReader_Version demonstrates reading the EPUB specification version.
func ExampleReader_Version() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	fmt.Println(book.Version())

	// Output:
	// 3.0
}

// ExampleReader_Description demonstrates reading the publication description.
func ExampleReader_Description() {
	book, err := epub.OpenReader("book.epub")
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(book.Description(), "\n"))
}

// ExampleReader_Metadata demonstrates reading the full metadata block and the
// metadata refinements in one pass.
func ExampleReader_Metadata() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	for key, value := range book.Metadata() {
		fmt.Printf("%s: %v\n", key, value)
	}
}

// ExampleReader_ListContentDocumentIds demonstrates listing the manifest IDs of
// all content documents.
func ExampleReader_ListContentDocumentIds() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	for _, id := range book.ListContentDocumentIds() {
		fmt.Println(id)
	}
}

// ExampleReader_ReadContentHTMLById demonstrates reading a content document as
// a parsed HTML node and extracting its text.
func ExampleReader_ReadContentHTMLById() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	node := book.ReadContentHTMLById("chapter-1.xhtml")
	fmt.Println(epub.GetTextContent(node))
}

// ExampleReader_ReadContentMarkdownById demonstrates reading a content document
// converted to Markdown.
func ExampleReader_ReadContentMarkdownById() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	fmt.Println(book.ReadContentMarkdownById("chapter-1.xhtml"))
}

// ExampleReader_ContentDocumentMarkdown demonstrates converting every content
// document to Markdown at once.
func ExampleReader_ContentDocumentMarkdown() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	for id, markdown := range book.ContentDocumentMarkdown() {
		fmt.Printf("== %s ==\n%s\n", id, markdown)
	}
}

// ExampleReader_Resources demonstrates iterating resources and selecting them
// by manifest ID and href.
func ExampleReader_Resources() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	for _, res := range book.Resources() {
		fmt.Printf("%s (%s)\n", res.ID, res.MIMEType)
	}

	if res := book.SelectResourceByHref("chapter-1.xhtml"); res != nil {
		fmt.Println("selected:", res.ID)
	}
}

// ExampleReader_Spine demonstrates reading the default reading order.
func ExampleReader_Spine() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	for _, res := range book.Spine() {
		fmt.Println(res.Href)
	}

	// Output:
	// chapter-1.xhtml
	// chapter-2.xhtml
}

// ExampleReader_TableOfContents demonstrates reading a nested table of
// contents.
func ExampleReader_TableOfContents() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	toc, err := book.TableOfContents()
	if err != nil {
		panic(err)
	}

	fmt.Println(toc.Title)
	for _, item := range toc.Items {
		fmt.Printf("- %s -> %s\n", item.Title, item.Href)
	}

	// Output:
	// Table of Contents
	// - Chapter 1 -> chapter-1.xhtml
	// - Chapter 2 -> chapter-2.xhtml
}

// ExampleTOC_JSON demonstrates serializing the table of contents to JSON.
func ExampleTOC_JSON() {
	book, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	toc, err := book.TableOfContents()
	if err != nil {
		panic(err)
	}

	data, err := toc.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

// ExampleReader_Landmarks demonstrates reading the landmarks navigation.
func ExampleReader_Landmarks() {
	book, err := epub.OpenReader("book.epub")
	if err != nil {
		panic(err)
	}

	for _, landmark := range book.Landmarks() {
		fmt.Printf("%s: %s\n", landmark.Type, landmark.Href)
	}
}

// ExampleReader_PageList demonstrates reading the page-list navigation.
func ExampleReader_PageList() {
	book, err := epub.OpenReader("book.epub")
	if err != nil {
		panic(err)
	}

	for _, page := range book.PageList() {
		fmt.Printf("%s -> %s\n", page.Title, page.Href)
	}
}

// ExampleReader_References demonstrates reading the structural guide
// references.
func ExampleReader_References() {
	book, err := epub.OpenReader("book.epub")
	if err != nil {
		panic(err)
	}

	for kind, node := range book.References() {
		fmt.Printf("%s: %v\n", kind, node != nil)
	}
}

// ExampleReader_SelectPackageRendition demonstrates switching between package
// renditions in a publication that ships multiple layouts.
func ExampleReader_SelectPackageRendition() {
	book, err := epub.OpenReader("book.epub")
	if err != nil {
		panic(err)
	}

	for _, rendition := range book.ListRenditions() {
		fmt.Println(rendition)
	}

	book.SelectPackageRendition("default")
	fmt.Println("Active package:", book.CurrentSelectedPackagePath())
}

// ExampleReader_Cover demonstrates extracting the cover image.
func ExampleReader_Cover() {
	book, err := epub.OpenReader("book.epub")
	if err != nil {
		panic(err)
	}

	cover := book.Cover()
	if cover == nil {
		fmt.Println("no cover")
		return
	}
	fmt.Println("cover bounds:", (*cover).Bounds())
}
