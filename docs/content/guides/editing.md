---
title: "Editing EPUB Publications"
description: "How to modify existing EPUBs in-place with concise fluent methods using epub.Editor."
---

One of the standout features of `raitucarp/epub` is the **`Editor`** workflow. 

Instead of rebuilding a complex EPUB from scratch when you only need to tweak metadata, inject an introductory chapter, or update an asset, you can transition directly from any `Reader` into an `Editor`.

## Basic Workflow: `reader.Edit()`

```go
package main

import (
    "log"
    "github.com/raitucarp/epub"
)

func main() {
    reader, err := epub.OpenReader("book.epub")
    if err != nil {
        log.Fatal(err)
    }

    // 1. Enter edit mode
    editor, err := reader.Edit()
    if err != nil {
        log.Fatal(err)
    }

    // 2. Chained metadata modifications
    editor.Title("New Title").
        Author("New Author").
        Description("Updated summary").
        Language("en")

    // 3. Save modified publication to disk
    if err := editor.SaveAs("book-updated.epub"); err != nil {
        log.Fatal(err)
    }
}
```

---

## Modifying Metadata

`Editor` provides methods that replace or append metadata entries cleanly:

| Method | Description |
| :--- | :--- |
| `Title("New Title")` | Replaces the title list. |
| `AddTitle("Subtitle")` | Appends an additional title. |
| `Author("Jane Doe")` | Replaces author/creator entries. |
| `AddAuthor("Co-Author")` | Appends a creator entry. |
| `Description("Summary")` | Replaces the publication description. |
| `Subject("Tag 1", "Tag 2")` | Replaces subjects/tags. |
| `Language("id", "en")` | Replaces language tags. |
| `Identifier("custom-id")` | Replaces identifiers. |
| `Date(time.Time)` | Updates publication date. |
| `Modified(time.Time)` | Sets `dcterms:modified` timestamp. |
| `Direction("rtl")` | Sets page progression direction. |

---

## Manipulating Content & Resources

You can update existing files, add new chapters, or remove items:

### Adding Content
```go
// Add a new chapter to manifest and spine
res, err := editor.AddContent("text/extra-chapter.xhtml", chapterHTMLBytes)

// Add an image
imgRes, err := editor.AddImage("images/new-photo.png", pngBytes)
```

### Updating Existing Documents
```go
// Update an existing chapter by ID or Href
err := editor.UpdateContent("text/chapter-1.xhtml", newChapterHTMLBytes)
```

### Removing Resources
```go
// Removes the resource from manifest, spine (if present), and the zip container
err := editor.RemoveResource("text/old-chapter.xhtml")
```

---

## Three Ways to Save Output

The `Editor` provides three versatile output methods:

### 1. `SaveAs(filename)`
Saves the EPUB directly to a file on disk:
```go
err := editor.SaveAs("final-book.epub")
```

### 2. `Save(io.Writer)`
Streams the EPUB into any standard `io.Writer` (such as `bytes.Buffer`, an HTTP response, or an S3 pipe):
```go
var buf bytes.Buffer
err := editor.Save(&buf)
```

### 3. `WriteBytes(&b)`
Populates a byte slice directly through a pointer:
```go
var b []byte
err := editor.WriteBytes(&b)
```

---

## Seamless Mode Switching: `editor.Reader()`

If you need to verify changes or read extracted data immediately in memory without saving to disk first, call `editor.Reader()`:

```go
// Mutate
editor.Title("Reviewed Title")

// Convert back to Reader mode
readerBack, err := editor.Reader()
if err != nil {
    log.Fatal(err)
}

// Inspect immediately!
fmt.Println(readerBack.Title()) // ["Reviewed Title"]
```
