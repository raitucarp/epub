package epub

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/raitucarp/epub/ncx"
	"golang.org/x/net/html"
)

// TOC represents the publication's table of contents in normalized form.
// The structure abstracts differences between NAV (EPUB 3) and NCX (EPUB 2)
// so higher-level code can work with a unified interface.
type TOC struct {
	Title string `json:"title,omitempty"`
	Href  string `json:"href,omitempty"`
	Items []TOC  `json:"items,omitempty"`

	reader *Reader
	ncx    *ncx.NCX
}

// JSON marshals the table of contents structure into JSON format. This is useful
// for external tools, logging, debugging, or serialization to other formats.
func (t *TOC) JSON() ([]byte, error) {
	return json.Marshal(t)
}

// ReadContentHTML returns the content document associated with the currently
// selected table of contents entry. The returned document is parsed into an
// html.Node tree. Behavior depends on TOC internal navigation selection state.
func (t *TOC) ReadContentHTML() (content *html.Node) {
	if t.Href != "" && t.reader != nil {
		return t.reader.ReadContentHTMLByHref(t.Href)
	}
	return
}

func (t *TOC) parseFromHTML(node *html.Node) error {
	navNode := t.findNavNode(node)
	if navNode == nil {
		return fmt.Errorf("nav element with id='toc' not found")
	}

	t.parseNav(navNode)
	return nil
}

func (t *TOC) findNavNode(node *html.Node) *html.Node {
	return findNavByType(node, "toc")
}

// findNavByType locates a nav element of the given epub:type within an HTML
// node tree. The toc nav may also be identified by its doc-toc role.
func findNavByType(node *html.Node, epubType string) *html.Node {
	return FindNode(node, func(n *html.Node) bool {
		if n.Type != html.ElementNode || n.Data != "nav" {
			return false
		}

		for _, attr := range n.Attr {
			if attr.Key == "epub:type" && attr.Val == epubType {
				return true
			}
			if epubType == "toc" && attr.Key == "role" && attr.Val == "doc-toc" {
				return true
			}
		}
		return false
	})
}

// navDocument returns the parsed HTML node of the EPUB navigation document,
// or nil if the publication does not declare one.
func (r *Reader) navDocument() *html.Node {
	resourceWithNavIndex := slices.IndexFunc(r.epub.resources, func(res PublicationResource) bool {
		return res.Properties == "nav"
	})
	if resourceWithNavIndex < 0 {
		return nil
	}

	return r.ReadContentHTMLById(r.epub.resources[resourceWithNavIndex].ID)
}

func (t *TOC) parseNav(navNode *html.Node) {
	// Extract title from the heading element (h1-h6), if present.
	for c := navNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && isHeading(c.Data) {
			t.Title = GetTextContent(c)
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

func isHeading(name string) bool {
	switch name {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return true
	default:
		return false
	}
}

func (t *TOC) parseList(listNode *html.Node) {
	var count int
	for li := listNode.FirstChild; li != nil; li = li.NextSibling {
		if li.Type == html.ElementNode && li.Data == "li" {
			count++
		}
	}

	if count > 0 {
		t.Items = slices.Grow(t.Items, count)
	}

	for li := listNode.FirstChild; li != nil; li = li.NextSibling {
		if li.Type == html.ElementNode && li.Data == "li" {
			item := t.parseListItem(li)
			t.Items = append(t.Items, item)
		}
	}
}

func (t *TOC) parseListItem(liNode *html.Node) TOC {
	item := TOC{}

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
				item.Href = attr.Val
				break
			}
		}

		// Extract title (text content, ignoring spans)
		item.Title = GetTextContent(anchor)
	}

	// Parse nested list if present
	if subList != nil {
		item.parseList(subList)
	}

	return item
}

// convertNavPointsToTOCItems recursively converts NavPoints to TOC items
func (t *TOC) convertNavPointsToTOCItems(navPoints []ncx.NavPoint) []TOC {
	if len(navPoints) == 0 {
		return nil
	}

	var tocItems []TOC
	for _, navPoint := range navPoints {
		tocItem := TOC{
			Title: navPoint.NavLabel.Text,
			Href:  navPoint.Content.Src,
		}

		// Recursively process nested navPoints
		if len(navPoint.NavPoints) > 0 {
			tocItem.Items = t.convertNavPointsToTOCItems(navPoint.NavPoints)
		}

		tocItems = append(tocItems, tocItem)
	}

	return tocItems
}

// flattenTOC returns a flat slice of all TOC items in depth-first order
func (t *TOC) flattenTOC() []TOC {
	if t.ncx == nil {
		return nil
	}

	var flatTOC []TOC
	t.flattenNavPoints(t.ncx.NavMap.NavPoints, &flatTOC)
	return flatTOC
}

// flattenNavPoints recursively flattens navPoints into a slice
func (t *TOC) flattenNavPoints(navPoints []ncx.NavPoint, result *[]TOC) {
	for _, navPoint := range navPoints {
		tocItem := TOC{
			Title: navPoint.NavLabel.Text,
			Href:  navPoint.Content.Src,
		}

		*result = append(*result, tocItem)

		// Recursively process nested navPoints
		if len(navPoint.NavPoints) > 0 {
			t.flattenNavPoints(navPoint.NavPoints, result)
		}
	}
}

// rangeNavMap iterates through all navPoints and executes a function for each
func (t *TOC) rangeNavMap(fn func(navPoint ncx.NavPoint, depth int)) {
	if t.ncx == nil {
		return
	}

	t.rangeNavPoints(t.ncx.NavMap.NavPoints, 0, fn)
}

// rangeNavPoints helper function for recursive iteration
func (t *TOC) rangeNavPoints(navPoints []ncx.NavPoint, depth int, fn func(navPoint ncx.NavPoint, depth int)) {
	for _, navPoint := range navPoints {
		fn(navPoint, depth)

		// Recursively process nested navPoints
		if len(navPoint.NavPoints) > 0 {
			t.rangeNavPoints(navPoint.NavPoints, depth+1, fn)
		}
	}
}

// GetTOCByLevel returns TOC items at a specific depth level
func (t *TOC) getTOCByLevel(level int) []TOC {
	if t.ncx == nil {
		return nil
	}

	var result []TOC
	t.rangeNavPoints(t.ncx.NavMap.NavPoints, 0, func(navPoint ncx.NavPoint, depth int) {
		if depth == level {
			result = append(result, TOC{
				Title: navPoint.NavLabel.Text,
				Href:  navPoint.Content.Src,
			})
		}
	})

	return result
}

func (t *TOC) parseNCX() {
	t.Title = t.ncx.DocTitle.Text
	t.Items = t.convertNavPointsToTOCItems(t.ncx.NavMap.NavPoints)
}

func visitTOC(toc *TOC, visitor func(*TOC, int)) {
	var visit func(*TOC, int)
	visit = func(item *TOC, depth int) {
		if item == nil {
			return
		}

		// Visit current node with depth information
		visitor(item, depth)

		// Recursively visit all items with increased depth
		for i := range item.Items {
			visit(&item.Items[i], depth+1)
		}
	}

	visit(toc, 0)
}

func tocToHTMLNode(toc TOC, lang []string) (*html.Node, error) {
	// Create the root html node
	doc := &html.Node{
		Type: html.DocumentNode,
	}

	// Create the HTML root element with all the necessary attributes
	htmlNode := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		Attr: []html.Attribute{
			{Key: "xmlns", Val: "http://www.w3.org/1999/xhtml"},
		},
	}

	for _, l := range lang {

		htmlNode.Attr = append(htmlNode.Attr, html.Attribute{Key: "lang", Val: l})
		htmlNode.Attr = append(htmlNode.Attr, html.Attribute{Key: "xml:lang", Val: l})
	}
	doc.AppendChild(htmlNode)

	// Create head section
	head := &html.Node{
		Type: html.ElementNode,
		Data: "head",
	}
	htmlNode.AppendChild(head)

	title := &html.Node{
		Type: html.ElementNode,
		Data: "title",
	}
	head.AppendChild(title)

	titleText := &html.Node{
		Type: html.TextNode,
		Data: toc.Title,
	}
	title.AppendChild(titleText)

	// Create body section
	body := &html.Node{
		Type: html.ElementNode,
		Data: "body",
		Attr: []html.Attribute{
			{Key: "epub:type", Val: "frontmatter"},
		},
	}
	htmlNode.AppendChild(body)

	// Create main TOC nav
	tocNav := createTOCNav(toc.Items)
	body.AppendChild(tocNav)

	// Create landmarks nav
	landmarksNav := createLandmarksNav(toc.Items)
	body.AppendChild(landmarksNav)

	return doc, nil
}

func createTOCNav(toc []TOC) *html.Node {
	nav := &html.Node{
		Type: html.ElementNode,
		Data: "nav",
		Attr: []html.Attribute{
			{Key: "id", Val: "toc"},
			{Key: "role", Val: "doc-toc"},
			{Key: "epub:type", Val: "toc"},
		},
	}

	// Add title
	h2 := &html.Node{
		Type: html.ElementNode,
		Data: "h2",
		Attr: []html.Attribute{
			{Key: "epub:type", Val: "title"},
		},
	}
	nav.AppendChild(h2)

	h2Text := &html.Node{
		Type: html.TextNode,
		Data: "Table of Contents",
	}
	h2.AppendChild(h2Text)

	// Create the main OL and recursively add TOC items
	ol := &html.Node{
		Type: html.ElementNode,
		Data: "ol",
	}
	nav.AppendChild(ol)

	addTOCItems(ol, toc)

	return nav
}

func createLandmarksNav(toc []TOC) *html.Node {
	nav := &html.Node{
		Type: html.ElementNode,
		Data: "nav",
		Attr: []html.Attribute{
			{Key: "id", Val: "landmarks"},
			{Key: "epub:type", Val: "landmarks"},
		},
	}

	// Add title
	h2 := &html.Node{
		Type: html.ElementNode,
		Data: "h2",
		Attr: []html.Attribute{
			{Key: "epub:type", Val: "title"},
		},
	}
	nav.AppendChild(h2)

	h2Text := &html.Node{
		Type: html.TextNode,
		Data: "Landmarks",
	}
	h2.AppendChild(h2Text)

	// Create landmarks OL
	ol := &html.Node{
		Type: html.ElementNode,
		Data: "ol",
	}
	nav.AppendChild(ol)

	// Find the first chapter for landmarks (assuming it's the first item with sub-items)
	var firstChapter *TOC
	for i := range toc {
		if len(toc[i].Items) > 0 {
			firstChapter = &toc[i]
			break
		}
	}

	if firstChapter != nil {
		li := &html.Node{
			Type: html.ElementNode,
			Data: "li",
		}
		ol.AppendChild(li)

		a := &html.Node{
			Type: html.ElementNode,
			Data: "a",
			Attr: []html.Attribute{
				{Key: "href", Val: firstChapter.Items[0].Href},
				{Key: "epub:type", Val: "bodymatter"},
			},
		}
		li.AppendChild(a)

		aText := &html.Node{
			Type: html.TextNode,
			Data: firstChapter.Title,
		}
		a.AppendChild(aText)
	}

	return nav
}

func addTOCItems(parent *html.Node, items []TOC) {
	for _, item := range items {
		li := &html.Node{
			Type: html.ElementNode,
			Data: "li",
		}
		parent.AppendChild(li)

		// Create anchor tag
		if item.Href != "" {
			a := &html.Node{
				Type: html.ElementNode,
				Data: "a",
				Attr: []html.Attribute{
					{Key: "href", Val: item.Href},
				},
			}

			a.FirstChild = &html.Node{Type: html.TextNode, Data: item.Title}
			li.AppendChild(a)

		} else {
			// If no href, just add the title as text
			titleText := &html.Node{
				Type: html.TextNode,
				Data: item.Title,
			}
			li.AppendChild(titleText)
		}

		// Recursively add sub-items if they exist
		if len(item.Items) > 0 {
			subOl := &html.Node{
				Type: html.ElementNode,
				Data: "ol",
			}
			li.AppendChild(subOl)
			addTOCItems(subOl, item.Items)
		}
	}
}

// TableOfContents returns the TOC version present (e.g., NAV or NCX).
// If both exist, behavior depends on publication version and priority rules.
func (r *Reader) TableOfContents() (toc TOC, err error) {
	toc.reader = r

	if doc := r.navDocument(); doc != nil {
		err = toc.parseFromHTML(doc)
		return
	}

	if r.epub.navigationCenterEXtended != nil {
		toc.ncx = r.epub.navigationCenterEXtended
		toc.parseNCX()
	}

	return
}

// Landmark represents an entry in the landmarks navigation of the EPUB
// navigation document. The Type field holds the epub:type semantic applied to
// the landmark (e.g., "cover", "bodymatter", "toc").
type Landmark struct {
	Title string `json:"title,omitempty"`
	Href  string `json:"href,omitempty"`
	Type  string `json:"type,omitempty"`
}

// navListItem is an internal representation of a flat navigation entry.
type navListItem struct {
	Title string
	Href  string
	Type  string
}

// Landmarks returns the landmarks navigation, which identifies the key
// structural points of the publication (e.g., cover, table of contents, and
// body matter).
func (r *Reader) Landmarks() []Landmark {
	doc := r.navDocument()
	if doc == nil {
		return nil
	}

	navNode := findNavByType(doc, "landmarks")
	if navNode == nil {
		return nil
	}

	var landmarks []Landmark
	for _, item := range parseFlatNavList(navNode) {
		landmarks = append(landmarks, Landmark{Title: item.Title, Href: item.Href, Type: item.Type})
	}
	return landmarks
}

// PageList returns the page-list navigation, which maps the print page
// breakpoints of the publication to their content locations.
func (r *Reader) PageList() []Landmark {
	doc := r.navDocument()
	if doc == nil {
		return nil
	}

	navNode := findNavByType(doc, "page-list")
	if navNode == nil {
		return nil
	}

	var pages []Landmark
	for _, item := range parseFlatNavList(navNode) {
		pages = append(pages, Landmark{Title: item.Title, Href: item.Href, Type: item.Type})
	}
	return pages
}

// parseFlatNavList parses the direct list items of a nav element (landmarks or
// page-list) into a flat slice, preserving the title, href, and epub:type.
func parseFlatNavList(navNode *html.Node) []navListItem {
	var list *html.Node
	for c := navNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "ol" || c.Data == "ul") {
			list = c
			break
		}
	}

	if list == nil {
		return nil
	}

	var items []navListItem
	for li := list.FirstChild; li != nil; li = li.NextSibling {
		if li.Type != html.ElementNode || li.Data != "li" {
			continue
		}

		item := navListItem{}
		for c := li.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode || c.Data != "a" {
				continue
			}

			for _, attr := range c.Attr {
				switch attr.Key {
				case "href":
					item.Href = attr.Val
				case "epub:type":
					item.Type = attr.Val
				}
			}
			item.Title = GetTextContent(c)
			break
		}

		items = append(items, item)
	}

	return items
}
