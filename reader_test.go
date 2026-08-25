package epub

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/raitucarp/epub/ocf"
	"github.com/raitucarp/epub/pkg"
)

func TestRenditionKey(t *testing.T) {
	tests := []struct {
		name     string
		rootFile ocf.RootFile
		index    int
		expected string
	}{
		{"first is always default", ocf.RootFile{Label: "main"}, 0, "default"},
		{"second with no attributes", ocf.RootFile{}, 1, "rendition-1"},
		{"second with layout", ocf.RootFile{Layout: "pre-paginated"}, 1, "pre-paginated"},
		{"second with multiple attributes", ocf.RootFile{Layout: "pre-paginated", Language: "en"}, 2, "pre-paginated_en"},
		{"second with all attributes", ocf.RootFile{Media: "media", Layout: "layout", Language: "en", AccessMode: "visual", Label: "Main"}, 1, "media_layout_en_visual_Main"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renditionKey(tt.rootFile, tt.index); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestReader_ListRenditions(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			packagePubs: map[string]*pkg.Package{
				"pre-paginated": {},
				"default":       {},
				"rendition-1":   {},
			},
		},
	}

	got := r.ListRenditions()
	if len(got) != 3 {
		t.Fatalf("expected 3 renditions, got %d", len(got))
	}
	if got[0] != "default" {
		t.Errorf("expected default rendition first, got %v", got)
	}
}

// writeTestEPUB builds a minimal valid EPUB file with the Writer, returning its
// path on disk and its raw bytes.
func writeTestEPUB(t *testing.T) (string, []byte) {
	t.Helper()

	w := New("urn:test:reader")
	w.Title("Reader Test")
	w.Author("Author")
	w.Languages("en")
	w.AddContent("chapter.xhtml", []byte(`<html><body><p>Chapter</p></body></html>`))

	toc := TOC{Title: "Contents", Items: []TOC{{Title: "Chapter", Href: "chapter.xhtml"}}}
	if err := w.TableOfContents("toc", toc); err != nil {
		t.Fatalf("TableOfContents error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "test.epub")
	if err := w.Write(path); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	return path, data
}

func TestReader_NewReader(t *testing.T) {
	t.Run("valid bytes", func(t *testing.T) {
		_, data := writeTestEPUB(t)

		r, err := NewReader(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := r.UID(); got != "urn:test:reader" {
			t.Errorf("expected uid %q, got %q", "urn:test:reader", got)
		}
	})

	t.Run("invalid bytes", func(t *testing.T) {
		if _, err := NewReader([]byte("not an epub")); err == nil {
			t.Error("expected error for invalid bytes, got nil")
		}
	})
}

func TestReader_OpenReader(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		path, _ := writeTestEPUB(t)

		r, err := OpenReader(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := r.Title(); len(got) != 1 || got[0] != "Reader Test" {
			t.Errorf("unexpected title: %v", got)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		if _, err := OpenReader(filepath.Join(t.TempDir(), "missing.epub")); err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})

	t.Run("invalid file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "invalid.epub")
		if err := os.WriteFile(path, []byte("garbage data"), 0o644); err != nil {
			t.Fatalf("failed to write invalid file: %v", err)
		}
		if _, err := OpenReader(path); err == nil {
			t.Error("expected error for invalid file, got nil")
		}
	})
}

func TestReader_MultipleRootFiles(t *testing.T) {
	z := ocf.NewOCFZipContainer()
	z.AddMimeType()

	mkPkg := func(id, href string) pkg.Package {
		return pkg.Package{
			Version:          "3.0",
			UniqueIdentifier: "id",
			Metadata: pkg.Metadata{
				Identifiers: []pkg.DCIdentifier{{ID: "id", Value: "urn:" + id}},
				Titles:      []pkg.DCTitle{{Value: id}},
				Languages:   []pkg.DCLanguage{{Value: "en"}},
			},
			Manifest: pkg.Manifest{Items: []pkg.Item{{ID: "doc", Href: href, MediaType: pkg.MediaTypeXHTML}}},
			Spine:    pkg.Spine{ItemRefs: []pkg.ItemRef{{IDRef: "doc"}}},
		}
	}

	if err := z.AddPackage("EPUB/package.opf", mkPkg("default", "a.xhtml")); err != nil {
		t.Fatalf("AddPackage: %v", err)
	}
	if err := z.AddPackage("EPUB/pre.opf", mkPkg("second", "b.xhtml")); err != nil {
		t.Fatalf("AddPackage: %v", err)
	}
	if err := z.AddContainerXML("EPUB/package.opf", "EPUB/pre.opf"); err != nil {
		t.Fatalf("AddContainerXML: %v", err)
	}
	z.AddFile("EPUB/a.xhtml", []byte(`<html><body><p>a</p></body></html>`))
	z.AddFile("EPUB/b.xhtml", []byte(`<html><body><p>b</p></body></html>`))

	path := filepath.Join(t.TempDir(), "multi.epub")
	if err := z.Write(path); err != nil {
		t.Fatalf("Write: %v", err)
	}

	r, err := OpenReader(path)
	if err != nil {
		t.Fatalf("OpenReader: %v", err)
	}

	renditions := r.ListRenditions()
	if len(renditions) != 2 {
		t.Fatalf("expected 2 renditions, got %v", renditions)
	}
	if renditions[0] != "default" {
		t.Errorf("expected default first, got %v", renditions)
	}

	if got := r.UID(); got != "urn:default" {
		t.Errorf("expected default uid, got %q", got)
	}

	r.SelectPackageRendition(renditions[1])
	if got := r.UID(); got != "urn:second" {
		t.Errorf("expected second rendition uid, got %q", got)
	}
}
