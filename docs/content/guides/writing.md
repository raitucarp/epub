---
title: "Creating EPUB Publications"
description: "How to construct new compliant EPUB 3 publications using epub.Writer."
---

The `epub.Writer` type enables assembling standards-compliant EPUB 3 publications from scratch.

## Initialization

Initialize a new writer with a publication identifier:

```go
w := epub.New("urn:uuid:12345678-1234-5678-1234-567812345678")
```

## Adding Metadata

Set standard and optional Dublin Core elements:

```go
w.Title("Primary Book Title", "Subtitle or Alternative Title")
w.Author("Jane Doe", "Co-Author")
w.Languages("en", "fr")
w.Description("A summary of the publication.")
w.Publisher("My Publishing House")
w.Subject("topic-1", "Go", "Programming", "Technology")
w.Rights("Copyright © 2026 Jane Doe. All rights reserved.")
w.Date(time.Now())
w.Modified(time.Now())
```

## Adding Cover Images

You can set cover images from byte slices, image files, or Go `image.Image` objects:

```go
// From raw bytes
w.Cover(imageBytes)

// From image.Image
w.CoverPNG(myRGBAImage)
w.CoverJPG(myRGBAImage)

// From local file
w.CoverFile("cover.jpg")
```

## Adding Assets and Content Documents

You can add XHTML documents, images, and arbitrary assets:

```go
// Add XHTML chapter content
w.AddContent("chapter-1.xhtml", []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Chapter 1</title></head>
<body>
  <h1>Chapter 1: The Beginning</h1>
  <p>Hello, reader!</p>
</body>
</html>`))

// Add an image resource
w.AddImage("chart.png", chartBytes)

// Add a content file directly from disk
w.AddContentFile("chapter-2.xhtml")
```

## Creating the Table of Contents

The EPUB specification requires every publication to provide a navigation structure. Call `TableOfContents`:

```go
toc := epub.TOC{
    Title: "Table of Contents",
    Items: []epub.TOC{
        {Title: "Chapter 1", Href: "chapter-1.xhtml"},
        {
            Title: "Part 2", 
            Href: "chapter-2.xhtml",
            Items: []epub.TOC{
                {Title: "Section 2.1", Href: "chapter-2.xhtml#s1"},
            },
        },
    },
}

if err := w.TableOfContents("toc", toc); err != nil {
    log.Fatalf("failed to create TOC: %v", err)
}
```

## Saving Output

Save your publication either to a file on disk or as raw in-memory bytes:

```go
// Write to file on disk
if err := w.Write("output.epub"); err != nil {
    log.Fatal(err)
}

// Or generate in-memory bytes without touching the filesystem
data, err := w.WriteBytes()
```
