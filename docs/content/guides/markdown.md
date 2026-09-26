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

## Compiling a Manuscript Directory

If you organize your book as a folder of `.md` files, you can compile the entire directory in one call:

```
manuscript/
├── 01-intro.md
├── 02-architecture.md
└── 03-conclusion.md
```

```go
w := epub.New("urn:uuid:my-manuscript")
w.Title("Full Manuscript")

// Automatically imports all .md files in natural sort order,
// converts them to XHTML, and extracts heading hierarchy for the TOC!
if err := w.AddMarkdownDirectory("manuscript"); err != nil {
    log.Fatalf("failed to compile directory: %v", err)
}

w.Write("manuscript.epub")
```
