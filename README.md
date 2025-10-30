# Go EPUB Library

This repository provides a Go (Golang) library for reading and writing EPUB publications. The implementation follows the EPUB 3.3 specification and associated resource processing guidelines as defined by the W3C:

[https://www.w3.org/TR/epub-33/](https://www.w3.org/TR/epub-33/)

The library offers functionality for parsing EPUB containers, extracting metadata, navigating publication resources, reading content documents (XHTML, SVG, Markdown), retrieving images, and working with package renditions.

---

## Features

* Reads EPUB files from disk or byte slices.
* Supports multiple package renditions.
* Extracts metadata such as title, author, and description.
* Provides navigation document and NCX access.
* Reads XHTML, SVG, and derived Markdown representations of content.
* Extracts embedded images and publication resources.
* Allows selecting resources by ID or href.
* Provides access to spine ordering and reading order contents.

---

## Installation

```bash
go get github.com/yourusername/go-epub
```

---

## Basic Usage

```go
package main

import (
	"log"
	"path"
	"github.com/yourusername/go-epub"
)

func main() {
	reader, err := epub.OpenReader(path.Join("books", "example.epub"))
	if err != nil {
		log.Fatalf("failed to open epub: %v", err)
	}

	title := reader.Title()
	author := reader.Author()
	mdDocs := reader.ContentDocumentMarkdown()

	log.Printf("Title: %s", title)
	log.Printf("Author: %s", author)
	for id, doc := range mdDocs {
		log.Printf("Document [%s]: %s", id, doc)
	}
}
```

---

## API Overview

### Types

```
type Epub
type PublicationResource
type Reader
type Writer
```

### Reader Construction

```go
func NewReader(b []byte) (reader Reader, err error)
func OpenReader(name string) (reader Reader, err error)
```

### Metadata and Publication Info

```go
func (r *Reader) Title() (title string)
func (r *Reader) Author() (author string)
func (r *Reader) Description() (description string)
func (r *Reader) UID() (identifier string)
func (r *Reader) Version() (version string)
func (r *Reader) Metadata() (metadata map[string]any)
```

### Content Documents

```go
func (r *Reader) ContentDocumentMarkdown() (documents map[string]string)
func (r *Reader) ContentDocumentSVG() (documents map[string]*html.Node)
func (r *Reader) ContentDocumentXHTML() (documents map[string]*html.Node)
func (r *Reader) ContentDocumentXHTMLString() (documents map[string]string)
func (r *Reader) ReadContentHTMLById(id string) (doc *html.Node)
func (r *Reader) ReadContentMarkdownById(id string) (md string)
```

### Images

```go
func (r *Reader) Cover() (cover *image.Image)
func (r *Reader) Images() (images map[string]image.Image)
func (r *Reader) ReadImageByHref(href string) (img *image.Image)
func (r *Reader) ReadImageById(id string) (img *image.Image)
```

### Navigation and Structure

```go
func (r *Reader) NavigationCenterExtended() *ncx.NCX
func (r *Reader) TableOfContents() (version string)
func (r *Reader) Spine() (orderedResources []PublicationResource)
func (r *Reader) Resources() (resources []PublicationResource)
func (r *Reader) Refines() (refines map[string]map[string][]string)
```

### Package and Resource Selection

```go
func (r *Reader) CurrentSelectedPackage() *pkg.Package
func (r *Reader) CurrentSelectedPackagePath() string
func (r *Reader) SelectPackageRendition(rendition string)
func (r *Reader) SelectResourceByHref(href string) (resource *PublicationResource)
func (r *Reader) SelectResourceById(id string) (resource *PublicationResource)
```

---

## Writer API

The `Writer` type is intended for constructing and exporting EPUB files. Documentation and usage examples will be expanded more soon.

---

## Example

```go
epub, err := epub.OpenReader(path.Join(epubPath, data.name))
if err != nil {
	t.Errorf("Something error %s", err)
	return
}
```

---

## Status

This project is under active development. Interfaces and behavior may change. Contributions, feedback, and issue reports are welcome.

## LICENSE

MIT License

Copyright (c) 2025 Ribhararnus Pracutiar

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
