---
title: "Documentation"
description: "High-performance, idiomatic Go library for reading, creating, and editing EPUB publications."
---

## Why `raitucarp/epub`?

Building and parsing EPUB files in Go has historically required juggling complex XML specifications (OCF, OPF, NCX, NAV XHTML), zip container quirks, and inconsistent schema versions.

**`github.com/raitucarp/epub`** was designed from the ground up for developer happiness, reliability, and speed:

- **Complete Workflow**: Read, Write, and now seamlessly **Edit** existing EPUB publications.
- **Spec-Compliant**: Full EPUB 2 and EPUB 3 compatibility, tested against W3C conformance test suites.
- **Battle-Tested**: Verified on production literature from [Standard Ebooks](https://standardebooks.org/).
- **Rich Media**: First-class support for PNG/JPEG covers, SVG diagrams, and font obfuscation (IDPF & Adobe algorithms).
- **Markdown Native**: Seamlessly write chapters in Markdown; get clean XHTML and automatic heading-based TOCs.

---

## 30-Second Quick Taste

### 1. Read an Existing Book
```go
reader, err := epub.OpenReader("book.epub")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Title:", reader.Title())
fmt.Println("Author:", reader.Author())

// Convert chapter directly to clean Markdown!
md := reader.ReadContentMarkdownById("chapter-1")
fmt.Println(md)
```

### 2. Create a Book from Scratch
```go
w := epub.New("urn:uuid:my-book-id")
w.Title("My Adventure")
w.Author("Jane Doe")
w.Languages("en")

w.AddContent("ch1.xhtml", []byte(`<h1>Chapter 1</h1><p>The journey begins.</p>`))
w.TableOfContents("toc", epub.TOC{
    Title: "Contents",
    Items: []epub.TOC{{Title: "Chapter 1", Href: "ch1.xhtml"}},
})

w.Write("my-book.epub")
```

### 3. Edit Any Existing Book
```go
reader, _ := epub.OpenReader("my-book.epub")

// Enter edit mode
editor, _ := reader.Edit()

editor.Title("My Adventure (Second Edition)").
    Author("Jane Doe, PhD").
    Description("Revised edition with new notes.")

// Add a chapter and save directly
editor.AddContent("ch2.xhtml", []byte(`<h1>Chapter 2</h1><p>New adventures.</p>`))
editor.SaveAs("my-book-v2.epub")
```
