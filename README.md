# epub

`epub` is a Go library for reading and writing [EPUB](https://www.w3.org/TR/epub-33/)
publications. It implements the EPUB 3.3 specification, including the Open
Container Format (OCF), the package document, the EPUB navigation document, and
the legacy NCX format used by EPUB 2.

The library provides two main entry points:

- [`Reader`](https://pkg.go.dev/github.com/raitucarp/epub#Reader) for inspecting
  the metadata, resources, navigation, and content of an existing publication.
- [`Writer`](https://pkg.go.dev/github.com/raitucarp/epub#Writer) for building
  new publications, including from Markdown sources.

## Installation

```sh
go get github.com/raitucarp/epub
```

Requires Go 1.25 or later.

## Quick start

### Read an EPUB

```go
book, err := epub.OpenReader("book.epub")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Title:", strings.Join(book.Title(), ", "))
fmt.Println("Author:", strings.Join(book.Author(), ", "))
fmt.Println("Language:", strings.Join(book.Language(), ", "))
fmt.Println("Identifier:", book.UID())

for _, id := range book.ListContentDocumentIds() {
    fmt.Println(book.ReadContentMarkdownById(id))
}
```

### Write an EPUB

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

### Write an EPUB from Markdown

A single Markdown document can be converted directly:

```go
w := epub.New("urn:isbn:9780000000001")
w.Title("A Markdown Book")
w.Languages("en")

if _, err := w.AddMarkdown("chapter-1.md", []byte("# Chapter 1\n\nOnce upon a time...\n")); err != nil {
    log.Fatal(err)
}
```

Or an entire directory of Markdown files can be assembled, with the table of
contents derived from the headings in each file:

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

`AddMarkdownDirectory` reads every `.md` and `.markdown` file in the directory,
adds them to the spine in file-name order, and builds a nested table of
contents from the `h1` through `h6` headings found in each document.

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

### Multiple renditions

Some publications ship more than one package document, for example a reflowable
and a fixed-layout rendition of the same content.

```go
for _, rendition := range book.ListRenditions() {
    fmt.Println(rendition)
}

book.SelectPackageRendition("pre-paginated")
book.CurrentSelectedPackage()      // *pkg.Package
book.CurrentSelectedPackagePath()  // path to the active package document
```

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
w.LongDescription("A longer description")
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

### Markdown

```go
w.AddMarkdown("chapter-1.md", []byte("# Chapter 1\n\nText.\n"))
w.AddMarkdownFile("chapter-1.md")
w.AddMarkdownDirectory("manuscript")
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

### Writing to disk or memory

```go
w.Write("book.epub")        // write to a file
data, err := w.WriteBytes() // write to an in-memory byte slice
```

## Supported media types

- Documents: `application/xhtml+xml`, `text/html`
- Navigation: `application/x-dtbncx+xml` (NCX), `application/xhtml+xml` (NAV)
- Styles: `text/css`
- Images: `image/jpeg`, `image/png`, `image/gif`, `image/webp`, `image/svg+xml`
- Fonts: `font/ttf`, `font/otf`, `font/woff`, `font/woff2`

## Examples

Runnable programs live under [`./examples`](./examples):

- [`read`](./examples/read) reads an EPUB file and prints its metadata and content.
- [`write`](./examples/write) builds a minimal EPUB from XHTML.
- [`markdown`](./examples/markdown) builds an EPUB from a directory of Markdown files.

The package also includes [Godoc examples](https://pkg.go.dev/github.com/raitucarp/epub#pkg-examples)
that progress from a basic write/read round-trip to a Markdown-driven build and
a read-modify-write workflow.

## Documentation

The full API reference is available on
[pkg.go.dev](https://pkg.go.dev/github.com/raitucarp/epub).

## Contributing

```sh
git clone https://github.com/raitucarp/epub.git
cd epub
go mod download
git config core.hooksPath .githooks
go test ./...
```

## License

[MIT](./LICENSE.md)
