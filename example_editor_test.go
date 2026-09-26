package epub_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/raitucarp/epub"
)

// ExampleEditor demonstrates loading an existing EPUB publication, editing its
// metadata using fluent method chaining, appending new content, and saving the result.
func ExampleEditor() {
	// Start with an existing publication
	r, err := epub.NewReader(buildSampleBook())
	if err != nil {
		panic(err)
	}

	// Switch from Reader to Editor mode
	editor, err := r.Edit()
	if err != nil {
		panic(err)
	}

	// Chain metadata modifications
	editor.Title("The Sample Book (Second Edition)").
		Author("Jane Doe", "John Smith").
		Subject("Software Engineering", "Go").
		Publisher("Tech Publishing").
		Language("en")

	// Append a new chapter
	if _, err := editor.AddContent("chapter-3.xhtml", sampleChapter("Chapter 3")); err != nil {
		panic(err)
	}

	// Switch back to Reader mode in memory to verify
	updatedReader, err := editor.Reader()
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(updatedReader.Title(), ", "))
	fmt.Println(strings.Join(updatedReader.Author(), ", "))
	fmt.Println("Total Spine Items:", len(updatedReader.Spine()))

	// Output:
	// The Sample Book (Second Edition)
	// Jane Doe, John Smith
	// Total Spine Items: 3
}

// ExampleReader_Edit demonstrates how to read an existing EPUB publication,
// edit its metadata using concise chained methods, and save the result to a byte slice.
func ExampleReader_Edit() {
	// 1. Create a base publication in memory
	w := epub.New("urn:example:orig")
	w.Title("Original Title")
	w.Author("Original Author")
	w.Languages("en")
	w.AddContent("ch1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Title: "Contents",
		Items: []epub.TOC{{Title: "Chapter 1", Href: "ch1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	data, err := w.WriteBytes()
	if err != nil {
		panic(err)
	}

	// 2. Open the publication using Reader
	reader, err := epub.NewReader(data)
	if err != nil {
		panic(err)
	}

	// 3. Obtain an Editor instance
	editor, err := reader.Edit()
	if err != nil {
		panic(err)
	}

	// 4. Modify metadata concisely
	editor.Title("Updated Title").
		Author("Updated Author").
		Language("id")

	// 5. Write the modified EPUB into a byte pointer
	var editedBytes []byte
	if err := editor.WriteBytes(&editedBytes); err != nil {
		panic(err)
	}

	// 6. Read back the modified EPUB to verify changes
	editedReader, err := epub.NewReader(editedBytes)
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(editedReader.Title(), ", "))
	fmt.Println(strings.Join(editedReader.Author(), ", "))
	fmt.Println(strings.Join(editedReader.Language(), ", "))

	// Output:
	// Updated Title
	// Updated Author
	// id
}

// ExampleEditor_Save demonstrates saving an edited publication to an io.Writer.
func ExampleEditor_Save() {
	w := epub.New("urn:example:stream")
	w.Title("Original Book")
	w.Author("Jane Doe")
	w.Languages("en")
	w.AddContent("ch1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Title: "Contents",
		Items: []epub.TOC{{Title: "Chapter 1", Href: "ch1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	data, _ := w.WriteBytes()
	reader, _ := epub.NewReader(data)

	editor, err := reader.Edit()
	if err != nil {
		panic(err)
	}

	// Update title and add a second chapter
	editor.Title("Expanded Edition")
	if _, err := editor.AddContent("ch2.xhtml", sampleChapter("Chapter 2")); err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := editor.Save(&buf); err != nil {
		panic(err)
	}

	editedReader, _ := epub.NewReader(buf.Bytes())
	fmt.Println(strings.Join(editedReader.Title(), ", "))
	fmt.Println("Spine items count:", len(editedReader.Spine()))

	// Output:
	// Expanded Edition
	// Spine items count: 2
}

// ExampleEditor_SaveAs demonstrates reading an existing EPUB, modifying its
// description and subject tags, and saving the result directly to disk.
func ExampleEditor_SaveAs() {
	tempDir, err := os.MkdirTemp("", "epub_example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	origPath := filepath.Join(tempDir, "original.epub")
	w := epub.New("urn:example:saveas")
	w.Title("File-backed Book")
	w.Author("Jane Doe")
	w.Languages("en")
	w.AddContent("ch1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Title: "Contents",
		Items: []epub.TOC{{Title: "Chapter 1", Href: "ch1.xhtml"}},
	}); err != nil {
		panic(err)
	}
	if err := w.Write(origPath); err != nil {
		panic(err)
	}

	// Open with OpenReader
	reader, err := epub.OpenReader(origPath)
	if err != nil {
		panic(err)
	}

	editor, err := reader.Edit()
	if err != nil {
		panic(err)
	}

	editor.Title("Updated File-backed Book").
		Description("A book edited and saved with SaveAs")

	savedPath := filepath.Join(tempDir, "updated.epub")
	if err := editor.SaveAs(savedPath); err != nil {
		panic(err)
	}

	updatedReader, err := epub.OpenReader(savedPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(updatedReader.Title(), ", "))

	// Output:
	// Updated File-backed Book
}

// ExampleEditor_Reader demonstrates switching back to Reader mode after performing edits.
func ExampleEditor_Reader() {
	w := epub.New("urn:example:reader-mode")
	w.Title("Draft Book")
	w.Author("Jane Doe")
	w.Languages("en")
	w.AddContent("ch1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Title: "Contents",
		Items: []epub.TOC{{Title: "Chapter 1", Href: "ch1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	data, _ := w.WriteBytes()
	book, _ := epub.NewReader(data)

	// Enter edit mode
	editor, err := book.Edit()
	if err != nil {
		panic(err)
	}

	editor.Title("Final Published Book").
		Author("Jane Doe, PhD")

	// Switch back to reader mode without saving to a separate file first
	readerBack, err := editor.Reader()
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(readerBack.Title(), ", "))
	fmt.Println(strings.Join(readerBack.Author(), ", "))

	// Output:
	// Final Published Book
	// Jane Doe, PhD
}

// ExampleEditor_WriteBytes demonstrates writing edited EPUB bytes into an existing buffer pointer.
func ExampleEditor_WriteBytes() {
	r, _ := epub.NewReader(buildSampleBook())
	editor, _ := r.Edit()

	editor.Title("Buffer Output Edition")

	var output []byte
	if err := editor.WriteBytes(&output); err != nil {
		panic(err)
	}

	updatedReader, _ := epub.NewReader(output)
	fmt.Println(strings.Join(updatedReader.Title(), ", "))

	// Output:
	// Buffer Output Edition
}

// ExampleEditor_AddMarkdown demonstrates appending a Markdown chapter directly to an existing EPUB.
func ExampleEditor_AddMarkdown() {
	r, _ := epub.NewReader(buildSampleBook())
	editor, _ := r.Edit()

	// Append Markdown content
	md := `# Appendix

This appendix was appended in Markdown!

- Item 1
- Item 2
`
	res, err := editor.AddMarkdown("appendix.md", []byte(md))
	if err != nil {
		panic(err)
	}

	updatedReader, _ := editor.Reader()
	fmt.Println("New resource href:", res.Href)
	fmt.Println("Total spine items:", len(updatedReader.Spine()))

	// Output:
	// New resource href: appendix.xhtml
	// Total spine items: 3
}



