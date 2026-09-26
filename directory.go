package epub

import (
	"bytes"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/raitucarp/epub/pkg"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v3"
)

// DirectoryMetadata represents metadata configured in a metadata.yml or metadata.yaml file.
// All fields are optional and support flexible scalar or list definitions.
type DirectoryMetadata struct {
	Identifier  string   `yaml:"identifier"`
	ID          string   `yaml:"id"`
	UUID        string   `yaml:"uuid"`
	Title       any      `yaml:"title"`
	Titles      []string `yaml:"titles"`
	Author      any      `yaml:"author"`
	Authors     []string `yaml:"authors"`
	Creator     any      `yaml:"creator"`
	Creators    []string `yaml:"creators"`
	Language    any      `yaml:"language"`
	Languages   []string `yaml:"languages"`
	Lang        string   `yaml:"lang"`
	Description string   `yaml:"description"`
	Desc        string   `yaml:"desc"`
	Publisher   any      `yaml:"publisher"`
	Publishers  []string `yaml:"publishers"`
	Rights      string   `yaml:"rights"`
	Subject     any      `yaml:"subject"`
	Subjects    []string `yaml:"subjects"`
	Tags        []string `yaml:"tags"`
	Date        string   `yaml:"date"`
	PubDate     string   `yaml:"pubdate"`
	Modified    string   `yaml:"modified"`
	Direction   string   `yaml:"direction"`
	Dir         string   `yaml:"dir"`
	Cover       string   `yaml:"cover"`
	CoverImage  string   `yaml:"cover_image"`
	TOCTitle    string   `yaml:"toc_title"`
}

func toStringSlice(val any) []string {
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			return []string{trimmed}
		}
		return nil
	case []string:
		return v
	case []any:
		var result []string
		for _, item := range v {
			if s, ok := item.(string); ok {
				trimmed := strings.TrimSpace(s)
				if trimmed != "" {
					result = append(result, trimmed)
				}
			} else if item != nil {
				trimmed := strings.TrimSpace(fmt.Sprint(item))
				if trimmed != "" {
					result = append(result, trimmed)
				}
			}
		}
		return result
	default:
		s := strings.TrimSpace(fmt.Sprint(val))
		if s != "" {
			return []string{s}
		}
		return nil
	}
}

// GetIdentifier returns the primary publication identifier.
func (m *DirectoryMetadata) GetIdentifier() string {
	for _, id := range []string{m.Identifier, m.ID, m.UUID} {
		if trimmed := strings.TrimSpace(id); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

// GetTitles returns all configured publication titles.
func (m *DirectoryMetadata) GetTitles() []string {
	if len(m.Titles) > 0 {
		return m.Titles
	}
	return toStringSlice(m.Title)
}

// GetAuthors returns all configured creators/authors.
func (m *DirectoryMetadata) GetAuthors() []string {
	if len(m.Authors) > 0 {
		return m.Authors
	}
	if len(m.Creators) > 0 {
		return m.Creators
	}
	if authors := toStringSlice(m.Author); len(authors) > 0 {
		return authors
	}
	return toStringSlice(m.Creator)
}

// GetLanguages returns all configured language codes.
func (m *DirectoryMetadata) GetLanguages() []string {
	if len(m.Languages) > 0 {
		return m.Languages
	}
	if langs := toStringSlice(m.Language); len(langs) > 0 {
		return langs
	}
	if trimmed := strings.TrimSpace(m.Lang); trimmed != "" {
		return []string{trimmed}
	}
	return nil
}

// GetDescription returns the publication description or summary.
func (m *DirectoryMetadata) GetDescription() string {
	if trimmed := strings.TrimSpace(m.Description); trimmed != "" {
		return trimmed
	}
	return strings.TrimSpace(m.Desc)
}

// GetPublishers returns all configured publishers.
func (m *DirectoryMetadata) GetPublishers() []string {
	if len(m.Publishers) > 0 {
		return m.Publishers
	}
	return toStringSlice(m.Publisher)
}

// GetRights returns the rights/license declaration.
func (m *DirectoryMetadata) GetRights() string {
	return strings.TrimSpace(m.Rights)
}

// GetSubjects returns all configured subjects or category tags.
func (m *DirectoryMetadata) GetSubjects() []string {
	if len(m.Subjects) > 0 {
		return m.Subjects
	}
	if len(m.Tags) > 0 {
		return m.Tags
	}
	return toStringSlice(m.Subject)
}

// GetDirection returns the reading progression direction ("ltr" or "rtl").
func (m *DirectoryMetadata) GetDirection() string {
	for _, d := range []string{m.Direction, m.Dir} {
		if trimmed := strings.TrimSpace(d); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

// GetCover returns the relative path of the cover image.
func (m *DirectoryMetadata) GetCover() string {
	for _, c := range []string{m.Cover, m.CoverImage} {
		if trimmed := strings.TrimSpace(c); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

// GetTOCTitle returns the table of contents title (defaults to "Table of Contents").
func (m *DirectoryMetadata) GetTOCTitle() string {
	if trimmed := strings.TrimSpace(m.TOCTitle); trimmed != "" {
		return trimmed
	}
	return "Table of Contents"
}

// GetDate parses and returns the publication date if specified.
func (m *DirectoryMetadata) GetDate() (time.Time, bool) {
	for _, d := range []string{m.Date, m.PubDate} {
		if t, ok := parseDate(d); ok {
			return t, true
		}
	}
	return time.Time{}, false
}

// GetModified parses and returns the last modified timestamp if specified.
func (m *DirectoryMetadata) GetModified() (time.Time, bool) {
	return parseDate(m.Modified)
}

func parseDate(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02",
		"02 Jan 2006",
		"Jan 02, 2006",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// parseMetadataFile looks for metadata.yml, metadata.yaml, book.yml, or book.yaml
// in the specified directory and unmarshals it into DirectoryMetadata.
func parseMetadataFile(dir string) (*DirectoryMetadata, error) {
	candidates := []string{
		"metadata.yml",
		"metadata.yaml",
		"book.yml",
		"book.yaml",
	}

	for _, cand := range candidates {
		fullPath := filepath.Join(dir, cand)
		data, err := os.ReadFile(fullPath)
		if err == nil {
			var meta DirectoryMetadata
			if err := yaml.Unmarshal(data, &meta); err != nil {
				return nil, fmt.Errorf("failed to parse %s: %w", cand, err)
			}
			return &meta, nil
		}
	}
	return nil, nil
}

func (w *Writer) applyMetadata(meta *DirectoryMetadata, dir string) error {
	if meta == nil {
		return nil
	}

	if id := meta.GetIdentifier(); id != "" {
		w.identifier = id
		p := w.epub.SelectedPackage()
		if p != nil {
			found := false
			for i, ident := range p.Metadata.Identifiers {
				if ident.ID == "pub-id" || ident.ID == p.UniqueIdentifier {
					p.Metadata.Identifiers[i].Value = id
					found = true
					break
				}
			}
			if !found {
				p.Metadata.Identifiers = append([]pkg.DCIdentifier{{ID: "pub-id", Value: id}}, p.Metadata.Identifiers...)
			}
		}
	}

	if titles := meta.GetTitles(); len(titles) > 0 {
		w.Title(titles...)
	}

	if authors := meta.GetAuthors(); len(authors) > 0 {
		w.Author(authors...)
	}

	if langs := meta.GetLanguages(); len(langs) > 0 {
		w.Languages(langs...)
	}

	if desc := meta.GetDescription(); desc != "" {
		w.Description(desc)
	}

	if pubs := meta.GetPublishers(); len(pubs) > 0 {
		w.Publisher(pubs...)
	}

	if rights := meta.GetRights(); rights != "" {
		w.Rights(rights)
	}

	if subjects := meta.GetSubjects(); len(subjects) > 0 {
		w.Subject("subject", subjects...)
	}

	if dirStr := meta.GetDirection(); dirStr != "" {
		w.Direction(dirStr)
	}

	if d, ok := meta.GetDate(); ok {
		w.Date(d)
	}

	if m, ok := meta.GetModified(); ok {
		w.Modified(m)
	}

	coverFile := meta.GetCover()
	if coverFile != "" {
		fullPath := filepath.Join(dir, coverFile)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to apply cover '%s': %w", coverFile, err)
		}
		base := filepath.Base(coverFile)
		w.addImageCover(base, data)
		w.MetaContent(map[string]string{"cover": base})
	} else {
		for _, cand := range []string{"cover.png", "cover.jpg", "cover.jpeg", "cover.webp"} {
			p := filepath.Join(dir, cand)
			if data, err := os.ReadFile(p); err == nil {
				w.addImageCover(cand, data)
				w.MetaContent(map[string]string{"cover": cand})
				break
			}
		}
	}

	return nil
}

func (e *Editor) applyMetadata(meta *DirectoryMetadata, dir string) error {
	if meta == nil {
		return nil
	}

	if id := meta.GetIdentifier(); id != "" {
		e.Identifier(id)
	}
	if titles := meta.GetTitles(); len(titles) > 0 {
		e.Title(titles...)
	}
	if authors := meta.GetAuthors(); len(authors) > 0 {
		e.Author(authors...)
	}
	if langs := meta.GetLanguages(); len(langs) > 0 {
		e.Language(langs...)
	}
	if desc := meta.GetDescription(); desc != "" {
		e.Description(desc)
	}
	if pubs := meta.GetPublishers(); len(pubs) > 0 {
		e.Publisher(pubs...)
	}
	if rights := meta.GetRights(); rights != "" {
		e.Rights(rights)
	}
	if subjects := meta.GetSubjects(); len(subjects) > 0 {
		e.Subject(subjects...)
	}
	if dirStr := meta.GetDirection(); dirStr != "" {
		e.Direction(dirStr)
	}
	if d, ok := meta.GetDate(); ok {
		e.Date(d)
	}
	if m, ok := meta.GetModified(); ok {
		e.Modified(m)
	}

	coverFile := meta.GetCover()
	if coverFile != "" {
		fullPath := filepath.Join(dir, coverFile)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to apply cover '%s': %w", coverFile, err)
		}
		if err := e.Cover(data); err != nil {
			return fmt.Errorf("failed to apply cover '%s': %w", coverFile, err)
		}
	} else {
		for _, cand := range []string{"cover.png", "cover.jpg", "cover.jpeg", "cover.webp"} {
			p := filepath.Join(dir, cand)
			if data, err := os.ReadFile(p); err == nil {
				_ = e.Cover(data)
				break
			}
		}
	}
	return nil
}

// detectMIMEType detects the MIME type of a publication resource.
func detectMIMEType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".xhtml", ".html", ".htm":
		return pkg.MediaTypeXHTML
	case ".css":
		return pkg.MediaTypeCSS
	case ".png":
		return pkg.MediaTypePNG
	case ".jpg", ".jpeg":
		return pkg.MediaTypeJPEG
	case ".gif":
		return pkg.MediaTypeGIF
	case ".svg":
		return pkg.MediaTypeSVG
	case ".webp":
		return pkg.MediaTypeWebP
	case ".otf":
		return "font/otf"
	case ".ttf":
		return "font/ttf"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ncx":
		return pkg.MediaTypeNCX
	default:
		t := mime.TypeByExtension(ext)
		if t != "" {
			return t
		}
		return "application/octet-stream"
	}
}

type htmlDocInfo struct {
	href     string
	title    string
	headings []headingInfo
}

// processXHTMLDocument parses an XHTML document, extracts title and headings,
// and ensures headings have id attributes for table of contents anchors.
func processXHTMLDocument(src []byte) ([]byte, []headingInfo, string, error) {
	doc, err := html.Parse(bytes.NewReader(fixSelfClosingTags(src)))
	if err != nil {
		return src, nil, "", err
	}

	title := ""
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "title" {
			title = strings.TrimSpace(GetTextContent(n))
			break
		}
	}

	body := getBody(doc)
	if body == nil {
		return src, nil, title, nil
	}

	var headings []headingInfo
	used := map[string]int{}
	modified := false

	for n := range body.Descendants() {
		level, ok := headingLevel(n.Data)
		if !ok {
			continue
		}
		text := strings.TrimSpace(GetTextContent(n))
		if text == "" {
			continue
		}

		id := ""
		idIndex := -1
		for i, attr := range n.Attr {
			if attr.Key == "id" {
				id = attr.Val
				idIndex = i
				break
			}
		}

		if id == "" {
			id = slugify(text)
			if used[id] > 0 {
				used[id]++
				id = fmt.Sprintf("%s-%d", id, used[id])
			} else {
				used[id] = 1
			}
			if idIndex >= 0 {
				n.Attr[idIndex].Val = id
			} else {
				n.Attr = append(n.Attr, html.Attribute{Key: "id", Val: id})
			}
			modified = true
		} else {
			used[id]++
		}

		headings = append(headings, headingInfo{
			level: level,
			text:  text,
			id:    id,
		})
	}

	if !modified {
		return src, headings, title, nil
	}

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return src, headings, title, nil
	}
	return buf.Bytes(), headings, title, nil
}

func buildHTMLTOC(docs []htmlDocInfo) TOC {
	var root TOC
	parents := []*TOC{&root}

	for _, doc := range docs {
		if len(doc.headings) == 0 {
			fallbackTitle := doc.title
			if fallbackTitle == "" {
				base := filepath.Base(doc.href)
				ext := filepath.Ext(base)
				fallbackTitle = strings.Title(strings.ReplaceAll(strings.TrimSuffix(base, ext), "-", " "))
			}
			root.Items = append(root.Items, TOC{
				Title: fallbackTitle,
				Href:  doc.href,
			})
			parents = []*TOC{&root}
			continue
		}

		for _, h := range doc.headings {
			item := TOC{Title: h.text, Href: doc.href + "#" + h.id}
			level := h.level
			if level > len(parents) {
				level = len(parents)
			}
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

// AddDirectory reads an EPUB content directory containing XHTML/HTML documents,
// assets (images, CSS, fonts), and an optional metadata.yml file.
//
// If metadata.yml (or metadata.yaml) is present, publication metadata (title,
// author, language, identifier, cover, etc.) is automatically configured.
//
// XHTML documents (.xhtml, .html, .htm) are added to the spine in filename order,
// and a table of contents is automatically generated from the document headings.
// Static assets (images, stylesheets, fonts) are added to the publication manifest.
func (w *Writer) AddDirectory(dir string) error {
	meta, err := parseMetadataFile(dir)
	if err != nil {
		return err
	}
	if err := w.applyMetadata(meta, dir); err != nil {
		return err
	}

	var xhtmlFiles []string
	var assetFiles []string

	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		relSlash := filepath.ToSlash(relPath)
		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") {
			return nil
		}

		lowerBase := strings.ToLower(base)
		if lowerBase == "metadata.yml" || lowerBase == "metadata.yaml" || lowerBase == "book.yml" || lowerBase == "book.yaml" {
			return nil
		}

		if meta != nil && meta.GetCover() != "" && relSlash == filepath.ToSlash(meta.GetCover()) {
			return nil
		}
		if (meta == nil || meta.GetCover() == "") && (lowerBase == "cover.png" || lowerBase == "cover.jpg" || lowerBase == "cover.jpeg" || lowerBase == "cover.webp") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".xhtml", ".html", ".htm":
			xhtmlFiles = append(xhtmlFiles, relPath)
		case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".css", ".otf", ".ttf", ".woff", ".woff2":
			assetFiles = append(assetFiles, relPath)
		}
		return nil
	})
	if err != nil {
		return err
	}

	sort.Strings(assetFiles)
	for _, rel := range assetFiles {
		data, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			return err
		}
		href := filepath.ToSlash(rel)
		mimeType := detectMIMEType(rel)
		if _, err := w.AddResource("", href, mimeType, pkg.NotProperty, data); err != nil {
			return err
		}
	}

	sort.Strings(xhtmlFiles)
	var docs []htmlDocInfo
	for _, rel := range xhtmlFiles {
		data, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			return err
		}
		annotatedData, headings, title, _ := processXHTMLDocument(data)
		href := filepath.ToSlash(rel)
		w.AddContent(href, annotatedData)
		docs = append(docs, htmlDocInfo{
			href:     href,
			title:    title,
			headings: headings,
		})
	}

	if len(docs) > 0 {
		toc := buildHTMLTOC(docs)
		toc.Title = meta.GetTOCTitle()
		if len(toc.Items) > 0 {
			if err := w.TableOfContents("toc", toc); err != nil {
				return err
			}
		}
	}

	return nil
}

// AddDirectory reads an EPUB content directory containing XHTML/HTML documents,
// assets, and optional metadata.yml, adding them to this Editor instance.
func (e *Editor) AddDirectory(dir string) error {
	meta, err := parseMetadataFile(dir)
	if err != nil {
		return err
	}
	if err := e.applyMetadata(meta, dir); err != nil {
		return err
	}

	var xhtmlFiles []string
	var assetFiles []string

	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		relSlash := filepath.ToSlash(relPath)
		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") {
			return nil
		}

		lowerBase := strings.ToLower(base)
		if lowerBase == "metadata.yml" || lowerBase == "metadata.yaml" || lowerBase == "book.yml" || lowerBase == "book.yaml" {
			return nil
		}

		if meta != nil && meta.GetCover() != "" && relSlash == filepath.ToSlash(meta.GetCover()) {
			return nil
		}
		if (meta == nil || meta.GetCover() == "") && (lowerBase == "cover.png" || lowerBase == "cover.jpg" || lowerBase == "cover.jpeg" || lowerBase == "cover.webp") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".xhtml", ".html", ".htm":
			xhtmlFiles = append(xhtmlFiles, relPath)
		case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".css", ".otf", ".ttf", ".woff", ".woff2":
			assetFiles = append(assetFiles, relPath)
		}
		return nil
	})
	if err != nil {
		return err
	}

	sort.Strings(assetFiles)
	for _, rel := range assetFiles {
		data, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			return err
		}
		href := filepath.ToSlash(rel)
		mimeType := detectMIMEType(rel)
		if _, err := e.AddResource("", href, mimeType, pkg.NotProperty, data); err != nil {
			return err
		}
	}

	sort.Strings(xhtmlFiles)
	var docs []htmlDocInfo
	for _, rel := range xhtmlFiles {
		data, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			return err
		}
		annotatedData, headings, title, _ := processXHTMLDocument(data)
		href := filepath.ToSlash(rel)
		if _, err := e.AddContent(href, annotatedData); err != nil {
			return err
		}
		docs = append(docs, htmlDocInfo{
			href:     href,
			title:    title,
			headings: headings,
		})
	}

	if len(docs) > 0 {
		toc := buildHTMLTOC(docs)
		toc.Title = meta.GetTOCTitle()
		if len(toc.Items) > 0 {
			if err := e.TableOfContents("toc", toc); err != nil {
				return err
			}
		}
	}

	return nil
}
