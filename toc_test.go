package epub

import (
	"strings"
	"testing"

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
					{
						Title: "Section 1.1",
						Href:  "chapter1.html#section1",
					},
					{
						Title: "Section 1.2",
						Href:  "chapter1.html#section2",
					},
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
			name: "no items",
			toc: &TOC{
				Title: "Chapter 2",
				Href:  "chapter2.html",
			},
			expected: `{"title":"Chapter 2","href":"chapter2.html"}`,
		},
		{
			name: "special characters",
			toc: &TOC{
				Title: "Chapter 3: \"The Awakening\" & <Others>",
				Href:  "chapter3.html?foo=bar&baz=qux",
			},
			expected: `{"title":"Chapter 3: \"The Awakening\" \u0026 \u003cOthers\u003e","href":"chapter3.html?foo=bar\u0026baz=qux"}`,
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
