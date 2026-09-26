---
title: "Showcase: The White Company"
description: "Live web showcase of Arthur Conan Doyle's The White Company rendered from an EPUB using raitucarp/epub."
---

<div class="book-header">
  <img src="/images/white-company-cover.png" alt="The White Company Cover" class="book-cover-img" />
  <div>
    <h2 style="margin-top: 0; border: none; padding: 0;">The White Company</h2>
    <p><strong>Author:</strong> Arthur Conan Doyle</p>
    <p><strong>Publisher:</strong> Standard Ebooks</p>
    <p><strong>Source File:</strong> <code>tests/data/arthur-conan-doyle_the-white-company.epub</code></p>
    <p style="color: var(--text-muted); font-size: 0.95rem;">
      This showcase demonstrates reading an authentic, production-grade Standard Ebooks publication directly using <code>epub.OpenReader()</code>, extracting its metadata and cover art, and rendering its chapters as responsive web pages inside Hugo.
    </p>
  </div>
</div>

## Read Chapters Online

The following sample chapters were converted from the original EPUB XHTML content documents directly into Markdown and rendered natively by Hugo:

- [**I**]({{< relref "chapter-1.md" >}})
- [**II**]({{< relref "chapter-2.md" >}})
- [**III**]({{< relref "chapter-3.md" >}})
- [**IV**]({{< relref "chapter-4.md" >}})
- [**V**]({{< relref "chapter-5.md" >}})

---

## How It Works

This entire section was generated with just a few lines of Go:

```go
reader, err := epub.OpenReader("tests/data/arthur-conan-doyle_the-white-company.epub")
if err != nil {
    log.Fatal(err)
}

// 1. Extract cover image
coverBytes, _ := reader.CoverBytes()
os.WriteFile("static/images/cover.png", coverBytes, 0644)

// 2. Read chapter directly as clean Markdown!
md := reader.ReadContentMarkdownByHref("text/chapter-1.xhtml")
```
