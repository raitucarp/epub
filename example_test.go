package epub_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"

	"github.com/raitucarp/epub"
)

// Example demonstrates the core write-then-read workflow. A publication is
// built in memory, serialized to bytes, and read back to verify its metadata.
func Example() {
	w := epub.New("urn:example:book")
	w.Title("The Example Book")
	w.Author("Jane Doe")
	w.Languages("en")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	if err := w.TableOfContents("toc", epub.TOC{
		Title: "Contents",
		Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
	}); err != nil {
		panic(err)
	}

	data, err := w.WriteBytes()
	if err != nil {
		panic(err)
	}

	book, err := epub.NewReader(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(book.Title(), ", "))
	fmt.Println(strings.Join(book.Author(), ", "))
	fmt.Println(book.UID())

	// Output:
	// The Example Book
	// Jane Doe
	// urn:example:book
}

// sampleChapter returns a minimal XHTML content document for use in examples.
func sampleChapter(title string) []byte {
	return []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<body>
<h1>` + title + `</h1>
<p>This is a sample paragraph.</p>
</body>
</html>`)
}

// samplePNG returns a small in-memory PNG image for use in the image and cover
// examples.
func samplePNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// buildSampleBook constructs a complete in-memory EPUB and returns its bytes.
// It is used by the reading examples so they can operate without a file on
// disk.
func buildSampleBook() []byte {
	w := epub.New("urn:example:sample")
	w.Title("The Sample Book")
	w.Author("Jane Doe")
	w.Languages("en")

	w.AddContent("chapter-1.xhtml", sampleChapter("Chapter 1"))
	w.AddContent("chapter-2.xhtml", sampleChapter("Chapter 2"))

	toc := epub.TOC{
		Title: "Contents",
		Items: []epub.TOC{
			{Title: "Chapter 1", Href: "chapter-1.xhtml"},
			{Title: "Chapter 2", Href: "chapter-2.xhtml"},
		},
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		panic(err)
	}

	data, err := w.WriteBytes()
	if err != nil {
		panic(err)
	}
	return data
}
