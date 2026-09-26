---
title: "Markdown Compilation"
description: "How to compile Markdown documents and whole directories into EPUB publications."
---

`raitucarp/epub` provides built-in support for GitHub Flavored Markdown (GFM), including tables, strikethrough, task lists, and autolinks.

## Converting Single Markdown Documents

Use `w.AddMarkdown`:

```go
w := epub.New("urn:uuid:my-markdown-book")
w.Title("Markdown Guide")

mdContent := []byte(`# Chapter 1: Introduction

This is regular paragraph text.

## Features

- **Bold** and *Italics*
- Code blocks:
  ~~~go
  fmt.Println("Hello!")
  ~~~
- GFM Tables:

| Feature | Supported |
| :--- | :--- |
| Tables | Yes |
| Lists | Yes |
`)

res, err := w.AddMarkdown("chapter-1.md", mdContent)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Added chapter with Href: %s\n", res.Href) // "chapter-1.xhtml"
```

## Compiling a Manuscript Directory with metadata.yml

If you organize your book as a folder of `.md` files, you can compile the entire directory in one call without declaring metadata in Go:

```
manuscript/
├── metadata.yml
├── cover.png
├── style.css
├── images/
│   └── diagram.png
├── 01-intro.md
├── 02-architecture.md
└── 03-conclusion.md
```

### Example `metadata.yml`

```yaml
identifier: urn:uuid:my-manuscript-uuid
title:
  - Full Manuscript Title
  - Secondary Subtitle
author:
  - Jane Doe
  - John Smith
language: en
description: A complete book compiled directly from Markdown with automated TOC and metadata.
publisher: Open Publishing
rights: CC-BY-SA 4.0
tags:
  - Go
  - EPUB
direction: ltr
cover: cover.png
toc_title: Table of Contents
```

### Compiling in Go

```go
w := epub.New("")

// Ingests metadata.yml, cover image, stylesheets, assets, and all .md files.
// Builds a nested TOC from heading hierarchies automatically!
if err := w.AddMarkdownDirectory("manuscript"); err != nil {
    log.Fatalf("failed to compile directory: %v", err)
}

if err := w.Write("manuscript.epub"); err != nil {
    log.Fatal(err)
}
```

## Compiling XHTML / HTML Directories

For publications already formatted as XHTML or HTML with assets (CSS, fonts, images), use `AddDirectory`:

```
book_sources/
├── metadata.yml
├── cover.jpg
├── css/
│   └── style.css
├── fonts/
│   └── custom.woff2
├── 01-chapter.xhtml
└── 02-chapter.xhtml
```

```go
w := epub.New("")

// Ingests XHTML/HTML files, static assets, and applies metadata.yml
if err := w.AddDirectory("book_sources"); err != nil {
    log.Fatalf("failed to compile XHTML directory: %v", err)
}

w.Write("book.epub")
```

Both `AddMarkdownDirectory` and `AddDirectory` are also available on `epub.Editor`.

