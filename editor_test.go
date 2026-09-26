package epub

import (
	"bytes"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/raitucarp/epub/pkg"
)

func createSampleEpub(t *testing.T) []byte {
	t.Helper()
	w := New("urn:sample:original")
	w.Title("Original Title")
	w.Author("Original Author")
	w.Languages("en")
	w.Description("Original Description")
	w.Cover(testPNGBytes(t, 2, 2))
	w.AddContent("ch1.xhtml", []byte(`<html><body><h1>Chapter 1</h1><p>Original content</p></body></html>`))
	toc := TOC{
		Title: "Contents",
		Items: []TOC{
			{Title: "Chapter 1", Href: "ch1.xhtml"},
		},
	}
	if err := w.TableOfContents("toc", toc); err != nil {
		t.Fatalf("TableOfContents: %v", err)
	}
	b, err := w.WriteBytes()
	if err != nil {
		t.Fatalf("WriteBytes: %v", err)
	}
	return b
}

func TestReader_Edit_Invalid(t *testing.T) {
	var r *Reader
	_, err := r.Edit()
	if err == nil {
		t.Fatal("expected error on nil reader Edit()")
	}

	emptyR := &Reader{}
	_, err = emptyR.Edit()
	if err == nil {
		t.Fatal("expected error on uninitialized reader Edit()")
	}
}

func TestEditor_EditMetadata_SaveAs(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	editor, err := r.Edit()
	if err != nil {
		t.Fatalf("reader.Edit(): %v", err)
	}

	// Mutate metadata concisely
	editor.Title("New Modified Title").
		Author("New Author 1", "New Author 2").
		Description("Updated description").
		Subject("Fiction", "Adventure").
		Publisher("Acme Publishing").
		Language("id").
		Rights("Creative Commons").
		Date(time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC)).
		Identifier("urn:uuid:new-id-123")

	outPath := filepath.Join(t.TempDir(), "edited.epub")
	if err := editor.SaveAs(outPath); err != nil {
		t.Fatalf("editor.SaveAs(): %v", err)
	}

	// Verify reading back the modified file
	savedR, err := OpenReader(outPath)
	if err != nil {
		t.Fatalf("OpenReader: %v", err)
	}

	if titles := savedR.Title(); len(titles) != 1 || titles[0] != "New Modified Title" {
		t.Errorf("unexpected titles: %v", titles)
	}

	authors := savedR.Author()
	if len(authors) != 2 || authors[0] != "New Author 1" || authors[1] != "New Author 2" {
		t.Errorf("unexpected authors: %v", authors)
	}

	if desc, ok := savedR.Metadata()["description"]; !ok || len(desc.([]string)) == 0 || desc.([]string)[0] != "Updated description" {
		t.Errorf("unexpected description: %v", desc)
	}

	if sub, ok := savedR.Metadata()["subject"]; !ok || len(sub.([]string)) != 2 || sub.([]string)[0] != "Fiction" {
		t.Errorf("unexpected subjects: %v", sub)
	}

	if pub, ok := savedR.Metadata()["publisher"]; !ok || len(pub.([]string)) != 1 || pub.([]string)[0] != "Acme Publishing" {
		t.Errorf("unexpected publisher: %v", pub)
	}

	if langs := savedR.Language(); len(langs) != 1 || langs[0] != "id" {
		t.Errorf("unexpected languages: %v", langs)
	}

	if ids := savedR.Identifier(); len(ids) != 1 || ids[0] != "urn:uuid:new-id-123" {
		t.Errorf("unexpected identifiers: %v", ids)
	}
}

func TestEditor_Save_IoWriter(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	editor, err := r.Edit()
	if err != nil {
		t.Fatalf("reader.Edit(): %v", err)
	}

	editor.Title("Writer Output Title")

	var buf bytes.Buffer
	if err := editor.Save(&buf); err != nil {
		t.Fatalf("editor.Save(&buf): %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("expected non-empty buffer from Save()")
	}

	savedR, err := NewReader(buf.Bytes())
	if err != nil {
		t.Fatalf("NewReader from saved buffer: %v", err)
	}

	if titles := savedR.Title(); len(titles) != 1 || titles[0] != "Writer Output Title" {
		t.Errorf("unexpected titles: %v", titles)
	}

	if err := editor.Save(nil); err == nil {
		t.Error("expected error on Save(nil)")
	}
}

func TestEditor_WriteBytes(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	editor, err := r.Edit()
	if err != nil {
		t.Fatalf("reader.Edit(): %v", err)
	}

	editor.Title("Byte Output Title")

	var b []byte
	if err := editor.WriteBytes(&b); err != nil {
		t.Fatalf("editor.WriteBytes(&b): %v", err)
	}

	if len(b) == 0 {
		t.Fatal("expected non-empty bytes in &b from WriteBytes()")
	}

	savedR, err := NewReader(b)
	if err != nil {
		t.Fatalf("NewReader from saved bytes: %v", err)
	}

	if titles := savedR.Title(); len(titles) != 1 || titles[0] != "Byte Output Title" {
		t.Errorf("unexpected titles: %v", titles)
	}

	if err := editor.WriteBytes(nil); err == nil {
		t.Error("expected error on WriteBytes(nil)")
	}

	// Test SaveByte alias
	var bAlias []byte
	if err := editor.SaveByte(&bAlias); err != nil || len(bAlias) == 0 {
		t.Errorf("SaveByte alias failed: %v", err)
	}

	// Also test SaveBytes helper
	bytesOut, err := editor.SaveBytes()
	if err != nil || len(bytesOut) == 0 {
		t.Errorf("SaveBytes failed: %v, len=%d", err, len(bytesOut))
	}
}

func TestEditor_AddContent_UpdateContent_RemoveResource(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	editor, err := r.Edit()
	if err != nil {
		t.Fatalf("reader.Edit(): %v", err)
	}

	// 1. Add content
	newRes, err := editor.AddContent("ch2.xhtml", []byte(`<html><body><h1>Chapter 2</h1><p>New content</p></body></html>`))
	if err != nil {
		t.Fatalf("AddContent: %v", err)
	}
	if newRes.Href != "ch2.xhtml" {
		t.Errorf("unexpected newRes href: %s", newRes.Href)
	}

	// 2. Add an image
	imgRes, err := editor.AddImage("photo.png", testPNGBytes(t, 4, 4))
	if err != nil {
		t.Fatalf("AddImage: %v", err)
	}
	if imgRes.Href != "photo.png" {
		t.Errorf("unexpected imgRes href: %s", imgRes.Href)
	}

	// 3. Update existing content
	updatedHTML := []byte(`<html><body><h1>Chapter 1 Modified</h1><p>Modified text</p></body></html>`)
	if err := editor.UpdateContent("ch1.xhtml", updatedHTML); err != nil {
		t.Fatalf("UpdateContent: %v", err)
	}

	var savedBytes []byte
	if err := editor.WriteBytes(&savedBytes); err != nil {
		t.Fatalf("SaveByte: %v", err)
	}

	savedR, err := NewReader(savedBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	// Check updated content
	ch1Res := savedR.SelectResourceByHref("ch1.xhtml")
	if ch1Res == nil || !strings.Contains(string(ch1Res.Content), "Chapter 1 Modified") {
		t.Errorf("expected updated content in ch1.xhtml, got: %v", ch1Res)
	}

	// Check newly added content
	ch2Str := savedR.ContentDocumentXHTMLString()[newRes.ID]
	if !strings.Contains(ch2Str, "Chapter 2") {
		t.Errorf("expected new content in ch2, got: %s", ch2Str)
	}

	// Check spine has 2 items
	spine := savedR.Spine()
	if len(spine) != 2 {
		t.Errorf("expected 2 spine items, got %d", len(spine))
	}

	// 4. Test remove resource
	editor2, err := savedR.Edit()
	if err != nil {
		t.Fatalf("Edit saved: %v", err)
	}
	if err := editor2.RemoveResource(newRes.ID); err != nil {
		t.Fatalf("RemoveResource: %v", err)
	}

	var finalBytes []byte
	if err := editor2.WriteBytes(&finalBytes); err != nil {
		t.Fatalf("WriteBytes after remove: %v", err)
	}

	finalR, err := NewReader(finalBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}
	if finalR.SelectResourceById(newRes.ID) != nil {
		t.Error("expected removed resource to be nil")
	}
	if len(finalR.Spine()) != 1 {
		t.Errorf("expected 1 spine item after remove, got %d", len(finalR.Spine()))
	}
}

func TestEditor_Cover(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	editor, err := r.Edit()
	if err != nil {
		t.Fatalf("reader.Edit(): %v", err)
	}

	newCoverBytes := testPNGBytes(t, 8, 8)
	if err := editor.Cover(newCoverBytes); err != nil {
		t.Fatalf("editor.Cover: %v", err)
	}

	var b []byte
	if err := editor.WriteBytes(&b); err != nil {
		t.Fatalf("SaveByte: %v", err)
	}

	savedR, err := NewReader(b)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	cover := savedR.Cover()
	if cover == nil {
		t.Fatal("expected cover to be found")
	}
	bounds := (*cover).Bounds()
	if bounds.Dx() != 8 || bounds.Dy() != 8 {
		t.Errorf("expected 8x8 cover image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestEditor_TableOfContents(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	editor, err := r.Edit()
	if err != nil {
		t.Fatalf("reader.Edit(): %v", err)
	}

	newTOC := TOC{
		Title: "New Updated TOC",
		Items: []TOC{
			{Title: "Updated Chapter 1", Href: "ch1.xhtml"},
		},
	}
	if err := editor.TableOfContents("new_toc", newTOC); err != nil {
		t.Fatalf("TableOfContents: %v", err)
	}

	var b []byte
	if err := editor.WriteBytes(&b); err != nil {
		t.Fatalf("SaveByte: %v", err)
	}

	savedR, err := NewReader(b)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	readTOC, err := savedR.TableOfContents()
	if err != nil {
		t.Fatalf("savedR.TableOfContents(): %v", err)
	}
	if len(readTOC.Items) != 1 || readTOC.Items[0].Title != "Updated Chapter 1" {
		t.Errorf("unexpected TOC: %+v", readTOC)
	}
}

func TestEditor_PackageDirectAccess_And_Renditions(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	editor, err := r.Edit()
	if err != nil {
		t.Fatalf("reader.Edit(): %v", err)
	}

	// Check renditions
	renditions := editor.ListRenditions()
	if len(renditions) == 0 {
		t.Fatal("expected at least 1 rendition")
	}

	// Direct package modification
	pkgDoc := editor.Package()
	if pkgDoc == nil {
		t.Fatal("expected non-nil Package()")
	}
	editor.Direction("rtl")
	editor.UniqueIdentifier("custom-uid")
	editor.Version("3.0")

	if editor.CurrentPackagePath() == "" {
		t.Error("expected non-empty CurrentPackagePath")
	}

	var b []byte
	if err := editor.WriteBytes(&b); err != nil {
		t.Fatalf("SaveByte: %v", err)
	}

	savedR, err := NewReader(b)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	if savedR.CurrentSelectedPackage().Dir != "rtl" {
		t.Errorf("expected dir rtl, got %s", savedR.CurrentSelectedPackage().Dir)
	}
	if savedR.UID() != "custom-uid" && savedR.CurrentSelectedPackage().UniqueIdentifier != "custom-uid" {
		t.Errorf("expected UID custom-uid, got %s", savedR.CurrentSelectedPackage().UniqueIdentifier)
	}
}

func TestEditor_Reader(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	editor, err := r.Edit()
	if err != nil {
		t.Fatalf("reader.Edit(): %v", err)
	}

	editor.Title("Transitioned to Reader Title").
		Author("Transitioned Author")

	// Transition back to reader mode directly
	newReader, err := editor.Reader()
	if err != nil {
		t.Fatalf("editor.Reader(): %v", err)
	}

	if titles := newReader.Title(); len(titles) != 1 || titles[0] != "Transitioned to Reader Title" {
		t.Errorf("unexpected titles: %v", titles)
	}

	if authors := newReader.Author(); len(authors) != 1 || authors[0] != "Transitioned Author" {
		t.Errorf("unexpected authors: %v", authors)
	}

	// Verify the new reader can also enter edit mode again
	editor2, err := newReader.Edit()
	if err != nil {
		t.Fatalf("newReader.Edit(): %v", err)
	}
	editor2.Title("Second Edit")
	finalReader, err := editor2.Reader()
	if err != nil {
		t.Fatalf("editor2.Reader(): %v", err)
	}
	if titles := finalReader.Title(); len(titles) != 1 || titles[0] != "Second Edit" {
		t.Errorf("unexpected titles: %v", titles)
	}
}

func TestEditor_MetadataMethods(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	ed, err := r.Edit()
	if err != nil {
		t.Fatalf("Edit(): %v", err)
	}

	ed.AddTitle("Secondary Title").
		AddAuthor("Second Author").
		Creator("illustrator", "Art Person").
		AddDescription("Additional descriptive details").
		AddPublisher("Second Publishing House").
		Contributor("editor", "Chief Editor").
		AddSubject("Science").
		AddLanguage("fr").
		AddIdentifier("isbn", "978-3-16-148410-0").
		DublinCore("source", "https://example.com/source").
		SetMeta("generator", "custom-generator").
		Meta(pkg.Meta{Name: "custom_meta", Content: "val"}).
		MetaContent(map[string]string{"foo": "bar"}).
		MetaProperty("prop-1", "custom-property", "prop-val").
		Refines("prop-1", "scheme", "custom-scheme").
		Modified(time.Now())

	// Verify DublinCore was added
	foundSource := false
	for _, dc := range ed.CurrentPackage().Metadata.OptionalDC {
		if dc.ID == "source" || dc.XMLName.Local == "source" {
			foundSource = true
			break
		}
	}
	if !foundSource {
		t.Error("expected source DublinCore metadata")
	}

	// Remove meta
	ed.RemoveMeta("custom_meta")
	ed.RemoveMetadata("source")

	b, err := ed.SaveBytes()
	if err != nil {
		t.Fatalf("SaveBytes(): %v", err)
	}

	r2, err := NewReader(b)
	if err != nil {
		t.Fatalf("NewReader r2: %v", err)
	}

	if len(r2.Title()) < 2 {
		t.Errorf("expected at least 2 titles, got %d", len(r2.Title()))
	}
	if len(r2.Author()) < 2 {
		t.Errorf("expected at least 2 authors, got %d", len(r2.Author()))
	}
}

func TestEditor_FilesAndAssets(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	ed, err := r.Edit()
	if err != nil {
		t.Fatalf("Edit(): %v", err)
	}

	// AddFile raw
	ed.AddFile("misc/data.txt", []byte("arbitrary text data"))

	// AddContentFile from disk
	tempDir := t.TempDir()
	contentFilePath := filepath.Join(tempDir, "ch2.xhtml")
	if err := os.WriteFile(contentFilePath, []byte("<html><body><h1>Ch2</h1></body></html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	resContent, err := ed.AddContentFile(contentFilePath)
	if err != nil {
		t.Fatalf("AddContentFile: %v", err)
	}
	if resContent.Href != "ch2.xhtml" {
		t.Errorf("expected href ch2.xhtml, got %q", resContent.Href)
	}

	// AddMarkdownFile from disk
	mdFilePath := filepath.Join(tempDir, "ch3.md")
	if err := os.WriteFile(mdFilePath, []byte("# Chapter 3\n\nMarkdown text."), 0o644); err != nil {
		t.Fatal(err)
	}
	resMD, err := ed.AddMarkdownFile(mdFilePath)
	if err != nil {
		t.Fatalf("AddMarkdownFile: %v", err)
	}
	if resMD.Href != "ch3.xhtml" {
		t.Errorf("expected href ch3.xhtml, got %q", resMD.Href)
	}

	// AddImageFile from disk
	img := image.NewRGBA(image.Rect(0, 0, 5, 5))
	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			img.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255})
		}
	}
	pngPath := filepath.Join(tempDir, "img.png")
	if err := os.WriteFile(pngPath, testPNGBytes(t, 5, 5), 0o644); err != nil {
		t.Fatal(err)
	}
	resImg, err := ed.AddImageFile(pngPath)
	if err != nil {
		t.Fatalf("AddImageFile: %v", err)
	}
	if resImg.Href != "img.png" {
		t.Errorf("expected href img.png, got %q", resImg.Href)
	}

	// CoverPNG, CoverJPG, CoverFile
	if err := ed.CoverPNG(img); err != nil {
		t.Fatalf("CoverPNG: %v", err)
	}
	if err := ed.CoverJPG(img); err != nil {
		t.Fatalf("CoverJPG: %v", err)
	}
	if err := ed.CoverFile(pngPath); err != nil {
		t.Fatalf("CoverFile: %v", err)
	}

	// Resources inspection
	allRes := ed.Resources()
	if len(allRes) == 0 {
		t.Error("expected non-empty Resources()")
	}
	if ed.SelectResourceById(resContent.ID) == nil {
		t.Errorf("expected to find %s by ID", resContent.ID)
	}
	if ed.SelectResourceByHref(resMD.Href) == nil {
		t.Errorf("expected to find %s by href", resMD.Href)
	}

	// AddSpineItem and RemoveSpineItem
	if err := ed.AddSpineItem("ch2.xhtml"); err != nil {
		t.Fatalf("AddSpineItem: %v", err)
	}
	if err := ed.RemoveSpineItem("ch2.xhtml"); err != nil {
		t.Fatalf("RemoveSpineItem: %v", err)
	}

	// SaveByte alias
	var b []byte
	if err := ed.SaveByte(&b); err != nil {
		t.Fatalf("SaveByte: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty byte slice from SaveByte")
	}
}

func TestEditor_ErrorBranches(t *testing.T) {
	sampleBytes := createSampleEpub(t)
	r, err := NewReader(sampleBytes)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	ed, err := r.Edit()
	if err != nil {
		t.Fatalf("Edit(): %v", err)
	}

	// Save(nil)
	if err := ed.Save(nil); err == nil {
		t.Error("expected error on Save(nil)")
	}

	// WriteBytes(nil)
	if err := ed.WriteBytes(nil); err == nil {
		t.Error("expected error on WriteBytes(nil)")
	}

	// SelectPackageRendition non-existent
	if err := ed.SelectPackageRendition("non-existent-rendition"); err == nil {
		t.Error("expected error on non-existent rendition")
	}

	// RemoveResource non-existent
	if err := ed.RemoveResource("completely-bogus-resource"); err == nil {
		t.Error("expected error on RemoveResource with non-existent ID")
	}

	// RemoveSpineItem non-existent
	if err := ed.RemoveSpineItem("completely-bogus-spine"); err == nil {
		t.Error("expected error on RemoveSpineItem with non-existent ID")
	}

	// AddSpineItem non-existent
	if err := ed.AddSpineItem("completely-bogus-spine"); err == nil {
		t.Error("expected error on AddSpineItem with non-existent ID")
	}

	// AddContentFile non-existent
	if _, err := ed.AddContentFile("non-existent-file-path.xhtml"); err == nil {
		t.Error("expected error on missing content file")
	}

	// AddMarkdownFile non-existent
	if _, err := ed.AddMarkdownFile("non-existent-file-path.md"); err == nil {
		t.Error("expected error on missing markdown file")
	}

	// AddImageFile non-existent
	if _, err := ed.AddImageFile("non-existent-file-path.png"); err == nil {
		t.Error("expected error on missing image file")
	}

	// CoverFile non-existent
	if err := ed.CoverFile("non-existent-file-path.png"); err == nil {
		t.Error("expected error on missing cover file")
	}
}

