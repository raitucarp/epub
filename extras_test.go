package epub

import (
	"strings"
	"testing"

	"github.com/raitucarp/epub/pkg"
	"golang.org/x/net/html"
)

// ---- Version / UID -----------------------------------------------------------

func TestReader_Version(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {Version: "3.0"},
		},
	}}

	if got := r.Version(); got != "3.0" {
		t.Errorf("expected '3.0', got %q", got)
	}
}

func TestReader_UID(t *testing.T) {
	t.Run("resolves unique-identifier", func(t *testing.T) {
		r := &Reader{epub: &Epub{
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": {
					UniqueIdentifier: "pub-id",
					Metadata: pkg.Metadata{
						Identifiers: []pkg.DCIdentifier{
							{ID: "secondary", Value: "urn:other"},
							{ID: "pub-id", Value: "urn:isbn:1234567890"},
						},
					},
				},
			},
		}}

		if got := r.UID(); got != "urn:isbn:1234567890" {
			t.Errorf("expected %q, got %q", "urn:isbn:1234567890", got)
		}
	})

	t.Run("falls back to first identifier", func(t *testing.T) {
		r := &Reader{epub: &Epub{
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": {
					Metadata: pkg.Metadata{
						Identifiers: []pkg.DCIdentifier{{ID: "id-1", Value: "urn:first"}},
					},
				},
			},
		}}

		if got := r.UID(); got != "urn:first" {
			t.Errorf("expected %q, got %q", "urn:first", got)
		}
	})
}

// ---- Language / Identifier ---------------------------------------------------

func TestReader_Language(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata: map[string]any{"language": []string{"en", "fr"}},
	}}

	if got := r.Language(); len(got) != 2 || got[0] != "en" {
		t.Errorf("unexpected languages: %v", got)
	}
}

func TestReader_Language_Missing(t *testing.T) {
	r := &Reader{epub: &Epub{metadata: map[string]any{}}}

	if got := r.Language(); got != nil {
		t.Errorf("expected nil for missing language, got %v", got)
	}
}

func TestReader_Identifier(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata: map[string]any{"identifiers": []string{"urn:x", "urn:y"}},
	}}

	if got := r.Identifier(); len(got) != 2 || got[0] != "urn:x" {
		t.Errorf("unexpected identifiers: %v", got)
	}
}

func TestReader_Identifier_Missing(t *testing.T) {
	r := &Reader{epub: &Epub{metadata: map[string]any{}}}

	if got := r.Identifier(); got != nil {
		t.Errorf("expected nil for missing identifiers, got %v", got)
	}
}

// ---- Title / Author fallbacks -------------------------------------------------

func TestReader_Title_FromMetadata(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata: map[string]any{"title": []string{"My Title"}},
	}}

	if got := r.Title(); len(got) != 1 || got[0] != "My Title" {
		t.Errorf("unexpected title: %v", got)
	}
}

func TestReader_Title_FromGuide(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata:  map[string]any{},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {
				Guide: &pkg.Guide{
					References: []pkg.GuideReference{{Type: pkg.GuideRefTitlePage, Href: "titlepage.xhtml"}},
				},
			},
		},
		resources: []PublicationResource{
			{ID: "tp", Href: "titlepage.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte(`<html xmlns:epub="http://www.idpf.org/2007/ops"><body><div epub:type="title">Guide Title</div></body></html>`)},
		},
	}}

	if got := r.Title(); len(got) != 1 || got[0] != "Guide Title" {
		t.Errorf("unexpected title: %v", got)
	}
}

func TestReader_Title_FromResource(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata:  map[string]any{},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {},
		},
		resources: []PublicationResource{
			{ID: "titlepage", Href: "titlepage.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte(`<html xmlns:epub="http://www.idpf.org/2007/ops"><body><div epub:type="title">Resource Title</div></body></html>`)},
		},
	}}

	if got := r.Title(); len(got) != 1 || got[0] != "Resource Title" {
		t.Errorf("unexpected title: %v", got)
	}
}

func TestReader_Author_Unknown(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata:  map[string]any{},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {},
		},
	}}

	if got := r.Author(); len(got) != 1 || got[0] != "Unknown" {
		t.Errorf("expected Unknown author, got %v", got)
	}
}

func TestReader_Author_FromMetadata(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata: map[string]any{"creator": []string{"Jane Doe"}},
	}}

	if got := r.Author(); len(got) != 1 || got[0] != "Jane Doe" {
		t.Errorf("unexpected author: %v", got)
	}
}

// ---- Cover detection ---------------------------------------------------------

func coverReader(resources []PublicationResource, metadata map[string]any, guide *pkg.Guide) *Reader {
	return &Reader{epub: &Epub{
		metadata:  metadata,
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {Guide: guide},
		},
		resources: resources,
	}}
}

func TestReader_GetCoverInMetadata(t *testing.T) {
	pngBytes := testPNGBytes(t, 2, 2)
	r := &Reader{epub: &Epub{
		metadata: map[string]any{"meta": map[string]any{"cover": "cover-id"}},
		resources: []PublicationResource{
			{ID: "cover-id", Href: "cover.png", MIMEType: pkg.MediaTypePNG, Content: pngBytes},
		},
	}}

	if cover := r.getCoverInMetadata(); cover == nil {
		t.Error("expected cover from metadata, got nil")
	}
}

func TestReader_GetCoverInResources(t *testing.T) {
	pngBytes := testPNGBytes(t, 2, 2)
	r := &Reader{epub: &Epub{
		metadata: map[string]any{},
		resources: []PublicationResource{
			{ID: "cover-img", Href: "cover.png", MIMEType: pkg.MediaTypePNG, Properties: pkg.CoverImageProperty, Content: pngBytes},
		},
	}}

	if cover := r.getCoverInResources(); cover == nil {
		t.Error("expected cover from resources, got nil")
	}
}

func TestReader_GetCoverInSpine(t *testing.T) {
	pngBytes := testPNGBytes(t, 2, 2)
	r := &Reader{epub: &Epub{
		metadata:  map[string]any{},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {Spine: pkg.Spine{ItemRefs: []pkg.ItemRef{{IDRef: "cover-image"}}}},
		},
		resources: []PublicationResource{
			{ID: "cover-image", Href: "cover.png", MIMEType: pkg.MediaTypePNG, Content: pngBytes},
		},
	}}

	if cover := r.getCoverInSpine(); cover == nil {
		t.Error("expected cover from spine, got nil")
	}
}

func TestFindFirstImg(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(`<html><body><div><img src="a.png"/></div></body></html>`))
	if img := findFirstImg(doc); img == nil || img.Data != "img" {
		t.Errorf("expected to find img, got %v", img)
	}
}

func TestGetImageSrc(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(`<html><body><img src="images/cover.png"/></body></html>`))
	img := findFirstImg(doc)
	if got := getImageSrc(img); got != "images/cover.png" {
		t.Errorf("expected src %q, got %q", "images/cover.png", got)
	}
}

func TestReader_Cover(t *testing.T) {
	pngBytes := testPNGBytes(t, 2, 2)
	r := coverReader(
		[]PublicationResource{
			{ID: "cover-img", Href: "cover.png", MIMEType: pkg.MediaTypePNG, Properties: pkg.CoverImageProperty, Content: pngBytes},
		},
		map[string]any{},
		nil,
	)

	if cover := r.Cover(); cover == nil {
		t.Error("expected cover, got nil")
	}
}

func TestReader_Cover_NoCover(t *testing.T) {
	r := coverReader(nil, map[string]any{}, nil)
	if cover := r.Cover(); cover != nil {
		t.Errorf("expected nil cover, got %v", cover)
	}
}

func TestReader_CoverBytes_NilCover(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata:    map[string]any{},
		rendition:   "default",
		packagePubs: map[string]*pkg.Package{"default": {}},
	}}

	if _, err := r.CoverBytes(); err == nil {
		t.Fatal("expected error when no cover is defined, got nil")
	}
}

func TestReader_CoverBytes(t *testing.T) {
	pngBytes := testPNGBytes(t, 2, 2)
	r := &Reader{epub: &Epub{
		metadata:    map[string]any{},
		rendition:   "default",
		packagePubs: map[string]*pkg.Package{"default": {}},
		resources: []PublicationResource{
			{ID: "cover", Href: "cover.png", MIMEType: pkg.MediaTypePNG, Properties: pkg.CoverImageProperty, Content: pngBytes},
		},
	}}

	got, err := r.CoverBytes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("expected non-empty cover bytes")
	}
}

// ---- Description -------------------------------------------------------------

func TestReader_Description_Metadata(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata: map[string]any{"description": []string{"Primary Description"}},
	}}

	if desc := r.Description(); strings.Join(desc, ", ") != "Primary Description" {
		t.Errorf("expected 'Primary Description', got %q", desc)
	}
}

func TestReader_Description_OptionalMeta(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata: map[string]any{
			"meta": map[string]any{"description": []any{"Secondary Description"}},
		},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {},
		},
	}}

	if desc := r.Description(); strings.Join(desc, ", ") != "Secondary Description" {
		t.Errorf("expected 'Secondary Description', got %q", desc)
	}
}

func TestReader_Description_SummaryMeta(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata: map[string]any{
			"meta": map[string]any{"summary": []any{"<p>Summary text</p>"}},
		},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {},
		},
	}}

	desc := r.Description()
	if len(desc) != 1 || !strings.Contains(desc[0], "Summary text") {
		t.Errorf("expected summary description, got %v", desc)
	}
}

func TestExtractDescriptionFromEpubType(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(`<html><body><div epub:type="abstract">Abstract content</div></body></html>`))

	if got := extractDescriptionFromEpubType("abstract", doc); !strings.Contains(got, "Abstract content") {
		t.Errorf("expected abstract content, got %q", got)
	}
}

func TestReader_Description_Spine(t *testing.T) {
	p := &pkg.Package{
		Spine: pkg.Spine{ItemRefs: []pkg.ItemRef{{IDRef: "res1"}}},
	}

	r := &Reader{epub: &Epub{
		metadata:  map[string]any{},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": p,
		},
		resources: []PublicationResource{
			{
				ID:       "res1",
				Href:     "intro.xhtml",
				Content:  []byte(`<html xmlns:epub="http://www.idpf.org/2007/ops"><body><div epub:type="introduction">Spine Description</div></body></html>`),
				MIMEType: pkg.MediaTypeXHTML,
			},
		},
	}}

	if desc := r.Description(); strings.Join(desc, ", ") != "Spine Description" {
		t.Errorf("expected 'Spine Description', got %q", desc)
	}
}

func TestReader_Description_References(t *testing.T) {
	p := &pkg.Package{
		Guide: &pkg.Guide{
			References: []pkg.GuideReference{{Type: pkg.GuideRefPreface, Href: "preface.xhtml"}},
		},
	}

	r := &Reader{epub: &Epub{
		metadata:  map[string]any{},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": p,
		},
		resources: []PublicationResource{
			{
				ID:       "res1",
				Href:     "preface.xhtml",
				Content:  []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body><p>Preface Description</p></body></html>`),
				MIMEType: pkg.MediaTypeXHTML,
			},
		},
	}}

	if desc := r.Description(); strings.Join(desc, ", ") != "Preface Description" {
		t.Errorf("expected 'Preface Description', got %q", desc)
	}
}

func TestReader_Description_TOC(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata:  map[string]any{},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {},
		},
		resources: []PublicationResource{
			{
				ID: "toc", Properties: pkg.NavProperty,
				Href:     "toc.xhtml",
				Content:  []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><nav epub:type="toc"><h1>TOC</h1><ol><li><a href="chap1.xhtml">Chapter 1</a></li></ol></nav></html>`),
				MIMEType: pkg.MediaTypeXHTML,
			},
			{
				ID:       "chap1",
				Href:     "chap1.xhtml",
				Content:  []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body><p>TOC Description</p></body></html>`),
				MIMEType: pkg.MediaTypeXHTML,
			},
		},
	}}

	if desc := r.Description(); strings.Join(desc, ", ") != "TOC Description" {
		t.Errorf("expected 'TOC Description', got %q", desc)
	}
}

// ---- References --------------------------------------------------------------

func TestReader_References(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {
				Guide: &pkg.Guide{
					References: []pkg.GuideReference{
						{Type: pkg.GuideRefCover, Href: "cover.xhtml"},
						{Type: pkg.GuideRefText, Href: "text.xhtml#chapter1"},
					},
				},
			},
		},
		resources: []PublicationResource{
			{ID: "cover", Href: "cover.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte("<html><body><p>cover</p></body></html>")},
			{ID: "text", Href: "text.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte("<html><body><p>text</p></body></html>")},
		},
	}}

	refs := r.References()
	if len(refs) != 2 {
		t.Fatalf("expected 2 references, got %d", len(refs))
	}
	if refs[pkg.GuideRefCover] == nil {
		t.Errorf("expected cover reference, got nil")
	}
	if refs[pkg.GuideRefText] == nil {
		t.Errorf("expected text reference (fragment stripped), got nil")
	}
}

func TestReader_References_NoGuide(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition:   "default",
		packagePubs: map[string]*pkg.Package{"default": {}},
	}}

	if refs := r.References(); len(refs) != 0 {
		t.Errorf("expected empty references, got %v", refs)
	}
}

func TestReader_Author_Fallbacks(t *testing.T) {
	// Fallback 1: Guide with titlepage
	r1 := &Reader{epub: &Epub{
		metadata:  map[string]any{},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {
				Guide: &pkg.Guide{
					References: []pkg.GuideReference{
						{Type: pkg.GuideRefTitlePage, Href: "titlepage.xhtml"},
					},
				},
			},
		},
		resources: []PublicationResource{
			{
				ID:       "tp",
				Href:     "titlepage.xhtml",
				Content:  []byte(`<html xmlns:epub="http://www.idpf.org/2007/ops"><body><div epub:type="author">Guide Author</div></body></html>`),
				MIMEType: pkg.MediaTypeXHTML,
			},
		},
	}}
	if authors := r1.Author(); len(authors) != 1 || authors[0] != "Guide Author" {
		t.Errorf("expected 'Guide Author', got %v", authors)
	}

	// Fallback 2: Resource ID matching title
	r2 := &Reader{epub: &Epub{
		metadata:  map[string]any{},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {},
		},
		resources: []PublicationResource{
			{
				ID:       "title-page",
				Href:     "tp.xhtml",
				Content:  []byte(`<html xmlns:epub="http://www.idpf.org/2007/ops"><body><div epub:type="author">Pattern Author</div></body></html>`),
				MIMEType: pkg.MediaTypeXHTML,
			},
		},
	}}
	if authors := r2.Author(); len(authors) != 1 || authors[0] != "Pattern Author" {
		t.Errorf("expected 'Pattern Author', got %v", authors)
	}

	// Fallback 3: Unknown
	r3 := &Reader{epub: &Epub{
		metadata:  map[string]any{},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {},
		},
	}}
	if authors := r3.Author(); len(authors) != 1 || authors[0] != "Unknown" {
		t.Errorf("expected 'Unknown', got %v", authors)
	}
}

func TestReader_CoverBytes_NoCover(t *testing.T) {
	r := &Reader{epub: &Epub{
		metadata:    map[string]any{},
		rendition:   "default",
		packagePubs: map[string]*pkg.Package{"default": {}},
	}}
	if _, err := r.CoverBytes(); err == nil {
		t.Error("expected error from CoverBytes() when no cover exists")
	}
}

func TestReader_CoverFromTOC(t *testing.T) {
	pngData := testPNGBytes(t, 10, 10)
	r := &Reader{epub: &Epub{
		metadata:  map[string]any{},
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {},
		},
		resources: []PublicationResource{
			{
				ID:       "toc",
				Href:     "toc.xhtml",
				Content:  []byte(`<html xmlns:epub="http://www.idpf.org/2007/ops"><nav epub:type="toc"><ol><li><a href="cover-page.xhtml">Cover</a></li></ol></nav></html>`),
				MIMEType: pkg.MediaTypeXHTML,
				Properties: pkg.NavProperty,
			},
			{
				ID:       "cover-page",
				Href:     "cover-page.xhtml",
				Content:  []byte(`<html><body><img src="cover-art.png"/></body></html>`),
				MIMEType: pkg.MediaTypeXHTML,
			},
			{
				ID:       "cover-art",
				Href:     "cover-art.png",
				Content:  pngData,
				MIMEType: pkg.MediaTypePNG,
			},
		},
	}}

	cov := r.Cover()
	if cov == nil {
		t.Error("expected cover to be found from TOC")
	}

	bytes, err := r.CoverBytes()
	if err != nil {
		t.Fatalf("expected CoverBytes to succeed: %v", err)
	}
	if len(bytes) == 0 {
		t.Error("expected non-empty cover bytes")
	}
}

func TestExtractDescriptionFromSummaryMeta(t *testing.T) {
	meta := map[string]any{
		"meta": map[string]any{
			"summary": []any{"This is a summary."},
		},
	}
	res := extractDescriptionFromSummaryMeta(meta)
	if len(res) != 1 || res[0] != "This is a summary." {
		t.Errorf("unexpected summary result: %v", res)
	}

	if res := extractDescriptionFromSummaryMeta(nil); res != nil {
		t.Errorf("expected nil for nil metadata, got %v", res)
	}
	if res := extractDescriptionFromSummaryMeta(map[string]any{"meta": "not-a-map"}); res != nil {
		t.Errorf("expected nil for non-map meta, got %v", res)
	}
}
