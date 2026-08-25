package epub

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/ocf"
	"github.com/raitucarp/epub/pkg"
	"golang.org/x/net/html"
)

// xmlNameCreator is the XML name of the dc:creator element.
var xmlNameCreator = xml.Name{Space: pkg.NamespaceDC, Local: "creator"}

// ---- SelectPackageRendition / CurrentSelectedPackage -------------------------

func newMultiRenditionReader() *Reader {
	z := ocf.NewOCFZipContainer()
	z.AddFile("mimetype", []byte("application/epub+zip"))
	z.AddFile("EPUB/a.xhtml", []byte(`<html><body><p>default content</p></body></html>`))
	z.AddFile("EPUB/b.xhtml", []byte(`<html><body><p>pre content</p></body></html>`))

	return &Reader{epub: &Epub{
		rendition: "default",
		packagePaths: map[string]string{
			"default":       "EPUB/package.opf",
			"pre-paginated": "EPUB/pre.opf",
		},
		packagePubs: map[string]*pkg.Package{
			"default": {
				Version: "3.0",
				Metadata: pkg.Metadata{
					Identifiers: []pkg.DCIdentifier{{ID: "pub-id", Value: "urn:default"}},
					Titles:      []pkg.DCTitle{{Value: "Default Rendition"}},
					Languages:   []pkg.DCLanguage{{Value: "en"}},
				},
				Manifest: pkg.Manifest{Items: []pkg.Item{{ID: "a", Href: "a.xhtml", MediaType: pkg.MediaTypeXHTML}}},
				Spine:    pkg.Spine{ItemRefs: []pkg.ItemRef{{IDRef: "a"}}},
			},
			"pre-paginated": {
				Version: "3.0",
				Metadata: pkg.Metadata{
					Identifiers: []pkg.DCIdentifier{{ID: "pub-id", Value: "urn:pre"}},
					Titles:      []pkg.DCTitle{{Value: "Pre-paginated Rendition"}},
					Languages:   []pkg.DCLanguage{{Value: "fr"}},
				},
				Manifest: pkg.Manifest{Items: []pkg.Item{{ID: "b", Href: "b.xhtml", MediaType: pkg.MediaTypeXHTML}}},
				Spine:    pkg.Spine{ItemRefs: []pkg.ItemRef{{IDRef: "b"}}},
			},
		},
		zipContainer: z,
	}}
}

func TestReader_SelectPackageRendition(t *testing.T) {
	r := newMultiRenditionReader()

	if got := r.CurrentSelectedPackagePath(); got != "EPUB/package.opf" {
		t.Fatalf("expected default package path, got %q", got)
	}
	if got := r.UID(); got != "urn:default" {
		t.Fatalf("expected default uid, got %q", got)
	}

	r.SelectPackageRendition("pre-paginated")

	if got := r.CurrentSelectedPackagePath(); got != "EPUB/pre.opf" {
		t.Errorf("expected pre-paginated package path, got %q", got)
	}
	if got := r.UID(); got != "urn:pre" {
		t.Errorf("expected pre uid, got %q", got)
	}
	if got := r.Title(); len(got) != 1 || got[0] != "Pre-paginated Rendition" {
		t.Errorf("unexpected title after switching: %v", got)
	}
	if ids := r.ListContentDocumentIds(); len(ids) != 1 || ids[0] != "b" {
		t.Errorf("expected content id 'b', got %v", ids)
	}
	if got := r.CurrentSelectedPackage(); got == nil || got.Version != "3.0" {
		t.Errorf("unexpected current package: %+v", got)
	}
}

func TestReader_CurrentSelectedPackagePath(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition:    "default",
		packagePaths: map[string]string{"default": "EPUB/package.opf"},
	}}

	if got := r.CurrentSelectedPackagePath(); got != "EPUB/package.opf" {
		t.Errorf("expected %q, got %q", "EPUB/package.opf", got)
	}
}

// ---- parseHTML / fixSelfClosingTags / isXHTMLContent -------------------------

func TestIsXHTMLContent(t *testing.T) {
	if !isXHTMLContent(pkg.MediaTypeXHTML) {
		t.Error("application/xhtml+xml should be XHTML content")
	}
	if !isXHTMLContent(pkg.MediaTypeHTML) {
		t.Error("text/html should be XHTML content")
	}
	if isXHTMLContent(pkg.MediaTypeSVG) {
		t.Error("image/svg+xml should not be XHTML content")
	}
}

func TestFixSelfClosingTags(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "expands non-void self-closing script",
			input: `<html><head><script type="text/javascript" src="x.js"/></head><body><nav epub:type="toc"><h1>T</h1></nav></body></html>`,
			want:  `<script type="text/javascript" src="x.js"></script>`,
		},
		{
			name:  "preserves void elements",
			input: `<img src="a.png"/><br/><meta charset="utf-8"/>`,
			want:  `<img src="a.png"/><br/><meta charset="utf-8"/>`,
		},
		{
			name:  "expands div",
			input: `<div class="x"/>`,
			want:  `<div class="x"></div>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(fixSelfClosingTags([]byte(tt.input)))
			if !strings.Contains(got, tt.want) {
				t.Errorf("expected output to contain %q, got %q", tt.want, got)
			}
		})
	}
}

func TestFixSelfClosingTags_ParsesNav(t *testing.T) {
	doc := `<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops"><head><script src="k.js"/></head><body><nav epub:type="toc"><ol><li><a href="c1.xhtml">Chapter 1</a></li></ol></nav></body></html>`

	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "nav", Href: "nav.xhtml", MIMEType: pkg.MediaTypeXHTML, Properties: pkg.NavProperty, Content: []byte(doc)},
		},
	}}

	toc, err := r.TableOfContents()
	if err != nil {
		t.Fatalf("unexpected error parsing toc: %v", err)
	}
	if len(toc.Items) != 1 || toc.Items[0].Title != "Chapter 1" {
		t.Errorf("unexpected toc: %+v", toc.Items)
	}
}

// ---- ListContentDocumentIds / ListImageIds ----------------------------------

func TestReader_ListContentDocumentIds(t *testing.T) {
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "doc1", Href: "a.xhtml", MIMEType: pkg.MediaTypeXHTML},
			{ID: "doc2", Href: "b.html", MIMEType: pkg.MediaTypeHTML},
			{ID: "img", Href: "c.png", MIMEType: pkg.MediaTypePNG},
		},
	}}

	ids := r.ListContentDocumentIds()
	if len(ids) != 2 {
		t.Fatalf("expected 2 content ids, got %d: %v", len(ids), ids)
	}
}

func TestReader_ListImageIds(t *testing.T) {
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "img1", Href: "a.png", MIMEType: pkg.MediaTypePNG},
			{ID: "img2", Href: "b.jpg", MIMEType: pkg.MediaTypeJPEG},
			{ID: "img3", Href: "c.svg", MIMEType: pkg.MediaTypeSVG},
			{ID: "doc", Href: "d.xhtml", MIMEType: pkg.MediaTypeXHTML},
			{ID: "css", Href: "e.css", MIMEType: pkg.MediaTypeCSS},
		},
	}}

	ids := r.ListImageIds()
	if len(ids) != 3 {
		t.Fatalf("expected 3 image ids, got %d: %v", len(ids), ids)
	}
}

// ---- ContentDocumentXHTML / ContentDocumentXHTMLString ----------------------

func TestReader_ContentDocumentXHTML(t *testing.T) {
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "doc", Href: "a.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte(`<html><body><p>Hi</p></body></html>`)},
			{ID: "img", Href: "a.png", MIMEType: pkg.MediaTypePNG, Content: []byte("not html")},
		},
	}}

	docs := r.ContentDocumentXHTML()
	if len(docs) != 1 {
		t.Fatalf("expected 1 xhtml document, got %d", len(docs))
	}
	if _, ok := docs["doc"]; !ok {
		t.Errorf("expected doc key, got %v", docs)
	}
}

func TestReader_ContentDocumentXHTMLString(t *testing.T) {
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "doc", Href: "doc.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte(`<html><body><h1>Title</h1><p>Body</p></body></html>`)},
		},
	}}

	docs := r.ContentDocumentXHTMLString()
	if len(docs) != 1 {
		t.Fatalf("expected 1 document, got %d", len(docs))
	}
	if !strings.Contains(docs["doc"], "Body") {
		t.Errorf("expected rendered html to contain 'Body', got %q", docs["doc"])
	}
}

// ---- cleanupHTML / extractTitle / getTextByEpubType --------------------------

func TestCleanupHTML(t *testing.T) {
	doc := `<html><head><title>Remove Me</title></head><body><a href="x.xhtml">Keep link</a><a>no href</a></body></html>`
	node, err := html.Parse(strings.NewReader(doc))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	cleanupHTML(node)

	if title := extractTitle(node); title != "" {
		t.Errorf("expected title to be removed, got %q", title)
	}

	var noHrefFound bool
	for desc := range node.Descendants() {
		if desc.Type == html.ElementNode && desc.Data == "a" {
			hasHref := false
			for _, a := range desc.Attr {
				if a.Key == "href" {
					hasHref = true
				}
			}
			if !hasHref {
				t.Error("expected a without href to become div")
			}
			noHrefFound = noHrefFound || !hasHref
		}
	}
	_ = noHrefFound
}

func TestExtractTitle(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(`<html><head><title>My Title</title></head><body></body></html>`))
	if got := extractTitle(node); got != "My Title" {
		t.Errorf("expected title %q, got %q", "My Title", got)
	}
}

func TestGetTextByEpubType(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(`<html><body><div epub:type="title">The Title</div></body></html>`))
	if got := getTextByEpubType(node, "title"); got != "The Title" {
		t.Errorf("expected %q, got %q", "The Title", got)
	}

	if got := getTextByEpubType(node, "nonexistent"); got != "" {
		t.Errorf("expected empty for no match, got %q", got)
	}
}

// ---- ContentDocumentMarkdown / nodeToMarkdown --------------------------------

func TestReader_ContentDocumentMarkdown(t *testing.T) {
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "doc", Href: "doc.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte(`<html><head><title>Doc Title</title></head><body><h1>Heading</h1><p>Paragraph text</p></body></html>`)},
		},
	}}

	docs := r.ContentDocumentMarkdown()
	if len(docs) != 1 {
		t.Fatalf("expected 1 markdown document, got %d", len(docs))
	}
	md := docs["doc"]
	if !strings.Contains(md, "Paragraph text") {
		t.Errorf("expected markdown to contain paragraph text, got %q", md)
	}
	if !strings.Contains(md, "Doc Title") {
		t.Errorf("expected markdown front matter to contain title, got %q", md)
	}
}

// ---- ReadContentHTMLById / ReadContentHTMLByHref -----------------------------

func TestReader_ReadContentHTMLById(t *testing.T) {
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "valid", Href: "valid.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte("<html><body><p>Valid</p></body></html>")},
			{ID: "html", Href: "valid.html", MIMEType: pkg.MediaTypeHTML, Content: []byte("<html><body><p>HTML</p></body></html>")},
			{ID: "image", Href: "image.jpg", MIMEType: pkg.MediaTypeJPEG, Content: []byte("fake-jpeg")},
		},
	}}

	t.Run("xhtml resource", func(t *testing.T) {
		doc := r.ReadContentHTMLById("valid")
		if doc == nil || doc.Type != html.DocumentNode {
			t.Fatalf("expected parsed DocumentNode, got %v", doc)
		}
	})

	t.Run("text/html resource", func(t *testing.T) {
		if doc := r.ReadContentHTMLById("html"); doc == nil {
			t.Fatalf("expected text/html content to be read, got nil")
		}
	})

	t.Run("non-existent id", func(t *testing.T) {
		if doc := r.ReadContentHTMLById("missing"); doc != nil {
			t.Fatalf("expected nil, got %v", doc)
		}
	})

	t.Run("mismatched mime type", func(t *testing.T) {
		if doc := r.ReadContentHTMLById("image"); doc != nil {
			t.Fatalf("expected nil for non-content mime type, got %v", doc)
		}
	})
}

func TestReader_ReadContentHTMLByHref(t *testing.T) {
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "doc", Href: "text/chapter.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte("<html><body><p>Hello</p></body></html>")},
		},
	}}

	if doc := r.ReadContentHTMLByHref("text/chapter.xhtml"); doc == nil {
		t.Fatal("expected document, got nil")
	}
	if doc := r.ReadContentHTMLByHref("text/missing.xhtml"); doc != nil {
		t.Errorf("expected nil for missing href, got %v", doc)
	}
}

// ---- ReadContentMarkdownById / ReadContentMarkdownByHref ---------------------

func TestReader_ReadContentMarkdownById(t *testing.T) {
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "doc", Href: "doc.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte(`<html><body><p>Markdown body</p></body></html>`)},
		},
	}}

	if md := r.ReadContentMarkdownById("doc"); !strings.Contains(md, "Markdown body") {
		t.Errorf("expected markdown content, got %q", md)
	}
	if md := r.ReadContentMarkdownById("missing"); md != "" {
		t.Errorf("expected empty for missing id, got %q", md)
	}
}

func TestReader_ReadContentMarkdownByHref(t *testing.T) {
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "doc", Href: "text/doc.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte(`<html><body><p>By href</p></body></html>`)},
		},
	}}

	if md := r.ReadContentMarkdownByHref("text/doc.xhtml"); !strings.Contains(md, "By href") {
		t.Errorf("expected markdown via href, got %q", md)
	}
	if md := r.ReadContentMarkdownByHref("missing.xhtml"); md != "" {
		t.Errorf("expected empty markdown for missing href, got %q", md)
	}
}

// ---- ReadImageById / ReadImageByHref / ReadImageBytes* -----------------------

func TestReader_ReadImageById(t *testing.T) {
	pngBytes := testPNGBytes(t, 3, 2)
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "img", Href: "img.png", MIMEType: pkg.MediaTypePNG, Content: pngBytes},
		},
	}}

	img := r.ReadImageById("img")
	if img == nil {
		t.Fatal("expected image, got nil")
	}
	if (*img).Bounds().Dx() != 3 || (*img).Bounds().Dy() != 2 {
		t.Errorf("unexpected bounds: %v", (*img).Bounds())
	}
	if img := r.ReadImageById("missing"); img != nil {
		t.Errorf("expected nil for missing id, got %v", img)
	}
}

func TestReader_ReadImageByHref(t *testing.T) {
	pngBytes := testPNGBytes(t, 1, 1)
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "img", Href: "images/covers/cover.png", MIMEType: pkg.MediaTypePNG, Content: pngBytes},
			{ID: "doc", Href: "text/doc.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte("<html><body></body></html>")},
		},
	}}

	if img := r.ReadImageByHref("images/covers/cover.png"); img == nil {
		t.Fatal("expected to resolve image via nested href, got nil")
	}
	if img := r.ReadImageByHref("images/covers/missing.png"); img != nil {
		t.Errorf("expected nil for missing href, got %v", img)
	}
	if img := r.ReadImageByHref("text/doc.xhtml"); img != nil {
		t.Errorf("expected nil for non-image href, got %v", img)
	}
}

func TestReader_ReadImageBytes(t *testing.T) {
	pngBytes := testPNGBytes(t, 2, 2)
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "img", Href: "images/a.png", MIMEType: pkg.MediaTypePNG, Content: pngBytes},
		},
	}}

	if got := r.ReadImageBytesById("img"); !bytes.Equal(got, pngBytes) {
		t.Error("ReadImageBytesById returned unexpected bytes")
	}
	if got := r.ReadImageBytesByHref("images/a.png"); !bytes.Equal(got, pngBytes) {
		t.Error("ReadImageBytesByHref returned unexpected bytes")
	}
	if got := r.ReadImageBytesById("missing"); got != nil {
		t.Error("expected nil for missing id")
	}
	if got := r.ReadImageBytesByHref("images/missing.png"); got != nil {
		t.Error("expected nil for missing href")
	}
}

// ---- ContentDocumentSVG / Images / ImageResources ----------------------------

func TestReader_ContentDocumentSVG(t *testing.T) {
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "svg1", Href: "a.svg", MIMEType: pkg.MediaTypeSVG, Content: []byte(`<svg xmlns="http://www.w3.org/2000/svg"><circle r="1"/></svg>`)},
			{ID: "doc", Href: "b.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte("<html><body></body></html>")},
		},
	}}

	docs := r.ContentDocumentSVG()
	if len(docs) != 1 {
		t.Fatalf("expected 1 svg document, got %d", len(docs))
	}
	if _, ok := docs["svg1"]; !ok {
		t.Errorf("expected svg1 in map, got %v", docs)
	}
}

func TestReader_Images(t *testing.T) {
	pngBytes := testPNGBytes(t, 2, 2)
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "img", Href: "a.png", MIMEType: pkg.MediaTypePNG, Content: pngBytes},
			{ID: "bad", Href: "b.png", MIMEType: pkg.MediaTypePNG, Content: []byte("not-an-image")},
		},
	}}

	images := r.Images()
	if len(images) != 1 {
		t.Fatalf("expected 1 valid image, got %d", len(images))
	}
	if _, ok := images["img"]; !ok {
		t.Errorf("expected img key, got %v", images)
	}
}

func TestReader_ImageResources(t *testing.T) {
	pngBytes := testPNGBytes(t, 2, 2)
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "img", Href: "a.png", MIMEType: pkg.MediaTypePNG, Content: pngBytes},
		},
	}}

	got := r.ImageResources()
	if len(got) != 1 {
		t.Fatalf("expected 1 image resource, got %d", len(got))
	}
	if !bytes.Equal(got["img"], pngBytes) {
		t.Errorf("expected image bytes to match original")
	}
}

// ---- Spine -------------------------------------------------------------------

func TestReader_Spine(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {
				Spine: pkg.Spine{ItemRefs: []pkg.ItemRef{
					{IDRef: "third"},
					{IDRef: "first"},
					{IDRef: "second"},
					{IDRef: "missing"},
				}},
			},
		},
		resources: []PublicationResource{
			{ID: "first", Href: "first.xhtml"},
			{ID: "second", Href: "second.xhtml"},
			{ID: "third", Href: "third.xhtml"},
		},
	}}

	spine := r.Spine()
	if len(spine) != 3 {
		t.Fatalf("expected 3 spine items (missing idref skipped), got %d", len(spine))
	}
	if spine[0].ID != "third" || spine[1].ID != "first" || spine[2].ID != "second" {
		t.Errorf("expected spine to preserve idref order, got %+v", spine)
	}
}

// ---- parseMetadata / Metadata / Refines --------------------------------------

func TestReader_Metadata(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {
				Metadata: pkg.Metadata{
					Identifiers: []pkg.DCIdentifier{{Value: "  urn:x  "}},
					Titles:      []pkg.DCTitle{{Value: " Title "}},
					Languages:   []pkg.DCLanguage{{Value: "en"}},
					OptionalDC:  []pkg.DCOptional{{XMLName: xmlNameCreator, Value: "  Jane  "}},
				},
			},
		},
	}}
	r.parseMetadata()

	md := r.Metadata()
	if md == nil {
		t.Fatal("expected metadata map")
	}
	if titles := md["title"].([]string); len(titles) != 1 || titles[0] != "Title" {
		t.Errorf("unexpected title metadata: %v", md["title"])
	}
	if creators := md["creator"].([]string); len(creators) != 1 || creators[0] != "Jane" {
		t.Errorf("unexpected creator metadata: %v", md["creator"])
	}
	if ids := md["identifiers"].([]string); len(ids) != 1 || ids[0] != "urn:x" {
		t.Errorf("unexpected identifiers: %v", md["identifiers"])
	}
}

func TestNormalizeWhitespace(t *testing.T) {
	if got := normalizeWhitespace("  hello   world  "); got != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", got)
	}
}

func TestReader_Refines(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {
				Metadata: pkg.Metadata{
					Titles: []pkg.DCTitle{{ID: "title", Value: "My Book"}},
					Meta:   []pkg.Meta{{Refines: "#title", Property: "title-type", Value: "main"}},
				},
			},
		},
	}}

	refines := r.Refines()
	titleRefines := refines["title"]
	if titleRefines == nil {
		t.Fatalf("expected refines for 'title', got %v", refines)
	}
	if vals := titleRefines["title-type"]; len(vals) != 1 || vals[0] != "main" {
		t.Errorf("unexpected title-type refine: %v", titleRefines)
	}
}

// ---- NavigationCenterExtended ------------------------------------------------

func TestReader_NavigationCenterExtended(t *testing.T) {
	ncxDoc := &ncx.NCX{}
	r := &Reader{epub: &Epub{navigationCenterEXtended: ncxDoc}}

	if got := r.NavigationCenterExtended(); got != ncxDoc {
		t.Errorf("expected NCX pointer, got %v", got)
	}
}
