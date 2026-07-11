package epub

import (
	"github.com/raitucarp/epub/pkg"
	"testing"
)

func TestReader_Version(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": {
					Version: "3.0",
				},
			},
		},
	}

	version := r.Version()
	if version != "3.0" {
		t.Errorf("expected '3.0', got %q", version)
	}
}

func TestReader_Description_Metadata(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			metadata: map[string]any{
				"description": []string{"Primary Description"},
			},
		},
	}

	desc := r.Description()
	if desc != "Primary Description" {
		t.Errorf("expected 'Primary Description', got %q", desc)
	}
}

func TestReader_Description_OptionalMeta(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			metadata: map[string]any{
				"meta": map[string]any{
					"description": []any{"Secondary Description"},
				},
			},
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": &pkg.Package{},
			},
		},
	}

	desc := r.Description()
	if desc != "Secondary Description" {
		t.Errorf("expected 'Secondary Description', got %q", desc)
	}
}

func TestReader_Description_Spine(t *testing.T) {
	p := &pkg.Package{
		Spine: pkg.Spine{
			ItemRefs: []pkg.ItemRef{
				{IDRef: "res1"},
			},
		},
	}

	r := &Reader{
		epub: &Epub{
			metadata:  map[string]any{},
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": p,
			},
			resources: []PublicationResource{
				{
					ID:       "res1",
					Href:     "intro.xhtml",
					Content:  []byte(`<html xmlns:epub="http://www.idpf.org/2007/ops"><body><div epub:type="introduction">Spine Description</div></body></html>`),
					MIMEType: pkg.MediaTypeXHTML,
				},
			},
		},
	}

	desc := r.Description()
	if desc != "Spine Description" {
		t.Errorf("expected 'Spine Description', got %q", desc)
	}
}

func TestReader_Description_References(t *testing.T) {
	p := &pkg.Package{
		Guide: &pkg.Guide{
			References: []pkg.GuideReference{
				{Type: "preface", Href: "preface.xhtml"},
			},
		},
	}

	r := &Reader{
		epub: &Epub{
			metadata:  map[string]any{},
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": p,
			},
			resources: []PublicationResource{
				{
					ID:       "res1",
					Href:     "preface.xhtml",
					Content:  []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body><p>Preface Description</p></body></html>`),
					MIMEType: pkg.MediaTypeXHTML,
				},
			},
		},
	}

	desc := r.Description()
	if desc != "Preface Description" {
		t.Errorf("expected 'Preface Description', got %q", desc)
	}
}

func TestReader_Description_TOC(t *testing.T) {
	// FirstFullContentTOCItem is used.
	r := &Reader{
		epub: &Epub{
			metadata:  map[string]any{},
			rendition: "default",
			packagePubs: map[string]*pkg.Package{
				"default": &pkg.Package{},
			},
			resources: []PublicationResource{
				{
					ID: "toc", Properties: "nav",
					Href:     "toc.xhtml",
					Content:  []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><nav epub:type="toc"><h1>TOC</h1><ol><li><a href="chap1.xhtml">Chapter 1</a></li></ol></nav></html>`),
					MIMEType: pkg.MediaTypeXHTML,
				},
				{
					ID:       "chap1",
					Href:     "chap1.xhtml",
					Content:  []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body><p>TOC Description</p></body></html>`),
					MIMEType: pkg.MediaTypeXHTML,
				},
			},
		},
	}

	desc := r.Description()
	if desc != "TOC Description" {
		t.Errorf("expected 'TOC Description', got %q", desc)
	}
}
