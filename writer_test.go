package epub

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/pkg"
)

// ---- Modified / ensureModifiedDate -------------------------------------------

func TestWriter_Modified(t *testing.T) {
	w := New("test-pub-id")
	w.Modified(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC))

	metas := w.epub.SelectedPackage().Metadata.Meta
	if len(metas) != 1 {
		t.Fatalf("expected 1 meta element, got %d", len(metas))
	}
	if metas[0].Property != "dcterms:modified" {
		t.Errorf("expected property dcterms:modified, got %q", metas[0].Property)
	}
	if metas[0].Value != "2024-01-02T03:04:05Z" {
		t.Errorf("expected value 2024-01-02T03:04:05Z, got %q", metas[0].Value)
	}
}

func TestWriter_EnsureModifiedDate(t *testing.T) {
	t.Run("adds when missing", func(t *testing.T) {
		w := New("test-pub-id")
		w.ensureModifiedDate()

		var found bool
		for _, meta := range w.epub.SelectedPackage().Metadata.Meta {
			if meta.Property == "dcterms:modified" && meta.Refines == "" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected dcterms:modified to be added")
		}
	})

	t.Run("does not duplicate", func(t *testing.T) {
		w := New("test-pub-id")
		w.Modified(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC))
		w.ensureModifiedDate()

		var count int
		for _, meta := range w.epub.SelectedPackage().Metadata.Meta {
			if meta.Property == "dcterms:modified" && meta.Refines == "" {
				count++
			}
		}
		if count != 1 {
			t.Errorf("expected exactly one dcterms:modified, got %d", count)
		}
	})
}

// ---- guardCheck --------------------------------------------------------------

func TestWriter_GuardCheck(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*pkg.Package, *Epub)
		expectedErr string
	}{
		{
			name: "Missing Identifiers",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = nil
			},
			expectedErr: "Package should have identifiers.",
		},
		{
			name: "Missing Titles",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = nil
			},
			expectedErr: "Package should have titles.",
		},
		{
			name: "Missing Languages",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = nil
			},
			expectedErr: "Package should have languages.",
		},
		{
			name: "Missing Manifest Items",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = []pkg.DCLanguage{{Value: "en"}}
				p.Manifest.Items = nil
			},
			expectedErr: "No content insides.",
		},
		{
			name: "Missing Text Content",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = []pkg.DCLanguage{{Value: "en"}}
				p.Manifest.Items = []pkg.Item{
					{MediaType: "image/jpeg", Properties: pkg.CoverImageProperty},
				}
			},
			expectedErr: "No text content insides.",
		},
		{
			name: "Cover Image Not Required",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = []pkg.DCLanguage{{Value: "en"}}
				p.Manifest.Items = []pkg.Item{
					{MediaType: pkg.MediaTypeXHTML},
				}
				e.navigationCenterEXtended = &ncx.NCX{}
			},
			expectedErr: "",
		},
		{
			name: "Missing Table of Contents",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = []pkg.DCLanguage{{Value: "en"}}
				p.Manifest.Items = []pkg.Item{
					{MediaType: pkg.MediaTypeXHTML},
					{MediaType: "image/jpeg", Properties: pkg.CoverImageProperty},
				}
				e.navigationCenterEXtended = nil
			},
			expectedErr: "No table of contents.",
		},
		{
			name: "Valid Package",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = []pkg.DCLanguage{{Value: "en"}}
				p.Manifest.Items = []pkg.Item{
					{MediaType: pkg.MediaTypeXHTML},
					{MediaType: "image/jpeg", Properties: pkg.CoverImageProperty},
				}
				e.navigationCenterEXtended = &ncx.NCX{}
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{
				epub: &Epub{
					packagePubs: map[string]*pkg.Package{
						"content": {Metadata: pkg.Metadata{}, Manifest: pkg.Manifest{}},
					},
				},
			}
			tt.setup(w.epub.packagePubs["content"], w.epub)
			err := w.guardCheck()
			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("Expected error '%s', got nil", tt.expectedErr)
				} else if err.Error() != tt.expectedErr {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
				}
			} else if err != nil {
				t.Errorf("Expected no error, got '%s'", err.Error())
			}
		})
	}
}

// ---- Title -------------------------------------------------------------------

func TestWriter_Title(t *testing.T) {
	w := New("multi-title")
	w.Title("Main Title", "Subtitle", "Alternate")

	titles := w.epub.SelectedPackage().Metadata.Titles
	if len(titles) != 3 {
		t.Fatalf("expected 3 titles, got %d", len(titles))
	}
	if titles[0].Value != "Main Title" || titles[0].ID != "title" {
		t.Errorf("unexpected primary title: %+v", titles[0])
	}
	if titles[1].Value != "Subtitle" || titles[2].Value != "Alternate" {
		t.Errorf("unexpected secondary titles: %+v", titles[1:])
	}
}

func TestWriter_Title_Empty(t *testing.T) {
	w := New("empty-title")
	w.Title()

	if len(w.epub.SelectedPackage().Metadata.Titles) != 0 {
		t.Errorf("expected no titles added, got %+v", w.epub.SelectedPackage().Metadata.Titles)
	}
}

// ---- Metadata methods --------------------------------------------------------

func TestWriter_MetadataMethods(t *testing.T) {
	w := New("meta-id")

	w.Description("A description", "Another description")
	w.Creator("creator-id", "Jane Doe")
	w.Contributor("editor", "John Editor")
	w.Subject("subject-id", "Fiction", "History")
	w.Publisher("Indie Press")
	w.Rights("All rights reserved")
	w.LongDescription("A longer description")
	w.Date(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC))
	w.Meta(pkg.Meta{Property: "custom", Value: "value"})
	w.MetaContent(map[string]string{"meta-name": "meta-content"})
	w.MetaProperty("meta-prop-id", "prop", "prop-value")
	w.Refines("#title", "file-as", "Doe, Jane")
	w.Identifiers("urn:other", "urn:another")

	p := w.epub.SelectedPackage()
	if len(p.Metadata.OptionalDC) == 0 {
		t.Fatal("expected optional dc elements")
	}
	if len(p.Metadata.Identifiers) != 3 {
		t.Errorf("expected 3 identifiers, got %d", len(p.Metadata.Identifiers))
	}
	if len(p.Metadata.Meta) == 0 {
		t.Fatal("expected meta elements")
	}
}

func TestWriter_Refines(t *testing.T) {
	w := New("refines-id")
	w.Refines("#title", "file-as", "Doe, Jane", pkg.Meta{Scheme: "x"})

	metas := w.epub.SelectedPackage().Metadata.Meta
	if len(metas) != 1 {
		t.Fatalf("expected 1 meta, got %d", len(metas))
	}
	m := metas[0]
	if m.Refines != "#title" || m.Property != "file-as" || m.Value != "Doe, Jane" {
		t.Errorf("unexpected refine meta: %+v", m)
	}
}

func TestWriter_Languages(t *testing.T) {
	w := New("lang-id")
	w.Languages("en", "fr")

	langs := w.epub.SelectedPackage().Metadata.Languages
	if len(langs) != 2 {
		t.Fatalf("expected 2 languages, got %d", len(langs))
	}
	if langs[0].Value != "en" || langs[1].Value != "fr" {
		t.Errorf("unexpected languages: %+v", langs)
	}
}

func TestWriter_Languages_Empty(t *testing.T) {
	w := New("lang-empty")
	w.Languages()

	if len(w.epub.SelectedPackage().Metadata.Languages) != 0 {
		t.Errorf("expected no languages added, got %+v", w.epub.SelectedPackage().Metadata.Languages)
	}
}

// ---- Direction / dirs --------------------------------------------------------

func TestWriter_DirectionAndDirs(t *testing.T) {
	w := New("dir-id")
	w.Direction("rtl")
	w.SetContentDir("content")
	w.SetTextDir("txt")
	w.SetImageDir("img")

	if w.epub.SelectedPackage().Dir != "rtl" {
		t.Errorf("expected rtl direction, got %q", w.epub.SelectedPackage().Dir)
	}
	if w.contentDir != "content" || w.textDir != "txt" || w.imagesDir != "img" {
		t.Errorf("unexpected dirs: %q %q %q", w.contentDir, w.textDir, w.imagesDir)
	}
}

// ---- AddGuide ----------------------------------------------------------------

func TestWriter_AddGuide(t *testing.T) {
	w := New("guide-id")
	w.AddGuide(pkg.GuideRefCover, "cover.xhtml", "Cover")

	guide := w.epub.SelectedPackage().Guide
	if guide == nil {
		t.Fatal("expected guide to be created")
	}
	if len(guide.References) != 1 {
		t.Fatalf("expected 1 reference, got %d", len(guide.References))
	}
	if guide.References[0].Type != pkg.GuideRefCover || guide.References[0].Href != "cover.xhtml" {
		t.Errorf("unexpected reference: %+v", guide.References[0])
	}
}

// ---- Content / image / spine -------------------------------------------------

func TestWriter_AddContentAndSpine(t *testing.T) {
	w := New("content-id")
	res := w.AddContent("chapter.xhtml", []byte("<html><body><p>Chapter</p></body></html>"))

	if res.ID == "" || res.Href != "chapter.xhtml" || res.MIMEType != pkg.MediaTypeXHTML {
		t.Errorf("unexpected resource: %+v", res)
	}

	spine := w.epub.SelectedPackage().Spine.ItemRefs
	if len(spine) != 1 || spine[0].IDRef != res.ID {
		t.Errorf("expected spine to contain added content, got %+v", spine)
	}

	img := w.AddImage("img.png", testPNGBytes(t, 1, 1))
	if img.MIMEType != pkg.MediaTypePNG {
		t.Errorf("expected png media type, got %q", img.MIMEType)
	}
}

func TestWriter_UniqueIDPreventsCollision(t *testing.T) {
	w := New("uid-collision")

	w.AddContent("same.xhtml", []byte("<html><body>one</body></html>"))
	w.AddContent("same.xhtml", []byte("<html><body>two</body></html>"))

	items := w.epub.SelectedPackage().Manifest.Items
	seen := map[string]bool{}
	for _, item := range items {
		if seen[item.ID] {
			t.Errorf("duplicate manifest item id %q", item.ID)
		}
		seen[item.ID] = true
	}
}

func TestWriter_NavNameDoesNotOverwriteContent(t *testing.T) {
	w := New("nav-collision")
	w.Title("T")
	w.Languages("en")

	w.AddContent("toc.xhtml", []byte("<html><body><p>real toc page</p></body></html>"))

	toc := TOC{Items: []TOC{{Title: "Chapter", Href: "toc.xhtml"}}}
	if err := w.TableOfContents("toc", toc); err != nil {
		t.Fatalf("TableOfContents error: %v", err)
	}

	var foundContent, navFound bool
	for _, res := range w.epub.resources {
		if res.Href == "toc.xhtml" && res.Properties == pkg.NotProperty {
			foundContent = true
			if string(res.Content) != "<html><body><p>real toc page</p></body></html>" {
				t.Errorf("content document was overwritten by nav: %q", string(res.Content))
			}
		}
		if res.Properties == pkg.NavProperty {
			navFound = true
			if res.Href == "toc.xhtml" {
				t.Errorf("nav document collided with content doc href")
			}
		}
	}
	if !foundContent {
		t.Error("expected the toc.xhtml content document to still exist")
	}
	if !navFound {
		t.Error("expected a nav document to be generated")
	}
}

// ---- detectContentMediaType --------------------------------------------------

func TestDetectContentMediaType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  []byte
		want     string
	}{
		{"svg extension", "img.svg", nil, pkg.MediaTypeSVG},
		{"xhtml extension", "doc.xhtml", nil, pkg.MediaTypeXHTML},
		{"html extension", "doc.html", nil, pkg.MediaTypeXHTML},
		{"htm extension", "doc.htm", nil, pkg.MediaTypeXHTML},
		{"xml extension", "doc.xml", nil, pkg.MediaTypeXHTML},
		{"unknown extension defaults to xhtml", "doc.txt", []byte("<html></html>"), pkg.MediaTypeXHTML},
		{"uppercase extension", "DOC.XHTML", nil, pkg.MediaTypeXHTML},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectContentMediaType(tt.filename, tt.content); got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

// ---- Cover methods -----------------------------------------------------------

func TestWriter_CoverRawBytes(t *testing.T) {
	w := New("cover-raw")
	if err := w.Cover(testPNGBytes(t, 1, 1)); err != nil {
		t.Fatalf("Cover returned error: %v", err)
	}

	if len(w.epub.resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(w.epub.resources))
	}
	res := w.epub.resources[0]
	if res.Properties != pkg.CoverImageProperty {
		t.Errorf("expected cover-image property, got %q", res.Properties)
	}
	if res.MIMEType != pkg.MediaTypePNG {
		t.Errorf("expected png media type, got %q", res.MIMEType)
	}
}

func TestWriter_CoverPNGAndJPG(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 1, G: 2, B: 3, A: 255})

	t.Run("png", func(t *testing.T) {
		w := New("cover-png")
		if err := w.CoverPNG(img); err != nil {
			t.Fatalf("CoverPNG returned error: %v", err)
		}
		if len(w.epub.resources) != 1 {
			t.Fatalf("expected 1 resource, got %d", len(w.epub.resources))
		}
	})

	t.Run("jpg", func(t *testing.T) {
		w := New("cover-jpg")
		if err := w.CoverJPG(img); err != nil {
			t.Fatalf("CoverJPG returned error: %v", err)
		}
		if len(w.epub.resources) != 1 {
			t.Fatalf("expected 1 resource, got %d", len(w.epub.resources))
		}
	})
}

// ---- TableOfContents ---------------------------------------------------------

func TestWriter_TableOfContents_Nested(t *testing.T) {
	w := New("nested-toc")
	w.Title("Nested")
	w.Languages("en")

	toc := TOC{
		Title: "Contents",
		Items: []TOC{
			{Title: "Chapter 1", Href: "ch1.xhtml", Items: []TOC{
				{Title: "Section 1.1", Href: "ch1.xhtml#s1", Items: []TOC{
					{Title: "Sub 1.1.1", Href: "ch1.xhtml#s1.1"},
				}},
			}},
			{Title: "Chapter 2", Href: "ch2.xhtml"},
		},
	}

	if err := w.TableOfContents("toc", toc); err != nil {
		t.Fatalf("TableOfContents returned error: %v", err)
	}

	got := w.epub.navigationCenterEXtended
	if got == nil {
		t.Fatal("expected NCX to be set")
	}
	if got.DocTitle.Text != "Contents" {
		t.Errorf("expected doc title %q, got %q", "Contents", got.DocTitle.Text)
	}

	points := got.NavMap.NavPoints
	if len(points) != 2 {
		t.Fatalf("expected 2 top-level navpoints, got %d", len(points))
	}
	if points[0].NavLabel.Text != "Chapter 1" || points[0].Content.Src != "ch1.xhtml" {
		t.Errorf("unexpected first navpoint: %+v", points[0])
	}
	if points[0].PlayOrder != "1" || points[0].ID != "nav-point-1" {
		t.Errorf("unexpected playOrder/id: %s/%s", points[0].PlayOrder, points[0].ID)
	}
	if len(points[0].NavPoints) != 1 {
		t.Fatalf("expected chapter 1 to nest 1 section, got %d", len(points[0].NavPoints))
	}
	if section := points[0].NavPoints[0]; section.NavLabel.Text != "Section 1.1" {
		t.Errorf("unexpected section navpoint: %+v", section)
	}
	if len(points[0].NavPoints[0].NavPoints) != 1 {
		t.Errorf("expected nested sub-section, got %+v", points[0].NavPoints[0].NavPoints)
	}
	if points[1].NavLabel.Text != "Chapter 2" || points[1].PlayOrder != "4" {
		t.Errorf("unexpected second navpoint: %+v", points[1])
	}
}

func TestWriter_TableOfContents_Flat(t *testing.T) {
	w := New("flat-toc")
	w.Title("Flat")
	w.Languages("en")

	toc := TOC{
		Items: []TOC{
			{Title: "One", Href: "one.xhtml"},
			{Title: "Two", Href: "two.xhtml"},
		},
	}

	if err := w.TableOfContents("toc", toc); err != nil {
		t.Fatalf("TableOfContents returned error: %v", err)
	}

	points := w.epub.navigationCenterEXtended.NavMap.NavPoints
	if len(points) != 2 {
		t.Fatalf("expected 2 navpoints, got %d", len(points))
	}
	if points[0].PlayOrder != "1" || points[1].PlayOrder != "2" {
		t.Errorf("unexpected play order: %s, %s", points[0].PlayOrder, points[1].PlayOrder)
	}
}

// ---- Path traversal / file access --------------------------------------------

func TestWriter_PathTraversal(t *testing.T) {
	w := New("test-pub-id")

	tests := []struct {
		name string
		path string
	}{
		{"Absolute Path", "/etc/passwd"},
		{"Parent Traversal", "../../../etc/passwd"},
		{"Hidden Traversal", "foo/../../etc/passwd"},
		{"Backslash Traversal", "..\\..\\etc\\passwd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := w.AddContentFile(tt.path); err == nil {
				t.Errorf("AddContentFile(%q) expected error, got nil", tt.path)
			}
			if err := w.CoverFile(tt.path); err == nil {
				t.Errorf("CoverFile(%q) expected error, got nil", tt.path)
			}
			if res := w.AddImageFile(tt.path); res.ID != "" {
				t.Errorf("AddImageFile(%q) expected empty resource on invalid path, got %v", tt.path, res)
			}
		})
	}
}

// withWorkingDir runs fn with the process working directory temporarily changed
// to a new temporary directory, restoring the original directory afterwards.
func withWorkingDir(t *testing.T, fn func(dir string)) {
	t.Helper()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	defer os.Chdir(oldWd)
	fn(dir)
}

func TestWriter_AddContentFile(t *testing.T) {
	withWorkingDir(t, func(dir string) {
		if err := os.WriteFile(filepath.Join(dir, "chapter.xhtml"), []byte("<html><body>file</body></html>"), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}

		w := New("id")
		res, err := w.AddContentFile("chapter.xhtml")
		if err != nil {
			t.Fatalf("AddContentFile error: %v", err)
		}
		if res.ID != "chapter.xhtml" || res.MIMEType != pkg.MediaTypeXHTML {
			t.Errorf("unexpected resource: %+v", res)
		}
	})
}

func TestWriter_CoverFile(t *testing.T) {
	withWorkingDir(t, func(dir string) {
		if err := os.WriteFile(filepath.Join(dir, "cover.png"), testPNGBytes(t, 1, 1), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}

		w := New("id")
		if err := w.CoverFile("cover.png"); err != nil {
			t.Fatalf("CoverFile error: %v", err)
		}
		if len(w.epub.resources) != 1 {
			t.Fatalf("expected 1 resource, got %d", len(w.epub.resources))
		}
		if w.epub.resources[0].Properties != pkg.CoverImageProperty {
			t.Errorf("expected cover-image property, got %q", w.epub.resources[0].Properties)
		}
	})
}

func TestWriter_AddImageFile(t *testing.T) {
	withWorkingDir(t, func(dir string) {
		if err := os.WriteFile(filepath.Join(dir, "img.png"), testPNGBytes(t, 1, 1), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}

		w := New("id")
		res := w.AddImageFile("img.png")
		if res.ID != "img.png" || res.MIMEType != pkg.MediaTypePNG {
			t.Errorf("unexpected resource: %+v", res)
		}
	})
}

func TestWriter_SymlinkTraversal(t *testing.T) {
	w := New("test-pub-id")

	safeDir := t.TempDir()
	sensitiveFile := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(sensitiveFile, []byte("super secret"), 0o644); err != nil {
		t.Fatalf("Failed to create sensitive file: %v", err)
	}

	withWorkingDir(t, func(dir string) {
		if err := os.Symlink(sensitiveFile, "link_to_secret"); err != nil {
			t.Skipf("cannot create symlink (requires privilege on Windows): %v", err)
		}

		if _, err := w.AddContentFile("link_to_secret"); err == nil {
			t.Errorf("AddContentFile: expected error when reading via escaping symlink, got nil")
		}
		if err := w.CoverFile("link_to_secret"); err == nil {
			t.Errorf("CoverFile: expected error when reading via escaping symlink, got nil")
		}
		if res := w.AddImageFile("link_to_secret"); res.ID != "" || len(res.Content) > 0 {
			t.Errorf("AddImageFile: expected empty resource when reading via escaping symlink, got %v", res)
		}
	})

	_ = safeDir
}

// ---- Write -------------------------------------------------------------------

func TestWriter_WriteReadRoundtrip(t *testing.T) {
	w := New("urn:roundtrip:123")
	w.Title("Roundtrip Book")
	w.Author("Jane Doe")
	w.Languages("en")
	w.Cover(testPNGBytes(t, 2, 2))
	w.AddContent("chapter.xhtml", []byte(`<html><body><h1>Chapter</h1><p>Text</p></body></html>`))

	toc := TOC{
		Title: "Contents",
		Items: []TOC{
			{Title: "Chapter", Href: "chapter.xhtml", Items: []TOC{
				{Title: "Section", Href: "chapter.xhtml#s1"},
			}},
		},
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		t.Fatalf("TableOfContents error: %v", err)
	}

	out := filepath.Join(t.TempDir(), "roundtrip.epub")
	if err := w.Write(out); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	r, err := OpenReader(out)
	if err != nil {
		t.Fatalf("failed to open written epub: %v", err)
	}

	if got := r.Title(); len(got) != 1 || got[0] != "Roundtrip Book" {
		t.Errorf("unexpected title: %v", got)
	}
	if got := r.Author(); len(got) != 1 || got[0] != "Jane Doe" {
		t.Errorf("unexpected author: %v", got)
	}
	if got := r.UID(); got != "urn:roundtrip:123" {
		t.Errorf("unexpected uid: %q", got)
	}

	readTOC, err := r.TableOfContents()
	if err != nil {
		t.Fatalf("TableOfContents read error: %v", err)
	}
	if len(readTOC.Items) != 1 || readTOC.Items[0].Title != "Chapter" {
		t.Fatalf("unexpected TOC: %+v", readTOC.Items)
	}
	if len(readTOC.Items[0].Items) != 1 || readTOC.Items[0].Items[0].Title != "Section" {
		t.Errorf("expected nested section to roundtrip, got %+v", readTOC.Items[0].Items)
	}
}

func TestWriter_Write_InvalidPath(t *testing.T) {
	w := New("write-err")
	w.Title("T")
	w.Languages("en")
	w.AddContent("chapter.xhtml", []byte("<html><body>c</body></html>"))
	if err := w.TableOfContents("toc", TOC{Items: []TOC{{Title: "C", Href: "chapter.xhtml"}}}); err != nil {
		t.Fatalf("TableOfContents: %v", err)
	}

	out := filepath.Join(t.TempDir(), "missing-dir", "out.epub")
	if err := w.Write(out); err == nil {
		t.Error("expected error writing to non-existent directory, got nil")
	}
}

func TestWriter_WriteBytes(t *testing.T) {
	w := New("urn:writebytes:123")
	w.Title("Bytes Book")
	w.Author("Jane Doe")
	w.Languages("en")
	w.AddContent("chapter.xhtml", []byte(`<html><body><p>Chapter</p></body></html>`))
	if err := w.TableOfContents("toc", TOC{Items: []TOC{{Title: "Chapter", Href: "chapter.xhtml"}}}); err != nil {
		t.Fatalf("TableOfContents: %v", err)
	}

	data, err := w.WriteBytes()
	if err != nil {
		t.Fatalf("WriteBytes error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty bytes")
	}

	r, err := NewReader(data)
	if err != nil {
		t.Fatalf("failed to read back bytes: %v", err)
	}
	if got := r.Title(); len(got) != 1 || got[0] != "Bytes Book" {
		t.Errorf("unexpected title: %v", got)
	}
	if got := r.UID(); got != "urn:writebytes:123" {
		t.Errorf("unexpected uid: %q", got)
	}
}
