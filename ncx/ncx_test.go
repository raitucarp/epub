package ncx

import (
	"testing"
)

func TestParse_Success(t *testing.T) {
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
  <head>
    <meta name="dtb:uid" content="urn:uuid:12345"/>
  </head>
  <docTitle>
    <text>Test Title</text>
  </docTitle>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel>
        <text>Chapter 1</text>
      </navLabel>
      <content src="chapter1.html"/>
    </navPoint>
  </navMap>
</ncx>`)

	ncx, err := Parse(xmlData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ncx == nil {
		t.Fatal("expected ncx to be non-nil")
	}

	if ncx.Version != "2005-1" {
		t.Errorf("expected version '2005-1', got '%v'", ncx.Version)
	}
	if ncx.DocTitle.Text != "Test Title" {
		t.Errorf("expected docTitle 'Test Title', got '%v'", ncx.DocTitle.Text)
	}
	if len(ncx.NavMap.NavPoints) != 1 {
		t.Fatalf("expected 1 navPoint, got %d", len(ncx.NavMap.NavPoints))
	}
	if ncx.NavMap.NavPoints[0].NavLabel.Text != "Chapter 1" {
		t.Errorf("expected NavLabel 'Chapter 1', got '%v'", ncx.NavMap.NavPoints[0].NavLabel.Text)
	}
	if ncx.NavMap.NavPoints[0].Content.Src != "chapter1.html" {
		t.Errorf("expected Content Src 'chapter1.html', got '%v'", ncx.NavMap.NavPoints[0].Content.Src)
	}
}

func TestParse_Error(t *testing.T) {
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<ncx>
  <head>
    <meta> <!-- unclosed tag -->
  </head>
</ncx>`)

	_, err := Parse(xmlData)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
