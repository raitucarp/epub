# Go EPUB Library

A Go library for reading and writing EPUB publications.
This library follows the EPUB 3.3 specification: https://www.w3.org/TR/epub-33/

[![Go Reference](https://pkg.go.dev/badge/github.com/raitucarp/epub.svg)](https://pkg.go.dev/github.com/raitucarp/epub)
[![Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/raitucarp)

---

## Features

- Parse EPUB 3.3 container, metadata, manifest, spine, guide, and navigation structures.
- Read content documents in:
  - XHTML
  - Markdown (auto-converted)
  - SVG
- Extract images in raw bytes or `image.Image`.
- Access package-level metadata and renditions.
- Generate new EPUB files programmatically.

---

## Installation

```
go get github.com/yourusername/epub
```

---

## Quick Start (Reading)

```go
package main

import (
	"fmt"
	"log"
	"github.com/yourusername/epub"
)

func main() {
	r, err := epub.OpenReader("example.epub")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Title:", r.Title())
	fmt.Println("Author:", r.Author())
	fmt.Println("Language:", r.Language())

	ids := r.ListContentDocumentIds()
	for _, id := range ids {
		html := r.ReadContentHTMLById(id)
		fmt.Printf("HTML Node for %s: %v\n", id, html)
	}
}
```

---

## Quick Start (Writing)

```go
package main

import (
	"log"
	"time"
	"github.com/yourusername/epub"
)

func main() {
	w := epub.New("pub-id-001")
	w.Title("Example Book")
	w.Author("John Doe")
	w.Language("en")
	w.Date(time.Now())

	w.AddContent("chapter1.xhtml", []byte(`<html><body><h1>Hello World</h1></body></html>`))
	w.Write("output.epub")
}
```

---

## Reader API Overview

```go
type Reader

func NewReader(b []byte) (reader Reader, err error)
func OpenReader(name string) (reader Reader, err error)

func (r *Reader) Title() string
func (r *Reader) Author() string
func (r *Reader) Identifier() string
func (r *Reader) Language() string
func (r *Reader) Metadata() map[string]any

func (r *Reader) ListContentDocumentIds() []string
func (r *Reader) ListImageIds() []string

func (r *Reader) ReadContentHTMLById(id string) *html.Node
func (r *Reader) ReadContentMarkdownById(id string) string

func (r *Reader) ReadImageById(id string) *image.Image
func (r *Reader) ImageResources() map[string][]byte

func (r *Reader) Spine() []PublicationResource
func (r *Reader) Resources() []PublicationResource

func (r *Reader) TableOfContents() (TOC, error)
```

---

## Writer API Overview

```go
type Writer

func New(pubId string) *Writer

func (w *Writer) Title(...string)
func (w *Writer) Author(string)
func (w *Writer) Languages(...string)
func (w *Writer) Date(time.Time)
func (w *Writer) Description(string)
func (w *Writer) Publisher(string)

func (w *Writer) AddContent(filename string, content []byte) PublicationResource
func (w *Writer) AddImage(name string, content []byte) PublicationResource
func (w *Writer) AddSpineItem(res PublicationResource)

func (w *Writer) Write(filename string) error
```

---

## License

MIT License.
