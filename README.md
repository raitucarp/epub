# epub

[![Go Reference](https://pkg.go.dev/badge/github.com/raitucarp/epub.svg)](https://pkg.go.dev/github.com/raitucarp/epub)
[![CI](https://github.com/raitucarp/epub/actions/workflows/ci.yml/badge.svg)](https://github.com/raitucarp/epub/actions/workflows/ci.yml)
[![Ko-fi](https://img.shields.io/badge/Ko--fi-Support-ff5e5b?logo=ko-fi&logoColor=white)](https://ko-fi.com/raitucarp)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE.md)

`epub` is a high-performance, developer-friendly Go library for reading, writing, and editing [EPUB](https://www.w3.org/TR/epub-33/) publications. It implements the EPUB 3.3 specification, including the Open Container Format (OCF), package documents, EPUB navigation documents (NAV), and the legacy NCX format used by EPUB 2.

Documentation is available at **[epub.raitucarp.name](https://epub.raitucarp.name)**.

The library provides three main entry points:

- [`Reader`](https://pkg.go.dev/github.com/raitucarp/epub#Reader) for inspecting metadata, resources, navigation, spine order, and converting content documents to Markdown.
- [`Writer`](https://pkg.go.dev/github.com/raitucarp/epub#Writer) for building new publications from scratch, including from Markdown manuscripts.
- [`Editor`](https://pkg.go.dev/github.com/raitucarp/epub#Editor) for mutating existing publications in place with concise fluent chaining (`reader.Edit()`).

---

## Installation

```sh
go get github.com/raitucarp/epub
```

Requires Go 1.25 or later.

---

## Quick start

### 1. Read an EPUB

```go
book, err := epub.OpenReader("book.epub")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Title:", strings.Join(book.Title(), ", "))
fmt.Println("Author:", strings.Join(book.Author(), ", "))
fmt.Println("Language:", strings.Join(book.Language(), ", "))
fmt.Println("Identifier:", book.UID())

// Convert content documents directly to clean Markdown!
for _, id := range book.ListContentDocumentIds() {
    fmt.Println(book.ReadContentMarkdownById(id))
}
```

### 2. Write an EPUB

```go
w := epub.New("urn:isbn:9780000000001")
w.Title("A Book")
w.Author("Jane Doe")
w.Languages("en")

w.AddContent("chapter-1.xhtml", []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body><h1>Chapter 1</h1><p>Hello, world.</p></body></html>`))

toc := epub.TOC{
    Title: "Contents",
    Items: []epub.TOC{{Title: "Chapter 1", Href: "chapter-1.xhtml"}},
}
if err := w.TableOfContents("toc", toc); err != nil {
    log.Fatal(err)
}

if err := w.Write("book.epub"); err != nil {
    log.Fatal(err)
}
```

### 3. Edit an Existing EPUB

```go
reader, err := epub.OpenReader("book.epub")
if err != nil {
    log.Fatal(err)
}

// Enter edit mode
editor, err := reader.Edit()
if err != nil {
    log.Fatal(err)
}

// Concise fluent method chaining
editor.Title("A Book (Second Edition)").
    Author("Jane Doe, PhD").
    Description("Revised edition with new notes.")

// Add or update chapters
editor.AddContent("chapter-2.xhtml", []byte(`<h1>Chapter 2</h1><p>Continuing the journey...</p>`))

// Save in one of three ways:
// 1. Direct to file
editor.SaveAs("book-v2.epub")

// 2. Stream to any io.Writer
var buf bytes.Buffer
editor.Save(&buf)

// 3. Populate a byte slice directly
var b []byte
editor.WriteBytes(&b)

// 4. Or transition back to Reader mode in memory:
readerBack, err := editor.Reader()
```

### 4. Write an EPUB from Markdown

A single Markdown document can be converted directly:

```go
w := epub.New("urn:isbn:9780000000001")
w.Title("A Markdown Book")
w.Languages("en")

if _, err := w.AddMarkdown("chapter-1.md", []byte("# Chapter 1\n\nOnce upon a time...\n")); err != nil {
    log.Fatal(err)
}
```

Or an entire directory of Markdown files can be assembled, with the table of contents automatically derived from the headings:

```go
w := epub.New("urn:isbn:9780000000001")
w.Title("A Markdown Book")
w.Author("Jane Doe")
w.Languages("en")

if err := w.AddMarkdownDirectory("manuscript"); err != nil {
    log.Fatal(err)
}

if err := w.Write("book.epub"); err != nil {
    log.Fatal(err)
}
```

---

## Reading

### Metadata

```go
book, _ := epub.OpenReader("book.epub")

book.Title()        // []string
book.Author()       // []string
book.Language()     // []string
book.Identifier()   // []string
book.UID()          // the resolved unique identifier
book.Version()      // the EPUB version, e.g. "3.0"
book.Description()  // []string
book.Metadata()     // map[string]any of the full metadata block
book.Refines()      // metadata refinements keyed by subject
```

### Content documents

```go
ids := book.ListContentDocumentIds() // manifest IDs of XHTML/SVG documents

book.ReadContentHTMLById(id)       // *html.Node
book.ReadContentHTMLByHref(href)   // *html.Node
book.ReadContentMarkdownById(id)   // string (Markdown)
book.ReadContentMarkdownByHref(h)  // string (Markdown)
book.ContentDocumentXHTML()        // map[string]*html.Node
book.ContentDocumentXHTMLString()  // map[string]string
book.ContentDocumentMarkdown()     // map[string]string
book.ContentDocumentSVG()          // map[string]*html.Node
```

### Resources and images

```go
book.Resources()                   // []PublicationResource
book.SelectResourceById(id)        // *PublicationResource
book.SelectResourceByHref(href)    // *PublicationResource

book.ListImageIds()                // manifest IDs of image resources
book.Images()                      // map[string]image.Image
book.ImageResources()              // map[string][]byte
book.ReadImageById(id)             // *image.Image
book.ReadImageByHref(href)         // *image.Image
book.ReadImageBytesById(id)        // []byte
book.ReadImageBytesByHref(href)    // []byte
```

### Navigation

```go
toc, _ := book.TableOfContents()   // TOC (NAV or NCX)
book.Landmarks()                   // []Landmark
book.PageList()                    // []Landmark

for _, item := range toc.Items {
    fmt.Println(item.Title, item.Href)
}

json, _ := toc.JSON()              // serialized table of contents
```

---

## Writing

### Metadata

```go
w := epub.New("urn:isbn:9780000000001")

w.Title("A Book")
w.Title("A Book", "Subtitle")       // additional titles
w.Author("Jane Doe")
w.Creator("creator-id", "Jane Doe")
w.Contributor("editor", "John Editor")
w.Languages("en")
w.Description("A short description")
w.Publisher("Example Press")
w.Subject("subject-id", "Fiction")
w.Rights("All rights reserved")
w.Date(time.Now())
w.Modified(time.Now())              // dcterms:modified
```

### Content and resources

```go
w.AddContent("chapter-1.xhtml", []byte("..."))   // XHTML or SVG
w.AddContentFile("chapter-1.xhtml")              // read from disk
w.AddImage("cover.png", imageBytes)
w.AddImageFile("cover.png")
w.Cover(imageBytes)                              // detect PNG/JPEG
w.CoverPNG(img)                                  // image.Image as PNG
w.CoverJPG(img)                                  // image.Image as JPEG
w.CoverFile("cover.png")
```

### Table of contents

```go
toc := epub.TOC{
    Title: "Contents",
    Items: []epub.TOC{
        {Title: "Chapter 1", Href: "chapter-1.xhtml", Items: []epub.TOC{
            {Title: "Section 1.1", Href: "chapter-1.xhtml#section-1-1"},
        }},
    },
}
w.TableOfContents("toc", toc)
```

---

## Editing

```go
editor, err := reader.Edit()

// Mutate metadata
editor.Title("Updated Title").
    Author("New Author").
    Language("en").
    Subject("Adventure", "Classics")

// Modify or inject content
editor.AddContent("text/extra.xhtml", contentBytes)
editor.UpdateContent("text/chapter-1.xhtml", updatedContentBytes)
editor.RemoveResource("text/obsolete.xhtml")

// Switch back to reader mode or save
readerBack, _ := editor.Reader()
editor.SaveAs("final.epub")
editor.Save(ioWriter)
editor.WriteBytes(&byteSlice)
```

---

## Supported media types

- Documents: `application/xhtml+xml`, `text/html`
- Navigation: `application/x-dtbncx+xml` (NCX), `application/xhtml+xml` (NAV)
- Styles: `text/css`
- Images: `image/jpeg`, `image/png`, `image/gif`, `image/webp`, `image/svg+xml`
- Fonts: `font/ttf`, `font/otf`, `font/woff`, `font/woff2`

---

## Examples

Runnable example programs live under [`./examples`](./examples):

- [`./examples/read`](./examples/read): Inspecting metadata, extracting covers, reading the spine order, and converting content documents to Markdown.
- [`./examples/edit`](./examples/edit): In-place editing of metadata, updating chapters, switching back to reader with `editor.Reader()`, and saving.
- [`./examples/write/basic`](./examples/write/basic): Minimal clean EPUB creation.
- [`./examples/write/with_cover`](./examples/write/with_cover): Adding PNG and JPEG cover art.
- [`./examples/write/multiple_chapters_and_assets`](./examples/write/multiple_chapters_and_assets): Multi-chapter books with images and SVG diagrams.
- [`./examples/write/nested_toc_and_guide`](./examples/write/nested_toc_and_guide): Hierarchical TOC and landmarks/guide references.
- [`./examples/write/markdown`](./examples/write/markdown): Compiling books from Markdown.
- [`./examples/write/advanced_multilingual_rtl`](./examples/write/advanced_multilingual_rtl): Multilingual publications with Right-To-Left (RTL) progression.
- [`./examples/write/reconstruct_standardebooks`](./examples/write/reconstruct_standardebooks): Reading and reconstructing production literature from Standard Ebooks.

---

## Documentation

- **Documentation Website**: [epub.raitucarp.name](https://epub.raitucarp.name)
- **GoDoc Reference**: [pkg.go.dev/github.com/raitucarp/epub](https://pkg.go.dev/github.com/raitucarp/epub)

---

## Support

If you find this project helpful or use it in your applications, consider supporting its maintenance and development:

[![Support on Ko-fi](https://img.shields.io/badge/Ko--fi-Support%20raitucarp-ff5e5b?style=for-the-badge&logo=ko-fi&logoColor=white)](https://ko-fi.com/raitucarp)

You can support the creator directly at **[ko-fi.com/raitucarp](https://ko-fi.com/raitucarp)**. Every coffee is greatly appreciated!

---

## Contributing

```sh
git clone https://github.com/raitucarp/epub.git
cd epub
go mod download
git config core.hooksPath .githooks
go test ./...
```

---

## License

[MIT](./LICENSE.md) © [Raitucarp](https://github.com/raitucarp)
