package epub

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

