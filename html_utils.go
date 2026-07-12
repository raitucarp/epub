package epub

import (
	"strings"

	"golang.org/x/net/html"
)

// FindNode traverses the HTML tree recursively and returns the first node that satisfies the predicate.
func FindNode(node *html.Node, predicate func(*html.Node) bool) *html.Node {
	if predicate(node) {
		return node
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if result := FindNode(c, predicate); result != nil {
			return result
		}
	}

	return nil
}

// GetTextContent extracts and concatenates all text content within an HTML node recursively.
func GetTextContent(node *html.Node) string {
	var text strings.Builder
	var extractText func(*html.Node)

	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(node)
	return strings.TrimSpace(text.String())
}
