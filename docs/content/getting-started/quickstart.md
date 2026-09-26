---
title: "Quickstart Guide"
description: "Learn how to read, write, and edit EPUB publications in under 5 minutes."
---

This quickstart introduces the three fundamental building blocks of the library:
1. **`epub.Reader`**: For inspecting and extracting data from existing `.epub` files.
2. **`epub.Writer`**: For constructing compliant new `.epub` files from scratch.
3. **`epub.Editor`**: For modifying existing `.epub` publications in place with fluent chaining.

---

## 1. Reading an EPUB

Opening a file from disk takes one line of code:

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/raitucarp/epub"
)

func main() {
    reader, err := epub.OpenReader("sample.epub")
    if err != nil {
        log.Fatalf("failed to open: %v", err)
    }

    // Read metadata
    fmt.Printf("Title: %s\n", strings.Join(reader.Title(), ", "))
    fmt.Printf("Author: %s\n", strings.Join(reader.Author(), ", "))

    // Extract cover image
    if cover := reader.Cover(); cover != nil {
        bounds := (*cover).Bounds()
        fmt.Printf("Cover Dimensions: %dx%d\n", bounds.Dx(), bounds.Dy())
    }

    // Inspect Table of Contents
    toc, _ := reader.TableOfContents()
    for _, item := range toc.Items {
        fmt.Printf("- %s (%s)\n", item.Title, item.Href)
    }
}
```

---

## 2. Writing a New EPUB

To create a new publication, initialize a `Writer` with a publication identifier:

```go
package main

import (
    "log"
    "github.com/raitucarp/epub"
)

func main() {
    w := epub.New("urn:uuid:my-book-001")
    w.Title("Adventures in Go")
    w.Author("Jane Doe")
    w.Languages("en")

    // Add content documents
    w.AddContent("chapter-1.xhtml", []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Chapter 1</title></head>
<body>
  <h1>Chapter 1: The First Step</h1>
  <p>Welcome to the story!</p>
</body>
</html>`))

    // Set Table of Contents (Required by EPUB specification)
    w.TableOfContents("toc", epub.TOC{
        Title: "Table of Contents",
        Items: []epub.TOC{
            {Title: "Chapter 1: The First Step", Href: "chapter-1.xhtml"},
        },
    })

    // Save to disk
    if err := w.Write("output.epub"); err != nil {
        log.Fatalf("failed to write epub: %v", err)
    }
}
```

---

## 3. Editing an Existing EPUB

When you need to update an existing book without recreating it from scratch, use `reader.Edit()`:

```go
package main

import (
    "log"
    "github.com/raitucarp/epub"
)

func main() {
    reader, err := epub.OpenReader("output.epub")
    if err != nil {
        log.Fatal(err)
    }

    // 1. Enter edit mode
    editor, err := reader.Edit()
    if err != nil {
        log.Fatal(err)
    }

    // 2. Mutate metadata with clean fluent chaining
    editor.Title("Adventures in Go (Second Edition)").
        Author("Jane Doe, PhD").
        Description("Revised edition with new chapters.")

    // 3. Add or update content
    editor.AddContent("chapter-2.xhtml", []byte(`<h1>Chapter 2</h1><p>Continuing the journey...</p>`))

    // 4. Save your changes
    if err := editor.SaveAs("output-v2.epub"); err != nil {
        log.Fatal(err)
    }
}
```
