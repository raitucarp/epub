package epub

import (
	"encoding/xml"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/raitucarp/epub/pkg"
)

func TestEditor_NilPackage(t *testing.T) {
	ed := &Editor{}
	now := time.Now()

	// Calling fluent metadata methods on nil/uninitialized package returns safely
	ed.Title()
	ed.Title("Title 1", "Title 2")
	ed.AddTitle("Title 3")
	ed.Author("Author 1", "Author 2")
	ed.AddAuthor("Author 3")
	ed.Creator("aut", "Creator 1")
	ed.Description("Desc 1")
	ed.AddDescription("Desc 2")
	ed.Publisher("Pub 1")
	ed.AddPublisher("Pub 2")
	ed.Contributor("edt", "Contributor 1")
	ed.Subject("Subj 1")
	ed.AddSubject("Subj 2")
	ed.Rights("Rights 1")
	ed.Date(now)
	ed.Modified(now)
	ed.Language("en")
	ed.AddLanguage("fr")
	ed.Identifier("id-1")
	ed.AddIdentifier("isbn", "123")
	ed.UniqueIdentifier("uid-1")
	ed.Version("3.0")
	ed.Direction("rtl")
	ed.DublinCore("custom", "val")
	ed.RemoveMetadata("custom")
	ed.Meta(pkg.Meta{})
	ed.SetMeta("prop", "val")
	ed.MetaContent(map[string]string{"k": "v"})
	ed.MetaProperty("id", "prop", "val")
	ed.Refines("id", "prop", "val")
	ed.RemoveMeta("prop")
	ed.AddFile("arbitrary.txt", []byte("data"))
	ed.AddResource("id", "href.txt", "text/plain", pkg.NotProperty, []byte("data"))
	ed.AddContent("page.xhtml", []byte("<html/>"))
	ed.AddMarkdown("page.md", []byte("# MD"))
	ed.AddSpineItem("page.xhtml")
	ed.RemoveSpineItem("page.xhtml")
	ed.Cover([]byte{1, 2, 3})
	ed.TableOfContents("toc", TOC{})

	if ed.Resources() != nil {
		t.Error("expected nil Resources() on empty editor")
	}
	if ed.Package() != nil {
		t.Error("expected nil Package() on empty editor")
	}
	if ed.CurrentPackage() != nil {
		t.Error("expected nil CurrentPackage() on empty editor")
	}
	if ed.CurrentPackagePath() != "" {
		t.Error("expected empty CurrentPackagePath() on empty editor")
	}
}

func TestEditor_AddDirectory_Comprehensive(t *testing.T) {
	// Build base epub
	w := New("urn:uuid:base")
	w.Title("Base")
	w.Author("Author")
	w.Languages("en")
	w.AddContent("base.xhtml", []byte("<h1>Base</h1>"))
	w.TableOfContents("toc", TOC{Title: "TOC", Items: []TOC{{Title: "Base", Href: "base.xhtml"}}})

	baseEpub := filepath.Join(t.TempDir(), "base.epub")
	if err := w.Write(baseEpub); err != nil {
		t.Fatalf("Write: %v", err)
	}

	r, err := OpenReader(baseEpub)
	if err != nil {
		t.Fatalf("OpenReader: %v", err)
	}
	ed, err := r.Edit()
	if err != nil {
		t.Fatalf("Edit(): %v", err)
	}

	// Prepare comprehensive directory
	dir := t.TempDir()
	metaYAML := `identifier: urn:uuid:comprehensive-dir
titles:
  - Comprehensive Book
  - Subtitle
authors:
  - Lead Author
  - Secondary Author
languages:
  - en
  - id
description: Full descriptive test
publishers:
  - Universal Press
rights: Open Access
subjects:
  - Computer Science
  - Software
direction: ltr
date: "2026-09-27T00:00:00Z"
modified: "2026-09-27T00:00:00Z"
cover: my_cover.png
toc_title: Full Contents
`
	if err := os.WriteFile(filepath.Join(dir, "metadata.yml"), []byte(metaYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "my_cover.png"), createDummyPNG(), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sec.xhtml"), []byte("<html><body><h1>Section</h1></body></html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := ed.AddDirectory(dir); err != nil {
		t.Fatalf("ed.AddDirectory: %v", err)
	}

	// Verify all metadata fields were applied to ed
	p := ed.CurrentPackage()
	if p == nil {
		t.Fatal("expected non-nil package")
	}
	if len(p.Metadata.Titles) < 2 {
		t.Errorf("expected 2 titles, got %d", len(p.Metadata.Titles))
	}
	if len(p.Metadata.Languages) < 2 {
		t.Errorf("expected 2 languages, got %d", len(p.Metadata.Languages))
	}
}

func TestDirectory_ApplyMetadata_Errors(t *testing.T) {
	w := New("")
	// Nil metadata is no-op
	if err := w.applyMetadata(nil, ""); err != nil {
		t.Errorf("expected nil error on nil metadata: %v", err)
	}

	ed := &Editor{}
	if err := ed.applyMetadata(nil, ""); err != nil {
		t.Errorf("expected nil error on nil metadata for editor: %v", err)
	}

	// Metadata with non-existent cover file returns error
	metaMissingCover := &DirectoryMetadata{Cover: "non-existent-cover.png"}
	if err := w.applyMetadata(metaMissingCover, t.TempDir()); err == nil {
		t.Error("expected error for missing cover file in w.applyMetadata")
	}
	if err := ed.applyMetadata(metaMissingCover, t.TempDir()); err == nil {
		t.Error("expected error for missing cover file in ed.applyMetadata")
	}
}

func TestReader_ListRenditions_DefaultReorder(t *testing.T) {
	r := &Reader{epub: &Epub{
		packagePubs: map[string]*pkg.Package{
			"alpha":   {},
			"default": {},
			"beta":    {},
		},
	}}

	renditions := r.ListRenditions()
	if len(renditions) != 3 {
		t.Fatalf("expected 3 renditions, got %d", len(renditions))
	}
	if renditions[0] != "default" {
		t.Errorf("expected 'default' to be ordered first, got %q", renditions[0])
	}
}

func TestContent_Refines_Comprehensive(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {
				Metadata: pkg.Metadata{
					Identifiers: []pkg.DCIdentifier{{ID: "id-1", XMLName: xml.Name{Local: "identifier"}, Value: "urn:1"}},
					Titles:      []pkg.DCTitle{{ID: "t-1", XMLName: xml.Name{Local: "title"}, Value: "Title 1"}},
					Languages:   []pkg.DCLanguage{{ID: "l-1", XMLName: xml.Name{Local: "language"}, Value: "en"}},
					OptionalDC:  []pkg.DCOptional{{ID: "c-1", XMLName: xml.Name{Local: "creator"}, Value: "Author 1"}},
					Meta: []pkg.Meta{
						{ID: "m-ref", Refines: "#c-1", Property: "role", Value: "aut"},
						{ID: "m-standalone", Property: "dcterms:modified", Value: "2026-09-27"},
					},
				},
			},
		},
	}}

	refines := r.Refines()
	if len(refines) == 0 {
		t.Fatal("expected non-empty refines map")
	}
	if refines["c-1"] == nil || refines["c-1"]["role"] == nil {
		t.Errorf("expected refinement for c-1 role, got %v", refines["c-1"])
	}
}

func TestContent_ReadMissing_And_SVG(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {
				Manifest: pkg.Manifest{
					Items: []pkg.Item{
						{ID: "diagram", Href: "diagram.svg", MediaType: pkg.MediaTypeSVG},
					},
				},
			},
		},
		resources: []PublicationResource{
			{ID: "diagram", Href: "diagram.svg", MIMEType: pkg.MediaTypeSVG, Content: []byte(`<svg xmlns="http://www.w3.org/2000/svg"><circle r="10"/></svg>`)},
		},
	}}

	// Missing ID lookups
	if doc := r.ReadContentHTMLById("no-such-id"); doc != nil {
		t.Error("expected nil for missing HTML ID")
	}
	if md := r.ReadContentMarkdownById("no-such-id"); md != "" {
		t.Errorf("expected empty string for missing Markdown ID, got %q", md)
	}

	// ContentDocumentSVG
	svgDocs := r.ContentDocumentSVG()
	if len(svgDocs) != 1 {
		t.Errorf("expected 1 SVG document, got %d", len(svgDocs))
	}
}

func TestWriter_UniqueID_Collisions(t *testing.T) {
	w := New("urn:uuid:test")
	// Requesting the same ID repeatedly appends -1, -2, etc.
	id1 := w.uniqueID("chapter")
	w.epub.SelectedPackage().Manifest.Items = append(w.epub.SelectedPackage().Manifest.Items, pkg.Item{ID: id1})

	id2 := w.uniqueID("chapter")
	w.epub.SelectedPackage().Manifest.Items = append(w.epub.SelectedPackage().Manifest.Items, pkg.Item{ID: id2})

	id3 := w.uniqueID("chapter")

	if id1 != "chapter" {
		t.Errorf("expected 'chapter', got %q", id1)
	}
	if id2 != "chapter-1" {
		t.Errorf("expected 'chapter-1', got %q", id2)
	}
	if id3 != "chapter-2" {
		t.Errorf("expected 'chapter-2', got %q", id3)
	}
}

func TestWriter_AddFiles_InvalidPaths(t *testing.T) {
	w := New("urn:uuid:test")
	// AddContentFile non-local path
	resContent, err := w.AddContentFile("../outside.xhtml")
	if err == nil && resContent.ID != "" {
		t.Error("expected error or empty resource for non-local AddContentFile")
	}

	// AddImageFile non-local path
	resImg := w.AddImageFile("../outside.png")
	if resImg.ID != "" {
		t.Error("expected empty resource for non-local AddImageFile")
	}
}

func TestEditor_CoverVariants(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}
	ed, err := r.Edit()
	if err != nil {
		t.Fatalf("Edit(): %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			img.Set(x, y, color.RGBA{R: 100, G: 100, B: 200, A: 255})
		}
	}

	if err := ed.CoverPNG(img); err != nil {
		t.Fatalf("CoverPNG: %v", err)
	}
	if err := ed.CoverJPG(img); err != nil {
		t.Fatalf("CoverJPG: %v", err)
	}
}

func TestDirectory_IgnoredDirs_And_MetadataVariants(t *testing.T) {
	dir := t.TempDir()
	// create hidden dir and node_modules
	_ = os.MkdirAll(filepath.Join(dir, ".hidden"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, ".hidden", "secret.xhtml"), []byte("<html></html>"), 0o644)
	_ = os.MkdirAll(filepath.Join(dir, "node_modules"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, "node_modules", "package.xhtml"), []byte("<html></html>"), 0o644)

	metaYAML := `id: urn:test:meta-variant
title: Variant Title
author: Variant Author
language: en
direction: rtl
date: 2026-09-27
modified: 2026-09-27
cover: cover.png
`
	_ = os.WriteFile(filepath.Join(dir, "metadata.yml"), []byte(metaYAML), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "cover.png"), createDummyPNG(), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "main.xhtml"), []byte("<html><body><h1>Main</h1></body></html>"), 0o644)

	w := New("existing-id")
	if err := w.AddDirectory(dir); err != nil {
		t.Fatalf("AddDirectory: %v", err)
	}

	ed := &Editor{}
	// Test update content on missing
	if err := ed.UpdateContent("missing", []byte("data")); err == nil {
		t.Error("expected error on UpdateContent for nil editor")
	}
}
