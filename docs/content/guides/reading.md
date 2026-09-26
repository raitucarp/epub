---
title: "Reading EPUB Publications"
description: "How to open, inspect, and extract contents and metadata using epub.Reader."
---

The `epub.Reader` type provides complete inspection and extraction capabilities for EPUB 2 and EPUB 3 publications.

## Opening a Publication

You can open a publication either directly from a file path on disk or from an in-memory byte slice:

```go
// From file path
reader, err := epub.OpenReader("path/to/book.epub")

// From in-memory byte slice
reader, err := epub.NewReader(rawBytes)
```

## Extracting Metadata

`Reader` normalizes metadata across both EPUB 2 and EPUB 3 formats:

```go
// Essential metadata
titles := reader.Title()         // []string
authors := reader.Author()       // []string (dc:creator)
languages := reader.Language()   // []string (dc:language)
identifiers := reader.Identifier()// []string
uid := reader.UID()              // Primary unique identifier
version := reader.Version()      // e.g. "3.0" or "2.0"

// Complete metadata block
metadataMap := reader.Metadata()
```

## Accessing the Cover Image

The cover image can be extracted as an `image.Image` or as raw bytes:

```go
// As decoded image.Image
if cover := reader.Cover(); cover != nil {
    bounds := (*cover).Bounds()
    fmt.Printf("Width: %d, Height: %d\n", bounds.Dx(), bounds.Dy())
}

// As raw encoded bytes (e.g. to save directly to disk or send over HTTP)
coverBytes, err := reader.CoverBytes()
if err == nil {
    os.WriteFile("cover.png", coverBytes, 0644)
}
```

## Navigating the Spine (Reading Order)

The spine defines the linear reading order of the publication:

```go
spine := reader.Spine()
for i, resource := range spine {
    fmt.Printf("Chapter %d: %s (ID: %s, MIME: %s)\n", i+1, resource.Href, resource.ID, resource.MIMEType)
}
```

## Reading Content Documents & Converting to Markdown

`Reader` includes a built-in HTML-to-Markdown engine:

```go
// Read document as parsed *html.Node
docNode := reader.ReadContentHTMLByHref("chapter-1.xhtml")

// Read document directly as clean Markdown with YAML frontmatter!
md := reader.ReadContentMarkdownByHref("chapter-1.xhtml")
fmt.Println(md)
```

## Table of Contents

The table of contents is unified into a normalized `epub.TOC` tree whether the source file uses EPUB 3 NAV or EPUB 2 NCX:

```go
toc, err := reader.TableOfContents()
if err == nil {
    for _, item := range toc.Items {
        fmt.Printf("%s -> %s\n", item.Title, item.Href)
        for _, subItem := range item.Items {
            fmt.Printf("  └─ %s -> %s\n", subItem.Title, subItem.Href)
        }
    }
}
```
