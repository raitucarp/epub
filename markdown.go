package epub

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/net/html"
)

// mdConverter renders Markdown to an HTML fragment. GitHub Flavored Markdown
// extensions (tables, task lists, strikethrough, and autolinks) are enabled.
var mdConverter = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

// headingInfo describes a single heading extracted from a Markdown document.
type headingInfo struct {
	level int
	text  string
	id    string
}

// markdownDocument is the result of converting a single Markdown file to an
// XHTML content document.
type markdownDocument struct {
	href     string
	headings []headingInfo
}

// AddMarkdown converts content from Markdown to XHTML and adds it to the
// publication as a spine content document. The extension of filename is
// replaced with .xhtml to derive the href of the resulting document.
func (w *Writer) AddMarkdown(filename string, content []byte) (PublicationResource, error) {
	res, _, err := w.addMarkdown(filename, content)
	return res, err
}

// AddMarkdownFile reads the Markdown file at name from disk and adds it to the
// publication as a spine content document.
func (w *Writer) AddMarkdownFile(name string) (PublicationResource, error) {
	if !filepath.IsLocal(name) {
		return PublicationResource{}, fmt.Errorf("invalid path: path must be local")
	}

	root, err := os.OpenRoot(".")
	if err != nil {
		return PublicationResource{}, err
	}
	defer root.Close()

	data, err := root.ReadFile(name)
	if err != nil {
		return PublicationResource{}, err
	}

	return w.AddMarkdown(filepath.Base(name), data)
}

// AddMarkdownDirectory reads every Markdown file in dir, converts each to an
// XHTML content document, and adds them to the spine in file-name order. A
// table of contents is generated from the heading structure of the documents
// and registered with the publication. The publication identifier, title,
// author, and language must be set on the Writer before calling Write.
func (w *Writer) AddMarkdownDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		switch strings.ToLower(filepath.Ext(e.Name())) {
		case ".md", ".markdown":
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	docs := make([]markdownDocument, 0, len(names))
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		_, headings, err := w.addMarkdown(name, data)
		if err != nil {
			return err
		}
		docs = append(docs, markdownDocument{
			href:     replaceExt(name, ".xhtml"),
			headings: headings,
		})
	}

	toc := buildMarkdownTOC(docs)
	if len(toc.Items) == 0 {
		return nil
	}
	return w.TableOfContents("toc", toc)
}

// addMarkdown converts content to XHTML, adds it as a content document, and
// returns the created resource together with the headings found in the
// document.
func (w *Writer) addMarkdown(filename string, content []byte) (PublicationResource, []headingInfo, error) {
	doc, err := markdownToXHTML(content)
	if err != nil {
		return PublicationResource{}, nil, err
	}

	res := w.AddContent(replaceExt(filename, ".xhtml"), []byte(doc.body))
	return res, doc.headings, nil
}

// markdownResult is the XHTML body and extracted headings of a Markdown file.
type markdownResult struct {
	body     string
	headings []headingInfo
}

// markdownToXHTML converts Markdown source to an XHTML fragment, assigning
// stable id attributes to headings and collecting them for table-of-contents
// generation.
func markdownToXHTML(src []byte) (markdownResult, error) {
	var buf bytes.Buffer
	if err := mdConverter.Convert(src, &buf); err != nil {
		return markdownResult{}, err
	}

	doc, err := html.Parse(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return markdownResult{}, err
	}

	body := getBody(doc)
	if body == nil {
		return markdownResult{body: buf.String()}, nil
	}

	var headings []headingInfo
	used := map[string]int{}
	for n := range body.Descendants() {
		level, ok := headingLevel(n.Data)
		if !ok {
			continue
		}
		text := strings.TrimSpace(GetTextContent(n))
		id := slugify(text)
		if used[id] > 0 {
			used[id]++
			id = fmt.Sprintf("%s-%d", id, used[id])
		} else {
			used[id] = 1
		}
		n.Attr = append(n.Attr, html.Attribute{Key: "id", Val: id})
		headings = append(headings, headingInfo{level: level, text: text, id: id})
	}

	return markdownResult{body: renderChildren(body), headings: headings}, nil
}

// headingLevel returns the 1-based heading level for an element name, or false
// if the element is not a heading.
func headingLevel(name string) (int, bool) {
	switch name {
	case "h1":
		return 1, true
	case "h2":
		return 2, true
	case "h3":
		return 3, true
	case "h4":
		return 4, true
	case "h5":
		return 5, true
	case "h6":
		return 6, true
	default:
		return 0, false
	}
}

// slugify derives a URL-friendly fragment identifier from a heading.
func slugify(s string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case !lastDash:
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

// renderChildren renders the child nodes of n as a single HTML string.
func renderChildren(n *html.Node) string {
	var b bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		html.Render(&b, c)
	}
	return b.String()
}

// replaceExt replaces the extension of filename with ext.
func replaceExt(filename, ext string) string {
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename)) + ext
}

// buildMarkdownTOC builds a nested TOC from the headings of a set of Markdown
// documents. Headings are nested depth-first by their level; a heading deeper
// than the current nesting is promoted to the nearest available level.
func buildMarkdownTOC(docs []markdownDocument) TOC {
	var root TOC
	// parents[i] points to the TOC item that is the current parent for a
	// heading at level i+1 (parents[0] is the root).
	parents := []*TOC{&root}
	for _, doc := range docs {
		for _, h := range doc.headings {
			item := TOC{Title: h.text, Href: doc.href + "#" + h.id}

			level := h.level
			if level > len(parents) {
				level = len(parents)
			}
			// Pop levels deeper than this heading.
			for len(parents) > level {
				parents = parents[:len(parents)-1]
			}

			parent := parents[len(parents)-1]
			parent.Items = append(parent.Items, item)
			parents = append(parents, &parent.Items[len(parent.Items)-1])
		}
	}
	return root
}
