package epub

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createDummyPNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func TestParseMetadataFile(t *testing.T) {
	dir := t.TempDir()
	metaContent := `identifier: urn:uuid:test-directory-book
title:
  - Directory Masterpiece
  - Secondary Title
author:
  - Jane Doe
  - John Smith
language: en
description: A comprehensive test book for directory parsing.
publisher: Raitucarp Publishing
rights: CC-BY-SA 4.0
tags:
  - Golang
  - EPUB
direction: ltr
cover: cover-custom.png
toc_title: Custom Contents
`
	if err := os.WriteFile(filepath.Join(dir, "metadata.yml"), []byte(metaContent), 0o644); err != nil {
		t.Fatalf("failed to write metadata.yml: %v", err)
	}

	meta, err := parseMetadataFile(dir)
	if err != nil {
		t.Fatalf("parseMetadataFile error: %v", err)
	}
	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}

	if meta.GetIdentifier() != "urn:uuid:test-directory-book" {
		t.Errorf("expected identifier urn:uuid:test-directory-book, got %q", meta.GetIdentifier())
	}
	titles := meta.GetTitles()
	if len(titles) != 2 || titles[0] != "Directory Masterpiece" || titles[1] != "Secondary Title" {
		t.Errorf("unexpected titles: %v", titles)
	}
	authors := meta.GetAuthors()
	if len(authors) != 2 || authors[0] != "Jane Doe" || authors[1] != "John Smith" {
		t.Errorf("unexpected authors: %v", authors)
	}
	langs := meta.GetLanguages()
	if len(langs) != 1 || langs[0] != "en" {
		t.Errorf("unexpected languages: %v", langs)
	}
	if meta.GetDescription() != "A comprehensive test book for directory parsing." {
		t.Errorf("unexpected description: %q", meta.GetDescription())
	}
	pubs := meta.GetPublishers()
	if len(pubs) != 1 || pubs[0] != "Raitucarp Publishing" {
		t.Errorf("unexpected publishers: %v", pubs)
	}
	if meta.GetRights() != "CC-BY-SA 4.0" {
		t.Errorf("unexpected rights: %q", meta.GetRights())
	}
	subjects := meta.GetSubjects()
	if len(subjects) != 2 || subjects[0] != "Golang" || subjects[1] != "EPUB" {
		t.Errorf("unexpected subjects: %v", subjects)
	}
	if meta.GetDirection() != "ltr" {
		t.Errorf("unexpected direction: %q", meta.GetDirection())
	}
	if meta.GetCover() != "cover-custom.png" {
		t.Errorf("unexpected cover: %q", meta.GetCover())
	}
	if meta.GetTOCTitle() != "Custom Contents" {
		t.Errorf("unexpected toc title: %q", meta.GetTOCTitle())
	}
}

func TestWriter_AddDirectory(t *testing.T) {
	dir := t.TempDir()

	metaContent := `id: urn:uuid:xhtml-dir-book
title: XHTML Directory Book
creator: Jane Developer
lang: en
cover: cover.png
toc_title: Table of Contents
`
	if err := os.WriteFile(filepath.Join(dir, "metadata.yml"), []byte(metaContent), 0o644); err != nil {
		t.Fatalf("failed to write metadata.yml: %v", err)
	}

	// Cover image
	if err := os.WriteFile(filepath.Join(dir, "cover.png"), createDummyPNG(), 0o644); err != nil {
		t.Fatalf("failed to write cover: %v", err)
	}

	// Asset: CSS stylesheet
	if err := os.WriteFile(filepath.Join(dir, "style.css"), []byte("body { font-family: serif; }"), 0o644); err != nil {
		t.Fatalf("failed to write style.css: %v", err)
	}

	// Subdirectory with image asset
	imgSubDir := filepath.Join(dir, "images")
	if err := os.MkdirAll(imgSubDir, 0o755); err != nil {
		t.Fatalf("failed to create images dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(imgSubDir, "figure.png"), createDummyPNG(), 0o644); err != nil {
		t.Fatalf("failed to write figure.png: %v", err)
	}

	// XHTML chapters
	ch1 := `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Chapter 1: The Beginning</title><link rel="stylesheet" href="style.css"/></head>
<body>
  <h1>The Beginning</h1>
  <p>First paragraph.</p>
  <h2>Section 1: Details</h2>
  <p>Second paragraph with <img src="images/figure.png" alt="figure"/></p>
</body>
</html>`
	if err := os.WriteFile(filepath.Join(dir, "chapter1.xhtml"), []byte(ch1), 0o644); err != nil {
		t.Fatalf("failed to write chapter1.xhtml: %v", err)
	}

	ch2 := `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Chapter 2: The Journey</title></head>
<body>
  <h1>The Journey</h1>
  <p>Further along.</p>
</body>
</html>`
	if err := os.WriteFile(filepath.Join(dir, "chapter2.xhtml"), []byte(ch2), 0o644); err != nil {
		t.Fatalf("failed to write chapter2.xhtml: %v", err)
	}

	// Create writer without manual Go metadata declarations
	w := New("")
	if err := w.AddDirectory(dir); err != nil {
		t.Fatalf("AddDirectory error: %v", err)
	}

	// Check package metadata
	pkgDoc := w.epub.SelectedPackage()
	if pkgDoc == nil {
		t.Fatal("selected package is nil")
	}
	if len(pkgDoc.Metadata.Titles) == 0 || pkgDoc.Metadata.Titles[0].Value != "XHTML Directory Book" {
		t.Errorf("expected title 'XHTML Directory Book', got: %+v", pkgDoc.Metadata.Titles)
	}
	if len(pkgDoc.Metadata.Languages) == 0 || pkgDoc.Metadata.Languages[0].Value != "en" {
		t.Errorf("expected language 'en', got: %+v", pkgDoc.Metadata.Languages)
	}

	// Check spine
	if len(pkgDoc.Spine.ItemRefs) != 2 {
		t.Fatalf("expected 2 spine items, got %d", len(pkgDoc.Spine.ItemRefs))
	}
	if pkgDoc.Spine.ItemRefs[0].IDRef != "chapter1.xhtml" || pkgDoc.Spine.ItemRefs[1].IDRef != "chapter2.xhtml" {
		t.Errorf("unexpected spine order: %s, %s", pkgDoc.Spine.ItemRefs[0].IDRef, pkgDoc.Spine.ItemRefs[1].IDRef)
	}

	// Check TOC
	ncxNav := w.epub.navigationCenterEXtended
	if ncxNav == nil {
		t.Fatal("expected NCX table of contents")
	}
	if len(ncxNav.NavMap.NavPoints) != 2 {
		t.Fatalf("expected 2 top-level navpoints, got %d", len(ncxNav.NavMap.NavPoints))
	}
	if ncxNav.NavMap.NavPoints[0].NavLabel.Text != "The Beginning" {
		t.Errorf("expected 'The Beginning', got %q", ncxNav.NavMap.NavPoints[0].NavLabel.Text)
	}
	if len(ncxNav.NavMap.NavPoints[0].NavPoints) != 1 || ncxNav.NavMap.NavPoints[0].NavPoints[0].NavLabel.Text != "Section 1: Details" {
		t.Errorf("expected nested section navpoint, got: %+v", ncxNav.NavMap.NavPoints[0].NavPoints)
	}

	// Verify writing out to EPUB file and reading back
	outEPUB := filepath.Join(t.TempDir(), "test.epub")
	if err := w.Write(outEPUB); err != nil {
		t.Fatalf("w.Write error: %v", err)
	}

	r, err := OpenReader(outEPUB)
	if err != nil {
		t.Fatalf("OpenReader error: %v", err)
	}

	titles := r.Title()
	if len(titles) == 0 || titles[0] != "XHTML Directory Book" {
		t.Errorf("expected read title 'XHTML Directory Book', got %v", titles)
	}
	authors := r.Author()
	if len(authors) == 0 || authors[0] != "Jane Developer" {
		t.Errorf("expected read author 'Jane Developer', got %v", authors)
	}
	langs := r.Language()
	if len(langs) == 0 || langs[0] != "en" {
		t.Errorf("expected read language 'en', got %v", langs)
	}
	if len(r.Spine()) != 2 {
		t.Errorf("expected 2 spine items, got %d", len(r.Spine()))
	}
}

func TestWriter_AddMarkdownDirectory_WithMetadata(t *testing.T) {
	dir := t.TempDir()

	metaContent := `identifier: urn:uuid:markdown-dir-book
title: Markdown Directory Book
author: Mark Downer
language: en
toc_title: Table of Contents
`
	if err := os.WriteFile(filepath.Join(dir, "metadata.yml"), []byte(metaContent), 0o644); err != nil {
		t.Fatalf("failed to write metadata.yml: %v", err)
	}

	// Automatic cover
	if err := os.WriteFile(filepath.Join(dir, "cover.png"), createDummyPNG(), 0o644); err != nil {
		t.Fatalf("failed to write cover: %v", err)
	}

	// Asset CSS
	if err := os.WriteFile(filepath.Join(dir, "custom.css"), []byte("p { color: #333; }"), 0o644); err != nil {
		t.Fatalf("failed to write custom.css: %v", err)
	}

	// Markdown chapters
	ch1 := `# Introduction

Welcome to the markdown book.

## Sub-intro

More details here.
`
	if err := os.WriteFile(filepath.Join(dir, "01-intro.md"), []byte(ch1), 0o644); err != nil {
		t.Fatalf("failed to write 01-intro.md: %v", err)
	}

	ch2 := `# Chapter Two

Concluding thoughts.
`
	if err := os.WriteFile(filepath.Join(dir, "02-conclusion.md"), []byte(ch2), 0o644); err != nil {
		t.Fatalf("failed to write 02-conclusion.md: %v", err)
	}

	// Create writer without manual Go metadata declarations
	w := New("")
	if err := w.AddMarkdownDirectory(dir); err != nil {
		t.Fatalf("AddMarkdownDirectory error: %v", err)
	}

	outEPUB := filepath.Join(t.TempDir(), "md_test.epub")
	if err := w.Write(outEPUB); err != nil {
		t.Fatalf("w.Write error: %v", err)
	}

	r, err := OpenReader(outEPUB)
	if err != nil {
		t.Fatalf("OpenReader error: %v", err)
	}

	titles := r.Title()
	if len(titles) == 0 || titles[0] != "Markdown Directory Book" {
		t.Errorf("expected read title 'Markdown Directory Book', got %v", titles)
	}
	authors := r.Author()
	if len(authors) == 0 || authors[0] != "Mark Downer" {
		t.Errorf("expected read author 'Mark Downer', got %v", authors)
	}
	langs := r.Language()
	if len(langs) == 0 || langs[0] != "en" {
		t.Errorf("expected read language 'en', got %v", langs)
	}
	if len(r.Spine()) != 2 {
		t.Errorf("expected 2 spine items, got %d", len(r.Spine()))
	}
}

func TestEditor_AddDirectory(t *testing.T) {
	// First build an initial epub
	w := New("urn:uuid:initial")
	w.Title("Initial Title")
	w.Author("Initial Author")
	w.Languages("en")
	w.AddContent("ch1.xhtml", []byte("<h1>Ch1</h1><p>initial</p>"))
	w.TableOfContents("toc", TOC{
		Title: "TOC",
		Items: []TOC{{Title: "Ch1", Href: "ch1.xhtml#ch1"}},
	})

	initialEPUB := filepath.Join(t.TempDir(), "initial.epub")
	if err := w.Write(initialEPUB); err != nil {
		t.Fatalf("failed to write initial: %v", err)
	}

	r, err := OpenReader(initialEPUB)
	if err != nil {
		t.Fatalf("OpenReader error: %v", err)
	}

	ed, err := r.Edit()
	if err != nil {
		t.Fatalf("r.Edit() error: %v", err)
	}

	// Prepare new directory with metadata.yml and new content
	dir := t.TempDir()
	metaContent := `title: Updated Title via Directory
author: Updated Author
`
	if err := os.WriteFile(filepath.Join(dir, "metadata.yml"), []byte(metaContent), 0o644); err != nil {
		t.Fatalf("failed to write metadata.yml: %v", err)
	}

	newCh := `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Appended Chapter</title></head>
<body><h1>Appended</h1><p>new content</p></body>
</html>`
	if err := os.WriteFile(filepath.Join(dir, "appended.xhtml"), []byte(newCh), 0o644); err != nil {
		t.Fatalf("failed to write appended.xhtml: %v", err)
	}

	if err := ed.AddDirectory(dir); err != nil {
		t.Fatalf("Editor.AddDirectory error: %v", err)
	}

	editedEPUB := filepath.Join(t.TempDir(), "edited.epub")
	if err := ed.Write(editedEPUB); err != nil {
		t.Fatalf("ed.Write error: %v", err)
	}

	r2, err := OpenReader(editedEPUB)
	if err != nil {
		t.Fatalf("OpenReader r2 error: %v", err)
	}

	titles := r2.Title()
	if len(titles) == 0 || titles[0] != "Updated Title via Directory" {
		t.Errorf("expected updated title, got %v", titles)
	}
	authors := r2.Author()
	if len(authors) == 0 || authors[0] != "Updated Author" {
		t.Errorf("expected updated author, got %v", authors)
	}
	if len(r2.Spine()) < 2 {
		t.Errorf("expected at least 2 spine items, got %d", len(r2.Spine()))
	}
}

func TestEditor_AddMarkdownDirectory(t *testing.T) {
	w := New("urn:uuid:initial-md")
	w.Title("Initial MD Title")
	w.Author("Initial MD Author")
	w.Languages("en")
	w.AddContent("ch1.xhtml", []byte("<h1>Ch1</h1><p>initial</p>"))
	w.TableOfContents("toc", TOC{
		Title: "TOC",
		Items: []TOC{{Title: "Ch1", Href: "ch1.xhtml#ch1"}},
	})

	initialEPUB := filepath.Join(t.TempDir(), "initial_md.epub")
	if err := w.Write(initialEPUB); err != nil {
		t.Fatalf("failed to write initial: %v", err)
	}

	r, err := OpenReader(initialEPUB)
	if err != nil {
		t.Fatalf("OpenReader error: %v", err)
	}

	ed, err := r.Edit()
	if err != nil {
		t.Fatalf("r.Edit() error: %v", err)
	}

	dir := t.TempDir()
	metaContent := `title: Updated via Markdown Directory
`
	if err := os.WriteFile(filepath.Join(dir, "metadata.yml"), []byte(metaContent), 0o644); err != nil {
		t.Fatalf("failed to write metadata.yml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "extra.md"), []byte("# Extra Chapter\n\nExtra markdown."), 0o644); err != nil {
		t.Fatalf("failed to write extra.md: %v", err)
	}

	if err := ed.AddMarkdownDirectory(dir); err != nil {
		t.Fatalf("Editor.AddMarkdownDirectory error: %v", err)
	}

	editedEPUB := filepath.Join(t.TempDir(), "edited_md.epub")
	if err := ed.Write(editedEPUB); err != nil {
		t.Fatalf("ed.Write error: %v", err)
	}

	r2, err := OpenReader(editedEPUB)
	if err != nil {
		t.Fatalf("OpenReader r2 error: %v", err)
	}

	titles := r2.Title()
	if len(titles) == 0 || titles[0] != "Updated via Markdown Directory" {
		t.Errorf("expected updated title, got %v", titles)
	}
}

func TestToStringSlice(t *testing.T) {
	if res := toStringSlice(nil); res != nil {
		t.Errorf("expected nil for nil input, got %v", res)
	}
	if res := toStringSlice("  single string  "); len(res) != 1 || res[0] != "single string" {
		t.Errorf("expected [single string], got %v", res)
	}
	if res := toStringSlice("   "); res != nil {
		t.Errorf("expected nil for blank string, got %v", res)
	}
	if res := toStringSlice([]string{"a", "b"}); len(res) != 2 {
		t.Errorf("expected 2 elements for []string, got %v", res)
	}
	mixed := []any{"one", 42, "  trimmed  ", nil, ""}
	res := toStringSlice(mixed)
	if len(res) != 3 || res[0] != "one" || res[1] != "42" || res[2] != "trimmed" {
		t.Errorf("unexpected result for mixed slice: %v", res)
	}
	if res := toStringSlice(12345); len(res) != 1 || res[0] != "12345" {
		t.Errorf("expected [12345], got %v", res)
	}
}

func TestParseDate(t *testing.T) {
	if _, ok := parseDate(""); ok {
		t.Error("expected false for empty string")
	}
	layouts := []string{
		"2026-09-27T04:00:00Z",
		"2026-09-27T04:00:00.123456Z",
		"2026-09-27T04:00:00",
		"2026-09-27 04:00:00",
		"2026-09-27",
		"2026/09/27",
		"27 Sep 2026",
		"Sep 27, 2026",
	}
	for _, l := range layouts {
		if _, ok := parseDate(l); !ok {
			t.Errorf("expected parseDate to succeed for %q", l)
		}
	}
	if _, ok := parseDate("invalid-date-format"); ok {
		t.Error("expected false for invalid date")
	}
}

func TestDetectMIMEType(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"doc.xhtml", "application/xhtml+xml"},
		{"doc.html", "application/xhtml+xml"},
		{"doc.htm", "application/xhtml+xml"},
		{"style.css", "text/css"},
		{"image.png", "image/png"},
		{"image.jpg", "image/jpeg"},
		{"image.jpeg", "image/jpeg"},
		{"image.gif", "image/gif"},
		{"image.svg", "image/svg+xml"},
		{"image.webp", "image/webp"},
		{"font.otf", "font/otf"},
		{"font.ttf", "font/ttf"},
		{"font.woff", "font/woff"},
		{"font.woff2", "font/woff2"},
		{"nav.ncx", "application/x-dtbncx+xml"},
		{"data.bin", "application/octet-stream"},
	}
	for _, tt := range tests {
		got := detectMIMEType(tt.path)
		if got != tt.want {
			t.Errorf("detectMIMEType(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestProcessXHTMLDocument(t *testing.T) {
	// Case 1: Document with headings missing IDs and duplicate titles
	src := []byte(`<?xml version="1.0"?>
<html>
<head><title>My Book Title</title></head>
<body>
  <h1>Introduction</h1>
  <h1>Introduction</h1>
  <h2>Sub 1</h2>
</body>
</html>`)
	annotated, headings, title, err := processXHTMLDocument(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "My Book Title" {
		t.Errorf("expected title 'My Book Title', got %q", title)
	}
	if len(headings) != 3 {
		t.Fatalf("expected 3 headings, got %d", len(headings))
	}
	if headings[0].id != "introduction" || headings[1].id != "introduction-2" {
		t.Errorf("expected deduped IDs, got %s and %s", headings[0].id, headings[1].id)
	}
	if !strings.Contains(string(annotated), `id="introduction-2"`) {
		t.Errorf("expected annotated HTML to contain id, got %s", string(annotated))
	}

	// Case 2: Document without body
	_, _, _, err = processXHTMLDocument([]byte(`<head><title>No body</title></head>`))
	if err != nil {
		t.Errorf("expected no error for document without body: %v", err)
	}

	// Case 3: Empty headings
	_, headings, _, _ = processXHTMLDocument([]byte(`<body><h1></h1><h1>Valid</h1></body>`))
	if len(headings) != 1 || headings[0].text != "Valid" {
		t.Errorf("expected 1 valid heading, got %d", len(headings))
	}
}

func TestBuildHTMLTOC(t *testing.T) {
	docs := []htmlDocInfo{
		{
			href:  "intro.xhtml",
			title: "", // Fallback from filename
		},
		{
			href:  "ch1.xhtml",
			title: "Chapter 1",
			headings: []headingInfo{
				{level: 1, text: "Chapter 1", id: "ch1"},
				{level: 2, text: "Section 1.1", id: "sec1"},
				{level: 4, text: "Deep sub", id: "deep"}, // Promoted level
			},
		},
	}
	toc := buildHTMLTOC(docs)
	if len(toc.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(toc.Items))
	}
	if toc.Items[0].Title != "Intro" {
		t.Errorf("expected fallback title 'Intro', got %q", toc.Items[0].Title)
	}
}

func TestDirectoryMetadata_Fallbacks(t *testing.T) {
	meta := &DirectoryMetadata{
		ID:         "fallback-id",
		Title:      "Scalar Title",
		Creator:    "Scalar Creator",
		Lang:       "id",
		Desc:       "Short desc",
		Publisher:  "Single Publisher",
		Rights:     "Public Domain",
		Subject:    "General",
		PubDate:    "2026-01-01",
		Modified:   "2026-09-27T00:00:00Z",
		Dir:        "rtl",
		CoverImage: "art.jpg",
	}

	if meta.GetIdentifier() != "fallback-id" {
		t.Errorf("expected fallback-id, got %q", meta.GetIdentifier())
	}
	if len(meta.GetTitles()) != 1 || meta.GetTitles()[0] != "Scalar Title" {
		t.Errorf("unexpected titles: %v", meta.GetTitles())
	}
	if len(meta.GetAuthors()) != 1 || meta.GetAuthors()[0] != "Scalar Creator" {
		t.Errorf("unexpected authors: %v", meta.GetAuthors())
	}
	if len(meta.GetLanguages()) != 1 || meta.GetLanguages()[0] != "id" {
		t.Errorf("unexpected languages: %v", meta.GetLanguages())
	}
	if meta.GetDescription() != "Short desc" {
		t.Errorf("unexpected description: %q", meta.GetDescription())
	}
	if len(meta.GetPublishers()) != 1 || meta.GetPublishers()[0] != "Single Publisher" {
		t.Errorf("unexpected publishers: %v", meta.GetPublishers())
	}
	if meta.GetRights() != "Public Domain" {
		t.Errorf("unexpected rights: %q", meta.GetRights())
	}
	if len(meta.GetSubjects()) != 1 || meta.GetSubjects()[0] != "General" {
		t.Errorf("unexpected subjects: %v", meta.GetSubjects())
	}
	if meta.GetDirection() != "rtl" {
		t.Errorf("unexpected direction: %q", meta.GetDirection())
	}
	if meta.GetCover() != "art.jpg" {
		t.Errorf("unexpected cover: %q", meta.GetCover())
	}
	if _, ok := meta.GetDate(); !ok {
		t.Error("expected valid PubDate")
	}
	if _, ok := meta.GetModified(); !ok {
		t.Error("expected valid Modified")
	}
}

func TestAddDirectory_InvalidDir(t *testing.T) {
	w := New("")
	if err := w.AddDirectory(filepath.Join(t.TempDir(), "non-existent-dir")); err == nil {
		t.Error("expected error for non-existent directory")
	}

	ed := &Editor{}
	if err := ed.AddDirectory(filepath.Join(t.TempDir(), "non-existent-dir")); err == nil {
		t.Error("expected error for non-existent directory on Editor")
	}
}
