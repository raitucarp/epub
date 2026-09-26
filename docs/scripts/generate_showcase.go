package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/raitucarp/epub"
)

func main() {
	epubPath := filepath.Join("tests", "data", "arthur-conan-doyle_the-white-company.epub")
	if _, err := os.Stat(epubPath); os.IsNotExist(err) {
		epubPath = filepath.Join("..", "..", "tests", "data", "arthur-conan-doyle_the-white-company.epub")
	}

	reader, err := epub.OpenReader(epubPath)
	if err != nil {
		log.Fatalf("Failed to open EPUB: %v", err)
	}

	// 1. Export cover image to static/images/
	staticImgDir := filepath.Join("docs", "static", "images")
	_ = os.MkdirAll(staticImgDir, 0o755)
	if coverBytes, err := reader.CoverBytes(); err == nil && len(coverBytes) > 0 {
		coverPath := filepath.Join(staticImgDir, "white-company-cover.png")
		_ = os.WriteFile(coverPath, coverBytes, 0o644)
		fmt.Printf("Exported cover to %s\n", coverPath)
	}

	showcaseDir := filepath.Join("docs", "content", "showcase")
	_ = os.MkdirAll(showcaseDir, 0o755)

	// 2. Read chapters and generate markdown pages
	spine := reader.Spine()
	var chapterList []string

	chapterCount := 0
	for _, res := range spine {
		if !strings.Contains(res.Href, "chapter-") {
			continue
		}
		chapterCount++
		// Limit to first 5 chapters for a crisp showcase sample
		if chapterCount > 5 {
			break
		}

		chNum := strings.TrimSuffix(strings.TrimPrefix(filepath.Base(res.Href), "chapter-"), ".xhtml")
		md := reader.ReadContentMarkdownByHref(res.Href)

		// Extract title from markdown if present
		title := fmt.Sprintf("Chapter %s", chNum)
		lines := strings.Split(md, "\n")
		contentBody := md
		if strings.HasPrefix(md, "---") {
			// strip original frontmatter
			parts := strings.SplitN(md, "---", 3)
			if len(parts) >= 3 {
				contentBody = strings.TrimSpace(parts[2])
			}
		}

		// Find first heading
		for _, l := range lines {
			lTrim := strings.TrimSpace(l)
			if strings.HasPrefix(lTrim, "## ") || strings.HasPrefix(lTrim, "# ") {
				title = strings.TrimLeft(lTrim, "# ")
				break
			}
		}

		chFileName := fmt.Sprintf("chapter-%s.md", chNum)
		chFilePath := filepath.Join(showcaseDir, chFileName)

		hugoPage := fmt.Sprintf(`---
title: %q
description: "The White Company by Arthur Conan Doyle - %s"
weight: %d
---

<div style="margin-bottom: 1.5rem; font-size: 0.9rem; color: var(--text-muted);">
  <a href="/showcase/">← Back to Book Overview</a>
</div>

%s

<div class="book-nav-buttons">
  <a href="/showcase/" class="btn-nav">Table of Contents</a>
</div>
`, title, title, chapterCount*10, contentBody)

		_ = os.WriteFile(chFilePath, []byte(hugoPage), 0o644)
		chapterList = append(chapterList, fmt.Sprintf("- [**%s**]({{< relref %q >}})", title, chFileName))
		fmt.Printf("Generated showcase page: %s (%s)\n", chFileName, title)
	}

	// 3. Generate showcase _index.md
	showcaseIndex := fmt.Sprintf(`---
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

%s

---

## How It Works

This entire section was generated with just a few lines of Go:

`+"```go"+`
reader, err := epub.OpenReader("tests/data/arthur-conan-doyle_the-white-company.epub")
if err != nil {
    log.Fatal(err)
}

// 1. Extract cover image
coverBytes, _ := reader.CoverBytes()
os.WriteFile("static/images/cover.png", coverBytes, 0644)

// 2. Read chapter directly as clean Markdown!
md := reader.ReadContentMarkdownByHref("text/chapter-1.xhtml")
`+"```"+`
`, strings.Join(chapterList, "\n"))

	_ = os.WriteFile(filepath.Join(showcaseDir, "_index.md"), []byte(showcaseIndex), 0o644)
	fmt.Println("Generated docs/content/showcase/_index.md successfully!")
}
