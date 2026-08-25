package epub

import (
	"strings"
	"testing"

	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/pkg"
	"golang.org/x/net/html"
)

func TestTOC_parseFromHTML(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		title      string
		itemCount  int
		firstHref  string
		firstTitle string
	}{
		{
			name:       "nav with heading and nested list",
			content:    `<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops"><body><nav epub:type="toc"><h2>Table of Contents</h2><ol><li><a href="chap1.xhtml">Chapter 1</a></li><li><a href="chap2.xhtml">Chapter 2</a><ol><li><a href="chap2.xhtml#s1">Section 2.1</a></li></ol></li></ol></nav></body></html>`,
			title:      "Table of Contents",
			itemCount:  2,
			firstHref:  "chap1.xhtml",
			firstTitle: "Chapter 1",
		},
		{
			name:       "nav without heading",
			content:    `<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops"><body><nav epub:type="toc"><ol><li><a href="c1.xhtml">One</a></li></ol></nav></body></html>`,
			title:      "",
			itemCount:  1,
			firstHref:  "c1.xhtml",
			firstTitle: "One",
		},
		{
			name:       "toc nav is selected over landmarks",
			content:    `<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops"><body><nav epub:type="landmarks"><ol><li><a href="cover.xhtml">Cover</a></li></ol></nav><nav epub:type="toc"><ol><li><a href="chap.xhtml">Chapter</a></li></ol></nav></body></html>`,
			title:      "",
			itemCount:  1,
			firstHref:  "chap.xhtml",
			firstTitle: "Chapter",
		},
		{
			name:       "toc nav identified by doc-toc role",
			content:    `<html xmlns="http://www.w3.org/1999/xhtml"><body><nav role="doc-toc"><ol><li><a href="c.xhtml">Role Toc</a></li></ol></nav></body></html>`,
			title:      "",
			itemCount:  1,
			firstHref:  "c.xhtml",
			firstTitle: "Role Toc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := html.Parse(strings.NewReader(tt.content))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var toc TOC
			if err := toc.parseFromHTML(node); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if toc.Title != tt.title {
				t.Errorf("expected title %q, got %q", tt.title, toc.Title)
			}
			if len(toc.Items) != tt.itemCount {
				t.Fatalf("expected %d items, got %d", tt.itemCount, len(toc.Items))
			}
			if toc.Items[0].Href != tt.firstHref {
				t.Errorf("expected first href %q, got %q", tt.firstHref, toc.Items[0].Href)
			}
			if toc.Items[0].Title != tt.firstTitle {
				t.Errorf("expected first title %q, got %q", tt.firstTitle, toc.Items[0].Title)
			}
		})
	}
}

func TestTOC_parseFromHTML_MissingToc(t *testing.T) {
	node, err := html.Parse(strings.NewReader(`<html><body><p>no nav</p></body></html>`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var toc TOC
	if err := toc.parseFromHTML(node); err == nil {
		t.Errorf("expected error when toc nav is missing")
	}
}

func TestTOC_ReadContentHTML(t *testing.T) {
	r := &Reader{epub: &Epub{
		resources: []PublicationResource{
			{ID: "doc", Href: "chap.xhtml", MIMEType: pkg.MediaTypeXHTML, Content: []byte("<html><body><p>Chapter</p></body></html>")},
		},
	}}

	toc := TOC{reader: r, Href: "chap.xhtml"}
	if doc := toc.ReadContentHTML(); doc == nil {
		t.Error("expected content document, got nil")
	}

	empty := TOC{}
	if doc := empty.ReadContentHTML(); doc != nil {
		t.Errorf("expected nil for empty toc, got %v", doc)
	}
}

func TestTOC_parseNCX(t *testing.T) {
	toc := TOC{ncx: &ncx.NCX{
		DocTitle: ncx.TextElement{Text: "NCX Title"},
		NavMap: ncx.NavMap{NavPoints: []ncx.NavPoint{
			{
				NavLabel: ncx.NavLabel{Text: "Chapter 1"},
				Content:  ncx.Content{Src: "ch1.xhtml"},
				NavPoints: []ncx.NavPoint{
					{NavLabel: ncx.NavLabel{Text: "Section"}, Content: ncx.Content{Src: "ch1.xhtml#s1"}},
				},
			},
		}},
	}}

	toc.parseNCX()

	if toc.Title != "NCX Title" {
		t.Errorf("expected title %q, got %q", "NCX Title", toc.Title)
	}
	if len(toc.Items) != 1 || toc.Items[0].Title != "Chapter 1" {
		t.Fatalf("unexpected items: %+v", toc.Items)
	}
	if len(toc.Items[0].Items) != 1 || toc.Items[0].Items[0].Title != "Section" {
		t.Errorf("unexpected nested item: %+v", toc.Items[0].Items)
	}
}

func TestTOC_convertNavPointsToTOCItems(t *testing.T) {
	toc := TOC{}
	items := toc.convertNavPointsToTOCItems([]ncx.NavPoint{
		{
			NavLabel: ncx.NavLabel{Text: "A"},
			Content:  ncx.Content{Src: "a.xhtml"},
			NavPoints: []ncx.NavPoint{
				{NavLabel: ncx.NavLabel{Text: "A1"}, Content: ncx.Content{Src: "a.xhtml#1"}},
			},
		},
		{NavLabel: ncx.NavLabel{Text: "B"}, Content: ncx.Content{Src: "b.xhtml"}},
	})

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Title != "A" || len(items[0].Items) != 1 || items[0].Items[0].Title != "A1" {
		t.Errorf("unexpected item: %+v", items[0])
	}
	if got := toc.convertNavPointsToTOCItems(nil); got != nil {
		t.Errorf("expected nil for empty navpoints, got %v", got)
	}
}

func TestTOC_flattenTOC(t *testing.T) {
	toc := TOC{ncx: &ncx.NCX{NavMap: ncx.NavMap{NavPoints: []ncx.NavPoint{
		{
			NavLabel: ncx.NavLabel{Text: "A"},
			Content:  ncx.Content{Src: "a.xhtml"},
			NavPoints: []ncx.NavPoint{
				{NavLabel: ncx.NavLabel{Text: "A1"}, Content: ncx.Content{Src: "a.xhtml#1"}},
			},
		},
		{NavLabel: ncx.NavLabel{Text: "B"}, Content: ncx.Content{Src: "b.xhtml"}},
	}}}}

	flat := toc.flattenTOC()
	if len(flat) != 3 {
		t.Fatalf("expected 3 flat items, got %d", len(flat))
	}
	if flat[0].Title != "A" || flat[1].Title != "A1" || flat[2].Title != "B" {
		t.Errorf("unexpected flat order: %+v", flat)
	}

	if got := (&TOC{}).flattenTOC(); got != nil {
		t.Errorf("expected nil for empty ncx, got %v", got)
	}
}

func TestTOC_rangeNavMap(t *testing.T) {
	toc := TOC{ncx: &ncx.NCX{NavMap: ncx.NavMap{NavPoints: []ncx.NavPoint{
		{
			NavLabel: ncx.NavLabel{Text: "A"},
			Content:  ncx.Content{Src: "a.xhtml"},
			NavPoints: []ncx.NavPoint{
				{NavLabel: ncx.NavLabel{Text: "A1"}, Content: ncx.Content{Src: "a.xhtml#1"}},
			},
		},
	}}}}

	var depths []int
	toc.rangeNavMap(func(np ncx.NavPoint, depth int) {
		depths = append(depths, depth)
	})

	if len(depths) != 2 || depths[0] != 0 || depths[1] != 1 {
		t.Errorf("unexpected depths: %v", depths)
	}
}

func TestTOC_getTOCByLevel(t *testing.T) {
	toc := TOC{ncx: &ncx.NCX{NavMap: ncx.NavMap{NavPoints: []ncx.NavPoint{
		{
			NavLabel: ncx.NavLabel{Text: "A"},
			Content:  ncx.Content{Src: "a.xhtml"},
			NavPoints: []ncx.NavPoint{
				{NavLabel: ncx.NavLabel{Text: "A1"}, Content: ncx.Content{Src: "a.xhtml#1"}},
			},
		},
		{NavLabel: ncx.NavLabel{Text: "B"}, Content: ncx.Content{Src: "b.xhtml"}},
	}}}}

	level0 := toc.getTOCByLevel(0)
	if len(level0) != 2 {
		t.Fatalf("expected 2 top-level items, got %d", len(level0))
	}
	level1 := toc.getTOCByLevel(1)
	if len(level1) != 1 || level1[0].Title != "A1" {
		t.Errorf("unexpected level 1 items: %+v", level1)
	}

	if got := (&TOC{}).getTOCByLevel(0); got != nil {
		t.Errorf("expected nil for empty ncx, got %v", got)
	}
}

func TestTOC_visitTOC(t *testing.T) {
	toc := TOC{
		Title: "Root",
		Items: []TOC{
			{Title: "A", Items: []TOC{{Title: "A1"}}},
			{Title: "B"},
		},
	}

	var visited []string
	visitTOC(&toc, func(t *TOC, depth int) {
		visited = append(visited, t.Title)
	})

	if len(visited) != 4 {
		t.Fatalf("expected 4 visited nodes, got %v", visited)
	}
	if visited[0] != "Root" || visited[1] != "A" || visited[2] != "A1" || visited[3] != "B" {
		t.Errorf("unexpected visit order: %v", visited)
	}
}

func TestReader_TableOfContents_NCXFallback(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition: "default",
		packagePubs: map[string]*pkg.Package{
			"default": {},
		},
		navigationCenterEXtended: &ncx.NCX{
			DocTitle: ncx.TextElement{Text: "NCX Only"},
			NavMap:   ncx.NavMap{NavPoints: []ncx.NavPoint{{NavLabel: ncx.NavLabel{Text: "Chapter"}, Content: ncx.Content{Src: "ch.xhtml"}}}},
		},
	}}

	toc, err := r.TableOfContents()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toc.Title != "NCX Only" || len(toc.Items) != 1 || toc.Items[0].Title != "Chapter" {
		t.Errorf("unexpected NCX toc: %+v", toc)
	}
}

func TestReader_TableOfContents_Empty(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition:   "default",
		packagePubs: map[string]*pkg.Package{"default": {}},
	}}

	if _, err := r.TableOfContents(); err != nil {
		t.Errorf("expected no error for empty toc, got %v", err)
	}
}

func TestReader_Landmarks(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": {},
			},
			resources: []PublicationResource{
				{
					ID:         "nav",
					Href:       "nav.xhtml",
					Properties: pkg.NavProperty,
					MIMEType:   pkg.MediaTypeXHTML,
					Content:    []byte(`<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops"><body><nav epub:type="landmarks"><ol><li><a epub:type="bodymatter" href="text/chapter-1.xhtml">Begin</a></li><li><a epub:type="toc" href="nav.xhtml">Contents</a></li></ol></nav></body></html>`),
				},
			},
		},
	}

	landmarks := r.Landmarks()
	if len(landmarks) != 2 {
		t.Fatalf("expected 2 landmarks, got %d", len(landmarks))
	}
	if landmarks[0].Title != "Begin" || landmarks[0].Href != "text/chapter-1.xhtml" || landmarks[0].Type != "bodymatter" {
		t.Errorf("unexpected first landmark: %+v", landmarks[0])
	}
	if landmarks[1].Type != "toc" {
		t.Errorf("unexpected second landmark type: %+v", landmarks[1])
	}
}

func TestReader_Landmarks_NoNav(t *testing.T) {
	r := &Reader{epub: &Epub{
		rendition:   "default",
		packagePubs: map[string]*pkg.Package{"default": {}},
	}}

	if got := r.Landmarks(); got != nil {
		t.Errorf("expected nil landmarks without nav, got %v", got)
	}
}

func TestReader_PageList(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": {},
			},
			resources: []PublicationResource{
				{
					ID:         "nav",
					Href:       "nav.xhtml",
					Properties: pkg.NavProperty,
					MIMEType:   pkg.MediaTypeXHTML,
					Content:    []byte(`<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops"><body><nav epub:type="page-list"><ol><li><a href="text/chapter-1.xhtml#page1">1</a></li><li><a href="text/chapter-2.xhtml#page2">2</a></li></ol></nav></body></html>`),
				},
			},
		},
	}

	pages := r.PageList()
	if len(pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(pages))
	}
	if pages[0].Title != "1" || pages[0].Href != "text/chapter-1.xhtml#page1" {
		t.Errorf("unexpected first page: %+v", pages[0])
	}
}

func TestTOC_JSON(t *testing.T) {
	tests := []struct {
		name     string
		toc      *TOC
		expected string
	}{
		{
			name: "nested items",
			toc: &TOC{
				Title: "Chapter 1",
				Href:  "chapter1.html",
				Items: []TOC{
					{Title: "Section 1.1", Href: "chapter1.html#section1"},
					{Title: "Section 1.2", Href: "chapter1.html#section2"},
				},
			},
			expected: `{"title":"Chapter 1","href":"chapter1.html","items":[{"title":"Section 1.1","href":"chapter1.html#section1"},{"title":"Section 1.2","href":"chapter1.html#section2"}]}`,
		},
		{
			name:     "empty TOC",
			toc:      &TOC{},
			expected: `{}`,
		},
		{
			name:     "no items",
			toc:      &TOC{Title: "Chapter 2", Href: "chapter2.html"},
			expected: `{"title":"Chapter 2","href":"chapter2.html"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.toc.JSON()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(b) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(b))
			}
		})
	}
}
