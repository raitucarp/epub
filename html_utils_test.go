package epub

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestFindNode(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(`<html><body><div id="a"><p class="target">Hi</p></div></body></html>`))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	t.Run("matches node", func(t *testing.T) {
		found := FindNode(doc, func(n *html.Node) bool {
			return n.Type == html.ElementNode && n.Data == "p"
		})
		if found == nil || found.Data != "p" {
			t.Errorf("expected to find p element, got %v", found)
		}
	})

	t.Run("returns nil when no match", func(t *testing.T) {
		found := FindNode(doc, func(n *html.Node) bool {
			return n.Type == html.ElementNode && n.Data == "article"
		})
		if found != nil {
			t.Errorf("expected nil, got %v", found)
		}
	})
}

func TestGetTextContent(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(`<html><body><div>Hello <span>nested <b>text</b></span> here</div></body></html>`))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	div := FindNode(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "div"
	})
	if div == nil {
		t.Fatal("expected to find div element")
	}

	if got := GetTextContent(div); got != "Hello nested text here" {
		t.Errorf("expected concatenated text, got %q", got)
	}
}
