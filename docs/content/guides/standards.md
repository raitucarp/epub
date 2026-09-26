---
title: "Working with Standard Ebooks & Advanced EPUBs"
description: "How to handle production-quality literature, font obfuscation, and re-packaging recipes."
---

Publications from [Standard Ebooks](https://standardebooks.org/) are widely regarded as the gold standard of digital book design and typography.

`raitucarp/epub` is designed to handle the full complexity of these production books, including multiple package renditions, font obfuscation, deep navigation trees, and complex metadata refinements.

## Recipe: Reconstructing a Standard Ebooks Publication

Here is how you can read a Standard Ebooks publication and reconstruct it using `epub.Writer`:

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "strings"
    
    "github.com/raitucarp/epub"
)

func main() {
    // 1. Open original publication
    reader, err := epub.OpenReader("arthur-conan-doyle_the-white-company.epub")
    if err != nil {
        log.Fatal(err)
    }

    // 2. Initialize a new Writer with the original UID
    w := epub.New(reader.UID())
    w.Title(reader.Title()...)
    w.Author(reader.Author()...)
    w.Languages(reader.Language()...)
    w.Description("Reconstructed edition")

    // 3. Transfer the cover
    if coverBytes, err := reader.CoverBytes(); err == nil {
        _ = w.Cover(coverBytes)
    }

    // 4. Transfer image resources
    for id, imgBytes := range reader.ImageResources() {
        if !strings.Contains(strings.ToLower(id), "cover") {
            w.AddImage(id+".png", imgBytes)
        }
    }

    // 5. Transfer content documents from spine
    var tocItems []epub.TOC
    for i, res := range reader.Spine() {
        filename := filepath.Base(res.Href)
        w.AddContent(filename, res.Content)
        
        tocItems = append(tocItems, epub.TOC{
            Title: fmt.Sprintf("Section %d", i+1),
            Href:  filename,
        })
    }

    // 6. Generate Table of Contents
    w.TableOfContents("toc", epub.TOC{
        Title: "Table of Contents",
        Items: tocItems,
    })

    // 7. Write to disk
    if err := w.Write("reconstructed.epub"); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Reconstruction completed successfully!")
}
```

## Font Obfuscation Support

EPUB publications often obfuscate embedded fonts to comply with licensing agreements. The library automatically detects obfuscation in `META-INF/encryption.xml` and de-obfuscates font data using both:
- **IDPF Algorithm**: `http://www.idpf.org/2008/embedding`
- **Adobe Algorithm**: `http://ns.adobe.com/pdf/enc#RC`
