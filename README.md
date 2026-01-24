# Go EPUB Library

<div align="center">

![EPUB](https://img.shields.io/badge/EPUB-3.3-blue)
[![Go Reference](https://pkg.go.dev/badge/github.com/raitucarp/epub.svg)](https://pkg.go.dev/github.com/raitucarp/epub)
[![Go Version](https://img.shields.io/badge/Go-1.25.3-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE.md)
[![GitHub](https://img.shields.io/badge/GitHub-raitucarp%2Fepub-black.svg)](https://github.com/raitucarp/epub)

A robust, feature-rich Go library for reading and writing EPUB publications with full support for the EPUB 3.3 specification.

[📚 Documentation](#reader-api-overview) • [🚀 Quick Start](#quick-start) • [✨ Features](#features) • [📦 Installation](#installation)

</div>

---

## 📋 Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Reading EPUB Files](#reading-epub-files)
  - [Writing EPUB Files](#writing-epub-files)
- [Core Capabilities](#core-capabilities)
  - [Metadata Access](#metadata-access)
  - [Content Processing](#content-processing)
  - [Resource Management](#resource-management)
  - [Navigation Support](#navigation-support)
  - [Image Handling](#image-handling)
- [API Overview](#reader-api-overview)
  - [Reader API](#reader-api)
  - [Writer API](#writer-api)
- [Advanced Usage](#advanced-usage)
- [Supported Content Types](#supported-content-types)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

---

## ✨ Features

### Core Reading Capabilities

- ✅ Parse EPUB 3.3 container, metadata, manifest, spine, guide, and navigation structures
- ✅ Full support for EPUB 2 and EPUB 3 specifications
- ✅ Multiple package rendition support for flexible reading layouts

### Content Processing

- 📄 Read content documents in multiple formats:
  - XHTML (native support)
  - Markdown (auto-converted from HTML)
  - SVG (scalable vector graphics)
- 🔄 Automatic HTML-to-Markdown conversion with customizable options
- 🎯 Precise HTML node parsing and manipulation

### Resource Management

- 📦 Extract and access all publication resources (images, stylesheets, fonts, etc.)
- 🖼️ Extract images in raw bytes or as `image.Image` objects
- 🔍 Query resources by ID or href reference

### Metadata & Navigation

- 🏷️ Complete metadata access (title, author, language, identifier, etc.)
- 📑 Table of Contents (TOC) support for both NAV (EPUB 3) and NCX (EPUB 2)
- 🗺️ Navigation structure abstraction for seamless cross-version compatibility
- 📊 JSON serializable TOC for external tools and integrations

### Writing Capabilities

- 🛠️ Generate new EPUB files programmatically
- 📝 Add content documents (XHTML/HTML)
- 🖼️ Embed images with automatic format detection
- ⚙️ Full package metadata configuration
- 📚 Automatic spine and manifest generation

---

## 📦 Installation

```bash
go get github.com/raitucarp/epub
```

**Requires Go 1.25.3 or higher**

---

## 🚀 Quick Start

### Reading EPUB Files

```go
package main

import (
	"fmt"
	"log"
	"github.com/raitucarp/epub"
)

func main() {
	// Open an EPUB file
	r, err := epub.OpenReader("example.epub")
	if err != nil {
		log.Fatal(err)
	}

	// Access basic metadata
	fmt.Println("Title:", r.Title())
	fmt.Println("Author:", r.Author())
	fmt.Println("Language:", r.Language())
	fmt.Println("Identifier:", r.Identifier())
	fmt.Println("Version:", r.Version())

	// Iterate through content documents
	ids := r.ListContentDocumentIds()
	for _, id := range ids {
		html := r.ReadContentHTMLById(id)
		fmt.Printf("Content for %s: %v\n", id, html)
	}

	// Access table of contents
	toc := r.TOC()
	fmt.Printf("TOC: %s\n", toc.Title)
}
```

### Writing EPUB Files

```go
package main

import (
	"log"
	"time"
	"github.com/raitucarp/epub"
)

func main() {
	// Create a new EPUB writer
	w := epub.New("pub-id-001")
	w.Title("My First Book")
	w.Author("Jane Doe")
	w.Language("en")
	w.Date(time.Now())
	w.Description("An example EPUB publication")
	w.Publisher("Indie Press")

	// Add HTML content
	w.AddContent("chapter1.xhtml", []byte(`
		<html>
			<body>
				<h1>Chapter 1: The Beginning</h1>
				<p>Once upon a time...</p>
			</body>
		</html>
	`))

	// Add an image
	imageData := []byte{/* ... */}
	w.AddImageContent("cover.png", imageData)

	// Write to disk
	if err := w.Write("output.epub"); err != nil {
		log.Fatal(err)
	}
}
```

---

## 🎯 Core Capabilities

### Metadata Access

Access comprehensive publication metadata:

```go
r, _ := epub.OpenReader("book.epub")

// Basic metadata
title := r.Title()
author := r.Author()
language := r.Language()
identifier := r.Identifier()
uid := r.UID()                    // Unique identifier
version := r.Version()            // EPUB version
metadata := r.Metadata()          // Full metadata map
cover := r.GetCover()             // Cover image
```

### Content Processing

Read and process content in multiple formats:

```go
// Get HTML content node
htmlNode := r.ReadContentHTMLById("chapter1")
htmlByHref := r.ReadContentHTMLByHref("text/chapter1.xhtml")

// Convert HTML to Markdown
markdown := r.ReadContentMarkdownById("chapter1")
markdownByHref := r.ReadContentMarkdownByHref("text/chapter1.xhtml")

// Access raw content
rawContent := r.ReadContentById("chapter1")
```

### Resource Management

Manage and access publication resources:

```go
// Get all resources
resources := r.Resources()

// Select specific resources
resource := r.SelectResourceById("img001")
resource := r.SelectResourceByHref("images/cover.jpg")

// Extract images
imageObj := r.ReadImageById("cover-image")
imageObj := r.ReadImageByHref("images/cover.jpg")
imageBytes := r.ReadImageBytesById("cover-image")
```

### Navigation Support

Work with table of contents:

```go
// Get TOC
toc := r.TOC()
fmt.Println("Title:", toc.Title)
fmt.Println("Href:", toc.Href)

// Access nested items
for _, item := range toc.Items {
	fmt.Println("- ", item.Title, "->", item.Href)
}

// Serialize to JSON
jsonBytes, _ := toc.JSON()

// Select multiple renditions
r.SelectPackageRendition("default")
r.SelectPackageRendition("alternative")
currentPackage := r.CurrentSelectedPackage()
```

### Image Handling

Work with images in multiple formats:

```go
// Supported formats: JPEG, PNG, GIF, WebP, SVG

// Get image as image.Image
img := r.ReadImageById("image001")

// Get raw bytes
bytes := r.ReadImageBytesById("image001")

// Get cover image
cover := r.GetCover()

// List image resources
for _, resource := range r.Resources() {
	if isImageType(resource.MIMEType) {
		img := r.ReadImageByHref(resource.Href)
		// Process image...
	}
}
```

---

## 📚 Reader API Overview

```go
type Reader

// Constructor functions
func NewReader(b []byte) (reader Reader, err error)
func OpenReader(name string) (reader Reader, err error)

// Metadata methods
func (r *Reader) Title() string
func (r *Reader) Author() string
func (r *Reader) Identifier() string
func (r *Reader) Language() string
func (r *Reader) UID() string
func (r *Reader) Version() string
func (r *Reader) Metadata() map[string]any
func (r *Reader) GetCover() *image.Image

// Content access
func (r *Reader) ReadContentHTMLById(id string) *html.Node
func (r *Reader) ReadContentHTMLByHref(href string) *html.Node
func (r *Reader) ReadContentMarkdownById(id string) string
func (r *Reader) ReadContentMarkdownByHref(href string) string
func (r *Reader) ReadContentById(id string) []byte
func (r *Reader) ListContentDocumentIds() []string

// Resource management
func (r *Reader) Resources() []PublicationResource
func (r *Reader) SelectResourceById(id string) *PublicationResource
func (r *Reader) SelectResourceByHref(href string) *PublicationResource

// Image handling
func (r *Reader) ReadImageById(id string) *image.Image
func (r *Reader) ReadImageByHref(href string) *image.Image
func (r *Reader) ReadImageBytesById(id string) []byte
func (r *Reader) ReadImageBytesByHref(href string) []byte

// Navigation
func (r *Reader) TOC() *TOC
func (r *Reader) SelectPackageRendition(rendition string)
func (r *Reader) CurrentSelectedPackage() *pkg.Package

// Package selection
func (r *Reader) SelectPackageRendition(rendition string)
func (r *Reader) CurrentSelectedPackagePath() string
```

---

## 📐 Writer API Overview

```go
type Writer

// Constructor
func New(pubId string) *Writer

// Metadata configuration
func (w *Writer) Title(title string) *Writer
func (w *Writer) Author(author string) *Writer
func (w *Writer) Language(lang string) *Writer
func (w *Writer) Description(desc string) *Writer
func (w *Writer) Publisher(pub string) *Writer
func (w *Writer) Date(date time.Time) *Writer

// Content management
func (w *Writer) AddContent(href string, content []byte) (id string, err error)
func (w *Writer) AddImageContent(href string, imageData []byte) (id string, err error)
func (w *Writer) AddCover(imagePath string) (id string, err error)

// Output
func (w *Writer) Write(filename string) error
func (w *Writer) WriteBytes() ([]byte, error)

// Advanced options
func (w *Writer) SetTextDirection(direction string) *Writer
func (w *Writer) SetContentDir(dir string) *Writer
```

---

## 💡 Advanced Usage

### Working with Multiple Renditions

Some EPUB files contain multiple renditions (different layouts, languages, etc.):

```go
r, _ := epub.OpenReader("book.epub")

// Switch between available renditions
r.SelectPackageRendition("default")
r.SelectPackageRendition("alternative")

// Get current package information
pkg := r.CurrentSelectedPackage()
fmt.Println("Manifest items:", len(pkg.Manifest.Items))
```

### Full Metadata Extraction

```go
metadata := r.Metadata()

// The metadata map contains:
// - Basic: title, author, language, identifier
// - Extended: publisher, rights, description, contributor
// - Meta: custom properties and relationships
// - Links: external references

for key, value := range metadata {
	fmt.Printf("%s: %v\n", key, value)
}
```

### Building EPUBs from Scratch

```go
w := epub.New("unique-pub-id")
w.Title("Novel: The Rising Sun")
w.Author("John Smith")
w.Language("en")
w.Publisher("Great Novels Inc")
w.Rights("© 2024 John Smith. All rights reserved.")
w.Description("An epic tale of adventure and discovery")

// Add chapters
chapters := []string{"chapter1.html", "chapter2.html", "chapter3.html"}
for _, ch := range chapters {
	content := readFile(ch)
	w.AddContent(ch, content)
}

// Add cover image
cover := readFile("cover.jpg")
w.AddCover(cover)

// Write final EPUB
w.Write("novel.epub")
```

---

## 📋 Supported Content Types

### Media Types

- **Documents**: `application/xhtml+xml` (XHTML)
- **Navigation**: `application/x-dtbncx+xml` (NCX), `application/nav+xml` (HTML NAV)
- **Styles**: `text/css` (CSS Stylesheets)
- **Images**: `image/jpeg`, `image/png`, `image/gif`, `image/webp`, `image/svg+xml`
- **Fonts**: `font/ttf`, `font/otf`, `application/font-woff`

### Spine Directions

- `ltr` - Left-to-Right (default)
- `rtl` - Right-to-Left
- `default` - Browser default

### Properties

- `nav` - Navigation document
- `cover-image` - Cover image
- `mathml` - MathML support
- `svg` - SVG support
- `scripted` - JavaScript support
- `remote-resources` - External resources
- `layout-pre-paginated` - Fixed layout

---

## 📖 Examples

### Example 1: Extract All Text from an EPUB

```go
package main

import (
	"fmt"
	"log"
	"strings"
	"github.com/raitucarp/epub"
	"golang.org/x/net/html"
)

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c) + " "
	}
	return text
}

func main() {
	r, err := epub.OpenReader("book.epub")
	if err != nil {
		log.Fatal(err)
	}

	for _, id := range r.ListContentDocumentIds() {
		node := r.ReadContentHTMLById(id)
		text := extractText(node)
		fmt.Println(text)
	}
}
```

### Example 2: Create an EPUB from Markdown Files

```go
package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
	"github.com/raitucarp/epub"
)

func main() {
	w := epub.New("markdown-book-001")
	w.Title("My Markdown Book")
	w.Author("Author Name")
	w.Language("en")
	w.Date(time.Now())

	// Convert markdown files to HTML and add to EPUB
	files, _ := filepath.Glob("chapters/*.md")
	for _, file := range files {
		content, _ := os.ReadFile(file)
		// In production, convert markdown to HTML
		w.AddContent(filepath.Base(file), content)
	}

	w.Write("output.epub")
}
```

### Example 3: Copy and Modify an EPUB

```go
package main

import (
	"log"
	"time"
	"github.com/raitucarp/epub"
)

func main() {
	// Read original
	r, err := epub.OpenReader("original.epub")
	if err != nil {
		log.Fatal(err)
	}

	// Create new EPUB with modified metadata
	w := epub.New("new-pub-id")
	w.Title(r.Title() + " (Edition 2)")
	w.Author(r.Author())
	w.Language(r.Language())
	w.Date(time.Now())

	// Copy all content
	for _, id := range r.ListContentDocumentIds() {
		content := r.ReadContentById(id)
		w.AddContent(id, content)
	}

	// Add new cover
	cover := r.GetCover()
	if cover != nil {
		w.AddCover("new-cover.png")
	}

	w.Write("modified.epub")
}
```

---

## 🔧 Dependencies

- **html-to-markdown**: HTML to Markdown conversion
- **golang.org/x/image**: Image processing (GIF, WebP support)
- **golang.org/x/net/html**: HTML parsing
- **golang.org/x/text**: Text normalization and Unicode handling

All dependencies are vendored and documented in `go.mod`.

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

```bash
git clone https://github.com/raitucarp/epub.git
cd epub
go mod download
go test ./...
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

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
````

---

## License

MIT License.
