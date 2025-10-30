package epub

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type toc struct {
	Title string
	Link  string
	Items []toc
}

func (t *toc) parseFromHTML(node *html.Node) error {
	navNode := t.findNavNode(node)
	if navNode == nil {
		return fmt.Errorf("nav element with id='toc' not found")
	}

	t.parseNav(navNode)
	return nil
}

func (t *toc) findNavNode(node *html.Node) *html.Node {
	var navNode *html.Node
	var findNav func(*html.Node)

	findNav = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Val == "toc" {
					navNode = n
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findNav(c)
		}
	}

	findNav(node)
	return navNode
}

func (t *toc) parseNav(navNode *html.Node) {
	// Extract title from h2
	for c := navNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			t.Title = t.getTextContent(c)
			break
		}
	}

	// Find and parse the main list (ol or ul)
	for c := navNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "ol" || c.Data == "ul") {
			t.parseList(c)
			break
		}
	}
}

func (t *toc) parseList(listNode *html.Node) {
	for li := listNode.FirstChild; li != nil; li = li.NextSibling {
		if li.Type == html.ElementNode && li.Data == "li" {
			item := t.parseListItem(li)
			t.Items = append(t.Items, item)
		}
	}
}

func (t *toc) parseListItem(liNode *html.Node) toc {
	item := toc{}

	// Find the first anchor tag and any nested list
	var anchor *html.Node
	var subList *html.Node

	for c := liNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.Data == "a" && anchor == nil {
				anchor = c
			} else if c.Data == "ol" || c.Data == "ul" {
				subList = c
			}
		}
	}

	if anchor != nil {
		// Extract link (href)
		for _, attr := range anchor.Attr {
			if attr.Key == "href" {
				item.Link = attr.Val
				break
			}
		}

		// Extract title (text content, ignoring spans)
		item.Title = t.getTextContent(anchor)
	}

	// Parse nested list if present
	if subList != nil {
		item.parseList(subList)
	}

	return item
}

func (t *toc) getTextContent(node *html.Node) string {
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
