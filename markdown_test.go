package epub

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Chapter 1", "chapter-1"},
		{"  The   Great   Escape!  ", "the-great-escape"},
		{"100% Pure", "100-pure"},
		{"Héllo Wörld", "h-llo-w-rld"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := slugify(tt.in); got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestMarkdownToXHTML(t *testing.T) {
	src := []byte("# Chapter One\n\nSome **bold** text.\n\n## Section 1.1\n\nMore text.\n")

	res, err := markdownToXHTML(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(res.body, "<h1 id=\"chapter-one\">") {
		t.Errorf("expected h1 with id, got %q", res.body)
	}
	if !strings.Contains(res.body, "<h2 id=\"section-1-1\">") {
		t.Errorf("expected h2 with id, got %q", res.body)
	}
	if !strings.Contains(res.body, "<strong>bold</strong>") {
		t.Errorf("expected bold text, got %q", res.body)
	}

	if len(res.headings) != 2 {
		t.Fatalf("expected 2 headings, got %d", len(res.headings))
	}
	if res.headings[0].text != "Chapter One" || res.headings[0].level != 1 || res.headings[0].id != "chapter-one" {
		t.Errorf("unexpected first heading: %+v", res.headings[0])
	}
	if res.headings[1].text != "Section 1.1" || res.headings[1].level != 2 {
		t.Errorf("unexpected second heading: %+v", res.headings[1])
	}
}

func TestBuildMarkdownTOC(t *testing.T) {
	docs := []markdownDocument{
		{
			href: "chapter-1.xhtml",
			headings: []headingInfo{
				{level: 1, text: "Chapter 1", id: "chapter-1"},
				{level: 2, text: "Section 1.1", id: "section-1-1"},
				{level: 2, text: "Section 1.2", id: "section-1-2"},
			},
		},
		{
			href: "chapter-2.xhtml",
			headings: []headingInfo{
				{level: 1, text: "Chapter 2", id: "chapter-2"},
			},
		},
	}

	toc := buildMarkdownTOC(docs)

	if len(toc.Items) != 2 {
		t.Fatalf("expected 2 top-level items, got %d", len(toc.Items))
	}
	if toc.Items[0].Title != "Chapter 1" || toc.Items[0].Href != "chapter-1.xhtml#chapter-1" {
		t.Errorf("unexpected first item: %+v", toc.Items[0])
	}
	if len(toc.Items[0].Items) != 2 {
		t.Fatalf("expected 2 sections under chapter 1, got %d", len(toc.Items[0].Items))
	}
	if toc.Items[0].Items[1].Href != "chapter-1.xhtml#section-1-2" {
		t.Errorf("unexpected nested href: %+v", toc.Items[0].Items[1])
	}
	if toc.Items[1].Title != "Chapter 2" {
		t.Errorf("unexpected second item: %+v", toc.Items[1])
	}
}

func TestBuildMarkdownTOC_DeepFirstHeading(t *testing.T) {
	docs := []markdownDocument{
		{
			href: "chapter-1.xhtml",
			headings: []headingInfo{
				{level: 2, text: "Section 1.1", id: "section-1-1"},
			},
		},
		{
			href: "chapter-2.xhtml",
			headings: []headingInfo{
				{level: 1, text: "Chapter 2", id: "chapter-2"},
			},
		},
	}

	toc := buildMarkdownTOC(docs)
	if len(toc.Items) != 2 {
		t.Fatalf("expected 2 top-level items, got %d", len(toc.Items))
	}
	if toc.Items[0].Title != "Section 1.1" || toc.Items[1].Title != "Chapter 2" {
		t.Errorf("unexpected items: %+v", toc.Items)
	}
}

func TestWriter_AddMarkdown(t *testing.T) {
	w := New("urn:md:test")
	w.Title("Markdown Book")
	w.Languages("en")

	res, err := w.AddMarkdown("chapter-1.md", []byte("# Chapter 1\n\nHello.\n"))
	if err != nil {
		t.Fatalf("AddMarkdown error: %v", err)
	}
	if res.Href != "chapter-1.xhtml" {
		t.Errorf("expected href chapter-1.xhtml, got %q", res.Href)
	}
	if res.MIMEType != "application/xhtml+xml" {
		t.Errorf("expected xhtml media type, got %q", res.MIMEType)
	}
	if !strings.Contains(string(res.Content), "Hello.") {
		t.Errorf("expected converted content, got %q", string(res.Content))
	}
}

func TestWriter_AddMarkdownFile(t *testing.T) {
	withWorkingDir(t, func(dir string) {
		if err := os.WriteFile(filepath.Join(dir, "chapter.md"), []byte("# Chapter\n\nText.\n"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		w := New("urn:md:file")
		w.Title("Book")
		w.Languages("en")

		res, err := w.AddMarkdownFile("chapter.md")
		if err != nil {
			t.Fatalf("AddMarkdownFile error: %v", err)
		}
		if res.Href != "chapter.xhtml" {
			t.Errorf("expected href chapter.xhtml, got %q", res.Href)
		}
	})
}

func TestWriter_AddMarkdownDirectory(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"chapter-2.md": "# Chapter 2\n\nSecond chapter.\n",
		"chapter-1.md": "# Chapter 1\n\nFirst chapter.\n\n## Section 1.1\n\nDetail.\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	w := New("urn:md:dir")
	w.Title("Directory Book")
	w.Languages("en")

	if err := w.AddMarkdownDirectory(dir); err != nil {
		t.Fatalf("AddMarkdownDirectory error: %v", err)
	}

	// Files are added in file-name order: chapter-1 before chapter-2.
	spine := w.epub.SelectedPackage().Spine.ItemRefs
	if len(spine) != 2 {
		t.Fatalf("expected 2 spine items, got %d", len(spine))
	}
	if spine[0].IDRef != "chapter-1.xhtml" || spine[1].IDRef != "chapter-2.xhtml" {
		t.Errorf("unexpected spine order: %s, %s", spine[0].IDRef, spine[1].IDRef)
	}

	toc := w.epub.navigationCenterEXtended
	if toc == nil {
		t.Fatal("expected NCX to be generated")
	}
	points := toc.NavMap.NavPoints
	if len(points) != 2 {
		t.Fatalf("expected 2 navpoints, got %d", len(points))
	}
	if points[0].NavLabel.Text != "Chapter 1" || points[0].Content.Src != "chapter-1.xhtml#chapter-1" {
		t.Errorf("unexpected first navpoint: %+v", points[0])
	}
	if len(points[0].NavPoints) != 1 || points[0].NavPoints[0].NavLabel.Text != "Section 1.1" {
		t.Errorf("unexpected nested navpoint: %+v", points[0].NavPoints)
	}
}

func TestWriter_AddMarkdownDirectory_MissingDir(t *testing.T) {
	w := New("urn:md:missing")
	if err := w.AddMarkdownDirectory(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Error("expected error for missing directory, got nil")
	}
}

func TestReplaceExt(t *testing.T) {
	if got := replaceExt("chapter-1.md", ".xhtml"); got != "chapter-1.xhtml" {
		t.Errorf("expected chapter-1.xhtml, got %q", got)
	}
	if got := replaceExt("dir/chapter.markdown", ".xhtml"); got != "chapter.xhtml" {
		t.Errorf("expected chapter.xhtml, got %q", got)
	}
}

func TestHeadingLevel(t *testing.T) {
	tests := []struct {
		tag  string
		lvl  int
		want bool
	}{
		{"h1", 1, true},
		{"h2", 2, true},
		{"h3", 3, true},
		{"h4", 4, true},
		{"h5", 5, true},
		{"h6", 6, true},
		{"h7", 0, false},
		{"p", 0, false},
		{"div", 0, false},
		{"header", 0, false},
	}
	for _, tt := range tests {
		lvl, ok := headingLevel(tt.tag)
		if ok != tt.want || lvl != tt.lvl {
			t.Errorf("headingLevel(%q) = (%d, %v), want (%d, %v)", tt.tag, lvl, ok, tt.lvl, tt.want)
		}
	}
}

func TestWriter_AddMarkdownFile_Errors(t *testing.T) {
	w := New("urn:md:err")
	// Non-local path
	if _, err := w.AddMarkdownFile("../outside.md"); err == nil {
		t.Error("expected error for non-local path in AddMarkdownFile")
	}
}

func TestMarkdownToXHTML_DuplicateHeadings(t *testing.T) {
	src := []byte("# Same Heading\n\nText\n\n# Same Heading\n\nMore text\n\n# Same Heading\n")
	res, err := markdownToXHTML(src)
	if err != nil {
		t.Fatalf("markdownToXHTML: %v", err)
	}
	if len(res.headings) != 3 {
		t.Fatalf("expected 3 headings, got %d", len(res.headings))
	}
	if res.headings[0].id != "same-heading" {
		t.Errorf("expected same-heading, got %s", res.headings[0].id)
	}
	if res.headings[1].id != "same-heading-2" {
		t.Errorf("expected same-heading-2, got %s", res.headings[1].id)
	}
	if res.headings[2].id != "same-heading-3" {
		t.Errorf("expected same-heading-3, got %s", res.headings[2].id)
	}
}
