package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/raitucarp/epub"
	"github.com/raitucarp/epub/pkg"
)

func main() {
	fmt.Println("=== Creating Advanced Multilingual & RTL EPUB ===")

	w := epub.New("urn:uuid:advanced-multilingual-rtl")
	w.Title("رحلة المعرفة", "A Journey of Knowledge")
	w.Author("ابن خلدون", "Ibn Khaldun")
	w.Languages("ar", "en")
	w.Publisher("Dar Al-Kitab")

	// 1. Set Right-to-Left (RTL) reading progression
	w.Direction("rtl")

	// 2. Set Dublin Core fields
	w.DublinCores(map[string]string{
		"rights": "Public Domain",
		"type":   "Text",
	})
	w.Date(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	w.Modified(time.Now())

	// 3. Metadata Refinements (title-type, display-seq, etc.)
	w.Refines("#title", "title-type", "main")
	w.Refines("#author", "role", "aut", pkg.Meta{Scheme: "marc:relators"})

	// 4. Custom meta properties
	w.MetaProperty("generator-id", "generator", "raitucarp/epub")
	w.MetaContent(map[string]string{
		"rendition:layout":      "reflowable",
		"rendition:orientation": "auto",
	})

	// 5. Add RTL Content
	chapterArabic := []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" dir="rtl" xml:lang="ar">
<head>
  <title>الفصل الأول</title>
  <style>
    body { direction: rtl; font-family: sans-serif; line-height: 1.8; text-align: right; }
    h1 { color: #2c3e50; }
  </style>
</head>
<body>
  <h1>الفصل الأول: المقدمة</h1>
  <p>هذا نموذج لكتاب إلكتروني يدعم اتجاه القراءة من اليمين إلى اليسار (RTL) وفق معايير EPUB 3.</p>
  <p>يتم دعم اللغة العربية والنصوص ثنائية الاتجاه بشكل كامل.</p>
</body>
</html>`)

	chapterEnglish := []byte(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" dir="ltr" xml:lang="en">
<head>
  <title>Chapter 2: English Notes</title>
</head>
<body>
  <h1>Chapter 2: English Translation &amp; Notes</h1>
  <p>This chapter is set in Left-to-Right (LTR) orientation while the overall book retains RTL progression.</p>
</body>
</html>`)

	w.AddContent("chapter-1-ar.xhtml", chapterArabic)
	w.AddContent("chapter-2-en.xhtml", chapterEnglish)

	// 6. Table of Contents
	toc := epub.TOC{
		Title: "فهرس المحتويات / Contents",
		Items: []epub.TOC{
			{Title: "الفصل الأول: المقدمة", Href: "chapter-1-ar.xhtml"},
			{Title: "Chapter 2: Translation", Href: "chapter-2-en.xhtml"},
		},
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		log.Fatalf("Failed to generate TOC: %v", err)
	}

	// 7. Save
	outDir := filepath.Join("temp", "examples_output")
	_ = os.MkdirAll(outDir, 0o755)
	outPath := filepath.Join(outDir, "multilingual-rtl.epub")

	if err := w.Write(outPath); err != nil {
		log.Fatalf("Write failed: %v", err)
	}

	fmt.Printf("Successfully created: %s\n", outPath)
}
